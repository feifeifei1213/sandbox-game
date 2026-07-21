package enum

const (
	SnapshotScopeGroup  = "GROUP"
	SnapshotScopeGlobal = "GLOBAL"
)

const (
	SnapshotTypeAuto   = "AUTO"
	SnapshotTypeManual = "MANUAL"
	SnapshotTypeSafety = "SAFETY"
)

const (
	SnapshotTriggerStageSubmitted       = "STAGE_SUBMITTED"
	SnapshotTriggerReportSubmitted      = "REPORT_SUBMITTED"
	SnapshotTriggerRollbackStageRetry   = "ROLLBACK_STAGE_RESUBMITTED"
	SnapshotTriggerRollbackReportRetry  = "ROLLBACK_REPORT_RESUBMITTED"
	SnapshotTriggerOrderPoolConfirmed   = "ORDER_POOL_CONFIRMED"
	SnapshotTriggerSegmentCompleted     = "SEGMENT_COMPLETED"
	SnapshotTriggerOpenNextYear         = "OPEN_NEXT_YEAR"
	SnapshotTriggerBeforeRollback       = "BEFORE_ROLLBACK"
	SnapshotTriggerAdjustmentBankruptcy = "ADJUSTMENT_BANKRUPTCY"
)

const (
	RollbackTypeUnlockRetry          = "UNLOCK_RETRY"
	RollbackTypeGroupSnapshotRestore = "GROUP_SNAPSHOT_RESTORE"
)

func IsValidSnapshotScope(value string) bool {
	switch value {
	case SnapshotScopeGroup, SnapshotScopeGlobal:
		return true
	default:
		return false
	}
}

func IsValidSnapshotType(value string) bool {
	switch value {
	case SnapshotTypeAuto, SnapshotTypeManual, SnapshotTypeSafety:
		return true
	default:
		return false
	}
}

func IsValidRollbackType(value string) bool {
	switch value {
	case RollbackTypeUnlockRetry, RollbackTypeGroupSnapshotRestore:
		return true
	default:
		return false
	}
}
