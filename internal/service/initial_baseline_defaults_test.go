package service

import (
	"testing"

	"sandbox-game/internal/model/payload"
)

func TestApplyInitialBaselineDefaultsToOperatingPayloadUsesBaselineWorkInConstructionWhenMissing(t *testing.T) {
	operatingPayload := payload.NewOperatingPayload()
	baseline := &payload.BaselinePayload{BaselineWorkInConstruction: 8}

	result := applyInitialBaselineDefaultsToOperatingPayload(operatingPayload, baseline)

	value, ok := result.YearEnd.AssetAdjustment["workInConstruction"].(float64)
	if !ok {
		t.Fatalf("expected workInConstruction default to be written, got %#v", result.YearEnd.AssetAdjustment["workInConstruction"])
	}
	if value != 8 {
		t.Fatalf("expected workInConstruction default to be 8, got %v", value)
	}
}

func TestApplyInitialBaselineDefaultsToOperatingPayloadPreservesExplicitZero(t *testing.T) {
	operatingPayload := payload.NewOperatingPayload()
	operatingPayload.YearEnd.AssetAdjustment["workInConstruction"] = 0.0
	baseline := &payload.BaselinePayload{BaselineWorkInConstruction: 8}

	result := applyInitialBaselineDefaultsToOperatingPayload(operatingPayload, baseline)

	value, ok := result.YearEnd.AssetAdjustment["workInConstruction"].(float64)
	if !ok {
		t.Fatalf("expected explicit workInConstruction to remain numeric, got %#v", result.YearEnd.AssetAdjustment["workInConstruction"])
	}
	if value != 0 {
		t.Fatalf("expected explicit zero to be preserved, got %v", value)
	}
}
