package dto

type AdminOrderGetControlConfigRequest struct {
	YearNo *int `form:"yearNo" binding:"required"`
}

type AdminOrderUpdateControlConfigRequest struct {
	YearNo *int                                `json:"yearNo" binding:"required"`
	Items  []AdminOrderUpdateControlConfigItem `json:"items" binding:"required"`
}

type AdminOrderUpdateControlConfigItem struct {
	MarketCode        string `json:"marketCode" binding:"required"`
	OrderType         string `json:"orderType" binding:"required"`
	OrderCount        int    `json:"orderCount"`
	ReleaseSequenceNo int    `json:"releaseSequenceNo" binding:"required"`
}

type AdminOrderGeneratePoolRequest struct {
	YearNo        *int   `json:"yearNo" binding:"required"`
	Overwrite     bool   `json:"overwrite"`
	SourceBatchID *int64 `json:"sourceBatchId"`
}

type AdminOrderGetOrderPoolRequest struct {
	YearNo     *int   `form:"yearNo" binding:"required"`
	MarketCode string `form:"marketCode" binding:"required"`
	OrderType  string `form:"orderType" binding:"required"`
}
