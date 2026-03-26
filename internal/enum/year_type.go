package enum

const (
	YearTypeInitial = "INITIAL"
	YearTypeDemo    = "DEMO"
	YearTypeFormal  = "FORMAL"
)

// IsValidYearType 判断年份类型是否合法。
func IsValidYearType(value string) bool {
	switch value {
	case YearTypeInitial, YearTypeDemo, YearTypeFormal:
		return true
	default:
		return false
	}
}
