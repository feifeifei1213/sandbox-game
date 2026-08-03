package service

import (
	"errors"
	"strings"
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

func TestValidateOperatingManualIntegersAllowsProjectProgressMetadata(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.YearEnd.ProjectProgressUpdate.Items[0]["factoryKey"] = "factoryA"
	operatingPayload.YearEnd.ProjectProgressUpdate.Items[0]["factoryLabel"] = "生产厂房 A"
	operatingPayload.YearEnd.ProjectProgressUpdate.Items[0]["quarterKey"] = "q1"
	operatingPayload.YearEnd.ProjectProgressUpdate.Items[0]["quarterLabel"] = "第一季度"
	operatingPayload.YearEnd.ProjectProgressUpdate.Items[0]["lineType"] = "人工"
	operatingPayload.YearEnd.ProjectProgressUpdate.Items[0]["progress"] = 4

	if err := validateOperatingManualIntegers(operatingPayload); err != nil {
		t.Fatalf("expected project progress metadata to be ignored by integer validation, got %v", err)
	}
}

func TestValidateOperatingManualIntegersRejectsInvalidProjectProgressLineType(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.YearEnd.ProjectProgressUpdate.Items[0]["lineType"] = "无人"

	err := validateOperatingManualIntegers(operatingPayload)
	if err == nil || !strings.Contains(err.Error(), "项目进度更新") {
		t.Fatalf("expected project progress line type validation error, got %v", err)
	}
}

func TestValidateOperatingManualIntegersRejectsInvalidProjectProgressValue(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.YearEnd.ProjectProgressUpdate.Items[0]["progress"] = 5

	err := validateOperatingManualIntegers(operatingPayload)
	if err == nil || !strings.Contains(err.Error(), "进度只能填写 0~4") {
		t.Fatalf("expected project progress value validation error, got %v", err)
	}
}

func TestValidateOperatingManualIntegersRejectsNegativeSupplyChainOrderQuantity(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.Quarter.SupplyChainOrderRecord["q1"]["basicProduct"] = -1

	err := validateOperatingManualIntegers(operatingPayload)
	if !errors.Is(err, ErrSupplyChainOrderQuantityInvalid) {
		t.Fatalf("expected ErrSupplyChainOrderQuantityInvalid, got %v", err)
	}
}

func TestValidateOperatingManualIntegersRejectsTextSupplyChainOrderQuantity(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.Quarter.SupplyChainOrderRecord["q1"]["basicProduct"] = "无"

	err := validateOperatingManualIntegers(operatingPayload)
	if !errors.Is(err, ErrSupplyChainOrderQuantityInvalid) {
		t.Fatalf("expected ErrSupplyChainOrderQuantityInvalid, got %v", err)
	}
}

func TestValidateOperatingManualIntegersRejectsInvalidMarketCultivationAnnualInvestment(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.YearEnd.MarketCultivation.Regional.AnnualInvestment = 2

	err := validateOperatingManualIntegers(operatingPayload)
	if err == nil || !strings.Contains(err.Error(), "新市场培育") {
		t.Fatalf("expected market cultivation validation error, got %v", err)
	}
}

func TestValidateOperatingManualIntegersRejectsLockedMarketCultivationInvestment(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.YearEnd.MarketCultivation.Regional.AnnualInvestment = 1
	operatingPayload.YearEnd.MarketCultivation.Regional.LockedByPrevious = true

	err := validateOperatingManualIntegers(operatingPayload)
	if err == nil || !strings.Contains(err.Error(), "已解锁") {
		t.Fatalf("expected locked market cultivation validation error, got %v", err)
	}
}

func TestValidateOperatingManualIntegersRejectsNonChineseQualificationStatus(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.YearEnd.QualificationCertification.HighTechEnterprise.Status = "unlocked"

	err := validateOperatingManualIntegers(operatingPayload)
	if err == nil || !strings.Contains(err.Error(), "资质认证") {
		t.Fatalf("expected qualification validation error, got %v", err)
	}
}

func TestValidateOperatingManualIntegersAllowsQualificationChineseStatus(t *testing.T) {
	operatingPayload := buildValidQ1OperatingPayload()
	operatingPayload.YearEnd.QualificationCertification.HighTechEnterprise.Status = "解锁"

	if err := validateOperatingManualIntegers(operatingPayload); err != nil {
		t.Fatalf("expected Chinese qualification status to pass, got %v", err)
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
