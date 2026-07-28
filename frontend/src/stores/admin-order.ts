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
  openNextAdminOrderRound,
  releaseNextAdminOrderSegment,
  updateAdminOrderControlConfig,
  updateAdminOrderForecastControl,
  updateAdminOrderMarketConfig,
  uploadAdminOrderExcel,
} from '@/api/sandbox-game/admin-order'
import { useDictionaryStore } from '@/stores/dictionary'
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
import type { AdminMarketSelectionStatus, AdminOrderSegmentStatus, OrderTemplateMeta } from '@/types/sandbox-game-order'
import { hasFractionInput } from '@/utils/manual-integer'
import { DEFAULT_PLAYER_ORDER_TEMPLATE } from '@/stores/player-order'

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
  const dictionaryStore = useDictionaryStore()
  const selectedYearNo = ref(1)
  const config = ref<OrderControlConfigResult | null>(null)
  const forecastControl = ref<OrderForecastControlResult | null>(null)
  const uploadResult = ref<UploadOrderExcelResult | null>(null)
  const orderPool = ref<OrderPoolResult | null>(null)
  const marketSelectionStatus = ref<AdminMarketSelectionStatus | null>(null)
  const marketSelectionStatusSnapshots = ref<AdminMarketSelectionStatus[]>([])
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
  const silentLoadingSelectionStatus = ref(false)
  const controllingMarket = ref(false)
  const releasingSegment = ref(false)
  const openingNextRound = ref(false)
  const skippingGroup = ref(false)
  const pageMessage = ref<PageMessage | null>(null)

  const poolFilter = reactive({
    marketCode: 'ALL' as OrderMarketCode | 'ALL',
    orderType: 'ALL' as AdminOrderType | 'ALL',
    selectedGroupKey: 'ALL',
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
  const forecastYearLocks = computed(() => forecastControl.value?.yearLocks ?? [])
  const enabledMarketCount = computed(() => editableMarketConfigs.value.filter((item) => item.enabled).length)
  const currentSegment = computed(() => {
    const status = marketSelectionStatus.value
    if (status?.currentSegment) {
      return status.currentSegment
    }
    return status?.segments.find((item) => item.segmentStatus === 'SELECTING')
      ?? status?.segments.find((item) => item.segmentStatus === 'ROUND_READY')
      ?? null
  })
  const orderTemplate = computed<OrderTemplateMeta>(() =>
    config.value?.orderTemplate
      ?? forecastControl.value?.orderTemplate
      ?? orderPool.value?.orderTemplate
      ?? marketSelectionStatus.value?.orderTemplate
      ?? DEFAULT_PLAYER_ORDER_TEMPLATE,
  )
  const marketOptions = computed(() => orderTemplate.value.markets ?? DEFAULT_PLAYER_ORDER_TEMPLATE.markets)
  const orderTypeOptions = computed(() => orderTemplate.value.orderTypes ?? DEFAULT_PLAYER_ORDER_TEMPLATE.orderTypes)
  const maxCardCount = computed(() => orderTemplate.value.maxCardCount ?? DEFAULT_PLAYER_ORDER_TEMPLATE.maxCardCount)
  const segmentCount = computed(() => orderTemplate.value.segmentCount || marketOptions.value.length * orderTypeOptions.value.length)

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
      await loadForecastControl({ silent: true })
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
    const sequenceInvalid = editableItems.value.find((item) => !Number.isInteger(Number(item.releaseSequenceNo)) || Number(item.releaseSequenceNo) < 1)
    if (sequenceInvalid) {
      pageMessage.value = {
        type: 'error',
        text: '释放顺序必须填写正整数',
      }
      return
    }
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
        orderTemplate: result.items.length > 0 ? orderTemplate.value : config.value?.orderTemplate ?? orderTemplate.value,
        latestBatchId: config.value?.latestBatchId ?? null,
        latestBatchUploadedAt: config.value?.latestBatchUploadedAt ?? null,
        forecast: forecastControl.value?.forecast ?? config.value?.forecast ?? { formulaVersion: orderTemplate.value.formulaVersion, orderTemplate: orderTemplate.value, stages: [] },
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
      await loadConfig({ silent: true })
      await loadPool({ silent: true })
      pageMessage.value = {
        type: 'success',
        text: `订单配置已保存，并已刷新 ${selectedYearNo.value} 年预览订单池。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '保存订单配置失败')
      throw error
    } finally {
      savingConfig.value = false
    }
  }

  async function saveForecastControl() {
    const forecastInvalid = editableForecastItems.value.find((item) => hasFractionInput(item.orderCount) || Number(item.orderCount) < 0)
    if (forecastInvalid) {
      pageMessage.value = {
        type: 'error',
        text: '订单数量控制台只能填写非负整数',
      }
      return
    }
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
      await loadPool({ silent: true })
      pageMessage.value = {
        type: 'success',
        text: result.autoPreviewYears?.length
          ? `多年订单数量控制台已保存，并已刷新 ${result.autoPreviewYears.map((yearNo) => `${yearNo}年`).join('、')} 预览订单池。`
          : `多年订单数量控制台已保存，操作人 ${result.updatedBy}。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '保存多年订单数量控制台失败')
      throw error
    } finally {
      savingForecastControl.value = false
    }
  }

  async function saveMarketConfig() {
    const marketConfigError = validateMarketConfigDraft(editableMarketConfigs.value)
    if (marketConfigError) {
      pageMessage.value = {
        type: 'error',
        text: marketConfigError,
      }
      return
    }
    savingMarketConfig.value = true
    pageMessage.value = null
    try {
      const result = await updateAdminOrderMarketConfig({
        yearNo: selectedYearNo.value,
        markets: editableMarketConfigs.value.map((item) => ({
          marketCode: item.marketCode,
          enabled: Boolean(item.enabled),
          marketInvestmentLimit: toNullableNumber(item.marketInvestmentLimit),
        })),
      })
      applyConfig({
        yearNo: result.yearNo,
        finalYear: config.value?.finalYear ?? selectedYearNo.value,
        orderTemplate: config.value?.orderTemplate ?? orderTemplate.value,
        latestBatchId: config.value?.latestBatchId ?? null,
        latestBatchUploadedAt: config.value?.latestBatchUploadedAt ?? null,
        forecast: forecastControl.value?.forecast ?? config.value?.forecast ?? { formulaVersion: orderTemplate.value.formulaVersion, orderTemplate: orderTemplate.value, stages: [] },
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
      await loadConfig({ silent: true })
      await loadPool({ silent: true })
      pageMessage.value = {
        type: 'success',
        text: `市场开启配置已保存，并已刷新 ${selectedYearNo.value} 年预览订单池。`,
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
    await loadConfig({ silent: true })
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
      await Promise.all([loadConfig({ silent: true }), loadPool({ silent: true }), loadSelectionStatus({ silent: true }), loadForecastControl({ silent: true })])
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
      const groupFilter = resolvePoolGroupFilter(poolFilter.selectedGroupKey)
      orderPool.value = await getAdminOrderPool(selectedYearNo.value, {
        marketCode: poolFilter.marketCode,
        orderType: poolFilter.orderType,
        ...groupFilter,
      })
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '获取订单池失败')
      throw error
    } finally {
      loadingPool.value = false
    }
  }

  async function loadSelectionStatus(options?: { silent?: boolean; skipFocus?: boolean }) {
    if (options?.silent) {
      if (silentLoadingSelectionStatus.value || loadingSelectionStatus.value) {
        return
      }
      silentLoadingSelectionStatus.value = true
    } else {
      loadingSelectionStatus.value = true
    }
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      const focusedStatus = options?.skipFocus ? null : await syncControlMarketToOrderFocus()
      marketSelectionStatus.value = focusedStatus?.marketCode === controlForm.marketCode
        ? focusedStatus
        : await getAdminMarketSelectionStatus(selectedYearNo.value, controlForm.marketCode)
    } catch (error) {
      marketSelectionStatus.value = null
      if (!options?.silent) {
        pageMessage.value = toErrorMessage(error, '获取市场选单状态失败')
      }
      if (!options?.silent) {
        throw error
      }
    } finally {
      if (options?.silent) {
        silentLoadingSelectionStatus.value = false
      } else {
        loadingSelectionStatus.value = false
      }
    }
  }

  async function syncControlMarketToOrderFocus() {
    const snapshots = await loadAllMarketSelectionStatusSnapshots()
    const focusSegment = resolveAdminFocusSegment(snapshots)
    if (!focusSegment) {
      return snapshots.find((item) => item.marketCode === controlForm.marketCode) ?? null
    }
    if (controlForm.marketCode !== focusSegment.marketCode) {
      controlForm.marketCode = focusSegment.marketCode
    }
    return snapshots.find((item) => item.marketCode === controlForm.marketCode) ?? null
  }

  async function loadAllMarketSelectionStatusSnapshots() {
    const markets = marketOptions.value
    if (markets.length === 0) {
      return []
    }
    const results = await Promise.allSettled(
      markets.map((market) => getAdminMarketSelectionStatus(selectedYearNo.value, market.code)),
    )
    const snapshots = results
      .filter((item): item is PromiseFulfilledResult<AdminMarketSelectionStatus> => item.status === 'fulfilled')
      .map((item) => item.value)
    if (snapshots.length > 0) {
      marketSelectionStatusSnapshots.value = snapshots
    }
    return snapshots
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
        text: `${dictionaryStore.marketName(controlForm.marketCode, templateMarketName(controlForm.marketCode))}投入已开放。`,
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
        text: `${dictionaryStore.marketName(controlForm.marketCode, templateMarketName(controlForm.marketCode))}投入已关闭，选单顺序已生成或市场已跳过。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '关闭市场投入失败')
      throw error
    } finally {
      controllingMarket.value = false
    }
  }

  async function releaseNextSegment() {
    const segment = currentSegment.value
    if (segment?.segmentStatus === 'SELECTING' || segment?.segmentStatus === 'ROUND_READY') {
      pageMessage.value = {
        type: 'error',
        text: segment.segmentStatus === 'ROUND_READY'
          ? '当前标段还有等待开启的下一轮，请先开启下一轮。'
          : '当前标段正在选单中，不能释放下一个标段。',
      }
      return
    }
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

  async function openNextRound() {
    const segment = currentSegment.value
    if (!segment || segment.segmentStatus !== 'ROUND_READY' || !segment.nextRoundNo) {
      pageMessage.value = { type: 'error', text: '当前没有等待开启的下一轮。' }
      return
    }
    openingNextRound.value = true
    pageMessage.value = null
    try {
      const result = await openNextAdminOrderRound({
        yearNo: selectedYearNo.value,
        marketCode: segment.marketCode,
        orderType: segment.orderType,
      })
      await loadSelectionStatus({ silent: true })
      pageMessage.value = {
        type: 'success',
        text: `已开启第 ${result.currentRoundNo ?? segment.nextRoundNo} 轮。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '开启下一轮失败')
      throw error
    } finally {
      openingNextRound.value = false
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
    editableMarketConfigs.value = (result.marketConfigs ?? defaultMarketConfigs(result.yearNo, result.orderTemplate ?? orderTemplate.value, dictionaryStore.marketName)).map((item) => ({ ...item }))
    const marketEnabledMap = new Map(editableMarketConfigs.value.map((item) => [item.marketCode, item.enabled]))
    editableItems.value = result.items.map((item) => ({
      ...item,
      marketEnabled: item.marketEnabled ?? marketEnabledMap.get(item.marketCode) ?? (result.orderTemplate ?? orderTemplate.value).markets.find((market) => market.code === item.marketCode)?.defaultEnabled ?? false,
    }))
    ensureControlMarketIsValid()
  }

  function applyForecastControl(result: OrderForecastControlResult) {
    forecastControl.value = result
    editableForecastItems.value = (result.items ?? []).map((item) => ({ ...item }))
    editableForecastNarratives.value = (result.narratives ?? []).map((item) => ({ ...item }))
    ensureControlMarketIsValid()
  }

  function setYear(yearNo: number) {
    selectedYearNo.value = Math.max(yearNo, 1)
  }

  function setControlMarket(marketCode: OrderMarketCode) {
    controlForm.marketCode = marketCode
  }

  function ensureControlMarketIsValid() {
    const markets = orderTemplate.value.markets ?? DEFAULT_PLAYER_ORDER_TEMPLATE.markets
    if (markets.some((item) => item.code === controlForm.marketCode)) {
      return
    }
    controlForm.marketCode = markets[0]?.code ?? 'LOCAL'
  }

  function templateMarketName(code: OrderMarketCode) {
    return orderTemplate.value.markets.find((item) => item.code === code)?.name ?? MARKET_OPTIONS.find((item) => item.code === code)?.name ?? code
  }

  return {
    selectedYearNo,
    config,
    forecastControl,
    uploadResult,
    orderPool,
    marketSelectionStatus,
    marketSelectionStatusSnapshots,
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
    openingNextRound,
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
    forecastYearLocks,
    totalOrderCount,
    totalGeneratedCount,
    hasLockedConfig,
    enabledMarketCount,
    currentSegment,
    orderTemplate,
    marketOptions,
    orderTypeOptions,
    maxCardCount,
    segmentCount,
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
    openNextRound,
    skipCurrentGroup,
    setYear,
    setControlMarket,
  }
})

function resolveAdminFocusSegment(statuses: AdminMarketSelectionStatus[]) {
  return statuses
    .flatMap((status) => status.segments)
    .map((segment) => ({
      segment,
      priority: adminSegmentFocusPriority(segment.segmentStatus),
    }))
    .filter((item) => item.priority > 0)
    .sort((a, b) => a.priority - b.priority || a.segment.releaseSequenceNo - b.segment.releaseSequenceNo)[0]?.segment ?? null
}

function adminSegmentFocusPriority(segmentStatus: AdminOrderSegmentStatus['segmentStatus']) {
  switch (segmentStatus) {
    case 'SELECTING':
      return 1
    case 'ROUND_READY':
      return 2
    case 'SEQUENCE_READY':
    case 'WAITING_RELEASE':
      return 3
    default:
      return 0
  }
}

function resolvePoolGroupFilter(selectedGroupKey: string) {
  if (selectedGroupKey === 'SELECTED') {
    return {
      selectedOnly: true,
      selectedGroupId: null,
    }
  }
  if (selectedGroupKey.startsWith('GROUP:')) {
    const groupId = Number(selectedGroupKey.slice('GROUP:'.length))
    return {
      selectedOnly: true,
      selectedGroupId: Number.isFinite(groupId) && groupId > 0 ? groupId : null,
    }
  }
  return {
    selectedOnly: false,
    selectedGroupId: null,
  }
}

function marketName(code: OrderMarketCode) {
  return MARKET_OPTIONS.find((item) => item.code === code)?.name ?? code
}

function defaultMarketConfigs(yearNo: number, template: OrderTemplateMeta, resolveMarketName: (code: string, fallback: string) => string): OrderMarketConfigItem[] {
  return (template.markets ?? DEFAULT_PLAYER_ORDER_TEMPLATE.markets).map((item) => ({
    yearNo,
    marketCode: item.code,
    marketName: resolveMarketName(item.code, item.name),
    enabled: item.defaultEnabled,
    marketInvestmentLimit: null,
    configStatus: 'DRAFT',
    lockedBatchId: null,
  }))
}

function validateMarketConfigDraft(items: OrderMarketConfigItem[]) {
  const invalid = items.find((item) => {
    const value = toNullableNumber(item.marketInvestmentLimit)
    return value !== null && (!Number.isInteger(value) || value < 0)
  })
  if (!invalid) {
    return ''
  }
  return `${invalid.marketName} 单市场投入上限必须为空或非负整数`
}

function toNullableNumber(value: number | null | undefined) {
  if (value === null || value === undefined || !Number.isFinite(Number(value))) {
    return null
  }
  return Number(value)
}

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}
