export interface AdminActionSummary {
  actionCode: string
  operatorName: string
  operateTime: string
  targetGroupId?: number | null
  targetYearNo?: number | null
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

export interface BaselinePayload {
  baselineSalesRevenue: number
  baselineDirectCost: number
  baselineComprehensiveCost: number
  baselineDepreciation: number
  baselineFinanceIncomeExpense: number
  baselineExtraIncomeExpense: number
  baselineIncomeTax: number
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

export function createEmptyBaselinePayload(): BaselinePayload {
  return {
    baselineSalesRevenue: 0,
    baselineDirectCost: 0,
    baselineComprehensiveCost: 0,
    baselineDepreciation: 0,
    baselineFinanceIncomeExpense: 0,
    baselineExtraIncomeExpense: 0,
    baselineIncomeTax: 0,
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
