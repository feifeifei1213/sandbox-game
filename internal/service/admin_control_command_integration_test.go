package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

func TestOpenNextYearOpensFormalYearForNormalGroupsOnly(t *testing.T) {
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
	markAllExistingGroupsBankrupt(t, ctx, tx, 0)

	bankruptYearNo := 0
	normalGroupID := createAdminIntegrationGroup(t, ctx, tx, "管理员开年-正常组", enum.BusinessStatusNormal, nil)
	bankruptGroupID := createAdminIntegrationGroup(t, ctx, tx, "管理员开年-破产组", enum.BusinessStatusBankrupt, &bankruptYearNo)

	createDetailedGroupYearState(t, ctx, tx, normalGroupID, 0, enum.YearTypeDemo, enum.YearStatusCompleted, enum.StageStatusYearEndOpen, enum.ReportStatusSubmitted, false, 5, 1)
	createDetailedGroupYearState(t, ctx, tx, bankruptGroupID, 0, enum.YearTypeDemo, enum.YearStatusLocked, enum.StageStatusQ1Open, enum.ReportStatusLocked, false, 0, 0)
	ensureFormalYearStatesForAllGroups(t, ctx, tx, 1, 1)

	service := NewAdminControlCommandService(tx)
	result, err := service.OpenNextYear(ctx, OpenNextYearCommand{
		TargetYearNo: 1,
		OperatorID:   90011,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("open next year: %v", err)
	}

	if result.PreviousOpenYear != 0 || result.CurrentOpenYear != 1 || result.OpenedYearNo != 1 {
		t.Fatalf("unexpected open next year result: %#v", result)
	}
	if result.FinalYear != 2 {
		t.Fatalf("expected final year 2, got %d", result.FinalYear)
	}
	if result.CanOpenNextYear {
		t.Fatalf("expected canOpenNextYear to be false immediately after opening year 1")
	}
	if result.OpenNextYearBlockedReason != openNextYearBlockedReasonUnfinishedReports {
		t.Fatalf("expected blocked reason %q, got %q", openNextYearBlockedReasonUnfinishedReports, result.OpenNextYearBlockedReason)
	}
	if result.LatestAdminAction == nil || result.LatestAdminAction.ActionCode != adminActionCodeOpenNextYear {
		t.Fatalf("expected latest admin action OPEN_NEXT_YEAR, got %#v", result.LatestAdminAction)
	}

	gameConfigRepo := repository.NewGameConfigRepository(tx)
	gameConfig, err := gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		t.Fatalf("reload game config: %v", err)
	}
	if gameConfig.CurrentOpenYear != 1 {
		t.Fatalf("expected current open year to become 1, got %d", gameConfig.CurrentOpenYear)
	}

	groupYearRepo := repository.NewGroupYearStateRepository(tx)
	normalYearOne, err := groupYearRepo.GetByGroupIDAndYear(ctx, normalGroupID, 1)
	if err != nil {
		t.Fatalf("reload normal group year 1 state: %v", err)
	}
	if normalYearOne.YearStatus != enum.YearStatusOperating || normalYearOne.StageStatus != enum.StageStatusQ1Open || normalYearOne.ReportStatus != enum.ReportStatusLocked {
		t.Fatalf("unexpected normal group year 1 state after open: %#v", normalYearOne)
	}

	bankruptYearOne, err := groupYearRepo.GetByGroupIDAndYear(ctx, bankruptGroupID, 1)
	if err != nil {
		t.Fatalf("reload bankrupt group year 1 state: %v", err)
	}
	if bankruptYearOne.YearStatus != enum.YearStatusLocked || bankruptYearOne.StageStatus != enum.StageStatusQ1Open || bankruptYearOne.ReportStatus != enum.ReportStatusLocked {
		t.Fatalf("expected bankrupt group to stay locked, got %#v", bankruptYearOne)
	}

	var actionCount int64
	if err := tx.WithContext(ctx).
		Model(&entity.AdminActionLog{}).
		Where("action_code = ?", adminActionCodeOpenNextYear).
		Count(&actionCount).Error; err != nil {
		t.Fatalf("count admin action logs: %v", err)
	}
	if actionCount != 1 {
		t.Fatalf("expected one OPEN_NEXT_YEAR action log, got %d", actionCount)
	}
}

func TestUnlockYearInvalidatesSubmittedArtifactsAndRecoversBankruptGroup(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureIntegrationGameConfig(t, ctx, tx, 2, 1, true)
	markAllExistingGroupsBankrupt(t, ctx, tx, 1)

	bankruptYearNo := 1
	groupID := createAdminIntegrationGroup(t, ctx, tx, "管理员解锁-破产组", enum.BusinessStatusBankrupt, &bankruptYearNo)
	createDetailedGroupYearState(t, ctx, tx, groupID, 1, enum.YearTypeFormal, enum.YearStatusCompleted, enum.StageStatusYearEndOpen, enum.ReportStatusSubmitted, true, 5, 1)
	createPreviousFormalReportRecord(t, ctx, tx, groupID, 1, time.Now())
	createSummarySnapshotFixture(t, ctx, tx, groupID, 1, 0, 0, 69, enum.BusinessStatusBankrupt, 69, true, 1)

	service := NewAdminControlCommandService(tx)
	result, err := service.UnlockYear(ctx, UnlockYearCommand{
		GroupID:          groupID,
		YearNo:           1,
		UnlockTargetType: unlockTargetTypeOperating,
		TargetStageCode:  state.StageCodeQ2,
		Reason:           "integration fix: rollback to Q2",
		OperatorID:       90012,
		OperatorName:     "integration-admin",
	})
	if err != nil {
		t.Fatalf("unlock year: %v", err)
	}

	if result.YearStatus != enum.YearStatusOperating {
		t.Fatalf("expected unlocked year status OPERATING, got %s", result.YearStatus)
	}
	if result.UnlockTargetType != unlockTargetTypeOperating {
		t.Fatalf("expected unlock target type OPERATING, got %s", result.UnlockTargetType)
	}
	if result.TargetStageCode == nil || *result.TargetStageCode != state.StageCodeQ2 {
		t.Fatalf("expected targetStageCode Q2, got %#v", result.TargetStageCode)
	}
	if result.EditableStageCode == nil || *result.EditableStageCode != state.StageCodeQ2 {
		t.Fatalf("expected editableStageCode Q2, got %#v", result.EditableStageCode)
	}
	if result.StageStatus != enum.StageStatusQ2Open {
		t.Fatalf("expected unlocked stage status Q2_OPEN, got %s", result.StageStatus)
	}
	if result.ReportStatus != enum.ReportStatusLocked {
		t.Fatalf("expected unlocked report status REPORT_LOCKED, got %s", result.ReportStatus)
	}
	if result.SummaryEffective {
		t.Fatalf("expected summaryEffective false after unlock")
	}
	if result.BusinessStatus != enum.BusinessStatusNormal {
		t.Fatalf("expected business status recovered to NORMAL, got %s", result.BusinessStatus)
	}
	if result.UnlockLogID <= 0 {
		t.Fatalf("expected unlock log id to be generated, got %d", result.UnlockLogID)
	}

	groupRepo := repository.NewGroupRepository(tx)
	group, err := groupRepo.GetByID(ctx, groupID)
	if err != nil {
		t.Fatalf("reload group: %v", err)
	}
	if group.BusinessStatus != enum.BusinessStatusNormal {
		t.Fatalf("expected group business status NORMAL after recovery, got %s", group.BusinessStatus)
	}
	if group.BankruptYearNo != nil {
		t.Fatalf("expected bankrupt year to be cleared, got %v", *group.BankruptYearNo)
	}

	groupYearRepo := repository.NewGroupYearStateRepository(tx)
	yearState, err := groupYearRepo.GetByGroupIDAndYear(ctx, groupID, 1)
	if err != nil {
		t.Fatalf("reload year state: %v", err)
	}
	if yearState.YearStatus != enum.YearStatusOperating || yearState.StageStatus != enum.StageStatusQ2Open || yearState.ReportStatus != enum.ReportStatusLocked || yearState.SummaryEffective {
		t.Fatalf("unexpected year state after unlock: %#v", yearState)
	}

	reportRepo := repository.NewReportRepository(tx)
	report, err := reportRepo.FindByGroupIDAndYear(ctx, groupID, 1)
	if err != nil {
		t.Fatalf("reload report after unlock: %v", err)
	}
	if report.SubmittedAt != nil {
		t.Fatalf("expected report submittedAt to be cleared after unlock")
	}
	if report.BalanceCheckPassed {
		t.Fatalf("expected report balance flag to be reset after unlock")
	}

	var summary entity.GroupSummarySnapshot
	if err := tx.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, 1).
		First(&summary).Error; err != nil {
		t.Fatalf("reload summary snapshot: %v", err)
	}
	if summary.SummaryEffective {
		t.Fatalf("expected summary snapshot to be withdrawn after unlock")
	}

	var unlockLogCount int64
	if err := tx.WithContext(ctx).
		Model(&entity.AdminUnlockLog{}).
		Where("group_id = ? AND year_no = ?", groupID, 1).
		Count(&unlockLogCount).Error; err != nil {
		t.Fatalf("count admin unlock logs: %v", err)
	}
	if unlockLogCount != 1 {
		t.Fatalf("expected one admin unlock log, got %d", unlockLogCount)
	}

	var actionLogCount int64
	if err := tx.WithContext(ctx).
		Model(&entity.AdminActionLog{}).
		Where("action_code = ? AND target_group_id = ? AND target_year_no = ?", adminActionCodeUnlockYear, groupID, 1).
		Count(&actionLogCount).Error; err != nil {
		t.Fatalf("count unlock action logs: %v", err)
	}
	if actionLogCount != 1 {
		t.Fatalf("expected one unlock action log, got %d", actionLogCount)
	}
}

func TestInitializeGameCreatesGroupsAccountsAndYearStates(t *testing.T) {
	db := openIntegrationMySQL(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() {
		_ = tx.Rollback().Error
	}()

	ctx := context.Background()
	ensureIntegrationGameConfig(t, ctx, tx, 3, 2, true)
	clearInitializationFixtures(t, ctx, tx)

	service := NewAdminControlCommandService(tx)
	result, err := service.InitializeGame(ctx, InitializeGameCommand{
		GroupCount:   3,
		OperatorID:   90013,
		OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("initialize game: %v", err)
	}

	if !result.Initialized || result.GroupCount != 3 {
		t.Fatalf("unexpected initialize result: %#v", result)
	}
	if result.CreatedGroupCount != 3 || result.CreatedAccountCount != 3 {
		t.Fatalf("unexpected created counts: %#v", result)
	}
	if result.CreatedYearStateCount != 12 {
		t.Fatalf("expected 12 year states for 3 groups x 4 years, got %d", result.CreatedYearStateCount)
	}

	groupRepo := repository.NewGroupRepository(tx)
	groups, err := groupRepo.ListAll(ctx)
	if err != nil {
		t.Fatalf("list groups after initialization: %v", err)
	}
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups after initialization, got %d", len(groups))
	}
	if groups[0].GroupNo != 1 || groups[0].GroupCode != "GROUP_01" || groups[0].GroupName != "第一组" {
		t.Fatalf("unexpected first group: %#v", groups[0])
	}

	accountRepo := repository.NewAccountRepository(tx)
	groupAccountCount, err := accountRepo.CountGroupAccounts(ctx)
	if err != nil {
		t.Fatalf("count group accounts after initialization: %v", err)
	}
	if groupAccountCount != 3 {
		t.Fatalf("expected 3 group accounts after initialization, got %d", groupAccountCount)
	}
	groupAccount, err := accountRepo.GetByUsername(ctx, "group03")
	if err != nil {
		t.Fatalf("load generated group03 account: %v", err)
	}
	if groupAccount.RoleType != enum.RoleTypeGroup || groupAccount.GroupID == nil {
		t.Fatalf("unexpected generated group account: %#v", groupAccount)
	}

	groupYearRepo := repository.NewGroupYearStateRepository(tx)
	yearStateCount, err := groupYearRepo.CountAll(ctx)
	if err != nil {
		t.Fatalf("count year states after initialization: %v", err)
	}
	if yearStateCount != 12 {
		t.Fatalf("expected 12 year states after initialization, got %d", yearStateCount)
	}
	groupOneDemo, err := groupYearRepo.GetByGroupIDAndYear(ctx, groups[0].ID, 0)
	if err != nil {
		t.Fatalf("load group one demo year state: %v", err)
	}
	if groupOneDemo.YearType != enum.YearTypeDemo || groupOneDemo.YearStatus != enum.YearStatusOperating {
		t.Fatalf("unexpected demo year state after initialization: %#v", groupOneDemo)
	}
	groupOneFormal, err := groupYearRepo.GetByGroupIDAndYear(ctx, groups[0].ID, 1)
	if err != nil {
		t.Fatalf("load group one formal year state: %v", err)
	}
	if groupOneFormal.YearType != enum.YearTypeFormal || groupOneFormal.YearStatus != enum.YearStatusLocked {
		t.Fatalf("unexpected formal year state after initialization: %#v", groupOneFormal)
	}

	gameConfig, err := repository.NewGameConfigRepository(tx).GetCurrent(ctx)
	if err != nil {
		t.Fatalf("reload game config after initialization: %v", err)
	}
	if gameConfig.CurrentOpenYear != 0 || gameConfig.InitialBaselineSubmitted {
		t.Fatalf("expected game config reset to year 0 and baseline unsubmitted, got %#v", gameConfig)
	}

	var actionLogCount int64
	if err := tx.WithContext(ctx).
		Model(&entity.AdminActionLog{}).
		Where("action_code = ?", adminActionCodeInitializeGame).
		Count(&actionLogCount).Error; err != nil {
		t.Fatalf("count initialize action logs: %v", err)
	}
	if actionLogCount != 1 {
		t.Fatalf("expected one initialize action log, got %d", actionLogCount)
	}
}

func markAllExistingGroupsBankrupt(t *testing.T, ctx context.Context, tx *gorm.DB, yearNo int) {
	t.Helper()

	now := time.Now()
	reason := "integration-neutralized"
	if err := tx.WithContext(ctx).
		Model(&entity.Group{}).
		Where("1 = 1").
		Updates(map[string]any{
			"business_status":  enum.BusinessStatusBankrupt,
			"bankrupt_year_no": yearNo,
			"bankrupt_reason":  reason,
			"updater":          "integration-test",
			"update_time":      now,
		}).Error; err != nil {
		t.Fatalf("mark existing groups bankrupt: %v", err)
	}
}

func clearInitializationFixtures(t *testing.T, ctx context.Context, tx *gorm.DB) {
	t.Helper()

	if err := tx.WithContext(ctx).Where("role_type = ?", enum.RoleTypeGroup).Delete(&entity.Account{}).Error; err != nil {
		t.Fatalf("clear group accounts: %v", err)
	}
	if err := tx.WithContext(ctx).Where("1 = 1").Delete(&entity.GroupYearState{}).Error; err != nil {
		t.Fatalf("clear group year states: %v", err)
	}
	if err := tx.WithContext(ctx).Where("1 = 1").Delete(&entity.Group{}).Error; err != nil {
		t.Fatalf("clear groups: %v", err)
	}
}

func ensureFormalYearStatesForAllGroups(t *testing.T, ctx context.Context, tx *gorm.DB, fromYear int, toYear int) {
	t.Helper()

	groupRepo := repository.NewGroupRepository(tx)
	groups, err := groupRepo.ListAll(ctx)
	if err != nil {
		t.Fatalf("list groups for formal year state ensure: %v", err)
	}
	groupIDs := make([]int64, 0, len(groups))
	for _, item := range groups {
		groupIDs = append(groupIDs, item.ID)
	}

	groupYearRepo := repository.NewGroupYearStateRepository(tx)
	if _, err := groupYearRepo.EnsureFormalYearStates(ctx, repository.EnsureFormalYearStatesCommand{
		GroupIDs:     groupIDs,
		FromYear:     fromYear,
		ToYear:       toYear,
		OperatorName: "integration-test",
		OperateTime:  time.Now(),
	}); err != nil {
		t.Fatalf("ensure formal year states: %v", err)
	}
}

func ensureIntegrationGameConfig(t *testing.T, ctx context.Context, tx *gorm.DB, finalYear int, currentOpenYear int, baselineSubmitted bool) int64 {
	t.Helper()

	now := time.Now()
	var item entity.GameConfig
	err := tx.WithContext(ctx).Order("id ASC").First(&item).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		item = entity.GameConfig{
			FinalYear:                finalYear,
			CurrentOpenYear:          currentOpenYear,
			RuleVersion:              "excel-final-2026-03-25",
			TemplateVersion:          "excel-page-v1",
			InitialBaselineSubmitted: baselineSubmitted,
			BaseEntity: entity.BaseEntity{
				Creator:    "integration-test",
				CreateTime: now,
				Updater:    "integration-test",
				UpdateTime: now,
			},
		}
		if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
			t.Fatalf("create game config: %v", err)
		}
	case err != nil:
		t.Fatalf("load game config: %v", err)
	default:
		if err := tx.WithContext(ctx).
			Model(&entity.GameConfig{}).
			Where("id = ?", item.ID).
			Updates(map[string]any{
				"final_year":                 finalYear,
				"current_open_year":          currentOpenYear,
				"rule_version":               "excel-final-2026-03-25",
				"template_version":           "excel-page-v1",
				"initial_baseline_submitted": baselineSubmitted,
				"updater":                    "integration-test",
				"update_time":                now,
			}).Error; err != nil {
			t.Fatalf("update game config: %v", err)
		}
	}

	return item.ID
}

func createAdminIntegrationGroup(t *testing.T, ctx context.Context, tx *gorm.DB, name string, businessStatus string, bankruptYearNo *int) int64 {
	t.Helper()

	now := time.Now()
	uniqueSeed := nextIntegrationUniqueSeed()
	group := entity.Group{
		GroupNo:        integrationGroupNoFromSeed(uniqueSeed),
		GroupCode:      fmt.Sprintf("IT_ADMIN_%d", uniqueSeed),
		GroupName:      name,
		BusinessStatus: businessStatus,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if businessStatus == enum.BusinessStatusBankrupt && bankruptYearNo != nil {
		year := *bankruptYearNo
		reason := "integration-test"
		group.BankruptYearNo = &year
		group.BankruptReason = &reason
	}
	if err := tx.WithContext(ctx).Create(&group).Error; err != nil {
		t.Fatalf("create admin integration group: %v", err)
	}
	return group.ID
}

func createDetailedGroupYearState(
	t *testing.T,
	ctx context.Context,
	tx *gorm.DB,
	groupID int64,
	yearNo int,
	yearType string,
	yearStatus string,
	stageStatus string,
	reportStatus string,
	summaryEffective bool,
	latestStageSubmitVersion int,
	latestReportSubmitVersion int,
) {
	t.Helper()

	now := time.Now()
	item := entity.GroupYearState{
		GroupID:                   groupID,
		YearNo:                    yearNo,
		YearType:                  yearType,
		YearStatus:                yearStatus,
		StageStatus:               stageStatus,
		ReportStatus:              reportStatus,
		SummaryEffective:          summaryEffective,
		LatestStageSubmitVersion:  latestStageSubmitVersion,
		LatestReportSubmitVersion: latestReportSubmitVersion,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		t.Fatalf("create group year state fixture: %v", err)
	}
}

func createSummarySnapshotFixture(
	t *testing.T,
	ctx context.Context,
	tx *gorm.DB,
	groupID int64,
	yearNo int,
	revenue float64,
	profit float64,
	equity float64,
	businessStatus string,
	rankingValue float64,
	summaryEffective bool,
	sourceReportSubmitVersion int,
) {
	t.Helper()

	now := time.Now()
	item := entity.GroupSummarySnapshot{
		GroupID:                   groupID,
		YearNo:                    yearNo,
		Revenue:                   revenue,
		Profit:                    profit,
		Equity:                    equity,
		BusinessStatus:            businessStatus,
		RankingValue:              rankingValue,
		SummaryEffective:          summaryEffective,
		SourceReportSubmitVersion: sourceReportSubmitVersion,
		BaseEntity: entity.BaseEntity{
			Creator:    "integration-test",
			CreateTime: now,
			Updater:    "integration-test",
			UpdateTime: now,
		},
	}
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		t.Fatalf("create summary snapshot fixture: %v", err)
	}
}
