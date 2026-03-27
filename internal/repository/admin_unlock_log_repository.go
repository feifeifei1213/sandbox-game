package repository

import (
	"context"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
)

type AdminUnlockLogRepository struct {
	db *gorm.DB
}

func NewAdminUnlockLogRepository(db *gorm.DB) *AdminUnlockLogRepository {
	return &AdminUnlockLogRepository{db: db}
}

func (r *AdminUnlockLogRepository) Create(ctx context.Context, item *entity.AdminUnlockLog) error {
	return r.db.WithContext(ctx).Create(item).Error
}
