package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
	calcctx "sandbox-game/internal/rules/context"
	reportrules "sandbox-game/internal/rules/report"
)

var directorScoreFieldTokens = [][]byte{
	[]byte(`"reportBestMarketDirectorBaseScore"`),
	[]byte(`"reportBestMarketDirectorScore"`),
	[]byte(`"reportBestTechnologyDirectorScore"`),
	[]byte(`"reportBestSalesDirectorScore"`),
	[]byte(`"reportBestCfoBaseScore"`),
	[]byte(`"reportBestCfoScore"`),
	[]byte(`"reportBestCeoScore"`),
}

func reportComputedHasDirectorScores(raw []byte) bool {
	if len(raw) == 0 {
		return false
	}
	for _, token := range directorScoreFieldTokens {
		if !bytes.Contains(raw, token) {
			return false
		}
	}
	return true
}

func overlayDirectorScores(target *payload.ReportComputedPayload, source payload.ReportComputedPayload) {
	target.ReportBestMarketDirectorBaseScore = source.ReportBestMarketDirectorBaseScore
	target.ReportBestMarketDirectorScore = source.ReportBestMarketDirectorScore
	target.ReportBestTechnologyDirectorScore = source.ReportBestTechnologyDirectorScore
	target.ReportBestSalesDirectorScore = source.ReportBestSalesDirectorScore
	target.ReportBestCfoBaseScore = source.ReportBestCfoBaseScore
	target.ReportBestCfoScore = source.ReportBestCfoScore
	target.ReportBestCeoScore = source.ReportBestCeoScore
}

func loadEffectiveReportWithDirectorScores(
	ctx context.Context,
	group entity.Group,
	gameConfig entity.GameConfig,
	yearNo int,
	groupYearRepo *repository.GroupYearStateRepository,
	operatingRepo *repository.OperatingRepository,
	initialBaselineRepo *repository.InitialBaselineRepository,
	reportRepo *repository.ReportRepository,
	playerNoticeService *PlayerNoticeService,
	calculator *reportrules.Calculator,
) (*payload.ReportComputedPayload, error) {
	report, err := reportRepo.FindEffectiveByGroupIDAndYear(ctx, group.ID, yearNo)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, nil
	case err != nil:
		return nil, err
	}

	var computed payload.ReportComputedPayload
	if len(report.ReportComputedPayload) > 0 {
		if unmarshalErr := json.Unmarshal(report.ReportComputedPayload, &computed); unmarshalErr != nil {
			return nil, unmarshalErr
		}
	}

	rebuilt, rebuildErr := rebuildReportDirectorScores(
		ctx,
		group,
		gameConfig,
		yearNo,
		groupYearRepo,
		operatingRepo,
		initialBaselineRepo,
		reportRepo,
		playerNoticeService,
		calculator,
	)
	if rebuildErr == nil {
		overlayDirectorScores(&computed, rebuilt)
	}

	return &computed, nil
}

func rebuildReportDirectorScores(
	ctx context.Context,
	group entity.Group,
	gameConfig entity.GameConfig,
	yearNo int,
	groupYearRepo *repository.GroupYearStateRepository,
	operatingRepo *repository.OperatingRepository,
	initialBaselineRepo *repository.InitialBaselineRepository,
	reportRepo *repository.ReportRepository,
	playerNoticeService *PlayerNoticeService,
	calculator *reportrules.Calculator,
) (payload.ReportComputedPayload, error) {
	cache := make(map[int]payload.ReportComputedPayload)
	return rebuildReportDirectorScoresRecursive(
		ctx,
		group,
		gameConfig,
		yearNo,
		groupYearRepo,
		operatingRepo,
		initialBaselineRepo,
		reportRepo,
		playerNoticeService,
		calculator,
		cache,
	)
}

func rebuildReportDirectorScoresRecursive(
	ctx context.Context,
	group entity.Group,
	gameConfig entity.GameConfig,
	yearNo int,
	groupYearRepo *repository.GroupYearStateRepository,
	operatingRepo *repository.OperatingRepository,
	initialBaselineRepo *repository.InitialBaselineRepository,
	reportRepo *repository.ReportRepository,
	playerNoticeService *PlayerNoticeService,
	calculator *reportrules.Calculator,
	cache map[int]payload.ReportComputedPayload,
) (payload.ReportComputedPayload, error) {
	if result, ok := cache[yearNo]; ok {
		return result, nil
	}

	yearState, err := groupYearRepo.GetByGroupIDAndYear(ctx, group.ID, yearNo)
	if err != nil {
		return payload.ReportComputedPayload{}, fmt.Errorf("load year state for director score rebuild: %w", err)
	}

	operatingPayload := payload.NewOperatingPayload()
	draft, draftErr := operatingRepo.FindDraft(ctx, group.ID, yearNo)
	switch {
	case draftErr == nil:
		if len(draft.OperatingPayload) > 0 {
			if unmarshalErr := json.Unmarshal(draft.OperatingPayload, &operatingPayload); unmarshalErr != nil {
				return payload.ReportComputedPayload{}, fmt.Errorf("unmarshal operating payload for director score rebuild: %w", unmarshalErr)
			}
		}
	case errors.Is(draftErr, gorm.ErrRecordNotFound):
	default:
		return payload.ReportComputedPayload{}, fmt.Errorf("load operating draft for director score rebuild: %w", draftErr)
	}

	operatingPayload = operatingPayload.Normalize().WithoutDerivedValues()
	operatingPayload, err = playerNoticeService.OverlayAdjustments(ctx, group.ID, yearNo, operatingPayload)
	if err != nil {
		return payload.ReportComputedPayload{}, fmt.Errorf("overlay adjustments for director score rebuild: %w", err)
	}

	calcContext := calcctx.NewCalculationContext(group, *yearState, gameConfig).
		WithOperatingPayload(&operatingPayload)

	manualPayload := payload.ReportManualPayload{}
	report, reportErr := reportRepo.FindByGroupIDAndYear(ctx, group.ID, yearNo)
	switch {
	case reportErr == nil:
		if len(report.ReportManualPayload) > 0 {
			if unmarshalErr := json.Unmarshal(report.ReportManualPayload, &manualPayload); unmarshalErr != nil {
				return payload.ReportComputedPayload{}, fmt.Errorf("unmarshal report manual payload for director score rebuild: %w", unmarshalErr)
			}
		}
	case errors.Is(reportErr, gorm.ErrRecordNotFound):
	default:
		return payload.ReportComputedPayload{}, fmt.Errorf("load report payload for director score rebuild: %w", reportErr)
	}

	calcContext = calcContext.WithReportManualPayload(&manualPayload)

	if yearNo == 0 {
		baseline, baselineErr := initialBaselineRepo.FindByGroupID(ctx, group.ID)
		switch {
		case baselineErr == nil:
			if len(baseline.BaselinePayload) > 0 {
				var baselinePayload payload.BaselinePayload
				if unmarshalErr := json.Unmarshal(baseline.BaselinePayload, &baselinePayload); unmarshalErr != nil {
					return payload.ReportComputedPayload{}, fmt.Errorf("unmarshal baseline for director score rebuild: %w", unmarshalErr)
				}
				calcContext = calcContext.WithInitialBaseline(&baselinePayload)
			}
		case errors.Is(baselineErr, gorm.ErrRecordNotFound):
		default:
			return payload.ReportComputedPayload{}, fmt.Errorf("load baseline for director score rebuild: %w", baselineErr)
		}
	} else {
		previous, rebuildErr := rebuildReportDirectorScoresRecursive(
			ctx,
			group,
			gameConfig,
			yearNo-1,
			groupYearRepo,
			operatingRepo,
			initialBaselineRepo,
			reportRepo,
			playerNoticeService,
			calculator,
			cache,
		)
		if rebuildErr != nil {
			return payload.ReportComputedPayload{}, rebuildErr
		}
		calcContext = calcContext.WithPreviousReport(&previous)
	}

	computedPayload, err := calculator.Calculate(calcContext)
	if err != nil {
		return payload.ReportComputedPayload{}, fmt.Errorf("calculate report payload for director score rebuild: %w", err)
	}

	cache[yearNo] = computedPayload
	return computedPayload, nil
}
