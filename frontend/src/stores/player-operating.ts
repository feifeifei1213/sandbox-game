import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { getCurrentGameConfig, getYearTabs } from '@/api/sandbox-game/game-config'
import { formatStageCode } from '@/utils/sandbox-game-display'
import { buildOperatingPreviewCalculation } from '@/utils/sandbox-game-operating-preview'
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
    buildOperatingPreviewCalculation({
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