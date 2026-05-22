package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
)

func TestPlayerOrderYearViewForDemoYearDoesNotRequireOrders(t *testing.T) {
	view, err := (&PlayerOrderQueryService{}).GetYearView(context.Background(), 101, 0)
	if err != nil {
		t.Fatalf("expected demo year order view to load without repositories, got %v", err)
	}
	if view.OrderRequired {
		t.Fatalf("expected 0 year to skip order flow")
	}
	if view.CanSubmitInvestment || view.InvestmentSubmitted {
		t.Fatalf("expected demo year to have no investment action, got %#v", view)
	}
	if view.PollingIntervalSeconds != orderPollingIntervalSeconds {
		t.Fatalf("expected polling interval %d, got %d", orderPollingIntervalSeconds, view.PollingIntervalSeconds)
	}
}

func TestValidateForecastControlCommandEnforcesOrderCountAndReleaseSequence(t *testing.T) {
	valid := UpdateOrderControlConfigCommand{
		YearNo: 1,
		Items:  buildTestControlConfigItems(1, map[string]int{testOrderSegmentKey(enum.MarketCodeLocal, enum.OrderTypeAgencyInspection): 15}),
	}
	if err := validateControlConfigCommand(valid); err != nil {
		t.Fatalf("expected valid release sequence config, got %v", err)
	}

	forecast := UpdateOrderForecastControlCommand{
		Items: buildTestForecastControlItems(map[string]int{
			testForecastSegmentKey(1, enum.MarketCodeLocal, enum.OrderTypeAgencyInspection): 15,
		}),
	}
	if err := validateForecastControlCommand(forecast); err != nil {
		t.Fatalf("expected valid 0~15 forecast control count, got %v", err)
	}

	invalidForecast := forecast
	invalidForecast.Items = append([]UpdateOrderForecastControlItem(nil), forecast.Items...)
	invalidForecast.Items[0].OrderCount = 16
	if err := validateForecastControlCommand(invalidForecast); !errors.Is(err, ErrAdminOrderForecastControlInvalid) {
		t.Fatalf("expected forecast order count > 15 to be invalid, got %v", err)
	}

	duplicatedEffectiveSequence := valid
	duplicatedEffectiveSequence.Items = append([]UpdateOrderControlConfigItem(nil), valid.Items...)
	duplicatedEffectiveSequence.Items[1].ReleaseSequenceNo = duplicatedEffectiveSequence.Items[0].ReleaseSequenceNo
	forecastCounts := forecastControlCountMap(1, []entity.OrderForecastControl{
		{YearNo: 1, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeAgencyInspection, OrderCount: 1},
		{YearNo: 1, MarketCode: enum.MarketCodeLocal, OrderType: enum.OrderTypeTwoCabinVIP, OrderCount: 1},
	})
	if err := validateEffectiveReleaseSequences(duplicatedEffectiveSequence.Items, buildMarketEnabledMap(nil), forecastCounts, 1); !errors.Is(err, ErrAdminOrderReleaseSequenceDuplicated) {
		t.Fatalf("expected duplicated effective release sequence to be rejected, got %v", err)
	}
}

func TestNormalizeMarketInvestmentInputsRequiresAllSixteenSegmentsAndAllowsZero(t *testing.T) {
	inputs := buildTestMarketInvestments(func(segment OrderSegmentDefinition) float64 {
		if segment.MarketCode == enum.MarketCodeLocal && segment.OrderType == enum.OrderTypeAgencyInspection {
			return 3
		}
		return 0
	})
	normalized, err := normalizeMarketInvestmentInputs(inputs)
	if err != nil {
		t.Fatalf("expected complete 16 segment investments to be valid, got %v", err)
	}
	if len(normalized) != requiredOrderInvestmentSegmentCount {
		t.Fatalf("expected %d normalized investments, got %d", requiredOrderInvestmentSegmentCount, len(normalized))
	}

	if _, err := normalizeMarketInvestmentInputs(inputs[:15]); !errors.Is(err, ErrOrderInvestmentInvalid) {
		t.Fatalf("expected missing segment to be invalid, got %v", err)
	}

	duplicated := append([]MarketInvestmentInput(nil), inputs...)
	duplicated[15] = duplicated[0]
	if _, err := normalizeMarketInvestmentInputs(duplicated); !errors.Is(err, ErrOrderInvestmentInvalid) {
		t.Fatalf("expected duplicated segment to be invalid, got %v", err)
	}

	negative := append([]MarketInvestmentInput(nil), inputs...)
	negative[0].MarketInvestment = -1
	if _, err := normalizeMarketInvestmentInputs(negative); !errors.Is(err, ErrOrderInvestmentInvalid) {
		t.Fatalf("expected negative investment to be invalid, got %v", err)
	}
}

func TestBuildGeneratedOrderPoolItemsIsDeterministicAndCapsCardCount(t *testing.T) {
	now := time.Date(2026, 5, 20, 9, 30, 0, 0, time.UTC)
	configs := []entity.OrderGenerationConfig{
		{
			YearNo:            1,
			MarketCode:        enum.MarketCodeLocal,
			OrderType:         enum.OrderTypeAgencyInspection,
			OrderCount:        2,
			ReleaseSequenceNo: 1,
		},
	}

	leftItems, leftDetails, leftParams, err := buildGeneratedOrderPoolItems(configs, 91, "stable-seed", "tester", now)
	if err != nil {
		t.Fatalf("generate left order pool: %v", err)
	}
	rightItems, rightDetails, rightParams, err := buildGeneratedOrderPoolItems(configs, 91, "stable-seed", "tester", now)
	if err != nil {
		t.Fatalf("generate right order pool: %v", err)
	}
	if leftParams.FormulaVersion != orderGenerationFormulaVersion || rightParams.MaxCardCount != 15 {
		t.Fatalf("unexpected generation parameters: left=%#v right=%#v", leftParams, rightParams)
	}
	if len(leftItems) != 2 || len(leftDetails) != 2 || len(rightItems) != 2 || len(rightDetails) != 2 {
		t.Fatalf("expected two generated orders, got left=%d/%d right=%d/%d", len(leftItems), len(leftDetails), len(rightItems), len(rightDetails))
	}
	for index := range leftItems {
		if leftItems[index].OrderAmount != rightItems[index].OrderAmount ||
			leftItems[index].OrderQuantity != rightItems[index].OrderQuantity ||
			leftItems[index].UnitPrice != rightItems[index].UnitPrice ||
			leftItems[index].AccountTerm != rightItems[index].AccountTerm ||
			leftItems[index].SourceRowKey != rightItems[index].SourceRowKey {
			t.Fatalf("expected deterministic generation at index %d, left=%#v right=%#v", index, leftItems[index], rightItems[index])
		}
	}
	if leftItems[0].SegmentCode != "LOCAL_AGENCY_INSPECTION" || leftItems[0].SourceSheetName != "系统生成" || leftItems[0].CardSequenceNo != 1 {
		t.Fatalf("unexpected generated order metadata: %#v", leftItems[0])
	}
	for index, item := range leftItems {
		if item.OrderAmount != float64(int64(item.OrderAmount)) {
			t.Fatalf("expected integer order amount at index %d, got %#v", index, item)
		}
		if !almostEqual(item.OrderAmount, item.OrderQuantity*item.UnitPrice) {
			t.Fatalf("expected amount to equal quantity * unit price at index %d, got %#v", index, item)
		}
		if item.OrderQuantity != float64(int64(item.OrderQuantity)) {
			t.Fatalf("expected integer quantity at index %d, got %#v", index, item)
		}
	}

	invalid := append([]entity.OrderGenerationConfig(nil), configs...)
	invalid[0].OrderCount = 16
	if _, _, _, err := buildGeneratedOrderPoolItems(invalid, 91, "stable-seed", "tester", now); !errors.Is(err, ErrAdminOrderConfigInvalid) {
		t.Fatalf("expected generated order count > max to be rejected, got %v", err)
	}
}

func TestValidIntegerOrderAmountAllowsDecimalUnitPriceWhenAmountIsInteger(t *testing.T) {
	if amount := validIntegerOrderAmount(10.8, 2); amount != 11 {
		t.Fatalf("expected nearest integer amount for quantity 2, got %v", amount)
	}
	if unitPrice := roundMoney(validIntegerOrderAmount(10.8, 2) / 2); unitPrice != 5.5 {
		t.Fatalf("expected decimal unit price to be allowed for quantity 2, got %v", unitPrice)
	}
	if amount := validIntegerOrderAmount(17.01, 3); amount != 18 {
		t.Fatalf("expected quantity 3 to avoid invalid two-decimal multiplication, got %v", amount)
	}
}

func almostEqual(left float64, right float64) bool {
	if left > right {
		return left-right < 0.000001
	}
	return right-left < 0.000001
}

func TestBuildMarketParticipantsHonorsLeaderAndPreviousAmountTieBreak(t *testing.T) {
	groups := []entity.Group{
		{ID: 1, GroupName: "一组", BusinessStatus: enum.BusinessStatusNormal},
		{ID: 2, GroupName: "二组", BusinessStatus: enum.BusinessStatusNormal},
		{ID: 3, GroupName: "三组", BusinessStatus: enum.BusinessStatusBankrupt},
	}

	leader := int64(1)
	participants, _, err := buildMarketParticipantsWithLeader(
		2,
		enum.MarketCodeLocal,
		groups,
		[]entity.GroupMarketBid{
			{GroupID: 1, MarketInvestment: 0},
			{GroupID: 2, MarketInvestment: 100},
			{GroupID: 3, MarketInvestment: 999},
		},
		map[int64]float64{1: 10, 2: 20, 3: 999},
		"leader-first",
		&leader,
	)
	if err != nil {
		t.Fatalf("build participants with leader: %v", err)
	}
	if len(participants) != 2 {
		t.Fatalf("expected bankrupt group to be excluded, got %#v", participants)
	}
	if participants[0].GroupID != 1 || !participants[0].IsMarketLeader {
		t.Fatalf("expected market leader to be first even with zero current investment, got %#v", participants)
	}

	tieParticipants, _, err := buildMarketParticipantsWithLeader(
		2,
		enum.MarketCodeLocal,
		groups,
		[]entity.GroupMarketBid{
			{GroupID: 1, MarketInvestment: 10},
			{GroupID: 2, MarketInvestment: 10},
		},
		map[int64]float64{1: 15, 2: 30},
		"previous-amount",
		nil,
	)
	if err != nil {
		t.Fatalf("build tie participants: %v", err)
	}
	if tieParticipants[0].GroupID != 2 {
		t.Fatalf("expected previous-year market order amount to break equal investment tie, got %#v", tieParticipants)
	}

	resolvedLeader := resolveMarketLeader(2, groups, map[int64]float64{1: 10, 2: 20, 3: 999}, rand.New(rand.NewSource(7)))
	if resolvedLeader == nil || *resolvedLeader != 2 {
		t.Fatalf("expected bankrupt previous leader to be invalid and group 2 to lead, got %v", resolvedLeader)
	}
}

func buildTestControlConfigItems(yearNo int, counts map[string]int) []UpdateOrderControlConfigItem {
	items := make([]UpdateOrderControlConfigItem, 0, len(defaultOrderSegments()))
	for index, segment := range defaultOrderSegments() {
		items = append(items, UpdateOrderControlConfigItem{
			MarketCode:        segment.MarketCode,
			OrderType:         segment.OrderType,
			OrderCount:        counts[testOrderSegmentKey(segment.MarketCode, segment.OrderType)],
			ReleaseSequenceNo: index + 1,
		})
	}
	return items
}

func buildTestForecastControlItems(counts map[string]int) []UpdateOrderForecastControlItem {
	items := make([]UpdateOrderForecastControlItem, 0, forecastControlMaxYear*len(defaultOrderSegments()))
	for yearNo := forecastControlMinYear; yearNo <= forecastControlMaxYear; yearNo++ {
		for _, segment := range defaultOrderSegments() {
			items = append(items, UpdateOrderForecastControlItem{
				YearNo:     yearNo,
				MarketCode: segment.MarketCode,
				OrderType:  segment.OrderType,
				OrderCount: counts[testForecastSegmentKey(yearNo, segment.MarketCode, segment.OrderType)],
			})
		}
	}
	return items
}

func buildTestMarketInvestments(resolve func(segment OrderSegmentDefinition) float64) []MarketInvestmentInput {
	inputs := make([]MarketInvestmentInput, 0, len(defaultOrderSegments()))
	for _, segment := range defaultOrderSegments() {
		inputs = append(inputs, MarketInvestmentInput{
			MarketCode:       segment.MarketCode,
			OrderType:        segment.OrderType,
			MarketInvestment: resolve(segment),
		})
	}
	return inputs
}

func testOrderSegmentKey(marketCode string, orderType string) string {
	return marketCode + "|" + orderType
}

func testForecastSegmentKey(yearNo int, marketCode string, orderType string) string {
	return fmt.Sprintf("%d|%s|%s", yearNo, marketCode, orderType)
}
