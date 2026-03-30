package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/http/dto"
	"sandbox-game/internal/http/middleware"
	"sandbox-game/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.AuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"登录参数不正确",
			err,
		))
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAuthInvalidCredentials):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnauthorized,
				enum.UnauthorizedCode,
				"用户名或密码错误",
				err,
			))
		case errors.Is(err, service.ErrAuthAccountDisabled):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusForbidden,
				enum.ForbiddenCode,
				"账号已禁用",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"登录失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	identity, ok := middleware.GetAuthIdentity(c)
	if !ok {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusUnauthorized,
			enum.UnauthorizedCode,
			"未获取到当前登录身份",
			nil,
		))
		return
	}

	result := h.authService.BuildCurrentUser(service.AuthenticatedUser{
		UserID:   identity.UserID,
		Username: identity.Username,
		RoleType: identity.RoleType,
		GroupID:  identity.GroupID,
	})
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if _, ok := middleware.GetAuthIdentity(c); !ok {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusUnauthorized,
			enum.UnauthorizedCode,
			"未获取到当前登录身份",
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, dto.Success(h.authService.BuildLogoutResult()))
}
