package repository

import (
	"strings"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/state"
)

func nullableStringValue(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func rollbackStagesAtOrAfter(stageCode string) []string {
	switch normalizeRollbackStageCode(stageCode) {
	case state.StageCodeQ1:
		return []string{state.StageCodeQ1, state.StageCodeQ2, state.StageCodeQ3, state.StageCodeQ4}
	case state.StageCodeQ2:
		return []string{state.StageCodeQ2, state.StageCodeQ3, state.StageCodeQ4}
	case state.StageCodeQ3:
		return []string{state.StageCodeQ3, state.StageCodeQ4}
	case state.StageCodeQ4:
		return []string{state.StageCodeQ4}
	default:
		return nil
	}
}

func normalizeRollbackStageCode(stageCode string) string {
	stageCode = strings.ToUpper(strings.TrimSpace(stageCode))
	switch stageCode {
	case enum.StageStatusQ1Open, state.StageCodeQ1:
		return state.StageCodeQ1
	case enum.StageStatusQ2Open, state.StageCodeQ2:
		return state.StageCodeQ2
	case enum.StageStatusQ3Open, state.StageCodeQ3:
		return state.StageCodeQ3
	case enum.StageStatusQ4Open, state.StageCodeQ4:
		return state.StageCodeQ4
	case enum.StageStatusYearEndOpen, state.StageCodeYearEnd:
		return state.StageCodeYearEnd
	case "REPORT", enum.ReportStatusOpen, enum.ReportStatusSubmitted, enum.ReportStatusLocked:
		return state.StageCodeYearEnd
	default:
		return stageCode
	}
}
