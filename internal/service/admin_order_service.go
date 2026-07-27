package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

const (
	adminActionCodeUploadOrderExcel      = "UPLOAD_ORDER_EXCEL"
	adminActionCodeUpdateForecastControl = "UPDATE_ORDER_FORECAST_CONTROL"
	adminActionCodeUpdateOrderConfig     = "UPDATE_ORDER_CONFIG"
	adminActionCodeUpdateMarketConfig    = "UPDATE_ORDER_MARKET_CONFIG"
	adminActionCodeGenerateOrderPool     = "GENERATE_ORDER_POOL"
	adminActionCodeGenerateOrderPreview  = "GENERATE_ORDER_PREVIEW"
	adminActionCodeConfirmOrderPool      = "CONFIRM_ORDER_POOL"
	defaultOrderReleaseSequenceStart     = 1
	defaultOrderParsePreviewLimit        = 20
	requiredOrderInvestmentSegmentCount  = 16
	forecastControlMinYear               = 1
	forecastControlMaxYear               = 8
)

const (
	orderYearLockReasonHistoricalYear   = "该年份已成为历史年份"
	orderYearLockReasonRollbackPending  = "该年份存在回退补提，订单事实保持锁定"
	orderYearLockReasonPoolConfirmed    = "订单池已确认"
	orderYearLockReasonSelectedOrders   = "该年份已产生订单历史"
	orderYearLockReasonGroupSelections  = "该年份已有小组选单记录"
	orderYearLockReasonMarketBids       = "该年份已有市场投入"
	orderYearLockReasonSelectionOrders  = "该年份已生成选单顺序"
	orderYearLockReasonWorkflowAdvanced = "该年份竞标流程已开始"
)

var (
	ErrAdminOrderYearInvalid                  = errors.New("admin order year invalid")
	ErrAdminOrderFileRequired                 = errors.New("admin order file required")
	ErrAdminOrderParseFailed                  = errors.New("admin order parse failed")
	ErrAdminOrderConfigInvalid                = errors.New("admin order config invalid")
	ErrAdminOrderReleaseSequenceDuplicated    = errors.New("admin order release sequence duplicated")
	ErrAdminOrderReleaseSequenceLocked        = errors.New("admin order release sequence locked")
	ErrAdminOrderPoolLocked                   = errors.New("admin order pool locked")
	ErrAdminOrderConfigNotFound               = errors.New("admin order config not found")
	ErrAdminOrderBatchNotFound                = errors.New("admin order batch not found")
	ErrAdminOrderSourceInsufficient           = errors.New("admin order source insufficient")
	ErrAdminOrderPreviewNotFound              = errors.New("admin order preview not found")
	ErrAdminOrderPoolNotConfirmed             = errors.New("admin order pool not confirmed")
	ErrAdminOrderInvestmentIncomplete         = errors.New("admin order investment incomplete")
	ErrAdminOrderSequenceAlreadyGenerated     = errors.New("admin order sequence already generated")
	ErrAdminOrderMarketConfigLocked           = errors.New("admin order market config locked")
	ErrAdminOrderForecastControlInvalid       = errors.New("admin order forecast control invalid")
	ErrAdminOrderMarketInvestmentLimitInvalid = errors.New("admin order market investment limit invalid")
)

type orderYearLockInfo struct {
	YearNo int
	Locked bool
	Reason string
}

type OrderSegmentDefinition struct {
	MarketCode    string
	MarketName    string
	OrderType     string
	OrderTypeName string
}

type UploadOrderExcelCommand struct {
	FileName     string
	FileSize     int64
	Content      []byte
	OperatorID   int64
	OperatorName string
}

type UploadOrderExcelResult struct {
	BatchID          int64                `json:"batchId"`
	OriginalFileName string               `json:"originalFileName"`
	ParseStatus      string               `json:"parseStatus"`
	ParsedOrderCount int                  `json:"parsedOrderCount"`
	Warnings         []string             `json:"warnings"`
	Summary          []OrderSourceSummary `json:"summary"`
	PreviewOrders    []ParsedOrderCard    `json:"previewOrders"`
	UploadedAt       string               `json:"uploadedAt"`
}

type OrderSourceSummary struct {
	YearNo         int    `json:"yearNo"`
	MarketCode     string `json:"marketCode"`
	MarketName     string `json:"marketName"`
	OrderType      string `json:"orderType"`
	OrderTypeName  string `json:"orderTypeName"`
	AvailableCount int    `json:"availableCount"`
}

type OrderForecastControlItem struct {
	YearNo            int    `json:"yearNo"`
	ForecastStageCode string `json:"forecastStageCode"`
	ForecastStageName string `json:"forecastStageName"`
	MarketCode        string `json:"marketCode"`
	MarketName        string `json:"marketName"`
	OrderType         string `json:"orderType"`
	OrderTypeName     string `json:"orderTypeName"`
	OrderCount        int    `json:"orderCount"`
}

type OrderForecastYearLock struct {
	YearNo int    `json:"yearNo"`
	Locked bool   `json:"locked"`
	Reason string `json:"reason,omitempty"`
}

type OrderForecastNarrativeItem struct {
	ForecastStageCode string `json:"forecastStageCode"`
	ForecastStageName string `json:"forecastStageName"`
	MarketCode        string `json:"marketCode"`
	MarketName        string `json:"marketName"`
	Content           string `json:"content"`
}

type OrderMarketForecastProduct struct {
	OrderType      string  `json:"orderType"`
	OrderTypeName  string  `json:"orderTypeName"`
	OrderCount     int     `json:"orderCount"`
	ForecastAmount float64 `json:"forecastAmount"`
}

type OrderMarketForecastYear struct {
	YearNo              int                          `json:"yearNo"`
	Products            []OrderMarketForecastProduct `json:"products"`
	TotalOrderCount     int                          `json:"totalOrderCount"`
	TotalForecastAmount float64                      `json:"totalForecastAmount"`
}

type OrderMarketForecastMarket struct {
	MarketCode string                    `json:"marketCode"`
	MarketName string                    `json:"marketName"`
	Narrative  string                    `json:"narrative"`
	Years      []OrderMarketForecastYear `json:"years"`
}

type OrderMarketForecastStage struct {
	ForecastStageCode string                      `json:"forecastStageCode"`
	ForecastStageName string                      `json:"forecastStageName"`
	YearRange         string                      `json:"yearRange"`
	Years             []int                       `json:"years"`
	Markets           []OrderMarketForecastMarket `json:"markets"`
}

type OrderMarketForecastResult struct {
	FormulaVersion string                     `json:"formulaVersion"`
	OrderTemplate  OrderTemplateView          `json:"orderTemplate"`
	Stages         []OrderMarketForecastStage `json:"stages"`
}

type OrderForecastControlResult struct {
	OrderTemplate OrderTemplateView            `json:"orderTemplate"`
	Items         []OrderForecastControlItem   `json:"items"`
	Narratives    []OrderForecastNarrativeItem `json:"narratives"`
	Forecast      OrderMarketForecastResult    `json:"forecast"`
	YearLocks     []OrderForecastYearLock      `json:"yearLocks"`
}

type OrderControlConfigItem struct {
	YearNo            int    `json:"yearNo"`
	MarketCode        string `json:"marketCode"`
	MarketName        string `json:"marketName"`
	MarketEnabled     bool   `json:"marketEnabled"`
	OrderType         string `json:"orderType"`
	OrderTypeName     string `json:"orderTypeName"`
	OrderCount        int    `json:"orderCount"`
	ReleaseSequenceNo int    `json:"releaseSequenceNo"`
	AvailableCount    int    `json:"availableCount"`
	ConfigStatus      string `json:"configStatus"`
	GeneratedCount    int    `json:"generatedCount"`
}

type OrderMarketConfigItem struct {
	YearNo                int      `json:"yearNo"`
	MarketCode            string   `json:"marketCode"`
	MarketName            string   `json:"marketName"`
	Enabled               bool     `json:"enabled"`
	MarketInvestmentLimit *float64 `json:"marketInvestmentLimit"`
	ConfigStatus          string   `json:"configStatus"`
	LockedBatchID         *int64   `json:"lockedBatchId,omitempty"`
}

type OrderGenerationBatchSummary struct {
	BatchID        int64   `json:"batchId"`
	BatchStatus    string  `json:"batchStatus"`
	FormulaVersion string  `json:"formulaVersion"`
	RandomSeed     string  `json:"randomSeed,omitempty"`
	GeneratedCount int     `json:"generatedCount"`
	GeneratedAt    string  `json:"generatedAt"`
	GeneratedBy    string  `json:"generatedBy"`
	ConfirmedAt    *string `json:"confirmedAt,omitempty"`
	ConfirmedBy    *string `json:"confirmedBy,omitempty"`
}

type OrderControlWarning struct {
	Level      string `json:"level"`
	Message    string `json:"message"`
	MarketCode string `json:"marketCode,omitempty"`
	OrderType  string `json:"orderType,omitempty"`
}

type OrderControlConfigResult struct {
	YearNo                int                          `json:"yearNo"`
	FinalYear             int                          `json:"finalYear"`
	OrderTemplate         OrderTemplateView            `json:"orderTemplate"`
	LatestBatchID         *int64                       `json:"latestBatchId"`
	LatestBatchUploadedAt *string                      `json:"latestBatchUploadedAt"`
	Forecast              OrderMarketForecastResult    `json:"forecast"`
	GenerationStatus      string                       `json:"generationStatus"`
	LatestPreviewBatch    *OrderGenerationBatchSummary `json:"latestPreviewBatch"`
	ConfirmedBatch        *OrderGenerationBatchSummary `json:"confirmedBatch"`
	CanUpdateConfig       bool                         `json:"canUpdateConfig"`
	CanGeneratePreview    bool                         `json:"canGeneratePreview"`
	CanConfirmPool        bool                         `json:"canConfirmPool"`
	MarketConfigs         []OrderMarketConfigItem      `json:"marketConfigs"`
	Items                 []OrderControlConfigItem     `json:"items"`
	Warnings              []OrderControlWarning        `json:"warnings"`
}

type UpdateOrderMarketConfigCommand struct {
	YearNo       int
	Markets      []UpdateOrderMarketConfigItem
	OperatorID   int64
	OperatorName string
}

type UpdateOrderMarketConfigItem struct {
	MarketCode            string
	Enabled               bool
	MarketInvestmentLimit *float64
}

type UpdateOrderMarketConfigResult struct {
	YearNo    int                      `json:"yearNo"`
	Markets   []OrderMarketConfigItem  `json:"markets"`
	Items     []OrderControlConfigItem `json:"items"`
	Warnings  []OrderControlWarning    `json:"warnings"`
	UpdatedAt string                   `json:"updatedAt"`
	UpdatedBy string                   `json:"updatedBy"`
}

type UpdateOrderForecastControlCommand struct {
	Items        []UpdateOrderForecastControlItem
	Narratives   []UpdateOrderForecastNarrativeItem
	OperatorID   int64
	OperatorName string
}

type UpdateOrderForecastControlItem struct {
	YearNo     int
	MarketCode string
	OrderType  string
	OrderCount int
}

type UpdateOrderForecastNarrativeItem struct {
	ForecastStageCode string
	MarketCode        string
	Content           string
}

type UpdateOrderForecastControlResult struct {
	OrderTemplate    OrderTemplateView            `json:"orderTemplate"`
	Items            []OrderForecastControlItem   `json:"items"`
	Narratives       []OrderForecastNarrativeItem `json:"narratives"`
	Forecast         OrderMarketForecastResult    `json:"forecast"`
	YearLocks        []OrderForecastYearLock      `json:"yearLocks"`
	AutoPreviewYears []int                        `json:"autoPreviewYears"`
	UpdatedAt        string                       `json:"updatedAt"`
	UpdatedBy        string                       `json:"updatedBy"`
}

type UpdateOrderControlConfigCommand struct {
	YearNo       int
	Items        []UpdateOrderControlConfigItem
	OperatorID   int64
	OperatorName string
}

type UpdateOrderControlConfigItem struct {
	MarketCode        string
	OrderType         string
	OrderCount        int
	ReleaseSequenceNo int
}

type UpdateOrderControlConfigResult struct {
	YearNo    int                      `json:"yearNo"`
	Items     []OrderControlConfigItem `json:"items"`
	Warnings  []OrderControlWarning    `json:"warnings"`
	UpdatedAt string                   `json:"updatedAt"`
	UpdatedBy string                   `json:"updatedBy"`
}

type GenerateOrderPoolCommand struct {
	YearNo        int
	Overwrite     bool
	SourceBatchID *int64
	OperatorID    int64
	OperatorName  string
}

type GenerateOrderPoolResult struct {
	YearNo         int                   `json:"yearNo"`
	SourceBatchID  int64                 `json:"sourceBatchId,omitempty"`
	BatchID        int64                 `json:"batchId"`
	BatchStatus    string                `json:"batchStatus"`
	RandomSeed     string                `json:"randomSeed"`
	FormulaVersion string                `json:"formulaVersion"`
	GeneratedCount int                   `json:"generatedCount"`
	SegmentCount   int                   `json:"segmentCount"`
	Warnings       []OrderControlWarning `json:"warnings"`
	GeneratedAt    string                `json:"generatedAt"`
	GeneratedBy    string                `json:"generatedBy"`
}

type ConfirmOrderPoolCommand struct {
	YearNo       int
	BatchID      int64
	OperatorID   int64
	OperatorName string
}

type ConfirmOrderPoolResult struct {
	YearNo         int    `json:"yearNo"`
	BatchID        int64  `json:"batchId"`
	GeneratedCount int    `json:"generatedCount"`
	SegmentCount   int    `json:"segmentCount"`
	ConfirmedAt    string `json:"confirmedAt"`
	ConfirmedBy    string `json:"confirmedBy"`
}

type OrderPoolItem struct {
	OrderID         int64          `json:"orderId"`
	BusinessOrderNo string         `json:"businessOrderNo"`
	CardSequenceNo  int            `json:"cardSequenceNo"`
	YearNo          int            `json:"yearNo"`
	MarketCode      string         `json:"marketCode"`
	MarketName      string         `json:"marketName"`
	OrderType       string         `json:"orderType"`
	OrderTypeName   string         `json:"orderTypeName"`
	OrderAmount     float64        `json:"orderAmount"`
	OrderQuantity   float64        `json:"orderQuantity"`
	UnitPrice       float64        `json:"unitPrice"`
	AccountTerm     int            `json:"accountTerm"`
	PoolStatus      string         `json:"poolStatus"`
	SelectedGroupID *int64         `json:"selectedGroupId"`
	SourceSheetName string         `json:"sourceSheetName"`
	SourceCell      string         `json:"sourceCell"`
	OrderPayload    map[string]any `json:"orderPayload,omitempty"`
}

type OrderPoolResult struct {
	YearNo        int               `json:"yearNo"`
	OrderTemplate OrderTemplateView `json:"orderTemplate"`
	MarketCode    string            `json:"marketCode"`
	MarketName    string            `json:"marketName"`
	OrderType     string            `json:"orderType"`
	OrderTypeName string            `json:"orderTypeName"`
	List          []OrderPoolItem   `json:"list"`
}

type AdminOrderQueryService struct {
	gameConfigRepo     *repository.GameConfigRepository
	groupRepo          *repository.GroupRepository
	importRepo         *repository.OrderImportBatchRepository
	batchRepo          *repository.OrderGenerationBatchRepository
	configRepo         *repository.OrderGenerationConfigRepository
	forecastRepo       *repository.OrderForecastControlRepository
	marketForecastRepo *repository.OrderMarketForecastRepository
	marketRepo         *repository.OrderMarketConfigRepository
	bidRepo            *repository.GroupMarketBidRepository
	poolRepo           *repository.OrderPoolRepository
	stateRepo          *repository.MarketBiddingStateRepository
	groupYearRepo      *repository.GroupYearStateRepository
	selectionRepo      *repository.GroupOrderSelectionRepository
	sequenceRepo       *repository.MarketSelectionOrderRepository
}

func NewAdminOrderQueryService(
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	importRepo *repository.OrderImportBatchRepository,
	batchRepo *repository.OrderGenerationBatchRepository,
	configRepo *repository.OrderGenerationConfigRepository,
	forecastRepo *repository.OrderForecastControlRepository,
	marketForecastRepo *repository.OrderMarketForecastRepository,
	marketRepo *repository.OrderMarketConfigRepository,
	bidRepo *repository.GroupMarketBidRepository,
	poolRepo *repository.OrderPoolRepository,
	stateRepo *repository.MarketBiddingStateRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	selectionRepo *repository.GroupOrderSelectionRepository,
	sequenceRepo *repository.MarketSelectionOrderRepository,
) *AdminOrderQueryService {
	return &AdminOrderQueryService{
		gameConfigRepo:     gameConfigRepo,
		groupRepo:          groupRepo,
		importRepo:         importRepo,
		batchRepo:          batchRepo,
		configRepo:         configRepo,
		forecastRepo:       forecastRepo,
		marketForecastRepo: marketForecastRepo,
		marketRepo:         marketRepo,
		bidRepo:            bidRepo,
		poolRepo:           poolRepo,
		stateRepo:          stateRepo,
		groupYearRepo:      groupYearRepo,
		selectionRepo:      selectionRepo,
		sequenceRepo:       sequenceRepo,
	}
}

func (s *AdminOrderQueryService) GetForecastControl(ctx context.Context) (*OrderForecastControlResult, error) {
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	template, err := ResolveOrderTemplateByGameConfig(*gameConfig)
	if err != nil {
		return nil, err
	}
	items, err := s.forecastRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list order forecast control: %w", err)
	}
	items = filterForecastControlsByTemplate(items, template)
	forecastRows, err := s.marketForecastRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list market forecast: %w", err)
	}
	forecastRows = filterMarketForecastsByTemplate(forecastRows, template)
	if len(items) == 0 {
		items = defaultForecastControlEntitiesForTemplate("system", time.Now(), template)
	}
	yearLocks, err := s.buildForecastYearLocks(ctx)
	if err != nil {
		return nil, err
	}
	forecast := buildMarketForecastResultForTemplate(items, forecastRows, template)
	return &OrderForecastControlResult{
		OrderTemplate: BuildOrderTemplateView(template),
		Items:         buildOrderForecastControlItemsForTemplate(items, template),
		Narratives:    buildOrderForecastNarrativesForTemplate(forecastRows, template),
		Forecast:      forecast,
		YearLocks:     yearLocks,
	}, nil
}

func (s *AdminOrderQueryService) GetControlConfig(ctx context.Context, yearNo int) (*OrderControlConfigResult, error) {
	if err := s.validateYear(ctx, yearNo); err != nil {
		return nil, err
	}
	configs, err := s.configRepo.ListByYear(ctx, yearNo)
	if err != nil {
		return nil, fmt.Errorf("list order configs: %w", err)
	}
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	template, err := ResolveOrderTemplateByGameConfig(*gameConfig)
	if err != nil {
		return nil, err
	}
	configs = filterGenerationConfigsByTemplate(configs, template)
	marketConfigs, err := s.marketRepo.ListByYear(ctx, yearNo)
	if err != nil {
		return nil, fmt.Errorf("list order market configs: %w", err)
	}
	marketConfigs = filterMarketConfigsByTemplate(marketConfigs, template)
	marketConfigMap := buildMarketEnabledMapForTemplate(marketConfigs, template)
	var latestBatchID *int64
	var latestBatchUploadedAt *string
	sourceCounts := defaultOrderSourceCountsForTemplate(yearNo, template)
	forecastControls, err := s.forecastRepo.ListByYear(ctx, yearNo)
	if err != nil {
		return nil, fmt.Errorf("list order forecast control: %w", err)
	}
	forecastControls = filterForecastControlsByTemplate(forecastControls, template)
	forecastControlCounts := forecastControlCountMapForTemplate(yearNo, forecastControls, template)
	allForecastControls, err := s.forecastRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list all order forecast control: %w", err)
	}
	allForecastControls = filterForecastControlsByTemplate(allForecastControls, template)
	forecastRows, err := s.marketForecastRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list order market forecast: %w", err)
	}
	forecastRows = filterMarketForecastsByTemplate(forecastRows, template)
	latestBatch, batchErr := s.importRepo.FindLatestSuccess(ctx)
	if batchErr == nil {
		latestBatchID = &latestBatch.ID
		uploadedAt := latestBatch.UploadedAt.Format(time.RFC3339)
		latestBatchUploadedAt = &uploadedAt
	} else if !repository.IsRecordNotFound(batchErr) {
		return nil, fmt.Errorf("load latest order batch: %w", batchErr)
	}
	poolCounts, err := s.countGeneratedByYear(ctx, yearNo)
	if err != nil {
		return nil, err
	}
	latestPreview, previewErr := s.batchRepo.FindLatestByYearStatuses(ctx, yearNo, []string{enum.OrderGenerationBatchStatusPreview})
	if previewErr != nil && !repository.IsRecordNotFound(previewErr) {
		return nil, fmt.Errorf("load latest preview batch: %w", previewErr)
	}
	confirmed, confirmedErr := s.batchRepo.FindConfirmedByYear(ctx, yearNo)
	if confirmedErr != nil && !repository.IsRecordNotFound(confirmedErr) {
		return nil, fmt.Errorf("load confirmed order batch: %w", confirmedErr)
	}
	states, err := s.stateRepo.ListByYear(ctx, yearNo)
	if err != nil {
		return nil, fmt.Errorf("list order states: %w", err)
	}
	lockInfo, err := resolveOrderYearLock(ctx, gameConfig, s.batchRepo, s.poolRepo, s.bidRepo, s.stateRepo, s.groupYearRepo, s.selectionRepo, s.sequenceRepo, yearNo)
	if err != nil {
		return nil, err
	}

	items := buildControlConfigItemsForTemplate(yearNo, configs, forecastControlCounts, sourceCounts, poolCounts, marketConfigMap, template)
	groupCount, err := s.groupRepo.CountAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("count groups: %w", err)
	}
	return &OrderControlConfigResult{
		YearNo:                yearNo,
		FinalYear:             gameConfig.FinalYear,
		OrderTemplate:         BuildOrderTemplateView(template),
		LatestBatchID:         latestBatchID,
		LatestBatchUploadedAt: latestBatchUploadedAt,
		Forecast:              buildMarketForecastResultForTemplate(mergeForecastControlDefaultsForTemplate(allForecastControls, template), forecastRows, template),
		GenerationStatus:      resolveOrderGenerationStatus(latestPreview, confirmed, states),
		LatestPreviewBatch:    buildOrderGenerationBatchSummary(latestPreview),
		ConfirmedBatch:        buildOrderGenerationBatchSummary(confirmed),
		CanUpdateConfig:       !lockInfo.Locked,
		CanGeneratePreview:    !lockInfo.Locked,
		CanConfirmPool:        latestPreview != nil && !lockInfo.Locked,
		MarketConfigs:         buildOrderMarketConfigItemsForTemplate(yearNo, marketConfigs, template),
		Items:                 items,
		Warnings:              buildOrderControlWarnings(items, int(groupCount)),
	}, nil
}

func (s *AdminOrderQueryService) GetOrderPool(ctx context.Context, yearNo int, marketCode string, orderType string) (*OrderPoolResult, error) {
	if err := s.validateYear(ctx, yearNo); err != nil {
		return nil, err
	}
	marketCode = strings.ToUpper(strings.TrimSpace(marketCode))
	orderType = strings.ToUpper(strings.TrimSpace(orderType))
	if marketCode == "ALL" {
		marketCode = ""
	}
	if orderType == "ALL" {
		orderType = ""
	}
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	template, err := ResolveOrderTemplateByGameConfig(*gameConfig)
	if err != nil {
		return nil, err
	}
	if marketCode != "" && !template.IsValidMarketCode(marketCode) {
		return nil, ErrAdminOrderConfigInvalid
	}
	if orderType != "" && !template.IsValidOrderType(orderType) {
		return nil, ErrAdminOrderConfigInvalid
	}
	items, err := s.poolRepo.ListByYearWithOptionalFilters(ctx, yearNo, marketCode, orderType)
	if err != nil {
		return nil, fmt.Errorf("list order pool: %w", err)
	}
	resultItems := make([]OrderPoolItem, 0, len(items))
	for _, item := range items {
		resultItems = append(resultItems, buildOrderPoolItem(item))
	}
	return &OrderPoolResult{
		YearNo:        yearNo,
		OrderTemplate: BuildOrderTemplateView(template),
		MarketCode:    marketCode,
		MarketName:    template.MarketName(marketCode),
		OrderType:     orderType,
		OrderTypeName: template.OrderTypeName(orderType),
		List:          resultItems,
	}, nil
}

type PlayerOrderForecastQueryService struct {
	gameConfigRepo     *repository.GameConfigRepository
	forecastRepo       *repository.OrderForecastControlRepository
	marketForecastRepo *repository.OrderMarketForecastRepository
}

func NewPlayerOrderForecastQueryService(
	gameConfigRepo *repository.GameConfigRepository,
	forecastRepo *repository.OrderForecastControlRepository,
	marketForecastRepo *repository.OrderMarketForecastRepository,
) *PlayerOrderForecastQueryService {
	return &PlayerOrderForecastQueryService{
		gameConfigRepo:     gameConfigRepo,
		forecastRepo:       forecastRepo,
		marketForecastRepo: marketForecastRepo,
	}
}

func (s *PlayerOrderForecastQueryService) GetMarketForecast(ctx context.Context) (*OrderMarketForecastResult, error) {
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	template, err := ResolveOrderTemplateByGameConfig(*gameConfig)
	if err != nil {
		return nil, err
	}
	items, err := s.forecastRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list forecast control: %w", err)
	}
	items = filterForecastControlsByTemplate(items, template)
	forecastRows, err := s.marketForecastRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list market forecast: %w", err)
	}
	forecastRows = filterMarketForecastsByTemplate(forecastRows, template)
	result := buildMarketForecastResultForTemplate(mergeForecastControlDefaultsForTemplate(items, template), forecastRows, template)
	return &result, nil
}

func (s *AdminOrderQueryService) validateYear(ctx context.Context, yearNo int) error {
	if yearNo < 1 {
		return ErrAdminOrderYearInvalid
	}
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return fmt.Errorf("load game config: %w", err)
	}
	if yearNo > gameConfig.FinalYear {
		return ErrAdminOrderYearInvalid
	}
	return nil
}

func (s *AdminOrderQueryService) buildForecastYearLocks(ctx context.Context) ([]OrderForecastYearLock, error) {
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	return buildForecastYearLocks(ctx, gameConfig, s.batchRepo, s.poolRepo, s.bidRepo, s.stateRepo, s.groupYearRepo, s.selectionRepo, s.sequenceRepo)
}

func buildForecastYearLocks(
	ctx context.Context,
	gameConfig *entity.GameConfig,
	batchRepo *repository.OrderGenerationBatchRepository,
	poolRepo *repository.OrderPoolRepository,
	bidRepo *repository.GroupMarketBidRepository,
	stateRepo *repository.MarketBiddingStateRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	selectionRepo *repository.GroupOrderSelectionRepository,
	sequenceRepo *repository.MarketSelectionOrderRepository,
) ([]OrderForecastYearLock, error) {
	result := make([]OrderForecastYearLock, 0, forecastControlMaxYear)
	for yearNo := forecastControlMinYear; yearNo <= forecastControlMaxYear; yearNo++ {
		lockInfo, err := resolveOrderYearLock(ctx, gameConfig, batchRepo, poolRepo, bidRepo, stateRepo, groupYearRepo, selectionRepo, sequenceRepo, yearNo)
		if err != nil {
			return nil, err
		}
		result = append(result, OrderForecastYearLock{
			YearNo: lockInfo.YearNo,
			Locked: lockInfo.Locked,
			Reason: lockInfo.Reason,
		})
	}
	return result, nil
}

func buildForecastYearLocksWithDB(ctx context.Context, db *gorm.DB) ([]OrderForecastYearLock, error) {
	gameConfigRepo := repository.NewGameConfigRepository(db)
	gameConfig, err := gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	return buildForecastYearLocks(
		ctx,
		gameConfig,
		repository.NewOrderGenerationBatchRepository(db),
		repository.NewOrderPoolRepository(db),
		repository.NewGroupMarketBidRepository(db),
		repository.NewMarketBiddingStateRepository(db),
		repository.NewGroupYearStateRepository(db),
		repository.NewGroupOrderSelectionRepository(db),
		repository.NewMarketSelectionOrderRepository(db),
	)
}

func resolveOrderYearLockWithDB(ctx context.Context, db *gorm.DB, yearNo int) (orderYearLockInfo, error) {
	gameConfigRepo := repository.NewGameConfigRepository(db)
	gameConfig, err := gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return orderYearLockInfo{}, fmt.Errorf("load game config: %w", err)
	}
	return resolveOrderYearLock(
		ctx,
		gameConfig,
		repository.NewOrderGenerationBatchRepository(db),
		repository.NewOrderPoolRepository(db),
		repository.NewGroupMarketBidRepository(db),
		repository.NewMarketBiddingStateRepository(db),
		repository.NewGroupYearStateRepository(db),
		repository.NewGroupOrderSelectionRepository(db),
		repository.NewMarketSelectionOrderRepository(db),
		yearNo,
	)
}

func resolveOrderYearLock(
	ctx context.Context,
	gameConfig *entity.GameConfig,
	batchRepo *repository.OrderGenerationBatchRepository,
	poolRepo *repository.OrderPoolRepository,
	bidRepo *repository.GroupMarketBidRepository,
	stateRepo *repository.MarketBiddingStateRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	selectionRepo *repository.GroupOrderSelectionRepository,
	sequenceRepo *repository.MarketSelectionOrderRepository,
	yearNo int,
) (orderYearLockInfo, error) {
	if gameConfig != nil && yearNo < gameConfig.CurrentOpenYear {
		return lockedOrderYear(yearNo, orderYearLockReasonHistoricalYear), nil
	}
	if groupYearRepo != nil {
		count, err := groupYearRepo.CountRollbackPendingByYear(ctx, yearNo)
		if err != nil {
			return orderYearLockInfo{}, fmt.Errorf("count rollback pending year: %w", err)
		}
		if count > 0 {
			return lockedOrderYear(yearNo, orderYearLockReasonRollbackPending), nil
		}
	}
	if batchRepo != nil {
		if _, err := batchRepo.FindConfirmedByYear(ctx, yearNo); err == nil {
			return lockedOrderYear(yearNo, orderYearLockReasonPoolConfirmed), nil
		} else if !repository.IsRecordNotFound(err) {
			return orderYearLockInfo{}, fmt.Errorf("load confirmed order batch: %w", err)
		}
	}
	if poolRepo != nil {
		hasSelected, err := poolRepo.HasSelectedByYear(ctx, yearNo)
		if err != nil {
			return orderYearLockInfo{}, fmt.Errorf("check selected order pool: %w", err)
		}
		if hasSelected {
			return lockedOrderYear(yearNo, orderYearLockReasonSelectedOrders), nil
		}
	}
	if selectionRepo != nil {
		count, err := selectionRepo.CountByYear(ctx, yearNo)
		if err != nil {
			return orderYearLockInfo{}, fmt.Errorf("count group order selections: %w", err)
		}
		if count > 0 {
			return lockedOrderYear(yearNo, orderYearLockReasonGroupSelections), nil
		}
	}
	if bidRepo != nil {
		count, err := bidRepo.CountByYear(ctx, yearNo)
		if err != nil {
			return orderYearLockInfo{}, fmt.Errorf("count market bids: %w", err)
		}
		if count > 0 {
			return lockedOrderYear(yearNo, orderYearLockReasonMarketBids), nil
		}
	}
	if sequenceRepo != nil {
		count, err := sequenceRepo.CountByYear(ctx, yearNo)
		if err != nil {
			return orderYearLockInfo{}, fmt.Errorf("count market selection orders: %w", err)
		}
		if count > 0 {
			return lockedOrderYear(yearNo, orderYearLockReasonSelectionOrders), nil
		}
	}
	if stateRepo != nil {
		started, err := stateRepo.HasWorkflowFactByYear(ctx, yearNo)
		if err != nil {
			return orderYearLockInfo{}, fmt.Errorf("check order workflow state: %w", err)
		}
		if started {
			return lockedOrderYear(yearNo, orderYearLockReasonWorkflowAdvanced), nil
		}
	}
	return orderYearLockInfo{YearNo: yearNo}, nil
}

func lockedOrderYear(yearNo int, reason string) orderYearLockInfo {
	return orderYearLockInfo{
		YearNo: yearNo,
		Locked: true,
		Reason: reason,
	}
}

type AdminOrderCommandService struct {
	db *gorm.DB
}

func NewAdminOrderCommandService(db *gorm.DB) *AdminOrderCommandService {
	return &AdminOrderCommandService{db: db}
}

func (s *AdminOrderCommandService) UploadExcel(ctx context.Context, cmd UploadOrderExcelCommand) (*UploadOrderExcelResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	if len(cmd.Content) == 0 {
		return nil, ErrAdminOrderFileRequired
	}

	parsed, parseErr := ParseOrderExcel(cmd.Content)
	status := enum.OrderImportStatusSuccess
	parseErrors := []string{}
	if parseErr != nil {
		status = enum.OrderImportStatusFailed
		parseErrors = append(parseErrors, parseErr.Error())
		parsed = &ParsedOrderWorkbook{}
	}
	if parsed == nil {
		parsed = &ParsedOrderWorkbook{}
	}
	now := time.Now()
	payloadJSON, err := marshalJSON(parsed)
	if err != nil {
		return nil, fmt.Errorf("marshal parsed order payload: %w", err)
	}
	parseErrorJSON, err := marshalJSON(append(parseErrors, parsed.Warnings...))
	if err != nil {
		return nil, fmt.Errorf("marshal parse warnings: %w", err)
	}

	item := &entity.OrderImportBatch{
		OriginalFileName: strings.TrimSpace(cmd.FileName),
		FileSize:         cmd.FileSize,
		ParseStatus:      status,
		ParsedOrderCount: len(parsed.Orders),
		ParseError:       parseErrorJSON,
		ParsedPayload:    payloadJSON,
		UploaderID:       cmd.OperatorID,
		UploaderName:     operatorName,
		UploadedAt:       now,
		BaseEntity: entity.BaseEntity{
			Creator:    operatorName,
			CreateTime: now,
			Updater:    operatorName,
			UpdateTime: now,
		},
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		importRepo := repository.NewOrderImportBatchRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)
		if err := importRepo.Create(ctx, item); err != nil {
			return fmt.Errorf("create order import batch: %w", err)
		}
		actionPayload, err := marshalJSON(map[string]any{
			"batchId":          item.ID,
			"originalFileName": item.OriginalFileName,
			"parseStatus":      item.ParseStatus,
			"parsedOrderCount": item.ParsedOrderCount,
		})
		if err != nil {
			return err
		}
		return actionRepo.Create(ctx, &entity.AdminActionLog{
			ActionCode:    adminActionCodeUploadOrderExcel,
			TargetGroupID: nil,
			TargetYearNo:  nil,
			ActionPayload: actionPayload,
			StateBefore:   []byte("{}"),
			StateAfter:    []byte("{}"),
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
		})
	})
	if err != nil {
		return nil, fmt.Errorf("upload order excel transaction: %w", err)
	}
	if parseErr != nil {
		return nil, ErrAdminOrderParseFailed
	}

	return buildUploadOrderExcelResult(item, parsed), nil
}

func (s *AdminOrderCommandService) UpdateForecastControl(ctx context.Context, cmd UpdateOrderForecastControlCommand) (*UpdateOrderForecastControlResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	now := time.Now()
	var result *UpdateOrderForecastControlResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		forecastRepo := repository.NewOrderForecastControlRepository(tx)
		marketForecastRepo := repository.NewOrderMarketForecastRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)

		gameConfig, err := gameConfigRepo.GetCurrent(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}
		template, err := ResolveOrderTemplateByGameConfig(*gameConfig)
		if err != nil {
			return err
		}
		if err := validateForecastControlCommandForTemplate(cmd, template); err != nil {
			return err
		}
		existingControls, err := forecastRepo.ListAll(ctx)
		if err != nil {
			return fmt.Errorf("list existing forecast control: %w", err)
		}
		existingControls = filterForecastControlsByTemplate(existingControls, template)
		changedYears := changedForecastControlYears(mergeForecastControlDefaultsForTemplate(existingControls, template), cmd.Items)
		autoPreviewYears := make([]int, 0)
		for yearNo := range changedYears {
			if yearNo > gameConfig.FinalYear {
				continue
			}
			lockInfo, err := resolveOrderYearLockWithDB(ctx, tx, yearNo)
			if err != nil {
				return err
			}
			if lockInfo.Locked {
				return fmt.Errorf("%w: %s", ErrAdminOrderPoolLocked, lockInfo.Reason)
			}
			autoPreviewYears = append(autoPreviewYears, yearNo)
		}
		sort.Ints(autoPreviewYears)

		existingForecasts, err := marketForecastRepo.ListAll(ctx)
		if err != nil {
			return fmt.Errorf("list existing market forecast: %w", err)
		}
		existingNarratives := forecastNarrativeMap(existingForecasts)

		items := make([]entity.OrderForecastControl, 0, len(cmd.Items))
		for _, reqItem := range cmd.Items {
			yearNo := reqItem.YearNo
			marketCode := normalizeMarketCode(reqItem.MarketCode)
			orderType := normalizeOrderType(reqItem.OrderType)
			items = append(items, entity.OrderForecastControl{
				OrderTemplateVersion: template.TemplateVersion,
				YearNo:               yearNo,
				ForecastStageCode:    forecastStageCodeForYear(yearNo),
				MarketCode:           marketCode,
				OrderType:            orderType,
				OrderCount:           reqItem.OrderCount,
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: now,
					Updater:    operatorName,
					UpdateTime: now,
				},
			})
		}
		if err := forecastRepo.UpsertBatch(ctx, items); err != nil {
			return fmt.Errorf("upsert forecast control: %w", err)
		}

		controlRows, err := forecastRepo.ListAll(ctx)
		if err != nil {
			return fmt.Errorf("reload forecast control: %w", err)
		}
		controlRows = mergeForecastControlDefaultsForTemplate(filterForecastControlsByTemplate(controlRows, template), template)
		narratives := mergeForecastNarratives(existingNarratives, cmd.Narratives)
		forecastRows, err := buildMarketForecastEntitiesForTemplate(controlRows, narratives, operatorName, now, template)
		if err != nil {
			return err
		}
		if err := marketForecastRepo.UpsertBatch(ctx, forecastRows); err != nil {
			return fmt.Errorf("upsert market forecast: %w", err)
		}
		forecastRows, err = marketForecastRepo.ListAll(ctx)
		if err != nil {
			return fmt.Errorf("reload market forecast: %w", err)
		}
		forecastRows = filterMarketForecastsByTemplate(forecastRows, template)
		for _, yearNo := range autoPreviewYears {
			if _, err := generateOrderPreviewInTx(ctx, tx, GenerateOrderPoolCommand{
				YearNo:       yearNo,
				Overwrite:    true,
				OperatorID:   cmd.OperatorID,
				OperatorName: operatorName,
			}, operatorName); err != nil {
				return err
			}
		}

		actionPayload, err := marshalJSON(map[string]any{
			"items":            cmd.Items,
			"narratives":       cmd.Narratives,
			"autoPreviewYears": autoPreviewYears,
		})
		if err != nil {
			return err
		}
		if err := actionRepo.Create(ctx, &entity.AdminActionLog{
			ActionCode:    adminActionCodeUpdateForecastControl,
			ActionPayload: actionPayload,
			StateBefore:   []byte("{}"),
			StateAfter:    []byte("{}"),
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
		}); err != nil {
			return fmt.Errorf("create admin action log: %w", err)
		}

		yearLocks, err := buildForecastYearLocksWithDB(ctx, tx)
		if err != nil {
			return err
		}
		result = &UpdateOrderForecastControlResult{
			OrderTemplate:    BuildOrderTemplateView(template),
			Items:            buildOrderForecastControlItemsForTemplate(controlRows, template),
			Narratives:       buildOrderForecastNarrativesForTemplate(forecastRows, template),
			Forecast:         buildMarketForecastResultForTemplate(controlRows, forecastRows, template),
			YearLocks:        yearLocks,
			AutoPreviewYears: autoPreviewYears,
			UpdatedAt:        now.Format(time.RFC3339),
			UpdatedBy:        operatorName,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("update order forecast control transaction: %w", err)
	}
	return result, nil
}

func (s *AdminOrderCommandService) UpdateMarketConfig(ctx context.Context, cmd UpdateOrderMarketConfigCommand) (*UpdateOrderMarketConfigResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)

	var result *UpdateOrderMarketConfigResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		batchRepo := repository.NewOrderGenerationBatchRepository(tx)
		configRepo := repository.NewOrderGenerationConfigRepository(tx)
		forecastRepo := repository.NewOrderForecastControlRepository(tx)
		marketRepo := repository.NewOrderMarketConfigRepository(tx)
		poolRepo := repository.NewOrderPoolRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)

		gameConfig, err := gameConfigRepo.GetCurrent(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}
		template, err := ResolveOrderTemplateByGameConfig(*gameConfig)
		if err != nil {
			return err
		}
		if err := validateMarketConfigCommandForTemplate(cmd, template); err != nil {
			return err
		}
		if cmd.YearNo < 1 || cmd.YearNo > gameConfig.FinalYear {
			return ErrAdminOrderYearInvalid
		}
		lockInfo, err := resolveOrderYearLockWithDB(ctx, tx, cmd.YearNo)
		if err != nil {
			return err
		}
		if lockInfo.Locked {
			return ErrAdminOrderMarketConfigLocked
		}

		now := time.Now()
		items := make([]entity.OrderMarketConfig, 0, len(cmd.Markets))
		for _, market := range cmd.Markets {
			items = append(items, entity.OrderMarketConfig{
				OrderTemplateVersion:  template.TemplateVersion,
				YearNo:                cmd.YearNo,
				MarketCode:            normalizeMarketCode(market.MarketCode),
				MarketEnabled:         market.Enabled,
				MarketInvestmentLimit: normalizeMarketInvestmentLimit(market.MarketInvestmentLimit),
				ConfigStatus:          enum.OrderConfigStatusDraft,
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: now,
					Updater:    operatorName,
					UpdateTime: now,
				},
			})
		}
		if err := marketRepo.UpsertBatch(ctx, items); err != nil {
			return fmt.Errorf("upsert market configs: %w", err)
		}
		if err := poolRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
			return fmt.Errorf("delete old order pool before market config update: %w", err)
		}
		if err := stateRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
			return fmt.Errorf("delete old order segment state before market config update: %w", err)
		}
		if err := batchRepo.VoidPreviewByYear(ctx, cmd.YearNo, operatorName, now); err != nil {
			return fmt.Errorf("void old preview batch: %w", err)
		}

		updatedMarkets, err := marketRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("reload market configs: %w", err)
		}
		updatedMarkets = filterMarketConfigsByTemplate(updatedMarkets, template)
		marketMap := buildMarketEnabledMapForTemplate(updatedMarkets, template)
		configs, err := configRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("reload order configs: %w", err)
		}
		forecastControls, err := forecastRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("list forecast control: %w", err)
		}
		forecastControls = filterForecastControlsByTemplate(forecastControls, template)
		forecastCounts := forecastControlCountMapForTemplate(cmd.YearNo, forecastControls, template)
		sourceCounts := defaultOrderSourceCountsForTemplate(cmd.YearNo, template)
		if _, err := generateOrderPreviewInTx(ctx, tx, GenerateOrderPoolCommand{
			YearNo:       cmd.YearNo,
			Overwrite:    true,
			OperatorID:   cmd.OperatorID,
			OperatorName: operatorName,
		}, operatorName); err != nil {
			return err
		}
		configs, err = configRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("reload order configs after preview: %w", err)
		}
		poolCounts, err := countGeneratedByYearWithRepoForTemplate(ctx, poolRepo, cmd.YearNo, template)
		if err != nil {
			return err
		}
		configs = filterGenerationConfigsByTemplate(configs, template)
		viewItems := buildControlConfigItemsForTemplate(cmd.YearNo, configs, forecastCounts, sourceCounts, poolCounts, marketMap, template)
		groupCount, err := groupRepo.CountAll(ctx)
		if err != nil {
			return fmt.Errorf("count groups: %w", err)
		}

		if err := createAdminOrderActionLog(ctx, actionRepo, adminActionCodeUpdateMarketConfig, cmd.OperatorID, operatorName, cmd.YearNo, map[string]any{
			"markets": cmd.Markets,
		}, now); err != nil {
			return err
		}

		result = &UpdateOrderMarketConfigResult{
			YearNo:    cmd.YearNo,
			Markets:   buildOrderMarketConfigItemsForTemplate(cmd.YearNo, updatedMarkets, template),
			Items:     viewItems,
			Warnings:  buildOrderControlWarnings(viewItems, int(groupCount)),
			UpdatedAt: now.Format(time.RFC3339),
			UpdatedBy: operatorName,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("update order market config transaction: %w", err)
	}
	return result, nil
}

func (s *AdminOrderCommandService) UpdateControlConfig(ctx context.Context, cmd UpdateOrderControlConfigCommand) (*UpdateOrderControlConfigResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)

	var result *UpdateOrderControlConfigResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		batchRepo := repository.NewOrderGenerationBatchRepository(tx)
		configRepo := repository.NewOrderGenerationConfigRepository(tx)
		forecastRepo := repository.NewOrderForecastControlRepository(tx)
		marketRepo := repository.NewOrderMarketConfigRepository(tx)
		poolRepo := repository.NewOrderPoolRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)

		gameConfig, err := gameConfigRepo.GetCurrent(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}
		template, err := ResolveOrderTemplateByGameConfig(*gameConfig)
		if err != nil {
			return err
		}
		if err := validateControlConfigCommandForTemplate(cmd, template); err != nil {
			return err
		}
		if cmd.YearNo < 1 || cmd.YearNo > gameConfig.FinalYear {
			return ErrAdminOrderYearInvalid
		}
		lockInfo, err := resolveOrderYearLockWithDB(ctx, tx, cmd.YearNo)
		if err != nil {
			return err
		}
		if lockInfo.Locked {
			return ErrAdminOrderPoolLocked
		}
		if err := poolRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
			return fmt.Errorf("delete old order pool before config update: %w", err)
		}
		if err := stateRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
			return fmt.Errorf("delete old order segment state before config update: %w", err)
		}
		if err := batchRepo.VoidPreviewByYear(ctx, cmd.YearNo, operatorName, time.Now()); err != nil {
			return fmt.Errorf("void old preview batch: %w", err)
		}
		marketConfigs, err := marketRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("list order market configs: %w", err)
		}
		if len(marketConfigs) == 0 {
			if err := ensureDefaultMarketConfigsForTemplate(ctx, marketRepo, cmd.YearNo, operatorName, time.Now(), template); err != nil {
				return fmt.Errorf("ensure default market configs: %w", err)
			}
			marketConfigs, err = marketRepo.ListByYear(ctx, cmd.YearNo)
			if err != nil {
				return fmt.Errorf("reload order market configs: %w", err)
			}
		}
		marketConfigs = filterMarketConfigsByTemplate(marketConfigs, template)
		marketMap := buildMarketEnabledMapForTemplate(marketConfigs, template)
		forecastControls, err := forecastRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("list forecast control: %w", err)
		}
		forecastControls = filterForecastControlsByTemplate(forecastControls, template)
		forecastCounts := forecastControlCountMapForTemplate(cmd.YearNo, forecastControls, template)
		if err := validateEffectiveReleaseSequencesForTemplate(cmd.Items, marketMap, forecastCounts, cmd.YearNo, template); err != nil {
			return err
		}
		sourceCounts := defaultOrderSourceCountsForTemplate(cmd.YearNo, template)

		now := time.Now()
		items := make([]entity.OrderGenerationConfig, 0, len(cmd.Items))
		for _, reqItem := range cmd.Items {
			marketCode := strings.ToUpper(strings.TrimSpace(reqItem.MarketCode))
			orderType := strings.ToUpper(strings.TrimSpace(reqItem.OrderType))
			orderCount := forecastCounts[segmentKey(cmd.YearNo, marketCode, orderType)]
			if !isMarketEnabled(marketMap, marketCode) {
				orderCount = 0
			}
			items = append(items, entity.OrderGenerationConfig{
				OrderTemplateVersion: template.TemplateVersion,
				YearNo:               cmd.YearNo,
				MarketCode:           marketCode,
				OrderType:            orderType,
				OrderCount:           orderCount,
				ReleaseSequenceNo:    reqItem.ReleaseSequenceNo,
				ConfigStatus:         enum.OrderConfigStatusDraft,
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: now,
					Updater:    operatorName,
					UpdateTime: now,
				},
			})
		}
		if err := configRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
			return fmt.Errorf("delete old order configs: %w", err)
		}
		if err := configRepo.CreateBatch(ctx, items); err != nil {
			return fmt.Errorf("create order configs: %w", err)
		}
		if _, err := generateOrderPreviewInTx(ctx, tx, GenerateOrderPoolCommand{
			YearNo:       cmd.YearNo,
			Overwrite:    true,
			OperatorID:   cmd.OperatorID,
			OperatorName: operatorName,
		}, operatorName); err != nil {
			return err
		}

		configs, err := configRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("reload order configs: %w", err)
		}
		poolCounts, err := countGeneratedByYearWithRepoForTemplate(ctx, poolRepo, cmd.YearNo, template)
		if err != nil {
			return err
		}
		configs = filterGenerationConfigsByTemplate(configs, template)
		viewItems := buildControlConfigItemsForTemplate(cmd.YearNo, configs, forecastCounts, sourceCounts, poolCounts, marketMap, template)
		groupCount, err := groupRepo.CountAll(ctx)
		if err != nil {
			return fmt.Errorf("count groups: %w", err)
		}
		warnings := buildOrderControlWarnings(viewItems, int(groupCount))

		targetYearNo := cmd.YearNo
		actionPayload, err := marshalJSON(map[string]any{
			"yearNo": cmd.YearNo,
			"items":  cmd.Items,
		})
		if err != nil {
			return err
		}
		if err := actionRepo.Create(ctx, &entity.AdminActionLog{
			ActionCode:    adminActionCodeUpdateOrderConfig,
			TargetYearNo:  &targetYearNo,
			ActionPayload: actionPayload,
			StateBefore:   []byte("{}"),
			StateAfter:    []byte("{}"),
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
		}); err != nil {
			return fmt.Errorf("create admin action log: %w", err)
		}

		result = &UpdateOrderControlConfigResult{
			YearNo:    cmd.YearNo,
			Items:     viewItems,
			Warnings:  warnings,
			UpdatedAt: now.Format(time.RFC3339),
			UpdatedBy: operatorName,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("update order control config transaction: %w", err)
	}
	return result, nil
}

func (s *AdminOrderCommandService) GenerateOrderPool(ctx context.Context, cmd GenerateOrderPoolCommand) (*GenerateOrderPoolResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	if cmd.YearNo < 1 {
		return nil, ErrAdminOrderYearInvalid
	}
	var result *GenerateOrderPoolResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		generated, err := generateOrderPreviewInTx(ctx, tx, cmd, operatorName)
		if err != nil {
			return err
		}
		result = generated
		return nil
	}); err != nil {
		return nil, fmt.Errorf("generate order pool transaction: %w", err)
	}
	return result, nil
}

func generateOrderPreviewInTx(ctx context.Context, tx *gorm.DB, cmd GenerateOrderPoolCommand, operatorName string) (*GenerateOrderPoolResult, error) {
	gameConfigRepo := repository.NewGameConfigRepository(tx)
	batchRepo := repository.NewOrderGenerationBatchRepository(tx)
	configRepo := repository.NewOrderGenerationConfigRepository(tx)
	forecastRepo := repository.NewOrderForecastControlRepository(tx)
	marketForecastRepo := repository.NewOrderMarketForecastRepository(tx)
	marketRepo := repository.NewOrderMarketConfigRepository(tx)
	poolRepo := repository.NewOrderPoolRepository(tx)
	stateRepo := repository.NewMarketBiddingStateRepository(tx)
	actionRepo := repository.NewAdminActionLogRepository(tx)

	gameConfig, err := gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	if cmd.YearNo < 1 || cmd.YearNo > gameConfig.FinalYear {
		return nil, ErrAdminOrderYearInvalid
	}
	template, err := ResolveOrderTemplateByGameConfig(*gameConfig)
	if err != nil {
		return nil, err
	}
	lockInfo, err := resolveOrderYearLockWithDB(ctx, tx, cmd.YearNo)
	if err != nil {
		return nil, err
	}
	if lockInfo.Locked {
		return nil, ErrAdminOrderPoolLocked
	}
	existingCount, err := poolRepo.CountByYear(ctx, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("count existing order pool: %w", err)
	}
	if existingCount > 0 && !cmd.Overwrite {
		return nil, ErrAdminOrderPoolLocked
	}

	configs, err := configRepo.ListByYear(ctx, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("list order configs: %w", err)
	}
	configs = filterGenerationConfigsByTemplate(configs, template)
	marketConfigs, err := marketRepo.ListByYear(ctx, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("list order market configs: %w", err)
	}
	marketConfigs = filterMarketConfigsByTemplate(marketConfigs, template)
	if len(marketConfigs) == 0 {
		if err := ensureDefaultMarketConfigsForTemplate(ctx, marketRepo, cmd.YearNo, operatorName, time.Now(), template); err != nil {
			return nil, fmt.Errorf("ensure default market configs: %w", err)
		}
		marketConfigs, err = marketRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return nil, fmt.Errorf("reload order market configs: %w", err)
		}
		marketConfigs = filterMarketConfigsByTemplate(marketConfigs, template)
	}
	marketMap := buildMarketEnabledMapForTemplate(marketConfigs, template)
	forecastControls, err := forecastRepo.ListByYear(ctx, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("list forecast control: %w", err)
	}
	forecastControls = filterForecastControlsByTemplate(forecastControls, template)
	forecastCounts := forecastControlCountMapForTemplate(cmd.YearNo, forecastControls, template)
	if len(configs) == 0 {
		configs = buildOrderGenerationConfigsFromForecastForTemplate(cmd.YearNo, forecastCounts, marketMap, operatorName, time.Now(), template)
		if err := configRepo.CreateBatch(ctx, configs); err != nil {
			return nil, fmt.Errorf("create default order release configs: %w", err)
		}
	} else {
		configs = applyForecastCountsToConfigs(configs, forecastCounts)
		configs = applyMarketEnabledToConfigsForTemplate(configs, marketMap, template)
		if err := configRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
			return nil, fmt.Errorf("delete old order configs before forecast snapshot: %w", err)
		}
		if err := configRepo.CreateBatch(ctx, configs); err != nil {
			return nil, fmt.Errorf("recreate order configs with forecast snapshot: %w", err)
		}
	}

	now := time.Now()
	if err := batchRepo.VoidPreviewByYear(ctx, cmd.YearNo, operatorName, now); err != nil {
		return nil, fmt.Errorf("void old preview batch: %w", err)
	}
	if err := poolRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
		return nil, fmt.Errorf("delete old order pool: %w", err)
	}
	if err := stateRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
		return nil, fmt.Errorf("delete old order segment state: %w", err)
	}
	seed := fmt.Sprintf("%d", now.UnixNano())
	controlSnapshot, err := marshalJSON(map[string]any{
		"orderTemplate":   BuildOrderTemplateView(template),
		"marketConfigs":   buildOrderMarketConfigItemsForTemplate(cmd.YearNo, marketConfigs, template),
		"configs":         configs,
		"forecastControl": buildOrderForecastControlItemsForTemplate(mergeForecastControlDefaultsForTemplate(forecastControls, template), template),
	})
	if err != nil {
		return nil, err
	}
	forecastRows, err := marketForecastRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list market forecast: %w", err)
	}
	forecastRows = filterMarketForecastsByTemplate(forecastRows, template)
	forecastSnapshot, err := marshalJSON(buildMarketForecastResultForTemplate(mergeForecastControlDefaultsForTemplate(forecastControls, template), forecastRows, template))
	if err != nil {
		return nil, err
	}
	params := defaultOrderGenerationParameters(template)
	parameterSnapshot, err := marshalJSON(params)
	if err != nil {
		return nil, err
	}
	batch := &entity.OrderGenerationBatch{
		OrderTemplateVersion: template.TemplateVersion,
		YearNo:               cmd.YearNo,
		BatchStatus:          enum.OrderGenerationBatchStatusPreview,
		FormulaVersion:       template.FormulaVersion,
		RandomSeed:           seed,
		ControlSnapshot:      controlSnapshot,
		ForecastSnapshot:     forecastSnapshot,
		ParameterSnapshot:    parameterSnapshot,
		OrderDetail:          []byte("[]"),
		GeneratedOrderCount:  0,
		GeneratedByID:        cmd.OperatorID,
		GeneratedByName:      operatorName,
		GeneratedAt:          now,
		BaseEntity: entity.BaseEntity{
			Creator:    operatorName,
			CreateTime: now,
			Updater:    operatorName,
			UpdateTime: now,
		},
	}
	if err := batchRepo.Create(ctx, batch); err != nil {
		return nil, fmt.Errorf("create order generation batch: %w", err)
	}
	poolItems, details, params, err := buildGeneratedOrderPoolItems(configs, batch.ID, seed, operatorName, now, template)
	if err != nil {
		return nil, err
	}
	detailJSON, err := marshalJSON(details)
	if err != nil {
		return nil, err
	}
	parameterSnapshot, err = marshalJSON(params)
	if err != nil {
		return nil, err
	}
	batch.ParameterSnapshot = parameterSnapshot
	if err := batchRepo.UpdateGenerationPayload(ctx, batch.ID, detailJSON, len(poolItems), operatorName, now); err != nil {
		return nil, fmt.Errorf("update order generation batch payload: %w", err)
	}
	if err := poolRepo.CreateBatch(ctx, poolItems); err != nil {
		return nil, fmt.Errorf("create order pool: %w", err)
	}
	states := buildSegmentStatesFromConfigsForTemplate(configs, marketMap, operatorName, now, template)
	if err := stateRepo.UpsertBatch(ctx, states); err != nil {
		return nil, fmt.Errorf("create segment states: %w", err)
	}

	targetYearNo := cmd.YearNo
	actionPayload, err := marshalJSON(map[string]any{
		"yearNo":         cmd.YearNo,
		"orderTemplate":  template.TemplateVersion,
		"batchId":        batch.ID,
		"generatedCount": len(poolItems),
		"segmentCount":   len(states),
		"overwrite":      cmd.Overwrite,
	})
	if err != nil {
		return nil, err
	}
	if err := actionRepo.Create(ctx, &entity.AdminActionLog{
		ActionCode:    adminActionCodeGenerateOrderPreview,
		TargetYearNo:  &targetYearNo,
		ActionPayload: actionPayload,
		StateBefore:   []byte("{}"),
		StateAfter:    []byte("{}"),
		OperatorID:    cmd.OperatorID,
		OperatorName:  operatorName,
		OperateTime:   now,
	}); err != nil {
		return nil, fmt.Errorf("create admin action log: %w", err)
	}

	return &GenerateOrderPoolResult{
		YearNo:         cmd.YearNo,
		BatchID:        batch.ID,
		BatchStatus:    batch.BatchStatus,
		RandomSeed:     seed,
		FormulaVersion: template.FormulaVersion,
		GeneratedCount: len(poolItems),
		SegmentCount:   len(states),
		Warnings:       nil,
		GeneratedAt:    now.Format(time.RFC3339),
		GeneratedBy:    operatorName,
	}, nil
}

func (s *AdminOrderCommandService) ConfirmOrderPool(ctx context.Context, cmd ConfirmOrderPoolCommand) (*ConfirmOrderPoolResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	if cmd.YearNo < 1 || cmd.BatchID <= 0 {
		return nil, ErrAdminOrderYearInvalid
	}
	var result *ConfirmOrderPoolResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		batchRepo := repository.NewOrderGenerationBatchRepository(tx)
		configRepo := repository.NewOrderGenerationConfigRepository(tx)
		marketRepo := repository.NewOrderMarketConfigRepository(tx)
		poolRepo := repository.NewOrderPoolRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)
		if err := validateFormalOrderYear(ctx, gameConfigRepo, cmd.YearNo); err != nil {
			return err
		}
		gameConfig, err := gameConfigRepo.GetCurrent(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}
		template, err := ResolveOrderTemplateByGameConfig(*gameConfig)
		if err != nil {
			return err
		}
		if _, err := batchRepo.FindConfirmedByYear(ctx, cmd.YearNo); err == nil {
			return ErrAdminOrderPoolLocked
		} else if !repository.IsRecordNotFound(err) {
			return fmt.Errorf("load confirmed order batch: %w", err)
		}
		batch, err := batchRepo.GetByIDForUpdate(ctx, cmd.BatchID)
		if err != nil {
			if repository.IsRecordNotFound(err) {
				return ErrAdminOrderPreviewNotFound
			}
			return fmt.Errorf("load preview batch: %w", err)
		}
		if batch.YearNo != cmd.YearNo || batch.BatchStatus != enum.OrderGenerationBatchStatusPreview {
			return ErrAdminOrderPreviewNotFound
		}
		lockInfo, err := resolveOrderYearLockWithDB(ctx, tx, cmd.YearNo)
		if err != nil {
			return err
		}
		if lockInfo.Locked {
			return ErrAdminOrderPoolLocked
		}
		now := time.Now()
		if err := batchRepo.Confirm(ctx, batch.ID, cmd.OperatorID, operatorName, now); err != nil {
			return fmt.Errorf("confirm order batch: %w", err)
		}
		if err := configRepo.UpdateStatusAndBatchByYear(ctx, cmd.YearNo, enum.OrderConfigStatusLocked, batch.ID, operatorName); err != nil {
			return fmt.Errorf("lock order configs: %w", err)
		}
		if err := ensureDefaultMarketConfigsForTemplate(ctx, marketRepo, cmd.YearNo, operatorName, now, template); err != nil {
			return fmt.Errorf("ensure default market configs: %w", err)
		}
		if err := marketRepo.UpdateStatusAndBatchByYear(ctx, cmd.YearNo, enum.OrderConfigStatusLocked, batch.ID, operatorName); err != nil {
			return fmt.Errorf("lock market configs: %w", err)
		}
		count, err := poolRepo.CountByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("count order pool: %w", err)
		}
		states, err := stateRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("list order states: %w", err)
		}
		targetYearNo := cmd.YearNo
		actionPayload, err := marshalJSON(map[string]any{
			"yearNo":         cmd.YearNo,
			"batchId":        batch.ID,
			"generatedCount": count,
			"segmentCount":   len(states),
		})
		if err != nil {
			return err
		}
		if err := actionRepo.Create(ctx, &entity.AdminActionLog{
			ActionCode:    adminActionCodeConfirmOrderPool,
			TargetYearNo:  &targetYearNo,
			ActionPayload: actionPayload,
			StateBefore:   []byte("{}"),
			StateAfter:    []byte("{}"),
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
		}); err != nil {
			return fmt.Errorf("create admin action log: %w", err)
		}
		if _, err := createGlobalSnapshot(ctx, tx, CreateGlobalSnapshotCommand{
			YearNo:       cmd.YearNo,
			SnapshotType: enum.SnapshotTypeAuto,
			TriggerCode:  enum.SnapshotTriggerOrderPoolConfirmed,
			Description:  "订单池确认后自动快照",
			OperatorID:   cmd.OperatorID,
			OperatorName: operatorName,
			OperateTime:  now,
		}); err != nil {
			return err
		}
		result = &ConfirmOrderPoolResult{
			YearNo:         cmd.YearNo,
			BatchID:        batch.ID,
			GeneratedCount: int(count),
			SegmentCount:   len(states),
			ConfirmedAt:    now.Format(time.RFC3339),
			ConfirmedBy:    operatorName,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("confirm order pool transaction: %w", err)
	}
	return result, nil
}

func defaultOrderSegments() []OrderSegmentDefinition {
	return defaultOrderTemplate().Segments()
}

type orderCodeName struct {
	code string
	name string
}

type forecastStageDefinition struct {
	code      string
	name      string
	yearRange string
	years     []int
}

func adminOrderMarkets() []orderCodeName {
	return orderMarketsForTemplate(defaultOrderTemplate())
}

func orderMarketsForTemplate(template OrderTemplateDefinition) []orderCodeName {
	items := make([]orderCodeName, 0, len(template.Markets))
	for _, market := range template.Markets {
		items = append(items, orderCodeName{market.Code, market.Name})
	}
	return items
}

func adminOrderTypes() []orderCodeName {
	return orderTypesForTemplate(defaultOrderTemplate())
}

func orderTypesForTemplate(template OrderTemplateDefinition) []orderCodeName {
	items := make([]orderCodeName, 0, len(template.OrderTypes))
	for _, orderType := range template.OrderTypes {
		items = append(items, orderCodeName{orderType.Code, orderType.Name})
	}
	return items
}

func forecastStages() []forecastStageDefinition {
	return []forecastStageDefinition{
		{code: "YEAR_1_3", name: "1~3年", yearRange: "1~3年", years: []int{1, 2, 3}},
		{code: "YEAR_4_5", name: "4~5年", yearRange: "4~5年", years: []int{4, 5}},
		{code: "YEAR_6_8", name: "6~8年", yearRange: "6~8年", years: []int{6, 7, 8}},
	}
}

func forecastStageCodeForYear(yearNo int) string {
	switch {
	case yearNo >= 1 && yearNo <= 3:
		return "YEAR_1_3"
	case yearNo >= 4 && yearNo <= 5:
		return "YEAR_4_5"
	default:
		return "YEAR_6_8"
	}
}

func forecastStageName(stageCode string) string {
	for _, stage := range forecastStages() {
		if stage.code == stageCode {
			return stage.name
		}
	}
	return stageCode
}

func isValidForecastStageCode(stageCode string) bool {
	for _, stage := range forecastStages() {
		if stage.code == stageCode {
			return true
		}
	}
	return false
}

func forecastNarrativeKey(stageCode string, marketCode string) string {
	return strings.ToUpper(strings.TrimSpace(stageCode)) + "|" + normalizeMarketCode(marketCode)
}

func forecastNarrativeMap(rows []entity.OrderMarketForecast) map[string]string {
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		result[forecastNarrativeKey(row.ForecastStageCode, row.MarketCode)] = row.Narrative
	}
	return result
}

func mergeForecastNarratives(existing map[string]string, updates []UpdateOrderForecastNarrativeItem) map[string]string {
	result := make(map[string]string, len(existing)+len(updates))
	for key, value := range existing {
		result[key] = value
	}
	for _, item := range updates {
		result[forecastNarrativeKey(item.ForecastStageCode, item.MarketCode)] = strings.TrimSpace(item.Content)
	}
	return result
}

func buildControlConfigItems(yearNo int, configs []entity.OrderGenerationConfig, forecastCounts map[string]int, sourceCounts map[string]int, poolCounts map[string]int, marketConfigMap map[string]bool) []OrderControlConfigItem {
	return buildControlConfigItemsForTemplate(yearNo, configs, forecastCounts, sourceCounts, poolCounts, marketConfigMap, defaultOrderTemplate())
}

func buildControlConfigItemsForTemplate(yearNo int, configs []entity.OrderGenerationConfig, forecastCounts map[string]int, sourceCounts map[string]int, poolCounts map[string]int, marketConfigMap map[string]bool, template OrderTemplateDefinition) []OrderControlConfigItem {
	configMap := make(map[string]entity.OrderGenerationConfig, len(configs))
	for _, item := range configs {
		configMap[segmentKey(item.YearNo, item.MarketCode, item.OrderType)] = item
	}
	segments := template.Segments()
	items := make([]OrderControlConfigItem, 0, len(segments))
	for index, segment := range segments {
		key := segmentKey(yearNo, segment.MarketCode, segment.OrderType)
		config, exists := configMap[key]
		orderCount := forecastCounts[key]
		releaseSequenceNo := defaultOrderReleaseSequenceStart + index
		configStatus := enum.OrderConfigStatusDraft
		if exists {
			releaseSequenceNo = config.ReleaseSequenceNo
			configStatus = config.ConfigStatus
		}
		marketEnabled := isMarketEnabledForTemplate(marketConfigMap, segment.MarketCode, template)
		if !marketEnabled {
			orderCount = 0
		}
		items = append(items, OrderControlConfigItem{
			YearNo:            yearNo,
			MarketCode:        segment.MarketCode,
			MarketName:        segment.MarketName,
			MarketEnabled:     marketEnabled,
			OrderType:         segment.OrderType,
			OrderTypeName:     segment.OrderTypeName,
			OrderCount:        orderCount,
			ReleaseSequenceNo: releaseSequenceNo,
			AvailableCount:    sourceCounts[key],
			ConfigStatus:      configStatus,
			GeneratedCount:    poolCounts[key],
		})
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ReleaseSequenceNo < items[j].ReleaseSequenceNo
	})
	return items
}

func buildOrderControlWarnings(items []OrderControlConfigItem, groupCount int) []OrderControlWarning {
	warnings := make([]OrderControlWarning, 0)
	total := 0
	marketTotals := map[string]int{}
	for _, item := range items {
		total += item.OrderCount
		marketTotals[item.MarketCode] += item.OrderCount
		if !item.MarketEnabled {
			continue
		}
		if item.AvailableCount > 0 && item.OrderCount > item.AvailableCount {
			warnings = append(warnings, OrderControlWarning{
				Level:      "WARN",
				Message:    fmt.Sprintf("%s-%s 配置数量 %d 大于首版单产品最大订单数 %d", item.MarketName, item.OrderTypeName, item.OrderCount, item.AvailableCount),
				MarketCode: item.MarketCode,
				OrderType:  item.OrderType,
			})
		}
		if item.OrderCount > 0 && item.OrderCount < groupCount {
			warnings = append(warnings, OrderControlWarning{
				Level:      "INFO",
				Message:    fmt.Sprintf("%s-%s 订单数少于当前小组数，系统只提示不强制保底", item.MarketName, item.OrderTypeName),
				MarketCode: item.MarketCode,
				OrderType:  item.OrderType,
			})
		}
	}
	if groupCount > 0 && total < groupCount {
		warnings = append(warnings, OrderControlWarning{
			Level:   "WARN",
			Message: fmt.Sprintf("全年订单总数 %d 少于当前小组数 %d，可能有小组拿不到订单", total, groupCount),
		})
	}
	seenMarkets := map[string]OrderControlConfigItem{}
	for _, item := range items {
		if _, exists := seenMarkets[item.MarketCode]; !exists {
			seenMarkets[item.MarketCode] = item
		}
	}
	for marketCode, item := range seenMarkets {
		if marketTotals[marketCode] > 0 && marketTotals[marketCode] < groupCount {
			warnings = append(warnings, OrderControlWarning{
				Level:      "INFO",
				Message:    fmt.Sprintf("%s 订单总数少于当前小组数，管理员可按现场规则决定是否继续", item.MarketName),
				MarketCode: marketCode,
			})
			delete(marketTotals, marketCode)
		}
	}
	return warnings
}

func validateControlConfigCommand(cmd UpdateOrderControlConfigCommand) error {
	return validateControlConfigCommandForTemplate(cmd, defaultOrderTemplate())
}

func validateControlConfigCommandForTemplate(cmd UpdateOrderControlConfigCommand, template OrderTemplateDefinition) error {
	if cmd.YearNo < 1 || len(cmd.Items) == 0 {
		return ErrAdminOrderConfigInvalid
	}
	seenSegments := map[string]bool{}
	for _, item := range cmd.Items {
		marketCode := strings.ToUpper(strings.TrimSpace(item.MarketCode))
		orderType := strings.ToUpper(strings.TrimSpace(item.OrderType))
		if !template.IsValidMarketCode(marketCode) || !template.IsValidOrderType(orderType) || item.ReleaseSequenceNo <= 0 {
			return ErrAdminOrderConfigInvalid
		}
		key := segmentKey(cmd.YearNo, marketCode, orderType)
		if seenSegments[key] {
			return ErrAdminOrderConfigInvalid
		}
		seenSegments[key] = true
	}
	return nil
}

func validateForecastControlCommand(cmd UpdateOrderForecastControlCommand) error {
	return validateForecastControlCommandForTemplate(cmd, defaultOrderTemplate())
}

func validateForecastControlCommandForTemplate(cmd UpdateOrderForecastControlCommand, template OrderTemplateDefinition) error {
	if len(cmd.Items) == 0 {
		return ErrAdminOrderForecastControlInvalid
	}
	params := defaultOrderGenerationParameters(template)
	seenSegments := map[string]bool{}
	for _, item := range cmd.Items {
		yearNo := item.YearNo
		marketCode := normalizeMarketCode(item.MarketCode)
		orderType := normalizeOrderType(item.OrderType)
		if yearNo < forecastControlMinYear ||
			yearNo > forecastControlMaxYear ||
			!template.IsValidMarketCode(marketCode) ||
			!template.IsValidOrderType(orderType) ||
			item.OrderCount < 0 ||
			(params.MaxCardCount > 0 && item.OrderCount > params.MaxCardCount) {
			return ErrAdminOrderForecastControlInvalid
		}
		key := segmentKey(yearNo, marketCode, orderType)
		if seenSegments[key] {
			return ErrAdminOrderForecastControlInvalid
		}
		seenSegments[key] = true
	}
	seenNarratives := map[string]bool{}
	for _, item := range cmd.Narratives {
		stageCode := strings.ToUpper(strings.TrimSpace(item.ForecastStageCode))
		marketCode := normalizeMarketCode(item.MarketCode)
		if !isValidForecastStageCode(stageCode) || !template.IsValidMarketCode(marketCode) {
			return ErrAdminOrderForecastControlInvalid
		}
		key := forecastNarrativeKey(stageCode, marketCode)
		if seenNarratives[key] {
			return ErrAdminOrderForecastControlInvalid
		}
		seenNarratives[key] = true
	}
	return nil
}

func validateEffectiveReleaseSequences(items []UpdateOrderControlConfigItem, marketConfigMap map[string]bool, forecastCounts map[string]int, yearNo int) error {
	return validateEffectiveReleaseSequencesForTemplate(items, marketConfigMap, forecastCounts, yearNo, defaultOrderTemplate())
}

func validateEffectiveReleaseSequencesForTemplate(items []UpdateOrderControlConfigItem, marketConfigMap map[string]bool, forecastCounts map[string]int, yearNo int, template OrderTemplateDefinition) error {
	seenSequences := map[int]bool{}
	for _, item := range items {
		marketCode := strings.ToUpper(strings.TrimSpace(item.MarketCode))
		orderType := strings.ToUpper(strings.TrimSpace(item.OrderType))
		if !isMarketEnabledForTemplate(marketConfigMap, marketCode, template) || forecastCounts[segmentKey(yearNo, marketCode, orderType)] <= 0 {
			continue
		}
		if item.ReleaseSequenceNo <= 0 {
			return ErrAdminOrderConfigInvalid
		}
		if seenSequences[item.ReleaseSequenceNo] {
			return ErrAdminOrderReleaseSequenceDuplicated
		}
		seenSequences[item.ReleaseSequenceNo] = true
	}
	return nil
}

func validateMarketConfigCommand(cmd UpdateOrderMarketConfigCommand) error {
	return validateMarketConfigCommandForTemplate(cmd, defaultOrderTemplate())
}

func validateMarketConfigCommandForTemplate(cmd UpdateOrderMarketConfigCommand, template OrderTemplateDefinition) error {
	if cmd.YearNo < 1 || len(cmd.Markets) == 0 {
		return ErrAdminOrderConfigInvalid
	}
	seen := map[string]bool{}
	for _, item := range cmd.Markets {
		marketCode := normalizeMarketCode(item.MarketCode)
		if !template.IsValidMarketCode(marketCode) {
			return ErrAdminOrderConfigInvalid
		}
		if item.MarketInvestmentLimit != nil && (*item.MarketInvestmentLimit < 0 || !isWholeNumber(*item.MarketInvestmentLimit)) {
			return ErrAdminOrderMarketInvestmentLimitInvalid
		}
		if seen[marketCode] {
			return ErrAdminOrderConfigInvalid
		}
		seen[marketCode] = true
	}
	return nil
}

func loadOrderBatch(ctx context.Context, repo *repository.OrderImportBatchRepository, sourceBatchID *int64) (*entity.OrderImportBatch, error) {
	if sourceBatchID != nil && *sourceBatchID > 0 {
		batch, err := repo.GetByID(ctx, *sourceBatchID)
		if err != nil {
			if repository.IsRecordNotFound(err) {
				return nil, ErrAdminOrderBatchNotFound
			}
			return nil, fmt.Errorf("load order batch: %w", err)
		}
		if batch.ParseStatus != enum.OrderImportStatusSuccess {
			return nil, ErrAdminOrderBatchNotFound
		}
		return batch, nil
	}
	batch, err := repo.FindLatestSuccess(ctx)
	if err != nil {
		if repository.IsRecordNotFound(err) {
			return nil, ErrAdminOrderBatchNotFound
		}
		return nil, fmt.Errorf("load latest order batch: %w", err)
	}
	return batch, nil
}

func buildOrderPoolItemsFromConfig(configs []entity.OrderGenerationConfig, parsedOrders []ParsedOrderCard, batchID int64, operatorName string, now time.Time) ([]entity.OrderPool, error) {
	ordersBySegment := make(map[string][]ParsedOrderCard)
	for _, order := range parsedOrders {
		key := segmentKey(order.YearNo, order.MarketCode, order.OrderType)
		ordersBySegment[key] = append(ordersBySegment[key], order)
	}
	for key := range ordersBySegment {
		sort.SliceStable(ordersBySegment[key], func(i, j int) bool {
			return ordersBySegment[key][i].SourceRowIndex < ordersBySegment[key][j].SourceRowIndex
		})
	}
	items := make([]entity.OrderPool, 0)
	for _, config := range configs {
		if config.OrderCount <= 0 {
			continue
		}
		key := segmentKey(config.YearNo, config.MarketCode, config.OrderType)
		candidates := ordersBySegment[key]
		if len(candidates) < config.OrderCount {
			return nil, &OrderSourceInsufficientError{
				YearNo:         config.YearNo,
				MarketName:     marketName(config.MarketCode),
				OrderTypeName:  orderTypeName(config.OrderType),
				RequiredCount:  config.OrderCount,
				AvailableCount: len(candidates),
			}
		}
		for index := 0; index < config.OrderCount; index++ {
			source := candidates[index]
			sourceBatchID := batchID
			sourceRowIndex := source.SourceRowIndex
			items = append(items, entity.OrderPool{
				YearNo:          source.YearNo,
				MarketCode:      source.MarketCode,
				OrderType:       source.OrderType,
				BusinessOrderNo: fmt.Sprintf("CARD-%02d", index+1),
				OrderAmount:     source.OrderAmount,
				OrderQuantity:   source.OrderQuantity,
				UnitPrice:       source.UnitPrice,
				AccountTerm:     source.AccountTerm,
				PoolStatus:      enum.OrderPoolStatusAvailable,
				SourceBatchID:   &sourceBatchID,
				SourceSheetName: source.SourceSheetName,
				SourceCell:      source.SourceCell,
				SourceRowIndex:  &sourceRowIndex,
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: now,
					Updater:    operatorName,
					UpdateTime: now,
				},
			})
		}
	}
	return items, nil
}

type OrderSourceInsufficientError struct {
	YearNo         int
	MarketName     string
	OrderTypeName  string
	RequiredCount  int
	AvailableCount int
}

func (e *OrderSourceInsufficientError) Error() string {
	return fmt.Sprintf("%d年%s-%s 可用订单不足：需要 %d 个，Excel 中只有 %d 个", e.YearNo, e.MarketName, e.OrderTypeName, e.RequiredCount, e.AvailableCount)
}

func (e *OrderSourceInsufficientError) Unwrap() error {
	return ErrAdminOrderSourceInsufficient
}

func buildSegmentStatesFromConfigs(configs []entity.OrderGenerationConfig, marketConfigMap map[string]bool, operatorName string, now time.Time) []entity.MarketBiddingState {
	return buildSegmentStatesFromConfigsForTemplate(configs, marketConfigMap, operatorName, now, defaultOrderTemplate())
}

func buildSegmentStatesFromConfigsForTemplate(configs []entity.OrderGenerationConfig, marketConfigMap map[string]bool, operatorName string, now time.Time, template OrderTemplateDefinition) []entity.MarketBiddingState {
	states := make([]entity.MarketBiddingState, 0, len(configs))
	for _, config := range configs {
		status := enum.OrderSegmentStatusWaitingInvestment
		if !isMarketEnabledForTemplate(marketConfigMap, config.MarketCode, template) {
			status = enum.OrderSegmentStatusMarketDisabled
		} else if config.OrderCount <= 0 {
			status = enum.OrderSegmentStatusNoOrderConfig
		}
		states = append(states, entity.MarketBiddingState{
			OrderTemplateVersion: template.TemplateVersion,
			YearNo:               config.YearNo,
			MarketCode:           config.MarketCode,
			OrderType:            config.OrderType,
			SegmentCode:          fmt.Sprintf("%s_%s", config.MarketCode, config.OrderType),
			ReleaseSequenceNo:    config.ReleaseSequenceNo,
			SegmentStatus:        status,
			BaseEntity: entity.BaseEntity{
				Creator:    operatorName,
				CreateTime: now,
				Updater:    operatorName,
				UpdateTime: now,
			},
		})
	}
	return states
}

func buildUploadOrderExcelResult(item *entity.OrderImportBatch, parsed *ParsedOrderWorkbook) *UploadOrderExcelResult {
	previewLimit := defaultOrderParsePreviewLimit
	if len(parsed.Orders) < previewLimit {
		previewLimit = len(parsed.Orders)
	}
	return &UploadOrderExcelResult{
		BatchID:          item.ID,
		OriginalFileName: item.OriginalFileName,
		ParseStatus:      item.ParseStatus,
		ParsedOrderCount: item.ParsedOrderCount,
		Warnings:         parsed.Warnings,
		Summary:          buildOrderSourceSummary(parsed.Orders),
		PreviewOrders:    parsed.Orders[:previewLimit],
		UploadedAt:       item.UploadedAt.Format(time.RFC3339),
	}
}

func buildOrderSourceSummary(orders []ParsedOrderCard) []OrderSourceSummary {
	counts := map[string]OrderSourceSummary{}
	for _, order := range orders {
		key := segmentKey(order.YearNo, order.MarketCode, order.OrderType)
		item := counts[key]
		item.YearNo = order.YearNo
		item.MarketCode = order.MarketCode
		item.MarketName = marketName(order.MarketCode)
		item.OrderType = order.OrderType
		item.OrderTypeName = orderTypeName(order.OrderType)
		item.AvailableCount++
		counts[key] = item
	}
	result := make([]OrderSourceSummary, 0, len(counts))
	for _, item := range counts {
		result = append(result, item)
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].YearNo != result[j].YearNo {
			return result[i].YearNo < result[j].YearNo
		}
		if result[i].MarketCode != result[j].MarketCode {
			return marketSortIndex(result[i].MarketCode) < marketSortIndex(result[j].MarketCode)
		}
		return orderTypeSortIndex(result[i].OrderType) < orderTypeSortIndex(result[j].OrderType)
	})
	return result
}

func buildOrderMarketConfigItems(yearNo int, configs []entity.OrderMarketConfig) []OrderMarketConfigItem {
	return buildOrderMarketConfigItemsForTemplate(yearNo, configs, defaultOrderTemplate())
}

func buildOrderMarketConfigItemsForTemplate(yearNo int, configs []entity.OrderMarketConfig, template OrderTemplateDefinition) []OrderMarketConfigItem {
	configMap := make(map[string]entity.OrderMarketConfig, len(configs))
	for _, item := range configs {
		configMap[normalizeMarketCode(item.MarketCode)] = item
	}
	items := make([]OrderMarketConfigItem, 0, len(template.Markets))
	for _, market := range template.Markets {
		config, exists := configMap[market.Code]
		enabled := template.DefaultMarketEnabled(market.Code)
		status := enum.OrderConfigStatusDraft
		var lockedBatchID *int64
		var marketInvestmentLimit *float64
		if exists {
			enabled = config.MarketEnabled
			status = config.ConfigStatus
			lockedBatchID = config.LockedBatchID
			marketInvestmentLimit = normalizeMarketInvestmentLimit(config.MarketInvestmentLimit)
		}
		items = append(items, OrderMarketConfigItem{
			YearNo:                yearNo,
			MarketCode:            market.Code,
			MarketName:            market.Name,
			Enabled:               enabled,
			MarketInvestmentLimit: marketInvestmentLimit,
			ConfigStatus:          status,
			LockedBatchID:         lockedBatchID,
		})
	}
	return items
}

func defaultForecastControlEntities(operatorName string, now time.Time) []entity.OrderForecastControl {
	return defaultForecastControlEntitiesForTemplate(operatorName, now, defaultOrderTemplate())
}

func defaultForecastControlEntitiesForTemplate(operatorName string, now time.Time, template OrderTemplateDefinition) []entity.OrderForecastControl {
	segments := template.Segments()
	items := make([]entity.OrderForecastControl, 0, forecastControlMaxYear*len(segments))
	for yearNo := forecastControlMinYear; yearNo <= forecastControlMaxYear; yearNo++ {
		for _, segment := range segments {
			items = append(items, entity.OrderForecastControl{
				OrderTemplateVersion: template.TemplateVersion,
				YearNo:               yearNo,
				ForecastStageCode:    forecastStageCodeForYear(yearNo),
				MarketCode:           segment.MarketCode,
				OrderType:            segment.OrderType,
				OrderCount:           0,
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: now,
					Updater:    operatorName,
					UpdateTime: now,
				},
			})
		}
	}
	return items
}

func mergeForecastControlDefaults(items []entity.OrderForecastControl) []entity.OrderForecastControl {
	return mergeForecastControlDefaultsForTemplate(items, defaultOrderTemplate())
}

func mergeForecastControlDefaultsForTemplate(items []entity.OrderForecastControl, template OrderTemplateDefinition) []entity.OrderForecastControl {
	now := time.Now()
	resultMap := make(map[string]entity.OrderForecastControl, len(items))
	for _, item := range items {
		if !template.IsValidMarketCode(item.MarketCode) || !template.IsValidOrderType(item.OrderType) {
			continue
		}
		resultMap[segmentKey(item.YearNo, item.MarketCode, item.OrderType)] = item
	}
	for _, item := range defaultForecastControlEntitiesForTemplate("system", now, template) {
		key := segmentKey(item.YearNo, item.MarketCode, item.OrderType)
		if _, exists := resultMap[key]; !exists {
			resultMap[key] = item
		}
	}
	result := make([]entity.OrderForecastControl, 0, len(resultMap))
	for _, item := range resultMap {
		result = append(result, item)
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].YearNo != result[j].YearNo {
			return result[i].YearNo < result[j].YearNo
		}
		if result[i].MarketCode != result[j].MarketCode {
			return template.MarketSortIndex(result[i].MarketCode) < template.MarketSortIndex(result[j].MarketCode)
		}
		return template.OrderTypeSortIndex(result[i].OrderType) < template.OrderTypeSortIndex(result[j].OrderType)
	})
	return result
}

func buildOrderForecastControlItems(items []entity.OrderForecastControl) []OrderForecastControlItem {
	return buildOrderForecastControlItemsForTemplate(items, defaultOrderTemplate())
}

func buildOrderForecastControlItemsForTemplate(items []entity.OrderForecastControl, template OrderTemplateDefinition) []OrderForecastControlItem {
	items = mergeForecastControlDefaultsForTemplate(items, template)
	result := make([]OrderForecastControlItem, 0, len(items))
	for _, item := range items {
		stageCode := forecastStageCodeForYear(item.YearNo)
		result = append(result, OrderForecastControlItem{
			YearNo:            item.YearNo,
			ForecastStageCode: stageCode,
			ForecastStageName: forecastStageName(stageCode),
			MarketCode:        item.MarketCode,
			MarketName:        template.MarketName(item.MarketCode),
			OrderType:         item.OrderType,
			OrderTypeName:     template.OrderTypeName(item.OrderType),
			OrderCount:        item.OrderCount,
		})
	}
	return result
}

func buildOrderForecastNarratives(rows []entity.OrderMarketForecast) []OrderForecastNarrativeItem {
	return buildOrderForecastNarrativesForTemplate(rows, defaultOrderTemplate())
}

func buildOrderForecastNarrativesForTemplate(rows []entity.OrderMarketForecast, template OrderTemplateDefinition) []OrderForecastNarrativeItem {
	narrativeMap := forecastNarrativeMap(rows)
	result := make([]OrderForecastNarrativeItem, 0, len(forecastStages())*len(template.Markets))
	for _, stage := range forecastStages() {
		for _, market := range template.Markets {
			key := forecastNarrativeKey(stage.code, market.Code)
			result = append(result, OrderForecastNarrativeItem{
				ForecastStageCode: stage.code,
				ForecastStageName: stage.name,
				MarketCode:        market.Code,
				MarketName:        market.Name,
				Content:           narrativeMap[key],
			})
		}
	}
	return result
}

func changedForecastControlYears(existing []entity.OrderForecastControl, updates []UpdateOrderForecastControlItem) map[int]bool {
	existingMap := forecastControlCountMapAll(existing)
	result := make(map[int]bool)
	for _, item := range updates {
		yearNo := item.YearNo
		marketCode := normalizeMarketCode(item.MarketCode)
		orderType := normalizeOrderType(item.OrderType)
		key := segmentKey(yearNo, marketCode, orderType)
		if existingMap[key] != item.OrderCount {
			result[yearNo] = true
		}
	}
	return result
}

func buildMarketForecastEntities(controlRows []entity.OrderForecastControl, narratives map[string]string, operatorName string, now time.Time) ([]entity.OrderMarketForecast, error) {
	return buildMarketForecastEntitiesForTemplate(controlRows, narratives, operatorName, now, defaultOrderTemplate())
}

func buildMarketForecastEntitiesForTemplate(controlRows []entity.OrderForecastControl, narratives map[string]string, operatorName string, now time.Time, template OrderTemplateDefinition) ([]entity.OrderMarketForecast, error) {
	controlSnapshot, err := marshalJSON(buildOrderForecastControlItemsForTemplate(controlRows, template))
	if err != nil {
		return nil, err
	}
	forecast := buildMarketForecastResultForTemplate(controlRows, nil, template)
	items := make([]entity.OrderMarketForecast, 0, len(forecast.Stages)*len(template.Markets))
	for _, stage := range forecast.Stages {
		for _, market := range stage.Markets {
			dataJSON, err := marshalJSON(market.Years)
			if err != nil {
				return nil, err
			}
			items = append(items, entity.OrderMarketForecast{
				OrderTemplateVersion: template.TemplateVersion,
				ForecastStageCode:    stage.ForecastStageCode,
				MarketCode:           market.MarketCode,
				ForecastData:         dataJSON,
				Narrative:            narratives[forecastNarrativeKey(stage.ForecastStageCode, market.MarketCode)],
				FormulaVersion:       template.FormulaVersion,
				RandomSeed:           "FORECAST_DETERMINISTIC",
				ControlSnapshotJSON:  controlSnapshot,
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: now,
					Updater:    operatorName,
					UpdateTime: now,
				},
			})
		}
	}
	return items, nil
}

func buildMarketForecastResult(controlRows []entity.OrderForecastControl, forecastRows []entity.OrderMarketForecast) OrderMarketForecastResult {
	return buildMarketForecastResultForTemplate(controlRows, forecastRows, defaultOrderTemplate())
}

func buildMarketForecastResultForTemplate(controlRows []entity.OrderForecastControl, forecastRows []entity.OrderMarketForecast, template OrderTemplateDefinition) OrderMarketForecastResult {
	controlRows = mergeForecastControlDefaultsForTemplate(controlRows, template)
	controlMap := forecastControlCountMapAll(controlRows)
	narrativeMap := forecastNarrativeMap(forecastRows)
	stages := make([]OrderMarketForecastStage, 0, len(forecastStages()))
	for _, stage := range forecastStages() {
		stageView := OrderMarketForecastStage{
			ForecastStageCode: stage.code,
			ForecastStageName: stage.name,
			YearRange:         stage.yearRange,
			Years:             append([]int(nil), stage.years...),
			Markets:           make([]OrderMarketForecastMarket, 0, len(template.Markets)),
		}
		for _, market := range template.Markets {
			marketView := OrderMarketForecastMarket{
				MarketCode: market.Code,
				MarketName: market.Name,
				Narrative:  narrativeMap[forecastNarrativeKey(stage.code, market.Code)],
				Years:      make([]OrderMarketForecastYear, 0, len(stage.years)),
			}
			for _, yearNo := range stage.years {
				yearView := OrderMarketForecastYear{
					YearNo:   yearNo,
					Products: make([]OrderMarketForecastProduct, 0, len(template.OrderTypes)),
				}
				for _, orderType := range template.OrderTypes {
					count := controlMap[segmentKey(yearNo, market.Code, orderType.Code)]
					amount := forecastAmountForProductForTemplate(yearNo, market.Code, orderType.Code, count, template)
					yearView.Products = append(yearView.Products, OrderMarketForecastProduct{
						OrderType:      orderType.Code,
						OrderTypeName:  orderType.Name,
						OrderCount:     count,
						ForecastAmount: amount,
					})
					yearView.TotalOrderCount += count
					yearView.TotalForecastAmount = roundIntegerMoney(yearView.TotalForecastAmount + amount)
				}
				marketView.Years = append(marketView.Years, yearView)
			}
			stageView.Markets = append(stageView.Markets, marketView)
		}
		stages = append(stages, stageView)
	}
	return OrderMarketForecastResult{
		FormulaVersion: template.FormulaVersion,
		OrderTemplate:  BuildOrderTemplateView(template),
		Stages:         stages,
	}
}

func forecastAmountForProduct(yearNo int, marketCode string, orderType string, orderCount int) float64 {
	return forecastAmountForProductForTemplate(yearNo, marketCode, orderType, orderCount, defaultOrderTemplate())
}

func forecastAmountForProductForTemplate(yearNo int, marketCode string, orderType string, orderCount int, template OrderTemplateDefinition) float64 {
	if orderCount <= 0 {
		return 0
	}
	params := defaultOrderGenerationParameters(template)
	if template.TemplateVersion == OrderTemplateVersionAirportV1 {
		loadFactor := 70.0
		flightCount := 10000.0
		routeFactor := 1.0
		if normalizeMarketCode(marketCode) == MarketCodeInternational {
			routeFactor = 1.7
		}
		runwayFactor := 1.0
		if normalizeOrderType(orderType) == OrderTypeWideBody {
			runwayFactor = 2.4
		}
		unitPrice := loadFactor / 200 * 0.012 * routeFactor * runwayFactor
		return roundIntegerMoney(float64(orderCount) * math.Round(unitPrice*flightCount))
	}
	avgPrice := params.AveragePrices[marketCode][orderType]
	avgQuantity := float64(params.MinQuantity+params.MaxQuantity) / 2
	return roundIntegerMoney(float64(orderCount) * avgPrice * avgQuantity)
}

func roundIntegerMoney(value float64) float64 {
	return math.Round(value)
}

func forecastControlCountMap(yearNo int, controls []entity.OrderForecastControl) map[string]int {
	return forecastControlCountMapForTemplate(yearNo, controls, defaultOrderTemplate())
}

func forecastControlCountMapForTemplate(yearNo int, controls []entity.OrderForecastControl, template OrderTemplateDefinition) map[string]int {
	segments := template.Segments()
	result := make(map[string]int, len(segments))
	for _, segment := range segments {
		result[segmentKey(yearNo, segment.MarketCode, segment.OrderType)] = 0
	}
	for _, item := range controls {
		if item.YearNo == yearNo && template.IsValidMarketCode(item.MarketCode) && template.IsValidOrderType(item.OrderType) {
			result[segmentKey(item.YearNo, item.MarketCode, item.OrderType)] = item.OrderCount
		}
	}
	return result
}

func forecastControlCountMapAll(controls []entity.OrderForecastControl) map[string]int {
	result := make(map[string]int, len(controls))
	for _, item := range controls {
		result[segmentKey(item.YearNo, item.MarketCode, item.OrderType)] = item.OrderCount
	}
	return result
}

func buildOrderGenerationConfigsFromForecast(yearNo int, forecastCounts map[string]int, marketConfigMap map[string]bool, operatorName string, now time.Time) []entity.OrderGenerationConfig {
	return buildOrderGenerationConfigsFromForecastForTemplate(yearNo, forecastCounts, marketConfigMap, operatorName, now, defaultOrderTemplate())
}

func buildOrderGenerationConfigsFromForecastForTemplate(yearNo int, forecastCounts map[string]int, marketConfigMap map[string]bool, operatorName string, now time.Time, template OrderTemplateDefinition) []entity.OrderGenerationConfig {
	segments := template.Segments()
	items := make([]entity.OrderGenerationConfig, 0, len(segments))
	for index, segment := range segments {
		count := forecastCounts[segmentKey(yearNo, segment.MarketCode, segment.OrderType)]
		if !isMarketEnabledForTemplate(marketConfigMap, segment.MarketCode, template) {
			count = 0
		}
		items = append(items, entity.OrderGenerationConfig{
			OrderTemplateVersion: template.TemplateVersion,
			YearNo:               yearNo,
			MarketCode:           segment.MarketCode,
			OrderType:            segment.OrderType,
			OrderCount:           count,
			ReleaseSequenceNo:    index + 1,
			ConfigStatus:         enum.OrderConfigStatusDraft,
			BaseEntity: entity.BaseEntity{
				Creator:    operatorName,
				CreateTime: now,
				Updater:    operatorName,
				UpdateTime: now,
			},
		})
	}
	return items
}

func applyForecastCountsToConfigs(configs []entity.OrderGenerationConfig, forecastCounts map[string]int) []entity.OrderGenerationConfig {
	result := make([]entity.OrderGenerationConfig, 0, len(configs))
	for _, item := range configs {
		item.OrderCount = forecastCounts[segmentKey(item.YearNo, item.MarketCode, item.OrderType)]
		result = append(result, item)
	}
	return result
}

func buildMarketEnabledMap(configs []entity.OrderMarketConfig) map[string]bool {
	return buildMarketEnabledMapForTemplate(configs, defaultOrderTemplate())
}

func buildMarketEnabledMapForTemplate(configs []entity.OrderMarketConfig, template OrderTemplateDefinition) map[string]bool {
	snapshotMap := buildMarketConfigSnapshotMapForTemplate(configs, template)
	result := make(map[string]bool, len(snapshotMap))
	for marketCode, item := range snapshotMap {
		result[marketCode] = item.Enabled
	}
	return result
}

type orderMarketConfigSnapshot struct {
	Enabled         bool
	InvestmentLimit *float64
}

func buildMarketConfigSnapshotMap(configs []entity.OrderMarketConfig) map[string]orderMarketConfigSnapshot {
	return buildMarketConfigSnapshotMapForTemplate(configs, defaultOrderTemplate())
}

func buildMarketConfigSnapshotMapForTemplate(configs []entity.OrderMarketConfig, template OrderTemplateDefinition) map[string]orderMarketConfigSnapshot {
	result := make(map[string]orderMarketConfigSnapshot, len(template.Markets))
	for _, market := range template.Markets {
		result[market.Code] = orderMarketConfigSnapshot{
			Enabled:         template.DefaultMarketEnabled(market.Code),
			InvestmentLimit: nil,
		}
	}
	for _, item := range configs {
		marketCode := normalizeMarketCode(item.MarketCode)
		if !template.IsValidMarketCode(marketCode) {
			continue
		}
		result[marketCode] = orderMarketConfigSnapshot{
			Enabled:         item.MarketEnabled,
			InvestmentLimit: normalizeMarketInvestmentLimit(item.MarketInvestmentLimit),
		}
	}
	return result
}

func defaultMarketEnabled(marketCode string) bool {
	return defaultOrderTemplate().DefaultMarketEnabled(marketCode)
}

func isMarketEnabled(marketConfigMap map[string]bool, marketCode string) bool {
	return isMarketEnabledForTemplate(marketConfigMap, marketCode, defaultOrderTemplate())
}

func isMarketEnabledForTemplate(marketConfigMap map[string]bool, marketCode string, template OrderTemplateDefinition) bool {
	marketCode = normalizeMarketCode(marketCode)
	if marketConfigMap == nil {
		return template.DefaultMarketEnabled(marketCode)
	}
	enabled, exists := marketConfigMap[marketCode]
	if !exists {
		return template.DefaultMarketEnabled(marketCode)
	}
	return enabled
}

func applyMarketEnabledToConfigs(configs []entity.OrderGenerationConfig, marketConfigMap map[string]bool) []entity.OrderGenerationConfig {
	return applyMarketEnabledToConfigsForTemplate(configs, marketConfigMap, defaultOrderTemplate())
}

func applyMarketEnabledToConfigsForTemplate(configs []entity.OrderGenerationConfig, marketConfigMap map[string]bool, template OrderTemplateDefinition) []entity.OrderGenerationConfig {
	result := make([]entity.OrderGenerationConfig, 0, len(configs))
	for _, item := range configs {
		if !isMarketEnabledForTemplate(marketConfigMap, item.MarketCode, template) {
			item.OrderCount = 0
		}
		item.OrderTemplateVersion = template.TemplateVersion
		result = append(result, item)
	}
	return result
}

func ensureDefaultMarketConfigs(ctx context.Context, repo *repository.OrderMarketConfigRepository, yearNo int, operatorName string, now time.Time) error {
	return ensureDefaultMarketConfigsForTemplate(ctx, repo, yearNo, operatorName, now, defaultOrderTemplate())
}

func ensureDefaultMarketConfigsForTemplate(ctx context.Context, repo *repository.OrderMarketConfigRepository, yearNo int, operatorName string, now time.Time, template OrderTemplateDefinition) error {
	existing, err := repo.ListByYear(ctx, yearNo)
	if err != nil {
		return err
	}
	existing = filterMarketConfigsByTemplate(existing, template)
	if len(existing) >= len(template.Markets) {
		return nil
	}
	existingMap := make(map[string]bool, len(existing))
	for _, item := range existing {
		existingMap[normalizeMarketCode(item.MarketCode)] = true
	}
	items := make([]entity.OrderMarketConfig, 0)
	for _, market := range template.Markets {
		if existingMap[market.Code] {
			continue
		}
		items = append(items, entity.OrderMarketConfig{
			OrderTemplateVersion:  template.TemplateVersion,
			YearNo:                yearNo,
			MarketCode:            market.Code,
			MarketEnabled:         template.DefaultMarketEnabled(market.Code),
			MarketInvestmentLimit: nil,
			ConfigStatus:          enum.OrderConfigStatusDraft,
			BaseEntity: entity.BaseEntity{
				Creator:    operatorName,
				CreateTime: now,
				Updater:    operatorName,
				UpdateTime: now,
			},
		})
	}
	return repo.UpsertBatch(ctx, items)
}

func normalizeMarketInvestmentLimit(limit *float64) *float64 {
	if limit == nil {
		return nil
	}
	value := *limit
	if value == 0 {
		normalized := 0.0
		return &normalized
	}
	normalized := math.Trunc(value)
	return &normalized
}

func isWholeNumber(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && math.Abs(value-math.Round(value)) < 0.000001
}

func summarizeParsedPayload(raw []byte) map[string]int {
	var parsed ParsedOrderWorkbook
	if len(raw) == 0 || json.Unmarshal(raw, &parsed) != nil {
		return map[string]int{}
	}
	counts := make(map[string]int)
	for _, order := range parsed.Orders {
		counts[segmentKey(order.YearNo, order.MarketCode, order.OrderType)]++
	}
	return counts
}

func defaultOrderSourceCounts(yearNo int) map[string]int {
	return defaultOrderSourceCountsForTemplate(yearNo, defaultOrderTemplate())
}

func defaultOrderSourceCountsForTemplate(yearNo int, template OrderTemplateDefinition) map[string]int {
	params := defaultOrderGenerationParameters(template)
	segments := template.Segments()
	counts := make(map[string]int, len(segments))
	for _, segment := range segments {
		counts[segmentKey(yearNo, segment.MarketCode, segment.OrderType)] = params.MaxCardCount
	}
	return counts
}

func (s *AdminOrderQueryService) countGeneratedByYear(ctx context.Context, yearNo int) (map[string]int, error) {
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	template, err := ResolveOrderTemplateByGameConfig(*gameConfig)
	if err != nil {
		return nil, err
	}
	return countGeneratedByYearWithRepoForTemplate(ctx, s.poolRepo, yearNo, template)
}

func countGeneratedByYearWithRepo(ctx context.Context, repo *repository.OrderPoolRepository, yearNo int) (map[string]int, error) {
	return countGeneratedByYearWithRepoForTemplate(ctx, repo, yearNo, defaultOrderTemplate())
}

func countGeneratedByYearWithRepoForTemplate(ctx context.Context, repo *repository.OrderPoolRepository, yearNo int, template OrderTemplateDefinition) (map[string]int, error) {
	counts := make(map[string]int)
	for _, segment := range template.Segments() {
		count, err := repo.CountBySegment(ctx, yearNo, segment.MarketCode, segment.OrderType)
		if err != nil {
			return nil, fmt.Errorf("count generated order pool: %w", err)
		}
		counts[segmentKey(yearNo, segment.MarketCode, segment.OrderType)] = int(count)
	}
	return counts, nil
}

func buildOrderPoolItem(item entity.OrderPool) OrderPoolItem {
	payload := map[string]any(nil)
	if len(item.OrderPayloadJSON) > 0 {
		_ = json.Unmarshal(item.OrderPayloadJSON, &payload)
	}
	return OrderPoolItem{
		OrderID:         item.ID,
		BusinessOrderNo: formatBusinessOrderNo(item),
		CardSequenceNo:  item.CardSequenceNo,
		YearNo:          item.YearNo,
		MarketCode:      item.MarketCode,
		MarketName:      marketName(item.MarketCode),
		OrderType:       item.OrderType,
		OrderTypeName:   orderTypeName(item.OrderType),
		OrderAmount:     item.OrderAmount,
		OrderQuantity:   item.OrderQuantity,
		UnitPrice:       item.UnitPrice,
		AccountTerm:     item.AccountTerm,
		PoolStatus:      item.PoolStatus,
		SelectedGroupID: item.SelectedGroupID,
		SourceSheetName: item.SourceSheetName,
		SourceCell:      item.SourceCell,
		OrderPayload:    payload,
	}
}

func formatBusinessOrderNo(item entity.OrderPool) string {
	if strings.TrimSpace(item.BusinessOrderNo) != "" {
		return strings.TrimSpace(item.BusinessOrderNo)
	}
	if strings.TrimSpace(item.SourceCell) != "" && strings.HasPrefix(strings.ToUpper(strings.TrimSpace(item.SourceCell)), "CARD-") {
		return strings.TrimSpace(item.SourceCell)
	}
	if item.CardSequenceNo > 0 {
		return fmt.Sprintf("CARD-%02d", item.CardSequenceNo)
	}
	return fmt.Sprintf("ORDER-%d", item.ID)
}

func orderTemplateVersionMatches(value string, template OrderTemplateDefinition) bool {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	return normalized == "" || normalized == template.TemplateVersion || (normalized == OrderTemplateVersionVIPServiceV1 && template.TemplateVersion == OrderTemplateVersionVIPServiceV1)
}

func filterForecastControlsByTemplate(items []entity.OrderForecastControl, template OrderTemplateDefinition) []entity.OrderForecastControl {
	result := make([]entity.OrderForecastControl, 0, len(items))
	for _, item := range items {
		if !orderTemplateVersionMatches(item.OrderTemplateVersion, template) {
			continue
		}
		if !template.IsValidMarketCode(item.MarketCode) || !template.IsValidOrderType(item.OrderType) {
			continue
		}
		item.OrderTemplateVersion = template.TemplateVersion
		result = append(result, item)
	}
	return result
}

func filterMarketForecastsByTemplate(items []entity.OrderMarketForecast, template OrderTemplateDefinition) []entity.OrderMarketForecast {
	result := make([]entity.OrderMarketForecast, 0, len(items))
	for _, item := range items {
		if !orderTemplateVersionMatches(item.OrderTemplateVersion, template) {
			continue
		}
		if !template.IsValidMarketCode(item.MarketCode) {
			continue
		}
		item.OrderTemplateVersion = template.TemplateVersion
		result = append(result, item)
	}
	return result
}

func filterGenerationConfigsByTemplate(items []entity.OrderGenerationConfig, template OrderTemplateDefinition) []entity.OrderGenerationConfig {
	result := make([]entity.OrderGenerationConfig, 0, len(items))
	for _, item := range items {
		if !orderTemplateVersionMatches(item.OrderTemplateVersion, template) {
			continue
		}
		if !template.IsValidMarketCode(item.MarketCode) || !template.IsValidOrderType(item.OrderType) {
			continue
		}
		item.OrderTemplateVersion = template.TemplateVersion
		result = append(result, item)
	}
	return result
}

func filterMarketConfigsByTemplate(items []entity.OrderMarketConfig, template OrderTemplateDefinition) []entity.OrderMarketConfig {
	result := make([]entity.OrderMarketConfig, 0, len(items))
	for _, item := range items {
		if !orderTemplateVersionMatches(item.OrderTemplateVersion, template) {
			continue
		}
		if !template.IsValidMarketCode(item.MarketCode) {
			continue
		}
		item.OrderTemplateVersion = template.TemplateVersion
		result = append(result, item)
	}
	return result
}

func buildOrderGenerationBatchSummary(item *entity.OrderGenerationBatch) *OrderGenerationBatchSummary {
	if item == nil {
		return nil
	}
	var confirmedAt *string
	if item.ConfirmedAt != nil {
		value := item.ConfirmedAt.Format(time.RFC3339)
		confirmedAt = &value
	}
	return &OrderGenerationBatchSummary{
		BatchID:        item.ID,
		BatchStatus:    item.BatchStatus,
		FormulaVersion: item.FormulaVersion,
		RandomSeed:     item.RandomSeed,
		GeneratedCount: item.GeneratedOrderCount,
		GeneratedAt:    item.GeneratedAt.Format(time.RFC3339),
		GeneratedBy:    item.GeneratedByName,
		ConfirmedAt:    confirmedAt,
		ConfirmedBy:    item.ConfirmedByName,
	}
}

func resolveOrderGenerationStatus(preview *entity.OrderGenerationBatch, confirmed *entity.OrderGenerationBatch, states []entity.MarketBiddingState) string {
	if len(states) > 0 {
		allFinished := true
		for _, state := range states {
			if state.SegmentStatus == enum.OrderSegmentStatusSelecting {
				return "SELECTING"
			}
			if !isTerminalOrderSegmentStatus(state.SegmentStatus) {
				allFinished = false
			}
		}
		if allFinished {
			return "COMPLETED"
		}
	}
	if confirmed != nil {
		return "POOL_CONFIRMED"
	}
	if preview != nil {
		return "PREVIEW_GENERATED"
	}
	return "NOT_GENERATED"
}

func isTerminalOrderSegmentStatus(status string) bool {
	switch status {
	case enum.OrderSegmentStatusCompleted,
		enum.OrderSegmentStatusSkipped,
		enum.OrderSegmentStatusMarketDisabled,
		enum.OrderSegmentStatusNoOrderConfig:
		return true
	default:
		return false
	}
}

func segmentKey(yearNo int, marketCode string, orderType string) string {
	return fmt.Sprintf("%d|%s|%s", yearNo, strings.ToUpper(strings.TrimSpace(marketCode)), strings.ToUpper(strings.TrimSpace(orderType)))
}

func marketSortIndex(code string) int {
	switch normalizeMarketCode(code) {
	case enum.MarketCodeLocal:
		return 1
	case enum.MarketCodeRegional:
		return 2
	case enum.MarketCodeNational:
		return 3
	case enum.MarketCodeGlobal:
		return 4
	default:
		return 99
	}
}

func orderTypeSortIndex(code string) int {
	switch normalizeOrderType(code) {
	case enum.OrderTypeAgencyInspection:
		return 1
	case enum.OrderTypeTwoCabinVIP:
		return 2
	case enum.OrderTypeBusinessVIP:
		return 3
	case enum.OrderTypeMemberCustom:
		return 4
	default:
		return 99
	}
}
