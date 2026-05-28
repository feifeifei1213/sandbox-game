package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/state"
)

type GroupYearStateRepository struct {
	db *gorm.DB
}

type EnsureFormalYearStatesCommand struct {
	GroupIDs     []int64
	FromYear     int
	ToYear       int
	OperatorName string
	OperateTime  time.Time
}

type groupYearStateKey struct {
	GroupID int64
	YearNo  int
}

type groupYearStateIdentity struct {
	GroupID int64 `gorm:"column:group_id"`
	YearNo  int   `gorm:"column:year_no"`
}

func NewGroupYearStateRepository(db *gorm.DB) *GroupYearStateRepository {
	return &GroupYearStateRepository{db: db}
}

func (r *GroupYearStateRepository) GetByGroupIDAndYear(ctx context.Context, groupID int64, yearNo int) (*entity.GroupYearState, error) {
	var item entity.GroupYearState
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GroupYearStateRepository) ListByYear(ctx context.Context, yearNo int) ([]entity.GroupYearState, error) {
	var items []entity.GroupYearState
	if err := r.db.WithContext(ctx).
		Where("year_no = ?", yearNo).
		Find(&items).Error; err != nil {
		return items, err
	}
	return items, nil
}

func (r *GroupYearStateRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.GroupYearState{}).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GroupYearStateRepository) CreateBatch(ctx context.Context, items []entity.GroupYearState) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(&items, 200).Error
}

func (r *GroupYearStateRepository) EnsureFormalYearStates(ctx context.Context, cmd EnsureFormalYearStatesCommand) (int64, error) {
	if len(cmd.GroupIDs) == 0 || cmd.FromYear > cmd.ToYear {
		return 0, nil
	}

	existingKeys, err := r.listExistingKeys(ctx, cmd.GroupIDs, cmd.FromYear, cmd.ToYear)
	if err != nil {
		return 0, err
	}

	items := buildMissingFormalYearStates(cmd.GroupIDs, cmd.FromYear, cmd.ToYear, cmd.OperatorName, cmd.OperateTime, existingKeys)
	if len(items) == 0 {
		return 0, nil
	}

	tx := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "group_id"}, {Name: "year_no"}},
			DoNothing: true,
		}).
		CreateInBatches(&items, 200)
	return tx.RowsAffected, tx.Error
}

func (r *GroupYearStateRepository) UpdateRuntimeState(ctx context.Context, groupID int64, yearNo int, runtimeState state.RuntimeState, operatorName string) error {
	return r.db.WithContext(ctx).
		Model(&entity.GroupYearState{}).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Updates(map[string]any{
			"year_status":                  runtimeState.YearStatus,
			"stage_status":                 runtimeState.StageStatus,
			"report_status":                runtimeState.ReportStatus,
			"summary_effective":            runtimeState.SummaryEffective,
			"latest_stage_submit_version":  runtimeState.LatestStageSubmitVersion,
			"latest_report_submit_version": runtimeState.LatestReportSubmitVersion,
			"updater":                      operatorName,
		}).Error
}

func (r *GroupYearStateRepository) MarkRollbackPending(ctx context.Context, groupID int64, yearNo int, targetYearNo int, targetStageCode string, rollbackLogID int64, operatorName string) error {
	return r.db.WithContext(ctx).
		Model(&entity.GroupYearState{}).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Updates(map[string]any{
			"rollback_pending":           true,
			"rollback_target_year_no":    targetYearNo,
			"rollback_target_stage_code": nullableStringValue(targetStageCode),
			"rollback_log_id":            rollbackLogID,
			"updater":                    operatorName,
		}).Error
}

func (r *GroupYearStateRepository) ClearRollbackPending(ctx context.Context, groupID int64, yearNo int, operatorName string) error {
	return r.db.WithContext(ctx).
		Model(&entity.GroupYearState{}).
		Where("group_id = ? AND year_no = ?", groupID, yearNo).
		Updates(map[string]any{
			"rollback_pending":           false,
			"rollback_target_year_no":    nil,
			"rollback_target_stage_code": nil,
			"rollback_log_id":            nil,
			"updater":                    operatorName,
		}).Error
}

func (r *GroupYearStateRepository) CountRollbackPending(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.GroupYearState{}).
		Where("rollback_pending = ?", true).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GroupYearStateRepository) listExistingKeys(ctx context.Context, groupIDs []int64, fromYear int, toYear int) (map[groupYearStateKey]struct{}, error) {
	var items []groupYearStateIdentity
	if err := r.db.WithContext(ctx).
		Model(&entity.GroupYearState{}).
		Select("group_id", "year_no").
		Where("group_id IN ? AND year_no BETWEEN ? AND ?", groupIDs, fromYear, toYear).
		Find(&items).Error; err != nil {
		return nil, err
	}

	result := make(map[groupYearStateKey]struct{}, len(items))
	for _, item := range items {
		result[groupYearStateKey{GroupID: item.GroupID, YearNo: item.YearNo}] = struct{}{}
	}
	return result, nil
}

func buildMissingFormalYearStates(groupIDs []int64, fromYear int, toYear int, operatorName string, operateTime time.Time, existingKeys map[groupYearStateKey]struct{}) []entity.GroupYearState {
	if len(groupIDs) == 0 || fromYear > toYear {
		return nil
	}

	items := make([]entity.GroupYearState, 0, len(groupIDs)*(toYear-fromYear+1))
	for _, groupID := range groupIDs {
		for yearNo := fromYear; yearNo <= toYear; yearNo++ {
			key := groupYearStateKey{GroupID: groupID, YearNo: yearNo}
			if _, exists := existingKeys[key]; exists {
				continue
			}
			items = append(items, entity.GroupYearState{
				GroupID:                   groupID,
				YearNo:                    yearNo,
				YearType:                  enum.YearTypeFormal,
				YearStatus:                enum.YearStatusLocked,
				StageStatus:               enum.StageStatusQ1Open,
				ReportStatus:              enum.ReportStatusLocked,
				SummaryEffective:          false,
				LatestStageSubmitVersion:  0,
				LatestReportSubmitVersion: 0,
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: operateTime,
					Updater:    operatorName,
					UpdateTime: operateTime,
				},
			})
		}
	}
	return items
}
