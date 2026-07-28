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
	yearNo := 2
	previousYearNo := yearNo - 1
	now := time.Now()

	ensureIntegrationGameConfig(t, ctx, tx, yearNo, yearNo, true)
	cleanupOrderIntegrationYears(t, ctx, tx, previousYearNo, yearNo)
	isolateOrderIntegrationGroups(t, ctx, tx)

	groupOneID := createOrderIntegrationGroup(t, ctx, tx, "I4 order leader", enum.BusinessStatusNormal, nil)
	groupTwoID := createOrderIntegrationGroup(t, ctx, tx, "I4 order follower", enum.BusinessStatusNormal, nil)
	bankruptYear := previousYearNo
	bankruptGroupID := createOrderIntegrationGroup(t, ctx, tx, "I4 bankrupt previous winner", enum.BusinessStatusBankrupt, &bankruptYear)

	createPreviousFormalReportRecord(t, ctx, tx, groupOneID, previousYearNo, now)
	createPreviousFormalReportRecord(t, ctx, tx, groupTwoID, previousYearNo, now)
	createGroupYearStateRecord(t, ctx, tx, groupOneID, yearNo, enum.YearTypeFormal, enum.YearStatusOperating, enum.StageStatusQ1Open, enum.ReportStatusLocked)
	createGroupYearStateRecord(t, ctx, tx, groupTwoID, yearNo, enum.YearTypeFormal, enum.YearStatusOperating, enum.StageStatusQ1Open, enum.ReportStatusLocked)
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

	if _, err := adminOrderService.UpdateForecastControl(ctx, UpdateOrderForecastControlCommand{
		Items: buildTestForecastControlItems(map[string]int{
			testForecastSegmentKey(yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection):    2,
			testForecastSegmentKey(yearNo, enum.MarketCodeRegional, enum.OrderTypeAgencyInspection): 1,
			testForecastSegmentKey(yearNo, enum.MarketCodeLocal, enum.OrderTypeTwoCabinVIP):         1,
		}),
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("update forecast control: %v", err)
	}

	configResult, err := adminOrderService.UpdateControlConfig(ctx, UpdateOrderControlConfigCommand{
		YearNo:       yearNo,
		Items:        buildTestControlConfigItems(yearNo, nil),
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

	// 预测配置会重建仍处于预览态的历史年份订单池，因此在当前年份订单池确认后再布置
	// 上一年已选订单事实，确保本用例验证的是龙头汇总规则而不是预览覆盖行为。
	createSelectedOrderAmountFixture(t, ctx, tx, groupOneID, previousYearNo, enum.MarketCodeLocal, 50)
	createSelectedOrderAmountFixture(t, ctx, tx, groupTwoID, previousYearNo, enum.MarketCodeLocal, 20)
	createSelectedOrderAmountFixture(t, ctx, tx, bankruptGroupID, previousYearNo, enum.MarketCodeLocal, 999)

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
	if _, err := adminControlService.GenerateSelectionSequence(ctx, GenerateSelectionSequenceCommand{
		YearNo: yearNo, OperatorID: 1, OperatorName: "integration-admin",
	}); !errors.Is(err, ErrAdminOrderSequenceAlreadyGenerated) {
		t.Fatalf("expected repeated sequence generation to be rejected without reordering, got %v", err)
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
	if len(localAgencySequence) != 8 {
		t.Fatalf("expected two qualified groups across four local agency rounds, got %#v", localAgencySequence)
	}
	for roundNo := 1; roundNo <= 4; roundNo++ {
		offset := (roundNo - 1) * 2
		if localAgencySequence[offset].RoundNo != roundNo || localAgencySequence[offset].GroupID != groupOneID || !localAgencySequence[offset].IsMarketLeader ||
			localAgencySequence[offset+1].RoundNo != roundNo || localAgencySequence[offset+1].GroupID != groupTwoID {
			t.Fatalf("expected every local agency round to filter the same leader-first base order, got %#v", localAgencySequence)
		}
	}
	localTwoCabinSequence, err := sequenceRepo.ListBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeTwoCabinVIP)
	if err != nil {
		t.Fatalf("list local two-cabin sequence: %v", err)
	}
	if len(localTwoCabinSequence) != 2 || localTwoCabinSequence[0].RoundNo != 1 || localTwoCabinSequence[1].RoundNo != 2 ||
		localTwoCabinSequence[0].GroupID != groupTwoID || localTwoCabinSequence[0].IsMarketLeader {
		t.Fatalf("expected zero-investment leader to be excluded and group two to qualify for two rounds, got %#v", localTwoCabinSequence)
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
	if completedAgencyState.SegmentStatus != enum.OrderSegmentStatusRoundReady || completedAgencyState.CurrentRoundNo != 1 {
		t.Fatalf("expected first round completion to wait for the administrator to open round two, got %#v", completedAgencyState)
	}
	localAgencyOrders, err = poolRepo.ListBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("reload local agency pool: %v", err)
	}
	if localAgencyOrders[0].PoolStatus != enum.OrderPoolStatusSelected || localAgencyOrders[1].PoolStatus != enum.OrderPoolStatusAvailable {
		t.Fatalf("expected selected order locked and unchosen order shared with the next round, got %#v", localAgencyOrders)
	}
	adminControlQueryService := NewAdminOrderControlQueryService(
		repository.NewGameConfigRepository(tx),
		repository.NewGroupRepository(tx),
		repository.NewGroupMarketBidRepository(tx),
		repository.NewMarketBiddingStateRepository(tx),
		repository.NewMarketSelectionOrderRepository(tx),
		repository.NewOrderPoolRepository(tx),
	)
	selectionStatus, err := adminControlQueryService.GetMarketSelectionStatus(ctx, yearNo, enum.MarketCodeLocal)
	if err != nil {
		t.Fatalf("query local market selection status after round one: %v", err)
	}
	if selectionStatus.CurrentSegment == nil {
		t.Fatalf("expected round-ready segment to remain the current operable segment")
	}
	if selectionStatus.CurrentSegment.SegmentStatus != enum.OrderSegmentStatusRoundReady {
		t.Fatalf("expected current segment status round-ready, got %#v", selectionStatus.CurrentSegment)
	}
	if selectionStatus.CurrentSegment.NextRoundNo == nil || *selectionStatus.CurrentSegment.NextRoundNo != 2 {
		t.Fatalf("expected round-ready segment to expose next round number, got %#v", selectionStatus.CurrentSegment)
	}
	openedRound, err := adminControlService.OpenNextRound(ctx, OpenNextOrderRoundCommand{
		YearNo:       yearNo,
		MarketCode:   enum.MarketCodeLocal,
		OrderType:    enum.OrderTypeAgencyInspection,
		OperatorID:   1,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("open local agency round two: %v", err)
	}
	if openedRound.SegmentStatus != enum.OrderSegmentStatusSelecting || openedRound.CurrentRoundNo != 2 || openedRound.CurrentGroupID == nil || *openedRound.CurrentGroupID != groupOneID {
		t.Fatalf("expected round two to start from the same leader-first base order, got %#v", openedRound)
	}
	if _, err := playerOrderService.SelectOrder(ctx, SelectOrderCommand{
		GroupID:      groupOneID,
		YearNo:       yearNo,
		MarketCode:   enum.MarketCodeLocal,
		OrderType:    enum.OrderTypeAgencyInspection,
		OrderID:      localAgencyOrders[1].ID,
		OperatorName: "group-one",
	}); err != nil {
		t.Fatalf("group one select final local agency order in round two: %v", err)
	}
	completedAgencyState, err = stateRepo.GetBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("reload exhausted local agency segment: %v", err)
	}
	if completedAgencyState.SegmentStatus != enum.OrderSegmentStatusCompleted || completedAgencyState.CompletionReason == nil || *completedAgencyState.CompletionReason != enum.OrderCompletionReasonPoolExhausted {
		t.Fatalf("expected order-pool exhaustion to complete the whole segment immediately, got %#v", completedAgencyState)
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
	if releasedSecond.MarketCode != enum.MarketCodeLocal || releasedSecond.OrderType != enum.OrderTypeTwoCabinVIP || releasedSecond.CurrentGroupID == nil || *releasedSecond.CurrentGroupID != groupTwoID {
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
		GroupID:      groupTwoID,
		YearNo:       yearNo,
		MarketCode:   enum.MarketCodeLocal,
		OrderType:    enum.OrderTypeTwoCabinVIP,
		OrderID:      twoCabinOrders[0].ID,
		OperatorName: "group-two",
	}); err != nil {
		t.Fatalf("group two select local two-cabin order: %v", err)
	}
	twoCabinAfterSelect, err := stateRepo.GetBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeTwoCabinVIP)
	if err != nil {
		t.Fatalf("reload local two-cabin after first selection: %v", err)
	}
	if twoCabinAfterSelect.SegmentStatus != enum.OrderSegmentStatusCompleted || twoCabinAfterSelect.CompletionReason == nil || *twoCabinAfterSelect.CompletionReason != enum.OrderCompletionReasonPoolExhausted {
		t.Fatalf("expected the final available order to complete the segment immediately, got %#v", twoCabinAfterSelect)
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

func TestOrderPassRetainsLaterRoundAndExpiresRemainingPool(t *testing.T) {
	db := openIntegrationMySQL(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() { _ = tx.Rollback().Error }()

	ctx := context.Background()
	yearNo := 1
	ensureIntegrationGameConfig(t, ctx, tx, yearNo, yearNo, true)
	cleanupOrderIntegrationYears(t, ctx, tx, yearNo)
	isolateOrderIntegrationGroups(t, ctx, tx)
	groupID := createOrderIntegrationGroup(t, ctx, tx, "I13 pass later round", enum.BusinessStatusNormal, nil)
	now := time.Now()
	template := MustResolveOrderTemplateVersion(OrderTemplateVersionVIPServiceV1)
	currentGroupID := groupID
	stateItem := entity.MarketBiddingState{
		OrderTemplateVersion: template.TemplateVersion,
		YearNo:               yearNo,
		MarketCode:           enum.MarketCodeLocal,
		OrderType:            enum.OrderTypeAgencyInspection,
		SegmentCode:          enum.MarketCodeLocal + "_" + enum.OrderTypeAgencyInspection,
		ReleaseSequenceNo:    1,
		SegmentStatus:        enum.OrderSegmentStatusSelecting,
		CurrentGroupID:       &currentGroupID,
		CurrentRoundNo:       1,
		BaseEntity: entity.BaseEntity{
			Creator: "integration-test", CreateTime: now, Updater: "integration-test", UpdateTime: now,
		},
	}
	if err := tx.Create(&stateItem).Error; err != nil {
		t.Fatalf("create round state: %v", err)
	}
	sequenceItems := []entity.MarketSelectionOrder{
		{OrderTemplateVersion: template.TemplateVersion, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, RoundNo: 1, SequenceNo: 1, GroupID: groupID, SelectionStatus: enum.OrderSelectionStatusCurrent, RankBasis: []byte(`{}`), BaseEntity: entity.BaseEntity{Creator: "integration-test", CreateTime: now, Updater: "integration-test", UpdateTime: now}},
		{OrderTemplateVersion: template.TemplateVersion, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, RoundNo: 2, SequenceNo: 1, GroupID: groupID, SelectionStatus: enum.OrderSelectionStatusWaiting, RankBasis: []byte(`{}`), BaseEntity: entity.BaseEntity{Creator: "integration-test", CreateTime: now, Updater: "integration-test", UpdateTime: now}},
	}
	if err := tx.Create(&sequenceItems).Error; err != nil {
		t.Fatalf("create round sequence: %v", err)
	}
	for cardNo := 1; cardNo <= 2; cardNo++ {
		if err := tx.Create(&entity.OrderPool{
			OrderTemplateVersion: template.TemplateVersion,
			YearNo:               yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection,
			SegmentCode:    enum.MarketCodeLocal + "_" + enum.OrderTypeAgencyInspection,
			CardSequenceNo: cardNo, BusinessOrderNo: fmt.Sprintf("I13-PASS-%02d", cardNo),
			OrderAmount: float64(cardNo * 10), OrderQuantity: 1, UnitPrice: float64(cardNo * 10), AccountTerm: 2,
			PoolStatus: enum.OrderPoolStatusAvailable, SourceSheetName: "integration-test", SourceCell: "A1", SourceRowKey: fmt.Sprintf("I13-PASS-%d", cardNo),
			BaseEntity: entity.BaseEntity{Creator: "integration-test", CreateTime: now, Updater: "integration-test", UpdateTime: now},
		}).Error; err != nil {
			t.Fatalf("create round pool card %d: %v", cardNo, err)
		}
	}

	playerService := NewPlayerOrderCommandService(tx)
	passed, err := playerService.PassSegment(ctx, PassOrderSegmentCommand{GroupID: groupID, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, OperatorName: "group-pass"})
	if err != nil {
		t.Fatalf("pass first round: %v", err)
	}
	if passed.RoundNo != 1 || passed.SegmentStatus != enum.OrderSegmentStatusRoundReady {
		t.Fatalf("expected pass to affect only round one and wait for round two, got %#v", passed)
	}
	stateRepo := repository.NewMarketBiddingStateRepository(tx)
	stateAfterPass, err := stateRepo.GetBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("reload state after pass: %v", err)
	}
	if stateAfterPass.CurrentRoundNo != 1 || stateAfterPass.CurrentGroupID != nil {
		t.Fatalf("expected round-ready state to clear current group, got %#v", stateAfterPass)
	}
	if _, err := playerService.PassSegment(ctx, PassOrderSegmentCommand{GroupID: groupID, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, OperatorName: "group-pass"}); !errors.Is(err, ErrOrderSegmentNotSelecting) {
		t.Fatalf("expected repeated pass while round is waiting to be rejected, got %v", err)
	}

	adminService := NewAdminOrderControlCommandService(tx)
	opened, err := adminService.OpenNextRound(ctx, OpenNextOrderRoundCommand{YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, OperatorID: 1, OperatorName: "integration-admin"})
	if err != nil {
		t.Fatalf("open second round after player pass: %v", err)
	}
	if opened.CurrentRoundNo != 2 || opened.SegmentStatus != enum.OrderSegmentStatusSelecting || opened.CurrentGroupID == nil || *opened.CurrentGroupID != groupID {
		t.Fatalf("expected same group to retain second-round eligibility, got %#v", opened)
	}
	if _, err := adminService.OpenNextRound(ctx, OpenNextOrderRoundCommand{YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, OperatorID: 1, OperatorName: "integration-admin"}); !errors.Is(err, ErrOrderRoundNotReady) {
		t.Fatalf("expected repeated next-round opening to be rejected, got %v", err)
	}

	completed, err := playerService.PassSegment(ctx, PassOrderSegmentCommand{GroupID: groupID, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, OperatorName: "group-pass"})
	if err != nil {
		t.Fatalf("pass final round: %v", err)
	}
	if completed.SegmentStatus != enum.OrderSegmentStatusCompleted {
		t.Fatalf("expected final eligible round to complete segment, got %#v", completed)
	}
	finalState, err := stateRepo.GetBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("reload completed state: %v", err)
	}
	if finalState.CompletionReason == nil || *finalState.CompletionReason != enum.OrderCompletionReasonAllRoundsCompleted {
		t.Fatalf("expected normal completion reason, got %#v", finalState.CompletionReason)
	}
	var expiredCount int64
	if err := tx.Model(&entity.OrderPool{}).Where("year_no = ? AND market_code = ? AND order_type = ? AND pool_status = ?", yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection, enum.OrderPoolStatusUnselectedExpired).Count(&expiredCount).Error; err != nil {
		t.Fatalf("count expired order pool: %v", err)
	}
	if expiredCount != 2 {
		t.Fatalf("expected all remaining pool cards to expire after normal completion, got %d", expiredCount)
	}
}

func TestOrderBankruptcyInvalidatesOnlyUnfinishedRounds(t *testing.T) {
	db := openIntegrationMySQL(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() { _ = tx.Rollback().Error }()

	ctx := context.Background()
	yearNo := 1
	ensureIntegrationGameConfig(t, ctx, tx, yearNo, yearNo, true)
	cleanupOrderIntegrationYears(t, ctx, tx, yearNo)
	isolateOrderIntegrationGroups(t, ctx, tx)
	groupID := createOrderIntegrationGroup(t, ctx, tx, "I13 bankrupt round", enum.BusinessStatusBankrupt, ptrIntForOrderTest(yearNo))
	now := time.Now()
	template := MustResolveOrderTemplateVersion(OrderTemplateVersionVIPServiceV1)
	currentGroupID := groupID
	stateItem := entity.MarketBiddingState{OrderTemplateVersion: template.TemplateVersion, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, SegmentCode: enum.MarketCodeLocal + "_" + enum.OrderTypeAgencyInspection, ReleaseSequenceNo: 1, SegmentStatus: enum.OrderSegmentStatusSelecting, CurrentGroupID: &currentGroupID, CurrentRoundNo: 2, BaseEntity: entity.BaseEntity{Creator: "integration-test", CreateTime: now, Updater: "integration-test", UpdateTime: now}}
	if err := tx.Create(&stateItem).Error; err != nil {
		t.Fatalf("create bankruptcy state: %v", err)
	}
	selectedSequence := entity.MarketSelectionOrder{OrderTemplateVersion: template.TemplateVersion, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, RoundNo: 1, SequenceNo: 1, GroupID: groupID, SelectionStatus: enum.OrderSelectionStatusSelected, RankBasis: []byte(`{}`), BaseEntity: entity.BaseEntity{Creator: "integration-test", CreateTime: now, Updater: "integration-test", UpdateTime: now}}
	currentSequence := entity.MarketSelectionOrder{OrderTemplateVersion: template.TemplateVersion, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, RoundNo: 2, SequenceNo: 1, GroupID: groupID, SelectionStatus: enum.OrderSelectionStatusCurrent, RankBasis: []byte(`{}`), BaseEntity: entity.BaseEntity{Creator: "integration-test", CreateTime: now, Updater: "integration-test", UpdateTime: now}}
	futureSequence := entity.MarketSelectionOrder{OrderTemplateVersion: template.TemplateVersion, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, RoundNo: 3, SequenceNo: 1, GroupID: groupID, SelectionStatus: enum.OrderSelectionStatusWaiting, RankBasis: []byte(`{}`), BaseEntity: entity.BaseEntity{Creator: "integration-test", CreateTime: now, Updater: "integration-test", UpdateTime: now}}
	if err := tx.Create(&[]entity.MarketSelectionOrder{selectedSequence, currentSequence, futureSequence}).Error; err != nil {
		t.Fatalf("create bankruptcy sequences: %v", err)
	}
	selectedGroupID := groupID
	selectedAt := now
	selectedPool := entity.OrderPool{OrderTemplateVersion: template.TemplateVersion, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, SegmentCode: enum.MarketCodeLocal + "_" + enum.OrderTypeAgencyInspection, CardSequenceNo: 1, BusinessOrderNo: "I13-BANKRUPT-SELECTED", OrderAmount: 30, OrderQuantity: 1, UnitPrice: 30, AccountTerm: 2, PoolStatus: enum.OrderPoolStatusSelected, SelectedGroupID: &selectedGroupID, SelectedAt: &selectedAt, SourceSheetName: "integration-test", SourceCell: "A1", SourceRowKey: "I13-BANKRUPT-SELECTED", BaseEntity: entity.BaseEntity{Creator: "integration-test", CreateTime: now, Updater: "integration-test", UpdateTime: now}}
	if err := tx.Create(&selectedPool).Error; err != nil {
		t.Fatalf("create selected bankruptcy order: %v", err)
	}
	selectedOrderID := selectedPool.ID
	selectionOrderID := selectedSequence.ID
	if err := tx.Model(&entity.MarketSelectionOrder{}).Where("id = ?", selectedSequence.ID).Update("selected_order_id", selectedOrderID).Error; err != nil {
		t.Fatalf("update selected sequence fixture: %v", err)
	}
	if err := tx.Create(&entity.GroupOrderSelection{OrderTemplateVersion: template.TemplateVersion, GroupID: groupID, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, RoundNo: 1, SelectionOrderID: &selectionOrderID, OrderID: selectedOrderID, SelectionStatus: enum.OrderSelectionStatusSelected, DeliveryStatus: enum.OrderDeliveryStatusSelected, DeliveryEffective: true, SelectedAt: now, BaseEntity: entity.BaseEntity{Creator: "integration-test", CreateTime: now, Updater: "integration-test", UpdateTime: now}}).Error; err != nil {
		t.Fatalf("create selected order history: %v", err)
	}
	if err := tx.Create(&entity.OrderPool{OrderTemplateVersion: template.TemplateVersion, YearNo: yearNo, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, SegmentCode: enum.MarketCodeLocal + "_" + enum.OrderTypeAgencyInspection, CardSequenceNo: 2, BusinessOrderNo: "I13-BANKRUPT-AVAILABLE", OrderAmount: 20, OrderQuantity: 1, UnitPrice: 20, AccountTerm: 2, PoolStatus: enum.OrderPoolStatusAvailable, SourceSheetName: "integration-test", SourceCell: "A2", SourceRowKey: "I13-BANKRUPT-AVAILABLE", BaseEntity: entity.BaseEntity{Creator: "integration-test", CreateTime: now, Updater: "integration-test", UpdateTime: now}}).Error; err != nil {
		t.Fatalf("create available bankruptcy order: %v", err)
	}
	if err := invalidateOrderParticipationAfterBankruptcy(ctx, tx, groupID, yearNo, 1, "bankruptcy-test", now); err != nil {
		t.Fatalf("invalidate unfinished bankruptcy rounds: %v", err)
	}
	stateRepo := repository.NewMarketBiddingStateRepository(tx)
	finalState, err := stateRepo.GetBySegment(ctx, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection)
	if err != nil {
		t.Fatalf("reload bankruptcy state: %v", err)
	}
	if finalState.SegmentStatus != enum.OrderSegmentStatusCompleted || finalState.CompletionReason == nil || *finalState.CompletionReason != enum.OrderCompletionReasonNoEligibleParticipants {
		t.Fatalf("expected no-eligible-participants completion after bankruptcy, got %#v", finalState)
	}
	var statuses []string
	if err := tx.Model(&entity.MarketSelectionOrder{}).Where("year_no = ? AND market_code = ? AND order_type = ?", yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection).Order("round_no ASC").Pluck("selection_status", &statuses).Error; err != nil {
		t.Fatalf("load bankruptcy sequence statuses: %v", err)
	}
	if len(statuses) != 3 || statuses[0] != enum.OrderSelectionStatusSelected || statuses[1] != enum.OrderSelectionStatusBankrupt || statuses[2] != enum.OrderSelectionStatusBankrupt {
		t.Fatalf("expected selected history preserved and unfinished rounds invalidated, got %#v", statuses)
	}
	var selectedHistoryCount int64
	if err := tx.Model(&entity.GroupOrderSelection{}).Where("group_id = ? AND year_no = ? AND order_id = ? AND delivery_effective = ?", groupID, yearNo, selectedOrderID, true).Count(&selectedHistoryCount).Error; err != nil {
		t.Fatalf("count selected order history: %v", err)
	}
	if selectedHistoryCount != 1 {
		t.Fatalf("expected previously selected order history to remain effective, got %d", selectedHistoryCount)
	}
	var availableCount, expiredCount int64
	if err := tx.Model(&entity.OrderPool{}).Where("year_no = ? AND pool_status = ?", yearNo, enum.OrderPoolStatusAvailable).Count(&availableCount).Error; err != nil {
		t.Fatalf("count available bankruptcy orders: %v", err)
	}
	if err := tx.Model(&entity.OrderPool{}).Where("year_no = ? AND pool_status = ?", yearNo, enum.OrderPoolStatusUnselectedExpired).Count(&expiredCount).Error; err != nil {
		t.Fatalf("count expired bankruptcy orders: %v", err)
	}
	if availableCount != 0 || expiredCount != 1 {
		t.Fatalf("expected remaining available pool to expire after bankruptcy, available=%d expired=%d", availableCount, expiredCount)
	}
}

func ptrIntForOrderTest(value int) *int {
	return &value
}

func TestOrderYearLockRejectsHistoricalForecastControlRefresh(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	historicalYearNo := 1
	currentOpenYear := 2

	ensureIntegrationGameConfig(t, ctx, tx, currentOpenYear, currentOpenYear, true)
	cleanupOrderIntegrationYears(t, ctx, tx, historicalYearNo, currentOpenYear)
	isolateOrderIntegrationGroups(t, ctx, tx)

	groupID := createOrderIntegrationGroup(t, ctx, tx, "I13 historical order owner", enum.BusinessStatusNormal, nil)
	createSelectedOrderAmountFixture(t, ctx, tx, groupID, historicalYearNo, enum.MarketCodeLocal, 88)

	adminOrderService := NewAdminOrderCommandService(tx)
	if _, err := adminOrderService.UpdateForecastControl(ctx, UpdateOrderForecastControlCommand{
		Items: buildTestForecastControlItems(map[string]int{
			testForecastSegmentKey(historicalYearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection): 1,
		}),
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); !errors.Is(err, ErrAdminOrderPoolLocked) {
		t.Fatalf("expected historical year forecast refresh to be locked, got %v", err)
	}

	var selectedCount int64
	if err := tx.WithContext(ctx).
		Model(&entity.OrderPool{}).
		Where("year_no = ? AND pool_status = ? AND selected_group_id = ?", historicalYearNo, enum.OrderPoolStatusSelected, groupID).
		Count(&selectedCount).Error; err != nil {
		t.Fatalf("count selected historical orders: %v", err)
	}
	if selectedCount != 1 {
		t.Fatalf("expected historical selected order to be preserved, got %d", selectedCount)
	}
}

func TestOrderYearLockProtectsSelectedPoolFromPreviewOverwrite(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	yearNo := 1

	ensureIntegrationGameConfig(t, ctx, tx, yearNo, yearNo, true)
	cleanupOrderIntegrationYears(t, ctx, tx, yearNo)
	isolateOrderIntegrationGroups(t, ctx, tx)

	groupID := createOrderIntegrationGroup(t, ctx, tx, "I13 selected order owner", enum.BusinessStatusNormal, nil)
	createSelectedOrderAmountFixture(t, ctx, tx, groupID, yearNo, enum.MarketCodeLocal, 66)

	adminOrderService := NewAdminOrderCommandService(tx)
	if _, err := adminOrderService.GenerateOrderPool(ctx, GenerateOrderPoolCommand{
		YearNo:       yearNo,
		Overwrite:    true,
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); !errors.Is(err, ErrAdminOrderPoolLocked) {
		t.Fatalf("expected selected order pool to reject preview overwrite, got %v", err)
	}

	var selectedCount int64
	if err := tx.WithContext(ctx).
		Model(&entity.OrderPool{}).
		Where("year_no = ? AND pool_status = ? AND selected_group_id = ?", yearNo, enum.OrderPoolStatusSelected, groupID).
		Count(&selectedCount).Error; err != nil {
		t.Fatalf("count selected orders after rejected overwrite: %v", err)
	}
	if selectedCount != 1 {
		t.Fatalf("expected selected order ownership to be preserved, got %d", selectedCount)
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
	yearNo := 1
	now := time.Now()

	ensureIntegrationGameConfig(t, ctx, tx, yearNo, yearNo, true)
	cleanupOrderIntegrationYears(t, ctx, tx, yearNo)
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
	if _, err := adminOrderService.UpdateForecastControl(ctx, UpdateOrderForecastControlCommand{
		Items: buildTestForecastControlItems(map[string]int{
			testForecastSegmentKey(yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection):    1,
			testForecastSegmentKey(yearNo, enum.MarketCodeRegional, enum.OrderTypeAgencyInspection): 1,
		}),
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("update forecast control: %v", err)
	}
	if _, err := adminOrderService.UpdateControlConfig(ctx, UpdateOrderControlConfigCommand{
		YearNo:       yearNo,
		Items:        buildTestControlConfigItems(yearNo, nil),
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

func TestOrderMarketInvestmentLimitAndConfigLock(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	yearNo := 1
	now := time.Now()

	ensureIntegrationGameConfig(t, ctx, tx, yearNo, yearNo, true)
	cleanupOrderIntegrationYears(t, ctx, tx, yearNo)
	isolateOrderIntegrationGroups(t, ctx, tx)

	groupID := createOrderIntegrationGroup(t, ctx, tx, "I6 limited market group", enum.BusinessStatusNormal, nil)
	createPreviousFormalReportRecord(t, ctx, tx, groupID, yearNo-1, now)
	createGroupYearStateRecord(t, ctx, tx, groupID, yearNo, enum.YearTypeFormal, enum.YearStatusOperating, enum.StageStatusQ1Open, enum.ReportStatusLocked)

	adminOrderService := NewAdminOrderCommandService(tx)
	playerOrderService := NewPlayerOrderCommandService(tx)
	localLimit := 10.0

	if _, err := adminOrderService.UpdateMarketConfig(ctx, UpdateOrderMarketConfigCommand{
		YearNo: yearNo,
		Markets: []UpdateOrderMarketConfigItem{
			{MarketCode: enum.MarketCodeLocal, Enabled: true, MarketInvestmentLimit: &localLimit},
			{MarketCode: enum.MarketCodeRegional, Enabled: false},
			{MarketCode: enum.MarketCodeNational, Enabled: false},
			{MarketCode: enum.MarketCodeGlobal, Enabled: false},
		},
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("update market config with limit: %v", err)
	}
	if _, err := adminOrderService.UpdateForecastControl(ctx, UpdateOrderForecastControlCommand{
		Items: buildTestForecastControlItems(map[string]int{
			testForecastSegmentKey(yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection): 1,
		}),
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("update forecast control: %v", err)
	}
	if _, err := adminOrderService.UpdateControlConfig(ctx, UpdateOrderControlConfigCommand{
		YearNo:       yearNo,
		Items:        buildTestControlConfigItems(yearNo, nil),
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
	if _, err := adminOrderService.ConfirmOrderPool(ctx, ConfirmOrderPoolCommand{
		YearNo:       yearNo,
		BatchID:      generateResult.BatchID,
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); err != nil {
		t.Fatalf("confirm order pool: %v", err)
	}

	overLimitInvestments := buildTestMarketInvestments(func(segment OrderSegmentDefinition) float64 {
		if segment.MarketCode == enum.MarketCodeLocal {
			if segment.OrderType == enum.OrderTypeAgencyInspection {
				return 6
			}
			if segment.OrderType == enum.OrderTypeTwoCabinVIP {
				return 5
			}
		}
		return 0
	})
	if _, err := playerOrderService.SubmitMarketInvestment(ctx, SubmitMarketInvestmentCommand{
		GroupID:      groupID,
		YearNo:       yearNo,
		Investments:  overLimitInvestments,
		OperatorName: "group-limited",
	}); !errors.Is(err, ErrOrderInvestmentLimitExceeded) {
		t.Fatalf("expected over-limit investment to be rejected, got %v", err)
	}

	withinLimitInvestments := buildTestMarketInvestments(func(segment OrderSegmentDefinition) float64 {
		if segment.MarketCode == enum.MarketCodeLocal && segment.OrderType == enum.OrderTypeAgencyInspection {
			return 10
		}
		return 0
	})
	if _, err := playerOrderService.SubmitMarketInvestment(ctx, SubmitMarketInvestmentCommand{
		GroupID:      groupID,
		YearNo:       yearNo,
		Investments:  withinLimitInvestments,
		OperatorName: "group-limited",
	}); err != nil {
		t.Fatalf("submit within limit investments: %v", err)
	}

	if _, err := adminOrderService.UpdateMarketConfig(ctx, UpdateOrderMarketConfigCommand{
		YearNo: yearNo,
		Markets: []UpdateOrderMarketConfigItem{
			{MarketCode: enum.MarketCodeLocal, Enabled: true, MarketInvestmentLimit: &localLimit},
			{MarketCode: enum.MarketCodeRegional, Enabled: false},
			{MarketCode: enum.MarketCodeNational, Enabled: false},
			{MarketCode: enum.MarketCodeGlobal, Enabled: false},
		},
		OperatorID:   1,
		OperatorName: "integration-admin",
	}); !errors.Is(err, ErrAdminOrderMarketConfigLocked) {
		t.Fatalf("expected market config to lock after investment submission, got %v", err)
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

func cleanupOrderIntegrationYears(t *testing.T, ctx context.Context, tx *gorm.DB, yearNos ...int) {
	t.Helper()

	if len(yearNos) == 0 {
		return
	}
	targets := make([]int, 0, len(yearNos))
	seen := map[int]bool{}
	for _, yearNo := range yearNos {
		if yearNo < 0 || seen[yearNo] {
			continue
		}
		targets = append(targets, yearNo)
		seen[yearNo] = true
	}
	if len(targets) == 0 {
		return
	}
	if err := tx.WithContext(ctx).
		Model(&entity.GroupYearState{}).
		Where("year_no IN ?", targets).
		Updates(map[string]any{
			"rollback_pending":           false,
			"rollback_target_year_no":    nil,
			"rollback_target_stage_code": nil,
			"rollback_log_id":            nil,
			"updater":                    "integration-test",
			"update_time":                time.Now(),
		}).Error; err != nil {
		t.Fatalf("cleanup rollback pending flags for order integration years %v: %v", targets, err)
	}
	entities := []any{
		&entity.GroupMarketBid{},
		&entity.MarketSelectionOrder{},
		&entity.GroupOrderSelection{},
		&entity.MarketBiddingState{},
		&entity.OrderPool{},
		&entity.OrderGenerationConfig{},
		&entity.OrderMarketConfig{},
		&entity.OrderGenerationBatch{},
		&entity.OrderForecastControl{},
	}
	for _, item := range entities {
		if err := tx.WithContext(ctx).Where("year_no IN ?", targets).Delete(item).Error; err != nil {
			t.Fatalf("cleanup order integration years %v for %T: %v", targets, item, err)
		}
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
