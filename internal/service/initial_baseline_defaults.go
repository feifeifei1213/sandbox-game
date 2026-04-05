package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/repository"
)

var initialBaselineWorkInConstructionAliases = []string{
	"workInConstruction",
	"unfinishedLineValue",
	"unfinishedProductionLineValue",
	"o52",
}

func loadInitialBaselinePayload(
	ctx context.Context,
	initialBaselineRepo *repository.InitialBaselineRepository,
	groupID int64,
) (*payload.BaselinePayload, error) {
	baseline, baselineErr := initialBaselineRepo.FindByGroupID(ctx, groupID)
	switch {
	case baselineErr == nil:
		if len(baseline.BaselinePayload) == 0 {
			return nil, nil
		}

		var baselinePayload payload.BaselinePayload
		if err := json.Unmarshal(baseline.BaselinePayload, &baselinePayload); err != nil {
			return nil, fmt.Errorf("unmarshal initial baseline: %w", err)
		}
		return &baselinePayload, nil
	case errors.Is(baselineErr, gorm.ErrRecordNotFound):
		return nil, nil
	default:
		return nil, fmt.Errorf("load initial baseline: %w", baselineErr)
	}
}

func applyInitialBaselineDefaultsToOperatingPayload(
	operatingPayload payload.OperatingPayload,
	baseline *payload.BaselinePayload,
) payload.OperatingPayload {
	normalized := operatingPayload.Normalize().WithoutDerivedValues()
	if baseline == nil {
		return normalized
	}

	if normalized.YearEnd.AssetAdjustment == nil {
		normalized.YearEnd.AssetAdjustment = map[string]any{}
	}
	if !hasExplicitMetricValue(normalized.YearEnd.AssetAdjustment, initialBaselineWorkInConstructionAliases...) {
		normalized.YearEnd.AssetAdjustment["workInConstruction"] = baseline.BaselineWorkInConstruction
	}

	return normalized
}

func hasExplicitMetricValue(source map[string]any, aliases ...string) bool {
	if len(source) == 0 || len(aliases) == 0 {
		return false
	}

	normalizedAliases := make([]string, 0, len(aliases))
	for _, alias := range aliases {
		normalized := normalizeBaselineMetricKey(alias)
		if normalized != "" {
			normalizedAliases = append(normalizedAliases, normalized)
		}
	}
	if len(normalizedAliases) == 0 {
		return false
	}

	for key, value := range source {
		if !baselineMetricKeyMatches(key, normalizedAliases) {
			continue
		}
		switch typed := value.(type) {
		case nil:
			continue
		case string:
			if strings.TrimSpace(typed) == "" {
				continue
			}
			return true
		default:
			return true
		}
	}
	return false
}

func baselineMetricKeyMatches(source string, normalizedAliases []string) bool {
	normalizedSource := normalizeBaselineMetricKey(source)
	if normalizedSource == "" {
		return false
	}
	for _, alias := range normalizedAliases {
		if normalizedSource == alias || strings.Contains(normalizedSource, alias) || strings.Contains(alias, normalizedSource) {
			return true
		}
	}
	return false
}

func normalizeBaselineMetricKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(value))
	for _, ch := range value {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}
