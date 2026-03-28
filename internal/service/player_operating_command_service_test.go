package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/assembler"
	appconfig "sandbox-game/internal/config"
	"sandbox-game/internal/enum"
	appdb "sandbox-game/internal/infra/db"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

func TestSubmitOperatingStageAdvancesEditableScopeFromQ1ToQ2(t *testing.T) {
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

	groupID := createIntegrationOperatingFixtures(t, ctx, tx)
	commandService, queryService := buildPlayerOperatingServices(tx)

	result, err := commandService.SubmitStage(ctx, SubmitOperatingStageCommand{
		GroupID:          groupID,
		YearNo:           0,
		StageCode:        state.StageCodeQ1,
		OperatingPayload: buildValidQ1OperatingPayload(),
		SubmitterID:      90001,
		OperatorName:     "integration-test",
	})
	if err != nil {
		t.Fatalf("submit stage: %v", err)
	}

	if result.StageCode != state.StageCodeQ1 {
		t.Fatalf("expected stage code Q1, got %s", result.StageCode)
	}
	if result.YearStatus != enum.YearStatusOperating {
		t.Fatalf("expected year status OPERATING, got %s", result.YearStatus)
	}
	if result.StageStatus != enum.StageStatusQ2Open {
		t.Fatalf("expected stage status Q2_OPEN, got %s", result.StageStatus)
	}
	if result.ReportStatus != enum.ReportStatusLocked {
		t.Fatalf("expected report status REPORT_LOCKED, got %s", result.ReportStatus)
	}
	if result.BusinessStatus != enum.BusinessStatusNormal {
		t.Fatalf("expected business status NORMAL, got %s", result.BusinessStatus)
	}
	if result.LatestStageSubmitVersion != 1 {
		t.Fatalf("expected latest stage submit version 1, got %d", result.LatestStageSubmitVersion)
	}
	if result.PeriodEndCash <= 0 {
		t.Fatalf("expected positive period end cash after Q1 submit, got %f", result.PeriodEndCash)
	}
	if result.SubmittedAt.IsZero() {
		t.Fatalf("expected submittedAt to be filled")
	}

	view, err := queryService.GetYearView(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("reload year view: %v", err)
	}

	if view.CurrentStageCode != state.StageCodeQ2 {
		t.Fatalf("expected current stage code Q2, got %s", view.CurrentStageCode)
	}
	if view.StageStatus != enum.StageStatusQ2Open {
		t.Fatalf("expected reloaded stage status Q2_OPEN, got %s", view.StageStatus)
	}
	if !view.CanEdit || !view.CanSubmit {
		t.Fatalf("expected reloaded view to remain editable/submittable")
	}
	if !slices.Equal(view.EditableScopes, []string{state.OperatingScopeQ2}) {
		t.Fatalf("expected editable scopes [Q2], got %#v", view.EditableScopes)
	}
	if !slices.Contains(view.ReadonlyScopes, state.OperatingScopeYearStart) || !slices.Contains(view.ReadonlyScopes, state.OperatingScopeQ1) {
		t.Fatalf("expected YEAR_START and Q1 to become readonly, got %#v", view.ReadonlyScopes)
	}
	if len(view.StageSubmitHistory) != 1 {
		t.Fatalf("expected one stage submission history item, got %d", len(view.StageSubmitHistory))
	}
	if view.StageSubmitHistory[0].StageCode != state.StageCodeQ1 || view.StageSubmitHistory[0].SubmitVersion != 1 {
		t.Fatalf("expected first history item to be Q1 version 1, got %#v", view.StageSubmitHistory[0])
	}
	if view.LastDraftSavedAt == nil || view.LastDraftSavedAt.IsZero() {
		t.Fatalf("expected last draft save time to be written on submit")
	}
}

func openIntegrationMySQL(t *testing.T) *gorm.DB {
	t.Helper()

	cfg, err := appconfig.Load(resolveProjectPath(t, "configs", "local.yaml"))
	if err != nil {
		t.Skipf("skip integration test: load local config failed: %v", err)
	}

	db, err := appdb.NewMySQL(cfg)
	if err != nil {
		t.Skipf("skip integration test: open mysql failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("extract sql db: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Skipf("skip integration test: ping mysql failed: %v", err)
	}

	return db
}

func buildPlayerOperatingServices(db *gorm.DB) (*PlayerOperatingCommandService, *PlayerOperatingQueryService) {
	gameConfigRepo := repository.NewGameConfigRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	groupYearRepo := repository.NewGroupYearStateRepository(db)
	operatingRepo := repository.NewOperatingRepository(db)
	initialBaselineRepo := repository.NewInitialBaselineRepository(db)
	reportRepo := repository.NewReportRepository(db)

	commandService := NewPlayerOperatingCommandService(
		db,
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaselineRepo,
		reportRepo,
	)
	queryService := NewPlayerOperatingQueryService(
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaselineRepo,
		reportRepo,
		assembler.NewPlayerOperatingAssembler(),
	)

	return commandService, queryService
}

func ensureGameConfigExists(t *testing.T, ctx context.Context, tx *gorm.DB) {
	t.Helper()

	var count int64
	if err := tx.WithContext(ctx).Model(&entity.GameConfig{}).Count(&count).Error; err != nil {
		t.Fatalf("count game config: %v", err)
	}
	if count > 0 {
		return
	}

	now := time.Now()
	item := entity.GameConfig{
		FinalYear:                8,
		CurrentOpenYear:          0,
		RuleVersion:              "integration-test",
		TemplateVersion:          "integration-test",
		InitialBaselineSubmitted: true,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		t.Fatalf("create fallback game config: %v", err)
	}
}

func createIntegrationOperatingFixtures(t *testing.T, ctx context.Context, tx *gorm.DB) int64 {
	t.Helper()

	now := time.Now()
	uniqueSeed := now.UnixNano()
	group := entity.Group{
		GroupNo:        int(900000 + uniqueSeed%100000),
		GroupCode:      fmt.Sprintf("IT_OPERATING_%d", uniqueSeed),
		GroupName:      "经营提交流程集成测试组",
		BusinessStatus: enum.BusinessStatusNormal,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&group).Error; err != nil {
		t.Fatalf("create group fixture: %v", err)
	}

	yearState := entity.GroupYearState{
		GroupID:                   group.ID,
		YearNo:                    0,
		YearType:                  enum.YearTypeDemo,
		YearStatus:                enum.YearStatusOperating,
		StageStatus:               enum.StageStatusQ1Open,
		ReportStatus:              enum.ReportStatusLocked,
		SummaryEffective:          false,
		LatestStageSubmitVersion:  0,
		LatestReportSubmitVersion: 0,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&yearState).Error; err != nil {
		t.Fatalf("create year state fixture: %v", err)
	}

	baselineJSON, err := json.Marshal(buildIntegrationBaselinePayload())
	if err != nil {
		t.Fatalf("marshal baseline payload: %v", err)
	}

	submittedAt := now
	baseline := entity.InitialBaseline{
		GroupID:         group.ID,
		BaselinePayload: baselineJSON,
		Submitted:       true,
		SubmitterID:     int64Ptr(1),
		SubmittedAt:     &submittedAt,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&baseline).Error; err != nil {
		t.Fatalf("create baseline fixture: %v", err)
	}

	return group.ID
}

func buildIntegrationBaselinePayload() payload.BaselinePayload {
	return payload.BaselinePayload{
		BaselineSalesRevenue:         32,
		BaselineDirectCost:           15,
		BaselineComprehensiveCost:    13,
		BaselineDepreciation:         1,
		BaselineFinanceIncomeExpense: 2,
		BaselineExtraIncomeExpense:   2,
		BaselineIncomeTax:            1,
		BaselineFactoryAsset:         40,
		BaselineLineResidual:         3,
		BaselineDepreciableAsset:     0,
		BaselineCash:                 36,
		BaselineReceivable:           0,
		BaselineWorkInProgress:       6,
		BaselineFinishedGoods:        4,
		BaselineRawMaterials:         1,
		BaselineShortTermLoan:        20,
		BaselineLongTermLoan:         0,
		BaselineShareCapital:         50,
		BaselineRetainedEarnings:     17,
	}
}

func buildValidQ1OperatingPayload() payload.OperatingPayload {
	p := payload.NewOperatingPayload()

	p.Beginning.TaxAndPlanning = map[string]any{
		"taxRate":       0,
		"marketBidCost": 0,
	}
	p.Beginning.MarketBid = []map[string]any{
		{
			"market":           "demo",
			"marketInvestment": 0,
			"orderAmount":      0,
		},
	}

	p.Quarter.ShortTermLoan = payload.OperatingQuarterMap{
		"q1": {
			"shortTermRepayment": 0,
			"shortTermInterest":  0,
			"newShortTermLoan":   0,
		},
	}
	p.Quarter.MaterialPayment = payload.OperatingQuarterMap{
		"q1": {
			"materialPayment": 0,
		},
	}
	p.Quarter.ProductionLineAdjust = payload.OperatingQuarterMap{
		"q1": {
			"changeProductCost":        0,
			"lineDismantleCost":        0,
			"lineSaleValue":            0,
			"newLineInstall":           0,
			"constructionToFixed":      0,
			"depreciableAssetIncrease": 0,
		},
	}
	p.Quarter.HumanResource = payload.OperatingQuarterMap{
		"q1": {
			"humanResourceCost": 0,
		},
	}
	p.Quarter.SalaryAndProduction = payload.OperatingQuarterMap{
		"q1": {
			"salaryAndProductionCost": 0,
		},
	}
	p.Quarter.ResearchAndManagement = payload.OperatingQuarterMap{
		"q1": {
			"researchCost":         0,
			"managementSystemCost": 0,
		},
	}
	p.Quarter.ReceivableUpdate = payload.OperatingQuarterMap{
		"q1": {
			"receivableRecovered": 0,
		},
	}
	p.Quarter.DeliverySettlement = payload.OperatingQuarterMap{
		"q1": {
			"salesRevenue":      0,
			"directCost":        0,
			"managementSalary":  0,
			"deliveryQuantity":  0,
			"receivableBalance": 0,
		},
	}
	p.Extra.IncomeAndPenalty = payload.OperatingQuarterMap{
		"q1": {
			"discountExpense":     0,
			"extraExpensePenalty": 0,
			"extraIncomeReward":   0,
		},
	}

	return p.Normalize()
}

func resolveProjectPath(t *testing.T, paths ...string) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("resolve project path: runtime caller not available")
	}

	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	fullPath := filepath.Join(append([]string{root}, paths...)...)
	if _, err := os.Stat(fullPath); err != nil {
		t.Fatalf("resolve project path %s: %v", fullPath, err)
	}
	return fullPath
}

func int64Ptr(value int64) *int64 {
	return &value
}
