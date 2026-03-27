package repository

import (
	"context"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) GetByID(ctx context.Context, groupID int64) (*entity.Group, error) {
	var item entity.Group
	if err := r.db.WithContext(ctx).Where("id = ?", groupID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GroupRepository) ListAll(ctx context.Context) ([]entity.Group, error) {
	var items []entity.Group
	if err := r.db.WithContext(ctx).
		Order("group_no ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupRepository) MarkBankrupt(ctx context.Context, groupID int64, yearNo int, reason string, operatorName string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Group{}).
		Where("id = ?", groupID).
		Updates(map[string]any{
			"business_status":  enum.BusinessStatusBankrupt,
			"bankrupt_year_no": yearNo,
			"bankrupt_reason":  reason,
			"updater":          operatorName,
		}).Error
}

func (r *GroupRepository) RecoverFromBankrupt(ctx context.Context, groupID int64, yearNo int, operatorName string) (bool, error) {
	tx := r.db.WithContext(ctx).
		Model(&entity.Group{}).
		Where("id = ? AND business_status = ? AND bankrupt_year_no = ?", groupID, enum.BusinessStatusBankrupt, yearNo).
		Updates(map[string]any{
			"business_status":  enum.BusinessStatusNormal,
			"bankrupt_year_no": nil,
			"bankrupt_reason":  nil,
			"updater":          operatorName,
		})
	return tx.RowsAffected > 0, tx.Error
}
