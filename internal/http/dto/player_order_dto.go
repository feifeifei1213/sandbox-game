package dto

type PlayerOrderGetYearViewRequest struct {
	YearNo *int `form:"yearNo" binding:"required"`
}

type PlayerOrderSubmitMarketInvestmentRequest struct {
	YearNo      *int                                   `json:"yearNo" binding:"required"`
	Investments []PlayerOrderMarketInvestmentInputItem `json:"investments" binding:"required"`
}

type PlayerOrderMarketInvestmentInputItem struct {
	MarketCode       string  `json:"marketCode" binding:"required"`
	OrderType        string  `json:"orderType" binding:"required"`
	MarketInvestment float64 `json:"marketInvestment"`
}

type PlayerOrderSelectOrderRequest struct {
	YearNo     *int   `json:"yearNo" binding:"required"`
	MarketCode string `json:"marketCode" binding:"required"`
	OrderType  string `json:"orderType" binding:"required"`
	OrderID    *int64 `json:"orderId" binding:"required"`
}

type PlayerOrderPassSegmentRequest struct {
	YearNo     *int   `json:"yearNo" binding:"required"`
	MarketCode string `json:"marketCode" binding:"required"`
	OrderType  string `json:"orderType" binding:"required"`
}

type PlayerOrderDeliverOrdersRequest struct {
	YearNo    *int    `json:"yearNo" binding:"required"`
	StageCode string  `json:"stageCode" binding:"required"`
	OrderIDs  []int64 `json:"orderIds" binding:"required"`
}

type AdminOrderOpenMarketBiddingRequest struct {
	YearNo     *int   `json:"yearNo" binding:"required"`
	MarketCode string `json:"marketCode" binding:"required"`
}

type AdminOrderCloseMarketBiddingRequest struct {
	YearNo     *int   `json:"yearNo" binding:"required"`
	MarketCode string `json:"marketCode" binding:"required"`
}

type AdminOrderGetMarketSelectionStatusRequest struct {
	YearNo     *int   `form:"yearNo" binding:"required"`
	MarketCode string `form:"marketCode" binding:"required"`
}

type AdminOrderReleaseNextSegmentRequest struct {
	YearNo *int `json:"yearNo" binding:"required"`
}

type AdminOrderOpenNextRoundRequest struct {
	YearNo     *int   `json:"yearNo" binding:"required"`
	MarketCode string `json:"marketCode" binding:"required"`
	OrderType  string `json:"orderType" binding:"required"`
}

type AdminOrderSkipCurrentGroupRequest struct {
	YearNo     *int   `json:"yearNo" binding:"required"`
	MarketCode string `json:"marketCode" binding:"required"`
	OrderType  string `json:"orderType" binding:"required"`
	GroupID    *int64 `json:"groupId" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
}
