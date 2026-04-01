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
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	carryforwardrules "sandbox-game/internal/rules/carryforward"
	calcctx "sandbox-game/internal/rules/context"
	operatingrules "sandbox-game/internal/rules/operating"
	reportrules "sandbox-game/internal/rules/report"
	"sandbox-game/internal/state"
)

type AdminGroupDataQueryService struct {
	gameConfigRepo      *repository.GameConfigRepository
	groupRepo           *repository.GroupRepository
	groupYearRepo       *repository.GroupYearStateRepository
	operatingRepo       *repository.OperatingRepository
	initialBaselineRepo *repository.InitialBaselineRepository
	reportRepo          *repository.ReportRepository
	operatingAssembler  *assembler.PlayerOperatingAssembler
	reportAssembler     *assembler.PlayerReportAssembler
	operatingCalculator *operatingrules.Calculator
	reportCalculator    *reportrules.Calculator
	carryForward        *carryforwardrules.Builder
	playerNoticeService *PlayerNoticeService
}

type AdminGroupOption struct {
	GroupID        int64  `json:"groupId"`
	GroupNo        int    `json:"groupNo"`
	GroupName      string `json:"groupName"`
	BusinessStatus string `json:"businessStatus"`
}

type ListAdminGroupsResult struct {
	List []AdminGroupOption `json:"list"`
}

func NewAdminGroupDataQueryService(
	gameConfigRepo *repository.GameConfigRepository,
	groupRepo *repository.GroupRepository,
	groupYearRepo *repository.GroupYearStateRepository,
	operatingRepo *repository.OperatingRepository,
	initialBaselineRepo *repository.InitialBaselineRepository,
	reportRepo *repository.ReportRepository,
	operatingAssembler *assembler.PlayerOperatingAssembler,
	reportAssembler *assembler.PlayerReportAssembler,
	playerNoticeService *PlayerNoticeService,
) *AdminGroupDataQueryService {
	return &AdminGroupDataQueryService{
		gameConfigRepo:      gameConfigRepo,
		groupRepo:           groupRepo,
		groupYearRepo:       groupYearRepo,
		operatingRepo:       operatingRepo,
		initialBaselineRepo: initialBaselineRepo,
		reportRepo:          reportRepo,
		operatingAssembler:  operatingAssembler,
		reportAssembler:     reportAssembler,
		operatingCalculator: operatingrules.NewCalculator(),
		reportCalculator:    reportrules.NewCalculator(),
		carryForward:        carryforwardrules.NewBuilder(),
		playerNoticeService: playerNoticeService,
	}
}

func (s *AdminGroupDataQueryService) ListGroups(ctx context.Context) (*ListAdminGroupsResult, error) {
	groups, err := s.groupRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("load groups: %w", err)
	}

	result := make([]AdminGroupOption, 0, len(groups))
	for _, item := range groups {
		result = append(result, AdminGroupOption{
			GroupID:        item.ID,
			GroupNo:        item.GroupNo,
			GroupName:      item.GroupName,
			BusinessStatus: item.BusinessStatus,
		})
	}

	return &ListAdminGroupsResult{List: result}, nil
}

func (s *AdminGroupDataQueryService) GetOperatingView(ctx context.Context, groupID int64, yearNo int) (*assembler.PlayerOperatingView, error) {
	calculationContext, draft, err := s.loadOperatingCalculationContext(ctx, groupID, yearNo)
	if err != nil {
		return nil, err
	}

	operatingResult, err := s.operatingCalculator.Calculate(calculationContext)
	if err != nil {
		return nil, fmt.Errorf("calculate operating view: %w", err)
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

	return s.operatingAssembler.Build(
		calculationContext,
		draft,
		stageSubmissions,
		operatingResult,
		buildAdminReadonlyOperatingPermission(calculationContext.State),
		carryForward,
		nil,
	), nil
}

func (s *AdminGroupDataQueryService) GetReportView(ctx context.Context, groupID int64, yearNo int) (*assembler.PlayerReportView, error) {
	calculationContext, err := s.loadReportCalculationContext(ctx, groupID, yearNo)
	if err != nil {
		return nil, err
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

	calculationContext = calculationContext.WithReportManualPayload(&manualPayload)
	calculatedPayload, err := s.reportCalculator.Calculate(calculationContext)
	if err != nil {
		return nil, fmt.Errorf("calculate report payload: %w", err)
	}

	if calculationContext.State.ReportStatus != enum.ReportStatusSubmitted || isZeroReportComputedPayload(computedPayload) {
		computedPayload = calculatedPayload
	}

	return s.reportAssembler.Build(
		calculationContext,
		computedPayload,
		manualPayload,
		lastDraftSavedAt,
		buildAdminReadonlyReportPermission(),
		nil,
	), nil
}

func (s *AdminGroupDataQueryService) loadOperatingCalculationContext(ctx context.Context, groupID int64, yearNo int) (calcctx.CalculationContext, *entity.GroupOperatingDraft, error) {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return calcctx.CalculationContext{}, nil, fmt.Errorf("load group: %w", err)
	}

	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return calcctx.CalculationContext{}, nil, fmt.Errorf("load game config: %w", err)
	}

	yearState, err := s.groupYearRepo.GetByGroupIDAndYear(ctx, groupID, yearNo)
	if err != nil {
		return calcctx.CalculationContext{}, nil, fmt.Errorf("load group year state: %w", err)
	}

	calculationContext := calcctx.NewCalculationContext(*group, *yearState, *gameConfig)

	operatingPayload := payload.NewOperatingPayload()
	draft, err := s.operatingRepo.FindDraft(ctx, groupID, yearNo)
	switch {
	case err == nil:
		if len(draft.OperatingPayload) > 0 {
			if unmarshalErr := json.Unmarshal(draft.OperatingPayload, &operatingPayload); unmarshalErr != nil {
				return calcctx.CalculationContext{}, nil, fmt.Errorf("unmarshal operating draft: %w", unmarshalErr)
			}
		}
	case errors.Is(err, gorm.ErrRecordNotFound):
		draft = nil
	default:
		return calcctx.CalculationContext{}, nil, fmt.Errorf("load operating draft: %w", err)
	}

	operatingPayload = operatingPayload.Normalize()
	operatingPayload, err = s.playerNoticeService.OverlayAdjustments(ctx, groupID, yearNo, operatingPayload)
	if err != nil {
		return calcctx.CalculationContext{}, nil, fmt.Errorf("overlay operating adjustments: %w", err)
	}
	calculationContext = calculationContext.WithOperatingPayload(&operatingPayload)

	calculationContext, err = s.attachCarrySource(ctx, calculationContext, groupID, yearNo)
	if err != nil {
		return calcctx.CalculationContext{}, nil, err
	}

	return calculationContext, draft, nil
}

func (s *AdminGroupDataQueryService) loadReportCalculationContext(ctx context.Context, groupID int64, yearNo int) (calcctx.CalculationContext, error) {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return calcctx.CalculationContext{}, fmt.Errorf("load group: %w", err)
	}

	gameConfig, err := s.gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return calcctx.CalculationContext{}, fmt.Errorf("load game config: %w", err)
	}

	yearState, err := s.groupYearRepo.GetByGroupIDAndYear(ctx, groupID, yearNo)
	if err != nil {
		return calcctx.CalculationContext{}, fmt.Errorf("load group year state: %w", err)
	}

	calculationContext := calcctx.NewCalculationContext(*group, *yearState, *gameConfig)

	operatingPayload := payload.NewOperatingPayload()
	draft, draftErr := s.operatingRepo.FindDraft(ctx, groupID, yearNo)
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

	operatingPayload = operatingPayload.Normalize()
	operatingPayload, err = s.playerNoticeService.OverlayAdjustments(ctx, groupID, yearNo, operatingPayload)
	if err != nil {
		return calcctx.CalculationContext{}, fmt.Errorf("overlay operating adjustments: %w", err)
	}
	calculationContext = calculationContext.WithOperatingPayload(&operatingPayload)

	calculationContext, err = s.attachCarrySource(ctx, calculationContext, groupID, yearNo)
	if err != nil {
		return calcctx.CalculationContext{}, err
	}

	return calculationContext, nil
}

func (s *AdminGroupDataQueryService) attachCarrySource(ctx context.Context, calculationContext calcctx.CalculationContext, groupID int64, yearNo int) (calcctx.CalculationContext, error) {
	if yearNo == 0 {
		baseline, baselineErr := s.initialBaselineRepo.FindByGroupID(ctx, groupID)
		switch {
		case baselineErr == nil:
			if len(baseline.BaselinePayload) > 0 {
				var baselinePayload payload.BaselinePayload
				if unmarshalErr := json.Unmarshal(baseline.BaselinePayload, &baselinePayload); unmarshalErr != nil {
					return calcctx.CalculationContext{}, fmt.Errorf("unmarshal initial baseline: %w", unmarshalErr)
				}
				calculationContext = calculationContext.WithInitialBaseline(&baselinePayload)
			}
		case errors.Is(baselineErr, gorm.ErrRecordNotFound):
		default:
			return calcctx.CalculationContext{}, fmt.Errorf("load initial baseline: %w", baselineErr)
		}
		return calculationContext, nil
	}

	previousReport, previousReportErr := s.reportRepo.FindEffectiveByGroupIDAndYear(ctx, groupID, yearNo-1)
	switch {
	case previousReportErr == nil:
		if len(previousReport.ReportComputedPayload) > 0 {
			var previous payload.ReportComputedPayload
			if unmarshalErr := json.Unmarshal(previousReport.ReportComputedPayload, &previous); unmarshalErr != nil {
				return calcctx.CalculationContext{}, fmt.Errorf("unmarshal previous report: %w", unmarshalErr)
			}
			calculationContext = calculationContext.WithPreviousReport(&previous)
		}
	case errors.Is(previousReportErr, gorm.ErrRecordNotFound):
	default:
		return calcctx.CalculationContext{}, fmt.Errorf("load previous report: %w", previousReportErr)
	}

	return calculationContext, nil
}

func buildAdminReadonlyOperatingPermission(current state.RuntimeState) state.OperatingPermission {
	return state.OperatingPermission{
		CanView:          true,
		CanEdit:          false,
		CanSubmit:        false,
		CurrentStageCode: state.CurrentStageCode(current.StageStatus),
		EditableScopes:   []string{},
		ReadonlyScopes: []string{
			state.OperatingScopeYearStart,
			state.OperatingScopeQ1,
			state.OperatingScopeQ2,
			state.OperatingScopeQ3,
			state.OperatingScopeQ4,
			state.OperatingScopeYearEnd,
		},
	}
}

func buildAdminReadonlyReportPermission() state.ReportPermission {
	return state.ReportPermission{
		CanView:   true,
		CanEdit:   false,
		CanSubmit: false,
	}
}
