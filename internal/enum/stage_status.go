package enum

const (
	StageStatusQ1Open      = "Q1_OPEN"
	StageStatusQ2Open      = "Q2_OPEN"
	StageStatusQ3Open      = "Q3_OPEN"
	StageStatusQ4Open      = "Q4_OPEN"
	StageStatusYearEndOpen = "YEAR_END_OPEN"
)

// IsValidStageStatus 判断经营阶段状态是否合法。
func IsValidStageStatus(value string) bool {
	switch value {
	case StageStatusQ1Open, StageStatusQ2Open, StageStatusQ3Open, StageStatusQ4Open, StageStatusYearEndOpen:
		return true
	default:
		return false
	}
}
