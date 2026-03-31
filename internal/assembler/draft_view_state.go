package assembler

import (
	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/state"
)

var orderedOperatingScopeCodes = []string{
	state.StageCodeQ1,
	state.StageCodeQ2,
	state.StageCodeQ3,
	state.StageCodeQ4,
	state.StageCodeYearEnd,
}

type OperatingDraftState struct {
	HasInvalidDraft        bool
	InvalidScopes          []string
	HasRetainedReportDraft bool
}

func buildOperatingDraftState(current state.RuntimeState, stageSubmissions []entity.GroupStageSubmission) OperatingDraftState {
	result := OperatingDraftState{
		InvalidScopes:          []string{},
		HasRetainedReportDraft: hasRetainedReportDraft(current),
	}

	if current.YearStatus != enum.YearStatusOperating {
		result.HasInvalidDraft = result.HasRetainedReportDraft
		return result
	}

	currentRank, ok := operatingStageRank(state.CurrentStageCode(current.StageStatus))
	if !ok {
		result.HasInvalidDraft = result.HasRetainedReportDraft
		return result
	}

	maxRank := maxSubmittedOperatingStageRankFromSubmissions(stageSubmissions)
	if maxRank <= currentRank {
		result.HasInvalidDraft = result.HasRetainedReportDraft
		return result
	}

	invalidScopes := make([]string, 0, len(orderedOperatingScopeCodes))
	for _, scopeCode := range orderedOperatingScopeCodes {
		rank, ok := operatingStageRank(scopeCode)
		if !ok {
			continue
		}
		if rank > currentRank && rank <= maxRank {
			invalidScopes = append(invalidScopes, scopeCode)
		}
	}

	result.InvalidScopes = invalidScopes
	result.HasInvalidDraft = len(invalidScopes) > 0 || result.HasRetainedReportDraft
	return result
}

func hasRetainedReportDraft(current state.RuntimeState) bool {
	return current.LatestReportSubmitVersion > 0 && current.ReportStatus != enum.ReportStatusSubmitted
}

func hasRetainedEditableReportDraft(current state.RuntimeState) bool {
	return current.LatestReportSubmitVersion > 0 && current.ReportStatus == enum.ReportStatusOpen
}

func maxSubmittedOperatingStageRankFromSubmissions(stageSubmissions []entity.GroupStageSubmission) int {
	maxRank := 0
	for _, item := range stageSubmissions {
		rank, ok := operatingStageRank(item.StageCode)
		if !ok {
			continue
		}
		if rank > maxRank {
			maxRank = rank
		}
	}
	return maxRank
}

func operatingStageRank(stageCode string) (int, bool) {
	switch stageCode {
	case state.StageCodeQ1:
		return 1, true
	case state.StageCodeQ2:
		return 2, true
	case state.StageCodeQ3:
		return 3, true
	case state.StageCodeQ4:
		return 4, true
	case state.StageCodeYearEnd:
		return 5, true
	default:
		return 0, false
	}
}
