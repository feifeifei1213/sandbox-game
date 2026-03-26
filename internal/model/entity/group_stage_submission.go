package entity

import "time"

type GroupStageSubmission struct {
	ID                       int64     `gorm:"column:id;primaryKey"`
	GroupID                  int64     `gorm:"column:group_id"`
	YearNo                   int       `gorm:"column:year_no"`
	StageCode                string    `gorm:"column:stage_code"`
	SubmitVersion            int       `gorm:"column:submit_version"`
	PeriodEndCash            float64   `gorm:"column:period_end_cash"`
	OperatingPayloadSnapshot []byte    `gorm:"column:operating_payload_snapshot_json"`
	StateBefore              []byte    `gorm:"column:state_before_json"`
	StateAfter               []byte    `gorm:"column:state_after_json"`
	SubmitterID              int64     `gorm:"column:submitter_id"`
	SubmitTime               time.Time `gorm:"column:submit_time"`
}

func (GroupStageSubmission) TableName() string {
	return "sg_group_stage_submission"
}
