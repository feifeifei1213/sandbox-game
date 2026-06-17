package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	carryforwardrules "sandbox-game/internal/rules/carryforward"
	calcctx "sandbox-game/internal/rules/context"
	operatingrules "sandbox-game/internal/rules/operating"
	"sandbox-game/internal/state"
)

var (
	ErrOperatingDraftNotEditable      = errors.New("operating draft is not editable")
	ErrOperatingDraftStageConflict    = errors.New("operating draft stage conflict")
	ErrOperatingDraftMissingCarryBase = errors.New("operating draft missing carry forward source")
	ErrOperatingStageSubmitInvalid    = errors.New("operating stage submit invalid")
)

type SaveOperatingDraftCommand struct {
	GroupID          int64
	YearNo           int
	StageStatus      string
	OperatingPayload payload.OperatingPayload
	OperatorName     string
}

type SaveOperatingDraftResult struct {
	GroupID          int64     `json:"groupId"`
	YearNo           int       `json:"yearNo"`
	StageStatus      string    `json:"stageStatus"`
	LastDraftSavedAt time.Time `json:"lastDraftSavedAt"`
}

type SubmitOperatingStageCommand struct {
	GroupID          int64
	YearNo           int
	StageCode        string
	OperatingPayload payload.OperatingPayload
	SubmitterID      int64
	OperatorName     string
}

type SubmitOperatingStageResult struct {
	GroupID                  int64     `json:"groupId"`
	YearNo                   int       `json:"yearNo"`
	StageCode                string    `json:"stageCode"`
	YearStatus               string    `json:"yearStatus"`
	StageStatus              string    `json:"stageStatus"`
	ReportStatus             string    `json:"reportStatus"`
	BusinessStatus           string    `json:"businessStatus"`
	LatestStageSubmitVersion int       `json:"latestStageSubmitVersion"`
	PeriodEndCash            float64   `json:"periodEndCash"`
	SubmittedAt              time.Time `json:"submittedAt"`
}

type PlayerOperatingCommandService struct {
	db                  *gorm.DB
	gameConfigRepo      *repository.GameConfigRepository
	groupRepo           *repository.GroupRepository
	groupYearRepo       *repository.GroupYearStateRepository
	operatingRepo       *repository.OperatingRepository
	initialBaseRepo     *repository.InitialBaselineRepository
	reportRepo          *repository.ReportRepository
	transitionGuard     *state.TransitionGuard
	carryForward        *carryforwardrules.Builder
	validator           *operatingrules.Validator
	calculator          *operatingrules.Calculator
	playerNoticeService *PlayerNoticeService
	orderLinkService    *OrderOperatingLinkService
}

func NewPlayerOperatingCommandService(
	db *gorm.DB,
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	operatingRepo *repository.OperatingRepository,
	initialBaseRepo *repository.InitialBaselineRepository,
	reportRepo *repository.ReportRepository,
	playerNoticeService *PlayerNoticeService,
	orderLinkService *OrderOperatingLinkService,
) *PlayerOperatingCommandService {
	return &PlayerOperatingCommandService{
		db:                  db,
		gameConfigRepo:      gameConfigRepo,
		groupRepo:           groupRepo,
		groupYearRepo:       groupYearRepo,
		operatingRepo:       operatingRepo,
		initialBaseRepo:     initialBaseRepo,
		reportRepo:          reportRepo,
		transitionGuard:     state.NewTransitionGuard(),
		carryForward:        carryforwardrules.NewBuilder(),
		validator:           operatingrules.NewValidator(),
		calculator:          operatingrules.NewCalculator(),
		playerNoticeService: playerNoticeService,
		orderLinkService:    orderLinkService,
	}
}

func (s *PlayerOperatingCommandService) SaveDraft(ctx context.Context, cmd SaveOperatingDraftCommand) (*SaveOperatingDraftResult, error) {
	group, err := s.groupRepo.GetByID(ctx, cmd.GroupID)
	if err != nil {
		return nil, fmt.Errorf("load group: %w", err)
	}

	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}

	yearState, err := s.groupYearRepo.GetByGroupIDAndYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load group year state: %w", err)
	}

	calcContext := calcctx.NewCalculationContext(*group, *yearState, *gameConfig)
	if cmd.YearNo == 0 {
		baselinePayload, baselineErr := loadInitialBaselinePayload(ctx, s.initialBaseRepo, cmd.GroupID)
		if baselineErr != nil {
			return nil, baselineErr
		}
		if baselinePayload != nil {
			calcContext = calcContext.WithInitialBaseline(baselinePayload)
		}
	}
	if cmd.YearNo > 0 {
		previousReport, previousReportErr := s.reportRepo.FindEffectiveByGroupIDAndYear(ctx, cmd.GroupID, cmd.YearNo-1)
		if previousReportErr == nil && len(previousReport.ReportComputedPayload) > 0 {
			var previous payload.ReportComputedPayload
			if unmarshalErr := json.Unmarshal(previousReport.ReportComputedPayload, &previous); unmarshalErr != nil {
				return nil, fmt.Errorf("unmarshal previous report: %w", unmarshalErr)
			}
			calcContext = calcContext.WithPreviousReport(&previous)
		} else if previousReportErr != nil && !errors.Is(previousReportErr, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("load previous report: %w", previousReportErr)
		}
	}

	if err := calcContext.Validate(); err != nil {
		return nil, err
	}
	if !calcContext.HasCarryForwardSource() {
		return nil, ErrOperatingDraftMissingCarryBase
	}

	permission := s.transitionGuard.BuildOperatingPermission(calcContext.State)
	if !permission.CanEdit {
		return nil, ErrOperatingDraftNotEditable
	}
	if cmd.StageStatus != calcContext.State.StageStatus {
		return nil, ErrOperatingDraftStageConflict
	}
	if !enum.IsValidStageStatus(cmd.StageStatus) {
		return nil, ErrOperatingDraftStageConflict
	}

	normalizedPayload := cmd.OperatingPayload.Normalize().WithoutDerivedValues()
	if err := validateOperatingManualIntegers(normalizedPayload); err != nil {
		return nil, err
	}
	normalizedPayload, err = s.playerNoticeService.OverlayAdjustments(ctx, cmd.GroupID, cmd.YearNo, normalizedPayload)
	if err != nil {
		return nil, fmt.Errorf("overlay operating adjustments: %w", err)
	}
	if cmd.YearNo > 0 && s.orderLinkService != nil {
		linkedPayload, _, linkErr := s.orderLinkService.ApplyFormalYearValues(ctx, cmd.GroupID, cmd.YearNo, normalizedPayload)
		if linkErr != nil {
			return nil, fmt.Errorf("apply order operating values: %w", linkErr)
		}
		normalizedPayload = linkedPayload
	}
	if cmd.YearNo == 0 && calcContext.InitialBaseline != nil {
		normalizedPayload = applyInitialBaselineDefaultsToOperatingPayload(normalizedPayload, calcContext.InitialBaseline)
	}

	now := time.Now()
	if err := s.operatingRepo.UpsertDraft(ctx, repository.UpsertOperatingDraftCommand{
		GroupID:          cmd.GroupID,
		YearNo:           cmd.YearNo,
		StageStatus:      calcContext.State.StageStatus,
		OperatingPayload: normalizedPayload,
		LastAutoSavedAt:  now,
		OperatorName:     cmd.OperatorName,
	}); err != nil {
		return nil, fmt.Errorf("save operating draft: %w", err)
	}

	return &SaveOperatingDraftResult{
		GroupID:          cmd.GroupID,
		YearNo:           cmd.YearNo,
		StageStatus:      calcContext.State.StageStatus,
		LastDraftSavedAt: now,
	}, nil
}

func (s *PlayerOperatingCommandService) SubmitStage(ctx context.Context, cmd SubmitOperatingStageCommand) (*SubmitOperatingStageResult, error) {
	group, err := s.groupRepo.GetByID(ctx, cmd.GroupID)
	if err != nil {
		return nil, fmt.Errorf("load group: %w", err)
	}

	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}

	yearState, err := s.groupYearRepo.GetByGroupIDAndYear(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load group year state: %w", err)
	}

	normalizedPayload := cmd.OperatingPayload.Normalize().WithoutDerivedValues()
	if err := validateOperatingManualIntegers(normalizedPayload); err != nil {
		return nil, err
	}
	normalizedPayload, err = s.playerNoticeService.OverlayAdjustments(ctx, cmd.GroupID, cmd.YearNo, normalizedPayload)
	if err != nil {
		return nil, fmt.Errorf("overlay operating adjustments: %w", err)
	}
	orderPrerequisiteCompleted := cmd.YearNo == 0 || s.orderLinkService == nil
	if cmd.YearNo > 0 && s.orderLinkService != nil {
		linkedPayload, values, linkErr := s.orderLinkService.ApplyFormalYearValues(ctx, cmd.GroupID, cmd.YearNo, normalizedPayload)
		if linkErr != nil {
			return nil, fmt.Errorf("apply order operating values: %w", linkErr)
		}
		normalizedPayload = linkedPayload
		orderPrerequisiteCompleted = values.PrerequisiteCompleted
	}
	calcContext := calcctx.NewCalculationContext(*group, *yearState, *gameConfig).
		WithOperatingPayload(&normalizedPayload)

	if cmd.YearNo == 0 {
		baselinePayload, baselineErr := loadInitialBaselinePayload(ctx, s.initialBaseRepo, cmd.GroupID)
		if baselineErr != nil {
			return nil, baselineErr
		}
		if baselinePayload != nil {
			normalizedPayload = applyInitialBaselineDefaultsToOperatingPayload(normalizedPayload, baselinePayload)
			calcContext = calcContext.WithOperatingPayload(&normalizedPayload)
			calcContext = calcContext.WithInitialBaseline(baselinePayload)
		}
	}
	if cmd.YearNo > 0 {
		previousReport, previousReportErr := s.reportRepo.FindEffectiveByGroupIDAndYear(ctx, cmd.GroupID, cmd.YearNo-1)
		if previousReportErr == nil && len(previousReport.ReportComputedPayload) > 0 {
			var previous payload.ReportComputedPayload
			if unmarshalErr := json.Unmarshal(previousReport.ReportComputedPayload, &previous); unmarshalErr != nil {
				return nil, fmt.Errorf("unmarshal previous report: %w", unmarshalErr)
			}
			calcContext = calcContext.WithPreviousReport(&previous)
		} else if previousReportErr != nil && !errors.Is(previousReportErr, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("load previous report: %w", previousReportErr)
		}
	}

	if err := calcContext.Validate(); err != nil {
		return nil, err
	}
	if !state.IsValidStageCode(cmd.StageCode) {
		return nil, ErrOperatingStageSubmitInvalid
	}
	if cmd.YearNo > 0 && cmd.StageCode == state.StageCodeQ1 && !orderPrerequisiteCompleted {
		return nil, ErrOrderPrerequisiteIncomplete
	}

	validation := s.validator.ValidateStageSubmit(calcContext, cmd.StageCode)
	if !validation.Passed {
		if currentStageCode := state.CurrentStageCode(calcContext.State.StageStatus); currentStageCode != "" && currentStageCode != cmd.StageCode {
			return nil, ErrOperatingDraftStageConflict
		}
		return nil, ErrOperatingStageSubmitInvalid
	}

	calculationResult, err := s.calculator.Calculate(calcContext)
	if err != nil {
		return nil, fmt.Errorf("calculate stage result: %w", err)
	}

	machine := state.NewStateMachine()
	stateBefore := state.BuildSnapshot(calcContext.State)
	nextState, err := machine.SubmitStage(calcContext.State, cmd.StageCode)
	if err != nil {
		return nil, err
	}
	historyMaxStageVersion, err := s.operatingRepo.MaxStageSubmitVersion(ctx, cmd.GroupID, cmd.YearNo)
	if err != nil {
		return nil, fmt.Errorf("load max stage submit version: %w", err)
	}
	if nextState.LatestStageSubmitVersion <= historyMaxStageVersion {
		nextState.LatestStageSubmitVersion = historyMaxStageVersion + 1
	}

	bankruptTriggered := calculationResult.PeriodEndCash < 0
	bankruptReason := ""
	if bankruptTriggered {
		nextState.BusinessStatus = enum.BusinessStatusBankrupt
		bankruptReason = "quarter period-end cash below zero"
	}
	stateAfter := state.BuildSnapshot(nextState)

	stateBeforeJSON, err := json.Marshal(stateBefore)
	if err != nil {
		return nil, fmt.Errorf("marshal state before: %w", err)
	}
	stateAfterJSON, err := json.Marshal(stateAfter)
	if err != nil {
		return nil, fmt.Errorf("marshal state after: %w", err)
	}

	submitTime := time.Now()
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txOperatingRepo := repository.NewOperatingRepository(tx)
		txYearRepo := repository.NewGroupYearStateRepository(tx)
		txGroupRepo := repository.NewGroupRepository(tx)
		txOrderSelectionRepo := repository.NewGroupOrderSelectionRepository(tx)

		if err := txOperatingRepo.UpsertDraft(ctx, repository.UpsertOperatingDraftCommand{
			GroupID:          cmd.GroupID,
			YearNo:           cmd.YearNo,
			StageStatus:      nextState.StageStatus,
			OperatingPayload: calculationResult.OperatingPayload.Normalize(),
			LastAutoSavedAt:  submitTime,
			OperatorName:     cmd.OperatorName,
		}); err != nil {
			return err
		}

		if err := txOperatingRepo.CreateStageSubmission(ctx, repository.CreateStageSubmissionCommand{
			GroupID:                  cmd.GroupID,
			YearNo:                   cmd.YearNo,
			StageCode:                cmd.StageCode,
			SubmitVersion:            nextState.LatestStageSubmitVersion,
			PeriodEndCash:            calculationResult.PeriodEndCash,
			OperatingPayloadSnapshot: calculationResult.OperatingPayload.Normalize(),
			StateBeforeJSON:          stateBeforeJSON,
			StateAfterJSON:           stateAfterJSON,
			SubmitterID:              cmd.SubmitterID,
			SubmitTime:               submitTime,
		}); err != nil {
			return err
		}

		if err := txYearRepo.UpdateRuntimeState(ctx, cmd.GroupID, cmd.YearNo, nextState, cmd.OperatorName); err != nil {
			return err
		}

		if bankruptTriggered {
			if err := txGroupRepo.MarkBankrupt(ctx, cmd.GroupID, cmd.YearNo, bankruptReason, cmd.OperatorName); err != nil {
				return err
			}
		}
		if cmd.YearNo > 0 && cmd.StageCode == state.StageCodeYearEnd {
			if err := txOrderSelectionRepo.MarkUnfinishedByGroupYear(ctx, cmd.GroupID, cmd.YearNo, cmd.OperatorName, submitTime); err != nil {
				return err
			}
		}
		if _, err := createGroupSnapshot(ctx, tx, CreateGroupSnapshotCommand{
			GroupID:       cmd.GroupID,
			YearNo:        cmd.YearNo,
			StageCode:     cmd.StageCode,
			SnapshotType:  enum.SnapshotTypeAuto,
			TriggerCode:   stageSubmitSnapshotTrigger(yearState.RollbackPending),
			Description:   stageSubmitSnapshotDescription(yearState.RollbackPending),
			OperatorID:    cmd.SubmitterID,
			OperatorName:  cmd.OperatorName,
			OperateTime:   submitTime,
			UseForRestore: true,
		}); err != nil {
			return err
		}

		return nil
	}, sqlTxOptionsReadCommitted); err != nil {
		return nil, fmt.Errorf("submit stage transaction: %w", err)
	}

	businessStatus := nextState.BusinessStatus
	return &SubmitOperatingStageResult{
		GroupID:                  cmd.GroupID,
		YearNo:                   cmd.YearNo,
		StageCode:                cmd.StageCode,
		YearStatus:               nextState.YearStatus,
		StageStatus:              nextState.StageStatus,
		ReportStatus:             nextState.ReportStatus,
		BusinessStatus:           businessStatus,
		LatestStageSubmitVersion: nextState.LatestStageSubmitVersion,
		PeriodEndCash:            calculationResult.PeriodEndCash,
		SubmittedAt:              submitTime,
	}, nil
}

func stageSubmitSnapshotTrigger(rollbackPending bool) string {
	if rollbackPending {
		return enum.SnapshotTriggerRollbackStageRetry
	}
	return enum.SnapshotTriggerStageSubmitted
}

func stageSubmitSnapshotDescription(rollbackPending bool) string {
	if rollbackPending {
		return "回退后重新提交经营阶段自动快照"
	}
	return "经营阶段提交后自动快照"
}

var sqlTxOptionsReadCommitted = &sql.TxOptions{}
