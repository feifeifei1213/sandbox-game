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
