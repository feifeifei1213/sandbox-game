package payload

// ReportComputedPayload 对应财报自动计算结果。
type ReportComputedPayload struct {
	ReportSalesRevenue          float64 `json:"reportSalesRevenue"`
	ReportDirectCost            float64 `json:"reportDirectCost"`
	ReportGrossProfit           float64 `json:"reportGrossProfit"`
	ReportComprehensiveCost     float64 `json:"reportComprehensiveCost"`
	ReportDepreciation          float64 `json:"reportDepreciation"`
	ReportOperatingProfit       float64 `json:"reportOperatingProfit"`
	ReportFinanceIncomeExpense  float64 `json:"reportFinanceIncomeExpense"`
	ReportExtraIncomeExpense    float64 `json:"reportExtraIncomeExpense"`
	ReportPreTaxProfit          float64 `json:"reportPreTaxProfit"`
	ReportIncomeTax             float64 `json:"reportIncomeTax"`
	ReportNetProfit             float64 `json:"reportNetProfit"`
	ReportWorkInProgress        float64 `json:"reportWorkInProgress"`
	ReportFinishedGoods         float64 `json:"reportFinishedGoods"`
	ReportRawMaterials          float64 `json:"reportRawMaterials"`
	ReportWorkInConstruction    float64 `json:"reportWorkInConstruction"`
	ReportFactoryAsset          float64 `json:"reportFactoryAsset"`
	ReportLineResidual          float64 `json:"reportLineResidual"`
	ReportDepreciableAsset      float64 `json:"reportDepreciableAsset"`
	ReportTotalNonCurrentAssets float64 `json:"reportTotalNonCurrentAssets"`
	ReportCash                  float64 `json:"reportCash"`
	ReportReceivable            float64 `json:"reportReceivable"`
	ReportPostTaxCash           float64 `json:"reportPostTaxCash"`
	ReportTotalCurrentAssets    float64 `json:"reportTotalCurrentAssets"`
	ReportTotalAssets           float64 `json:"reportTotalAssets"`
	ReportShortTermLiability    float64 `json:"reportShortTermLiability"`
	ReportLongTermLiability     float64 `json:"reportLongTermLiability"`
	ReportTotalLiability        float64 `json:"reportTotalLiability"`
	ReportShareCapital          float64 `json:"reportShareCapital"`
	ReportRetainedEarnings      float64 `json:"reportRetainedEarnings"`
	ReportTotalEquity           float64 `json:"reportTotalEquity"`
	ReportTotalLiabilityEquity  float64 `json:"reportTotalLiabilityEquity"`
}

// BalanceGap 返回资产负债平衡差额。
func (p ReportComputedPayload) BalanceGap() float64 {
	return p.ReportTotalAssets - p.ReportTotalLiabilityEquity
}
