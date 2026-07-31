package dto

type AdminGroupDataGetOperatingViewRequest struct {
	GroupID *int64 `form:"groupId" binding:"required"`
	YearNo  *int   `form:"yearNo" binding:"required"`
}

type AdminGroupDataGetReportViewRequest struct {
	GroupID *int64 `form:"groupId" binding:"required"`
	YearNo  *int   `form:"yearNo" binding:"required"`
}

type AdminGroupDataGetOperationContextRequest struct {
	GroupID *int64 `form:"groupId" binding:"required"`
}
