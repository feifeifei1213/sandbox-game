package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

func TestOrderWorkflowCoversGenerationSequenceSelectionDeliveryAndUnfinished(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	yearNo := int(nextIntegrationUniqueSeed()%100000) + 1000
	previousYearNo := yearNo - 1
	now := time.Now()

	ensureIntegrationGameConfig(t, ctx, tx, yearNo, yearNo, true)
	isolateOrderIntegrationGroups(t, ctx, tx)

	groupOneID := createOrderIntegrationGroup(t, ctx, tx, "I4 order leader", enum.BusinessStatusNormal, nil)
	groupTwoID := createOrderIntegrationGroup(t, ctx, tx, "I4 order follower", enum.BusinessStatusNormal, nil)
	bankruptYear := previousYearNo
	bankruptGroupID := createOrderIntegrationGroup(t, ctx, tx, "I4 bankrupt previous winner", enum.BusinessStatusBankrupt, &bankruptYear)

	createPreviousFormalReportRecord(t, ctx, tx, groupOneID, previousYearNo, now)
	createPreviousFormalReportRecord(t, ctx, tx, groupTwoID, previousYearNo, now)
	createGroupYearStateRecord(t, ctx, tx, groupOneID, yearNo, enum.YearTypeFormal, enum.YearStatusOperating, enum.StageStatusQ1Open, enum.ReportStatusLocked)
	createGroupYearStateRecord(t, ctx, tx, groupTwoID, yearNo, enum.YearTypeFormal, enum.YearStatusOperating, enum.StageStatusQ1Open, enum.ReportStatusLocked)
	createSelectedOrderAmountFixture(t, ctx, tx, groupOneID, previousYearNo, enum.MarketCodeLocal, 50)
	createSelectedOrderAmountFixture(t, ctx, tx, groupTwoID, previousYearNo, enum.MarketCodeLocal, 20)
	createSelectedOrderAmountFixture(t, ctx, tx, bankruptGroupID, previousYearNo, enum.MarketCodeLocal, 999)

	adminOrderService := NewAdminOrderCommandService(tx)
	adminControlService := NewAdminOrderControlCommandService(tx)
	playerOrderService := NewPlayerOrderCommandService(tx)

	if _, err := adminOrderService.UpdateMarketConfig(ctx, UpdateOrderMarketConfigCommand{
		YearNo: yearNo,
		Markets: []UpdateOrderMarketConfigItem{
			{MarketCode: enum.MarketCodeLocal, Enabled: true},
			{MarketCode: enum.MarketCodeRegional, Enabled: true},
			{MarketCode: enum.MarketCodeNational, Enabled: false},
			{MarketCode: enum.MarketCodeGlobal, Enabled: false},
		},
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("update market enabled config: %v", err)
	}

	configResult, err := adminOrderService.UpdateControlConfig(ctx, UpdateOrderControlConfigCommand{
		YearNo: yearNo,
		Items: buildTestControlConfigItems(yearNo, map[string]int{
			testOrderSegmentKey(enum.MarketCodeLocal, enum.OrderTypeAgencyInspection):    2,
			testOrderSegmentKey(enum.MarketCodeRegional, enum.OrderTypeAgencyInspection): 1,
			testOrderSegmentKey(enum.MarketCodeLocal, enum.OrderTypeTwoCabinVIP):         1,
		}),
		OperatorID:   1,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("update order control config: %v", err)
	}
	if len(configResult.Warnings) == 0 {
		t.Fatalf("expected low-order-count warnings to help admin control pool redundancy")
	}

	generateResult, err := adminOrderService.GenerateOrderPool(ctx, GenerateOrderPoolCommand{
		YearNo:       yearNo,
		Overwrite:    true,
		OperatorID:   1,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("generate order preview pool: %v", err)
	}
	if generateResult.BatchStatus != enum.OrderGenerationBatchStatusPreview || generateResult.RandomSeed == "" || generateResult.FormulaVersion != orderGenerationFormulaVersion {
		t.Fatalf("expected preview batch with fixed seed and formula version, got %#v", generateResult)
	}
	if generateResult.GeneratedCount != 4 || generateResult.SegmentCount != requiredOrderInvestmentSegmentCount {
		t.Fatalf("expected 4 generated orders across 16 segments, got %#v", generateResult)
	}

	confirmed, err := adminOrderService.ConfirmOrderPool(ctx, ConfirmOrderPoolCommand{
		YearNo:       yearNo,
		BatchID:      generateResult.BatchID,
		OperatorID:   1,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("confirm order pool: %v", err)
	}
	if confirmed.BatchID != generateResult.BatchID || confirmed.GeneratedCount != 4 {
		t.Fatalf("unexpected confirmed pool result: %#v", confirmed)
	}

	if _, err := playerOrderService.SubmitMarketInvestment(ctx, SubmitMarketInvestmentCommand{
		GroupID: groupOneID,
		YearNo:  yearNo,
		Investments: buildTestMarketInvestments(func(segment OrderSegmentDefinition) float64 {
			if segment.MarketCode == enum.MarketCodeLocal && segment.OrderType == enum.OrderTypeAgencyInspection {
				return 10
			}
			return 0
		}),
		OperatorName: "group-one",
	}); err != nil {
		t.Fatalf("group one submit market investments: %v", err)
	}
	if _, err := playerOrderService.SubmitMarketInvestment(ctx, SubmitMarketInvestmentCommand{
		GroupID: groupTwoID,
		YearNo:  yearNo,
		Investments: buildTestMarketInvestments(func(segment OrderSegmentDefinition) float64 {
			switch {
			case segment.MarketCode == enum.MarketCodeLocal && segment.OrderType == enum.OrderTypeAgencyInspection:
				return 10
			case segment.MarketCode == enum.MarketCodeLocal && segment.OrderType == enum.OrderTypeTwoCabinVIP:
				return 5
			default:
				return 0
			}
		}),
		OperatorName: "group-two",
	}); err != nil {
		t.Fatalf("group two submit market investments: %v", err)
	}

	sequenceResult, err := adminControlService.GenerateSelectionSequence(ctx, GenerateSelectionSequenceCommand{
		YearNo:       yearNo,
		OperatorID:   1,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("generate selection sequence: %v", err)
	}
	if sequenceResult.AffectedCount != 3 {
		t.Fatalf("expected configured enabled segments to become ready/skipped, got %#v", sequenceResult)
	}

	stateRepo := repository.NewMarketBiddingStateRepository(tx)
	sequenceRepo := repository.NewMarketSelectionOrderRepository(tx)
	poolRepo := repository.NewOrderPoolRepository(tx)
	selectionRepo := repository.NewGroupOrderSelectionRepository(tx)

	localAgencyState, err := stateRepo.GetBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("load local agency state: %v", err)
	}
	if localAgencyState.SegmentStatus != enum.OrderSegmentStatusSequenceReady || localAgencyState.LeaderGroupID == nil || *localAgencyState.LeaderGroupID != groupOneID {
		t.Fatalf("expected non-bankrupt previous local leader to lead local agency, got %#v", localAgencyState)
	}
	regionalAgencyState, err := stateRepo.GetBySegment(ctx, yearNo, enum.MarketCodeRegional, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("load regional agency state: %v", err)
	}
	if regionalAgencyState.SegmentStatus != enum.OrderSegmentStatusSkipped {
		t.Fatalf("expected no-investment segment to be skipped, got %#v", regionalAgencyState)
	}
	nationalAgencyState, err := stateRepo.GetBySegment(ctx, yearNo, enum.MarketCodeNational, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("load national agency state: %v", err)
	}
	if nationalAgencyState.SegmentStatus != enum.OrderSegmentStatusMarketDisabled {
		t.Fatalf("expected disabled national market to be marked disabled, got %#v", nationalAgencyState)
	}
	localTwoCabinState, err := stateRepo.GetBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeTwoCabinVIP)
	if err != nil {
		t.Fatalf("load local two-cabin state: %v", err)
	}
	if localTwoCabinState.LeaderGroupID == nil || *localTwoCabinState.LeaderGroupID != groupOneID {
		t.Fatalf("expected the same local market leader to be reused across local product segments, got %#v", localTwoCabinState)
	}

	localAgencySequence, err := sequenceRepo.ListBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("list local agency sequence: %v", err)
	}
	if len(localAgencySequence) != 2 || localAgencySequence[0].GroupID != groupOneID || !localAgencySequence[0].IsMarketLeader {
		t.Fatalf("expected local leader to be first for local agency, got %#v", localAgencySequence)
	}
	localTwoCabinSequence, err := sequenceRepo.ListBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeTwoCabinVIP)
	if err != nil {
		t.Fatalf("list local two-cabin sequence: %v", err)
	}
	if len(localTwoCabinSequence) != 2 || localTwoCabinSequence[0].GroupID != groupOneID || !localTwoCabinSequence[0].IsMarketLeader {
		t.Fatalf("expected local leader with zero segment investment to still be first, got %#v", localTwoCabinSequence)
	}

	released, err := adminControlService.ReleaseNextSegment(ctx, ReleaseNextSegmentCommand{
		YearNo:       yearNo,
		OperatorID:   1,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("release first segment: %v", err)
	}
	if released.MarketCode != enum.MarketCodeLocal || released.OrderType != enum.OrderTypeAgencyInspection || released.CurrentGroupID == nil || *released.CurrentGroupID != groupOneID {
		t.Fatalf("expected first released segment to follow configured order and current group, got %#v", released)
	}
	if _, err := adminControlService.ReleaseNextSegment(ctx, ReleaseNextSegmentCommand{
		YearNo:       yearNo,
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); !errors.Is(err, ErrOrderSegmentReleaseBlocked) {
		t.Fatalf("expected release to be blocked while current segment is selecting, got %v", err)
	}

	localAgencyOrders, err := poolRepo.ListBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("list local agency pool: %v", err)
	}
	if len(localAgencyOrders) != 2 {
		t.Fatalf("expected two local agency orders, got %d", len(localAgencyOrders))
	}
	if _, err := playerOrderService.SelectOrder(ctx, SelectOrderCommand{
		GroupID:      groupTwoID,
		YearNo:       yearNo,
		MarketCode:   enum.MarketCodeLocal,
		OrderType:    enum.OrderTypeAgencyInspection,
		OrderID:      localAgencyOrders[0].ID,
		OperatorName: "group-two",
	}); !errors.Is(err, ErrOrderSegmentCurrentGroupMismatch) {
		t.Fatalf("expected off-turn group to be blocked, got %v", err)
	}
	if _, err := playerOrderService.SelectOrder(ctx, SelectOrderCommand{
		GroupID:      groupOneID,
		YearNo:       yearNo,
		MarketCode:   enum.MarketCodeLocal,
		OrderType:    enum.OrderTypeAgencyInspection,
		OrderID:      localAgencyOrders[0].ID,
		OperatorName: "group-one",
	}); err != nil {
		t.Fatalf("group one select local agency order: %v", err)
	}
	afterSelectState, err := stateRepo.GetBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("reload local agency after select: %v", err)
	}
	if afterSelectState.SegmentStatus != enum.OrderSegmentStatusSelecting || afterSelectState.CurrentGroupID == nil || *afterSelectState.CurrentGroupID != groupTwoID {
		t.Fatalf("expected selecting segment to advance to group two, got %#v", afterSelectState)
	}
	if _, err := adminControlService.AdminSkipCurrentGroup(ctx, AdminSkipCurrentGroupCommand{
		YearNo:       yearNo,
		MarketCode:   enum.MarketCodeLocal,
		OrderType:    enum.OrderTypeAgencyInspection,
		GroupID:      groupTwoID,
		Reason:       "integration: absent from table",
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("admin skip current group: %v", err)
	}
	completedAgencyState, err := stateRepo.GetBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("reload local agency after skip: %v", err)
	}
	if completedAgencyState.SegmentStatus != enum.OrderSegmentStatusCompleted {
		t.Fatalf("expected one-round segment to complete after sequence is exhausted, got %#v", completedAgencyState)
	}
	localAgencyOrders, err = poolRepo.ListBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("reload local agency pool: %v", err)
	}
	if localAgencyOrders[0].PoolStatus != enum.OrderPoolStatusSelected || localAgencyOrders[1].PoolStatus != enum.OrderPoolStatusAvailable {
		t.Fatalf("expected selected order locked and unchosen order retained available after one round, got %#v", localAgencyOrders)
	}

	operatingCommandService := buildOrderLinkedOperatingCommandService(tx)
	if _, err := operatingCommandService.SubmitStage(ctx, SubmitOperatingStageCommand{
		GroupID:          groupOneID,
		YearNo:           yearNo,
		StageCode:        state.StageCodeQ1,
		OperatingPayload: buildValidQ1OperatingPayload(),
		SubmitterID:      groupOneID,
		OperatorName:     "group-one",
	}); !errors.Is(err, ErrOrderPrerequisiteIncomplete) {
		t.Fatalf("expected Q1 submit to be blocked before all order segments finish, got %v", err)
	}

	releasedSecond, err := adminControlService.ReleaseNextSegment(ctx, ReleaseNextSegmentCommand{
		YearNo:       yearNo,
		OperatorID:   1,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("release next non-skipped segment: %v", err)
	}
	if releasedSecond.MarketCode != enum.MarketCodeLocal || releasedSecond.OrderType != enum.OrderTypeTwoCabinVIP || releasedSecond.CurrentGroupID == nil || *releasedSecond.CurrentGroupID != groupOneID {
		t.Fatalf("expected skipped regional segment to be bypassed and local two-cabin released, got %#v", releasedSecond)
	}
	twoCabinOrders, err := poolRepo.ListBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeTwoCabinVIP)
	if err != nil {
		t.Fatalf("list local two-cabin pool: %v", err)
	}
	if len(twoCabinOrders) != 1 {
		t.Fatalf("expected one local two-cabin order, got %d", len(twoCabinOrders))
	}
	if _, err := playerOrderService.SelectOrder(ctx, SelectOrderCommand{
		GroupID:      groupOneID,
		YearNo:       yearNo,
		MarketCode:   enum.MarketCodeLocal,
		OrderType:    enum.OrderTypeTwoCabinVIP,
		OrderID:      twoCabinOrders[0].ID,
		OperatorName: "group-one",
	}); err != nil {
		t.Fatalf("group one select local two-cabin order: %v", err)
	}
	if completed, err := NewOrderOperatingLinkService(
		repository.NewGroupMarketBidRepository(tx),
		stateRepo,
		selectionRepo,
	).IsPrerequisiteCompleted(ctx, yearNo); err != nil || !completed {
		t.Fatalf("expected order prerequisite to be complete after all segments completed/skipped, completed=%t err=%v", completed, err)
	}

	selectedOrders, err := selectionRepo.ListByGroupYear(ctx, groupOneID, yearNo)
	if err != nil {
		t.Fatalf("list group one selected orders: %v", err)
	}
	if len(selectedOrders) != 2 {
		t.Fatalf("expected group one to have two selected orders, got %#v", selectedOrders)
	}
	firstSelection, err := selectionRepo.GetByGroupSegment(ctx, groupOneID, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("load local agency selection: %v", err)
	}
	firstDeliveredOrderID := firstSelection.OrderID
	firstDeliveredAmount := orderAmountByID(t, localAgencyOrders, firstDeliveredOrderID)
	saveOperatingDraft(t, ctx, tx, groupOneID, yearNo, enum.StageStatusQ1Open, buildQ1PayloadWithRevenue(0))
	if _, err := playerOrderService.DeliverOrders(ctx, DeliverOrdersCommand{
		GroupID:      groupOneID,
		YearNo:       yearNo,
		StageCode:    state.StageCodeQ1,
		OrderIDs:     []int64{firstDeliveredOrderID},
		OperatorName: "group-one",
	}); !errors.Is(err, ErrOrderDeliveryRevenueMismatch) {
		t.Fatalf("expected delivery revenue mismatch to be rejected, got %v", err)
	}
	saveOperatingDraft(t, ctx, tx, groupOneID, yearNo, enum.StageStatusQ1Open, buildQ1PayloadWithRevenue(firstDeliveredAmount))
	if _, err := playerOrderService.DeliverOrders(ctx, DeliverOrdersCommand{
		GroupID:      groupOneID,
		YearNo:       yearNo,
		StageCode:    state.StageCodeQ1,
		OrderIDs:     []int64{firstDeliveredOrderID},
		OperatorName: "group-one",
	}); err != nil {
		t.Fatalf("deliver order after matching revenue: %v", err)
	}

	if _, err := operatingCommandService.SubmitStage(ctx, SubmitOperatingStageCommand{
		GroupID:          groupOneID,
		YearNo:           yearNo,
		StageCode:        state.StageCodeQ1,
		OperatingPayload: buildQ1PayloadWithRevenue(firstDeliveredAmount),
		SubmitterID:      groupOneID,
		OperatorName:     "group-one",
	}); err != nil {
		t.Fatalf("submit Q1 after order prerequisite completed: %v", err)
	}

	setGroupYearStage(t, ctx, tx, groupOneID, yearNo, enum.YearStatusOperating, enum.StageStatusYearEndOpen, enum.ReportStatusLocked, 4)
	if _, err := operatingCommandService.SubmitStage(ctx, SubmitOperatingStageCommand{
		GroupID:          groupOneID,
		YearNo:           yearNo,
		StageCode:        state.StageCodeYearEnd,
		OperatingPayload: buildYearEndOperatingPayload(),
		SubmitterID:      groupOneID,
		OperatorName:     "group-one",
	}); err != nil {
		t.Fatalf("submit year end to mark unfinished orders: %v", err)
	}

	finalSelections, err := selectionRepo.ListByGroupYear(ctx, groupOneID, yearNo)
	if err != nil {
		t.Fatalf("reload group selections after year end: %v", err)
	}
	deliveryStatusByOrderID := make(map[int64]string, len(finalSelections))
	for _, item := range finalSelections {
		deliveryStatusByOrderID[item.OrderID] = item.DeliveryStatus
	}
	if deliveryStatusByOrderID[firstDeliveredOrderID] != enum.OrderDeliveryStatusDelivered {
		t.Fatalf("expected delivered order to stay delivered, got %#v", finalSelections)
	}
	var unfinishedFound bool
	for _, item := range finalSelections {
		if item.OrderID != firstDeliveredOrderID && item.DeliveryStatus == enum.OrderDeliveryStatusUnfinished {
			unfinishedFound = true
		}
	}
	if !unfinishedFound {
		t.Fatalf("expected unsubmitted selected order to be retained as unfinished, got %#v", finalSelections)
	}
}

func TestOrderMarketDisabledRequiresZeroInvestment(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	yearNo := int(nextIntegrationUniqueSeed()%100000) + 200000
	now := time.Now()

	ensureIntegrationGameConfig(t, ctx, tx, yearNo, yearNo, true)
	isolateOrderIntegrationGroups(t, ctx, tx)

	groupID := createOrderIntegrationGroup(t, ctx, tx, "I4 disabled market group", enum.BusinessStatusNormal, nil)
	createPreviousFormalReportRecord(t, ctx, tx, groupID, yearNo-1, now)
	createGroupYearStateRecord(t, ctx, tx, groupID, yearNo, enum.YearTypeFormal, enum.YearStatusOperating, enum.StageStatusQ1Open, enum.ReportStatusLocked)

	adminOrderService := NewAdminOrderCommandService(tx)
	playerOrderService := NewPlayerOrderCommandService(tx)

	if _, err := adminOrderService.UpdateMarketConfig(ctx, UpdateOrderMarketConfigCommand{
		YearNo: yearNo,
		Markets: []UpdateOrderMarketConfigItem{
			{MarketCode: enum.MarketCodeLocal, Enabled: true},
			{MarketCode: enum.MarketCodeRegional, Enabled: false},
			{MarketCode: enum.MarketCodeNational, Enabled: false},
			{MarketCode: enum.MarketCodeGlobal, Enabled: false},
		},
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("update market config: %v", err)
	}
	if _, err := adminOrderService.UpdateControlConfig(ctx, UpdateOrderControlConfigCommand{
		YearNo: yearNo,
		Items: buildTestControlConfigItems(yearNo, map[string]int{
			testOrderSegmentKey(enum.MarketCodeLocal, enum.OrderTypeAgencyInspection):    1,
			testOrderSegmentKey(enum.MarketCodeRegional, enum.OrderTypeAgencyInspection): 1,
		}),
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("update order config: %v", err)
	}
	generateResult, err := adminOrderService.GenerateOrderPool(ctx, GenerateOrderPoolCommand{
		YearNo:       yearNo,
		Overwrite:    true,
		OperatorID:   1,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("generate order pool: %v", err)
	}
	if generateResult.GeneratedCount != 1 {
		t.Fatalf("expected disabled regional market not to generate orders, got %#v", generateResult)
	}
	if _, err := adminOrderService.ConfirmOrderPool(ctx, ConfirmOrderPoolCommand{
		YearNo:       yearNo,
		BatchID:      generateResult.BatchID,
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("confirm order pool: %v", err)
	}

	invalidInvestments := buildTestMarketInvestments(func(segment OrderSegmentDefinition) float64 {
		if segment.MarketCode == enum.MarketCodeRegional && segment.OrderType == enum.OrderTypeAgencyInspection {
			return 1
		}
		return 0
	})
	if _, err := playerOrderService.SubmitMarketInvestment(ctx, SubmitMarketInvestmentCommand{
		GroupID:      groupID,
		YearNo:       yearNo,
		Investments:  invalidInvestments,
		OperatorName: "group-disabled",
	}); !errors.Is(err, ErrOrderMarketDisabledInvestment) {
		t.Fatalf("expected disabled market non-zero investment to be rejected, got %v", err)
	}
}

func isolateOrderIntegrationGroups(t *testing.T, ctx context.Context, tx *gorm.DB) {
	t.Helper()

	now := time.Now()
	if err := tx.WithContext(ctx).
		Model(&entity.Group{}).
		Where("1 = 1").
		Updates(map[string]any{
			"business_status": enum.BusinessStatusBankrupt,
			"updater":         "integration-test",
			"update_time":     now,
		}).Error; err != nil {
		t.Fatalf("isolate existing groups: %v", err)
	}
}

func createOrderIntegrationGroup(t *testing.T, ctx context.Context, tx *gorm.DB, name string, businessStatus string, bankruptYearNo *int) int64 {
	t.Helper()

	now := time.Now()
	uniqueSeed := nextIntegrationUniqueSeed()
	group := entity.Group{
		GroupNo:        integrationGroupNoFromSeed(uniqueSeed),
		GroupCode:      fmt.Sprintf("IT_ORDER_%d", uniqueSeed),
		GroupName:      name,
		BusinessStatus: businessStatus,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if businessStatus == enum.BusinessStatusBankrupt && bankruptYearNo != nil {
		year := *bankruptYearNo
		reason := "integration-test"
		group.BankruptYearNo = &year
		group.BankruptReason = &reason
	}
	if err := tx.WithContext(ctx).Create(&group).Error; err != nil {
		t.Fatalf("create order integration group: %v", err)
	}
	return group.ID
}

func createSelectedOrderAmountFixture(t *testing.T, ctx context.Context, tx *gorm.DB, groupID int64, yearNo int, marketCode string, amount float64) {
	t.Helper()

	now := time.Now()
	sourceIndex := 1
	selectedGroupID := groupID
	selectedAt := now
	item := entity.OrderPool{
		YearNo:          yearNo,
		MarketCode:      marketCode,
		OrderType:       enum.OrderTypeAgencyInspection,
		SegmentCode:     marketCode + "_" + enum.OrderTypeAgencyInspection,
		CardSequenceNo:  1,
		OrderAmount:     amount,
		OrderQuantity:   1,
		UnitPrice:       amount,
		AccountTerm:     2,
		PoolStatus:      enum.OrderPoolStatusSelected,
		SelectedGroupID: &selectedGroupID,
		SelectedAt:      &selectedAt,
		SourceSheetName: "integration-test",
		SourceCell:      "PREV",
		SourceRowIndex:  &sourceIndex,
		SourceRowKey:    fmt.Sprintf("prev_%d_%s_%d", yearNo, marketCode, groupID),
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		t.Fatalf("create previous selected order amount fixture: %v", err)
	}
}

func buildOrderLinkedOperatingCommandService(db *gorm.DB) *PlayerOperatingCommandService {
	gameConfigRepo := repository.NewGameConfigRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	groupYearRepo := repository.NewGroupYearStateRepository(db)
	operatingRepo := repository.NewOperatingRepository(db)
	initialBaselineRepo := repository.NewInitialBaselineRepository(db)
	reportRepo := repository.NewReportRepository(db)
	noticeRepo := repository.NewNoticeRepository(db)
	adjustmentRepo := repository.NewGroupAdjustmentRepository(db)
	playerNoticeService := NewPlayerNoticeService(noticeRepo, adjustmentRepo)
	orderLinkService := NewOrderOperatingLinkService(
		repository.NewGroupMarketBidRepository(db),
		repository.NewMarketBiddingStateRepository(db),
		repository.NewGroupOrderSelectionRepository(db),
	)
	return NewPlayerOperatingCommandService(
		db,
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaselineRepo,
		reportRepo,
		playerNoticeService,
		orderLinkService,
	)
}

func saveOperatingDraft(t *testing.T, ctx context.Context, tx *gorm.DB, groupID int64, yearNo int, stageStatus string, operatingPayload payload.OperatingPayload) {
	t.Helper()

	if err := repository.NewOperatingRepository(tx).UpsertDraft(ctx, repository.UpsertOperatingDraftCommand{
		GroupID:          groupID,
		YearNo:           yearNo,
		StageStatus:      stageStatus,
		OperatingPayload: operatingPayload,
		LastAutoSavedAt:  time.Now(),
		OperatorName:     "integration-test",
	}); err != nil {
		t.Fatalf("save operating draft fixture: %v", err)
	}
}

func buildQ1PayloadWithRevenue(salesRevenue float64) payload.OperatingPayload {
	p := buildValidQ1OperatingPayload()
	p.Quarter.DeliverySettlement["q1"]["salesRevenue"] = salesRevenue
	return p.Normalize()
}

func buildYearEndOperatingPayload() payload.OperatingPayload {
	p := payload.NewOperatingPayload()
	p.YearEnd.LongTermLoan = map[string]any{
		"longTermInterest":  0,
		"longTermRepayment": 0,
		"newLongTermLoan":   0,
	}
	p.YearEnd.AssetAdjustment = map[string]any{
		"lineMaintenance":    0,
		"factoryPurchase":    0,
		"factorySale":        0,
		"factoryRent":        0,
		"workInConstruction": 0,
		"marketCultivation":  0,
	}
	return p.Normalize()
}

func setGroupYearStage(t *testing.T, ctx context.Context, tx *gorm.DB, groupID int64, yearNo int, yearStatus string, stageStatus string, reportStatus string, latestStageVersion int) {
	t.Helper()

	if err := tx.WithContext(ctx).
		Model(&entity.GroupYearState{}).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Updates(map[string]any{
			"year_status":                 yearStatus,
			"stage_status":                stageStatus,
			"report_status":               reportStatus,
			"latest_stage_submit_version": latestStageVersion,
			"updater":                     "integration-test",
			"update_time":                 time.Now(),
		}).Error; err != nil {
		t.Fatalf("set group year stage: %v", err)
	}
}

func orderAmountByID(t *testing.T, knownOrders []entity.OrderPool, orderID int64) float64 {
	t.Helper()

	for _, item := range knownOrders {
		if item.ID == orderID {
			return item.OrderAmount
		}
	}
	t.Fatalf("order %d not found in %#v", orderID, knownOrders)
	return 0
}
