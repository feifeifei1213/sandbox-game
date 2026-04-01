package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/assembler"
	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
)

func TestSavePlayerReportDraftPersistsManualPayload(t *testing.T) {
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
	commandService, queryService := buildPlayerReportServices(tx)
	manualPayload := buildBalancedReportManualPayload()

	result, err := commandService.SaveDraft(ctx, SavePlayerReportDraftCommand{
		GroupID:             groupID,
		YearNo:              0,
		ReportManualPayload: manualPayload,
		OperatorName:        "integration-test",
	})
	if err != nil {
		t.Fatalf("save report draft: %v", err)
	}

	if result.GroupID != groupID || result.YearNo != 0 {
		t.Fatalf("unexpected save draft identity: %#v", result)
	}
	if result.YearStatus != enum.YearStatusReportPending {
		t.Fatalf("expected year status REPORT_PENDING, got %s", result.YearStatus)
	}
	if result.ReportStatus != enum.ReportStatusOpen {
		t.Fatalf("expected report status REPORT_OPEN, got %s", result.ReportStatus)
	}
	if result.LastDraftSavedAt.IsZero() {
		t.Fatalf("expected last draft saved at to be filled")
	}

	reportRepo := repository.NewReportRepository(tx)
	report, err := reportRepo.FindByGroupIDAndYear(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("find saved report draft: %v", err)
	}
	storedManual, err := unmarshalReportManualPayload(report.ReportManualPayload)
	if err != nil {
		t.Fatalf("unmarshal saved manual payload: %v", err)
	}
	assertReportManualPayloadEqual(t, manualPayload, storedManual)
	if report.SubmittedAt != nil {
		t.Fatalf("expected submittedAt to stay nil after draft save")
	}

	view, err := queryService.GetView(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("load report view after save: %v", err)
	}
	assertReportManualPayloadEqual(t, manualPayload, view.ReportManualPayload)
	if !view.CanEdit || !view.CanSubmit {
		t.Fatalf("expected draft view to remain editable/submittable")
	}
}

func TestSubmitPlayerReportCompletesFormalYearAndWritesSummary(t *testing.T) {
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

	groupID := createFormalReportFixtures(t, ctx, tx)
	commandService, queryService := buildPlayerReportServices(tx)
	manualPayload := buildBalancedReportManualPayload()

	result, err := commandService.Submit(ctx, SubmitPlayerReportCommand{
		GroupID:             groupID,
		YearNo:              1,
		ReportManualPayload: manualPayload,
		SubmitterID:         90002,
		OperatorName:        "integration-test",
	})
	if err != nil {
		t.Fatalf("submit report: %v", err)
	}

	if result.YearStatus != enum.YearStatusCompleted {
		t.Fatalf("expected year status COMPLETED, got %s", result.YearStatus)
	}
	if result.ReportStatus != enum.ReportStatusSubmitted {
		t.Fatalf("expected report status REPORT_SUBMITTED, got %s", result.ReportStatus)
	}
	if !result.SummaryEffective {
		t.Fatalf("expected summaryEffective true for formal year")
	}
	if !result.BalanceCheckPassed {
		t.Fatalf("expected balance check passed")
	}
	if result.LatestReportSubmitVersion != 1 {
		t.Fatalf("expected latest report submit version 1, got %d", result.LatestReportSubmitVersion)
	}
	if result.BusinessStatus != enum.BusinessStatusNormal {
		t.Fatalf("expected business status NORMAL, got %s", result.BusinessStatus)
	}
	if result.SubmittedAt.IsZero() {
		t.Fatalf("expected submittedAt to be filled")
	}

	yearRepo := repository.NewGroupYearStateRepository(tx)
	yearState, err := yearRepo.GetByGroupIDAndYear(ctx, groupID, 1)
	if err != nil {
		t.Fatalf("reload year state: %v", err)
	}
	if yearState.YearStatus != enum.YearStatusCompleted || yearState.ReportStatus != enum.ReportStatusSubmitted {
		t.Fatalf("unexpected reloaded year state: year=%s report=%s", yearState.YearStatus, yearState.ReportStatus)
	}
	if !yearState.SummaryEffective {
		t.Fatalf("expected persisted summaryEffective true")
	}
	if yearState.LatestReportSubmitVersion != 1 {
		t.Fatalf("expected persisted latest report version 1, got %d", yearState.LatestReportSubmitVersion)
	}

	reportRepo := repository.NewReportRepository(tx)
	report, err := reportRepo.FindByGroupIDAndYear(ctx, groupID, 1)
	if err != nil {
		t.Fatalf("reload persisted report: %v", err)
	}
	if !report.BalanceCheckPassed || report.SubmittedAt == nil {
		t.Fatalf("expected persisted report to be submitted with balance passed")
	}
	storedManual, err := unmarshalReportManualPayload(report.ReportManualPayload)
	if err != nil {
		t.Fatalf("unmarshal persisted manual payload: %v", err)
	}
	assertReportManualPayloadEqual(t, manualPayload, storedManual)

	var storedComputed payload.ReportComputedPayload
	if err := json.Unmarshal(report.ReportComputedPayload, &storedComputed); err != nil {
		t.Fatalf("unmarshal persisted computed payload: %v", err)
	}
	if gap := storedComputed.BalanceGap(); gap != 0 {
		t.Fatalf("expected persisted computed payload balance gap 0, got %f", gap)
	}
	if storedComputed.ReportTotalEquity != 69 {
		t.Fatalf("expected report total equity 69, got %f", storedComputed.ReportTotalEquity)
	}

	var snapshot entity.GroupSummarySnapshot
	if err := tx.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, 1).
		First(&snapshot).Error; err != nil {
		t.Fatalf("load summary snapshot: %v", err)
	}
	if !snapshot.SummaryEffective {
		t.Fatalf("expected summary snapshot effective")
	}
	if snapshot.Equity != 69 || snapshot.RankingValue != 69 {
		t.Fatalf("unexpected summary equity/ranking: equity=%f ranking=%f", snapshot.Equity, snapshot.RankingValue)
	}
	if snapshot.SourceReportSubmitVersion != 1 {
		t.Fatalf("expected source report submit version 1, got %d", snapshot.SourceReportSubmitVersion)
	}

	var submission entity.GroupReportSubmission
	if err := tx.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, 1).
		First(&submission).Error; err != nil {
		t.Fatalf("load report submission: %v", err)
	}
	if !submission.BalanceCheckPassed || submission.SubmitVersion != 1 {
		t.Fatalf("unexpected report submission persisted state: balance=%t version=%d", submission.BalanceCheckPassed, submission.SubmitVersion)
	}

	view, err := queryService.GetView(ctx, groupID, 1)
	if err != nil {
		t.Fatalf("load report view after submit: %v", err)
	}
	if view.CanEdit || view.CanSubmit {
		t.Fatalf("expected submitted report view to be readonly")
	}
	if view.YearStatus != enum.YearStatusCompleted || view.ReportStatus != enum.ReportStatusSubmitted {
		t.Fatalf("unexpected submitted view status: year=%s report=%s", view.YearStatus, view.ReportStatus)
	}
	assertReportManualPayloadEqual(t, manualPayload, view.ReportManualPayload)
	if view.ReportComputedPayload.BalanceGap() != 0 {
		t.Fatalf("expected view computed payload balance gap 0, got %f", view.ReportComputedPayload.BalanceGap())
	}
}

func buildPlayerReportServices(db *gorm.DB) (*PlayerReportCommandService, *PlayerReportQueryService) {
	gameConfigRepo := repository.NewGameConfigRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	groupYearRepo := repository.NewGroupYearStateRepository(db)
	operatingRepo := repository.NewOperatingRepository(db)
	initialBaselineRepo := repository.NewInitialBaselineRepository(db)
	reportRepo := repository.NewReportRepository(db)
	noticeRepo := repository.NewNoticeRepository(db)
	adjustmentRepo := repository.NewGroupAdjustmentRepository(db)
	playerNoticeService := NewPlayerNoticeService(noticeRepo, adjustmentRepo)

	commandService := NewPlayerReportCommandService(
		db,
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaselineRepo,
		reportRepo,
		playerNoticeService,
	)
	queryService := NewPlayerReportQueryService(
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaselineRepo,
		reportRepo,
		assembler.NewPlayerReportAssembler(),
		playerNoticeService,
	)

	return commandService, queryService
}

func createDemoReportFixtures(t *testing.T, ctx context.Context, tx *gorm.DB) int64 {
	t.Helper()

	now := time.Now()
	groupID := createIntegrationGroup(t, ctx, tx, now, "IT_REPORT_DEMO")
	createInitialBaselineRecord(t, ctx, tx, groupID, now)
	createGroupYearStateRecord(t, ctx, tx, groupID, 0, enum.YearTypeDemo, enum.YearStatusReportPending, enum.StageStatusYearEndOpen, enum.ReportStatusOpen)

	return groupID
}

func createFormalReportFixtures(t *testing.T, ctx context.Context, tx *gorm.DB) int64 {
	t.Helper()

	now := time.Now()
	groupID := createIntegrationGroup(t, ctx, tx, now, "IT_REPORT_FORMAL")

	createGroupYearStateRecord(t, ctx, tx, groupID, 0, enum.YearTypeDemo, enum.YearStatusCompleted, enum.StageStatusYearEndOpen, enum.ReportStatusSubmitted)
	createPreviousFormalReportRecord(t, ctx, tx, groupID, 0, now)
	createGroupYearStateRecord(t, ctx, tx, groupID, 1, enum.YearTypeFormal, enum.YearStatusReportPending, enum.StageStatusYearEndOpen, enum.ReportStatusOpen)

	return groupID
}

func createIntegrationGroup(t *testing.T, ctx context.Context, tx *gorm.DB, now time.Time, codePrefix string) int64 {
	t.Helper()

	uniqueSeed := nextIntegrationUniqueSeed()
	groupCode := fmt.Sprintf("ITR%d", uniqueSeed)
	group := entity.Group{
		GroupNo:        integrationGroupNoFromSeed(uniqueSeed),
		GroupCode:      groupCode,
		GroupName:      codePrefix,
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
	return group.ID
}

func createInitialBaselineRecord(t *testing.T, ctx context.Context, tx *gorm.DB, groupID int64, now time.Time) {
	t.Helper()

	baselineJSON, err := json.Marshal(buildIntegrationBaselinePayload())
	if err != nil {
		t.Fatalf("marshal baseline payload: %v", err)
	}
	submittedAt := now
	baseline := entity.InitialBaseline{
		GroupID:         groupID,
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
}

func createGroupYearStateRecord(t *testing.T, ctx context.Context, tx *gorm.DB, groupID int64, yearNo int, yearType string, yearStatus string, stageStatus string, reportStatus string) {
	t.Helper()

	now := time.Now()
	item := entity.GroupYearState{
		GroupID:                   groupID,
		YearNo:                    yearNo,
		YearType:                  yearType,
		YearStatus:                yearStatus,
		StageStatus:               stageStatus,
		ReportStatus:              reportStatus,
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
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		t.Fatalf("create year state fixture: %v", err)
	}
}

func createPreviousFormalReportRecord(t *testing.T, ctx context.Context, tx *gorm.DB, groupID int64, yearNo int, now time.Time) {
	t.Helper()

	manual := buildBalancedReportManualPayload()
	computed := buildPreviousFormalReportComputedPayload()
	manualJSON, err := json.Marshal(manual)
	if err != nil {
		t.Fatalf("marshal previous manual payload: %v", err)
	}
	computedJSON, err := json.Marshal(computed)
	if err != nil {
		t.Fatalf("marshal previous computed payload: %v", err)
	}
	submittedAt := now
	item := entity.GroupReport{
		GroupID:               groupID,
		YearNo:                yearNo,
		ReportManualPayload:   manualJSON,
		ReportComputedPayload: computedJSON,
		BalanceCheckPassed:    true,
		LastAutoSavedAt:       &now,
		SubmittedAt:           &submittedAt,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		t.Fatalf("create previous report fixture: %v", err)
	}
}

func buildBalancedReportManualPayload() payload.ReportManualPayload {
	return payload.ReportManualPayload{
		WorkInProgress: float64Ptr(6),
		FinishedGoods:  float64Ptr(4),
		RawMaterials:   float64Ptr(1),
		IncomeTaxRate:  float64Ptr(0),
	}
}

func buildPreviousFormalReportComputedPayload() payload.ReportComputedPayload {
	return payload.ReportComputedPayload{
		ReportSalesRevenue:          0,
		ReportDirectCost:            0,
		ReportGrossProfit:           0,
		ReportComprehensiveCost:     0,
		ReportDepreciation:          0,
		ReportOperatingProfit:       0,
		ReportFinanceIncomeExpense:  0,
		ReportExtraIncomeExpense:    0,
		ReportPreTaxProfit:          0,
		ReportIncomeTax:             0,
		ReportNetProfit:             0,
		ReportWorkInProgress:        6,
		ReportFinishedGoods:         4,
		ReportRawMaterials:          1,
		ReportWorkInConstruction:    0,
		ReportFactoryAsset:          40,
		ReportLineResidual:          3,
		ReportDepreciableAsset:      0,
		ReportTotalNonCurrentAssets: 43,
		ReportCash:                  35,
		ReportReceivable:            0,
		ReportPostTaxCash:           35,
		ReportTotalCurrentAssets:    46,
		ReportTotalAssets:           89,
		ReportShortTermLiability:    20,
		ReportLongTermLiability:     0,
		ReportTotalLiability:        20,
		ReportShareCapital:          50,
		ReportRetainedEarnings:      19,
		ReportTotalEquity:           69,
		ReportTotalLiabilityEquity:  89,
	}
}

func unmarshalReportManualPayload(raw []byte) (payload.ReportManualPayload, error) {
	var result payload.ReportManualPayload
	if len(raw) == 0 {
		return result, nil
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return payload.ReportManualPayload{}, err
	}
	return result, nil
}

func assertReportManualPayloadEqual(t *testing.T, expected payload.ReportManualPayload, actual payload.ReportManualPayload) {
	t.Helper()
	assertFloat64PointerEqual(t, "workInProgress", expected.WorkInProgress, actual.WorkInProgress)
	assertFloat64PointerEqual(t, "finishedGoods", expected.FinishedGoods, actual.FinishedGoods)
	assertFloat64PointerEqual(t, "rawMaterials", expected.RawMaterials, actual.RawMaterials)
	assertFloat64PointerEqual(t, "incomeTaxRate", expected.IncomeTaxRate, actual.IncomeTaxRate)
}

func assertFloat64PointerEqual(t *testing.T, label string, expected *float64, actual *float64) {
	t.Helper()
	if expected == nil || actual == nil {
		if expected != actual {
			t.Fatalf("expected %s pointer equality, got expected=%v actual=%v", label, expected, actual)
		}
		return
	}
	if *expected != *actual {
		t.Fatalf("expected %s=%f, got %f", label, *expected, *actual)
	}
}

func float64Ptr(value float64) *float64 {
	return &value
}
