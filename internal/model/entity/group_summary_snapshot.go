package entity

import "time"

type GroupSummarySnapshot struct {
	ID                        int64      `gorm:"column:id;primaryKey"`
	GroupID                   int64      `gorm:"column:group_id"`
	YearNo                    int        `gorm:"column:year_no"`
	Revenue                   float64    `gorm:"column:revenue"`
	Profit                    float64    `gorm:"column:profit"`
	Equity                    float64    `gorm:"column:equity"`
	BusinessStatus            string     `gorm:"column:business_status"`
	RankingValue              float64    `gorm:"column:ranking_value"`
	SummaryEffective          bool       `gorm:"column:summary_effective"`
	SourceReportSubmitVersion int        `gorm:"column:source_report_submit_version"`
	InvalidatedByRollbackID   *int64     `gorm:"column:invalidated_by_rollback_id"`
	InvalidatedAt             *time.Time `gorm:"column:invalidated_at"`
	BaseEntity
}

func (GroupSummarySnapshot) TableName() string {
	return "sg_group_summary_snapshot"
}
