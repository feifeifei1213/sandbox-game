package service

import (
	"context"
	"fmt"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
)

type OrderOperatingLinkService struct {
	bidRepo       *repository.GroupMarketBidRepository
	stateRepo     *repository.MarketBiddingStateRepository
	selectionRepo *repository.GroupOrderSelectionRepository
}

func NewOrderOperatingLinkService(
	bidRepo *repository.GroupMarketBidRepository,
	stateRepo *repository.MarketBiddingStateRepository,
	selectionRepo *repository.GroupOrderSelectionRepository,
) *OrderOperatingLinkService {
	return &OrderOperatingLinkService{
		bidRepo:       bidRepo,
		stateRepo:     stateRepo,
		selectionRepo: selectionRepo,
	}
}

type FormalYearOrderLinkValues struct {
	MarketInvestmentTotal float64
	OrderTotal            float64
	MarketRows            []map[string]any
	PrerequisiteCompleted bool
}

func (s *OrderOperatingLinkService) LoadFormalYearValues(ctx context.Context, groupID int64, yearNo int) (FormalYearOrderLinkValues, error) {
	result := FormalYearOrderLinkValues{
		MarketRows:            buildEmptyOrderMarketRows(),
		PrerequisiteCompleted: yearNo == 0,
	}
	if yearNo == 0 {
		return result, nil
	}
	if s == nil {
		return result, nil
	}

	bids, err := s.bidRepo.ListByGroupYear(ctx, groupID, yearNo)
	if err != nil {
		return result, fmt.Errorf("list group market bids: %w", err)
	}
	for _, bid := range bids {
		result.MarketInvestmentTotal += bid.MarketInvestment
	}

	segmentAmounts, err := s.selectionRepo.SumSelectedAmountByGroupYearSegment(ctx, groupID, yearNo)
	if err != nil {
		return result, fmt.Errorf("sum selected order amount: %w", err)
	}
	for _, item := range segmentAmounts {
		result.OrderTotal += item.Amount
		fillOrderMarketRowAmount(result.MarketRows, item.MarketCode, item.OrderType, item.Amount)
	}

	completed, err := s.IsPrerequisiteCompleted(ctx, yearNo)
	if err != nil {
		return result, err
	}
	result.PrerequisiteCompleted = completed
	return result, nil
}

func (s *OrderOperatingLinkService) IsPrerequisiteCompleted(ctx context.Context, yearNo int) (bool, error) {
	if yearNo == 0 {
		return true, nil
	}
	states, err := s.stateRepo.ListByYear(ctx, yearNo)
	if err != nil {
		return false, fmt.Errorf("list order segment states: %w", err)
	}
	if len(states) == 0 {
		return false, nil
	}
	for _, item := range states {
		if !isTerminalOrderSegmentStatus(item.SegmentStatus) {
			return false, nil
		}
	}
	return true, nil
}

func (s *OrderOperatingLinkService) ApplyFormalYearValues(ctx context.Context, groupID int64, yearNo int, source payload.OperatingPayload) (payload.OperatingPayload, FormalYearOrderLinkValues, error) {
	normalized := source.Normalize().WithoutDerivedValues()
	values, err := s.LoadFormalYearValues(ctx, groupID, yearNo)
	if err != nil {
		return normalized, values, err
	}
	if yearNo == 0 {
		return normalized, values, nil
	}
	return ApplyOrderLinkValuesToOperatingPayload(normalized, values), values, nil
}

func ApplyOrderLinkValuesToOperatingPayload(source payload.OperatingPayload, values FormalYearOrderLinkValues) payload.OperatingPayload {
	next := source.Normalize()
	next.Beginning.MarketBid = values.MarketRows
	next.Beginning.TaxAndPlanning["orderLinked"] = true
	next.Beginning.TaxAndPlanning["marketInvestmentTotal"] = values.MarketInvestmentTotal
	next.Beginning.TaxAndPlanning["marketBidCost"] = values.MarketInvestmentTotal
	next.Beginning.TaxAndPlanning["orderTotal"] = values.OrderTotal
	next.Derived.Values["marketInvestmentTotal"] = values.MarketInvestmentTotal
	next.Derived.Values["marketBidCost"] = values.MarketInvestmentTotal
	next.Derived.Values["orderTotal"] = values.OrderTotal
	return next
}

func buildEmptyOrderMarketRows() []map[string]any {
	rows := make([]map[string]any, 0, len(defaultOrderMarkets()))
	for _, market := range defaultOrderMarkets() {
		rows = append(rows, map[string]any{
			"marketCode":              market.code,
			"basicProductTotal":       0,
			"standardProductTotal":    0,
			"precisionProductTotal":   0,
			"intelligentProductTotal": 0,
			"agencyInspectionTotal":   0,
			"twoCabinVipTotal":        0,
			"businessVipTotal":        0,
			"memberCustomTotal":       0,
			"orderAmount":             0,
		})
	}
	return rows
}

func fillOrderMarketRowAmount(rows []map[string]any, marketCode string, orderType string, amount float64) {
	index := orderMarketRowIndex(marketCode)
	if index < 0 || index >= len(rows) {
		return
	}
	key := orderTypeMarketBidKey(orderType)
	if key != "" {
		rows[index][key] = amount
	}
	rows[index]["orderAmount"] = toFloat(rows[index]["orderAmount"]) + amount
}

func orderMarketRowIndex(marketCode string) int {
	switch marketCode {
	case enum.MarketCodeLocal:
		return 0
	case enum.MarketCodeRegional:
		return 1
	case enum.MarketCodeNational:
		return 2
	case enum.MarketCodeGlobal:
		return 3
	default:
		return -1
	}
}

func orderTypeMarketBidKey(orderType string) string {
	switch orderType {
	case enum.OrderTypeAgencyInspection:
		return "agencyInspectionTotal"
	case enum.OrderTypeTwoCabinVIP:
		return "twoCabinVipTotal"
	case enum.OrderTypeBusinessVIP:
		return "businessVipTotal"
	case enum.OrderTypeMemberCustom:
		return "memberCustomTotal"
	default:
		return ""
	}
}

func toFloat(value any) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case float32:
		return float64(typed)
	case int:
		return float64(typed)
	case int64:
		return float64(typed)
	case int32:
		return float64(typed)
	default:
		return 0
	}
}
