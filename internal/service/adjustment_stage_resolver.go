package service

import (
	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/state"
)

// resolveAdjustmentStage 只依据服务端当前年度状态决定奖惩归属，管理员请求不得覆盖该结果。
func resolveAdjustmentStage(yearState entity.GroupYearState) (string, error) {
	if yearState.YearStatus == enum.YearStatusLocked {
		return "", ErrAdminAdjustmentYearNotOpen
	}
	if yearState.YearStatus == enum.YearStatusCompleted || yearState.ReportStatus == enum.ReportStatusSubmitted {
		return "", ErrAdminAdjustmentStageLocked
	}

	switch yearState.YearStatus {
	case enum.YearStatusOperating:
		stageCode := state.CurrentStageCode(yearState.StageStatus)
		if !state.IsValidStageCode(stageCode) {
			return "", ErrAdminAdjustmentStageInvalid
		}
		return stageCode, nil
	case enum.YearStatusReportPending, enum.YearStatusReporting:
		if yearState.ReportStatus != enum.ReportStatusOpen {
			return "", ErrAdminAdjustmentStageLocked
		}
		return state.StageCodeYearEnd, nil
	default:
		return "", ErrAdminAdjustmentStageInvalid
	}
}
