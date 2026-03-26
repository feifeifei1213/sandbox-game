package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/http/dto"
)

// Recovery 在 panic 时返回统一错误格式，避免裸栈暴露到接口层。
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("panic recovered", zap.Any("recovered", recovered))
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.Failure(enum.InternalServerErrorCode, "系统异常"))
	})
}
