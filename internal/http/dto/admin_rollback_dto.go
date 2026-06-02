package dto

type AdminRollbackListSnapshotsRequest struct {
	SnapshotScope *string `form:"snapshotScope"`
	SnapshotType  *string `form:"snapshotType"`
	GroupID       *int64  `form:"groupId"`
	YearNo        *int    `form:"yearNo"`
	StageCode     *string `form:"stageCode"`
	PageNo        *int    `form:"pageNo"`
	PageSize      *int    `form:"pageSize"`
}

type AdminRollbackGetSnapshotDetailRequest struct {
	SnapshotID *int64 `form:"snapshotId"`
}

type AdminRollbackCreateSnapshotRequest struct {
	SnapshotScope string `json:"snapshotScope"`
	GroupID       *int64 `json:"groupId"`
	YearNo        *int   `json:"yearNo"`
	StageCode     string `json:"stageCode"`
	Description   string `json:"description"`
}

type AdminRollbackRestoreGroupSnapshotRequest struct {
	SnapshotID  *int64 `json:"snapshotId"`
	Reason      string `json:"reason"`
	ConfirmText string `json:"confirmText"`
}
