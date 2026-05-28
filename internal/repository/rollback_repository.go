package repository

import (
	"context"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
)

type StateSnapshotRepository struct {
	db *gorm.DB
}

type StateSnapshotListFilter struct {
	SnapshotScope string
	SnapshotType  string
	GroupID       *int64
	YearNo        *int
	StageCode     string
	Limit         int
	Offset        int
}

func NewStateSnapshotRepository(db *gorm.DB) *StateSnapshotRepository {
	return &StateSnapshotRepository{db: db}
}

func (r *StateSnapshotRepository) CreateWithPayload(ctx context.Context, snapshot *entity.StateSnapshot, payload *entity.StateSnapshotPayload) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(snapshot).Error; err != nil {
			return err
		}
		payload.SnapshotID = snapshot.ID
		return tx.Create(payload).Error
	})
}

func (r *StateSnapshotRepository) GetByID(ctx context.Context, id int64) (*entity.StateSnapshot, error) {
	var item entity.StateSnapshot
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *StateSnapshotRepository) GetPayloadBySnapshotID(ctx context.Context, snapshotID int64) (*entity.StateSnapshotPayload, error) {
	var item entity.StateSnapshotPayload
	if err := r.db.WithContext(ctx).
		Where("snapshot_id = ?", snapshotID).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *StateSnapshotRepository) List(ctx context.Context, filter StateSnapshotListFilter) ([]entity.StateSnapshot, error) {
	query := r.db.WithContext(ctx).Model(&entity.StateSnapshot{})
	query = applyStateSnapshotFilter(query, filter)
	query = query.Order("created_at DESC").Order("id DESC")
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var items []entity.StateSnapshot
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *StateSnapshotRepository) Count(ctx context.Context, filter StateSnapshotListFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.StateSnapshot{})
	query = applyStateSnapshotFilter(query, filter)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func applyStateSnapshotFilter(query *gorm.DB, filter StateSnapshotListFilter) *gorm.DB {
	if filter.SnapshotScope != "" {
		query = query.Where("snapshot_scope = ?", filter.SnapshotScope)
	}
	if filter.SnapshotType != "" {
		query = query.Where("snapshot_type = ?", filter.SnapshotType)
	}
	if filter.GroupID != nil {
		query = query.Where("target_group_id = ?", *filter.GroupID)
	}
	if filter.YearNo != nil {
		query = query.Where("target_year_no = ?", *filter.YearNo)
	}
	if filter.StageCode != "" {
		query = query.Where("target_stage_code = ?", filter.StageCode)
	}
	return query
}

type RollbackLogRepository struct {
	db *gorm.DB
}

func NewRollbackLogRepository(db *gorm.DB) *RollbackLogRepository {
	return &RollbackLogRepository{db: db}
}

func (r *RollbackLogRepository) Create(ctx context.Context, item *entity.RollbackLog) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *RollbackLogRepository) GetByID(ctx context.Context, id int64) (*entity.RollbackLog, error) {
	var item entity.RollbackLog
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RollbackLogRepository) ListByGroupYear(ctx context.Context, groupID int64, yearNo int) ([]entity.RollbackLog, error) {
	var items []entity.RollbackLog
	if err := r.db.WithContext(ctx).
		Where("target_group_id = ? AND target_year_no = ?", groupID, yearNo).
		Order("operate_time DESC").
		Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
