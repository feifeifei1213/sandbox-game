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
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	calcctx "sandbox-game/internal/rules/context"
	reportrules "sandbox-game/internal/rules/report"
	summaryrules "sandbox-game/internal/rules/summary"
	"sandbox-game/internal/state"
)

var (
	ErrPlayerReportDraftNotEditable = errors.New("player report draft is not editable")
	ErrPlayerReportDraftInvalid     = errors.New("player report draft invalid")
	ErrPlayerReportSubmitConflict   = errors.New("player report submit conflict")
	ErrPlayerReportSubmitInvalid    = errors.New("player report submit invalid")
)

type SavePlayerReportDraftCommand struct {
	GroupID             int64
	YearNo              int
	ReportManualPayload payload.ReportManualPayload
	OperatorName        string
}

type SavePlayerReportDraftResult struct {
	GroupID          int64     `json:"groupId"`
	YearNo           int       `json:"yearNo"`
	YearStatus       string    `json:"yearStatus"`
	ReportStatus     string    `json:"reportStatus"`
	LastDraftSavedAt time.Time `json:"lastDraftSavedAt"`
}

type SubmitPlayerReportCommand struct {
	GroupID             int64
	YearNo              int
	ReportManualPayload payload.ReportManualPayload
	SubmitterID         int64
	OperatorName        string
}

type SubmitPlayerReportResult struct {
	GroupID                   int64     `json:"groupId"`
	YearNo                    int       `json:"yearNo"`
	YearStatus                string    `json:"yearStatus"`
	ReportStatus              string    `json:"reportStatus"`
	BusinessStatus            string    `json:"businessStatus"`
	SummaryEffective          bool      `json:"summaryEffective"`
	LatestReportSubmitVersion int       `json:"latestReportSubmitVersion"`
	BalanceCheckPassed        bool      `json:"balanceCheckPassed"`
	SubmittedAt               time.Time `json:"submittedAt"`
}

type PlayerReportCommandService struct {
	db                  *gorm.DB
	gameConfigRepo      *repository.GameConfigRepository
	groupRepo           *repository.GroupRepository
	groupYearRepo       *repository.GroupYearStateRepository
	operatingRepo       *repository.OperatingRepository
	initialBaseRepo     *repository.InitialBaselineRepository
	reportRepo          *repository.ReportRepository
	transitionGuard     *state.TransitionGuard
	validator           *reportrules.Validator
	calculator          *reportrules.Calculator
	summaryBuilder      *summaryrules.Builder
	playerNoticeService *PlayerNoticeService
}

func NewPlayerReportCommandService(
	db *gorm.DB,
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	operatingRepo *repository.OperatingRepository,
	initialBaseRepo *repository.InitialBaselineRepository,
	reportRepo *repository.ReportRepository,
	playerNoticeService *PlayerNoticeService,
) *PlayerReportCommandService {
	return &PlayerReportCommandService{
		db:                  db,
		gameConfigRepo:      gameConfigRepo,
		groupRepo:           groupRepo,
		groupYearRepo:       groupYearRepo,
		operatingRepo:       operatingRepo,
		initialBaseRepo:     initialBaseRepo,
		reportRepo:          reportRepo,
		transitionGuard:     state.NewTransitionGuard(),
		validator:           reportrules.NewValidator(),
		calculator:          reportrules.NewCalculator(),
		summaryBuilder:      summaryrules.NewBuilder(),
		playerNoticeService: playerNoticeService,
	}
}

func (s *PlayerReportCommandService) SaveDraft(ctx context.Context, cmd SavePlayerReportDraftCommand) (*SavePlayerReportDraftResult, error) {
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
	if err := calcContext.Validate(); err != nil {
		return nil, err
	}

	permission := s.transitionGuard.BuildReportPermission(calcContext.State)
	if !permission.CanEdit {
		return nil, ErrPlayerReportDraftNotEditable
	}

	validation := s.validator.ValidateDraft(&cmd.ReportManualPayload)
	if !validation.Passed {
		return nil, ErrPlayerReportDraftInvalid
	}

	now := time.Now()
	if err := s.reportRepo.UpsertDraft(ctx, repository.UpsertReportDraftCommand{
		GroupID:             cmd.GroupID,
		YearNo:              cmd.YearNo,
		ReportManualPayload: cmd.ReportManualPayload,
		LastAutoSavedAt:     now,
		OperatorName:        cmd.OperatorName,
	}); err != nil {
		return nil, fmt.Errorf("save player report draft: %w", err)
	}

	return &SavePlayerReportDraftResult{
		GroupID:          cmd.GroupID,
		YearNo:           cmd.YearNo,
		YearStatus:       calcContext.State.YearStatus,
		ReportStatus:     calcContext.State.ReportStatus,
		LastDraftSavedAt: now,
	}, nil
}

func (s *PlayerReportCommandService) Submit(ctx context.Context, cmd SubmitPlayerReportCommand) (*SubmitPlayerReportResult, error) {
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

	calcContext, err := s.buildCalculationContext(ctx, *group, *yearState, *gameConfig, cmd.ReportManualPayload)
	if err != nil {
		return nil, err
	}

	if calcContext.State.YearStatus == enum.YearStatusCompleted || calcContext.State.ReportStatus == enum.ReportStatusSubmitted {
		return nil, ErrPlayerReportSubmitConflict
	}
	if !s.transitionGuard.BuildReportPermission(calcContext.State).CanSubmit {
		return nil, ErrPlayerReportSubmitInvalid
	}

	computedPayload, err := s.calculator.Calculate(calcContext)
	if err != nil {
		return nil, fmt.Errorf("calculate report payload: %w", err)
	}

	validation := s.validator.ValidateSubmit(calcContext, computedPayload)
	if !validation.Passed {
		return nil, ErrPlayerReportSubmitInvalid
	}

	machine := state.NewStateMachine()
	stateBefore := state.BuildSnapshot(calcContext.State)
	nextState, err := machine.SubmitReport(calcContext.State)
	if err != nil {
		if errors.Is(err, state.ErrReportCannotSubmit) {
			return nil, ErrPlayerReportSubmitInvalid
		}
		return nil, err
	}

	bankruptTriggered := cmd.YearNo == gameConfig.FinalYear && computedPayload.ReportTotalEquity < 0
	bankruptReason := ""
	if bankruptTriggered {
		nextState.BusinessStatus = enum.BusinessStatusBankrupt
		bankruptReason = "final year equity below zero"
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
	submittedAt := submitTime

	var summaryResult summaryrules.SummaryResult
	if !calcContext.IsDemoYear() {
		nextContext, buildErr := buildNextReportContext(calcContext, nextState)
		if buildErr != nil {
			return nil, buildErr
		}
		summaryResult, err = s.summaryBuilder.Build(nextContext, computedPayload)
		if err != nil {
			return nil, fmt.Errorf("build summary snapshot: %w", err)
		}
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txReportRepo := repository.NewReportRepository(tx)
		txYearRepo := repository.NewGroupYearStateRepository(tx)
		txGroupRepo := repository.NewGroupRepository(tx)
		txSummaryRepo := repository.NewSummarySnapshotRepository(tx)

		if err := txReportRepo.UpsertCurrentReport(ctx, repository.UpsertCurrentReportCommand{
			GroupID:               cmd.GroupID,
			YearNo:                cmd.YearNo,
			ReportManualPayload:   cmd.ReportManualPayload,
			ReportComputedPayload: computedPayload,
			BalanceCheckPassed:    true,
			LastAutoSavedAt:       submitTime,
			SubmittedAt:           &submittedAt,
			OperatorName:          cmd.OperatorName,
		}); err != nil {
			return err
		}

		if err := txReportRepo.CreateSubmission(ctx, repository.CreateReportSubmissionCommand{
			GroupID:                cmd.GroupID,
			YearNo:                 cmd.YearNo,
			SubmitVersion:          nextState.LatestReportSubmitVersion,
			ReportManualSnapshot:   cmd.ReportManualPayload,
			ReportComputedSnapshot: computedPayload,
			BalanceCheckPassed:     true,
			StateBeforeJSON:        stateBeforeJSON,
			StateAfterJSON:         stateAfterJSON,
			SubmitterID:            cmd.SubmitterID,
			SubmitTime:             submitTime,
		}); err != nil {
			return err
		}

		if err := txYearRepo.UpdateRuntimeState(ctx, cmd.GroupID, cmd.YearNo, nextState, cmd.OperatorName); err != nil {
			return err
		}

		if !calcContext.IsDemoYear() {
			if err := txSummaryRepo.UpsertSnapshot(ctx, repository.UpsertSummarySnapshotCommand{
				GroupID:                   cmd.GroupID,
				YearNo:                    cmd.YearNo,
				SummaryResult:             summaryResult,
				SourceReportSubmitVersion: nextState.LatestReportSubmitVersion,
				OperatorName:              cmd.OperatorName,
				OperateTime:               submitTime,
			}); err != nil {
				return err
			}
		}

		if bankruptTriggered {
			if err := txGroupRepo.MarkBankrupt(ctx, cmd.GroupID, cmd.YearNo, bankruptReason, cmd.OperatorName); err != nil {
				return err
			}
		}

		return nil
	}, &sql.TxOptions{}); err != nil {
		return nil, fmt.Errorf("submit report transaction: %w", err)
	}

	return &SubmitPlayerReportResult{
		GroupID:                   cmd.GroupID,
		YearNo:                    cmd.YearNo,
		YearStatus:                nextState.YearStatus,
		ReportStatus:              nextState.ReportStatus,
		BusinessStatus:            nextState.BusinessStatus,
		SummaryEffective:          nextState.SummaryEffective,
		LatestReportSubmitVersion: nextState.LatestReportSubmitVersion,
		BalanceCheckPassed:        true,
		SubmittedAt:               submitTime,
	}, nil
}

func (s *PlayerReportCommandService) buildCalculationContext(
	ctx context.Context,
	group entity.Group,
	yearState entity.GroupYearState,
	gameConfig entity.GameConfig,
	manualPayload payload.ReportManualPayload,
) (calcctx.CalculationContext, error) {
	calcContext := calcctx.NewCalculationContext(group, yearState, gameConfig).
		WithReportManualPayload(&manualPayload)

	operatingPayload := payload.NewOperatingPayload()
	draft, draftErr := s.operatingRepo.FindDraft(ctx, group.ID, yearState.YearNo)
	switch {
	case draftErr == nil:
		if len(draft.OperatingPayload) > 0 {
			if unmarshalErr := json.Unmarshal(draft.OperatingPayload, &operatingPayload); unmarshalErr != nil {
				return calcctx.CalculationContext{}, fmt.Errorf("unmarshal operating draft: %w", unmarshalErr)
			}
		}
	case errors.Is(draftErr, gorm.ErrRecordNotFound):
	default:
		return calcctx.CalculationContext{}, fmt.Errorf("load operating draft: %w", draftErr)
	}
	operatingPayload = operatingPayload.Normalize().WithoutDerivedValues()
	operatingPayload, err := s.playerNoticeService.OverlayAdjustments(ctx, group.ID, yearState.YearNo, operatingPayload)
	if err != nil {
		return calcctx.CalculationContext{}, fmt.Errorf("overlay operating adjustments: %w", err)
	}

	if yearState.YearNo == 0 {
		baselinePayload, baselineErr := loadInitialBaselinePayload(ctx, s.initialBaseRepo, group.ID)
		if baselineErr != nil {
			return calcctx.CalculationContext{}, baselineErr
		}
		if baselinePayload != nil {
			operatingPayload = applyInitialBaselineDefaultsToOperatingPayload(operatingPayload, baselinePayload)
			calcContext = calcContext.WithInitialBaseline(baselinePayload)
		}
	} else {
		previous, previousErr := loadEffectiveReportWithDirectorScores(
			ctx,
			group,
			gameConfig,
			yearState.YearNo-1,
			s.groupYearRepo,
			s.operatingRepo,
			s.initialBaseRepo,
			s.reportRepo,
			s.playerNoticeService,
			s.calculator,
		)
		if previousErr != nil {
			return calcctx.CalculationContext{}, fmt.Errorf("load previous report: %w", previousErr)
		}
		if previous != nil {
			calcContext = calcContext.WithPreviousReport(previous)
		}
	}
	calcContext = calcContext.WithOperatingPayload(&operatingPayload)

	if err := calcContext.Validate(); err != nil {
		return calcctx.CalculationContext{}, err
	}
	return calcContext, nil
}

func buildNextReportContext(
	current calcctx.CalculationContext,
	nextState state.RuntimeState,
) (calcctx.CalculationContext, error) {
	nextGroup := current.Group
	nextGroup.BusinessStatus = nextState.BusinessStatus

	nextYearState := current.YearState
	nextYearState.YearStatus = nextState.YearStatus
	nextYearState.StageStatus = nextState.StageStatus
	nextYearState.ReportStatus = nextState.ReportStatus
	nextYearState.SummaryEffective = nextState.SummaryEffective
	nextYearState.LatestStageSubmitVersion = nextState.LatestStageSubmitVersion
	nextYearState.LatestReportSubmitVersion = nextState.LatestReportSubmitVersion

	nextContext := calcctx.NewCalculationContext(nextGroup, nextYearState, current.GameConfig).
		WithOperatingPayload(current.OperatingPayload).
		WithReportManualPayload(current.ReportManualPayload)
	if current.InitialBaseline != nil {
		nextContext = nextContext.WithInitialBaseline(current.InitialBaseline)
	}
	if current.PreviousReport != nil {
		nextContext = nextContext.WithPreviousReport(current.PreviousReport)
	}
	if err := nextContext.Validate(); err != nil {
		return calcctx.CalculationContext{}, err
	}
	return nextContext, nil
}
