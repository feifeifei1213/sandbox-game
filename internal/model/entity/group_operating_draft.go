package entity

import "time"

type GroupOperatingDraft struct {
	ID               int64      `gorm:"column:id;primaryKey"`
	GroupID          int64      `gorm:"column:group_id"`
	YearNo           int        `gorm:"column:year_no"`
	StageStatus      string     `gorm:"column:stage_status"`
	OperatingPayload []byte     `gorm:"column:operating_payload_json"`
	LastAutoSavedAt  *time.Time `gorm:"column:last_auto_saved_at"`
	BaseEntity
}

func (GroupOperatingDraft) TableName() string {
	return "sg_group_operating_draft"
}
