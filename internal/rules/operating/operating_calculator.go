package operating

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
)

var quarterSequence = []string{"Q1", "Q2", "Q3", "Q4"}

type operatingCarryBase struct {
	previousIncomeTax        float64
	previousShortTermLoan    float64
	previousLongTermLoan     float64
	previousLineResidual     float64
	previousDepreciableAsset float64
	previousCash             float64
}

type annualOperatingMetrics struct {
	marketBidCost        float64
	orderTotal           float64
	shortTermRepayment   float64
	shortTermInterest    float64
	newShortTermLoan     float64
	materialPayment      float64
	changeProductCost    float64
	lineDismantleCost    float64
	lineSaleValue        float64
	newLineInstall       float64
	transferToFixed      float64
	depreciableAssetIncr float64
	humanResourceCost    float64
	salaryAndProduction  float64
	researchCost         float64
	managementSystemCost float64
	receivableRecovered  float64
	salesRevenue         float64
	directCost           float64
	managementSalary     float64
	longTermInterest     float64
	longTermRepayment    float64
	newLongTermLoan      float64
	lineMaintenance      float64
	factoryPurchase      float64
	factorySale          float64
	factoryRent          float64
	workInConstruction   float64
	marketCultivation    float64
	discountExpense      float64
	extraExpensePenalty  float64
	extraIncomeReward    float64
}

type quarterOperatingMetrics struct {
	shortTermRepayment   float64
	shortTermInterest    float64
	newShortTermLoan     float64
	materialPayment      float64
	changeProductCost    float64
	lineDismantleCost    float64
	lineSaleValue        float64
	newLineInstall       float64
	humanResourceCost    float64
	salaryAndProduction  float64
	researchCost         float64
	managementSystemCost float64
	receivableRecovered  float64
	managementSalary     float64
	discountExpense      float64
	extraExpensePenalty  float64
	extraIncomeReward    float64
}

// CalculationResult 统一承接经营页规则层输出。
type CalculationResult struct {
	OperatingPayload  payload.OperatingPayload `json:"operatingPayload"`
	QuarterCashChecks map[string]float64       `json:"quarterCashChecks"`
	DerivedValues     map[string]float64       `json:"derivedValues"`
	PeriodEndCash     float64                  `json:"periodEndCash"`
}

// Calculator 是经营规则的统一入口。
type Calculator struct{}

// NewCalculator 创建经营计算器。
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Calculate 按当前 Excel 经营页规则链返回派生值、季度现金与提交节点现金。
func (c *Calculator) Calculate(ctx calcctx.CalculationContext) (CalculationResult, error) {
	if err := ctx.Validate(); err != nil {
		return CalculationResult{}, err
	}
	if ctx.OperatingPayload == nil {
		return CalculationResult{}, fmt.Errorf("missing operating payload for operating calculation")
	}

	carryBase, err := buildCarryBase(ctx)
	if err != nil {
		return CalculationResult{}, err
	}

	normalizedPayload := ctx.OperatingPayload.Normalize()
	annualMetrics := extractAnnualMetrics(normalizedPayload)
	quarterMetrics := extractQuarterMetrics(normalizedPayload)
	quarterCashChecks := calculateQuarterCashChecks(carryBase, annualMetrics.marketBidCost, quarterMetrics)
	derivedValues := buildDerivedValues(carryBase, annualMetrics, quarterCashChecks)

	normalizedPayload.Derived.Values = mergeDerivedValues(normalizedPayload.Derived.Values, derivedValues)

	return CalculationResult{
		OperatingPayload:  normalizedPayload,
		QuarterCashChecks: quarterCashChecks,
		DerivedValues:     derivedValues,
		PeriodEndCash:     selectPeriodEndCash(ctx.State.StageStatus, quarterCashChecks, derivedValues["periodEndCash"]),
	}, nil
}

func buildCarryBase(ctx calcctx.CalculationContext) (operatingCarryBase, error) {
	if ctx.IsDemoYear() {
		if ctx.InitialBaseline == nil {
			return operatingCarryBase{}, fmt.Errorf("missing initial baseline for demo year operating")
		}

		return operatingCarryBase{
			previousIncomeTax:        ctx.InitialBaseline.BaselineIncomeTax,
			previousShortTermLoan:    ctx.InitialBaseline.BaselineShortTermLoan,
			previousLongTermLoan:     ctx.InitialBaseline.BaselineLongTermLoan,
			previousLineResidual:     ctx.InitialBaseline.BaselineLineResidual,
			previousDepreciableAsset: ctx.InitialBaseline.BaselineDepreciableAsset,
			previousCash:             ctx.InitialBaseline.BaselineCash,
		}, nil
	}

	if ctx.PreviousReport == nil {
		return operatingCarryBase{}, fmt.Errorf("missing previous report for formal year operating")
	}

	return operatingCarryBase{
		previousIncomeTax:        ctx.PreviousReport.ReportIncomeTax,
		previousShortTermLoan:    ctx.PreviousReport.ReportShortTermLiability,
		previousLongTermLoan:     ctx.PreviousReport.ReportLongTermLiability,
		previousLineResidual:     ctx.PreviousReport.ReportLineResidual,
		previousDepreciableAsset: ctx.PreviousReport.ReportDepreciableAsset,
		previousCash:             ctx.PreviousReport.ReportCash,
	}, nil
}

func extractAnnualMetrics(value payload.OperatingPayload) annualOperatingMetrics {
	derived := value.Derived.Values

	return annualOperatingMetrics{
		marketBidCost: firstNonZero(
			lookupPrioritizedNumber(derived, "marketBidCost", "marketInvestmentTotal", "p5"),
			lookupPrioritizedNumber(value.Beginning.TaxAndPlanning, "marketBidCost", "marketInvestmentTotal", "p5"),
			extractMetric(value.Beginning.MarketBid, false, "marketBidCost", "marketInvestment", "marketInvestmentTotal", "investment", "bidInvestment", "p5"),
		),
		orderTotal: firstNonZero(
			lookupPrioritizedNumber(derived, "orderTotal", "orderAmountTotal", "o5"),
			lookupPrioritizedNumber(value.Beginning.TaxAndPlanning, "orderTotal", "orderAmountTotal", "o5"),
			extractMetric(value.Beginning.MarketBid, false, "orderTotal", "orderAmount", "totalOrderAmount", "annualOrderTotal", "o5"),
		),
		shortTermRepayment: firstNonZero(
			lookupAnyNumber(derived, "shortTermRepayment", "o10"),
			sumQuarterMetric(value.Quarter.ShortTermLoan, false, "shortTermRepayment", "dueRepayment", "repayment", "o10"),
		),
		shortTermInterest: firstNonZero(
			lookupAnyNumber(derived, "shortTermInterest", "financeShortTermInterest", "o11"),
			sumQuarterMetric(value.Quarter.ShortTermLoan, false, "shortTermInterest", "interest", "loanInterest", "o11"),
		),
		newShortTermLoan: firstNonZero(
			lookupAnyNumber(derived, "newShortTermLoan", "shortTermNewLoan", "o12"),
			sumQuarterMetric(value.Quarter.ShortTermLoan, false, "newShortTermLoan", "newLoan", "additionalLoan", "o12"),
		),
		materialPayment: firstNonZero(
			lookupAnyNumber(derived, "materialPayment", "materialPaymentTotal", "o18"),
			sumQuarterMetric(value.Quarter.MaterialPayment, true, "materialPayment", "materialPaymentTotal", "materialCost", "o18"),
		),
		changeProductCost: firstNonZero(
			lookupAnyNumber(derived, "changeProductCost", "o21"),
			sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "changeProductCost", "changeProduct", "productChange", "switchProduct", "o21"),
		),
		lineDismantleCost: firstNonZero(
			lookupAnyNumber(derived, "lineDismantleCost", "o22"),
			sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "lineDismantleCost", "dismantleCost", "removeCost", "o22"),
		),
		lineSaleValue: firstNonZero(
			lookupAnyNumber(derived, "lineSaleValue", "o23"),
			sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "lineSaleValue", "lineSale", "sellValue", "saleValue", "o23"),
		),
		newLineInstall: firstNonZero(
			lookupAnyNumber(derived, "newLineInstall", "o26"),
			sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "newLineInstall", "installCost", "newLineInvestment", "o26"),
		),
		transferToFixed: firstNonZero(
			lookupAnyNumber(derived, "transferToFixed", "lineResidualIncrease", "o27"),
			sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "transferToFixed", "constructionToFixed", "completedTransfer", "lineResidualIncrease", "o27"),
		),
		depreciableAssetIncr: firstNonZero(
			lookupAnyNumber(derived, "depreciableAssetIncrease", "o28"),
			sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "depreciableAssetIncrease", "assetCapitalization", "newDepreciableAsset", "o28"),
		),
		humanResourceCost: firstNonZero(
			lookupAnyNumber(derived, "humanResourceCost", "o29"),
			sumQuarterMetric(value.Quarter.HumanResource, true, "humanResourceCost", "humanResourceFee", "staffCost", "o29"),
		),
		salaryAndProduction: firstNonZero(
			lookupAnyNumber(derived, "salaryAndProductionCost", "o31"),
			sumQuarterMetric(value.Quarter.SalaryAndProduction, true, "salaryAndProductionCost", "salaryCost", "workerSalary", "productionSalary", "o31"),
		),
		researchCost: firstNonZero(
			lookupAnyNumber(derived, "researchCost", "o33"),
			sumQuarterMetric(value.Quarter.ResearchAndManagement, false, "researchCost", "technologyResearch", "rdCost", "o33"),
		),
		managementSystemCost: firstNonZero(
			lookupAnyNumber(derived, "managementSystemCost", "o34"),
			sumQuarterMetric(value.Quarter.ResearchAndManagement, false, "managementSystemCost", "managementSystem", "qhseCost", "o34"),
		),
		receivableRecovered: firstNonZero(
			lookupAnyNumber(derived, "receivableRecovered", "o36"),
			sumQuarterMetric(value.Quarter.ReceivableUpdate, false, "receivableRecovered", "receivableCollection", "cashCollection", "o36"),
		),
		salesRevenue: firstNonZero(
			lookupAnyNumber(derived, "salesRevenue", "deliverySalesRevenue", "o38"),
			sumQuarterMetric(value.Quarter.DeliverySettlement, false, "salesRevenue", "deliverySalesRevenue", "orderSalesRevenue", "o38"),
		),
		directCost: firstNonZero(
			lookupAnyNumber(derived, "directCost", "deliveryDirectCost", "o39"),
			sumQuarterMetric(value.Quarter.DeliverySettlement, false, "directCost", "deliveryDirectCost", "orderCost", "o39"),
		),
		managementSalary: firstNonZero(
			lookupAnyNumber(derived, "managementSalary", "o40"),
			sumQuarterMetric(value.Quarter.DeliverySettlement, false, "managementSalary", "managementStaffCost", "adminSalary", "o40"),
		),
		longTermInterest: firstNonZero(
			lookupAnyNumber(derived, "longTermInterest", "o42"),
			extractMetric(value.YearEnd.LongTermLoan, false, "longTermInterest", "interest", "loanInterest", "o42"),
		),
		longTermRepayment: firstNonZero(
			lookupAnyNumber(derived, "longTermRepayment", "o43"),
			extractMetric(value.YearEnd.LongTermLoan, false, "longTermRepayment", "repayment", "dueRepayment", "o43"),
		),
		newLongTermLoan: firstNonZero(
			lookupAnyNumber(derived, "newLongTermLoan", "o44"),
			extractMetric(value.YearEnd.LongTermLoan, false, "newLongTermLoan", "newLoan", "additionalLoan", "o44"),
		),
		lineMaintenance: firstNonZero(
			lookupAnyNumber(derived, "lineMaintenance", "o45"),
			extractMetric(value.YearEnd.AssetAdjustment, false, "lineMaintenance", "maintenanceCost", "annualMaintenance", "o45"),
		),
		factoryPurchase: firstNonZero(
			lookupAnyNumber(derived, "factoryPurchase", "o46"),
			extractMetric(value.YearEnd.AssetAdjustment, false, "factoryPurchase", "purchase", "assetPurchase", "o46"),
		),
		factorySale: firstNonZero(
			lookupAnyNumber(derived, "factorySale", "o47"),
			extractMetric(value.YearEnd.AssetAdjustment, false, "factorySale", "sale", "assetSale", "o47"),
		),
		factoryRent: firstNonZero(
			lookupAnyNumber(derived, "factoryRent", "o48"),
			extractMetric(value.YearEnd.AssetAdjustment, false, "factoryRent", "rent", "rentalCost", "o48"),
		),
		workInConstruction: firstNonZero(
			lookupAnyNumber(derived, "workInConstruction", "unfinishedLineValue", "o52"),
			extractMetric(value.YearEnd.AssetAdjustment, false, "workInConstruction", "unfinishedLineValue", "unfinishedProductionLineValue", "o52"),
		),
		marketCultivation: extractMarketCultivationMetric(value, derived),
		discountExpense: firstNonZero(
			lookupAnyNumber(derived, "discountExpense", "o55"),
			sumQuarterMetric(value.Extra.IncomeAndPenalty, false, "discountExpense", "factoringExpense", "discountCost", "o55"),
		),
		extraExpensePenalty: firstNonZero(
			lookupAnyNumber(derived, "extraExpensePenalty", "o57"),
			sumAdjustmentMetric(value.Extra.IncomeAndPenalty, false, "extraExpensePenalty", "penalty", "fine", "o57"),
		),
		extraIncomeReward: firstNonZero(
			lookupAnyNumber(derived, "extraIncomeReward", "o58"),
			sumAdjustmentMetric(value.Extra.IncomeAndPenalty, false, "extraIncomeReward", "reward", "bonus", "o58"),
		),
	}
}

func extractMarketCultivationMetric(value payload.OperatingPayload, derived map[string]any) float64 {
	amount, ok := payload.MarketCultivationAnnualAmount(value.YearEnd.MarketCultivation)
	if ok {
		return amount
	}
	return firstNonZero(
		lookupAnyNumber(derived, "marketCultivation", "o53"),
		extractMetric(value.YearEnd.AssetAdjustment, false, "marketCultivation", "newMarketCultivation", "marketDevelopment", "o53"),
	)
}

func extractQuarterMetrics(value payload.OperatingPayload) map[string]quarterOperatingMetrics {
	result := map[string]quarterOperatingMetrics{}
	for _, quarter := range quarterSequence {
		result[quarter] = quarterOperatingMetrics{
			shortTermRepayment:   quarterMetric(value.Quarter.ShortTermLoan, quarter, false, "shortTermRepayment", "dueRepayment", "repayment", "o10"),
			shortTermInterest:    quarterMetric(value.Quarter.ShortTermLoan, quarter, false, "shortTermInterest", "interest", "loanInterest", "o11"),
			newShortTermLoan:     quarterMetric(value.Quarter.ShortTermLoan, quarter, false, "newShortTermLoan", "newLoan", "additionalLoan", "o12"),
			materialPayment:      quarterMetric(value.Quarter.MaterialPayment, quarter, true, "materialPayment", "materialPaymentTotal", "materialCost", "o18"),
			changeProductCost:    quarterMetric(value.Quarter.ProductionLineAdjust, quarter, false, "changeProductCost", "changeProduct", "productChange", "switchProduct", "o21"),
			lineDismantleCost:    quarterMetric(value.Quarter.ProductionLineAdjust, quarter, false, "lineDismantleCost", "dismantleCost", "removeCost", "o22"),
			lineSaleValue:        quarterMetric(value.Quarter.ProductionLineAdjust, quarter, false, "lineSaleValue", "lineSale", "sellValue", "saleValue", "o23"),
			newLineInstall:       quarterMetric(value.Quarter.ProductionLineAdjust, quarter, false, "newLineInstall", "installCost", "newLineInvestment", "o26"),
			humanResourceCost:    quarterMetric(value.Quarter.HumanResource, quarter, true, "humanResourceCost", "humanResourceFee", "staffCost", "o29"),
			salaryAndProduction:  quarterMetric(value.Quarter.SalaryAndProduction, quarter, true, "salaryAndProductionCost", "salaryCost", "workerSalary", "productionSalary", "o31"),
			researchCost:         quarterMetric(value.Quarter.ResearchAndManagement, quarter, false, "researchCost", "technologyResearch", "rdCost", "o33"),
			managementSystemCost: quarterMetric(value.Quarter.ResearchAndManagement, quarter, false, "managementSystemCost", "managementSystem", "qhseCost", "o34"),
			receivableRecovered:  quarterMetric(value.Quarter.ReceivableUpdate, quarter, false, "receivableRecovered", "receivableCollection", "cashCollection", "o36"),
			managementSalary:     quarterMetric(value.Quarter.DeliverySettlement, quarter, false, "managementSalary", "managementStaffCost", "adminSalary", "o40"),
			discountExpense:      quarterMetric(value.Extra.IncomeAndPenalty, quarter, false, "discountExpense", "factoringExpense", "discountCost", "o55"),
			extraExpensePenalty:  quarterMetric(value.Extra.IncomeAndPenalty, quarter, false, "extraExpensePenalty", "penalty", "fine", "o57"),
			extraIncomeReward:    quarterMetric(value.Extra.IncomeAndPenalty, quarter, false, "extraIncomeReward", "reward", "bonus", "o58"),
		}
	}
	return result
}

func calculateQuarterCashChecks(carryBase operatingCarryBase, marketBidCost float64, quarterMetrics map[string]quarterOperatingMetrics) map[string]float64 {
	result := map[string]float64{"Q1": 0, "Q2": 0, "Q3": 0, "Q4": 0}

	runningCash := carryBase.previousCash - carryBase.previousIncomeTax - marketBidCost
	for _, quarter := range quarterSequence {
		metrics := quarterMetrics[quarter]
		runningCash += metrics.newShortTermLoan + metrics.lineSaleValue + metrics.receivableRecovered + metrics.extraIncomeReward
		runningCash -= metrics.shortTermRepayment + metrics.shortTermInterest + metrics.materialPayment + metrics.changeProductCost + metrics.lineDismantleCost + metrics.newLineInstall + metrics.humanResourceCost + metrics.salaryAndProduction + metrics.researchCost + metrics.managementSystemCost + metrics.managementSalary + metrics.discountExpense + metrics.extraExpensePenalty
		result[quarter] = runningCash
	}

	return result
}

func buildDerivedValues(carryBase operatingCarryBase, metrics annualOperatingMetrics, quarterCashChecks map[string]float64) map[string]float64 {
	comprehensiveCost := metrics.marketBidCost +
		metrics.changeProductCost +
		metrics.lineDismantleCost +
		metrics.humanResourceCost +
		metrics.researchCost +
		metrics.managementSystemCost +
		metrics.managementSalary +
		metrics.lineMaintenance +
		metrics.factoryRent +
		metrics.marketCultivation

	shortTermLoanBalance := carryBase.previousShortTermLoan + metrics.newShortTermLoan - metrics.shortTermRepayment
	longTermLoanBalance := carryBase.previousLongTermLoan + metrics.newLongTermLoan - metrics.longTermRepayment
	lineResidual := carryBase.previousLineResidual + metrics.transferToFixed - metrics.lineSaleValue
	depreciableAssetTotal := carryBase.previousDepreciableAsset + metrics.depreciableAssetIncr
	depreciation := excelRoundToInteger(depreciableAssetTotal / 3)
	financeIncomeExpense := metrics.shortTermInterest + metrics.longTermInterest + metrics.discountExpense
	factoryAssetChange := metrics.factoryPurchase - metrics.factorySale
	receivableChange := metrics.salesRevenue - metrics.receivableRecovered
	extraIncomeExpense := metrics.extraIncomeReward - metrics.extraExpensePenalty
	lineResidualChange := metrics.transferToFixed - metrics.lineSaleValue
	depreciableAssetChange := metrics.depreciableAssetIncr - depreciation
	marketReturnRatio := safeDivide(metrics.orderTotal, metrics.marketBidCost)
	researchIntensity := safeDivide(metrics.researchCost, metrics.salesRevenue+metrics.extraIncomeReward)
	laborProductivity := calculateLaborProductivity(carryBase, metrics)
	periodEndCash := carryBase.previousCash +
		(metrics.newShortTermLoan + metrics.lineSaleValue + metrics.receivableRecovered + metrics.newLongTermLoan + metrics.factorySale + metrics.extraIncomeReward) -
		(carryBase.previousIncomeTax + metrics.marketBidCost + metrics.shortTermRepayment + metrics.shortTermInterest + metrics.materialPayment + metrics.changeProductCost + metrics.lineDismantleCost + metrics.newLineInstall + metrics.humanResourceCost + metrics.salaryAndProduction + metrics.researchCost + metrics.managementSystemCost + metrics.managementSalary + metrics.lineMaintenance + metrics.factoryPurchase + metrics.longTermInterest + metrics.longTermRepayment + metrics.factoryRent + metrics.marketCultivation + metrics.discountExpense + metrics.extraExpensePenalty)

	return map[string]float64{
		"marketBidCost":            metrics.marketBidCost,
		"orderTotal":               metrics.orderTotal,
		"comprehensiveCostTotal":   comprehensiveCost,
		"shortTermRepayment":       metrics.shortTermRepayment,
		"shortTermInterest":        metrics.shortTermInterest,
		"newShortTermLoan":         metrics.newShortTermLoan,
		"shortTermLoanBalance":     shortTermLoanBalance,
		"materialPayment":          metrics.materialPayment,
		"changeProductCost":        metrics.changeProductCost,
		"lineDismantleCost":        metrics.lineDismantleCost,
		"lineSaleValue":            metrics.lineSaleValue,
		"newLineInstall":           metrics.newLineInstall,
		"transferToFixed":          metrics.transferToFixed,
		"depreciableAssetIncrease": metrics.depreciableAssetIncr,
		"humanResourceCost":        metrics.humanResourceCost,
		"salaryAndProductionCost":  metrics.salaryAndProduction,
		"researchCost":             metrics.researchCost,
		"managementSystemCost":     metrics.managementSystemCost,
		"receivableRecovered":      metrics.receivableRecovered,
		"salesRevenue":             metrics.salesRevenue,
		"directCost":               metrics.directCost,
		"managementSalary":         metrics.managementSalary,
		"longTermInterest":         metrics.longTermInterest,
		"longTermRepayment":        metrics.longTermRepayment,
		"newLongTermLoan":          metrics.newLongTermLoan,
		"longTermLoanBalance":      longTermLoanBalance,
		"lineMaintenance":          metrics.lineMaintenance,
		"factoryPurchase":          metrics.factoryPurchase,
		"factorySale":              metrics.factorySale,
		"factoryAssetChange":       factoryAssetChange,
		"factoryRent":              metrics.factoryRent,
		"workInConstruction":       metrics.workInConstruction,
		"marketCultivation":        metrics.marketCultivation,
		"marketReturnRatio":        marketReturnRatio,
		"researchIntensity":        researchIntensity,
		"laborProductivity":        laborProductivity,
		"lineResidual":             lineResidual,
		"depreciableAssetTotal":    depreciableAssetTotal,
		"depreciation":             depreciation,
		"financeIncomeExpense":     financeIncomeExpense,
		"receivableChange":         receivableChange,
		"extraIncomeExpense":       extraIncomeExpense,
		"lineResidualChange":       lineResidualChange,
		"depreciableAssetChange":   depreciableAssetChange,
		"discountExpense":          metrics.discountExpense,
		"extraExpensePenalty":      metrics.extraExpensePenalty,
		"extraIncomeReward":        metrics.extraIncomeReward,
		"periodEndCash":            periodEndCash,
		"q1QuarterEndCashCheck":    quarterCashChecks["Q1"],
		"q2QuarterEndCashCheck":    quarterCashChecks["Q2"],
		"q3QuarterEndCashCheck":    quarterCashChecks["Q3"],
		"q4QuarterEndCashCheck":    quarterCashChecks["Q4"],
	}
}

// BuildReportIndicatorValues 基于财报口径返回经营页底部需要展示的收益率指标。
func BuildReportIndicatorValues(report payload.ReportComputedPayload) map[string]float64 {
	return map[string]float64{
		"netAssetYield":   safeDivide(report.ReportNetProfit, report.ReportTotalEquity),
		"totalAssetYield": safeDivide(report.ReportNetProfit, report.ReportTotalAssets),
		"netProfitRate":   safeDivide(report.ReportNetProfit, report.ReportSalesRevenue),
		"grossMarginRate": 1 - safeDivide(report.ReportDirectCost, report.ReportSalesRevenue),
	}
}

func calculateLaborProductivity(carryBase operatingCarryBase, metrics annualOperatingMetrics) float64 {
	numerator := metrics.salesRevenue -
		metrics.directCost +
		metrics.shortTermInterest -
		metrics.marketBidCost -
		carryBase.previousIncomeTax -
		metrics.changeProductCost -
		metrics.lineDismantleCost +
		metrics.humanResourceCost -
		metrics.researchCost -
		metrics.managementSystemCost -
		metrics.longTermInterest -
		metrics.lineMaintenance -
		metrics.factoryRent -
		metrics.marketCultivation -
		metrics.extraExpensePenalty +
		metrics.extraIncomeReward

	denominator := metrics.salaryAndProduction*7 + 5
	return safeDivide(numerator, denominator) * 100
}

func safeDivide(numerator float64, denominator float64) float64 {
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

func selectPeriodEndCash(stageStatus string, quarterCashChecks map[string]float64, yearEndCash float64) float64 {
	switch stageStatus {
	case enum.StageStatusQ1Open:
		return quarterCashChecks["Q1"]
	case enum.StageStatusQ2Open:
		return quarterCashChecks["Q2"]
	case enum.StageStatusQ3Open:
		return quarterCashChecks["Q3"]
	case enum.StageStatusQ4Open:
		return quarterCashChecks["Q4"]
	case enum.StageStatusYearEndOpen:
		return yearEndCash
	default:
		return yearEndCash
	}
}

func mergeDerivedValues(existing map[string]any, computed map[string]float64) map[string]any {
	result := map[string]any{}
	for key, value := range existing {
		result[key] = value
	}
	for key, value := range computed {
		result[key] = value
	}
	return result
}

func sumQuarterMetric(source payload.OperatingQuarterMap, allowFallbackTotal bool, keys ...string) float64 {
	var total float64
	for _, quarter := range quarterSequence {
		total += quarterMetric(source, quarter, allowFallbackTotal, keys...)
	}
	return total
}

func sumAdjustmentMetric(source payload.OperatingQuarterMap, allowFallbackTotal bool, keys ...string) float64 {
	return sumQuarterMetric(source, allowFallbackTotal, keys...) +
		quarterMetric(source, "YEAR_END", allowFallbackTotal, keys...)
}

func quarterMetric(source payload.OperatingQuarterMap, quarter string, allowFallbackTotal bool, keys ...string) float64 {
	if len(source) == 0 {
		return 0
	}
	for key, value := range source {
		if normalizeSearchKey(key) != normalizeSearchKey(quarter) {
			continue
		}
		return extractMetric(value, allowFallbackTotal, keys...)
	}
	return 0
}

func extractMetric(source any, allowFallbackTotal bool, keys ...string) float64 {
	if value, ok := sumMatchedNumbers(source, keys...); ok {
		return value
	}
	if allowFallbackTotal {
		return sumAllNumbers(source)
	}
	return 0
}

func lookupAnyNumber(source map[string]any, keys ...string) float64 {
	if len(source) == 0 {
		return 0
	}
	if value, ok := sumMatchedNumbers(source, keys...); ok {
		return value
	}
	return 0
}

func lookupPrioritizedNumber(source map[string]any, keys ...string) float64 {
	if len(source) == 0 {
		return 0
	}
	for _, key := range keys {
		if value, ok := findMatchedNumber(source, key); ok {
			return value
		}
	}
	return 0
}

func findMatchedNumber(source any, key string) (float64, bool) {
	normalized := normalizeSearchKey(key)
	if normalized == "" {
		return 0, false
	}
	return walkFirstMatchedNumber(source, normalized)
}

func walkFirstMatchedNumber(source any, key string) (float64, bool) {
	switch typed := source.(type) {
	case map[string]any:
		for currentKey, value := range typed {
			if keyMatches(currentKey, []string{key}) {
				if number, ok := numberValue(value); ok {
					return number, true
				}
				if nested, nestedFound := sumAllNumbersWithFound(value); nestedFound {
					return nested, true
				}
			}

			if nested, nestedFound := walkFirstMatchedNumber(value, key); nestedFound {
				return nested, true
			}
		}
		return 0, false
	case []map[string]any:
		for _, item := range typed {
			if nested, nestedFound := walkFirstMatchedNumber(item, key); nestedFound {
				return nested, true
			}
		}
		return 0, false
	case []any:
		for _, item := range typed {
			if nested, nestedFound := walkFirstMatchedNumber(item, key); nestedFound {
				return nested, true
			}
		}
		return 0, false
	default:
		return 0, false
	}
}

func sumMatchedNumbers(source any, keys ...string) (float64, bool) {
	normalizedKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		normalized := normalizeSearchKey(key)
		if normalized == "" {
			continue
		}
		normalizedKeys = append(normalizedKeys, normalized)
	}
	if len(normalizedKeys) == 0 {
		return 0, false
	}
	return walkMatchedNumbers(source, normalizedKeys)
}

func walkMatchedNumbers(source any, keys []string) (float64, bool) {
	switch typed := source.(type) {
	case map[string]any:
		var total float64
		found := false
		for key, value := range typed {
			if keyMatches(key, keys) {
				if number, ok := numberValue(value); ok {
					total += number
					found = true
					continue
				}
				if nested, nestedFound := sumAllNumbersWithFound(value); nestedFound {
					total += nested
					found = true
					continue
				}
			}

			if nested, nestedFound := walkMatchedNumbers(value, keys); nestedFound {
				total += nested
				found = true
			}
		}
		return total, found
	case []map[string]any:
		var total float64
		found := false
		for _, item := range typed {
			if nested, nestedFound := walkMatchedNumbers(item, keys); nestedFound {
				total += nested
				found = true
			}
		}
		return total, found
	case []any:
		var total float64
		found := false
		for _, item := range typed {
			if nested, nestedFound := walkMatchedNumbers(item, keys); nestedFound {
				total += nested
				found = true
			}
		}
		return total, found
	default:
		return 0, false
	}
}

func sumAllNumbers(source any) float64 {
	total, _ := sumAllNumbersWithFound(source)
	return total
}

func sumAllNumbersWithFound(source any) (float64, bool) {
	switch typed := source.(type) {
	case map[string]any:
		var total float64
		found := false
		for _, value := range typed {
			if nested, nestedFound := sumAllNumbersWithFound(value); nestedFound {
				total += nested
				found = true
			}
		}
		return total, found
	case []map[string]any:
		var total float64
		found := false
		for _, item := range typed {
			if nested, nestedFound := sumAllNumbersWithFound(item); nestedFound {
				total += nested
				found = true
			}
		}
		return total, found
	case []any:
		var total float64
		found := false
		for _, item := range typed {
			if nested, nestedFound := sumAllNumbersWithFound(item); nestedFound {
				total += nested
				found = true
			}
		}
		return total, found
	default:
		number, ok := numberValue(source)
		return number, ok
	}
}

func numberValue(source any) (float64, bool) {
	switch typed := source.(type) {
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case string:
		value := strings.TrimSpace(typed)
		if value == "" {
			return 0, false
		}
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func keyMatches(source string, keys []string) bool {
	normalizedSource := normalizeSearchKey(source)
	if normalizedSource == "" {
		return false
	}
	for _, key := range keys {
		if normalizedSource == key || strings.Contains(normalizedSource, key) || strings.Contains(key, normalizedSource) {
			return true
		}
	}
	return false
}

func normalizeSearchKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}
	var builder strings.Builder
	builder.Grow(len(value))
	for _, ch := range value {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}

func excelRoundToInteger(value float64) float64 {
	return math.Round(value)
}

func firstNonZero(values ...float64) float64 {
	for _, value := range values {
		if value != 0 {
			return value
		}
	}
	return 0
}
