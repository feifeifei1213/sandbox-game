package dto

type AdminOrderGetControlConfigRequest struct {
	YearNo *int `form:"yearNo" binding:"required"`
}

type AdminOrderUpdateForecastControlRequest struct {
	Items      []AdminOrderUpdateForecastControlItem   `json:"items" binding:"required"`
	Narratives []AdminOrderUpdateForecastNarrativeItem `json:"narratives"`
}

type AdminOrderUpdateForecastControlItem struct {
	YearNo     int    `json:"yearNo" binding:"required"`
	MarketCode string `json:"marketCode" binding:"required"`
	OrderType  string `json:"orderType" binding:"required"`
	OrderCount int    `json:"orderCount"`
}

type AdminOrderUpdateForecastNarrativeItem struct {
	ForecastStageCode string `json:"forecastStageCode" binding:"required"`
	MarketCode        string `json:"marketCode" binding:"required"`
	Content           string `json:"content"`
}

type AdminOrderUpdateControlConfigRequest struct {
	YearNo *int                                `json:"yearNo" binding:"required"`
	Items  []AdminOrderUpdateControlConfigItem `json:"items" binding:"required"`
}

type AdminOrderUpdateMarketConfigRequest struct {
	YearNo  *int                               `json:"yearNo" binding:"required"`
	Markets []AdminOrderUpdateMarketConfigItem `json:"markets" binding:"required"`
}

type AdminOrderUpdateMarketConfigItem struct {
	MarketCode            string   `json:"marketCode" binding:"required"`
	Enabled               bool     `json:"enabled"`
	MarketInvestmentLimit *float64 `json:"marketInvestmentLimit"`
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

type AdminOrderConfirmPoolRequest struct {
	YearNo  *int   `json:"yearNo" binding:"required"`
	BatchID *int64 `json:"batchId" binding:"required"`
}

type AdminOrderGenerateSelectionSequenceRequest struct {
	YearNo *int `json:"yearNo" binding:"required"`
}

type AdminOrderGetOrderPoolRequest struct {
	YearNo     *int   `form:"yearNo" binding:"required"`
	MarketCode string `form:"marketCode"`
	OrderType  string `form:"orderType"`
}
