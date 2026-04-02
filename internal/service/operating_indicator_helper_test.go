package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	calcctx "sandbox-game/internal/rules/context"
	operatingrules "sandbox-game/internal/rules/operating"
	reportrules "sandbox-game/internal/rules/report"
)

func TestPlayerOperatingViewIncludesReportPreviewIndicatorsWhenReportDraftExists(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureGameConfigExists(t, ctx, tx)

	groupID := createDemoReportFixtures(t, ctx, tx)
	_, queryService := buildPlayerOperatingServices(tx)

	operatingPayload := buildIndicatorPreviewOperatingPayload()
	createOperatingDraftRecord(t, ctx, tx, groupID, 0, enum.StageStatusYearEndOpen, operatingPayload)

	reportManual := buildBalancedReportManualPayload()
	createReportDraftRecord(t, ctx, tx, groupID, 0, reportManual)

	view, err := queryService.GetYearView(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("load operating view: %v", err)
	}

	expected := buildExpectedPreviewReportIndicators(t, ctx, tx, groupID, 0, operatingPayload, reportManual)

	for key, want := range expected {
		got, ok := view.DerivedValues[key]
		if !ok {
			t.Fatalf("expected derived value %s to exist", key)
		}
		assertFloatAlmostEquals(t, key, got, want)
	}
}

func TestPlayerOperatingViewIncludesReportIndicatorsWhenReportDraftMissing(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureGameConfigExists(t, ctx, tx)

	groupID := createDemoReportFixtures(t, ctx, tx)
	_, queryService := buildPlayerOperatingServices(tx)

	operatingPayload := buildIndicatorPreviewOperatingPayload()
	createOperatingDraftRecord(t, ctx, tx, groupID, 0, enum.StageStatusYearEndOpen, operatingPayload)

	view, err := queryService.GetYearView(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("load operating view: %v", err)
	}

	expected := buildExpectedPreviewReportIndicators(t, ctx, tx, groupID, 0, operatingPayload, payload.ReportManualPayload{})

	for key, want := range expected {
		got, ok := view.DerivedValues[key]
		if !ok {
			t.Fatalf("expected derived value %s to exist when report draft is missing", key)
		}
		assertFloatAlmostEquals(t, key, got, want)
	}
}

func TestPlayerOperatingViewUsesSubmittedReportIndicatorsFirst(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureGameConfigExists(t, ctx, tx)

	groupID := createDemoReportFixtures(t, ctx, tx)
	_, queryService := buildPlayerOperatingServices(tx)

	operatingPayload := buildIndicatorPreviewOperatingPayload()
	createOperatingDraftRecord(t, ctx, tx, groupID, 0, enum.StageStatusYearEndOpen, operatingPayload)

	submittedComputed := payload.ReportComputedPayload{
		ReportSalesRevenue: 100,
		ReportDirectCost:   20,
		ReportNetProfit:    30,
		ReportTotalEquity:  60,
		ReportTotalAssets:  120,
	}
	createSubmittedReportRecord(t, ctx, tx, groupID, 0, submittedComputed)

	view, err := queryService.GetYearView(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("load operating view: %v", err)
	}

	expected := operatingrules.BuildReportIndicatorValues(submittedComputed)
	for key, want := range expected {
		got, ok := view.DerivedValues[key]
		if !ok {
			t.Fatalf("expected derived value %s to exist for submitted report", key)
		}
		assertFloatAlmostEquals(t, key, got, want)
	}
}

func TestPlayerOperatingViewIncludesFormalYearReportIndicatorsUsingPreviousReportCarryForward(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureGameConfigExists(t, ctx, tx)

	now := time.Now()
	groupID := createIntegrationGroup(t, ctx, tx, now, "IT_OPERATING_FORMAL_INDICATOR")
	createGroupYearStateRecord(t, ctx, tx, groupID, 0, enum.YearTypeDemo, enum.YearStatusCompleted, enum.StageStatusYearEndOpen, enum.ReportStatusSubmitted)
	createPreviousFormalReportRecord(t, ctx, tx, groupID, 0, now)
	createGroupYearStateRecord(t, ctx, tx, groupID, 1, enum.YearTypeFormal, enum.YearStatusOperating, enum.StageStatusYearEndOpen, enum.ReportStatusLocked)

	_, queryService := buildPlayerOperatingServices(tx)

	operatingPayload := buildIndicatorPreviewOperatingPayload()
	createOperatingDraftRecord(t, ctx, tx, groupID, 1, enum.StageStatusYearEndOpen, operatingPayload)

	view, err := queryService.GetYearView(ctx, groupID, 1)
	if err != nil {
		t.Fatalf("load formal year operating view: %v", err)
	}

	expected := buildExpectedPreviewReportIndicators(t, ctx, tx, groupID, 1, operatingPayload, payload.ReportManualPayload{})

	for key, want := range expected {
		got, ok := view.DerivedValues[key]
		if !ok {
			t.Fatalf("expected formal year derived value %s to exist", key)
		}
		assertFloatAlmostEquals(t, key, got, want)
	}
}

func buildIndicatorPreviewOperatingPayload() payload.OperatingPayload {
	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.Beginning.TaxAndPlanning = map[string]any{
		"marketInvestmentTotal": 1.0,
		"orderAmountTotal":      11.0,
	}
	operatingPayload.Quarter.ShortTermLoan = payload.OperatingQuarterMap{
		"Q1": {"interest": 1.0, "newLoan": 20.0},
	}
	operatingPayload.Quarter.MaterialPayment = payload.OperatingQuarterMap{
		"Q1": {"materialCost": 20.0},
	}
	operatingPayload.Quarter.ProductionLineAdjust = payload.OperatingQuarterMap{
		"Q1": {
			"changeProduct":       2.0,
			"dismantleCost":       1.0,
			"lineSale":            5.0,
			"newLineInstall":      2.0,
			"constructionToFixed": 10.0,
			"newDepreciableAsset": 19.0,
		},
	}
	operatingPayload.Quarter.HumanResource = payload.OperatingQuarterMap{
		"Q1": {"staffCost": 4.0},
	}
	operatingPayload.Quarter.SalaryAndProduction = payload.OperatingQuarterMap{
		"Q1": {"salaryCost": 14.0},
	}
	operatingPayload.Quarter.ResearchAndManagement = payload.OperatingQuarterMap{
		"Q1": {"technologyResearch": 5.0, "managementSystem": 6.0},
	}
	operatingPayload.Quarter.ReceivableUpdate = payload.OperatingQuarterMap{
		"Q1": {"receivableCollection": 11.0},
	}
	operatingPayload.Quarter.DeliverySettlement = payload.OperatingQuarterMap{
		"Q1": {"salesRevenue": 11.0, "directCost": 4.0},
	}
	operatingPayload.YearEnd.LongTermLoan = map[string]any{
		"interest": 1.0,
		"newLoan":  140.0,
	}
	operatingPayload.YearEnd.AssetAdjustment = map[string]any{
		"lineMaintenance":   1.0,
		"purchase":          20.0,
		"sale":              15.0,
		"rent":              1.0,
		"marketCultivation": 1.0,
	}

	return operatingPayload.Normalize().WithoutDerivedValues()
}

func createOperatingDraftRecord(
	t *testing.T,
	ctx context.Context,
	tx *gorm.DB,
	groupID int64,
	yearNo int,
	stageStatus string,
	operatingPayload payload.OperatingPayload,
) {
	t.Helper()

	rawPayload, err := json.Marshal(operatingPayload.Normalize().WithoutDerivedValues())
	if err != nil {
		t.Fatalf("marshal operating payload: %v", err)
	}

	now := time.Now()
	item := entity.GroupOperatingDraft{
		GroupID:          groupID,
		YearNo:           yearNo,
		StageStatus:      stageStatus,
		OperatingPayload: rawPayload,
		LastAutoSavedAt:  &now,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		t.Fatalf("create operating draft fixture: %v", err)
	}
}

func createReportDraftRecord(
	t *testing.T,
	ctx context.Context,
	tx *gorm.DB,
	groupID int64,
	yearNo int,
	manualPayload payload.ReportManualPayload,
) {
	t.Helper()

	rawPayload, err := json.Marshal(manualPayload)
	if err != nil {
		t.Fatalf("marshal report manual payload: %v", err)
	}

	now := time.Now()
	item := entity.GroupReport{
		GroupID:               groupID,
		YearNo:                yearNo,
		ReportManualPayload:   rawPayload,
		ReportComputedPayload: []byte("{}"),
		LastAutoSavedAt:       &now,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		t.Fatalf("create report draft fixture: %v", err)
	}
}

func createSubmittedReportRecord(
	t *testing.T,
	ctx context.Context,
	tx *gorm.DB,
	groupID int64,
	yearNo int,
	computedPayload payload.ReportComputedPayload,
) {
	t.Helper()

	rawComputed, err := json.Marshal(computedPayload)
	if err != nil {
		t.Fatalf("marshal report computed payload: %v", err)
	}

	now := time.Now()
	item := entity.GroupReport{
		GroupID:               groupID,
		YearNo:                yearNo,
		ReportManualPayload:   []byte("{}"),
		ReportComputedPayload: rawComputed,
		BalanceCheckPassed:    true,
		LastAutoSavedAt:       &now,
		SubmittedAt:           &now,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		t.Fatalf("create submitted report fixture: %v", err)
	}
}

func reportIndicatorKeys() []string {
	return []string{
		"netAssetYield",
		"totalAssetYield",
		"netProfitRate",
		"grossMarginRate",
	}
}

func assertFloatAlmostEquals(t *testing.T, label string, got float64, want float64) {
	t.Helper()

	const tolerance = 1e-9
	diff := got - want
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		t.Fatalf("expected %s=%f, got %f", label, want, got)
	}
}

func buildExpectedPreviewReportIndicators(
	t *testing.T,
	ctx context.Context,
	tx *gorm.DB,
	groupID int64,
	yearNo int,
	operatingPayload payload.OperatingPayload,
	reportManual payload.ReportManualPayload,
) map[string]float64 {
	t.Helper()

	groupRepo := repository.NewGroupRepository(tx)
	gameConfigRepo := repository.NewGameConfigRepository(tx)
	groupYearRepo := repository.NewGroupYearStateRepository(tx)
	initialBaselineRepo := repository.NewInitialBaselineRepository(tx)
	reportRepo := repository.NewReportRepository(tx)

	group, err := groupRepo.GetByID(ctx, groupID)
	if err != nil {
		t.Fatalf("load group fixture: %v", err)
	}
	gameConfig, err := gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		t.Fatalf("load game config fixture: %v", err)
	}
	yearState, err := groupYearRepo.GetByGroupIDAndYear(ctx, groupID, yearNo)
	if err != nil {
		t.Fatalf("load year state fixture: %v", err)
	}
	calculationContext := calcctx.NewCalculationContext(*group, *yearState, *gameConfig).
		WithOperatingPayload(&operatingPayload).
		WithReportManualPayload(&reportManual)

	if yearNo == 0 {
		initialBaseline, err := initialBaselineRepo.FindByGroupID(ctx, groupID)
		if err != nil {
			t.Fatalf("load initial baseline fixture: %v", err)
		}

		var baselinePayload payload.BaselinePayload
		if err := json.Unmarshal(initialBaseline.BaselinePayload, &baselinePayload); err != nil {
			t.Fatalf("unmarshal initial baseline payload: %v", err)
		}

		calculationContext = calculationContext.WithInitialBaseline(&baselinePayload)
	} else {
		previousReport, err := reportRepo.FindEffectiveByGroupIDAndYear(ctx, groupID, yearNo-1)
		if err != nil {
			t.Fatalf("load previous report fixture: %v", err)
		}

		var previousReportPayload payload.ReportComputedPayload
		if err := json.Unmarshal(previousReport.ReportComputedPayload, &previousReportPayload); err != nil {
			t.Fatalf("unmarshal previous report payload: %v", err)
		}

		calculationContext = calculationContext.WithPreviousReport(&previousReportPayload)
	}

	previewPayload, err := reportrules.NewCalculator().Calculate(calculationContext)
	if err != nil {
		t.Fatalf("calculate report preview: %v", err)
	}

	return operatingrules.BuildReportIndicatorValues(previewPayload)
}
