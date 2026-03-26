package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/http/dto"
)

// AppError 表示业务或基础设施层抛出的结构化错误。
type AppError struct {
	HTTPStatus int
	Code       int
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// NewAppError 创建结构化错误。
func NewAppError(httpStatus int, code int, message string, err error) *AppError {
	if code == 0 {
		code = httpStatus
	}
	return &AppError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// AbortWithAppError 用于在 Handler/Service 层中断当前请求。
func AbortWithAppError(c *gin.Context, appErr *AppError) {
	_ = c.Error(appErr)
	c.Abort()
}

// ErrorHandler 统一将错误转换为标准响应体。
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		lastErr := c.Errors.Last().Err
		var appErr *AppError
		if errors.As(lastErr, &appErr) {
			c.JSON(appErr.HTTPStatus, dto.Failure(appErr.Code, appErr.Message))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.Failure(enum.InternalServerErrorCode, "系统异常"))
	}
}
