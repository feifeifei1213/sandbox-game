package enum

const (
	ReportStatusLocked    = "REPORT_LOCKED"
	ReportStatusOpen      = "REPORT_OPEN"
	ReportStatusSubmitted = "REPORT_SUBMITTED"
)

// IsValidReportStatus 判断财报状态是否合法。
func IsValidReportStatus(value string) bool {
	switch value {
	case ReportStatusLocked, ReportStatusOpen, ReportStatusSubmitted:
		return true
	default:
		return false
	}
}
