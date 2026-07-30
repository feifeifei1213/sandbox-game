import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { getCurrentGameConfig, getYearTabs } from '@/api/sandbox-game/game-config'
import { getPlayerAdjustmentSync } from '@/api/sandbox-game/player-notice'
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
  type PlayerAdjustmentSyncResult,
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
  const adjustmentSyncMessage = ref('')
  let adjustmentSyncMessageTimer: ReturnType<typeof setTimeout> | null = null

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
  const operatingFeatureIssues = computed(() => collectOperatingFeatureIssues(draftPayload.value))

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
    if (operatingFeatureIssues.value.length > 0) {
      pageMessage.value = {
        type: 'error',
        text: `经营页新增记录填写不符合规则：${operatingFeatureIssues.value.slice(0, 3).join('、')}`,
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
    if (operatingFeatureIssues.value.length > 0) {
      pageMessage.value = {
        type: 'error',
        text: `经营页新增记录填写不符合规则：${operatingFeatureIssues.value.slice(0, 3).join('、')}`,
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

  async function syncAdjustments() {
    const view = currentView.value
    const yearNo = selectedYear.value
    if (!view || yearViewLoading.value) return
    try {
      const result = await getPlayerAdjustmentSync(yearNo, view.adjustmentRevision)
      if (result.notModified || selectedYear.value !== yearNo || currentView.value !== view) return
      applyAdjustmentSync(result)
      showAdjustmentSyncMessage(result.bankrupt ? '奖惩已同步，小组已进入破产状态' : '奖惩已更新，经营结果已同步')
    } catch {
      // 轮询失败保持静默，下一轮继续检查。
    }
  }

  function applyAdjustmentSync(result: PlayerAdjustmentSyncResult) {
    const view = currentView.value
    if (!view || !result.incomeAndPenalty) return
    draftPayload.value.extra.incomeAndPenalty = structuredClone(result.incomeAndPenalty)
    currentView.value = {
      ...view,
      adjustmentRevision: result.adjustmentRevision,
      businessStatus: result.businessStatus ?? view.businessStatus,
      noticeBoard: result.noticeBoard ?? view.noticeBoard,
      operatingPayload: {
        ...view.operatingPayload,
        extra: { ...view.operatingPayload.extra, incomeAndPenalty: structuredClone(result.incomeAndPenalty) },
      },
      derivedValues: result.derivedValues ?? view.derivedValues,
      quarterCashChecks: result.quarterCashChecks ?? view.quarterCashChecks,
      periodEndCash: result.periodEndCash ?? view.periodEndCash,
    }
  }

  function showAdjustmentSyncMessage(message: string) {
    adjustmentSyncMessage.value = message
    if (adjustmentSyncMessageTimer) clearTimeout(adjustmentSyncMessageTimer)
    adjustmentSyncMessageTimer = setTimeout(() => {
      adjustmentSyncMessage.value = ''
      adjustmentSyncMessageTimer = null
    }, 5000)
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
    adjustmentSyncMessage,
    reportEnabled,
    currentTab,
    previewCalculation,
    manualIntegerIssues,
    operatingFeatureIssues,
    bootstrap,
    loadYearView,
    updateDraft,
    saveDraft,
    submitCurrentStage,
    refreshTabs,
    clearMessage,
    syncAdjustments,
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
  collectManualIntegerIssues('项目进度', normalized.yearEnd.projectProgressUpdate.items, issues)
  collectManualIntegerIssues('其他收支', normalized.extra.incomeAndPenalty, issues)
  return issues
}

function collectOperatingFeatureIssues(payload: OperatingPayload) {
  const normalized = cloneOperatingPayload(payload)
  const issues: string[] = []
  const markets = [
    { label: '区域市场培育', item: normalized.yearEnd.marketCultivation.regional },
    { label: '全国市场培育', item: normalized.yearEnd.marketCultivation.national },
    { label: '全球市场培育', item: normalized.yearEnd.marketCultivation.global },
  ]
  for (const market of markets) {
    const raw = market.item.annualInvestment
    if (raw === '' || raw === undefined || raw === null) {
      continue
    }
    const parsed = Number(raw)
    if (!Number.isFinite(parsed) || !Number.isInteger(parsed) || (parsed !== 0 && parsed !== 1)) {
      issues.push(`${market.label}只能填0或1`)
      continue
    }
    if (market.item.lockedByPrevious && parsed !== 0) {
      issues.push(`${market.label}已解锁不能继续投入`)
    }
  }

  const qualifications = [
    { label: '质量、环境健康体系认证企业', item: normalized.yearEnd.qualificationCertification.qualityEnvironmentalHealth },
    { label: '高新技术企业', item: normalized.yearEnd.qualificationCertification.highTechEnterprise },
    { label: '专精特新小巨人', item: normalized.yearEnd.qualificationCertification.specializedInnovation },
    { label: '上市企业', item: normalized.yearEnd.qualificationCertification.listedCompany },
  ]
  for (const qualification of qualifications) {
    const status = qualification.item.status
    if (status !== '未解锁' && status !== '解锁') {
      issues.push(`${qualification.label}只能选择未解锁或解锁`)
      continue
    }
    if (qualification.item.lockedByPrevious && status !== '解锁') {
      issues.push(`${qualification.label}已继承解锁不能改回未解锁`)
    }
  }
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
