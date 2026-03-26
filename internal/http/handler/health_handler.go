package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appconfig "sandbox-game/internal/config"
	"sandbox-game/internal/http/dto"
)

// HealthHandler 提供最小健康检查接口。
type HealthHandler struct {
	config *appconfig.Config
}

func NewHealthHandler(config *appconfig.Config) *HealthHandler {
	return &HealthHandler{config: config}
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, dto.Success(gin.H{
		"name":    h.config.App.Name,
		"version": h.config.App.Version,
		"status":  "ok",
	}))
}
