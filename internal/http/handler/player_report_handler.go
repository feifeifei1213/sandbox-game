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

type PlayerReportHandler struct {
	queryService   *service.PlayerReportQueryService
	commandService *service.PlayerReportCommandService
}

func NewPlayerReportHandler(
	queryService *service.PlayerReportQueryService,
	commandService *service.PlayerReportCommandService,
) *PlayerReportHandler {
	return &PlayerReportHandler{
		queryService:   queryService,
		commandService: commandService,
	}
}

func (h *PlayerReportHandler) GetView(c *gin.Context) {
	var req dto.PlayerReportGetViewRequest
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
			"当前身份无权访问玩家财报页",
			nil,
		))
		return
	}

	view, err := h.queryService.GetView(c.Request.Context(), *identity.GroupID, *req.YearNo)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到对应财报数据",
				err,
			))
		case errors.Is(err, service.ErrPlayerReportNotOpen):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"当前财报页尚未开放",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"获取财报页视图失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(view))
}

// SaveDraft 保存当前登录玩家所属组的财报手工项草稿。
func (h *PlayerReportHandler) SaveDraft(c *gin.Context) {
	var req dto.PlayerReportSaveDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"保存财报草稿参数不正确",
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
			"当前身份无权保存财报草稿",
			nil,
		))
		return
	}

	result, err := h.commandService.SaveDraft(c.Request.Context(), service.SavePlayerReportDraftCommand{
		GroupID:             *identity.GroupID,
		YearNo:              *req.YearNo,
		ReportManualPayload: req.ReportManualPayload,
		OperatorName:        identity.Username,
	})
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到对应财报数据",
				err,
			))
		case errors.Is(err, service.ErrPlayerReportDraftNotEditable):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"当前年份不允许保存财报草稿",
				err,
			))
		case errors.Is(err, service.ErrPlayerReportDraftInvalid):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"当前财报草稿不符合保存要求",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"保存财报草稿失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.PlayerReportSaveDraftResponse{
		GroupID:          result.GroupID,
		YearNo:           result.YearNo,
		YearStatus:       result.YearStatus,
		ReportStatus:     result.ReportStatus,
		LastDraftSavedAt: result.LastDraftSavedAt,
	}))
}

// Submit 提交当前登录玩家所属组的财报。
func (h *PlayerReportHandler) Submit(c *gin.Context) {
	var req dto.PlayerReportSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"财报提交参数不正确",
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
			"当前身份无权提交财报",
			nil,
		))
		return
	}

	result, err := h.commandService.Submit(c.Request.Context(), service.SubmitPlayerReportCommand{
		GroupID:             *identity.GroupID,
		YearNo:              *req.YearNo,
		ReportManualPayload: req.ReportManualPayload,
		SubmitterID:         identity.UserID,
		OperatorName:        identity.Username,
	})
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到对应财报数据",
				err,
			))
		case errors.Is(err, service.ErrPlayerReportSubmitConflict):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusConflict,
				enum.ConflictCode,
				"本年度财报已提交，请刷新页面后重试",
				err,
			))
		case errors.Is(err, service.ErrPlayerReportSubmitInvalid):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"当前财报不满足提交要求",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"财报提交失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.PlayerReportSubmitResponse{
		GroupID:                   result.GroupID,
		YearNo:                    result.YearNo,
		YearStatus:                result.YearStatus,
		ReportStatus:              result.ReportStatus,
		BusinessStatus:            result.BusinessStatus,
		SummaryEffective:          result.SummaryEffective,
		LatestReportSubmitVersion: result.LatestReportSubmitVersion,
		BalanceCheckPassed:        result.BalanceCheckPassed,
		SubmittedAt:               result.SubmittedAt,
	}))
}
