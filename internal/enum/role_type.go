package enum

const (
	RoleTypeAdmin = "ADMIN"
	RoleTypeGroup = "GROUP"
)

// IsValidRoleType 判断角色类型是否合法。
func IsValidRoleType(value string) bool {
	switch value {
	case RoleTypeAdmin, RoleTypeGroup:
		return true
	default:
		return false
	}
}
