package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/http/dto"
	"sandbox-game/internal/http/middleware"
	"sandbox-game/internal/service"
)

type AdminOrderHandler struct {
	queryService   *service.AdminOrderQueryService
	commandService *service.AdminOrderCommandService
}

func NewAdminOrderHandler(
	queryService *service.AdminOrderQueryService,
	commandService *service.AdminOrderCommandService,
) *AdminOrderHandler {
	return &AdminOrderHandler{
		queryService:   queryService,
		commandService: commandService,
	}
}

func (h *AdminOrderHandler) GetControlConfig(c *gin.Context) {
	var req dto.AdminOrderGetControlConfigRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 参数不正确",
			err,
		))
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权查看订单配置") {
		return
	}

	result, err := h.queryService.GetControlConfig(c.Request.Context(), *req.YearNo)
	if err != nil {
		abortAdminOrderError(c, err, "获取订单配置失败")
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminOrderHandler) UploadExcel(c *gin.Context) {
	if !ensureAdminIdentity(c, "当前身份无权上传订单 Excel") {
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"请上传订单 Excel 文件",
			err,
		))
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"读取订单 Excel 文件失败",
			err,
		))
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"读取订单 Excel 文件失败",
			err,
		))
		return
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.UploadExcel(c.Request.Context(), service.UploadOrderExcelCommand{
		FileName:     fileHeader.Filename,
		FileSize:     fileHeader.Size,
		Content:      content,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortAdminOrderError(c, err, "上传订单 Excel 失败")
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminOrderHandler) UpdateControlConfig(c *gin.Context) {
	var req dto.AdminOrderUpdateControlConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"订单数量配置参数不正确",
			err,
		))
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 || len(req.Items) == 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 或 items 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权更新订单配置") {
		return
	}

	items := make([]service.UpdateOrderControlConfigItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, service.UpdateOrderControlConfigItem{
			MarketCode:        item.MarketCode,
			OrderType:         item.OrderType,
			OrderCount:        item.OrderCount,
			ReleaseSequenceNo: item.ReleaseSequenceNo,
		})
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.UpdateControlConfig(c.Request.Context(), service.UpdateOrderControlConfigCommand{
		YearNo:       *req.YearNo,
		Items:        items,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortAdminOrderError(c, err, "更新订单配置失败")
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminOrderHandler) GenerateOrderPool(c *gin.Context) {
	var req dto.AdminOrderGeneratePoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"生成订单池参数不正确",
			err,
		))
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权生成订单池") {
		return
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.GenerateOrderPool(c.Request.Context(), service.GenerateOrderPoolCommand{
		YearNo:        *req.YearNo,
		Overwrite:     req.Overwrite,
		SourceBatchID: req.SourceBatchID,
		OperatorID:    identity.UserID,
		OperatorName:  identity.Username,
	})
	if err != nil {
		abortAdminOrderError(c, err, "生成订单池失败")
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminOrderHandler) GetOrderPool(c *gin.Context) {
	var req dto.AdminOrderGetOrderPoolRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"订单池查询参数不正确",
			err,
		))
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权查看订单池") {
		return
	}

	result, err := h.queryService.GetOrderPool(c.Request.Context(), *req.YearNo, req.MarketCode, req.OrderType)
	if err != nil {
		abortAdminOrderError(c, err, "获取订单池失败")
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func abortAdminOrderError(c *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusNotFound,
			enum.NotFoundCode,
			"未找到订单相关数据",
			err,
		))
	case errors.Is(err, service.ErrAdminOrderFileRequired):
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"请上传订单 Excel 文件",
			err,
		))
	case errors.Is(err, service.ErrAdminOrderReleaseSequenceLocked),
		errors.Is(err, service.ErrAdminOrderPoolLocked):
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusConflict,
			enum.ConflictCode,
			resolveAdminOrderErrorMessage(err),
			err,
		))
	case errors.Is(err, service.ErrAdminOrderBatchNotFound),
		errors.Is(err, service.ErrAdminOrderConfigNotFound):
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusNotFound,
			enum.NotFoundCode,
			resolveAdminOrderErrorMessage(err),
			err,
		))
	case errors.Is(err, service.ErrAdminOrderYearInvalid),
		errors.Is(err, service.ErrAdminOrderParseFailed),
		errors.Is(err, service.ErrAdminOrderConfigInvalid),
		errors.Is(err, service.ErrAdminOrderReleaseSequenceDuplicated),
		errors.Is(err, service.ErrAdminOrderSourceInsufficient):
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusUnprocessableEntity,
			enum.UnprocessableEntityCode,
			resolveAdminOrderErrorMessage(err),
			err,
		))
	default:
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusInternalServerError,
			enum.InternalServerErrorCode,
			fallbackMessage,
			err,
		))
	}
}

func resolveAdminOrderErrorMessage(err error) string {
	switch {
	case errors.Is(err, service.ErrAdminOrderYearInvalid):
		return "年份不在可配置订单范围内"
	case errors.Is(err, service.ErrAdminOrderParseFailed):
		return "订单 Excel 解析失败，请确认文件格式与订单推算模板一致"
	case errors.Is(err, service.ErrAdminOrderConfigInvalid):
		return "订单数量或标段释放顺序配置不合法"
	case errors.Is(err, service.ErrAdminOrderReleaseSequenceDuplicated):
		return "同一年内标段释放顺序不能重复"
	case errors.Is(err, service.ErrAdminOrderReleaseSequenceLocked):
		return "该年份已有标段进入开标流程，不能修改释放顺序"
	case errors.Is(err, service.ErrAdminOrderPoolLocked):
		return "该年份订单池已开放或已有订单被选择，不能覆盖"
	case errors.Is(err, service.ErrAdminOrderConfigNotFound):
		return "请先保存订单数量与标段释放顺序配置"
	case errors.Is(err, service.ErrAdminOrderBatchNotFound):
		return "请先上传并成功解析订单 Excel"
	case errors.Is(err, service.ErrAdminOrderSourceInsufficient):
		return err.Error()
	default:
		return "订单管理处理失败"
	}
}
