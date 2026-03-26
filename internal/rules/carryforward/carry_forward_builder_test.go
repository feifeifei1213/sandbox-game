package carryforward

import (
	"testing"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
)

func TestBuildUsesInitialBaselineForDemoYear(t *testing.T) {
	t.Parallel()

	builder := NewBuilder()
	ctx := newCarryForwardContext(0).WithInitialBaseline(&payload.BaselinePayload{
		BaselineSalesRevenue:         100,
		BaselineDirectCost:           60,
		BaselineComprehensiveCost:    10,
		BaselineDepreciation:         9,
		BaselineFinanceIncomeExpense: 2,
		BaselineExtraIncomeExpense:   1,
		BaselineIncomeTax:            1,
		BaselineShortTermLoan:        20,
		BaselineLongTermLoan:         3,
		BaselineLineResidual:         4,
		BaselineDepreciableAsset:     5,
		BaselineCash:                 6,
		BaselineReceivable:           7,
		BaselineShareCapital:         8,
		BaselineRetainedEarnings:     9,
	})

	result, err := builder.Build(ctx)
	if err != nil {
		t.Fatalf("build carry forward failed: %v", err)
	}

	if result.PreviousCash != 6 {
		t.Fatalf("expected previous cash 6, got %v", result.PreviousCash)
	}
	if result.RetainedEarnings != 28 {
		t.Fatalf("expected retained earnings 28, got %v", result.RetainedEarnings)
	}
}

func TestBuildUsesPreviousReportForFormalYear(t *testing.T) {
	t.Parallel()

	builder := NewBuilder()
	ctx := newCarryForwardContext(1).WithPreviousReport(&payload.ReportComputedPayload{
		ReportIncomeTax:          1,
		ReportShortTermLiability: 2,
		ReportLongTermLiability:  3,
		ReportLineResidual:       4,
		ReportDepreciableAsset:   5,
		ReportCash:               6,
		ReportReceivable:         7,
		ReportShareCapital:       8,
		ReportRetainedEarnings:   9,
		ReportNetProfit:          10,
	})

	result, err := builder.Build(ctx)
	if err != nil {
		t.Fatalf("build carry forward failed: %v", err)
	}

	if result.RetainedEarnings != 19 {
		t.Fatalf("expected retained earnings 19, got %v", result.RetainedEarnings)
	}
}

func newCarryForwardContext(yearNo int) calcctx.CalculationContext {
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
		YearStatus:   enum.YearStatusOperating,
		StageStatus:  enum.StageStatusQ1Open,
		ReportStatus: enum.ReportStatusLocked,
	}
	gameConfig := entity.GameConfig{
		FinalYear:       8,
		CurrentOpenYear: yearNo,
	}

	return calcctx.NewCalculationContext(group, yearState, gameConfig)
}
