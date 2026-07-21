package dto

type PlayerAdjustmentSyncRequest struct {
	YearNo        *int   `form:"yearNo" binding:"required"`
	KnownRevision *int64 `form:"knownRevision" binding:"required"`
}
