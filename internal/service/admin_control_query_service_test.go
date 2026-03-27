package service

import (
	"testing"
	"time"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
)

func TestEvaluateOpenNextYearStatusBlocksWhenReachedFinalYear(t *testing.T) {
	result := evaluateOpenNextYearStatus(&entity.GameConfig{
		CurrentOpenYear:          3,
		FinalYear:                3,
		InitialBaselineSubmitted: true,
	}, nil, nil)

	if result.CanOpenNextYear {
		t.Fatalf("expected canOpenNextYear to be false")
	}
	if result.BlockedReason != openNextYearBlockedReasonFinalYearReached {
		t.Fatalf("expected blocked reason %q, got %q", openNextYearBlockedReasonFinalYearReached, result.BlockedReason)
	}
	if result.NextOpenableYear != 0 {
		t.Fatalf("expected nextOpenableYear to be 0, got %d", result.NextOpenableYear)
	}
}

func TestEvaluateOpenNextYearStatusBlocksWhenBaselineNotSubmitted(t *testing.T) {
	result := evaluateOpenNextYearStatus(&entity.GameConfig{
		CurrentOpenYear:          0,
		FinalYear:                3,
		InitialBaselineSubmitted: false,
	}, nil, nil)

	if result.CanOpenNextYear {
		t.Fatalf("expected canOpenNextYear to be false")
	}
	if result.BlockedReason != openNextYearBlockedReasonBaselineUnsubmitted {
		t.Fatalf("expected blocked reason %q, got %q", openNextYearBlockedReasonBaselineUnsubmitted, result.BlockedReason)
	}
	if result.NextOpenableYear != 1 {
		t.Fatalf("expected nextOpenableYear to be 1, got %d", result.NextOpenableYear)
	}
}

func TestEvaluateOpenNextYearStatusRequiresAllNonBankruptGroupsCompleted(t *testing.T) {
	groups := []entity.Group{
		{ID: 1, BusinessStatus: enum.BusinessStatusNormal},
		{ID: 2, BusinessStatus: enum.BusinessStatusNormal},
	}
	states := []entity.GroupYearState{
		{GroupID: 1, YearStatus: enum.YearStatusCompleted},
		{GroupID: 2, YearStatus: enum.YearStatusReporting},
	}

	result := evaluateOpenNextYearStatus(&entity.GameConfig{
		CurrentOpenYear:          1,
		FinalYear:                3,
		InitialBaselineSubmitted: true,
	}, groups, states)

	if result.CanOpenNextYear {
		t.Fatalf("expected canOpenNextYear to be false")
	}
	if result.BlockedReason != openNextYearBlockedReasonUnfinishedReports {
		t.Fatalf("expected blocked reason %q, got %q", openNextYearBlockedReasonUnfinishedReports, result.BlockedReason)
	}
}

func TestEvaluateOpenNextYearStatusIgnoresBankruptGroups(t *testing.T) {
	groups := []entity.Group{
		{ID: 1, BusinessStatus: enum.BusinessStatusNormal},
		{ID: 2, BusinessStatus: enum.BusinessStatusBankrupt},
	}
	states := []entity.GroupYearState{
		{GroupID: 1, YearStatus: enum.YearStatusCompleted},
		{GroupID: 2, YearStatus: enum.YearStatusLocked},
	}

	result := evaluateOpenNextYearStatus(&entity.GameConfig{
		CurrentOpenYear:          1,
		FinalYear:                3,
		InitialBaselineSubmitted: true,
	}, groups, states)

	if !result.CanOpenNextYear {
		t.Fatalf("expected canOpenNextYear to be true")
	}
	if result.BlockedReason != "" {
		t.Fatalf("expected empty blocked reason, got %q", result.BlockedReason)
	}
	if result.NextOpenableYear != 2 {
		t.Fatalf("expected nextOpenableYear to be 2, got %d", result.NextOpenableYear)
	}
}

func TestBuildInitialBaselineViewResultSetsEditableBeforeSubmit(t *testing.T) {
	view := buildInitialBaselineViewResult(false, payload.BaselinePayload{BaselineCash: 36}, 0, nil, nil)
	if !view.Editable {
		t.Fatalf("expected editable to be true before submit")
	}
	if view.Submitted {
		t.Fatalf("expected submitted to be false before submit")
	}
	if view.AppliedGroupCount != 0 {
		t.Fatalf("expected appliedGroupCount to be 0, got %d", view.AppliedGroupCount)
	}
	if view.BaselinePayload.BaselineCash != 36 {
		t.Fatalf("expected baseline cash to be preserved")
	}
}

func TestBuildInitialBaselineViewResultFormatsSubmittedMeta(t *testing.T) {
	now := time.Date(2026, 3, 26, 10, 30, 0, 0, time.FixedZone("CST", 8*3600))
	submitter := "admin"
	view := buildInitialBaselineViewResult(true, payload.BaselinePayload{}, 10, &submitter, &now)
	if view.Editable {
		t.Fatalf("expected editable to be false after submit")
	}
	if !view.Submitted {
		t.Fatalf("expected submitted to be true")
	}
	if view.SubmittedAt == nil || *view.SubmittedAt != now.Format(time.RFC3339) {
		t.Fatalf("expected submittedAt to be formatted RFC3339")
	}
	if view.SubmitterName == nil || *view.SubmitterName != submitter {
		t.Fatalf("expected submitterName to be preserved")
	}
}
