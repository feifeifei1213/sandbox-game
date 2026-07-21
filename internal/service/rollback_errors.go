package service

import "errors"

var (
	ErrRollbackSnapshotNotFound         = errors.New("rollback snapshot not found")
	ErrRollbackSnapshotScopeUnsupported = errors.New("rollback snapshot scope unsupported")
	ErrRollbackTargetInvalid            = errors.New("rollback target invalid")
	ErrRollbackReasonRequired           = errors.New("rollback reason required")
	ErrRollbackConfirmRequired          = errors.New("rollback confirm required")
	ErrRollbackStateChanged             = errors.New("rollback state changed")
	ErrRollbackGroupWriteLocked         = errors.New("rollback group write locked")
	ErrRollbackSnapshotNotRestorable    = errors.New("rollback snapshot not restorable")
)
