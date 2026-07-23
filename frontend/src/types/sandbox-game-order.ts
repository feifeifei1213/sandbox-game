export type OrderMarketCode = string
export type OrderTypeCode = string
export type OrderPoolStatus = 'AVAILABLE' | 'SELECTED' | 'UNSELECTED_EXPIRED' | 'VOID'
export type OrderForecastStageCode = 'YEAR_1_3' | 'YEAR_4_5' | 'YEAR_6_8' | string
export type OrderSegmentStatus =
  | 'MARKET_DISABLED'
  | 'NO_ORDER_CONFIG'
  | 'WAITING_INVESTMENT'
  | 'BID_OPEN'
  | 'BID_CLOSED'
  | 'SEQUENCE_READY'
  | 'WAITING_RELEASE'
  | 'SELECTING'
  | 'ROUND_READY'
  | 'COMPLETED'
  | 'SKIPPED'
  | string
export type OrderSelectionStatus =
  | 'INELIGIBLE'
  | 'WAITING'
  | 'CURRENT'
  | 'SELECTED'
  | 'PASSED'
  | 'ADMIN_SKIPPED'
  | 'INELIGIBLE_BANKRUPT'
  | string
export type OrderDeliveryStatus = 'SELECTED' | 'DELIVERED' | 'UNFINISHED' | string
export type OrderDeliveryStageCode = 'Q1' | 'Q2' | 'Q3' | 'Q4'

export interface OrderTemplateMarket {
  code: OrderMarketCode
  name: string
  defaultEnabled: boolean
  sortOrder: number
}

export interface OrderTemplateOrderType {
  code: OrderTypeCode
  name: string
  sortOrder: number
}

export interface OrderTemplateField {
  code: string
  name: string
  valueType: string
  unit?: string
  order: number
}

export interface OrderTemplateMeta {
  templateVersion: string
  templateName: string
  formulaVersion: string
  segmentCount: number
  maxCardCount: number
  maxSelectionRounds: number
  deliveryEnabled: boolean
  markets: OrderTemplateMarket[]
  orderTypes: OrderTemplateOrderType[]
  fields: OrderTemplateField[]
}

export interface OrderMarketForecastProduct {
  orderType: OrderTypeCode
  orderTypeName: string
  orderCount: number
  forecastAmount: number
}

export interface OrderMarketForecastYear {
  yearNo: number
  products: OrderMarketForecastProduct[]
  totalOrderCount: number
  totalForecastAmount: number
}

export interface OrderMarketForecastMarket {
  marketCode: OrderMarketCode
  marketName: string
  narrative: string
  years: OrderMarketForecastYear[]
}

export interface OrderMarketForecastStage {
  forecastStageCode: OrderForecastStageCode
  forecastStageName: string
  yearRange: string
  years: number[]
  markets: OrderMarketForecastMarket[]
}

export interface OrderMarketForecastResult {
  formulaVersion: string
  orderTemplate: OrderTemplateMeta
  stages: OrderMarketForecastStage[]
}

export interface PlayerOrderPoolItem {
  orderId: number
  businessOrderNo: string
  cardSequenceNo: number
  orderAmount: number
  orderQuantity: number
  unitPrice: number
  accountTerm: number
  poolStatus: OrderPoolStatus
  roundNo?: number
  deliveryStatus?: OrderDeliveryStatus
  deliveredStageCode?: string | null
  orderPayload?: Record<string, unknown>
}

export interface PlayerOrderSequenceView {
  sequenceNo: number
  groupId: number
  groupName: string
  selectionStatus: OrderSelectionStatus
  isSelf: boolean
  isMarketLeader: boolean
}

export interface PlayerOrderSegmentView {
  marketCode: OrderMarketCode
  marketName: string
  orderType: OrderTypeCode
  orderTypeName: string
  marketEnabled: boolean
  marketInvestmentLimit: number | null
  marketInvestment: number
  investmentSubmitted: boolean
  releaseSequenceNo: number
  segmentStatus: OrderSegmentStatus
  currentRoundNo: number
  nextRoundNo?: number | null
  selfRoundStatus: OrderSelectionStatus | ''
  ordersVisible: boolean
  selectionOrder: PlayerOrderSequenceView[]
  currentGroupId: number | null
  availableOrders: PlayerOrderPoolItem[]
  lockedOrders: PlayerOrderPoolItem[]
  selectedOrder: PlayerOrderPoolItem | null
  selectedOrders: PlayerOrderPoolItem[]
  deliveryStatus: OrderDeliveryStatus | ''
  canSelectOrder: boolean
  canPassSegment: boolean
}

export interface PlayerOrderMarketView {
  marketCode: OrderMarketCode
  marketName: string
  marketBidStatus: OrderSegmentStatus | ''
  investmentSubmitted: boolean
  marketInvestment: number
  canSubmitInvestment: boolean
  canSelectOrder: boolean
  selectionSequenceNo: number | null
  isMarketLeader: boolean
  marketEnabled: boolean
  marketInvestmentLimit: number | null
  segments: PlayerOrderSegmentView[]
}

export interface PlayerOrderYearView {
  groupId: number
  yearNo: number
  orderTemplate: OrderTemplateMeta
  orderRequired: boolean
  investmentSubmitted: boolean
  canSubmitInvestment: boolean
  markets: PlayerOrderMarketView[] | null
  pollingIntervalSeconds: number
}

export interface MarketInvestmentInput {
  marketCode: OrderMarketCode
  orderType: OrderTypeCode
  marketInvestment: number
}

export interface SubmitMarketInvestmentRequest {
  yearNo: number
  investments: MarketInvestmentInput[]
}

export interface SubmitMarketInvestmentResult {
  groupId: number
  yearNo: number
  submittedCount: number
  submittedAt: string
}

export interface SelectOrderRequest {
  yearNo: number
  marketCode: OrderMarketCode
  orderType: OrderTypeCode
  orderId: number
}

export interface SelectOrderResult {
  selectedOrderId: number
  marketCode: OrderMarketCode
  orderType: OrderTypeCode
  roundNo: number
  selectionSequenceNo: number
  nextGroupId: number | null
  segmentStatus: OrderSegmentStatus
}

export interface PassOrderSegmentRequest {
  yearNo: number
  marketCode: OrderMarketCode
  orderType: OrderTypeCode
}

export interface PassOrderSegmentResult {
  marketCode: OrderMarketCode
  orderType: OrderTypeCode
  nextGroupId: number | null
  segmentStatus: OrderSegmentStatus
}

export interface DeliverOrdersRequest {
  yearNo: number
  stageCode: OrderDeliveryStageCode
  orderIds: number[]
}

export interface DeliverOrdersResult {
  groupId: number
  yearNo: number
  stageCode: string
  orderIds: number[]
  deliveredAmount: number
  deliveredAt: string
}

export interface AdminOrderControlResult {
  yearNo: number
  marketCode?: OrderMarketCode
  marketName?: string
  orderType?: OrderTypeCode
  orderTypeName?: string
  segmentStatus?: OrderSegmentStatus
  currentGroupId?: number | null
  currentRoundNo?: number
  affectedCount: number
  operatedAt: string
  operatedBy: string
}

export interface AdminMarketBidView {
  groupId: number
  groupNo: number
  groupName: string
  marketInvestment: number
  submitted: boolean
  businessStatus: string
}

export interface AdminSelectionOrderView {
  roundNo: number
  sequenceNo: number
  groupId: number
  groupName: string
  marketInvestment: number
  previousMarketOrderAmount: number
  isMarketLeader: boolean
  selectionStatus: OrderSelectionStatus
  selectedOrderId: number | null
  selectedOrderNo?: string
}

export interface AdminOrderSegmentStatus {
  marketCode: OrderMarketCode
  marketName: string
  orderType: OrderTypeCode
  orderTypeName: string
  releaseSequenceNo: number
  segmentStatus: OrderSegmentStatus
  currentGroupId: number | null
  currentRoundNo: number
  nextRoundNo?: number | null
  completionReason?: string | null
  selectionOrder: AdminSelectionOrderView[]
  availableCount: number
  selectedCount: number
  theoreticalMaxSelections: number
  warnings: string[]
}

export interface AdminMarketSelectionStatus {
  yearNo: number
  orderTemplate: OrderTemplateMeta
  marketCode: OrderMarketCode
  marketName: string
  marketBidStatus: OrderSegmentStatus | ''
  leaderGroupId: number | null
  bids: AdminMarketBidView[]
  segments: AdminOrderSegmentStatus[]
  currentSegment: AdminOrderSegmentStatus | null
}
