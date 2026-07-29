package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

const rollbackSnapshotPayloadVersion = "I7_V1"

type SnapshotOperator struct {
	OperatorID   int64
	OperatorName string
}

type CreateGroupSnapshotCommand struct {
	GroupID       int64
	YearNo        int
	StageCode     string
	SnapshotType  string
	TriggerCode   string
	Description   string
	OperatorID    int64
	OperatorName  string
	OperateTime   time.Time
	UseForRestore bool
	Metadata      map[string]any
}

type CreateGlobalSnapshotCommand struct {
	YearNo       int
	SnapshotType string
	TriggerCode  string
	Description  string
	OperatorID   int64
	OperatorName string
	OperateTime  time.Time
}

type SnapshotPayloadEnvelope struct {
	Version         string                     `json:"version"`
	CapturedAt      time.Time                  `json:"capturedAt"`
	SnapshotScope   string                     `json:"snapshotScope"`
	SnapshotType    string                     `json:"snapshotType"`
	TriggerCode     string                     `json:"triggerCode"`
	TargetGroupID   *int64                     `json:"targetGroupId,omitempty"`
	TargetYearNo    *int                       `json:"targetYearNo,omitempty"`
	TargetStageCode *string                    `json:"targetStageCode,omitempty"`
	GroupState      *GroupSnapshotState        `json:"groupState,omitempty"`
	GlobalState     *GlobalSnapshotState       `json:"globalState,omitempty"`
	StateSummary    map[string]any             `json:"stateSummary"`
	PayloadPreview  map[string]any             `json:"payloadPreview"`
	Metadata        map[string]any             `json:"metadata,omitempty"`
	RawBlocks       map[string]json.RawMessage `json:"rawBlocks,omitempty"`
}

type GroupSnapshotState struct {
	Group                entity.Group                   `json:"group"`
	TargetYearState      entity.GroupYearState          `json:"targetYearState"`
	YearStatesFromTarget []entity.GroupYearState        `json:"yearStatesFromTarget"`
	OperatingDrafts      []entity.GroupOperatingDraft   `json:"operatingDrafts"`
	StageSubmissions     []entity.GroupStageSubmission  `json:"stageSubmissions"`
	Reports              []entity.GroupReport           `json:"reports"`
	ReportSubmissions    []entity.GroupReportSubmission `json:"reportSubmissions"`
	SummarySnapshots     []entity.GroupSummarySnapshot  `json:"summarySnapshots"`
	Adjustments          []entity.GroupAdjustment       `json:"adjustments"`
	MarketBids           []entity.GroupMarketBid        `json:"marketBids"`
	OrderSelections      []entity.GroupOrderSelection   `json:"orderSelections"`
}

type GlobalSnapshotState struct {
	GameConfig             *entity.GameConfig             `json:"gameConfig,omitempty"`
	Groups                 []entity.Group                 `json:"groups"`
	CurrentYearStates      []entity.GroupYearState        `json:"currentYearStates"`
	OrderGenerationConfigs []entity.OrderGenerationConfig `json:"orderGenerationConfigs"`
	OrderMarketConfigs     []entity.OrderMarketConfig     `json:"orderMarketConfigs"`
	OrderGenerationBatches []entity.OrderGenerationBatch  `json:"orderGenerationBatches"`
	OrderPools             []entity.OrderPool             `json:"orderPools"`
	MarketBiddingStates    []entity.MarketBiddingState    `json:"marketBiddingStates"`
	MarketSelectionOrders  []entity.MarketSelectionOrder  `json:"marketSelectionOrders"`
}

type SnapshotService struct {
	db *gorm.DB
}

func NewSnapshotService(db *gorm.DB) *SnapshotService {
	return &SnapshotService{db: db}
}

func (s *SnapshotService) CreateGroupSnapshot(ctx context.Context, cmd CreateGroupSnapshotCommand) (*entity.StateSnapshot, error) {
	return createGroupSnapshot(ctx, s.db, cmd)
}

func (s *SnapshotService) CreateGlobalSnapshot(ctx context.Context, cmd CreateGlobalSnapshotCommand) (*entity.StateSnapshot, error) {
	return createGlobalSnapshot(ctx, s.db, cmd)
}

func createGroupSnapshot(ctx context.Context, db *gorm.DB, cmd CreateGroupSnapshotCommand) (*entity.StateSnapshot, error) {
	cmd.SnapshotType = normalizeSnapshotType(cmd.SnapshotType)
	cmd.TriggerCode = strings.TrimSpace(cmd.TriggerCode)
	cmd.StageCode = normalizeRollbackTargetStageCode(cmd.StageCode)
	cmd.OperatorName = normalizeAdminOperatorName(cmd.OperatorName)
	if cmd.OperateTime.IsZero() {
		cmd.OperateTime = time.Now()
	}
	if cmd.SnapshotType == "" {
		cmd.SnapshotType = enum.SnapshotTypeAuto
	}
	if cmd.TriggerCode == "" {
		cmd.TriggerCode = enum.SnapshotTriggerStageSubmitted
	}
	if cmd.GroupID <= 0 || cmd.YearNo < 0 {
		return nil, ErrRollbackTargetInvalid
	}

	envelope, err := buildGroupSnapshotPayload(ctx, db, cmd)
	if err != nil {
		return nil, err
	}
	payloadBytes, err := marshalSnapshotPayload(envelope)
	if err != nil {
		return nil, err
	}
	payloadHash := hashSnapshotPayload(payloadBytes)

	groupID := cmd.GroupID
	yearNo := cmd.YearNo
	var stageCode *string
	if cmd.StageCode != "" {
		stageCode = &cmd.StageCode
	}
	reportStatus := envelope.GroupState.TargetYearState.ReportStatus
	desc := strings.TrimSpace(cmd.Description)
	var description *string
	if desc != "" {
		description = &desc
	}
	snapshot := &entity.StateSnapshot{
		SnapshotScope:      enum.SnapshotScopeGroup,
		SnapshotType:       cmd.SnapshotType,
		TriggerCode:        cmd.TriggerCode,
		TargetGroupID:      &groupID,
		TargetYearNo:       &yearNo,
		TargetStageCode:    stageCode,
		TargetReportStatus: &reportStatus,
		Description:        description,
		PayloadHash:        payloadHash,
		CreatedByID:        cmd.OperatorID,
		CreatedByName:      cmd.OperatorName,
		CreatedAt:          cmd.OperateTime,
	}
	payloadItem := &entity.StateSnapshotPayload{
		PayloadVersion: rollbackSnapshotPayloadVersion,
		PayloadJSON:    payloadBytes,
		PayloadSize:    len(payloadBytes),
		BaseEntity: entity.BaseEntity{
			Creator:    cmd.OperatorName,
			CreateTime: cmd.OperateTime,
			Updater:    cmd.OperatorName,
			UpdateTime: cmd.OperateTime,
		},
	}
	if err := repository.NewStateSnapshotRepository(db).CreateWithPayload(ctx, snapshot, payloadItem); err != nil {
		return nil, fmt.Errorf("create group snapshot: %w", err)
	}
	return snapshot, nil
}

func createGlobalSnapshot(ctx context.Context, db *gorm.DB, cmd CreateGlobalSnapshotCommand) (*entity.StateSnapshot, error) {
	cmd.SnapshotType = normalizeSnapshotType(cmd.SnapshotType)
	cmd.TriggerCode = strings.TrimSpace(cmd.TriggerCode)
	cmd.OperatorName = normalizeAdminOperatorName(cmd.OperatorName)
	if cmd.OperateTime.IsZero() {
		cmd.OperateTime = time.Now()
	}
	if cmd.SnapshotType == "" {
		cmd.SnapshotType = enum.SnapshotTypeAuto
	}
	if cmd.TriggerCode == "" {
		cmd.TriggerCode = enum.SnapshotTriggerOpenNextYear
	}

	envelope, err := buildGlobalSnapshotPayload(ctx, db, cmd)
	if err != nil {
		return nil, err
	}
	payloadBytes, err := marshalSnapshotPayload(envelope)
	if err != nil {
		return nil, err
	}
	payloadHash := hashSnapshotPayload(payloadBytes)

	var yearNo *int
	if cmd.YearNo >= 0 {
		yearNoValue := cmd.YearNo
		yearNo = &yearNoValue
	}
	desc := strings.TrimSpace(cmd.Description)
	var description *string
	if desc != "" {
		description = &desc
	}
	snapshot := &entity.StateSnapshot{
		SnapshotScope: enum.SnapshotScopeGlobal,
		SnapshotType:  cmd.SnapshotType,
		TriggerCode:   cmd.TriggerCode,
		TargetYearNo:  yearNo,
		Description:   description,
		PayloadHash:   payloadHash,
		CreatedByID:   cmd.OperatorID,
		CreatedByName: cmd.OperatorName,
		CreatedAt:     cmd.OperateTime,
	}
	payloadItem := &entity.StateSnapshotPayload{
		PayloadVersion: rollbackSnapshotPayloadVersion,
		PayloadJSON:    payloadBytes,
		PayloadSize:    len(payloadBytes),
		BaseEntity: entity.BaseEntity{
			Creator:    cmd.OperatorName,
			CreateTime: cmd.OperateTime,
			Updater:    cmd.OperatorName,
			UpdateTime: cmd.OperateTime,
		},
	}
	if err := repository.NewStateSnapshotRepository(db).CreateWithPayload(ctx, snapshot, payloadItem); err != nil {
		return nil, fmt.Errorf("create global snapshot: %w", err)
	}
	return snapshot, nil
}

func buildGroupSnapshotPayload(ctx context.Context, db *gorm.DB, cmd CreateGroupSnapshotCommand) (*SnapshotPayloadEnvelope, error) {
	groupRepo := repository.NewGroupRepository(db)
	yearRepo := repository.NewGroupYearStateRepository(db)
	operatingRepo := repository.NewOperatingRepository(db)
	reportRepo := repository.NewReportRepository(db)
	summaryRepo := repository.NewSummarySnapshotRepository(db)
	adjustmentRepo := repository.NewGroupAdjustmentRepository(db)
	marketBidRepo := repository.NewGroupMarketBidRepository(db)
	orderSelectionRepo := repository.NewGroupOrderSelectionRepository(db)

	group, err := groupRepo.GetByID(ctx, cmd.GroupID)
	if err != nil {
		return nil, fmt.Errorf("load snapshot group: %w", err)
	}
	targetYearState, err := yearRepo.GetByGroupIDAndYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load snapshot target year state: %w", err)
	}
	yearStates, err := yearRepo.ListByGroupIDFromYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load snapshot year states: %w", err)
	}
	operatingDrafts, err := operatingRepo.ListDraftsByGroupFromYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load snapshot operating drafts: %w", err)
	}
	stageSubmissions, err := operatingRepo.ListStageSubmissionsByGroupFromYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load snapshot stage submissions: %w", err)
	}
	reports, err := reportRepo.ListByGroupFromYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load snapshot reports: %w", err)
	}
	reportSubmissions, err := reportRepo.ListSubmissionsByGroupFromYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load snapshot report submissions: %w", err)
	}
	summarySnapshots, err := summaryRepo.ListByGroupFromYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load snapshot summaries: %w", err)
	}
	adjustments, err := adjustmentRepo.ListAllByGroupFromYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load snapshot adjustments: %w", err)
	}
	marketBids, err := marketBidRepo.ListByGroupFromYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load snapshot market bids: %w", err)
	}
	orderSelections, err := orderSelectionRepo.ListByGroupFromYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load snapshot order selections: %w", err)
	}
	orderLinkService := NewOrderOperatingLinkService(
		marketBidRepo,
		repository.NewMarketBiddingStateRepository(db),
		orderSelectionRepo,
	)
	linkedOperatingDrafts, linkedStageSubmissions, err := applyOrderSalesRevenueToSnapshotPayloads(ctx, orderLinkService, cmd.GroupID, operatingDrafts, stageSubmissions)
	if err != nil {
		return nil, err
	}

	targetStageCode := cmd.StageCode
	if targetStageCode == "" {
		targetStageCode = state.CurrentStageCode(targetYearState.StageStatus)
	}
	groupState := &GroupSnapshotState{
		Group:                *group,
		TargetYearState:      *targetYearState,
		YearStatesFromTarget: yearStates,
		OperatingDrafts:      linkedOperatingDrafts,
		StageSubmissions:     linkedStageSubmissions,
		Reports:              reports,
		ReportSubmissions:    reportSubmissions,
		SummarySnapshots:     summarySnapshots,
		Adjustments:          adjustments,
		MarketBids:           marketBids,
		OrderSelections:      orderSelections,
	}
	return &SnapshotPayloadEnvelope{
		Version:         rollbackSnapshotPayloadVersion,
		CapturedAt:      cmd.OperateTime,
		SnapshotScope:   enum.SnapshotScopeGroup,
		SnapshotType:    cmd.SnapshotType,
		TriggerCode:     cmd.TriggerCode,
		TargetGroupID:   &cmd.GroupID,
		TargetYearNo:    &cmd.YearNo,
		TargetStageCode: nullableString(targetStageCode),
		GroupState:      groupState,
		StateSummary:    buildGroupSnapshotStateSummary(*group, *targetYearState, targetStageCode),
		PayloadPreview:  buildGroupSnapshotPayloadPreview(groupState),
		Metadata:        mergeSnapshotMetadata(cmd.Metadata, map[string]any{"useForRestore": cmd.UseForRestore}),
	}, nil
}

func applyOrderSalesRevenueToSnapshotPayloads(
	ctx context.Context,
	orderLinkService *OrderOperatingLinkService,
	groupID int64,
	operatingDrafts []entity.GroupOperatingDraft,
	stageSubmissions []entity.GroupStageSubmission,
) ([]entity.GroupOperatingDraft, []entity.GroupStageSubmission, error) {
	linkedDrafts := make([]entity.GroupOperatingDraft, len(operatingDrafts))
	copy(linkedDrafts, operatingDrafts)
	for i := range linkedDrafts {
		if linkedDrafts[i].YearNo <= 0 || len(linkedDrafts[i].OperatingPayload) == 0 {
			continue
		}
		linkedPayload, err := applyOrderSalesRevenueToSnapshotPayload(ctx, orderLinkService, groupID, linkedDrafts[i].YearNo, linkedDrafts[i].OperatingPayload)
		if err != nil {
			return nil, nil, fmt.Errorf("link snapshot operating draft revenue: %w", err)
		}
		linkedDrafts[i].OperatingPayload = linkedPayload
	}

	linkedSubmissions := make([]entity.GroupStageSubmission, len(stageSubmissions))
	copy(linkedSubmissions, stageSubmissions)
	for i := range linkedSubmissions {
		if linkedSubmissions[i].YearNo <= 0 || len(linkedSubmissions[i].OperatingPayloadSnapshot) == 0 {
			continue
		}
		linkedPayload, err := applyOrderSalesRevenueToSnapshotPayload(ctx, orderLinkService, groupID, linkedSubmissions[i].YearNo, linkedSubmissions[i].OperatingPayloadSnapshot)
		if err != nil {
			return nil, nil, fmt.Errorf("link snapshot stage submission revenue: %w", err)
		}
		linkedSubmissions[i].OperatingPayloadSnapshot = linkedPayload
	}

	return linkedDrafts, linkedSubmissions, nil
}

func applyOrderSalesRevenueToSnapshotPayload(ctx context.Context, orderLinkService *OrderOperatingLinkService, groupID int64, yearNo int, rawPayload []byte) ([]byte, error) {
	var operatingPayload payload.OperatingPayload
	if err := json.Unmarshal(rawPayload, &operatingPayload); err != nil {
		return nil, fmt.Errorf("unmarshal operating payload: %w", err)
	}
	linkedPayload, _, err := orderLinkService.ApplyFormalYearValues(ctx, groupID, yearNo, operatingPayload)
	if err != nil {
		return nil, err
	}
	next, err := json.Marshal(linkedPayload.Normalize())
	if err != nil {
		return nil, fmt.Errorf("marshal linked operating payload: %w", err)
	}
	return next, nil
}

func mergeSnapshotMetadata(source map[string]any, required map[string]any) map[string]any {
	result := make(map[string]any, len(source)+len(required))
	for key, value := range source {
		result[key] = value
	}
	for key, value := range required {
		result[key] = value
	}
	return result
}

func buildGlobalSnapshotPayload(ctx context.Context, db *gorm.DB, cmd CreateGlobalSnapshotCommand) (*SnapshotPayloadEnvelope, error) {
	gameConfigRepo := repository.NewGameConfigRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	yearRepo := repository.NewGroupYearStateRepository(db)
	configRepo := repository.NewOrderGenerationConfigRepository(db)
	marketConfigRepo := repository.NewOrderMarketConfigRepository(db)
	batchRepo := repository.NewOrderGenerationBatchRepository(db)
	poolRepo := repository.NewOrderPoolRepository(db)
	marketStateRepo := repository.NewMarketBiddingStateRepository(db)
	selectionOrderRepo := repository.NewMarketSelectionOrderRepository(db)

	var gameConfig *entity.GameConfig
	config, err := gameConfigRepo.GetCurrent(ctx)
	if err == nil {
		gameConfig = config
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("load snapshot game config: %w", err)
	}
	groups, err := groupRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("load snapshot groups: %w", err)
	}
	currentYear := cmd.YearNo
	if currentYear < 0 && gameConfig != nil {
		currentYear = gameConfig.CurrentOpenYear
	}
	currentStates := []entity.GroupYearState{}
	if currentYear >= 0 {
		currentStates, err = yearRepo.ListByYear(ctx, currentYear)
		if err != nil {
			return nil, fmt.Errorf("load snapshot current year states: %w", err)
		}
	}
	orderGenerationConfigs, err := configRepo.ListAllFromYear(ctx, maxInt(currentYear, 0))
	if err != nil {
		return nil, fmt.Errorf("load snapshot order configs: %w", err)
	}
	orderMarketConfigs, err := marketConfigRepo.ListAllFromYear(ctx, maxInt(currentYear, 0))
	if err != nil {
		return nil, fmt.Errorf("load snapshot order market configs: %w", err)
	}
	orderGenerationBatches, err := batchRepo.ListFromYear(ctx, maxInt(currentYear, 0))
	if err != nil {
		return nil, fmt.Errorf("load snapshot order batches: %w", err)
	}
	orderPools, err := poolRepo.ListFromYear(ctx, maxInt(currentYear, 0))
	if err != nil {
		return nil, fmt.Errorf("load snapshot order pool: %w", err)
	}
	marketBiddingStates, err := marketStateRepo.ListFromYear(ctx, maxInt(currentYear, 0))
	if err != nil {
		return nil, fmt.Errorf("load snapshot market bidding states: %w", err)
	}
	marketSelectionOrders, err := selectionOrderRepo.ListFromYear(ctx, maxInt(currentYear, 0))
	if err != nil {
		return nil, fmt.Errorf("load snapshot market selection orders: %w", err)
	}

	globalState := &GlobalSnapshotState{
		GameConfig:             gameConfig,
		Groups:                 groups,
		CurrentYearStates:      currentStates,
		OrderGenerationConfigs: orderGenerationConfigs,
		OrderMarketConfigs:     orderMarketConfigs,
		OrderGenerationBatches: orderGenerationBatches,
		OrderPools:             orderPools,
		MarketBiddingStates:    marketBiddingStates,
		MarketSelectionOrders:  marketSelectionOrders,
	}
	var targetYearNo *int
	if currentYear >= 0 {
		targetYearNo = &currentYear
	}
	return &SnapshotPayloadEnvelope{
		Version:        rollbackSnapshotPayloadVersion,
		CapturedAt:     cmd.OperateTime,
		SnapshotScope:  enum.SnapshotScopeGlobal,
		SnapshotType:   cmd.SnapshotType,
		TriggerCode:    cmd.TriggerCode,
		TargetYearNo:   targetYearNo,
		GlobalState:    globalState,
		StateSummary:   buildGlobalSnapshotStateSummary(gameConfig, currentYear, groups, currentStates),
		PayloadPreview: buildGlobalSnapshotPayloadPreview(globalState),
	}, nil
}

func buildGroupSnapshotStateSummary(group entity.Group, yearState entity.GroupYearState, targetStageCode string) map[string]any {
	return map[string]any{
		"groupId":                   group.ID,
		"groupNo":                   group.GroupNo,
		"groupName":                 group.GroupName,
		"yearNo":                    yearState.YearNo,
		"yearStatus":                yearState.YearStatus,
		"stageStatus":               yearState.StageStatus,
		"targetStageCode":           targetStageCode,
		"reportStatus":              yearState.ReportStatus,
		"businessStatus":            group.BusinessStatus,
		"summaryEffective":          yearState.SummaryEffective,
		"latestStageSubmitVersion":  yearState.LatestStageSubmitVersion,
		"latestReportSubmitVersion": yearState.LatestReportSubmitVersion,
		"rollbackPending":           yearState.RollbackPending,
	}
}

func buildGroupSnapshotPayloadPreview(groupState *GroupSnapshotState) map[string]any {
	return map[string]any{
		"yearStateCount":        len(groupState.YearStatesFromTarget),
		"operatingDraftCount":   len(groupState.OperatingDrafts),
		"stageSubmissionCount":  len(groupState.StageSubmissions),
		"reportCount":           len(groupState.Reports),
		"reportSubmissionCount": len(groupState.ReportSubmissions),
		"summarySnapshotCount":  len(groupState.SummarySnapshots),
		"adjustmentCount":       len(groupState.Adjustments),
		"marketBidCount":        len(groupState.MarketBids),
		"orderSelectionCount":   len(groupState.OrderSelections),
	}
}

func buildGlobalSnapshotStateSummary(gameConfig *entity.GameConfig, currentYear int, groups []entity.Group, currentStates []entity.GroupYearState) map[string]any {
	result := map[string]any{
		"currentYear": currentYear,
		"groupCount":  len(groups),
		"stateCount":  len(currentStates),
	}
	if gameConfig != nil {
		result["currentOpenYear"] = gameConfig.CurrentOpenYear
		result["finalYear"] = gameConfig.FinalYear
		result["initialBaselineSubmitted"] = gameConfig.InitialBaselineSubmitted
	}
	return result
}

func buildGlobalSnapshotPayloadPreview(globalState *GlobalSnapshotState) map[string]any {
	return map[string]any{
		"groupCount":                 len(globalState.Groups),
		"currentYearStateCount":      len(globalState.CurrentYearStates),
		"orderGenerationConfigCount": len(globalState.OrderGenerationConfigs),
		"orderMarketConfigCount":     len(globalState.OrderMarketConfigs),
		"orderGenerationBatchCount":  len(globalState.OrderGenerationBatches),
		"orderPoolCount":             len(globalState.OrderPools),
		"marketBiddingStateCount":    len(globalState.MarketBiddingStates),
		"marketSelectionOrderCount":  len(globalState.MarketSelectionOrders),
	}
}

func marshalSnapshotPayload(payload *SnapshotPayloadEnvelope) ([]byte, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal snapshot payload: %w", err)
	}
	return raw, nil
}

func hashSnapshotPayload(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func normalizeSnapshotType(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func normalizeRollbackTargetStageCode(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	switch value {
	case enum.StageStatusQ1Open:
		return state.StageCodeQ1
	case enum.StageStatusQ2Open:
		return state.StageCodeQ2
	case enum.StageStatusQ3Open:
		return state.StageCodeQ3
	case enum.StageStatusQ4Open:
		return state.StageCodeQ4
	case enum.StageStatusYearEndOpen:
		return state.StageCodeYearEnd
	default:
		return value
	}
}

func unmarshalSnapshotPayload(raw []byte) (*SnapshotPayloadEnvelope, error) {
	if len(raw) == 0 {
		return nil, ErrRollbackTargetInvalid
	}
	var payload SnapshotPayloadEnvelope
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot payload: %w", err)
	}
	return &payload, nil
}

func maxInt(left int, right int) int {
	if left > right {
		return left
	}
	return right
}
