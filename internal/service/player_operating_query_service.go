package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"sandbox-game/internal/assembler"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	carryforwardrules "sandbox-game/internal/rules/carryforward"
	calcctx "sandbox-game/internal/rules/context"
	operatingrules "sandbox-game/internal/rules/operating"
	reportrules "sandbox-game/internal/rules/report"
	"sandbox-game/internal/state"
)

type PlayerOperatingQueryService struct {
	gameConfigRepo      *repository.GameConfigRepository
	groupRepo           *repository.GroupRepository
	groupYearRepo       *repository.GroupYearStateRepository
	operatingRepo       *repository.OperatingRepository
	initialBaseRepo     *repository.InitialBaselineRepository
	reportRepo          *repository.ReportRepository
	assembler           *assembler.PlayerOperatingAssembler
	calculator          *operatingrules.Calculator
	reportCalculator    *reportrules.Calculator
	carryForward        *carryforwardrules.Builder
	transitionGuard     *state.TransitionGuard
	playerNoticeService *PlayerNoticeService
}

func NewPlayerOperatingQueryService(
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	operatingRepo *repository.OperatingRepository,
	initialBaseRepo *repository.InitialBaselineRepository,
	reportRepo *repository.ReportRepository,
	assembler *assembler.PlayerOperatingAssembler,
	playerNoticeService *PlayerNoticeService,
) *PlayerOperatingQueryService {
	return &PlayerOperatingQueryService{
		gameConfigRepo:      gameConfigRepo,
		groupRepo:           groupRepo,
		groupYearRepo:       groupYearRepo,
		operatingRepo:       operatingRepo,
		initialBaseRepo:     initialBaseRepo,
		reportRepo:          reportRepo,
		assembler:           assembler,
		calculator:          operatingrules.NewCalculator(),
		reportCalculator:    reportrules.NewCalculator(),
		carryForward:        carryforwardrules.NewBuilder(),
		transitionGuard:     state.NewTransitionGuard(),
		playerNoticeService: playerNoticeService,
	}
}

func (s *PlayerOperatingQueryService) GetYearView(ctx context.Context, groupID int64, yearNo int) (*assembler.PlayerOperatingView, error) {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("load group: %w", err)
	}

	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, fmt.Errorf("load game config: %w", err)
	}

	yearState, err := s.groupYearRepo.GetByGroupIDAndYear(ctx, groupID, yearNo)
	if err != nil {
		return nil, fmt.Errorf("load group year state: %w", err)
	}

	calculationContext := calcctx.NewCalculationContext(*group, *yearState, *gameConfig)

	operatingPayload := payload.NewOperatingPayload()
	draft, err := s.operatingRepo.FindDraft(ctx, groupID, yearNo)
	switch {
	case err == nil:
		if len(draft.OperatingPayload) > 0 {
			if unmarshalErr := json.Unmarshal(draft.OperatingPayload, &operatingPayload); unmarshalErr != nil {
				return nil, fmt.Errorf("unmarshal operating draft: %w", unmarshalErr)
			}
		}
	case errors.Is(err, gorm.ErrRecordNotFound):
	default:
		return nil, fmt.Errorf("load operating draft: %w", err)
	}

	operatingPayload = operatingPayload.Normalize().WithoutDerivedValues()
	operatingPayload, err = s.playerNoticeService.OverlayAdjustments(ctx, groupID, yearNo, operatingPayload)
	if err != nil {
		return nil, fmt.Errorf("overlay operating adjustments: %w", err)
	}
	calculationContext = calculationContext.WithOperatingPayload(&operatingPayload)

	if yearNo == 0 {
		baseline, baselineErr := s.initialBaseRepo.FindByGroupID(ctx, groupID)
		switch {
		case baselineErr == nil:
			if len(baseline.BaselinePayload) > 0 {
				var baselinePayload payload.BaselinePayload
				if unmarshalErr := json.Unmarshal(baseline.BaselinePayload, &baselinePayload); unmarshalErr != nil {
					return nil, fmt.Errorf("unmarshal initial baseline: %w", unmarshalErr)
				}
				calculationContext = calculationContext.WithInitialBaseline(&baselinePayload)
			}
		case errors.Is(baselineErr, gorm.ErrRecordNotFound):
		default:
			return nil, fmt.Errorf("load initial baseline: %w", baselineErr)
		}
	} else {
		previousReport, previousReportErr := s.reportRepo.FindEffectiveByGroupIDAndYear(ctx, groupID, yearNo-1)
		switch {
		case previousReportErr == nil:
			if len(previousReport.ReportComputedPayload) > 0 {
				var previous payload.ReportComputedPayload
				if unmarshalErr := json.Unmarshal(previousReport.ReportComputedPayload, &previous); unmarshalErr != nil {
					return nil, fmt.Errorf("unmarshal previous report: %w", unmarshalErr)
				}
				calculationContext = calculationContext.WithPreviousReport(&previous)
			}
		case errors.Is(previousReportErr, gorm.ErrRecordNotFound):
		default:
			return nil, fmt.Errorf("load previous report: %w", previousReportErr)
		}
	}

	operatingResult, err := s.calculator.Calculate(calculationContext)
	if err != nil {
		return nil, fmt.Errorf("calculate operating view: %w", err)
	}
	if err := enrichOperatingDerivedValuesWithReportMetrics(ctx, s.reportRepo, s.reportCalculator, calculationContext, groupID, yearNo, &operatingResult); err != nil {
		return nil, err
	}

	stageSubmissions, err := s.operatingRepo.ListStageSubmissions(ctx, groupID, yearNo)
	if err != nil {
		return nil, fmt.Errorf("load stage submissions: %w", err)
	}

	var carryForward *carryforwardrules.CarryForwardResult
	if calculationContext.HasCarryForwardSource() {
		result, buildErr := s.carryForward.Build(calculationContext)
		if buildErr != nil {
			return nil, fmt.Errorf("build carry forward: %w", buildErr)
		}
		carryForward = &result
	}

	noticeBoard, err := s.playerNoticeService.BuildBoard(ctx, groupID, yearNo)
	if err != nil {
		return nil, fmt.Errorf("build notice board: %w", err)
	}

	permission := s.transitionGuard.BuildOperatingPermission(calculationContext.State)

	return s.assembler.Build(
		calculationContext,
		draft,
		stageSubmissions,
		operatingResult,
		permission,
		carryForward,
		noticeBoard,
	), nil
}
