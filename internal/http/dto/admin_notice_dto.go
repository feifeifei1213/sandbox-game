package dto

type AdminNoticeGetRecordsRequest struct {
	Limit *int `form:"limit"`
}

type AdminNoticeSendGeneralRequest struct {
	TargetScope   string `json:"targetScope"`
	TargetGroupID *int64 `json:"targetGroupId"`
	Content       string `json:"content"`
	Pinned        bool   `json:"pinned"`
}

type AdminNoticeSendAdjustmentRequest struct {
	GroupID        *int64   `json:"groupId" binding:"required"`
	YearNo         *int     `json:"yearNo" binding:"required"`
	StageCode      string   `json:"stageCode"`
	AdjustmentType string   `json:"adjustmentType"`
	Amount         *float64 `json:"amount" binding:"required"`
	Reason         string   `json:"reason"`
}
