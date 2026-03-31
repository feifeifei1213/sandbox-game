package service

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"testing"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

func TestOperatingViewAfterUnlockShowsInvalidDraftScopesAndRetainedReportDraft(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureIntegrationGameConfig(t, ctx, tx, 2, 0, true)

	now := time.Now()
	groupID := createIntegrationGroup(t, ctx, tx, now, "I1_OPERATING_UNLOCK_VIEW")
	createInitialBaselineRecord(t, ctx, tx, groupID, now)
	createDetailedGroupYearState(t, ctx, tx, groupID, 0, enum.YearTypeDemo, enum.YearStatusCompleted, enum.StageStatusYearEndOpen, enum.ReportStatusSubmitted, false, 5, 1)
	createCompletedOperatingArtifacts(t, ctx, tx, groupID, 0, now)
	createSubmittedReportFixture(t, ctx, tx, groupID, 0, now)

	adminService := NewAdminControlCommandService(tx)
	if _, err := adminService.UnlockYear(ctx, UnlockYearCommand{
		GroupID:          groupID,
		YearNo:           0,
		UnlockTargetType: "OPERATING",
		TargetStageCode:  state.StageCodeQ2,
		Reason:           "integration: rollback to Q2",
		OperatorID:       91001,
		OperatorName:     "integration-admin",
	}); err != nil {
		t.Fatalf("unlock operating year: %v", err)
	}

	_, queryService := buildPlayerOperatingServices(tx)
	view, err := queryService.GetYearView(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("load operating view after unlock: %v", err)
	}

	if !view.HasInvalidDraft {
		t.Fatalf("expected operating view to expose invalid draft state")
	}
	if !slices.Equal(view.InvalidScopes, []string{state.StageCodeQ3, state.StageCodeQ4, state.StageCodeYearEnd}) {
		t.Fatalf("unexpected invalid scopes: %#v", view.InvalidScopes)
	}
	if !view.HasRetainedReportDraft {
		t.Fatalf("expected retained report draft flag after operating unlock")
	}
	if view.CurrentStageCode != state.StageCodeQ2 {
		t.Fatalf("expected current stage code Q2, got %s", view.CurrentStageCode)
	}
	if !slices.Equal(view.EditableScopes, []string{state.OperatingScopeQ2}) {
		t.Fatalf("unexpected editable scopes after unlock: %#v", view.EditableScopes)
	}
}

func TestReportViewAfterUnlockShowsInvalidDraft(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureIntegrationGameConfig(t, ctx, tx, 2, 0, true)

	now := time.Now()
	groupID := createIntegrationGroup(t, ctx, tx, now, "I1_REPORT_UNLOCK_VIEW")
	createInitialBaselineRecord(t, ctx, tx, groupID, now)
	createDetailedGroupYearState(t, ctx, tx, groupID, 0, enum.YearTypeDemo, enum.YearStatusCompleted, enum.StageStatusYearEndOpen, enum.ReportStatusSubmitted, false, 5, 1)
	createCompletedOperatingArtifacts(t, ctx, tx, groupID, 0, now)
	createSubmittedReportFixture(t, ctx, tx, groupID, 0, now)

	adminService := NewAdminControlCommandService(tx)
	if _, err := adminService.UnlockYear(ctx, UnlockYearCommand{
		GroupID:          groupID,
		YearNo:           0,
		UnlockTargetType: "REPORT",
		Reason:           "integration: reopen report",
		OperatorID:       91002,
		OperatorName:     "integration-admin",
	}); err != nil {
		t.Fatalf("unlock report year: %v", err)
	}

	_, queryService := buildPlayerReportServices(tx)
	view, err := queryService.GetView(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("load report view after unlock: %v", err)
	}

	if !view.HasInvalidDraft {
		t.Fatalf("expected report view to expose invalid draft state")
	}
	if !view.CanEdit || !view.CanSubmit {
		t.Fatalf("expected report view to become editable after unlock")
	}
	if view.ReportStatus != enum.ReportStatusOpen {
		t.Fatalf("expected report status REPORT_OPEN after unlock, got %s", view.ReportStatus)
	}
	assertReportManualPayloadEqual(t, buildBalancedReportManualPayload(), view.ReportManualPayload)
}

func TestFindEffectiveReportIgnoresInvalidatedSubmittedRecord(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	now := time.Now()
	groupID := createIntegrationGroup(t, ctx, tx, now, "I1_EFFECTIVE_REPORT")
	createSubmittedReportFixture(t, ctx, tx, groupID, 0, now)

	reportRepo := repository.NewReportRepository(tx)
	effective, err := reportRepo.FindEffectiveByGroupIDAndYear(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("find effective report before invalidation: %v", err)
	}
	if effective.SubmittedAt == nil || !effective.BalanceCheckPassed {
		t.Fatalf("expected effective report to stay submitted before invalidation")
	}

	invalidated, err := reportRepo.InvalidateSubmission(ctx, groupID, 0, "integration-admin", now)
	if err != nil {
		t.Fatalf("invalidate report submission: %v", err)
	}
	if !invalidated {
		t.Fatalf("expected invalidate submission to affect one row")
	}

	if _, err := reportRepo.FindEffectiveByGroupIDAndYear(ctx, groupID, 0); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected effective report lookup to ignore invalidated row, got %v", err)
	}

	raw, err := reportRepo.FindByGroupIDAndYear(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("find raw report after invalidation: %v", err)
	}
	if raw.SubmittedAt != nil || raw.BalanceCheckPassed {
		t.Fatalf("expected raw report row to stay retained but not effective: %#v", raw)
	}
}

func createCompletedOperatingArtifacts(t *testing.T, ctx context.Context, tx *gorm.DB, groupID int64, yearNo int, now time.Time) {
	t.Helper()

	operatingRepo := repository.NewOperatingRepository(tx)
	operatingPayload := payload.NewOperatingPayload().Normalize()
	operatingPayload.Quarter.DeliverySettlement = payload.OperatingQuarterMap{
		"q3": {"salesRevenue": 18},
		"q4": {"salesRevenue": 26},
	}

	if err := operatingRepo.UpsertDraft(ctx, repository.UpsertOperatingDraftCommand{
		GroupID:          groupID,
		YearNo:           yearNo,
		StageStatus:      enum.StageStatusYearEndOpen,
		OperatingPayload: operatingPayload,
		LastAutoSavedAt:  now,
		OperatorName:     "integration-test",
	}); err != nil {
		t.Fatalf("upsert operating draft fixture: %v", err)
	}

	stageCodes := []string{state.StageCodeQ1, state.StageCodeQ2, state.StageCodeQ3, state.StageCodeQ4, state.StageCodeYearEnd}
	for index, stageCode := range stageCodes {
		if err := operatingRepo.CreateStageSubmission(ctx, repository.CreateStageSubmissionCommand{
			GroupID:                  groupID,
			YearNo:                   yearNo,
			StageCode:                stageCode,
			SubmitVersion:            index + 1,
			PeriodEndCash:            float64(30 + index),
			OperatingPayloadSnapshot: operatingPayload,
			StateBeforeJSON:          []byte("{}"),
			StateAfterJSON:           []byte("{}"),
			SubmitterID:              92000 + int64(index),
			SubmitTime:               now.Add(time.Duration(index) * time.Minute),
		}); err != nil {
			t.Fatalf("create operating stage submission fixture: %v", err)
		}
	}
}

func createSubmittedReportFixture(t *testing.T, ctx context.Context, tx *gorm.DB, groupID int64, yearNo int, now time.Time) {
	t.Helper()

	manual := buildBalancedReportManualPayload()
	computed := buildPreviousFormalReportComputedPayload()
	rawManual, err := json.Marshal(manual)
	if err != nil {
		t.Fatalf("marshal report manual fixture: %v", err)
	}
	rawComputed, err := json.Marshal(computed)
	if err != nil {
		t.Fatalf("marshal report computed fixture: %v", err)
	}

	item := entity.GroupReport{
		GroupID:               groupID,
		YearNo:                yearNo,
		ReportManualPayload:   rawManual,
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
