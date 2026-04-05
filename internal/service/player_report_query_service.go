package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/assembler"
	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	calcctx "sandbox-game/internal/rules/context"
	reportrules "sandbox-game/internal/rules/report"
	"sandbox-game/internal/state"
)

var ErrPlayerReportNotOpen = errors.New("player report not open")

type PlayerReportQueryService struct {
	gameConfigRepo      *repository.GameConfigRepository
	groupRepo           *repository.GroupRepository
	groupYearRepo       *repository.GroupYearStateRepository
	operatingRepo       *repository.OperatingRepository
	initialBaseRepo     *repository.InitialBaselineRepository
	reportRepo          *repository.ReportRepository
	assembler           *assembler.PlayerReportAssembler
	calculator          *reportrules.Calculator
	transitionGuard     *state.TransitionGuard
	playerNoticeService *PlayerNoticeService
}

func NewPlayerReportQueryService(
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	operatingRepo *repository.OperatingRepository,
	initialBaseRepo *repository.InitialBaselineRepository,
	reportRepo *repository.ReportRepository,
	assembler *assembler.PlayerReportAssembler,
	playerNoticeService *PlayerNoticeService,
) *PlayerReportQueryService {
	return &PlayerReportQueryService{
		gameConfigRepo:      gameConfigRepo,
		groupRepo:           groupRepo,
		groupYearRepo:       groupYearRepo,
		operatingRepo:       operatingRepo,
		initialBaseRepo:     initialBaseRepo,
		reportRepo:          reportRepo,
		assembler:           assembler,
		calculator:          reportrules.NewCalculator(),
		transitionGuard:     state.NewTransitionGuard(),
		playerNoticeService: playerNoticeService,
	}
}

func (s *PlayerReportQueryService) GetView(ctx context.Context, groupID int64, yearNo int) (*assembler.PlayerReportView, error) {
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

	calcContext := calcctx.NewCalculationContext(*group, *yearState, *gameConfig)
	permission := s.transitionGuard.BuildReportPermission(calcContext.State)
	if !permission.CanView {
		return nil, ErrPlayerReportNotOpen
	}

	operatingPayload := payload.NewOperatingPayload()
	draft, draftErr := s.operatingRepo.FindDraft(ctx, groupID, yearNo)
	switch {
	case draftErr == nil:
		if len(draft.OperatingPayload) > 0 {
			if unmarshalErr := json.Unmarshal(draft.OperatingPayload, &operatingPayload); unmarshalErr != nil {
				return nil, fmt.Errorf("unmarshal operating draft: %w", unmarshalErr)
			}
		}
	case errors.Is(draftErr, gorm.ErrRecordNotFound):
	default:
		return nil, fmt.Errorf("load operating draft: %w", draftErr)
	}
	operatingPayload = operatingPayload.Normalize().WithoutDerivedValues()
	operatingPayload, err = s.playerNoticeService.OverlayAdjustments(ctx, groupID, yearNo, operatingPayload)
	if err != nil {
		return nil, fmt.Errorf("overlay operating adjustments: %w", err)
	}

	if yearNo == 0 {
		baselinePayload, baselineErr := loadInitialBaselinePayload(ctx, s.initialBaseRepo, groupID)
		if baselineErr != nil {
			return nil, baselineErr
		}
		if baselinePayload != nil {
			operatingPayload = applyInitialBaselineDefaultsToOperatingPayload(operatingPayload, baselinePayload)
			calcContext = calcContext.WithInitialBaseline(baselinePayload)
		}
	} else {
		previous, previousErr := loadEffectiveReportWithDirectorScores(
			ctx,
			*group,
			*gameConfig,
			yearNo-1,
			s.groupYearRepo,
			s.operatingRepo,
			s.initialBaseRepo,
			s.reportRepo,
			s.playerNoticeService,
			s.calculator,
		)
		if previousErr != nil {
			return nil, fmt.Errorf("load previous report: %w", previousErr)
		}
		if previous != nil {
			calcContext = calcContext.WithPreviousReport(previous)
		}
	}
	calcContext = calcContext.WithOperatingPayload(&operatingPayload)

	manualPayload := payload.ReportManualPayload{}
	computedPayload := payload.ReportComputedPayload{}
	var lastDraftSavedAt *time.Time
	report, reportErr := s.reportRepo.FindByGroupIDAndYear(ctx, groupID, yearNo)
	switch {
	case reportErr == nil:
		if len(report.ReportManualPayload) > 0 {
			if unmarshalErr := json.Unmarshal(report.ReportManualPayload, &manualPayload); unmarshalErr != nil {
				return nil, fmt.Errorf("unmarshal report manual payload: %w", unmarshalErr)
			}
		}
		if len(report.ReportComputedPayload) > 0 {
			if unmarshalErr := json.Unmarshal(report.ReportComputedPayload, &computedPayload); unmarshalErr != nil {
				return nil, fmt.Errorf("unmarshal report computed payload: %w", unmarshalErr)
			}
		}
		lastDraftSavedAt = report.LastAutoSavedAt
	case errors.Is(reportErr, gorm.ErrRecordNotFound):
	default:
		return nil, fmt.Errorf("load group report: %w", reportErr)
	}

	calcContext = calcContext.WithReportManualPayload(&manualPayload)
	calculatedPayload, err := s.calculator.Calculate(calcContext)
	if err != nil {
		return nil, fmt.Errorf("calculate report payload: %w", err)
	}

	if calcContext.State.ReportStatus != enum.ReportStatusSubmitted || isZeroReportComputedPayload(computedPayload) {
		computedPayload = calculatedPayload
	} else {
		overlayDirectorScores(&computedPayload, calculatedPayload)
	}

	noticeBoard, err := s.playerNoticeService.BuildBoard(ctx, groupID, yearNo)
	if err != nil {
		return nil, fmt.Errorf("build report notice board: %w", err)
	}

	return s.assembler.Build(calcContext, computedPayload, manualPayload, lastDraftSavedAt, permission, noticeBoard), nil
}
