package enum

const (
	BusinessStatusNormal   = "NORMAL"
	BusinessStatusBankrupt = "BANKRUPT"
)

// IsValidBusinessStatus 判断经营状态是否合法。
func IsValidBusinessStatus(value string) bool {
	switch value {
	case BusinessStatusNormal, BusinessStatusBankrupt:
		return true
	default:
		return false
	}
}
