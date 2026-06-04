package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/service"
)

const (
	orderScenarioOperatorID   int64 = 1
	orderScenarioOperatorName       = "I4-07订单专项造数"
)

type seededOrderScenario struct {
	GroupIDs []int64
}

func seedOrderScenario(ctx context.Context, db *gorm.DB) error {
	var seeded seededOrderScenario
	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		seeded.GroupIDs, err = seedOrderScenarioGroups(ctx, tx)
		if err != nil {
			return err
		}
		if err := seedOrderScenarioHistoricalYears(ctx, tx, seeded.GroupIDs); err != nil {
			return err
		}
		if err := seedOrderScenarioPreviousOrders(ctx, tx, seeded.GroupIDs); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("seed order scenario transaction: %w", err)
	}

	fmt.Printf(
		"order scenario seeded: groups=%d groupIDs=%v stop=2年已开放、1年订单历史已写入、2年订单配置未开始\n",
		len(seeded.GroupIDs),
		seeded.GroupIDs,
	)
	return nil
}

func seedOrderScenarioGroups(ctx context.Context, tx *gorm.DB) ([]int64, error) {
	now := time.Now()
	groups := []entity.Group{
		buildScenarioGroup(1, "GROUP_01", "第一组", now),
		buildScenarioGroup(2, "GROUP_02", "第二组", now),
		buildScenarioGroup(3, "GROUP_03", "第三组", now),
	}
	if err := tx.WithContext(ctx).Create(&groups).Error; err != nil {
		return nil, fmt.Errorf("create scenario groups: %w", err)
	}

	accounts := make([]entity.Account, 0, len(groups))
	for index, group := range groups {
		groupID := group.ID
		accounts = append(accounts, entity.Account{
			ID:           int64(101 + index),
			Username:     fmt.Sprintf("group%02d", group.GroupNo),
			PasswordHash: "{sha256}8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
			RoleType:     enum.RoleTypeGroup,
			GroupID:      &groupID,
			Status:       enum.AccountStatusEnabled,
			BaseEntity:   buildScenarioBase(now),
		})
	}
	if err := tx.WithContext(ctx).Create(&accounts).Error; err != nil {
		return nil, fmt.Errorf("create scenario accounts: %w", err)
	}

	if err := tx.WithContext(ctx).
		Model(&entity.GameConfig{}).
		Where("id = ?", 1).
		Updates(map[string]any{
			"final_year":                 8,
			"current_open_year":          2,
			"edition_code":               service.GameEditionVIPServiceV1,
			"edition_name":               service.GameEditionVIPServiceV1Name,
			"rule_version":               service.FormulaVersionCommonV1,
			"template_version":           service.TemplateVersionVIPServiceV1,
			"operating_template_version": service.OperatingTemplateVersionVIPServiceV1,
			"report_template_version":    service.ReportTemplateVersionVIPServiceV1,
			"order_template_version":     service.OrderTemplateVersionVIPServiceV1,
			"process_rule_version":       service.ProcessRuleVersionCommonV1,
			"initial_baseline_submitted": true,
			"updater":                    orderScenarioOperatorName,
			"update_time":                now,
		}).Error; err != nil {
		return nil, fmt.Errorf("update scenario game config: %w", err)
	}

	groupIDs := make([]int64, 0, len(groups))
	for _, group := range groups {
		groupIDs = append(groupIDs, group.ID)
	}
	return groupIDs, nil
}

func buildScenarioGroup(no int, code string, name string, now time.Time) entity.Group {
	return entity.Group{
		ID:             int64(no),
		GroupNo:        no,
		GroupCode:      code,
		GroupName:      name,
		BusinessStatus: enum.BusinessStatusNormal,
		BaseEntity:     buildScenarioBase(now),
	}
}

func seedOrderScenarioHistoricalYears(ctx context.Context, tx *gorm.DB, groupIDs []int64) error {
	now := time.Now()
	manualJSON, computedJSON, err := buildScenarioReportJSON()
	if err != nil {
		return err
	}
	operatingJSON, err := marshalScenarioJSON(buildScenarioOperatingPayload())
	if err != nil {
		return err
	}
	baselineJSON, err := marshalScenarioJSON(buildScenarioBaselinePayload())
	if err != nil {
		return err
	}
	for _, groupID := range groupIDs {
		submittedAt := now
		submitterID := orderScenarioOperatorID
		baseline := entity.InitialBaseline{
			GroupID:         groupID,
			BaselinePayload: baselineJSON,
			Submitted:       true,
			SubmitterID:     &submitterID,
			SubmittedAt:     &submittedAt,
			BaseEntity:      buildScenarioBase(now),
		}
		if err := tx.WithContext(ctx).Create(&baseline).Error; err != nil {
			return fmt.Errorf("create baseline for group %d: %w", groupID, err)
		}

		for _, yearNo := range []int{0, 1} {
			yearType := enum.YearTypeFormal
			if yearNo == 0 {
				yearType = enum.YearTypeDemo
			}
			yearState := entity.GroupYearState{
				GroupID:                   groupID,
				YearNo:                    yearNo,
				YearType:                  yearType,
				YearStatus:                enum.YearStatusCompleted,
				StageStatus:               enum.StageStatusYearEndOpen,
				ReportStatus:              enum.ReportStatusSubmitted,
				SummaryEffective:          yearNo > 0,
				LatestStageSubmitVersion:  5,
				LatestReportSubmitVersion: 1,
				BaseEntity:                buildScenarioBase(now),
			}
			if err := tx.WithContext(ctx).Create(&yearState).Error; err != nil {
				return fmt.Errorf("create year state group=%d year=%d: %w", groupID, yearNo, err)
			}
			if err := seedScenarioOperatingDraftAndSubmissions(ctx, tx, groupID, yearNo, operatingJSON, now); err != nil {
				return err
			}
			if err := seedScenarioReport(ctx, tx, groupID, yearNo, manualJSON, computedJSON, now); err != nil {
				return err
			}
			if yearNo > 0 {
				if err := seedScenarioSummary(ctx, tx, groupID, yearNo, now); err != nil {
					return err
				}
			}
		}

		yearTwoState := entity.GroupYearState{
			GroupID:                   groupID,
			YearNo:                    2,
			YearType:                  enum.YearTypeFormal,
			YearStatus:                enum.YearStatusOperating,
			StageStatus:               enum.StageStatusQ1Open,
			ReportStatus:              enum.ReportStatusLocked,
			SummaryEffective:          false,
			LatestStageSubmitVersion:  0,
			LatestReportSubmitVersion: 0,
			BaseEntity:                buildScenarioBase(now),
		}
		if err := tx.WithContext(ctx).Create(&yearTwoState).Error; err != nil {
			return fmt.Errorf("create year 2 state group=%d: %w", groupID, err)
		}
	}
	return nil
}

func seedScenarioOperatingDraftAndSubmissions(ctx context.Context, tx *gorm.DB, groupID int64, yearNo int, operatingJSON []byte, now time.Time) error {
	draft := entity.GroupOperatingDraft{
		GroupID:          groupID,
		YearNo:           yearNo,
		StageStatus:      enum.StageStatusYearEndOpen,
		OperatingPayload: append([]byte(nil), operatingJSON...),
		LastAutoSavedAt:  &now,
		BaseEntity:       buildScenarioBase(now),
	}
	if err := tx.WithContext(ctx).Create(&draft).Error; err != nil {
		return fmt.Errorf("create operating draft group=%d year=%d: %w", groupID, yearNo, err)
	}
	stages := []string{"Q1", "Q2", "Q3", "Q4", "YEAR_END"}
	for index, stageCode := range stages {
		submission := entity.GroupStageSubmission{
			GroupID:                  groupID,
			YearNo:                   yearNo,
			StageCode:                stageCode,
			SubmitVersion:            1,
			PeriodEndCash:            80,
			OperatingPayloadSnapshot: append([]byte(nil), operatingJSON...),
			StateBefore:              []byte("{}"),
			StateAfter:               []byte("{}"),
			SubmitterID:              groupID,
			SubmitTime:               now.Add(time.Duration(index) * time.Minute),
		}
		if err := tx.WithContext(ctx).Create(&submission).Error; err != nil {
			return fmt.Errorf("create operating submission group=%d year=%d stage=%s: %w", groupID, yearNo, stageCode, err)
		}
	}
	return nil
}

func seedScenarioReport(ctx context.Context, tx *gorm.DB, groupID int64, yearNo int, manualJSON []byte, computedJSON []byte, now time.Time) error {
	submittedAt := now
	report := entity.GroupReport{
		GroupID:               groupID,
		YearNo:                yearNo,
		ReportManualPayload:   append([]byte(nil), manualJSON...),
		ReportComputedPayload: append([]byte(nil), computedJSON...),
		BalanceCheckPassed:    true,
		LastAutoSavedAt:       &now,
		SubmittedAt:           &submittedAt,
		BaseEntity:            buildScenarioBase(now),
	}
	if err := tx.WithContext(ctx).Create(&report).Error; err != nil {
		return fmt.Errorf("create report group=%d year=%d: %w", groupID, yearNo, err)
	}
	submission := entity.GroupReportSubmission{
		GroupID:                groupID,
		YearNo:                 yearNo,
		SubmitVersion:          1,
		ReportManualSnapshot:   append([]byte(nil), manualJSON...),
		ReportComputedSnapshot: append([]byte(nil), computedJSON...),
		BalanceCheckPassed:     true,
		StateBefore:            []byte("{}"),
		StateAfter:             []byte("{}"),
		SubmitterID:            groupID,
		SubmitTime:             now,
	}
	if err := tx.WithContext(ctx).Create(&submission).Error; err != nil {
		return fmt.Errorf("create report submission group=%d year=%d: %w", groupID, yearNo, err)
	}
	return nil
}

func seedScenarioSummary(ctx context.Context, tx *gorm.DB, groupID int64, yearNo int, now time.Time) error {
	summary := entity.GroupSummarySnapshot{
		GroupID:                   groupID,
		YearNo:                    yearNo,
		Revenue:                   0,
		Profit:                    0,
		Equity:                    69,
		BusinessStatus:            enum.BusinessStatusNormal,
		RankingValue:              69,
		SummaryEffective:          true,
		SourceReportSubmitVersion: 1,
		BaseEntity:                buildScenarioBase(now),
	}
	if err := tx.WithContext(ctx).Create(&summary).Error; err != nil {
		return fmt.Errorf("create summary group=%d year=%d: %w", groupID, yearNo, err)
	}
	return nil
}

func seedOrderScenarioPreviousOrders(ctx context.Context, tx *gorm.DB, groupIDs []int64) error {
	if len(groupIDs) < 3 {
		return fmt.Errorf("seed previous orders requires 3 groups")
	}
	now := time.Now()
	type leaderCase struct {
		marketCode string
		amounts    []float64
	}
	cases := []leaderCase{
		{marketCode: enum.MarketCodeLocal, amounts: []float64{120, 80, 60}},
		{marketCode: enum.MarketCodeRegional, amounts: []float64{70, 130, 90}},
		{marketCode: enum.MarketCodeNational, amounts: []float64{65, 75, 150}},
		{marketCode: enum.MarketCodeGlobal, amounts: []float64{140, 100, 110}},
	}
	orders := make([]entity.OrderPool, 0, len(cases)*len(groupIDs))
	for _, item := range cases {
		for index, groupID := range groupIDs {
			selectedGroupID := groupID
			selectedAt := now
			sourceIndex := index + 1
			amount := item.amounts[index]
			orders = append(orders, entity.OrderPool{
				YearNo:          1,
				MarketCode:      item.marketCode,
				OrderType:       enum.OrderTypeAgencyInspection,
				SegmentCode:     fmt.Sprintf("%s_%s", item.marketCode, enum.OrderTypeAgencyInspection),
				CardSequenceNo:  index + 1,
				BusinessOrderNo: fmt.Sprintf("Y1-%s-G%d", item.marketCode, index+1),
				OrderAmount:     amount,
				OrderQuantity:   1,
				UnitPrice:       amount,
				AccountTerm:     2,
				PoolStatus:      enum.OrderPoolStatusSelected,
				SelectedGroupID: &selectedGroupID,
				SelectedAt:      &selectedAt,
				SourceSheetName: "I4-07造数",
				SourceCell:      "PREV",
				SourceRowIndex:  &sourceIndex,
				SourceRowKey:    fmt.Sprintf("i4_07_prev_%s_%d", item.marketCode, groupID),
				BaseEntity:      buildScenarioBase(now),
			})
		}
	}
	if err := tx.WithContext(ctx).Create(&orders).Error; err != nil {
		return fmt.Errorf("create previous selected orders: %w", err)
	}
	return nil
}

func buildScenarioReportJSON() ([]byte, []byte, error) {
	manual := payload.ReportManualPayload{
		WorkInProgress:               float64Ptr(6),
		FinishedGoods:                float64Ptr(4),
		RawMaterials:                 float64Ptr(1),
		IncomeTaxRate:                float64Ptr(0),
		EnterpriseCertificationScore: float64Ptr(0),
		ProductionHumanScore:         float64Ptr(0),
		ClosingSpeedScore:            float64Ptr(0),
	}
	computed := payload.ReportComputedPayload{
		ReportWorkInProgress:        6,
		ReportFinishedGoods:         4,
		ReportRawMaterials:          1,
		ReportFactoryAsset:          40,
		ReportLineResidual:          3,
		ReportTotalNonCurrentAssets: 43,
		ReportCash:                  35,
		ReportPostTaxCash:           35,
		ReportTotalCurrentAssets:    46,
		ReportTotalAssets:           89,
		ReportShortTermLiability:    20,
		ReportTotalLiability:        20,
		ReportShareCapital:          50,
		ReportRetainedEarnings:      19,
		ReportTotalEquity:           69,
		ReportTotalLiabilityEquity:  89,
	}
	manualJSON, err := marshalScenarioJSON(manual)
	if err != nil {
		return nil, nil, err
	}
	computedJSON, err := marshalScenarioJSON(computed)
	if err != nil {
		return nil, nil, err
	}
	return manualJSON, computedJSON, nil
}

func buildScenarioOperatingPayload() payload.OperatingPayload {
	p := payload.NewOperatingPayload()
	p.Beginning.TaxAndPlanning = map[string]any{
		"taxRate":       0,
		"marketBidCost": 0,
	}
	p.Beginning.MarketBid = []map[string]any{
		{"market": "LOCAL", "marketInvestment": 0, "orderAmount": 0},
	}
	for _, quarter := range []string{"q1", "q2", "q3", "q4"} {
		p.Quarter.ShortTermLoan[quarter] = map[string]any{
			"shortTermRepayment": 0,
			"shortTermInterest":  0,
			"newShortTermLoan":   0,
		}
		p.Quarter.MaterialPayment[quarter] = map[string]any{"materialPayment": 0}
		p.Quarter.ProductionLineAdjust[quarter] = map[string]any{
			"changeProductCost":        0,
			"lineDismantleCost":        0,
			"lineSaleValue":            0,
			"newLineInstall":           0,
			"constructionToFixed":      0,
			"depreciableAssetIncrease": 0,
		}
		p.Quarter.HumanResource[quarter] = map[string]any{"humanResourceCost": 0}
		p.Quarter.SalaryAndProduction[quarter] = map[string]any{"salaryAndProductionCost": 0}
		p.Quarter.ResearchAndManagement[quarter] = map[string]any{
			"researchCost":         0,
			"managementSystemCost": 0,
		}
		p.Quarter.ReceivableUpdate[quarter] = map[string]any{"receivableRecovered": 0}
		p.Quarter.DeliverySettlement[quarter] = map[string]any{
			"salesRevenue":      0,
			"directCost":        0,
			"managementSalary":  0,
			"deliveryQuantity":  0,
			"receivableBalance": 0,
		}
		p.Extra.IncomeAndPenalty[quarter] = map[string]any{
			"discountExpense":     0,
			"extraExpensePenalty": 0,
			"extraIncomeReward":   0,
		}
	}
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

func buildScenarioBaselinePayload() payload.BaselinePayload {
	return payload.BaselinePayload{
		BaselineSalesRevenue:         32,
		BaselineDirectCost:           15,
		BaselineComprehensiveCost:    13,
		BaselineDepreciation:         1,
		BaselineFinanceIncomeExpense: 2,
		BaselineExtraIncomeExpense:   2,
		BaselineIncomeTax:            1,
		BaselineFactoryAsset:         40,
		BaselineLineResidual:         3,
		BaselineCash:                 36,
		BaselineWorkInProgress:       6,
		BaselineFinishedGoods:        4,
		BaselineRawMaterials:         1,
		BaselineShortTermLoan:        20,
		BaselineShareCapital:         50,
		BaselineRetainedEarnings:     17,
	}
}

func buildScenarioBase(now time.Time) entity.BaseEntity {
	return entity.BaseEntity{
		Creator:    orderScenarioOperatorName,
		CreateTime: now,
		Updater:    orderScenarioOperatorName,
		UpdateTime: now,
	}
}

func marshalScenarioJSON(value any) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal scenario json: %w", err)
	}
	return raw, nil
}

func float64Ptr(value float64) *float64 {
	return &value
}
