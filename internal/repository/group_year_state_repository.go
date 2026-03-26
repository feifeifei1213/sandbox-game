package repository

import (
	"context"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/state"
)

type GroupYearStateRepository struct {
	db *gorm.DB
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
		return nil, err
	}
	return items, nil
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
