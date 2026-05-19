package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
)

type OrderImportBatchRepository struct {
	db *gorm.DB
}

func NewOrderImportBatchRepository(db *gorm.DB) *OrderImportBatchRepository {
	return &OrderImportBatchRepository{db: db}
}

func (r *OrderImportBatchRepository) Create(ctx context.Context, item *entity.OrderImportBatch) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *OrderImportBatchRepository) GetByID(ctx context.Context, id int64) (*entity.OrderImportBatch, error) {
	var item entity.OrderImportBatch
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *OrderImportBatchRepository) FindLatestSuccess(ctx context.Context) (*entity.OrderImportBatch, error) {
	var item entity.OrderImportBatch
	if err := r.db.WithContext(ctx).
		Where("parse_status = ?", enum.OrderImportStatusSuccess).
		Order("uploaded_at DESC").
		Order("id DESC").
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

type OrderGenerationConfigRepository struct {
	db *gorm.DB
}

func NewOrderGenerationConfigRepository(db *gorm.DB) *OrderGenerationConfigRepository {
	return &OrderGenerationConfigRepository{db: db}
}

func (r *OrderGenerationConfigRepository) ListByYear(ctx context.Context, yearNo int) ([]entity.OrderGenerationConfig, error) {
	var items []entity.OrderGenerationConfig
	if err := r.db.WithContext(ctx).
		Where("year_no = ?", yearNo).
		Order("release_sequence_no ASC").
		Order("market_code ASC").
		Order("order_type ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderGenerationConfigRepository) DeleteByYear(ctx context.Context, yearNo int) error {
	return r.db.WithContext(ctx).Where("year_no = ?", yearNo).Delete(&entity.OrderGenerationConfig{}).Error
}

func (r *OrderGenerationConfigRepository) CreateBatch(ctx context.Context, items []entity.OrderGenerationConfig) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(&items, 100).Error
}

func (r *OrderGenerationConfigRepository) UpsertBatch(ctx context.Context, items []entity.OrderGenerationConfig) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "year_no"},
			{Name: "market_code"},
			{Name: "order_type"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"order_count",
			"release_sequence_no",
			"source_batch_id",
			"config_status",
			"updater",
			"update_time",
		}),
	}).Create(&items).Error
}

func (r *OrderGenerationConfigRepository) UpdateStatusByYear(ctx context.Context, yearNo int, status string, updater string) error {
	return r.db.WithContext(ctx).
		Model(&entity.OrderGenerationConfig{}).
		Where("year_no = ?", yearNo).
		Updates(map[string]any{
			"config_status": status,
			"updater":       updater,
		}).Error
}

type OrderPoolRepository struct {
	db *gorm.DB
}

func NewOrderPoolRepository(db *gorm.DB) *OrderPoolRepository {
	return &OrderPoolRepository{db: db}
}

func (r *OrderPoolRepository) ListBySegment(ctx context.Context, yearNo int, marketCode string, orderType string) ([]entity.OrderPool, error) {
	var items []entity.OrderPool
	if err := r.db.WithContext(ctx).
		Where("year_no = ? AND market_code = ? AND order_type = ?", yearNo, marketCode, orderType).
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderPoolRepository) CountByYear(ctx context.Context, yearNo int) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.OrderPool{}).
		Where("year_no = ?", yearNo).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderPoolRepository) CountBySegment(ctx context.Context, yearNo int, marketCode string, orderType string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.OrderPool{}).
		Where("year_no = ? AND market_code = ? AND order_type = ?", yearNo, marketCode, orderType).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderPoolRepository) HasSelectedByYear(ctx context.Context, yearNo int) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.OrderPool{}).
		Where("year_no = ? AND (pool_status = ? OR selected_group_id IS NOT NULL)", yearNo, enum.OrderPoolStatusSelected).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *OrderPoolRepository) DeleteByYear(ctx context.Context, yearNo int) error {
	return r.db.WithContext(ctx).Where("year_no = ?", yearNo).Delete(&entity.OrderPool{}).Error
}

func (r *OrderPoolRepository) CreateBatch(ctx context.Context, items []entity.OrderPool) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(&items, 200).Error
}

type MarketBiddingStateRepository struct {
	db *gorm.DB
}

func NewMarketBiddingStateRepository(db *gorm.DB) *MarketBiddingStateRepository {
	return &MarketBiddingStateRepository{db: db}
}

func (r *MarketBiddingStateRepository) ListByYear(ctx context.Context, yearNo int) ([]entity.MarketBiddingState, error) {
	var items []entity.MarketBiddingState
	if err := r.db.WithContext(ctx).
		Where("year_no = ?", yearNo).
		Order("release_sequence_no ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MarketBiddingStateRepository) HasStartedByYear(ctx context.Context, yearNo int) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.MarketBiddingState{}).
		Where(
			"year_no = ? AND (segment_status IN ? OR (segment_status = ? AND (released_at IS NOT NULL OR completed_at IS NOT NULL)))",
			yearNo,
			[]string{enum.OrderSegmentStatusSelecting, enum.OrderSegmentStatusCompleted},
			enum.OrderSegmentStatusSkipped,
		).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MarketBiddingStateRepository) DeleteByYear(ctx context.Context, yearNo int) error {
	return r.db.WithContext(ctx).Where("year_no = ?", yearNo).Delete(&entity.MarketBiddingState{}).Error
}

func (r *MarketBiddingStateRepository) UpsertBatch(ctx context.Context, items []entity.MarketBiddingState) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "year_no"},
			{Name: "market_code"},
			{Name: "order_type"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"segment_code",
			"release_sequence_no",
			"segment_status",
			"updater",
			"update_time",
		}),
	}).Create(&items).Error
}

func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
