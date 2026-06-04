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
	adminActionCodeInitializeGame        = "INITIALIZE_GAME"
	adminActionCodeUpdateFinalYear       = "UPDATE_FINAL_YEAR"
	adminActionCodeOpenNextYear          = "OPEN_NEXT_YEAR"
	adminActionCodeSubmitInitialBaseline = "SUBMIT_INITIAL_BASELINE"
	adminActionCodeUnlockYear            = "UNLOCK_YEAR"
	unlockTargetTypeOperating            = "OPERATING"
	unlockTargetTypeReport               = "REPORT"
	defaultUnlockRetryReason             = "管理员退回重提"
	initializeGameMaxGroupCount          = 10
	initializeGameDefaultPassword        = "123456"
)

var (
	ErrAdminControlAlreadyInitialized       = errors.New("admin control already initialized")
	ErrAdminControlInitializeInvalid        = errors.New("admin control initialize invalid")
	ErrAdminControlFinalYearTooSmall        = errors.New("admin control final year too small")
	ErrAdminControlTargetYearMismatch       = errors.New("admin control target year mismatch")
	ErrAdminControlFinalYearReached         = errors.New("admin control final year reached")
	ErrAdminControlOpenNextYearBlocked      = errors.New("admin control open next year blocked")
	ErrAdminControlInitialBaselineSubmitted = errors.New("admin control initial baseline already submitted")
	ErrAdminControlInitialBaselineInvalid   = errors.New("admin control initial baseline invalid")
	ErrAdminControlEditionRequired          = errors.New("admin control edition required")
	ErrAdminControlEditionInvalid           = errors.New("admin control edition invalid")
	ErrAdminControlUnlockReasonRequired     = errors.New("admin control unlock reason required")
	ErrAdminControlUnlockTargetTypeRequired = errors.New("admin control unlock target type required")
	ErrAdminControlUnlockTargetTypeInvalid  = errors.New("admin control unlock target type invalid")
	ErrAdminControlUnlockStageRequired      = errors.New("admin control unlock stage required")
	ErrAdminControlUnlockStageInvalid       = errors.New("admin control unlock stage invalid")
	ErrAdminControlUnlockNotAllowed         = errors.New("admin control unlock not allowed")
)

const (
	unlockYearBlockedReasonNextYearOpened     = "下一年已开放，不能再解锁本年"
	unlockYearBlockedReasonOperatingEditable  = "经营页当前无需解锁"
	unlockYearBlockedReasonReportEditable     = "财报页当前无需解锁"
	unlockYearBlockedReasonTargetNotSubmitted = "目标尚未正式提交，不能解锁"
	unlockYearBlockedReasonInvalidState       = "当前目标不满足异常解锁条件"
)

type InitializeInvalidError struct {
	Reason string
}

func (e *InitializeInvalidError) Error() string {
	if e == nil || strings.TrimSpace(e.Reason) == "" {
		return ErrAdminControlInitializeInvalid.Error()
	}
	return e.Reason
}

func (e *InitializeInvalidError) Unwrap() error {
	return ErrAdminControlInitializeInvalid
}

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

type InitializeGameCommand struct {
	GroupCount         int
	EditionCode        string
	DictionarySchemeID *int64
	DictionaryItems    []DictionaryItemInput
	OperatorID         int64
	OperatorName       string
}

type InitializeGameResult struct {
	Initialized           bool   `json:"initialized"`
	GroupCount            int    `json:"groupCount"`
	EditionCode           string `json:"editionCode"`
	EditionName           string `json:"editionName"`
	RuleVersion           string `json:"ruleVersion"`
	TemplateVersion       string `json:"templateVersion"`
	DictionaryRevision    int    `json:"dictionaryRevision"`
	CreatedGroupCount     int    `json:"createdGroupCount"`
	CreatedAccountCount   int    `json:"createdAccountCount"`
	CreatedYearStateCount int    `json:"createdYearStateCount"`
	InitializedAt         string `json:"initializedAt"`
	InitializedBy         string `json:"initializedBy"`
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
	GroupID          int64
	YearNo           int
	UnlockTargetType string
	TargetStageCode  string
	Reason           string
	OperatorID       int64
	OperatorName     string
}

type UnlockYearResult struct {
	GroupID           int64   `json:"groupId"`
	YearNo            int     `json:"yearNo"`
	UnlockTargetType  string  `json:"unlockTargetType"`
	TargetStageCode   *string `json:"targetStageCode"`
	EditableStageCode *string `json:"editableStageCode"`
	YearStatus        string  `json:"yearStatus"`
	StageStatus       string  `json:"stageStatus"`
	ReportStatus      string  `json:"reportStatus"`
	SummaryEffective  bool    `json:"summaryEffective"`
	BusinessStatus    string  `json:"businessStatus"`
	UnlockLogID       int64   `json:"unlockLogId"`
}

type AdminControlCommandService struct {
	db *gorm.DB
}

func NewAdminControlCommandService(db *gorm.DB) *AdminControlCommandService {
	return &AdminControlCommandService{db: db}
}

func (s *AdminControlCommandService) InitializeGame(ctx context.Context, cmd InitializeGameCommand) (*InitializeGameResult, error) {
	operatorName := normalizeAdminOperatorName(cmd.OperatorName)
	if err := validateInitializeGameInput(cmd.GroupCount); err != nil {
		return nil, err
	}
	edition, err := validateInitializeGameEdition(cmd.EditionCode)
	if err != nil {
		return nil, err
	}

	var result *InitializeGameResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txGameConfigRepo := repository.NewGameConfigRepository(tx)
		txGroupRepo := repository.NewGroupRepository(tx)
		txAccountRepo := repository.NewAccountRepository(tx)
		txGroupYearRepo := repository.NewGroupYearStateRepository(tx)
		txAdminActionLogRepo := repository.NewAdminActionLogRepository(tx)
		txDictionaryRepo := repository.NewDictionaryRepository(tx)

		gameConfig, err := txGameConfigRepo.GetCurrentForUpdate(ctx)
		if err != nil {
			return fmt.Errorf("load game config: %w", err)
		}

		existingGroupCount, err := txGroupRepo.CountAll(ctx)
		if err != nil {
			return fmt.Errorf("count groups: %w", err)
		}
		if existingGroupCount > 0 {
			return ErrAdminControlAlreadyInitialized
		}

		existingGroupAccountCount, err := txAccountRepo.CountGroupAccounts(ctx)
		if err != nil {
			return fmt.Errorf("count group accounts: %w", err)
		}
		existingYearStateCount, err := txGroupYearRepo.CountAll(ctx)
		if err != nil {
			return fmt.Errorf("count group year states: %w", err)
		}
		if err := ensureInitializeGameEnvironmentClean(existingGroupAccountCount, existingYearStateCount); err != nil {
			return err
		}

		now := time.Now()
		groups := buildInitializeGameGroups(cmd.GroupCount, operatorName, now)
		if err := txGroupRepo.CreateBatch(ctx, groups); err != nil {
			return fmt.Errorf("create groups: %w", err)
		}

		accounts := buildInitializeGameAccounts(groups, operatorName, now)
		if err := txAccountRepo.CreateBatch(ctx, accounts); err != nil {
			return fmt.Errorf("create group accounts: %w", err)
		}

		yearStates := buildInitializeGameYearStates(groups, gameConfig.FinalYear, operatorName, now)
		if err := txGroupYearRepo.CreateBatch(ctx, yearStates); err != nil {
			return fmt.Errorf("create group year states: %w", err)
		}

		dictionaryItems, dictionarySchemeID, err := resolveInitializeDictionaryItems(ctx, txDictionaryRepo, edition.EditionCode, cmd.DictionarySchemeID, cmd.DictionaryItems)
		if err != nil {
			return err
		}
		dictionaryRevision := 1
		currentDictionaryItems := currentDictionaryEntitiesFromResults(edition.EditionCode, dictionaryItems, operatorName, now)
		if err := txDictionaryRepo.ReplaceCurrentItems(ctx, currentDictionaryItems); err != nil {
			return fmt.Errorf("initialize current dictionary: %w", err)
		}

		if err := txGameConfigRepo.PrepareForInitialization(ctx, gameConfig.ID, repository.PrepareGameConfigInitializationCommand{
			EditionCode:              edition.EditionCode,
			EditionName:              edition.EditionName,
			RuleVersion:              edition.RuleVersion,
			TemplateVersion:          edition.TemplateVersion,
			OperatingTemplateVersion: edition.OperatingTemplateVersion,
			ReportTemplateVersion:    edition.ReportTemplateVersion,
			OrderTemplateVersion:     edition.OrderTemplateVersion,
			ProcessRuleVersion:       edition.ProcessRuleVersion,
			DictionaryRevision:       dictionaryRevision,
			OperatorName:             operatorName,
			OperateTime:              now,
		}); err != nil {
			return fmt.Errorf("prepare game config for initialization: %w", err)
		}

		stateBefore, err := buildInitializeGameStateSnapshot(int(existingGroupCount), gameConfig.CurrentOpenYear, gameConfig.FinalYear, gameConfig.InitialBaselineSubmitted, normalizeGameConfigEdition(*gameConfig))
		if err != nil {
			return fmt.Errorf("build initialize game before state: %w", err)
		}
		stateAfter, err := buildInitializeGameStateSnapshot(len(groups), 0, gameConfig.FinalYear, false, edition)
		if err != nil {
			return fmt.Errorf("build initialize game after state: %w", err)
		}
		actionPayload, err := json.Marshal(map[string]any{
			"groupCount":            len(groups),
			"editionCode":           edition.EditionCode,
			"editionName":           edition.EditionName,
			"ruleVersion":           edition.RuleVersion,
			"templateVersion":       edition.TemplateVersion,
			"dictionaryRevision":    dictionaryRevision,
			"dictionarySchemeId":    dictionarySchemeID,
			"createdGroupCount":     len(groups),
			"createdAccountCount":   len(accounts),
			"createdYearStateCount": len(yearStates),
		})
		if err != nil {
			return fmt.Errorf("marshal initialize game payload: %w", err)
		}

		logItem := &entity.AdminActionLog{
			ActionCode:    adminActionCodeInitializeGame,
			ActionPayload: actionPayload,
			StateBefore:   stateBefore,
			StateAfter:    stateAfter,
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
		}
		if err := txAdminActionLogRepo.Create(ctx, logItem); err != nil {
			return fmt.Errorf("create initialize game action log: %w", err)
		}
		if err := createDictionaryChangeLog(ctx, txDictionaryRepo, edition.EditionCode, dictionaryChangeTypeInitialize, dictionarySchemeID, "初始化比赛时生成当前比赛字典快照", nil, dictionaryItems, dictionaryRevision, cmd.OperatorID, operatorName, now); err != nil {
			return err
		}

		result = &InitializeGameResult{
			Initialized:           true,
			GroupCount:            len(groups),
			EditionCode:           edition.EditionCode,
			EditionName:           edition.EditionName,
			RuleVersion:           edition.RuleVersion,
			TemplateVersion:       edition.TemplateVersion,
			DictionaryRevision:    dictionaryRevision,
			CreatedGroupCount:     len(groups),
			CreatedAccountCount:   len(accounts),
			CreatedYearStateCount: len(yearStates),
			InitializedAt:         now.Format(time.RFC3339),
			InitializedBy:         operatorName,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
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
		if _, err := createGlobalSnapshot(ctx, tx, CreateGlobalSnapshotCommand{
			YearNo:       cmd.TargetYearNo,
			SnapshotType: enum.SnapshotTypeAuto,
			TriggerCode:  enum.SnapshotTriggerOpenNextYear,
			Description:  "开放下一年后自动快照",
			OperatorID:   cmd.OperatorID,
			OperatorName: operatorName,
			OperateTime:  now,
		}); err != nil {
			return err
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
	if reason == "" {
		reason = defaultUnlockRetryReason
	}
	targetType := normalizeUnlockTargetType(cmd.UnlockTargetType)
	targetStageCode := normalizeUnlockTargetStageCode(cmd.TargetStageCode)

	var result *UnlockYearResult
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txGameConfigRepo := repository.NewGameConfigRepository(tx)
		txGroupRepo := repository.NewGroupRepository(tx)
		txGroupYearRepo := repository.NewGroupYearStateRepository(tx)
		txReportRepo := repository.NewReportRepository(tx)
		txSummaryRepo := repository.NewSummarySnapshotRepository(tx)
		txAdminUnlockLogRepo := repository.NewAdminUnlockLogRepository(tx)
		txAdminActionLogRepo := repository.NewAdminActionLogRepository(tx)
		txRollbackLogRepo := repository.NewRollbackLogRepository(tx)
		txAdjustmentRepo := repository.NewGroupAdjustmentRepository(tx)
		txOrderSelectionRepo := repository.NewGroupOrderSelectionRepository(tx)

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
		nextState, err := ensureUnlockYearAllowed(currentState, reason, targetType, targetStageCode, gameConfig.CurrentOpenYear > cmd.YearNo)
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
		rollbackTargetStageCode := targetStageCode
		if targetType == unlockTargetTypeReport {
			rollbackTargetStageCode = rollbackStageReport
		}
		safetySnapshot, err := createGroupSnapshot(ctx, tx, CreateGroupSnapshotCommand{
			GroupID:       cmd.GroupID,
			YearNo:        cmd.YearNo,
			StageCode:     rollbackTargetStageCode,
			SnapshotType:  enum.SnapshotTypeSafety,
			TriggerCode:   enum.SnapshotTriggerBeforeRollback,
			Description:   "退回重提前自动安全快照",
			OperatorID:    cmd.OperatorID,
			OperatorName:  operatorName,
			OperateTime:   now,
			UseForRestore: false,
		})
		if err != nil {
			return fmt.Errorf("create unlock safety snapshot: %w", err)
		}
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

		rollbackLogItem := &entity.RollbackLog{
			RollbackType:     enum.RollbackTypeUnlockRetry,
			TargetGroupID:    cmd.GroupID,
			TargetYearNo:     cmd.YearNo,
			TargetStageCode:  nullableString(rollbackTargetStageCode),
			SafetySnapshotID: safetySnapshot.ID,
			Reason:           reason,
			StateBefore:      stateBeforeJSON,
			StateAfter:       stateAfterJSON,
			OperatorID:       cmd.OperatorID,
			OperatorName:     operatorName,
			OperateTime:      now,
		}
		if err := txRollbackLogRepo.Create(ctx, rollbackLogItem); err != nil {
			return fmt.Errorf("create unlock rollback log: %w", err)
		}
		if _, err := txAdjustmentRepo.MarkInvalidAfterTarget(ctx, cmd.GroupID, cmd.YearNo, rollbackTargetStageCode, rollbackLogItem.ID, reason, operatorName, now); err != nil {
			return fmt.Errorf("invalidate adjustments after unlock: %w", err)
		}
		if _, err := txOrderSelectionRepo.InvalidateDeliveryAfterTarget(ctx, cmd.GroupID, cmd.YearNo, rollbackTargetStageCode, rollbackLogItem.ID, operatorName, now); err != nil {
			return fmt.Errorf("invalidate order delivery after unlock: %w", err)
		}
		if err := txGroupYearRepo.MarkRollbackPending(ctx, cmd.GroupID, cmd.YearNo, cmd.YearNo, rollbackTargetStageCode, rollbackLogItem.ID, operatorName); err != nil {
			return fmt.Errorf("mark unlock rollback pending: %w", err)
		}

		unlockLogItem := &entity.AdminUnlockLog{
			GroupID:          cmd.GroupID,
			YearNo:           cmd.YearNo,
			Reason:           reason,
			UnlockTargetType: targetType,
			TargetStageCode:  nullableString(targetStageCode),
			SafetySnapshotID: &safetySnapshot.ID,
			StateBefore:      stateBeforeJSON,
			StateAfter:       stateAfterJSON,
			OperatorID:       cmd.OperatorID,
			OperatorName:     operatorName,
			OperateTime:      now,
		}
		if err := txAdminUnlockLogRepo.Create(ctx, unlockLogItem); err != nil {
			return fmt.Errorf("create admin unlock log: %w", err)
		}

		targetGroupID := cmd.GroupID
		targetYearNo := cmd.YearNo
		actionPayload, err := json.Marshal(map[string]any{
			"groupId":           cmd.GroupID,
			"yearNo":            cmd.YearNo,
			"unlockTargetType":  targetType,
			"targetStageCode":   nullableString(targetStageCode),
			"reason":            reason,
			"reportInvalidated": reportInvalidated,
			"summaryWithdrawn":  summaryWithdrawn,
			"businessRecovered": businessRecovered,
			"rollbackLogId":     rollbackLogItem.ID,
			"safetySnapshotId":  safetySnapshot.ID,
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
			GroupID:           cmd.GroupID,
			YearNo:            cmd.YearNo,
			UnlockTargetType:  targetType,
			TargetStageCode:   nullableString(targetStageCode),
			EditableStageCode: buildEditableStageCode(targetType, nextState),
			YearStatus:        nextState.YearStatus,
			StageStatus:       nextState.StageStatus,
			ReportStatus:      nextState.ReportStatus,
			SummaryEffective:  nextState.SummaryEffective,
			BusinessStatus:    nextState.BusinessStatus,
			UnlockLogID:       unlockLogItem.ID,
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

func validateInitializeGameInput(groupCount int) error {
	if groupCount < 1 || groupCount > initializeGameMaxGroupCount {
		return &InitializeInvalidError{Reason: fmt.Sprintf("小组数量必须在 1 到 %d 之间", initializeGameMaxGroupCount)}
	}
	return nil
}

func validateInitializeGameEdition(editionCode string) (GameEdition, error) {
	if strings.TrimSpace(editionCode) == "" {
		return GameEdition{}, ErrAdminControlEditionRequired
	}
	edition, ok := FindGameEdition(editionCode)
	if !ok {
		return GameEdition{}, ErrAdminControlEditionInvalid
	}
	return edition, nil
}

func resolveInitializeDictionaryItems(
	ctx context.Context,
	repo *repository.DictionaryRepository,
	editionCode string,
	schemeID *int64,
	inputs []DictionaryItemInput,
) ([]DictionaryItemResult, *int64, error) {
	if schemeID != nil && *schemeID > 0 {
		scheme, err := repo.GetSchemeByID(ctx, *schemeID)
		if err != nil {
			return nil, nil, fmt.Errorf("load dictionary scheme: %w", err)
		}
		if scheme.EditionCode != editionCode {
			return nil, nil, ErrDictionarySchemeCrossEdition
		}
		if len(inputs) == 0 {
			items, err := repo.ListSchemeItems(ctx, scheme.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("load dictionary scheme items: %w", err)
			}
			return buildDictionaryItemResultsFromSchemeItems(items), &scheme.ID, nil
		}
	}

	items, err := buildDictionaryItemsForEdition(editionCode, inputs)
	if err != nil {
		return nil, nil, err
	}
	if schemeID != nil && *schemeID > 0 {
		return items, schemeID, nil
	}
	return items, nil, nil
}

func ensureInitializeGameEnvironmentClean(existingGroupAccountCount int64, existingYearStateCount int64) error {
	if existingGroupAccountCount > 0 {
		return &InitializeInvalidError{Reason: "检测到未清理的玩家账号数据，请清理后再初始化比赛"}
	}
	if existingYearStateCount > 0 {
		return &InitializeInvalidError{Reason: "检测到未清理的年份状态数据，请清理后再初始化比赛"}
	}
	return nil
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
	if err := validateBaselineManualIntegers(*baselinePayload); err != nil {
		return err
	}
	return nil
}

func ensureUnlockYearAllowed(current state.RuntimeState, reason string, unlockTargetType string, targetStageCode string, nextYearAlreadyOpened bool) (state.RuntimeState, error) {
	if nextYearAlreadyOpened {
		return current, &UnlockNotAllowedError{Reason: unlockYearBlockedReasonNextYearOpened}
	}

	switch unlockTargetType {
	case "":
		return current, ErrAdminControlUnlockTargetTypeRequired
	case unlockTargetTypeOperating:
		if targetStageCode == "" {
			return current, ErrAdminControlUnlockStageRequired
		}
		if !state.IsValidStageCode(targetStageCode) {
			return current, ErrAdminControlUnlockStageInvalid
		}
		if isOperatingTargetAlreadyEditable(current, targetStageCode) {
			return current, &UnlockNotAllowedError{Reason: unlockYearBlockedReasonOperatingEditable}
		}
		if !hasSubmittedOperatingTarget(current, targetStageCode) {
			return current, &UnlockNotAllowedError{Reason: unlockYearBlockedReasonTargetNotSubmitted}
		}

		nextState, err := state.NewStateMachine().UnlockOperatingYear(current, targetStageCode, false)
		if err != nil {
			if errors.Is(err, state.ErrYearCannotUnlock) || errors.Is(err, state.ErrInvalidStageCode) {
				return current, &UnlockNotAllowedError{Reason: unlockYearBlockedReasonInvalidState}
			}
			return current, err
		}
		return nextState, nil
	case unlockTargetTypeReport:
		if isReportTargetAlreadyEditable(current) {
			return current, &UnlockNotAllowedError{Reason: unlockYearBlockedReasonReportEditable}
		}
		if !isReportTargetSubmitted(current) {
			return current, &UnlockNotAllowedError{Reason: unlockYearBlockedReasonTargetNotSubmitted}
		}

		nextState, err := state.NewStateMachine().UnlockReportYear(current, false)
		if err != nil {
			if errors.Is(err, state.ErrYearCannotUnlock) {
				return current, &UnlockNotAllowedError{Reason: unlockYearBlockedReasonInvalidState}
			}
			return current, err
		}
		return nextState, nil
	default:
		return current, ErrAdminControlUnlockTargetTypeInvalid
	}
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

func buildInitializeGameStateSnapshot(groupCount int, currentOpenYear int, finalYear int, baselineSubmitted bool, edition GameEdition) ([]byte, error) {
	return json.Marshal(map[string]any{
		"initialized":              groupCount > 0,
		"groupCount":               groupCount,
		"currentOpenYear":          currentOpenYear,
		"finalYear":                finalYear,
		"initialBaselineSubmitted": baselineSubmitted,
		"editionCode":              edition.EditionCode,
		"editionName":              edition.EditionName,
		"ruleVersion":              edition.RuleVersion,
		"templateVersion":          edition.TemplateVersion,
	})
}

func normalizeAdminOperatorName(operatorName string) string {
	operatorName = strings.TrimSpace(operatorName)
	if operatorName == "" {
		return "admin"
	}
	return operatorName
}
func normalizeUnlockTargetType(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func normalizeUnlockTargetStageCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func buildEditableStageCode(unlockTargetType string, nextState state.RuntimeState) *string {
	if unlockTargetType != unlockTargetTypeOperating {
		return nil
	}
	stageCode := state.CurrentStageCode(nextState.StageStatus)
	if stageCode == "" {
		return nil
	}
	return &stageCode
}

func nullableString(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}

func isOperatingTargetAlreadyEditable(current state.RuntimeState, targetStageCode string) bool {
	guard := state.NewTransitionGuard()
	permission := guard.BuildOperatingPermission(current)
	return permission.CanEdit && permission.CurrentStageCode == targetStageCode
}

func isReportTargetAlreadyEditable(current state.RuntimeState) bool {
	guard := state.NewTransitionGuard()
	return guard.BuildReportPermission(current).CanEdit
}

func hasSubmittedOperatingTarget(current state.RuntimeState, targetStageCode string) bool {
	targetRank, ok := stageCodeRank(targetStageCode)
	if !ok {
		return false
	}
	return targetRank <= maxSubmittedOperatingStageRank(current)
}

func isReportTargetSubmitted(current state.RuntimeState) bool {
	return current.ReportStatus == enum.ReportStatusSubmitted || current.YearStatus == enum.YearStatusCompleted
}

func maxSubmittedOperatingStageRank(current state.RuntimeState) int {
	switch current.YearStatus {
	case enum.YearStatusReportPending, enum.YearStatusReporting, enum.YearStatusCompleted:
		return 5
	case enum.YearStatusOperating:
		switch current.StageStatus {
		case enum.StageStatusQ1Open:
			return 0
		case enum.StageStatusQ2Open:
			return 1
		case enum.StageStatusQ3Open:
			return 2
		case enum.StageStatusQ4Open:
			return 3
		case enum.StageStatusYearEndOpen:
			return 4
		default:
			return 0
		}
	default:
		return 0
	}
}

func stageCodeRank(stageCode string) (int, bool) {
	switch stageCode {
	case state.StageCodeQ1:
		return 1, true
	case state.StageCodeQ2:
		return 2, true
	case state.StageCodeQ3:
		return 3, true
	case state.StageCodeQ4:
		return 4, true
	case state.StageCodeYearEnd:
		return 5, true
	default:
		return 0, false
	}
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

func buildInitializeGameGroups(groupCount int, operatorName string, operateTime time.Time) []entity.Group {
	items := make([]entity.Group, 0, groupCount)
	for groupNo := 1; groupNo <= groupCount; groupNo++ {
		items = append(items, entity.Group{
			GroupNo:        groupNo,
			GroupCode:      fmt.Sprintf("GROUP_%02d", groupNo),
			GroupName:      buildInitializeGroupName(groupNo),
			BusinessStatus: enum.BusinessStatusNormal,
			BaseEntity: entity.BaseEntity{
				Creator:    operatorName,
				CreateTime: operateTime,
				Updater:    operatorName,
				UpdateTime: operateTime,
			},
		})
	}
	return items
}

func buildInitializeGameAccounts(groups []entity.Group, operatorName string, operateTime time.Time) []entity.Account {
	items := make([]entity.Account, 0, len(groups))
	for _, group := range groups {
		groupID := group.ID
		items = append(items, entity.Account{
			Username:     fmt.Sprintf("group%02d", group.GroupNo),
			PasswordHash: hashSHA256Password(initializeGameDefaultPassword),
			RoleType:     enum.RoleTypeGroup,
			GroupID:      &groupID,
			Status:       enum.AccountStatusEnabled,
			BaseEntity: entity.BaseEntity{
				Creator:    operatorName,
				CreateTime: operateTime,
				Updater:    operatorName,
				UpdateTime: operateTime,
			},
		})
	}
	return items
}

func buildInitializeGameYearStates(groups []entity.Group, finalYear int, operatorName string, operateTime time.Time) []entity.GroupYearState {
	if len(groups) == 0 {
		return nil
	}

	maxYear := finalYear
	if maxYear < 0 {
		maxYear = 0
	}

	items := make([]entity.GroupYearState, 0, len(groups)*(maxYear+1))
	for _, group := range groups {
		for yearNo := 0; yearNo <= maxYear; yearNo++ {
			yearType := enum.YearTypeFormal
			yearStatus := enum.YearStatusLocked
			if yearNo == 0 {
				yearType = enum.YearTypeDemo
				yearStatus = enum.YearStatusOperating
			}
			items = append(items, entity.GroupYearState{
				GroupID:                   group.ID,
				YearNo:                    yearNo,
				YearType:                  yearType,
				YearStatus:                yearStatus,
				StageStatus:               enum.StageStatusQ1Open,
				ReportStatus:              enum.ReportStatusLocked,
				SummaryEffective:          false,
				LatestStageSubmitVersion:  0,
				LatestReportSubmitVersion: 0,
				BaseEntity: entity.BaseEntity{
					Creator:    operatorName,
					CreateTime: operateTime,
					Updater:    operatorName,
					UpdateTime: operateTime,
				},
			})
		}
	}
	return items
}

func buildInitializeGroupName(groupNo int) string {
	chinese := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}
	if groupNo >= 1 && groupNo <= 10 {
		return "第" + chinese[groupNo] + "组"
	}
	return fmt.Sprintf("第%d组", groupNo)
}
