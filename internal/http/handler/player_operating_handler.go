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

type PlayerOperatingHandler struct {
	queryService   *service.PlayerOperatingQueryService
	commandService *service.PlayerOperatingCommandService
}

func NewPlayerOperatingHandler(
	queryService *service.PlayerOperatingQueryService,
	commandService *service.PlayerOperatingCommandService,
) *PlayerOperatingHandler {
	return &PlayerOperatingHandler{
		queryService:   queryService,
		commandService: commandService,
	}
}

func (h *PlayerOperatingHandler) GetYearView(c *gin.Context) {
	var req dto.PlayerOperatingGetYearViewRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 参数不正确",
			err,
		))
		return
	}
	if req.YearNo == nil || *req.YearNo < 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 参数不正确",
			nil,
		))
		return
	}

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

	if identity.RoleType != enum.RoleTypeGroup || identity.GroupID == nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusForbidden,
			enum.ForbiddenCode,
			"当前身份无权访问玩家经营页",
			nil,
		))
		return
	}

	view, err := h.queryService.GetYearView(c.Request.Context(), *identity.GroupID, *req.YearNo)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到对应经营页数据",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"获取经营页视图失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(view))
}

func (h *PlayerOperatingHandler) SaveDraft(c *gin.Context) {
	var req dto.PlayerOperatingSaveDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"保存草稿参数不正确",
			err,
		))
		return
	}
	if req.YearNo == nil || *req.YearNo < 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 参数不正确",
			nil,
		))
		return
	}

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
	if identity.RoleType != enum.RoleTypeGroup || identity.GroupID == nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusForbidden,
			enum.ForbiddenCode,
			"当前身份无权保存经营页草稿",
			nil,
		))
		return
	}

	result, err := h.commandService.SaveDraft(c.Request.Context(), service.SaveOperatingDraftCommand{
		GroupID:          *identity.GroupID,
		YearNo:           *req.YearNo,
		StageStatus:      req.StageStatus,
		OperatingPayload: req.OperatingPayload,
		OperatorName:     identity.Username,
	})
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到对应经营页数据",
				err,
			))
		case errors.Is(err, service.ErrOperatingDraftStageConflict):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusConflict,
				enum.ConflictCode,
				"当前阶段已变化，请刷新页面后重试",
				err,
			))
		case errors.Is(err, service.ErrOperatingDraftNotEditable), errors.Is(err, service.ErrOperatingDraftMissingCarryBase):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"当前年份不允许保存经营页草稿",
				err,
			))
		case errors.Is(err, service.ErrManualNumberNotInteger), errors.Is(err, service.ErrSupplyChainOrderQuantityInvalid):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				err.Error(),
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"保存经营页草稿失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.PlayerOperatingSaveDraftResponse{
		GroupID:          result.GroupID,
		YearNo:           result.YearNo,
		StageStatus:      result.StageStatus,
		LastDraftSavedAt: result.LastDraftSavedAt,
	}))
}

func (h *PlayerOperatingHandler) SubmitStage(c *gin.Context) {
	var req dto.PlayerOperatingSubmitStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"阶段提交参数不正确",
			err,
		))
		return
	}
	if req.YearNo == nil || *req.YearNo < 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 参数不正确",
			nil,
		))
		return
	}

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
	if identity.RoleType != enum.RoleTypeGroup || identity.GroupID == nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusForbidden,
			enum.ForbiddenCode,
			"当前身份无权提交经营阶段",
			nil,
		))
		return
	}

	result, err := h.commandService.SubmitStage(c.Request.Context(), service.SubmitOperatingStageCommand{
		GroupID:          *identity.GroupID,
		YearNo:           *req.YearNo,
		StageCode:        req.StageCode,
		OperatingPayload: req.OperatingPayload,
		SubmitterID:      identity.UserID,
		OperatorName:     identity.Username,
	})
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到对应经营页数据",
				err,
			))
		case errors.Is(err, service.ErrOperatingDraftStageConflict):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusConflict,
				enum.ConflictCode,
				"当前阶段已变化或已提交，请刷新页面后重试",
				err,
			))
		case errors.Is(err, service.ErrOperatingStageSubmitInvalid), errors.Is(err, service.ErrOperatingDraftMissingCarryBase):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"当前阶段不满足提交流程要求",
				err,
			))
		case errors.Is(err, service.ErrOrderPrerequisiteIncomplete):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"本年订单选择尚未完成，不能提交 Q1",
				err,
			))
		case errors.Is(err, service.ErrManualNumberNotInteger), errors.Is(err, service.ErrSupplyChainOrderQuantityInvalid):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				err.Error(),
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"经营阶段提交失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.PlayerOperatingSubmitStageResponse{
		GroupID:                  result.GroupID,
		YearNo:                   result.YearNo,
		StageCode:                result.StageCode,
		YearStatus:               result.YearStatus,
		StageStatus:              result.StageStatus,
		ReportStatus:             result.ReportStatus,
		BusinessStatus:           result.BusinessStatus,
		LatestStageSubmitVersion: result.LatestStageSubmitVersion,
		PeriodEndCash:            result.PeriodEndCash,
		SubmittedAt:              result.SubmittedAt,
	}))
}
