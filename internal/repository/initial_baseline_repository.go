package repository

import (
	"context"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
)

type InitialBaselineRepository struct {
	db *gorm.DB
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
