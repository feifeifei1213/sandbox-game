package enum

const (
	YearStatusLocked        = "LOCKED"
	YearStatusOperating     = "OPERATING"
	YearStatusReportPending = "REPORT_PENDING"
	YearStatusReporting     = "REPORTING"
	YearStatusCompleted     = "COMPLETED"
)

// IsValidYearStatus 判断年份主状态是否合法。
func IsValidYearStatus(value string) bool {
	switch value {
	case YearStatusLocked, YearStatusOperating, YearStatusReportPending, YearStatusReporting, YearStatusCompleted:
		return true
	default:
		return false
	}
}
