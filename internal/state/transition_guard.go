package state

import "sandbox-game/internal/enum"

const (
	OperatingScopeYearStart = "YEAR_START"
	OperatingScopeQ1        = "Q1"
	OperatingScopeQ2        = "Q2"
	OperatingScopeQ3        = "Q3"
	OperatingScopeQ4        = "Q4"
	OperatingScopeYearEnd   = "YEAR_END"
)

var orderedOperatingScopes = []string{
	OperatingScopeYearStart,
	OperatingScopeQ1,
	OperatingScopeQ2,
	OperatingScopeQ3,
	OperatingScopeQ4,
	OperatingScopeYearEnd,
}

// OperatingPermission 表示经营页的可编辑/可提交结果。
type OperatingPermission struct {
	CanView          bool     `json:"canView"`
	CanEdit          bool     `json:"canEdit"`
	CanSubmit        bool     `json:"canSubmit"`
	CurrentStageCode string   `json:"currentStageCode"`
	EditableScopes   []string `json:"editableScopes"`
	ReadonlyScopes   []string `json:"readonlyScopes"`
}

// ReportPermission 表示财报页的可查看/可编辑/可提交结果。
type ReportPermission struct {
	CanView   bool `json:"canView"`
	CanEdit   bool `json:"canEdit"`
	CanSubmit bool `json:"canSubmit"`
}

// TransitionGuard 负责状态允许性判断，不直接修改状态。
type TransitionGuard struct{}

// NewTransitionGuard 创建状态判断器。
func NewTransitionGuard() *TransitionGuard {
	return &TransitionGuard{}
}

// CanOpenYear 判断该年份是否允许被管理员开放。
func (g *TransitionGuard) CanOpenYear(current RuntimeState) bool {
	return current.YearStatus == enum.YearStatusLocked && current.BusinessStatus != enum.BusinessStatusBankrupt
}

// CanSubmitStage 判断当前阶段是否允许正式提交。
func (g *TransitionGuard) CanSubmitStage(current RuntimeState, stageCode string) bool {
	if current.BusinessStatus == enum.BusinessStatusBankrupt || current.YearStatus != enum.YearStatusOperating {
		return false
	}

	expectedStatus, ok := stageStatusByCode[stageCode]
	if !ok {
		return false
	}

	return current.StageStatus == expectedStatus
}

// CanSubmitReport 判断当前财报是否允许正式提交。
func (g *TransitionGuard) CanSubmitReport(current RuntimeState) bool {
	if current.BusinessStatus == enum.BusinessStatusBankrupt {
		return false
	}
	if current.ReportStatus != enum.ReportStatusOpen {
		return false
	}

	return current.YearStatus == enum.YearStatusReportPending || current.YearStatus == enum.YearStatusReporting
}

// CanUnlockYear 判断管理员是否可对该年份执行异常解锁。
func (g *TransitionGuard) CanUnlockYear(current RuntimeState, nextYearAlreadyOpened bool) bool {
	if nextYearAlreadyOpened {
		return false
	}

	switch current.YearStatus {
	case enum.YearStatusOperating, enum.YearStatusReportPending, enum.YearStatusReporting, enum.YearStatusCompleted:
		return true
	default:
		return false
	}
}

// BuildOperatingPermission 输出经营页的只读/可编辑范围。
func (g *TransitionGuard) BuildOperatingPermission(current RuntimeState) OperatingPermission {
	permission := OperatingPermission{
		CanView:          true,
		CanEdit:          false,
		CanSubmit:        false,
		CurrentStageCode: CurrentStageCode(current.StageStatus),
		EditableScopes:   []string{},
		ReadonlyScopes:   append([]string(nil), orderedOperatingScopes...),
	}

	if current.YearStatus != enum.YearStatusOperating || current.BusinessStatus == enum.BusinessStatusBankrupt {
		return permission
	}

	editableScopes := editableScopesForStage(current.StageStatus)
	permission.CanEdit = len(editableScopes) > 0
	permission.CanSubmit = len(editableScopes) > 0
	permission.EditableScopes = editableScopes
	permission.ReadonlyScopes = readonlyScopes(editableScopes)

	return permission
}

// BuildReportPermission 输出财报页的可查看/可编辑能力。
func (g *TransitionGuard) BuildReportPermission(current RuntimeState) ReportPermission {
	canView := current.ReportStatus == enum.ReportStatusOpen ||
		current.ReportStatus == enum.ReportStatusSubmitted ||
		current.YearStatus == enum.YearStatusReportPending ||
		current.YearStatus == enum.YearStatusReporting ||
		current.YearStatus == enum.YearStatusCompleted

	canEdit := g.CanSubmitReport(current)

	return ReportPermission{
		CanView:   canView,
		CanEdit:   canEdit,
		CanSubmit: canEdit,
	}
}

func editableScopesForStage(stageStatus string) []string {
	switch stageStatus {
	case enum.StageStatusQ1Open:
		return []string{OperatingScopeYearStart, OperatingScopeQ1}
	case enum.StageStatusQ2Open:
		return []string{OperatingScopeQ2}
	case enum.StageStatusQ3Open:
		return []string{OperatingScopeQ3}
	case enum.StageStatusQ4Open:
		return []string{OperatingScopeQ4}
	case enum.StageStatusYearEndOpen:
		return []string{OperatingScopeYearEnd}
	default:
		return []string{}
	}
}

func readonlyScopes(editableScopes []string) []string {
	if len(editableScopes) == 0 {
		return append([]string(nil), orderedOperatingScopes...)
	}

	editableMap := make(map[string]struct{}, len(editableScopes))
	for _, scope := range editableScopes {
		editableMap[scope] = struct{}{}
	}

	result := make([]string, 0, len(orderedOperatingScopes)-len(editableScopes))
	for _, scope := range orderedOperatingScopes {
		if _, ok := editableMap[scope]; ok {
			continue
		}
		result = append(result, scope)
	}

	return result
}
