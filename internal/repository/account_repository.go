package repository

import (
	"context"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) GetByID(ctx context.Context, accountID int64) (*entity.Account, error) {
	var item entity.Account
	if err := r.db.WithContext(ctx).
		Where("id = ?", accountID).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
