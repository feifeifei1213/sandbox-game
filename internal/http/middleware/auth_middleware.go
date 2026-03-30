package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/service"
)

const authIdentityContextKey = "auth_identity"

type AuthIdentity struct {
	UserID   int64
	RoleType string
	GroupID  *int64
	Username string
}

// AuthBypass 是开发联调阶段保留的认证占位实现。
// 当配置显式切换到 bypass 模式时，可继续复用旧的请求头透传方式。
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

func RequireAuth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			abortUnauthorized(c, "未登录或登录态已失效", err)
			return
		}

		authenticatedUser, err := authService.AuthenticateAccessToken(c.Request.Context(), accessToken)
		if err != nil {
			abortUnauthorized(c, "未登录或登录态已失效", err)
			return
		}

		c.Set(authIdentityContextKey, AuthIdentity{
			UserID:   authenticatedUser.UserID,
			RoleType: authenticatedUser.RoleType,
			GroupID:  authenticatedUser.GroupID,
			Username: authenticatedUser.Username,
		})
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

func abortUnauthorized(c *gin.Context, message string, err error) {
	AbortWithAppError(c, NewAppError(
		http.StatusUnauthorized,
		enum.UnauthorizedCode,
		message,
		err,
	))
}

func extractBearerToken(authHeader string) (string, error) {
	fields := strings.Fields(strings.TrimSpace(authHeader))
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") || strings.TrimSpace(fields[1]) == "" {
		return "", errors.New("invalid authorization header")
	}
	return fields[1], nil
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
