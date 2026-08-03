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

func TestOperatingProjectProgressDefaultsUseFactorySlotQuarterStructure(t *testing.T) {
	value := NewOperatingProjectProgressPayload()
	if len(value.Items) != 32 {
		t.Fatalf("expected 32 project progress records, got %d", len(value.Items))
	}

	first := value.Items[0]
	if first["projectName"] != "生产厂房 A-槽位1-第一季度" || first["factoryKey"] != "factoryA" || first["slotNo"] != 1 || first["quarterKey"] != "q1" {
		t.Fatalf("unexpected first project progress item: %#v", first)
	}

	last := value.Items[len(value.Items)-1]
	if last["projectName"] != "生产厂房 C-槽位1-第四季度" || last["factoryKey"] != "factoryC" || last["slotNo"] != 1 || last["quarterKey"] != "q4" {
		t.Fatalf("unexpected last project progress item: %#v", last)
	}
}

func TestOperatingProjectProgressNormalizeMapsLegacyTwelveItemsToFirstSlot(t *testing.T) {
	value := NewOperatingProjectProgressPayload()
	value.Items = []map[string]any{
		{"projectName": "生产厂房 A-第一季度", "lineType": "人工", "progress": 1},
		{"projectName": "生产厂房 A-第二季度", "lineType": "半自动", "progress": 2},
		{"projectName": "生产厂房 A-第三季度", "lineType": "", "progress": ""},
		{"projectName": "生产厂房 A-第四季度", "lineType": "", "progress": ""},
		{"projectName": "生产厂房 B-第一季度", "lineType": "自动", "progress": 3},
		{"projectName": "生产厂房 B-第二季度", "lineType": "", "progress": ""},
		{"projectName": "生产厂房 B-第三季度", "lineType": "", "progress": ""},
		{"projectName": "生产厂房 B-第四季度", "lineType": "", "progress": ""},
		{"projectName": "生产厂房 C-第一季度", "lineType": "智能", "progress": 4},
		{"projectName": "生产厂房 C-第二季度", "lineType": "", "progress": ""},
		{"projectName": "生产厂房 C-第三季度", "lineType": "", "progress": ""},
		{"projectName": "生产厂房 C-第四季度", "lineType": "", "progress": ""},
	}

	normalized := normalizeProjectProgressPayload(value)
	if len(normalized.Items) != 32 {
		t.Fatalf("expected 32 normalized project progress records, got %d", len(normalized.Items))
	}
	if got := normalized.Items[0]["lineType"]; got != "人工" {
		t.Fatalf("expected legacy A Q1 to map to A slot 1 Q1, got %#v", got)
	}
	if got := normalized.Items[4]["lineType"]; got != "" {
		t.Fatalf("expected A slot 2 Q1 to stay empty, got %#v", got)
	}
	if got := normalized.Items[16]["lineType"]; got != "自动" {
		t.Fatalf("expected legacy B Q1 to map to B slot 1 Q1, got %#v", got)
	}
	if got := normalized.Items[28]["lineType"]; got != "智能" {
		t.Fatalf("expected legacy C Q1 to map to C slot 1 Q1, got %#v", got)
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
