package repository

import (
	"context"
	"time"

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

func (r *AccountRepository) GetByUsername(ctx context.Context, username string) (*entity.Account, error) {
	var item entity.Account
	if err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *AccountRepository) UpdateLastLoginTime(ctx context.Context, accountID int64, loginTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.Account{}).
		Where("id = ?", accountID).
		Update("last_login_time", loginTime).
		Error
}
