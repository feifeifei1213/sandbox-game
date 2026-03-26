package repository

import (
	"context"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
)

type GameConfigRepository struct {
	db *gorm.DB
}

func NewGameConfigRepository(db *gorm.DB) *GameConfigRepository {
	return &GameConfigRepository{db: db}
}

func (r *GameConfigRepository) GetCurrent(ctx context.Context) (*entity.GameConfig, error) {
	var item entity.GameConfig
	if err := r.db.WithContext(ctx).Order("id ASC").First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
