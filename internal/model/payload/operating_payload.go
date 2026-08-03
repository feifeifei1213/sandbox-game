package payload

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

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
	LongTermLoan               map[string]any                      `json:"longTermLoan"`
	AssetAdjustment            map[string]any                      `json:"assetAdjustment"`
	ProjectProgressUpdate      OperatingProjectProgressPayload     `json:"projectProgressUpdate"`
	MarketCultivation          OperatingMarketCultivationPayload   `json:"marketCultivation"`
	QualificationCertification OperatingQualificationCertification `json:"qualificationCertification"`
}

type OperatingProjectProgressPayload struct {
	Items []map[string]any `json:"items"`
}

type OperatingMarketCultivationPayload struct {
	Regional     OperatingMarketCultivationItem `json:"regional"`
	National     OperatingMarketCultivationItem `json:"national"`
	Global       OperatingMarketCultivationItem `json:"global"`
	AnnualTotal  float64                        `json:"annualTotal"`
	StateApplied bool                           `json:"stateApplied"`
}

type OperatingMarketCultivationItem struct {
	AnnualInvestment          any     `json:"annualInvestment"`
	EffectiveAnnualInvestment float64 `json:"effectiveAnnualInvestment"`
	PreviousCumulative        float64 `json:"previousCumulative"`
	CumulativeInvestment      float64 `json:"cumulativeInvestment"`
	Unlocked                  bool    `json:"unlocked"`
	Locked                    bool    `json:"locked"`
	LockedByPrevious          bool    `json:"lockedByPrevious"`
	WillUnlock                bool    `json:"willUnlock"`
}

type OperatingQualificationCertification struct {
	QualityEnvironmentalHealth OperatingQualificationItem `json:"qualityEnvironmentalHealth"`
	HighTechEnterprise         OperatingQualificationItem `json:"highTechEnterprise"`
	SpecializedInnovation      OperatingQualificationItem `json:"specializedInnovation"`
	ListedCompany              OperatingQualificationItem `json:"listedCompany"`
}

type OperatingQualificationItem struct {
	Status           any  `json:"status"`
	Unlocked         bool `json:"unlocked"`
	Locked           bool `json:"locked"`
	LockedByPrevious bool `json:"lockedByPrevious"`
}

type MarketCultivationCarryState struct {
	RegionalCumulative float64
	NationalCumulative float64
	GlobalCumulative   float64
	RegionalUnlocked   bool
	NationalUnlocked   bool
	GlobalUnlocked     bool
}

type QualificationCarryState struct {
	QualityEnvironmentalHealthUnlocked bool
	HighTechEnterpriseUnlocked         bool
	SpecializedInnovationUnlocked      bool
	ListedCompanyUnlocked              bool
}

type projectProgressFactoryDefinition struct {
	Key       string
	Label     string
	SlotCount int
}

type projectProgressQuarterDefinition struct {
	Key   string
	Label string
}

var projectProgressFactories = []projectProgressFactoryDefinition{
	{Key: "factoryA", Label: "生产厂房 A", SlotCount: 4},
	{Key: "factoryB", Label: "生产厂房 B", SlotCount: 3},
	{Key: "factoryC", Label: "生产厂房 C", SlotCount: 1},
}

var projectProgressQuarters = []projectProgressQuarterDefinition{
	{Key: "q1", Label: "第一季度"},
	{Key: "q2", Label: "第二季度"},
	{Key: "q3", Label: "第三季度"},
	{Key: "q4", Label: "第四季度"},
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
			LongTermLoan:               map[string]any{},
			AssetAdjustment:            map[string]any{},
			ProjectProgressUpdate:      NewOperatingProjectProgressPayload(),
			MarketCultivation:          NewOperatingMarketCultivationPayload(),
			QualificationCertification: NewOperatingQualificationCertification(),
		},
		Extra: OperatingExtraPayload{
			IncomeAndPenalty: OperatingQuarterMap{},
		},
		Derived: OperatingDerivedPayload{
			Values: map[string]any{},
		},
	}
}

func NewOperatingProjectProgressPayload() OperatingProjectProgressPayload {
	return OperatingProjectProgressPayload{Items: newOperatingProjectProgressItems()}
}

func NewOperatingMarketCultivationPayload() OperatingMarketCultivationPayload {
	return OperatingMarketCultivationPayload{
		Regional: newOperatingMarketCultivationItem(),
		National: newOperatingMarketCultivationItem(),
		Global:   newOperatingMarketCultivationItem(),
	}
}

func newOperatingMarketCultivationItem() OperatingMarketCultivationItem {
	return OperatingMarketCultivationItem{AnnualInvestment: ""}
}

func NewOperatingQualificationCertification() OperatingQualificationCertification {
	return OperatingQualificationCertification{
		QualityEnvironmentalHealth: newOperatingQualificationItem(),
		HighTechEnterprise:         newOperatingQualificationItem(),
		SpecializedInnovation:      newOperatingQualificationItem(),
		ListedCompany:              newOperatingQualificationItem(),
	}
}

func newOperatingQualificationItem() OperatingQualificationItem {
	return OperatingQualificationItem{Status: "未解锁"}
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
	normalized.YearEnd.ProjectProgressUpdate = normalizeProjectProgressPayload(normalized.YearEnd.ProjectProgressUpdate)
	normalized.YearEnd.MarketCultivation = normalizeMarketCultivationPayload(normalized.YearEnd.MarketCultivation)
	normalized.YearEnd.QualificationCertification = normalizeQualificationCertification(normalized.YearEnd.QualificationCertification)

	normalized.Extra.IncomeAndPenalty = normalizeQuarterMap(normalized.Extra.IncomeAndPenalty)

	if normalized.Derived.Values == nil {
		normalized.Derived.Values = map[string]any{}
	}

	return normalized
}

func normalizeProjectProgressPayload(value OperatingProjectProgressPayload) OperatingProjectProgressPayload {
	if value.Items == nil {
		return NewOperatingProjectProgressPayload()
	}
	fallback := newOperatingProjectProgressItems()
	keyedItems := make(map[string]map[string]any)
	for _, item := range value.Items {
		factoryKey := projectProgressString(item["factoryKey"])
		slotNo, slotOK := parseWholeNumber(item["slotNo"])
		quarterKey := projectProgressString(item["quarterKey"])
		if factoryKey != "" && slotOK && quarterKey != "" {
			keyedItems[projectProgressRecordKey(factoryKey, slotNo, quarterKey)] = item
		}
	}

	legacyItems := make([]map[string]any, 0)
	for _, item := range value.Items {
		if projectProgressString(item["factoryKey"]) == "" && item["slotNo"] == nil && projectProgressString(item["quarterKey"]) == "" {
			legacyItems = append(legacyItems, item)
		}
	}

	normalized := make([]map[string]any, len(fallback))
	for index, item := range fallback {
		source := keyedItems[projectProgressRecordKey(toString(item["factoryKey"]), intNumber(item["slotNo"]), toString(item["quarterKey"]))]
		if source == nil && intNumber(item["slotNo"]) == 1 {
			factoryIndex := projectProgressFactoryIndex(toString(item["factoryKey"]))
			quarterIndex := projectProgressQuarterIndex(toString(item["quarterKey"]))
			legacyIndex := factoryIndex*len(projectProgressQuarters) + quarterIndex
			if factoryIndex >= 0 && quarterIndex >= 0 && legacyIndex >= 0 && legacyIndex < len(legacyItems) {
				source = legacyItems[legacyIndex]
			}
		}
		if source == nil && len(value.Items) >= len(fallback) && index < len(value.Items) {
			source = value.Items[index]
		}
		if source != nil {
			item["lineType"] = source["lineType"]
			item["progress"] = source["progress"]
		}
		normalized[index] = item
	}

	return OperatingProjectProgressPayload{Items: normalized}
}

func newOperatingProjectProgressItems() []map[string]any {
	items := make([]map[string]any, 0, 32)
	for _, factory := range projectProgressFactories {
		for slotNo := 1; slotNo <= factory.SlotCount; slotNo++ {
			for _, quarter := range projectProgressQuarters {
				items = append(items, map[string]any{
					"projectName":  fmt.Sprintf("%s-槽位%d-%s", factory.Label, slotNo, quarter.Label),
					"factoryKey":   factory.Key,
					"factoryLabel": factory.Label,
					"slotNo":       slotNo,
					"quarterKey":   quarter.Key,
					"quarterLabel": quarter.Label,
					"lineType":     "",
					"progress":     "",
				})
			}
		}
	}
	return items
}

func projectProgressRecordKey(factoryKey string, slotNo int, quarterKey string) string {
	return fmt.Sprintf("%s:%d:%s", factoryKey, slotNo, quarterKey)
}

func projectProgressString(value any) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(toString(value))
}

func projectProgressFactoryIndex(factoryKey string) int {
	for index, factory := range projectProgressFactories {
		if factory.Key == factoryKey {
			return index
		}
	}
	return -1
}

func projectProgressQuarterIndex(quarterKey string) int {
	for index, quarter := range projectProgressQuarters {
		if quarter.Key == quarterKey {
			return index
		}
	}
	return -1
}

func normalizeMarketCultivationPayload(value OperatingMarketCultivationPayload) OperatingMarketCultivationPayload {
	value.Regional = normalizeMarketCultivationItem(value.Regional)
	value.National = normalizeMarketCultivationItem(value.National)
	value.Global = normalizeMarketCultivationItem(value.Global)
	if value.StateApplied && hasMarketCultivationPayloadSignal(value) {
		value.AnnualTotal = value.Regional.EffectiveAnnualInvestment + value.National.EffectiveAnnualInvestment + value.Global.EffectiveAnnualInvestment
	}
	return value
}

func normalizeMarketCultivationItem(value OperatingMarketCultivationItem) OperatingMarketCultivationItem {
	if value.AnnualInvestment == nil {
		value.AnnualInvestment = ""
	}
	return value
}

func normalizeQualificationCertification(value OperatingQualificationCertification) OperatingQualificationCertification {
	value.QualityEnvironmentalHealth = normalizeQualificationItem(value.QualityEnvironmentalHealth)
	value.HighTechEnterprise = normalizeQualificationItem(value.HighTechEnterprise)
	value.SpecializedInnovation = normalizeQualificationItem(value.SpecializedInnovation)
	value.ListedCompany = normalizeQualificationItem(value.ListedCompany)
	return value
}

func normalizeQualificationItem(value OperatingQualificationItem) OperatingQualificationItem {
	if value.Status == nil || strings.TrimSpace(toString(value.Status)) == "" {
		value.Status = "未解锁"
	}
	return value
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

func MarketCultivationAnnualAmount(value OperatingMarketCultivationPayload) (float64, bool) {
	value = normalizeMarketCultivationPayload(value)
	if value.StateApplied {
		if !hasMarketCultivationPayloadSignal(value) {
			return 0, false
		}
		return value.Regional.EffectiveAnnualInvestment +
			value.National.EffectiveAnnualInvestment +
			value.Global.EffectiveAnnualInvestment, true
	}

	total := 0.0
	found := false
	for _, item := range []OperatingMarketCultivationItem{value.Regional, value.National, value.Global} {
		number, ok := ParseManualNumber(item.AnnualInvestment)
		if !ok {
			continue
		}
		total += number
		found = true
	}
	return total, found
}

func ExtractMarketCultivationCarryState(source OperatingPayload) MarketCultivationCarryState {
	normalized := source.Normalize()
	market := normalized.YearEnd.MarketCultivation
	regionalCumulative := carryCumulative(market.Regional, 1)
	nationalCumulative := carryCumulative(market.National, 2)
	globalCumulative := carryCumulative(market.Global, 3)
	return MarketCultivationCarryState{
		RegionalCumulative: regionalCumulative,
		NationalCumulative: nationalCumulative,
		GlobalCumulative:   globalCumulative,
		RegionalUnlocked:   market.Regional.Unlocked || regionalCumulative >= 1,
		NationalUnlocked:   market.National.Unlocked || nationalCumulative >= 2,
		GlobalUnlocked:     market.Global.Unlocked || globalCumulative >= 3,
	}
}

func carryCumulative(item OperatingMarketCultivationItem, threshold float64) float64 {
	if item.CumulativeInvestment > 0 || item.Unlocked || item.Locked || item.LockedByPrevious {
		if item.CumulativeInvestment < threshold && item.Unlocked {
			return threshold
		}
		return item.CumulativeInvestment
	}
	number, ok := ParseManualNumber(item.AnnualInvestment)
	if !ok {
		return 0
	}
	return number
}

func ApplyMarketCultivationState(
	current OperatingMarketCultivationPayload,
	previous MarketCultivationCarryState,
	currentYearEndEffective bool,
	preserveLockedAnnual bool,
) OperatingMarketCultivationPayload {
	current = normalizeMarketCultivationPayload(current)
	stateApplied := hasMarketCultivationPayloadSignal(current) || hasMarketCultivationCarryState(previous)
	current.Regional = applyMarketCultivationItemState(current.Regional, previous.RegionalCumulative, previous.RegionalUnlocked, 1, currentYearEndEffective, preserveLockedAnnual)
	current.National = applyMarketCultivationItemState(current.National, previous.NationalCumulative, previous.NationalUnlocked, 2, currentYearEndEffective, preserveLockedAnnual)
	current.Global = applyMarketCultivationItemState(current.Global, previous.GlobalCumulative, previous.GlobalUnlocked, 3, currentYearEndEffective, preserveLockedAnnual)
	if stateApplied {
		current.AnnualTotal = current.Regional.EffectiveAnnualInvestment + current.National.EffectiveAnnualInvestment + current.Global.EffectiveAnnualInvestment
	} else {
		current.AnnualTotal = 0
	}
	current.StateApplied = stateApplied
	return current
}

func hasMarketCultivationCarryState(previous MarketCultivationCarryState) bool {
	return previous.RegionalCumulative != 0 ||
		previous.NationalCumulative != 0 ||
		previous.GlobalCumulative != 0 ||
		previous.RegionalUnlocked ||
		previous.NationalUnlocked ||
		previous.GlobalUnlocked
}

func hasMarketCultivationPayloadSignal(value OperatingMarketCultivationPayload) bool {
	return hasMarketCultivationItemSignal(value.Regional) ||
		hasMarketCultivationItemSignal(value.National) ||
		hasMarketCultivationItemSignal(value.Global)
}

func hasMarketCultivationItemSignal(item OperatingMarketCultivationItem) bool {
	return hasExplicitManualValue(item.AnnualInvestment) ||
		item.EffectiveAnnualInvestment != 0 ||
		item.PreviousCumulative != 0 ||
		item.CumulativeInvestment != 0 ||
		item.Unlocked ||
		item.Locked ||
		item.LockedByPrevious ||
		item.WillUnlock
}

func hasExplicitManualValue(value any) bool {
	return value != nil && strings.TrimSpace(toString(value)) != ""
}

func applyMarketCultivationItemState(
	current OperatingMarketCultivationItem,
	previousCumulative float64,
	previousUnlocked bool,
	threshold float64,
	currentYearEndEffective bool,
	preserveLockedAnnual bool,
) OperatingMarketCultivationItem {
	current = normalizeMarketCultivationItem(current)
	annualInvestment, annualFound := ParseManualNumber(current.AnnualInvestment)
	if !annualFound {
		annualInvestment = 0
	}

	if previousUnlocked {
		if !preserveLockedAnnual {
			current.AnnualInvestment = 0
		}
		if previousCumulative < threshold {
			previousCumulative = threshold
		}
		current.EffectiveAnnualInvestment = 0
		current.PreviousCumulative = previousCumulative
		current.CumulativeInvestment = previousCumulative
		current.Unlocked = true
		current.Locked = true
		current.LockedByPrevious = true
		current.WillUnlock = false
		return current
	}

	cumulative := previousCumulative + annualInvestment
	current.EffectiveAnnualInvestment = annualInvestment
	current.PreviousCumulative = previousCumulative
	current.CumulativeInvestment = cumulative
	current.Unlocked = currentYearEndEffective && cumulative >= threshold
	current.Locked = current.Unlocked
	current.LockedByPrevious = false
	current.WillUnlock = annualInvestment > 0 && cumulative >= threshold && !currentYearEndEffective
	return current
}

func ExtractQualificationCarryState(source OperatingPayload) QualificationCarryState {
	qualification := source.Normalize().YearEnd.QualificationCertification
	return QualificationCarryState{
		QualityEnvironmentalHealthUnlocked: qualification.QualityEnvironmentalHealth.Unlocked || IsQualificationUnlocked(qualification.QualityEnvironmentalHealth.Status),
		HighTechEnterpriseUnlocked:         qualification.HighTechEnterprise.Unlocked || IsQualificationUnlocked(qualification.HighTechEnterprise.Status),
		SpecializedInnovationUnlocked:      qualification.SpecializedInnovation.Unlocked || IsQualificationUnlocked(qualification.SpecializedInnovation.Status),
		ListedCompanyUnlocked:              qualification.ListedCompany.Unlocked || IsQualificationUnlocked(qualification.ListedCompany.Status),
	}
}

func ApplyQualificationState(current OperatingQualificationCertification, previous QualificationCarryState, currentYearEndEffective bool) OperatingQualificationCertification {
	current = normalizeQualificationCertification(current)
	current.QualityEnvironmentalHealth = applyQualificationItemState(current.QualityEnvironmentalHealth, previous.QualityEnvironmentalHealthUnlocked, currentYearEndEffective)
	current.HighTechEnterprise = applyQualificationItemState(current.HighTechEnterprise, previous.HighTechEnterpriseUnlocked, currentYearEndEffective)
	current.SpecializedInnovation = applyQualificationItemState(current.SpecializedInnovation, previous.SpecializedInnovationUnlocked, currentYearEndEffective)
	current.ListedCompany = applyQualificationItemState(current.ListedCompany, previous.ListedCompanyUnlocked, currentYearEndEffective)
	return current
}

func applyQualificationItemState(current OperatingQualificationItem, previousUnlocked bool, currentYearEndEffective bool) OperatingQualificationItem {
	current = normalizeQualificationItem(current)
	if previousUnlocked {
		current.Status = "解锁"
		current.Unlocked = true
		current.Locked = true
		current.LockedByPrevious = true
		return current
	}
	selectedUnlocked := IsQualificationUnlocked(current.Status)
	current.Unlocked = selectedUnlocked
	current.Locked = currentYearEndEffective && selectedUnlocked
	current.LockedByPrevious = false
	return current
}

func IsQualificationUnlocked(value any) bool {
	normalized := strings.ToLower(strings.TrimSpace(toString(value)))
	switch normalized {
	case "解锁", "unlocked", "unlock", "true", "1":
		return true
	default:
		return false
	}
}

func ParseManualNumber(value any) (float64, bool) {
	switch typed := value.(type) {
	case nil:
		return 0, false
	case int:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case float32:
		return float64(typed), true
	case float64:
		return typed, true
	case json.Number:
		parsed, err := strconv.ParseFloat(typed.String(), 64)
		return parsed, err == nil
	case string:
		trimmed := strings.TrimSpace(typed)
		if trimmed == "" {
			return 0, false
		}
		parsed, err := strconv.ParseFloat(trimmed, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}

func parseWholeNumber(value any) (int, bool) {
	number, ok := ParseManualNumber(value)
	if !ok {
		return 0, false
	}
	whole := int(number)
	if number != float64(whole) {
		return 0, false
	}
	return whole, true
}

func intNumber(value any) int {
	number, _ := parseWholeNumber(value)
	return number
}

func toString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case json.Number:
		return typed.String()
	default:
		return fmt.Sprint(value)
	}
}
