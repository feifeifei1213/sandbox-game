package service

import (
	"errors"
	"strings"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
)

const (
	OrderTemplateVersionVIPServiceV1 = "VIP_ORDER_TEMPLATE_V1"
	OrderTemplateVersionAirportV1    = "AIRPORT_ORDER_TEMPLATE_V1"

	OrderFormulaVersionVIPServiceV1 = "ORDER_GEN_SERVICE_V1"
	OrderFormulaVersionAirportV1    = "AIRPORT_ORDER_FORMULA_V1"

	MarketCodeDomestic      = "DOMESTIC"
	MarketCodeInternational = "INTERNATIONAL"

	OrderTypeNarrowBody = "NARROW_BODY"
	OrderTypeWideBody   = "WIDE_BODY"
)

var ErrOrderTemplateUnsupported = errors.New("order template unsupported")

type OrderTemplateDefinition struct {
	TemplateVersion string
	TemplateName    string
	FormulaVersion  string
	MaxCardCount    int
	DeliveryEnabled bool
	Markets         []OrderMarketDefinition
	OrderTypes      []OrderTypeDefinition
	Fields          []OrderFieldDefinition
}

type OrderMarketDefinition struct {
	Code           string
	Name           string
	DefaultEnabled bool
	SortOrder      int
}

type OrderTypeDefinition struct {
	Code      string
	Name      string
	SortOrder int
}

type OrderFieldDefinition struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	ValueType string `json:"valueType"`
	Unit      string `json:"unit,omitempty"`
	Order     int    `json:"order"`
}

type OrderTemplateView struct {
	TemplateVersion string                 `json:"templateVersion"`
	TemplateName    string                 `json:"templateName"`
	FormulaVersion  string                 `json:"formulaVersion"`
	SegmentCount    int                    `json:"segmentCount"`
	MaxCardCount    int                    `json:"maxCardCount"`
	DeliveryEnabled bool                   `json:"deliveryEnabled"`
	Markets         []OrderMarketView      `json:"markets"`
	OrderTypes      []OrderTypeView        `json:"orderTypes"`
	Fields          []OrderFieldDefinition `json:"fields"`
}

type OrderMarketView struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	DefaultEnabled bool   `json:"defaultEnabled"`
	SortOrder      int    `json:"sortOrder"`
}

type OrderTypeView struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	SortOrder int    `json:"sortOrder"`
}

func ResolveOrderTemplateByGameConfig(config entity.GameConfig) (OrderTemplateDefinition, error) {
	edition := normalizeGameConfigEdition(config)
	return ResolveOrderTemplateVersion(edition.OrderTemplateVersion)
}

func ResolveOrderTemplateVersion(version string) (OrderTemplateDefinition, error) {
	normalized := strings.ToUpper(strings.TrimSpace(version))
	if normalized == "" {
		normalized = OrderTemplateVersionVIPServiceV1
	}
	for _, template := range builtInOrderTemplates() {
		if template.TemplateVersion == normalized {
			return template, nil
		}
	}
	return OrderTemplateDefinition{}, ErrOrderTemplateUnsupported
}

func MustResolveOrderTemplateVersion(version string) OrderTemplateDefinition {
	template, err := ResolveOrderTemplateVersion(version)
	if err != nil {
		return builtInOrderTemplates()[0]
	}
	return template
}

func (t OrderTemplateDefinition) Segments() []OrderSegmentDefinition {
	segments := make([]OrderSegmentDefinition, 0, len(t.Markets)*len(t.OrderTypes))
	for _, market := range t.Markets {
		for _, orderType := range t.OrderTypes {
			segments = append(segments, OrderSegmentDefinition{
				MarketCode:    market.Code,
				MarketName:    market.Name,
				OrderType:     orderType.Code,
				OrderTypeName: orderType.Name,
			})
		}
	}
	return segments
}

func (t OrderTemplateDefinition) IsValidMarketCode(code string) bool {
	_, ok := t.MarketByCode(code)
	return ok
}

func (t OrderTemplateDefinition) IsValidOrderType(code string) bool {
	_, ok := t.OrderTypeByCode(code)
	return ok
}

func (t OrderTemplateDefinition) MarketByCode(code string) (OrderMarketDefinition, bool) {
	normalized := normalizeMarketCode(code)
	for _, market := range t.Markets {
		if market.Code == normalized {
			return market, true
		}
	}
	return OrderMarketDefinition{}, false
}

func (t OrderTemplateDefinition) OrderTypeByCode(code string) (OrderTypeDefinition, bool) {
	normalized := normalizeOrderType(code)
	for _, item := range t.OrderTypes {
		if item.Code == normalized {
			return item, true
		}
	}
	return OrderTypeDefinition{}, false
}

func (t OrderTemplateDefinition) MarketName(code string) string {
	if market, ok := t.MarketByCode(code); ok {
		return market.Name
	}
	return strings.ToUpper(strings.TrimSpace(code))
}

func (t OrderTemplateDefinition) OrderTypeName(code string) string {
	if orderType, ok := t.OrderTypeByCode(code); ok {
		return orderType.Name
	}
	return strings.ToUpper(strings.TrimSpace(code))
}

func (t OrderTemplateDefinition) MarketSortIndex(code string) int {
	if market, ok := t.MarketByCode(code); ok {
		return market.SortOrder
	}
	return 99
}

func (t OrderTemplateDefinition) OrderTypeSortIndex(code string) int {
	if orderType, ok := t.OrderTypeByCode(code); ok {
		return orderType.SortOrder
	}
	return 99
}

func (t OrderTemplateDefinition) DefaultMarketEnabled(code string) bool {
	if market, ok := t.MarketByCode(code); ok {
		return market.DefaultEnabled
	}
	return false
}

func BuildOrderTemplateView(template OrderTemplateDefinition) OrderTemplateView {
	markets := make([]OrderMarketView, 0, len(template.Markets))
	for _, market := range template.Markets {
		markets = append(markets, OrderMarketView{
			Code:           market.Code,
			Name:           market.Name,
			DefaultEnabled: market.DefaultEnabled,
			SortOrder:      market.SortOrder,
		})
	}
	orderTypes := make([]OrderTypeView, 0, len(template.OrderTypes))
	for _, orderType := range template.OrderTypes {
		orderTypes = append(orderTypes, OrderTypeView{
			Code:      orderType.Code,
			Name:      orderType.Name,
			SortOrder: orderType.SortOrder,
		})
	}
	fields := make([]OrderFieldDefinition, len(template.Fields))
	copy(fields, template.Fields)
	return OrderTemplateView{
		TemplateVersion: template.TemplateVersion,
		TemplateName:    template.TemplateName,
		FormulaVersion:  template.FormulaVersion,
		SegmentCount:    len(template.Markets) * len(template.OrderTypes),
		MaxCardCount:    template.MaxCardCount,
		DeliveryEnabled: template.DeliveryEnabled,
		Markets:         markets,
		OrderTypes:      orderTypes,
		Fields:          fields,
	}
}

func builtInOrderTemplates() []OrderTemplateDefinition {
	return []OrderTemplateDefinition{
		{
			TemplateVersion: OrderTemplateVersionVIPServiceV1,
			TemplateName:    "贵宾服务订单模板 V1",
			FormulaVersion:  OrderFormulaVersionVIPServiceV1,
			MaxCardCount:    15,
			DeliveryEnabled: true,
			Markets: []OrderMarketDefinition{
				{Code: enum.MarketCodeLocal, Name: "本地市场", DefaultEnabled: true, SortOrder: 1},
				{Code: enum.MarketCodeRegional, Name: "区域市场", DefaultEnabled: false, SortOrder: 2},
				{Code: enum.MarketCodeNational, Name: "全国市场", DefaultEnabled: false, SortOrder: 3},
				{Code: enum.MarketCodeGlobal, Name: "全球市场", DefaultEnabled: false, SortOrder: 4},
			},
			OrderTypes: []OrderTypeDefinition{
				{Code: enum.OrderTypeAgencyInspection, Name: "代办过检", SortOrder: 1},
				{Code: enum.OrderTypeTwoCabinVIP, Name: "两舱贵宾", SortOrder: 2},
				{Code: enum.OrderTypeBusinessVIP, Name: "商务贵宾", SortOrder: 3},
				{Code: enum.OrderTypeMemberCustom, Name: "会员定制", SortOrder: 4},
			},
			Fields: []OrderFieldDefinition{
				{Code: "orderAmount", Name: "金额", ValueType: "money", Unit: "M", Order: 10},
				{Code: "orderQuantity", Name: "数量", ValueType: "number", Order: 20},
				{Code: "unitPrice", Name: "单价", ValueType: "money", Order: 30},
				{Code: "accountTerm", Name: "账期", ValueType: "number", Unit: "季度", Order: 40},
			},
		},
		{
			TemplateVersion: OrderTemplateVersionAirportV1,
			TemplateName:    "机场沙盘订单模板 V1",
			FormulaVersion:  OrderFormulaVersionAirportV1,
			MaxCardCount:    28,
			DeliveryEnabled: false,
			Markets: []OrderMarketDefinition{
				{Code: MarketCodeDomestic, Name: "国内市场", DefaultEnabled: true, SortOrder: 1},
				{Code: MarketCodeInternational, Name: "国际市场", DefaultEnabled: false, SortOrder: 2},
			},
			OrderTypes: []OrderTypeDefinition{
				{Code: OrderTypeNarrowBody, Name: "窄体", SortOrder: 1},
				{Code: OrderTypeWideBody, Name: "宽体", SortOrder: 2},
			},
			Fields: []OrderFieldDefinition{
				{Code: "routeArea", Name: "航线区域", ValueType: "text", Order: 10},
				{Code: "aircraftType", Name: "飞机类型", ValueType: "text", Order: 20},
				{Code: "routeType", Name: "航线类型", ValueType: "text", Order: 30},
				{Code: "runwayRequirement", Name: "跑道要求", ValueType: "text", Order: 40},
				{Code: "flightCount", Name: "架次", ValueType: "number", Order: 50},
				{Code: "throughput", Name: "吞吐量", ValueType: "number", Unit: "M", Order: 60},
				{Code: "loadFactor", Name: "客座率", ValueType: "percent", Unit: "%", Order: 70},
				{Code: "orderAmount", Name: "总收入", ValueType: "money", Unit: "M", Order: 80},
				{Code: "unitPrice", Name: "单价", ValueType: "money", Order: 90},
				{Code: "accountTerm", Name: "账期", ValueType: "number", Unit: "季度", Order: 100},
			},
		},
	}
}

func defaultOrderTemplate() OrderTemplateDefinition {
	return MustResolveOrderTemplateVersion(OrderTemplateVersionVIPServiceV1)
}
