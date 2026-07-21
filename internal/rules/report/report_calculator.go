package report

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
)

type reportCarryBase struct {
	previousIncomeTax        float64
	previousShortTermLoan    float64
	previousLongTermLoan     float64
	previousFactoryAsset     float64
	previousLineResidual     float64
	previousDepreciableAsset float64
	previousCash             float64
	previousReceivable       float64
	shareCapital             float64
	retainedEarnings         float64
	previousTotalEquity      float64
	previousBestMarketScore  float64
	previousBestTechScore    float64
	previousBestSalesScore   float64
	previousBestCfoBaseScore float64
}

type operatingReportMetrics struct {
	orderTotal           float64
	marketBidCost        float64
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

// Calculator 是财报计算统一入口。
type Calculator struct{}

// NewCalculator 创建财报计算器。
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Calculate 按已确认的 Excel 公式链计算财报自动结果。
func (c *Calculator) Calculate(ctx calcctx.CalculationContext) (payload.ReportComputedPayload, error) {
	if err := ctx.Validate(); err != nil {
		return payload.ReportComputedPayload{}, err
	}
	if ctx.OperatingPayload == nil {
		return payload.ReportComputedPayload{}, fmt.Errorf("missing operating payload for report calculation")
	}

	carryBase, err := buildCarryBase(ctx)
	if err != nil {
		return payload.ReportComputedPayload{}, err
	}

	metrics := extractOperatingMetrics(*ctx.OperatingPayload)
	manual := normalizeManualPayload(ctx.ReportManualPayload)

	reportComprehensiveCost := metrics.marketBidCost +
		metrics.changeProductCost +
		metrics.lineDismantleCost +
		metrics.humanResourceCost +
		metrics.researchCost +
		metrics.managementSystemCost +
		metrics.managementSalary +
		metrics.lineMaintenance +
		metrics.factoryRent +
		metrics.marketCultivation

	reportDepreciation := excelRoundToInteger((carryBase.previousDepreciableAsset + metrics.depreciableAssetIncr) / 3)
	reportGrossProfit := metrics.salesRevenue - metrics.directCost
	reportOperatingProfit := reportGrossProfit - reportComprehensiveCost - reportDepreciation
	reportFinanceIncomeExpense := metrics.shortTermInterest + metrics.longTermInterest + metrics.discountExpense
	reportExtraIncomeExpense := metrics.extraIncomeReward - metrics.extraExpensePenalty
	reportPreTaxProfit := reportOperatingProfit - reportFinanceIncomeExpense + reportExtraIncomeExpense

	incomeTaxRate := manual.incomeTaxRate
	reportIncomeTax := math.Max(excelRoundToInteger(reportPreTaxProfit*incomeTaxRate), 0)
	reportNetProfit := reportPreTaxProfit - reportIncomeTax

	reportFactoryAsset := carryBase.previousFactoryAsset + (metrics.factoryPurchase - metrics.factorySale)
	reportLineResidual := carryBase.previousLineResidual + metrics.transferToFixed - metrics.lineSaleValue
	reportDepreciableAsset := carryBase.previousDepreciableAsset + metrics.depreciableAssetIncr - reportDepreciation
	reportTotalNonCurrentAssets := metrics.workInConstruction + reportFactoryAsset + reportLineResidual + reportDepreciableAsset

	reportCash := carryBase.previousCash +
		(metrics.newShortTermLoan + metrics.lineSaleValue + metrics.receivableRecovered + metrics.newLongTermLoan + metrics.factorySale + metrics.extraIncomeReward) -
		(carryBase.previousIncomeTax + metrics.marketBidCost + metrics.shortTermRepayment + metrics.shortTermInterest + metrics.materialPayment + metrics.changeProductCost + metrics.lineDismantleCost + metrics.newLineInstall + metrics.humanResourceCost + metrics.salaryAndProduction + metrics.researchCost + metrics.managementSystemCost + metrics.managementSalary + metrics.lineMaintenance + metrics.factoryPurchase + metrics.longTermInterest + metrics.longTermRepayment + metrics.factoryRent + metrics.marketCultivation + metrics.discountExpense + metrics.extraExpensePenalty)

	reportReceivable := carryBase.previousReceivable + metrics.salesRevenue - metrics.receivableRecovered
	reportPostTaxCash := reportCash - reportIncomeTax
	reportTotalCurrentAssets := reportReceivable + manual.workInProgress + manual.finishedGoods + manual.rawMaterials + reportPostTaxCash
	reportTotalAssets := reportTotalNonCurrentAssets + reportTotalCurrentAssets

	reportShortTermLiability := carryBase.previousShortTermLoan + metrics.newShortTermLoan - metrics.shortTermRepayment
	reportLongTermLiability := carryBase.previousLongTermLoan + metrics.newLongTermLoan - metrics.longTermRepayment
	reportTotalLiability := reportShortTermLiability + reportLongTermLiability
	reportTotalEquity := carryBase.shareCapital + carryBase.retainedEarnings + reportNetProfit
	reportTotalLiabilityEquity := reportTotalLiability + reportTotalEquity
	reportBestMarketDirectorBaseScore := carryBase.previousBestMarketScore + metrics.marketCultivation
	reportBestMarketDirectorScore := reportBestMarketDirectorBaseScore + manual.enterpriseCertificationScore
	reportBestTechnologyDirectorScore := carryBase.previousBestTechScore + metrics.researchCost
	// 销售总监得分按截至当年的累计订单总额整体除以 10；
	// 上一年得分已经完成缩放，因此递推时只缩放本年新增订单总额。
	reportBestSalesDirectorScore := carryBase.previousBestSalesScore + metrics.orderTotal/10
	reportBestCfoBaseScore := 0.0
	if !ctx.IsDemoYear() {
		// Excel 的 CFO 黄色分始终以 0 年财报权益为基线；
		// 这里用“上一年黄色基础分 + 本年权益增量的一半”做等价累计。
		reportBestCfoBaseScore = carryBase.previousBestCfoBaseScore + 0.5*(reportTotalEquity-carryBase.previousTotalEquity)
	}
	reportBestCfoScore := reportBestCfoBaseScore + manual.closingSpeedScore
	reportBestCeoScore := reportBestMarketDirectorScore +
		reportBestTechnologyDirectorScore +
		manual.productionHumanScore +
		reportBestSalesDirectorScore +
		reportBestCfoScore +
		reportTotalEquity/2

	return payload.ReportComputedPayload{
		ReportSalesRevenue:                metrics.salesRevenue,
		ReportDirectCost:                  metrics.directCost,
		ReportGrossProfit:                 reportGrossProfit,
		ReportComprehensiveCost:           reportComprehensiveCost,
		ReportDepreciation:                reportDepreciation,
		ReportOperatingProfit:             reportOperatingProfit,
		ReportFinanceIncomeExpense:        reportFinanceIncomeExpense,
		ReportExtraIncomeExpense:          reportExtraIncomeExpense,
		ReportPreTaxProfit:                reportPreTaxProfit,
		ReportIncomeTax:                   reportIncomeTax,
		ReportNetProfit:                   reportNetProfit,
		ReportWorkInProgress:              manual.workInProgress,
		ReportFinishedGoods:               manual.finishedGoods,
		ReportRawMaterials:                manual.rawMaterials,
		ReportWorkInConstruction:          metrics.workInConstruction,
		ReportFactoryAsset:                reportFactoryAsset,
		ReportLineResidual:                reportLineResidual,
		ReportDepreciableAsset:            reportDepreciableAsset,
		ReportTotalNonCurrentAssets:       reportTotalNonCurrentAssets,
		ReportCash:                        reportCash,
		ReportReceivable:                  reportReceivable,
		ReportPostTaxCash:                 reportPostTaxCash,
		ReportTotalCurrentAssets:          reportTotalCurrentAssets,
		ReportTotalAssets:                 reportTotalAssets,
		ReportShortTermLiability:          reportShortTermLiability,
		ReportLongTermLiability:           reportLongTermLiability,
		ReportTotalLiability:              reportTotalLiability,
		ReportShareCapital:                carryBase.shareCapital,
		ReportRetainedEarnings:            carryBase.retainedEarnings,
		ReportTotalEquity:                 reportTotalEquity,
		ReportTotalLiabilityEquity:        reportTotalLiabilityEquity,
		ReportBestMarketDirectorBaseScore: reportBestMarketDirectorBaseScore,
		ReportBestMarketDirectorScore:     reportBestMarketDirectorScore,
		ReportBestTechnologyDirectorScore: reportBestTechnologyDirectorScore,
		ReportBestSalesDirectorScore:      reportBestSalesDirectorScore,
		ReportBestCfoBaseScore:            reportBestCfoBaseScore,
		ReportBestCfoScore:                reportBestCfoScore,
		ReportBestCeoScore:                reportBestCeoScore,
	}, nil
}

type normalizedManualPayload struct {
	workInProgress               float64
	finishedGoods                float64
	rawMaterials                 float64
	incomeTaxRate                float64
	enterpriseCertificationScore float64
	productionHumanScore         float64
	closingSpeedScore            float64
}

func normalizeManualPayload(manual *payload.ReportManualPayload) normalizedManualPayload {
	result := normalizedManualPayload{}
	if manual == nil {
		return result
	}
	if manual.WorkInProgress != nil {
		result.workInProgress = *manual.WorkInProgress
	}
	if manual.FinishedGoods != nil {
		result.finishedGoods = *manual.FinishedGoods
	}
	if manual.RawMaterials != nil {
		result.rawMaterials = *manual.RawMaterials
	}
	if manual.IncomeTaxRate != nil {
		result.incomeTaxRate = *manual.IncomeTaxRate
	}
	if manual.EnterpriseCertificationScore != nil {
		result.enterpriseCertificationScore = *manual.EnterpriseCertificationScore
	}
	if manual.ProductionHumanScore != nil {
		result.productionHumanScore = *manual.ProductionHumanScore
	}
	if manual.ClosingSpeedScore != nil {
		result.closingSpeedScore = *manual.ClosingSpeedScore
	}
	return result
}

func buildCarryBase(ctx calcctx.CalculationContext) (reportCarryBase, error) {
	if ctx.IsDemoYear() {
		if ctx.InitialBaseline == nil {
			return reportCarryBase{}, fmt.Errorf("missing initial baseline for demo year report")
		}
		return reportCarryBase{
			previousIncomeTax:        ctx.InitialBaseline.BaselineIncomeTax,
			previousShortTermLoan:    ctx.InitialBaseline.BaselineShortTermLoan,
			previousLongTermLoan:     ctx.InitialBaseline.BaselineLongTermLoan,
			previousFactoryAsset:     ctx.InitialBaseline.BaselineFactoryAsset,
			previousLineResidual:     ctx.InitialBaseline.BaselineLineResidual,
			previousDepreciableAsset: ctx.InitialBaseline.BaselineDepreciableAsset,
			previousCash:             ctx.InitialBaseline.BaselineCash,
			previousReceivable:       ctx.InitialBaseline.BaselineReceivable,
			shareCapital:             ctx.InitialBaseline.BaselineShareCapital,
			retainedEarnings:         ctx.InitialBaseline.BaselineRetainedEarnings + calculateBaselineNetProfit(*ctx.InitialBaseline),
			previousTotalEquity:      ctx.InitialBaseline.BaselineShareCapital + ctx.InitialBaseline.BaselineRetainedEarnings + calculateBaselineNetProfit(*ctx.InitialBaseline),
		}, nil
	}

	if ctx.PreviousReport == nil {
		return reportCarryBase{}, fmt.Errorf("missing previous report for formal year report")
	}

	return reportCarryBase{
		previousIncomeTax:        ctx.PreviousReport.ReportIncomeTax,
		previousShortTermLoan:    ctx.PreviousReport.ReportShortTermLiability,
		previousLongTermLoan:     ctx.PreviousReport.ReportLongTermLiability,
		previousFactoryAsset:     ctx.PreviousReport.ReportFactoryAsset,
		previousLineResidual:     ctx.PreviousReport.ReportLineResidual,
		previousDepreciableAsset: ctx.PreviousReport.ReportDepreciableAsset,
		previousCash:             ctx.PreviousReport.ReportCash,
		previousReceivable:       ctx.PreviousReport.ReportReceivable,
		shareCapital:             ctx.PreviousReport.ReportShareCapital,
		retainedEarnings:         ctx.PreviousReport.ReportRetainedEarnings + ctx.PreviousReport.ReportNetProfit,
		previousTotalEquity:      ctx.PreviousReport.ReportTotalEquity,
		previousBestMarketScore:  ctx.PreviousReport.ReportBestMarketDirectorBaseScore,
		previousBestTechScore:    ctx.PreviousReport.ReportBestTechnologyDirectorScore,
		previousBestSalesScore:   ctx.PreviousReport.ReportBestSalesDirectorScore,
		previousBestCfoBaseScore: ctx.PreviousReport.ReportBestCfoBaseScore,
	}, nil
}

func calculateBaselineNetProfit(baseline payload.BaselinePayload) float64 {
	grossProfit := baseline.BaselineSalesRevenue - baseline.BaselineDirectCost
	operatingProfit := grossProfit - baseline.BaselineComprehensiveCost - baseline.BaselineDepreciation
	preTaxProfit := operatingProfit - baseline.BaselineFinanceIncomeExpense + baseline.BaselineExtraIncomeExpense
	return preTaxProfit - baseline.BaselineIncomeTax
}

func extractOperatingMetrics(value payload.OperatingPayload) operatingReportMetrics {
	derived := value.Derived.Values
	result := operatingReportMetrics{}

	result.orderTotal = firstNonZero(
		lookupPrioritizedNumber(derived, "orderTotal", "marketOrderTotal", "o5"),
		lookupPrioritizedNumber(value.Beginning.TaxAndPlanning, "orderTotal", "marketOrderTotal", "o5"),
		extractMetric(value.Beginning.MarketBid, true, "orderTotal", "orderAmount", "marketOrderTotal", "o5"),
	)
	result.marketBidCost = firstNonZero(
		lookupPrioritizedNumber(derived, "marketBidCost", "marketInvestmentTotal", "p5"),
		lookupPrioritizedNumber(value.Beginning.TaxAndPlanning, "marketBidCost", "marketInvestmentTotal", "p5"),
		extractMetric(value.Beginning.MarketBid, true, "marketBidCost", "marketInvestment", "marketInvestmentTotal", "investment", "bidInvestment", "p5"),
	)
	result.shortTermRepayment = firstNonZero(
		lookupAnyNumber(derived, "shortTermRepayment", "o10"),
		sumQuarterMetric(value.Quarter.ShortTermLoan, false, "shortTermRepayment", "dueRepayment", "repayment", "o10"),
	)
	result.shortTermInterest = firstNonZero(
		lookupAnyNumber(derived, "shortTermInterest", "financeShortTermInterest", "o11"),
		sumQuarterMetric(value.Quarter.ShortTermLoan, false, "shortTermInterest", "interest", "loanInterest", "o11"),
	)
	result.newShortTermLoan = firstNonZero(
		lookupAnyNumber(derived, "newShortTermLoan", "shortTermNewLoan", "o12"),
		sumQuarterMetric(value.Quarter.ShortTermLoan, false, "newShortTermLoan", "newLoan", "additionalLoan", "o12"),
	)
	result.materialPayment = firstNonZero(
		lookupAnyNumber(derived, "materialPayment", "materialPaymentTotal", "o18"),
		sumQuarterMetric(value.Quarter.MaterialPayment, true, "materialPayment", "materialPaymentTotal", "materialCost", "o18"),
	)
	result.changeProductCost = firstNonZero(
		lookupAnyNumber(derived, "changeProductCost", "o21"),
		sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "changeProductCost", "changeProduct", "productChange", "switchProduct", "o21"),
	)
	result.lineDismantleCost = firstNonZero(
		lookupAnyNumber(derived, "lineDismantleCost", "o22"),
		sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "lineDismantleCost", "dismantleCost", "removeCost", "o22"),
	)
	result.lineSaleValue = firstNonZero(
		lookupAnyNumber(derived, "lineSaleValue", "o23"),
		sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "lineSaleValue", "lineSale", "sellValue", "saleValue", "o23"),
	)
	result.newLineInstall = firstNonZero(
		lookupAnyNumber(derived, "newLineInstall", "o26"),
		sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "newLineInstall", "installCost", "newLineInvestment", "o26"),
	)
	result.transferToFixed = firstNonZero(
		lookupAnyNumber(derived, "transferToFixed", "lineResidualIncrease", "o27"),
		sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "transferToFixed", "constructionToFixed", "completedTransfer", "lineResidualIncrease", "o27"),
	)
	result.depreciableAssetIncr = firstNonZero(
		lookupAnyNumber(derived, "depreciableAssetIncrease", "o28"),
		sumQuarterMetric(value.Quarter.ProductionLineAdjust, false, "depreciableAssetIncrease", "assetCapitalization", "newDepreciableAsset", "o28"),
	)
	result.humanResourceCost = firstNonZero(
		lookupAnyNumber(derived, "humanResourceCost", "o29"),
		sumQuarterMetric(value.Quarter.HumanResource, true, "humanResourceCost", "humanResourceFee", "staffCost", "o29"),
	)
	result.salaryAndProduction = firstNonZero(
		lookupAnyNumber(derived, "salaryAndProductionCost", "o31"),
		sumQuarterMetric(value.Quarter.SalaryAndProduction, true, "salaryAndProductionCost", "salaryCost", "workerSalary", "productionSalary", "o31"),
	)
	result.researchCost = firstNonZero(
		lookupAnyNumber(derived, "researchCost", "o33"),
		sumQuarterMetric(value.Quarter.ResearchAndManagement, false, "researchCost", "technologyResearch", "rdCost", "o33"),
	)
	result.managementSystemCost = firstNonZero(
		lookupAnyNumber(derived, "managementSystemCost", "o34"),
		sumQuarterMetric(value.Quarter.ResearchAndManagement, false, "managementSystemCost", "managementSystem", "qhseCost", "o34"),
	)
	result.receivableRecovered = firstNonZero(
		lookupAnyNumber(derived, "receivableRecovered", "o36"),
		sumQuarterMetric(value.Quarter.ReceivableUpdate, false, "receivableRecovered", "receivableCollection", "cashCollection", "o36"),
	)
	result.salesRevenue = firstNonZero(
		lookupAnyNumber(derived, "salesRevenue", "deliverySalesRevenue", "o38"),
		sumQuarterMetric(value.Quarter.DeliverySettlement, false, "salesRevenue", "deliverySalesRevenue", "orderSalesRevenue", "o38"),
	)
	result.directCost = firstNonZero(
		lookupAnyNumber(derived, "directCost", "deliveryDirectCost", "o39"),
		sumQuarterMetric(value.Quarter.DeliverySettlement, false, "directCost", "deliveryDirectCost", "orderCost", "o39"),
	)
	result.managementSalary = firstNonZero(
		lookupAnyNumber(derived, "managementSalary", "o40"),
		sumQuarterMetric(value.Quarter.DeliverySettlement, false, "managementSalary", "managementStaffCost", "adminSalary", "o40"),
	)
	result.longTermInterest = firstNonZero(
		lookupAnyNumber(derived, "longTermInterest", "o42"),
		extractMetric(value.YearEnd.LongTermLoan, false, "longTermInterest", "interest", "loanInterest", "o42"),
	)
	result.longTermRepayment = firstNonZero(
		lookupAnyNumber(derived, "longTermRepayment", "o43"),
		extractMetric(value.YearEnd.LongTermLoan, false, "longTermRepayment", "repayment", "dueRepayment", "o43"),
	)
	result.newLongTermLoan = firstNonZero(
		lookupAnyNumber(derived, "newLongTermLoan", "o44"),
		extractMetric(value.YearEnd.LongTermLoan, false, "newLongTermLoan", "newLoan", "additionalLoan", "o44"),
	)
	result.lineMaintenance = firstNonZero(
		lookupAnyNumber(derived, "lineMaintenance", "o45"),
		extractMetric(value.YearEnd.AssetAdjustment, false, "lineMaintenance", "maintenanceCost", "annualMaintenance", "o45"),
	)
	result.factoryPurchase = firstNonZero(
		lookupAnyNumber(derived, "factoryPurchase", "o46"),
		extractMetric(value.YearEnd.AssetAdjustment, false, "factoryPurchase", "purchase", "assetPurchase", "o46"),
	)
	result.factorySale = firstNonZero(
		lookupAnyNumber(derived, "factorySale", "o47"),
		extractMetric(value.YearEnd.AssetAdjustment, false, "factorySale", "sale", "assetSale", "o47"),
	)
	result.factoryRent = firstNonZero(
		lookupAnyNumber(derived, "factoryRent", "o48"),
		extractMetric(value.YearEnd.AssetAdjustment, false, "factoryRent", "rent", "rentalCost", "o48"),
	)
	result.workInConstruction = firstNonZero(
		lookupAnyNumber(derived, "workInConstruction", "unfinishedLineValue", "o52"),
		extractMetric(value.YearEnd.AssetAdjustment, false, "workInConstruction", "unfinishedLineValue", "unfinishedProductionLineValue", "o52"),
	)
	result.marketCultivation = firstNonZero(
		lookupAnyNumber(derived, "marketCultivation", "o53"),
		extractMetric(value.YearEnd.AssetAdjustment, false, "marketCultivation", "newMarketCultivation", "marketDevelopment", "o53"),
	)
	result.discountExpense = firstNonZero(
		lookupAnyNumber(derived, "discountExpense", "o55"),
		sumQuarterMetric(value.Extra.IncomeAndPenalty, false, "discountExpense", "factoringExpense", "discountCost", "o55"),
	)
	result.extraExpensePenalty = firstNonZero(
		lookupAnyNumber(derived, "extraExpensePenalty", "o57"),
		sumAdjustmentMetric(value.Extra.IncomeAndPenalty, false, "extraExpensePenalty", "penalty", "fine", "o57"),
	)
	result.extraIncomeReward = firstNonZero(
		lookupAnyNumber(derived, "extraIncomeReward", "o58"),
		sumAdjustmentMetric(value.Extra.IncomeAndPenalty, false, "extraIncomeReward", "reward", "bonus", "o58"),
	)

	return result
}

func sumQuarterMetric(source payload.OperatingQuarterMap, allowFallbackTotal bool, keys ...string) float64 {
	var total float64
	for _, quarter := range []string{"Q1", "Q2", "Q3", "Q4"} {
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
