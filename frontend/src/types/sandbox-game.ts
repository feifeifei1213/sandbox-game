export type NumericCellValue = number | string | ''
export type CellValue = number | string | boolean | null
export type QuarterValueMap = Record<string, Record<string, NumericCellValue>>

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
  demoYearEnabled: boolean
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
    receivableUpdate: QuarterValueMap
    deliverySettlement: QuarterValueMap
  }
  yearEnd: {
    longTermLoan: Record<string, NumericCellValue>
    assetAdjustment: Record<string, NumericCellValue>
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
      receivableUpdate: {},
      deliverySettlement: {},
    },
    yearEnd: {
      longTermLoan: {},
      assetAdjustment: {},
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
  return JSON.parse(JSON.stringify(payload)) as OperatingPayload
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
  const reportBestCeoScore =
    reportBestMarketDirectorScore +
    base.reportBestTechnologyDirectorScore +
    productionHumanScore +
    base.reportBestSalesDirectorScore +
    reportBestCfoScore

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
