package summary

import (
	"testing"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
)

func TestBuildMarksFormalCompletedYearAsEffective(t *testing.T) {
	t.Parallel()

	builder := NewBuilder()
	result, err := builder.Build(newSummaryContext(1), payload.ReportComputedPayload{
		ReportSalesRevenue: 100,
		ReportNetProfit:    30,
		ReportTotalEquity:  88,
	})
	if err != nil {
		t.Fatalf("build summary failed: %v", err)
	}

	if !result.SummaryEffective {
		t.Fatal("expected formal completed year to be summary effective")
	}
	if result.RankingValue != 88 {
		t.Fatalf("expected ranking value 88, got %v", result.RankingValue)
	}
}

func TestBuildKeepsDemoYearOutOfSummary(t *testing.T) {
	t.Parallel()

	builder := NewBuilder()
	result, err := builder.Build(newSummaryContext(0), payload.ReportComputedPayload{
		ReportSalesRevenue: 100,
		ReportNetProfit:    30,
		ReportTotalEquity:  88,
	})
	if err != nil {
		t.Fatalf("build summary failed: %v", err)
	}

	if result.SummaryEffective {
		t.Fatal("expected demo year not to be summary effective")
	}
}

func newSummaryContext(yearNo int) calcctx.CalculationContext {
	group := entity.Group{
		ID:             1,
		BusinessStatus: enum.BusinessStatusNormal,
	}
	yearType := enum.YearTypeFormal
	if yearNo == 0 {
		yearType = enum.YearTypeDemo
	}
	yearState := entity.GroupYearState{
		GroupID:      1,
		YearNo:       yearNo,
		YearType:     yearType,
		YearStatus:   enum.YearStatusCompleted,
		StageStatus:  enum.StageStatusYearEndOpen,
		ReportStatus: enum.ReportStatusSubmitted,
	}
	gameConfig := entity.GameConfig{
		FinalYear:       8,
		CurrentOpenYear: yearNo,
	}

	return calcctx.NewCalculationContext(group, yearState, gameConfig)
}
