package dto

import "sandbox-game/internal/model/payload"

type AdminControlInitializeGameRequest struct {
	GroupCount *int `json:"groupCount" binding:"required"`
}

type AdminControlUpdateFinalYearRequest struct {
	FinalYear *int `json:"finalYear" binding:"required"`
}

type AdminControlOpenNextYearRequest struct {
	TargetYearNo *int `json:"targetYearNo" binding:"required"`
}

type AdminControlSubmitInitialBaselineRequest struct {
	BaselinePayload *payload.BaselinePayload `json:"baselinePayload" binding:"required"`
}

type AdminControlUnlockYearRequest struct {
	GroupID          *int64  `json:"groupId" binding:"required"`
	YearNo           *int    `json:"yearNo" binding:"required"`
	UnlockTargetType string  `json:"unlockTargetType"`
	TargetStageCode  *string `json:"targetStageCode"`
	Reason           string  `json:"reason"`
}
