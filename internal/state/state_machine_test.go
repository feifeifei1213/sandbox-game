package state

import (
	"errors"
	"testing"

	"sandbox-game/internal/enum"
)

func TestValidateRejectsInvalidStateCombination(t *testing.T) {
	t.Parallel()

	machine := NewStateMachine()
	err := machine.Validate(RuntimeState{
		YearNo:         1,
		YearType:       enum.YearTypeFormal,
		YearStatus:     enum.YearStatusOperating,
		StageStatus:    enum.StageStatusQ1Open,
		ReportStatus:   enum.ReportStatusOpen,
		BusinessStatus: enum.BusinessStatusNormal,
	})
	if err == nil {
		t.Fatal("expected invalid combination error, got nil")
	}
	if !errors.Is(err, ErrInvalidRuntimeState) {
		t.Fatalf("expected ErrInvalidRuntimeState, got %v", err)
	}
}

func TestSubmitStageTransitionsToNextStage(t *testing.T) {
	t.Parallel()

	machine := NewStateMachine()
	next, err := machine.SubmitStage(RuntimeState{
		YearNo:                    1,
		YearType:                  enum.YearTypeFormal,
		YearStatus:                enum.YearStatusOperating,
		StageStatus:               enum.StageStatusQ1Open,
		ReportStatus:              enum.ReportStatusLocked,
		BusinessStatus:            enum.BusinessStatusNormal,
		LatestStageSubmitVersion:  2,
		LatestReportSubmitVersion: 0,
	}, StageCodeQ1)
	if err != nil {
		t.Fatalf("submit q1 failed: %v", err)
	}

	if next.StageStatus != enum.StageStatusQ2Open {
		t.Fatalf("expected stageStatus=%s, got %s", enum.StageStatusQ2Open, next.StageStatus)
	}
	if next.LatestStageSubmitVersion != 3 {
		t.Fatalf("expected latestStageSubmitVersion=3, got %d", next.LatestStageSubmitVersion)
	}
}

func TestSubmitYearEndTransitionsToReportPending(t *testing.T) {
	t.Parallel()

	machine := NewStateMachine()
	next, err := machine.SubmitStage(RuntimeState{
		YearNo:         1,
		YearType:       enum.YearTypeFormal,
		YearStatus:     enum.YearStatusOperating,
		StageStatus:    enum.StageStatusYearEndOpen,
		ReportStatus:   enum.ReportStatusLocked,
		BusinessStatus: enum.BusinessStatusNormal,
	}, StageCodeYearEnd)
	if err != nil {
		t.Fatalf("submit year end failed: %v", err)
	}

	if next.YearStatus != enum.YearStatusReportPending {
		t.Fatalf("expected yearStatus=%s, got %s", enum.YearStatusReportPending, next.YearStatus)
	}
	if next.ReportStatus != enum.ReportStatusOpen {
		t.Fatalf("expected reportStatus=%s, got %s", enum.ReportStatusOpen, next.ReportStatus)
	}
}

func TestSubmitReportMarksFormalYearSummaryEffective(t *testing.T) {
	t.Parallel()

	machine := NewStateMachine()
	next, err := machine.SubmitReport(RuntimeState{
		YearNo:         1,
		YearType:       enum.YearTypeFormal,
		YearStatus:     enum.YearStatusReportPending,
		StageStatus:    enum.StageStatusYearEndOpen,
		ReportStatus:   enum.ReportStatusOpen,
		BusinessStatus: enum.BusinessStatusNormal,
	})
	if err != nil {
		t.Fatalf("submit report failed: %v", err)
	}

	if next.YearStatus != enum.YearStatusCompleted {
		t.Fatalf("expected completed year, got %s", next.YearStatus)
	}
	if !next.SummaryEffective {
		t.Fatal("expected formal year to be summary effective")
	}
}

func TestSubmitReportKeepsDemoYearOutOfSummary(t *testing.T) {
	t.Parallel()

	machine := NewStateMachine()
	next, err := machine.SubmitReport(RuntimeState{
		YearNo:         0,
		YearType:       enum.YearTypeDemo,
		YearStatus:     enum.YearStatusReporting,
		StageStatus:    enum.StageStatusYearEndOpen,
		ReportStatus:   enum.ReportStatusOpen,
		BusinessStatus: enum.BusinessStatusNormal,
	})
	if err != nil {
		t.Fatalf("submit report failed: %v", err)
	}

	if next.SummaryEffective {
		t.Fatal("expected demo year not to enter summary")
	}
}

func TestBuildOperatingPermissionForQ1(t *testing.T) {
	t.Parallel()

	guard := NewTransitionGuard()
	permission := guard.BuildOperatingPermission(RuntimeState{
		YearNo:         1,
		YearType:       enum.YearTypeFormal,
		YearStatus:     enum.YearStatusOperating,
		StageStatus:    enum.StageStatusQ1Open,
		ReportStatus:   enum.ReportStatusLocked,
		BusinessStatus: enum.BusinessStatusNormal,
	})

	if !permission.CanEdit || !permission.CanSubmit {
		t.Fatal("expected q1 operating page to be editable and submittable")
	}
	if len(permission.EditableScopes) != 2 {
		t.Fatalf("expected 2 editable scopes, got %d", len(permission.EditableScopes))
	}
	if permission.EditableScopes[0] != OperatingScopeYearStart || permission.EditableScopes[1] != OperatingScopeQ1 {
		t.Fatalf("unexpected editable scopes: %#v", permission.EditableScopes)
	}
}

func TestUnlockCompletedYearFallsBackToOperating(t *testing.T) {
	t.Parallel()

	machine := NewStateMachine()
	next, err := machine.UnlockYear(RuntimeState{
		YearNo:           1,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusCompleted,
		StageStatus:      enum.StageStatusYearEndOpen,
		ReportStatus:     enum.ReportStatusSubmitted,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: true,
	}, false)
	if err != nil {
		t.Fatalf("unlock year failed: %v", err)
	}

	if next.YearStatus != enum.YearStatusOperating {
		t.Fatalf("expected yearStatus=%s, got %s", enum.YearStatusOperating, next.YearStatus)
	}
	if next.ReportStatus != enum.ReportStatusLocked {
		t.Fatalf("expected reportStatus=%s, got %s", enum.ReportStatusLocked, next.ReportStatus)
	}
	if next.SummaryEffective {
		t.Fatal("expected summary to be withdrawn after unlock")
	}
}

func TestUnlockRejectedAfterNextYearOpened(t *testing.T) {
	t.Parallel()

	machine := NewStateMachine()
	_, err := machine.UnlockYear(RuntimeState{
		YearNo:         1,
		YearType:       enum.YearTypeFormal,
		YearStatus:     enum.YearStatusCompleted,
		StageStatus:    enum.StageStatusYearEndOpen,
		ReportStatus:   enum.ReportStatusSubmitted,
		BusinessStatus: enum.BusinessStatusNormal,
	}, true)
	if err == nil {
		t.Fatal("expected unlock to be rejected")
	}
	if !errors.Is(err, ErrYearCannotUnlock) {
		t.Fatalf("expected ErrYearCannotUnlock, got %v", err)
	}
}
