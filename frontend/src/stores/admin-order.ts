import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'

import {
  generateAdminOrderPool,
  getAdminOrderControlConfig,
  getAdminOrderPool,
  updateAdminOrderControlConfig,
  uploadAdminOrderExcel,
} from '@/api/sandbox-game/admin-order'
import type {
  AdminOrderType,
  OrderControlConfigItem,
  OrderControlConfigResult,
  OrderMarketCode,
  OrderPoolResult,
  UploadOrderExcelResult,
} from '@/types/sandbox-game-admin'

type MessageType = 'success' | 'error' | 'info'

export interface PageMessage {
  type: MessageType
  text: string
}

export const MARKET_OPTIONS: Array<{ code: OrderMarketCode; name: string }> = [
  { code: 'LOCAL', name: '本地市场' },
  { code: 'REGIONAL', name: '区域市场' },
  { code: 'NATIONAL', name: '全国市场' },
  { code: 'GLOBAL', name: '全球市场' },
]

export const ORDER_TYPE_OPTIONS: Array<{ code: AdminOrderType; name: string }> = [
  { code: 'AGENCY_INSPECTION', name: '代办过检' },
  { code: 'TWO_CABIN_VIP', name: '两舱贵宾' },
  { code: 'BUSINESS_VIP', name: '商务贵宾' },
  { code: 'MEMBER_CUSTOM', name: '会员定制' },
]

export const useAdminOrderStore = defineStore('sandbox-admin-order', () => {
  const selectedYearNo = ref(1)
  const config = ref<OrderControlConfigResult | null>(null)
  const uploadResult = ref<UploadOrderExcelResult | null>(null)
  const orderPool = ref<OrderPoolResult | null>(null)
  const loading = ref(false)
  const uploading = ref(false)
  const savingConfig = ref(false)
  const generatingPool = ref(false)
  const loadingPool = ref(false)
  const pageMessage = ref<PageMessage | null>(null)

  const poolFilter = reactive({
    marketCode: 'LOCAL' as OrderMarketCode,
    orderType: 'AGENCY_INSPECTION' as AdminOrderType,
  })

  const editableItems = ref<OrderControlConfigItem[]>([])

  const totalOrderCount = computed(() => editableItems.value.reduce((sum, item) => sum + Number(item.orderCount || 0), 0))
  const totalGeneratedCount = computed(() => editableItems.value.reduce((sum, item) => sum + Number(item.generatedCount || 0), 0))
  const hasLockedConfig = computed(() => editableItems.value.some((item) => item.configStatus === 'LOCKED'))
  const sortedItems = computed(() => [...editableItems.value].sort((a, b) => a.releaseSequenceNo - b.releaseSequenceNo))

  async function bootstrap(yearNo: number) {
    selectedYearNo.value = Math.max(yearNo, 1)
    await loadConfig()
  }

  async function loadConfig(options?: { silent?: boolean }) {
    loading.value = true
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      const result = await getAdminOrderControlConfig(selectedYearNo.value)
      applyConfig(result)
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '获取订单配置失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  async function uploadExcel(file: File) {
    uploading.value = true
    pageMessage.value = null
    try {
      const result = await uploadAdminOrderExcel(file)
      uploadResult.value = result
      await loadConfig({ silent: true })
      pageMessage.value = {
        type: 'success',
        text: `订单 Excel 已解析：有效订单 ${result.parsedOrderCount} 个。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '上传订单 Excel 失败')
      throw error
    } finally {
      uploading.value = false
    }
  }

  async function saveConfig() {
    savingConfig.value = true
    pageMessage.value = null
    try {
      const result = await updateAdminOrderControlConfig({
        yearNo: selectedYearNo.value,
        items: editableItems.value.map((item) => ({
          marketCode: item.marketCode,
          orderType: item.orderType,
          orderCount: Number(item.orderCount || 0),
          releaseSequenceNo: Number(item.releaseSequenceNo || 0),
        })),
      })
      applyConfig({
        yearNo: result.yearNo,
        finalYear: config.value?.finalYear ?? selectedYearNo.value,
        latestBatchId: config.value?.latestBatchId ?? null,
        latestBatchUploadedAt: config.value?.latestBatchUploadedAt ?? null,
        items: result.items,
        warnings: result.warnings,
      })
      pageMessage.value = {
        type: 'success',
        text: `订单配置已保存，操作人 ${result.updatedBy}。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '保存订单配置失败')
      throw error
    } finally {
      savingConfig.value = false
    }
  }

  async function generatePool(overwrite: boolean) {
    generatingPool.value = true
    pageMessage.value = null
    try {
      const result = await generateAdminOrderPool({
        yearNo: selectedYearNo.value,
        overwrite,
      })
      await loadConfig({ silent: true })
      await loadPool({ silent: true })
      pageMessage.value = {
        type: 'success',
        text: `${selectedYearNo.value} 年订单池已生成 ${result.generatedCount} 个订单，锁定 ${result.segmentCount} 个标段配置。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '生成订单池失败')
      throw error
    } finally {
      generatingPool.value = false
    }
  }

  async function loadPool(options?: { silent?: boolean }) {
    loadingPool.value = true
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      orderPool.value = await getAdminOrderPool(selectedYearNo.value, poolFilter.marketCode, poolFilter.orderType)
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '获取订单池失败')
      throw error
    } finally {
      loadingPool.value = false
    }
  }

  function applyConfig(result: OrderControlConfigResult) {
    config.value = result
    editableItems.value = result.items.map((item) => ({ ...item }))
  }

  function setYear(yearNo: number) {
    selectedYearNo.value = Math.max(yearNo, 1)
  }

  return {
    selectedYearNo,
    config,
    uploadResult,
    orderPool,
    loading,
    uploading,
    savingConfig,
    generatingPool,
    loadingPool,
    pageMessage,
    poolFilter,
    editableItems,
    sortedItems,
    totalOrderCount,
    totalGeneratedCount,
    hasLockedConfig,
    bootstrap,
    loadConfig,
    uploadExcel,
    saveConfig,
    generatePool,
    loadPool,
    setYear,
  }
})

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}
