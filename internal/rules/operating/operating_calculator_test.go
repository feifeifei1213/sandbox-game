package operating

import (
	"math"
	"testing"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
)

func TestCalculateDemoYearBuildsQuarterCashChecksAndYearEndCash(t *testing.T) {
	t.Parallel()

	calculator := NewCalculator()
	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Beginning.MarketBid = []map[string]any{
		{"marketInvestment": 7.0},
	}
	operatingPayload.Quarter.ShortTermLoan = payload.OperatingQuarterMap{
		"q1": {"dueRepayment": 5.0, "interest": 2.0, "newLoan": 11.0},
		"q2": {"dueRepayment": 3.0, "interest": 1.0, "newLoan": 0.0},
		"q3": {"dueRepayment": 0.0, "interest": 2.0, "newLoan": 5.0},
		"q4": {"dueRepayment": 4.0, "interest": 1.0, "newLoan": 0.0},
	}
	operatingPayload.Quarter.MaterialPayment = payload.OperatingQuarterMap{
		"q1": {"itemA": 3.0, "itemB": 4.0},
		"q2": {"materialCost": 5.0},
		"q3": {"materialCost": 6.0},
		"q4": {"materialCost": 7.0},
	}
	operatingPayload.Quarter.ProductionLineAdjust = payload.OperatingQuarterMap{
		"q1": {"changeProduct": 1.0, "dismantleCost": 2.0, "lineSale": 3.0, "newLineInstall": 4.0, "constructionToFixed": 5.0, "newDepreciableAsset": 6.0},
		"q2": {"changeProduct": 0.0, "dismantleCost": 1.0, "lineSale": 0.0, "newLineInstall": 2.0, "constructionToFixed": 0.0, "newDepreciableAsset": 3.0},
		"q3": {"changeProduct": 2.0, "dismantleCost": 0.0, "lineSale": 1.0, "newLineInstall": 0.0, "constructionToFixed": 4.0, "newDepreciableAsset": 0.0},
		"q4": {"changeProduct": 1.0, "dismantleCost": 1.0, "lineSale": 2.0, "newLineInstall": 1.0, "constructionToFixed": 0.0, "newDepreciableAsset": 2.0},
	}
	operatingPayload.Quarter.HumanResource = payload.OperatingQuarterMap{
		"q1": {"staffCost": 2.0},
		"q2": {"staffCost": 1.0},
		"q3": {"staffCost": 3.0},
		"q4": {"staffCost": 4.0},
	}
	operatingPayload.Quarter.SalaryAndProduction = payload.OperatingQuarterMap{
		"q1": {"salaryCost": 8.0},
		"q2": {"salaryCost": 7.0},
		"q3": {"salaryCost": 6.0},
		"q4": {"salaryCost": 5.0},
	}
	operatingPayload.Quarter.ResearchAndManagement = payload.OperatingQuarterMap{
		"q1": {"technologyResearch": 2.0, "managementSystem": 1.0},
		"q2": {"technologyResearch": 1.0, "managementSystem": 2.0},
		"q3": {"technologyResearch": 3.0, "managementSystem": 1.0},
		"q4": {"technologyResearch": 2.0, "managementSystem": 2.0},
	}
	operatingPayload.Quarter.ReceivableUpdate = payload.OperatingQuarterMap{
		"q1": {"receivableCollection": 9.0},
		"q2": {"receivableCollection": 10.0},
		"q3": {"receivableCollection": 11.0},
		"q4": {"receivableCollection": 12.0},
	}
	operatingPayload.Quarter.DeliverySettlement = payload.OperatingQuarterMap{
		"q1": {"salesRevenue": 13.0, "directCost": 4.0, "managementStaffCost": 1.0},
		"q2": {"salesRevenue": 14.0, "directCost": 5.0, "managementStaffCost": 2.0},
		"q3": {"salesRevenue": 15.0, "directCost": 6.0, "managementStaffCost": 1.0},
		"q4": {"salesRevenue": 16.0, "directCost": 7.0, "managementStaffCost": 2.0},
	}
	operatingPayload.YearEnd.LongTermLoan = map[string]any{
		"interest":  4.0,
		"repayment": 8.0,
		"newLoan":   12.0,
	}
	operatingPayload.YearEnd.AssetAdjustment = map[string]any{
		"lineMaintenance":    2.0,
		"purchase":           5.0,
		"sale":               1.0,
		"rent":               3.0,
		"workInConstruction": 7.0,
		"marketCultivation":  9.0,
	}
	operatingPayload.Extra.IncomeAndPenalty = payload.OperatingQuarterMap{
		"q1": {"discountExpense": 1.0, "extraExpensePenalty": 2.0, "extraIncomeReward": 0.0},
		"q2": {"discountExpense": 0.0, "extraExpensePenalty": 0.0, "extraIncomeReward": 1.0},
		"q3": {"discountExpense": 2.0, "extraExpensePenalty": 1.0, "extraIncomeReward": 0.0},
		"q4": {"discountExpense": 1.0, "extraExpensePenalty": 0.0, "extraIncomeReward": 2.0},
	}

	ctx := newOperatingCalculationContext(0, enum.StageStatusYearEndOpen).
		WithInitialBaseline(&payload.BaselinePayload{
			BaselineIncomeTax:        10,
			BaselineShortTermLoan:    20,
			BaselineLongTermLoan:     30,
			BaselineLineResidual:     40,
			BaselineDepreciableAsset: 90,
			BaselineCash:             100,
		}).
		WithOperatingPayload(&operatingPayload)

	result, err := calculator.Calculate(ctx)
	if err != nil {
		t.Fatalf("calculate operating failed: %v", err)
	}

	assertFloatEquals(t, result.QuarterCashChecks["Q1"], 68)
	assertFloatEquals(t, result.QuarterCashChecks["Q2"], 54)
	assertFloatEquals(t, result.QuarterCashChecks["Q3"], 44)
	assertFloatEquals(t, result.QuarterCashChecks["Q4"], 29)
	assertFloatEquals(t, result.PeriodEndCash, 11)
	assertFloatEquals(t, result.DerivedValues["marketBidCost"], 7)
	assertFloatEquals(t, result.DerivedValues["comprehensiveCostTotal"], 59)
	assertFloatEquals(t, result.DerivedValues["shortTermLoanBalance"], 24)
	assertFloatEquals(t, result.DerivedValues["longTermLoanBalance"], 34)
	assertFloatEquals(t, result.DerivedValues["lineResidual"], 43)
	assertFloatEquals(t, result.DerivedValues["depreciableAssetTotal"], 101)
	assertFloatEquals(t, result.DerivedValues["depreciation"], 34)
	assertFloatEquals(t, result.DerivedValues["financeIncomeExpense"], 14)
	assertFloatEquals(t, result.DerivedValues["receivableChange"], 16)
	assertFloatEquals(t, result.DerivedValues["periodEndCash"], 11)
	assertFloatEquals(t, mustFloat64(t, result.OperatingPayload.Derived.Values["q4QuarterEndCashCheck"]), 29)
}

func TestCalculateFormalYearUsesPreviousReportCarryForwardAndCurrentStageCash(t *testing.T) {
	t.Parallel()

	calculator := NewCalculator()
	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Beginning.MarketBid = []map[string]any{
		{"marketInvestment": 4.0},
	}
	operatingPayload.Quarter.ShortTermLoan = payload.OperatingQuarterMap{
		"Q1": {"dueRepayment": 2.0, "interest": 1.0, "newLoan": 3.0},
		"Q2": {"dueRepayment": 1.0, "interest": 1.0, "newLoan": 0.0},
	}
	operatingPayload.Quarter.MaterialPayment = payload.OperatingQuarterMap{
		"Q1": {"materialCost": 1.0},
		"Q2": {"materialCost": 2.0},
	}
	operatingPayload.Quarter.ProductionLineAdjust = payload.OperatingQuarterMap{
		"Q1": {"changeProduct": 1.0, "lineSale": 2.0, "newLineInstall": 1.0, "constructionToFixed": 4.0, "newDepreciableAsset": 6.0},
		"Q2": {"dismantleCost": 1.0, "newDepreciableAsset": 3.0},
	}
	operatingPayload.Quarter.HumanResource = payload.OperatingQuarterMap{
		"Q1": {"staffCost": 1.0},
		"Q2": {"staffCost": 2.0},
	}
	operatingPayload.Quarter.SalaryAndProduction = payload.OperatingQuarterMap{
		"Q1": {"salaryCost": 2.0},
		"Q2": {"salaryCost": 3.0},
	}
	operatingPayload.Quarter.ResearchAndManagement = payload.OperatingQuarterMap{
		"Q1": {"technologyResearch": 1.0, "managementSystem": 1.0},
		"Q2": {"technologyResearch": 0.0, "managementSystem": 1.0},
	}
	operatingPayload.Quarter.ReceivableUpdate = payload.OperatingQuarterMap{
		"Q1": {"receivableCollection": 5.0},
		"Q2": {"receivableCollection": 6.0},
	}
	operatingPayload.Quarter.DeliverySettlement = payload.OperatingQuarterMap{
		"Q1": {"salesRevenue": 10.0, "directCost": 2.0, "managementStaffCost": 1.0},
		"Q2": {"salesRevenue": 12.0, "directCost": 3.0, "managementStaffCost": 1.0},
	}
	operatingPayload.YearEnd.LongTermLoan = map[string]any{
		"interest":  2.0,
		"repayment": 4.0,
		"newLoan":   5.0,
	}
	operatingPayload.YearEnd.AssetAdjustment = map[string]any{
		"lineMaintenance":    1.0,
		"purchase":           2.0,
		"sale":               0.0,
		"rent":               1.0,
		"workInConstruction": 3.0,
		"marketCultivation":  2.0,
	}
	operatingPayload.Extra.IncomeAndPenalty = payload.OperatingQuarterMap{
		"Q1": {"discountExpense": 1.0, "extraExpensePenalty": 0.0, "extraIncomeReward": 0.0},
		"Q2": {"discountExpense": 0.0, "extraExpensePenalty": 1.0, "extraIncomeReward": 2.0},
	}

	ctx := newOperatingCalculationContext(1, enum.StageStatusQ2Open).
		WithPreviousReport(&payload.ReportComputedPayload{
			ReportIncomeTax:          5,
			ReportShortTermLiability: 7,
			ReportLongTermLiability:  20,
			ReportLineResidual:       30,
			ReportDepreciableAsset:   33,
			ReportCash:               60,
		}).
		WithOperatingPayload(&operatingPayload)

	result, err := calculator.Calculate(ctx)
	if err != nil {
		t.Fatalf("calculate operating failed: %v", err)
	}

	assertFloatEquals(t, result.QuarterCashChecks["Q1"], 48)
	assertFloatEquals(t, result.QuarterCashChecks["Q2"], 43)
	assertFloatEquals(t, result.PeriodEndCash, 43)
	assertFloatEquals(t, result.DerivedValues["periodEndCash"], 36)
	assertFloatEquals(t, result.DerivedValues["shortTermLoanBalance"], 7)
	assertFloatEquals(t, result.DerivedValues["longTermLoanBalance"], 21)
	assertFloatEquals(t, result.DerivedValues["lineResidual"], 32)
	assertFloatEquals(t, result.DerivedValues["depreciableAssetTotal"], 42)
	assertFloatEquals(t, result.DerivedValues["depreciation"], 14)
	assertFloatEquals(t, result.DerivedValues["comprehensiveCostTotal"], 18)
}

func newOperatingCalculationContext(yearNo int, stageStatus string) calcctx.CalculationContext {
	group := entity.Group{
		ID:             1,
		BusinessStatus: enum.BusinessStatusNormal,
	}
	yearType := enum.YearTypeFormal
	if yearNo == 0 {
		yearType = enum.YearTypeDemo
	}
	yearState := entity.GroupYearState{
		GroupID:      1,
		YearNo:       yearNo,
		YearType:     yearType,
		YearStatus:   enum.YearStatusOperating,
		StageStatus:  stageStatus,
		ReportStatus: enum.ReportStatusLocked,
	}
	gameConfig := entity.GameConfig{
		FinalYear:       8,
		CurrentOpenYear: yearNo,
	}

	return calcctx.NewCalculationContext(group, yearState, gameConfig)
}

func assertFloatEquals(t *testing.T, actual float64, expected float64) {
	t.Helper()

	if math.Abs(actual-expected) > 1e-9 {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func mustFloat64(t *testing.T, value any) float64 {
	t.Helper()

	number, ok := value.(float64)
	if !ok {
		t.Fatalf("expected float64, got %T", value)
	}
	return number
}

func TestCalculateDoesNotDoubleCountAliasedMarketInvestment(t *testing.T) {
	t.Parallel()

	calculator := NewCalculator()
	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Beginning.TaxAndPlanning = map[string]any{
		"marketInvestmentTotal": 1.0,
		"marketBidCost":         1.0,
	}

	ctx := newOperatingCalculationContext(0, enum.StageStatusQ1Open).
		WithInitialBaseline(&payload.BaselinePayload{
			BaselineIncomeTax:        1,
			BaselineShortTermLoan:    0,
			BaselineLongTermLoan:     0,
			BaselineLineResidual:     0,
			BaselineDepreciableAsset: 0,
			BaselineCash:             36,
		}).
		WithOperatingPayload(&operatingPayload)

	result, err := calculator.Calculate(ctx)
	if err != nil {
		t.Fatalf("calculate operating failed: %v", err)
	}

	assertFloatEquals(t, result.DerivedValues["marketBidCost"], 1)
	assertFloatEquals(t, result.QuarterCashChecks["Q1"], 34)
	assertFloatEquals(t, result.PeriodEndCash, 34)
}
