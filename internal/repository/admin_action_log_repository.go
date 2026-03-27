package repository

import (
	"context"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
)

type AdminActionLogRepository struct {
	db *gorm.DB
}

func NewAdminActionLogRepository(db *gorm.DB) *AdminActionLogRepository {
	return &AdminActionLogRepository{db: db}
}

func (r *AdminActionLogRepository) Create(ctx context.Context, item *entity.AdminActionLog) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *AdminActionLogRepository) FindLatest(ctx context.Context) (*entity.AdminActionLog, error) {
	var item entity.AdminActionLog
	if err := r.db.WithContext(ctx).
		Order("operate_time DESC").
		Order("id DESC").
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
