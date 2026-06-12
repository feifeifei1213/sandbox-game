package main

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/service"
)

const (
	airportOrderScenarioOperatorID   int64 = 1
	airportOrderScenarioOperatorName       = "I9机场订单专项造数"
)

func seedAirportOrderScenario(ctx context.Context, db *gorm.DB, withLeaderHistory bool) error {
	controlService := service.NewAdminControlCommandService(db)
	if _, err := controlService.InitializeGame(ctx, service.InitializeGameCommand{
		GroupCount:   3,
		EditionCode:  service.GameEditionAirportV1,
		OperatorID:   airportOrderScenarioOperatorID,
		OperatorName: airportOrderScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("initialize airport order scenario game: %w", err)
	}
	if _, err := controlService.SubmitInitialBaseline(ctx, service.SubmitInitialBaselineCommand{
		BaselinePayload: airportScenarioBaselinePayload(),
		OperatorID:      airportOrderScenarioOperatorID,
		OperatorName:    airportOrderScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("submit airport scenario initial baseline: %w", err)
	}

	groups, err := repository.NewGroupRepository(db).ListAll(ctx)
	if err != nil {
		return fmt.Errorf("list airport scenario groups: %w", err)
	}
	groupIDs := make([]int64, 0, len(groups))
	for _, group := range groups {
		groupIDs = append(groupIDs, group.ID)
	}

	if err := markAirportScenarioYearCompleted(ctx, db, groupIDs, 0); err != nil {
		return err
	}
	if _, err := controlService.OpenNextYear(ctx, service.OpenNextYearCommand{
		TargetYearNo: 1,
		OperatorID:   airportOrderScenarioOperatorID,
		OperatorName: airportOrderScenarioOperatorName,
	}); err != nil {
		return fmt.Errorf("open airport scenario year 1: %w", err)
	}

	stopText := "1年已开放，停在1年订单配置前"
	if withLeaderHistory {
		if err := seedAirportScenarioLeaderHistory(ctx, db, groupIDs); err != nil {
			return err
		}
		if err := markAirportScenarioYearCompleted(ctx, db, groupIDs, 1); err != nil {
			return err
		}
		if _, err := controlService.OpenNextYear(ctx, service.OpenNextYearCommand{
			TargetYearNo: 2,
			OperatorID:   airportOrderScenarioOperatorID,
			OperatorName: airportOrderScenarioOperatorName,
		}); err != nil {
			return fmt.Errorf("open airport scenario year 2: %w", err)
		}
		stopText = "2年已开放，1年国内/国际已选订单历史已写入，停在2年订单配置前"
	}

	fmt.Printf(
		"airport order scenario seeded: groups=%d groupIDs=%v edition=%s stop=%s\n",
		len(groupIDs),
		groupIDs,
		service.GameEditionAirportV1,
		stopText,
	)
	return nil
}

func markAirportScenarioYearCompleted(ctx context.Context, db *gorm.DB, groupIDs []int64, yearNo int) error {
	now := time.Now()
	yearType := enum.YearTypeFormal
	if yearNo == 0 {
		yearType = enum.YearTypeDemo
	}
	for _, groupID := range groupIDs {
		if err := db.WithContext(ctx).
			Model(&entity.GroupYearState{}).
			Where("group_id = ? AND year_no = ?", groupID, yearNo).
			Updates(map[string]any{
				"year_type":                    yearType,
				"year_status":                  enum.YearStatusCompleted,
				"stage_status":                 enum.StageStatusYearEndOpen,
				"report_status":                enum.ReportStatusSubmitted,
				"summary_effective":            yearNo > 0,
				"latest_stage_submit_version":  5,
				"latest_report_submit_version": 1,
				"updater":                      airportOrderScenarioOperatorName,
				"update_time":                  now,
			}).Error; err != nil {
			return fmt.Errorf("complete airport scenario year group=%d year=%d: %w", groupID, yearNo, err)
		}
	}
	return nil
}

func seedAirportScenarioLeaderHistory(ctx context.Context, db *gorm.DB, groupIDs []int64) error {
	if len(groupIDs) < 3 {
		return fmt.Errorf("airport leader scenario expects at least 3 groups")
	}
	now := time.Now()
	type leaderCase struct {
		marketCode string
		orderType  string
		amounts    []float64
	}
	cases := []leaderCase{
		{marketCode: service.MarketCodeDomestic, orderType: service.OrderTypeNarrowBody, amounts: []float64{90, 130, 70}},
		{marketCode: service.MarketCodeDomestic, orderType: service.OrderTypeWideBody, amounts: []float64{60, 120, 55}},
		{marketCode: service.MarketCodeInternational, orderType: service.OrderTypeNarrowBody, amounts: []float64{80, 65, 145}},
		{marketCode: service.MarketCodeInternational, orderType: service.OrderTypeWideBody, amounts: []float64{75, 60, 155}},
	}
	orders := make([]entity.OrderPool, 0, len(cases)*len(groupIDs))
	for _, item := range cases {
		for index, groupID := range groupIDs {
			selectedGroupID := groupID
			selectedAt := now
			sourceIndex := index + 1
			amount := item.amounts[index]
			orders = append(orders, entity.OrderPool{
				OrderTemplateVersion: service.OrderTemplateVersionAirportV1,
				YearNo:               1,
				MarketCode:           item.marketCode,
				OrderType:            item.orderType,
				SegmentCode:          fmt.Sprintf("%s_%s", item.marketCode, item.orderType),
				CardSequenceNo:       index + 1,
				BusinessOrderNo:      fmt.Sprintf("AIR-Y1-%s-%s-G%d", item.marketCode, item.orderType, index+1),
				OrderAmount:          amount,
				OrderQuantity:        10000,
				UnitPrice:            0.01,
				AccountTerm:          2,
				PoolStatus:           enum.OrderPoolStatusSelected,
				SelectedGroupID:      &selectedGroupID,
				SelectedAt:           &selectedAt,
				SourceSheetName:      "I9机场订单专项造数",
				SourceCell:           "PREV",
				SourceRowIndex:       &sourceIndex,
				SourceRowKey:         fmt.Sprintf("i9_airport_prev_%s_%s_%d", item.marketCode, item.orderType, groupID),
				BaseEntity:           airportScenarioBase(now),
			})
		}
	}
	if err := db.WithContext(ctx).Create(&orders).Error; err != nil {
		return fmt.Errorf("create airport previous selected orders: %w", err)
	}
	return nil
}

func airportScenarioBaselinePayload() *payload.BaselinePayload {
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

func airportScenarioBase(now time.Time) entity.BaseEntity {
	return entity.BaseEntity{
		Creator:    airportOrderScenarioOperatorName,
		CreateTime: now,
		Updater:    airportOrderScenarioOperatorName,
		UpdateTime: now,
	}
}
