package app

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"sandbox-game/internal/assembler"
	appconfig "sandbox-game/internal/config"
	"sandbox-game/internal/http/handler"
	"sandbox-game/internal/http/middleware"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/service"
)

// NewRouter 注册当前阶段最小可运行的基础路由。
func NewRouter(cfg *appconfig.Config, logger *zap.Logger, db *gorm.DB) *gin.Engine {
	gin.SetMode(resolveGinMode(cfg.Server.Mode))

	engine := gin.New()
	engine.Use(middleware.RequestLogger(logger))
	engine.Use(middleware.Recovery(logger))
	engine.Use(middleware.ErrorHandler())

	logger.Info("router initialized",
		zap.String("mode", gin.Mode()),
		zap.String("appName", cfg.App.Name),
	)

	healthHandler := handler.NewHealthHandler(cfg)
	playerOperatingAssembler := assembler.NewPlayerOperatingAssembler()
	playerReportAssembler := assembler.NewPlayerReportAssembler()
	gameConfigRepo := repository.NewGameConfigRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	groupYearRepo := repository.NewGroupYearStateRepository(db)
	operatingRepo := repository.NewOperatingRepository(db)
	initialBaseRepo := repository.NewInitialBaselineRepository(db)
	reportRepo := repository.NewReportRepository(db)
	summaryRepo := repository.NewSummarySnapshotRepository(db)

	playerOperatingQueryService := service.NewPlayerOperatingQueryService(
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaseRepo,
		reportRepo,
		playerOperatingAssembler,
	)
	playerOperatingCommandService := service.NewPlayerOperatingCommandService(
		db,
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaseRepo,
		reportRepo,
	)
	playerOperatingHandler := handler.NewPlayerOperatingHandler(
		playerOperatingQueryService,
		playerOperatingCommandService,
	)
	playerReportQueryService := service.NewPlayerReportQueryService(
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaseRepo,
		reportRepo,
		playerReportAssembler,
	)
	playerReportCommandService := service.NewPlayerReportCommandService(
		db,
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaseRepo,
		reportRepo,
	)
	playerReportHandler := handler.NewPlayerReportHandler(
		playerReportQueryService,
		playerReportCommandService,
	)
	adminSummaryQueryService := service.NewAdminSummaryQueryService(
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		summaryRepo,
	)
	adminSummaryHandler := handler.NewAdminSummaryHandler(adminSummaryQueryService)

	engine.GET("/healthz", healthHandler.GetHealth)

	apiV1 := engine.Group("/api/v1/sandbox-game")
	apiV1.Use(middleware.AuthBypass())
	{
		playerOperating := apiV1.Group("/player-operating")
		playerOperating.GET("/get-year-view", playerOperatingHandler.GetYearView)
		playerOperating.PUT("/save-draft", playerOperatingHandler.SaveDraft)
		playerOperating.POST("/submit-stage", playerOperatingHandler.SubmitStage)

		playerReport := apiV1.Group("/player-report")
		playerReport.GET("/get-view", playerReportHandler.GetView)
		playerReport.PUT("/save-draft", playerReportHandler.SaveDraft)
		playerReport.POST("/submit", playerReportHandler.Submit)

		adminSummary := apiV1.Group("/admin-summary")
		adminSummary.GET("/get-year-summary", adminSummaryHandler.GetYearSummary)
		adminSummary.GET("/get-final-ranking", adminSummaryHandler.GetFinalRanking)
	}

	return engine
}

func resolveGinMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case gin.ReleaseMode:
		return gin.ReleaseMode
	case gin.TestMode:
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}
