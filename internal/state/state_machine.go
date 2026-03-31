package state

import (
	"errors"
	"fmt"

	"sandbox-game/internal/enum"
)

const (
	StageCodeQ1      = "Q1"
	StageCodeQ2      = "Q2"
	StageCodeQ3      = "Q3"
	StageCodeQ4      = "Q4"
	StageCodeYearEnd = "YEAR_END"
)

var (
	ErrInvalidRuntimeState = errors.New("invalid runtime state")
	ErrYearCannotOpen      = errors.New("year cannot open")
	ErrStageCannotSubmit   = errors.New("stage cannot submit")
	ErrReportCannotStart   = errors.New("report cannot start")
	ErrReportCannotSubmit  = errors.New("report cannot submit")
	ErrYearCannotUnlock    = errors.New("year cannot unlock")
	ErrInvalidStageCode    = errors.New("invalid stage code")
)

var stageStatusByCode = map[string]string{
	StageCodeQ1:      enum.StageStatusQ1Open,
	StageCodeQ2:      enum.StageStatusQ2Open,
	StageCodeQ3:      enum.StageStatusQ3Open,
	StageCodeQ4:      enum.StageStatusQ4Open,
	StageCodeYearEnd: enum.StageStatusYearEndOpen,
}

var stageCodeByStatus = map[string]string{
	enum.StageStatusQ1Open:      StageCodeQ1,
	enum.StageStatusQ2Open:      StageCodeQ2,
	enum.StageStatusQ3Open:      StageCodeQ3,
	enum.StageStatusQ4Open:      StageCodeQ4,
	enum.StageStatusYearEndOpen: StageCodeYearEnd,
}

type RuntimeState struct {
	YearNo                    int
	YearType                  string
	YearStatus                string
	StageStatus               string
	ReportStatus              string
	BusinessStatus            string
	SummaryEffective          bool
	LatestStageSubmitVersion  int
	LatestReportSubmitVersion int
}

type StateMachine struct {
	guard *TransitionGuard
}

func NewStateMachine() *StateMachine {
	return &StateMachine{guard: NewTransitionGuard()}
}

func (m *StateMachine) Validate(current RuntimeState) error {
	if !enum.IsValidYearType(current.YearType) {
		return fmt.Errorf("%w: unsupported yearType=%s", ErrInvalidRuntimeState, current.YearType)
	}
	if !enum.IsValidYearStatus(current.YearStatus) {
		return fmt.Errorf("%w: unsupported yearStatus=%s", ErrInvalidRuntimeState, current.YearStatus)
	}
	if !enum.IsValidStageStatus(current.StageStatus) {
		return fmt.Errorf("%w: unsupported stageStatus=%s", ErrInvalidRuntimeState, current.StageStatus)
	}
	if !enum.IsValidReportStatus(current.ReportStatus) {
		return fmt.Errorf("%w: unsupported reportStatus=%s", ErrInvalidRuntimeState, current.ReportStatus)
	}
	if !enum.IsValidBusinessStatus(current.BusinessStatus) {
		return fmt.Errorf("%w: unsupported businessStatus=%s", ErrInvalidRuntimeState, current.BusinessStatus)
	}
	if current.YearNo < 0 {
		return fmt.Errorf("%w: yearNo cannot be negative", ErrInvalidRuntimeState)
	}
	if current.LatestStageSubmitVersion < 0 || current.LatestReportSubmitVersion < 0 {
		return fmt.Errorf("%w: submit version cannot be negative", ErrInvalidRuntimeState)
	}

	switch current.YearStatus {
	case enum.YearStatusLocked, enum.YearStatusOperating:
		if current.ReportStatus != enum.ReportStatusLocked {
			return fmt.Errorf("%w: yearStatus=%s requires reportStatus=%s", ErrInvalidRuntimeState, current.YearStatus, enum.ReportStatusLocked)
		}
	case enum.YearStatusReportPending, enum.YearStatusReporting:
		if current.ReportStatus != enum.ReportStatusOpen {
			return fmt.Errorf("%w: yearStatus=%s requires reportStatus=%s", ErrInvalidRuntimeState, current.YearStatus, enum.ReportStatusOpen)
		}
	case enum.YearStatusCompleted:
		if current.ReportStatus != enum.ReportStatusSubmitted {
			return fmt.Errorf("%w: yearStatus=%s requires reportStatus=%s", ErrInvalidRuntimeState, current.YearStatus, enum.ReportStatusSubmitted)
		}
	}

	if current.SummaryEffective {
		if current.YearStatus != enum.YearStatusCompleted {
			return fmt.Errorf("%w: summary effective requires completed year", ErrInvalidRuntimeState)
		}
		if current.YearNo == 0 {
			return fmt.Errorf("%w: demo year cannot be summary effective", ErrInvalidRuntimeState)
		}
	}

	return nil
}

func (m *StateMachine) OpenYear(current RuntimeState) (RuntimeState, error) {
	if err := m.Validate(current); err != nil {
		return current, err
	}
	if !m.guard.CanOpenYear(current) {
		return current, fmt.Errorf("%w: yearStatus=%s, businessStatus=%s", ErrYearCannotOpen, current.YearStatus, current.BusinessStatus)
	}

	next := current
	next.YearStatus = enum.YearStatusOperating
	next.StageStatus = enum.StageStatusQ1Open
	next.ReportStatus = enum.ReportStatusLocked
	next.SummaryEffective = false

	return next, nil
}

func (m *StateMachine) SubmitStage(current RuntimeState, stageCode string) (RuntimeState, error) {
	if err := m.Validate(current); err != nil {
		return current, err
	}
	if !m.guard.CanSubmitStage(current, stageCode) {
		return current, fmt.Errorf("%w: currentStage=%s, stageCode=%s", ErrStageCannotSubmit, current.StageStatus, stageCode)
	}

	next := current
	next.LatestStageSubmitVersion++

	switch stageCode {
	case StageCodeQ1:
		next.StageStatus = enum.StageStatusQ2Open
	case StageCodeQ2:
		next.StageStatus = enum.StageStatusQ3Open
	case StageCodeQ3:
		next.StageStatus = enum.StageStatusQ4Open
	case StageCodeQ4:
		next.StageStatus = enum.StageStatusYearEndOpen
	case StageCodeYearEnd:
		next.YearStatus = enum.YearStatusReportPending
		next.StageStatus = enum.StageStatusYearEndOpen
		next.ReportStatus = enum.ReportStatusOpen
	default:
		return current, fmt.Errorf("%w: stageCode=%s", ErrInvalidStageCode, stageCode)
	}

	return next, nil
}

func (m *StateMachine) StartReport(current RuntimeState) (RuntimeState, error) {
	if err := m.Validate(current); err != nil {
		return current, err
	}
	if current.BusinessStatus == enum.BusinessStatusBankrupt {
		return current, fmt.Errorf("%w: bankrupt group cannot start report", ErrReportCannotStart)
	}

	switch current.YearStatus {
	case enum.YearStatusReportPending:
		next := current
		next.YearStatus = enum.YearStatusReporting
		return next, nil
	case enum.YearStatusReporting:
		return current, nil
	default:
		return current, fmt.Errorf("%w: yearStatus=%s", ErrReportCannotStart, current.YearStatus)
	}
}

func (m *StateMachine) SubmitReport(current RuntimeState) (RuntimeState, error) {
	if err := m.Validate(current); err != nil {
		return current, err
	}
	if !m.guard.CanSubmitReport(current) {
		return current, fmt.Errorf("%w: yearStatus=%s, reportStatus=%s", ErrReportCannotSubmit, current.YearStatus, current.ReportStatus)
	}

	next := current
	next.YearStatus = enum.YearStatusCompleted
	next.ReportStatus = enum.ReportStatusSubmitted
	next.SummaryEffective = current.YearNo > 0
	next.LatestReportSubmitVersion++

	return next, nil
}

func (m *StateMachine) UnlockOperatingYear(current RuntimeState, targetStageCode string, nextYearAlreadyOpened bool) (RuntimeState, error) {
	if err := m.Validate(current); err != nil {
		return current, err
	}
	if !m.guard.CanUnlockYear(current, nextYearAlreadyOpened) {
		return current, fmt.Errorf("%w: yearStatus=%s, nextYearAlreadyOpened=%t", ErrYearCannotUnlock, current.YearStatus, nextYearAlreadyOpened)
	}

	targetStageStatus, ok := StageStatusByCode(targetStageCode)
	if !ok {
		return current, fmt.Errorf("%w: stageCode=%s", ErrInvalidStageCode, targetStageCode)
	}

	next := current
	next.YearStatus = enum.YearStatusOperating
	next.StageStatus = targetStageStatus
	next.ReportStatus = enum.ReportStatusLocked
	next.SummaryEffective = false

	return next, nil
}

func (m *StateMachine) UnlockReportYear(current RuntimeState, nextYearAlreadyOpened bool) (RuntimeState, error) {
	if err := m.Validate(current); err != nil {
		return current, err
	}
	if !m.guard.CanUnlockYear(current, nextYearAlreadyOpened) {
		return current, fmt.Errorf("%w: yearStatus=%s, nextYearAlreadyOpened=%t", ErrYearCannotUnlock, current.YearStatus, nextYearAlreadyOpened)
	}

	next := current
	next.YearStatus = enum.YearStatusReportPending
	next.StageStatus = enum.StageStatusYearEndOpen
	next.ReportStatus = enum.ReportStatusOpen
	next.SummaryEffective = false

	return next, nil
}

func CurrentStageCode(stageStatus string) string {
	return stageCodeByStatus[stageStatus]
}

func IsValidStageCode(stageCode string) bool {
	_, ok := stageStatusByCode[stageCode]
	return ok
}

func StageStatusByCode(stageCode string) (string, bool) {
	stageStatus, ok := stageStatusByCode[stageCode]
	return stageStatus, ok
}
