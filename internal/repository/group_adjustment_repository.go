package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sandbox-game/internal/model/entity"
)

type GroupAdjustmentRepository struct {
	db *gorm.DB
}

func NewGroupAdjustmentRepository(db *gorm.DB) *GroupAdjustmentRepository {
	return &GroupAdjustmentRepository{db: db}
}

func (r *GroupAdjustmentRepository) Create(ctx context.Context, item *entity.GroupAdjustment) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GroupAdjustmentRepository) GetByID(ctx context.Context, adjustmentID int64) (*entity.GroupAdjustment, error) {
	var item entity.GroupAdjustment
	if err := r.db.WithContext(ctx).First(&item, adjustmentID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GroupAdjustmentRepository) GetByIDForUpdate(ctx context.Context, adjustmentID int64) (*entity.GroupAdjustment, error) {
	var item entity.GroupAdjustment
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&item, adjustmentID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GroupAdjustmentRepository) ListByGroupIDAndYear(ctx context.Context, groupID int64, yearNo int) ([]entity.GroupAdjustment, error) {
	items := make([]entity.GroupAdjustment, 0)
	err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ? AND effective = ?", groupID, yearNo, true).
		Order("published_at DESC").
		Order("id DESC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupAdjustmentRepository) ListAllByGroupIDAndYear(ctx context.Context, groupID int64, yearNo int) ([]entity.GroupAdjustment, error) {
	items := make([]entity.GroupAdjustment, 0)
	err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Order("published_at DESC").
		Order("id DESC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupAdjustmentRepository) Void(ctx context.Context, adjustmentID int64, operatorID int64, operatorName string, reason string, operateTime time.Time) (bool, error) {
	tx := r.db.WithContext(ctx).
		Model(&entity.GroupAdjustment{}).
		Where("id = ? AND effective = ?", adjustmentID, true).
		Updates(map[string]any{
			"effective":      false,
			"voided_by_id":   operatorID,
			"voided_by_name": operatorName,
			"void_reason":    reason,
			"voided_at":      operateTime,
			"updater":        operatorName,
			"update_time":    operateTime,
		})
	return tx.RowsAffected > 0, tx.Error
}

// RestoreEffectiveState 按快照恢复单条奖惩的有效状态。作废与既有回退审计字段不会被物理清除。
func (r *GroupAdjustmentRepository) RestoreEffectiveState(ctx context.Context, adjustmentID int64, effective bool, rollbackID int64, reason string, operatorName string, operateTime time.Time) (bool, error) {
	updates := map[string]any{
		"effective":   effective,
		"updater":     operatorName,
		"update_time": operateTime,
	}
	if !effective {
		updates["invalidated_by_rollback_id"] = rollbackID
		updates["invalid_reason"] = reason
		updates["invalidated_at"] = operateTime
	}
	tx := r.db.WithContext(ctx).
		Model(&entity.GroupAdjustment{}).
		Where("id = ? AND effective <> ?", adjustmentID, effective).
		Updates(updates)
	return tx.RowsAffected > 0, tx.Error
}

func (r *GroupAdjustmentRepository) ListAllByGroupFromYear(ctx context.Context, groupID int64, fromYearNo int) ([]entity.GroupAdjustment, error) {
	var items []entity.GroupAdjustment
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no >= ?", groupID, fromYearNo).
		Order("year_no ASC").
		Order("published_at DESC").
		Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupAdjustmentRepository) MarkInvalidAfterTarget(ctx context.Context, groupID int64, targetYearNo int, targetStageCode string, rollbackID int64, reason string, operatorName string, operateTime time.Time) (int64, error) {
	query := r.db.WithContext(ctx).
		Model(&entity.GroupAdjustment{}).
		Where("group_id = ? AND effective = ?", groupID, true)
	if targetStageCode == "" {
		query = query.Where("year_no > ?", targetYearNo)
	} else if afterStages := rollbackStagesAtOrAfter(targetStageCode); len(afterStages) > 0 {
		query = query.Where("(year_no > ? OR (year_no = ? AND stage_code IN ?))", targetYearNo, targetYearNo, afterStages)
	} else {
		query = query.Where("year_no > ?", targetYearNo)
	}
	tx := query.Updates(map[string]any{
		"effective":                  false,
		"invalidated_by_rollback_id": rollbackID,
		"invalid_reason":             reason,
		"invalidated_at":             operateTime,
		"updater":                    operatorName,
		"update_time":                operateTime,
	})
	return tx.RowsAffected, tx.Error
}

func (r *GroupAdjustmentRepository) ListRecent(ctx context.Context, limit int) ([]entity.GroupAdjustment, error) {
	if limit <= 0 {
		limit = 20
	}

	items := make([]entity.GroupAdjustment, 0, limit)
	err := r.db.WithContext(ctx).
		Order("published_at DESC").
		Order("id DESC").
		Limit(limit).
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
