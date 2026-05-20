package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

const (
	adminActionCodeOpenOrderMarket           = "OPEN_ORDER_MARKET"
	adminActionCodeCloseOrderMarket          = "CLOSE_ORDER_MARKET"
	adminActionCodeGenerateSelectionSequence = "GENERATE_ORDER_SELECTION_SEQUENCE"
	adminActionCodeReleaseOrderSegment       = "RELEASE_ORDER_SEGMENT"
	adminActionCodeSkipOrderCurrentGroup     = "SKIP_ORDER_CURRENT_GROUP"
	orderPollingIntervalSeconds              = 3
)

var (
	ErrOrderNotRequiredForDemoYear      = errors.New("order not required for demo year")
	ErrOrderYearInvalid                 = errors.New("order year invalid")
	ErrOrderGroupUnavailable            = errors.New("order group unavailable")
	ErrOrderMarketInvalid               = errors.New("order market invalid")
	ErrOrderTypeInvalid                 = errors.New("order type invalid")
	ErrOrderPoolNotGenerated            = errors.New("order pool not generated")
	ErrOrderMarketNotOpen               = errors.New("order market not open")
	ErrOrderMarketAlreadyOpen           = errors.New("order market already open")
	ErrOrderMarketAlreadyClosed         = errors.New("order market already closed")
	ErrOrderInvestmentAlreadySubmitted  = errors.New("order investment already submitted")
	ErrOrderInvestmentInvalid           = errors.New("order investment invalid")
	ErrOrderMarketDisabledInvestment    = errors.New("order market disabled investment must be zero")
	ErrOrderSegmentNotReady             = errors.New("order segment not ready")
	ErrOrderSegmentNotSelecting         = errors.New("order segment not selecting")
	ErrOrderSegmentReleaseBlocked       = errors.New("order segment release blocked")
	ErrOrderSegmentCurrentGroupMismatch = errors.New("order segment current group mismatch")
	ErrOrderSelectionNotEligible        = errors.New("order selection not eligible")
	ErrOrderAlreadySelected             = errors.New("order already selected")
	ErrOrderCannotSelect                = errors.New("order cannot select")
	ErrOrderAdminSkipReasonRequired     = errors.New("order admin skip reason required")
	ErrOrderDeliveryStageInvalid        = errors.New("order delivery stage invalid")
	ErrOrderDeliveryOrderInvalid        = errors.New("order delivery order invalid")
	ErrOrderDeliveryRevenueMismatch     = errors.New("order delivery revenue mismatch")
	ErrOrderPrerequisiteIncomplete      = errors.New("order prerequisite incomplete")
)

type PlayerOrderQueryService struct {
	gameConfigRepo *repository.GameConfigRepository
	groupRepo      *repository.GroupRepository
	groupYearRepo  *repository.GroupYearStateRepository
	batchRepo      *repository.OrderGenerationBatchRepository
	bidRepo        *repository.GroupMarketBidRepository
	stateRepo      *repository.MarketBiddingStateRepository
	sequenceRepo   *repository.MarketSelectionOrderRepository
	poolRepo       *repository.OrderPoolRepository
	selectionRepo  *repository.GroupOrderSelectionRepository
}

func NewPlayerOrderQueryService(
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	batchRepo *repository.OrderGenerationBatchRepository,
	bidRepo *repository.GroupMarketBidRepository,
	stateRepo *repository.MarketBiddingStateRepository,
	sequenceRepo *repository.MarketSelectionOrderRepository,
	poolRepo *repository.OrderPoolRepository,
	selectionRepo *repository.GroupOrderSelectionRepository,
) *PlayerOrderQueryService {
	return &PlayerOrderQueryService{
		gameConfigRepo: gameConfigRepo,
		groupRepo:      groupRepo,
		groupYearRepo:  groupYearRepo,
		batchRepo:      batchRepo,
		bidRepo:        bidRepo,
		stateRepo:      stateRepo,
		sequenceRepo:   sequenceRepo,
		poolRepo:       poolRepo,
		selectionRepo:  selectionRepo,
	}
}

type PlayerOrderCommandService struct {
	db *gorm.DB
}

func NewPlayerOrderCommandService(db *gorm.DB) *PlayerOrderCommandService {
	return &PlayerOrderCommandService{db: db}
}

type AdminOrderControlQueryService struct {
	gameConfigRepo *repository.GameConfigRepository
	groupRepo      *repository.GroupRepository
	bidRepo        *repository.GroupMarketBidRepository
	stateRepo      *repository.MarketBiddingStateRepository
	sequenceRepo   *repository.MarketSelectionOrderRepository
	poolRepo       *repository.OrderPoolRepository
}

func NewAdminOrderControlQueryService(
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	bidRepo *repository.GroupMarketBidRepository,
	stateRepo *repository.MarketBiddingStateRepository,
	sequenceRepo *repository.MarketSelectionOrderRepository,
	poolRepo *repository.OrderPoolRepository,
) *AdminOrderControlQueryService {
	return &AdminOrderControlQueryService{
		gameConfigRepo: gameConfigRepo,
		groupRepo:      groupRepo,
		bidRepo:        bidRepo,
		stateRepo:      stateRepo,
		sequenceRepo:   sequenceRepo,
		poolRepo:       poolRepo,
	}
}

type AdminOrderControlCommandService struct {
	db *gorm.DB
}

func NewAdminOrderControlCommandService(db *gorm.DB) *AdminOrderControlCommandService {
	return &AdminOrderControlCommandService{db: db}
}

type PlayerOrderYearView struct {
	GroupID                int64                   `json:"groupId"`
	YearNo                 int                     `json:"yearNo"`
	OrderRequired          bool                    `json:"orderRequired"`
	InvestmentSubmitted    bool                    `json:"investmentSubmitted"`
	CanSubmitInvestment    bool                    `json:"canSubmitInvestment"`
	Markets                []PlayerOrderMarketView `json:"markets"`
	PollingIntervalSeconds int                     `json:"pollingIntervalSeconds"`
}

type PlayerOrderMarketView struct {
	MarketCode          string                   `json:"marketCode"`
	MarketName          string                   `json:"marketName"`
	MarketBidStatus     string                   `json:"marketBidStatus"`
	InvestmentSubmitted bool                     `json:"investmentSubmitted"`
	MarketInvestment    float64                  `json:"marketInvestment"`
	CanSubmitInvestment bool                     `json:"canSubmitInvestment"`
	CanSelectOrder      bool                     `json:"canSelectOrder"`
	SelectionSequenceNo *int                     `json:"selectionSequenceNo"`
	IsMarketLeader      bool                     `json:"isMarketLeader"`
	MarketEnabled       bool                     `json:"marketEnabled"`
	Segments            []PlayerOrderSegmentView `json:"segments"`
}

type PlayerOrderSegmentView struct {
	MarketCode          string                    `json:"marketCode"`
	MarketName          string                    `json:"marketName"`
	OrderType           string                    `json:"orderType"`
	OrderTypeName       string                    `json:"orderTypeName"`
	MarketEnabled       bool                      `json:"marketEnabled"`
	MarketInvestment    float64                   `json:"marketInvestment"`
	InvestmentSubmitted bool                      `json:"investmentSubmitted"`
	ReleaseSequenceNo   int                       `json:"releaseSequenceNo"`
	SegmentStatus       string                    `json:"segmentStatus"`
	SelectionOrder      []PlayerOrderSequenceView `json:"selectionOrder"`
	CurrentGroupID      *int64                    `json:"currentGroupId"`
	AvailableOrders     []PlayerOrderPoolItem     `json:"availableOrders"`
	LockedOrders        []PlayerOrderPoolItem     `json:"lockedOrders"`
	SelectedOrder       *PlayerOrderPoolItem      `json:"selectedOrder"`
	DeliveryStatus      string                    `json:"deliveryStatus"`
	CanSelectOrder      bool                      `json:"canSelectOrder"`
	CanPassSegment      bool                      `json:"canPassSegment"`
}

type PlayerOrderSequenceView struct {
	SequenceNo      int    `json:"sequenceNo"`
	GroupID         int64  `json:"groupId"`
	GroupName       string `json:"groupName"`
	SelectionStatus string `json:"selectionStatus"`
	IsSelf          bool   `json:"isSelf"`
	IsMarketLeader  bool   `json:"isMarketLeader"`
}

type PlayerOrderPoolItem struct {
	OrderID            int64   `json:"orderId"`
	BusinessOrderNo    string  `json:"businessOrderNo"`
	CardSequenceNo     int     `json:"cardSequenceNo"`
	OrderAmount        float64 `json:"orderAmount"`
	OrderQuantity      float64 `json:"orderQuantity"`
	UnitPrice          float64 `json:"unitPrice"`
	AccountTerm        int     `json:"accountTerm"`
	PoolStatus         string  `json:"poolStatus"`
	DeliveryStatus     string  `json:"deliveryStatus,omitempty"`
	DeliveredStageCode *string `json:"deliveredStageCode,omitempty"`
}

type SubmitMarketInvestmentCommand struct {
	GroupID      int64
	YearNo       int
	Investments  []MarketInvestmentInput
	OperatorName string
}

type MarketInvestmentInput struct {
	MarketCode       string
	OrderType        string
	MarketInvestment float64
}

type SubmitMarketInvestmentResult struct {
	GroupID        int64     `json:"groupId"`
	YearNo         int       `json:"yearNo"`
	SubmittedCount int       `json:"submittedCount"`
	SubmittedAt    time.Time `json:"submittedAt"`
}

type SelectOrderCommand struct {
	GroupID      int64
	YearNo       int
	MarketCode   string
	OrderType    string
	OrderID      int64
	OperatorName string
}

type SelectOrderResult struct {
	SelectedOrderID     int64  `json:"selectedOrderId"`
	MarketCode          string `json:"marketCode"`
	OrderType           string `json:"orderType"`
	SelectionSequenceNo int    `json:"selectionSequenceNo"`
	NextGroupID         *int64 `json:"nextGroupId"`
	SegmentStatus       string `json:"segmentStatus"`
}

type PassOrderSegmentCommand struct {
	GroupID      int64
	YearNo       int
	MarketCode   string
	OrderType    string
	OperatorName string
}

type PassOrderSegmentResult struct {
	MarketCode    string `json:"marketCode"`
	OrderType     string `json:"orderType"`
	NextGroupID   *int64 `json:"nextGroupId"`
	SegmentStatus string `json:"segmentStatus"`
}

type DeliverOrdersCommand struct {
	GroupID      int64
	YearNo       int
	StageCode    string
	OrderIDs     []int64
	OperatorName string
}

type DeliverOrdersResult struct {
	GroupID         int64     `json:"groupId"`
	YearNo          int       `json:"yearNo"`
	StageCode       string    `json:"stageCode"`
	OrderIDs        []int64   `json:"orderIds"`
	DeliveredAmount float64   `json:"deliveredAmount"`
	DeliveredAt     time.Time `json:"deliveredAt"`
}

type OpenMarketBiddingCommand struct {
	YearNo       int
	MarketCode   string
	OperatorID   int64
	OperatorName string
}

type CloseMarketBiddingCommand struct {
	YearNo       int
	MarketCode   string
	OperatorID   int64
	OperatorName string
}

type ReleaseNextSegmentCommand struct {
	YearNo       int
	OperatorID   int64
	OperatorName string
}

type GenerateSelectionSequenceCommand struct {
	YearNo       int
	OperatorID   int64
	OperatorName string
}

type AdminSkipCurrentGroupCommand struct {
	YearNo       int
	MarketCode   string
	OrderType    string
	GroupID      int64
	Reason       string
	OperatorID   int64
	OperatorName string
}

type AdminOrderControlResult struct {
	YearNo         int    `json:"yearNo"`
	MarketCode     string `json:"marketCode,omitempty"`
	MarketName     string `json:"marketName,omitempty"`
	OrderType      string `json:"orderType,omitempty"`
	OrderTypeName  string `json:"orderTypeName,omitempty"`
	SegmentStatus  string `json:"segmentStatus,omitempty"`
	CurrentGroupID *int64 `json:"currentGroupId,omitempty"`
	AffectedCount  int64  `json:"affectedCount"`
	OperatedAt     string `json:"operatedAt"`
	OperatedBy     string `json:"operatedBy"`
}

type AdminMarketSelectionStatus struct {
	YearNo          int                       `json:"yearNo"`
	MarketCode      string                    `json:"marketCode"`
	MarketName      string                    `json:"marketName"`
	MarketBidStatus string                    `json:"marketBidStatus"`
	LeaderGroupID   *int64                    `json:"leaderGroupId"`
	Bids            []AdminMarketBidView      `json:"bids"`
	Segments        []AdminOrderSegmentStatus `json:"segments"`
	CurrentSegment  *AdminOrderSegmentStatus  `json:"currentSegment"`
}

type AdminMarketBidView struct {
	GroupID          int64   `json:"groupId"`
	GroupNo          int     `json:"groupNo"`
	GroupName        string  `json:"groupName"`
	MarketInvestment float64 `json:"marketInvestment"`
	Submitted        bool    `json:"submitted"`
	BusinessStatus   string  `json:"businessStatus"`
}

type AdminOrderSegmentStatus struct {
	MarketCode        string                    `json:"marketCode"`
	MarketName        string                    `json:"marketName"`
	OrderType         string                    `json:"orderType"`
	OrderTypeName     string                    `json:"orderTypeName"`
	ReleaseSequenceNo int                       `json:"releaseSequenceNo"`
	SegmentStatus     string                    `json:"segmentStatus"`
	CurrentGroupID    *int64                    `json:"currentGroupId"`
	SelectionOrder    []AdminSelectionOrderView `json:"selectionOrder"`
	AvailableCount    int                       `json:"availableCount"`
	SelectedCount     int                       `json:"selectedCount"`
}

type AdminSelectionOrderView struct {
	SequenceNo       int     `json:"sequenceNo"`
	GroupID          int64   `json:"groupId"`
	GroupName        string  `json:"groupName"`
	MarketInvestment float64 `json:"marketInvestment"`
	IsMarketLeader   bool    `json:"isMarketLeader"`
	SelectionStatus  string  `json:"selectionStatus"`
	SelectedOrderID  *int64  `json:"selectedOrderId"`
	SelectedOrderNo  string  `json:"selectedOrderNo,omitempty"`
}

func (s *PlayerOrderQueryService) GetYearView(ctx context.Context, groupID int64, yearNo int) (*PlayerOrderYearView, error) {
	if yearNo == 0 {
		return &PlayerOrderYearView{
			GroupID:                groupID,
			YearNo:                 yearNo,
			OrderRequired:          false,
			Markets:                nil,
			PollingIntervalSeconds: orderPollingIntervalSeconds,
		}, nil
	}
	if err := validateFormalOrderYear(ctx, s.gameConfigRepo, yearNo); err != nil {
		return nil, err
	}
	if _, err := s.groupRepo.GetByID(ctx, groupID); err != nil {
		return nil, fmt.Errorf("load group: %w", err)
	}
	if _, err := s.groupYearRepo.GetByGroupIDAndYear(ctx, groupID, yearNo); err != nil {
		return nil, fmt.Errorf("load group year state: %w", err)
	}

	_, confirmedErr := s.batchRepo.FindConfirmedByYear(ctx, yearNo)
	if confirmedErr != nil {
		if repository.IsRecordNotFound(confirmedErr) {
			return &PlayerOrderYearView{
				GroupID:                groupID,
				YearNo:                 yearNo,
				OrderRequired:          true,
				InvestmentSubmitted:    false,
				CanSubmitInvestment:    false,
				Markets:                []PlayerOrderMarketView{},
				PollingIntervalSeconds: orderPollingIntervalSeconds,
			}, nil
		}
		return nil, fmt.Errorf("load confirmed order batch: %w", confirmedErr)
	}

	states, err := s.stateRepo.ListByYear(ctx, yearNo)
	if err != nil {
		return nil, fmt.Errorf("list order segment states: %w", err)
	}
	bids, err := s.bidRepo.ListByGroupYear(ctx, groupID, yearNo)
	if err != nil {
		return nil, fmt.Errorf("list self market investments: %w", err)
	}
	selfBids := map[string]entity.GroupMarketBid{}
	for _, bid := range bids {
		selfBids[segmentKey(bid.YearNo, bid.MarketCode, bid.OrderType)] = bid
	}
	sequences, err := s.sequenceRepo.ListByYear(ctx, yearNo)
	if err != nil {
		return nil, fmt.Errorf("list selection order: %w", err)
	}
	selected, err := s.selectionRepo.ListByGroupYear(ctx, groupID, yearNo)
	if err != nil {
		return nil, fmt.Errorf("list group order selections: %w", err)
	}
	groups, err := s.groupRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}

	return buildPlayerOrderYearView(groupID, yearNo, states, selfBids, sequences, selected, groups, true, s.poolRepo, ctx)
}

func (s *PlayerOrderCommandService) SubmitMarketInvestment(ctx context.Context, cmd SubmitMarketInvestmentCommand) (*SubmitMarketInvestmentResult, error) {
	normalizedInvestments, err := normalizeMarketInvestmentInputs(cmd.Investments)
	if err != nil {
		return nil, err
	}
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	now := time.Now()
	var result *SubmitMarketInvestmentResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		groupYearRepo := repository.NewGroupYearStateRepository(tx)
		batchRepo := repository.NewOrderGenerationBatchRepository(tx)
		marketRepo := repository.NewOrderMarketConfigRepository(tx)
		bidRepo := repository.NewGroupMarketBidRepository(tx)

		if err := validateFormalOrderYear(ctx, gameConfigRepo, cmd.YearNo); err != nil {
			return err
		}
		if _, err := batchRepo.FindConfirmedByYear(ctx, cmd.YearNo); err != nil {
			if repository.IsRecordNotFound(err) {
				return ErrOrderPoolNotGenerated
			}
			return fmt.Errorf("load confirmed order batch: %w", err)
		}
		group, err := groupRepo.GetByID(ctx, cmd.GroupID)
		if err != nil {
			return fmt.Errorf("load group: %w", err)
		}
		if group.BusinessStatus == enum.BusinessStatusBankrupt {
			return ErrOrderGroupUnavailable
		}
		yearState, err := groupYearRepo.GetByGroupIDAndYear(ctx, cmd.GroupID, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("load group year state: %w", err)
		}
		if yearState.YearStatus != enum.YearStatusOperating {
			return ErrOrderYearInvalid
		}
		existingCount, err := bidRepo.CountByGroupYear(ctx, cmd.GroupID, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("count market investments: %w", err)
		}
		if existingCount > 0 {
			return ErrOrderInvestmentAlreadySubmitted
		}
		marketConfigs, err := marketRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("list order market configs: %w", err)
		}
		if err := validateDisabledMarketInvestmentsAreZero(normalizedInvestments, buildMarketEnabledMap(marketConfigs)); err != nil {
			return err
		}
		items := make([]entity.GroupMarketBid, 0, len(normalizedInvestments))
		for _, investment := range normalizedInvestments {
			items = append(items, entity.GroupMarketBid{
				GroupID:          cmd.GroupID,
				YearNo:           cmd.YearNo,
				MarketCode:       investment.MarketCode,
				OrderType:        investment.OrderType,
				MarketInvestment: investment.MarketInvestment,
				BidStatus:        enum.OrderBidStatusSubmitted,
				SubmittedAt:      now,
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: now,
					Updater:    operatorName,
					UpdateTime: now,
				},
			})
		}
		if err := bidRepo.CreateBatch(ctx, items); err != nil {
			return fmt.Errorf("create market investments: %w", err)
		}
		result = &SubmitMarketInvestmentResult{
			GroupID:        cmd.GroupID,
			YearNo:         cmd.YearNo,
			SubmittedCount: len(items),
			SubmittedAt:    now,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("submit market investment transaction: %w", err)
	}
	return result, nil
}

func (s *PlayerOrderCommandService) SelectOrder(ctx context.Context, cmd SelectOrderCommand) (*SelectOrderResult, error) {
	marketCode := normalizeMarketCode(cmd.MarketCode)
	orderType := normalizeOrderType(cmd.OrderType)
	if !enum.IsValidMarketCode(marketCode) {
		return nil, ErrOrderMarketInvalid
	}
	if !enum.IsValidOrderType(orderType) {
		return nil, ErrOrderTypeInvalid
	}
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	now := time.Now()
	var result *SelectOrderResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		sequenceRepo := repository.NewMarketSelectionOrderRepository(tx)
		poolRepo := repository.NewOrderPoolRepository(tx)
		selectionRepo := repository.NewGroupOrderSelectionRepository(tx)

		if err := validateFormalOrderYear(ctx, gameConfigRepo, cmd.YearNo); err != nil {
			return err
		}
		if err := validateSelectableGroup(ctx, groupRepo, cmd.GroupID); err != nil {
			return err
		}
		segment, err := stateRepo.GetBySegmentForUpdate(ctx, cmd.YearNo, marketCode, orderType)
		if err != nil {
			return fmt.Errorf("load segment: %w", err)
		}
		if segment.SegmentStatus != enum.OrderSegmentStatusSelecting {
			return ErrOrderSegmentNotSelecting
		}
		if segment.CurrentGroupID == nil || *segment.CurrentGroupID != cmd.GroupID {
			return ErrOrderSegmentCurrentGroupMismatch
		}
		sequence, err := sequenceRepo.GetByGroupSegmentForUpdate(ctx, cmd.GroupID, cmd.YearNo, marketCode, orderType)
		if err != nil {
			return fmt.Errorf("load current selection order: %w", err)
		}
		if sequence.SelectionStatus != enum.OrderSelectionStatusCurrent {
			return ErrOrderSegmentCurrentGroupMismatch
		}
		if _, err := selectionRepo.GetByGroupSegment(ctx, cmd.GroupID, cmd.YearNo, marketCode, orderType); err == nil {
			return ErrOrderAlreadySelected
		} else if !repository.IsRecordNotFound(err) {
			return fmt.Errorf("load existing selection: %w", err)
		}
		order, err := poolRepo.GetByIDForUpdate(ctx, cmd.OrderID)
		if err != nil {
			return fmt.Errorf("load order pool: %w", err)
		}
		if order.YearNo != cmd.YearNo || order.MarketCode != marketCode || order.OrderType != orderType || order.PoolStatus != enum.OrderPoolStatusAvailable {
			return ErrOrderCannotSelect
		}
		if err := poolRepo.MarkSelected(ctx, order.ID, cmd.GroupID, operatorName, now); err != nil {
			return fmt.Errorf("mark order selected: %w", err)
		}
		if err := selectionRepo.Create(ctx, &entity.GroupOrderSelection{
			GroupID:         cmd.GroupID,
			YearNo:          cmd.YearNo,
			MarketCode:      marketCode,
			OrderType:       orderType,
			OrderID:         order.ID,
			SelectionStatus: enum.OrderSelectionStatusSelected,
			DeliveryStatus:  enum.OrderDeliveryStatusSelected,
			SelectedAt:      now,
			BaseEntity: entity.BaseEntity{
				Creator:    operatorName,
				CreateTime: now,
				Updater:    operatorName,
				UpdateTime: now,
			},
		}); err != nil {
			return fmt.Errorf("create group order selection: %w", err)
		}
		if err := sequenceRepo.SetStatus(ctx, sequence.ID, enum.OrderSelectionStatusSelected, &order.ID, nil, nil, &now, operatorName); err != nil {
			return fmt.Errorf("mark selection order selected: %w", err)
		}
		nextGroupID, status, err := advanceOrderSegment(ctx, stateRepo, sequenceRepo, poolRepo, *segment, operatorName, now)
		if err != nil {
			return err
		}
		result = &SelectOrderResult{
			SelectedOrderID:     order.ID,
			MarketCode:          marketCode,
			OrderType:           orderType,
			SelectionSequenceNo: sequence.SequenceNo,
			NextGroupID:         nextGroupID,
			SegmentStatus:       status,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("select order transaction: %w", err)
	}
	return result, nil
}

func (s *PlayerOrderCommandService) PassSegment(ctx context.Context, cmd PassOrderSegmentCommand) (*PassOrderSegmentResult, error) {
	marketCode := normalizeMarketCode(cmd.MarketCode)
	orderType := normalizeOrderType(cmd.OrderType)
	if !enum.IsValidMarketCode(marketCode) {
		return nil, ErrOrderMarketInvalid
	}
	if !enum.IsValidOrderType(orderType) {
		return nil, ErrOrderTypeInvalid
	}
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	now := time.Now()
	var result *PassOrderSegmentResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		sequenceRepo := repository.NewMarketSelectionOrderRepository(tx)
		poolRepo := repository.NewOrderPoolRepository(tx)

		if err := validateFormalOrderYear(ctx, gameConfigRepo, cmd.YearNo); err != nil {
			return err
		}
		if err := validateSelectableGroup(ctx, groupRepo, cmd.GroupID); err != nil {
			return err
		}
		segment, err := stateRepo.GetBySegmentForUpdate(ctx, cmd.YearNo, marketCode, orderType)
		if err != nil {
			return fmt.Errorf("load segment: %w", err)
		}
		if segment.SegmentStatus != enum.OrderSegmentStatusSelecting {
			return ErrOrderSegmentNotSelecting
		}
		if segment.CurrentGroupID == nil || *segment.CurrentGroupID != cmd.GroupID {
			return ErrOrderSegmentCurrentGroupMismatch
		}
		sequence, err := sequenceRepo.GetByGroupSegmentForUpdate(ctx, cmd.GroupID, cmd.YearNo, marketCode, orderType)
		if err != nil {
			return fmt.Errorf("load current selection order: %w", err)
		}
		if sequence.SelectionStatus != enum.OrderSelectionStatusCurrent {
			return ErrOrderSegmentCurrentGroupMismatch
		}
		if err := sequenceRepo.SetStatus(ctx, sequence.ID, enum.OrderSelectionStatusPassed, nil, nil, nil, &now, operatorName); err != nil {
			return fmt.Errorf("mark selection passed: %w", err)
		}
		nextGroupID, status, err := advanceOrderSegment(ctx, stateRepo, sequenceRepo, poolRepo, *segment, operatorName, now)
		if err != nil {
			return err
		}
		result = &PassOrderSegmentResult{
			MarketCode:    marketCode,
			OrderType:     orderType,
			NextGroupID:   nextGroupID,
			SegmentStatus: status,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("pass order segment transaction: %w", err)
	}
	return result, nil
}

func (s *PlayerOrderCommandService) DeliverOrders(ctx context.Context, cmd DeliverOrdersCommand) (*DeliverOrdersResult, error) {
	stageCode := normalizeOrderStageCode(cmd.StageCode)
	if !isOrderQuarterStageCode(stageCode) {
		return nil, ErrOrderDeliveryStageInvalid
	}
	orderIDs := uniquePositiveOrderIDs(cmd.OrderIDs)
	if len(orderIDs) == 0 {
		return nil, ErrOrderDeliveryOrderInvalid
	}
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	now := time.Now()
	var result *DeliverOrdersResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		groupYearRepo := repository.NewGroupYearStateRepository(tx)
		operatingRepo := repository.NewOperatingRepository(tx)
		selectionRepo := repository.NewGroupOrderSelectionRepository(tx)

		if err := validateFormalOrderYear(ctx, gameConfigRepo, cmd.YearNo); err != nil {
			return err
		}
		if err := validateSelectableGroup(ctx, groupRepo, cmd.GroupID); err != nil {
			return err
		}
		yearState, err := groupYearRepo.GetByGroupIDAndYear(ctx, cmd.GroupID, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("load group year state: %w", err)
		}
		if state.CurrentStageCode(yearState.StageStatus) != stageCode {
			return ErrOrderDeliveryStageInvalid
		}

		details, err := selectionRepo.ListSelectedOrderDetailsForUpdate(ctx, cmd.GroupID, cmd.YearNo, orderIDs)
		if err != nil {
			return fmt.Errorf("load selected orders: %w", err)
		}
		if len(details) != len(orderIDs) {
			return ErrOrderDeliveryOrderInvalid
		}
		selectionIDs := make([]int64, 0, len(details))
		deliveredAmount := 0.0
		for _, item := range details {
			if item.DeliveryStatus != enum.OrderDeliveryStatusSelected {
				return ErrOrderDeliveryOrderInvalid
			}
			if item.DeliveredStageCode != nil {
				return ErrOrderDeliveryOrderInvalid
			}
			selectionIDs = append(selectionIDs, item.SelectionID)
			deliveredAmount += item.OrderAmount
		}

		draft, err := operatingRepo.FindDraft(ctx, cmd.GroupID, cmd.YearNo)
		if err != nil {
			if repository.IsRecordNotFound(err) {
				return ErrOrderDeliveryRevenueMismatch
			}
			return fmt.Errorf("load operating draft: %w", err)
		}
		var operatingPayload payload.OperatingPayload
		if len(draft.OperatingPayload) > 0 {
			if err := json.Unmarshal(draft.OperatingPayload, &operatingPayload); err != nil {
				return fmt.Errorf("unmarshal operating draft: %w", err)
			}
		} else {
			operatingPayload = payload.NewOperatingPayload()
		}
		currentRevenue := extractOperatingQuarterSalesRevenue(operatingPayload.Normalize(), stageCode)
		deliveredInStage, err := selectionRepo.SumDeliveredAmountByStage(ctx, cmd.GroupID, cmd.YearNo, stageCode)
		if err != nil {
			return fmt.Errorf("sum delivered stage amount: %w", err)
		}
		if !sameMoneyAmount(currentRevenue, deliveredInStage+deliveredAmount) {
			return ErrOrderDeliveryRevenueMismatch
		}
		if err := selectionRepo.MarkDeliveredBySelectionIDs(ctx, selectionIDs, stageCode, operatorName, now); err != nil {
			return fmt.Errorf("mark orders delivered: %w", err)
		}
		result = &DeliverOrdersResult{
			GroupID:         cmd.GroupID,
			YearNo:          cmd.YearNo,
			StageCode:       stageCode,
			OrderIDs:        orderIDs,
			DeliveredAmount: deliveredAmount,
			DeliveredAt:     now,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("deliver orders transaction: %w", err)
	}
	return result, nil
}

func (s *AdminOrderControlQueryService) GetMarketSelectionStatus(ctx context.Context, yearNo int, marketCode string) (*AdminMarketSelectionStatus, error) {
	marketCode = normalizeMarketCode(marketCode)
	if !enum.IsValidMarketCode(marketCode) {
		return nil, ErrOrderMarketInvalid
	}
	if err := validateFormalOrderYear(ctx, s.gameConfigRepo, yearNo); err != nil {
		return nil, err
	}
	states, err := s.stateRepo.ListByYearAndMarket(ctx, yearNo, marketCode)
	if err != nil {
		return nil, fmt.Errorf("list market segment states: %w", err)
	}
	groups, err := s.groupRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}
	bids, err := s.bidRepo.ListByYearMarket(ctx, yearNo, marketCode)
	if err != nil {
		return nil, fmt.Errorf("list market bids: %w", err)
	}
	bidMap := make(map[int64][]entity.GroupMarketBid, len(bids))
	for _, bid := range bids {
		bidMap[bid.GroupID] = append(bidMap[bid.GroupID], bid)
	}
	bidViews := make([]AdminMarketBidView, 0, len(groups))
	for _, group := range groups {
		groupBids := bidMap[group.ID]
		investmentTotal := 0.0
		for _, bid := range groupBids {
			investmentTotal += bid.MarketInvestment
		}
		bidViews = append(bidViews, AdminMarketBidView{
			GroupID:          group.ID,
			GroupNo:          group.GroupNo,
			GroupName:        group.GroupName,
			MarketInvestment: investmentTotal,
			Submitted:        len(groupBids) >= 4,
			BusinessStatus:   group.BusinessStatus,
		})
	}
	groupNames := buildGroupNameMap(groups)
	segments := make([]AdminOrderSegmentStatus, 0, len(states))
	var currentSegment *AdminOrderSegmentStatus
	var leaderGroupID *int64
	for _, state := range states {
		if leaderGroupID == nil && state.LeaderGroupID != nil {
			leaderGroupID = state.LeaderGroupID
		}
		segment, err := buildAdminSegmentStatus(ctx, s.sequenceRepo, s.poolRepo, state, groupNames)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
		if state.SegmentStatus == enum.OrderSegmentStatusSelecting {
			copySegment := segment
			currentSegment = &copySegment
		}
	}
	return &AdminMarketSelectionStatus{
		YearNo:          yearNo,
		MarketCode:      marketCode,
		MarketName:      marketName(marketCode),
		MarketBidStatus: resolveMarketBidStatus(states),
		LeaderGroupID:   leaderGroupID,
		Bids:            bidViews,
		Segments:        segments,
		CurrentSegment:  currentSegment,
	}, nil
}

func (s *AdminOrderControlCommandService) OpenMarketBidding(ctx context.Context, cmd OpenMarketBiddingCommand) (*AdminOrderControlResult, error) {
	marketCode := normalizeMarketCode(cmd.MarketCode)
	if !enum.IsValidMarketCode(marketCode) {
		return nil, ErrOrderMarketInvalid
	}
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	now := time.Now()
	var result *AdminOrderControlResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)
		if err := validateFormalOrderYear(ctx, gameConfigRepo, cmd.YearNo); err != nil {
			return err
		}
		states, err := stateRepo.ListByYearAndMarket(ctx, cmd.YearNo, marketCode)
		if err != nil {
			return fmt.Errorf("list market states: %w", err)
		}
		if len(states) == 0 {
			return ErrOrderPoolNotGenerated
		}
		if hasSegmentStatus(states, enum.OrderSegmentStatusBidOpen) {
			return ErrOrderMarketAlreadyOpen
		}
		if hasAnySegmentStatus(states, []string{enum.OrderSegmentStatusBidClosed, enum.OrderSegmentStatusSequenceReady, enum.OrderSegmentStatusSelecting, enum.OrderSegmentStatusCompleted}) {
			return ErrOrderMarketAlreadyClosed
		}
		affected, err := stateRepo.UpdateMarketStatus(ctx, cmd.YearNo, marketCode, []string{enum.OrderSegmentStatusWaitingRelease}, enum.OrderSegmentStatusBidOpen, operatorName, now)
		if err != nil {
			return fmt.Errorf("open market bidding: %w", err)
		}
		if affected == 0 {
			return ErrOrderPoolNotGenerated
		}
		if err := createAdminOrderActionLog(ctx, actionRepo, adminActionCodeOpenOrderMarket, cmd.OperatorID, operatorName, cmd.YearNo, map[string]any{
			"marketCode": marketCode,
			"affected":   affected,
		}, now); err != nil {
			return err
		}
		result = &AdminOrderControlResult{
			YearNo:        cmd.YearNo,
			MarketCode:    marketCode,
			MarketName:    marketName(marketCode),
			AffectedCount: affected,
			OperatedAt:    now.Format(time.RFC3339),
			OperatedBy:    operatorName,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("open market bidding transaction: %w", err)
	}
	return result, nil
}

func (s *AdminOrderControlCommandService) CloseMarketBidding(ctx context.Context, cmd CloseMarketBiddingCommand) (*AdminOrderControlResult, error) {
	marketCode := normalizeMarketCode(cmd.MarketCode)
	if !enum.IsValidMarketCode(marketCode) {
		return nil, ErrOrderMarketInvalid
	}
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	now := time.Now()
	var result *AdminOrderControlResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		bidRepo := repository.NewGroupMarketBidRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		sequenceRepo := repository.NewMarketSelectionOrderRepository(tx)
		poolRepo := repository.NewOrderPoolRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)
		if err := validateFormalOrderYear(ctx, gameConfigRepo, cmd.YearNo); err != nil {
			return err
		}
		states, err := stateRepo.ListByYearAndMarket(ctx, cmd.YearNo, marketCode)
		if err != nil {
			return fmt.Errorf("list market states: %w", err)
		}
		if !hasSegmentStatus(states, enum.OrderSegmentStatusBidOpen) {
			return ErrOrderMarketNotOpen
		}
		groups, err := groupRepo.ListAll(ctx)
		if err != nil {
			return fmt.Errorf("list groups: %w", err)
		}
		bids, err := bidRepo.ListByYearMarket(ctx, cmd.YearNo, marketCode)
		if err != nil {
			return fmt.Errorf("list bids: %w", err)
		}
		previousAmounts, err := poolRepo.SumSelectedAmountByYearMarket(ctx, cmd.YearNo-1, marketCode)
		if err != nil {
			return fmt.Errorf("sum previous market amount: %w", err)
		}
		seed := fmt.Sprintf("%d", now.UnixNano())
		participants, leaderGroupID, rankBasis, err := buildMarketParticipants(cmd.YearNo, marketCode, groups, bids, previousAmounts, seed)
		if err != nil {
			return err
		}
		if err := sequenceRepo.DeleteByYearMarket(ctx, cmd.YearNo, marketCode); err != nil {
			return fmt.Errorf("delete old sequence order: %w", err)
		}
		affected, err := stateRepo.UpdateMarketStatus(ctx, cmd.YearNo, marketCode, []string{enum.OrderSegmentStatusBidOpen}, enum.OrderSegmentStatusBidClosed, operatorName, now)
		if err != nil {
			return fmt.Errorf("close market bidding: %w", err)
		}
		if len(participants) == 0 {
			if _, err := stateRepo.MarkMarketSkipped(ctx, cmd.YearNo, marketCode, operatorName, now); err != nil {
				return fmt.Errorf("skip market segments: %w", err)
			}
		} else {
			sequenceItems := buildSelectionOrderItems(states, participants, operatorName, now)
			if err := sequenceRepo.CreateBatch(ctx, sequenceItems); err != nil {
				return fmt.Errorf("create selection order: %w", err)
			}
			if _, err := stateRepo.MarkMarketSequenceReady(ctx, cmd.YearNo, marketCode, leaderGroupID, rankBasis, seed, operatorName, now); err != nil {
				return fmt.Errorf("mark market sequence ready: %w", err)
			}
		}
		if err := createAdminOrderActionLog(ctx, actionRepo, adminActionCodeCloseOrderMarket, cmd.OperatorID, operatorName, cmd.YearNo, map[string]any{
			"marketCode":       marketCode,
			"affected":         affected,
			"leaderGroupId":    leaderGroupID,
			"participantCount": len(participants),
		}, now); err != nil {
			return err
		}
		result = &AdminOrderControlResult{
			YearNo:        cmd.YearNo,
			MarketCode:    marketCode,
			MarketName:    marketName(marketCode),
			AffectedCount: affected,
			OperatedAt:    now.Format(time.RFC3339),
			OperatedBy:    operatorName,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("close market bidding transaction: %w", err)
	}
	return result, nil
}

func (s *AdminOrderControlCommandService) GenerateSelectionSequence(ctx context.Context, cmd GenerateSelectionSequenceCommand) (*AdminOrderControlResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	now := time.Now()
	var result *AdminOrderControlResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		batchRepo := repository.NewOrderGenerationBatchRepository(tx)
		bidRepo := repository.NewGroupMarketBidRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		sequenceRepo := repository.NewMarketSelectionOrderRepository(tx)
		poolRepo := repository.NewOrderPoolRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)
		if err := validateFormalOrderYear(ctx, gameConfigRepo, cmd.YearNo); err != nil {
			return err
		}
		if _, err := batchRepo.FindConfirmedByYear(ctx, cmd.YearNo); err != nil {
			if repository.IsRecordNotFound(err) {
				return ErrAdminOrderPoolNotConfirmed
			}
			return fmt.Errorf("load confirmed order batch: %w", err)
		}
		states, err := stateRepo.ListByYear(ctx, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("list order states: %w", err)
		}
		if len(states) == 0 {
			return ErrOrderPoolNotGenerated
		}
		if hasAnySegmentStatus(states, []string{enum.OrderSegmentStatusSequenceReady, enum.OrderSegmentStatusSelecting, enum.OrderSegmentStatusCompleted}) {
			return ErrAdminOrderSequenceAlreadyGenerated
		}
		if err := ensureAllActiveGroupsSubmittedInvestments(ctx, groupRepo, bidRepo, cmd.YearNo); err != nil {
			return err
		}
		groups, err := groupRepo.ListAll(ctx)
		if err != nil {
			return fmt.Errorf("list groups: %w", err)
		}
		if err := sequenceRepo.DeleteByYear(ctx, cmd.YearNo); err != nil {
			return fmt.Errorf("delete old sequence order: %w", err)
		}
		leaderResolved := make(map[string]bool)
		leaderCache := make(map[string]*int64)
		previousAmountCache := make(map[string]map[int64]float64)
		readyCount := int64(0)
		skippedCount := int64(0)
		for _, item := range states {
			if item.SegmentStatus == enum.OrderSegmentStatusSkipped ||
				item.SegmentStatus == enum.OrderSegmentStatusMarketDisabled ||
				item.SegmentStatus == enum.OrderSegmentStatusNoOrderConfig {
				if item.SegmentStatus == enum.OrderSegmentStatusSkipped {
					skippedCount++
				}
				continue
			}
			if item.SegmentStatus != enum.OrderSegmentStatusWaitingInvestment && item.SegmentStatus != enum.OrderSegmentStatusWaitingRelease && item.SegmentStatus != enum.OrderSegmentStatusBidClosed {
				return ErrAdminOrderSequenceAlreadyGenerated
			}
			availableCount, err := poolRepo.CountAvailableBySegment(ctx, item.YearNo, item.MarketCode, item.OrderType)
			if err != nil {
				return fmt.Errorf("count available order: %w", err)
			}
			seed := fmt.Sprintf("%d-%s-%s-%d", now.UnixNano(), item.MarketCode, item.OrderType, item.ID)
			bids, err := bidRepo.ListByYearSegment(ctx, cmd.YearNo, item.MarketCode, item.OrderType)
			if err != nil {
				return fmt.Errorf("list segment bids: %w", err)
			}
			previousAmounts, exists := previousAmountCache[item.MarketCode]
			if !exists {
				previousAmounts, err = poolRepo.SumSelectedAmountByYearMarket(ctx, cmd.YearNo-1, item.MarketCode)
				if err != nil {
					return fmt.Errorf("sum previous market amount: %w", err)
				}
				previousAmountCache[item.MarketCode] = previousAmounts
			}
			if !leaderResolved[item.MarketCode] {
				leaderSeed := fmt.Sprintf("%d-%s-leader", now.UnixNano(), item.MarketCode)
				leaderRNG := rand.New(rand.NewSource(parseSeed(leaderSeed)))
				leaderCache[item.MarketCode] = resolveMarketLeader(cmd.YearNo, groups, previousAmounts, leaderRNG)
				leaderResolved[item.MarketCode] = true
			}
			leaderGroupID := leaderCache[item.MarketCode]
			participants, rankBasis, err := buildMarketParticipantsWithLeader(cmd.YearNo, item.MarketCode, groups, bids, previousAmounts, seed, leaderGroupID)
			if err != nil {
				return err
			}
			if availableCount <= 0 || !hasPositiveSegmentInvestment(groups, bids) || len(participants) == 0 {
				if err := stateRepo.MarkSegmentSkipped(ctx, item.ID, leaderGroupID, rankBasis, seed, operatorName, now); err != nil {
					return fmt.Errorf("skip segment: %w", err)
				}
				skippedCount++
				continue
			}
			sequenceItems := buildSelectionOrderItems([]entity.MarketBiddingState{item}, participants, operatorName, now)
			if err := sequenceRepo.CreateBatch(ctx, sequenceItems); err != nil {
				return fmt.Errorf("create selection order: %w", err)
			}
			if err := stateRepo.MarkSegmentSequenceReady(ctx, item.ID, leaderGroupID, rankBasis, seed, operatorName, now); err != nil {
				return fmt.Errorf("mark segment sequence ready: %w", err)
			}
			readyCount++
		}
		if err := createAdminOrderActionLog(ctx, actionRepo, adminActionCodeGenerateSelectionSequence, cmd.OperatorID, operatorName, cmd.YearNo, map[string]any{
			"readyCount":   readyCount,
			"skippedCount": skippedCount,
		}, now); err != nil {
			return err
		}
		result = &AdminOrderControlResult{
			YearNo:        cmd.YearNo,
			AffectedCount: readyCount + skippedCount,
			OperatedAt:    now.Format(time.RFC3339),
			OperatedBy:    operatorName,
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("generate selection sequence transaction: %w", err)
	}
	return result, nil
}

func (s *AdminOrderControlCommandService) ReleaseNextSegment(ctx context.Context, cmd ReleaseNextSegmentCommand) (*AdminOrderControlResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	now := time.Now()
	var result *AdminOrderControlResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		sequenceRepo := repository.NewMarketSelectionOrderRepository(tx)
		poolRepo := repository.NewOrderPoolRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)
		if err := validateFormalOrderYear(ctx, gameConfigRepo, cmd.YearNo); err != nil {
			return err
		}
		if _, err := stateRepo.GetCurrentSelecting(ctx, cmd.YearNo); err == nil {
			return ErrOrderSegmentReleaseBlocked
		} else if !repository.IsRecordNotFound(err) {
			return fmt.Errorf("load current selecting segment: %w", err)
		}
		next, err := stateRepo.GetNextReleasable(ctx, cmd.YearNo)
		if err != nil {
			if repository.IsRecordNotFound(err) {
				return ErrOrderSegmentNotReady
			}
			return fmt.Errorf("load next releasable segment: %w", err)
		}
		blocked, err := stateRepo.ExistsSelectingBefore(ctx, cmd.YearNo, next.ReleaseSequenceNo)
		if err != nil {
			return fmt.Errorf("check prior segments: %w", err)
		}
		if blocked {
			return ErrOrderSegmentReleaseBlocked
		}
		availableCount, err := poolRepo.CountAvailableBySegment(ctx, next.YearNo, next.MarketCode, next.OrderType)
		if err != nil {
			return fmt.Errorf("count available order: %w", err)
		}
		if availableCount <= 0 {
			if err := stateRepo.CompleteSegment(ctx, next.ID, operatorName, now); err != nil {
				return fmt.Errorf("complete empty segment: %w", err)
			}
			result = buildAdminOrderControlResultFromState(*next, nil, enum.OrderSegmentStatusCompleted, 1, operatorName, now)
			return nil
		}
		first, err := sequenceRepo.GetNextWaiting(ctx, next.YearNo, next.MarketCode, next.OrderType)
		if err != nil {
			if repository.IsRecordNotFound(err) {
				if err := stateRepo.CompleteSegment(ctx, next.ID, operatorName, now); err != nil {
					return fmt.Errorf("complete no sequence segment: %w", err)
				}
				result = buildAdminOrderControlResultFromState(*next, nil, enum.OrderSegmentStatusCompleted, 1, operatorName, now)
				return nil
			}
			return fmt.Errorf("load first waiting group: %w", err)
		}
		if err := sequenceRepo.SetCurrent(ctx, first.ID, operatorName, now); err != nil {
			return fmt.Errorf("set first current: %w", err)
		}
		currentGroupID := first.GroupID
		if err := stateRepo.ReleaseSegment(ctx, next.ID, &currentGroupID, operatorName, now); err != nil {
			return fmt.Errorf("release segment: %w", err)
		}
		if err := createAdminOrderActionLog(ctx, actionRepo, adminActionCodeReleaseOrderSegment, cmd.OperatorID, operatorName, cmd.YearNo, map[string]any{
			"marketCode":        next.MarketCode,
			"orderType":         next.OrderType,
			"releaseSequenceNo": next.ReleaseSequenceNo,
			"currentGroupId":    currentGroupID,
		}, now); err != nil {
			return err
		}
		result = buildAdminOrderControlResultFromState(*next, &currentGroupID, enum.OrderSegmentStatusSelecting, 1, operatorName, now)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("release next segment transaction: %w", err)
	}
	return result, nil
}

func (s *AdminOrderControlCommandService) AdminSkipCurrentGroup(ctx context.Context, cmd AdminSkipCurrentGroupCommand) (*AdminOrderControlResult, error) {
	marketCode := normalizeMarketCode(cmd.MarketCode)
	orderType := normalizeOrderType(cmd.OrderType)
	reason := strings.TrimSpace(cmd.Reason)
	if !enum.IsValidMarketCode(marketCode) {
		return nil, ErrOrderMarketInvalid
	}
	if !enum.IsValidOrderType(orderType) {
		return nil, ErrOrderTypeInvalid
	}
	if cmd.GroupID <= 0 || reason == "" {
		return nil, ErrOrderAdminSkipReasonRequired
	}
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	now := time.Now()
	var result *AdminOrderControlResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		stateRepo := repository.NewMarketBiddingStateRepository(tx)
		sequenceRepo := repository.NewMarketSelectionOrderRepository(tx)
		poolRepo := repository.NewOrderPoolRepository(tx)
		actionRepo := repository.NewAdminActionLogRepository(tx)
		if err := validateFormalOrderYear(ctx, gameConfigRepo, cmd.YearNo); err != nil {
			return err
		}
		segment, err := stateRepo.GetBySegmentForUpdate(ctx, cmd.YearNo, marketCode, orderType)
		if err != nil {
			return fmt.Errorf("load segment: %w", err)
		}
		if segment.SegmentStatus != enum.OrderSegmentStatusSelecting {
			return ErrOrderSegmentNotSelecting
		}
		if segment.CurrentGroupID == nil || *segment.CurrentGroupID != cmd.GroupID {
			return ErrOrderSegmentCurrentGroupMismatch
		}
		sequence, err := sequenceRepo.GetByGroupSegmentForUpdate(ctx, cmd.GroupID, cmd.YearNo, marketCode, orderType)
		if err != nil {
			return fmt.Errorf("load sequence order: %w", err)
		}
		if sequence.SelectionStatus != enum.OrderSelectionStatusCurrent {
			return ErrOrderSegmentCurrentGroupMismatch
		}
		if err := sequenceRepo.SetStatus(ctx, sequence.ID, enum.OrderSelectionStatusAdminSkipped, nil, &cmd.OperatorID, &reason, &now, operatorName); err != nil {
			return fmt.Errorf("mark admin skipped: %w", err)
		}
		nextGroupID, status, err := advanceOrderSegment(ctx, stateRepo, sequenceRepo, poolRepo, *segment, operatorName, now)
		if err != nil {
			return err
		}
		if err := createAdminOrderActionLog(ctx, actionRepo, adminActionCodeSkipOrderCurrentGroup, cmd.OperatorID, operatorName, cmd.YearNo, map[string]any{
			"marketCode": marketCode,
			"orderType":  orderType,
			"groupId":    cmd.GroupID,
			"reason":     reason,
		}, now); err != nil {
			return err
		}
		result = buildAdminOrderControlResultFromState(*segment, nextGroupID, status, 1, operatorName, now)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("admin skip current group transaction: %w", err)
	}
	return result, nil
}

type marketParticipant struct {
	GroupID                   int64   `json:"groupId"`
	GroupName                 string  `json:"groupName"`
	MarketInvestment          float64 `json:"marketInvestment"`
	PreviousMarketOrderAmount float64 `json:"previousMarketOrderAmount"`
	IsMarketLeader            bool    `json:"isMarketLeader"`
	RandomRank                int     `json:"randomRank"`
}

func buildMarketParticipants(yearNo int, marketCode string, groups []entity.Group, bids []entity.GroupMarketBid, previousAmounts map[int64]float64, seed string) ([]marketParticipant, *int64, []byte, error) {
	rng := rand.New(rand.NewSource(parseSeed(seed)))
	leaderGroupID := resolveMarketLeader(yearNo, groups, previousAmounts, rng)
	participants, rankBasis, err := buildMarketParticipantsWithLeader(yearNo, marketCode, groups, bids, previousAmounts, seed, leaderGroupID)
	return participants, leaderGroupID, rankBasis, err
}

func buildMarketParticipantsWithLeader(yearNo int, marketCode string, groups []entity.Group, bids []entity.GroupMarketBid, previousAmounts map[int64]float64, seed string, leaderGroupID *int64) ([]marketParticipant, []byte, error) {
	bidMap := make(map[int64]float64, len(bids))
	bidSubmitted := make(map[int64]bool, len(bids))
	for _, bid := range bids {
		bidMap[bid.GroupID] = bid.MarketInvestment
		bidSubmitted[bid.GroupID] = true
	}
	rng := rand.New(rand.NewSource(parseSeed(seed)))
	participants := make([]marketParticipant, 0, len(groups))
	for _, group := range groups {
		if group.BusinessStatus == enum.BusinessStatusBankrupt {
			continue
		}
		investment := bidMap[group.ID]
		isLeader := leaderGroupID != nil && *leaderGroupID == group.ID
		if !bidSubmitted[group.ID] {
			continue
		}
		if investment <= 0 && !isLeader {
			continue
		}
		participants = append(participants, marketParticipant{
			GroupID:                   group.ID,
			GroupName:                 group.GroupName,
			MarketInvestment:          investment,
			PreviousMarketOrderAmount: previousAmounts[group.ID],
			IsMarketLeader:            isLeader,
			RandomRank:                rng.Int(),
		})
	}
	sort.SliceStable(participants, func(i, j int) bool {
		left := participants[i]
		right := participants[j]
		if left.IsMarketLeader != right.IsMarketLeader {
			return left.IsMarketLeader
		}
		if left.MarketInvestment != right.MarketInvestment {
			return left.MarketInvestment > right.MarketInvestment
		}
		if yearNo >= 2 && left.PreviousMarketOrderAmount != right.PreviousMarketOrderAmount {
			return left.PreviousMarketOrderAmount > right.PreviousMarketOrderAmount
		}
		return left.RandomRank < right.RandomRank
	})
	rankBasis, err := marshalJSON(map[string]any{
		"yearNo":        yearNo,
		"marketCode":    marketCode,
		"seed":          seed,
		"leaderGroupId": leaderGroupID,
		"participants":  participants,
	})
	if err != nil {
		return nil, nil, err
	}
	return participants, rankBasis, nil
}

func resolveMarketLeader(yearNo int, groups []entity.Group, previousAmounts map[int64]float64, rng *rand.Rand) *int64 {
	if yearNo < 2 {
		return nil
	}
	maxAmount := 0.0
	candidates := make([]int64, 0)
	groupExists := make(map[int64]bool, len(groups))
	for _, group := range groups {
		if group.BusinessStatus != enum.BusinessStatusBankrupt {
			groupExists[group.ID] = true
		}
	}
	for groupID, amount := range previousAmounts {
		if !groupExists[groupID] || amount <= 0 {
			continue
		}
		if amount > maxAmount {
			maxAmount = amount
			candidates = []int64{groupID}
			continue
		}
		if amount == maxAmount {
			candidates = append(candidates, groupID)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i] < candidates[j] })
	selected := candidates[rng.Intn(len(candidates))]
	return &selected
}

func buildSelectionOrderItems(states []entity.MarketBiddingState, participants []marketParticipant, operatorName string, now time.Time) []entity.MarketSelectionOrder {
	items := make([]entity.MarketSelectionOrder, 0, len(states)*len(participants))
	for _, state := range states {
		if isTerminalOrderSegmentStatus(state.SegmentStatus) {
			continue
		}
		for index, participant := range participants {
			rankBasis, _ := marshalJSON(participant)
			items = append(items, entity.MarketSelectionOrder{
				YearNo:                    state.YearNo,
				MarketCode:                state.MarketCode,
				OrderType:                 state.OrderType,
				SequenceNo:                index + 1,
				GroupID:                   participant.GroupID,
				MarketInvestment:          participant.MarketInvestment,
				PreviousMarketOrderAmount: participant.PreviousMarketOrderAmount,
				IsMarketLeader:            participant.IsMarketLeader,
				RankBasis:                 rankBasis,
				SelectionStatus:           enum.OrderSelectionStatusWaiting,
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

func advanceOrderSegment(ctx context.Context, stateRepo *repository.MarketBiddingStateRepository, sequenceRepo *repository.MarketSelectionOrderRepository, poolRepo *repository.OrderPoolRepository, segment entity.MarketBiddingState, operatorName string, now time.Time) (*int64, string, error) {
	availableCount, err := poolRepo.CountAvailableBySegment(ctx, segment.YearNo, segment.MarketCode, segment.OrderType)
	if err != nil {
		return nil, "", fmt.Errorf("count available order: %w", err)
	}
	if availableCount <= 0 {
		if err := stateRepo.CompleteSegment(ctx, segment.ID, operatorName, now); err != nil {
			return nil, "", fmt.Errorf("complete empty segment: %w", err)
		}
		return nil, enum.OrderSegmentStatusCompleted, nil
	}
	next, err := sequenceRepo.GetNextWaiting(ctx, segment.YearNo, segment.MarketCode, segment.OrderType)
	if err != nil {
		if repository.IsRecordNotFound(err) {
			if err := stateRepo.CompleteSegment(ctx, segment.ID, operatorName, now); err != nil {
				return nil, "", fmt.Errorf("complete segment: %w", err)
			}
			return nil, enum.OrderSegmentStatusCompleted, nil
		}
		return nil, "", fmt.Errorf("load next waiting group: %w", err)
	}
	if err := sequenceRepo.SetCurrent(ctx, next.ID, operatorName, now); err != nil {
		return nil, "", fmt.Errorf("set next current: %w", err)
	}
	nextGroupID := next.GroupID
	if err := stateRepo.UpdateCurrentGroup(ctx, segment.ID, &nextGroupID, operatorName, now); err != nil {
		return nil, "", fmt.Errorf("update current group: %w", err)
	}
	return &nextGroupID, enum.OrderSegmentStatusSelecting, nil
}

func normalizeMarketInvestmentInputs(inputs []MarketInvestmentInput) ([]MarketInvestmentInput, error) {
	if len(inputs) != requiredOrderInvestmentSegmentCount {
		return nil, ErrOrderInvestmentInvalid
	}
	seen := make(map[string]bool, requiredOrderInvestmentSegmentCount)
	result := make([]MarketInvestmentInput, 0, requiredOrderInvestmentSegmentCount)
	for _, input := range inputs {
		marketCode := normalizeMarketCode(input.MarketCode)
		orderType := normalizeOrderType(input.OrderType)
		if !enum.IsValidMarketCode(marketCode) || !enum.IsValidOrderType(orderType) || input.MarketInvestment < 0 {
			return nil, ErrOrderInvestmentInvalid
		}
		key := segmentKey(0, marketCode, orderType)
		if seen[key] {
			return nil, ErrOrderInvestmentInvalid
		}
		seen[key] = true
		result = append(result, MarketInvestmentInput{
			MarketCode:       marketCode,
			OrderType:        orderType,
			MarketInvestment: input.MarketInvestment,
		})
	}
	for _, market := range defaultOrderMarkets() {
		for _, orderType := range defaultOrderTypes() {
			if !seen[segmentKey(0, market.code, orderType.code)] {
				return nil, ErrOrderInvestmentInvalid
			}
		}
	}
	return result, nil
}

func validateDisabledMarketInvestmentsAreZero(inputs []MarketInvestmentInput, marketConfigMap map[string]bool) error {
	for _, input := range inputs {
		if !isMarketEnabled(marketConfigMap, input.MarketCode) && input.MarketInvestment != 0 {
			return ErrOrderMarketDisabledInvestment
		}
	}
	return nil
}

func ensureAllActiveGroupsSubmittedInvestments(ctx context.Context, groupRepo *repository.GroupRepository, bidRepo *repository.GroupMarketBidRepository, yearNo int) error {
	groups, err := groupRepo.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("list groups: %w", err)
	}
	counts, err := bidRepo.CountSubmittedSegmentsByYear(ctx, yearNo)
	if err != nil {
		return fmt.Errorf("count submitted investments: %w", err)
	}
	for _, group := range groups {
		if group.BusinessStatus == enum.BusinessStatusBankrupt {
			continue
		}
		if counts[group.ID] < requiredOrderInvestmentSegmentCount {
			return ErrAdminOrderInvestmentIncomplete
		}
	}
	return nil
}

func hasPositiveSegmentInvestment(groups []entity.Group, bids []entity.GroupMarketBid) bool {
	activeGroups := make(map[int64]bool, len(groups))
	for _, group := range groups {
		if group.BusinessStatus != enum.BusinessStatusBankrupt {
			activeGroups[group.ID] = true
		}
	}
	for _, bid := range bids {
		if activeGroups[bid.GroupID] && bid.MarketInvestment > 0 {
			return true
		}
	}
	return false
}

func buildPlayerOrderYearView(groupID int64, yearNo int, states []entity.MarketBiddingState, bids map[string]entity.GroupMarketBid, sequences []entity.MarketSelectionOrder, selected []entity.GroupOrderSelection, groups []entity.Group, poolConfirmed bool, poolRepo *repository.OrderPoolRepository, ctx context.Context) (*PlayerOrderYearView, error) {
	groupNames := buildGroupNameMap(groups)
	investmentSubmitted := len(bids) >= requiredOrderInvestmentSegmentCount
	canSubmitInvestment := poolConfirmed && !investmentSubmitted
	sequenceMap := make(map[string][]entity.MarketSelectionOrder)
	selfSequence := make(map[string]entity.MarketSelectionOrder)
	for _, sequence := range sequences {
		key := segmentKey(sequence.YearNo, sequence.MarketCode, sequence.OrderType)
		sequenceMap[key] = append(sequenceMap[key], sequence)
		if sequence.GroupID == groupID {
			selfSequence[key] = sequence
		}
	}
	selectedMap := make(map[string]entity.GroupOrderSelection, len(selected))
	for _, item := range selected {
		selectedMap[segmentKey(item.YearNo, item.MarketCode, item.OrderType)] = item
	}
	stateMap := make(map[string][]entity.MarketBiddingState)
	for _, state := range states {
		stateMap[state.MarketCode] = append(stateMap[state.MarketCode], state)
	}
	markets := make([]PlayerOrderMarketView, 0, 4)
	for _, market := range defaultOrderMarkets() {
		segments := make([]PlayerOrderSegmentView, 0)
		canSelectMarket := false
		var selfMarketSequenceNo *int
		isMarketLeader := false
		marketInvestment := 0.0
		marketSubmittedCount := 0
		marketEnabled := true
		if len(stateMap[market.code]) > 0 {
			marketEnabled = !allMarketStatesDisabled(stateMap[market.code])
		}
		for _, state := range stateMap[market.code] {
			key := segmentKey(state.YearNo, state.MarketCode, state.OrderType)
			bid, bidSubmitted := bids[key]
			if bidSubmitted {
				marketSubmittedCount++
				marketInvestment += bid.MarketInvestment
			}
			selfSeq, hasSelfSeq := selfSequence[key]
			if hasSelfSeq {
				seqNo := selfSeq.SequenceNo
				selfMarketSequenceNo = &seqNo
				isMarketLeader = isMarketLeader || selfSeq.IsMarketLeader
			}
			pools, err := poolRepo.ListBySegment(ctx, state.YearNo, state.MarketCode, state.OrderType)
			if err != nil {
				return nil, err
			}
			segment := buildPlayerSegmentView(groupID, state, bid, bidSubmitted, sequenceMap[key], selfSeq, hasSelfSeq, selectedMap[key], pools, groupNames)
			canSelectMarket = canSelectMarket || segment.CanSelectOrder
			segments = append(segments, segment)
		}
		sort.SliceStable(segments, func(i, j int) bool {
			return segments[i].ReleaseSequenceNo < segments[j].ReleaseSequenceNo
		})
		markets = append(markets, PlayerOrderMarketView{
			MarketCode:          market.code,
			MarketName:          market.name,
			MarketBidStatus:     resolveMarketBidStatus(stateMap[market.code]),
			InvestmentSubmitted: marketSubmittedCount >= len(defaultOrderTypes()),
			MarketInvestment:    marketInvestment,
			CanSubmitInvestment: canSubmitInvestment,
			CanSelectOrder:      canSelectMarket,
			SelectionSequenceNo: selfMarketSequenceNo,
			IsMarketLeader:      isMarketLeader,
			MarketEnabled:       marketEnabled,
			Segments:            segments,
		})
	}
	return &PlayerOrderYearView{
		GroupID:                groupID,
		YearNo:                 yearNo,
		OrderRequired:          true,
		InvestmentSubmitted:    investmentSubmitted,
		CanSubmitInvestment:    canSubmitInvestment,
		Markets:                markets,
		PollingIntervalSeconds: orderPollingIntervalSeconds,
	}, nil
}

func allMarketStatesDisabled(states []entity.MarketBiddingState) bool {
	if len(states) == 0 {
		return false
	}
	for _, state := range states {
		if state.SegmentStatus != enum.OrderSegmentStatusMarketDisabled {
			return false
		}
	}
	return true
}

func buildPlayerSegmentView(groupID int64, state entity.MarketBiddingState, bid entity.GroupMarketBid, bidSubmitted bool, sequences []entity.MarketSelectionOrder, selfSequence entity.MarketSelectionOrder, hasSelfSequence bool, selected entity.GroupOrderSelection, pools []entity.OrderPool, groupNames map[int64]string) PlayerOrderSegmentView {
	available := make([]PlayerOrderPoolItem, 0)
	locked := make([]PlayerOrderPoolItem, 0)
	var selectedOrder *PlayerOrderPoolItem
	for _, pool := range pools {
		item := PlayerOrderPoolItem{
			OrderID:         pool.ID,
			BusinessOrderNo: formatBusinessOrderNo(pool),
			CardSequenceNo:  pool.CardSequenceNo,
			OrderAmount:     pool.OrderAmount,
			OrderQuantity:   pool.OrderQuantity,
			UnitPrice:       pool.UnitPrice,
			AccountTerm:     pool.AccountTerm,
			PoolStatus:      pool.PoolStatus,
		}
		if pool.PoolStatus == enum.OrderPoolStatusAvailable {
			available = append(available, item)
		} else {
			locked = append(locked, item)
		}
		if selected.OrderID == pool.ID {
			copyItem := item
			copyItem.DeliveryStatus = selected.DeliveryStatus
			copyItem.DeliveredStageCode = selected.DeliveredStageCode
			selectedOrder = &copyItem
		}
	}
	sequenceViews := make([]PlayerOrderSequenceView, 0, len(sequences))
	for _, sequence := range sequences {
		sequenceViews = append(sequenceViews, PlayerOrderSequenceView{
			SequenceNo:      sequence.SequenceNo,
			GroupID:         sequence.GroupID,
			GroupName:       groupNames[sequence.GroupID],
			SelectionStatus: sequence.SelectionStatus,
			IsSelf:          sequence.GroupID == groupID,
			IsMarketLeader:  sequence.IsMarketLeader,
		})
	}
	deliveryStatus := ""
	if selected.ID > 0 {
		deliveryStatus = selected.DeliveryStatus
	}
	canAct := state.SegmentStatus == enum.OrderSegmentStatusSelecting && state.CurrentGroupID != nil && *state.CurrentGroupID == groupID && hasSelfSequence && selfSequence.SelectionStatus == enum.OrderSelectionStatusCurrent
	return PlayerOrderSegmentView{
		MarketCode:          state.MarketCode,
		MarketName:          marketName(state.MarketCode),
		OrderType:           state.OrderType,
		OrderTypeName:       orderTypeName(state.OrderType),
		MarketEnabled:       state.SegmentStatus != enum.OrderSegmentStatusMarketDisabled,
		MarketInvestment:    bid.MarketInvestment,
		InvestmentSubmitted: bidSubmitted,
		ReleaseSequenceNo:   state.ReleaseSequenceNo,
		SegmentStatus:       state.SegmentStatus,
		SelectionOrder:      sequenceViews,
		CurrentGroupID:      state.CurrentGroupID,
		AvailableOrders:     available,
		LockedOrders:        locked,
		SelectedOrder:       selectedOrder,
		DeliveryStatus:      deliveryStatus,
		CanSelectOrder:      canAct && len(available) > 0,
		CanPassSegment:      canAct,
	}
}

func buildAdminSegmentStatus(ctx context.Context, sequenceRepo *repository.MarketSelectionOrderRepository, poolRepo *repository.OrderPoolRepository, state entity.MarketBiddingState, groupNames map[int64]string) (AdminOrderSegmentStatus, error) {
	sequences, err := sequenceRepo.ListBySegment(ctx, state.YearNo, state.MarketCode, state.OrderType)
	if err != nil {
		return AdminOrderSegmentStatus{}, fmt.Errorf("list segment sequence: %w", err)
	}
	pools, err := poolRepo.ListBySegment(ctx, state.YearNo, state.MarketCode, state.OrderType)
	if err != nil {
		return AdminOrderSegmentStatus{}, fmt.Errorf("list segment pool: %w", err)
	}
	availableCount := 0
	selectedCount := 0
	for _, pool := range pools {
		if pool.PoolStatus == enum.OrderPoolStatusAvailable {
			availableCount++
		}
		if pool.PoolStatus == enum.OrderPoolStatusSelected {
			selectedCount++
		}
	}
	sequenceViews := make([]AdminSelectionOrderView, 0, len(sequences))
	orderNoMap := make(map[int64]string, len(pools))
	for _, pool := range pools {
		orderNoMap[pool.ID] = formatBusinessOrderNo(pool)
	}
	for _, sequence := range sequences {
		selectedOrderNo := ""
		if sequence.SelectedOrderID != nil {
			selectedOrderNo = orderNoMap[*sequence.SelectedOrderID]
		}
		sequenceViews = append(sequenceViews, AdminSelectionOrderView{
			SequenceNo:       sequence.SequenceNo,
			GroupID:          sequence.GroupID,
			GroupName:        groupNames[sequence.GroupID],
			MarketInvestment: sequence.MarketInvestment,
			IsMarketLeader:   sequence.IsMarketLeader,
			SelectionStatus:  sequence.SelectionStatus,
			SelectedOrderID:  sequence.SelectedOrderID,
			SelectedOrderNo:  selectedOrderNo,
		})
	}
	return AdminOrderSegmentStatus{
		MarketCode:        state.MarketCode,
		MarketName:        marketName(state.MarketCode),
		OrderType:         state.OrderType,
		OrderTypeName:     orderTypeName(state.OrderType),
		ReleaseSequenceNo: state.ReleaseSequenceNo,
		SegmentStatus:     state.SegmentStatus,
		CurrentGroupID:    state.CurrentGroupID,
		SelectionOrder:    sequenceViews,
		AvailableCount:    availableCount,
		SelectedCount:     selectedCount,
	}, nil
}

func buildAdminOrderControlResultFromState(state entity.MarketBiddingState, currentGroupID *int64, status string, affected int64, operatorName string, now time.Time) *AdminOrderControlResult {
	return &AdminOrderControlResult{
		YearNo:         state.YearNo,
		MarketCode:     state.MarketCode,
		MarketName:     marketName(state.MarketCode),
		OrderType:      state.OrderType,
		OrderTypeName:  orderTypeName(state.OrderType),
		SegmentStatus:  status,
		CurrentGroupID: currentGroupID,
		AffectedCount:  affected,
		OperatedAt:     now.Format(time.RFC3339),
		OperatedBy:     operatorName,
	}
}

func validateFormalOrderYear(ctx context.Context, repo *repository.GameConfigRepository, yearNo int) error {
	if yearNo == 0 {
		return ErrOrderNotRequiredForDemoYear
	}
	if yearNo < 1 {
		return ErrOrderYearInvalid
	}
	config, err := repo.GetCurrent(ctx)
	if err != nil {
		return fmt.Errorf("load game config: %w", err)
	}
	if yearNo > config.FinalYear {
		return ErrOrderYearInvalid
	}
	return nil
}

func validateSelectableGroup(ctx context.Context, repo *repository.GroupRepository, groupID int64) error {
	group, err := repo.GetByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("load group: %w", err)
	}
	if group.BusinessStatus == enum.BusinessStatusBankrupt {
		return ErrOrderGroupUnavailable
	}
	return nil
}

func normalizeOrderStageCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func isOrderQuarterStageCode(value string) bool {
	switch value {
	case state.StageCodeQ1, state.StageCodeQ2, state.StageCodeQ3, state.StageCodeQ4:
		return true
	default:
		return false
	}
}

func uniquePositiveOrderIDs(source []int64) []int64 {
	seen := map[int64]bool{}
	result := make([]int64, 0, len(source))
	for _, id := range source {
		if id <= 0 || seen[id] {
			continue
		}
		seen[id] = true
		result = append(result, id)
	}
	return result
}

func extractOperatingQuarterSalesRevenue(value payload.OperatingPayload, stageCode string) float64 {
	quarterKey := strings.ToLower(stageCode)
	for key, record := range value.Quarter.DeliverySettlement {
		if strings.EqualFold(key, quarterKey) {
			return sumMatchedDeliveryNumbers(record, "salesRevenue", "deliverySalesRevenue", "orderSalesRevenue", "o38")
		}
	}
	return 0
}

func sumMatchedDeliveryNumbers(source map[string]any, keys ...string) float64 {
	if len(source) == 0 {
		return 0
	}
	normalizedKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		normalizedKeys = append(normalizedKeys, normalizeSearchableOrderKey(key))
	}
	total := 0.0
	for key, value := range source {
		if !searchableOrderKeyMatches(key, normalizedKeys) {
			continue
		}
		total += numericOrderValue(value)
	}
	return total
}

func searchableOrderKeyMatches(key string, normalizedKeys []string) bool {
	normalized := normalizeSearchableOrderKey(key)
	if normalized == "" {
		return false
	}
	for _, item := range normalizedKeys {
		if normalized == item || strings.Contains(normalized, item) || strings.Contains(item, normalized) {
			return true
		}
	}
	return false
}

func normalizeSearchableOrderKey(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	for _, ch := range value {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}

func numericOrderValue(value any) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case float32:
		return float64(typed)
	case int:
		return float64(typed)
	case int64:
		return float64(typed)
	case int32:
		return float64(typed)
	case json.Number:
		result, _ := typed.Float64()
		return result
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		if err != nil {
			return 0
		}
		return parsed
	default:
		return 0
	}
}

func sameMoneyAmount(left float64, right float64) bool {
	return math.Abs(left-right) < 0.000001
}

func defaultOrderMarkets() []struct {
	code string
	name string
} {
	return []struct {
		code string
		name string
	}{
		{enum.MarketCodeLocal, "本地市场"},
		{enum.MarketCodeRegional, "区域市场"},
		{enum.MarketCodeNational, "全国市场"},
		{enum.MarketCodeGlobal, "全球市场"},
	}
}

func defaultOrderTypes() []struct {
	code string
	name string
} {
	return []struct {
		code string
		name string
	}{
		{enum.OrderTypeAgencyInspection, "代办过检"},
		{enum.OrderTypeTwoCabinVIP, "两舱贵宾"},
		{enum.OrderTypeBusinessVIP, "商务贵宾"},
		{enum.OrderTypeMemberCustom, "会员定制"},
	}
}

func buildGroupNameMap(groups []entity.Group) map[int64]string {
	result := make(map[int64]string, len(groups))
	for _, group := range groups {
		if strings.TrimSpace(group.GroupName) != "" {
			result[group.ID] = group.GroupName
		} else {
			result[group.ID] = fmt.Sprintf("第%d组", group.GroupNo)
		}
	}
	return result
}

func resolveMarketBidStatus(states []entity.MarketBiddingState) string {
	if allMarketStatesDisabled(states) {
		return enum.OrderSegmentStatusMarketDisabled
	}
	switch {
	case hasSegmentStatus(states, enum.OrderSegmentStatusBidOpen):
		return enum.OrderSegmentStatusBidOpen
	case hasSegmentStatus(states, enum.OrderSegmentStatusBidClosed):
		return enum.OrderSegmentStatusBidClosed
	case hasSegmentStatus(states, enum.OrderSegmentStatusSequenceReady),
		hasSegmentStatus(states, enum.OrderSegmentStatusWaitingRelease),
		hasSegmentStatus(states, enum.OrderSegmentStatusSelecting):
		return enum.OrderSegmentStatusSequenceReady
	case len(states) > 0:
		return states[0].SegmentStatus
	default:
		return ""
	}
}

func hasSegmentStatus(states []entity.MarketBiddingState, status string) bool {
	for _, state := range states {
		if state.SegmentStatus == status {
			return true
		}
	}
	return false
}

func hasAnySegmentStatus(states []entity.MarketBiddingState, statuses []string) bool {
	statusMap := make(map[string]bool, len(statuses))
	for _, status := range statuses {
		statusMap[status] = true
	}
	for _, state := range states {
		if statusMap[state.SegmentStatus] {
			return true
		}
	}
	return false
}

func normalizeMarketCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func normalizeOrderType(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func parseSeed(seed string) int64 {
	var value int64
	for _, ch := range seed {
		value = value*31 + int64(ch)
	}
	if value == 0 {
		return time.Now().UnixNano()
	}
	return value
}

func createAdminOrderActionLog(ctx context.Context, repo *repository.AdminActionLogRepository, actionCode string, operatorID int64, operatorName string, yearNo int, payload map[string]any, operateTime time.Time) error {
	targetYearNo := yearNo
	rawPayload, err := marshalJSON(payload)
	if err != nil {
		return err
	}
	return repo.Create(ctx, &entity.AdminActionLog{
		ActionCode:    actionCode,
		TargetYearNo:  &targetYearNo,
		ActionPayload: rawPayload,
		StateBefore:   []byte("{}"),
		StateAfter:    []byte("{}"),
		OperatorID:    operatorID,
		OperatorName:  operatorName,
		OperateTime:   operateTime,
	})
}
