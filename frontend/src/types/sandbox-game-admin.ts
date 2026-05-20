export type AdminGroupDataPageType = 'operating' | 'report'

export type UnlockTargetType = 'OPERATING' | 'REPORT'
export type UnlockStageCode = 'Q1' | 'Q2' | 'Q3' | 'Q4' | 'YEAR_END'
export type NoticeTargetScope = 'ALL' | 'GROUP'
export type AdjustmentType = 'REWARD' | 'PENALTY'
import type { OrderMarketCode, OrderPoolStatus, OrderTypeCode } from '@/types/sandbox-game-order'

export type { OrderMarketCode, OrderPoolStatus }
export type AdminOrderType = OrderTypeCode
export type OrderConfigStatus = 'DRAFT' | 'LOCKED'
export type OrderGenerationStatus = 'NOT_GENERATED' | 'PREVIEW_GENERATED' | 'POOL_CONFIRMED' | 'SELECTING' | 'COMPLETED' | string
export type OrderGenerationBatchStatus = 'PREVIEW' | 'CONFIRMED' | 'VOID' | string

export interface AdminActionSummary {
  actionCode: string
  operatorName: string
  operateTime: string
  targetGroupId?: number | null
  targetYearNo?: number | null
}

export interface AdminControlSetupStatusResult {
  initialized: boolean
  groupCount: number
  finalYear: number
  currentOpenYear: number
  initialBaselineSubmitted: boolean
  defaultRoute: string
}

export interface AdminControlConfigResult {
  finalYear: number
  currentOpenYear: number
  canOpenNextYear: boolean
  nextOpenableYear: number
  openNextYearBlockedReason: string
  ruleVersion: string
  templateVersion: string
  initialBaselineSubmitted: boolean
  initialBaselineSubmittedAt: string | null
  initialBaselineSubmitterName: string | null
  latestAdminAction: AdminActionSummary | null
}

export interface InitializeGameResult {
  initialized: boolean
  groupCount: number
  createdGroupCount: number
  createdAccountCount: number
  createdYearStateCount: number
  initializedAt: string
  initializedBy: string
}

export interface UpdateFinalYearResult {
  finalYear: number
  currentOpenYear: number
  initializedFromYear: number | null
  initializedToYear: number | null
  initializedYearCount: number
  updatedAt: string
  updatedBy: string
}

export interface OpenNextYearResult {
  previousOpenYear: number
  currentOpenYear: number
  openedYearNo: number
  finalYear: number
  canOpenNextYear: boolean
  openNextYearBlockedReason: string
  latestAdminAction: AdminActionSummary | null
}

export interface UnlockYearRequest {
  groupId: number
  yearNo: number
  unlockTargetType: UnlockTargetType
  targetStageCode?: UnlockStageCode | null
  reason: string
}

export interface UnlockYearResult {
  groupId: number
  yearNo: number
  unlockTargetType: UnlockTargetType
  targetStageCode: UnlockStageCode | null
  editableStageCode: UnlockStageCode | null
  yearStatus: string
  stageStatus: string
  reportStatus: string
  summaryEffective: boolean
  businessStatus: string
  unlockLogId: number
}

export interface BaselinePayload {
  baselineSalesRevenue: number
  baselineDirectCost: number
  baselineComprehensiveCost: number
  baselineDepreciation: number
  baselineFinanceIncomeExpense: number
  baselineExtraIncomeExpense: number
  baselineIncomeTax: number
  baselineWorkInConstruction: number
  baselineFactoryAsset: number
  baselineLineResidual: number
  baselineDepreciableAsset: number
  baselineCash: number
  baselineReceivable: number
  baselineWorkInProgress: number
  baselineFinishedGoods: number
  baselineRawMaterials: number
  baselineShortTermLoan: number
  baselineLongTermLoan: number
  baselineShareCapital: number
  baselineRetainedEarnings: number
}

export interface InitialBaselineViewResult {
  submitted: boolean
  editable: boolean
  baselinePayload: BaselinePayload
  appliedGroupCount: number
  submitterName: string | null
  submittedAt: string | null
}

export interface SubmitInitialBaselineRequest {
  baselinePayload: BaselinePayload
}

export interface SubmitInitialBaselineResult {
  submitted: boolean
  appliedGroupCount: number
  submittedAt: string
  submitterName: string
}

export interface AdminYearSummaryItem {
  groupId: number
  groupNo: number
  groupName: string
  revenue: number
  profit: number
  equity: number
  businessStatus: string
  ranking?: number | null
}

export interface AdminYearSummaryResult {
  yearNo: number
  list: AdminYearSummaryItem[]
}

export interface AdminFinalRankingItem {
  ranking: number
  groupId: number
  groupNo: number
  groupName: string
  equity: number
  revenue: number
  profit: number
  businessStatus: string
}

export interface AdminFinalRankingResult {
  finalYear: number
  list: AdminFinalRankingItem[]
}

export interface AdminGroupOption {
  groupId: number
  groupNo: number
  groupName: string
  businessStatus: string
}

export interface ListAdminGroupsResult {
  list: AdminGroupOption[]
}

export interface AdminGeneralNoticeRecord {
  id: number
  targetScope: NoticeTargetScope
  targetGroupId: number | null
  targetGroupName: string | null
  content: string
  pinned: boolean
  publishedAt: string
  operatorName: string
}

export interface AdminAdjustmentRecord {
  id: number
  groupId: number
  groupNo: number
  groupName: string
  yearNo: number
  stageCode: 'Q1' | 'Q2' | 'Q3' | 'Q4'
  adjustmentType: AdjustmentType
  amount: number
  reason: string
  publishedAt: string
  operatorName: string
}

export interface AdminNoticeRecordsResult {
  generalNotices: AdminGeneralNoticeRecord[]
  adjustments: AdminAdjustmentRecord[]
}

export interface SendAdminGeneralNoticeRequest {
  targetScope: NoticeTargetScope
  targetGroupId?: number | null
  content: string
  pinned: boolean
}

export interface SendAdminGeneralNoticeResult {
  noticeId: number
  targetScope: NoticeTargetScope
  targetGroupId: number | null
  content: string
  pinned: boolean
  publishedAt: string
}

export interface SendAdminAdjustmentRequest {
  groupId: number
  yearNo: number
  stageCode: 'Q1' | 'Q2' | 'Q3' | 'Q4'
  adjustmentType: AdjustmentType
  amount: number
  reason: string
}

export interface SendAdminAdjustmentResult {
  adjustmentId: number
  groupId: number
  yearNo: number
  stageCode: 'Q1' | 'Q2' | 'Q3' | 'Q4'
  adjustmentType: AdjustmentType
  amount: number
  reason: string
  publishedAt: string
}

export interface ParsedOrderCard {
  yearNo: number
  marketCode: OrderMarketCode
  marketName: string
  orderType: AdminOrderType
  orderTypeName: string
  orderAmount: number
  orderQuantity: number
  unitPrice: number
  accountTerm: number
  sourceSheetName: string
  sourceCell: string
  sourceRowIndex: number
}

export interface OrderSourceSummary {
  yearNo: number
  marketCode: OrderMarketCode
  marketName: string
  orderType: AdminOrderType
  orderTypeName: string
  availableCount: number
}

export interface UploadOrderExcelResult {
  batchId: number
  originalFileName: string
  parseStatus: string
  parsedOrderCount: number
  warnings: string[]
  summary: OrderSourceSummary[]
  previewOrders: ParsedOrderCard[]
  uploadedAt: string
}

export interface OrderControlConfigItem {
  yearNo: number
  marketCode: OrderMarketCode
  marketName: string
  marketEnabled: boolean
  orderType: AdminOrderType
  orderTypeName: string
  orderCount: number
  releaseSequenceNo: number
  availableCount: number
  configStatus: OrderConfigStatus
  generatedCount: number
}

export interface OrderMarketConfigItem {
  yearNo: number
  marketCode: OrderMarketCode
  marketName: string
  enabled: boolean
  configStatus: OrderConfigStatus
  lockedBatchId?: number | null
}

export interface OrderControlWarning {
  level: 'INFO' | 'WARN' | 'ERROR' | string
  message: string
  marketCode?: OrderMarketCode
  orderType?: AdminOrderType
}

export interface OrderGenerationBatchSummary {
  batchId: number
  batchStatus: OrderGenerationBatchStatus
  formulaVersion: string
  randomSeed?: string
  generatedCount: number
  generatedAt: string
  generatedBy: string
  confirmedAt?: string | null
  confirmedBy?: string | null
}

export interface OrderControlConfigResult {
  yearNo: number
  finalYear: number
  latestBatchId: number | null
  latestBatchUploadedAt: string | null
  generationStatus: OrderGenerationStatus
  latestPreviewBatch: OrderGenerationBatchSummary | null
  confirmedBatch: OrderGenerationBatchSummary | null
  canUpdateConfig: boolean
  canGeneratePreview: boolean
  canConfirmPool: boolean
  marketConfigs: OrderMarketConfigItem[]
  items: OrderControlConfigItem[]
  warnings: OrderControlWarning[]
}

export interface UpdateOrderMarketConfigRequest {
  yearNo: number
  markets: Array<{
    marketCode: OrderMarketCode
    enabled: boolean
  }>
}

export interface UpdateOrderMarketConfigResult {
  yearNo: number
  markets: OrderMarketConfigItem[]
  items: OrderControlConfigItem[]
  warnings: OrderControlWarning[]
  updatedAt: string
  updatedBy: string
}

export interface UpdateOrderControlConfigRequest {
  yearNo: number
  items: Array<{
    marketCode: OrderMarketCode
    orderType: AdminOrderType
    orderCount: number
    releaseSequenceNo: number
  }>
}

export interface UpdateOrderControlConfigResult {
  yearNo: number
  items: OrderControlConfigItem[]
  warnings: OrderControlWarning[]
  updatedAt: string
  updatedBy: string
}

export interface GenerateOrderPoolRequest {
  yearNo: number
  overwrite: boolean
  sourceBatchId?: number | null
}

export interface GenerateOrderPoolResult {
  yearNo: number
  sourceBatchId?: number
  batchId: number
  batchStatus: OrderGenerationBatchStatus
  randomSeed: string
  formulaVersion: string
  generatedCount: number
  segmentCount: number
  warnings: OrderControlWarning[]
  generatedAt: string
  generatedBy: string
}

export interface ConfirmOrderPoolRequest {
  yearNo: number
  batchId: number
}

export interface ConfirmOrderPoolResult {
  yearNo: number
  batchId: number
  generatedCount: number
  segmentCount: number
  confirmedAt: string
  confirmedBy: string
}

export interface OrderPoolItem {
  orderId: number
  businessOrderNo: string
  cardSequenceNo: number
  yearNo: number
  marketCode: OrderMarketCode
  marketName: string
  orderType: AdminOrderType
  orderTypeName: string
  orderAmount: number
  orderQuantity: number
  unitPrice: number
  accountTerm: number
  poolStatus: OrderPoolStatus
  selectedGroupId: number | null
  sourceSheetName: string
  sourceCell: string
}

export interface OrderPoolResult {
  yearNo: number
  marketCode: OrderMarketCode | ''
  marketName: string
  orderType: AdminOrderType | ''
  orderTypeName: string
  list: OrderPoolItem[]
}

export function createEmptyBaselinePayload(): BaselinePayload {
  return {
    baselineSalesRevenue: 0,
    baselineDirectCost: 0,
    baselineComprehensiveCost: 0,
    baselineDepreciation: 0,
    baselineFinanceIncomeExpense: 0,
    baselineExtraIncomeExpense: 0,
    baselineIncomeTax: 0,
    baselineWorkInConstruction: 0,
    baselineFactoryAsset: 0,
    baselineLineResidual: 0,
    baselineDepreciableAsset: 0,
    baselineCash: 0,
    baselineReceivable: 0,
    baselineWorkInProgress: 0,
    baselineFinishedGoods: 0,
    baselineRawMaterials: 0,
    baselineShortTermLoan: 0,
    baselineLongTermLoan: 0,
    baselineShareCapital: 0,
    baselineRetainedEarnings: 0,
  }
}

export function cloneBaselinePayload(payload?: BaselinePayload | null): BaselinePayload {
  return {
    ...createEmptyBaselinePayload(),
    ...(payload ?? {}),
  }
}
