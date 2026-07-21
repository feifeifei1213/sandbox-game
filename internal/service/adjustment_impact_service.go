package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	calcctx "sandbox-game/internal/rules/context"
	operatingrules "sandbox-game/internal/rules/operating"
	reportrules "sandbox-game/internal/rules/report"
	"sandbox-game/internal/state"
)

const (
	adjustmentOperationCreate  = "CREATE"
	adjustmentOperationVoid    = "VOID"
	adjustmentOperationCurrent = "CURRENT"
)

type AdjustmentImpactRequest struct {
	Operation string
	GroupID   int64
	YearNo    int
	Candidate *entity.GroupAdjustment
	VoidID    int64
}

type AdjustmentImpactResult struct {
	ResolvedStageCode       string                           `json:"resolvedStageCode"`
	CashBefore              float64                          `json:"cashBefore"`
	CashAfter               float64                          `json:"cashAfter"`
	PreTaxProfitAfter       float64                          `json:"preTaxProfitAfter"`
	IncomeTaxAfter          float64                          `json:"incomeTaxAfter"`
	NetProfitAfter          float64                          `json:"netProfitAfter"`
	TotalEquityAfter        float64                          `json:"totalEquityAfter"`
	WillBankrupt            bool                             `json:"willBankrupt"`
	CalculationBasisSavedAt *time.Time                       `json:"calculationBasisSavedAt"`
	OperatingAfter          operatingrules.CalculationResult `json:"-"`
	ReportAfter             payload.ReportComputedPayload    `json:"-"`
	ReportBefore            payload.ReportComputedPayload    `json:"-"`
	ReportManual            payload.ReportManualPayload      `json:"-"`
}

type adjustmentImpactCalculator struct {
	db *gorm.DB
}

func newAdjustmentImpactCalculator(db *gorm.DB) *adjustmentImpactCalculator {
	return &adjustmentImpactCalculator{db: db}
}

func (c *adjustmentImpactCalculator) Calculate(ctx context.Context, req AdjustmentImpactRequest) (*AdjustmentImpactResult, error) {
	groupRepo := repository.NewGroupRepository(c.db)
	yearRepo := repository.NewGroupYearStateRepository(c.db)
	gameConfigRepo := repository.NewGameConfigRepository(c.db)
	operatingRepo := repository.NewOperatingRepository(c.db)
	reportRepo := repository.NewReportRepository(c.db)
	initialBaseRepo := repository.NewInitialBaselineRepository(c.db)
	adjustmentRepo := repository.NewGroupAdjustmentRepository(c.db)

	group, err := groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	yearState, err := yearRepo.GetByGroupIDAndYear(ctx, req.GroupID, req.YearNo)
	if err != nil {
		return nil, err
	}
	resolvedStageCode := ""
	if req.Operation == adjustmentOperationCurrent {
		resolvedStageCode = "YEAR_END"
		if current := strings.TrimSpace(state.CurrentStageCode(yearState.StageStatus)); current != "" {
			resolvedStageCode = current
		}
	} else {
		resolvedStageCode, err = resolveAdjustmentStage(*yearState)
		if err != nil {
			return nil, err
		}
	}
	if req.Operation == adjustmentOperationCreate && req.Candidate != nil {
		candidate := *req.Candidate
		candidate.StageCode = resolvedStageCode
		req.Candidate = &candidate
	}
	gameConfig, err := gameConfigRepo.GetCurrent(ctx)
	if err != nil {
		return nil, err
	}

	active, err := adjustmentRepo.ListByGroupIDAndYear(ctx, req.GroupID, req.YearNo)
	if err != nil {
		return nil, fmt.Errorf("list active adjustments: %w", err)
	}
	afterAdjustments, err := applyAdjustmentOperation(active, req)
	if err != nil {
		return nil, err
	}

	operatingPayload := payload.NewOperatingPayload()
	var basisSavedAt *time.Time
	draft, draftErr := operatingRepo.FindDraft(ctx, req.GroupID, req.YearNo)
	switch {
	case draftErr == nil:
		if len(draft.OperatingPayload) > 0 {
			if err := json.Unmarshal(draft.OperatingPayload, &operatingPayload); err != nil {
				return nil, fmt.Errorf("unmarshal operating draft: %w", err)
			}
		}
		basisSavedAt = latestSavedAt(basisSavedAt, draft.LastAutoSavedAt)
	case errors.Is(draftErr, gorm.ErrRecordNotFound):
	default:
		return nil, fmt.Errorf("load operating draft: %w", draftErr)
	}
	operatingPayload = operatingPayload.Normalize().WithoutDerivedValues()

	if req.YearNo > 0 {
		orderLinkService := NewOrderOperatingLinkService(
			repository.NewGroupMarketBidRepository(c.db),
			repository.NewMarketBiddingStateRepository(c.db),
			repository.NewGroupOrderSelectionRepository(c.db),
		)
		linkedPayload, _, linkErr := orderLinkService.ApplyFormalYearValues(ctx, req.GroupID, req.YearNo, operatingPayload)
		if linkErr != nil {
			return nil, fmt.Errorf("apply order operating values: %w", linkErr)
		}
		operatingPayload = linkedPayload
	}

	reportManual := payload.ReportManualPayload{}
	reportItem, reportErr := reportRepo.FindByGroupIDAndYear(ctx, req.GroupID, req.YearNo)
	switch {
	case reportErr == nil:
		if len(reportItem.ReportManualPayload) > 0 {
			if err := json.Unmarshal(reportItem.ReportManualPayload, &reportManual); err != nil {
				return nil, fmt.Errorf("unmarshal report manual payload: %w", err)
			}
		}
		basisSavedAt = latestSavedAt(basisSavedAt, reportItem.LastAutoSavedAt)
	case errors.Is(reportErr, gorm.ErrRecordNotFound):
	default:
		return nil, fmt.Errorf("load report draft: %w", reportErr)
	}

	baseContext := calcctx.NewCalculationContext(*group, *yearState, *gameConfig).
		WithReportManualPayload(&reportManual)
	if req.YearNo == 0 {
		baseline, baselineErr := loadInitialBaselinePayload(ctx, initialBaseRepo, req.GroupID)
		if baselineErr != nil {
			return nil, baselineErr
		}
		if baseline != nil {
			operatingPayload = applyInitialBaselineDefaultsToOperatingPayload(operatingPayload, baseline)
			baseContext = baseContext.WithInitialBaseline(baseline)
		}
	} else {
		playerNoticeService := NewPlayerNoticeService(repository.NewNoticeRepository(c.db), adjustmentRepo)
		previous, previousErr := loadEffectiveReportWithDirectorScores(
			ctx, *group, *gameConfig, req.YearNo-1,
			yearRepo, operatingRepo, initialBaseRepo, reportRepo,
			playerNoticeService, reportrules.NewCalculator(),
		)
		if previousErr != nil {
			return nil, fmt.Errorf("load previous report: %w", previousErr)
		}
		if previous != nil {
			baseContext = baseContext.WithPreviousReport(previous)
		}
	}

	beforeBase, err := cloneAdjustmentOperatingPayload(operatingPayload)
	if err != nil {
		return nil, err
	}
	afterBase, err := cloneAdjustmentOperatingPayload(operatingPayload)
	if err != nil {
		return nil, err
	}
	beforeOperatingPayload := overlayAdjustmentValues(beforeBase, active)
	afterOperatingPayload := overlayAdjustmentValues(afterBase, afterAdjustments)
	operatingCalculator := operatingrules.NewCalculator()
	beforeOperating, err := operatingCalculator.Calculate(baseContext.WithOperatingPayload(&beforeOperatingPayload))
	if err != nil {
		return nil, fmt.Errorf("calculate operating impact before: %w", err)
	}
	afterOperating, err := operatingCalculator.Calculate(baseContext.WithOperatingPayload(&afterOperatingPayload))
	if err != nil {
		return nil, fmt.Errorf("calculate operating impact after: %w", err)
	}

	reportCalculator := reportrules.NewCalculator()
	beforeReport, err := reportCalculator.Calculate(baseContext.WithOperatingPayload(&beforeOperating.OperatingPayload))
	if err != nil {
		return nil, fmt.Errorf("calculate report impact before: %w", err)
	}
	afterReport, err := reportCalculator.Calculate(baseContext.WithOperatingPayload(&afterOperating.OperatingPayload))
	if err != nil {
		return nil, fmt.Errorf("calculate report impact after: %w", err)
	}

	return &AdjustmentImpactResult{
		ResolvedStageCode:       resolvedStageCode,
		CashBefore:              beforeReport.ReportPostTaxCash,
		CashAfter:               afterReport.ReportPostTaxCash,
		PreTaxProfitAfter:       afterReport.ReportPreTaxProfit,
		IncomeTaxAfter:          afterReport.ReportIncomeTax,
		NetProfitAfter:          afterReport.ReportNetProfit,
		TotalEquityAfter:        afterReport.ReportTotalEquity,
		WillBankrupt:            afterReport.ReportPostTaxCash < 0,
		CalculationBasisSavedAt: basisSavedAt,
		OperatingAfter:          afterOperating,
		ReportAfter:             afterReport,
		ReportBefore:            beforeReport,
		ReportManual:            reportManual,
	}, nil
}

func cloneAdjustmentOperatingPayload(source payload.OperatingPayload) (payload.OperatingPayload, error) {
	raw, err := json.Marshal(source)
	if err != nil {
		return payload.OperatingPayload{}, fmt.Errorf("marshal operating impact base: %w", err)
	}
	var cloned payload.OperatingPayload
	if err := json.Unmarshal(raw, &cloned); err != nil {
		return payload.OperatingPayload{}, fmt.Errorf("unmarshal operating impact base: %w", err)
	}
	return cloned.Normalize(), nil
}

func applyAdjustmentOperation(active []entity.GroupAdjustment, req AdjustmentImpactRequest) ([]entity.GroupAdjustment, error) {
	result := append([]entity.GroupAdjustment(nil), active...)
	switch req.Operation {
	case adjustmentOperationCurrent:
		return result, nil
	case adjustmentOperationCreate:
		if req.Candidate == nil {
			return nil, ErrAdminAdjustmentTypeInvalid
		}
		return append(result, *req.Candidate), nil
	case adjustmentOperationVoid:
		found := false
		filtered := make([]entity.GroupAdjustment, 0, len(result))
		for _, item := range result {
			if item.ID == req.VoidID {
				found = true
				continue
			}
			filtered = append(filtered, item)
		}
		if !found {
			return nil, ErrAdminAdjustmentNotEffective
		}
		return filtered, nil
	default:
		return nil, ErrAdminAdjustmentOperationInvalid
	}
}

func latestSavedAt(current *time.Time, candidate *time.Time) *time.Time {
	if candidate == nil {
		return current
	}
	if current == nil || candidate.After(*current) {
		copied := *candidate
		return &copied
	}
	return current
}

func reportBalancePassed(value payload.ReportComputedPayload) bool {
	return math.Abs(value.BalanceGap()) < 0.005
}
