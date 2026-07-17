package payload

// OperatingQuarterMap 统一承接 Q1/Q2/Q3/Q4 的区块化输入。
type OperatingQuarterMap map[string]map[string]any

// OperatingPayload 对应经营页的区块化 JSON 结构。
type OperatingPayload struct {
	Beginning OperatingBeginningPayload `json:"beginning"`
	Quarter   OperatingQuarterPayload   `json:"quarter"`
	YearEnd   OperatingYearEndPayload   `json:"yearEnd"`
	Extra     OperatingExtraPayload     `json:"extra"`
	Derived   OperatingDerivedPayload   `json:"derived"`
}

// OperatingBeginningPayload 对应经营页年初区。
type OperatingBeginningPayload struct {
	TaxAndPlanning map[string]any   `json:"taxAndPlanning"`
	MarketBid      []map[string]any `json:"marketBid"`
}

// OperatingQuarterPayload 对应经营页季度区。
type OperatingQuarterPayload struct {
	ShortTermLoan          OperatingQuarterMap `json:"shortTermLoan"`
	MaterialPayment        OperatingQuarterMap `json:"materialPayment"`
	ProductionLineAdjust   OperatingQuarterMap `json:"productionLineAdjustment"`
	HumanResource          OperatingQuarterMap `json:"humanResource"`
	SalaryAndProduction    OperatingQuarterMap `json:"salaryAndProduction"`
	ResearchAndManagement  OperatingQuarterMap `json:"researchAndManagement"`
	SupplyChainOrderRecord OperatingQuarterMap `json:"supplyChainOrderRecord"`
	ReceivableUpdate       OperatingQuarterMap `json:"receivableUpdate"`
	DeliverySettlement     OperatingQuarterMap `json:"deliverySettlement"`
}

// OperatingYearEndPayload 对应经营页年末区。
type OperatingYearEndPayload struct {
	LongTermLoan    map[string]any `json:"longTermLoan"`
	AssetAdjustment map[string]any `json:"assetAdjustment"`
}

// OperatingExtraPayload 对应额外收入/罚款区。
type OperatingExtraPayload struct {
	IncomeAndPenalty OperatingQuarterMap `json:"incomeAndPenalty"`
}

// OperatingDerivedPayload 对应经营页只读派生区。
type OperatingDerivedPayload struct {
	Values map[string]any `json:"values"`
}

// NewOperatingPayload 提供统一的空载体，避免后续各层自行判空。
func NewOperatingPayload() OperatingPayload {
	return OperatingPayload{
		Beginning: OperatingBeginningPayload{
			TaxAndPlanning: map[string]any{},
			MarketBid:      []map[string]any{},
		},
		Quarter: OperatingQuarterPayload{
			ShortTermLoan:          OperatingQuarterMap{},
			MaterialPayment:        OperatingQuarterMap{},
			ProductionLineAdjust:   OperatingQuarterMap{},
			HumanResource:          OperatingQuarterMap{},
			SalaryAndProduction:    OperatingQuarterMap{},
			ResearchAndManagement:  OperatingQuarterMap{},
			SupplyChainOrderRecord: OperatingQuarterMap{},
			ReceivableUpdate:       OperatingQuarterMap{},
			DeliverySettlement:     OperatingQuarterMap{},
		},
		YearEnd: OperatingYearEndPayload{
			LongTermLoan:    map[string]any{},
			AssetAdjustment: map[string]any{},
		},
		Extra: OperatingExtraPayload{
			IncomeAndPenalty: OperatingQuarterMap{},
		},
		Derived: OperatingDerivedPayload{
			Values: map[string]any{},
		},
	}
}

// Normalize 会把空 map/slice 归一化为稳定结构，避免接口回包出现大量 null。
func (p OperatingPayload) Normalize() OperatingPayload {
	normalized := p

	if normalized.Beginning.TaxAndPlanning == nil {
		normalized.Beginning.TaxAndPlanning = map[string]any{}
	}
	if normalized.Beginning.MarketBid == nil {
		normalized.Beginning.MarketBid = []map[string]any{}
	}

	normalized.Quarter.ShortTermLoan = normalizeQuarterMap(normalized.Quarter.ShortTermLoan)
	normalized.Quarter.MaterialPayment = normalizeQuarterMap(normalized.Quarter.MaterialPayment)
	normalized.Quarter.ProductionLineAdjust = normalizeQuarterMap(normalized.Quarter.ProductionLineAdjust)
	normalized.Quarter.HumanResource = normalizeQuarterMap(normalized.Quarter.HumanResource)
	normalized.Quarter.SalaryAndProduction = normalizeQuarterMap(normalized.Quarter.SalaryAndProduction)
	normalized.Quarter.ResearchAndManagement = normalizeQuarterMap(normalized.Quarter.ResearchAndManagement)
	normalized.Quarter.SupplyChainOrderRecord = normalizeQuarterMap(normalized.Quarter.SupplyChainOrderRecord)
	normalized.Quarter.ReceivableUpdate = normalizeQuarterMap(normalized.Quarter.ReceivableUpdate)
	normalized.Quarter.DeliverySettlement = normalizeQuarterMap(normalized.Quarter.DeliverySettlement)

	if normalized.YearEnd.LongTermLoan == nil {
		normalized.YearEnd.LongTermLoan = map[string]any{}
	}
	if normalized.YearEnd.AssetAdjustment == nil {
		normalized.YearEnd.AssetAdjustment = map[string]any{}
	}

	normalized.Extra.IncomeAndPenalty = normalizeQuarterMap(normalized.Extra.IncomeAndPenalty)

	if normalized.Derived.Values == nil {
		normalized.Derived.Values = map[string]any{}
	}

	return normalized
}

func normalizeQuarterMap(value OperatingQuarterMap) OperatingQuarterMap {
	if value == nil {
		return OperatingQuarterMap{}
	}
	return value
}

// WithoutDerivedValues 返回移除派生区输入后的 payload，避免旧派生值反向污染实时计算。
func (p OperatingPayload) WithoutDerivedValues() OperatingPayload {
	normalized := p.Normalize()
	normalized.Derived.Values = map[string]any{}
	return normalized
}
