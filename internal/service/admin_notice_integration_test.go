package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

func TestAdminNoticeWorkflowBuildsPlayerBoardAndOverlaysReward(t *testing.T) {
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

	groupID := createAdminNoticeFixtures(t, ctx, tx, enum.YearStatusOperating, enum.StageStatusQ1Open)
	commandService := NewAdminNoticeCommandService(tx)
	playerNoticeService := buildPlayerNoticeService(tx)

	generalResult, err := commandService.SendGeneralNotice(ctx, SendGeneralNoticeCommand{
		TargetScope:   noticeScopeGroup,
		TargetGroupID: &groupID,
		Content:       "I2 notice integration general",
		Pinned:        true,
		OperatorID:    1,
		OperatorName:  "integration-test",
	})
	if err != nil {
		t.Fatalf("send general notice: %v", err)
	}
	if generalResult.NoticeID <= 0 {
		t.Fatalf("expected generated notice id, got %#v", generalResult)
	}

	adjustmentResult, err := commandService.SendAdjustment(ctx, SendAdjustmentCommand{
		GroupID:        groupID,
		YearNo:         0,
		StageCode:      state.StageCodeQ1,
		AdjustmentType: adjustmentTypeReward,
		Amount:         88,
		Reason:         "I2 notice integration reward",
		OperatorID:     1,
		OperatorName:   "integration-test",
	})
	if err != nil {
		t.Fatalf("send reward adjustment: %v", err)
	}
	if adjustmentResult.AdjustmentID <= 0 {
		t.Fatalf("expected generated adjustment id, got %#v", adjustmentResult)
	}
	syncService := NewPlayerAdjustmentSyncService(tx, buildPlayerNoticeService(tx))
	syncResult, err := syncService.GetSync(ctx, groupID, 0, 0)
	if err != nil {
		t.Fatalf("sync adjustment revision: %v", err)
	}
	if syncResult.NotModified || syncResult.AdjustmentRevision != 1 {
		t.Fatalf("expected changed revision 1, got %#v", syncResult)
	}
	notModified, err := syncService.GetSync(ctx, groupID, 0, 1)
	if err != nil {
		t.Fatalf("sync unchanged adjustment revision: %v", err)
	}
	if !notModified.NotModified || notModified.AdjustmentRevision != 1 {
		t.Fatalf("expected notModified revision 1, got %#v", notModified)
	}

	overlaidPayload, err := playerNoticeService.OverlayAdjustments(ctx, groupID, 0, buildValidQ1OperatingPayload())
	if err != nil {
		t.Fatalf("overlay adjustments: %v", err)
	}

	q1Extra := overlaidPayload.Normalize().Extra.IncomeAndPenalty["q1"]
	if reward, ok := q1Extra["extraIncomeReward"].(float64); !ok || reward != 88 {
		t.Fatalf("expected q1 reward 88, got %#v", q1Extra["extraIncomeReward"])
	}
	if penalty, ok := q1Extra["extraExpensePenalty"].(float64); !ok || penalty != 0 {
		t.Fatalf("expected q1 penalty 0, got %#v", q1Extra["extraExpensePenalty"])
	}

	board, err := playerNoticeService.BuildBoard(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("build player notice board: %v", err)
	}
	if board.PinnedNotice == nil {
		t.Fatalf("expected pinned notice to exist")
	}
	if board.PinnedNotice.Content != "I2 notice integration general" {
		t.Fatalf("unexpected pinned notice content: %#v", board.PinnedNotice)
	}
	if len(board.RecentList) != 1 {
		t.Fatalf("expected one adjustment notice item, got %d", len(board.RecentList))
	}
	if board.RecentList[0].Kind != adjustmentTypeReward {
		t.Fatalf("expected reward notice item, got %#v", board.RecentList[0])
	}
	if board.RecentList[0].Amount == nil || *board.RecentList[0].Amount != 88 {
		t.Fatalf("expected reward amount 88, got %#v", board.RecentList[0].Amount)
	}

	preview, err := commandService.PreviewAdjustment(ctx, PreviewAdjustmentCommand{
		Operation: adjustmentOperationVoid, AdjustmentID: adjustmentResult.AdjustmentID,
	})
	if err != nil {
		t.Fatalf("preview void adjustment: %v", err)
	}
	if preview.CashAfter >= preview.CashBefore {
		t.Fatalf("expected voiding reward to reduce cash, before=%v after=%v", preview.CashBefore, preview.CashAfter)
	}
	voidResult, err := commandService.VoidAdjustment(ctx, VoidAdjustmentCommand{
		AdjustmentID: adjustmentResult.AdjustmentID,
		Reason:       "integration correction", OperatorID: 2, OperatorName: "integration-admin",
	})
	if err != nil {
		t.Fatalf("void adjustment: %v", err)
	}
	if voidResult.AdjustmentRevision != 2 {
		t.Fatalf("expected revision 2 after void, got %d", voidResult.AdjustmentRevision)
	}
	voidedItem, err := repository.NewGroupAdjustmentRepository(tx).GetByID(ctx, adjustmentResult.AdjustmentID)
	if err != nil {
		t.Fatalf("reload voided adjustment: %v", err)
	}
	if voidedItem.Effective || voidedItem.VoidedByID == nil || *voidedItem.VoidedByID != 2 ||
		voidedItem.VoidedByName == nil || *voidedItem.VoidedByName != "integration-admin" ||
		voidedItem.VoidReason == nil || *voidedItem.VoidReason != "integration correction" || voidedItem.VoidedAt == nil {
		t.Fatalf("expected complete void audit fields, got %#v", voidedItem)
	}
	syncResult, err = syncService.GetSync(ctx, groupID, 0, 1)
	if err != nil {
		t.Fatalf("sync voided adjustment revision: %v", err)
	}
	if syncResult.NotModified || syncResult.AdjustmentRevision != 2 {
		t.Fatalf("expected changed revision 2 after void, got %#v", syncResult)
	}
	overlaidPayload, err = playerNoticeService.OverlayAdjustments(ctx, groupID, 0, buildValidQ1OperatingPayload())
	if err != nil {
		t.Fatalf("overlay adjustments after void: %v", err)
	}
	if reward := overlaidPayload.Normalize().Extra.IncomeAndPenalty["q1"]["extraIncomeReward"]; reward != float64(0) {
		t.Fatalf("expected voided reward excluded from calculation, got %#v", reward)
	}
	board, err = playerNoticeService.BuildBoard(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("build player notice board after void: %v", err)
	}
	if board.RecentList[0].Status == nil || *board.RecentList[0].Status != "VOIDED" {
		t.Fatalf("expected voided adjustment status, got %#v", board.RecentList[0].Status)
	}
}

func TestAdminNoticeCommandServiceResolvesCurrentStageInsteadOfHistoricalSubmission(t *testing.T) {
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

	groupID := createAdminNoticeFixtures(t, ctx, tx, enum.YearStatusOperating, enum.StageStatusQ2Open)
	createLockedStageSubmission(t, ctx, tx, groupID, 0, state.StageCodeQ1)
	commandService := NewAdminNoticeCommandService(tx)

	result, err := commandService.SendAdjustment(ctx, SendAdjustmentCommand{
		GroupID:        groupID,
		YearNo:         0,
		StageCode:      state.StageCodeQ1,
		AdjustmentType: adjustmentTypePenalty,
		Amount:         66,
		Reason:         "I2 locked stage should reject",
		OperatorID:     1,
		OperatorName:   "integration-test",
	})
	if err != nil {
		t.Fatalf("send adjustment at current q2 stage: %v", err)
	}
	if result.StageCode != state.StageCodeQ2 {
		t.Fatalf("expected service-resolved Q2 stage, got %s", result.StageCode)
	}
	if result.AdjustmentRevision != 1 {
		t.Fatalf("expected revision 1, got %d", result.AdjustmentRevision)
	}
}

func TestAdminNoticeCommandServiceRejectsDecimalAdjustmentAmount(t *testing.T) {
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

	groupID := createAdminNoticeFixtures(t, ctx, tx, enum.YearStatusOperating, enum.StageStatusQ1Open)
	commandService := NewAdminNoticeCommandService(tx)

	_, err := commandService.SendAdjustment(ctx, SendAdjustmentCommand{
		GroupID:        groupID,
		YearNo:         0,
		StageCode:      state.StageCodeQ1,
		AdjustmentType: adjustmentTypeReward,
		Amount:         6.5,
		Reason:         "decimal amount should reject",
		OperatorID:     1,
		OperatorName:   "integration-test",
	})
	if !errors.Is(err, ErrManualNumberNotInteger) {
		t.Fatalf("expected ErrManualNumberNotInteger, got %v", err)
	}
}

func TestAdminNoticeCommandServiceAssignsReportDraftToYearEndAndRejectsSubmittedReport(t *testing.T) {
	db := openIntegrationMySQL(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() { _ = tx.Rollback().Error }()

	ctx := context.Background()
	ensureGameConfigExists(t, ctx, tx)
	groupID := createAdminNoticeFixtures(t, ctx, tx, enum.YearStatusOperating, enum.StageStatusYearEndOpen)
	if err := tx.WithContext(ctx).Model(&entity.GroupYearState{}).
		Where("group_id = ? AND year_no = ?", groupID, 0).
		Updates(map[string]any{"year_status": enum.YearStatusReporting, "report_status": enum.ReportStatusOpen}).Error; err != nil {
		t.Fatalf("open report draft state: %v", err)
	}

	commandService := NewAdminNoticeCommandService(tx)
	result, err := commandService.SendAdjustment(ctx, SendAdjustmentCommand{
		GroupID: groupID, YearNo: 0, AdjustmentType: adjustmentTypeReward, Amount: 10,
		Reason: "report draft reward", OperatorID: 1, OperatorName: "integration-test",
	})
	if err != nil {
		t.Fatalf("send report draft adjustment: %v", err)
	}
	if result.StageCode != state.StageCodeYearEnd {
		t.Fatalf("expected report draft adjustment at YEAR_END, got %s", result.StageCode)
	}

	if err := tx.WithContext(ctx).Model(&entity.GroupYearState{}).
		Where("group_id = ? AND year_no = ?", groupID, 0).
		Updates(map[string]any{"year_status": enum.YearStatusCompleted, "report_status": enum.ReportStatusSubmitted}).Error; err != nil {
		t.Fatalf("mark report submitted: %v", err)
	}
	_, err = commandService.SendAdjustment(ctx, SendAdjustmentCommand{
		GroupID: groupID, YearNo: 0, AdjustmentType: adjustmentTypeReward, Amount: 10,
		Reason: "submitted report reward", OperatorID: 1, OperatorName: "integration-test",
	})
	if !errors.Is(err, ErrAdminAdjustmentStageLocked) {
		t.Fatalf("expected submitted report to reject adjustment, got %v", err)
	}
	_, err = commandService.VoidAdjustment(ctx, VoidAdjustmentCommand{
		AdjustmentID: result.AdjustmentID, Reason: "submitted report void",
		OperatorID: 1, OperatorName: "integration-test",
	})
	if !errors.Is(err, ErrAdminAdjustmentStageLocked) {
		t.Fatalf("expected submitted report to reject void, got %v", err)
	}
}

func TestAdjustmentPenaltyCreatesReadOnlyBankruptcySnapshotAndFreezesFurtherChanges(t *testing.T) {
	db := openIntegrationMySQL(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() { _ = tx.Rollback().Error }()

	ctx := context.Background()
	ensureGameConfigExists(t, ctx, tx)
	groupID := createAdminNoticeFixtures(t, ctx, tx, enum.YearStatusOperating, enum.StageStatusQ1Open)
	commandService := NewAdminNoticeCommandService(tx)
	result, err := commandService.SendAdjustment(ctx, SendAdjustmentCommand{
		GroupID: groupID, YearNo: 0, AdjustmentType: adjustmentTypePenalty, Amount: 1000,
		Reason: "bankruptcy penalty", OperatorID: 9, OperatorName: "bankruptcy-admin",
	})
	if err != nil {
		t.Fatalf("send bankruptcy penalty: %v", err)
	}
	if !result.Impact.WillBankrupt || result.Impact.CashAfter >= 0 {
		t.Fatalf("expected negative cash bankruptcy impact, got %#v", result.Impact)
	}

	group, err := repository.NewGroupRepository(tx).GetByID(ctx, groupID)
	if err != nil {
		t.Fatalf("reload bankrupt group: %v", err)
	}
	if group.BusinessStatus != enum.BusinessStatusBankrupt || group.BankruptYearNo == nil || *group.BankruptYearNo != 0 {
		t.Fatalf("expected irreversible bankrupt group state, got %#v", group)
	}

	yearNo := 0
	snapshots, err := repository.NewStateSnapshotRepository(tx).List(ctx, repository.StateSnapshotListFilter{
		SnapshotScope: enum.SnapshotScopeGroup, GroupID: &groupID, YearNo: &yearNo, Limit: 10,
	})
	if err != nil {
		t.Fatalf("list bankruptcy snapshots: %v", err)
	}
	if len(snapshots) != 1 || snapshots[0].TriggerCode != enum.SnapshotTriggerAdjustmentBankruptcy {
		t.Fatalf("expected one adjustment bankruptcy snapshot, got %#v", snapshots)
	}
	payloadItem, err := repository.NewStateSnapshotRepository(tx).GetPayloadBySnapshotID(ctx, snapshots[0].ID)
	if err != nil {
		t.Fatalf("load bankruptcy snapshot payload: %v", err)
	}
	envelope, err := unmarshalSnapshotPayload(payloadItem.PayloadJSON)
	if err != nil {
		t.Fatalf("unmarshal bankruptcy snapshot payload: %v", err)
	}
	if allowed, ok := envelope.Metadata["useForRestore"].(bool); !ok || allowed {
		t.Fatalf("expected useForRestore=false, got %#v", envelope.Metadata)
	}
	if operation, _ := envelope.Metadata["operation"].(string); operation != adjustmentOperationCreate {
		t.Fatalf("expected CREATE bankruptcy metadata, got %#v", envelope.Metadata)
	}

	queryService := NewAdminRollbackQueryService(repository.NewStateSnapshotRepository(tx), repository.NewGroupRepository(tx))
	detail, err := queryService.GetSnapshotDetail(ctx, snapshots[0].ID)
	if err != nil {
		t.Fatalf("load bankruptcy snapshot detail: %v", err)
	}
	if detail.Snapshot.CanRestore {
		t.Fatalf("expected bankruptcy snapshot to be read-only")
	}
	_, err = NewAdminRollbackCommandService(tx).RestoreGroupSnapshot(ctx, RestoreGroupSnapshotCommand{
		SnapshotID: snapshots[0].ID, Reason: "must reject", ConfirmText: rollbackConfirmText,
		OperatorID: 9, OperatorName: "bankruptcy-admin",
	})
	if !errors.Is(err, ErrRollbackSnapshotNotRestorable) {
		t.Fatalf("expected bankruptcy snapshot restore rejection, got %v", err)
	}

	var reportSubmissionCount int64
	if err := tx.WithContext(ctx).Model(&entity.GroupReportSubmission{}).
		Where("group_id = ? AND year_no = ?", groupID, 0).Count(&reportSubmissionCount).Error; err != nil {
		t.Fatalf("count report submissions: %v", err)
	}
	var summaryCount int64
	if err := tx.WithContext(ctx).Model(&entity.GroupSummarySnapshot{}).
		Where("group_id = ? AND year_no = ?", groupID, 0).Count(&summaryCount).Error; err != nil {
		t.Fatalf("count summary snapshots: %v", err)
	}
	if reportSubmissionCount != 0 || summaryCount != 0 {
		t.Fatalf("bankruptcy snapshot must not create formal report or summary, reports=%d summaries=%d", reportSubmissionCount, summaryCount)
	}

	_, err = commandService.SendAdjustment(ctx, SendAdjustmentCommand{
		GroupID: groupID, YearNo: 0, AdjustmentType: adjustmentTypeReward, Amount: 1,
		Reason: "must freeze", OperatorID: 9, OperatorName: "bankruptcy-admin",
	})
	if !errors.Is(err, ErrAdminAdjustmentGroupNotAvailable) {
		t.Fatalf("expected bankrupt group send rejection, got %v", err)
	}
	_, err = commandService.VoidAdjustment(ctx, VoidAdjustmentCommand{
		AdjustmentID: result.AdjustmentID, Reason: "must freeze",
		OperatorID: 9, OperatorName: "bankruptcy-admin",
	})
	if !errors.Is(err, ErrAdminAdjustmentGroupNotAvailable) {
		t.Fatalf("expected bankrupt group void rejection, got %v", err)
	}
}

func TestVoidingRewardCanTriggerIrreversibleBankruptcy(t *testing.T) {
	db := openIntegrationMySQL(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() { _ = tx.Rollback().Error }()

	ctx := context.Background()
	ensureGameConfigExists(t, ctx, tx)
	groupID := createAdminNoticeFixtures(t, ctx, tx, enum.YearStatusOperating, enum.StageStatusQ1Open)
	baseline := buildIntegrationBaselinePayload()
	baseline.BaselineCash = -100
	baselineJSON, err := json.Marshal(baseline)
	if err != nil {
		t.Fatalf("marshal negative cash baseline: %v", err)
	}
	if err := tx.WithContext(ctx).Model(&entity.InitialBaseline{}).
		Where("group_id = ?", groupID).Update("baseline_payload_json", baselineJSON).Error; err != nil {
		t.Fatalf("update negative cash baseline: %v", err)
	}

	commandService := NewAdminNoticeCommandService(tx)
	sendResult, err := commandService.SendAdjustment(ctx, SendAdjustmentCommand{
		GroupID: groupID, YearNo: 0, AdjustmentType: adjustmentTypeReward, Amount: 200,
		Reason: "temporary rescue reward", OperatorID: 7, OperatorName: "void-admin",
	})
	if err != nil {
		t.Fatalf("send rescue reward: %v", err)
	}
	if sendResult.Impact.WillBankrupt {
		t.Fatalf("reward should keep group solvent before void, got %#v", sendResult.Impact)
	}
	voidResult, err := commandService.VoidAdjustment(ctx, VoidAdjustmentCommand{
		AdjustmentID: sendResult.AdjustmentID, Reason: "remove rescue reward",
		OperatorID: 7, OperatorName: "void-admin",
	})
	if err != nil {
		t.Fatalf("void rescue reward: %v", err)
	}
	if !voidResult.Impact.WillBankrupt || voidResult.Impact.CashAfter >= 0 {
		t.Fatalf("expected void to trigger bankruptcy, got %#v", voidResult.Impact)
	}

	group, err := repository.NewGroupRepository(tx).GetByID(ctx, groupID)
	if err != nil {
		t.Fatalf("reload void-bankrupt group: %v", err)
	}
	if group.BusinessStatus != enum.BusinessStatusBankrupt {
		t.Fatalf("expected BANKRUPT after void, got %s", group.BusinessStatus)
	}
	yearNo := 0
	snapshots, err := repository.NewStateSnapshotRepository(tx).List(ctx, repository.StateSnapshotListFilter{
		SnapshotScope: enum.SnapshotScopeGroup, GroupID: &groupID, YearNo: &yearNo, Limit: 10,
	})
	if err != nil || len(snapshots) != 1 {
		t.Fatalf("expected one void bankruptcy snapshot, snapshots=%#v err=%v", snapshots, err)
	}
	payloadItem, err := repository.NewStateSnapshotRepository(tx).GetPayloadBySnapshotID(ctx, snapshots[0].ID)
	if err != nil {
		t.Fatalf("load void bankruptcy snapshot payload: %v", err)
	}
	envelope, err := unmarshalSnapshotPayload(payloadItem.PayloadJSON)
	if err != nil {
		t.Fatalf("unmarshal void bankruptcy snapshot payload: %v", err)
	}
	if operation, _ := envelope.Metadata["operation"].(string); operation != adjustmentOperationVoid {
		t.Fatalf("expected VOID bankruptcy metadata, got %#v", envelope.Metadata)
	}
}

func TestSnapshotRestoreRestoresAdjustmentStateAndRevisionButPreservesBankruptcy(t *testing.T) {
	db := openIntegrationMySQL(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	defer func() { _ = tx.Rollback().Error }()

	ctx := context.Background()
	ensureGameConfigExists(t, ctx, tx)
	groupID := createAdminNoticeFixtures(t, ctx, tx, enum.YearStatusOperating, enum.StageStatusQ1Open)
	commandService := NewAdminNoticeCommandService(tx)
	sendResult, err := commandService.SendAdjustment(ctx, SendAdjustmentCommand{
		GroupID: groupID, YearNo: 0, AdjustmentType: adjustmentTypeReward, Amount: 20,
		Reason: "captured reward", OperatorID: 5, OperatorName: "snapshot-admin",
	})
	if err != nil {
		t.Fatalf("send captured reward: %v", err)
	}
	snapshot, err := createGroupSnapshot(ctx, tx, CreateGroupSnapshotCommand{
		GroupID: groupID, YearNo: 0, StageCode: state.StageCodeQ1,
		SnapshotType: enum.SnapshotTypeManual, TriggerCode: enum.SnapshotTypeManual,
		Description: "adjustment restore fixture", OperatorID: 5, OperatorName: "snapshot-admin",
		OperateTime: time.Now(), UseForRestore: true,
	})
	if err != nil {
		t.Fatalf("create restorable adjustment snapshot: %v", err)
	}
	voidResult, err := commandService.VoidAdjustment(ctx, VoidAdjustmentCommand{
		AdjustmentID: sendResult.AdjustmentID, Reason: "void after snapshot",
		OperatorID: 5, OperatorName: "snapshot-admin",
	})
	if err != nil {
		t.Fatalf("void captured reward: %v", err)
	}
	if voidResult.AdjustmentRevision != 2 {
		t.Fatalf("expected revision 2 before restore, got %d", voidResult.AdjustmentRevision)
	}
	if err := repository.NewGroupRepository(tx).MarkBankrupt(ctx, groupID, 0, "confirmed bankruptcy", "snapshot-admin"); err != nil {
		t.Fatalf("mark group bankrupt before restore: %v", err)
	}

	restoreResult, err := NewAdminRollbackCommandService(tx).RestoreGroupSnapshot(ctx, RestoreGroupSnapshotCommand{
		SnapshotID: snapshot.ID, Reason: "restore adjustment state", ConfirmText: rollbackConfirmText,
		OperatorID: 5, OperatorName: "snapshot-admin",
	})
	if err != nil {
		t.Fatalf("restore adjustment snapshot: %v", err)
	}
	if restoreResult.BusinessStatus != enum.BusinessStatusBankrupt {
		t.Fatalf("snapshot restore must preserve confirmed bankruptcy, got %s", restoreResult.BusinessStatus)
	}
	restoredAdjustment, err := repository.NewGroupAdjustmentRepository(tx).GetByID(ctx, sendResult.AdjustmentID)
	if err != nil {
		t.Fatalf("reload restored adjustment: %v", err)
	}
	if !restoredAdjustment.Effective {
		t.Fatalf("expected snapshot to restore adjustment effective state, got %#v", restoredAdjustment)
	}
	revision, err := repository.NewGroupAdjustmentRevisionRepository(tx).Get(ctx, groupID, 0)
	if err != nil {
		t.Fatalf("load revision after snapshot restore: %v", err)
	}
	if revision != 3 {
		t.Fatalf("expected revision 3 after snapshot restore, got %d", revision)
	}
	group, err := repository.NewGroupRepository(tx).GetByID(ctx, groupID)
	if err != nil {
		t.Fatalf("reload group after snapshot restore: %v", err)
	}
	if group.BusinessStatus != enum.BusinessStatusBankrupt || group.BankruptYearNo == nil || *group.BankruptYearNo != 0 {
		t.Fatalf("expected bankruptcy fields preserved after snapshot restore, got %#v", group)
	}
}

func buildPlayerNoticeService(db *gorm.DB) *PlayerNoticeService {
	return NewPlayerNoticeService(
		repository.NewNoticeRepository(db),
		repository.NewGroupAdjustmentRepository(db),
	)
}

func createAdminNoticeFixtures(t *testing.T, ctx context.Context, tx *gorm.DB, yearStatus string, stageStatus string) int64 {
	t.Helper()

	now := time.Now()
	groupID := createIntegrationGroup(t, ctx, tx, now, "IT_NOTICE")
	createInitialBaselineRecord(t, ctx, tx, groupID, now)
	createGroupYearStateRecord(t, ctx, tx, groupID, 0, enum.YearTypeDemo, yearStatus, stageStatus, enum.ReportStatusLocked)

	return groupID
}

func createLockedStageSubmission(t *testing.T, ctx context.Context, tx *gorm.DB, groupID int64, yearNo int, stageCode string) {
	t.Helper()

	now := time.Now()
	item := entity.GroupStageSubmission{
		GroupID:                  groupID,
		YearNo:                   yearNo,
		StageCode:                stageCode,
		SubmitVersion:            1,
		PeriodEndCash:            100,
		OperatingPayloadSnapshot: []byte(`{}`),
		StateBefore:              []byte(`{}`),
		StateAfter:               []byte(`{}`),
		SubmitterID:              1,
		SubmitTime:               now,
	}
	if err := tx.WithContext(ctx).Create(&item).Error; err != nil {
		t.Fatalf("create locked stage submission: %v", err)
	}
}
