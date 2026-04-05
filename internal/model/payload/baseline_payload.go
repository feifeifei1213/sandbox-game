package payload

// BaselinePayload 对应管理员录入的初始年财报基线。
type BaselinePayload struct {
	BaselineSalesRevenue         float64 `json:"baselineSalesRevenue"`
	BaselineDirectCost           float64 `json:"baselineDirectCost"`
	BaselineComprehensiveCost    float64 `json:"baselineComprehensiveCost"`
	BaselineDepreciation         float64 `json:"baselineDepreciation"`
	BaselineFinanceIncomeExpense float64 `json:"baselineFinanceIncomeExpense"`
	BaselineExtraIncomeExpense   float64 `json:"baselineExtraIncomeExpense"`
	BaselineIncomeTax            float64 `json:"baselineIncomeTax"`
	BaselineWorkInConstruction   float64 `json:"baselineWorkInConstruction"`
	BaselineFactoryAsset         float64 `json:"baselineFactoryAsset"`
	BaselineLineResidual         float64 `json:"baselineLineResidual"`
	BaselineDepreciableAsset     float64 `json:"baselineDepreciableAsset"`
	BaselineCash                 float64 `json:"baselineCash"`
	BaselineReceivable           float64 `json:"baselineReceivable"`
	BaselineWorkInProgress       float64 `json:"baselineWorkInProgress"`
	BaselineFinishedGoods        float64 `json:"baselineFinishedGoods"`
	BaselineRawMaterials         float64 `json:"baselineRawMaterials"`
	BaselineShortTermLoan        float64 `json:"baselineShortTermLoan"`
	BaselineLongTermLoan         float64 `json:"baselineLongTermLoan"`
	BaselineShareCapital         float64 `json:"baselineShareCapital"`
	BaselineRetainedEarnings     float64 `json:"baselineRetainedEarnings"`
}
