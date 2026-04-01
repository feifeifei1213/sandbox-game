package service

import (
	"context"
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
}

func TestAdminNoticeCommandServiceRejectsAdjustmentForLockedStage(t *testing.T) {
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

	_, err := commandService.SendAdjustment(ctx, SendAdjustmentCommand{
		GroupID:        groupID,
		YearNo:         0,
		StageCode:      state.StageCodeQ1,
		AdjustmentType: adjustmentTypePenalty,
		Amount:         66,
		Reason:         "I2 locked stage should reject",
		OperatorID:     1,
		OperatorName:   "integration-test",
	})
	if !errors.Is(err, ErrAdminAdjustmentStageLocked) {
		t.Fatalf("expected ErrAdminAdjustmentStageLocked, got %v", err)
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
