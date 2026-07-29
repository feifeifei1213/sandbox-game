package enum

const (
	MarketCodeLocal    = "LOCAL"
	MarketCodeRegional = "REGIONAL"
	MarketCodeNational = "NATIONAL"
	MarketCodeGlobal   = "GLOBAL"
)

const (
	OrderTypeAgencyInspection = "AGENCY_INSPECTION"
	OrderTypeTwoCabinVIP      = "TWO_CABIN_VIP"
	OrderTypeBusinessVIP      = "BUSINESS_VIP"
	OrderTypeMemberCustom     = "MEMBER_CUSTOM"
)

const (
	OrderImportStatusSuccess = "SUCCESS"
	OrderImportStatusFailed  = "FAILED"
)

const (
	OrderConfigStatusDraft  = "DRAFT"
	OrderConfigStatusLocked = "LOCKED"
)

const (
	OrderGenerationBatchStatusPreview   = "PREVIEW"
	OrderGenerationBatchStatusConfirmed = "CONFIRMED"
	OrderGenerationBatchStatusVoid      = "VOID"
)

const (
	OrderPoolStatusAvailable         = "AVAILABLE"
	OrderPoolStatusSelected          = "SELECTED"
	OrderPoolStatusUnselectedExpired = "UNSELECTED_EXPIRED"
	OrderPoolStatusVoid              = "VOID"
)

const (
	OrderSegmentStatusMarketDisabled    = "MARKET_DISABLED"
	OrderSegmentStatusNoOrderConfig     = "NO_ORDER_CONFIG"
	OrderSegmentStatusWaitingInvestment = "WAITING_INVESTMENT"
	OrderSegmentStatusBidOpen           = "BID_OPEN"
	OrderSegmentStatusBidClosed         = "BID_CLOSED"
	OrderSegmentStatusSequenceReady     = "SEQUENCE_READY"
	OrderSegmentStatusWaitingRelease    = "WAITING_RELEASE"
	OrderSegmentStatusSelecting         = "SELECTING"
	OrderSegmentStatusRoundReady        = "ROUND_READY"
	OrderSegmentStatusCompleted         = "COMPLETED"
	OrderSegmentStatusSkipped           = "SKIPPED"
)

const (
	OrderBidStatusSubmitted = "SUBMITTED"
	OrderBidStatusLocked    = "LOCKED"
)

const (
	OrderSelectionStatusIneligible   = "INELIGIBLE"
	OrderSelectionStatusWaiting      = "WAITING"
	OrderSelectionStatusCurrent      = "CURRENT"
	OrderSelectionStatusSelected     = "SELECTED"
	OrderSelectionStatusPassed       = "PASSED"
	OrderSelectionStatusAdminSkipped = "ADMIN_SKIPPED"
	OrderSelectionStatusBankrupt     = "INELIGIBLE_BANKRUPT"
)

const (
	OrderCompletionReasonAllRoundsCompleted     = "ALL_ROUNDS_COMPLETED"
	OrderCompletionReasonPoolExhausted          = "ORDER_POOL_EXHAUSTED"
	OrderCompletionReasonNoEligibleParticipants = "NO_ELIGIBLE_PARTICIPANTS"
)

const (
	OrderDeliveryStatusSelected   = "SELECTED"
	OrderDeliveryStatusDelivered  = "DELIVERED"
	OrderDeliveryStatusUnfinished = "UNFINISHED"
)

const (
	OrderDeliveryRevisionTypeDelivered   = "DELIVERED"
	OrderDeliveryRevisionTypeInvalidated = "INVALIDATED"
)

func IsValidMarketCode(value string) bool {
	switch value {
	case MarketCodeLocal, MarketCodeRegional, MarketCodeNational, MarketCodeGlobal:
		return true
	default:
		return false
	}
}

func IsValidOrderType(value string) bool {
	switch value {
	case OrderTypeAgencyInspection, OrderTypeTwoCabinVIP, OrderTypeBusinessVIP, OrderTypeMemberCustom:
		return true
	default:
		return false
	}
}
