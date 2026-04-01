package assembler

import (
	"time"

	"sandbox-game/internal/model/payload"
	calcctx "sandbox-game/internal/rules/context"
	"sandbox-game/internal/state"
)

type PlayerReportManualFieldOptions struct {
	IncomeTaxRateOptions []float64 `json:"incomeTaxRateOptions"`
}

type PlayerReportView struct {
	GroupID               int64                          `json:"groupId"`
	YearNo                int                            `json:"yearNo"`
	YearStatus            string                         `json:"yearStatus"`
	ReportStatus          string                         `json:"reportStatus"`
	BusinessStatus        string                         `json:"businessStatus"`
	CanView               bool                           `json:"canView"`
	CanEdit               bool                           `json:"canEdit"`
	CanSubmit             bool                           `json:"canSubmit"`
	HasInvalidDraft       bool                           `json:"hasInvalidDraft"`
	ReportComputedPayload payload.ReportComputedPayload  `json:"reportComputedPayload"`
	ReportManualPayload   payload.ReportManualPayload    `json:"reportManualPayload"`
	ManualFieldOptions    PlayerReportManualFieldOptions `json:"manualFieldOptions"`
	LastDraftSavedAt      *time.Time                     `json:"lastDraftSavedAt"`
	NoticeBoard           *PlayerNoticeBoard             `json:"noticeBoard"`
}

type PlayerReportAssembler struct{}

func NewPlayerReportAssembler() *PlayerReportAssembler {
	return &PlayerReportAssembler{}
}

func (a *PlayerReportAssembler) Build(
	ctx calcctx.CalculationContext,
	computed payload.ReportComputedPayload,
	manual payload.ReportManualPayload,
	lastDraftSavedAt *time.Time,
	permission state.ReportPermission,
	noticeBoard *PlayerNoticeBoard,
) *PlayerReportView {
	return &PlayerReportView{
		GroupID:               ctx.Group.ID,
		YearNo:                ctx.YearState.YearNo,
		YearStatus:            ctx.State.YearStatus,
		ReportStatus:          ctx.State.ReportStatus,
		BusinessStatus:        ctx.State.BusinessStatus,
		CanView:               permission.CanView,
		CanEdit:               permission.CanEdit,
		CanSubmit:             permission.CanSubmit,
		HasInvalidDraft:       hasRetainedEditableReportDraft(ctx.State),
		ReportComputedPayload: computed,
		ReportManualPayload:   manual,
		ManualFieldOptions: PlayerReportManualFieldOptions{
			IncomeTaxRateOptions: payload.AllowedIncomeTaxRates(),
		},
		LastDraftSavedAt: lastDraftSavedAt,
		NoticeBoard:      noticeBoard,
	}
}
