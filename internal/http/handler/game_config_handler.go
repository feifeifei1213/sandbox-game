package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/http/dto"
	"sandbox-game/internal/http/middleware"
	"sandbox-game/internal/service"
)

type GameConfigHandler struct {
	queryService *service.GameConfigQueryService
}

func NewGameConfigHandler(queryService *service.GameConfigQueryService) *GameConfigHandler {
	return &GameConfigHandler{queryService: queryService}
}

func (h *GameConfigHandler) GetCurrent(c *gin.Context) {
	result, err := h.queryService.GetCurrent(c.Request.Context())
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到游戏配置",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"获取游戏配置失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *GameConfigHandler) ListGameEditions(c *gin.Context) {
	if !ensureAdminIdentity(c, "当前身份无权查看沙盘版本包") {
		return
	}

	result, err := h.queryService.ListGameEditions(c.Request.Context())
	if err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusInternalServerError,
			enum.InternalServerErrorCode,
			"获取沙盘版本包失败",
			err,
		))
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *GameConfigHandler) GetYearTabs(c *gin.Context) {
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
	if identity.RoleType == enum.RoleTypeGroup && identity.GroupID == nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusForbidden,
			enum.ForbiddenCode,
			"当前身份缺少小组信息，无法获取年份标签",
			nil,
		))
		return
	}

	result, err := h.queryService.GetYearTabs(c.Request.Context(), identity.RoleType, identity.GroupID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到年份标签配置",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"获取年份标签状态失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}
