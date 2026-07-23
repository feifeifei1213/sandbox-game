package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sandbox-game/internal/model/entity"
)

type GroupAdjustmentRevisionRepository struct {
	db *gorm.DB
}

func NewGroupAdjustmentRevisionRepository(db *gorm.DB) *GroupAdjustmentRevisionRepository {
	return &GroupAdjustmentRevisionRepository{db: db}
}

func (r *GroupAdjustmentRevisionRepository) Get(ctx context.Context, groupID int64, yearNo int) (int64, error) {
	var items []entity.GroupAdjustmentRevision
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Order("id").
		Limit(1).
		Find(&items).Error; err != nil {
		return 0, err
	}
	if len(items) == 0 {
		return 0, nil
	}
	return items[0].Revision, nil
}

// Increment 在调用方事务中原子递增奖惩版本，新年份第一次变更后的版本为 1。
func (r *GroupAdjustmentRevisionRepository) Increment(ctx context.Context, groupID int64, yearNo int, operateTime time.Time) (int64, error) {
	item := entity.GroupAdjustmentRevision{
		GroupID:   groupID,
		YearNo:    yearNo,
		Revision:  1,
		UpdatedAt: operateTime,
	}
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "group_id"}, {Name: "year_no"}},
			DoUpdates: clause.Assignments(map[string]any{
				"revision":   gorm.Expr("revision + 1"),
				"updated_at": operateTime,
			}),
		}).
		Create(&item).Error; err != nil {
		return 0, err
	}
	return r.Get(ctx, groupID, yearNo)
}
