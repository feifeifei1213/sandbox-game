package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	appconfig "sandbox-game/internal/config"
	appdb "sandbox-game/internal/infra/db"
	applogger "sandbox-game/internal/infra/logger"
)

// Application 表示首版 Go 单体应用的最小运行时容器。
type Application struct {
	Config *appconfig.Config
	Logger *zap.Logger
	Router *gin.Engine
	DB     *gorm.DB
}

// Bootstrap 负责完成配置、日志和路由的最小装配。
func Bootstrap(configPath string) (*Application, error) {
	cfg, err := appconfig.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	logger, err := applogger.Build(cfg.Log.Level)
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	db, err := appdb.NewMySQL(cfg)
	if err != nil {
		return nil, fmt.Errorf("init mysql: %w", err)
	}

	router := NewRouter(cfg, logger, db)

	return &Application{
		Config: cfg,
		Logger: logger,
		Router: router,
		DB:     db,
	}, nil
}
