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
	var items []entity.AdminActionLog
	result := r.db.WithContext(ctx).
		Order("operate_time DESC").
		Order("id DESC").
		Limit(1).
		Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &items[0], nil
}
