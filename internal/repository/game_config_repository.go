package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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

func (r *GameConfigRepository) GetCurrentForUpdate(ctx context.Context) (*entity.GameConfig, error) {
	var item entity.GameConfig
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Order("id ASC").
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GameConfigRepository) UpdateFinalYear(ctx context.Context, id int64, finalYear int, operatorName string, updateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.GameConfig{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"final_year":  finalYear,
			"updater":     operatorName,
			"update_time": updateTime,
		}).Error
}

func (r *GameConfigRepository) UpdateCurrentOpenYear(ctx context.Context, id int64, currentOpenYear int, operatorName string, updateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.GameConfig{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"current_open_year": currentOpenYear,
			"updater":           operatorName,
			"update_time":       updateTime,
		}).Error
}

func (r *GameConfigRepository) UpdateInitialBaselineSubmitted(ctx context.Context, id int64, submitted bool, operatorName string, updateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.GameConfig{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"initial_baseline_submitted": submitted,
			"updater":                    operatorName,
			"update_time":                updateTime,
		}).Error
}
