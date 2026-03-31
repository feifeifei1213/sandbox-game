package assembler

import (
	"time"

	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	carryforwardrules "sandbox-game/internal/rules/carryforward"
	calcctx "sandbox-game/internal/rules/context"
	operatingrules "sandbox-game/internal/rules/operating"
	"sandbox-game/internal/state"
)

type PlayerOperatingView struct {
	GroupID                int64                                 `json:"groupId"`
	YearNo                 int                                   `json:"yearNo"`
	YearStatus             string                                `json:"yearStatus"`
	StageStatus            string                                `json:"stageStatus"`
	ReportStatus           string                                `json:"reportStatus"`
	BusinessStatus         string                                `json:"businessStatus"`
	CurrentStageCode       string                                `json:"currentStageCode"`
	CanView                bool                                  `json:"canView"`
	CanEdit                bool                                  `json:"canEdit"`
	CanSubmit              bool                                  `json:"canSubmit"`
	HasInvalidDraft        bool                                  `json:"hasInvalidDraft"`
	InvalidScopes          []string                              `json:"invalidScopes"`
	HasRetainedReportDraft bool                                  `json:"hasRetainedReportDraft"`
	OperatingPayload       payload.OperatingPayload              `json:"operatingPayload"`
	EditableScopes         []string                              `json:"editableScopes"`
	ReadonlyScopes         []string                              `json:"readonlyScopes"`
	StageSubmitHistory     []PlayerOperatingStageSubmitHistory   `json:"stageSubmitHistory"`
	LastDraftSavedAt       *time.Time                            `json:"lastDraftSavedAt"`
	QuarterCashChecks      map[string]float64                    `json:"quarterCashChecks"`
	DerivedValues          map[string]float64                    `json:"derivedValues"`
	PeriodEndCash          float64                               `json:"periodEndCash"`
	CarryForward           *carryforwardrules.CarryForwardResult `json:"carryForward,omitempty"`
}

type PlayerOperatingStageSubmitHistory struct {
	StageCode     string    `json:"stageCode"`
	SubmitVersion int       `json:"submitVersion"`
	PeriodEndCash float64   `json:"periodEndCash"`
	SubmitTime    time.Time `json:"submitTime"`
}

type PlayerOperatingAssembler struct{}

func NewPlayerOperatingAssembler() *PlayerOperatingAssembler {
	return &PlayerOperatingAssembler{}
}

func (a *PlayerOperatingAssembler) Build(
	ctx calcctx.CalculationContext,
	draft *entity.GroupOperatingDraft,
	stageSubmissions []entity.GroupStageSubmission,
	result operatingrules.CalculationResult,
	permission state.OperatingPermission,
	carryForward *carryforwardrules.CarryForwardResult,
) *PlayerOperatingView {
	history := make([]PlayerOperatingStageSubmitHistory, 0, len(stageSubmissions))
	for _, item := range stageSubmissions {
		history = append(history, PlayerOperatingStageSubmitHistory{
			StageCode:     item.StageCode,
			SubmitVersion: item.SubmitVersion,
			PeriodEndCash: item.PeriodEndCash,
			SubmitTime:    item.SubmitTime,
		})
	}

	var lastDraftSavedAt *time.Time
	if draft != nil {
		lastDraftSavedAt = draft.LastAutoSavedAt
	}

	draftState := buildOperatingDraftState(ctx.State, stageSubmissions)

	return &PlayerOperatingView{
		GroupID:                ctx.Group.ID,
		YearNo:                 ctx.YearState.YearNo,
		YearStatus:             ctx.State.YearStatus,
		StageStatus:            ctx.State.StageStatus,
		ReportStatus:           ctx.State.ReportStatus,
		BusinessStatus:         ctx.State.BusinessStatus,
		CurrentStageCode:       permission.CurrentStageCode,
		CanView:                permission.CanView,
		CanEdit:                permission.CanEdit,
		CanSubmit:              permission.CanSubmit,
		HasInvalidDraft:        draftState.HasInvalidDraft,
		InvalidScopes:          draftState.InvalidScopes,
		HasRetainedReportDraft: draftState.HasRetainedReportDraft,
		OperatingPayload:       result.OperatingPayload,
		EditableScopes:         permission.EditableScopes,
		ReadonlyScopes:         permission.ReadonlyScopes,
		StageSubmitHistory:     history,
		LastDraftSavedAt:       lastDraftSavedAt,
		QuarterCashChecks:      result.QuarterCashChecks,
		DerivedValues:          result.DerivedValues,
		PeriodEndCash:          result.PeriodEndCash,
		CarryForward:           carryForward,
	}
}
