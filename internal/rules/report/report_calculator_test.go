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
	enterpriseCertificationScore := 4.0
	productionHumanScore := 5.0
	closingSpeedScore := 6.0

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
			WorkInProgress:               &workInProgress,
			FinishedGoods:                &finishedGoods,
			RawMaterials:                 &rawMaterials,
			IncomeTaxRate:                &incomeTaxRate,
			EnterpriseCertificationScore: &enterpriseCertificationScore,
			ProductionHumanScore:         &productionHumanScore,
			ClosingSpeedScore:            &closingSpeedScore,
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
	assertFloatEquals(t, result.ReportBestMarketDirectorBaseScore, 107)
	assertFloatEquals(t, result.ReportBestMarketDirectorScore, 111)
	assertFloatEquals(t, result.ReportBestTechnologyDirectorScore, 47)
	assertFloatEquals(t, result.ReportBestSalesDirectorScore, 0)
	assertFloatEquals(t, result.ReportBestCfoBaseScore, 0)
	assertFloatEquals(t, result.ReportBestCfoScore, 6)
	assertFloatEquals(t, result.ReportBestCeoScore, 169)
	assertFloatEquals(t, result.BalanceGap(), 0)
}

func TestCalculateUsesPreviousReportForFormalYear(t *testing.T) {
	t.Parallel()

	calculator := NewCalculator()
	workInProgress := 20.0
	finishedGoods := 30.0
	rawMaterials := 40.0
	incomeTaxRate := 0.15
	enterpriseCertificationScore := 8.0
	productionHumanScore := 9.0
	closingSpeedScore := 10.0

	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Derived.Values = map[string]any{
		"orderTotal":               140.0,
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
			ReportIncomeTax:                   10,
			ReportShortTermLiability:          20,
			ReportLongTermLiability:           30,
			ReportFactoryAsset:                40,
			ReportLineResidual:                50,
			ReportDepreciableAsset:            75,
			ReportCash:                        60,
			ReportReceivable:                  70,
			ReportShareCapital:                284,
			ReportRetainedEarnings:            90,
			ReportNetProfit:                   10,
			ReportTotalEquity:                 69,
			ReportBestMarketDirectorBaseScore: 3,
			ReportBestMarketDirectorScore:     3,
			ReportBestTechnologyDirectorScore: 7,
			ReportBestSalesDirectorScore:      100,
			ReportBestCfoBaseScore:            0,
		}).
		WithOperatingPayload(&operatingPayload).
		WithReportManualPayload(&payload.ReportManualPayload{
			WorkInProgress:               &workInProgress,
			FinishedGoods:                &finishedGoods,
			RawMaterials:                 &rawMaterials,
			IncomeTaxRate:                &incomeTaxRate,
			EnterpriseCertificationScore: &enterpriseCertificationScore,
			ProductionHumanScore:         &productionHumanScore,
			ClosingSpeedScore:            &closingSpeedScore,
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
	assertFloatEquals(t, result.ReportBestMarketDirectorBaseScore, 23)
	assertFloatEquals(t, result.ReportBestMarketDirectorScore, 31)
	assertFloatEquals(t, result.ReportBestTechnologyDirectorScore, 16)
	assertFloatEquals(t, result.ReportBestSalesDirectorScore, 240)
	assertFloatEquals(t, result.ReportBestCfoBaseScore, 201.5)
	assertFloatEquals(t, result.ReportBestCfoScore, 211.5)
	assertFloatEquals(t, result.ReportBestCeoScore, 507.5)
	assertFloatEquals(t, result.BalanceGap(), 0)
}

func TestCalculateFormalYearAccumulatesCfoBaseAgainstDemoYearReference(t *testing.T) {
	t.Parallel()

	calculator := NewCalculator()
	workInProgress := 20.0
	finishedGoods := 30.0
	rawMaterials := 40.0
	incomeTaxRate := 0.15
	enterpriseCertificationScore := 8.0
	productionHumanScore := 9.0
	closingSpeedScore := 10.0

	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Derived.Values = map[string]any{
		"orderTotal":               140.0,
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

	ctx := newReportCalculationContext(2).
		WithPreviousReport(&payload.ReportComputedPayload{
			ReportIncomeTax:                   16,
			ReportShortTermLiability:          17,
			ReportLongTermLiability:           31,
			ReportFactoryAsset:                57,
			ReportLineResidual:                56,
			ReportDepreciableAsset:            60,
			ReportCash:                        -105,
			ReportReceivable:                  359,
			ReportShareCapital:                284,
			ReportRetainedEarnings:            100,
			ReportNetProfit:                   88,
			ReportTotalEquity:                 472,
			ReportBestMarketDirectorBaseScore: 23,
			ReportBestMarketDirectorScore:     31,
			ReportBestTechnologyDirectorScore: 16,
			ReportBestSalesDirectorScore:      240,
			ReportBestCfoBaseScore:            201.5,
			ReportBestCfoScore:                211.5,
		}).
		WithOperatingPayload(&operatingPayload).
		WithReportManualPayload(&payload.ReportManualPayload{
			WorkInProgress:               &workInProgress,
			FinishedGoods:                &finishedGoods,
			RawMaterials:                 &rawMaterials,
			IncomeTaxRate:                &incomeTaxRate,
			EnterpriseCertificationScore: &enterpriseCertificationScore,
			ProductionHumanScore:         &productionHumanScore,
			ClosingSpeedScore:            &closingSpeedScore,
		})

	result, err := calculator.Calculate(ctx)
	if err != nil {
		t.Fatalf("calculate report failed: %v", err)
	}

	assertFloatEquals(t, result.ReportTotalEquity, 565)
	assertFloatEquals(t, result.ReportBestCfoBaseScore, 248)
	assertFloatEquals(t, result.ReportBestCfoScore, 258)
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

func TestCalculateDemoYearMatchesExcelFinalSample(t *testing.T) {
	t.Parallel()

	calculator := NewCalculator()
	operatingPayload := buildExcelFinalSampleOperatingPayload()
	workInProgress := 8.0
	finishedGoods := 6.0
	rawMaterials := 0.0
	incomeTaxRate := 0.25
	enterpriseCertificationScore := 4.0
	productionHumanScore := 6.0
	closingSpeedScore := 2.0

	ctx := newReportCalculationContext(0).
		WithInitialBaseline(buildExcelFinalSampleBaseline()).
		WithOperatingPayload(&operatingPayload).
		WithReportManualPayload(&payload.ReportManualPayload{
			WorkInProgress:               &workInProgress,
			FinishedGoods:                &finishedGoods,
			RawMaterials:                 &rawMaterials,
			IncomeTaxRate:                &incomeTaxRate,
			EnterpriseCertificationScore: &enterpriseCertificationScore,
			ProductionHumanScore:         &productionHumanScore,
			ClosingSpeedScore:            &closingSpeedScore,
		})

	result, err := calculator.Calculate(ctx)
	if err != nil {
		t.Fatalf("calculate report failed: %v", err)
	}

	assertFloatEquals(t, result.ReportSalesRevenue, 11)
	assertFloatEquals(t, result.ReportDirectCost, 4)
	assertFloatEquals(t, result.ReportGrossProfit, 7)
	assertFloatEquals(t, result.ReportComprehensiveCost, 22)
	assertFloatEquals(t, result.ReportDepreciation, 6)
	assertFloatEquals(t, result.ReportOperatingProfit, -21)
	assertFloatEquals(t, result.ReportFinanceIncomeExpense, 2)
	assertFloatEquals(t, result.ReportExtraIncomeExpense, 0)
	assertFloatEquals(t, result.ReportPreTaxProfit, -23)
	assertFloatEquals(t, result.ReportIncomeTax, 0)
	assertFloatEquals(t, result.ReportNetProfit, -23)
	assertFloatEquals(t, result.ReportWorkInConstruction, 0)
	assertFloatEquals(t, result.ReportFactoryAsset, 45)
	assertFloatEquals(t, result.ReportLineResidual, 8)
	assertFloatEquals(t, result.ReportDepreciableAsset, 13)
	assertFloatEquals(t, result.ReportTotalNonCurrentAssets, 66)
	assertFloatEquals(t, result.ReportCash, 146)
	assertFloatEquals(t, result.ReportReceivable, 0)
	assertFloatEquals(t, result.ReportPostTaxCash, 146)
	assertFloatEquals(t, result.ReportTotalCurrentAssets, 160)
	assertFloatEquals(t, result.ReportTotalAssets, 226)
	assertFloatEquals(t, result.ReportShortTermLiability, 40)
	assertFloatEquals(t, result.ReportLongTermLiability, 140)
	assertFloatEquals(t, result.ReportTotalLiability, 180)
	assertFloatEquals(t, result.ReportShareCapital, 50)
	assertFloatEquals(t, result.ReportRetainedEarnings, 19)
	assertFloatEquals(t, result.ReportTotalEquity, 46)
	assertFloatEquals(t, result.ReportTotalLiabilityEquity, 226)
	assertFloatEquals(t, result.ReportBestMarketDirectorBaseScore, 1)
	assertFloatEquals(t, result.ReportBestMarketDirectorScore, 5)
	assertFloatEquals(t, result.ReportBestTechnologyDirectorScore, 5)
	assertFloatEquals(t, result.ReportBestSalesDirectorScore, 32)
	assertFloatEquals(t, result.ReportBestCfoBaseScore, 0)
	assertFloatEquals(t, result.ReportBestCfoScore, 2)
	assertFloatEquals(t, result.ReportBestCeoScore, 50)
	assertFloatEquals(t, result.BalanceGap(), 0)
}

func buildExcelFinalSampleBaseline() *payload.BaselinePayload {
	return &payload.BaselinePayload{
		BaselineSalesRevenue:         32,
		BaselineDirectCost:           15,
		BaselineComprehensiveCost:    13,
		BaselineDepreciation:         1,
		BaselineFinanceIncomeExpense: 2,
		BaselineExtraIncomeExpense:   2,
		BaselineIncomeTax:            1,
		BaselineFactoryAsset:         40,
		BaselineLineResidual:         3,
		BaselineDepreciableAsset:     0,
		BaselineCash:                 36,
		BaselineReceivable:           0,
		BaselineWorkInProgress:       6,
		BaselineFinishedGoods:        4,
		BaselineRawMaterials:         1,
		BaselineShortTermLoan:        20,
		BaselineLongTermLoan:         0,
		BaselineShareCapital:         50,
		BaselineRetainedEarnings:     17,
	}
}

func buildExcelFinalSampleOperatingPayload() payload.OperatingPayload {
	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Beginning.TaxAndPlanning = map[string]any{
		"marketInvestmentTotal": 1.0,
		"orderTotal":            32.0,
	}
	operatingPayload.Quarter.ShortTermLoan = payload.OperatingQuarterMap{
		"Q1": {"interest": 1.0, "newLoan": 20.0},
	}
	operatingPayload.Quarter.MaterialPayment = payload.OperatingQuarterMap{
		"Q1": {"materialCost": 20.0},
	}
	operatingPayload.Quarter.ProductionLineAdjust = payload.OperatingQuarterMap{
		"Q1": {
			"changeProduct":       2.0,
			"dismantleCost":       1.0,
			"lineSale":            5.0,
			"newLineInstall":      2.0,
			"constructionToFixed": 10.0,
			"newDepreciableAsset": 19.0,
		},
	}
	operatingPayload.Quarter.HumanResource = payload.OperatingQuarterMap{
		"Q1": {"staffCost": 4.0},
	}
	operatingPayload.Quarter.SalaryAndProduction = payload.OperatingQuarterMap{
		"Q1": {"salaryCost": 14.0},
	}
	operatingPayload.Quarter.ResearchAndManagement = payload.OperatingQuarterMap{
		"Q1": {"technologyResearch": 5.0, "managementSystem": 6.0},
	}
	operatingPayload.Quarter.ReceivableUpdate = payload.OperatingQuarterMap{
		"Q1": {"receivableCollection": 11.0},
	}
	operatingPayload.Quarter.DeliverySettlement = payload.OperatingQuarterMap{
		"Q1": {"salesRevenue": 11.0, "directCost": 4.0},
	}
	operatingPayload.YearEnd.LongTermLoan = map[string]any{
		"interest": 1.0,
		"newLoan":  140.0,
	}
	operatingPayload.YearEnd.AssetAdjustment = map[string]any{
		"lineMaintenance":   1.0,
		"purchase":          20.0,
		"sale":              15.0,
		"rent":              1.0,
		"marketCultivation": 1.0,
	}
	return operatingPayload
}
