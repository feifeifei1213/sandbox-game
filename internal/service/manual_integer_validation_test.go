package service

import (
	"errors"
	"testing"

	"sandbox-game/internal/model/payload"
)

func TestValidateOperatingManualIntegersRejectsDecimalInput(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.Quarter.ShortTermLoan["q1"]["newShortTermLoan"] = 1.5

	err := validateOperatingManualIntegers(operatingPayload)
	if !errors.Is(err, ErrManualNumberNotInteger) {
		t.Fatalf("expected ErrManualNumberNotInteger, got %v", err)
	}
}

func TestValidateOperatingManualIntegersIgnoresDerivedDecimalValues(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.Derived.Values["netProfitRate"] = 0.1234

	if err := validateOperatingManualIntegers(operatingPayload); err != nil {
		t.Fatalf("expected derived decimal to be ignored, got %v", err)
	}
}

func TestValidateOperatingManualIntegersIgnoresMarketMetadata(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.Beginning.MarketBid = []map[string]any{
		{
			"marketCode":              "LOCAL",
			"basicProductTotal":       1.0,
			"standardProductTotal":    2.0,
			"precisionProductTotal":   3.0,
			"intelligentProductTotal": 4.0,
			"orderAmount":             10.0,
		},
	}

	if err := validateOperatingManualIntegers(operatingPayload); err != nil {
		t.Fatalf("expected market metadata to be ignored, got %v", err)
	}
}

func TestValidateReportManualIntegersAllowsIncomeTaxRateDecimal(t *testing.T) {
	manual := payload.ReportManualPayload{
		WorkInProgress:               float64Ptr(6),
		FinishedGoods:                float64Ptr(4),
		RawMaterials:                 float64Ptr(1),
		IncomeTaxRate:                float64Ptr(0.25),
		EnterpriseCertificationScore: float64Ptr(0),
		ProductionHumanScore:         float64Ptr(0),
		ClosingSpeedScore:            float64Ptr(0),
	}

	if err := validateReportManualIntegers(manual); err != nil {
		t.Fatalf("expected tax rate decimal to be ignored, got %v", err)
	}
}

func TestValidateReportManualIntegersRejectsManualScoreDecimal(t *testing.T) {
	manual := payload.ReportManualPayload{
		WorkInProgress:               float64Ptr(6),
		FinishedGoods:                float64Ptr(4),
		RawMaterials:                 float64Ptr(1),
		IncomeTaxRate:                float64Ptr(0.25),
		EnterpriseCertificationScore: float64Ptr(0.5),
		ProductionHumanScore:         float64Ptr(0),
		ClosingSpeedScore:            float64Ptr(0),
	}

	err := validateReportManualIntegers(manual)
	if !errors.Is(err, ErrManualNumberNotInteger) {
		t.Fatalf("expected ErrManualNumberNotInteger, got %v", err)
	}
}

func TestValidateBaselineManualIntegersRejectsDecimalInput(t *testing.T) {
	baselinePayload := buildIntegrationBaselinePayload()
	baselinePayload.BaselineCash = 36.5

	err := validateBaselineManualIntegers(baselinePayload)
	if !errors.Is(err, ErrManualNumberNotInteger) {
		t.Fatalf("expected ErrManualNumberNotInteger, got %v", err)
	}
}
