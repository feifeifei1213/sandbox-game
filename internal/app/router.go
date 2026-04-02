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
	accountRepo := repository.NewAccountRepository(db)
	adminActionLogRepo := repository.NewAdminActionLogRepository(db)
	noticeRepo := repository.NewNoticeRepository(db)
	adjustmentRepo := repository.NewGroupAdjustmentRepository(db)

	authService := service.NewAuthService(accountRepo, groupRepo, cfg.Auth)
	authHandler := handler.NewAuthHandler(authService)
	playerNoticeService := service.NewPlayerNoticeService(noticeRepo, adjustmentRepo)
	gameConfigQueryService := service.NewGameConfigQueryService(gameConfigRepo, groupRepo, groupYearRepo)
	gameConfigHandler := handler.NewGameConfigHandler(gameConfigQueryService)
	playerOperatingQueryService := service.NewPlayerOperatingQueryService(
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaseRepo,
		reportRepo,
		playerOperatingAssembler,
		playerNoticeService,
	)
	playerOperatingCommandService := service.NewPlayerOperatingCommandService(
		db,
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaseRepo,
		reportRepo,
		playerNoticeService,
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
		playerNoticeService,
	)
	playerReportCommandService := service.NewPlayerReportCommandService(
		db,
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaseRepo,
		reportRepo,
		playerNoticeService,
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
	adminGroupDataQueryService := service.NewAdminGroupDataQueryService(
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		operatingRepo,
		initialBaseRepo,
		reportRepo,
		playerOperatingAssembler,
		playerReportAssembler,
		playerNoticeService,
	)
	adminGroupDataHandler := handler.NewAdminGroupDataHandler(adminGroupDataQueryService)
	adminControlQueryService := service.NewAdminControlQueryService(
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		initialBaseRepo,
		accountRepo,
		adminActionLogRepo,
	)
	adminControlCommandService := service.NewAdminControlCommandService(db)
	adminControlHandler := handler.NewAdminControlHandler(
		adminControlQueryService,
		adminControlCommandService,
	)
	adminNoticeQueryService := service.NewAdminNoticeQueryService(
		noticeRepo,
		adjustmentRepo,
		groupRepo,
	)
	adminNoticeCommandService := service.NewAdminNoticeCommandService(db)
	adminNoticeHandler := handler.NewAdminNoticeHandler(
		adminNoticeQueryService,
		adminNoticeCommandService,
	)

	engine.GET("/healthz", healthHandler.GetHealth)

	apiV1 := engine.Group("/api/v1/sandbox-game")
	publicAuth := apiV1.Group("/auth")
	publicAuth.POST("/login", authHandler.Login)

	protected := apiV1.Group("")
	protected.Use(resolveSandboxAuthMiddleware(cfg, authService))
	{
		authGroup := protected.Group("/auth")
		authGroup.GET("/get-current-user", authHandler.GetCurrentUser)
		authGroup.POST("/logout", authHandler.Logout)

		gameConfig := protected.Group("/game-config")
		gameConfig.GET("/get-current", gameConfigHandler.GetCurrent)
		gameConfig.GET("/get-year-tabs", gameConfigHandler.GetYearTabs)

		playerOperating := protected.Group("/player-operating")
		playerOperating.GET("/get-year-view", playerOperatingHandler.GetYearView)
		playerOperating.PUT("/save-draft", playerOperatingHandler.SaveDraft)
		playerOperating.POST("/submit-stage", playerOperatingHandler.SubmitStage)

		playerReport := protected.Group("/player-report")
		playerReport.GET("/get-view", playerReportHandler.GetView)
		playerReport.PUT("/save-draft", playerReportHandler.SaveDraft)
		playerReport.POST("/submit", playerReportHandler.Submit)

		adminSummary := protected.Group("/admin-summary")
		adminSummary.GET("/get-year-summary", adminSummaryHandler.GetYearSummary)
		adminSummary.GET("/get-final-ranking", adminSummaryHandler.GetFinalRanking)

		adminGroupData := protected.Group("/admin-group-data")
		adminGroupData.GET("/list-groups", adminGroupDataHandler.ListGroups)
		adminGroupData.GET("/get-operating-view", adminGroupDataHandler.GetOperatingView)
		adminGroupData.GET("/get-report-view", adminGroupDataHandler.GetReportView)

		adminControl := protected.Group("/admin-control")
		adminControl.GET("/get-setup-status", adminControlHandler.GetSetupStatus)
		adminControl.POST("/initialize-game", adminControlHandler.InitializeGame)
		adminControl.GET("/get-config", adminControlHandler.GetConfig)
		adminControl.PUT("/update-final-year", adminControlHandler.UpdateFinalYear)
		adminControl.POST("/open-next-year", adminControlHandler.OpenNextYear)
		adminControl.GET("/get-initial-baseline", adminControlHandler.GetInitialBaseline)
		adminControl.POST("/submit-initial-baseline", adminControlHandler.SubmitInitialBaseline)
		adminControl.POST("/unlock-year", adminControlHandler.UnlockYear)

		adminNotice := protected.Group("/admin-notice")
		adminNotice.GET("/get-records", adminNoticeHandler.GetRecords)
		adminNotice.POST("/send-general", adminNoticeHandler.SendGeneral)
		adminNotice.POST("/send-adjustment", adminNoticeHandler.SendAdjustment)
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

func resolveSandboxAuthMiddleware(cfg *appconfig.Config, authService *service.AuthService) gin.HandlerFunc {
	if strings.EqualFold(strings.TrimSpace(cfg.Auth.Mode), "bypass") {
		return middleware.AuthBypass()
	}
	return middleware.RequireAuth(authService)
}
