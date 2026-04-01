package repository

import (
	"context"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
)

type NoticeRepository struct {
	db *gorm.DB
}

func NewNoticeRepository(db *gorm.DB) *NoticeRepository {
	return &NoticeRepository{db: db}
}

func (r *NoticeRepository) Create(ctx context.Context, item *entity.Notice) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *NoticeRepository) ListForGroup(ctx context.Context, groupID int64, limit int) ([]entity.Notice, error) {
	if limit <= 0 {
		limit = 10
	}

	items := make([]entity.Notice, 0, limit)
	err := r.db.WithContext(ctx).
		Where("(target_scope = ? OR (target_scope = ? AND target_group_id = ?))", "ALL", "GROUP", groupID).
		Order("pinned DESC").
		Order("published_at DESC").
		Order("id DESC").
		Limit(limit).
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *NoticeRepository) ListRecent(ctx context.Context, limit int) ([]entity.Notice, error) {
	if limit <= 0 {
		limit = 20
	}

	items := make([]entity.Notice, 0, limit)
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
