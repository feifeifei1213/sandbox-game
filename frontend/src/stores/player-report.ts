import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { getCurrentGameConfig, getYearTabs } from '@/api/sandbox-game/game-config'
import { getPlayerReportView, savePlayerReportDraft, submitPlayerReport } from '@/api/sandbox-game/player-report'
import { resolveReportRequiredFieldLabels } from '@/configs/sandbox-game-service-labels'
import {
  buildReportComputedPreview,
  cloneReportManualPayload,
  createEmptyReportManualPayload,
  reportBalanceGap,
  type CurrentGameConfigResult,
  type PlayerReportView,
  type ReportManualPayload,
  type YearTabItem,
  type YearTabsResult,
} from '@/types/sandbox-game'

type MessageType = 'success' | 'error' | 'info'

const balanceTolerance = 0.000001

export interface PageMessage {
  type: MessageType
  text: string
}

export const usePlayerReportStore = defineStore('sandbox-player-report', () => {
  const currentConfig = ref<CurrentGameConfigResult | null>(null)
  const yearTabs = ref<YearTabItem[]>([])
  const currentView = ref<PlayerReportView | null>(null)
  const draftManualPayload = ref<ReportManualPayload>(createEmptyReportManualPayload())
  const selectedYear = ref(0)
  const loading = ref(false)
  const yearViewLoading = ref(false)
  const saving = ref(false)
  const submitting = ref(false)
  const dirty = ref(false)
  const pageMessage = ref<PageMessage | null>(null)

  const currentTab = computed(() => yearTabs.value.find((item) => item.yearNo === selectedYear.value) ?? null)
  const previewComputedPayload = computed(() =>
    buildReportComputedPreview(currentView.value?.reportComputedPayload, draftManualPayload.value),
  )
  const balanceGap = computed(() => reportBalanceGap(previewComputedPayload.value))
  const balancePassed = computed(() => Math.abs(balanceGap.value) <= balanceTolerance)
  const missingFields = computed(() => {
    const missing: string[] = []
    const requiredFieldLabels = resolveReportRequiredFieldLabels(currentConfig.value?.editionCode)
    if (draftManualPayload.value.workInProgress === null) {
      missing.push(requiredFieldLabels.workInProgress)
    }
    if (draftManualPayload.value.finishedGoods === null) {
      missing.push(requiredFieldLabels.finishedGoods)
    }
    if (draftManualPayload.value.rawMaterials === null) {
      missing.push(requiredFieldLabels.rawMaterials)
    }
    if (draftManualPayload.value.incomeTaxRate === null) {
      missing.push(requiredFieldLabels.incomeTaxRate)
    }
    if (draftManualPayload.value.enterpriseCertificationScore === null) {
      missing.push(requiredFieldLabels.enterpriseCertificationScore)
    }
    if (draftManualPayload.value.productionHumanScore === null) {
      missing.push(requiredFieldLabels.productionHumanScore)
    }
    if (draftManualPayload.value.closingSpeedScore === null) {
      missing.push(requiredFieldLabels.closingSpeedScore)
    }
    return missing
  })
  const submitReady = computed(() => Boolean(currentView.value?.canSubmit) && missingFields.value.length === 0 && balancePassed.value)

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
      pageMessage.value = toErrorMessage(error, '初始化玩家财报页失败')
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
      const view = await getPlayerReportView(yearNo)
      currentView.value = view
      draftManualPayload.value = cloneReportManualPayload(view.reportManualPayload)
      selectedYear.value = yearNo
      dirty.value = false
    } catch (error) {
      pageMessage.value = toErrorMessage(error, `读取 ${yearNo} 年财报页失败`)
      throw error
    } finally {
      yearViewLoading.value = false
    }
  }

  function updateDraft(nextPayload: ReportManualPayload) {
    draftManualPayload.value = cloneReportManualPayload(nextPayload)
    dirty.value = true
  }

  async function saveDraft(options?: { silent?: boolean }) {
    const view = currentView.value
    if (!view) {
      return
    }
    saving.value = true
    try {
      const result = await savePlayerReportDraft({
        yearNo: selectedYear.value,
        reportManualPayload: draftManualPayload.value,
        clientSaveTime: new Date().toISOString(),
      })
      currentView.value = {
        ...view,
        yearStatus: result.yearStatus,
        reportStatus: result.reportStatus,
        reportManualPayload: cloneReportManualPayload(draftManualPayload.value),
        lastDraftSavedAt: result.lastDraftSavedAt,
      }
      dirty.value = false
      if (!options?.silent) {
        pageMessage.value = {
          type: 'success',
          text: `财报草稿已保存，时间 ${formatDateTime(result.lastDraftSavedAt)}`,
        }
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '保存财报草稿失败')
      throw error
    } finally {
      saving.value = false
    }
  }

  async function submitCurrentReport() {
    const view = currentView.value
    if (!view) {
      return
    }
    submitting.value = true
    try {
      const result = await submitPlayerReport({
        yearNo: selectedYear.value,
        reportManualPayload: draftManualPayload.value,
      })
      const resultText = result.businessStatus === 'BANKRUPT' ? '。该组已进入破产状态。' : '。'
      pageMessage.value = {
        type: 'success',
        text: `${selectedYear.value} 年财报提交成功${resultText}`,
      }
      await Promise.all([refreshTabs(), loadYearView(selectedYear.value)])
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '财报提交失败')
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
    draftManualPayload,
    selectedYear,
    loading,
    yearViewLoading,
    saving,
    submitting,
    dirty,
    pageMessage,
    currentTab,
    previewComputedPayload,
    balanceGap,
    balancePassed,
    missingFields,
    submitReady,
    bootstrap,
    loadYearView,
    updateDraft,
    saveDraft,
    submitCurrentReport,
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
