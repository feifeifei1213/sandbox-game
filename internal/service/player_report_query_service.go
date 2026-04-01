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
	calcContext = calcContext.WithOperatingPayload(&operatingPayload)

	if yearNo == 0 {
		baseline, baselineErr := s.initialBaseRepo.FindByGroupID(ctx, groupID)
		switch {
		case baselineErr == nil:
			if len(baseline.BaselinePayload) > 0 {
				var baselinePayload payload.BaselinePayload
				if unmarshalErr := json.Unmarshal(baseline.BaselinePayload, &baselinePayload); unmarshalErr != nil {
					return nil, fmt.Errorf("unmarshal initial baseline: %w", unmarshalErr)
				}
				calcContext = calcContext.WithInitialBaseline(&baselinePayload)
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
				calcContext = calcContext.WithPreviousReport(&previous)
			}
		case errors.Is(previousReportErr, gorm.ErrRecordNotFound):
		default:
			return nil, fmt.Errorf("load previous report: %w", previousReportErr)
		}
	}

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
	}

	noticeBoard, err := s.playerNoticeService.BuildBoard(ctx, groupID, yearNo)
	if err != nil {
		return nil, fmt.Errorf("build report notice board: %w", err)
	}

	return s.assembler.Build(calcContext, computedPayload, manualPayload, lastDraftSavedAt, permission, noticeBoard), nil
}

func isZeroReportComputedPayload(value payload.ReportComputedPayload) bool {
	return value == (payload.ReportComputedPayload{})
}
