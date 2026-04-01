package report

import (
	"math"
	"testing"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
)

func TestCalculateUsesInitialBaselineForDemoYear(t *testing.T) {
	t.Parallel()

	calculator := NewCalculator()
	workInProgress := 10.0
	finishedGoods := 20.0
	rawMaterials := 30.0
	incomeTaxRate := 0.25

	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Derived.Values = map[string]any{
		"marketBidCost":            11.0,
		"shortTermRepayment":       7.0,
		"shortTermInterest":        2.0,
		"newShortTermLoan":         13.0,
		"materialPayment":          17.0,
		"changeProductCost":        19.0,
		"lineDismantleCost":        23.0,
		"lineSaleValue":            29.0,
		"newLineInstall":           31.0,
		"transferToFixed":          37.0,
		"depreciableAssetIncrease": 45.0,
		"humanResourceCost":        41.0,
		"salaryAndProductionCost":  43.0,
		"researchCost":             47.0,
		"managementSystemCost":     53.0,
		"receivableRecovered":      59.0,
		"salesRevenue":             61.0,
		"directCost":               17.0,
		"managementSalary":         67.0,
		"longTermInterest":         71.0,
		"longTermRepayment":        73.0,
		"newLongTermLoan":          79.0,
		"lineMaintenance":          83.0,
		"factoryPurchase":          89.0,
		"factorySale":              97.0,
		"factoryRent":              101.0,
		"workInConstruction":       103.0,
		"marketCultivation":        107.0,
		"discountExpense":          109.0,
		"extraExpensePenalty":      113.0,
		"extraIncomeReward":        127.0,
	}

	ctx := newReportCalculationContext(0).
		WithInitialBaseline(&payload.BaselinePayload{
			BaselineSalesRevenue:         100,
			BaselineDirectCost:           40,
			BaselineComprehensiveCost:    10,
			BaselineDepreciation:         20,
			BaselineFinanceIncomeExpense: 5,
			BaselineExtraIncomeExpense:   3,
			BaselineIncomeTax:            7,
			BaselineFactoryAsset:         30,
			BaselineLineResidual:         40,
			BaselineDepreciableAsset:     90,
			BaselineCash:                 100,
			BaselineReceivable:           50,
			BaselineShortTermLoan:        10,
			BaselineLongTermLoan:         20,
			BaselineShareCapital:         363,
			BaselineRetainedEarnings:     60,
		}).
		WithOperatingPayload(&operatingPayload).
		WithReportManualPayload(&payload.ReportManualPayload{
			WorkInProgress: &workInProgress,
			FinishedGoods:  &finishedGoods,
			RawMaterials:   &rawMaterials,
			IncomeTaxRate:  &incomeTaxRate,
		})

	result, err := calculator.Calculate(ctx)
	if err != nil {
		t.Fatalf("calculate report failed: %v", err)
	}

	assertFloatEquals(t, result.ReportComprehensiveCost, 552)
	assertFloatEquals(t, result.ReportDepreciation, 45)
	assertFloatEquals(t, result.ReportGrossProfit, 44)
	assertFloatEquals(t, result.ReportOperatingProfit, -553)
	assertFloatEquals(t, result.ReportFinanceIncomeExpense, 182)
	assertFloatEquals(t, result.ReportExtraIncomeExpense, 14)
	assertFloatEquals(t, result.ReportPreTaxProfit, -721)
	assertFloatEquals(t, result.ReportIncomeTax, 0)
	assertFloatEquals(t, result.ReportNetProfit, -721)
	assertFloatEquals(t, result.ReportFactoryAsset, 22)
	assertFloatEquals(t, result.ReportLineResidual, 48)
	assertFloatEquals(t, result.ReportDepreciableAsset, 90)
	assertFloatEquals(t, result.ReportTotalNonCurrentAssets, 263)
	assertFloatEquals(t, result.ReportCash, -610)
	assertFloatEquals(t, result.ReportReceivable, 52)
	assertFloatEquals(t, result.ReportPostTaxCash, -610)
	assertFloatEquals(t, result.ReportTotalCurrentAssets, -498)
	assertFloatEquals(t, result.ReportTotalAssets, -235)
	assertFloatEquals(t, result.ReportShortTermLiability, 16)
	assertFloatEquals(t, result.ReportLongTermLiability, 26)
	assertFloatEquals(t, result.ReportTotalLiability, 42)
	assertFloatEquals(t, result.ReportShareCapital, 363)
	assertFloatEquals(t, result.ReportRetainedEarnings, 81)
	assertFloatEquals(t, result.ReportTotalEquity, -277)
	assertFloatEquals(t, result.ReportTotalLiabilityEquity, -235)
	assertFloatEquals(t, result.BalanceGap(), 0)
}

func TestCalculateUsesPreviousReportForFormalYear(t *testing.T) {
	t.Parallel()

	calculator := NewCalculator()
	workInProgress := 20.0
	finishedGoods := 30.0
	rawMaterials := 40.0
	incomeTaxRate := 0.15

	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Derived.Values = map[string]any{
		"marketBidCost":            10.0,
		"shortTermRepayment":       5.0,
		"shortTermInterest":        1.0,
		"newShortTermLoan":         2.0,
		"materialPayment":          3.0,
		"changeProductCost":        4.0,
		"lineDismantleCost":        0.0,
		"lineSaleValue":            0.0,
		"newLineInstall":           0.0,
		"transferToFixed":          6.0,
		"depreciableAssetIncrease": 15.0,
		"humanResourceCost":        7.0,
		"salaryAndProductionCost":  8.0,
		"researchCost":             9.0,
		"managementSystemCost":     10.0,
		"receivableRecovered":      11.0,
		"salesRevenue":             300.0,
		"directCost":               30.0,
		"managementSalary":         12.0,
		"longTermInterest":         13.0,
		"longTermRepayment":        14.0,
		"newLongTermLoan":          15.0,
		"lineMaintenance":          16.0,
		"factoryPurchase":          17.0,
		"factorySale":              0.0,
		"factoryRent":              18.0,
		"workInConstruction":       19.0,
		"marketCultivation":        20.0,
		"discountExpense":          21.0,
		"extraExpensePenalty":      4.0,
		"extraIncomeReward":        9.0,
	}

	ctx := newReportCalculationContext(1).
		WithPreviousReport(&payload.ReportComputedPayload{
			ReportIncomeTax:          10,
			ReportShortTermLiability: 20,
			ReportLongTermLiability:  30,
			ReportFactoryAsset:       40,
			ReportLineResidual:       50,
			ReportDepreciableAsset:   75,
			ReportCash:               60,
			ReportReceivable:         70,
			ReportShareCapital:       284,
			ReportRetainedEarnings:   90,
			ReportNetProfit:          10,
		}).
		WithOperatingPayload(&operatingPayload).
		WithReportManualPayload(&payload.ReportManualPayload{
			WorkInProgress: &workInProgress,
			FinishedGoods:  &finishedGoods,
			RawMaterials:   &rawMaterials,
			IncomeTaxRate:  &incomeTaxRate,
		})

	result, err := calculator.Calculate(ctx)
	if err != nil {
		t.Fatalf("calculate report failed: %v", err)
	}

	assertFloatEquals(t, result.ReportComprehensiveCost, 106)
	assertFloatEquals(t, result.ReportDepreciation, 30)
	assertFloatEquals(t, result.ReportPreTaxProfit, 104)
	assertFloatEquals(t, result.ReportIncomeTax, 16)
	assertFloatEquals(t, result.ReportNetProfit, 88)
	assertFloatEquals(t, result.ReportFactoryAsset, 57)
	assertFloatEquals(t, result.ReportLineResidual, 56)
	assertFloatEquals(t, result.ReportDepreciableAsset, 60)
	assertFloatEquals(t, result.ReportCash, -105)
	assertFloatEquals(t, result.ReportReceivable, 359)
	assertFloatEquals(t, result.ReportShortTermLiability, 17)
	assertFloatEquals(t, result.ReportLongTermLiability, 31)
	assertFloatEquals(t, result.ReportShareCapital, 284)
	assertFloatEquals(t, result.ReportRetainedEarnings, 100)
	assertFloatEquals(t, result.ReportTotalAssets, 520)
	assertFloatEquals(t, result.ReportTotalEquity, 472)
	assertFloatEquals(t, result.ReportTotalLiabilityEquity, 520)
	assertFloatEquals(t, result.BalanceGap(), 0)
}

func newReportCalculationContext(yearNo int) calcctx.CalculationContext {
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
		YearStatus:   enum.YearStatusReportPending,
		StageStatus:  enum.StageStatusYearEndOpen,
		ReportStatus: enum.ReportStatusOpen,
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

func TestCalculateDoesNotDoubleCountAliasedMarketInvestment(t *testing.T) {
	t.Parallel()

	calculator := NewCalculator()
	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Beginning.TaxAndPlanning = map[string]any{
		"marketInvestmentTotal": 1.0,
		"marketBidCost":         1.0,
	}

	ctx := newReportCalculationContext(0).
		WithInitialBaseline(&payload.BaselinePayload{
			BaselineIncomeTax:        1,
			BaselineShortTermLoan:    0,
			BaselineLongTermLoan:     0,
			BaselineFactoryAsset:     0,
			BaselineLineResidual:     0,
			BaselineDepreciableAsset: 0,
			BaselineCash:             36,
			BaselineReceivable:       0,
			BaselineShareCapital:     0,
			BaselineRetainedEarnings: 0,
		}).
		WithOperatingPayload(&operatingPayload)

	result, err := calculator.Calculate(ctx)
	if err != nil {
		t.Fatalf("calculate report failed: %v", err)
	}

	assertFloatEquals(t, result.ReportComprehensiveCost, 1)
	assertFloatEquals(t, result.ReportCash, 34)
}
