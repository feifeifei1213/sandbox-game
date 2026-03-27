package service

import (
	"reflect"
	"testing"

	"sandbox-game/internal/enum"
	"sandbox-game/internal/state"
)

func TestBuildAdminReadonlyOperatingPermissionAlwaysReadonly(t *testing.T) {
	permission := buildAdminReadonlyOperatingPermission(state.RuntimeState{
		StageStatus: enum.StageStatusQ3Open,
	})

	if !permission.CanView {
		t.Fatalf("expected admin operating view to be visible")
	}
	if permission.CanEdit {
		t.Fatalf("expected admin operating view to be readonly")
	}
	if permission.CanSubmit {
		t.Fatalf("expected admin operating view to disable submit")
	}
	if permission.CurrentStageCode != "Q3" {
		t.Fatalf("expected current stage code Q3, got %q", permission.CurrentStageCode)
	}

	expectedReadonlyScopes := []string{
		state.OperatingScopeYearStart,
		state.OperatingScopeQ1,
		state.OperatingScopeQ2,
		state.OperatingScopeQ3,
		state.OperatingScopeQ4,
		state.OperatingScopeYearEnd,
	}
	if !reflect.DeepEqual(permission.ReadonlyScopes, expectedReadonlyScopes) {
		t.Fatalf("expected readonly scopes %v, got %v", expectedReadonlyScopes, permission.ReadonlyScopes)
	}
	if len(permission.EditableScopes) != 0 {
		t.Fatalf("expected no editable scopes, got %v", permission.EditableScopes)
	}
}

func TestBuildAdminReadonlyReportPermissionAlwaysReadonly(t *testing.T) {
	permission := buildAdminReadonlyReportPermission()

	if !permission.CanView {
		t.Fatalf("expected admin report view to be visible")
	}
	if permission.CanEdit {
		t.Fatalf("expected admin report view to be readonly")
	}
	if permission.CanSubmit {
		t.Fatalf("expected admin report view to disable submit")
	}
}
