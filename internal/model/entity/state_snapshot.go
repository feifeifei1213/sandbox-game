package entity

import "time"

type StateSnapshot struct {
	ID                 int64     `gorm:"column:id;primaryKey"`
	SnapshotScope      string    `gorm:"column:snapshot_scope"`
	SnapshotType       string    `gorm:"column:snapshot_type"`
	TriggerCode        string    `gorm:"column:trigger_code"`
	TargetGroupID      *int64    `gorm:"column:target_group_id"`
	TargetYearNo       *int      `gorm:"column:target_year_no"`
	TargetStageCode    *string   `gorm:"column:target_stage_code"`
	TargetReportStatus *string   `gorm:"column:target_report_status"`
	Description        *string   `gorm:"column:description"`
	PayloadHash        string    `gorm:"column:payload_hash"`
	CreatedByID        int64     `gorm:"column:created_by_id"`
	CreatedByName      string    `gorm:"column:created_by_name"`
	CreatedAt          time.Time `gorm:"column:created_at"`
}

func (StateSnapshot) TableName() string {
	return "sg_state_snapshot"
}

type StateSnapshotPayload struct {
	ID             int64  `gorm:"column:id;primaryKey"`
	SnapshotID     int64  `gorm:"column:snapshot_id"`
	PayloadVersion string `gorm:"column:payload_version"`
	PayloadJSON    []byte `gorm:"column:payload_json"`
	PayloadSize    int    `gorm:"column:payload_size"`
	BaseEntity
}

func (StateSnapshotPayload) TableName() string {
	return "sg_state_snapshot_payload"
}

type RollbackLog struct {
	ID               int64     `gorm:"column:id;primaryKey"`
	RollbackType     string    `gorm:"column:rollback_type"`
	TargetGroupID    int64     `gorm:"column:target_group_id"`
	TargetYearNo     int       `gorm:"column:target_year_no"`
	TargetStageCode  *string   `gorm:"column:target_stage_code"`
	SnapshotID       *int64    `gorm:"column:snapshot_id"`
	SafetySnapshotID int64     `gorm:"column:safety_snapshot_id"`
	Reason           string    `gorm:"column:reason"`
	StateBefore      []byte    `gorm:"column:state_before_json"`
	StateAfter       []byte    `gorm:"column:state_after_json"`
	OperatorID       int64     `gorm:"column:operator_id"`
	OperatorName     string    `gorm:"column:operator_name"`
	OperateTime      time.Time `gorm:"column:operate_time"`
}

func (RollbackLog) TableName() string {
	return "sg_rollback_log"
}
