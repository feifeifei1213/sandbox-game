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

type AdminNoticeHandler struct {
	queryService   *service.AdminNoticeQueryService
	commandService *service.AdminNoticeCommandService
}

func NewAdminNoticeHandler(
	queryService *service.AdminNoticeQueryService,
	commandService *service.AdminNoticeCommandService,
) *AdminNoticeHandler {
	return &AdminNoticeHandler{
		queryService:   queryService,
		commandService: commandService,
	}
}

func (h *AdminNoticeHandler) GetRecords(c *gin.Context) {
	var req dto.AdminNoticeGetRecordsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"limit 参数不正确",
			err,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权查看通知与奖惩记录") {
		return
	}

	limit := 20
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	result, err := h.queryService.GetRecords(c.Request.Context(), limit)
	if err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusInternalServerError,
			enum.InternalServerErrorCode,
			"获取通知与奖惩记录失败",
			err,
		))
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminNoticeHandler) SendGeneral(c *gin.Context) {
	var req dto.AdminNoticeSendGeneralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"普通通知参数不正确",
			err,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权发送普通通知") {
		return
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.SendGeneralNotice(c.Request.Context(), service.SendGeneralNoticeCommand{
		TargetScope:   req.TargetScope,
		TargetGroupID: req.TargetGroupID,
		Content:       req.Content,
		Pinned:        req.Pinned,
		OperatorID:    identity.UserID,
		OperatorName:  identity.Username,
	})
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"目标小组不存在",
				err,
			))
		case errors.Is(err, service.ErrAdminNoticeContentRequired),
			errors.Is(err, service.ErrAdminNoticeTargetScopeInvalid),
			errors.Is(err, service.ErrAdminNoticeTargetGroupRequired):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				resolveAdminNoticeErrorMessage(err),
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"发送普通通知失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminNoticeHandler) SendAdjustment(c *gin.Context) {
	var req dto.AdminNoticeSendAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"奖惩参数不正确",
			err,
		))
		return
	}
	if req.GroupID == nil || *req.GroupID <= 0 || req.YearNo == nil || *req.YearNo < 0 || req.Amount == nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"groupId、yearNo 或 amount 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权下发奖惩") {
		return
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.SendAdjustment(c.Request.Context(), service.SendAdjustmentCommand{
		GroupID:        *req.GroupID,
		YearNo:         *req.YearNo,
		AdjustmentType: req.AdjustmentType,
		Amount:         *req.Amount,
		Reason:         req.Reason,
		OperatorID:     identity.UserID,
		OperatorName:   identity.Username,
	})
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"目标小组或年份不存在",
				err,
			))
		case errors.Is(err, service.ErrAdminAdjustmentStageInvalid),
			errors.Is(err, service.ErrAdminAdjustmentTypeInvalid),
			errors.Is(err, service.ErrAdminAdjustmentAmountInvalid),
			errors.Is(err, service.ErrManualNumberNotInteger),
			errors.Is(err, service.ErrAdminAdjustmentReasonRequired),
			errors.Is(err, service.ErrAdminAdjustmentYearNotOpen),
			errors.Is(err, service.ErrAdminAdjustmentStageLocked),
			errors.Is(err, service.ErrAdminAdjustmentGroupNotAvailable):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				resolveAdminNoticeErrorMessage(err),
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"下发奖惩失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminNoticeHandler) PreviewAdjustment(c *gin.Context) {
	var req dto.AdminNoticePreviewAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusBadRequest, enum.BadRequestCode, "奖惩预览参数不正确", err))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权预览奖惩影响") {
		return
	}
	cmd := service.PreviewAdjustmentCommand{Operation: req.Operation, AdjustmentType: req.AdjustmentType, Reason: req.Reason}
	if req.AdjustmentID != nil {
		cmd.AdjustmentID = *req.AdjustmentID
	}
	if req.GroupID != nil {
		cmd.GroupID = *req.GroupID
	}
	if req.YearNo != nil {
		cmd.YearNo = *req.YearNo
	}
	if req.Amount != nil {
		cmd.Amount = *req.Amount
	}
	result, err := h.commandService.PreviewAdjustment(c.Request.Context(), cmd)
	if err != nil {
		h.abortAdjustmentError(c, err, "预览奖惩影响失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminNoticeHandler) VoidAdjustment(c *gin.Context) {
	var req dto.AdminNoticeVoidAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.AdjustmentID == nil || *req.AdjustmentID <= 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusBadRequest, enum.BadRequestCode, "奖惩作废参数不正确", err))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权作废奖惩") {
		return
	}
	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.VoidAdjustment(c.Request.Context(), service.VoidAdjustmentCommand{
		AdjustmentID: *req.AdjustmentID, Reason: req.Reason,
		OperatorID: identity.UserID, OperatorName: identity.Username,
	})
	if err != nil {
		h.abortAdjustmentError(c, err, "作废奖惩失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminNoticeHandler) abortAdjustmentError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusNotFound, enum.NotFoundCode, "目标奖惩、小组或年份不存在", err))
	case errors.Is(err, service.ErrAdminAdjustmentNotEffective):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusConflict, enum.ConflictCode, "该奖惩已失效，不能重复操作", err))
	case isAdminAdjustmentBusinessError(err):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, resolveAdminNoticeErrorMessage(err), err))
	default:
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusInternalServerError, enum.InternalServerErrorCode, fallback, err))
	}
}

func isAdminAdjustmentBusinessError(err error) bool {
	return errors.Is(err, service.ErrAdminAdjustmentOperationInvalid) ||
		errors.Is(err, service.ErrAdminAdjustmentStageInvalid) ||
		errors.Is(err, service.ErrAdminAdjustmentTypeInvalid) ||
		errors.Is(err, service.ErrAdminAdjustmentAmountInvalid) ||
		errors.Is(err, service.ErrManualNumberNotInteger) ||
		errors.Is(err, service.ErrAdminAdjustmentReasonRequired) ||
		errors.Is(err, service.ErrAdminAdjustmentVoidReasonRequired) ||
		errors.Is(err, service.ErrAdminAdjustmentYearNotOpen) ||
		errors.Is(err, service.ErrAdminAdjustmentStageLocked) ||
		errors.Is(err, service.ErrAdminAdjustmentGroupNotAvailable)
}

func resolveAdminNoticeErrorMessage(err error) string {
	switch {
	case errors.Is(err, service.ErrAdminAdjustmentOperationInvalid):
		return "奖惩预览操作必须为 CREATE 或 VOID"
	case errors.Is(err, service.ErrAdminAdjustmentVoidReasonRequired):
		return "作废原因不能为空"
	case errors.Is(err, service.ErrAdminAdjustmentNotEffective):
		return "该奖惩当前不是有效状态"
	case errors.Is(err, service.ErrAdminNoticeContentRequired):
		return "普通通知内容不能为空"
	case errors.Is(err, service.ErrAdminNoticeTargetScopeInvalid):
		return "普通通知目标范围不合法"
	case errors.Is(err, service.ErrAdminNoticeTargetGroupRequired):
		return "单组通知必须选择目标小组"
	case errors.Is(err, service.ErrAdminAdjustmentStageInvalid):
		return "当前年度阶段无法自动归属奖罚"
	case errors.Is(err, service.ErrAdminAdjustmentTypeInvalid):
		return "奖惩类型不合法"
	case errors.Is(err, service.ErrAdminAdjustmentAmountInvalid):
		return "奖惩金额必须大于 0"
	case errors.Is(err, service.ErrManualNumberNotInteger):
		return "奖惩金额必须为整数"
	case errors.Is(err, service.ErrAdminAdjustmentReasonRequired):
		return "奖惩原因不能为空"
	case errors.Is(err, service.ErrAdminAdjustmentYearNotOpen):
		return "目标年份当前未开放，不能直接下发奖惩"
	case errors.Is(err, service.ErrAdminAdjustmentStageLocked):
		return "目标季度已经锁定，请先执行异常解锁再处理奖惩"
	case errors.Is(err, service.ErrAdminAdjustmentGroupNotAvailable):
		return "目标小组当前不可直接下发奖惩"
	default:
		return "通知与奖惩处理失败"
	}
}
