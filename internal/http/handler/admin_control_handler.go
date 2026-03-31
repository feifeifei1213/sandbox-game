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

type AdminControlHandler struct {
	queryService   *service.AdminControlQueryService
	commandService *service.AdminControlCommandService
}

func NewAdminControlHandler(
	queryService *service.AdminControlQueryService,
	commandService *service.AdminControlCommandService,
) *AdminControlHandler {
	return &AdminControlHandler{
		queryService:   queryService,
		commandService: commandService,
	}
}

func (h *AdminControlHandler) GetConfig(c *gin.Context) {
	if !ensureAdminIdentity(c, "当前身份无权查看年度控制配置") {
		return
	}

	result, err := h.queryService.GetConfig(c.Request.Context())
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到年度控制配置",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"获取年度控制配置失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminControlHandler) UpdateFinalYear(c *gin.Context) {
	var req dto.AdminControlUpdateFinalYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"更新最终年份参数不正确",
			err,
		))
		return
	}
	if req.FinalYear == nil || *req.FinalYear < 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"finalYear 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权更新最终年份") {
		return
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.UpdateFinalYear(c.Request.Context(), service.UpdateFinalYearCommand{
		FinalYear:    *req.FinalYear,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到游戏配置",
				err,
			))
		case errors.Is(err, service.ErrAdminControlFinalYearTooSmall):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"最终年份不能小于当前开放年份",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"更新最终年份失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminControlHandler) OpenNextYear(c *gin.Context) {
	var req dto.AdminControlOpenNextYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"开放下一年参数不正确",
			err,
		))
		return
	}
	if req.TargetYearNo == nil || *req.TargetYearNo < 1 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"targetYearNo 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权开放下一年") {
		return
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.OpenNextYear(c.Request.Context(), service.OpenNextYearCommand{
		TargetYearNo: *req.TargetYearNo,
		OperatorID:   identity.UserID,
		OperatorName: identity.Username,
	})
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到需要开放的年份数据",
				err,
			))
		case errors.Is(err, service.ErrAdminControlTargetYearMismatch):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusConflict,
				enum.ConflictCode,
				"开放年份与当前状态不一致",
				err,
			))
		case errors.Is(err, service.ErrAdminControlFinalYearReached):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusConflict,
				enum.ConflictCode,
				"已达到最终年份，无法继续开放",
				err,
			))
		case errors.Is(err, service.ErrAdminControlOpenNextYearBlocked):
			msg := err.Error()
			if msg == "" || msg == service.ErrAdminControlOpenNextYearBlocked.Error() {
				msg = "当前仍有未完成财报的小组"
			}
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				msg,
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"开放下一年失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminControlHandler) GetInitialBaseline(c *gin.Context) {
	if !ensureAdminIdentity(c, "当前身份无权查看初始基线") {
		return
	}

	result, err := h.queryService.GetInitialBaseline(c.Request.Context())
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到初始基线数据",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"获取初始基线失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminControlHandler) SubmitInitialBaseline(c *gin.Context) {
	var req dto.AdminControlSubmitInitialBaselineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"提交初始基线参数不正确",
			err,
		))
		return
	}
	if req.BaselinePayload == nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"baselinePayload 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权提交初始基线") {
		return
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.SubmitInitialBaseline(c.Request.Context(), service.SubmitInitialBaselineCommand{
		BaselinePayload: req.BaselinePayload,
		OperatorID:      identity.UserID,
		OperatorName:    identity.Username,
	})
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到游戏配置",
				err,
			))
		case errors.Is(err, service.ErrAdminControlInitialBaselineSubmitted):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusConflict,
				enum.ConflictCode,
				"初始基线已提交，不能重复提交",
				err,
			))
		case errors.Is(err, service.ErrAdminControlInitialBaselineInvalid):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"初始基线数据不完整",
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"提交初始基线失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}

func (h *AdminControlHandler) UnlockYear(c *gin.Context) {
	var req dto.AdminControlUnlockYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"异常解锁参数不正确",
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
	if req.YearNo == nil || *req.YearNo < 0 {
		middleware.AbortWithAppError(c, middleware.NewAppError(
			http.StatusBadRequest,
			enum.BadRequestCode,
			"yearNo 参数不正确",
			nil,
		))
		return
	}
	if !ensureAdminIdentity(c, "当前身份无权执行异常解锁") {
		return
	}

	targetStageCode := ""
	if req.TargetStageCode != nil {
		targetStageCode = *req.TargetStageCode
	}

	identity, _ := middleware.GetAuthIdentity(c)
	result, err := h.commandService.UnlockYear(c.Request.Context(), service.UnlockYearCommand{
		GroupID:          *req.GroupID,
		YearNo:           *req.YearNo,
		UnlockTargetType: req.UnlockTargetType,
		TargetStageCode:  targetStageCode,
		Reason:           req.Reason,
		OperatorID:       identity.UserID,
		OperatorName:     identity.Username,
	})
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusNotFound,
				enum.NotFoundCode,
				"未找到目标小组或年份数据",
				err,
			))
		case errors.Is(err, service.ErrAdminControlUnlockReasonRequired):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"解锁原因不能为空",
				err,
			))
		case errors.Is(err, service.ErrAdminControlUnlockTargetTypeRequired):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"解锁目标不能为空",
				err,
			))
		case errors.Is(err, service.ErrAdminControlUnlockTargetTypeInvalid):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"解锁目标不合法",
				err,
			))
		case errors.Is(err, service.ErrAdminControlUnlockStageRequired):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"经营页回退阶段不能为空",
				err,
			))
		case errors.Is(err, service.ErrAdminControlUnlockStageInvalid):
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusUnprocessableEntity,
				enum.UnprocessableEntityCode,
				"经营页回退阶段不合法",
				err,
			))
		case errors.Is(err, service.ErrAdminControlUnlockNotAllowed):
			msg := err.Error()
			if msg == "" || msg == service.ErrAdminControlUnlockNotAllowed.Error() {
				msg = "当前目标不满足异常解锁条件"
			}
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusConflict,
				enum.ConflictCode,
				msg,
				err,
			))
		default:
			middleware.AbortWithAppError(c, middleware.NewAppError(
				http.StatusInternalServerError,
				enum.InternalServerErrorCode,
				"异常解锁失败",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}
