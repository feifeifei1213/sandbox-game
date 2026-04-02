package payload

var allowedIncomeTaxRates = []float64{0.25, 0.15, 0}

// ReportManualPayload 对应财报页绿色手工输入项。
// 使用指针是为了区分“未填写”和“显式填写 0”。
type ReportManualPayload struct {
	WorkInProgress               *float64 `json:"workInProgress"`
	FinishedGoods                *float64 `json:"finishedGoods"`
	RawMaterials                 *float64 `json:"rawMaterials"`
	IncomeTaxRate                *float64 `json:"incomeTaxRate"`
	EnterpriseCertificationScore *float64 `json:"enterpriseCertificationScore"`
	ProductionHumanScore         *float64 `json:"productionHumanScore"`
	ClosingSpeedScore            *float64 `json:"closingSpeedScore"`
}

// AllowedIncomeTaxRates 返回当前允许的税率列表。
func AllowedIncomeTaxRates() []float64 {
	result := make([]float64, len(allowedIncomeTaxRates))
	copy(result, allowedIncomeTaxRates)
	return result
}
