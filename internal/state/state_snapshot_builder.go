package state

import (
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
)

// BuildSnapshot 将运行时状态转换为可落库的最小快照。
func BuildSnapshot(current RuntimeState) payload.StateSnapshot {
	return payload.StateSnapshot{
		YearStatus:       current.YearStatus,
		StageStatus:      current.StageStatus,
		ReportStatus:     current.ReportStatus,
		BusinessStatus:   current.BusinessStatus,
		SummaryEffective: current.SummaryEffective,
	}
}

// NewRuntimeStateFromEntities 从组表和组年度状态表构建运行时状态。
func NewRuntimeStateFromEntities(group entity.Group, yearState entity.GroupYearState) RuntimeState {
	return RuntimeState{
		YearNo:                    yearState.YearNo,
		YearType:                  yearState.YearType,
		YearStatus:                yearState.YearStatus,
		StageStatus:               yearState.StageStatus,
		ReportStatus:              yearState.ReportStatus,
		BusinessStatus:            group.BusinessStatus,
		SummaryEffective:          yearState.SummaryEffective,
		LatestStageSubmitVersion:  yearState.LatestStageSubmitVersion,
		LatestReportSubmitVersion: yearState.LatestReportSubmitVersion,
	}
}
