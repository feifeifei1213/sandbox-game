package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	calcctx "sandbox-game/internal/rules/context"
	reportrules "sandbox-game/internal/rules/report"
	"sandbox-game/internal/service"
	"sandbox-game/internal/state"
)

const (
	rollbackScenarioOperatorID   int64 = 1
	rollbackScenarioOperatorName       = "I7回退演练造数"
)

type seededRollbackScenario struct {
	GroupIDs            []int64
	SnapshotCount       int64
	GroupSnapshotCount  int64
	GlobalSnapshotCount int64
	SelectionCount      int64
	DeliveredCount      int64
}

func seedRollbackScenario(ctx context.Context, db *gorm.DB) error {
	controlService := service.NewAdminControlCommandService(db)
	adminOrderService := service.NewAdminOrderCommandService(db)
	adminOrderControlService := service.NewAdminOrderControlCommandService(db)
	playerOrderService := service.NewPlayerOrderCommandService(db)
	operatingService := buildRollbackOperatingService(db)
	reportService := buildRollbackReportService(db)

	if _, err := controlService.InitializeGame(ctx, service.InitializeGameCommand{
		GroupCount:   3,
		EditionCode:  service.GameEditionVIPServiceV1,
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("initialize rollback scenario game: %w", err)
	}
	if _, err := controlService.SubmitInitialBaseline(ctx, service.SubmitInitialBaselineCommand{
		BaselinePayload: rollbackScenarioBaselinePayload(),
		OperatorID:      rollbackScenarioOperatorID,
		OperatorName:    rollbackScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("submit rollback scenario initial baseline: %w", err)
	}

	groups, err := repository.NewGroupRepository(db).ListAll(ctx)
	if err != nil {
		return fmt.Errorf("list rollback scenario groups: %w", err)
	}
	groupIDs := make([]int64, 0, len(groups))
	for _, group := range groups {
		groupIDs = append(groupIDs, group.ID)
	}
	if len(groupIDs) != 3 {
		return fmt.Errorf("rollback scenario expects 3 groups, got %d", len(groupIDs))
	}

	if err := seedRollbackCompletedHistoricalYear(ctx, db, groupIDs, 0, operatingService, reportService); err != nil {
		return err
	}
	if _, err := controlService.OpenNextYear(ctx, service.OpenNextYearCommand{
		TargetYearNo: 1,
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("open rollback scenario year 1: %w", err)
	}
	if err := seedRollbackHistoricalOrderPrerequisite(ctx, adminOrderService, 1); err != nil {
		return err
	}
	if err := seedRollbackPreviousSelectedOrders(ctx, db, groupIDs, time.Now()); err != nil {
		return err
	}
	if err := seedRollbackCompletedHistoricalYear(ctx, db, groupIDs, 1, operatingService, reportService); err != nil {
		return err
	}
	if _, err := controlService.OpenNextYear(ctx, service.OpenNextYearCommand{
		TargetYearNo: 2,
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("open rollback scenario year 2: %w", err)
	}
	if err := seedRollbackYearTwoOrders(ctx, db, groupIDs, adminOrderService, adminOrderControlService, playerOrderService); err != nil {
		return err
	}
	if err := seedRollbackYearTwoOperations(ctx, db, groupIDs, playerOrderService, operatingService, reportService); err != nil {
		return err
	}
	if err := ensureRollbackScenarioSnapshotCoverage(ctx, db, groupIDs, []int{0, 1, 2}); err != nil {
		return err
	}

	var seeded seededRollbackScenario
	seeded.GroupIDs = groupIDs
	if err := db.WithContext(ctx).Model(&entity.StateSnapshot{}).Count(&seeded.SnapshotCount).Error; err != nil {
		return fmt.Errorf("count rollback scenario snapshots: %w", err)
	}
	if err := db.WithContext(ctx).Model(&entity.StateSnapshot{}).
		Where("snapshot_scope = ?", enum.SnapshotScopeGroup).
		Count(&seeded.GroupSnapshotCount).Error; err != nil {
		return fmt.Errorf("count rollback scenario group snapshots: %w", err)
	}
	if err := db.WithContext(ctx).Model(&entity.StateSnapshot{}).
		Where("snapshot_scope = ?", enum.SnapshotScopeGlobal).
		Count(&seeded.GlobalSnapshotCount).Error; err != nil {
		return fmt.Errorf("count rollback scenario global snapshots: %w", err)
	}
	if err := db.WithContext(ctx).Model(&entity.GroupOrderSelection{}).Count(&seeded.SelectionCount).Error; err != nil {
		return fmt.Errorf("count rollback scenario selections: %w", err)
	}
	if err := db.WithContext(ctx).
		Model(&entity.GroupOrderSelection{}).
		Where("delivery_status = ? AND delivery_effective = ?", enum.OrderDeliveryStatusDelivered, true).
		Count(&seeded.DeliveredCount).Error; err != nil {
		return fmt.Errorf("count rollback scenario delivered orders: %w", err)
	}

	fmt.Printf(
		"rollback scenario seeded: groups=%d groupIDs=%v snapshots=%d groupSnapshots=%d globalSnapshots=%d selections=%d delivered=%d stop=2年已完成、尚未开放3年，可测试0年/1年/2年跨年恢复、失效草稿保留、退回重提/年度阻断\n",
		len(seeded.GroupIDs),
		seeded.GroupIDs,
		seeded.SnapshotCount,
		seeded.GroupSnapshotCount,
		seeded.GlobalSnapshotCount,
		seeded.SelectionCount,
		seeded.DeliveredCount,
	)
	return nil
}

func seedRollbackCompletedHistoricalYear(
	ctx context.Context,
	db *gorm.DB,
	groupIDs []int64,
	yearNo int,
	operatingService *service.PlayerOperatingCommandService,
	reportService *service.PlayerReportCommandService,
) error {
	for index, groupID := range groupIDs {
		baseCost := float64(6 + yearNo*2 + index*2)
		if err := submitRollbackFullOperatingYear(ctx, operatingService, groupID, yearNo, index, baseCost, 0); err != nil {
			return err
		}
		reportManual, err := rollbackScenarioReportManualPayloadForGroup(ctx, db, groupID, yearNo)
		if err != nil {
			return err
		}
		if _, err := reportService.Submit(ctx, service.SubmitPlayerReportCommand{
			GroupID:             groupID,
			YearNo:              yearNo,
			ReportManualPayload: reportManual,
			SubmitterID:         groupID,
			OperatorName:        fmt.Sprintf("group%02d", index+1),
		}); err != nil {
			return fmt.Errorf("submit rollback historical report group=%d year=%d: %w", groupID, yearNo, err)
		}
	}
	return nil
}

func seedRollbackHistoricalOrderPrerequisite(ctx context.Context, adminOrderService *service.AdminOrderCommandService, yearNo int) error {
	if _, err := adminOrderService.UpdateForecastControl(ctx, service.UpdateOrderForecastControlCommand{
		Items:        rollbackZeroForecastControlItems(yearNo),
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("update rollback historical zero forecast control year=%d: %w", yearNo, err)
	}
	generated, err := adminOrderService.GenerateOrderPool(ctx, service.GenerateOrderPoolCommand{
		YearNo:       yearNo,
		Overwrite:    true,
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	})
	if err != nil {
		return fmt.Errorf("generate rollback historical order pool year=%d: %w", yearNo, err)
	}
	if _, err := adminOrderService.ConfirmOrderPool(ctx, service.ConfirmOrderPoolCommand{
		YearNo:       yearNo,
		BatchID:      generated.BatchID,
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("confirm rollback historical order pool year=%d: %w", yearNo, err)
	}
	return nil
}

func rollbackZeroForecastControlItems(yearNo int) []service.UpdateOrderForecastControlItem {
	items := make([]service.UpdateOrderForecastControlItem, 0, len(rollbackOrderSegments()))
	for _, segment := range rollbackOrderSegments() {
		items = append(items, service.UpdateOrderForecastControlItem{
			YearNo:     yearNo,
			MarketCode: segment.MarketCode,
			OrderType:  segment.OrderType,
			OrderCount: 0,
		})
	}
	return items
}

func ensureRollbackScenarioSnapshotCoverage(ctx context.Context, db *gorm.DB, groupIDs []int64, yearNos []int) error {
	expectedStages := []string{
		state.StageCodeQ1,
		state.StageCodeQ2,
		state.StageCodeQ3,
		state.StageCodeQ4,
		state.StageCodeYearEnd,
		"REPORT",
	}
	for _, groupID := range groupIDs {
		for _, yearNo := range yearNos {
			for _, stageCode := range expectedStages {
				var count int64
				if err := db.WithContext(ctx).Model(&entity.StateSnapshot{}).
					Where("snapshot_scope = ? AND snapshot_type = ? AND target_group_id = ? AND target_year_no = ? AND target_stage_code = ?",
						enum.SnapshotScopeGroup, enum.SnapshotTypeAuto, groupID, yearNo, stageCode).
					Count(&count).Error; err != nil {
					return fmt.Errorf("count rollback snapshot coverage group=%d year=%d stage=%s: %w", groupID, yearNo, stageCode, err)
				}
				if count == 0 {
					return fmt.Errorf("rollback scenario missing group snapshot: group=%d year=%d stage=%s", groupID, yearNo, stageCode)
				}
			}
		}
	}
	return nil
}

func seedRollbackPreviousSelectedOrders(ctx context.Context, db *gorm.DB, groupIDs []int64, now time.Time) error {
	amounts := map[string][]float64{
		enum.MarketCodeLocal:    []float64{130, 95, 80},
		enum.MarketCodeRegional: []float64{60, 125, 90},
		enum.MarketCodeNational: []float64{70, 80, 150},
		enum.MarketCodeGlobal:   []float64{140, 105, 115},
	}
	orders := make([]entity.OrderPool, 0, len(amounts)*len(groupIDs))
	for marketCode, values := range amounts {
		for index, groupID := range groupIDs {
			selectedGroupID := groupID
			selectedAt := now
			sourceIndex := index + 1
			amount := values[index]
			orders = append(orders, entity.OrderPool{
				YearNo:          1,
				MarketCode:      marketCode,
				OrderType:       enum.OrderTypeAgencyInspection,
				SegmentCode:     fmt.Sprintf("%s_%s", marketCode, enum.OrderTypeAgencyInspection),
				CardSequenceNo:  index + 1,
				BusinessOrderNo: fmt.Sprintf("RB-Y1-%s-G%d", marketCode, index+1),
				OrderAmount:     amount,
				OrderQuantity:   1,
				UnitPrice:       amount,
				AccountTerm:     2,
				PoolStatus:      enum.OrderPoolStatusSelected,
				SelectedGroupID: &selectedGroupID,
				SelectedAt:      &selectedAt,
				SourceSheetName: "I7回退演练造数",
				SourceCell:      "PREV",
				SourceRowIndex:  &sourceIndex,
				SourceRowKey:    fmt.Sprintf("i7_prev_%s_%d", marketCode, groupID),
				BaseEntity:      rollbackScenarioBase(now),
			})
		}
	}
	if err := db.WithContext(ctx).Create(&orders).Error; err != nil {
		return fmt.Errorf("create rollback previous selected orders: %w", err)
	}
	return nil
}

func seedRollbackYearTwoOrders(
	ctx context.Context,
	db *gorm.DB,
	groupIDs []int64,
	adminOrderService *service.AdminOrderCommandService,
	adminOrderControlService *service.AdminOrderControlCommandService,
	playerOrderService *service.PlayerOrderCommandService,
) error {
	yearNo := 2
	if _, err := adminOrderService.UpdateMarketConfig(ctx, service.UpdateOrderMarketConfigCommand{
		YearNo: yearNo,
		Markets: []service.UpdateOrderMarketConfigItem{
			{MarketCode: enum.MarketCodeLocal, Enabled: true},
			{MarketCode: enum.MarketCodeRegional, Enabled: false},
			{MarketCode: enum.MarketCodeNational, Enabled: false},
			{MarketCode: enum.MarketCodeGlobal, Enabled: false},
		},
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("update rollback market config: %w", err)
	}
	if _, err := adminOrderService.UpdateForecastControl(ctx, service.UpdateOrderForecastControlCommand{
		Items:        rollbackForecastControlItems(yearNo),
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("update rollback forecast control: %w", err)
	}
	if _, err := adminOrderService.UpdateControlConfig(ctx, service.UpdateOrderControlConfigCommand{
		YearNo:       yearNo,
		Items:        rollbackControlConfigItems(yearNo),
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("update rollback order control config: %w", err)
	}
	generated, err := adminOrderService.GenerateOrderPool(ctx, service.GenerateOrderPoolCommand{
		YearNo:       yearNo,
		Overwrite:    true,
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	})
	if err != nil {
		return fmt.Errorf("generate rollback order pool: %w", err)
	}
	if _, err := adminOrderService.ConfirmOrderPool(ctx, service.ConfirmOrderPoolCommand{
		YearNo:       yearNo,
		BatchID:      generated.BatchID,
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("confirm rollback order pool: %w", err)
	}

	investments := [][]float64{
		{12, 3, 0, 0},
		{10, 2, 0, 0},
		{7, 1, 0, 0},
	}
	for index, groupID := range groupIDs {
		if _, err := playerOrderService.SubmitMarketInvestment(ctx, service.SubmitMarketInvestmentCommand{
			GroupID: groupID,
			YearNo:  yearNo,
			Investments: rollbackMarketInvestments(func(segment service.OrderSegmentDefinition) float64 {
				if segment.MarketCode != enum.MarketCodeLocal {
					return 0
				}
				switch segment.OrderType {
				case enum.OrderTypeAgencyInspection:
					return investments[index][0]
				case enum.OrderTypeTwoCabinVIP:
					return investments[index][1]
				default:
					return 0
				}
			}),
			OperatorName: fmt.Sprintf("group%02d", index+1),
		}); err != nil {
			return fmt.Errorf("submit rollback market investments group=%d: %w", groupID, err)
		}
	}
	if _, err := adminOrderControlService.GenerateSelectionSequence(ctx, service.GenerateSelectionSequenceCommand{
		YearNo:       yearNo,
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("generate rollback selection sequence: %w", err)
	}
	if err := rollbackCompleteSegment(ctx, db, adminOrderControlService, playerOrderService, yearNo, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection, true); err != nil {
		return err
	}
	if err := rollbackCompleteSegment(ctx, db, adminOrderControlService, playerOrderService, yearNo, enum.MarketCodeLocal, enum.OrderTypeTwoCabinVIP, false); err != nil {
		return err
	}
	return nil
}

func rollbackCompleteSegment(
	ctx context.Context,
	db *gorm.DB,
	adminOrderControlService *service.AdminOrderControlCommandService,
	playerOrderService *service.PlayerOrderCommandService,
	yearNo int,
	marketCode string,
	orderType string,
	selectFirst bool,
) error {
	released, err := adminOrderControlService.ReleaseNextSegment(ctx, service.ReleaseNextSegmentCommand{
		YearNo:       yearNo,
		OperatorID:   rollbackScenarioOperatorID,
		OperatorName: rollbackScenarioOperatorName,
	})
	if err != nil {
		return fmt.Errorf("release rollback segment %s/%s: %w", marketCode, orderType, err)
	}
	if released.MarketCode != marketCode || released.OrderType != orderType || released.CurrentGroupID == nil {
		return fmt.Errorf("unexpected rollback released segment: got %#v want %s/%s", released, marketCode, orderType)
	}
	currentGroupID := *released.CurrentGroupID
	if selectFirst {
		orders, err := repository.NewOrderPoolRepository(db).ListBySegment(ctx, yearNo, marketCode, orderType)
		if err != nil {
			return fmt.Errorf("list rollback segment orders %s/%s: %w", marketCode, orderType, err)
		}
		if len(orders) == 0 {
			return fmt.Errorf("rollback segment %s/%s has no orders", marketCode, orderType)
		}
		if _, err := playerOrderService.SelectOrder(ctx, service.SelectOrderCommand{
			GroupID:      currentGroupID,
			YearNo:       yearNo,
			MarketCode:   marketCode,
			OrderType:    orderType,
			OrderID:      orders[0].ID,
			OperatorID:   currentGroupID,
			OperatorName: fmt.Sprintf("group-%d", currentGroupID),
		}); err != nil {
			return fmt.Errorf("select rollback order %s/%s: %w", marketCode, orderType, err)
		}
	}
	for {
		stateItem, err := repository.NewMarketBiddingStateRepository(db).GetBySegment(ctx, yearNo, marketCode, orderType)
		if err != nil {
			return fmt.Errorf("reload rollback segment %s/%s: %w", marketCode, orderType, err)
		}
		if stateItem.SegmentStatus == enum.OrderSegmentStatusCompleted {
			return nil
		}
		if stateItem.SegmentStatus != enum.OrderSegmentStatusSelecting || stateItem.CurrentGroupID == nil {
			return fmt.Errorf("rollback segment %s/%s unexpected status after select/pass: %#v", marketCode, orderType, stateItem)
		}
		if _, err := playerOrderService.PassSegment(ctx, service.PassOrderSegmentCommand{
			GroupID:      *stateItem.CurrentGroupID,
			YearNo:       yearNo,
			MarketCode:   marketCode,
			OrderType:    orderType,
			OperatorID:   *stateItem.CurrentGroupID,
			OperatorName: fmt.Sprintf("group-%d", *stateItem.CurrentGroupID),
		}); err != nil {
			return fmt.Errorf("pass rollback segment %s/%s group=%d: %w", marketCode, orderType, *stateItem.CurrentGroupID, err)
		}
	}
}

func seedRollbackYearTwoOperations(
	ctx context.Context,
	db *gorm.DB,
	groupIDs []int64,
	playerOrderService *service.PlayerOrderCommandService,
	operatingService *service.PlayerOperatingCommandService,
	reportService *service.PlayerReportCommandService,
) error {
	for index, groupID := range groupIDs {
		operatorName := fmt.Sprintf("group%02d", index+1)
		baseCost := float64(10 + index*2)
		deliveryRevenue, deliveredOrderIDs, err := rollbackDeliveryRevenue(ctx, db, groupID, 2, index == 0)
		if err != nil {
			return err
		}
		operatingPayload := rollbackScenarioOperatingPayload(2, index, baseCost, deliveryRevenue, state.StageCodeQ1)
		if len(deliveredOrderIDs) > 0 {
			if _, err := operatingService.SaveDraft(ctx, service.SaveOperatingDraftCommand{
				GroupID:          groupID,
				YearNo:           2,
				StageStatus:      enum.StageStatusQ1Open,
				OperatingPayload: operatingPayload,
				OperatorName:     operatorName,
			}); err != nil {
				return fmt.Errorf("save rollback delivery draft group=%d: %w", groupID, err)
			}
			if _, err := playerOrderService.DeliverOrders(ctx, service.DeliverOrdersCommand{
				GroupID:      groupID,
				YearNo:       2,
				StageCode:    state.StageCodeQ1,
				OrderIDs:     deliveredOrderIDs,
				OperatorName: operatorName,
			}); err != nil {
				return fmt.Errorf("deliver rollback orders group=%d: %w", groupID, err)
			}
		}
		if err := submitRollbackFullOperatingYear(ctx, operatingService, groupID, 2, index, baseCost, deliveryRevenue); err != nil {
			return err
		}
		reportManual, err := rollbackScenarioReportManualPayloadForGroup(ctx, db, groupID, 2)
		if err != nil {
			return err
		}
		if _, err := reportService.Submit(ctx, service.SubmitPlayerReportCommand{
			GroupID:             groupID,
			YearNo:              2,
			ReportManualPayload: reportManual,
			SubmitterID:         groupID,
			OperatorName:        operatorName,
		}); err != nil {
			return fmt.Errorf("submit rollback year 2 report group=%d: %w", groupID, err)
		}
	}
	return nil
}

func rollbackDeliveryRevenue(ctx context.Context, db *gorm.DB, groupID int64, yearNo int, shouldDeliver bool) (float64, []int64, error) {
	if !shouldDeliver {
		return 0, nil, nil
	}
	selections, err := repository.NewGroupOrderSelectionRepository(db).ListByGroupYear(ctx, groupID, yearNo)
	if err != nil {
		return 0, nil, fmt.Errorf("list rollback selected orders group=%d year=%d: %w", groupID, yearNo, err)
	}
	orderIDs := make([]int64, 0, len(selections))
	for _, selection := range selections {
		if selection.DeliveryStatus == enum.OrderDeliveryStatusSelected &&
			selection.DeliveryEffective &&
			selection.DeliveredStageCode == nil {
			orderIDs = append(orderIDs, selection.OrderID)
		}
	}
	details, err := repository.NewGroupOrderSelectionRepository(db).ListSelectedOrderDetailsForUpdate(ctx, groupID, yearNo, orderIDs)
	if err != nil {
		return 0, nil, fmt.Errorf("load rollback selected order details group=%d year=%d: %w", groupID, yearNo, err)
	}
	if len(details) == 0 {
		return 0, nil, fmt.Errorf("rollback scenario expected group=%d year=%d to have a deliverable selected order", groupID, yearNo)
	}
	return details[0].OrderAmount, []int64{details[0].OrderID}, nil
}

func submitRollbackFullOperatingYear(ctx context.Context, operatingService *service.PlayerOperatingCommandService, groupID int64, yearNo int, groupIndex int, baseCost float64, deliveryRevenue float64) error {
	payloads := map[string]payload.OperatingPayload{
		state.StageCodeQ1:      rollbackScenarioOperatingPayload(yearNo, groupIndex, baseCost, deliveryRevenue, state.StageCodeQ1),
		state.StageCodeQ2:      rollbackScenarioOperatingPayload(yearNo, groupIndex, baseCost, deliveryRevenue, state.StageCodeQ2),
		state.StageCodeQ3:      rollbackScenarioOperatingPayload(yearNo, groupIndex, baseCost, deliveryRevenue, state.StageCodeQ3),
		state.StageCodeQ4:      rollbackScenarioOperatingPayload(yearNo, groupIndex, baseCost, deliveryRevenue, state.StageCodeQ4),
		state.StageCodeYearEnd: rollbackScenarioOperatingPayload(yearNo, groupIndex, baseCost, deliveryRevenue, state.StageCodeYearEnd),
	}
	for _, stageCode := range []string{state.StageCodeQ1, state.StageCodeQ2, state.StageCodeQ3, state.StageCodeQ4, state.StageCodeYearEnd} {
		if _, err := operatingService.SubmitStage(ctx, service.SubmitOperatingStageCommand{
			GroupID:          groupID,
			YearNo:           yearNo,
			StageCode:        stageCode,
			OperatingPayload: payloads[stageCode],
			SubmitterID:      groupID,
			OperatorName:     fmt.Sprintf("group-%d", groupID),
		}); err != nil {
			return fmt.Errorf("submit rollback operating group=%d year=%d stage=%s: %w", groupID, yearNo, stageCode, err)
		}
	}
	return nil
}

func rollbackScenarioReportManualPayloadForGroup(ctx context.Context, db *gorm.DB, groupID int64, yearNo int) (payload.ReportManualPayload, error) {
	group, err := repository.NewGroupRepository(db).GetByID(ctx, groupID)
	if err != nil {
		return payload.ReportManualPayload{}, fmt.Errorf("load rollback report group=%d: %w", groupID, err)
	}
	gameConfig, err := repository.NewGameConfigRepository(db).GetCurrent(ctx)
	if err != nil {
		return payload.ReportManualPayload{}, fmt.Errorf("load rollback report game config: %w", err)
	}
	yearState, err := repository.NewGroupYearStateRepository(db).GetByGroupIDAndYear(ctx, groupID, yearNo)
	if err != nil {
		return payload.ReportManualPayload{}, fmt.Errorf("load rollback report year state group=%d year=%d: %w", groupID, yearNo, err)
	}
	draft, err := repository.NewOperatingRepository(db).FindDraft(ctx, groupID, yearNo)
	if err != nil {
		return payload.ReportManualPayload{}, fmt.Errorf("load rollback report operating draft group=%d year=%d: %w", groupID, yearNo, err)
	}
	var operatingPayload payload.OperatingPayload
	if err := json.Unmarshal(draft.OperatingPayload, &operatingPayload); err != nil {
		return payload.ReportManualPayload{}, fmt.Errorf("unmarshal rollback report operating draft group=%d year=%d: %w", groupID, yearNo, err)
	}
	operatingPayload = operatingPayload.Normalize().WithoutDerivedValues()
	noticeService := service.NewPlayerNoticeService(
		repository.NewNoticeRepository(db),
		repository.NewGroupAdjustmentRepository(db),
	)
	operatingPayload, err = noticeService.OverlayAdjustments(ctx, groupID, yearNo, operatingPayload)
	if err != nil {
		return payload.ReportManualPayload{}, fmt.Errorf("overlay rollback report adjustments group=%d year=%d: %w", groupID, yearNo, err)
	}
	orderLinkService := service.NewOrderOperatingLinkService(
		repository.NewGroupMarketBidRepository(db),
		repository.NewMarketBiddingStateRepository(db),
		repository.NewGroupOrderSelectionRepository(db),
	)
	operatingPayload, _, err = orderLinkService.ApplyFormalYearValues(ctx, groupID, yearNo, operatingPayload)
	if err != nil {
		return payload.ReportManualPayload{}, fmt.Errorf("apply rollback report order values group=%d year=%d: %w", groupID, yearNo, err)
	}
	zeroInventoryManual := payload.ReportManualPayload{
		WorkInProgress:               float64Ptr(0),
		FinishedGoods:                float64Ptr(0),
		RawMaterials:                 float64Ptr(0),
		IncomeTaxRate:                float64Ptr(0),
		EnterpriseCertificationScore: float64Ptr(1),
		ProductionHumanScore:         float64Ptr(1),
		ClosingSpeedScore:            float64Ptr(1),
	}
	calcContext := calcctx.NewCalculationContext(*group, *yearState, *gameConfig)
	if yearNo == 0 {
		calcContext = calcContext.WithInitialBaseline(rollbackScenarioBaselinePayload())
	} else {
		previousReport, err := repository.NewReportRepository(db).FindEffectiveByGroupIDAndYear(ctx, groupID, yearNo-1)
		if err != nil {
			return payload.ReportManualPayload{}, fmt.Errorf("load rollback previous report group=%d year=%d: %w", groupID, yearNo-1, err)
		}
		var previousComputed payload.ReportComputedPayload
		if err := json.Unmarshal(previousReport.ReportComputedPayload, &previousComputed); err != nil {
			return payload.ReportManualPayload{}, fmt.Errorf("unmarshal rollback previous report group=%d year=%d: %w", groupID, yearNo-1, err)
		}
		calcContext = calcContext.WithPreviousReport(&previousComputed)
	}
	calcContext = calcContext.
		WithOperatingPayload(&operatingPayload).
		WithReportManualPayload(&zeroInventoryManual)
	if err := calcContext.Validate(); err != nil {
		return payload.ReportManualPayload{}, fmt.Errorf("validate rollback report context group=%d year=%d: %w", groupID, yearNo, err)
	}
	computed, err := reportrules.NewCalculator().Calculate(calcContext)
	if err != nil {
		return payload.ReportManualPayload{}, fmt.Errorf("calculate rollback zero-inventory report group=%d year=%d: %w", groupID, yearNo, err)
	}
	inventoryTotal := math.Round(-computed.BalanceGap())
	if inventoryTotal < 0 {
		return payload.ReportManualPayload{}, fmt.Errorf("rollback report balance requires negative inventory group=%d year=%d gap=%.2f", groupID, yearNo, computed.BalanceGap())
	}
	finishedGoods := float64(yearNo + group.GroupNo)
	rawMaterials := float64(group.GroupNo)
	workInProgress := inventoryTotal - finishedGoods - rawMaterials
	if workInProgress < 0 {
		workInProgress = inventoryTotal
		finishedGoods = 0
		rawMaterials = 0
	}
	result := payload.ReportManualPayload{
		WorkInProgress:               float64Ptr(workInProgress),
		FinishedGoods:                float64Ptr(finishedGoods),
		RawMaterials:                 float64Ptr(rawMaterials),
		IncomeTaxRate:                float64Ptr(0),
		EnterpriseCertificationScore: float64Ptr(float64(yearNo + 1)),
		ProductionHumanScore:         float64Ptr(float64(group.GroupNo + 1)),
		ClosingSpeedScore:            float64Ptr(float64(yearNo + group.GroupNo + 1)),
	}
	checkContext := calcContext.WithReportManualPayload(&result)
	checkComputed, err := reportrules.NewCalculator().Calculate(checkContext)
	if err != nil {
		return payload.ReportManualPayload{}, fmt.Errorf("check rollback report balance group=%d year=%d: %w", groupID, yearNo, err)
	}
	if math.Abs(checkComputed.BalanceGap()) > 0.000001 {
		return payload.ReportManualPayload{}, fmt.Errorf("rollback report still unbalanced group=%d year=%d gap=%.6f inventory=%.2f", groupID, yearNo, checkComputed.BalanceGap(), inventoryTotal)
	}
	return result, nil
}

func buildRollbackOperatingService(db *gorm.DB) *service.PlayerOperatingCommandService {
	noticeService := service.NewPlayerNoticeService(
		repository.NewNoticeRepository(db),
		repository.NewGroupAdjustmentRepository(db),
	)
	orderLinkService := service.NewOrderOperatingLinkService(
		repository.NewGroupMarketBidRepository(db),
		repository.NewMarketBiddingStateRepository(db),
		repository.NewGroupOrderSelectionRepository(db),
	)
	return service.NewPlayerOperatingCommandService(
		db,
		repository.NewGameConfigRepository(db),
		repository.NewGroupRepository(db),
		repository.NewGroupYearStateRepository(db),
		repository.NewOperatingRepository(db),
		repository.NewInitialBaselineRepository(db),
		repository.NewReportRepository(db),
		noticeService,
		orderLinkService,
	)
}

func buildRollbackReportService(db *gorm.DB) *service.PlayerReportCommandService {
	noticeService := service.NewPlayerNoticeService(
		repository.NewNoticeRepository(db),
		repository.NewGroupAdjustmentRepository(db),
	)
	orderLinkService := service.NewOrderOperatingLinkService(
		repository.NewGroupMarketBidRepository(db),
		repository.NewMarketBiddingStateRepository(db),
		repository.NewGroupOrderSelectionRepository(db),
	)
	return service.NewPlayerReportCommandService(
		db,
		repository.NewGameConfigRepository(db),
		repository.NewGroupRepository(db),
		repository.NewGroupYearStateRepository(db),
		repository.NewOperatingRepository(db),
		repository.NewInitialBaselineRepository(db),
		repository.NewReportRepository(db),
		noticeService,
		orderLinkService,
	)
}

func rollbackScenarioBaselinePayload() *payload.BaselinePayload {
	item := payload.BaselinePayload{
		BaselineSalesRevenue:         32,
		BaselineDirectCost:           15,
		BaselineComprehensiveCost:    13,
		BaselineDepreciation:         1,
		BaselineFinanceIncomeExpense: 2,
		BaselineExtraIncomeExpense:   2,
		BaselineIncomeTax:            1,
		BaselineWorkInConstruction:   0,
		BaselineFactoryAsset:         40,
		BaselineLineResidual:         3,
		BaselineDepreciableAsset:     0,
		BaselineCash:                 90,
		BaselineReceivable:           0,
		BaselineWorkInProgress:       6,
		BaselineFinishedGoods:        4,
		BaselineRawMaterials:         1,
		BaselineShortTermLoan:        20,
		BaselineLongTermLoan:         0,
		BaselineShareCapital:         50,
		BaselineRetainedEarnings:     17,
	}
	return &item
}

func rollbackScenarioOperatingPayload(yearNo int, groupIndex int, baseCost float64, deliveryRevenue float64, throughStageCode string) payload.OperatingPayload {
	p := payload.NewOperatingPayload()
	stageLimit := rollbackScenarioStageLimit(throughStageCode)
	beginningInvestment := float64((yearNo + 1) * (groupIndex + 1))
	p.Beginning.TaxAndPlanning = map[string]any{
		"taxRate":       0,
		"marketBidCost": beginningInvestment,
	}
	p.Beginning.MarketBid = []map[string]any{
		{"marketCode": enum.MarketCodeLocal, "marketInvestment": beginningInvestment, "orderAmount": deliveryRevenue},
	}
	for quarterIndex, quarter := range []string{"q1", "q2", "q3", "q4"} {
		if quarterIndex > stageLimit {
			continue
		}
		quarterOffset := float64(quarterIndex + 1)
		groupOffset := float64(groupIndex + 1)
		yearOffset := float64(yearNo + 1)
		cashSafetyLoan := float64(40 + yearNo*5 + groupIndex*10)
		balanceBuffer := 0.0
		cashBuffer := 0.0
		if quarter == "q4" {
			balanceBuffer = float64(20 + yearNo*3 + groupIndex*2)
			cashBuffer = float64(80 + yearNo*5 + groupIndex*10)
		}
		p.Quarter.ShortTermLoan[quarter] = map[string]any{
			"shortTermRepayment": 0,
			"shortTermInterest":  0,
			"newShortTermLoan":   groupOffset + quarterOffset + cashSafetyLoan + balanceBuffer + cashBuffer,
		}
		p.Quarter.MaterialPayment[quarter] = map[string]any{"materialPayment": baseCost + quarterOffset + balanceBuffer}
		p.Quarter.ProductionLineAdjust[quarter] = map[string]any{
			"changeProductCost":        yearOffset,
			"lineDismantleCost":        groupOffset,
			"lineSaleValue":            0,
			"newLineInstall":           quarterOffset,
			"constructionToFixed":      0,
			"depreciableAssetIncrease": quarterOffset,
		}
		p.Quarter.HumanResource[quarter] = map[string]any{"humanResourceCost": groupOffset}
		p.Quarter.SalaryAndProduction[quarter] = map[string]any{"salaryAndProductionCost": yearOffset}
		p.Quarter.ResearchAndManagement[quarter] = map[string]any{
			"researchCost":         quarterOffset,
			"managementSystemCost": groupOffset,
		}
		p.Quarter.ReceivableUpdate[quarter] = map[string]any{"receivableRecovered": 0}
		salesRevenue := 0.0
		directCost := 0.0
		if quarter == "q1" && deliveryRevenue > 0 {
			salesRevenue = deliveryRevenue
			directCost = math.Floor(deliveryRevenue / 2)
		}
		p.Quarter.DeliverySettlement[quarter] = map[string]any{
			"salesRevenue":      salesRevenue,
			"directCost":        directCost,
			"managementSalary":  groupOffset,
			"deliveryQuantity":  quarterOffset,
			"receivableBalance": 0,
		}
		p.Extra.IncomeAndPenalty[quarter] = map[string]any{
			"discountExpense":     0,
			"extraExpensePenalty": groupOffset,
			"extraIncomeReward":   8 + yearOffset + quarterOffset,
		}
	}
	if stageLimit >= 4 {
		groupOffset := float64(groupIndex + 1)
		yearOffset := float64(yearNo + 1)
		p.YearEnd.LongTermLoan = map[string]any{
			"longTermInterest":  groupOffset,
			"longTermRepayment": 0,
			"newLongTermLoan":   4 + yearOffset + groupOffset,
		}
		p.YearEnd.AssetAdjustment = map[string]any{
			"lineMaintenance":    groupOffset,
			"factoryPurchase":    0,
			"factorySale":        0,
			"factoryRent":        yearOffset,
			"workInConstruction": yearOffset + groupOffset,
			"marketCultivation":  groupOffset + 1,
		}
	}
	return p.Normalize()
}

func rollbackScenarioStageLimit(stageCode string) int {
	switch stageCode {
	case state.StageCodeQ1:
		return 0
	case state.StageCodeQ2:
		return 1
	case state.StageCodeQ3:
		return 2
	case state.StageCodeQ4:
		return 3
	case state.StageCodeYearEnd:
		return 4
	default:
		return 4
	}
}

func rollbackForecastControlItems(yearNo int) []service.UpdateOrderForecastControlItem {
	items := make([]service.UpdateOrderForecastControlItem, 0, 16)
	for _, segment := range rollbackOrderSegments() {
		count := 0
		if segment.MarketCode == enum.MarketCodeLocal {
			switch segment.OrderType {
			case enum.OrderTypeAgencyInspection:
				count = 2
			case enum.OrderTypeTwoCabinVIP:
				count = 1
			}
		}
		items = append(items, service.UpdateOrderForecastControlItem{
			YearNo:     yearNo,
			MarketCode: segment.MarketCode,
			OrderType:  segment.OrderType,
			OrderCount: count,
		})
	}
	return items
}

func rollbackControlConfigItems(yearNo int) []service.UpdateOrderControlConfigItem {
	items := make([]service.UpdateOrderControlConfigItem, 0, 16)
	for index, segment := range rollbackOrderSegments() {
		items = append(items, service.UpdateOrderControlConfigItem{
			MarketCode:        segment.MarketCode,
			OrderType:         segment.OrderType,
			OrderCount:        0,
			ReleaseSequenceNo: index + 1,
		})
	}
	return items
}

func rollbackMarketInvestments(resolve func(segment service.OrderSegmentDefinition) float64) []service.MarketInvestmentInput {
	items := make([]service.MarketInvestmentInput, 0, 16)
	for _, segment := range rollbackOrderSegments() {
		items = append(items, service.MarketInvestmentInput{
			MarketCode:       segment.MarketCode,
			OrderType:        segment.OrderType,
			MarketInvestment: resolve(segment),
		})
	}
	return items
}

func rollbackOrderSegments() []service.OrderSegmentDefinition {
	markets := []struct {
		code string
		name string
	}{
		{enum.MarketCodeLocal, "本地市场"},
		{enum.MarketCodeRegional, "区域市场"},
		{enum.MarketCodeNational, "全国市场"},
		{enum.MarketCodeGlobal, "全球市场"},
	}
	orderTypes := []struct {
		code string
		name string
	}{
		{enum.OrderTypeAgencyInspection, "代办过检"},
		{enum.OrderTypeTwoCabinVIP, "两舱贵宾"},
		{enum.OrderTypeBusinessVIP, "商务贵宾"},
		{enum.OrderTypeMemberCustom, "会员定制"},
	}
	items := make([]service.OrderSegmentDefinition, 0, 16)
	for _, market := range markets {
		for _, orderType := range orderTypes {
			items = append(items, service.OrderSegmentDefinition{
				MarketCode:    market.code,
				MarketName:    market.name,
				OrderType:     orderType.code,
				OrderTypeName: orderType.name,
			})
		}
	}
	return items
}

func rollbackScenarioBase(now time.Time) entity.BaseEntity {
	return entity.BaseEntity{
		Creator:    rollbackScenarioOperatorName,
		CreateTime: now,
		Updater:    rollbackScenarioOperatorName,
		UpdateTime: now,
	}
}

func rollbackMarshalJSON(value any) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal rollback scenario json: %w", err)
	}
	return raw, nil
}
