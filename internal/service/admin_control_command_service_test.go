package service

import (
	"errors"
	"testing"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/model/payload"
	"sandbox-game/internal/state"
)

func TestBuildFinalYearExpansionPlanReturnsEmptyWhenNotExpanded(t *testing.T) {
	plan := buildFinalYearExpansionPlan(3, 3)
	if plan.InitializedFromYear != nil {
		t.Fatalf("expected InitializedFromYear to be nil, got %v", *plan.InitializedFromYear)
	}
	if plan.InitializedToYear != nil {
		t.Fatalf("expected InitializedToYear to be nil, got %v", *plan.InitializedToYear)
	}
	if plan.InitializedYearCount != 0 {
		t.Fatalf("expected InitializedYearCount to be 0, got %d", plan.InitializedYearCount)
	}
}

func TestBuildFinalYearExpansionPlanBuildsRangeWhenExpanded(t *testing.T) {
	plan := buildFinalYearExpansionPlan(3, 5)
	if plan.InitializedFromYear == nil || *plan.InitializedFromYear != 4 {
		t.Fatalf("expected InitializedFromYear to be 4, got %v", plan.InitializedFromYear)
	}
	if plan.InitializedToYear == nil || *plan.InitializedToYear != 5 {
		t.Fatalf("expected InitializedToYear to be 5, got %v", plan.InitializedToYear)
	}
	if plan.InitializedYearCount != 2 {
		t.Fatalf("expected InitializedYearCount to be 2, got %d", plan.InitializedYearCount)
	}
}

func TestValidateFinalYearChangeRejectsSmallerThanCurrentOpenYear(t *testing.T) {
	err := validateFinalYearChange(3, 2)
	if !errors.Is(err, ErrAdminControlFinalYearTooSmall) {
		t.Fatalf("expected ErrAdminControlFinalYearTooSmall, got %v", err)
	}
}

func TestValidateFinalYearChangeAllowsSameOrLargerYear(t *testing.T) {
	if err := validateFinalYearChange(3, 3); err != nil {
		t.Fatalf("expected nil error for equal year, got %v", err)
	}
	if err := validateFinalYearChange(3, 5); err != nil {
		t.Fatalf("expected nil error for larger year, got %v", err)
	}
}

func TestValidateOpenNextYearTargetRejectsWhenReachedFinalYear(t *testing.T) {
	err := validateOpenNextYearTarget(3, 3, 4)
	if !errors.Is(err, ErrAdminControlFinalYearReached) {
		t.Fatalf("expected ErrAdminControlFinalYearReached, got %v", err)
	}
}

func TestValidateOpenNextYearTargetRejectsTargetMismatch(t *testing.T) {
	err := validateOpenNextYearTarget(1, 4, 3)
	if !errors.Is(err, ErrAdminControlTargetYearMismatch) {
		t.Fatalf("expected ErrAdminControlTargetYearMismatch, got %v", err)
	}
}

func TestValidateOpenNextYearTargetAllowsDirectNextYear(t *testing.T) {
	if err := validateOpenNextYearTarget(1, 4, 2); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestEnsureOpenNextYearAllowedRejectsBlockedStatus(t *testing.T) {
	err := ensureOpenNextYearAllowed(openNextYearStatus{CanOpenNextYear: false, BlockedReason: openNextYearBlockedReasonUnfinishedReports})
	if !errors.Is(err, ErrAdminControlOpenNextYearBlocked) {
		t.Fatalf("expected ErrAdminControlOpenNextYearBlocked, got %v", err)
	}
	if err.Error() != openNextYearBlockedReasonUnfinishedReports {
		t.Fatalf("expected blocked reason to propagate, got %q", err.Error())
	}
}

func TestValidateInitialBaselineSubmissionRejectsNilPayload(t *testing.T) {
	err := validateInitialBaselineSubmission(false, nil)
	if !errors.Is(err, ErrAdminControlInitialBaselineInvalid) {
		t.Fatalf("expected ErrAdminControlInitialBaselineInvalid, got %v", err)
	}
}

func TestValidateInitialBaselineSubmissionRejectsDuplicateSubmit(t *testing.T) {
	err := validateInitialBaselineSubmission(true, &payload.BaselinePayload{})
	if !errors.Is(err, ErrAdminControlInitialBaselineSubmitted) {
		t.Fatalf("expected ErrAdminControlInitialBaselineSubmitted, got %v", err)
	}
}

func TestValidateInitialBaselineSubmissionAllowsFreshSubmit(t *testing.T) {
	if err := validateInitialBaselineSubmission(false, &payload.BaselinePayload{BaselineCash: 36}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestEnsureUnlockYearAllowedRejectsBlankReason(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusCompleted,
		StageStatus:      enum.StageStatusYearEndOpen,
		ReportStatus:     enum.ReportStatusSubmitted,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: true,
	}

	_, err := ensureUnlockYearAllowed(current, "   ", false)
	if !errors.Is(err, ErrAdminControlUnlockReasonRequired) {
		t.Fatalf("expected ErrAdminControlUnlockReasonRequired, got %v", err)
	}
}

func TestEnsureUnlockYearAllowedRejectsWhenNextYearAlreadyOpened(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusCompleted,
		StageStatus:      enum.StageStatusYearEndOpen,
		ReportStatus:     enum.ReportStatusSubmitted,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: true,
	}

	_, err := ensureUnlockYearAllowed(current, "reason text", true)
	if !errors.Is(err, ErrAdminControlUnlockNotAllowed) {
		t.Fatalf("expected ErrAdminControlUnlockNotAllowed, got %v", err)
	}
	if err.Error() != unlockYearBlockedReasonNextYearOpened {
		t.Fatalf("expected next year opened reason, got %q", err.Error())
	}
}

func TestEnsureUnlockYearAllowedRejectsAlreadyEditableYear(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusOperating,
		StageStatus:      enum.StageStatusQ2Open,
		ReportStatus:     enum.ReportStatusLocked,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: false,
	}

	_, err := ensureUnlockYearAllowed(current, "need reopen locked section", false)
	if !errors.Is(err, ErrAdminControlUnlockNotAllowed) {
		t.Fatalf("expected ErrAdminControlUnlockNotAllowed, got %v", err)
	}
	if err.Error() != unlockYearBlockedReasonAlreadyEditing {
		t.Fatalf("expected already editable reason, got %q", err.Error())
	}
}

func TestEnsureUnlockYearAllowedAllowsCompletedYear(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusCompleted,
		StageStatus:      enum.StageStatusYearEndOpen,
		ReportStatus:     enum.ReportStatusSubmitted,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: true,
	}

	next, err := ensureUnlockYearAllowed(current, "fix year end report", false)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if next.YearStatus != enum.YearStatusOperating {
		t.Fatalf("expected year status OPERATING, got %s", next.YearStatus)
	}
	if next.StageStatus != enum.StageStatusYearEndOpen {
		t.Fatalf("expected stage status YEAR_END_OPEN, got %s", next.StageStatus)
	}
	if next.ReportStatus != enum.ReportStatusLocked {
		t.Fatalf("expected report status REPORT_LOCKED, got %s", next.ReportStatus)
	}
	if next.SummaryEffective {
		t.Fatalf("expected summaryEffective to be false after unlock")
	}
}

func TestEnsureUnlockYearAllowedAllowsBankruptOperatingYear(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusOperating,
		StageStatus:      enum.StageStatusQ3Open,
		ReportStatus:     enum.ReportStatusLocked,
		BusinessStatus:   enum.BusinessStatusBankrupt,
		SummaryEffective: false,
	}

	next, err := ensureUnlockYearAllowed(current, "resume editing after bankrupt", false)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if next.StageStatus != enum.StageStatusQ3Open {
		t.Fatalf("expected stage status to remain Q3_OPEN, got %s", next.StageStatus)
	}
	if next.BusinessStatus != enum.BusinessStatusBankrupt {
		t.Fatalf("expected business status to remain BANKRUPT before recovery, got %s", next.BusinessStatus)
	}
}

func TestShouldRecoverFromBankruptMatchesBankruptYear(t *testing.T) {
	bankruptYearNo := 2
	group := &entity.Group{
		BusinessStatus: enum.BusinessStatusBankrupt,
		BankruptYearNo: &bankruptYearNo,
	}

	if !shouldRecoverFromBankrupt(group, 2) {
		t.Fatalf("expected recover decision to be true")
	}
	if shouldRecoverFromBankrupt(group, 1) {
		t.Fatalf("expected recover decision to be false for mismatched year")
	}
}
