package service

import (
	"testing"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
)

func TestBuildAdminYearTabs(t *testing.T) {
	tabs := buildAdminYearTabs(1, 3)
	if len(tabs) != 4 {
		t.Fatalf("expected 4 tabs, got %d", len(tabs))
	}
	if tabs[0].TabStatus != yearTabStatusEnterable || !tabs[0].CanEnter {
		t.Fatalf("expected year 0 to be enterable")
	}
	if tabs[1].TabStatus != yearTabStatusEnterable || !tabs[1].IsCurrentOpenYear {
		t.Fatalf("expected year 1 to be current enterable tab")
	}
	if tabs[2].TabStatus != yearTabStatusLocked || tabs[2].CanEnter {
		t.Fatalf("expected year 2 to be locked")
	}
	if tabs[3].Label != "3年" || !tabs[3].IsFormalYear {
		t.Fatalf("expected year 3 label/formal-year metadata to be correct")
	}
}

func TestBuildGroupYearTabsNormalGroup(t *testing.T) {
	stateMap := map[int]entity.GroupYearState{
		0: {YearNo: 0, YearStatus: enum.YearStatusCompleted},
		1: {YearNo: 1, YearStatus: enum.YearStatusCompleted},
	}

	tabs := buildGroupYearTabs(2, 4, enum.BusinessStatusNormal, nil, stateMap)
	if tabs[0].TabStatus != yearTabStatusCompleted {
		t.Fatalf("expected year 0 completed, got %s", tabs[0].TabStatus)
	}
	if tabs[1].TabStatus != yearTabStatusCompleted {
		t.Fatalf("expected year 1 completed, got %s", tabs[1].TabStatus)
	}
	if tabs[2].TabStatus != yearTabStatusEnterable || !tabs[2].CanEnter {
		t.Fatalf("expected year 2 enterable")
	}
	if tabs[3].TabStatus != yearTabStatusLocked || tabs[4].TabStatus != yearTabStatusLocked {
		t.Fatalf("expected future years locked")
	}
}

func TestBuildGroupYearTabsBankruptGroup(t *testing.T) {
	bankruptYearNo := 2
	stateMap := map[int]entity.GroupYearState{
		1: {YearNo: 1, YearStatus: enum.YearStatusCompleted},
	}

	tabs := buildGroupYearTabs(3, 4, enum.BusinessStatusBankrupt, &bankruptYearNo, stateMap)
	if tabs[1].TabStatus != yearTabStatusCompleted {
		t.Fatalf("expected historical completed year to stay completed")
	}
	if tabs[2].TabStatus != yearTabStatusBankruptReadOnly || !tabs[2].CanEnter {
		t.Fatalf("expected bankrupt year to be readonly and enterable")
	}
	if tabs[3].TabStatus != yearTabStatusLocked || tabs[4].TabStatus != yearTabStatusLocked {
		t.Fatalf("expected years after bankruptcy to remain locked")
	}
}
