export type NumericCellValue = number | ''
export type QuarterValueMap = Record<string, Record<string, NumericCellValue>>

export interface CurrentGameConfigResult {
  currentOpenYear: number
  finalYear: number
  ruleVersion: string
  templateVersion: string
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
    taxAndPlanning: Record<string, NumericCellValue>
    marketBid: Array<Record<string, NumericCellValue>>
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
  operatingPayload: OperatingPayload
  editableScopes: string[]
  readonlyScopes: string[]
  stageSubmitHistory: PlayerOperatingStageSubmitHistory[]
  lastDraftSavedAt: string | null
  quarterCashChecks: Record<string, number>
  derivedValues: Record<string, number>
  periodEndCash: number
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