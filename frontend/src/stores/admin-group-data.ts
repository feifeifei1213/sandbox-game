import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

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
} from '@/types/sandbox-game-admin'

export const useAdminGroupDataStore = defineStore('sandbox-admin-group-data', () => {
  const groups = ref<AdminGroupOption[]>([])
  const selectedGroupId = ref<number | null>(null)
  const selectedYear = ref(0)
  const selectedPageType = ref<AdminGroupDataPageType>('operating')
  const operatingView = ref<PlayerOperatingView | null>(null)
  const reportView = ref<PlayerReportView | null>(null)
  const loading = ref(false)
  const pageMessage = ref<PageMessage | null>(null)

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

  return {
    groups,
    selectedGroupId,
    selectedYear,
    selectedPageType,
    operatingView,
    reportView,
    loading,
    pageMessage,
    selectedGroup,
    activeView,
    bootstrap,
    normalizeSelection,
    loadCurrentView,
    setSelectedGroupId,
    setSelectedYear,
    setSelectedPageType,
  }
})

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}
