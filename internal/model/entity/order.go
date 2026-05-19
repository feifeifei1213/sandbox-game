package entity

import "time"

type OrderImportBatch struct {
	ID               int64     `gorm:"column:id;primaryKey"`
	OriginalFileName string    `gorm:"column:original_file_name"`
	FileSize         int64     `gorm:"column:file_size"`
	ParseStatus      string    `gorm:"column:parse_status"`
	ParsedOrderCount int       `gorm:"column:parsed_order_count"`
	ParseError       []byte    `gorm:"column:parse_error_json"`
	ParsedPayload    []byte    `gorm:"column:parsed_payload_json"`
	UploaderID       int64     `gorm:"column:uploader_id"`
	UploaderName     string    `gorm:"column:uploader_name"`
	UploadedAt       time.Time `gorm:"column:uploaded_at"`
	BaseEntity
}

func (OrderImportBatch) TableName() string {
	return "sg_order_import_batch"
}

type OrderGenerationConfig struct {
	ID                int64  `gorm:"column:id;primaryKey"`
	YearNo            int    `gorm:"column:year_no"`
	MarketCode        string `gorm:"column:market_code"`
	OrderType         string `gorm:"column:order_type"`
	OrderCount        int    `gorm:"column:order_count"`
	ReleaseSequenceNo int    `gorm:"column:release_sequence_no"`
	SourceBatchID     *int64 `gorm:"column:source_batch_id"`
	ConfigStatus      string `gorm:"column:config_status"`
	BaseEntity
}

func (OrderGenerationConfig) TableName() string {
	return "sg_order_generation_config"
}

type OrderPool struct {
	ID              int64   `gorm:"column:id;primaryKey"`
	YearNo          int     `gorm:"column:year_no"`
	MarketCode      string  `gorm:"column:market_code"`
	OrderType       string  `gorm:"column:order_type"`
	OrderAmount     float64 `gorm:"column:order_amount"`
	OrderQuantity   float64 `gorm:"column:order_quantity"`
	UnitPrice       float64 `gorm:"column:unit_price"`
	AccountTerm     int     `gorm:"column:account_term"`
	PoolStatus      string  `gorm:"column:pool_status"`
	SelectedGroupID *int64  `gorm:"column:selected_group_id"`
	SourceBatchID   *int64  `gorm:"column:source_batch_id"`
	SourceSheetName string  `gorm:"column:source_sheet_name"`
	SourceCell      string  `gorm:"column:source_cell"`
	SourceRowIndex  *int    `gorm:"column:source_row_index"`
	BaseEntity
}

func (OrderPool) TableName() string {
	return "sg_order_pool"
}

type MarketBiddingState struct {
	ID                int64      `gorm:"column:id;primaryKey"`
	YearNo            int        `gorm:"column:year_no"`
	MarketCode        string     `gorm:"column:market_code"`
	OrderType         string     `gorm:"column:order_type"`
	SegmentCode       string     `gorm:"column:segment_code"`
	ReleaseSequenceNo int        `gorm:"column:release_sequence_no"`
	SegmentStatus     string     `gorm:"column:segment_status"`
	LeaderGroupID     *int64     `gorm:"column:leader_group_id"`
	LeaderRule        []byte     `gorm:"column:leader_rule_json"`
	RandomSeed        *string    `gorm:"column:random_seed"`
	CurrentGroupID    *int64     `gorm:"column:current_group_id"`
	OpenedAt          *time.Time `gorm:"column:opened_at"`
	ClosedAt          *time.Time `gorm:"column:closed_at"`
	ReleasedAt        *time.Time `gorm:"column:released_at"`
	CompletedAt       *time.Time `gorm:"column:completed_at"`
	BaseEntity
}

func (MarketBiddingState) TableName() string {
	return "sg_market_bidding_state"
}

type GroupMarketBid struct {
	ID               int64     `gorm:"column:id;primaryKey"`
	GroupID          int64     `gorm:"column:group_id"`
	YearNo           int       `gorm:"column:year_no"`
	MarketCode       string    `gorm:"column:market_code"`
	MarketInvestment float64   `gorm:"column:market_investment"`
	SubmittedAt      time.Time `gorm:"column:submitted_at"`
	BaseEntity
}

func (GroupMarketBid) TableName() string {
	return "sg_group_market_bid"
}

type MarketSelectionOrder struct {
	ID                        int64      `gorm:"column:id;primaryKey"`
	YearNo                    int        `gorm:"column:year_no"`
	MarketCode                string     `gorm:"column:market_code"`
	OrderType                 string     `gorm:"column:order_type"`
	SequenceNo                int        `gorm:"column:sequence_no"`
	GroupID                   int64      `gorm:"column:group_id"`
	MarketInvestment          float64    `gorm:"column:market_investment"`
	PreviousMarketOrderAmount float64    `gorm:"column:previous_market_order_amount"`
	IsMarketLeader            bool       `gorm:"column:is_market_leader"`
	RankBasis                 []byte     `gorm:"column:rank_basis_json"`
	SelectionStatus           string     `gorm:"column:selection_status"`
	SelectedOrderID           *int64     `gorm:"column:selected_order_id"`
	SelectedAt                *time.Time `gorm:"column:selected_at"`
	SkippedByAdminID          *int64     `gorm:"column:skipped_by_admin_id"`
	SkippedReason             *string    `gorm:"column:skipped_reason"`
	SkippedAt                 *time.Time `gorm:"column:skipped_at"`
	BaseEntity
}

func (MarketSelectionOrder) TableName() string {
	return "sg_market_selection_order"
}

type GroupOrderSelection struct {
	ID                 int64      `gorm:"column:id;primaryKey"`
	GroupID            int64      `gorm:"column:group_id"`
	YearNo             int        `gorm:"column:year_no"`
	MarketCode         string     `gorm:"column:market_code"`
	OrderType          string     `gorm:"column:order_type"`
	OrderID            int64      `gorm:"column:order_id"`
	SelectionStatus    string     `gorm:"column:selection_status"`
	DeliveryStatus     string     `gorm:"column:delivery_status"`
	DeliveredStageCode *string    `gorm:"column:delivered_stage_code"`
	DeliveredAt        *time.Time `gorm:"column:delivered_at"`
	SelectedAt         time.Time  `gorm:"column:selected_at"`
	BaseEntity
}

func (GroupOrderSelection) TableName() string {
	return "sg_group_order_selection"
}
