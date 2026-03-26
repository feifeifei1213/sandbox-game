package entity

import "time"

type InitialBaseline struct {
	ID              int64      `gorm:"column:id;primaryKey"`
	GroupID         int64      `gorm:"column:group_id"`
	BaselinePayload []byte     `gorm:"column:baseline_payload_json"`
	Submitted       bool       `gorm:"column:submitted"`
	SubmitterID     *int64     `gorm:"column:submitter_id"`
	SubmittedAt     *time.Time `gorm:"column:submitted_at"`
	BaseEntity
}

func (InitialBaseline) TableName() string {
	return "sg_initial_baseline"
}
