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
	orderImportRepo := repository.NewOrderImportBatchRepository(db)
	orderBatchRepo := repository.NewOrderGenerationBatchRepository(db)
	orderConfigRepo := repository.NewOrderGenerationConfigRepository(db)
	orderForecastRepo := repository.NewOrderForecastControlRepository(db)
	orderMarketForecastRepo := repository.NewOrderMarketForecastRepository(db)
	orderMarketConfigRepo := repository.NewOrderMarketConfigRepository(db)
	orderPoolRepo := repository.NewOrderPoolRepository(db)
	marketBidRepo := repository.NewGroupMarketBidRepository(db)
	marketStateRepo := repository.NewMarketBiddingStateRepository(db)
	marketSequenceRepo := repository.NewMarketSelectionOrderRepository(db)
	groupOrderSelectionRepo := repository.NewGroupOrderSelectionRepository(db)
	stateSnapshotRepo := repository.NewStateSnapshotRepository(db)

	authService := service.NewAuthService(accountRepo, groupRepo, cfg.Auth)
	authHandler := handler.NewAuthHandler(authService)
	playerNoticeService := service.NewPlayerNoticeService(noticeRepo, adjustmentRepo, repository.NewGroupAdjustmentRevisionRepository(db))
	orderLinkService := service.NewOrderOperatingLinkService(
		marketBidRepo,
		marketStateRepo,
		groupOrderSelectionRepo,
	)
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
		orderLinkService,
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
		orderLinkService,
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
		orderLinkService,
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
		orderLinkService,
	)
	playerReportHandler := handler.NewPlayerReportHandler(
		playerReportQueryService,
		playerReportCommandService,
	)
	playerNoticeHandler := handler.NewPlayerNoticeHandler(service.NewPlayerAdjustmentSyncService(db, playerNoticeService))
	playerOrderQueryService := service.NewPlayerOrderQueryService(
		gameConfigRepo,
		groupRepo,
		groupYearRepo,
		orderBatchRepo,
		orderMarketConfigRepo,
		marketBidRepo,
		marketStateRepo,
		marketSequenceRepo,
		orderPoolRepo,
		groupOrderSelectionRepo,
	)
	playerOrderForecastQueryService := service.NewPlayerOrderForecastQueryService(
		gameConfigRepo,
		orderForecastRepo,
		orderMarketForecastRepo,
	)
	playerOrderCommandService := service.NewPlayerOrderCommandService(db)
	playerOrderHandler := handler.NewPlayerOrderHandler(
		playerOrderQueryService,
		playerOrderForecastQueryService,
		playerOrderCommandService,
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
		orderLinkService,
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
	adminDictionaryService := service.NewAdminDictionaryService(db)
	adminDictionaryHandler := handler.NewAdminDictionaryHandler(adminDictionaryService)
	adminRollbackQueryService := service.NewAdminRollbackQueryService(
		stateSnapshotRepo,
		groupRepo,
	)
	adminRollbackCommandService := service.NewAdminRollbackCommandService(db)
	adminRollbackHandler := handler.NewAdminRollbackHandler(
		adminRollbackQueryService,
		adminRollbackCommandService,
	)
	adminNoticeQueryService := service.NewAdminNoticeQueryService(
		noticeRepo,
		adjustmentRepo,
		groupRepo,
		repository.NewGroupAdjustmentRevisionRepository(db),
	)
	adminNoticeCommandService := service.NewAdminNoticeCommandService(db)
	adminNoticeHandler := handler.NewAdminNoticeHandler(
		adminNoticeQueryService,
		adminNoticeCommandService,
	)
	adminOrderQueryService := service.NewAdminOrderQueryService(
		gameConfigRepo,
		groupRepo,
		orderImportRepo,
		orderBatchRepo,
		orderConfigRepo,
		orderForecastRepo,
		orderMarketForecastRepo,
		orderMarketConfigRepo,
		marketBidRepo,
		orderPoolRepo,
		marketStateRepo,
	)
	adminOrderCommandService := service.NewAdminOrderCommandService(db)
	adminOrderControlQueryService := service.NewAdminOrderControlQueryService(
		gameConfigRepo,
		groupRepo,
		marketBidRepo,
		marketStateRepo,
		marketSequenceRepo,
		orderPoolRepo,
	)
	adminOrderControlCommandService := service.NewAdminOrderControlCommandService(db)
	adminOrderHandler := handler.NewAdminOrderHandler(
		adminOrderQueryService,
		adminOrderCommandService,
		adminOrderControlQueryService,
		adminOrderControlCommandService,
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
		gameConfig.GET("/list-game-editions", gameConfigHandler.ListGameEditions)
		gameConfig.GET("/get-year-tabs", gameConfigHandler.GetYearTabs)

		playerOperating := protected.Group("/player-operating")
		playerOperating.GET("/get-year-view", playerOperatingHandler.GetYearView)
		playerOperating.PUT("/save-draft", playerOperatingHandler.SaveDraft)
		playerOperating.POST("/submit-stage", playerOperatingHandler.SubmitStage)

		playerReport := protected.Group("/player-report")
		playerReport.GET("/get-view", playerReportHandler.GetView)
		playerReport.PUT("/save-draft", playerReportHandler.SaveDraft)
		playerReport.POST("/submit", playerReportHandler.Submit)

		playerNotice := protected.Group("/player-notice")
		playerNotice.GET("/get-adjustment-sync", playerNoticeHandler.GetAdjustmentSync)

		playerOrder := protected.Group("/player-order")
		playerOrder.GET("/get-year-view", playerOrderHandler.GetYearView)
		playerOrder.GET("/get-market-forecast", playerOrderHandler.GetMarketForecast)
		playerOrder.POST("/submit-market-investment", playerOrderHandler.SubmitMarketInvestment)
		playerOrder.POST("/submit-market-investments", playerOrderHandler.SubmitMarketInvestment)
		playerOrder.POST("/select-order", playerOrderHandler.SelectOrder)
		playerOrder.POST("/pass-segment", playerOrderHandler.PassSegment)
		playerOrder.POST("/deliver-orders", playerOrderHandler.DeliverOrders)

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

		adminDictionary := protected.Group("/admin-dictionary")
		adminDictionary.GET("/get-current", adminDictionaryHandler.GetCurrent)
		adminDictionary.GET("/list-schemes", adminDictionaryHandler.ListSchemes)
		adminDictionary.GET("/get-scheme-detail", adminDictionaryHandler.GetSchemeDetail)
		adminDictionary.POST("/save-scheme", adminDictionaryHandler.SaveScheme)
		adminDictionary.DELETE("/delete-scheme", adminDictionaryHandler.DeleteScheme)
		adminDictionary.PUT("/update-current", adminDictionaryHandler.UpdateCurrent)
		adminDictionary.POST("/apply-scheme-to-current", adminDictionaryHandler.ApplySchemeToCurrent)
		adminDictionary.POST("/restore-current-default", adminDictionaryHandler.RestoreCurrentDefault)
		adminDictionary.GET("/page-change-logs", adminDictionaryHandler.PageChangeLogs)
		adminDictionary.GET("/get-revision", adminDictionaryHandler.GetRevision)

		adminRollback := protected.Group("/admin-rollback")
		adminRollback.GET("/list-snapshots", adminRollbackHandler.ListSnapshots)
		adminRollback.GET("/get-snapshot-detail", adminRollbackHandler.GetSnapshotDetail)
		adminRollback.POST("/create-snapshot", adminRollbackHandler.CreateSnapshot)
		adminRollback.POST("/restore-group-snapshot", adminRollbackHandler.RestoreGroupSnapshot)

		adminNotice := protected.Group("/admin-notice")
		adminNotice.GET("/get-records", adminNoticeHandler.GetRecords)
		adminNotice.POST("/send-general", adminNoticeHandler.SendGeneral)
		adminNotice.POST("/send-adjustment", adminNoticeHandler.SendAdjustment)
		adminNotice.POST("/preview-adjustment", adminNoticeHandler.PreviewAdjustment)
		adminNotice.POST("/void-adjustment", adminNoticeHandler.VoidAdjustment)

		adminOrder := protected.Group("/admin-order")
		adminOrder.GET("/get-forecast-control", adminOrderHandler.GetForecastControl)
		adminOrder.PUT("/update-forecast-control", adminOrderHandler.UpdateForecastControl)
		adminOrder.GET("/get-control-config", adminOrderHandler.GetControlConfig)
		adminOrder.POST("/upload-excel", adminOrderHandler.UploadExcel)
		adminOrder.PUT("/update-market-enabled-config", adminOrderHandler.UpdateMarketConfig)
		adminOrder.PUT("/update-control-config", adminOrderHandler.UpdateControlConfig)
		adminOrder.POST("/generate-order-pool", adminOrderHandler.GenerateOrderPool)
		adminOrder.POST("/generate-preview-pool", adminOrderHandler.GenerateOrderPool)
		adminOrder.POST("/confirm-order-pool", adminOrderHandler.ConfirmOrderPool)
		adminOrder.POST("/generate-selection-sequence", adminOrderHandler.GenerateSelectionSequence)
		adminOrder.GET("/get-order-pool", adminOrderHandler.GetOrderPool)
		adminOrder.POST("/open-market-bidding", adminOrderHandler.OpenMarketBidding)
		adminOrder.POST("/close-market-bidding", adminOrderHandler.CloseMarketBidding)
		adminOrder.GET("/get-market-selection-status", adminOrderHandler.GetMarketSelectionStatus)
		adminOrder.POST("/release-next-segment", adminOrderHandler.ReleaseNextSegment)
		adminOrder.POST("/open-next-round", adminOrderHandler.OpenNextRound)
		adminOrder.POST("/admin-skip-current-group", adminOrderHandler.AdminSkipCurrentGroup)
	}

	registerFrontendStaticRoutes(engine, cfg, logger)

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
