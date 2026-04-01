package repository

import (
	"context"

	"gorm.io/gorm"

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

func (r *GroupAdjustmentRepository) ListByGroupIDAndYear(ctx context.Context, groupID int64, yearNo int) ([]entity.GroupAdjustment, error) {
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
