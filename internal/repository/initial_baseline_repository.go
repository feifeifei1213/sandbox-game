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

type InitialBaselineRepository struct {
	db *gorm.DB
}

type UpsertSharedInitialBaselineCommand struct {
	GroupIDs        []int64
	BaselinePayload payload.BaselinePayload
	Submitted       bool
	SubmitterID     *int64
	SubmittedAt     *time.Time
	OperatorName    string
	OperateTime     time.Time
}

func NewInitialBaselineRepository(db *gorm.DB) *InitialBaselineRepository {
	return &InitialBaselineRepository{db: db}
}

func (r *InitialBaselineRepository) FindByGroupID(ctx context.Context, groupID int64) (*entity.InitialBaseline, error) {
	var item entity.InitialBaseline
	err := r.db.WithContext(ctx).
		Where("group_id = ?", groupID).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *InitialBaselineRepository) FindLatestSubmitted(ctx context.Context) (*entity.InitialBaseline, error) {
	var item entity.InitialBaseline
	err := r.db.WithContext(ctx).
		Where("submitted = ?", true).
		Order("submitted_at DESC").
		Order("id DESC").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *InitialBaselineRepository) FindTemplateSample(ctx context.Context) (*entity.InitialBaseline, error) {
	var item entity.InitialBaseline
	err := r.db.WithContext(ctx).
		Order("group_id ASC").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *InitialBaselineRepository) CountSubmitted(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.InitialBaseline{}).
		Where("submitted = ?", true).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *InitialBaselineRepository) UpsertSharedTemplate(ctx context.Context, cmd UpsertSharedInitialBaselineCommand) error {
	if len(cmd.GroupIDs) == 0 {
		return nil
	}

	rawPayload, err := json.Marshal(cmd.BaselinePayload)
	if err != nil {
		return fmt.Errorf("marshal baseline payload: %w", err)
	}

	items := make([]entity.InitialBaseline, 0, len(cmd.GroupIDs))
	for _, groupID := range cmd.GroupIDs {
		items = append(items, entity.InitialBaseline{
			GroupID:         groupID,
			BaselinePayload: rawPayload,
			Submitted:       cmd.Submitted,
			SubmitterID:     cmd.SubmitterID,
			SubmittedAt:     cmd.SubmittedAt,
			BaseEntity: entity.BaseEntity{
				Creator:    cmd.OperatorName,
				CreateTime: cmd.OperateTime,
				Updater:    cmd.OperatorName,
				UpdateTime: cmd.OperateTime,
			},
		})
	}

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "group_id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"baseline_payload_json": rawPayload,
				"submitted":             cmd.Submitted,
				"submitter_id":          cmd.SubmitterID,
				"submitted_at":          cmd.SubmittedAt,
				"updater":               cmd.OperatorName,
				"update_time":           cmd.OperateTime,
			}),
		}).
		Create(&items).Error
}
