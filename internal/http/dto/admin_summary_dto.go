package dto

type AdminSummaryGetYearSummaryRequest struct {
	YearNo *int `form:"yearNo" binding:"required"`
}
