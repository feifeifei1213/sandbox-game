import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { unlockAdminYear } from '@/api/sandbox-game/admin-control'
import {
  getAdminGroupOperatingView,
  getAdminGroupReportView,
  listAdminGroups,
} from '@/api/sandbox-game/admin-group-data'
import type { PageMessage } from '@/stores/admin-shell'
import type { PlayerOperatingView, PlayerReportView } from '@/types/sandbox-game'
import type {
  AdminGroupDataPageType,
  AdminGroupOption,
  UnlockStageCode,
  UnlockTargetType,
  UnlockYearResult,
} from '@/types/sandbox-game-admin'

const DEFAULT_OPERATING_STAGE: UnlockStageCode = 'YEAR_END'

export const useAdminGroupDataStore = defineStore('sandbox-admin-group-data', () => {
  const groups = ref<AdminGroupOption[]>([])
  const selectedGroupId = ref<number | null>(null)
  const selectedYear = ref(0)
  const selectedPageType = ref<AdminGroupDataPageType>('operating')
  const operatingView = ref<PlayerOperatingView | null>(null)
  const reportView = ref<PlayerReportView | null>(null)
  const loading = ref(false)
  const unlocking = ref(false)
  const pageMessage = ref<PageMessage | null>(null)
  const unlockDialogVisible = ref(false)
  const unlockReason = ref('')
  const unlockTargetType = ref<UnlockTargetType>('OPERATING')
  const unlockTargetStageCode = ref<UnlockStageCode | null>(DEFAULT_OPERATING_STAGE)
  const latestUnlockResult = ref<UnlockYearResult | null>(null)
  const latestUnlockReason = ref('')

  const selectedGroup = computed(() => groups.value.find((item) => item.groupId === selectedGroupId.value) ?? null)
  const activeView = computed(() => (selectedPageType.value === 'operating' ? operatingView.value : reportView.value))

  async function bootstrap(finalYear: number, currentOpenYear: number) {
    loading.value = true
    pageMessage.value = null
    try {
      const result = await listAdminGroups()
      groups.value = result.list
      normalizeSelection(finalYear, currentOpenYear)
      if (selectedGroupId.value !== null) {
        await loadCurrentView({ silent: true })
      } else {
        operatingView.value = null
        reportView.value = null
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '读取组数据入口失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  function normalizeSelection(finalYear: number, currentOpenYear: number) {
    const maxYear = Math.max(finalYear, 0)
    const hasCurrentGroup = groups.value.some((item) => item.groupId === selectedGroupId.value)
    if (!hasCurrentGroup) {
      selectedGroupId.value = groups.value[0]?.groupId ?? null
    }
    if (selectedYear.value > maxYear) {
      selectedYear.value = Math.min(Math.max(currentOpenYear, 0), maxYear)
    }
    if (selectedYear.value < 0) {
      selectedYear.value = 0
    }
  }

  async function loadCurrentView(options?: { silent?: boolean }) {
    if (selectedGroupId.value === null) {
      operatingView.value = null
      reportView.value = null
      return
    }

    loading.value = true
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      if (selectedPageType.value === 'operating') {
        operatingView.value = await getAdminGroupOperatingView(selectedGroupId.value, selectedYear.value)
        reportView.value = null
      } else {
        reportView.value = await getAdminGroupReportView(selectedGroupId.value, selectedYear.value)
        operatingView.value = null
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '读取组数据详情失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  function setSelectedGroupId(groupId: number) {
    selectedGroupId.value = groupId
  }

  function setSelectedYear(yearNo: number) {
    selectedYear.value = yearNo
  }

  function setSelectedPageType(pageType: AdminGroupDataPageType) {
    selectedPageType.value = pageType
  }

  function openUnlockDialog() {
    unlockReason.value = ''
    unlockTargetType.value = selectedPageType.value === 'report' ? 'REPORT' : 'OPERATING'
    unlockTargetStageCode.value = resolveDefaultUnlockStageCode()
    unlockDialogVisible.value = true
  }

  function closeUnlockDialog() {
    unlockDialogVisible.value = false
  }

  function setUnlockReason(reason: string) {
    unlockReason.value = reason
  }

  function setUnlockTargetType(targetType: UnlockTargetType) {
    unlockTargetType.value = targetType
    if (targetType === 'OPERATING' && !unlockTargetStageCode.value) {
      unlockTargetStageCode.value = resolveDefaultUnlockStageCode()
    }
    if (targetType === 'REPORT') {
      unlockTargetStageCode.value = null
    }
  }

  function setUnlockTargetStageCode(stageCode: UnlockStageCode) {
    unlockTargetStageCode.value = stageCode
  }

  async function submitUnlock() {
    if (selectedGroupId.value === null) {
      throw new Error('请先选择目标小组')
    }
    if (!unlockReason.value.trim()) {
      throw new Error('解锁原因不能为空')
    }
    if (unlockTargetType.value === 'OPERATING' && !unlockTargetStageCode.value) {
      throw new Error('请选择要回退的经营阶段')
    }

    unlocking.value = true
    try {
      const result = await unlockAdminYear({
        groupId: selectedGroupId.value,
        yearNo: selectedYear.value,
        unlockTargetType: unlockTargetType.value,
        targetStageCode: unlockTargetType.value === 'OPERATING' ? unlockTargetStageCode.value : null,
        reason: unlockReason.value.trim(),
      })
      latestUnlockResult.value = result
      latestUnlockReason.value = unlockReason.value.trim()
      unlockDialogVisible.value = false
      pageMessage.value = {
        type: 'success',
        text: `异常解锁已提交，第 ${selectedGroup.value?.groupNo ?? '--'} 组 ${selectedYear.value} 年已进入新的可编辑状态。`,
      }
      await loadCurrentView({ silent: true })
      return result
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '提交异常解锁失败')
      throw error
    } finally {
      unlocking.value = false
    }
  }

  function resolveDefaultUnlockStageCode(): UnlockStageCode {
    const currentStageCode = operatingView.value?.currentStageCode as UnlockStageCode | undefined
    if (currentStageCode && ['Q1', 'Q2', 'Q3', 'Q4', 'YEAR_END'].includes(currentStageCode)) {
      return currentStageCode
    }
    return DEFAULT_OPERATING_STAGE
  }

  return {
    groups,
    selectedGroupId,
    selectedYear,
    selectedPageType,
    operatingView,
    reportView,
    loading,
    unlocking,
    pageMessage,
    unlockDialogVisible,
    unlockReason,
    unlockTargetType,
    unlockTargetStageCode,
    latestUnlockResult,
    latestUnlockReason,
    selectedGroup,
    activeView,
    bootstrap,
    normalizeSelection,
    loadCurrentView,
    setSelectedGroupId,
    setSelectedYear,
    setSelectedPageType,
    openUnlockDialog,
    closeUnlockDialog,
    setUnlockReason,
    setUnlockTargetType,
    setUnlockTargetStageCode,
    submitUnlock,
  }
})

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}