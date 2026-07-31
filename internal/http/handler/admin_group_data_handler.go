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

type AdminGroupDataHandler struct {
	queryService *service.AdminGroupDataQueryService
}

func NewAdminGroupDataHandler(queryService *service.AdminGroupDataQueryService) *AdminGroupDataHandler {
	return &AdminGroupDataHandler{queryService: queryService}
}

func (h *AdminGroupDataHandler) ListGroups(c *gin.Context) {
	if !ensureAdminIdentity(c, "当前身份无权查看小组列表") {
		return
	}

	result, err := h.queryService.ListGroups(c.Request.Context())
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到小组列表数据",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"获取小组列表失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminGroupDataHandler) GetOperationContext(c *gin.Context) {
	var req dto.AdminGroupDataGetOperationContextRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"groupId 参数不正确",
			err,
		))
		return
	}
	if req.GroupID == nil || *req.GroupID <= 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"groupId 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权查看小组操作上下文") {
		return
	}

	result, err := h.queryService.GetOperationContext(c.Request.Context(), *req.GroupID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到目标小组或年份状态",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"获取小组操作上下文失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminGroupDataHandler) GetOperatingView(c *gin.Context) {
	var req dto.AdminGroupDataGetOperatingViewRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"groupId 或 yearNo 参数不正确",
			err,
		))
		return
	}
	if req.GroupID == nil || *req.GroupID <= 0 || req.YearNo == nil || *req.YearNo < 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"groupId 或 yearNo 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权查看组经营数据") {
		return
	}

	result, err := h.queryService.GetOperatingView(c.Request.Context(), *req.GroupID, *req.YearNo)
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
				"获取组经营页视图失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminGroupDataHandler) GetReportView(c *gin.Context) {
	var req dto.AdminGroupDataGetReportViewRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"groupId 或 yearNo 参数不正确",
			err,
		))
		return
	}
	if req.GroupID == nil || *req.GroupID <= 0 || req.YearNo == nil || *req.YearNo < 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"groupId 或 yearNo 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权查看组财报数据") {
		return
	}

	result, err := h.queryService.GetReportView(c.Request.Context(), *req.GroupID, *req.YearNo)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到对应财报数据",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"获取组财报页视图失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}
