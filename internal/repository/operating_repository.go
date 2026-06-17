package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
)

type OperatingRepository struct {
	db *gorm.DB
}

func NewOperatingRepository(db *gorm.DB) *OperatingRepository {
	return &OperatingRepository{db: db}
}

func (r *OperatingRepository) FindDraft(ctx context.Context, groupID int64, yearNo int) (*entity.GroupOperatingDraft, error) {
	var item entity.GroupOperatingDraft
	err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *OperatingRepository) ListDraftsByGroupFromYear(ctx context.Context, groupID int64, fromYearNo int) ([]entity.GroupOperatingDraft, error) {
	var items []entity.GroupOperatingDraft
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no >= ?", groupID, fromYearNo).
		Order("year_no ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OperatingRepository) ListStageSubmissions(ctx context.Context, groupID int64, yearNo int) ([]entity.GroupStageSubmission, error) {
	items := make([]entity.GroupStageSubmission, 0)
	err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Order("submit_time ASC, id ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OperatingRepository) ListStageSubmissionsByGroupFromYear(ctx context.Context, groupID int64, fromYearNo int) ([]entity.GroupStageSubmission, error) {
	var items []entity.GroupStageSubmission
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no >= ?", groupID, fromYearNo).
		Order("year_no ASC, submit_time ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OperatingRepository) HasStageSubmission(ctx context.Context, groupID int64, yearNo int, stageCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.GroupStageSubmission{}).
		Where("group_id = ? AND year_no = ? AND stage_code = ?", groupID, yearNo, stageCode).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *OperatingRepository) MaxStageSubmitVersion(ctx context.Context, groupID int64, yearNo int) (int, error) {
	var maxVersion int
	err := r.db.WithContext(ctx).
		Model(&entity.GroupStageSubmission{}).
		Select("COALESCE(MAX(submit_version), 0)").
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Scan(&maxVersion).Error
	return maxVersion, err
}

type UpsertOperatingDraftCommand struct {
	GroupID          int64
	YearNo           int
	StageStatus      string
	OperatingPayload payload.OperatingPayload
	LastAutoSavedAt  time.Time
	OperatorName     string
}

func (r *OperatingRepository) UpsertDraft(ctx context.Context, cmd UpsertOperatingDraftCommand) error {
	rawPayload, err := json.Marshal(cmd.OperatingPayload)
	if err != nil {
		return fmt.Errorf("marshal operating payload: %w", err)
	}

	item := entity.GroupOperatingDraft{
		GroupID:          cmd.GroupID,
		YearNo:           cmd.YearNo,
		StageStatus:      cmd.StageStatus,
		OperatingPayload: rawPayload,
		LastAutoSavedAt:  &cmd.LastAutoSavedAt,
		BaseEntity: entity.BaseEntity{
			Creator:    cmd.OperatorName,
			CreateTime: cmd.LastAutoSavedAt,
			Updater:    cmd.OperatorName,
			UpdateTime: cmd.LastAutoSavedAt,
		},
	}

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "group_id"},
				{Name: "year_no"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"stage_status":           item.StageStatus,
				"operating_payload_json": item.OperatingPayload,
				"last_auto_saved_at":     item.LastAutoSavedAt,
				"updater":                item.Updater,
				"update_time":            item.UpdateTime,
			}),
		}).
		Create(&item).Error
}

type CreateStageSubmissionCommand struct {
	GroupID                  int64
	YearNo                   int
	StageCode                string
	SubmitVersion            int
	PeriodEndCash            float64
	OperatingPayloadSnapshot payload.OperatingPayload
	StateBeforeJSON          []byte
	StateAfterJSON           []byte
	SubmitterID              int64
	SubmitTime               time.Time
}

func (r *OperatingRepository) CreateStageSubmission(ctx context.Context, cmd CreateStageSubmissionCommand) error {
	rawPayload, err := json.Marshal(cmd.OperatingPayloadSnapshot)
	if err != nil {
		return fmt.Errorf("marshal operating payload snapshot: %w", err)
	}

	item := entity.GroupStageSubmission{
		GroupID:                  cmd.GroupID,
		YearNo:                   cmd.YearNo,
		StageCode:                cmd.StageCode,
		SubmitVersion:            cmd.SubmitVersion,
		PeriodEndCash:            cmd.PeriodEndCash,
		OperatingPayloadSnapshot: rawPayload,
		StateBefore:              cmd.StateBeforeJSON,
		StateAfter:               cmd.StateAfterJSON,
		SubmitterID:              cmd.SubmitterID,
		SubmitTime:               cmd.SubmitTime,
	}

	return r.db.WithContext(ctx).Create(&item).Error
}
