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

type AdminSummaryHandler struct {
	queryService *service.AdminSummaryQueryService
}

func NewAdminSummaryHandler(queryService *service.AdminSummaryQueryService) *AdminSummaryHandler {
	return &AdminSummaryHandler{queryService: queryService}
}

func (h *AdminSummaryHandler) GetYearSummary(c *gin.Context) {
	var req dto.AdminSummaryGetYearSummaryRequest
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

	if !ensureAdminIdentity(c, "当前身份无权查看管理员汇总") {
		return
	}

	result, err := h.queryService.GetYearSummary(c.Request.Context(), *req.YearNo)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到管理员汇总数据",
				err,
			))
		case errors.Is(err, service.ErrAdminSummaryYearInvalid):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"当前年份不属于正式汇总范围",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"获取年度汇总失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminSummaryHandler) GetFinalRanking(c *gin.Context) {
	if !ensureAdminIdentity(c, "当前身份无权查看最终排名") {
		return
	}

	result, err := h.queryService.GetFinalRanking(c.Request.Context())
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到最终排名数据",
				err,
			))
		case errors.Is(err, service.ErrFinalRankingNotReady):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"最终年份尚未整体完成，暂不可查看最终排名",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"获取最终排名失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func ensureAdminIdentity(c *gin.Context, forbiddenMsg string) bool {
	identity, ok := middleware.GetAuthIdentity(c)
	if !ok {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusUnauthorized,
			enum.UnauthorizedCode,
			"未获取到当前登录身份",
			nil,
		))
		return false
	}
	if identity.RoleType != enum.RoleTypeAdmin {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusForbidden,
			enum.ForbiddenCode,
			forbiddenMsg,
			nil,
		))
		return false
	}
	return true
}
