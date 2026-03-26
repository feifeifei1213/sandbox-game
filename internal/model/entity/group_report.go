package entity

import "time"

type GroupReport struct {
	ID                    int64      `gorm:"column:id;primaryKey"`
	GroupID               int64      `gorm:"column:group_id"`
	YearNo                int        `gorm:"column:year_no"`
	ReportManualPayload   []byte     `gorm:"column:report_manual_payload_json"`
	ReportComputedPayload []byte     `gorm:"column:report_computed_payload_json"`
	BalanceCheckPassed    bool       `gorm:"column:balance_check_passed"`
	LastAutoSavedAt       *time.Time `gorm:"column:last_auto_saved_at"`
	SubmittedAt           *time.Time `gorm:"column:submitted_at"`
	BaseEntity
}

func (GroupReport) TableName() string {
	return "sg_group_report"
}
