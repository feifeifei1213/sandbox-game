package service

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
)

const orderGenerationFormulaVersion = OrderFormulaVersionVIPServiceV1

type orderGenerationParameterSnapshot struct {
	FormulaVersion string                         `json:"formulaVersion"`
	Volatility     float64                        `json:"volatility"`
	MinQuantity    int                            `json:"minQuantity"`
	MaxQuantity    int                            `json:"maxQuantity"`
	MinAccountTerm int                            `json:"minAccountTerm"`
	MaxAccountTerm int                            `json:"maxAccountTerm"`
	MaxCardCount   int                            `json:"maxCardCount"`
	AveragePrices  map[string]map[string]float64  `json:"averagePrices"`
	Airport        *airportOrderParameterSnapshot `json:"airport,omitempty"`
}

type generatedOrderDetail struct {
	YearNo         int            `json:"yearNo"`
	MarketCode     string         `json:"marketCode"`
	MarketName     string         `json:"marketName"`
	OrderType      string         `json:"orderType"`
	OrderTypeName  string         `json:"orderTypeName"`
	CardSequenceNo int            `json:"cardSequenceNo"`
	OrderAmount    float64        `json:"orderAmount"`
	OrderQuantity  float64        `json:"orderQuantity"`
	UnitPrice      float64        `json:"unitPrice"`
	AccountTerm    int            `json:"accountTerm"`
	SourceRowKey   string         `json:"sourceRowKey"`
	OrderPayload   map[string]any `json:"orderPayload,omitempty"`
}

type airportOrderParameterSnapshot struct {
	MinFlightCount int                `json:"minFlightCount"`
	MaxFlightCount int                `json:"maxFlightCount"`
	MinLoadFactor  int                `json:"minLoadFactor"`
	MaxLoadFactor  int                `json:"maxLoadFactor"`
	RunwayFactors  map[string]float64 `json:"runwayFactors"`
	RouteFactors   map[string]float64 `json:"routeFactors"`
}

type airportOrderPayload struct {
	RouteArea         string  `json:"routeArea"`
	AircraftType      string  `json:"aircraftType"`
	RouteType         string  `json:"routeType"`
	RunwayRequirement string  `json:"runwayRequirement"`
	FlightCount       int     `json:"flightCount"`
	Throughput        float64 `json:"throughput"`
	DisplayThroughput float64 `json:"displayThroughput"`
	LoadFactor        int     `json:"loadFactor"`
	RawRevenue        float64 `json:"rawRevenue"`
	DisplayUnitPrice  float64 `json:"displayUnitPrice"`
	RouteFactor       float64 `json:"routeFactor"`
	RunwayFactor      float64 `json:"runwayFactor"`
}

func buildGeneratedOrderPoolItems(configs []entity.OrderGenerationConfig, batchID int64, seed string, operatorName string, now time.Time, templates ...OrderTemplateDefinition) ([]entity.OrderPool, []generatedOrderDetail, orderGenerationParameterSnapshot, error) {
	template := defaultOrderTemplate()
	if len(templates) > 0 {
		template = templates[0]
	}
	params := defaultOrderGenerationParameters(template)
	rng := rand.New(rand.NewSource(parseSeed(seed)))
	poolItems := make([]entity.OrderPool, 0)
	details := make([]generatedOrderDetail, 0)
	for _, config := range configs {
		if config.OrderCount < 0 {
			return nil, nil, params, ErrAdminOrderConfigInvalid
		}
		if config.OrderCount == 0 {
			continue
		}
		if params.MaxCardCount > 0 && config.OrderCount > params.MaxCardCount {
			return nil, nil, params, ErrAdminOrderConfigInvalid
		}
		for index := 1; index <= config.OrderCount; index++ {
			item, detail, err := buildGeneratedOrderPoolItem(config, batchID, index, rng, params, operatorName, now, template)
			if err != nil {
				return nil, nil, params, err
			}
			poolItems = append(poolItems, item)
			details = append(details, detail)
		}
	}
	return poolItems, details, params, nil
}

func buildGeneratedOrderPoolItem(config entity.OrderGenerationConfig, batchID int64, index int, rng *rand.Rand, params orderGenerationParameterSnapshot, operatorName string, now time.Time, template OrderTemplateDefinition) (entity.OrderPool, generatedOrderDetail, error) {
	if template.TemplateVersion == OrderTemplateVersionAirportV1 {
		return buildAirportGeneratedOrderPoolItem(config, batchID, index, rng, params, operatorName, now, template)
	}
	return buildVIPGeneratedOrderPoolItem(config, batchID, index, rng, params, operatorName, now, template)
}

func buildVIPGeneratedOrderPoolItem(config entity.OrderGenerationConfig, batchID int64, index int, rng *rand.Rand, params orderGenerationParameterSnapshot, operatorName string, now time.Time, template OrderTemplateDefinition) (entity.OrderPool, generatedOrderDetail, error) {
	avgPrice, ok := params.AveragePrices[config.MarketCode][config.OrderType]
	if !ok {
		return entity.OrderPool{}, generatedOrderDetail{}, ErrAdminOrderConfigInvalid
	}
	quantity := randomIntInclusive(rng, params.MinQuantity, params.MaxQuantity)
	accountTerm := randomIntInclusive(rng, params.MinAccountTerm, params.MaxAccountTerm)
	rawAmount := avgPrice*float64(quantity) - params.Volatility*(float64(quantity)-float64(params.MinQuantity+params.MaxQuantity)/2)
	amount := validIntegerOrderAmount(rawAmount, quantity)
	unitPrice := roundMoney(amount / float64(quantity))
	return buildGeneratedOrderPoolEntity(config, batchID, index, amount, float64(quantity), unitPrice, accountTerm, nil, operatorName, now, template), generatedOrderDetail{
		YearNo:         config.YearNo,
		MarketCode:     config.MarketCode,
		MarketName:     template.MarketName(config.MarketCode),
		OrderType:      config.OrderType,
		OrderTypeName:  template.OrderTypeName(config.OrderType),
		CardSequenceNo: index,
		OrderAmount:    amount,
		OrderQuantity:  float64(quantity),
		UnitPrice:      unitPrice,
		AccountTerm:    accountTerm,
		SourceRowKey:   sourceRowKey(config, index),
	}, nil
}

func buildAirportGeneratedOrderPoolItem(config entity.OrderGenerationConfig, batchID int64, index int, rng *rand.Rand, params orderGenerationParameterSnapshot, operatorName string, now time.Time, template OrderTemplateDefinition) (entity.OrderPool, generatedOrderDetail, error) {
	if params.Airport == nil {
		return entity.OrderPool{}, generatedOrderDetail{}, ErrAdminOrderConfigInvalid
	}
	flightCount := randomIntInclusive(rng, params.Airport.MinFlightCount, params.Airport.MaxFlightCount)
	loadFactor := randomIntInclusive(rng, params.Airport.MinLoadFactor, params.Airport.MaxLoadFactor)
	routeType := template.MarketName(config.MarketCode)
	routeArea := "国内"
	routeFactor := params.Airport.RouteFactors[MarketCodeDomestic]
	if config.MarketCode == MarketCodeInternational {
		areas := []string{"亚太", "欧非", "美洲"}
		routeArea = areas[rng.Intn(len(areas))]
		routeFactor = params.Airport.RouteFactors[MarketCodeInternational]
	}
	aircraftType := template.OrderTypeName(config.OrderType)
	runwayRequirement := "4C"
	if config.OrderType == OrderTypeWideBody {
		runways := []string{"4D", "4E", "4F"}
		runwayRequirement = runways[rng.Intn(len(runways))]
	}
	runwayFactor := params.Airport.RunwayFactors[runwayRequirement]
	throughput := 160 * float64(flightCount) * float64(loadFactor) / 1000000
	unitPrice := float64(loadFactor) / 200 * 0.012 * routeFactor * runwayFactor
	rawRevenue := unitPrice * float64(flightCount)
	amount := math.Round(rawRevenue)
	accountTerm := randomIntInclusive(rng, params.MinAccountTerm, params.MaxAccountTerm)
	displayUnitPrice := math.Round(unitPrice*10000) / 10000
	payload := airportOrderPayload{
		RouteArea:         routeArea,
		AircraftType:      aircraftType,
		RouteType:         routeType,
		RunwayRequirement: runwayRequirement,
		FlightCount:       flightCount,
		Throughput:        throughput,
		DisplayThroughput: math.Round(throughput),
		LoadFactor:        loadFactor,
		RawRevenue:        rawRevenue,
		DisplayUnitPrice:  displayUnitPrice,
		RouteFactor:       routeFactor,
		RunwayFactor:      runwayFactor,
	}
	payloadJSON, payloadMap := marshalOrderPayload(payload)
	return buildGeneratedOrderPoolEntity(config, batchID, index, amount, float64(flightCount), displayUnitPrice, accountTerm, payloadJSON, operatorName, now, template), generatedOrderDetail{
		YearNo:         config.YearNo,
		MarketCode:     config.MarketCode,
		MarketName:     template.MarketName(config.MarketCode),
		OrderType:      config.OrderType,
		OrderTypeName:  template.OrderTypeName(config.OrderType),
		CardSequenceNo: index,
		OrderAmount:    amount,
		OrderQuantity:  float64(flightCount),
		UnitPrice:      displayUnitPrice,
		AccountTerm:    accountTerm,
		SourceRowKey:   sourceRowKey(config, index),
		OrderPayload:   payloadMap,
	}, nil
}

func buildGeneratedOrderPoolEntity(config entity.OrderGenerationConfig, batchID int64, index int, amount float64, quantity float64, unitPrice float64, accountTerm int, payloadJSON []byte, operatorName string, now time.Time, template OrderTemplateDefinition) entity.OrderPool {
	generationBatchID := batchID
	rowIndex := index
	businessOrderNo := fmt.Sprintf("CARD-%02d", index)
	return entity.OrderPool{
		OrderTemplateVersion: template.TemplateVersion,
		YearNo:               config.YearNo,
		MarketCode:           config.MarketCode,
		OrderType:            config.OrderType,
		SegmentCode:          fmt.Sprintf("%s_%s", config.MarketCode, config.OrderType),
		CardSequenceNo:       index,
		BusinessOrderNo:      businessOrderNo,
		OrderAmount:          amount,
		OrderQuantity:        quantity,
		UnitPrice:            unitPrice,
		AccountTerm:          accountTerm,
		PoolStatus:           enum.OrderPoolStatusAvailable,
		GenerationBatchID:    &generationBatchID,
		SourceSheetName:      "系统生成",
		SourceCell:           businessOrderNo,
		SourceRowIndex:       &rowIndex,
		SourceRowKey:         sourceRowKey(config, index),
		OrderPayloadJSON:     payloadJSON,
		BaseEntity: entity.BaseEntity{
			Creator:    operatorName,
			CreateTime: now,
			Updater:    operatorName,
			UpdateTime: now,
		},
	}
}

func sourceRowKey(config entity.OrderGenerationConfig, index int) string {
	return fmt.Sprintf("%d_%s_%s_%02d", config.YearNo, config.MarketCode, config.OrderType, index)
}

func marshalOrderPayload(value any) ([]byte, map[string]any) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, nil
	}
	result := map[string]any{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return raw, nil
	}
	return raw, result
}

func defaultOrderGenerationParameters(template OrderTemplateDefinition) orderGenerationParameterSnapshot {
	params := orderGenerationParameterSnapshot{
		FormulaVersion: template.FormulaVersion,
		Volatility:     1,
		MinQuantity:    2,
		MaxQuantity:    4,
		MinAccountTerm: 2,
		MaxAccountTerm: 4,
		MaxCardCount:   template.MaxCardCount,
	}
	if template.TemplateVersion == OrderTemplateVersionAirportV1 {
		params.Airport = &airportOrderParameterSnapshot{
			MinFlightCount: 9000,
			MaxFlightCount: 11000,
			MinLoadFactor:  60,
			MaxLoadFactor:  80,
			RunwayFactors: map[string]float64{
				"4C": 1,
				"4D": 2,
				"4E": 2.4,
				"4F": 3,
			},
			RouteFactors: map[string]float64{
				MarketCodeDomestic:      1,
				MarketCodeInternational: 1.7,
			},
		}
		return params
	}
	params.AveragePrices = map[string]map[string]float64{
		enum.MarketCodeLocal: {
			enum.OrderTypeAgencyInspection: 5,
			enum.OrderTypeTwoCabinVIP:      6,
			enum.OrderTypeBusinessVIP:      8,
			enum.OrderTypeMemberCustom:     11,
		},
		enum.MarketCodeRegional: {
			enum.OrderTypeAgencyInspection: 5.5,
			enum.OrderTypeTwoCabinVIP:      6.5,
			enum.OrderTypeBusinessVIP:      8.5,
			enum.OrderTypeMemberCustom:     11.5,
		},
		enum.MarketCodeNational: {
			enum.OrderTypeAgencyInspection: 6,
			enum.OrderTypeTwoCabinVIP:      7,
			enum.OrderTypeBusinessVIP:      9,
			enum.OrderTypeMemberCustom:     12,
		},
		enum.MarketCodeGlobal: {
			enum.OrderTypeAgencyInspection: 7,
			enum.OrderTypeTwoCabinVIP:      8,
			enum.OrderTypeBusinessVIP:      10,
			enum.OrderTypeMemberCustom:     13,
		},
	}
	return params
}

func defaultVIPOrderGenerationParameters() orderGenerationParameterSnapshot {
	return orderGenerationParameterSnapshot{
		FormulaVersion: OrderFormulaVersionVIPServiceV1,
		Volatility:     1,
		MinQuantity:    2,
		MaxQuantity:    4,
		MinAccountTerm: 2,
		MaxAccountTerm: 4,
		MaxCardCount:   0,
		AveragePrices: map[string]map[string]float64{
			enum.MarketCodeLocal: {
				enum.OrderTypeAgencyInspection: 5,
				enum.OrderTypeTwoCabinVIP:      6,
				enum.OrderTypeBusinessVIP:      8,
				enum.OrderTypeMemberCustom:     11,
			},
			enum.MarketCodeRegional: {
				enum.OrderTypeAgencyInspection: 5.5,
				enum.OrderTypeTwoCabinVIP:      6.5,
				enum.OrderTypeBusinessVIP:      8.5,
				enum.OrderTypeMemberCustom:     11.5,
			},
			enum.MarketCodeNational: {
				enum.OrderTypeAgencyInspection: 6,
				enum.OrderTypeTwoCabinVIP:      7,
				enum.OrderTypeBusinessVIP:      9,
				enum.OrderTypeMemberCustom:     12,
			},
			enum.MarketCodeGlobal: {
				enum.OrderTypeAgencyInspection: 7,
				enum.OrderTypeTwoCabinVIP:      8,
				enum.OrderTypeBusinessVIP:      10,
				enum.OrderTypeMemberCustom:     13,
			},
		},
	}
}

func randomIntInclusive(rng *rand.Rand, min int, max int) int {
	if max <= min {
		return min
	}
	return min + rng.Intn(max-min+1)
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}

func validIntegerOrderAmount(rawAmount float64, quantity int) float64 {
	target := int(math.Round(rawAmount))
	if target < 1 {
		target = 1
	}
	if quantity <= 0 {
		return float64(target)
	}
	for offset := 0; ; offset++ {
		if candidate := target - offset; candidate >= 1 && isValidIntegerAmountForQuantity(candidate, quantity) {
			return float64(candidate)
		}
		if candidate := target + offset; candidate >= 1 && isValidIntegerAmountForQuantity(candidate, quantity) {
			return float64(candidate)
		}
	}
}

func isValidIntegerAmountForQuantity(amount int, quantity int) bool {
	if quantity <= 0 {
		return false
	}
	return (amount*100)%quantity == 0
}
