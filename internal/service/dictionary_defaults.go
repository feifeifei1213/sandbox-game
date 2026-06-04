package service

import (
	"encoding/json"
	"fmt"
)

const (
	DictionaryCategoryMarket    = "MARKET"
	DictionaryCategoryOrderType = "ORDER_TYPE"
	DictionaryCategoryOperating = "OPERATING"
	DictionaryCategoryReport    = "REPORT"
	DictionaryCategoryBaseline  = "BASELINE"
)

type DictionaryDefinition struct {
	ItemCode       string
	ItemCategory   string
	DefaultName    string
	DisplayOrder   int
	Editable       bool
	RelatedPayload map[string]any
}

func BuiltInDictionaryDefinitions(editionCode string) []DictionaryDefinition {
	edition := ResolveGameEdition(editionCode)
	items := make([]DictionaryDefinition, 0, 80)
	items = append(items, builtInCommonOrderDictionaryDefinitions()...)
	items = append(items, builtInOperatingDictionaryDefinitions(edition.EditionCode)...)
	items = append(items, builtInReportDictionaryDefinitions(edition.EditionCode)...)
	items = append(items, builtInBaselineDictionaryDefinitions(edition.EditionCode)...)
	return items
}

func builtInCommonOrderDictionaryDefinitions() []DictionaryDefinition {
	return []DictionaryDefinition{
		dictDef("market.LOCAL", DictionaryCategoryMarket, "本地市场", 10, "LOCAL"),
		dictDef("market.REGIONAL", DictionaryCategoryMarket, "区域市场", 20, "REGIONAL"),
		dictDef("market.NATIONAL", DictionaryCategoryMarket, "全国市场", 30, "NATIONAL"),
		dictDef("market.GLOBAL", DictionaryCategoryMarket, "全球市场", 40, "GLOBAL"),
		dictDef("orderType.AGENCY_INSPECTION", DictionaryCategoryOrderType, "代办过检", 100, "AGENCY_INSPECTION"),
		dictDef("orderType.TWO_CABIN_VIP", DictionaryCategoryOrderType, "两舱贵宾", 110, "TWO_CABIN_VIP"),
		dictDef("orderType.BUSINESS_VIP", DictionaryCategoryOrderType, "商务贵宾", 120, "BUSINESS_VIP"),
		dictDef("orderType.MEMBER_CUSTOM", DictionaryCategoryOrderType, "会员定制", 130, "MEMBER_CUSTOM"),
	}
}

func builtInOperatingDictionaryDefinitions(editionCode string) []DictionaryDefinition {
	if editionCode == GameEditionProductionV1 {
		return []DictionaryDefinition{
			opDef("taxPolicyNote", "一般企业 25%，高新企业 15%", 200),
			opDef("materialGroupTitle", "2. 支付材料费用", 210),
			opDef("productionAdjustmentGroupTitle", "3. 生产线调整", 220),
			opDef("humanResourceGroupTitle", "4. 人力资源", 230),
			opDef("humanResourceLabel", "人力资源费用", 240),
			opDef("salaryGroupTitle", "5. 工资与生产", 250),
			opDef("salaryLabel", "工资与生产费用", 260),
			opDef("researchAndManagementGroupTitle", "6. 研发与管理", 270),
			opDef("newOrderReminder", "下新供应链订单", 280),
			opDef("receivableGroupTitle", "7. 应收更新", 290),
			opDef("receivableLabel", "应收回款", 300),
			opDef("receivableReminder", "请将应收账款向左移动一格", 310),
			opDef("deliveryGroupTitle", "8. 产品验收交付", 320),
			opDef("deliveryRevenueLabel", "销售收入", 330),
			opDef("deliveryCostLabel", "直接成本", 340),
			opDef("managementStaffGroupTitle", "9. 支付管理人员费用", 350),
			opDef("lineMaintenanceGroupTitle", "2. 支付生产线年度维护费", 360),
			opDef("lineMaintenanceNote", "1M / 条", 370),
			opDef("assetGroupTitle", "3. 生产厂房", 380),
			opDef("assetValueNote", "A / B / C 厂房价值（40M、32M、16M）", 390),
			opDef("rentGroupTitle", "4. 生产厂房租金", 400),
			opDef("rentValueNote", "A / B / C 厂房租金（4M、3M、2M）", 410),
			opDef("residualGroupTitle", "5. 生产线残值", 420),
			opDef("depreciationGroupTitle", "6. 生产线折旧", 430),
			opDef("workInConstructionLabel", "未完工生产线价值", 440),
			opArrayDef("marketProductFields", "agencyInspectionTotal", "代办过检服务总价", 450),
			opArrayDef("marketProductFields", "twoCabinVipTotal", "两舱贵宾服务总价", 460),
			opArrayDef("marketProductFields", "businessVipTotal", "商务贵宾服务总价", 470),
			opArrayDef("marketProductFields", "memberCustomTotal", "会员定制服务总价", 480),
			opArrayDef("materialFields", "basicProduct", "基础产品", 490),
			opArrayDef("materialFields", "standardProduct", "标准产品", 500),
			opArrayDef("materialFields", "precisionProduct", "精密产品", 510),
			opArrayDef("materialFields", "intelligentProduct", "智能产品", 520),
			opArrayDef("productionLineRows", "changeProduct", "变更产品", 530),
			opArrayDef("productionLineRows", "dismantleCost", "生产线拆除 1M / 条", 540),
			opArrayDef("productionLineRows", "lineSale", "生产线出售", 550),
			opArrayDef("productionLineRows", "newLineInstall", "新生产线安装", 560),
			opArrayDef("productionLineRows", "constructionToFixed", "生产线残值", 570),
			opArrayGroupDef("productionLineRows", "constructionToFixed", "建成生产线转固", 575),
			opArrayDef("productionLineRows", "newDepreciableAsset", "待折资产", 580),
			opArrayDef("researchFields", "technologyResearch", "技术研发", 590),
			opArrayDef("researchFields", "managementSystem", "管理体系", 600),
			opDef("progressUpdateReminder", "产品进度更新", 610),
		}
	}

	return []DictionaryDefinition{
		opDef("taxPolicyNote", "一般企业 25%，支线机场贵宾 15%", 200),
		opDef("materialGroupTitle", "2. 支付前期花费费用", 210),
		opDef("productionAdjustmentGroupTitle", "3. 贵宾厅调整", 220),
		opDef("humanResourceGroupTitle", "5. 人力资源管理", 230),
		opDef("humanResourceLabel", "支付招聘、辞退、培训、待岗费用", 240),
		opDef("salaryGroupTitle", "6. 贵宾厅空位开始服务，发放服务人员工资", 250),
		opDef("salaryLabel", "支付贵宾厅的服务人员工资", 260),
		opDef("researchAndManagementGroupTitle", "7. 新服务研创 / 8. 管理体系投入", 270),
		opDef("newOrderReminder", "9. 下新服务订单", 280),
		opDef("receivableGroupTitle", "10. 更新应收账款", 290),
		opDef("receivableLabel", "如应收账款进入现金区做记录", 300),
		opDef("receivableReminder", "所有应收向左移动一格", 310),
		opDef("deliveryGroupTitle", "11. 服务验收交付", 320),
		opDef("deliveryRevenueLabel", "服务订单销售额", 330),
		opDef("deliveryCostLabel", "服务订单成本", 340),
		opDef("managementStaffGroupTitle", "12. 支付管理人员费用", 350),
		opDef("lineMaintenanceGroupTitle", "2. 支付贵宾厅年度维护费", 360),
		opDef("lineMaintenanceNote", "1M / 间", 370),
		opDef("assetGroupTitle", "3. 贵宾区资产", 380),
		opDef("assetValueNote", "贵宾A / B / C区 价值（40M、32M、16M）", 390),
		opDef("rentGroupTitle", "4. 贵宾区租金", 400),
		opDef("rentValueNote", "贵宾A / B / C区 租金（4M、3M、2M）", 410),
		opDef("residualGroupTitle", "5. 贵宾厅残值", 420),
		opDef("depreciationGroupTitle", "6. 贵宾厅折旧", 430),
		opDef("workInConstructionLabel", "未完工贵宾厅价值", 440),
		opArrayDef("marketProductFields", "agencyInspectionTotal", "代办过检服务总价", 450),
		opArrayDef("marketProductFields", "twoCabinVipTotal", "两舱贵宾服务总价", 460),
		opArrayDef("marketProductFields", "businessVipTotal", "商务贵宾服务总价", 470),
		opArrayDef("marketProductFields", "memberCustomTotal", "会员定制服务总价", 480),
		opArrayDef("materialFields", "basicProduct", "代办过检服务", 490),
		opArrayDef("materialFields", "standardProduct", "两舱贵宾服务", 500),
		opArrayDef("materialFields", "precisionProduct", "商务贵宾服务", 510),
		opArrayDef("materialFields", "intelligentProduct", "会员定制服务", 520),
		opArrayDef("productionLineRows", "changeProduct", "变更服务", 530),
		opArrayDef("productionLineRows", "dismantleCost", "贵宾厅拆除 1M / 间", 540),
		opArrayDef("productionLineRows", "lineSale", "贵宾厅出售", 550),
		opArrayDef("productionLineRows", "newLineInstall", "新贵宾厅建造", 560),
		opArrayDef("productionLineRows", "constructionToFixed", "贵宾厅残值", 570),
		opArrayGroupDef("productionLineRows", "constructionToFixed", "建成贵宾厅转固", 575),
		opArrayDef("productionLineRows", "newDepreciableAsset", "待折资产", 580),
		opArrayDef("researchFields", "technologyResearch", "新服务研创", 590),
		opArrayDef("researchFields", "managementSystem", "质量环境健康费用", 600),
		opDef("progressUpdateReminder", "服务进度更新", 610),
	}
}

func builtInReportDictionaryDefinitions(editionCode string) []DictionaryDefinition {
	if editionCode == GameEditionProductionV1 {
		return []DictionaryDefinition{
			reportDef("workInConstruction", "在建生产线", 700),
			reportDef("factoryAsset", "厂房", 710),
			reportDef("lineResidual", "生产线残值", 720),
			reportDef("workInProgress", "在制品", 730),
			reportDef("finishedGoods", "成品", 740),
			reportDef("rawMaterials", "材料", 750),
			reportDef("row19Label", "税前利润", 760),
			reportDef("totalLiabilityEquity", "总负债和权益", 770),
			reportDef("bestProductionHumanDirector", "最佳生产人力总监得分", 780),
			reportDef("taxRateEnterprise15", "0.15 高新企业", 790),
			reportRequiredDef("workInProgress", "在制品", 800),
			reportRequiredDef("finishedGoods", "成品", 810),
			reportRequiredDef("rawMaterials", "材料", 820),
			reportRequiredDef("incomeTaxRate", "所得税税率", 830),
			reportRequiredDef("enterpriseCertificationScore", "企业认证得分", 840),
			reportRequiredDef("productionHumanScore", "最佳生产人力总监得分", 850),
			reportRequiredDef("closingSpeedScore", "关账速度得分", 860),
		}
	}
	return []DictionaryDefinition{
		reportDef("workInConstruction", "在建贵宾厅", 700),
		reportDef("factoryAsset", "贵宾区", 710),
		reportDef("lineResidual", "贵宾厅残值", 720),
		reportDef("workInProgress", "在途服务", 730),
		reportDef("finishedGoods", "服务验收", 740),
		reportDef("rawMaterials", "前期花费", 750),
		reportDef("row19Label", "3.服务进度更新", 760),
		reportDef("totalLiabilityEquity", "总负债权益", 770),
		reportDef("bestProductionHumanDirector", "最佳服务人力总监得分", 780),
		reportDef("taxRateEnterprise15", "0.15 支线机场贵宾", 790),
		reportRequiredDef("workInProgress", "在途服务", 800),
		reportRequiredDef("finishedGoods", "服务验收", 810),
		reportRequiredDef("rawMaterials", "前期花费", 820),
		reportRequiredDef("incomeTaxRate", "所得税税率", 830),
		reportRequiredDef("enterpriseCertificationScore", "企业认证得分", 840),
		reportRequiredDef("productionHumanScore", "最佳服务人力总监得分", 850),
		reportRequiredDef("closingSpeedScore", "关账速度得分", 860),
	}
}

func builtInBaselineDictionaryDefinitions(editionCode string) []DictionaryDefinition {
	commonDefinitions := []DictionaryDefinition{
		baselineDef("salesRevenue", "销售收入", 900),
		baselineDef("directCost", "直接成本", 910),
		baselineDef("comprehensiveCost", "综合费用", 920),
		baselineDef("depreciation", "折旧", 930),
		baselineDef("financeIncomeExpense", "财务收入/支出", 940),
		baselineDef("extraIncomeExpense", "额外收入/支出", 950),
		baselineDef("incomeTax", "所得税", 960),
		baselineDef("depreciableAsset", "待折资产", 1000),
		baselineDef("cash", "现金", 1010),
		baselineDef("receivable", "应收款", 1020),
		baselineDef("shortTermLoan", "短期负债", 1060),
		baselineDef("longTermLoan", "长期负债", 1070),
		baselineDef("shareCapital", "股东资本", 1080),
		baselineDef("retainedEarnings", "利润留存", 1090),
	}
	if editionCode == GameEditionProductionV1 {
		return append(commonDefinitions,
			baselineDef("workInConstruction", "在建生产线", 970),
			baselineDef("factoryAsset", "厂房", 980),
			baselineDef("lineResidual", "生产线残值", 990),
			baselineDef("workInProgress", "在制品", 1030),
			baselineDef("finishedGoods", "成品", 1040),
			baselineDef("rawMaterials", "材料", 1050),
		)
	}
	return append(commonDefinitions,
		baselineDef("workInConstruction", "在建贵宾厅", 970),
		baselineDef("factoryAsset", "贵宾区", 980),
		baselineDef("lineResidual", "贵宾厅残值", 990),
		baselineDef("workInProgress", "在途服务", 1030),
		baselineDef("finishedGoods", "服务验收", 1040),
		baselineDef("rawMaterials", "前期花费", 1050),
	)
}

func dictDef(code string, category string, defaultName string, order int, businessCode string) DictionaryDefinition {
	return DictionaryDefinition{
		ItemCode:     code,
		ItemCategory: category,
		DefaultName:  defaultName,
		DisplayOrder: order,
		Editable:     true,
		RelatedPayload: map[string]any{
			"businessCode": businessCode,
		},
	}
}

func opDef(key string, defaultName string, order int) DictionaryDefinition {
	return DictionaryDefinition{
		ItemCode:     fmt.Sprintf("operating.%s", key),
		ItemCategory: DictionaryCategoryOperating,
		DefaultName:  defaultName,
		DisplayOrder: order,
		Editable:     true,
		RelatedPayload: map[string]any{
			"labelKey": key,
		},
	}
}

func opArrayDef(arrayName string, key string, defaultName string, order int) DictionaryDefinition {
	return DictionaryDefinition{
		ItemCode:     fmt.Sprintf("operating.%s.%s", arrayName, key),
		ItemCategory: DictionaryCategoryOperating,
		DefaultName:  defaultName,
		DisplayOrder: order,
		Editable:     true,
		RelatedPayload: map[string]any{
			"arrayName": arrayName,
			"rowKey":    key,
			"field":     "label",
		},
	}
}

func opArrayGroupDef(arrayName string, key string, defaultName string, order int) DictionaryDefinition {
	return DictionaryDefinition{
		ItemCode:     fmt.Sprintf("operating.%s.%s.groupLabel", arrayName, key),
		ItemCategory: DictionaryCategoryOperating,
		DefaultName:  defaultName,
		DisplayOrder: order,
		Editable:     true,
		RelatedPayload: map[string]any{
			"arrayName": arrayName,
			"rowKey":    key,
			"field":     "groupLabel",
		},
	}
}

func reportDef(key string, defaultName string, order int) DictionaryDefinition {
	return DictionaryDefinition{
		ItemCode:     fmt.Sprintf("report.%s", key),
		ItemCategory: DictionaryCategoryReport,
		DefaultName:  defaultName,
		DisplayOrder: order,
		Editable:     true,
		RelatedPayload: map[string]any{
			"labelKey": key,
		},
	}
}

func reportRequiredDef(key string, defaultName string, order int) DictionaryDefinition {
	return DictionaryDefinition{
		ItemCode:     fmt.Sprintf("reportRequired.%s", key),
		ItemCategory: DictionaryCategoryReport,
		DefaultName:  defaultName,
		DisplayOrder: order,
		Editable:     true,
		RelatedPayload: map[string]any{
			"labelKey": key,
		},
	}
}

func baselineDef(key string, defaultName string, order int) DictionaryDefinition {
	return DictionaryDefinition{
		ItemCode:     fmt.Sprintf("baseline.%s", key),
		ItemCategory: DictionaryCategoryBaseline,
		DefaultName:  defaultName,
		DisplayOrder: order,
		Editable:     true,
		RelatedPayload: map[string]any{
			"labelKey": key,
		},
	}
}

func marshalDictionaryRelatedPayload(value map[string]any) []byte {
	if len(value) == 0 {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return raw
}
