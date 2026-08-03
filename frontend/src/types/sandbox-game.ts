export type NumericCellValue = number | string | ''
export type CellValue = number | string | boolean | null
export type QuarterValueMap = Record<string, Record<string, NumericCellValue>>

export interface ProjectProgressItem {
  projectName: string
  factoryKey?: string
  factoryLabel?: string
  slotNo?: number
  quarterKey?: string
  quarterLabel?: string
  lineType: string
  progress: NumericCellValue
}

export interface ProjectProgressUpdatePayload {
  items: ProjectProgressItem[]
}

export interface MarketCultivationItem {
  annualInvestment: NumericCellValue
  effectiveAnnualInvestment: number
  previousCumulative: number
  cumulativeInvestment: number
  unlocked: boolean
  locked: boolean
  lockedByPrevious: boolean
  willUnlock: boolean
}

export interface MarketCultivationPayload {
  regional: MarketCultivationItem
  national: MarketCultivationItem
  global: MarketCultivationItem
  annualTotal: number
  stateApplied: boolean
}

export interface QualificationItem {
  status: string
  unlocked: boolean
  locked: boolean
  lockedByPrevious: boolean
}

export interface QualificationCertificationPayload {
  qualityEnvironmentalHealth: QualificationItem
  highTechEnterprise: QualificationItem
  specializedInnovation: QualificationItem
  listedCompany: QualificationItem
}

export interface PlayerNoticeItem {
  id: number
  kind: string
  title: string
  content: string
  pinned: boolean
  publishedAt: string
  yearNo?: number | null
  stageCode?: string | null
  amount?: number | null
  status?: 'EFFECTIVE' | 'VOIDED' | 'SNAPSHOT_INACTIVE' | null
}

export interface PlayerNoticeBoard {
  pinnedNotice: PlayerNoticeItem | null
  recentList: PlayerNoticeItem[]
}

export interface CurrentGameConfigResult {
  currentOpenYear: number
  finalYear: number
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
  demoYearEnabled: boolean
}

export interface PlayerAdjustmentSyncResult {
  notModified: boolean
  adjustmentRevision: number
  incomeAndPenalty?: QuarterValueMap
  derivedValues?: Record<string, number>
  quarterCashChecks?: Record<string, number>
  periodEndCash?: number
  reportComputedPayload?: ReportComputedPayload
  balanceGap?: number
  balanceCheckPassed?: boolean
  noticeBoard?: PlayerNoticeBoard | null
  businessStatus?: string
  bankrupt?: boolean
}

export interface YearTabItem {
  yearNo: number
  label: string
  tabStatus: string
  canEnter: boolean
  isCurrentOpenYear: boolean
  isFormalYear: boolean
}

export interface YearTabsResult {
  currentOpenYear: number
  finalYear: number
  demoYearEnabled: boolean
  tabs: YearTabItem[]
}

export interface OperatingPayload {
  beginning: {
    taxAndPlanning: Record<string, CellValue>
    marketBid: Array<Record<string, CellValue>>
  }
  quarter: {
    shortTermLoan: QuarterValueMap
    materialPayment: QuarterValueMap
    productionLineAdjustment: QuarterValueMap
    humanResource: QuarterValueMap
    salaryAndProduction: QuarterValueMap
    researchAndManagement: QuarterValueMap
    supplyChainOrderRecord: QuarterValueMap
    receivableUpdate: QuarterValueMap
    deliverySettlement: QuarterValueMap
  }
  yearEnd: {
    longTermLoan: Record<string, NumericCellValue>
    assetAdjustment: Record<string, NumericCellValue>
    projectProgressUpdate: ProjectProgressUpdatePayload
    marketCultivation: MarketCultivationPayload
    qualificationCertification: QualificationCertificationPayload
  }
  extra: {
    incomeAndPenalty: QuarterValueMap
  }
  derived: {
    values: Record<string, NumericCellValue>
  }
}

export interface PlayerOperatingStageSubmitHistory {
  stageCode: string
  submitVersion: number
  periodEndCash: number
  submitTime: string
}

export interface OperatingCarryForward {
  previousIncomeTax: number
  previousShortTermLoan: number
  previousLongTermLoan: number
  previousEquipmentResidual: number
  previousDepreciableAsset: number
  previousCash: number
  previousReceivable: number
  shareholderCapital: number
  retainedEarnings: number
}

export interface PlayerOperatingView {
  groupId: number
  yearNo: number
  yearStatus: string
  stageStatus: string
  reportStatus: string
  businessStatus: string
  currentStageCode: string
  canView: boolean
  canEdit: boolean
  canSubmit: boolean
  hasInvalidDraft: boolean
  invalidScopes: string[]
  hasRetainedReportDraft: boolean
  rollbackPending: boolean
  rollbackTargetYearNo?: number | null
  rollbackTargetStageCode?: string | null
  rollbackNotice: string
  operatingPayload: OperatingPayload
  editableScopes: string[]
  readonlyScopes: string[]
  stageSubmitHistory: PlayerOperatingStageSubmitHistory[]
  lastDraftSavedAt: string | null
  quarterCashChecks: Record<string, number>
  derivedValues: Record<string, number>
  periodEndCash: number
  carryForward?: OperatingCarryForward | null
  noticeBoard: PlayerNoticeBoard | null
  adjustmentRevision: number
}

export interface SaveDraftRequest {
  yearNo: number
  stageStatus: string
  operatingPayload: OperatingPayload
  clientSaveTime: string
}

export interface SaveDraftResponse {
  groupId: number
  yearNo: number
  stageStatus: string
  lastDraftSavedAt: string
}

export interface SubmitStageRequest {
  yearNo: number
  stageCode: string
  operatingPayload: OperatingPayload
}

export interface SubmitStageResponse {
  groupId: number
  yearNo: number
  stageCode: string
  yearStatus: string
  stageStatus: string
  reportStatus: string
  businessStatus: string
  latestStageSubmitVersion: number
  periodEndCash: number
  submittedAt: string
}

const projectProgressFactoryDefinitions = [
  { key: 'factoryA', label: '生产厂房 A', slotCount: 4 },
  { key: 'factoryB', label: '生产厂房 B', slotCount: 3 },
  { key: 'factoryC', label: '生产厂房 C', slotCount: 1 },
] as const

const projectProgressQuarterDefinitions = [
  { key: 'q1', label: '第一季度' },
  { key: 'q2', label: '第二季度' },
  { key: 'q3', label: '第三季度' },
  { key: 'q4', label: '第四季度' },
] as const

export function createEmptyOperatingPayload(): OperatingPayload {
  return {
    beginning: {
      taxAndPlanning: {},
      marketBid: [{ marketInvestment: '', orderAmount: '' }],
    },
    quarter: {
      shortTermLoan: {},
      materialPayment: {},
      productionLineAdjustment: {},
      humanResource: {},
      salaryAndProduction: {},
      researchAndManagement: {},
      supplyChainOrderRecord: {},
      receivableUpdate: {},
      deliverySettlement: {},
    },
    yearEnd: {
      longTermLoan: {},
      assetAdjustment: {},
      projectProgressUpdate: createEmptyProjectProgressUpdatePayload(),
      marketCultivation: createEmptyMarketCultivationPayload(),
      qualificationCertification: createEmptyQualificationCertificationPayload(),
    },
    extra: {
      incomeAndPenalty: {},
    },
    derived: {
      values: {},
    },
  }
}

export function cloneOperatingPayload(payload?: OperatingPayload | null): OperatingPayload {
  if (!payload) {
    return createEmptyOperatingPayload()
  }
  const cloned = JSON.parse(JSON.stringify(payload)) as OperatingPayload
  cloned.quarter.supplyChainOrderRecord = cloned.quarter.supplyChainOrderRecord ?? {}
  cloned.yearEnd.projectProgressUpdate = normalizeProjectProgressUpdatePayload(cloned.yearEnd.projectProgressUpdate)
  cloned.yearEnd.marketCultivation = normalizeMarketCultivationPayload(cloned.yearEnd.marketCultivation)
  cloned.yearEnd.qualificationCertification = normalizeQualificationCertificationPayload(cloned.yearEnd.qualificationCertification)
  return cloned
}

export function createEmptyProjectProgressUpdatePayload(): ProjectProgressUpdatePayload {
  return {
    items: createDefaultProjectProgressItems(),
  }
}

export function createEmptyMarketCultivationPayload(): MarketCultivationPayload {
  return {
    regional: createEmptyMarketCultivationItem(),
    national: createEmptyMarketCultivationItem(),
    global: createEmptyMarketCultivationItem(),
    annualTotal: 0,
    stateApplied: false,
  }
}

export function createEmptyMarketCultivationItem(): MarketCultivationItem {
  return {
    annualInvestment: '',
    effectiveAnnualInvestment: 0,
    previousCumulative: 0,
    cumulativeInvestment: 0,
    unlocked: false,
    locked: false,
    lockedByPrevious: false,
    willUnlock: false,
  }
}

export function createEmptyQualificationCertificationPayload(): QualificationCertificationPayload {
  return {
    qualityEnvironmentalHealth: createEmptyQualificationItem(),
    highTechEnterprise: createEmptyQualificationItem(),
    specializedInnovation: createEmptyQualificationItem(),
    listedCompany: createEmptyQualificationItem(),
  }
}

export function createEmptyQualificationItem(): QualificationItem {
  return {
    status: '未解锁',
    unlocked: false,
    locked: false,
    lockedByPrevious: false,
  }
}

function normalizeProjectProgressUpdatePayload(payload?: ProjectProgressUpdatePayload | null): ProjectProgressUpdatePayload {
  if (!payload || !Array.isArray(payload.items)) {
    return createEmptyProjectProgressUpdatePayload()
  }
  const fallback = createDefaultProjectProgressItems()
  const currentItems = payload.items
  const keyedItems = new Map<string, ProjectProgressItem>()
  currentItems.forEach((item) => {
    const factoryKey = typeof item.factoryKey === 'string' ? item.factoryKey : ''
    const slotNo = typeof item.slotNo === 'number' ? item.slotNo : Number(item.slotNo)
    const quarterKey = typeof item.quarterKey === 'string' ? item.quarterKey : ''
    if (factoryKey && Number.isInteger(slotNo) && quarterKey) {
      keyedItems.set(projectProgressKey(factoryKey, slotNo, quarterKey), item)
    }
  })

  const legacyItems = currentItems.filter((item) => !item.factoryKey && item.slotNo === undefined && !item.quarterKey)
  const legacyFactories: string[] = projectProgressFactoryDefinitions.map((factory) => factory.key)
  const legacyQuarters: string[] = projectProgressQuarterDefinitions.map((quarter) => quarter.key)

  return {
    items: fallback.map((item, index) => {
      const keyed = keyedItems.get(projectProgressKey(item.factoryKey ?? '', item.slotNo ?? 0, item.quarterKey ?? ''))
      if (keyed) {
        return {
          ...item,
          lineType: keyed.lineType ?? '',
          progress: keyed.progress ?? '',
        }
      }

      const factoryIndex = legacyFactories.indexOf(item.factoryKey ?? '')
      const quarterIndex = legacyQuarters.indexOf(item.quarterKey ?? '')
      const legacyIndex = factoryIndex >= 0 && quarterIndex >= 0 && item.slotNo === 1 ? factoryIndex * legacyQuarters.length + quarterIndex : -1
      const legacy = legacyIndex >= 0 ? legacyItems[legacyIndex] : undefined
      const positional = currentItems.length >= fallback.length ? currentItems[index] : undefined
      const source = legacy ?? positional
      return {
        ...item,
        lineType: source?.lineType ?? '',
        progress: source?.progress ?? '',
      }
    }),
  }
}

function createDefaultProjectProgressItems(): ProjectProgressItem[] {
  return projectProgressFactoryDefinitions.flatMap((factory) =>
    Array.from({ length: factory.slotCount }, (_, slotIndex) => slotIndex + 1).flatMap((slotNo) =>
      projectProgressQuarterDefinitions.map((quarter) => ({
        projectName: `${factory.label}-槽位${slotNo}-${quarter.label}`,
        factoryKey: factory.key,
        factoryLabel: factory.label,
        slotNo,
        quarterKey: quarter.key,
        quarterLabel: quarter.label,
        lineType: '',
        progress: '',
      })),
    ),
  )
}

function projectProgressKey(factoryKey: string, slotNo: number, quarterKey: string) {
  return `${factoryKey}:${slotNo}:${quarterKey}`
}

function normalizeMarketCultivationPayload(payload?: MarketCultivationPayload | null): MarketCultivationPayload {
  const empty = createEmptyMarketCultivationPayload()
  if (!payload) {
    return empty
  }
  return {
    regional: normalizeMarketCultivationItem(payload.regional),
    national: normalizeMarketCultivationItem(payload.national),
    global: normalizeMarketCultivationItem(payload.global),
    annualTotal: normalizeNumber(payload.annualTotal),
    stateApplied: payload.stateApplied === true,
  }
}

function normalizeMarketCultivationItem(item?: MarketCultivationItem | null): MarketCultivationItem {
  return {
    ...createEmptyMarketCultivationItem(),
    ...(item ?? {}),
    effectiveAnnualInvestment: normalizeNumber(item?.effectiveAnnualInvestment),
    previousCumulative: normalizeNumber(item?.previousCumulative),
    cumulativeInvestment: normalizeNumber(item?.cumulativeInvestment),
    unlocked: item?.unlocked === true,
    locked: item?.locked === true,
    lockedByPrevious: item?.lockedByPrevious === true,
    willUnlock: item?.willUnlock === true,
  }
}

function normalizeQualificationCertificationPayload(payload?: QualificationCertificationPayload | null): QualificationCertificationPayload {
  const empty = createEmptyQualificationCertificationPayload()
  if (!payload) {
    return empty
  }
  return {
    qualityEnvironmentalHealth: normalizeQualificationItem(payload.qualityEnvironmentalHealth),
    highTechEnterprise: normalizeQualificationItem(payload.highTechEnterprise),
    specializedInnovation: normalizeQualificationItem(payload.specializedInnovation),
    listedCompany: normalizeQualificationItem(payload.listedCompany),
  }
}

function normalizeQualificationItem(item?: QualificationItem | null): QualificationItem {
  return {
    status: item?.status === '解锁' ? '解锁' : '未解锁',
    unlocked: item?.unlocked === true,
    locked: item?.locked === true,
    lockedByPrevious: item?.lockedByPrevious === true,
  }
}
export interface ReportManualPayload {
  workInProgress: number | null
  finishedGoods: number | null
  rawMaterials: number | null
  incomeTaxRate: number | null
  enterpriseCertificationScore: number | null
  productionHumanScore: number | null
  closingSpeedScore: number | null
}

export interface ReportComputedPayload {
  reportSalesRevenue: number
  reportDirectCost: number
  reportGrossProfit: number
  reportComprehensiveCost: number
  reportDepreciation: number
  reportOperatingProfit: number
  reportFinanceIncomeExpense: number
  reportExtraIncomeExpense: number
  reportPreTaxProfit: number
  reportIncomeTax: number
  reportNetProfit: number
  reportWorkInProgress: number
  reportFinishedGoods: number
  reportRawMaterials: number
  reportWorkInConstruction: number
  reportFactoryAsset: number
  reportLineResidual: number
  reportDepreciableAsset: number
  reportTotalNonCurrentAssets: number
  reportCash: number
  reportReceivable: number
  reportPostTaxCash: number
  reportTotalCurrentAssets: number
  reportTotalAssets: number
  reportShortTermLiability: number
  reportLongTermLiability: number
  reportTotalLiability: number
  reportShareCapital: number
  reportRetainedEarnings: number
  reportTotalEquity: number
  reportTotalLiabilityEquity: number
  reportBestMarketDirectorBaseScore: number
  reportBestMarketDirectorScore: number
  reportBestTechnologyDirectorScore: number
  reportBestSalesDirectorScore: number
  reportBestCfoBaseScore: number
  reportBestCfoScore: number
  reportBestCeoScore: number
}

export interface PlayerReportManualFieldOptions {
  incomeTaxRateOptions: number[]
}

export interface PlayerReportView {
  groupId: number
  yearNo: number
  yearStatus: string
  reportStatus: string
  businessStatus: string
  canView: boolean
  canEdit: boolean
  canSubmit: boolean
  hasInvalidDraft: boolean
  rollbackPending: boolean
  rollbackTargetYearNo?: number | null
  rollbackTargetStageCode?: string | null
  rollbackNotice: string
  reportComputedPayload: ReportComputedPayload
  reportManualPayload: ReportManualPayload
  manualFieldOptions: PlayerReportManualFieldOptions
  lastDraftSavedAt: string | null
  noticeBoard: PlayerNoticeBoard | null
  adjustmentRevision: number
}

export interface SavePlayerReportDraftRequest {
  yearNo: number
  reportManualPayload: ReportManualPayload
  clientSaveTime: string
}

export interface SavePlayerReportDraftResponse {
  groupId: number
  yearNo: number
  yearStatus: string
  reportStatus: string
  lastDraftSavedAt: string
}

export interface SubmitPlayerReportRequest {
  yearNo: number
  reportManualPayload: ReportManualPayload
}

export interface SubmitPlayerReportResponse {
  groupId: number
  yearNo: number
  yearStatus: string
  reportStatus: string
  businessStatus: string
  summaryEffective: boolean
  latestReportSubmitVersion: number
  balanceCheckPassed: boolean
  submittedAt: string
}

export function createEmptyReportManualPayload(): ReportManualPayload {
  return {
    workInProgress: null,
    finishedGoods: null,
    rawMaterials: null,
    incomeTaxRate: null,
    enterpriseCertificationScore: null,
    productionHumanScore: null,
    closingSpeedScore: null,
  }
}

export function cloneReportManualPayload(payload?: ReportManualPayload | null): ReportManualPayload {
  return {
    workInProgress: normalizeNullableNumber(payload?.workInProgress),
    finishedGoods: normalizeNullableNumber(payload?.finishedGoods),
    rawMaterials: normalizeNullableNumber(payload?.rawMaterials),
    incomeTaxRate: normalizeNullableNumber(payload?.incomeTaxRate),
    enterpriseCertificationScore: normalizeNullableNumber(payload?.enterpriseCertificationScore),
    productionHumanScore: normalizeNullableNumber(payload?.productionHumanScore),
    closingSpeedScore: normalizeNullableNumber(payload?.closingSpeedScore),
  }
}

export function createEmptyReportComputedPayload(): ReportComputedPayload {
  return {
    reportSalesRevenue: 0,
    reportDirectCost: 0,
    reportGrossProfit: 0,
    reportComprehensiveCost: 0,
    reportDepreciation: 0,
    reportOperatingProfit: 0,
    reportFinanceIncomeExpense: 0,
    reportExtraIncomeExpense: 0,
    reportPreTaxProfit: 0,
    reportIncomeTax: 0,
    reportNetProfit: 0,
    reportWorkInProgress: 0,
    reportFinishedGoods: 0,
    reportRawMaterials: 0,
    reportWorkInConstruction: 0,
    reportFactoryAsset: 0,
    reportLineResidual: 0,
    reportDepreciableAsset: 0,
    reportTotalNonCurrentAssets: 0,
    reportCash: 0,
    reportReceivable: 0,
    reportPostTaxCash: 0,
    reportTotalCurrentAssets: 0,
    reportTotalAssets: 0,
    reportShortTermLiability: 0,
    reportLongTermLiability: 0,
    reportTotalLiability: 0,
    reportShareCapital: 0,
    reportRetainedEarnings: 0,
    reportTotalEquity: 0,
    reportTotalLiabilityEquity: 0,
    reportBestMarketDirectorBaseScore: 0,
    reportBestMarketDirectorScore: 0,
    reportBestTechnologyDirectorScore: 0,
    reportBestSalesDirectorScore: 0,
    reportBestCfoBaseScore: 0,
    reportBestCfoScore: 0,
    reportBestCeoScore: 0,
  }
}

export function buildReportComputedPreview(
  basePayload?: ReportComputedPayload | null,
  manualPayload?: ReportManualPayload | null,
): ReportComputedPayload {
  const base = basePayload ? { ...basePayload } : createEmptyReportComputedPayload()
  const manual = cloneReportManualPayload(manualPayload)

  const workInProgress = manual.workInProgress ?? 0
  const finishedGoods = manual.finishedGoods ?? 0
  const rawMaterials = manual.rawMaterials ?? 0
  const incomeTaxRate = manual.incomeTaxRate ?? 0
  const enterpriseCertificationScore = manual.enterpriseCertificationScore ?? 0
  const productionHumanScore = manual.productionHumanScore ?? 0
  const closingSpeedScore = manual.closingSpeedScore ?? 0
  const reportIncomeTax = Math.max(Math.round(base.reportPreTaxProfit * incomeTaxRate), 0)
  const reportNetProfit = base.reportPreTaxProfit - reportIncomeTax
  const reportPostTaxCash = base.reportCash - reportIncomeTax
  const reportTotalCurrentAssets = base.reportReceivable + workInProgress + finishedGoods + rawMaterials + reportPostTaxCash
  const reportTotalAssets = base.reportTotalNonCurrentAssets + reportTotalCurrentAssets
  const reportTotalEquity = base.reportShareCapital + base.reportRetainedEarnings + reportNetProfit
  const reportTotalLiabilityEquity = base.reportTotalLiability + reportTotalEquity
  const reportBestMarketDirectorScore = base.reportBestMarketDirectorBaseScore + enterpriseCertificationScore
  const reportBestCfoScore = base.reportBestCfoBaseScore + closingSpeedScore
  // 销售总监得分已由后端按跨年累计规则完成缩放，财报手工预览只需复用该结果。
  const reportBestCeoScore =
    reportBestMarketDirectorScore +
    base.reportBestTechnologyDirectorScore +
    productionHumanScore +
    base.reportBestSalesDirectorScore +
    reportBestCfoScore +
    reportTotalEquity / 2

  return {
    ...base,
    reportIncomeTax,
    reportNetProfit,
    reportWorkInProgress: workInProgress,
    reportFinishedGoods: finishedGoods,
    reportRawMaterials: rawMaterials,
    reportPostTaxCash,
    reportTotalCurrentAssets,
    reportTotalAssets,
    reportTotalEquity,
    reportTotalLiabilityEquity,
    reportBestMarketDirectorScore,
    reportBestCfoScore,
    reportBestCeoScore,
  }
}

export function reportBalanceGap(payload?: ReportComputedPayload | null) {
  if (!payload) {
    return 0
  }
  return payload.reportTotalAssets - payload.reportTotalLiabilityEquity
}

function normalizeNullableNumber(value: number | null | undefined) {
  if (typeof value !== 'number' || Number.isNaN(value)) {
    return null
  }
  return value
}

function normalizeNumber(value: number | null | undefined) {
  if (typeof value !== 'number' || Number.isNaN(value)) {
    return 0
  }
  return value
}
