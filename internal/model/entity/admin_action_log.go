package entity

import "time"

type AdminActionLog struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	ActionCode    string    `gorm:"column:action_code"`
	TargetGroupID *int64    `gorm:"column:target_group_id"`
	TargetYearNo  *int      `gorm:"column:target_year_no"`
	ActionPayload []byte    `gorm:"column:action_payload_json"`
	StateBefore   []byte    `gorm:"column:state_before_json"`
	StateAfter    []byte    `gorm:"column:state_after_json"`
	OperatorID    int64     `gorm:"column:operator_id"`
	OperatorName  string    `gorm:"column:operator_name"`
	OperateTime   time.Time `gorm:"column:operate_time"`
}

func (AdminActionLog) TableName() string {
	return "sg_admin_action_log"
}
