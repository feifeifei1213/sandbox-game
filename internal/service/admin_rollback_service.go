package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

const (
	defaultRestoreGroupSnapshotReason = "管理员恢复快照"
	rollbackConfirmText               = "确认恢复"
	rollbackStageReport               = "REPORT"
)

type SnapshotSummaryResult struct {
	ID            int64   `json:"id"`
	SnapshotScope string  `json:"snapshotScope"`
	SnapshotType  string  `json:"snapshotType"`
	TriggerCode   string  `json:"triggerCode"`
	GroupID       *int64  `json:"groupId"`
	GroupName     *string `json:"groupName"`
	YearNo        *int    `json:"yearNo"`
	StageCode     *string `json:"stageCode"`
	ReportStatus  *string `json:"reportStatus"`
	Description   string  `json:"description"`
	PayloadHash   string  `json:"payloadHash"`
	CreatedByName string  `json:"createdByName"`
	CreatedAt     string  `json:"createdAt"`
	CanRestore    bool    `json:"canRestore"`
}

type SnapshotListResult struct {
	List     []SnapshotSummaryResult `json:"list"`
	PageNo   int                     `json:"pageNo"`
	PageSize int                     `json:"pageSize"`
	Total    int64                   `json:"total"`
}

type ListSnapshotsCommand struct {
	SnapshotScope string
	SnapshotType  string
	GroupID       *int64
	YearNo        *int
	StageCode     string
	PageNo        int
	PageSize      int
}

type SnapshotDetailResult struct {
	Snapshot       SnapshotSummaryResult `json:"snapshot"`
	StateSummary   map[string]any        `json:"stateSummary"`
	PayloadPreview map[string]any        `json:"payloadPreview"`
	PayloadVersion string                `json:"payloadVersion"`
	PayloadSize    int                   `json:"payloadSize"`
}

type CreateManualSnapshotCommand struct {
	SnapshotScope string
	GroupID       *int64
	YearNo        int
	StageCode     string
	Description   string
	OperatorID    int64
	OperatorName  string
}

type RestoreGroupSnapshotCommand struct {
	SnapshotID   int64
	Reason       string
	ConfirmText  string
	OperatorID   int64
	OperatorName string
}

type RestoreGroupSnapshotResult struct {
	RollbackLogID      int64  `json:"rollbackLogId"`
	SafetySnapshotID   int64  `json:"safetySnapshotId"`
	GroupID            int64  `json:"groupId"`
	TargetYearNo       int    `json:"targetYearNo"`
	TargetStageCode    string `json:"targetStageCode"`
	YearStatus         string `json:"yearStatus"`
	StageStatus        string `json:"stageStatus"`
	ReportStatus       string `json:"reportStatus"`
	BusinessStatus     string `json:"businessStatus"`
	HasRollbackPending bool   `json:"hasRollbackPending"`
}

type AdminRollbackQueryService struct {
	snapshotRepo *repository.StateSnapshotRepository
	groupRepo    *repository.GroupRepository
}

func NewAdminRollbackQueryService(
	snapshotRepo *repository.StateSnapshotRepository,
	groupRepo *repository.GroupRepository,
) *AdminRollbackQueryService {
	return &AdminRollbackQueryService{
		snapshotRepo: snapshotRepo,
		groupRepo:    groupRepo,
	}
}

func (s *AdminRollbackQueryService) ListSnapshots(ctx context.Context, cmd ListSnapshotsCommand) (*SnapshotListResult, error) {
	pageNo := cmd.PageNo
	if pageNo <= 0 {
		pageNo = 1
	}
	pageSize := cmd.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	filter := repository.StateSnapshotListFilter{
		SnapshotScope: strings.ToUpper(strings.TrimSpace(cmd.SnapshotScope)),
		SnapshotType:  strings.ToUpper(strings.TrimSpace(cmd.SnapshotType)),
		GroupID:       cmd.GroupID,
		YearNo:        cmd.YearNo,
		StageCode:     normalizeRollbackTargetStageCode(cmd.StageCode),
		Limit:         pageSize,
		Offset:        (pageNo - 1) * pageSize,
	}
	items, err := s.snapshotRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list snapshots: %w", err)
	}
	total, err := s.snapshotRepo.Count(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("count snapshots: %w", err)
	}
	groupNameMap, err := s.loadSnapshotGroupNames(ctx, items)
	if err != nil {
		return nil, err
	}

	result := make([]SnapshotSummaryResult, 0, len(items))
	for _, item := range items {
		result = append(result, buildSnapshotSummaryResult(item, groupNameMap))
	}
	return &SnapshotListResult{
		List:     result,
		PageNo:   pageNo,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *AdminRollbackQueryService) GetSnapshotDetail(ctx context.Context, snapshotID int64) (*SnapshotDetailResult, error) {
	snapshot, err := s.snapshotRepo.GetByID(ctx, snapshotID)
	if err != nil {
		return nil, fmt.Errorf("load snapshot: %w", err)
	}
	payload, err := s.snapshotRepo.GetPayloadBySnapshotID(ctx, snapshot.ID)
	if err != nil {
		return nil, fmt.Errorf("load snapshot payload: %w", err)
	}
	envelope, err := unmarshalSnapshotPayload(payload.PayloadJSON)
	if err != nil {
		return nil, err
	}
	groupNameMap, err := s.loadSnapshotGroupNames(ctx, []entity.StateSnapshot{*snapshot})
	if err != nil {
		return nil, err
	}
	return &SnapshotDetailResult{
		Snapshot:       buildSnapshotSummaryResult(*snapshot, groupNameMap),
		StateSummary:   envelope.StateSummary,
		PayloadPreview: envelope.PayloadPreview,
		PayloadVersion: payload.PayloadVersion,
		PayloadSize:    payload.PayloadSize,
	}, nil
}

func (s *AdminRollbackQueryService) loadSnapshotGroupNames(ctx context.Context, items []entity.StateSnapshot) (map[int64]string, error) {
	result := make(map[int64]string)
	for _, item := range items {
		if item.TargetGroupID == nil || result[*item.TargetGroupID] != "" {
			continue
		}
		group, err := s.groupRepo.GetByID(ctx, *item.TargetGroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, fmt.Errorf("load snapshot group name: %w", err)
		}
		result[group.ID] = group.GroupName
	}
	return result, nil
}

type AdminRollbackCommandService struct {
	db *gorm.DB
}

func NewAdminRollbackCommandService(db *gorm.DB) *AdminRollbackCommandService {
	return &AdminRollbackCommandService{db: db}
}

func (s *AdminRollbackCommandService) CreateManualSnapshot(ctx context.Context, cmd CreateManualSnapshotCommand) (*SnapshotSummaryResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	description := strings.TrimSpace(cmd.Description)
	if description == "" {
		return nil, ErrRollbackReasonRequired
	}
	scope := strings.ToUpper(strings.TrimSpace(cmd.SnapshotScope))
	now := time.Now()

	var snapshot *entity.StateSnapshot
	var err error
	switch scope {
	case enum.SnapshotScopeGroup:
		if cmd.GroupID == nil || *cmd.GroupID <= 0 || cmd.YearNo < 0 {
			return nil, ErrRollbackTargetInvalid
		}
		snapshot, err = createGroupSnapshot(ctx, s.db, CreateGroupSnapshotCommand{
			GroupID:       *cmd.GroupID,
			YearNo:        cmd.YearNo,
			StageCode:     cmd.StageCode,
			SnapshotType:  enum.SnapshotTypeManual,
			TriggerCode:   enum.SnapshotTypeManual,
			Description:   description,
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
			UseForRestore: true,
		})
	case enum.SnapshotScopeGlobal:
		snapshot, err = createGlobalSnapshot(ctx, s.db, CreateGlobalSnapshotCommand{
			YearNo:       cmd.YearNo,
			SnapshotType: enum.SnapshotTypeManual,
			TriggerCode:  enum.SnapshotTypeManual,
			Description:  description,
			OperatorID:   cmd.OperatorID,
			OperatorName: operatorName,
			OperateTime:  now,
		})
	default:
		return nil, ErrRollbackSnapshotScopeUnsupported
	}
	if err != nil {
		return nil, err
	}
	groupNameMap := map[int64]string{}
	if snapshot.TargetGroupID != nil {
		if group, groupErr := repository.NewGroupRepository(s.db).GetByID(ctx, *snapshot.TargetGroupID); groupErr == nil {
			groupNameMap[group.ID] = group.GroupName
		}
	}
	result := buildSnapshotSummaryResult(*snapshot, groupNameMap)
	return &result, nil
}

func (s *AdminRollbackCommandService) RestoreGroupSnapshot(ctx context.Context, cmd RestoreGroupSnapshotCommand) (*RestoreGroupSnapshotResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	reason := strings.TrimSpace(cmd.Reason)
	if reason == "" {
		reason = defaultRestoreGroupSnapshotReason
	}
	if cmd.SnapshotID <= 0 {
		return nil, ErrRollbackSnapshotNotFound
	}

	var result *RestoreGroupSnapshotResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		snapshotRepo := repository.NewStateSnapshotRepository(tx)
		rollbackRepo := repository.NewRollbackLogRepository(tx)
		gameConfigRepo := repository.NewGameConfigRepository(tx)
		groupRepo := repository.NewGroupRepository(tx)
		yearRepo := repository.NewGroupYearStateRepository(tx)
		reportRepo := repository.NewReportRepository(tx)
		summaryRepo := repository.NewSummarySnapshotRepository(tx)
		adjustmentRepo := repository.NewGroupAdjustmentRepository(tx)
		orderSelectionRepo := repository.NewGroupOrderSelectionRepository(tx)
		deliveryRevisionRepo := repository.NewGroupOrderDeliveryRevisionRepository(tx)

		snapshot, err := snapshotRepo.GetByID(ctx, cmd.SnapshotID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRollbackSnapshotNotFound
			}
			return fmt.Errorf("load target snapshot: %w", err)
		}
		if snapshot.SnapshotScope != enum.SnapshotScopeGroup || snapshot.TargetGroupID == nil || snapshot.TargetYearNo == nil {
			return ErrRollbackSnapshotScopeUnsupported
		}
		payloadItem, err := snapshotRepo.GetPayloadBySnapshotID(ctx, snapshot.ID)
		if err != nil {
			return fmt.Errorf("load target snapshot payload: %w", err)
		}
		envelope, err := unmarshalSnapshotPayload(payloadItem.PayloadJSON)
		if err != nil {
			return err
		}
		if snapshot.TriggerCode == enum.SnapshotTriggerAdjustmentBankruptcy || !snapshotAllowsRestore(envelope) {
			return ErrRollbackSnapshotNotRestorable
		}
		if envelope.GroupState == nil {
			return ErrRollbackTargetInvalid
		}
		gameConfig, err := gameConfigRepo.GetCurrent(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}

		targetGroupID := *snapshot.TargetGroupID
		targetYearNo := *snapshot.TargetYearNo
		targetStageCode := resolveRollbackTargetStageCode(snapshot, envelope)
		if targetStageCode == "" {
			return ErrRollbackTargetInvalid
		}
		currentGroup, err := groupRepo.GetByIDForUpdate(ctx, targetGroupID)
		if err != nil {
			return fmt.Errorf("lock target group: %w", err)
		}
		currentYearState, err := yearRepo.GetByGroupIDAndYearForUpdate(ctx, targetGroupID, targetYearNo)
		if err != nil {
			return fmt.Errorf("lock target year state: %w", err)
		}
		currentRuntime := state.NewRuntimeStateFromEntities(*currentGroup, *currentYearState)
		stateBeforeJSON, err := json.Marshal(state.BuildSnapshot(currentRuntime))
		if err != nil {
			return fmt.Errorf("marshal rollback before state: %w", err)
		}

		now := time.Now()
		safetySnapshot, err := createGroupSnapshot(ctx, tx, CreateGroupSnapshotCommand{
			GroupID:       targetGroupID,
			YearNo:        targetYearNo,
			StageCode:     targetStageCode,
			SnapshotType:  enum.SnapshotTypeSafety,
			TriggerCode:   enum.SnapshotTriggerBeforeRollback,
			Description:   "恢复快照前自动安全快照",
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
			UseForRestore: false,
		})
		if err != nil {
			return fmt.Errorf("create safety snapshot: %w", err)
		}

		nextYearState, nextBusinessStatus, err := buildRestoredRuntimeStateFromSnapshot(*currentGroup, envelope.GroupState.TargetYearState, targetStageCode, targetYearNo)
		if err != nil {
			return err
		}
		if err := groupRepo.RestoreBusinessState(ctx, targetGroupID, nextBusinessStatus, nextYearState.GroupBankruptYearNo, nextYearState.GroupBankruptReason, operatorName, now); err != nil {
			return fmt.Errorf("restore group business state: %w", err)
		}
		if err := yearRepo.RestoreState(ctx, nextYearState.YearState, operatorName, now); err != nil {
			return fmt.Errorf("restore target year state: %w", err)
		}
		if err := yearRepo.ResetAfterRollbackTarget(ctx, targetGroupID, targetYearNo, gameConfig.CurrentOpenYear, operatorName, now); err != nil {
			return fmt.Errorf("reset states after target: %w", err)
		}

		stateAfterJSON, err := json.Marshal(state.BuildSnapshot(state.NewRuntimeStateFromEntities(entity.Group{
			ID:             targetGroupID,
			BusinessStatus: nextBusinessStatus,
		}, nextYearState.YearState)))
		if err != nil {
			return fmt.Errorf("marshal rollback after state: %w", err)
		}

		rollbackLog := &entity.RollbackLog{
			RollbackType:     enum.RollbackTypeGroupSnapshotRestore,
			TargetGroupID:    targetGroupID,
			TargetYearNo:     targetYearNo,
			TargetStageCode:  nullableString(targetStageCode),
			SnapshotID:       &snapshot.ID,
			SafetySnapshotID: safetySnapshot.ID,
			Reason:           reason,
			StateBefore:      stateBeforeJSON,
			StateAfter:       stateAfterJSON,
			OperatorID:       cmd.OperatorID,
			OperatorName:     operatorName,
			OperateTime:      now,
		}
		if err := rollbackRepo.Create(ctx, rollbackLog); err != nil {
			return fmt.Errorf("create rollback log: %w", err)
		}

		if _, err := reportRepo.InvalidateAfterTarget(ctx, targetGroupID, targetYearNo, targetStageCode, operatorName, now); err != nil {
			return fmt.Errorf("invalidate reports after target: %w", err)
		}
		if _, err := summaryRepo.WithdrawByRollbackAfterTarget(ctx, targetGroupID, targetYearNo, rollbackLog.ID, operatorName, now); err != nil {
			return fmt.Errorf("withdraw summaries after target: %w", err)
		}
		if err := restoreAdjustmentStatesFromSnapshot(ctx, adjustmentRepo, repository.NewGroupAdjustmentRevisionRepository(tx), targetGroupID, targetYearNo, envelope.GroupState.Adjustments, rollbackLog.ID, reason, operatorName, now); err != nil {
			return fmt.Errorf("restore adjustments from snapshot: %w", err)
		}
		deliveryDetails, err := orderSelectionRepo.ListEffectiveDeliveryDetailsAfterTarget(ctx, targetGroupID, targetYearNo, targetStageCode)
		if err != nil {
			return fmt.Errorf("load order deliveries before snapshot restore invalidation: %w", err)
		}
		if _, err := orderSelectionRepo.InvalidateDeliveryAfterTarget(ctx, targetGroupID, targetYearNo, targetStageCode, rollbackLog.ID, operatorName, now); err != nil {
			return fmt.Errorf("invalidate order deliveries after target: %w", err)
		}
		if err := createOrderDeliveryInvalidationRevisions(ctx, deliveryRevisionRepo, deliveryDetails, rollbackLog.ID, cmd.OperatorID, operatorName, now, "恢复快照导致原交付失效"); err != nil {
			return fmt.Errorf("create order delivery invalidation revisions: %w", err)
		}
		if err := yearRepo.MarkRollbackPendingRange(ctx, targetGroupID, targetYearNo, gameConfig.CurrentOpenYear, targetYearNo, targetStageCode, rollbackLog.ID, operatorName); err != nil {
			return fmt.Errorf("mark rollback pending: %w", err)
		}

		result = &RestoreGroupSnapshotResult{
			RollbackLogID:      rollbackLog.ID,
			SafetySnapshotID:   safetySnapshot.ID,
			GroupID:            targetGroupID,
			TargetYearNo:       targetYearNo,
			TargetStageCode:    targetStageCode,
			YearStatus:         nextYearState.YearState.YearStatus,
			StageStatus:        nextYearState.YearState.StageStatus,
			ReportStatus:       nextYearState.YearState.ReportStatus,
			BusinessStatus:     nextBusinessStatus,
			HasRollbackPending: true,
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

type restoredGroupYearState struct {
	YearState           entity.GroupYearState
	GroupBankruptYearNo *int
	GroupBankruptReason *string
}

func buildRestoredRuntimeStateFromSnapshot(currentGroup entity.Group, snapshotYearState entity.GroupYearState, targetStageCode string, targetYearNo int) (restoredGroupYearState, string, error) {
	restored := snapshotYearState
	restored.RollbackPending = false
	restored.RollbackTargetYearNo = nil
	restored.RollbackTargetStageCode = nil
	restored.RollbackLogID = nil

	targetStageCode = normalizeRollbackTargetStageCode(targetStageCode)
	switch {
	case targetStageCode == rollbackStageReport:
		restored.YearStatus = enum.YearStatusReportPending
		restored.StageStatus = enum.StageStatusYearEndOpen
		restored.ReportStatus = enum.ReportStatusOpen
		restored.SummaryEffective = false
	case state.IsValidStageCode(targetStageCode):
		stageStatus, _ := state.StageStatusByCode(targetStageCode)
		restored.YearStatus = enum.YearStatusOperating
		restored.StageStatus = stageStatus
		restored.ReportStatus = enum.ReportStatusLocked
		restored.SummaryEffective = false
	default:
		return restoredGroupYearState{}, "", ErrRollbackTargetInvalid
	}

	businessStatus := currentGroup.BusinessStatus
	bankruptYearNo := currentGroup.BankruptYearNo
	bankruptReason := currentGroup.BankruptReason
	if businessStatus == "" {
		businessStatus = enum.BusinessStatusNormal
	}

	runtimeState := state.NewRuntimeStateFromEntities(entity.Group{
		ID:             currentGroup.ID,
		BusinessStatus: businessStatus,
	}, restored)
	if err := state.NewStateMachine().Validate(runtimeState); err != nil {
		return restoredGroupYearState{}, "", fmt.Errorf("validate restored state: %w", err)
	}

	return restoredGroupYearState{
		YearState:           restored,
		GroupBankruptYearNo: bankruptYearNo,
		GroupBankruptReason: bankruptReason,
	}, businessStatus, nil
}

func snapshotAllowsRestore(envelope *SnapshotPayloadEnvelope) bool {
	if envelope == nil || envelope.Metadata == nil {
		return true
	}
	value, exists := envelope.Metadata["useForRestore"]
	if !exists {
		return true
	}
	allowed, ok := value.(bool)
	return ok && allowed
}

func restoreAdjustmentStatesFromSnapshot(
	ctx context.Context,
	adjustmentRepo *repository.GroupAdjustmentRepository,
	revisionRepo *repository.GroupAdjustmentRevisionRepository,
	groupID int64,
	fromYearNo int,
	snapshotItems []entity.GroupAdjustment,
	rollbackID int64,
	reason string,
	operatorName string,
	operateTime time.Time,
) error {
	desired := make(map[int64]bool, len(snapshotItems))
	for _, item := range snapshotItems {
		desired[item.ID] = item.Effective
	}
	currentItems, err := adjustmentRepo.ListAllByGroupFromYear(ctx, groupID, fromYearNo)
	if err != nil {
		return err
	}
	changedYears := map[int]struct{}{}
	for _, item := range currentItems {
		effective, existedAtSnapshot := desired[item.ID]
		if !existedAtSnapshot {
			effective = false
		}
		changed, err := adjustmentRepo.RestoreEffectiveState(ctx, item.ID, effective, rollbackID, reason, operatorName, operateTime)
		if err != nil {
			return err
		}
		if changed {
			changedYears[item.YearNo] = struct{}{}
		}
	}
	for yearNo := range changedYears {
		if _, err := revisionRepo.Increment(ctx, groupID, yearNo, operateTime); err != nil {
			return err
		}
	}
	return nil
}

func resolveRollbackTargetStageCode(snapshot *entity.StateSnapshot, payload *SnapshotPayloadEnvelope) string {
	if snapshot != nil && snapshot.TargetStageCode != nil {
		return normalizeRollbackTargetStageCode(*snapshot.TargetStageCode)
	}
	if payload != nil && payload.TargetStageCode != nil {
		return normalizeRollbackTargetStageCode(*payload.TargetStageCode)
	}
	if snapshot != nil && snapshot.TargetReportStatus != nil && *snapshot.TargetReportStatus == enum.ReportStatusSubmitted {
		return rollbackStageReport
	}
	if payload != nil && payload.GroupState != nil {
		return normalizeRollbackTargetStageCode(state.CurrentStageCode(payload.GroupState.TargetYearState.StageStatus))
	}
	return ""
}

func buildSnapshotSummaryResult(item entity.StateSnapshot, groupNameMap map[int64]string) SnapshotSummaryResult {
	description := ""
	if item.Description != nil {
		description = *item.Description
	}
	var groupName *string
	if item.TargetGroupID != nil {
		if name := groupNameMap[*item.TargetGroupID]; name != "" {
			groupName = &name
		}
	}
	return SnapshotSummaryResult{
		ID:            item.ID,
		SnapshotScope: item.SnapshotScope,
		SnapshotType:  item.SnapshotType,
		TriggerCode:   item.TriggerCode,
		GroupID:       item.TargetGroupID,
		GroupName:     groupName,
		YearNo:        item.TargetYearNo,
		StageCode:     item.TargetStageCode,
		ReportStatus:  item.TargetReportStatus,
		Description:   description,
		PayloadHash:   item.PayloadHash,
		CreatedByName: item.CreatedByName,
		CreatedAt:     item.CreatedAt.Format(time.RFC3339),
		CanRestore: item.SnapshotScope == enum.SnapshotScopeGroup &&
			item.SnapshotType != enum.SnapshotTypeSafety &&
			item.TriggerCode != enum.SnapshotTriggerAdjustmentBankruptcy,
	}
}
