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

type PlayerOrderHandler struct {
	queryService   *service.PlayerOrderQueryService
	commandService *service.PlayerOrderCommandService
}

func NewPlayerOrderHandler(
	queryService *service.PlayerOrderQueryService,
	commandService *service.PlayerOrderCommandService,
) *PlayerOrderHandler {
	return &PlayerOrderHandler{
		queryService:   queryService,
		commandService: commandService,
	}
}

func (h *PlayerOrderHandler) GetYearView(c *gin.Context) {
	var req dto.PlayerOrderGetYearViewRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		abortPlayerOrderBadRequest(c, "yearNo 参数不正确", err)
		return
	}
	if req.YearNo == nil || *req.YearNo < 0 {
		abortPlayerOrderBadRequest(c, "yearNo 参数不正确", nil)
		return
	}
	identity, ok := requireGroupIdentity(c, "当前身份无权访问玩家订单页")
	if !ok {
		return
	}
	result, err := h.queryService.GetYearView(c.Request.Context(), *identity.GroupID, *req.YearNo)
	if err != nil {
		abortOrderError(c, err, "获取玩家订单页失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *PlayerOrderHandler) SubmitMarketInvestment(c *gin.Context) {
	var req dto.PlayerOrderSubmitMarketInvestmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortPlayerOrderBadRequest(c, "市场投入参数不正确", err)
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 {
		abortPlayerOrderBadRequest(c, "yearNo 参数不正确", nil)
		return
	}
	identity, ok := requireGroupIdentity(c, "当前身份无权提交市场投入")
	if !ok {
		return
	}
	result, err := h.commandService.SubmitMarketInvestment(c.Request.Context(), service.SubmitMarketInvestmentCommand{
		GroupID:          *identity.GroupID,
		YearNo:           *req.YearNo,
		MarketCode:       req.MarketCode,
		MarketInvestment: req.MarketInvestment,
		OperatorName:     identity.Username,
	})
	if err != nil {
		abortOrderError(c, err, "提交市场投入失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *PlayerOrderHandler) SelectOrder(c *gin.Context) {
	var req dto.PlayerOrderSelectOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortPlayerOrderBadRequest(c, "选择订单参数不正确", err)
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 || req.OrderID == nil || *req.OrderID <= 0 {
		abortPlayerOrderBadRequest(c, "yearNo 或 orderId 参数不正确", nil)
		return
	}
	identity, ok := requireGroupIdentity(c, "当前身份无权选择订单")
	if !ok {
		return
	}
	result, err := h.commandService.SelectOrder(c.Request.Context(), service.SelectOrderCommand{
		GroupID:      *identity.GroupID,
		YearNo:       *req.YearNo,
		MarketCode:   req.MarketCode,
		OrderType:    req.OrderType,
		OrderID:      *req.OrderID,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortOrderError(c, err, "选择订单失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *PlayerOrderHandler) PassSegment(c *gin.Context) {
	var req dto.PlayerOrderPassSegmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortPlayerOrderBadRequest(c, "放弃标段参数不正确", err)
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 {
		abortPlayerOrderBadRequest(c, "yearNo 参数不正确", nil)
		return
	}
	identity, ok := requireGroupIdentity(c, "当前身份无权放弃标段")
	if !ok {
		return
	}
	result, err := h.commandService.PassSegment(c.Request.Context(), service.PassOrderSegmentCommand{
		GroupID:      *identity.GroupID,
		YearNo:       *req.YearNo,
		MarketCode:   req.MarketCode,
		OrderType:    req.OrderType,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortOrderError(c, err, "放弃标段失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *PlayerOrderHandler) DeliverOrders(c *gin.Context) {
	var req dto.PlayerOrderDeliverOrdersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		abortPlayerOrderBadRequest(c, "交付订单参数不正确", err)
		return
	}
	if req.YearNo == nil || *req.YearNo < 1 || len(req.OrderIDs) == 0 {
		abortPlayerOrderBadRequest(c, "yearNo 或 orderIds 参数不正确", nil)
		return
	}
	identity, ok := requireGroupIdentity(c, "当前身份无权交付订单")
	if !ok {
		return
	}
	result, err := h.commandService.DeliverOrders(c.Request.Context(), service.DeliverOrdersCommand{
		GroupID:      *identity.GroupID,
		YearNo:       *req.YearNo,
		StageCode:    req.StageCode,
		OrderIDs:     req.OrderIDs,
		OperatorName: identity.Username,
	})
	if err != nil {
		abortOrderError(c, err, "交付订单失败")
		return
	}
	c.JSON(http.StatusOK, dto.Success(result))
}

func requireGroupIdentity(c *gin.Context, forbiddenMsg string) (middleware.AuthIdentity, bool) {
	identity, ok := middleware.GetAuthIdentity(c)
	if !ok {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusUnauthorized,
			enum.UnauthorizedCode,
			"未获取到当前登录身份",
			nil,
		))
		return middleware.AuthIdentity{}, false
	}
	if identity.RoleType != enum.RoleTypeGroup || identity.GroupID == nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusForbidden,
			enum.ForbiddenCode,
			forbiddenMsg,
			nil,
		))
		return middleware.AuthIdentity{}, false
	}
	return identity, true
}

func abortPlayerOrderBadRequest(c *gin.Context, message string, err error) {
	middleware.AbortWithAppError(c, middleware.NewAppError(
		http.StatusBadRequest,
		enum.BadRequestCode,
		message,
		err,
	))
}

func abortOrderError(c *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusNotFound, enum.NotFoundCode, "未找到订单相关数据", err))
	case errors.Is(err, service.ErrOrderInvestmentAlreadySubmitted),
		errors.Is(err, service.ErrOrderMarketAlreadyOpen),
		errors.Is(err, service.ErrOrderMarketAlreadyClosed),
		errors.Is(err, service.ErrOrderAlreadySelected),
		errors.Is(err, service.ErrOrderSegmentReleaseBlocked),
		errors.Is(err, service.ErrOrderSegmentCurrentGroupMismatch):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusConflict, enum.ConflictCode, resolveOrderErrorMessage(err), err))
	case errors.Is(err, service.ErrOrderNotRequiredForDemoYear),
		errors.Is(err, service.ErrOrderYearInvalid),
		errors.Is(err, service.ErrOrderGroupUnavailable),
		errors.Is(err, service.ErrOrderMarketInvalid),
		errors.Is(err, service.ErrOrderTypeInvalid),
		errors.Is(err, service.ErrOrderPoolNotGenerated),
		errors.Is(err, service.ErrOrderMarketNotOpen),
		errors.Is(err, service.ErrOrderInvestmentInvalid),
		errors.Is(err, service.ErrOrderSegmentNotReady),
		errors.Is(err, service.ErrOrderSegmentNotSelecting),
		errors.Is(err, service.ErrOrderSelectionNotEligible),
		errors.Is(err, service.ErrOrderCannotSelect),
		errors.Is(err, service.ErrOrderAdminSkipReasonRequired),
		errors.Is(err, service.ErrOrderDeliveryStageInvalid),
		errors.Is(err, service.ErrOrderDeliveryOrderInvalid),
		errors.Is(err, service.ErrOrderDeliveryRevenueMismatch),
		errors.Is(err, service.ErrOrderPrerequisiteIncomplete):
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusUnprocessableEntity, enum.UnprocessableEntityCode, resolveOrderErrorMessage(err), err))
	default:
		middleware.AbortWithAppError(c, middleware.NewAppError(http.StatusInternalServerError, enum.InternalServerErrorCode, fallbackMessage, err))
	}
}

func resolveOrderErrorMessage(err error) string {
	switch {
	case errors.Is(err, service.ErrOrderNotRequiredForDemoYear):
		return "0年不需要订单"
	case errors.Is(err, service.ErrOrderYearInvalid):
		return "年份不在订单流程范围内"
	case errors.Is(err, service.ErrOrderGroupUnavailable):
		return "当前小组不可参与订单流程"
	case errors.Is(err, service.ErrOrderMarketInvalid):
		return "市场不合法"
	case errors.Is(err, service.ErrOrderTypeInvalid):
		return "订单类型不合法"
	case errors.Is(err, service.ErrOrderPoolNotGenerated):
		return "订单池尚未生成"
	case errors.Is(err, service.ErrOrderMarketNotOpen):
		return "该市场暂未开放投入"
	case errors.Is(err, service.ErrOrderMarketAlreadyOpen):
		return "该市场已开放投入"
	case errors.Is(err, service.ErrOrderMarketAlreadyClosed):
		return "该市场投入已关闭，不能重复操作"
	case errors.Is(err, service.ErrOrderInvestmentAlreadySubmitted):
		return "该市场投入已提交，不能修改"
	case errors.Is(err, service.ErrOrderInvestmentInvalid):
		return "市场投入不能为负数"
	case errors.Is(err, service.ErrOrderSegmentNotReady):
		return "当前没有可释放标段"
	case errors.Is(err, service.ErrOrderSegmentNotSelecting):
		return "当前标段未处于选单中"
	case errors.Is(err, service.ErrOrderSegmentReleaseBlocked):
		return "当前标段未完成，不能释放下一个标段"
	case errors.Is(err, service.ErrOrderSegmentCurrentGroupMismatch):
		return "当前未轮到该小组操作"
	case errors.Is(err, service.ErrOrderSelectionNotEligible):
		return "当前小组不具备该标段选单资格"
	case errors.Is(err, service.ErrOrderAlreadySelected):
		return "当前小组已在该标段选择订单"
	case errors.Is(err, service.ErrOrderCannotSelect):
		return "该订单不可选择"
	case errors.Is(err, service.ErrOrderAdminSkipReasonRequired):
		return "管理员跳过必须选择当前小组并填写原因"
	case errors.Is(err, service.ErrOrderDeliveryStageInvalid):
		return "交付季度必须等于当前经营季度"
	case errors.Is(err, service.ErrOrderDeliveryOrderInvalid):
		return "只能交付本组本年已选且未交付订单"
	case errors.Is(err, service.ErrOrderDeliveryRevenueMismatch):
		return "本季度销售收入必须等于交付订单金额合计"
	case errors.Is(err, service.ErrOrderPrerequisiteIncomplete):
		return "本年订单选择尚未完成"
	default:
		return "订单流程处理失败"
	}
}
