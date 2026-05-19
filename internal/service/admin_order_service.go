package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
)

const (
	adminActionCodeUploadOrderExcel  = "UPLOAD_ORDER_EXCEL"
	adminActionCodeUpdateOrderConfig = "UPDATE_ORDER_CONFIG"
	adminActionCodeGenerateOrderPool = "GENERATE_ORDER_POOL"
	defaultOrderReleaseSequenceStart = 1
	defaultOrderParsePreviewLimit    = 20
)

var (
	ErrAdminOrderYearInvalid               = errors.New("admin order year invalid")
	ErrAdminOrderFileRequired              = errors.New("admin order file required")
	ErrAdminOrderParseFailed               = errors.New("admin order parse failed")
	ErrAdminOrderConfigInvalid             = errors.New("admin order config invalid")
	ErrAdminOrderReleaseSequenceDuplicated = errors.New("admin order release sequence duplicated")
	ErrAdminOrderReleaseSequenceLocked     = errors.New("admin order release sequence locked")
	ErrAdminOrderPoolLocked                = errors.New("admin order pool locked")
	ErrAdminOrderConfigNotFound            = errors.New("admin order config not found")
	ErrAdminOrderBatchNotFound             = errors.New("admin order batch not found")
	ErrAdminOrderSourceInsufficient        = errors.New("admin order source insufficient")
)

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

type OrderControlConfigItem struct {
	YearNo            int    `json:"yearNo"`
	MarketCode        string `json:"marketCode"`
	MarketName        string `json:"marketName"`
	OrderType         string `json:"orderType"`
	OrderTypeName     string `json:"orderTypeName"`
	OrderCount        int    `json:"orderCount"`
	ReleaseSequenceNo int    `json:"releaseSequenceNo"`
	AvailableCount    int    `json:"availableCount"`
	ConfigStatus      string `json:"configStatus"`
	GeneratedCount    int    `json:"generatedCount"`
}

type OrderControlWarning struct {
	Level      string `json:"level"`
	Message    string `json:"message"`
	MarketCode string `json:"marketCode,omitempty"`
	OrderType  string `json:"orderType,omitempty"`
}

type OrderControlConfigResult struct {
	YearNo                int                      `json:"yearNo"`
	FinalYear             int                      `json:"finalYear"`
	LatestBatchID         *int64                   `json:"latestBatchId"`
	LatestBatchUploadedAt *string                  `json:"latestBatchUploadedAt"`
	Items                 []OrderControlConfigItem `json:"items"`
	Warnings              []OrderControlWarning    `json:"warnings"`
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
	SourceBatchID  int64                 `json:"sourceBatchId"`
	GeneratedCount int                   `json:"generatedCount"`
	SegmentCount   int                   `json:"segmentCount"`
	Warnings       []OrderControlWarning `json:"warnings"`
	GeneratedAt    string                `json:"generatedAt"`
	GeneratedBy    string                `json:"generatedBy"`
}

type OrderPoolItem struct {
	OrderID         int64   `json:"orderId"`
	YearNo          int     `json:"yearNo"`
	MarketCode      string  `json:"marketCode"`
	MarketName      string  `json:"marketName"`
	OrderType       string  `json:"orderType"`
	OrderTypeName   string  `json:"orderTypeName"`
	OrderAmount     float64 `json:"orderAmount"`
	OrderQuantity   float64 `json:"orderQuantity"`
	UnitPrice       float64 `json:"unitPrice"`
	AccountTerm     int     `json:"accountTerm"`
	PoolStatus      string  `json:"poolStatus"`
	SelectedGroupID *int64  `json:"selectedGroupId"`
	SourceSheetName string  `json:"sourceSheetName"`
	SourceCell      string  `json:"sourceCell"`
}

type OrderPoolResult struct {
	YearNo        int             `json:"yearNo"`
	MarketCode    string          `json:"marketCode"`
	MarketName    string          `json:"marketName"`
	OrderType     string          `json:"orderType"`
	OrderTypeName string          `json:"orderTypeName"`
	List          []OrderPoolItem `json:"list"`
}

type AdminOrderQueryService struct {
	gameConfigRepo *repository.GameConfigRepository
	groupRepo      *repository.GroupRepository
	importRepo     *repository.OrderImportBatchRepository
	configRepo     *repository.OrderGenerationConfigRepository
	poolRepo       *repository.OrderPoolRepository
}

func NewAdminOrderQueryService(
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	importRepo *repository.OrderImportBatchRepository,
	configRepo *repository.OrderGenerationConfigRepository,
	poolRepo *repository.OrderPoolRepository,
) *AdminOrderQueryService {
	return &AdminOrderQueryService{
		gameConfigRepo: gameConfigRepo,
		groupRepo:      groupRepo,
		importRepo:     importRepo,
		configRepo:     configRepo,
		poolRepo:       poolRepo,
	}
}

func (s *AdminOrderQueryService) GetControlConfig(ctx context.Context, yearNo int) (*OrderControlConfigResult, error) {
	if err := s.validateYear(ctx, yearNo); err != nil {
		return nil, err
	}
	configs, err := s.configRepo.ListByYear(ctx, yearNo)
	if err != nil {
		return nil, fmt.Errorf("list order configs: %w", err)
	}
	latestBatch, batchErr := s.importRepo.FindLatestSuccess(ctx)
	if batchErr != nil && !repository.IsRecordNotFound(batchErr) {
		return nil, fmt.Errorf("load latest order batch: %w", batchErr)
	}
	var latestBatchID *int64
	var latestBatchUploadedAt *string
	sourceCounts := map[string]int{}
	if batchErr == nil {
		latestBatchID = &latestBatch.ID
		uploadedAt := latestBatch.UploadedAt.Format(time.RFC3339)
		latestBatchUploadedAt = &uploadedAt
		sourceCounts = summarizeParsedPayload(latestBatch.ParsedPayload)
	}
	poolCounts, err := s.countGeneratedByYear(ctx, yearNo)
	if err != nil {
		return nil, err
	}

	items := buildControlConfigItems(yearNo, configs, sourceCounts, poolCounts)
	groupCount, err := s.groupRepo.CountAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("count groups: %w", err)
	}
	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}
	return &OrderControlConfigResult{
		YearNo:                yearNo,
		FinalYear:             gameConfig.FinalYear,
		LatestBatchID:         latestBatchID,
		LatestBatchUploadedAt: latestBatchUploadedAt,
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
	if !enum.IsValidMarketCode(marketCode) || !enum.IsValidOrderType(orderType) {
		return nil, ErrAdminOrderConfigInvalid
	}
	items, err := s.poolRepo.ListBySegment(ctx, yearNo, marketCode, orderType)
	if err != nil {
		return nil, fmt.Errorf("list order pool: %w", err)
	}
	resultItems := make([]OrderPoolItem, 0, len(items))
	for _, item := range items {
		resultItems = append(resultItems, buildOrderPoolItem(item))
	}
	return &OrderPoolResult{
		YearNo:        yearNo,
		MarketCode:    marketCode,
		MarketName:    marketName(marketCode),
		OrderType:     orderType,
		OrderTypeName: orderTypeName(orderType),
		List:          resultItems,
	}, nil
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

func (s *AdminOrderCommandService) UpdateControlConfig(ctx context.Context, cmd UpdateOrderControlConfigCommand) (*UpdateOrderControlConfigResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	if err := validateControlConfigCommand(cmd); err != nil {
		return nil, err
	}

	var result *UpdateOrderControlConfigResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		importRepo := repository.NewOrderImportBatchRepository(tx)
		configRepo := repository.NewOrderGenerationConfigRepository(tx)
		poolRepo := repository.NewOrderPoolRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)

		gameConfig, err := gameConfigRepo.GetCurrent(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}
		if cmd.YearNo < 1 || cmd.YearNo > gameConfig.FinalYear {
			return ErrAdminOrderYearInvalid
		}
		started, err := stateRepo.HasStartedByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("check order segment started: %w", err)
		}
		if started {
			return ErrAdminOrderReleaseSequenceLocked
		}
		selected, err := poolRepo.HasSelectedByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("check selected order pool: %w", err)
		}
		if selected {
			return ErrAdminOrderPoolLocked
		}
		if err := poolRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
			return fmt.Errorf("delete old order pool before config update: %w", err)
		}
		if err := stateRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
			return fmt.Errorf("delete old order segment state before config update: %w", err)
		}

		latestBatch, batchErr := importRepo.FindLatestSuccess(ctx)
		var sourceBatchID *int64
		sourceCounts := map[string]int{}
		if batchErr == nil {
			sourceBatchID = &latestBatch.ID
			sourceCounts = summarizeParsedPayload(latestBatch.ParsedPayload)
		} else if !repository.IsRecordNotFound(batchErr) {
			return fmt.Errorf("load latest order batch: %w", batchErr)
		}

		now := time.Now()
		items := make([]entity.OrderGenerationConfig, 0, len(cmd.Items))
		for _, reqItem := range cmd.Items {
			items = append(items, entity.OrderGenerationConfig{
				YearNo:            cmd.YearNo,
				MarketCode:        strings.ToUpper(strings.TrimSpace(reqItem.MarketCode)),
				OrderType:         strings.ToUpper(strings.TrimSpace(reqItem.OrderType)),
				OrderCount:        reqItem.OrderCount,
				ReleaseSequenceNo: reqItem.ReleaseSequenceNo,
				SourceBatchID:     sourceBatchID,
				ConfigStatus:      enum.OrderConfigStatusDraft,
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

		configs, err := configRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("reload order configs: %w", err)
		}
		poolCounts, err := countGeneratedByYearWithRepo(ctx, poolRepo, cmd.YearNo)
		if err != nil {
			return err
		}
		viewItems := buildControlConfigItems(cmd.YearNo, configs, sourceCounts, poolCounts)
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
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		importRepo := repository.NewOrderImportBatchRepository(tx)
		configRepo := repository.NewOrderGenerationConfigRepository(tx)
		poolRepo := repository.NewOrderPoolRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)

		gameConfig, err := gameConfigRepo.GetCurrent(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}
		if cmd.YearNo < 1 || cmd.YearNo > gameConfig.FinalYear {
			return ErrAdminOrderYearInvalid
		}
		started, err := stateRepo.HasStartedByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("check order segment started: %w", err)
		}
		if started {
			return ErrAdminOrderPoolLocked
		}
		selected, err := poolRepo.HasSelectedByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("check selected order pool: %w", err)
		}
		if selected {
			return ErrAdminOrderPoolLocked
		}
		existingCount, err := poolRepo.CountByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("count existing order pool: %w", err)
		}
		if existingCount > 0 && !cmd.Overwrite {
			return ErrAdminOrderPoolLocked
		}

		configs, err := configRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("list order configs: %w", err)
		}
		if len(configs) == 0 {
			return ErrAdminOrderConfigNotFound
		}

		batch, err := loadOrderBatch(ctx, importRepo, cmd.SourceBatchID)
		if err != nil {
			return err
		}
		var parsed ParsedOrderWorkbook
		if err := json.Unmarshal(batch.ParsedPayload, &parsed); err != nil {
			return fmt.Errorf("unmarshal order parsed payload: %w", err)
		}

		now := time.Now()
		poolItems, err := buildOrderPoolItemsFromConfig(configs, parsed.Orders, batch.ID, operatorName, now)
		if err != nil {
			return err
		}
		if err := poolRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
			return fmt.Errorf("delete old order pool: %w", err)
		}
		if err := stateRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
			return fmt.Errorf("delete old order segment state: %w", err)
		}
		if err := poolRepo.CreateBatch(ctx, poolItems); err != nil {
			return fmt.Errorf("create order pool: %w", err)
		}
		states := buildSegmentStatesFromConfigs(configs, operatorName, now)
		if err := stateRepo.UpsertBatch(ctx, states); err != nil {
			return fmt.Errorf("create segment states: %w", err)
		}
		if err := configRepo.UpdateStatusByYear(ctx, cmd.YearNo, enum.OrderConfigStatusLocked, operatorName); err != nil {
			return fmt.Errorf("lock order configs: %w", err)
		}

		targetYearNo := cmd.YearNo
		actionPayload, err := marshalJSON(map[string]any{
			"yearNo":         cmd.YearNo,
			"sourceBatchId":  batch.ID,
			"generatedCount": len(poolItems),
			"segmentCount":   len(states),
			"overwrite":      cmd.Overwrite,
		})
		if err != nil {
			return err
		}
		if err := actionRepo.Create(ctx, &entity.AdminActionLog{
			ActionCode:    adminActionCodeGenerateOrderPool,
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

		result = &GenerateOrderPoolResult{
			YearNo:         cmd.YearNo,
			SourceBatchID:  batch.ID,
			GeneratedCount: len(poolItems),
			SegmentCount:   len(states),
			Warnings:       nil,
			GeneratedAt:    now.Format(time.RFC3339),
			GeneratedBy:    operatorName,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("generate order pool transaction: %w", err)
	}
	return result, nil
}

func defaultOrderSegments() []OrderSegmentDefinition {
	markets := []struct {
		code string
		name string
	}{
		{enum.MarketCodeLocal, "本地市场"},
		{enum.MarketCodeRegional, "区域市场"},
		{enum.MarketCodeNational, "全国市场"},
		{enum.MarketCodeGlobal, "全球市场"},
	}
	orderTypes := []struct {
		code string
		name string
	}{
		{enum.OrderTypeAgencyInspection, "代办过检"},
		{enum.OrderTypeTwoCabinVIP, "两舱贵宾"},
		{enum.OrderTypeBusinessVIP, "商务贵宾"},
		{enum.OrderTypeMemberCustom, "会员定制"},
	}
	segments := make([]OrderSegmentDefinition, 0, len(markets)*len(orderTypes))
	for _, market := range markets {
		for _, orderType := range orderTypes {
			segments = append(segments, OrderSegmentDefinition{
				MarketCode:    market.code,
				MarketName:    market.name,
				OrderType:     orderType.code,
				OrderTypeName: orderType.name,
			})
		}
	}
	return segments
}

func buildControlConfigItems(yearNo int, configs []entity.OrderGenerationConfig, sourceCounts map[string]int, poolCounts map[string]int) []OrderControlConfigItem {
	configMap := make(map[string]entity.OrderGenerationConfig, len(configs))
	for _, item := range configs {
		configMap[segmentKey(item.YearNo, item.MarketCode, item.OrderType)] = item
	}
	items := make([]OrderControlConfigItem, 0, len(defaultOrderSegments()))
	for index, segment := range defaultOrderSegments() {
		key := segmentKey(yearNo, segment.MarketCode, segment.OrderType)
		config, exists := configMap[key]
		orderCount := 0
		releaseSequenceNo := defaultOrderReleaseSequenceStart + index
		configStatus := enum.OrderConfigStatusDraft
		if exists {
			orderCount = config.OrderCount
			releaseSequenceNo = config.ReleaseSequenceNo
			configStatus = config.ConfigStatus
		}
		items = append(items, OrderControlConfigItem{
			YearNo:            yearNo,
			MarketCode:        segment.MarketCode,
			MarketName:        segment.MarketName,
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
		if item.OrderCount > item.AvailableCount && item.AvailableCount > 0 {
			warnings = append(warnings, OrderControlWarning{
				Level:      "WARN",
				Message:    fmt.Sprintf("%s-%s 配置数量 %d 大于 Excel 可用订单 %d", item.MarketName, item.OrderTypeName, item.OrderCount, item.AvailableCount),
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
	for _, segment := range defaultOrderSegments() {
		if marketTotals[segment.MarketCode] > 0 && marketTotals[segment.MarketCode] < groupCount {
			warnings = append(warnings, OrderControlWarning{
				Level:      "INFO",
				Message:    fmt.Sprintf("%s 订单总数少于当前小组数，管理员可按现场规则决定是否继续", segment.MarketName),
				MarketCode: segment.MarketCode,
			})
			delete(marketTotals, segment.MarketCode)
		}
	}
	return warnings
}

func validateControlConfigCommand(cmd UpdateOrderControlConfigCommand) error {
	if cmd.YearNo < 1 || len(cmd.Items) == 0 {
		return ErrAdminOrderConfigInvalid
	}
	seenSegments := map[string]bool{}
	seenSequences := map[int]bool{}
	for _, item := range cmd.Items {
		marketCode := strings.ToUpper(strings.TrimSpace(item.MarketCode))
		orderType := strings.ToUpper(strings.TrimSpace(item.OrderType))
		if !enum.IsValidMarketCode(marketCode) || !enum.IsValidOrderType(orderType) || item.OrderCount < 0 || item.ReleaseSequenceNo <= 0 {
			return ErrAdminOrderConfigInvalid
		}
		key := segmentKey(cmd.YearNo, marketCode, orderType)
		if seenSegments[key] {
			return ErrAdminOrderConfigInvalid
		}
		if seenSequences[item.ReleaseSequenceNo] {
			return ErrAdminOrderReleaseSequenceDuplicated
		}
		seenSegments[key] = true
		seenSequences[item.ReleaseSequenceNo] = true
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

func buildSegmentStatesFromConfigs(configs []entity.OrderGenerationConfig, operatorName string, now time.Time) []entity.MarketBiddingState {
	states := make([]entity.MarketBiddingState, 0, len(configs))
	for _, config := range configs {
		status := enum.OrderSegmentStatusWaitingRelease
		if config.OrderCount <= 0 {
			status = enum.OrderSegmentStatusSkipped
		}
		states = append(states, entity.MarketBiddingState{
			YearNo:            config.YearNo,
			MarketCode:        config.MarketCode,
			OrderType:         config.OrderType,
			SegmentCode:       fmt.Sprintf("%s_%s", config.MarketCode, config.OrderType),
			ReleaseSequenceNo: config.ReleaseSequenceNo,
			SegmentStatus:     status,
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

func (s *AdminOrderQueryService) countGeneratedByYear(ctx context.Context, yearNo int) (map[string]int, error) {
	return countGeneratedByYearWithRepo(ctx, s.poolRepo, yearNo)
}

func countGeneratedByYearWithRepo(ctx context.Context, repo *repository.OrderPoolRepository, yearNo int) (map[string]int, error) {
	counts := make(map[string]int)
	for _, segment := range defaultOrderSegments() {
		count, err := repo.CountBySegment(ctx, yearNo, segment.MarketCode, segment.OrderType)
		if err != nil {
			return nil, fmt.Errorf("count generated order pool: %w", err)
		}
		counts[segmentKey(yearNo, segment.MarketCode, segment.OrderType)] = int(count)
	}
	return counts, nil
}

func buildOrderPoolItem(item entity.OrderPool) OrderPoolItem {
	return OrderPoolItem{
		OrderID:         item.ID,
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
	}
}

func segmentKey(yearNo int, marketCode string, orderType string) string {
	return fmt.Sprintf("%d|%s|%s", yearNo, strings.ToUpper(strings.TrimSpace(marketCode)), strings.ToUpper(strings.TrimSpace(orderType)))
}

func marketSortIndex(code string) int {
	switch code {
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
	switch code {
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
