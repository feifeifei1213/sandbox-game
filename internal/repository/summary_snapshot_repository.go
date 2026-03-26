package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/rules/summary"
)

type SummarySnapshotRepository struct {
	db *gorm.DB
}

func NewSummarySnapshotRepository(db *gorm.DB) *SummarySnapshotRepository {
	return &SummarySnapshotRepository{db: db}
}

type UpsertSummarySnapshotCommand struct {
	GroupID                   int64
	YearNo                    int
	SummaryResult             summary.SummaryResult
	SourceReportSubmitVersion int
	OperatorName              string
	OperateTime               time.Time
}

type SummarySnapshotWithGroup struct {
	GroupID        int64   `gorm:"column:group_id"`
	GroupNo        int     `gorm:"column:group_no"`
	GroupName      string  `gorm:"column:group_name"`
	Revenue        float64 `gorm:"column:revenue"`
	Profit         float64 `gorm:"column:profit"`
	Equity         float64 `gorm:"column:equity"`
	BusinessStatus string  `gorm:"column:business_status"`
	RankingValue   float64 `gorm:"column:ranking_value"`
}

// UpsertSnapshot 保存正式年份汇总快照，后续重提时可覆盖当前快照。
func (r *SummarySnapshotRepository) UpsertSnapshot(ctx context.Context, cmd UpsertSummarySnapshotCommand) error {
	item := entity.GroupSummarySnapshot{
		GroupID:                   cmd.GroupID,
		YearNo:                    cmd.YearNo,
		Revenue:                   cmd.SummaryResult.Revenue,
		Profit:                    cmd.SummaryResult.Profit,
		Equity:                    cmd.SummaryResult.Equity,
		BusinessStatus:            cmd.SummaryResult.BusinessStatus,
		RankingValue:              cmd.SummaryResult.RankingValue,
		SummaryEffective:          cmd.SummaryResult.SummaryEffective,
		SourceReportSubmitVersion: cmd.SourceReportSubmitVersion,
		BaseEntity: entity.BaseEntity{
			Creator:    cmd.OperatorName,
			CreateTime: cmd.OperateTime,
			Updater:    cmd.OperatorName,
			UpdateTime: cmd.OperateTime,
		},
	}

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "group_id"},
				{Name: "year_no"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"revenue":                      item.Revenue,
				"profit":                       item.Profit,
				"equity":                       item.Equity,
				"business_status":              item.BusinessStatus,
				"ranking_value":                item.RankingValue,
				"summary_effective":            item.SummaryEffective,
				"source_report_submit_version": item.SourceReportSubmitVersion,
				"updater":                      item.Updater,
				"update_time":                  item.UpdateTime,
			}),
		}).
		Create(&item).Error
}

func (r *SummarySnapshotRepository) ListEffectiveByYear(ctx context.Context, yearNo int) ([]SummarySnapshotWithGroup, error) {
	var items []SummarySnapshotWithGroup
	if err := r.db.WithContext(ctx).
		Table("sg_group_summary_snapshot AS s").
		Select("s.group_id, g.group_no, g.group_name, s.revenue, s.profit, s.equity, s.business_status, s.ranking_value").
		Joins("INNER JOIN sg_group AS g ON g.id = s.group_id").
		Where("s.year_no = ? AND s.summary_effective = ?", yearNo, true).
		Order("g.group_no ASC").
		Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *SummarySnapshotRepository) ListEffectiveRankingByYear(ctx context.Context, yearNo int) ([]SummarySnapshotWithGroup, error) {
	var items []SummarySnapshotWithGroup
	if err := r.db.WithContext(ctx).
		Table("sg_group_summary_snapshot AS s").
		Select("s.group_id, g.group_no, g.group_name, s.revenue, s.profit, s.equity, s.business_status, s.ranking_value").
		Joins("INNER JOIN sg_group AS g ON g.id = s.group_id").
		Where("s.year_no = ? AND s.summary_effective = ?", yearNo, true).
		Order("s.ranking_value DESC, g.group_no ASC").
		Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
