import { defineStore } from 'pinia'
import { computed, reactive, ref } from 'vue'

import { getCurrentGameConfig, getYearTabs } from '@/api/sandbox-game/game-config'
import { useDictionaryStore } from '@/stores/dictionary'
import {
  deliverPlayerOrders,
  getPlayerMarketForecast,
  getPlayerOrderYearView,
  passPlayerOrderSegment,
  selectPlayerOrder,
  submitPlayerMarketInvestment,
} from '@/api/sandbox-game/player-order'
import type { CurrentGameConfigResult, YearTabItem, YearTabsResult } from '@/types/sandbox-game'
import type {
  OrderDeliveryStageCode,
  OrderMarketCode,
  OrderMarketForecastResult,
  OrderTypeCode,
  PlayerOrderMarketView,
  PlayerOrderSegmentView,
  PlayerOrderYearView,
} from '@/types/sandbox-game-order'

type MessageType = 'success' | 'error' | 'info'

export interface PageMessage {
  type: MessageType
  text: string
}

export const PLAYER_ORDER_MARKETS: Array<{ code: OrderMarketCode; name: string }> = [
  { code: 'LOCAL', name: '本地市场' },
  { code: 'REGIONAL', name: '区域市场' },
  { code: 'NATIONAL', name: '全国市场' },
  { code: 'GLOBAL', name: '全球市场' },
]

export const PLAYER_ORDER_TYPES: Array<{ code: OrderTypeCode; name: string }> = [
  { code: 'AGENCY_INSPECTION', name: '代办过检' },
  { code: 'TWO_CABIN_VIP', name: '两舱贵宾' },
  { code: 'BUSINESS_VIP', name: '商务贵宾' },
  { code: 'MEMBER_CUSTOM', name: '会员定制' },
]

export const usePlayerOrderStore = defineStore('sandbox-player-order', () => {
  const dictionaryStore = useDictionaryStore()
  const currentConfig = ref<CurrentGameConfigResult | null>(null)
  const yearTabs = ref<YearTabItem[]>([])
  const currentView = ref<PlayerOrderYearView | null>(null)
  const marketForecast = ref<OrderMarketForecastResult | null>(null)
  const selectedYear = ref(1)
  const selectedMarketCode = ref<OrderMarketCode>('LOCAL')
  const loading = ref(false)
  const yearViewLoading = ref(false)
  const submittingInvestment = ref(false)
  const selectingOrder = ref(false)
  const passingSegment = ref(false)
  const deliveringOrders = ref(false)
  const pageMessage = ref<PageMessage | null>(null)
  const investmentDraft = reactive<Record<string, number | null>>(createEmptyInvestmentDraft())
  const investmentDraftYear = ref<number | null>(null)
  const investmentErrors = ref<Record<string, string>>({})

  const markets = computed(() => currentView.value?.markets ?? [])
  const selectedMarket = computed<PlayerOrderMarketView | null>(
    () => markets.value.find((item) => item.marketCode === selectedMarketCode.value) ?? markets.value[0] ?? null,
  )
  const currentSegment = computed<PlayerOrderSegmentView | null>(() => {
    const segments = selectedMarket.value?.segments ?? []
    return segments.find((item) => item.segmentStatus === 'SELECTING') ?? segments.find((item) => item.canSelectOrder || item.canPassSegment) ?? null
  })
  const selectableSegments = computed(() =>
    markets.value.flatMap((market) => market.segments).filter((segment) => segment.segmentStatus === 'SELECTING'),
  )
  const currentTab = computed(() => yearTabs.value.find((item) => item.yearNo === selectedYear.value) ?? null)

  async function bootstrap(preferredYear?: number) {
    loading.value = true
    pageMessage.value = null
    try {
      const [config, tabsResult, forecast] = await Promise.all([getCurrentGameConfig(), getYearTabs(), getPlayerMarketForecast()])
      currentConfig.value = config
      yearTabs.value = tabsResult.tabs
      marketForecast.value = forecast
      const nextYear = Math.max(resolveInitialYear(tabsResult, preferredYear), 0)
      await loadYearView(nextYear, { silent: true })
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '初始化年度订单页失败')
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

  async function loadYearView(yearNo: number, options?: { silent?: boolean }) {
    const showLoading = !options?.silent
    if (showLoading) {
      yearViewLoading.value = true
    }
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      const view = await getPlayerOrderYearView(yearNo)
      currentView.value = view
      selectedYear.value = yearNo
      applyInvestmentDraft(view, investmentDraft, investmentDraftYear.value === yearNo)
      investmentDraftYear.value = yearNo
      const selectedExists = view.markets?.some((item) => item.marketCode === selectedMarketCode.value)
      if (!selectedExists && view.markets?.[0]) {
        selectedMarketCode.value = view.markets[0].marketCode
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, `读取 ${yearNo} 年订单页失败`)
      throw error
    } finally {
      if (showLoading) {
        yearViewLoading.value = false
      }
    }
  }

  async function submitInvestments() {
    const validation = validateInvestmentDraft(currentView.value, investmentDraft, dictionaryStore.orderTypeName)
    investmentErrors.value = validation.errors
    if (validation.message) {
      pageMessage.value = {
        type: 'error',
        text: validation.message,
      }
      return
    }
    submittingInvestment.value = true
    pageMessage.value = null
    try {
      await submitPlayerMarketInvestment({
        yearNo: selectedYear.value,
        investments: buildInvestmentPayload(investmentDraft, currentView.value),
      })
      pageMessage.value = {
        type: 'success',
        text: '16 项市场投入已提交。',
      }
      investmentErrors.value = {}
      await loadYearView(selectedYear.value, { silent: true })
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '提交市场投入失败')
      throw error
    } finally {
      submittingInvestment.value = false
    }
  }

  async function selectOrder(segment: PlayerOrderSegmentView, orderId: number) {
    selectingOrder.value = true
    pageMessage.value = null
    try {
      await selectPlayerOrder({
        yearNo: selectedYear.value,
        marketCode: segment.marketCode,
        orderType: segment.orderType,
        orderId,
      })
      pageMessage.value = {
        type: 'success',
        text: `${segment.marketName} ${segment.orderTypeName} 已选择订单 ${formatOrderNo(segment.availableOrders.find((item) => item.orderId === orderId)?.businessOrderNo, orderId)}。`,
      }
      await loadYearView(selectedYear.value, { silent: true })
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '选择订单失败')
      throw error
    } finally {
      selectingOrder.value = false
    }
  }

  async function passSegment(segment: PlayerOrderSegmentView) {
    passingSegment.value = true
    pageMessage.value = null
    try {
      await passPlayerOrderSegment({
        yearNo: selectedYear.value,
        marketCode: segment.marketCode,
        orderType: segment.orderType,
      })
      pageMessage.value = {
        type: 'success',
        text: `${segment.marketName} ${segment.orderTypeName} 已放弃。`,
      }
      await loadYearView(selectedYear.value, { silent: true })
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '放弃标段失败')
      throw error
    } finally {
      passingSegment.value = false
    }
  }

  async function deliverSelectedOrder(segment: PlayerOrderSegmentView, stageCode: OrderDeliveryStageCode) {
    if (!segment.selectedOrder) {
      return
    }
    await deliverOrders([segment.selectedOrder.orderId], stageCode)
  }

  async function deliverOrders(orderIds: number[], stageCode: OrderDeliveryStageCode) {
    const uniqueOrderIds = Array.from(new Set(orderIds.filter((item) => Number.isFinite(item) && item > 0)))
    if (uniqueOrderIds.length === 0) {
      return
    }
    deliveringOrders.value = true
    pageMessage.value = null
    try {
      await deliverPlayerOrders({
        yearNo: selectedYear.value,
        stageCode,
        orderIds: uniqueOrderIds,
      })
      pageMessage.value = {
        type: 'success',
        text: `已交付 ${uniqueOrderIds.length} 个订单。`,
      }
      await loadYearView(selectedYear.value, { silent: true })
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '交付订单失败')
      throw error
    } finally {
      deliveringOrders.value = false
    }
  }

  function setSelectedMarket(marketCode: OrderMarketCode) {
    selectedMarketCode.value = marketCode
  }

  function setInvestmentDraft(marketCode: OrderMarketCode, orderType: OrderTypeCode, value: number | null) {
    investmentDraft[investmentKey(marketCode, orderType)] = value
    investmentErrors.value = {}
  }

  return {
    currentConfig,
    yearTabs,
    currentView,
    marketForecast,
    selectedYear,
    selectedMarketCode,
    loading,
    yearViewLoading,
    submittingInvestment,
    selectingOrder,
    passingSegment,
    deliveringOrders,
    pageMessage,
    investmentDraft,
    investmentErrors,
    markets,
    selectedMarket,
    currentSegment,
    selectableSegments,
    currentTab,
    bootstrap,
    refreshTabs,
    loadYearView,
    submitInvestments,
    selectOrder,
    passSegment,
    deliverSelectedOrder,
    deliverOrders,
    setSelectedMarket,
    setInvestmentDraft,
  }
})

function applyInvestmentDraft(view: PlayerOrderYearView, draft: Record<string, number | null>, preserveUnsavedDraft: boolean) {
  const previousDraft = preserveUnsavedDraft ? { ...draft } : {}
  for (const key of Object.keys(draft)) {
    draft[key] = null
  }
  for (const market of view.markets ?? []) {
    for (const segment of market.segments) {
      const key = investmentKey(segment.marketCode, segment.orderType)
      if (!segment.marketEnabled) {
        draft[key] = 0
      } else {
        draft[key] = segment.investmentSubmitted ? segment.marketInvestment : previousDraft[key] ?? 0
      }
    }
  }
}

function resolveInitialYear(result: YearTabsResult, preferredYear?: number) {
  if (typeof preferredYear === 'number' && Number.isFinite(preferredYear)) {
    const matched = result.tabs.find((item) => item.yearNo === preferredYear)
    if (matched?.canEnter) {
      return preferredYear
    }
  }

  const currentOpen = result.tabs.find((item) => item.isCurrentOpenYear && item.yearNo > 0)
  if (currentOpen?.canEnter) {
    return currentOpen.yearNo
  }

  const firstFormal = result.tabs.find((item) => item.canEnter && item.yearNo > 0)
  if (firstFormal) {
    return firstFormal.yearNo
  }

  const firstEnterable = result.tabs.find((item) => item.canEnter)
  return firstEnterable?.yearNo ?? 0
}

export function investmentKey(marketCode: OrderMarketCode, orderType: OrderTypeCode) {
  return `${marketCode}|${orderType}`
}

function createEmptyInvestmentDraft() {
  const result: Record<string, number | null> = {}
  for (const market of PLAYER_ORDER_MARKETS) {
    for (const orderType of PLAYER_ORDER_TYPES) {
      result[investmentKey(market.code, orderType.code)] = null
    }
  }
  return result
}

function buildInvestmentPayload(draft: Record<string, number | null>, view: PlayerOrderYearView | null) {
  const marketEnabledMap = new Map((view?.markets ?? []).map((market) => [market.marketCode, market.marketEnabled]))
  return PLAYER_ORDER_MARKETS.flatMap((market) =>
    PLAYER_ORDER_TYPES.map((orderType) => {
      const raw = Number(draft[investmentKey(market.code, orderType.code)] ?? 0)
      if (marketEnabledMap.get(market.code) === false) {
        return {
          marketCode: market.code,
          orderType: orderType.code,
          marketInvestment: 0,
        }
      }
      return {
        marketCode: market.code,
        orderType: orderType.code,
        marketInvestment: Number.isFinite(raw) ? raw : 0,
      }
    }),
  )
}

function validateInvestmentDraft(
  view: PlayerOrderYearView | null,
  draft: Record<string, number | null>,
  resolveOrderTypeName: (code: string, fallback: string) => string,
) {
  const errors: Record<string, string> = {}
  if (!view?.markets) {
    return { errors, message: '当前没有可提交的市场投入数据。' }
  }
  for (const market of view.markets) {
    let marketTotal = 0
    for (const orderType of PLAYER_ORDER_TYPES) {
      const key = investmentKey(market.marketCode, orderType.code)
      const rawValue = market.marketEnabled ? draft[key] : 0
      const value = Number(rawValue ?? 0)
      if (!Number.isFinite(value) || value < 0) {
        errors[key] = `${market.marketName} ${resolveOrderTypeName(orderType.code, orderType.name)} 投入必须为非负整数`
        continue
      }
      if (!Number.isInteger(value)) {
        errors[key] = `${market.marketName} ${resolveOrderTypeName(orderType.code, orderType.name)} 投入必须为整数`
        continue
      }
      marketTotal += value
    }
    if (market.marketEnabled && market.marketInvestmentLimit !== null && market.marketInvestmentLimit !== undefined && marketTotal > market.marketInvestmentLimit) {
      for (const orderType of PLAYER_ORDER_TYPES) {
        errors[investmentKey(market.marketCode, orderType.code)] = `${market.marketName} 4项投入合计不能超过 ${market.marketInvestmentLimit}M`
      }
    }
  }
  const firstMessage = Object.values(errors)[0] ?? ''
  return { errors, message: firstMessage }
}

function formatOrderNo(businessOrderNo: string | undefined, orderId: number) {
  return businessOrderNo || `#${orderId}`
}

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}
