package entity

import "time"

type GroupAdjustment struct {
	ID                      int64      `gorm:"column:id;primaryKey"`
	GroupID                 int64      `gorm:"column:group_id"`
	YearNo                  int        `gorm:"column:year_no"`
	StageCode               string     `gorm:"column:stage_code"`
	AdjustmentType          string     `gorm:"column:adjustment_type"`
	Amount                  float64    `gorm:"column:amount"`
	Reason                  string     `gorm:"column:reason"`
	Effective               bool       `gorm:"column:effective"`
	InvalidatedByRollbackID *int64     `gorm:"column:invalidated_by_rollback_id"`
	InvalidReason           *string    `gorm:"column:invalid_reason"`
	InvalidatedAt           *time.Time `gorm:"column:invalidated_at"`
	PublishedAt             time.Time  `gorm:"column:published_at"`
	OperatorID              int64      `gorm:"column:operator_id"`
	OperatorName            string     `gorm:"column:operator_name"`
	BaseEntity
}

func (GroupAdjustment) TableName() string {
	return "sg_group_adjustment"
}
