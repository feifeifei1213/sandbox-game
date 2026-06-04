package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/http/dto"
	"sandbox-game/internal/http/middleware"
	"sandbox-game/internal/service"
)

type AdminDictionaryHandler struct {
	service *service.AdminDictionaryService
}

func NewAdminDictionaryHandler(service *service.AdminDictionaryService) *AdminDictionaryHandler {
	return &AdminDictionaryHandler{service: service}
}

func (h *AdminDictionaryHandler) GetCurrent(c *gin.Context) {
	result, err := h.service.GetCurrent(c.Request.Context(), c.Query("editionCode"))
	if err != nil {
		abortDictionaryError(c, err, "获取业务显示字典失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminDictionaryHandler) ListSchemes(c *gin.Context) {
	if !ensureAdminIdentity(c, "当前身份无权查看字典方案") {
		return
	}
	result, err := h.service.ListSchemes(c.Request.Context(), c.Query("editionCode"))
	if err != nil {
		abortDictionaryError(c, err, "获取字典方案失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminDictionaryHandler) GetSchemeDetail(c *gin.Context) {
	id, err := parseInt64Query(c, "id")
	if err != nil || id <= 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"id 参数不正确",
			err,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权查看字典方案") {
		return
	}
	result, err := h.service.GetSchemeDetail(c.Request.Context(), id)
	if err != nil {
		abortDictionaryError(c, err, "获取字典方案详情失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminDictionaryHandler) SaveScheme(c *gin.Context) {
	var req dto.AdminDictionarySaveSchemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"保存字典方案参数不正确",
			err,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权保存字典方案") {
		return
	}
	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.service.SaveScheme(c.Request.Context(), service.SaveDictionarySchemeCommand{
		SchemeID:     req.SchemeID,
		EditionCode:  req.EditionCode,
		SchemeName:   req.SchemeName,
		Description:  req.Description,
		Items:        toDictionaryItemInputs(req.Items),
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortDictionaryError(c, err, "保存字典方案失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminDictionaryHandler) DeleteScheme(c *gin.Context) {
	id, err := parseInt64Query(c, "id")
	if err != nil || id <= 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"id 参数不正确",
			err,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权删除字典方案") {
		return
	}
	identity, _ := middleware.GetAuthIdentity(c)
	if err := h.service.DeleteScheme(c.Request.Context(), service.DeleteDictionarySchemeCommand{
		SchemeID:     id,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	}); err != nil {
		abortDictionaryError(c, err, "删除字典方案失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(map[string]any{"deleted": true}))
}

func (h *AdminDictionaryHandler) UpdateCurrent(c *gin.Context) {
	var req dto.AdminDictionaryUpdateCurrentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"更新当前比赛字典参数不正确",
			err,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权修改当前比赛字典") {
		return
	}
	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.service.UpdateCurrent(c.Request.Context(), service.UpdateCurrentDictionaryCommand{
		Items:        toDictionaryItemInputs(req.Items),
		Reason:       req.Reason,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortDictionaryError(c, err, "更新当前比赛字典失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminDictionaryHandler) ApplySchemeToCurrent(c *gin.Context) {
	var req dto.AdminDictionaryApplySchemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"应用字典方案参数不正确",
			err,
		))
		return
	}
	if req.SchemeID == nil || *req.SchemeID <= 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"schemeId 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权应用字典方案") {
		return
	}
	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.service.ApplySchemeToCurrent(c.Request.Context(), service.ApplyDictionarySchemeCommand{
		SchemeID:     *req.SchemeID,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortDictionaryError(c, err, "应用字典方案失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminDictionaryHandler) RestoreCurrentDefault(c *gin.Context) {
	if !ensureAdminIdentity(c, "当前身份无权恢复默认显示名称") {
		return
	}
	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.service.RestoreCurrentDefault(c.Request.Context(), service.RestoreCurrentDictionaryDefaultCommand{
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortDictionaryError(c, err, "恢复默认显示名称失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminDictionaryHandler) PageChangeLogs(c *gin.Context) {
	if !ensureAdminIdentity(c, "当前身份无权查看字典修改记录") {
		return
	}
	pageNo, _ := strconv.Atoi(c.DefaultQuery("pageNo", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.service.PageChangeLogs(c.Request.Context(), pageNo, pageSize)
	if err != nil {
		abortDictionaryError(c, err, "获取字典修改记录失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminDictionaryHandler) GetRevision(c *gin.Context) {
	result, err := h.service.GetRevision(c.Request.Context())
	if err != nil {
		abortDictionaryError(c, err, "获取字典版本失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func abortDictionaryError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusNotFound, enum.NotFoundCode, "未找到字典数据", err))
	case errors.Is(err, service.ErrDictionaryEditionInvalid):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, "沙盘版本不合法", err))
	case errors.Is(err, service.ErrDictionarySchemeNameRequired):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, "字典方案名称不能为空", err))
	case errors.Is(err, service.ErrDictionarySchemeInvalid):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, "字典方案不允许执行该操作", err))
	case errors.Is(err, service.ErrDictionarySchemeCrossEdition):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, "字典方案不属于当前沙盘版本", err))
	case errors.Is(err, service.ErrDictionaryDisplayNameRequired):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, "字典显示名称不能为空", err))
	case errors.Is(err, service.ErrDictionaryItemInvalid):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, "字典项不合法", err))
	case errors.Is(err, service.ErrDictionaryNotInitialized):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusConflict, enum.ConflictCode, "比赛尚未初始化，不能修改当前比赛字典", err))
	default:
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusInternalServerError, enum.InternalServerErrorCode, fallback, err))
	}
}

func parseInt64Query(c *gin.Context, name string) (int64, error) {
	return strconv.ParseInt(c.Query(name), 10, 64)
}
