package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	calcctx "sandbox-game/internal/rules/context"
	operatingrules "sandbox-game/internal/rules/operating"
	reportrules "sandbox-game/internal/rules/report"
)

func enrichOperatingDerivedValuesWithReportMetrics(
	ctx context.Context,
	reportRepo *repository.ReportRepository,
	reportCalculator *reportrules.Calculator,
	calculationContext calcctx.CalculationContext,
	groupID int64,
	yearNo int,
	result *operatingrules.CalculationResult,
) error {
	reportPayload, ok, err := loadCurrentYearReportPreviewOrResult(ctx, reportRepo, reportCalculator, calculationContext, groupID, yearNo)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	for key, value := range operatingrules.BuildReportIndicatorValues(reportPayload) {
		result.DerivedValues[key] = value
		result.OperatingPayload.Derived.Values[key] = value
	}
	return nil
}

func loadCurrentYearReportPreviewOrResult(
	ctx context.Context,
	reportRepo *repository.ReportRepository,
	reportCalculator *reportrules.Calculator,
	calculationContext calcctx.CalculationContext,
	groupID int64,
	yearNo int,
) (payload.ReportComputedPayload, bool, error) {
	manualPayload := payload.ReportManualPayload{}
	computedPayload := payload.ReportComputedPayload{}

	report, err := reportRepo.FindByGroupIDAndYear(ctx, groupID, yearNo)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return calculateReportPreviewPayload(reportCalculator, calculationContext, manualPayload)
	case err != nil:
		return payload.ReportComputedPayload{}, false, fmt.Errorf("load group report: %w", err)
	}

	if len(report.ReportManualPayload) > 0 {
		if unmarshalErr := json.Unmarshal(report.ReportManualPayload, &manualPayload); unmarshalErr != nil {
			return payload.ReportComputedPayload{}, false, fmt.Errorf("unmarshal report manual payload: %w", unmarshalErr)
		}
	}

	if len(report.ReportComputedPayload) > 0 {
		if unmarshalErr := json.Unmarshal(report.ReportComputedPayload, &computedPayload); unmarshalErr != nil {
			return payload.ReportComputedPayload{}, false, fmt.Errorf("unmarshal report computed payload: %w", unmarshalErr)
		}
	}

	if report.SubmittedAt != nil && !isZeroReportComputedPayload(computedPayload) {
		return computedPayload, true, nil
	}

	return calculateReportPreviewPayload(reportCalculator, calculationContext, manualPayload)
}

func calculateReportPreviewPayload(
	reportCalculator *reportrules.Calculator,
	calculationContext calcctx.CalculationContext,
	manualPayload payload.ReportManualPayload,
) (payload.ReportComputedPayload, bool, error) {
	reportContext := calculationContext.WithReportManualPayload(&manualPayload)
	previewPayload, calcErr := reportCalculator.Calculate(reportContext)
	if calcErr != nil {
		return payload.ReportComputedPayload{}, false, fmt.Errorf("calculate report preview: %w", calcErr)
	}
	return previewPayload, true, nil
}

func isZeroReportComputedPayload(value payload.ReportComputedPayload) bool {
	return value == (payload.ReportComputedPayload{})
}
