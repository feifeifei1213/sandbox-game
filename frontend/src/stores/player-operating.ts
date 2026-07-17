import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { getCurrentGameConfig, getYearTabs } from '@/api/sandbox-game/game-config'
import { formatStageCode } from '@/utils/sandbox-game-display'
import { buildOperatingPreviewCalculation } from '@/utils/sandbox-game-operating-preview'
import { hasFractionInput } from '@/utils/manual-integer'
import {
  getPlayerOperatingYearView,
  savePlayerOperatingDraft,
  submitPlayerOperatingStage,
} from '@/api/sandbox-game/player-operating'
import {
  cloneOperatingPayload,
  createEmptyOperatingPayload,
  type CurrentGameConfigResult,
  type OperatingPayload,
  type PlayerOperatingView,
  type YearTabItem,
  type YearTabsResult,
} from '@/types/sandbox-game'

type MessageType = 'success' | 'error' | 'info'

export interface PageMessage {
  type: MessageType
  text: string
}

export const usePlayerOperatingStore = defineStore('sandbox-player-operating', () => {
  const currentConfig = ref<CurrentGameConfigResult | null>(null)
  const yearTabs = ref<YearTabItem[]>([])
  const currentView = ref<PlayerOperatingView | null>(null)
  const draftPayload = ref<OperatingPayload>(createEmptyOperatingPayload())
  const selectedYear = ref(0)
  const loading = ref(false)
  const yearViewLoading = ref(false)
  const saving = ref(false)
  const submitting = ref(false)
  const dirty = ref(false)
  const pageMessage = ref<PageMessage | null>(null)

  const reportEnabled = computed(() => {
    const view = currentView.value
    if (!view) {
      return false
    }
    return view.reportStatus !== 'REPORT_LOCKED' || ['REPORT_PENDING', 'REPORTING', 'COMPLETED'].includes(view.yearStatus)
  })

  const currentTab = computed(() => yearTabs.value.find((item) => item.yearNo === selectedYear.value) ?? null)
  const previewCalculation = computed(() =>
    currentView.value?.rollbackPending && !currentView.value.carryForward
      ? {
          quarterCashChecks: currentView.value.quarterCashChecks,
          derivedValues: currentView.value.derivedValues,
          periodEndCash: currentView.value.periodEndCash,
        }
      : buildOperatingPreviewCalculation({
          payload: draftPayload.value,
          carryForward: currentView.value?.carryForward,
          currentStageCode: currentView.value?.currentStageCode,
          fallback: currentView.value
            ? {
                quarterCashChecks: currentView.value.quarterCashChecks,
                derivedValues: currentView.value.derivedValues,
                periodEndCash: currentView.value.periodEndCash,
              }
            : null,
        }),
  )
  const manualIntegerIssues = computed(() => collectOperatingIntegerIssues(draftPayload.value))
  const supplyChainOrderQuantityIssues = computed(() => collectSupplyChainOrderQuantityIssues(draftPayload.value))

  async function bootstrap(preferredYear?: number) {
    loading.value = true
    pageMessage.value = null
    try {
      const [config, tabsResult] = await Promise.all([getCurrentGameConfig(), getYearTabs()])
      currentConfig.value = config
      yearTabs.value = tabsResult.tabs
      const nextYear = resolveInitialYear(tabsResult, preferredYear)
      await loadYearView(nextYear)
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '初始化玩家经营页失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  async function refreshTabs() {
    const tabsResult = await getYearTabs()
    yearTabs.value = tabsResult.tabs
    return tabsResult
  }

  async function loadYearView(yearNo: number) {
    yearViewLoading.value = true
    pageMessage.value = null
    try {
      const view = await getPlayerOperatingYearView(yearNo)
      currentView.value = view
      draftPayload.value = cloneOperatingPayload(view.operatingPayload)
      selectedYear.value = yearNo
      dirty.value = false
    } catch (error) {
      pageMessage.value = toErrorMessage(error, `读取 ${yearNo} 年经营页失败`)
      throw error
    } finally {
      yearViewLoading.value = false
    }
  }

  function updateDraft(nextPayload: OperatingPayload) {
    draftPayload.value = cloneOperatingPayload(nextPayload)
    dirty.value = true
  }

  async function saveDraft(options?: { silent?: boolean }) {
    const view = currentView.value
    if (!view) {
      return
    }
    if (manualIntegerIssues.value.length > 0) {
      pageMessage.value = {
        type: 'error',
        text: buildIntegerIssueMessage('经营页手工数字', manualIntegerIssues.value),
      }
      return
    }
    if (supplyChainOrderQuantityIssues.value.length > 0) {
      pageMessage.value = {
        type: 'error',
        text: `订单数量必须为非负整数：${supplyChainOrderQuantityIssues.value.slice(0, 3).join('、')}`,
      }
      return
    }
    saving.value = true
    try {
      const result = await savePlayerOperatingDraft({
        yearNo: selectedYear.value,
        stageStatus: view.stageStatus,
        operatingPayload: draftPayload.value,
        clientSaveTime: new Date().toISOString(),
      })
      currentView.value = {
        ...view,
        lastDraftSavedAt: result.lastDraftSavedAt,
      }
      dirty.value = false
      if (!options?.silent) {
        pageMessage.value = {
          type: 'success',
          text: `草稿已保存，时间 ${formatDateTime(result.lastDraftSavedAt)}`,
        }
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '保存草稿失败')
      throw error
    } finally {
      saving.value = false
    }
  }

  async function submitCurrentStage() {
    const view = currentView.value
    if (!view) {
      return
    }
    if (manualIntegerIssues.value.length > 0) {
      pageMessage.value = {
        type: 'error',
        text: buildIntegerIssueMessage('经营页手工数字', manualIntegerIssues.value),
      }
      return
    }
    if (supplyChainOrderQuantityIssues.value.length > 0) {
      pageMessage.value = {
        type: 'error',
        text: `订单数量必须为非负整数：${supplyChainOrderQuantityIssues.value.slice(0, 3).join('、')}`,
      }
      return
    }
    submitting.value = true
    try {
      const result = await submitPlayerOperatingStage({
        yearNo: selectedYear.value,
        stageCode: view.currentStageCode,
        operatingPayload: draftPayload.value,
      })
      pageMessage.value = {
        type: 'success',
        text: `${result.stageCode} 提交成功，期末现金 ${result.periodEndCash.toLocaleString('zh-CN')}。请关注贷款更新。`,
      }
      await Promise.all([refreshTabs(), loadYearView(selectedYear.value)])
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '阶段提交失败')
      throw error
    } finally {
      submitting.value = false
    }
  }

  function clearMessage() {
    pageMessage.value = null
  }

  return {
    currentConfig,
    yearTabs,
    currentView,
    draftPayload,
    selectedYear,
    loading,
    yearViewLoading,
    saving,
    submitting,
    dirty,
    pageMessage,
    reportEnabled,
    currentTab,
    previewCalculation,
    manualIntegerIssues,
    bootstrap,
    loadYearView,
    updateDraft,
    saveDraft,
    submitCurrentStage,
    refreshTabs,
    clearMessage,
  }
})

function resolveInitialYear(result: YearTabsResult, preferredYear?: number) {
  if (typeof preferredYear === 'number' && Number.isFinite(preferredYear)) {
    const matched = result.tabs.find((item) => item.yearNo === preferredYear)
    if (matched?.canEnter) {
      return preferredYear
    }
  }

  const currentOpen = result.tabs.find((item) => item.isCurrentOpenYear)
  if (currentOpen?.canEnter) {
    return currentOpen.yearNo
  }

  const firstEnterable = result.tabs.find((item) => item.canEnter)
  if (firstEnterable) {
    return firstEnterable.yearNo
  }

  return 0
}

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}

function formatDateTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString('zh-CN', { hour12: false })
}

function collectOperatingIntegerIssues(payload: OperatingPayload) {
  const issues: string[] = []
  const normalized = cloneOperatingPayload(payload)
  collectManualIntegerIssues('年初规划', normalized.beginning.taxAndPlanning, issues)
  collectManualIntegerIssues('年初市场竞标', normalized.beginning.marketBid, issues)
  collectManualIntegerIssues('短期贷款', normalized.quarter.shortTermLoan, issues)
  collectManualIntegerIssues('材料费', normalized.quarter.materialPayment, issues)
  collectManualIntegerIssues('生产线调整', normalized.quarter.productionLineAdjustment, issues)
  collectManualIntegerIssues('人力资源', normalized.quarter.humanResource, issues)
  collectManualIntegerIssues('工资与生产', normalized.quarter.salaryAndProduction, issues)
  collectManualIntegerIssues('研发与管理', normalized.quarter.researchAndManagement, issues)
  collectManualIntegerIssues('订单数量留痕', normalized.quarter.supplyChainOrderRecord, issues)
  collectManualIntegerIssues('应收更新', normalized.quarter.receivableUpdate, issues)
  collectManualIntegerIssues('交货结算', normalized.quarter.deliverySettlement, issues)
  collectManualIntegerIssues('长期贷款', normalized.yearEnd.longTermLoan, issues)
  collectManualIntegerIssues('资产调整', normalized.yearEnd.assetAdjustment, issues)
  collectManualIntegerIssues('其他收支', normalized.extra.incomeAndPenalty, issues)
  return issues
}

function collectSupplyChainOrderQuantityIssues(payload: OperatingPayload) {
  const issues: string[] = []
  const fieldLabels: Record<string, string> = {
    basicProduct: '第一类',
    standardProduct: '第二类',
    precisionProduct: '第三类',
    intelligentProduct: '第四类',
  }
  for (const [quarterKey, values] of Object.entries(payload.quarter.supplyChainOrderRecord ?? {})) {
    for (const [fieldKey, label] of Object.entries(fieldLabels)) {
      const raw = values?.[fieldKey]
      if (raw === '' || raw === undefined || raw === null) {
        continue
      }
      const parsed = Number(raw)
      if (!Number.isFinite(parsed) || !Number.isInteger(parsed) || parsed < 0) {
        issues.push(`${quarterKey.toUpperCase()}.${label}`)
      }
    }
  }
  return issues
}

function collectManualIntegerIssues(label: string, value: unknown, issues: string[]) {
  if (value === null || value === undefined || value === '') {
    return
  }
  if (typeof value === 'number') {
    if (hasFractionInput(value)) {
      issues.push(label)
    }
    return
  }
  if (typeof value === 'string') {
    const trimmed = value.trim()
    if (trimmed === '' || !Number.isFinite(Number(trimmed))) {
      return
    }
    if (hasFractionInput(value)) {
      issues.push(label)
    }
    return
  }
  if (Array.isArray(value)) {
    value.forEach((item, index) => collectManualIntegerIssues(`${label}${index + 1}`, item, issues))
    return
  }
  if (typeof value === 'object') {
    for (const [key, child] of Object.entries(value as Record<string, unknown>)) {
      if (key === 'orderLinked') {
        continue
      }
      collectManualIntegerIssues(`${label}.${key}`, child, issues)
    }
  }
}

function buildIntegerIssueMessage(scope: string, issues: string[]) {
  const uniqueIssues = Array.from(new Set(issues))
  const preview = uniqueIssues.slice(0, 3).join('、')
  const suffix = uniqueIssues.length > 3 ? `等 ${uniqueIssues.length} 项` : ''
  return `${scope}必须填写整数，发现小数：${preview}${suffix}`
}
