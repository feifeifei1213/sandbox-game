package entity

import "time"

type AdminUnlockLog struct {
	ID               int64     `gorm:"column:id;primaryKey"`
	GroupID          int64     `gorm:"column:group_id"`
	YearNo           int       `gorm:"column:year_no"`
	Reason           string    `gorm:"column:reason"`
	UnlockTargetType string    `gorm:"column:unlock_target_type"`
	TargetStageCode  *string   `gorm:"column:target_stage_code"`
	SafetySnapshotID *int64    `gorm:"column:safety_snapshot_id"`
	StateBefore      []byte    `gorm:"column:state_before_json"`
	StateAfter       []byte    `gorm:"column:state_after_json"`
	OperatorID       int64     `gorm:"column:operator_id"`
	OperatorName     string    `gorm:"column:operator_name"`
	OperateTime      time.Time `gorm:"column:operate_time"`
}

func (AdminUnlockLog) TableName() string {
	return "sg_admin_unlock_log"
}
