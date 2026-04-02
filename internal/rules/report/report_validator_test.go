package report

import (
	"testing"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
)

func TestValidateSubmitRejectsMissingManualFields(t *testing.T) {
	t.Parallel()

	validator := NewValidator()
	ctx := newReportContext().WithReportManualPayload(&payload.ReportManualPayload{})

	result := validator.ValidateSubmit(ctx, payload.ReportComputedPayload{
		ReportTotalAssets:          100,
		ReportTotalLiabilityEquity: 100,
	})
	if result.Passed {
		t.Fatal("expected validation to fail when manual fields are missing")
	}
}

func TestValidateSubmitRejectsIllegalTaxRate(t *testing.T) {
	t.Parallel()

	workInProgress := 1.0
	finishedGoods := 2.0
	rawMaterials := 3.0
	incomeTaxRate := 0.2
	enterpriseCertificationScore := 0.0
	productionHumanScore := 0.0
	closingSpeedScore := 0.0

	validator := NewValidator()
	ctx := newReportContext().WithReportManualPayload(&payload.ReportManualPayload{
		WorkInProgress:               &workInProgress,
		FinishedGoods:                &finishedGoods,
		RawMaterials:                 &rawMaterials,
		IncomeTaxRate:                &incomeTaxRate,
		EnterpriseCertificationScore: &enterpriseCertificationScore,
		ProductionHumanScore:         &productionHumanScore,
		ClosingSpeedScore:            &closingSpeedScore,
	})

	result := validator.ValidateSubmit(ctx, payload.ReportComputedPayload{
		ReportTotalAssets:          100,
		ReportTotalLiabilityEquity: 100,
	})
	if result.Passed {
		t.Fatal("expected validation to fail when tax rate is illegal")
	}
}

func TestValidateSubmitAcceptsBalancedReport(t *testing.T) {
	t.Parallel()

	workInProgress := 0.0
	finishedGoods := 0.0
	rawMaterials := 0.0
	incomeTaxRate := 0.25
	enterpriseCertificationScore := 0.0
	productionHumanScore := 0.0
	closingSpeedScore := 0.0

	validator := NewValidator()
	ctx := newReportContext().WithReportManualPayload(&payload.ReportManualPayload{
		WorkInProgress:               &workInProgress,
		FinishedGoods:                &finishedGoods,
		RawMaterials:                 &rawMaterials,
		IncomeTaxRate:                &incomeTaxRate,
		EnterpriseCertificationScore: &enterpriseCertificationScore,
		ProductionHumanScore:         &productionHumanScore,
		ClosingSpeedScore:            &closingSpeedScore,
	})

	result := validator.ValidateSubmit(ctx, payload.ReportComputedPayload{
		ReportTotalAssets:          100,
		ReportTotalLiabilityEquity: 100,
	})
	if !result.Passed {
		t.Fatalf("expected validation to pass, got %#v", result.Issues)
	}
}

func TestValidateDraftAllowsPartialManualFields(t *testing.T) {
	t.Parallel()

	incomeTaxRate := 0.25
	validator := NewValidator()

	result := validator.ValidateDraft(&payload.ReportManualPayload{
		IncomeTaxRate: &incomeTaxRate,
	})
	if !result.Passed {
		t.Fatalf("expected draft validation to pass, got %#v", result.Issues)
	}
}

func TestValidateDraftRejectsIllegalTaxRate(t *testing.T) {
	t.Parallel()

	incomeTaxRate := 0.2
	validator := NewValidator()

	result := validator.ValidateDraft(&payload.ReportManualPayload{
		IncomeTaxRate: &incomeTaxRate,
	})
	if result.Passed {
		t.Fatal("expected draft validation to fail when tax rate is illegal")
	}
}

func newReportContext() calcctx.CalculationContext {
	group := entity.Group{
		ID:             1,
		BusinessStatus: enum.BusinessStatusNormal,
	}
	yearState := entity.GroupYearState{
		GroupID:      1,
		YearNo:       1,
		YearType:     enum.YearTypeFormal,
		YearStatus:   enum.YearStatusReportPending,
		StageStatus:  enum.StageStatusYearEndOpen,
		ReportStatus: enum.ReportStatusOpen,
	}
	gameConfig := entity.GameConfig{
		FinalYear:       8,
		CurrentOpenYear: 1,
	}

	return calcctx.NewCalculationContext(group, yearState, gameConfig)
}
