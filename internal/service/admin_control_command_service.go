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
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	"sandbox-game/internal/state"
)

const (
	adminActionCodeUpdateFinalYear       = "UPDATE_FINAL_YEAR"
	adminActionCodeOpenNextYear          = "OPEN_NEXT_YEAR"
	adminActionCodeSubmitInitialBaseline = "SUBMIT_INITIAL_BASELINE"
	adminActionCodeUnlockYear            = "UNLOCK_YEAR"
)

var (
	ErrAdminControlFinalYearTooSmall        = errors.New("admin control final year too small")
	ErrAdminControlTargetYearMismatch       = errors.New("admin control target year mismatch")
	ErrAdminControlFinalYearReached         = errors.New("admin control final year reached")
	ErrAdminControlOpenNextYearBlocked      = errors.New("admin control open next year blocked")
	ErrAdminControlInitialBaselineSubmitted = errors.New("admin control initial baseline already submitted")
	ErrAdminControlInitialBaselineInvalid   = errors.New("admin control initial baseline invalid")
	ErrAdminControlUnlockReasonRequired     = errors.New("admin control unlock reason required")
	ErrAdminControlUnlockNotAllowed         = errors.New("admin control unlock not allowed")
)

const (
	unlockYearBlockedReasonNextYearOpened = "\u4e0b\u4e00\u5e74\u5df2\u5f00\u653e\uff0c\u4e0d\u80fd\u518d\u89e3\u9501\u672c\u5e74"
	unlockYearBlockedReasonAlreadyEditing = "\u5f53\u524d\u5e74\u4efd\u4ecd\u5904\u4e8e\u53ef\u7f16\u8f91\u72b6\u6001\uff0c\u65e0\u9700\u89e3\u9501"
	unlockYearBlockedReasonInvalidState   = "\u5f53\u524d\u5e74\u4efd\u4e0d\u6ee1\u8db3\u5f02\u5e38\u89e3\u9501\u6761\u4ef6"
)

type OpenNextYearBlockedError struct {
	Reason string
}

func (e *OpenNextYearBlockedError) Error() string {
	if e == nil || strings.TrimSpace(e.Reason) == "" {
		return ErrAdminControlOpenNextYearBlocked.Error()
	}
	return e.Reason
}

func (e *OpenNextYearBlockedError) Unwrap() error {
	return ErrAdminControlOpenNextYearBlocked
}

type UnlockNotAllowedError struct {
	Reason string
}

func (e *UnlockNotAllowedError) Error() string {
	if e == nil || strings.TrimSpace(e.Reason) == "" {
		return ErrAdminControlUnlockNotAllowed.Error()
	}
	return e.Reason
}

func (e *UnlockNotAllowedError) Unwrap() error {
	return ErrAdminControlUnlockNotAllowed
}

type UpdateFinalYearCommand struct {
	FinalYear    int
	OperatorID   int64
	OperatorName string
}

type UpdateFinalYearResult struct {
	FinalYear            int    `json:"finalYear"`
	CurrentOpenYear      int    `json:"currentOpenYear"`
	InitializedFromYear  *int   `json:"initializedFromYear"`
	InitializedToYear    *int   `json:"initializedToYear"`
	InitializedYearCount int    `json:"initializedYearCount"`
	UpdatedAt            string `json:"updatedAt"`
	UpdatedBy            string `json:"updatedBy"`
}

type finalYearExpansionPlan struct {
	InitializedFromYear  *int
	InitializedToYear    *int
	InitializedYearCount int
}

type OpenNextYearCommand struct {
	TargetYearNo int
	OperatorID   int64
	OperatorName string
}

type OpenNextYearResult struct {
	PreviousOpenYear          int                 `json:"previousOpenYear"`
	CurrentOpenYear           int                 `json:"currentOpenYear"`
	OpenedYearNo              int                 `json:"openedYearNo"`
	FinalYear                 int                 `json:"finalYear"`
	CanOpenNextYear           bool                `json:"canOpenNextYear"`
	OpenNextYearBlockedReason string              `json:"openNextYearBlockedReason"`
	LatestAdminAction         *AdminActionSummary `json:"latestAdminAction"`
}

type SubmitInitialBaselineCommand struct {
	BaselinePayload *payload.BaselinePayload
	OperatorID      int64
	OperatorName    string
}

type SubmitInitialBaselineResult struct {
	Submitted         bool   `json:"submitted"`
	AppliedGroupCount int    `json:"appliedGroupCount"`
	SubmittedAt       string `json:"submittedAt"`
	SubmitterName     string `json:"submitterName"`
}

type UnlockYearCommand struct {
	GroupID      int64
	YearNo       int
	Reason       string
	OperatorID   int64
	OperatorName string
}

type UnlockYearResult struct {
	GroupID          int64  `json:"groupId"`
	YearNo           int    `json:"yearNo"`
	YearStatus       string `json:"yearStatus"`
	StageStatus      string `json:"stageStatus"`
	ReportStatus     string `json:"reportStatus"`
	SummaryEffective bool   `json:"summaryEffective"`
	BusinessStatus   string `json:"businessStatus"`
	UnlockLogID      int64  `json:"unlockLogId"`
}

type AdminControlCommandService struct {
	db *gorm.DB
}

func NewAdminControlCommandService(db *gorm.DB) *AdminControlCommandService {
	return &AdminControlCommandService{db: db}
}

func (s *AdminControlCommandService) UpdateFinalYear(ctx context.Context, cmd UpdateFinalYearCommand) (*UpdateFinalYearResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)

	var result *UpdateFinalYearResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txGameConfigRepo := repository.NewGameConfigRepository(tx)
		txGroupRepo := repository.NewGroupRepository(tx)
		txGroupYearRepo := repository.NewGroupYearStateRepository(tx)
		txAdminActionLogRepo := repository.NewAdminActionLogRepository(tx)

		gameConfig, err := txGameConfigRepo.GetCurrentForUpdate(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}
		if err := validateFinalYearChange(gameConfig.CurrentOpenYear, cmd.FinalYear); err != nil {
			return err
		}

		expansionPlan := buildFinalYearExpansionPlan(gameConfig.FinalYear, cmd.FinalYear)
		now := time.Now()
		if expansionPlan.InitializedYearCount > 0 {
			groups, err := txGroupRepo.ListAll(ctx)
			if err != nil {
				return fmt.Errorf("load groups: %w", err)
			}
			groupIDs := make([]int64, 0, len(groups))
			for _, group := range groups {
				groupIDs = append(groupIDs, group.ID)
			}
			if _, err := txGroupYearRepo.EnsureFormalYearStates(ctx, repository.EnsureFormalYearStatesCommand{
				GroupIDs:     groupIDs,
				FromYear:     *expansionPlan.InitializedFromYear,
				ToYear:       *expansionPlan.InitializedToYear,
				OperatorName: operatorName,
				OperateTime:  now,
			}); err != nil {
				return fmt.Errorf("initialize future year states: %w", err)
			}
		}

		if err := txGameConfigRepo.UpdateFinalYear(ctx, gameConfig.ID, cmd.FinalYear, operatorName, now); err != nil {
			return fmt.Errorf("update final year: %w", err)
		}

		stateBefore, err := buildFinalYearStateSnapshot(gameConfig.FinalYear, gameConfig.CurrentOpenYear)
		if err != nil {
			return fmt.Errorf("build final year before state: %w", err)
		}
		stateAfter, err := buildFinalYearStateSnapshot(cmd.FinalYear, gameConfig.CurrentOpenYear)
		if err != nil {
			return fmt.Errorf("build final year after state: %w", err)
		}
		actionPayload, err := json.Marshal(map[string]any{
			"finalYear":            cmd.FinalYear,
			"currentOpenYear":      gameConfig.CurrentOpenYear,
			"initializedFromYear":  expansionPlan.InitializedFromYear,
			"initializedToYear":    expansionPlan.InitializedToYear,
			"initializedYearCount": expansionPlan.InitializedYearCount,
		})
		if err != nil {
			return fmt.Errorf("marshal update final year payload: %w", err)
		}

		logItem := &entity.AdminActionLog{
			ActionCode:    adminActionCodeUpdateFinalYear,
			ActionPayload: actionPayload,
			StateBefore:   stateBefore,
			StateAfter:    stateAfter,
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
		}
		if err := txAdminActionLogRepo.Create(ctx, logItem); err != nil {
			return fmt.Errorf("create admin action log: %w", err)
		}

		result = &UpdateFinalYearResult{
			FinalYear:            cmd.FinalYear,
			CurrentOpenYear:      gameConfig.CurrentOpenYear,
			InitializedFromYear:  expansionPlan.InitializedFromYear,
			InitializedToYear:    expansionPlan.InitializedToYear,
			InitializedYearCount: expansionPlan.InitializedYearCount,
			UpdatedAt:            now.Format(time.RFC3339),
			UpdatedBy:            operatorName,
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *AdminControlCommandService) OpenNextYear(ctx context.Context, cmd OpenNextYearCommand) (*OpenNextYearResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	stateMachine := state.NewStateMachine()

	var result *OpenNextYearResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txGameConfigRepo := repository.NewGameConfigRepository(tx)
		txGroupRepo := repository.NewGroupRepository(tx)
		txGroupYearRepo := repository.NewGroupYearStateRepository(tx)
		txAdminActionLogRepo := repository.NewAdminActionLogRepository(tx)

		gameConfig, err := txGameConfigRepo.GetCurrentForUpdate(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}
		if err := validateOpenNextYearTarget(gameConfig.CurrentOpenYear, gameConfig.FinalYear, cmd.TargetYearNo); err != nil {
			return err
		}

		groups, err := txGroupRepo.ListAll(ctx)
		if err != nil {
			return fmt.Errorf("load groups: %w", err)
		}
		currentStates, err := txGroupYearRepo.ListByYear(ctx, gameConfig.CurrentOpenYear)
		if err != nil {
			return fmt.Errorf("load current year states: %w", err)
		}
		currentOpenStatus := evaluateOpenNextYearStatus(gameConfig, groups, currentStates)
		if err := ensureOpenNextYearAllowed(currentOpenStatus); err != nil {
			return err
		}

		targetStates, err := txGroupYearRepo.ListByYear(ctx, cmd.TargetYearNo)
		if err != nil {
			return fmt.Errorf("load target year states: %w", err)
		}
		targetStateMap := make(map[int64]entity.GroupYearState, len(targetStates))
		for _, item := range targetStates {
			targetStateMap[item.GroupID] = item
		}

		updatedTargetStates := make([]entity.GroupYearState, 0, len(targetStates))
		for _, group := range groups {
			targetState, ok := targetStateMap[group.ID]
			if !ok {
				return gorm.ErrRecordNotFound
			}
			if group.BusinessStatus == enum.BusinessStatusBankrupt {
				updatedTargetStates = append(updatedTargetStates, targetState)
				continue
			}

			runtimeState := buildRuntimeState(targetState, group.BusinessStatus)
			nextState, err := stateMachine.OpenYear(runtimeState)
			if err != nil {
				return fmt.Errorf("open target year for group %d: %w", group.ID, err)
			}
			if err := txGroupYearRepo.UpdateRuntimeState(ctx, group.ID, cmd.TargetYearNo, nextState, operatorName); err != nil {
				return fmt.Errorf("update target year state for group %d: %w", group.ID, err)
			}
			targetState.YearStatus = nextState.YearStatus
			targetState.StageStatus = nextState.StageStatus
			targetState.ReportStatus = nextState.ReportStatus
			targetState.SummaryEffective = nextState.SummaryEffective
			targetState.LatestStageSubmitVersion = nextState.LatestStageSubmitVersion
			targetState.LatestReportSubmitVersion = nextState.LatestReportSubmitVersion
			updatedTargetStates = append(updatedTargetStates, targetState)
		}

		now := time.Now()
		if err := txGameConfigRepo.UpdateCurrentOpenYear(ctx, gameConfig.ID, cmd.TargetYearNo, operatorName, now); err != nil {
			return fmt.Errorf("update current open year: %w", err)
		}

		nextConfig := *gameConfig
		nextConfig.CurrentOpenYear = cmd.TargetYearNo
		nextOpenStatus := evaluateOpenNextYearStatus(&nextConfig, groups, updatedTargetStates)
		stateBefore, err := buildOpenNextYearStateSnapshot(gameConfig.CurrentOpenYear, gameConfig.FinalYear, currentOpenStatus)
		if err != nil {
			return fmt.Errorf("build open next year before state: %w", err)
		}
		stateAfter, err := buildOpenNextYearStateSnapshot(cmd.TargetYearNo, gameConfig.FinalYear, nextOpenStatus)
		if err != nil {
			return fmt.Errorf("build open next year after state: %w", err)
		}
		actionPayload, err := json.Marshal(map[string]any{
			"previousOpenYear": gameConfig.CurrentOpenYear,
			"openedYearNo":     cmd.TargetYearNo,
			"finalYear":        gameConfig.FinalYear,
		})
		if err != nil {
			return fmt.Errorf("marshal open next year payload: %w", err)
		}

		targetYearNo := cmd.TargetYearNo
		logItem := &entity.AdminActionLog{
			ActionCode:    adminActionCodeOpenNextYear,
			TargetYearNo:  &targetYearNo,
			ActionPayload: actionPayload,
			StateBefore:   stateBefore,
			StateAfter:    stateAfter,
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
		}
		if err := txAdminActionLogRepo.Create(ctx, logItem); err != nil {
			return fmt.Errorf("create admin action log: %w", err)
		}

		result = &OpenNextYearResult{
			PreviousOpenYear:          gameConfig.CurrentOpenYear,
			CurrentOpenYear:           cmd.TargetYearNo,
			OpenedYearNo:              cmd.TargetYearNo,
			FinalYear:                 gameConfig.FinalYear,
			CanOpenNextYear:           nextOpenStatus.CanOpenNextYear,
			OpenNextYearBlockedReason: nextOpenStatus.BlockedReason,
			LatestAdminAction:         buildAdminActionSummary(logItem),
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *AdminControlCommandService) SubmitInitialBaseline(ctx context.Context, cmd SubmitInitialBaselineCommand) (*SubmitInitialBaselineResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	if err := validateInitialBaselineSubmission(false, cmd.BaselinePayload); err != nil && !errors.Is(err, ErrAdminControlInitialBaselineSubmitted) {
		return nil, err
	}

	var result *SubmitInitialBaselineResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txGameConfigRepo := repository.NewGameConfigRepository(tx)
		txGroupRepo := repository.NewGroupRepository(tx)
		txInitialBaselineRepo := repository.NewInitialBaselineRepository(tx)
		txAdminActionLogRepo := repository.NewAdminActionLogRepository(tx)

		gameConfig, err := txGameConfigRepo.GetCurrentForUpdate(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}
		if err := validateInitialBaselineSubmission(gameConfig.InitialBaselineSubmitted, cmd.BaselinePayload); err != nil {
			return err
		}

		groups, err := txGroupRepo.ListAll(ctx)
		if err != nil {
			return fmt.Errorf("load groups: %w", err)
		}
		groupIDs := make([]int64, 0, len(groups))
		for _, group := range groups {
			groupIDs = append(groupIDs, group.ID)
		}

		beforeCount, err := txInitialBaselineRepo.CountSubmitted(ctx)
		if err != nil {
			return fmt.Errorf("count submitted initial baseline before update: %w", err)
		}

		now := time.Now()
		submitterID := cmd.OperatorID
		if err := txInitialBaselineRepo.UpsertSharedTemplate(ctx, repository.UpsertSharedInitialBaselineCommand{
			GroupIDs:        groupIDs,
			BaselinePayload: *cmd.BaselinePayload,
			Submitted:       true,
			SubmitterID:     &submitterID,
			SubmittedAt:     &now,
			OperatorName:    operatorName,
			OperateTime:     now,
		}); err != nil {
			return fmt.Errorf("upsert shared initial baseline: %w", err)
		}

		if err := txGameConfigRepo.UpdateInitialBaselineSubmitted(ctx, gameConfig.ID, true, operatorName, now); err != nil {
			return fmt.Errorf("update initial baseline submitted flag: %w", err)
		}

		stateBefore, err := buildInitialBaselineStateSnapshot(gameConfig.InitialBaselineSubmitted, int(beforeCount))
		if err != nil {
			return fmt.Errorf("build initial baseline before state: %w", err)
		}
		stateAfter, err := buildInitialBaselineStateSnapshot(true, len(groupIDs))
		if err != nil {
			return fmt.Errorf("build initial baseline after state: %w", err)
		}
		actionPayload, err := json.Marshal(map[string]any{
			"baselinePayload":   cmd.BaselinePayload,
			"appliedGroupCount": len(groupIDs),
		})
		if err != nil {
			return fmt.Errorf("marshal submit initial baseline payload: %w", err)
		}

		logItem := &entity.AdminActionLog{
			ActionCode:    adminActionCodeSubmitInitialBaseline,
			ActionPayload: actionPayload,
			StateBefore:   stateBefore,
			StateAfter:    stateAfter,
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
		}
		if err := txAdminActionLogRepo.Create(ctx, logItem); err != nil {
			return fmt.Errorf("create admin action log: %w", err)
		}

		result = &SubmitInitialBaselineResult{
			Submitted:         true,
			AppliedGroupCount: len(groupIDs),
			SubmittedAt:       now.Format(time.RFC3339),
			SubmitterName:     operatorName,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *AdminControlCommandService) UnlockYear(ctx context.Context, cmd UnlockYearCommand) (*UnlockYearResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	reason := strings.TrimSpace(cmd.Reason)

	var result *UnlockYearResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txGameConfigRepo := repository.NewGameConfigRepository(tx)
		txGroupRepo := repository.NewGroupRepository(tx)
		txGroupYearRepo := repository.NewGroupYearStateRepository(tx)
		txReportRepo := repository.NewReportRepository(tx)
		txSummaryRepo := repository.NewSummarySnapshotRepository(tx)
		txAdminUnlockLogRepo := repository.NewAdminUnlockLogRepository(tx)
		txAdminActionLogRepo := repository.NewAdminActionLogRepository(tx)

		gameConfig, err := txGameConfigRepo.GetCurrent(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}
		group, err := txGroupRepo.GetByID(ctx, cmd.GroupID)
		if err != nil {
			return fmt.Errorf("load group: %w", err)
		}
		yearState, err := txGroupYearRepo.GetByGroupIDAndYear(ctx, cmd.GroupID, cmd.YearNo)
		if err != nil {
			return fmt.Errorf("load group year state: %w", err)
		}

		currentState := state.NewRuntimeStateFromEntities(*group, *yearState)
		nextState, err := ensureUnlockYearAllowed(currentState, reason, gameConfig.CurrentOpenYear > cmd.YearNo)
		if err != nil {
			return err
		}

		businessRecovered := false
		if shouldRecoverFromBankrupt(group, cmd.YearNo) {
			recovered, recoverErr := txGroupRepo.RecoverFromBankrupt(ctx, cmd.GroupID, cmd.YearNo, operatorName)
			if recoverErr != nil {
				return fmt.Errorf("recover group from bankrupt: %w", recoverErr)
			}
			if recovered {
				businessRecovered = true
				nextState.BusinessStatus = enum.BusinessStatusNormal
			}
		}

		stateBeforeJSON, err := json.Marshal(state.BuildSnapshot(currentState))
		if err != nil {
			return fmt.Errorf("marshal unlock state before: %w", err)
		}
		stateAfterJSON, err := json.Marshal(state.BuildSnapshot(nextState))
		if err != nil {
			return fmt.Errorf("marshal unlock state after: %w", err)
		}

		now := time.Now()
		if err := txGroupYearRepo.UpdateRuntimeState(ctx, cmd.GroupID, cmd.YearNo, nextState, operatorName); err != nil {
			return fmt.Errorf("update group year state: %w", err)
		}

		reportInvalidated, err := txReportRepo.InvalidateSubmission(ctx, cmd.GroupID, cmd.YearNo, operatorName, now)
		if err != nil {
			return fmt.Errorf("invalidate report submission: %w", err)
		}
		summaryWithdrawn, err := txSummaryRepo.WithdrawEffective(ctx, cmd.GroupID, cmd.YearNo, operatorName, now)
		if err != nil {
			return fmt.Errorf("withdraw summary snapshot: %w", err)
		}

		unlockLogItem := &entity.AdminUnlockLog{
			GroupID:      cmd.GroupID,
			YearNo:       cmd.YearNo,
			Reason:       reason,
			StateBefore:  stateBeforeJSON,
			StateAfter:   stateAfterJSON,
			OperatorID:   cmd.OperatorID,
			OperatorName: operatorName,
			OperateTime:  now,
		}
		if err := txAdminUnlockLogRepo.Create(ctx, unlockLogItem); err != nil {
			return fmt.Errorf("create admin unlock log: %w", err)
		}

		targetGroupID := cmd.GroupID
		targetYearNo := cmd.YearNo
		actionPayload, err := json.Marshal(map[string]any{
			"groupId":           cmd.GroupID,
			"yearNo":            cmd.YearNo,
			"reason":            reason,
			"reportInvalidated": reportInvalidated,
			"summaryWithdrawn":  summaryWithdrawn,
			"businessRecovered": businessRecovered,
		})
		if err != nil {
			return fmt.Errorf("marshal unlock year payload: %w", err)
		}

		actionLogItem := &entity.AdminActionLog{
			ActionCode:    adminActionCodeUnlockYear,
			TargetGroupID: &targetGroupID,
			TargetYearNo:  &targetYearNo,
			ActionPayload: actionPayload,
			StateBefore:   stateBeforeJSON,
			StateAfter:    stateAfterJSON,
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
		}
		if err := txAdminActionLogRepo.Create(ctx, actionLogItem); err != nil {
			return fmt.Errorf("create admin action log: %w", err)
		}

		result = &UnlockYearResult{
			GroupID:          cmd.GroupID,
			YearNo:           cmd.YearNo,
			YearStatus:       nextState.YearStatus,
			StageStatus:      nextState.StageStatus,
			ReportStatus:     nextState.ReportStatus,
			SummaryEffective: nextState.SummaryEffective,
			BusinessStatus:   nextState.BusinessStatus,
			UnlockLogID:      unlockLogItem.ID,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func buildFinalYearExpansionPlan(previousFinalYear int, nextFinalYear int) finalYearExpansionPlan {
	if nextFinalYear <= previousFinalYear {
		return finalYearExpansionPlan{}
	}

	fromYear := previousFinalYear + 1
	toYear := nextFinalYear
	return finalYearExpansionPlan{
		InitializedFromYear:  intPointer(fromYear),
		InitializedToYear:    intPointer(toYear),
		InitializedYearCount: toYear - fromYear + 1,
	}
}

func intPointer(value int) *int {
	return &value
}

func validateFinalYearChange(currentOpenYear int, finalYear int) error {
	if finalYear < currentOpenYear {
		return ErrAdminControlFinalYearTooSmall
	}
	return nil
}

func validateOpenNextYearTarget(currentOpenYear int, finalYear int, targetYearNo int) error {
	if currentOpenYear >= finalYear {
		return ErrAdminControlFinalYearReached
	}
	if targetYearNo != currentOpenYear+1 {
		return ErrAdminControlTargetYearMismatch
	}
	if targetYearNo > finalYear {
		return ErrAdminControlFinalYearReached
	}
	return nil
}

func ensureOpenNextYearAllowed(status openNextYearStatus) error {
	if status.CanOpenNextYear {
		return nil
	}
	if status.BlockedReason == openNextYearBlockedReasonFinalYearReached {
		return ErrAdminControlFinalYearReached
	}
	return &OpenNextYearBlockedError{Reason: status.BlockedReason}
}

func validateInitialBaselineSubmission(alreadySubmitted bool, baselinePayload *payload.BaselinePayload) error {
	if baselinePayload == nil {
		return ErrAdminControlInitialBaselineInvalid
	}
	if alreadySubmitted {
		return ErrAdminControlInitialBaselineSubmitted
	}
	return nil
}

func ensureUnlockYearAllowed(current state.RuntimeState, reason string, nextYearAlreadyOpened bool) (state.RuntimeState, error) {
	if strings.TrimSpace(reason) == "" {
		return current, ErrAdminControlUnlockReasonRequired
	}
	if nextYearAlreadyOpened {
		return current, &UnlockNotAllowedError{Reason: unlockYearBlockedReasonNextYearOpened}
	}
	if isUnlockAlreadyEditable(current) {
		return current, &UnlockNotAllowedError{Reason: unlockYearBlockedReasonAlreadyEditing}
	}

	nextState, err := state.NewStateMachine().UnlockYear(current, false)
	if err != nil {
		if errors.Is(err, state.ErrYearCannotUnlock) {
			return current, &UnlockNotAllowedError{Reason: unlockYearBlockedReasonInvalidState}
		}
		return current, err
	}
	return nextState, nil
}

func buildFinalYearStateSnapshot(finalYear int, currentOpenYear int) ([]byte, error) {
	return json.Marshal(map[string]any{
		"finalYear":       finalYear,
		"currentOpenYear": currentOpenYear,
	})
}

func buildOpenNextYearStateSnapshot(currentOpenYear int, finalYear int, status openNextYearStatus) ([]byte, error) {
	return json.Marshal(map[string]any{
		"currentOpenYear":           currentOpenYear,
		"finalYear":                 finalYear,
		"canOpenNextYear":           status.CanOpenNextYear,
		"nextOpenableYear":          status.NextOpenableYear,
		"openNextYearBlockedReason": status.BlockedReason,
	})
}

func buildInitialBaselineStateSnapshot(submitted bool, appliedGroupCount int) ([]byte, error) {
	return json.Marshal(map[string]any{
		"initialBaselineSubmitted": submitted,
		"appliedGroupCount":        appliedGroupCount,
	})
}

func normalizeAdminOperatorName(operatorName string) string {
	operatorName = strings.TrimSpace(operatorName)
	if operatorName == "" {
		return "admin"
	}
	return operatorName
}

func isUnlockAlreadyEditable(current state.RuntimeState) bool {
	guard := state.NewTransitionGuard()
	return guard.BuildOperatingPermission(current).CanEdit || guard.BuildReportPermission(current).CanEdit
}

func shouldRecoverFromBankrupt(group *entity.Group, yearNo int) bool {
	if group == nil || group.BusinessStatus != enum.BusinessStatusBankrupt || group.BankruptYearNo == nil {
		return false
	}
	return *group.BankruptYearNo == yearNo
}

func buildRuntimeState(item entity.GroupYearState, businessStatus string) state.RuntimeState {
	return state.RuntimeState{
		YearNo:                    item.YearNo,
		YearType:                  item.YearType,
		YearStatus:                item.YearStatus,
		StageStatus:               item.StageStatus,
		ReportStatus:              item.ReportStatus,
		BusinessStatus:            businessStatus,
		SummaryEffective:          item.SummaryEffective,
		LatestStageSubmitVersion:  item.LatestStageSubmitVersion,
		LatestReportSubmitVersion: item.LatestReportSubmitVersion,
	}
}
