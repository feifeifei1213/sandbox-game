package payload

import (
	"encoding/json"
	"testing"
)

func TestOperatingPayloadNormalizeKeepsLegacySupplyChainOrderRecordEmpty(t *testing.T) {
	var legacy OperatingPayload
	if err := json.Unmarshal([]byte(`{
		"beginning":{"taxAndPlanning":{},"marketBid":[]},
		"quarter":{"shortTermLoan":{},"receivableUpdate":{}},
		"yearEnd":{"longTermLoan":{},"assetAdjustment":{}},
		"extra":{"incomeAndPenalty":{}},
		"derived":{"values":{}}
	}`), &legacy); err != nil {
		t.Fatalf("unmarshal legacy operating payload: %v", err)
	}

	normalized := legacy.Normalize()
	if normalized.Quarter.SupplyChainOrderRecord == nil {
		t.Fatal("expected legacy payload to normalize supplyChainOrderRecord to an empty map")
	}
	if len(normalized.Quarter.SupplyChainOrderRecord) != 0 {
		t.Fatalf("expected legacy supplyChainOrderRecord to stay empty, got %#v", normalized.Quarter.SupplyChainOrderRecord)
	}
}

func TestOperatingPayloadNormalizePreservesSupplyChainOrderRecord(t *testing.T) {
	value := NewOperatingPayload()
	value.Quarter.SupplyChainOrderRecord = OperatingQuarterMap{
		"q1": {
			"basicProduct":       1,
			"standardProduct":    2,
			"precisionProduct":   0,
			"intelligentProduct": 0,
		},
	}

	normalized := value.Normalize()
	if got := normalized.Quarter.SupplyChainOrderRecord["q1"]["standardProduct"]; got != 2 {
		t.Fatalf("expected supply chain order record to be preserved, got %#v", got)
	}
}

func TestApplyMarketCultivationStateLeavesLegacyFallbackWhenNewFieldsAreEmpty(t *testing.T) {
	value := NewOperatingMarketCultivationPayload()

	applied := ApplyMarketCultivationState(value, MarketCultivationCarryState{}, false, true)
	if applied.StateApplied {
		t.Fatalf("expected empty new market cultivation structure to stay unapplied")
	}
	if amount, ok := MarketCultivationAnnualAmount(applied); ok || amount != 0 {
		t.Fatalf("expected empty new market cultivation to allow legacy fallback, got amount=%v ok=%v", amount, ok)
	}
}

func TestMarketCultivationExplicitZeroSuppressesLegacyFallback(t *testing.T) {
	value := NewOperatingMarketCultivationPayload()
	value.Regional.AnnualInvestment = 0

	applied := ApplyMarketCultivationState(value, MarketCultivationCarryState{}, false, true)
	if !applied.StateApplied {
		t.Fatal("expected explicit 0 input to mark new market cultivation structure as applied")
	}
	amount, ok := MarketCultivationAnnualAmount(applied)
	if !ok {
		t.Fatal("expected explicit 0 input to produce an authoritative annual amount")
	}
	if amount != 0 {
		t.Fatalf("expected explicit 0 input to produce 0 annual amount, got %v", amount)
	}
}

func TestApplyMarketCultivationStateKeepsUnlockedPreviousMarketFreeOfRepeatedCharge(t *testing.T) {
	value := NewOperatingMarketCultivationPayload()
	value.Regional.AnnualInvestment = 1

	applied := ApplyMarketCultivationState(value, MarketCultivationCarryState{
		RegionalCumulative: 1,
		RegionalUnlocked:   true,
	}, false, false)

	if !applied.StateApplied {
		t.Fatal("expected previous unlocked market to mark state as applied")
	}
	if !applied.Regional.Locked || !applied.Regional.LockedByPrevious || !applied.Regional.Unlocked {
		t.Fatalf("expected previous unlocked regional market to be locked and inherited, got %#v", applied.Regional)
	}
	if applied.Regional.EffectiveAnnualInvestment != 0 {
		t.Fatalf("expected inherited unlocked market to avoid repeated charge, got %v", applied.Regional.EffectiveAnnualInvestment)
	}
	amount, ok := MarketCultivationAnnualAmount(applied)
	if !ok || amount != 0 {
		t.Fatalf("expected inherited unlocked market annual amount 0, got amount=%v ok=%v", amount, ok)
	}
}
