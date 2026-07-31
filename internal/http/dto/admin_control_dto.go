package dto

import "sandbox-game/internal/model/payload"

type AdminControlInitializeGameRequest struct {
	GroupCount         *int                       `json:"groupCount" binding:"required"`
	EditionCode        string                     `json:"editionCode"`
	DictionarySchemeID *int64                     `json:"dictionarySchemeId"`
	DictionaryItems    []AdminDictionaryItemInput `json:"dictionaryItems"`
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
	YearNo           *int    `json:"yearNo"`
	UnlockTargetType string  `json:"unlockTargetType"`
	TargetStageCode  *string `json:"targetStageCode"`
	Reason           string  `json:"reason"`
}
