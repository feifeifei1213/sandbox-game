package dto

import (
	"time"

	"sandbox-game/internal/model/payload"
)

type PlayerOperatingGetYearViewRequest struct {
	YearNo *int `form:"yearNo" binding:"required"`
}

type PlayerOperatingSaveDraftRequest struct {
	YearNo           *int                     `json:"yearNo" binding:"required"`
	StageStatus      string                   `json:"stageStatus" binding:"required"`
	OperatingPayload payload.OperatingPayload `json:"operatingPayload"`
	ClientSaveTime   *time.Time               `json:"clientSaveTime"`
}

type PlayerOperatingSaveDraftResponse struct {
	GroupID          int64     `json:"groupId"`
	YearNo           int       `json:"yearNo"`
	StageStatus      string    `json:"stageStatus"`
	LastDraftSavedAt time.Time `json:"lastDraftSavedAt"`
}

type PlayerOperatingSubmitStageRequest struct {
	YearNo           *int                     `json:"yearNo" binding:"required"`
	StageCode        string                   `json:"stageCode" binding:"required"`
	OperatingPayload payload.OperatingPayload `json:"operatingPayload"`
}

type PlayerOperatingSubmitStageResponse struct {
	GroupID                  int64     `json:"groupId"`
	YearNo                   int       `json:"yearNo"`
	StageCode                string    `json:"stageCode"`
	YearStatus               string    `json:"yearStatus"`
	StageStatus              string    `json:"stageStatus"`
	ReportStatus             string    `json:"reportStatus"`
	BusinessStatus           string    `json:"businessStatus"`
	LatestStageSubmitVersion int       `json:"latestStageSubmitVersion"`
	PeriodEndCash            float64   `json:"periodEndCash"`
	SubmittedAt              time.Time `json:"submittedAt"`
}

type PlayerOperatingStageSubmitHistory struct {
	StageCode     string    `json:"stageCode"`
	SubmitVersion int       `json:"submitVersion"`
	PeriodEndCash float64   `json:"periodEndCash"`
	SubmitTime    time.Time `json:"submitTime"`
}
