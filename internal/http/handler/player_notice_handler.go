package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/http/dto"
	"sandbox-game/internal/http/middleware"
	"sandbox-game/internal/service"
)

type PlayerNoticeHandler struct {
	syncService *service.PlayerAdjustmentSyncService
}

func NewPlayerNoticeHandler(syncService *service.PlayerAdjustmentSyncService) *PlayerNoticeHandler {
	return &PlayerNoticeHandler{syncService: syncService}
}

func (h *PlayerNoticeHandler) GetAdjustmentSync(c *gin.Context) {
	var req dto.PlayerAdjustmentSyncRequest
	if err := c.ShouldBindQuery(&req); err != nil || req.YearNo == nil || *req.YearNo < 0 || req.KnownRevision == nil || *req.KnownRevision < 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusBadRequest, enum.BadRequestCode, "yearNo 或 knownRevision 参数不正确", err))
		return
	}
	identity, ok := middleware.GetAuthIdentity(c)
	if !ok || identity.RoleType != enum.RoleTypeGroup || identity.GroupID == nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusForbidden, enum.ForbiddenCode, "当前身份无权同步玩家奖惩", nil))
		return
	}
	result, err := h.syncService.GetSync(c.Request.Context(), *identity.GroupID, *req.YearNo, *req.KnownRevision)
	if err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusInternalServerError, enum.InternalServerErrorCode, "同步奖惩失败", err))
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}
