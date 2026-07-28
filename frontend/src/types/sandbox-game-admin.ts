export type AdminGroupDataPageType = 'operating' | 'report'

export type UnlockTargetType = 'OPERATING' | 'REPORT'
export type UnlockStageCode = 'Q1' | 'Q2' | 'Q3' | 'Q4' | 'YEAR_END'
export type SnapshotScope = 'GROUP' | 'GLOBAL'
export type SnapshotType = 'AUTO' | 'MANUAL' | 'SAFETY'
export type NoticeTargetScope = 'ALL' | 'GROUP'
export type AdjustmentType = 'REWARD' | 'PENALTY'
import type {
  OrderForecastStageCode,
  OrderMarketCode,
  OrderMarketForecastResult,
  OrderPoolStatus,
  OrderTemplateMeta,
  OrderTypeCode,
} from '@/types/sandbox-game-order'

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

export interface GameEdition {
  editionCode: string
  editionName: string
  description: string
  defaultEdition: boolean
  ruleVersion: string
  formulaVersion: string
  templateVersion: string
  operatingTemplateVersion: string
  reportTemplateVersion: string
  orderTemplateVersion: string
  processRuleVersion: string
}

export interface AdminControlSetupStatusResult {
  initialized: boolean
  groupCount: number
  finalYear: number
  currentOpenYear: number
  editionCode: string
  editionName: string
  ruleVersion: string
  formulaVersion: string
  templateVersion: string
  operatingTemplateVersion: string
  reportTemplateVersion: string
  orderTemplateVersion: string
  processRuleVersion: string
  dictionaryRevision: number
  availableEditions: GameEdition[]
  initialBaselineSubmitted: boolean
  defaultRoute: string
}

export interface AdminControlConfigResult {
  finalYear: number
  currentOpenYear: number
  canOpenNextYear: boolean
  nextOpenableYear: number
  openNextYearBlockedReason: string
  editionCode: string
  editionName: string
  ruleVersion: string
  formulaVersion: string
  templateVersion: string
  operatingTemplateVersion: string
  reportTemplateVersion: string
  orderTemplateVersion: string
  processRuleVersion: string
  dictionaryRevision: number
  initialBaselineSubmitted: boolean
  initialBaselineSubmittedAt: string | null
  initialBaselineSubmitterName: string | null
  latestAdminAction: AdminActionSummary | null
}

export interface InitializeGameResult {
  initialized: boolean
  groupCount: number
  editionCode: string
  editionName: string
  ruleVersion: string
  templateVersion: string
  dictionaryRevision: number
  createdGroupCount: number
  createdAccountCount: number
  createdYearStateCount: number
  initializedAt: string
  initializedBy: string
}

export interface AdminDictionaryItem {
  itemCode: string
  itemCategory: string
  defaultName: string
  displayName: string
  displayOrder: number
  editable: boolean
  relatedPayload?: unknown
}

export interface AdminDictionaryCurrentResult {
  initialized: boolean
  editionCode: string
  editionName: string
  dictionaryRevision: number
  canUpdateCurrent: boolean
  items: AdminDictionaryItem[]
}

export interface AdminDictionarySchemeSummary {
  id: number
  editionCode: string
  schemeName: string
  description: string
  builtIn: boolean
  itemCount: number
  updatedAt: string
  updatedBy: string
}

export interface AdminDictionarySchemeListResult {
  editionCode: string
  list: AdminDictionarySchemeSummary[]
}

export interface AdminDictionarySchemeDetailResult {
  id: number
  editionCode: string
  schemeName: string
  description: string
  builtIn: boolean
  items: AdminDictionaryItem[]
  updatedAt: string
  updatedBy: string
}

export interface AdminDictionaryItemInput {
  itemCode: string
  displayName: string
}

export interface AdminDictionarySaveSchemeRequest {
  schemeId?: number | null
  editionCode: string
  schemeName: string
  description: string
  items: AdminDictionaryItemInput[]
}

export interface AdminDictionaryUpdateCurrentRequest {
  items: AdminDictionaryItemInput[]
  reason: string
}

export interface AdminDictionaryApplySchemeRequest {
  schemeId: number
}

export interface AdminDictionaryRevisionResult {
  editionCode: string
  dictionaryRevision: number
}

export interface AdminDictionaryChangeLogItem {
  id: number
  editionCode: string
  changeType: string
  schemeId: number | null
  reason: string
  revision: number
  operatorName: string
  operateTime: string
  changedSummary: string
}

export interface AdminDictionaryChangeLogPageResult {
  list: AdminDictionaryChangeLogItem[]
  pageNo: number
  pageSize: number
  total: number
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

export interface AdminRollbackPendingGroup {
  groupId: number
  groupNo: number
  groupName: string
  businessStatus: string
  rollbackTargetYearNo: number | null
  rollbackTargetStageCode: string | null
}

export interface AdminYearSummaryResult {
  yearNo: number
  list: AdminYearSummaryItem[]
  pendingGroups: AdminRollbackPendingGroup[]
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

export type AdjustmentStageCode = 'Q1' | 'Q2' | 'Q3' | 'Q4' | 'YEAR_END'
export type AdjustmentRecordStatus = 'EFFECTIVE' | 'VOIDED' | 'SNAPSHOT_INACTIVE'

export interface AdminAdjustmentRecord {
  id: number
  groupId: number
  groupNo: number
  groupName: string
  yearNo: number
  stageCode: AdjustmentStageCode
  adjustmentType: AdjustmentType
  amount: number
  reason: string
  publishedAt: string
  operatorName: string
  status: AdjustmentRecordStatus
  canVoid: boolean
  adjustmentRevision: number
  voidedByName: string | null
  voidReason: string | null
  voidedAt: string | null
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
  adjustmentType: AdjustmentType
  amount: number
  reason: string
}

export interface AdjustmentImpactResult {
  resolvedStageCode: AdjustmentStageCode
  cashBefore: number
  cashAfter: number
  preTaxProfitAfter: number
  incomeTaxAfter: number
  netProfitAfter: number
  totalEquityAfter: number
  willBankrupt: boolean
  calculationBasisSavedAt: string | null
}

export type PreviewAdminAdjustmentRequest =
  | ({ operation: 'CREATE' } & SendAdminAdjustmentRequest)
  | { operation: 'VOID'; adjustmentId: number }

export interface VoidAdminAdjustmentRequest {
  adjustmentId: number
  reason: string
}

export interface SendAdminAdjustmentResult {
  adjustmentId: number
  groupId: number
  yearNo: number
  stageCode: AdjustmentStageCode
  adjustmentType: AdjustmentType
  amount: number
  reason: string
  publishedAt: string
  adjustmentRevision: number
  impact: AdjustmentImpactResult
}

export interface VoidAdminAdjustmentResult {
  adjustmentId: number
  groupId: number
  yearNo: number
  adjustmentRevision: number
  voidedAt: string
  impact: AdjustmentImpactResult
}

export interface SnapshotSummary {
  id: number
  snapshotScope: SnapshotScope | string
  snapshotType: SnapshotType | string
  triggerCode: string
  groupId: number | null
  groupName: string | null
  yearNo: number | null
  stageCode: string | null
  reportStatus: string | null
  description: string
  payloadHash: string
  createdByName: string
  createdAt: string
  canRestore: boolean
}

export interface SnapshotListResult {
  list: SnapshotSummary[]
  pageNo: number
  pageSize: number
  total: number
}

export interface SnapshotDetailResult {
  snapshot: SnapshotSummary
  stateSummary: Record<string, unknown>
  payloadPreview: Record<string, unknown>
  payloadVersion: string
  payloadSize: number
}

export interface CreateSnapshotRequest {
  snapshotScope: SnapshotScope
  groupId?: number | null
  yearNo?: number | null
  stageCode?: string
  description: string
}

export interface RestoreGroupSnapshotRequest {
  snapshotId: number
  reason: string
  confirmText: string
}

export interface RestoreGroupSnapshotResult {
  rollbackLogId: number
  safetySnapshotId: number
  groupId: number
  targetYearNo: number
  targetStageCode: string
  yearStatus: string
  stageStatus: string
  reportStatus: string
  businessStatus: string
  hasRollbackPending: boolean
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

export interface OrderForecastControlItem {
  yearNo: number
  forecastStageCode: OrderForecastStageCode
  forecastStageName: string
  marketCode: OrderMarketCode
  marketName: string
  orderType: AdminOrderType
  orderTypeName: string
  orderCount: number
}

export interface OrderForecastNarrativeItem {
  forecastStageCode: OrderForecastStageCode
  forecastStageName: string
  marketCode: OrderMarketCode
  marketName: string
  content: string
}

export interface OrderForecastYearLock {
  yearNo: number
  locked: boolean
  reason?: string
}

export interface OrderForecastControlResult {
  orderTemplate: OrderTemplateMeta
  items: OrderForecastControlItem[]
  narratives: OrderForecastNarrativeItem[]
  forecast: OrderMarketForecastResult
  yearLocks: OrderForecastYearLock[]
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
  marketInvestmentLimit: number | null
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
  orderTemplate: OrderTemplateMeta
  latestBatchId: number | null
  latestBatchUploadedAt: string | null
  forecast: OrderMarketForecastResult
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
    marketInvestmentLimit: number | null
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

export interface UpdateOrderForecastControlRequest {
  items: Array<{
    yearNo: number
    marketCode: OrderMarketCode
    orderType: AdminOrderType
    orderCount: number
  }>
  narratives: Array<{
    forecastStageCode: OrderForecastStageCode
    marketCode: OrderMarketCode
    content: string
  }>
}

export interface UpdateOrderForecastControlResult extends OrderForecastControlResult {
  autoPreviewYears: number[]
  updatedAt: string
  updatedBy: string
}

export interface UpdateOrderControlConfigRequest {
  yearNo: number
  items: Array<{
    marketCode: OrderMarketCode
    orderType: AdminOrderType
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
  selectedGroupName: string
  sourceSheetName: string
  sourceCell: string
  orderPayload?: Record<string, unknown>
}

export interface OrderPoolResult {
  yearNo: number
  orderTemplate: OrderTemplateMeta
  marketCode: OrderMarketCode | ''
  marketName: string
  orderType: AdminOrderType | ''
  orderTypeName: string
  selectedOnly: boolean
  selectedGroupId: number | null
  groupOptions: AdminGroupOption[]
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
