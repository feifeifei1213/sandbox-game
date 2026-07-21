package service

import (
	"errors"
	"testing"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/model/entity"
	"sandbox-game/internal/state"
)

func TestResolveAdjustmentStageFromCurrentRuntimeState(t *testing.T) {
	tests := []struct {
		name         string
		yearStatus   string
		stageStatus  string
		reportStatus string
		want         string
	}{
		{"q1", enum.YearStatusOperating, enum.StageStatusQ1Open, enum.ReportStatusLocked, state.StageCodeQ1},
		{"q2", enum.YearStatusOperating, enum.StageStatusQ2Open, enum.ReportStatusLocked, state.StageCodeQ2},
		{"q3", enum.YearStatusOperating, enum.StageStatusQ3Open, enum.ReportStatusLocked, state.StageCodeQ3},
		{"q4", enum.YearStatusOperating, enum.StageStatusQ4Open, enum.ReportStatusLocked, state.StageCodeQ4},
		{"year end", enum.YearStatusOperating, enum.StageStatusYearEndOpen, enum.ReportStatusLocked, state.StageCodeYearEnd},
		{"report pending", enum.YearStatusReportPending, enum.StageStatusYearEndOpen, enum.ReportStatusOpen, state.StageCodeYearEnd},
		{"reporting", enum.YearStatusReporting, enum.StageStatusYearEndOpen, enum.ReportStatusOpen, state.StageCodeYearEnd},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveAdjustmentStage(entity.GroupYearState{
				YearStatus:   tt.yearStatus,
				StageStatus:  tt.stageStatus,
				ReportStatus: tt.reportStatus,
			})
			if err != nil {
				t.Fatalf("resolve adjustment stage: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}

func TestResolveAdjustmentStageRejectsClosedStates(t *testing.T) {
	tests := []struct {
		name      string
		state     entity.GroupYearState
		wantError error
	}{
		{"locked year", entity.GroupYearState{YearStatus: enum.YearStatusLocked, ReportStatus: enum.ReportStatusLocked}, ErrAdminAdjustmentYearNotOpen},
		{"completed year", entity.GroupYearState{YearStatus: enum.YearStatusCompleted, ReportStatus: enum.ReportStatusSubmitted}, ErrAdminAdjustmentStageLocked},
		{"submitted report", entity.GroupYearState{YearStatus: enum.YearStatusReporting, ReportStatus: enum.ReportStatusSubmitted}, ErrAdminAdjustmentStageLocked},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := resolveAdjustmentStage(tt.state)
			if !errors.Is(err, tt.wantError) {
				t.Fatalf("expected %v, got %v", tt.wantError, err)
			}
		})
	}
}
