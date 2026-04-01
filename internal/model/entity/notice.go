package entity

import "time"

type Notice struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	TargetScope   string    `gorm:"column:target_scope"`
	TargetGroupID *int64    `gorm:"column:target_group_id"`
	Content       string    `gorm:"column:content"`
	Pinned        bool      `gorm:"column:pinned"`
	PublishedAt   time.Time `gorm:"column:published_at"`
	OperatorID    int64     `gorm:"column:operator_id"`
	OperatorName  string    `gorm:"column:operator_name"`
	BaseEntity
}

func (Notice) TableName() string {
	return "sg_notice"
}
