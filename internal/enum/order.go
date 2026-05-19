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
	OrderPoolStatusAvailable = "AVAILABLE"
	OrderPoolStatusSelected  = "SELECTED"
	OrderPoolStatusVoid      = "VOID"
)

const (
	OrderSegmentStatusWaitingRelease = "WAITING_RELEASE"
	OrderSegmentStatusSelecting      = "SELECTING"
	OrderSegmentStatusCompleted      = "COMPLETED"
	OrderSegmentStatusSkipped        = "SKIPPED"
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
