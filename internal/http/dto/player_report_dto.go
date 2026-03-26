package dto

import (
	"time"

	"sandbox-game/internal/model/payload"
)

type PlayerReportGetViewRequest struct {
	YearNo *int `form:"yearNo" binding:"required"`
}

type PlayerReportSaveDraftRequest struct {
	YearNo              *int                        `json:"yearNo" binding:"required"`
	ReportManualPayload payload.ReportManualPayload `json:"reportManualPayload"`
	ClientSaveTime      *time.Time                  `json:"clientSaveTime"`
}

type PlayerReportSaveDraftResponse struct {
	GroupID          int64     `json:"groupId"`
	YearNo           int       `json:"yearNo"`
	YearStatus       string    `json:"yearStatus"`
	ReportStatus     string    `json:"reportStatus"`
	LastDraftSavedAt time.Time `json:"lastDraftSavedAt"`
}

type PlayerReportSubmitRequest struct {
	YearNo              *int                        `json:"yearNo" binding:"required"`
	ReportManualPayload payload.ReportManualPayload `json:"reportManualPayload"`
}

type PlayerReportSubmitResponse struct {
	GroupID                   int64     `json:"groupId"`
	YearNo                    int       `json:"yearNo"`
	YearStatus                string    `json:"yearStatus"`
	ReportStatus              string    `json:"reportStatus"`
	BusinessStatus            string    `json:"businessStatus"`
	SummaryEffective          bool      `json:"summaryEffective"`
	LatestReportSubmitVersion int       `json:"latestReportSubmitVersion"`
	BalanceCheckPassed        bool      `json:"balanceCheckPassed"`
	SubmittedAt               time.Time `json:"submittedAt"`
}
