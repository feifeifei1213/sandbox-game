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

type AdminRollbackHandler struct {
	queryService   *service.AdminRollbackQueryService
	commandService *service.AdminRollbackCommandService
}

func NewAdminRollbackHandler(
	queryService *service.AdminRollbackQueryService,
	commandService *service.AdminRollbackCommandService,
) *AdminRollbackHandler {
	return &AdminRollbackHandler{
		queryService:   queryService,
		commandService: commandService,
	}
}

func (h *AdminRollbackHandler) ListSnapshots(c *gin.Context) {
	var req dto.AdminRollbackListSnapshotsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusBadRequest, enum.BadRequestCode, "快照查询参数不正确", err))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权查看状态快照") {
		return
	}

	cmd := service.ListSnapshotsCommand{
		GroupID: req.GroupID,
		YearNo:  req.YearNo,
	}
	if req.SnapshotScope != nil {
		cmd.SnapshotScope = *req.SnapshotScope
	}
	if req.SnapshotType != nil {
		cmd.SnapshotType = *req.SnapshotType
	}
	if req.StageCode != nil {
		cmd.StageCode = *req.StageCode
	}
	if req.PageNo != nil {
		cmd.PageNo = *req.PageNo
	}
	if req.PageSize != nil {
		cmd.PageSize = *req.PageSize
	}
	result, err := h.queryService.ListSnapshots(c.Request.Context(), cmd)
	if err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusInternalServerError, enum.InternalServerErrorCode, "查询状态快照失败", err))
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminRollbackHandler) GetSnapshotDetail(c *gin.Context) {
	var req dto.AdminRollbackGetSnapshotDetailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusBadRequest, enum.BadRequestCode, "快照详情参数不正确", err))
		return
	}
	if req.SnapshotID == nil || *req.SnapshotID <= 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusBadRequest, enum.BadRequestCode, "snapshotId 参数不正确", nil))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权查看状态快照详情") {
		return
	}

	result, err := h.queryService.GetSnapshotDetail(c.Request.Context(), *req.SnapshotID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusNotFound, enum.NotFoundCode, "未找到状态快照", err))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusInternalServerError, enum.InternalServerErrorCode, "查询状态快照详情失败", err))
		}
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminRollbackHandler) CreateSnapshot(c *gin.Context) {
	var req dto.AdminRollbackCreateSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusBadRequest, enum.BadRequestCode, "创建快照参数不正确", err))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权创建状态快照") {
		return
	}

	yearNo := 0
	if req.YearNo != nil {
		yearNo = *req.YearNo
	}
	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.CreateManualSnapshot(c.Request.Context(), service.CreateManualSnapshotCommand{
		SnapshotScope: req.SnapshotScope,
		GroupID:       req.GroupID,
		YearNo:        yearNo,
		StageCode:     req.StageCode,
		Description:   req.Description,
		OperatorID:    identity.UserID,
		OperatorName:  identity.Username,
	})
	if err != nil {
		writeAdminRollbackError(c, err, "创建状态快照失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminRollbackHandler) RestoreGroupSnapshot(c *gin.Context) {
	var req dto.AdminRollbackRestoreGroupSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusBadRequest, enum.BadRequestCode, "恢复快照参数不正确", err))
		return
	}
	if req.SnapshotID == nil || *req.SnapshotID <= 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusBadRequest, enum.BadRequestCode, "snapshotId 参数不正确", nil))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权恢复状态快照") {
		return
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.RestoreGroupSnapshot(c.Request.Context(), service.RestoreGroupSnapshotCommand{
		SnapshotID:   *req.SnapshotID,
		Reason:       req.Reason,
		ConfirmText:  req.ConfirmText,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		writeAdminRollbackError(c, err, "恢复状态快照失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func writeAdminRollbackError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound), errors.Is(err, service.ErrRollbackSnapshotNotFound):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusNotFound, enum.NotFoundCode, "未找到目标快照或小组年份数据", err))
	case errors.Is(err, service.ErrRollbackReasonRequired):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, "原因或说明不能为空", err))
	case errors.Is(err, service.ErrRollbackConfirmRequired):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, "请填写确认文本：确认恢复", err))
	case errors.Is(err, service.ErrRollbackSnapshotScopeUnsupported):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, "首版只支持单组快照恢复", err))
	case errors.Is(err, service.ErrRollbackTargetInvalid):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, "回退目标不合法", err))
	case errors.Is(err, service.ErrRollbackStateChanged), errors.Is(err, service.ErrRollbackGroupWriteLocked):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusConflict, enum.ConflictCode, "状态已变化，请刷新后重试", err))
	default:
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusInternalServerError, enum.InternalServerErrorCode, fallback, err))
	}
}
