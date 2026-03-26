package entity

import "time"

type GroupReportSubmission struct {
	ID                     int64     `gorm:"column:id;primaryKey"`
	GroupID                int64     `gorm:"column:group_id"`
	YearNo                 int       `gorm:"column:year_no"`
	SubmitVersion          int       `gorm:"column:submit_version"`
	ReportManualSnapshot   []byte    `gorm:"column:report_manual_snapshot_json"`
	ReportComputedSnapshot []byte    `gorm:"column:report_computed_snapshot_json"`
	BalanceCheckPassed     bool      `gorm:"column:balance_check_passed"`
	StateBefore            []byte    `gorm:"column:state_before_json"`
	StateAfter             []byte    `gorm:"column:state_after_json"`
	SubmitterID            int64     `gorm:"column:submitter_id"`
	SubmitTime             time.Time `gorm:"column:submit_time"`
}

func (GroupReportSubmission) TableName() string {
	return "sg_group_report_submission"
}
