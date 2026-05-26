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

func TestValidateInitializeGameEditionRequiresKnownEdition(t *testing.T) {
	_, err := validateInitializeGameEdition("")
	if !errors.Is(err, ErrAdminControlEditionRequired) {
		t.Fatalf("expected ErrAdminControlEditionRequired, got %v", err)
	}

	_, err = validateInitializeGameEdition("UNKNOWN")
	if !errors.Is(err, ErrAdminControlEditionInvalid) {
		t.Fatalf("expected ErrAdminControlEditionInvalid, got %v", err)
	}

	edition, err := validateInitializeGameEdition("vip_service_v1")
	if err != nil {
		t.Fatalf("expected known edition to pass, got %v", err)
	}
	if edition.EditionCode != GameEditionVIPServiceV1 || edition.RuleVersion != FormulaVersionCommonV1 {
		t.Fatalf("unexpected edition: %#v", edition)
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

	_, err := ensureUnlockYearAllowed(current, "   ", unlockTargetTypeOperating, state.StageCodeYearEnd, false)
	if !errors.Is(err, ErrAdminControlUnlockReasonRequired) {
		t.Fatalf("expected ErrAdminControlUnlockReasonRequired, got %v", err)
	}
}

func TestEnsureUnlockYearAllowedRejectsMissingTargetType(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusCompleted,
		StageStatus:      enum.StageStatusYearEndOpen,
		ReportStatus:     enum.ReportStatusSubmitted,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: true,
	}

	_, err := ensureUnlockYearAllowed(current, "reason text", "", "", false)
	if !errors.Is(err, ErrAdminControlUnlockTargetTypeRequired) {
		t.Fatalf("expected ErrAdminControlUnlockTargetTypeRequired, got %v", err)
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

	_, err := ensureUnlockYearAllowed(current, "reason text", unlockTargetTypeOperating, state.StageCodeYearEnd, true)
	if !errors.Is(err, ErrAdminControlUnlockNotAllowed) {
		t.Fatalf("expected ErrAdminControlUnlockNotAllowed, got %v", err)
	}
	if err.Error() != unlockYearBlockedReasonNextYearOpened {
		t.Fatalf("expected next year opened reason, got %q", err.Error())
	}
}

func TestEnsureUnlockYearAllowedRejectsOperatingUnlockWithoutStage(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusCompleted,
		StageStatus:      enum.StageStatusYearEndOpen,
		ReportStatus:     enum.ReportStatusSubmitted,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: true,
	}

	_, err := ensureUnlockYearAllowed(current, "reason text", unlockTargetTypeOperating, "", false)
	if !errors.Is(err, ErrAdminControlUnlockStageRequired) {
		t.Fatalf("expected ErrAdminControlUnlockStageRequired, got %v", err)
	}
}

func TestEnsureUnlockYearAllowedRejectsOperatingTargetAlreadyEditable(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusOperating,
		StageStatus:      enum.StageStatusQ2Open,
		ReportStatus:     enum.ReportStatusLocked,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: false,
	}

	_, err := ensureUnlockYearAllowed(current, "need reopen q2", unlockTargetTypeOperating, state.StageCodeQ2, false)
	if !errors.Is(err, ErrAdminControlUnlockNotAllowed) {
		t.Fatalf("expected ErrAdminControlUnlockNotAllowed, got %v", err)
	}
	if err.Error() != unlockYearBlockedReasonOperatingEditable {
		t.Fatalf("expected operating already editable reason, got %q", err.Error())
	}
}

func TestEnsureUnlockYearAllowedRejectsReportTargetAlreadyEditable(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusReportPending,
		StageStatus:      enum.StageStatusYearEndOpen,
		ReportStatus:     enum.ReportStatusOpen,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: false,
	}

	_, err := ensureUnlockYearAllowed(current, "need reopen report", unlockTargetTypeReport, "", false)
	if !errors.Is(err, ErrAdminControlUnlockNotAllowed) {
		t.Fatalf("expected ErrAdminControlUnlockNotAllowed, got %v", err)
	}
	if err.Error() != unlockYearBlockedReasonReportEditable {
		t.Fatalf("expected report already editable reason, got %q", err.Error())
	}
}

func TestEnsureUnlockYearAllowedRejectsUnsubmittedOperatingTarget(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusOperating,
		StageStatus:      enum.StageStatusQ3Open,
		ReportStatus:     enum.ReportStatusLocked,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: false,
	}

	_, err := ensureUnlockYearAllowed(current, "need reopen year end", unlockTargetTypeOperating, state.StageCodeYearEnd, false)
	if !errors.Is(err, ErrAdminControlUnlockNotAllowed) {
		t.Fatalf("expected ErrAdminControlUnlockNotAllowed, got %v", err)
	}
	if err.Error() != unlockYearBlockedReasonTargetNotSubmitted {
		t.Fatalf("expected target not submitted reason, got %q", err.Error())
	}
}

func TestEnsureUnlockYearAllowedAllowsCompletedYearToQ2(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusCompleted,
		StageStatus:      enum.StageStatusYearEndOpen,
		ReportStatus:     enum.ReportStatusSubmitted,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: true,
	}

	next, err := ensureUnlockYearAllowed(current, "fix q2 result", unlockTargetTypeOperating, state.StageCodeQ2, false)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if next.YearStatus != enum.YearStatusOperating {
		t.Fatalf("expected year status OPERATING, got %s", next.YearStatus)
	}
	if next.StageStatus != enum.StageStatusQ2Open {
		t.Fatalf("expected stage status Q2_OPEN, got %s", next.StageStatus)
	}
	if next.ReportStatus != enum.ReportStatusLocked {
		t.Fatalf("expected report status REPORT_LOCKED, got %s", next.ReportStatus)
	}
	if next.SummaryEffective {
		t.Fatalf("expected summaryEffective to be false after unlock")
	}
}

func TestEnsureUnlockYearAllowedAllowsCompletedYearReportUnlock(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusCompleted,
		StageStatus:      enum.StageStatusYearEndOpen,
		ReportStatus:     enum.ReportStatusSubmitted,
		BusinessStatus:   enum.BusinessStatusNormal,
		SummaryEffective: true,
	}

	next, err := ensureUnlockYearAllowed(current, "fix report", unlockTargetTypeReport, "", false)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if next.YearStatus != enum.YearStatusReportPending {
		t.Fatalf("expected year status REPORT_PENDING, got %s", next.YearStatus)
	}
	if next.ReportStatus != enum.ReportStatusOpen {
		t.Fatalf("expected report status REPORT_OPEN, got %s", next.ReportStatus)
	}
}

func TestEnsureUnlockYearAllowedAllowsBankruptOperatingYearRollback(t *testing.T) {
	current := state.RuntimeState{
		YearNo:           2,
		YearType:         enum.YearTypeFormal,
		YearStatus:       enum.YearStatusOperating,
		StageStatus:      enum.StageStatusQ3Open,
		ReportStatus:     enum.ReportStatusLocked,
		BusinessStatus:   enum.BusinessStatusBankrupt,
		SummaryEffective: false,
	}

	next, err := ensureUnlockYearAllowed(current, "resume editing after bankrupt", unlockTargetTypeOperating, state.StageCodeQ2, false)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if next.StageStatus != enum.StageStatusQ2Open {
		t.Fatalf("expected stage status Q2_OPEN, got %s", next.StageStatus)
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
