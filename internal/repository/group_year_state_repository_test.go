package repository

import (
	"testing"
	"time"

	"sandbox-game/internal/enum"
)

func TestBuildMissingFormalYearStatesCreatesLockedFormalYears(t *testing.T) {
	operateTime := time.Date(2026, 3, 27, 10, 0, 0, 0, time.UTC)
	items := buildMissingFormalYearStates([]int64{1}, 4, 5, "admin", operateTime, map[groupYearStateKey]struct{}{})
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	first := items[0]
	if first.GroupID != 1 || first.YearNo != 4 {
		t.Fatalf("unexpected first identity: group=%d year=%d", first.GroupID, first.YearNo)
	}
	if first.YearType != enum.YearTypeFormal {
		t.Fatalf("expected YearType FORMAL, got %s", first.YearType)
	}
	if first.YearStatus != enum.YearStatusLocked {
		t.Fatalf("expected YearStatus LOCKED, got %s", first.YearStatus)
	}
	if first.StageStatus != enum.StageStatusQ1Open {
		t.Fatalf("expected StageStatus Q1_OPEN, got %s", first.StageStatus)
	}
	if first.ReportStatus != enum.ReportStatusLocked {
		t.Fatalf("expected ReportStatus REPORT_LOCKED, got %s", first.ReportStatus)
	}
	if first.SummaryEffective {
		t.Fatalf("expected SummaryEffective false")
	}
	if first.LatestStageSubmitVersion != 0 || first.LatestReportSubmitVersion != 0 {
		t.Fatalf("expected submit versions to default to 0, got stage=%d report=%d", first.LatestStageSubmitVersion, first.LatestReportSubmitVersion)
	}
	if first.Creator != "admin" || first.Updater != "admin" {
		t.Fatalf("expected creator/updater admin, got creator=%s updater=%s", first.Creator, first.Updater)
	}
	if !first.CreateTime.Equal(operateTime) || !first.UpdateTime.Equal(operateTime) {
		t.Fatalf("expected audit timestamps to equal operateTime")
	}
}

func TestBuildMissingFormalYearStatesSkipsExistingPairs(t *testing.T) {
	items := buildMissingFormalYearStates(
		[]int64{1, 2},
		4,
		4,
		"admin",
		time.Now(),
		map[groupYearStateKey]struct{}{
			{GroupID: 1, YearNo: 4}: {},
		},
	)
	if len(items) != 1 {
		t.Fatalf("expected 1 item after skipping existing pair, got %d", len(items))
	}
	if items[0].GroupID != 2 || items[0].YearNo != 4 {
		t.Fatalf("expected only group 2 year 4 to remain, got group=%d year=%d", items[0].GroupID, items[0].YearNo)
	}
}

func TestNormalizeRollbackStageCodeAcceptsStatusAndReport(t *testing.T) {
	cases := map[string]string{
		"q2_open":                  "Q2",
		" q3 ":                     "Q3",
		enum.ReportStatusSubmitted: "YEAR_END",
	}

	for input, expected := range cases {
		if actual := normalizeRollbackStageCode(input); actual != expected {
			t.Fatalf("expected %q -> %q, got %q", input, expected, actual)
		}
	}
}
