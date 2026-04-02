import type { OperatingCarryForward, OperatingPayload } from '@/types/sandbox-game'

const quarterSequence = ['Q1', 'Q2', 'Q3', 'Q4'] as const
const marketProductKeys = ['basicProductTotal', 'standardProductTotal', 'precisionProductTotal', 'intelligentProductTotal'] as const

type QuarterCode = (typeof quarterSequence)[number]

type NumericRecord = Record<string, unknown>
type QuarterSource = Record<string, NumericRecord>

export interface OperatingPreviewCalculation {
  quarterCashChecks: Record<QuarterCode, number>
  derivedValues: Record<string, number>
  periodEndCash: number
}

interface BuildOperatingPreviewParams {
  payload?: OperatingPayload | null
  carryForward?: OperatingCarryForward | null
  currentStageCode?: string | null
  fallback?: Partial<OperatingPreviewCalculation> | null
}

export function buildOperatingPreviewCalculation(params: BuildOperatingPreviewParams): OperatingPreviewCalculation {
  const fallback = normalizeFallback(params.fallback)
  if (!params.payload || !params.carryForward) {
    return fallback
  }

  const payload = params.payload
  const carryForward = params.carryForward

  const marketBidCost = firstNonZero(
    toNumber(payload.beginning.taxAndPlanning.marketInvestmentTotal),
    toNumber(payload.beginning.taxAndPlanning.marketBidCost),
  )
  const orderTotal = computeOrderTotal(payload)
  const shortTermRepayment = sumQuarterField(payload.quarter.shortTermLoan, 'dueRepayment')
  const shortTermInterest = sumQuarterField(payload.quarter.shortTermLoan, 'interest')
  const newShortTermLoan = sumQuarterField(payload.quarter.shortTermLoan, 'newLoan')
  const materialPayment = sumQuarterValues(payload.quarter.materialPayment)
  const changeProductCost = sumQuarterField(payload.quarter.productionLineAdjustment, 'changeProduct')
  const lineDismantleCost = sumQuarterField(payload.quarter.productionLineAdjustment, 'dismantleCost')
  const lineSaleValue = sumQuarterField(payload.quarter.productionLineAdjustment, 'lineSale')
  const newLineInstall = sumQuarterField(payload.quarter.productionLineAdjustment, 'newLineInstall')
  const transferToFixed = sumQuarterField(payload.quarter.productionLineAdjustment, 'constructionToFixed')
  const depreciableAssetIncrease = sumQuarterField(payload.quarter.productionLineAdjustment, 'newDepreciableAsset')
  const humanResourceCost = sumQuarterField(payload.quarter.humanResource, 'staffCost')
  const salaryAndProductionCost = sumQuarterField(payload.quarter.salaryAndProduction, 'salaryCost')
  const researchCost = sumQuarterField(payload.quarter.researchAndManagement, 'technologyResearch')
  const managementSystemCost = sumQuarterField(payload.quarter.researchAndManagement, 'managementSystem')
  const receivableRecovered = sumQuarterField(payload.quarter.receivableUpdate, 'receivableCollection')
  const salesRevenue = sumQuarterField(payload.quarter.deliverySettlement, 'salesRevenue')
  const directCost = sumQuarterField(payload.quarter.deliverySettlement, 'directCost')
  const managementSalary = sumQuarterField(payload.quarter.deliverySettlement, 'managementStaffCost')
  const longTermInterest = toNumber(payload.yearEnd.longTermLoan.interest)
  const longTermRepayment = toNumber(payload.yearEnd.longTermLoan.repayment)
  const newLongTermLoan = toNumber(payload.yearEnd.longTermLoan.newLoan)
  const lineMaintenance = toNumber(payload.yearEnd.assetAdjustment.lineMaintenance)
  const factoryPurchase = toNumber(payload.yearEnd.assetAdjustment.purchase)
  const factorySale = toNumber(payload.yearEnd.assetAdjustment.sale)
  const factoryRent = toNumber(payload.yearEnd.assetAdjustment.rent)
  const workInConstruction = toNumber(payload.yearEnd.assetAdjustment.workInConstruction)
  const marketCultivation = toNumber(payload.yearEnd.assetAdjustment.marketCultivation)
  const discountExpense = sumQuarterField(payload.extra.incomeAndPenalty, 'discountExpense')
  const extraExpensePenalty = sumQuarterField(payload.extra.incomeAndPenalty, 'extraExpensePenalty')
  const extraIncomeReward = sumQuarterField(payload.extra.incomeAndPenalty, 'extraIncomeReward')

  const quarterCashChecks = {
    Q1: 0,
    Q2: 0,
    Q3: 0,
    Q4: 0,
  } satisfies Record<QuarterCode, number>

  let runningCash = carryForward.previousCash - carryForward.previousIncomeTax - marketBidCost
  quarterSequence.forEach((quarter) => {
    runningCash +=
      quarterValue(payload.quarter.shortTermLoan, quarter, 'newLoan') +
      quarterValue(payload.quarter.productionLineAdjustment, quarter, 'lineSale') +
      quarterValue(payload.quarter.receivableUpdate, quarter, 'receivableCollection') +
      quarterValue(payload.extra.incomeAndPenalty, quarter, 'extraIncomeReward')

    runningCash -=
      quarterValue(payload.quarter.shortTermLoan, quarter, 'dueRepayment') +
      quarterValue(payload.quarter.shortTermLoan, quarter, 'interest') +
      quarterRecordTotal(payload.quarter.materialPayment, quarter) +
      quarterValue(payload.quarter.productionLineAdjustment, quarter, 'changeProduct') +
      quarterValue(payload.quarter.productionLineAdjustment, quarter, 'dismantleCost') +
      quarterValue(payload.quarter.productionLineAdjustment, quarter, 'newLineInstall') +
      quarterValue(payload.quarter.humanResource, quarter, 'staffCost') +
      quarterValue(payload.quarter.salaryAndProduction, quarter, 'salaryCost') +
      quarterValue(payload.quarter.researchAndManagement, quarter, 'technologyResearch') +
      quarterValue(payload.quarter.researchAndManagement, quarter, 'managementSystem') +
      quarterValue(payload.quarter.deliverySettlement, quarter, 'managementStaffCost') +
      quarterValue(payload.extra.incomeAndPenalty, quarter, 'discountExpense') +
      quarterValue(payload.extra.incomeAndPenalty, quarter, 'extraExpensePenalty')

    quarterCashChecks[quarter] = runningCash
  })

  const comprehensiveCostTotal =
    marketBidCost +
    changeProductCost +
    lineDismantleCost +
    humanResourceCost +
    researchCost +
    managementSystemCost +
    managementSalary +
    lineMaintenance +
    factoryRent +
    marketCultivation

  const shortTermLoanBalance = carryForward.previousShortTermLoan + newShortTermLoan - shortTermRepayment
  const longTermLoanBalance = carryForward.previousLongTermLoan + newLongTermLoan - longTermRepayment
  const lineResidual = carryForward.previousEquipmentResidual + transferToFixed - lineSaleValue
  const depreciableAssetTotal = carryForward.previousDepreciableAsset + depreciableAssetIncrease
  const depreciation = Math.round(depreciableAssetTotal / 3)
  const financeIncomeExpense = shortTermInterest + longTermInterest + discountExpense
  const factoryAssetChange = factoryPurchase - factorySale
  const receivableChange = salesRevenue - receivableRecovered
  const extraIncomeExpense = extraIncomeReward - extraExpensePenalty
  const lineResidualChange = transferToFixed - lineSaleValue
  const depreciableAssetChange = depreciableAssetIncrease - depreciation
  const marketReturnRatio = safeDivide(orderTotal, marketBidCost)
  const researchIntensity = safeDivide(researchCost, salesRevenue + extraIncomeReward)
  const laborProductivity =
    safeDivide(
      salesRevenue -
        directCost +
        shortTermInterest -
        marketBidCost -
        carryForward.previousIncomeTax -
        changeProductCost -
        lineDismantleCost +
        humanResourceCost -
        researchCost -
        managementSystemCost -
        longTermInterest -
        lineMaintenance -
        factoryRent -
        marketCultivation -
        extraExpensePenalty +
        extraIncomeReward,
      salaryAndProductionCost * 7 + 5,
    ) * 100
  const yearEndCash =
    carryForward.previousCash +
    (newShortTermLoan + lineSaleValue + receivableRecovered + newLongTermLoan + factorySale + extraIncomeReward) -
    (carryForward.previousIncomeTax +
      marketBidCost +
      shortTermRepayment +
      shortTermInterest +
      materialPayment +
      changeProductCost +
      lineDismantleCost +
      newLineInstall +
      humanResourceCost +
      salaryAndProductionCost +
      researchCost +
      managementSystemCost +
      managementSalary +
      lineMaintenance +
      factoryPurchase +
      longTermInterest +
      longTermRepayment +
      factoryRent +
      marketCultivation +
      discountExpense +
      extraExpensePenalty)

  const derivedValues: Record<string, number> = {
    marketBidCost,
    orderTotal,
    comprehensiveCostTotal,
    shortTermRepayment,
    shortTermInterest,
    newShortTermLoan,
    shortTermLoanBalance,
    materialPayment,
    changeProductCost,
    lineDismantleCost,
    lineSaleValue,
    newLineInstall,
    transferToFixed,
    depreciableAssetIncrease,
    humanResourceCost,
    salaryAndProductionCost,
    researchCost,
    managementSystemCost,
    receivableRecovered,
    salesRevenue,
    directCost,
    managementSalary,
    longTermInterest,
    longTermRepayment,
    newLongTermLoan,
    longTermLoanBalance,
    lineMaintenance,
    factoryPurchase,
    factorySale,
    factoryAssetChange,
    factoryRent,
    workInConstruction,
    marketCultivation,
    marketReturnRatio,
    researchIntensity,
    laborProductivity,
    lineResidual,
    depreciableAssetTotal,
    depreciation,
    financeIncomeExpense,
    receivableChange,
    extraIncomeExpense,
    lineResidualChange,
    depreciableAssetChange,
    discountExpense,
    extraExpensePenalty,
    extraIncomeReward,
    periodEndCash: yearEndCash,
    q1QuarterEndCashCheck: quarterCashChecks.Q1,
    q2QuarterEndCashCheck: quarterCashChecks.Q2,
    q3QuarterEndCashCheck: quarterCashChecks.Q3,
    q4QuarterEndCashCheck: quarterCashChecks.Q4,
  }

  return {
    quarterCashChecks,
    derivedValues: {
      ...fallback.derivedValues,
      ...derivedValues,
    },
    periodEndCash: selectPeriodEndCash(params.currentStageCode, quarterCashChecks, yearEndCash),
  }
}

function normalizeFallback(fallback?: Partial<OperatingPreviewCalculation> | null): OperatingPreviewCalculation {
  return {
    quarterCashChecks: {
      Q1: fallback?.quarterCashChecks?.Q1 ?? 0,
      Q2: fallback?.quarterCashChecks?.Q2 ?? 0,
      Q3: fallback?.quarterCashChecks?.Q3 ?? 0,
      Q4: fallback?.quarterCashChecks?.Q4 ?? 0,
    },
    derivedValues: { ...(fallback?.derivedValues ?? {}) },
    periodEndCash: fallback?.periodEndCash ?? 0,
  }
}

function selectPeriodEndCash(currentStageCode: string | null | undefined, quarterCashChecks: Record<QuarterCode, number>, yearEndCash: number) {
  switch ((currentStageCode ?? '').toUpperCase()) {
    case 'Q1':
      return quarterCashChecks.Q1
    case 'Q2':
      return quarterCashChecks.Q2
    case 'Q3':
      return quarterCashChecks.Q3
    case 'Q4':
      return quarterCashChecks.Q4
    default:
      return yearEndCash
  }
}

function computeOrderTotal(payload: OperatingPayload): number {
  return payload.beginning.marketBid.reduce((sum: number, row) => {
    const typedRow = row as NumericRecord
    const rowTotal = marketProductKeys.reduce((total: number, key) => total + toNumber(typedRow[key]), 0)
    return sum + rowTotal
  }, 0)
}

function sumQuarterField(source: QuarterSource, fieldKey: string): number {
  return quarterSequence.reduce((sum: number, quarter) => sum + quarterValue(source, quarter, fieldKey), 0)
}

function sumQuarterValues(source: QuarterSource): number {
  return quarterSequence.reduce((sum: number, quarter) => sum + quarterRecordTotal(source, quarter), 0)
}

function quarterValue(source: QuarterSource, quarter: QuarterCode, fieldKey: string): number {
  const record = getQuarterRecord(source, quarter)
  return toNumber(record?.[fieldKey])
}

function quarterRecordTotal(source: QuarterSource, quarter: QuarterCode): number {
  const record = getQuarterRecord(source, quarter)
  if (!record) {
    return 0
  }
  return Object.values(record as NumericRecord).reduce((sum: number, value) => sum + toNumber(value), 0)
}

function getQuarterRecord(source: QuarterSource, quarter: QuarterCode): NumericRecord | undefined {
  const target = quarter.toLowerCase()
  return Object.entries(source).find(([key]) => key.toLowerCase() === target)?.[1]
}

function firstNonZero(...values: number[]) {
  return values.find((value) => value !== 0) ?? 0
}

function safeDivide(numerator: number, denominator: number) {
  if (denominator === 0) {
    return 0
  }
  return numerator / denominator
}

function toNumber(value: unknown) {
  if (typeof value === 'number' && !Number.isNaN(value)) {
    return value
  }
  if (typeof value === 'string') {
    const trimmed = value.trim()
    if (trimmed === '') {
      return 0
    }
    const parsed = Number(trimmed)
    return Number.isNaN(parsed) ? 0 : parsed
  }
  return 0
}
