package service

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
)

const orderGenerationFormulaVersion = "ORDER_GEN_SERVICE_V1"

type orderGenerationParameterSnapshot struct {
	FormulaVersion string                        `json:"formulaVersion"`
	Volatility     float64                       `json:"volatility"`
	MinQuantity    int                           `json:"minQuantity"`
	MaxQuantity    int                           `json:"maxQuantity"`
	MinAccountTerm int                           `json:"minAccountTerm"`
	MaxAccountTerm int                           `json:"maxAccountTerm"`
	MaxCardCount   int                           `json:"maxCardCount"`
	AveragePrices  map[string]map[string]float64 `json:"averagePrices"`
}

type generatedOrderDetail struct {
	YearNo         int     `json:"yearNo"`
	MarketCode     string  `json:"marketCode"`
	MarketName     string  `json:"marketName"`
	OrderType      string  `json:"orderType"`
	OrderTypeName  string  `json:"orderTypeName"`
	CardSequenceNo int     `json:"cardSequenceNo"`
	OrderAmount    float64 `json:"orderAmount"`
	OrderQuantity  float64 `json:"orderQuantity"`
	UnitPrice      float64 `json:"unitPrice"`
	AccountTerm    int     `json:"accountTerm"`
	SourceRowKey   string  `json:"sourceRowKey"`
}

func buildGeneratedOrderPoolItems(configs []entity.OrderGenerationConfig, batchID int64, seed string, operatorName string, now time.Time) ([]entity.OrderPool, []generatedOrderDetail, orderGenerationParameterSnapshot, error) {
	params := defaultOrderGenerationParameters()
	rng := rand.New(rand.NewSource(parseSeed(seed)))
	poolItems := make([]entity.OrderPool, 0)
	details := make([]generatedOrderDetail, 0)
	for _, config := range configs {
		if config.OrderCount <= 0 {
			continue
		}
		if config.OrderCount > params.MaxCardCount {
			return nil, nil, params, ErrAdminOrderConfigInvalid
		}
		avgPrice, ok := params.AveragePrices[config.MarketCode][config.OrderType]
		if !ok {
			return nil, nil, params, ErrAdminOrderConfigInvalid
		}
		for index := 1; index <= config.OrderCount; index++ {
			quantity := randomIntInclusive(rng, params.MinQuantity, params.MaxQuantity)
			accountTerm := randomIntInclusive(rng, params.MinAccountTerm, params.MaxAccountTerm)
			rawAmount := avgPrice*float64(quantity) - params.Volatility*(float64(quantity)-float64(params.MinQuantity+params.MaxQuantity)/2)
			amount := validIntegerOrderAmount(rawAmount, quantity)
			unitPrice := roundMoney(amount / float64(quantity))
			sourceRowKey := fmt.Sprintf("%d_%s_%s_%02d", config.YearNo, config.MarketCode, config.OrderType, index)
			segmentCode := fmt.Sprintf("%s_%s", config.MarketCode, config.OrderType)
			businessOrderNo := fmt.Sprintf("CARD-%02d", index)
			generationBatchID := batchID
			rowIndex := index
			item := entity.OrderPool{
				YearNo:            config.YearNo,
				MarketCode:        config.MarketCode,
				OrderType:         config.OrderType,
				SegmentCode:       segmentCode,
				CardSequenceNo:    index,
				BusinessOrderNo:   businessOrderNo,
				OrderAmount:       amount,
				OrderQuantity:     float64(quantity),
				UnitPrice:         unitPrice,
				AccountTerm:       accountTerm,
				PoolStatus:        enum.OrderPoolStatusAvailable,
				GenerationBatchID: &generationBatchID,
				SourceSheetName:   "系统生成",
				SourceCell:        businessOrderNo,
				SourceRowIndex:    &rowIndex,
				SourceRowKey:      sourceRowKey,
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: now,
					Updater:    operatorName,
					UpdateTime: now,
				},
			}
			poolItems = append(poolItems, item)
			details = append(details, generatedOrderDetail{
				YearNo:         config.YearNo,
				MarketCode:     config.MarketCode,
				MarketName:     marketName(config.MarketCode),
				OrderType:      config.OrderType,
				OrderTypeName:  orderTypeName(config.OrderType),
				CardSequenceNo: index,
				OrderAmount:    amount,
				OrderQuantity:  float64(quantity),
				UnitPrice:      unitPrice,
				AccountTerm:    accountTerm,
				SourceRowKey:   sourceRowKey,
			})
		}
	}
	return poolItems, details, params, nil
}

func defaultOrderGenerationParameters() orderGenerationParameterSnapshot {
	return orderGenerationParameterSnapshot{
		FormulaVersion: orderGenerationFormulaVersion,
		Volatility:     1,
		MinQuantity:    2,
		MaxQuantity:    4,
		MinAccountTerm: 2,
		MaxAccountTerm: 4,
		MaxCardCount:   15,
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
