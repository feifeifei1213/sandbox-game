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
	queryService          *service.AdminOrderQueryService
	commandService        *service.AdminOrderCommandService
	controlQueryService   *service.AdminOrderControlQueryService
	controlCommandService *service.AdminOrderControlCommandService
}

func NewAdminOrderHandler(
	queryService *service.AdminOrderQueryService,
	commandService *service.AdminOrderCommandService,
	controlQueryService *service.AdminOrderControlQueryService,
	controlCommandService *service.AdminOrderControlCommandService,
) *AdminOrderHandler {
	return &AdminOrderHandler{
		queryService:          queryService,
		commandService:        commandService,
		controlQueryService:   controlQueryService,
		controlCommandService: controlCommandService,
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

func (h *AdminOrderHandler) GetForecastControl(c *gin.Context) {
	if !ensureAdminIdentity(c, "当前身份无权查看订单数量控制台") {
		return
	}

	result, err := h.queryService.GetForecastControl(c.Request.Context())
	if err != nil {
		abortAdminOrderError(c, err, "获取订单数量控制台失败")
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

func (h *AdminOrderHandler) UpdateForecastControl(c *gin.Context) {
	var req dto.AdminOrderUpdateForecastControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"订单数量控制台参数不正确",
			err,
		))
		return
	}
	if len(req.Items) == 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"items 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权更新订单数量控制台") {
		return
	}

	items := make([]service.UpdateOrderForecastControlItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, service.UpdateOrderForecastControlItem{
			YearNo:     item.YearNo,
			MarketCode: item.MarketCode,
			OrderType:  item.OrderType,
			OrderCount: item.OrderCount,
		})
	}
	narratives := make([]service.UpdateOrderForecastNarrativeItem, 0, len(req.Narratives))
	for _, item := range req.Narratives {
		narratives = append(narratives, service.UpdateOrderForecastNarrativeItem{
			ForecastStageCode: item.ForecastStageCode,
			MarketCode:        item.MarketCode,
			Content:           item.Content,
		})
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.UpdateForecastControl(c.Request.Context(), service.UpdateOrderForecastControlCommand{
		Items:        items,
		Narratives:   narratives,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortAdminOrderError(c, err, "更新订单数量控制台失败")
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

func (h *AdminOrderHandler) UpdateMarketConfig(c *gin.Context) {
	var req dto.AdminOrderUpdateMarketConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"市场开启配置参数不正确",
			err,
		))
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 || len(req.Markets) == 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 或 markets 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权更新市场开启配置") {
		return
	}

	markets := make([]service.UpdateOrderMarketConfigItem, 0, len(req.Markets))
	for _, item := range req.Markets {
		markets = append(markets, service.UpdateOrderMarketConfigItem{
			MarketCode:            item.MarketCode,
			Enabled:               item.Enabled,
			MarketInvestmentLimit: item.MarketInvestmentLimit,
		})
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.UpdateMarketConfig(c.Request.Context(), service.UpdateOrderMarketConfigCommand{
		YearNo:       *req.YearNo,
		Markets:      markets,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortAdminOrderError(c, err, "更新市场开启配置失败")
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

func (h *AdminOrderHandler) ConfirmOrderPool(c *gin.Context) {
	var req dto.AdminOrderConfirmPoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"确认订单池参数不正确",
			err,
		))
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 || req.BatchID == nil || *req.BatchID <= 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 或 batchId 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权确认订单池") {
		return
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.ConfirmOrderPool(c.Request.Context(), service.ConfirmOrderPoolCommand{
		YearNo:       *req.YearNo,
		BatchID:      *req.BatchID,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortAdminOrderError(c, err, "确认订单池失败")
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminOrderHandler) GenerateSelectionSequence(c *gin.Context) {
	var req dto.AdminOrderGenerateSelectionSequenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortPlayerOrderBadRequest(c, "生成选单顺序参数不正确", err)
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 {
		abortPlayerOrderBadRequest(c, "yearNo 参数不正确", nil)
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权生成选单顺序") {
		return
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.controlCommandService.GenerateSelectionSequence(c.Request.Context(), service.GenerateSelectionSequenceCommand{
		YearNo:       *req.YearNo,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortAdminOrderError(c, err, "生成选单顺序失败")
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

func (h *AdminOrderHandler) OpenMarketBidding(c *gin.Context) {
	var req dto.AdminOrderOpenMarketBiddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortPlayerOrderBadRequest(c, "开放市场投入参数不正确", err)
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 {
		abortPlayerOrderBadRequest(c, "yearNo 参数不正确", nil)
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权开放市场投入") {
		return
	}
	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.controlCommandService.OpenMarketBidding(c.Request.Context(), service.OpenMarketBiddingCommand{
		YearNo:       *req.YearNo,
		MarketCode:   req.MarketCode,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortOrderError(c, err, "开放市场投入失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminOrderHandler) CloseMarketBidding(c *gin.Context) {
	var req dto.AdminOrderCloseMarketBiddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortPlayerOrderBadRequest(c, "关闭市场投入参数不正确", err)
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 {
		abortPlayerOrderBadRequest(c, "yearNo 参数不正确", nil)
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权关闭市场投入") {
		return
	}
	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.controlCommandService.CloseMarketBidding(c.Request.Context(), service.CloseMarketBiddingCommand{
		YearNo:       *req.YearNo,
		MarketCode:   req.MarketCode,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortOrderError(c, err, "关闭市场投入失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminOrderHandler) GetMarketSelectionStatus(c *gin.Context) {
	var req dto.AdminOrderGetMarketSelectionStatusRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		abortPlayerOrderBadRequest(c, "市场选单状态查询参数不正确", err)
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 {
		abortPlayerOrderBadRequest(c, "yearNo 参数不正确", nil)
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权查看市场选单状态") {
		return
	}
	result, err := h.controlQueryService.GetMarketSelectionStatus(c.Request.Context(), *req.YearNo, req.MarketCode)
	if err != nil {
		abortOrderError(c, err, "获取市场选单状态失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminOrderHandler) ReleaseNextSegment(c *gin.Context) {
	var req dto.AdminOrderReleaseNextSegmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortPlayerOrderBadRequest(c, "释放标段参数不正确", err)
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 {
		abortPlayerOrderBadRequest(c, "yearNo 参数不正确", nil)
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权释放标段") {
		return
	}
	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.controlCommandService.ReleaseNextSegment(c.Request.Context(), service.ReleaseNextSegmentCommand{
		YearNo:       *req.YearNo,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortOrderError(c, err, "释放下一个标段失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminOrderHandler) AdminSkipCurrentGroup(c *gin.Context) {
	var req dto.AdminOrderSkipCurrentGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortPlayerOrderBadRequest(c, "管理员跳过参数不正确", err)
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 || req.GroupID == nil || *req.GroupID <= 0 {
		abortPlayerOrderBadRequest(c, "yearNo 或 groupId 参数不正确", nil)
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权跳过当前小组") {
		return
	}
	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.controlCommandService.AdminSkipCurrentGroup(c.Request.Context(), service.AdminSkipCurrentGroupCommand{
		YearNo:       *req.YearNo,
		MarketCode:   req.MarketCode,
		OrderType:    req.OrderType,
		GroupID:      *req.GroupID,
		Reason:       req.Reason,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortOrderError(c, err, "跳过当前小组失败")
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
		errors.Is(err, service.ErrAdminOrderPoolLocked),
		errors.Is(err, service.ErrAdminOrderSequenceAlreadyGenerated),
		errors.Is(err, service.ErrAdminOrderMarketConfigLocked):
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusConflict,
			enum.ConflictCode,
			resolveAdminOrderErrorMessage(err),
			err,
		))
	case errors.Is(err, service.ErrAdminOrderBatchNotFound),
		errors.Is(err, service.ErrAdminOrderConfigNotFound),
		errors.Is(err, service.ErrAdminOrderPreviewNotFound),
		errors.Is(err, service.ErrAdminOrderPoolNotConfirmed):
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusNotFound,
			enum.NotFoundCode,
			resolveAdminOrderErrorMessage(err),
			err,
		))
	case errors.Is(err, service.ErrAdminOrderYearInvalid),
		errors.Is(err, service.ErrAdminOrderParseFailed),
		errors.Is(err, service.ErrAdminOrderConfigInvalid),
		errors.Is(err, service.ErrAdminOrderMarketInvestmentLimitInvalid),
		errors.Is(err, service.ErrAdminOrderForecastControlInvalid),
		errors.Is(err, service.ErrAdminOrderReleaseSequenceDuplicated),
		errors.Is(err, service.ErrAdminOrderSourceInsufficient),
		errors.Is(err, service.ErrAdminOrderInvestmentIncomplete):
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
	case errors.Is(err, service.ErrAdminOrderMarketConfigLocked):
		return "订单池已确认或标段已进入开标流程，不能修改市场开启状态"
	case errors.Is(err, service.ErrAdminOrderYearInvalid):
		return "年份不在可配置订单范围内"
	case errors.Is(err, service.ErrAdminOrderParseFailed):
		return "订单 Excel 解析失败，请确认文件格式与订单推算模板一致"
	case errors.Is(err, service.ErrAdminOrderConfigInvalid):
		return "订单数量或标段释放顺序配置不合法"
	case errors.Is(err, service.ErrAdminOrderMarketInvestmentLimitInvalid):
		return "市场投入上限必须为空或非负整数"
	case errors.Is(err, service.ErrAdminOrderForecastControlInvalid):
		return "订单数量控制台配置不合法"
	case errors.Is(err, service.ErrAdminOrderReleaseSequenceDuplicated):
		return "同一年内标段释放顺序不能重复"
	case errors.Is(err, service.ErrAdminOrderReleaseSequenceLocked):
		return "该年份已有标段进入开标流程，不能修改释放顺序"
	case errors.Is(err, service.ErrAdminOrderPoolLocked):
		return "该年份订单池已确认或已有订单被选择，不能覆盖"
	case errors.Is(err, service.ErrAdminOrderConfigNotFound):
		return "请先保存订单数量与标段释放顺序配置"
	case errors.Is(err, service.ErrAdminOrderBatchNotFound):
		return "未找到可用订单批次"
	case errors.Is(err, service.ErrAdminOrderPreviewNotFound):
		return "请先生成本年度预览订单池"
	case errors.Is(err, service.ErrAdminOrderPoolNotConfirmed):
		return "请先确认本年度订单池"
	case errors.Is(err, service.ErrAdminOrderInvestmentIncomplete):
		if err.Error() != service.ErrAdminOrderInvestmentIncomplete.Error() {
			return err.Error()
		}
		return "仍有未破产小组没有提交完整 16 项市场投入"
	case errors.Is(err, service.ErrAdminOrderSequenceAlreadyGenerated):
		return "该年份已生成选单顺序，不能重复生成"
	case errors.Is(err, service.ErrAdminOrderSourceInsufficient):
		return err.Error()
	default:
		return "订单管理处理失败"
	}
}
