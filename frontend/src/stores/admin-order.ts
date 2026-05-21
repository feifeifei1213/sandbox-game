import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'

import {
  adminSkipCurrentOrderGroup,
  closeAdminMarketBidding,
  confirmAdminOrderPool,
  generateAdminSelectionSequence,
  generateAdminOrderPool,
  getAdminMarketSelectionStatus,
  getAdminOrderForecastControl,
  getAdminOrderControlConfig,
  getAdminOrderPool,
  openAdminMarketBidding,
  releaseNextAdminOrderSegment,
  updateAdminOrderControlConfig,
  updateAdminOrderForecastControl,
  updateAdminOrderMarketConfig,
  uploadAdminOrderExcel,
} from '@/api/sandbox-game/admin-order'
import type {
  AdminOrderType,
  OrderControlConfigItem,
  OrderControlConfigResult,
  OrderForecastControlItem,
  OrderForecastControlResult,
  OrderForecastNarrativeItem,
  OrderMarketConfigItem,
  OrderMarketCode,
  OrderPoolResult,
  UploadOrderExcelResult,
} from '@/types/sandbox-game-admin'
import type { AdminMarketSelectionStatus } from '@/types/sandbox-game-order'

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
  const forecastControl = ref<OrderForecastControlResult | null>(null)
  const uploadResult = ref<UploadOrderExcelResult | null>(null)
  const orderPool = ref<OrderPoolResult | null>(null)
  const marketSelectionStatus = ref<AdminMarketSelectionStatus | null>(null)
  const loading = ref(false)
  const uploading = ref(false)
  const savingConfig = ref(false)
  const savingForecastControl = ref(false)
  const generatingPool = ref(false)
  const confirmingPool = ref(false)
  const generatingSequence = ref(false)
  const loadingPool = ref(false)
  const savingMarketConfig = ref(false)
  const loadingSelectionStatus = ref(false)
  const controllingMarket = ref(false)
  const releasingSegment = ref(false)
  const skippingGroup = ref(false)
  const pageMessage = ref<PageMessage | null>(null)

  const poolFilter = reactive({
    marketCode: 'ALL' as OrderMarketCode | 'ALL',
    orderType: 'ALL' as AdminOrderType | 'ALL',
  })
  const controlForm = reactive({
    marketCode: 'LOCAL' as OrderMarketCode,
    skipReason: '',
  })

  const editableItems = ref<OrderControlConfigItem[]>([])
  const editableForecastItems = ref<OrderForecastControlItem[]>([])
  const editableForecastNarratives = ref<OrderForecastNarrativeItem[]>([])
  const editableMarketConfigs = ref<OrderMarketConfigItem[]>([])

  const totalOrderCount = computed(() => editableItems.value.reduce((sum, item) => sum + Number(item.orderCount || 0), 0))
  const totalGeneratedCount = computed(() => editableItems.value.reduce((sum, item) => sum + Number(item.generatedCount || 0), 0))
  const hasLockedConfig = computed(() => editableItems.value.some((item) => item.configStatus === 'LOCKED'))
  const sortedItems = computed(() => [...editableItems.value].sort((a, b) => a.releaseSequenceNo - b.releaseSequenceNo))
  const forecastStages = computed(() => forecastControl.value?.forecast.stages ?? config.value?.forecast.stages ?? [])
  const enabledMarketCount = computed(() => editableMarketConfigs.value.filter((item) => item.enabled).length)
  const currentSegment = computed(() => marketSelectionStatus.value?.currentSegment ?? null)

  async function bootstrap(yearNo: number) {
    selectedYearNo.value = Math.max(yearNo, 1)
    await Promise.all([loadForecastControl({ silent: true }), loadConfig()])
  }

  async function loadForecastControl(options?: { silent?: boolean }) {
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      const result = await getAdminOrderForecastControl()
      applyForecastControl(result)
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '获取订单数量控制台失败')
      throw error
    }
  }

  async function loadConfig(options?: { silent?: boolean }) {
    loading.value = true
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      const result = await getAdminOrderControlConfig(selectedYearNo.value)
      applyConfig(result)
      await loadSelectionStatus({ silent: true })
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
          releaseSequenceNo: Number(item.releaseSequenceNo || 0),
        })),
      })
      applyConfig({
        yearNo: result.yearNo,
        finalYear: config.value?.finalYear ?? selectedYearNo.value,
        latestBatchId: config.value?.latestBatchId ?? null,
        latestBatchUploadedAt: config.value?.latestBatchUploadedAt ?? null,
        forecast: forecastControl.value?.forecast ?? config.value?.forecast ?? { formulaVersion: '', stages: [] },
        generationStatus: config.value?.generationStatus ?? 'NOT_GENERATED',
        latestPreviewBatch: config.value?.latestPreviewBatch ?? null,
        confirmedBatch: config.value?.confirmedBatch ?? null,
        canUpdateConfig: config.value?.canUpdateConfig ?? true,
        canGeneratePreview: config.value?.canGeneratePreview ?? true,
        canConfirmPool: config.value?.canConfirmPool ?? false,
        marketConfigs: config.value?.marketConfigs ?? editableMarketConfigs.value,
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

  async function saveForecastControl() {
    savingForecastControl.value = true
    pageMessage.value = null
    try {
      const result = await updateAdminOrderForecastControl({
        items: editableForecastItems.value.map((item) => ({
          yearNo: item.yearNo,
          marketCode: item.marketCode,
          orderType: item.orderType,
          orderCount: Number(item.orderCount || 0),
        })),
        narratives: editableForecastNarratives.value.map((item) => ({
          forecastStageCode: item.forecastStageCode,
          marketCode: item.marketCode,
          content: item.content,
        })),
      })
      applyForecastControl(result)
      await loadConfig({ silent: true })
      pageMessage.value = {
        type: 'success',
        text: `多年订单数量控制台已保存，操作人 ${result.updatedBy}。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '保存多年订单数量控制台失败')
      throw error
    } finally {
      savingForecastControl.value = false
    }
  }

  async function saveMarketConfig() {
    savingMarketConfig.value = true
    pageMessage.value = null
    try {
      const result = await updateAdminOrderMarketConfig({
        yearNo: selectedYearNo.value,
        markets: editableMarketConfigs.value.map((item) => ({
          marketCode: item.marketCode,
          enabled: Boolean(item.enabled),
        })),
      })
      applyConfig({
        yearNo: result.yearNo,
        finalYear: config.value?.finalYear ?? selectedYearNo.value,
        latestBatchId: config.value?.latestBatchId ?? null,
        latestBatchUploadedAt: config.value?.latestBatchUploadedAt ?? null,
        forecast: forecastControl.value?.forecast ?? config.value?.forecast ?? { formulaVersion: '', stages: [] },
        generationStatus: config.value?.generationStatus ?? 'NOT_GENERATED',
        latestPreviewBatch: config.value?.latestPreviewBatch ?? null,
        confirmedBatch: config.value?.confirmedBatch ?? null,
        canUpdateConfig: config.value?.canUpdateConfig ?? true,
        canGeneratePreview: config.value?.canGeneratePreview ?? true,
        canConfirmPool: config.value?.canConfirmPool ?? false,
        marketConfigs: result.markets,
        items: result.items,
        warnings: result.warnings,
      })
      pageMessage.value = {
        type: 'success',
        text: `市场开启配置已保存，操作人 ${result.updatedBy}。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '保存市场开启配置失败')
      throw error
    } finally {
      savingMarketConfig.value = false
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
      await loadSelectionStatus({ silent: true })
      pageMessage.value = {
        type: 'success',
        text: `${selectedYearNo.value} 年预览订单池已生成 ${result.generatedCount} 个订单，批次 #${result.batchId}。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '生成订单池失败')
      throw error
    } finally {
      generatingPool.value = false
    }
  }

  async function confirmPool() {
    const batchId = config.value?.latestPreviewBatch?.batchId
    if (!batchId) {
      pageMessage.value = {
        type: 'error',
        text: '请先生成预览订单池。',
      }
      return
    }
    confirmingPool.value = true
    pageMessage.value = null
    try {
      const result = await confirmAdminOrderPool({
        yearNo: selectedYearNo.value,
        batchId,
      })
      await Promise.all([loadConfig({ silent: true }), loadPool({ silent: true }), loadSelectionStatus({ silent: true })])
      pageMessage.value = {
        type: 'success',
        text: `${selectedYearNo.value} 年订单池已确认，固化 ${result.generatedCount} 个订单。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '确认订单池失败')
      throw error
    } finally {
      confirmingPool.value = false
    }
  }

  async function generateSelectionSequence() {
    generatingSequence.value = true
    pageMessage.value = null
    try {
      const result = await generateAdminSelectionSequence({
        yearNo: selectedYearNo.value,
      })
      await Promise.all([loadConfig({ silent: true }), loadSelectionStatus({ silent: true })])
      pageMessage.value = {
        type: 'success',
        text: `已生成选单顺序，处理 ${result.affectedCount} 个标段。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '生成选单顺序失败')
      throw error
    } finally {
      generatingSequence.value = false
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

  async function loadSelectionStatus(options?: { silent?: boolean }) {
    loadingSelectionStatus.value = true
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      marketSelectionStatus.value = await getAdminMarketSelectionStatus(selectedYearNo.value, controlForm.marketCode)
    } catch (error) {
      marketSelectionStatus.value = null
      if (!options?.silent) {
        pageMessage.value = toErrorMessage(error, '获取市场选单状态失败')
      }
      if (!options?.silent) {
        throw error
      }
    } finally {
      loadingSelectionStatus.value = false
    }
  }

  async function openMarket() {
    controllingMarket.value = true
    pageMessage.value = null
    try {
      await openAdminMarketBidding({
        yearNo: selectedYearNo.value,
        marketCode: controlForm.marketCode,
      })
      await Promise.all([loadConfig({ silent: true }), loadSelectionStatus({ silent: true })])
      pageMessage.value = {
        type: 'success',
        text: `${marketName(controlForm.marketCode)}投入已开放。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '开放市场投入失败')
      throw error
    } finally {
      controllingMarket.value = false
    }
  }

  async function closeMarket() {
    controllingMarket.value = true
    pageMessage.value = null
    try {
      await closeAdminMarketBidding({
        yearNo: selectedYearNo.value,
        marketCode: controlForm.marketCode,
      })
      await Promise.all([loadConfig({ silent: true }), loadSelectionStatus({ silent: true })])
      pageMessage.value = {
        type: 'success',
        text: `${marketName(controlForm.marketCode)}投入已关闭，选单顺序已生成或市场已跳过。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '关闭市场投入失败')
      throw error
    } finally {
      controllingMarket.value = false
    }
  }

  async function releaseNextSegment() {
    releasingSegment.value = true
    pageMessage.value = null
    try {
      const result = await releaseNextAdminOrderSegment({
        yearNo: selectedYearNo.value,
      })
      if (result.marketCode) {
        controlForm.marketCode = result.marketCode
      }
      await Promise.all([loadConfig({ silent: true }), loadSelectionStatus({ silent: true })])
      pageMessage.value = {
        type: 'success',
        text: result.orderTypeName ? `已释放 ${result.marketName} ${result.orderTypeName}。` : '已处理下一个标段。',
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '释放下一个标段失败')
      throw error
    } finally {
      releasingSegment.value = false
    }
  }

  async function skipCurrentGroup() {
    const segment = currentSegment.value
    if (!segment?.currentGroupId) {
      pageMessage.value = {
        type: 'error',
        text: '当前没有可跳过的小组。',
      }
      return
    }
    skippingGroup.value = true
    pageMessage.value = null
    try {
      await adminSkipCurrentOrderGroup({
        yearNo: selectedYearNo.value,
        marketCode: segment.marketCode,
        orderType: segment.orderType,
        groupId: segment.currentGroupId,
        reason: controlForm.skipReason.trim(),
      })
      controlForm.skipReason = ''
      await loadSelectionStatus({ silent: true })
      pageMessage.value = {
        type: 'success',
        text: '已跳过当前小组。',
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '跳过当前小组失败')
      throw error
    } finally {
      skippingGroup.value = false
    }
  }

  function applyConfig(result: OrderControlConfigResult) {
    config.value = result
    editableMarketConfigs.value = (result.marketConfigs ?? defaultMarketConfigs(result.yearNo)).map((item) => ({ ...item }))
    const marketEnabledMap = new Map(editableMarketConfigs.value.map((item) => [item.marketCode, item.enabled]))
    editableItems.value = result.items.map((item) => ({
      ...item,
      marketEnabled: item.marketEnabled ?? marketEnabledMap.get(item.marketCode) ?? item.marketCode === 'LOCAL',
    }))
  }

  function applyForecastControl(result: OrderForecastControlResult) {
    forecastControl.value = result
    editableForecastItems.value = (result.items ?? []).map((item) => ({ ...item }))
    editableForecastNarratives.value = (result.narratives ?? []).map((item) => ({ ...item }))
  }

  function setYear(yearNo: number) {
    selectedYearNo.value = Math.max(yearNo, 1)
  }

  function setControlMarket(marketCode: OrderMarketCode) {
    controlForm.marketCode = marketCode
  }

  return {
    selectedYearNo,
    config,
    forecastControl,
    uploadResult,
    orderPool,
    marketSelectionStatus,
    loading,
    uploading,
    savingConfig,
    savingForecastControl,
    savingMarketConfig,
    generatingPool,
    confirmingPool,
    generatingSequence,
    loadingPool,
    loadingSelectionStatus,
    controllingMarket,
    releasingSegment,
    skippingGroup,
    pageMessage,
    poolFilter,
    controlForm,
    editableItems,
    editableForecastItems,
    editableForecastNarratives,
    editableMarketConfigs,
    sortedItems,
    forecastStages,
    totalOrderCount,
    totalGeneratedCount,
    hasLockedConfig,
    enabledMarketCount,
    currentSegment,
    bootstrap,
    loadForecastControl,
    loadConfig,
    uploadExcel,
    saveMarketConfig,
    saveForecastControl,
    saveConfig,
    generatePool,
    confirmPool,
    generateSelectionSequence,
    loadPool,
    loadSelectionStatus,
    openMarket,
    closeMarket,
    releaseNextSegment,
    skipCurrentGroup,
    setYear,
    setControlMarket,
  }
})

function marketName(code: OrderMarketCode) {
  return MARKET_OPTIONS.find((item) => item.code === code)?.name ?? code
}

function defaultMarketConfigs(yearNo: number): OrderMarketConfigItem[] {
  return MARKET_OPTIONS.map((item) => ({
    yearNo,
    marketCode: item.code,
    marketName: item.name,
    enabled: item.code === 'LOCAL',
    configStatus: 'DRAFT',
    lockedBatchId: null,
  }))
}

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}
