package entity

type GroupYearState struct {
	ID                        int64   `gorm:"column:id;primaryKey"`
	GroupID                   int64   `gorm:"column:group_id"`
	YearNo                    int     `gorm:"column:year_no"`
	YearType                  string  `gorm:"column:year_type"`
	YearStatus                string  `gorm:"column:year_status"`
	StageStatus               string  `gorm:"column:stage_status"`
	ReportStatus              string  `gorm:"column:report_status"`
	SummaryEffective          bool    `gorm:"column:summary_effective"`
	LatestStageSubmitVersion  int     `gorm:"column:latest_stage_submit_version"`
	LatestReportSubmitVersion int     `gorm:"column:latest_report_submit_version"`
	RollbackPending           bool    `gorm:"column:rollback_pending"`
	RollbackTargetYearNo      *int    `gorm:"column:rollback_target_year_no"`
	RollbackTargetStageCode   *string `gorm:"column:rollback_target_stage_code"`
	RollbackLogID             *int64  `gorm:"column:rollback_log_id"`
	BaseEntity
}

func (GroupYearState) TableName() string {
	return "sg_group_year_state"
}
