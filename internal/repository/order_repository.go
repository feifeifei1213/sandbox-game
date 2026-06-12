package repository

import (
	"context"
	"errors"
	"time"

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

func (r *OrderGenerationConfigRepository) ListAllFromYear(ctx context.Context, fromYearNo int) ([]entity.OrderGenerationConfig, error) {
	var items []entity.OrderGenerationConfig
	if err := r.db.WithContext(ctx).
		Where("year_no >= ?", fromYearNo).
		Order("year_no ASC").
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

func (r *OrderGenerationConfigRepository) UpdateStatusAndBatchByYear(ctx context.Context, yearNo int, status string, batchID int64, updater string) error {
	return r.db.WithContext(ctx).
		Model(&entity.OrderGenerationConfig{}).
		Where("year_no = ?", yearNo).
		Updates(map[string]any{
			"config_status":       status,
			"generation_batch_id": batchID,
			"updater":             updater,
		}).Error
}

type OrderForecastControlRepository struct {
	db *gorm.DB
}

func NewOrderForecastControlRepository(db *gorm.DB) *OrderForecastControlRepository {
	return &OrderForecastControlRepository{db: db}
}

func (r *OrderForecastControlRepository) ListAll(ctx context.Context) ([]entity.OrderForecastControl, error) {
	var items []entity.OrderForecastControl
	if err := r.db.WithContext(ctx).
		Order("year_no ASC").
		Order("market_code ASC").
		Order("order_type ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderForecastControlRepository) ListByYear(ctx context.Context, yearNo int) ([]entity.OrderForecastControl, error) {
	var items []entity.OrderForecastControl
	if err := r.db.WithContext(ctx).
		Where("year_no = ?", yearNo).
		Order("market_code ASC").
		Order("order_type ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderForecastControlRepository) UpsertBatch(ctx context.Context, items []entity.OrderForecastControl) error {
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
			"forecast_stage_code",
			"order_count",
			"updater",
			"update_time",
		}),
	}).Create(&items).Error
}

type OrderMarketForecastRepository struct {
	db *gorm.DB
}

func NewOrderMarketForecastRepository(db *gorm.DB) *OrderMarketForecastRepository {
	return &OrderMarketForecastRepository{db: db}
}

func (r *OrderMarketForecastRepository) ListAll(ctx context.Context) ([]entity.OrderMarketForecast, error) {
	var items []entity.OrderMarketForecast
	if err := r.db.WithContext(ctx).
		Order("forecast_stage_code ASC").
		Order("market_code ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderMarketForecastRepository) UpsertBatch(ctx context.Context, items []entity.OrderMarketForecast) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "forecast_stage_code"},
			{Name: "market_code"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"forecast_data_json",
			"narrative",
			"formula_version",
			"random_seed",
			"control_snapshot_json",
			"updater",
			"update_time",
		}),
	}).Create(&items).Error
}

type OrderMarketConfigRepository struct {
	db *gorm.DB
}

func NewOrderMarketConfigRepository(db *gorm.DB) *OrderMarketConfigRepository {
	return &OrderMarketConfigRepository{db: db}
}

func (r *OrderMarketConfigRepository) ListByYear(ctx context.Context, yearNo int) ([]entity.OrderMarketConfig, error) {
	var items []entity.OrderMarketConfig
	if err := r.db.WithContext(ctx).
		Where("year_no = ?", yearNo).
		Order("market_code ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderMarketConfigRepository) ListAllFromYear(ctx context.Context, fromYearNo int) ([]entity.OrderMarketConfig, error) {
	var items []entity.OrderMarketConfig
	if err := r.db.WithContext(ctx).
		Where("year_no >= ?", fromYearNo).
		Order("year_no ASC").
		Order("market_code ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderMarketConfigRepository) UpsertBatch(ctx context.Context, items []entity.OrderMarketConfig) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "year_no"},
			{Name: "market_code"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"market_enabled",
			"market_investment_limit",
			"config_status",
			"locked_batch_id",
			"updater",
			"update_time",
		}),
	}).Create(&items).Error
}

func (r *OrderMarketConfigRepository) UpdateStatusAndBatchByYear(ctx context.Context, yearNo int, status string, batchID int64, updater string) error {
	return r.db.WithContext(ctx).
		Model(&entity.OrderMarketConfig{}).
		Where("year_no = ?", yearNo).
		Updates(map[string]any{
			"config_status":   status,
			"locked_batch_id": batchID,
			"updater":         updater,
		}).Error
}

type OrderGenerationBatchRepository struct {
	db *gorm.DB
}

func NewOrderGenerationBatchRepository(db *gorm.DB) *OrderGenerationBatchRepository {
	return &OrderGenerationBatchRepository{db: db}
}

func (r *OrderGenerationBatchRepository) Create(ctx context.Context, item *entity.OrderGenerationBatch) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *OrderGenerationBatchRepository) GetByIDForUpdate(ctx context.Context, id int64) (*entity.OrderGenerationBatch, error) {
	var item entity.OrderGenerationBatch
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *OrderGenerationBatchRepository) FindLatestByYearStatuses(ctx context.Context, yearNo int, statuses []string) (*entity.OrderGenerationBatch, error) {
	var item entity.OrderGenerationBatch
	if err := r.db.WithContext(ctx).
		Where("year_no = ? AND batch_status IN ?", yearNo, statuses).
		Order("generated_at DESC").
		Order("id DESC").
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *OrderGenerationBatchRepository) FindConfirmedByYear(ctx context.Context, yearNo int) (*entity.OrderGenerationBatch, error) {
	return r.FindLatestByYearStatuses(ctx, yearNo, []string{enum.OrderGenerationBatchStatusConfirmed})
}

func (r *OrderGenerationBatchRepository) ListFromYear(ctx context.Context, fromYearNo int) ([]entity.OrderGenerationBatch, error) {
	var items []entity.OrderGenerationBatch
	if err := r.db.WithContext(ctx).
		Where("year_no >= ?", fromYearNo).
		Order("year_no ASC").
		Order("generated_at DESC").
		Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderGenerationBatchRepository) VoidPreviewByYear(ctx context.Context, yearNo int, updater string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.OrderGenerationBatch{}).
		Where("year_no = ? AND batch_status = ?", yearNo, enum.OrderGenerationBatchStatusPreview).
		Updates(map[string]any{
			"batch_status": enum.OrderGenerationBatchStatusVoid,
			"updater":      updater,
			"update_time":  operateTime,
		}).Error
}

func (r *OrderGenerationBatchRepository) Confirm(ctx context.Context, id int64, operatorID int64, operatorName string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.OrderGenerationBatch{}).
		Where("id = ? AND batch_status = ?", id, enum.OrderGenerationBatchStatusPreview).
		Updates(map[string]any{
			"batch_status":      enum.OrderGenerationBatchStatusConfirmed,
			"confirmed_by_id":   operatorID,
			"confirmed_by_name": operatorName,
			"confirmed_at":      operateTime,
			"updater":           operatorName,
			"update_time":       operateTime,
		}).Error
}

func (r *OrderGenerationBatchRepository) UpdateGenerationPayload(ctx context.Context, id int64, detail []byte, generatedCount int, updater string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.OrderGenerationBatch{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"order_detail_json":     detail,
			"generated_order_count": generatedCount,
			"updater":               updater,
			"update_time":           operateTime,
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
		Order("market_code ASC").
		Order("order_type ASC").
		Order("card_sequence_no ASC").
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderPoolRepository) ListByYearWithOptionalFilters(ctx context.Context, yearNo int, marketCode string, orderType string) ([]entity.OrderPool, error) {
	query := r.db.WithContext(ctx).Where("year_no = ?", yearNo)
	if marketCode != "" {
		query = query.Where("market_code = ?", marketCode)
	}
	if orderType != "" {
		query = query.Where("order_type = ?", orderType)
	}
	var items []entity.OrderPool
	if err := query.
		Order("market_code ASC").
		Order("order_type ASC").
		Order("card_sequence_no ASC").
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderPoolRepository) ListFromYear(ctx context.Context, fromYearNo int) ([]entity.OrderPool, error) {
	var items []entity.OrderPool
	if err := r.db.WithContext(ctx).
		Where("year_no >= ?", fromYearNo).
		Order("year_no ASC").
		Order("market_code ASC").
		Order("order_type ASC").
		Order("card_sequence_no ASC").
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *OrderPoolRepository) GetByIDForUpdate(ctx context.Context, id int64) (*entity.OrderPool, error) {
	var item entity.OrderPool
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
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

func (r *OrderPoolRepository) CountAvailableBySegment(ctx context.Context, yearNo int, marketCode string, orderType string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.OrderPool{}).
		Where("year_no = ? AND market_code = ? AND order_type = ? AND pool_status = ?", yearNo, marketCode, orderType, enum.OrderPoolStatusAvailable).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrderPoolRepository) SumSelectedAmountByYearMarket(ctx context.Context, yearNo int, marketCode string) (map[int64]float64, error) {
	type row struct {
		GroupID int64   `gorm:"column:selected_group_id"`
		Amount  float64 `gorm:"column:amount"`
	}
	var rows []row
	if err := r.db.WithContext(ctx).
		Model(&entity.OrderPool{}).
		Select("selected_group_id, SUM(order_amount) AS amount").
		Where("year_no = ? AND market_code = ? AND pool_status = ? AND selected_group_id IS NOT NULL", yearNo, marketCode, enum.OrderPoolStatusSelected).
		Group("selected_group_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[int64]float64, len(rows))
	for _, item := range rows {
		result[item.GroupID] = item.Amount
	}
	return result, nil
}

func (r *OrderPoolRepository) MarkSelected(ctx context.Context, orderID int64, groupID int64, operatorName string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.OrderPool{}).
		Where("id = ? AND pool_status = ?", orderID, enum.OrderPoolStatusAvailable).
		Updates(map[string]any{
			"pool_status":       enum.OrderPoolStatusSelected,
			"selected_group_id": groupID,
			"selected_at":       operateTime,
			"updater":           operatorName,
			"update_time":       operateTime,
		}).Error
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

func (r *MarketBiddingStateRepository) ListFromYear(ctx context.Context, fromYearNo int) ([]entity.MarketBiddingState, error) {
	var items []entity.MarketBiddingState
	if err := r.db.WithContext(ctx).
		Where("year_no >= ?", fromYearNo).
		Order("year_no ASC").
		Order("release_sequence_no ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MarketBiddingStateRepository) ListByYearAndMarket(ctx context.Context, yearNo int, marketCode string) ([]entity.MarketBiddingState, error) {
	var items []entity.MarketBiddingState
	if err := r.db.WithContext(ctx).
		Where("year_no = ? AND market_code = ?", yearNo, marketCode).
		Order("release_sequence_no ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MarketBiddingStateRepository) GetBySegment(ctx context.Context, yearNo int, marketCode string, orderType string) (*entity.MarketBiddingState, error) {
	var item entity.MarketBiddingState
	if err := r.db.WithContext(ctx).
		Where("year_no = ? AND market_code = ? AND order_type = ?", yearNo, marketCode, orderType).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MarketBiddingStateRepository) GetBySegmentForUpdate(ctx context.Context, yearNo int, marketCode string, orderType string) (*entity.MarketBiddingState, error) {
	var item entity.MarketBiddingState
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("year_no = ? AND market_code = ? AND order_type = ?", yearNo, marketCode, orderType).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MarketBiddingStateRepository) GetCurrentSelecting(ctx context.Context, yearNo int) (*entity.MarketBiddingState, error) {
	var item entity.MarketBiddingState
	if err := r.db.WithContext(ctx).
		Where("year_no = ? AND segment_status = ?", yearNo, enum.OrderSegmentStatusSelecting).
		Order("release_sequence_no ASC").
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MarketBiddingStateRepository) GetNextReleasable(ctx context.Context, yearNo int) (*entity.MarketBiddingState, error) {
	var item entity.MarketBiddingState
	if err := r.db.WithContext(ctx).
		Where("year_no = ? AND segment_status = ?", yearNo, enum.OrderSegmentStatusSequenceReady).
		Order("release_sequence_no ASC").
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
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

func (r *MarketBiddingStateRepository) ExistsSelectingBefore(ctx context.Context, yearNo int, releaseSequenceNo int) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.MarketBiddingState{}).
		Where("year_no = ? AND release_sequence_no < ? AND segment_status NOT IN ?", yearNo, releaseSequenceNo, []string{
			enum.OrderSegmentStatusCompleted,
			enum.OrderSegmentStatusSkipped,
			enum.OrderSegmentStatusMarketDisabled,
			enum.OrderSegmentStatusNoOrderConfig,
		}).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MarketBiddingStateRepository) DeleteByYear(ctx context.Context, yearNo int) error {
	return r.db.WithContext(ctx).Where("year_no = ?", yearNo).Delete(&entity.MarketBiddingState{}).Error
}

func (r *MarketBiddingStateRepository) UpdateMarketStatus(ctx context.Context, yearNo int, marketCode string, fromStatuses []string, toStatus string, operatorName string, operateTime time.Time) (int64, error) {
	updates := map[string]any{
		"segment_status": toStatus,
		"updater":        operatorName,
		"update_time":    operateTime,
	}
	if toStatus == enum.OrderSegmentStatusBidOpen {
		updates["opened_at"] = operateTime
	}
	if toStatus == enum.OrderSegmentStatusBidClosed {
		updates["closed_at"] = operateTime
	}
	tx := r.db.WithContext(ctx).
		Model(&entity.MarketBiddingState{}).
		Where("year_no = ? AND market_code = ? AND segment_status IN ?", yearNo, marketCode, fromStatuses).
		Updates(updates)
	return tx.RowsAffected, tx.Error
}

func (r *MarketBiddingStateRepository) MarkMarketSequenceReady(ctx context.Context, yearNo int, marketCode string, leaderGroupID *int64, leaderRule []byte, randomSeed string, operatorName string, operateTime time.Time) (int64, error) {
	tx := r.db.WithContext(ctx).
		Model(&entity.MarketBiddingState{}).
		Where("year_no = ? AND market_code = ? AND segment_status IN ?", yearNo, marketCode, []string{enum.OrderSegmentStatusBidClosed, enum.OrderSegmentStatusWaitingInvestment}).
		Updates(map[string]any{
			"segment_status":   enum.OrderSegmentStatusSequenceReady,
			"leader_group_id":  leaderGroupID,
			"leader_rule_json": leaderRule,
			"random_seed":      randomSeed,
			"updater":          operatorName,
			"update_time":      operateTime,
		})
	return tx.RowsAffected, tx.Error
}

func (r *MarketBiddingStateRepository) MarkMarketSkipped(ctx context.Context, yearNo int, marketCode string, operatorName string, operateTime time.Time) (int64, error) {
	tx := r.db.WithContext(ctx).
		Model(&entity.MarketBiddingState{}).
		Where("year_no = ? AND market_code = ? AND segment_status IN ?", yearNo, marketCode, []string{enum.OrderSegmentStatusBidClosed, enum.OrderSegmentStatusWaitingInvestment}).
		Updates(map[string]any{
			"segment_status": enum.OrderSegmentStatusSkipped,
			"completed_at":   operateTime,
			"updater":        operatorName,
			"update_time":    operateTime,
		})
	return tx.RowsAffected, tx.Error
}

func (r *MarketBiddingStateRepository) MarkSegmentSequenceReady(ctx context.Context, id int64, leaderGroupID *int64, leaderRule []byte, randomSeed string, operatorName string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.MarketBiddingState{}).
		Where("id = ? AND segment_status = ?", id, enum.OrderSegmentStatusWaitingInvestment).
		Updates(map[string]any{
			"segment_status":   enum.OrderSegmentStatusSequenceReady,
			"leader_group_id":  leaderGroupID,
			"leader_rule_json": leaderRule,
			"random_seed":      randomSeed,
			"updater":          operatorName,
			"update_time":      operateTime,
		}).Error
}

func (r *MarketBiddingStateRepository) MarkSegmentSkipped(ctx context.Context, id int64, leaderGroupID *int64, leaderRule []byte, randomSeed string, operatorName string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.MarketBiddingState{}).
		Where("id = ? AND segment_status = ?", id, enum.OrderSegmentStatusWaitingInvestment).
		Updates(map[string]any{
			"segment_status":   enum.OrderSegmentStatusSkipped,
			"leader_group_id":  leaderGroupID,
			"leader_rule_json": leaderRule,
			"random_seed":      randomSeed,
			"completed_at":     operateTime,
			"updater":          operatorName,
			"update_time":      operateTime,
		}).Error
}

func (r *MarketBiddingStateRepository) ReleaseSegment(ctx context.Context, id int64, currentGroupID *int64, operatorName string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.MarketBiddingState{}).
		Where("id = ? AND segment_status = ?", id, enum.OrderSegmentStatusSequenceReady).
		Updates(map[string]any{
			"segment_status":   enum.OrderSegmentStatusSelecting,
			"current_group_id": currentGroupID,
			"released_at":      operateTime,
			"updater":          operatorName,
			"update_time":      operateTime,
		}).Error
}

func (r *MarketBiddingStateRepository) CompleteSegment(ctx context.Context, id int64, operatorName string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.MarketBiddingState{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"segment_status":   enum.OrderSegmentStatusCompleted,
			"current_group_id": nil,
			"completed_at":     operateTime,
			"updater":          operatorName,
			"update_time":      operateTime,
		}).Error
}

func (r *MarketBiddingStateRepository) UpdateCurrentGroup(ctx context.Context, id int64, currentGroupID *int64, operatorName string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.MarketBiddingState{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"current_group_id": currentGroupID,
			"updater":          operatorName,
			"update_time":      operateTime,
		}).Error
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

type GroupMarketBidRepository struct {
	db *gorm.DB
}

func NewGroupMarketBidRepository(db *gorm.DB) *GroupMarketBidRepository {
	return &GroupMarketBidRepository{db: db}
}

func (r *GroupMarketBidRepository) Upsert(ctx context.Context, item *entity.GroupMarketBid) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "group_id"},
			{Name: "year_no"},
			{Name: "market_code"},
			{Name: "order_type"},
		},
		DoNothing: true,
	}).Create(item).Error
}

func (r *GroupMarketBidRepository) GetByGroupYearMarket(ctx context.Context, groupID int64, yearNo int, marketCode string) (*entity.GroupMarketBid, error) {
	var item entity.GroupMarketBid
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ? AND market_code = ?", groupID, yearNo, marketCode).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GroupMarketBidRepository) GetByGroupYearSegment(ctx context.Context, groupID int64, yearNo int, marketCode string, orderType string) (*entity.GroupMarketBid, error) {
	var item entity.GroupMarketBid
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ? AND market_code = ? AND order_type = ?", groupID, yearNo, marketCode, orderType).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GroupMarketBidRepository) ListByYearMarket(ctx context.Context, yearNo int, marketCode string) ([]entity.GroupMarketBid, error) {
	var items []entity.GroupMarketBid
	if err := r.db.WithContext(ctx).
		Where("year_no = ? AND market_code = ?", yearNo, marketCode).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupMarketBidRepository) ListByYear(ctx context.Context, yearNo int) ([]entity.GroupMarketBid, error) {
	var items []entity.GroupMarketBid
	if err := r.db.WithContext(ctx).
		Where("year_no = ?", yearNo).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupMarketBidRepository) ListByYearSegment(ctx context.Context, yearNo int, marketCode string, orderType string) ([]entity.GroupMarketBid, error) {
	var items []entity.GroupMarketBid
	if err := r.db.WithContext(ctx).
		Where("year_no = ? AND market_code = ? AND order_type = ?", yearNo, marketCode, orderType).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupMarketBidRepository) ListByGroupYear(ctx context.Context, groupID int64, yearNo int) ([]entity.GroupMarketBid, error) {
	var items []entity.GroupMarketBid
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupMarketBidRepository) ListByGroupFromYear(ctx context.Context, groupID int64, fromYearNo int) ([]entity.GroupMarketBid, error) {
	var items []entity.GroupMarketBid
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no >= ?", groupID, fromYearNo).
		Order("year_no ASC").
		Order("market_code ASC").
		Order("order_type ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupMarketBidRepository) CreateBatch(ctx context.Context, items []entity.GroupMarketBid) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(&items, 100).Error
}

func (r *GroupMarketBidRepository) CountByGroupYear(ctx context.Context, groupID int64, yearNo int) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.GroupMarketBid{}).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GroupMarketBidRepository) CountByYear(ctx context.Context, yearNo int) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.GroupMarketBid{}).
		Where("year_no = ?", yearNo).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GroupMarketBidRepository) CountSubmittedSegmentsByYear(ctx context.Context, yearNo int) (map[int64]int, error) {
	type row struct {
		GroupID int64 `gorm:"column:group_id"`
		Count   int   `gorm:"column:count"`
	}
	var rows []row
	if err := r.db.WithContext(ctx).
		Model(&entity.GroupMarketBid{}).
		Select("group_id, COUNT(*) AS count").
		Where("year_no = ?", yearNo).
		Group("group_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[int64]int, len(rows))
	for _, item := range rows {
		result[item.GroupID] = item.Count
	}
	return result, nil
}

type MarketSelectionOrderRepository struct {
	db *gorm.DB
}

func NewMarketSelectionOrderRepository(db *gorm.DB) *MarketSelectionOrderRepository {
	return &MarketSelectionOrderRepository{db: db}
}

func (r *MarketSelectionOrderRepository) CreateBatch(ctx context.Context, items []entity.MarketSelectionOrder) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(&items, 200).Error
}

func (r *MarketSelectionOrderRepository) DeleteByYearMarket(ctx context.Context, yearNo int, marketCode string) error {
	return r.db.WithContext(ctx).Where("year_no = ? AND market_code = ?", yearNo, marketCode).Delete(&entity.MarketSelectionOrder{}).Error
}

func (r *MarketSelectionOrderRepository) DeleteByYear(ctx context.Context, yearNo int) error {
	return r.db.WithContext(ctx).Where("year_no = ?", yearNo).Delete(&entity.MarketSelectionOrder{}).Error
}

func (r *MarketSelectionOrderRepository) ListBySegment(ctx context.Context, yearNo int, marketCode string, orderType string) ([]entity.MarketSelectionOrder, error) {
	var items []entity.MarketSelectionOrder
	if err := r.db.WithContext(ctx).
		Where("year_no = ? AND market_code = ? AND order_type = ?", yearNo, marketCode, orderType).
		Order("sequence_no ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MarketSelectionOrderRepository) ListByYear(ctx context.Context, yearNo int) ([]entity.MarketSelectionOrder, error) {
	var items []entity.MarketSelectionOrder
	if err := r.db.WithContext(ctx).
		Where("year_no = ?", yearNo).
		Order("market_code ASC").
		Order("order_type ASC").
		Order("sequence_no ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MarketSelectionOrderRepository) ListFromYear(ctx context.Context, fromYearNo int) ([]entity.MarketSelectionOrder, error) {
	var items []entity.MarketSelectionOrder
	if err := r.db.WithContext(ctx).
		Where("year_no >= ?", fromYearNo).
		Order("year_no ASC").
		Order("market_code ASC").
		Order("order_type ASC").
		Order("sequence_no ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MarketSelectionOrderRepository) GetByGroupSegmentForUpdate(ctx context.Context, groupID int64, yearNo int, marketCode string, orderType string) (*entity.MarketSelectionOrder, error) {
	var item entity.MarketSelectionOrder
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("group_id = ? AND year_no = ? AND market_code = ? AND order_type = ?", groupID, yearNo, marketCode, orderType).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MarketSelectionOrderRepository) GetNextWaiting(ctx context.Context, yearNo int, marketCode string, orderType string) (*entity.MarketSelectionOrder, error) {
	var item entity.MarketSelectionOrder
	if err := r.db.WithContext(ctx).
		Where("year_no = ? AND market_code = ? AND order_type = ? AND selection_status = ?", yearNo, marketCode, orderType, enum.OrderSelectionStatusWaiting).
		Order("sequence_no ASC").
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MarketSelectionOrderRepository) SetStatus(ctx context.Context, id int64, status string, selectedOrderID *int64, skippedByAdminID *int64, skippedReason *string, operateTime *time.Time, operatorName string) error {
	updates := map[string]any{
		"selection_status":  status,
		"selected_order_id": selectedOrderID,
		"updater":           operatorName,
	}
	if operateTime != nil {
		updates["selected_at"] = nil
		updates["skipped_at"] = nil
		switch status {
		case enum.OrderSelectionStatusSelected:
			updates["selected_at"] = *operateTime
		case enum.OrderSelectionStatusAdminSkipped:
			updates["skipped_at"] = *operateTime
		}
		updates["update_time"] = *operateTime
	}
	if skippedByAdminID != nil {
		updates["skipped_by_admin_id"] = *skippedByAdminID
	}
	if skippedReason != nil {
		updates["skipped_reason"] = *skippedReason
	}
	return r.db.WithContext(ctx).
		Model(&entity.MarketSelectionOrder{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *MarketSelectionOrderRepository) SetCurrent(ctx context.Context, id int64, operatorName string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.MarketSelectionOrder{}).
		Where("id = ? AND selection_status = ?", id, enum.OrderSelectionStatusWaiting).
		Updates(map[string]any{
			"selection_status": enum.OrderSelectionStatusCurrent,
			"updater":          operatorName,
			"update_time":      operateTime,
		}).Error
}

type GroupOrderSelectionRepository struct {
	db *gorm.DB
}

func NewGroupOrderSelectionRepository(db *gorm.DB) *GroupOrderSelectionRepository {
	return &GroupOrderSelectionRepository{db: db}
}

func (r *GroupOrderSelectionRepository) Create(ctx context.Context, item *entity.GroupOrderSelection) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GroupOrderSelectionRepository) ListByGroupYear(ctx context.Context, groupID int64, yearNo int) ([]entity.GroupOrderSelection, error) {
	var items []entity.GroupOrderSelection
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupOrderSelectionRepository) ListByGroupFromYear(ctx context.Context, groupID int64, fromYearNo int) ([]entity.GroupOrderSelection, error) {
	var items []entity.GroupOrderSelection
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no >= ?", groupID, fromYearNo).
		Order("year_no ASC").
		Order("market_code ASC").
		Order("order_type ASC").
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GroupOrderSelectionRepository) GetByGroupSegment(ctx context.Context, groupID int64, yearNo int, marketCode string, orderType string) (*entity.GroupOrderSelection, error) {
	var item entity.GroupOrderSelection
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ? AND market_code = ? AND order_type = ?", groupID, yearNo, marketCode, orderType).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

type GroupOrderAmountSummary struct {
	MarketCode string  `gorm:"column:market_code"`
	OrderType  string  `gorm:"column:order_type"`
	Amount     float64 `gorm:"column:amount"`
}

func (r *GroupOrderSelectionRepository) SumSelectedAmountByGroupYearSegment(ctx context.Context, groupID int64, yearNo int) ([]GroupOrderAmountSummary, error) {
	var rows []GroupOrderAmountSummary
	if err := r.db.WithContext(ctx).
		Table("sg_group_order_selection AS s").
		Select("s.market_code, s.order_type, COALESCE(SUM(p.order_amount), 0) AS amount").
		Joins("JOIN sg_order_pool AS p ON p.id = s.order_id").
		Where("s.group_id = ? AND s.year_no = ? AND s.selection_status = ?", groupID, yearNo, enum.OrderSelectionStatusSelected).
		Group("s.market_code, s.order_type").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

type GroupSelectedOrderDetail struct {
	SelectionID        int64   `gorm:"column:selection_id"`
	OrderID            int64   `gorm:"column:order_id"`
	MarketCode         string  `gorm:"column:market_code"`
	OrderType          string  `gorm:"column:order_type"`
	DeliveryStatus     string  `gorm:"column:delivery_status"`
	DeliveryEffective  bool    `gorm:"column:delivery_effective"`
	DeliveredStageCode *string `gorm:"column:delivered_stage_code"`
	OrderAmount        float64 `gorm:"column:order_amount"`
}

func (r *GroupOrderSelectionRepository) ListSelectedOrderDetailsForUpdate(ctx context.Context, groupID int64, yearNo int, orderIDs []int64) ([]GroupSelectedOrderDetail, error) {
	if len(orderIDs) == 0 {
		return []GroupSelectedOrderDetail{}, nil
	}
	var rows []GroupSelectedOrderDetail
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Table("sg_group_order_selection AS s").
		Select("s.id AS selection_id, s.order_id, s.market_code, s.order_type, s.delivery_status, s.delivery_effective, s.delivered_stage_code, p.order_amount").
		Joins("JOIN sg_order_pool AS p ON p.id = s.order_id").
		Where("s.group_id = ? AND s.year_no = ? AND s.order_id IN ? AND s.delivery_effective = ?", groupID, yearNo, orderIDs, true).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *GroupOrderSelectionRepository) SumDeliveredAmountByStage(ctx context.Context, groupID int64, yearNo int, stageCode string) (float64, error) {
	type row struct {
		Amount float64 `gorm:"column:amount"`
	}
	var result row
	if err := r.db.WithContext(ctx).
		Table("sg_group_order_selection AS s").
		Select("COALESCE(SUM(p.order_amount), 0) AS amount").
		Joins("JOIN sg_order_pool AS p ON p.id = s.order_id").
		Where("s.group_id = ? AND s.year_no = ? AND s.delivery_status = ? AND s.delivered_stage_code = ? AND s.delivery_effective = ?", groupID, yearNo, enum.OrderDeliveryStatusDelivered, stageCode, true).
		Scan(&result).Error; err != nil {
		return 0, err
	}
	return result.Amount, nil
}

func (r *GroupOrderSelectionRepository) MarkDeliveredBySelectionIDs(ctx context.Context, selectionIDs []int64, stageCode string, operatorName string, operateTime time.Time) error {
	if len(selectionIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&entity.GroupOrderSelection{}).
		Where("id IN ? AND delivery_status = ?", selectionIDs, enum.OrderDeliveryStatusSelected).
		Updates(map[string]any{
			"delivery_status":            enum.OrderDeliveryStatusDelivered,
			"delivery_effective":         true,
			"delivered_stage_code":       stageCode,
			"delivered_at":               operateTime,
			"invalidated_by_rollback_id": gorm.Expr("NULL"),
			"invalidated_at":             gorm.Expr("NULL"),
			"updater":                    operatorName,
			"update_time":                operateTime,
		}).Error
}

func (r *GroupOrderSelectionRepository) MarkUnfinishedByGroupYear(ctx context.Context, groupID int64, yearNo int, operatorName string, operateTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.GroupOrderSelection{}).
		Where("group_id = ? AND year_no = ? AND delivery_status = ?", groupID, yearNo, enum.OrderDeliveryStatusSelected).
		Updates(map[string]any{
			"delivery_status": enum.OrderDeliveryStatusUnfinished,
			"updater":         operatorName,
			"update_time":     operateTime,
		}).Error
}

func (r *GroupOrderSelectionRepository) InvalidateDeliveryAfterTarget(ctx context.Context, groupID int64, targetYearNo int, targetStageCode string, rollbackID int64, operatorName string, operateTime time.Time) (int64, error) {
	query := r.db.WithContext(ctx).
		Model(&entity.GroupOrderSelection{}).
		Where("group_id = ? AND delivery_effective = ? AND delivery_status = ?", groupID, true, enum.OrderDeliveryStatusDelivered)
	if targetStageCode == "" {
		query = query.Where("year_no > ?", targetYearNo)
	} else if afterStages := rollbackStagesAtOrAfter(targetStageCode); len(afterStages) > 0 {
		query = query.Where("(year_no > ? OR (year_no = ? AND delivered_stage_code IN ?))", targetYearNo, targetYearNo, afterStages)
	} else {
		query = query.Where("year_no > ?", targetYearNo)
	}
	tx := query.Updates(map[string]any{
		"delivery_effective":         false,
		"invalidated_by_rollback_id": rollbackID,
		"invalidated_at":             operateTime,
		"updater":                    operatorName,
		"update_time":                operateTime,
	})
	return tx.RowsAffected, tx.Error
}

func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
