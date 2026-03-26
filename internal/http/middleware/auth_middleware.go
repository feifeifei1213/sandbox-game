package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"sandbox-game/internal/enum"
)

const authIdentityContextKey = "auth_identity"

type AuthIdentity struct {
	UserID   int64
	RoleType string
	GroupID  *int64
	Username string
}

// AuthBypass 是 M1-02 阶段的认证占位实现。
// 后续接入真实登录态后，再替换为正式认证中间件。
func AuthBypass() gin.HandlerFunc {
	return func(c *gin.Context) {
		identity, err := buildBypassIdentity(c)
		if err != nil {
			AbortWithAppError(c, NewAppError(
				http.StatusBadRequest,
				enum.BadRequestCode,
				err.Error(),
				err,
			))
			return
		}
		c.Set(authIdentityContextKey, identity)
		c.Next()
	}
}

func GetAuthIdentity(c *gin.Context) (AuthIdentity, bool) {
	value, ok := c.Get(authIdentityContextKey)
	if !ok {
		return AuthIdentity{}, false
	}
	identity, ok := value.(AuthIdentity)
	return identity, ok
}

func buildBypassIdentity(c *gin.Context) (AuthIdentity, error) {
	roleType := strings.ToUpper(strings.TrimSpace(c.GetHeader("X-Role-Type")))
	if roleType == "" {
		roleType = enum.RoleTypeGroup
	}
	if !enum.IsValidRoleType(roleType) {
		return AuthIdentity{}, strconv.ErrSyntax
	}

	userID := int64(101)
	if value := strings.TrimSpace(c.GetHeader("X-User-Id")); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return AuthIdentity{}, err
		}
		userID = parsed
	}

	username := strings.TrimSpace(c.GetHeader("X-Username"))
	if username == "" {
		if roleType == enum.RoleTypeAdmin {
			username = "admin"
		} else {
			username = "group01"
		}
	}

	if roleType == enum.RoleTypeAdmin {
		return AuthIdentity{
			UserID:   userID,
			RoleType: roleType,
			GroupID:  nil,
			Username: username,
		}, nil
	}

	groupID := int64(1)
	if value := strings.TrimSpace(c.GetHeader("X-Group-Id")); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return AuthIdentity{}, err
		}
		groupID = parsed
	}

	return AuthIdentity{
		UserID:   userID,
		RoleType: roleType,
		GroupID:  &groupID,
		Username: username,
	}, nil
}
