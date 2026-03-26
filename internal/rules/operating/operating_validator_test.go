package operating

import (
	"testing"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
	"sandbox-game/internal/state"
)

func TestValidateStageSubmitRejectsMissingQ1Scopes(t *testing.T) {
	t.Parallel()

	validator := NewValidator()
	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Beginning.TaxAndPlanning = map[string]any{
		"plannedRevenue": 100.0,
	}

	ctx := newOperatingValidationContext(enum.StageStatusQ1Open).
		WithInitialBaseline(&payload.BaselinePayload{
			BaselineCash: 1,
		}).
		WithOperatingPayload(&operatingPayload)

	result := validator.ValidateStageSubmit(ctx, state.StageCodeQ1)
	if result.Passed {
		t.Fatal("expected validation to fail when required q1 scopes are missing")
	}

	assertHasIssueField(t, result.Issues, "beginning.marketBid")
	assertHasIssueField(t, result.Issues, "quarter.shortTermLoan.q1")
	assertHasIssueField(t, result.Issues, "extra.incomeAndPenalty.q1")
}

func TestValidateStageSubmitAcceptsExplicitZeroValues(t *testing.T) {
	t.Parallel()

	validator := NewValidator()
	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Quarter.ShortTermLoan = payload.OperatingQuarterMap{
		"q2": {"dueRepayment": 0.0, "interest": 0.0, "newLoan": 0.0},
	}
	operatingPayload.Quarter.MaterialPayment = payload.OperatingQuarterMap{
		"q2": {"productA": 0.0, "productB": 0.0},
	}
	operatingPayload.Quarter.ProductionLineAdjust = payload.OperatingQuarterMap{
		"q2": {"changeProduct": 0.0, "dismantleCost": 0.0, "lineSale": 0.0},
	}
	operatingPayload.Quarter.HumanResource = payload.OperatingQuarterMap{
		"q2": {"staffCost": 0.0},
	}
	operatingPayload.Quarter.SalaryAndProduction = payload.OperatingQuarterMap{
		"q2": {"salaryCost": 0.0},
	}
	operatingPayload.Quarter.ResearchAndManagement = payload.OperatingQuarterMap{
		"q2": {"technologyResearch": 0.0, "managementSystem": 0.0},
	}
	operatingPayload.Quarter.ReceivableUpdate = payload.OperatingQuarterMap{
		"q2": {"receivableCollection": 0.0},
	}
	operatingPayload.Quarter.DeliverySettlement = payload.OperatingQuarterMap{
		"q2": {"salesRevenue": 0.0, "directCost": 0.0, "managementStaffCost": 0.0},
	}
	operatingPayload.Extra.IncomeAndPenalty = payload.OperatingQuarterMap{
		"q2": {"discountExpense": 0.0, "extraExpensePenalty": 0.0, "extraIncomeReward": 0.0},
	}

	ctx := newOperatingValidationContext(enum.StageStatusQ2Open).
		WithPreviousReport(&payload.ReportComputedPayload{
			ReportCash: 1,
		}).
		WithOperatingPayload(&operatingPayload)

	result := validator.ValidateStageSubmit(ctx, state.StageCodeQ2)
	if !result.Passed {
		t.Fatalf("expected explicit zero values to pass, got %#v", result.Issues)
	}
}

func TestValidateStageSubmitYearEndOnlyChecksYearEndScopes(t *testing.T) {
	t.Parallel()

	validator := NewValidator()
	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.YearEnd.LongTermLoan = map[string]any{
		"interest":  0.0,
		"repayment": 0.0,
		"newLoan":   0.0,
	}
	operatingPayload.YearEnd.AssetAdjustment = map[string]any{
		"lineMaintenance":    0.0,
		"purchase":           0.0,
		"sale":               0.0,
		"rent":               0.0,
		"workInConstruction": 0.0,
		"marketCultivation":  0.0,
	}

	ctx := newOperatingValidationContext(enum.StageStatusYearEndOpen).
		WithPreviousReport(&payload.ReportComputedPayload{
			ReportCash: 1,
		}).
		WithOperatingPayload(&operatingPayload)

	result := validator.ValidateStageSubmit(ctx, state.StageCodeYearEnd)
	if !result.Passed {
		t.Fatalf("expected year end validation to ignore prior quarter scopes, got %#v", result.Issues)
	}
}

func TestValidateStageSubmitRejectsBlankLeafValue(t *testing.T) {
	t.Parallel()

	validator := NewValidator()
	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Quarter.ShortTermLoan = payload.OperatingQuarterMap{
		"q3": {"dueRepayment": "", "interest": 0.0, "newLoan": 0.0},
	}
	operatingPayload.Quarter.MaterialPayment = payload.OperatingQuarterMap{
		"q3": {"productA": 0.0},
	}
	operatingPayload.Quarter.ProductionLineAdjust = payload.OperatingQuarterMap{
		"q3": {"changeProduct": 0.0},
	}
	operatingPayload.Quarter.HumanResource = payload.OperatingQuarterMap{
		"q3": {"staffCost": 0.0},
	}
	operatingPayload.Quarter.SalaryAndProduction = payload.OperatingQuarterMap{
		"q3": {"salaryCost": 0.0},
	}
	operatingPayload.Quarter.ResearchAndManagement = payload.OperatingQuarterMap{
		"q3": {"technologyResearch": 0.0},
	}
	operatingPayload.Quarter.ReceivableUpdate = payload.OperatingQuarterMap{
		"q3": {"receivableCollection": 0.0},
	}
	operatingPayload.Quarter.DeliverySettlement = payload.OperatingQuarterMap{
		"q3": {"salesRevenue": 0.0},
	}
	operatingPayload.Extra.IncomeAndPenalty = payload.OperatingQuarterMap{
		"q3": {"discountExpense": 0.0},
	}

	ctx := newOperatingValidationContext(enum.StageStatusQ3Open).
		WithPreviousReport(&payload.ReportComputedPayload{
			ReportCash: 1,
		}).
		WithOperatingPayload(&operatingPayload)

	result := validator.ValidateStageSubmit(ctx, state.StageCodeQ3)
	if result.Passed {
		t.Fatal("expected validation to fail when a leaf value is blank")
	}

	assertHasIssueField(t, result.Issues, "quarter.shortTermLoan.q3.dueRepayment")
}

func newOperatingValidationContext(stageStatus string) calcctx.CalculationContext {
	group := entity.Group{
		ID:             1,
		BusinessStatus: enum.BusinessStatusNormal,
	}
	yearState := entity.GroupYearState{
		GroupID:      1,
		YearNo:       1,
		YearType:     enum.YearTypeFormal,
		YearStatus:   enum.YearStatusOperating,
		StageStatus:  stageStatus,
		ReportStatus: enum.ReportStatusLocked,
	}
	if stageStatus == enum.StageStatusQ1Open {
		yearState.YearNo = 0
		yearState.YearType = enum.YearTypeDemo
	}
	gameConfig := entity.GameConfig{
		FinalYear:       8,
		CurrentOpenYear: yearState.YearNo,
	}

	return calcctx.NewCalculationContext(group, yearState, gameConfig)
}

func assertHasIssueField(t *testing.T, issues []ValidationIssue, field string) {
	t.Helper()

	for _, issue := range issues {
		if issue.Field == field {
			return
		}
	}

	t.Fatalf("expected issue field %s, got %#v", field, issues)
}
