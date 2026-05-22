<template>
  <div class="player-page">
    <div class="shell">
      <header class="page-header">
        <div>
          <p class="eyebrow">Sandbox Game / Player</p>
          <h1>玩家订单页</h1>
          <p class="subtext">按管理员释放的市场标段提交投入、查看顺序并选择订单。</p>
        </div>
        <div class="header-pills">
          <span class="pill">组别：{{ currentView?.groupId ?? '--' }}</span>
          <span class="pill">开放年份：{{ currentConfig?.currentOpenYear ?? '--' }}</span>
          <span class="pill">最终年份：{{ currentConfig?.finalYear ?? '--' }}</span>
          <span v-if="currentUser" class="pill">账号：{{ currentUser.username }}</span>
          <button class="logout-button" type="button" @click="handleLogout">退出登录</button>
        </div>
      </header>

      <section class="toolbar-card">
        <div class="toolbar-top">
          <YearTabs :tabs="yearTabs" :active-year="selectedYear" @select="handleYearSelect" />
          <PageModeSwitch active-mode="order" :report-enabled="true" @operating="goOperating" @report="goReport" />
        </div>
      </section>

      <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
        {{ pageMessage.text }}
      </section>

      <main class="order-workspace">
        <section class="status-grid">
          <article class="status-card">
            <span>当前年份</span>
            <strong>{{ selectedYear }} 年</strong>
          </article>
          <article class="status-card">
            <span>订单前置</span>
            <strong>{{ currentView?.orderRequired === false ? '无需订单' : '正式流程' }}</strong>
          </article>
          <article class="status-card">
            <span>当前标段</span>
            <strong>{{ activeSegmentTitle }}</strong>
          </article>
          <article class="status-card">
            <span>市场投入</span>
            <strong>{{ currentView?.investmentSubmitted ? '已提交' : currentView?.canSubmitInvestment ? '待提交' : '未开放' }}</strong>
          </article>
          <article class="status-card">
            <span>轮询间隔</span>
            <strong>{{ pollingText }}</strong>
          </article>
        </section>

        <section v-if="loading || yearViewLoading" class="loading-card">正在加载订单页数据...</section>

        <section v-else-if="currentView?.orderRequired === false" class="empty-state">
          <strong>0 年不需要线上订单</strong>
          <span>0 年是引导年，可以直接进入经营页熟悉规则。</span>
        </section>

        <template v-else-if="currentView">
          <section v-if="marketForecast" class="forecast-panel">
            <div class="panel-head">
              <div>
                <strong>市场预测</strong>
                <span>订单预测按 1~3年、4~5年、6~8年展示，供开标前判断投入方向。</span>
              </div>
              <span class="formula-pill">{{ marketForecast.formulaVersion || '--' }}</span>
            </div>
            <div class="forecast-stage-grid">
              <article v-for="stage in marketForecast.stages" :key="stage.forecastStageCode" class="forecast-stage-card">
                <header>
                  <strong>{{ stage.forecastStageName }}</strong>
                  <span>{{ stage.yearRange }}</span>
                </header>
                <div class="forecast-market-grid">
                  <section v-for="market in stage.markets" :key="market.marketCode" class="forecast-market-card">
                    <div class="forecast-market-head">
                      <strong>{{ market.marketName }}</strong>
                      <span>{{ formatIntegerAmount(stageMarketAmount(market)) }}</span>
                    </div>
                    <p v-if="market.narrative">{{ market.narrative }}</p>
                    <div class="forecast-chart">
                      <div v-for="year in market.years" :key="year.yearNo" class="forecast-year-group">
                        <div class="forecast-bars">
                          <span
                            v-for="product in year.products"
                            :key="product.orderType"
                            class="forecast-bar"
                            :class="`type-${product.orderType}`"
                            :style="{ height: `${forecastBarHeight(product.forecastAmount, market)}%` }"
                            :title="`${year.yearNo}年 ${product.orderTypeName}: ${product.orderCount}单 / ${formatIntegerAmount(product.forecastAmount)}`"
                          />
                        </div>
                        <span class="forecast-year-label">{{ year.yearNo }}年</span>
                        <em>{{ year.totalOrderCount }}单 / {{ formatIntegerAmount(year.totalForecastAmount) }}</em>
                      </div>
                    </div>
                    <div class="forecast-legend">
                      <span v-for="orderType in orderTypeOptions" :key="orderType.code" :class="`type-${orderType.code}`">
                        {{ orderType.name }}
                      </span>
                    </div>
                  </section>
                </div>
              </article>
            </div>
          </section>

          <section class="investment-panel">
            <div class="panel-head">
              <div>
                <strong>年度市场投入</strong>
                <span>开标前一次性提交 16 项投入，0 可以提交，提交后不能修改。</span>
              </div>
              <button
                type="button"
                class="btn primary"
                :disabled="!currentView.canSubmitInvestment || submittingInvestment"
                @click="handleSubmitInvestments"
              >
                {{ currentView.investmentSubmitted ? '已提交' : submittingInvestment ? '提交中...' : '提交 16 项投入' }}
              </button>
            </div>
            <div class="investment-grid">
              <div class="investment-grid-head">市场</div>
              <div v-for="orderType in orderTypeOptions" :key="orderType.code" class="investment-grid-head">{{ orderType.name }}</div>
              <template v-for="market in marketOptions" :key="market.code">
                <div class="investment-market-name">
                  {{ market.name }}
                  <small v-if="isMarketDisabled(market.code)">未开启</small>
                </div>
                <label v-for="orderType in orderTypeOptions" :key="`${market.code}-${orderType.code}`" class="investment-cell">
                  <input
                    :value="store.investmentDraft[investmentKey(market.code, orderType.code)] ?? 0"
                    type="number"
                    min="0"
                    step="0.01"
                    :disabled="!currentView.canSubmitInvestment || submittingInvestment"
                    @input="handleInvestmentInput(market.code, orderType.code, $event)"
                  >
                </label>
              </template>
            </div>
          </section>

          <section class="market-strip">
            <button
              v-for="market in markets"
              :key="market.marketCode"
              type="button"
              class="market-tab"
              :class="{ active: market.marketCode === selectedMarketCode }"
              @click="store.setSelectedMarket(market.marketCode)"
            >
              <strong>{{ market.marketName }}</strong>
              <span>{{ formatSegmentStatus(market.marketBidStatus) }}</span>
            </button>
          </section>

          <section v-if="selectedMarket" class="market-layout">
            <aside class="market-panel">
              <div class="panel-head compact">
                <div>
                  <strong>{{ selectedMarket.marketName }}状态</strong>
                  <span>{{ selectedMarket.isMarketLeader ? '本组为该市场龙头' : '按标段投入决定参与资格' }}</span>
                </div>
              </div>
              <div class="investment-box">
                <p class="hint">{{ selectedMarket.investmentSubmitted ? `该市场 4 个标段投入合计 ${formatAmount(selectedMarket.marketInvestment)}` : '等待年度 16 项投入提交。' }}</p>
              </div>

              <div class="sequence-box">
                <strong>当前市场标段</strong>
                <button
                  v-for="segment in selectedMarket.segments"
                  :key="`${segment.marketCode}-${segment.orderType}`"
                  type="button"
                  class="segment-row"
                  :class="{ current: segment.segmentStatus === 'SELECTING' }"
                >
                  <span>#{{ segment.releaseSequenceNo }} {{ segment.orderTypeName }}</span>
                  <em>{{ formatSegmentStatus(segment.segmentStatus) }}</em>
                </button>
              </div>
            </aside>

            <section class="segment-panel">
              <div class="panel-head">
                <div>
                  <strong>{{ currentSegmentTitle }}</strong>
                  <span>当前释放标段按顺序逐组选择；已被选择的订单会置灰。</span>
                </div>
                <button type="button" class="btn" :disabled="yearViewLoading" @click="handleRefresh">
                  刷新
                </button>
              </div>

              <section v-if="!visibleSegment" class="empty-state nested">
                <strong>当前市场暂无正在选择的标段</strong>
                <span>等待管理员开放投入、关闭投入并释放对应标段。</span>
              </section>

              <template v-else>
                <div class="segment-summary">
                  <span>释放顺序 #{{ visibleSegment.releaseSequenceNo }}</span>
                  <span>{{ formatSegmentStatus(visibleSegment.segmentStatus) }}</span>
                  <span>可选 {{ visibleSegment.availableOrders.length }} 单</span>
                  <span>已锁定 {{ visibleSegment.lockedOrders.length }} 单</span>
                </div>

                <div class="sequence-list">
                  <article
                    v-for="item in visibleSegment.selectionOrder"
                    :key="item.groupId"
                    class="sequence-card"
                    :class="{ self: item.isSelf, current: item.selectionStatus === 'CURRENT' }"
                  >
                    <span>#{{ item.sequenceNo }}</span>
                    <strong>{{ item.groupName }}</strong>
                    <em>{{ formatSelectionStatus(item.selectionStatus) }}</em>
                    <small v-if="item.isMarketLeader">市场龙头</small>
                  </article>
                </div>

                <div v-if="visibleSegment.selectedOrder" class="selected-order">
                  <div>
                    <strong>本组已选订单</strong>
                    <span>{{ formatOrderNo(visibleSegment.selectedOrder) }} · 金额 {{ formatIntegerAmount(visibleSegment.selectedOrder.orderAmount) }} · 账期 {{ visibleSegment.selectedOrder.accountTerm }} 季度 · {{ formatDeliveryStatus(visibleSegment.deliveryStatus) }}</span>
                  </div>
                  <div v-if="visibleSegment.deliveryStatus === 'SELECTED'" class="delivery-actions">
                    <select v-model="deliveryStageDraft[visibleSegment.selectedOrder.orderId]">
                      <option value="Q1">第一季度</option>
                      <option value="Q2">第二季度</option>
                      <option value="Q3">第三季度</option>
                      <option value="Q4">第四季度</option>
                    </select>
                    <button
                      type="button"
                      class="btn primary"
                      :disabled="deliveringOrders"
                      @click="handleDeliverOrder(visibleSegment)"
                    >
                      {{ deliveringOrders ? '交付中...' : '交付订单' }}
                    </button>
                  </div>
                  <em v-else-if="visibleSegment.selectedOrder.deliveredStageCode" class="delivery-note">
                    已在 {{ formatDeliveryStage(visibleSegment.selectedOrder.deliveredStageCode) }} 交付
                  </em>
                </div>

                <section class="order-grid">
                  <article
                    v-for="order in visibleSegment.availableOrders"
                    :key="order.orderId"
                    class="order-card"
                  >
                    <div class="order-title">
                      <strong>{{ formatOrderNo(order) }}</strong>
                      <span>可选</span>
                    </div>
                    <dl>
                      <div><dt>金额</dt><dd>{{ formatIntegerAmount(order.orderAmount) }}</dd></div>
                      <div><dt>数量</dt><dd>{{ formatQuantity(order.orderQuantity) }}</dd></div>
                      <div><dt>单价</dt><dd>{{ formatUnitPrice(order.unitPrice) }}</dd></div>
                      <div><dt>账期</dt><dd>{{ order.accountTerm }} 季度</dd></div>
                    </dl>
                    <button
                      type="button"
                      class="btn primary full"
                      :disabled="!visibleSegment.canSelectOrder || selectingOrder"
                      @click="handleSelectOrder(visibleSegment, order.orderId)"
                    >
                      {{ selectingOrder ? '选择中...' : '选择订单' }}
                    </button>
                  </article>
                  <article
                    v-for="order in visibleSegment.lockedOrders"
                    :key="order.orderId"
                    class="order-card locked"
                  >
                    <div class="order-title">
                      <strong>{{ formatOrderNo(order) }}</strong>
                      <span>已锁定</span>
                    </div>
                    <dl>
                      <div><dt>金额</dt><dd>{{ formatIntegerAmount(order.orderAmount) }}</dd></div>
                      <div><dt>数量</dt><dd>{{ formatQuantity(order.orderQuantity) }}</dd></div>
                      <div><dt>单价</dt><dd>{{ formatUnitPrice(order.unitPrice) }}</dd></div>
                      <div><dt>账期</dt><dd>{{ order.accountTerm }} 季度</dd></div>
                    </dl>
                    <button type="button" class="btn full" disabled>不可选择</button>
                  </article>
                </section>

                <div class="action-line">
                  <button
                    type="button"
                    class="btn danger"
                    :disabled="!visibleSegment.canPassSegment || passingSegment"
                    @click="handlePassSegment(visibleSegment)"
                  >
                    {{ passingSegment ? '放弃中...' : '放弃本标段' }}
                  </button>
                  <span>{{ visibleSegment.canPassSegment ? '当前轮到本组，可选择或放弃。' : '未轮到本组时只能查看。' }}</span>
                </div>
              </template>
            </section>

            <section class="delivery-panel">
              <div class="panel-head">
                <div>
                  <strong>本组待交付订单</strong>
                  <span>同一季度可勾选多个完整订单一次交付。</span>
                </div>
              </div>
              <section v-if="pendingDeliveryOrders.length === 0" class="empty-state nested">
                <strong>暂无待交付订单</strong>
                <span>已选订单会在这里集中展示。</span>
              </section>
              <template v-else>
                <div class="delivery-toolbar">
                  <select v-model="bulkDeliveryStage">
                    <option value="Q1">第一季度</option>
                    <option value="Q2">第二季度</option>
                    <option value="Q3">第三季度</option>
                    <option value="Q4">第四季度</option>
                  </select>
                  <button
                    type="button"
                    class="btn primary"
                    :disabled="selectedDeliveryOrderIds.length === 0 || deliveringOrders"
                    @click="handleBulkDeliverOrders"
                  >
                    {{ deliveringOrders ? '交付中...' : `交付 ${selectedDeliveryOrderIds.length} 单` }}
                  </button>
                </div>
                <div class="delivery-list">
                  <label
                    v-for="item in pendingDeliveryOrders"
                    :key="item.orderId"
                    class="delivery-row"
                  >
                    <input v-model="selectedDeliveryOrderIds" type="checkbox" :value="item.orderId">
                    <span>{{ item.marketName }} · {{ item.orderTypeName }}</span>
                    <strong>{{ item.businessOrderNo || `#${item.orderId}` }}</strong>
                    <em>{{ formatIntegerAmount(item.orderAmount) }}</em>
                  </label>
                </div>
              </template>
            </section>
          </section>
        </template>

        <section v-else class="empty-state">
          <strong>当前没有可展示的订单数据</strong>
          <span>请确认管理员已生成本年度订单池。</span>
        </section>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'

import YearTabs from '@/components/sandbox-game/common/YearTabs.vue'
import PageModeSwitch from '@/components/sandbox-game/player/PageModeSwitch.vue'
import { useAuthStore } from '@/stores/auth'
import { investmentKey, PLAYER_ORDER_MARKETS, PLAYER_ORDER_TYPES, usePlayerOrderStore } from '@/stores/player-order'
import type { OrderDeliveryStageCode, OrderMarketCode, OrderMarketForecastMarket, OrderTypeCode, PlayerOrderPoolItem, PlayerOrderSegmentView } from '@/types/sandbox-game-order'

const route = useRoute()
const router = useRouter()
const store = usePlayerOrderStore()
const authStore = useAuthStore()
const {
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
  markets,
  selectedMarket,
  currentSegment,
} = storeToRefs(store)
const { currentUser } = storeToRefs(authStore)

let initialized = false
let pollTimer = 0
const deliveryStageDraft = ref<Record<number, OrderDeliveryStageCode>>({})
const bulkDeliveryStage = ref<OrderDeliveryStageCode>('Q1')
const selectedDeliveryOrderIds = ref<number[]>([])
const marketOptions = PLAYER_ORDER_MARKETS
const orderTypeOptions = PLAYER_ORDER_TYPES

const visibleSegment = computed(() => currentSegment.value ?? selectedMarket.value?.segments[0] ?? null)
const activeSegmentTitle = computed(() => {
  const segment = markets.value.flatMap((market) => market.segments).find((item) => item.segmentStatus === 'SELECTING')
  return segment ? `${segment.marketName} ${segment.orderTypeName}` : '暂无'
})
const currentSegmentTitle = computed(() => {
  if (!visibleSegment.value) {
    return '标段选单'
  }
  return `${visibleSegment.value.marketName} · ${visibleSegment.value.orderTypeName}`
})
const pollingText = computed(() => {
  const seconds = currentView.value?.pollingIntervalSeconds ?? 3
  return `${seconds} 秒`
})
const pendingDeliveryOrders = computed(() =>
  markets.value.flatMap((market) =>
    market.segments
      .filter((segment) => segment.selectedOrder && segment.deliveryStatus === 'SELECTED')
      .map((segment) => ({
        marketName: segment.marketName,
        orderTypeName: segment.orderTypeName,
        orderId: segment.selectedOrder!.orderId,
        businessOrderNo: segment.selectedOrder!.businessOrderNo,
        orderAmount: segment.selectedOrder!.orderAmount,
      })),
  ),
)

onMounted(async () => {
  await store.bootstrap(readRouteYear())
  initialized = true
  if (selectedYear.value !== readRouteYear()) {
    syncRouteYear(selectedYear.value)
  }
  startPolling()
})

watch(
  () => route.query.yearNo,
  async () => {
    const targetYear = readRouteYear()
    if (!initialized || targetYear === undefined || targetYear === selectedYear.value) {
      return
    }
    try {
      await store.loadYearView(targetYear)
      restartPolling()
    } catch {
      syncRouteYear(selectedYear.value)
    }
  },
)

onBeforeUnmount(() => {
  stopPolling()
})

function readRouteYear() {
  const raw = Array.isArray(route.query.yearNo) ? route.query.yearNo[0] : route.query.yearNo
  if (!raw) {
    return undefined
  }
  const parsed = Number(raw)
  return Number.isFinite(parsed) ? parsed : undefined
}

function syncRouteYear(yearNo: number) {
  router.replace({
    path: route.path,
    query: { ...route.query, yearNo: String(yearNo) },
  })
}

function handleYearSelect(yearNo: number) {
  if (yearNo === selectedYear.value) {
    return
  }
  syncRouteYear(yearNo)
}

async function handleRefresh() {
  try {
    await store.loadYearView(selectedYear.value)
  } catch {
    // 页面消息由 store 统一处理。
  }
}

function handleInvestmentInput(marketCode: OrderMarketCode, orderType: OrderTypeCode, event: Event) {
  const input = event.target as HTMLInputElement
  const next = input.value === '' ? 0 : Number(input.value)
  store.setInvestmentDraft(marketCode, orderType, Number.isFinite(next) ? Math.max(next, 0) : 0)
}

function isMarketDisabled(marketCode: OrderMarketCode) {
  const market = markets.value.find((item) => item.marketCode === marketCode)
  return market?.marketEnabled === false
}

async function handleSubmitInvestments() {
  try {
    await store.submitInvestments()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleSelectOrder(segment: PlayerOrderSegmentView, orderId: number) {
  try {
    await store.selectOrder(segment, orderId)
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handlePassSegment(segment: PlayerOrderSegmentView) {
  try {
    await store.passSegment(segment)
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleDeliverOrder(segment: PlayerOrderSegmentView) {
  if (!segment.selectedOrder) {
    return
  }
  const stageCode = deliveryStageDraft.value[segment.selectedOrder.orderId] ?? 'Q1'
  try {
    await store.deliverSelectedOrder(segment, stageCode)
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleBulkDeliverOrders() {
  try {
    await store.deliverOrders(selectedDeliveryOrderIds.value, bulkDeliveryStage.value)
    selectedDeliveryOrderIds.value = []
  } catch {
    // 页面消息由 store 统一处理。
  }
}

function startPolling() {
  stopPolling()
  const seconds = currentView.value?.pollingIntervalSeconds ?? 3
  pollTimer = window.setInterval(async () => {
    if (loading.value || yearViewLoading.value || submittingInvestment.value || selectingOrder.value || passingSegment.value || deliveringOrders.value) {
      return
    }
    try {
      const scrollY = window.scrollY
      await store.loadYearView(selectedYear.value, { silent: true })
      await nextTick()
      window.scrollTo({ top: scrollY, behavior: 'auto' })
    } catch {
      // 自动轮询失败时保留当前页面状态。
    }
  }, Math.max(seconds, 1) * 1000)
}

function restartPolling() {
  startPolling()
}

function stopPolling() {
  if (pollTimer) {
    window.clearInterval(pollTimer)
    pollTimer = 0
  }
}

function goOperating() {
  router.push({
    path: '/sandbox-game/player/operating',
    query: { yearNo: String(selectedYear.value) },
  })
}

function goReport() {
  router.push({
    path: '/sandbox-game/player/report',
    query: { yearNo: String(selectedYear.value) },
  })
}

async function handleLogout() {
  try {
    await authStore.logout()
  } finally {
    await router.replace('/sandbox-game/login')
  }
}

function formatAmount(value: number) {
  return Number(value || 0).toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

function formatIntegerAmount(value: number) {
  return Math.round(Number(value || 0)).toLocaleString('zh-CN', { maximumFractionDigits: 0 })
}

function formatQuantity(value: number) {
  return Math.round(Number(value || 0)).toLocaleString('zh-CN', { maximumFractionDigits: 0 })
}

function formatUnitPrice(value: number) {
  return Number(value || 0).toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

function stageMarketAmount(market: OrderMarketForecastMarket) {
  return market.years.reduce((sum, year) => sum + Number(year.totalForecastAmount || 0), 0)
}

function stageMarketMaxAmount(market: OrderMarketForecastMarket) {
  return Math.max(
    ...market.years.flatMap((year) => year.products.map((product) => Number(product.forecastAmount || 0))),
    0,
  )
}

function forecastBarHeight(amount: number, market: OrderMarketForecastMarket) {
  const max = stageMarketMaxAmount(market)
  if (max <= 0) {
    return 0
  }
  return Math.max(4, Math.round((Number(amount || 0) / max) * 100))
}

function formatOrderNo(order: PlayerOrderPoolItem) {
  return order.businessOrderNo || `#${order.orderId}`
}

function formatSegmentStatus(value?: string) {
  const map: Record<string, string> = {
    MARKET_DISABLED: '市场未开启',
    NO_ORDER_CONFIG: '未配置订单',
    BID_OPEN: '投入开放',
    WAITING_INVESTMENT: '等待投入',
    BID_CLOSED: '投入关闭',
    SEQUENCE_READY: '顺序已生成',
    WAITING_RELEASE: '等待释放',
    SELECTING: '选单中',
    COMPLETED: '已完成',
    SKIPPED: '已跳过',
  }
  return value ? map[value] ?? value : '未开始'
}

function formatSelectionStatus(value: string) {
  const map: Record<string, string> = {
    INELIGIBLE: '无资格',
    WAITING: '待选择',
    CURRENT: '当前选择',
    SELECTED: '已选择',
    PASSED: '已放弃',
    ADMIN_SKIPPED: '管理员跳过',
  }
  return map[value] ?? value
}

function formatDeliveryStatus(value: string) {
  const map: Record<string, string> = {
    SELECTED: '待交付',
    DELIVERED: '已交付',
    UNFINISHED: '未完成',
  }
  return value ? map[value] ?? value : '待交付'
}

function formatDeliveryStage(value: string) {
  const map: Record<string, string> = {
    Q1: '第一季度',
    Q2: '第二季度',
    Q3: '第三季度',
    Q4: '第四季度',
  }
  return map[value] ?? value
}
</script>

<style scoped>
.player-page {
  min-height: 100vh;
  padding: 20px;
}

.shell {
  width: min(1760px, calc(100vw - 24px));
  margin: 0 auto;
  background: var(--shell-bg);
  border: 1px solid #dbe2ea;
  border-radius: 22px;
  box-shadow: var(--shadow);
  overflow: hidden;
}

.page-header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 22px 24px 18px;
  background: linear-gradient(180deg, #fbfdff 0%, #eef4fb 100%);
  border-bottom: 1px solid #dde5ef;
}

.page-header h1 {
  margin: 6px 0 8px;
  font-size: 28px;
}

.eyebrow {
  margin: 0;
  color: var(--accent);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.subtext {
  margin: 0;
  color: var(--muted);
}

.header-pills {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  align-content: flex-start;
  gap: 10px;
}

.pill {
  display: inline-flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: #ffffff;
  font-size: 13px;
}

.logout-button,
.btn {
  border: 1px solid var(--line);
  background: #ffffff;
  color: var(--text);
}

.logout-button {
  height: 36px;
  padding: 0 14px;
  border-radius: 999px;
}

.toolbar-card {
  padding: 16px 18px;
  border-bottom: 1px solid var(--line);
  background: #f8fafc;
}

.toolbar-top {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
}

.message-bar {
  margin: 16px 18px 0;
  padding: 12px 14px;
  border-radius: 14px;
  font-size: 14px;
}

.message-bar.success {
  background: #edfdf3;
  color: var(--success);
  border: 1px solid #b7e2c5;
}

.message-bar.error {
  background: #fff5f5;
  color: var(--danger);
  border: 1px solid #efc4c4;
}

.message-bar.info {
  background: var(--accent-soft);
  color: var(--accent);
  border: 1px solid #cbdcff;
}

.order-workspace {
  display: grid;
  gap: 16px;
  padding: 16px 18px 20px;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.status-card,
.market-panel,
.segment-panel,
.empty-state,
.loading-card {
  border: 1px solid var(--line);
  border-radius: 16px;
  background: #ffffff;
}

.status-card {
  display: grid;
  gap: 6px;
  padding: 14px 16px;
}

.status-card span,
.hint,
.panel-head span,
.action-line span {
  color: var(--muted);
  font-size: 13px;
}

.status-card strong {
  font-size: 18px;
}

.market-strip {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.investment-panel {
  border: 1px solid var(--line);
  border-radius: 16px;
  background: #ffffff;
  overflow: hidden;
}

.forecast-panel {
  border: 1px solid var(--line);
  border-radius: 16px;
  background: #ffffff;
  overflow: hidden;
}

.formula-pill {
  border: 1px solid var(--line);
  border-radius: 999px;
  background: #ffffff;
  padding: 7px 10px;
  color: var(--muted);
  font-size: 12px;
}

.forecast-stage-grid {
  display: grid;
  gap: 14px;
  padding: 16px;
}

.forecast-stage-card {
  border: 1px solid var(--line);
  border-radius: 14px;
  background: #fbfcfe;
  overflow: hidden;
}

.forecast-stage-card header {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  padding: 11px 13px;
  border-bottom: 1px solid var(--line);
}

.forecast-stage-card header span {
  color: var(--muted);
  font-size: 13px;
}

.forecast-market-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
  padding: 12px;
}

.forecast-market-card {
  display: grid;
  gap: 10px;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #ffffff;
  padding: 11px 12px;
}

.forecast-market-head {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.forecast-market-head span {
  color: var(--success);
  font-weight: 700;
}

.forecast-market-card p {
  margin: 0;
  color: var(--muted);
  font-size: 12px;
  line-height: 1.5;
}

.forecast-chart {
  display: grid;
  grid-template-columns: repeat(3, minmax(74px, 1fr));
  gap: 8px;
  align-items: end;
  min-height: 172px;
  border: 1px solid var(--line);
  border-radius: 10px;
  background:
    linear-gradient(to top, rgba(148, 163, 184, 0.18) 1px, transparent 1px) 0 0 / 100% 25%,
    #fbfcfe;
  padding: 12px 10px 10px;
}

.forecast-year-group {
  display: grid;
  grid-template-rows: 112px auto auto;
  gap: 5px;
  min-width: 0;
  text-align: center;
}

.forecast-bars {
  display: flex;
  gap: 4px;
  align-items: end;
  justify-content: center;
  height: 112px;
  border-bottom: 1px solid #cbd5e1;
}

.forecast-bar {
  width: 10px;
  min-height: 0;
  border-radius: 3px 3px 0 0;
  background: #3b82f6;
}

.forecast-year-label,
.forecast-year-group em {
  overflow: hidden;
  color: var(--muted);
  font-size: 11px;
  font-style: normal;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.forecast-year-label {
  color: var(--text);
  font-weight: 700;
}

.forecast-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 10px;
}

.forecast-legend span {
  display: inline-flex;
  gap: 5px;
  align-items: center;
  color: var(--muted);
  font-size: 11px;
}

.forecast-legend span::before {
  width: 8px;
  height: 8px;
  border-radius: 2px;
  background: currentColor;
  content: "";
}

.type-AGENCY_INSPECTION {
  background-color: #2563eb;
  color: #2563eb;
}

.type-TWO_CABIN_VIP {
  background-color: #16a34a;
  color: #16a34a;
}

.type-BUSINESS_VIP {
  background-color: #d97706;
  color: #d97706;
}

.type-MEMBER_CUSTOM {
  background-color: #7c3aed;
  color: #7c3aed;
}

.investment-grid {
  display: grid;
  grid-template-columns: 120px repeat(4, minmax(120px, 1fr));
  gap: 0;
  padding: 16px;
}

.investment-grid-head,
.investment-market-name,
.investment-cell {
  min-height: 48px;
  border: 1px solid var(--line);
  margin: -1px 0 0 -1px;
}

.investment-grid-head,
.investment-market-name {
  display: flex;
  gap: 6px;
  align-items: center;
  padding: 10px 12px;
  background: #f7f9fc;
  font-weight: 700;
}

.investment-market-name small {
  color: var(--muted);
  font-size: 12px;
  font-weight: 400;
}

.investment-grid-head {
  justify-content: center;
  color: var(--muted);
  font-size: 13px;
}

.investment-cell {
  display: flex;
  align-items: center;
  padding: 7px;
  background: #ffffff;
}

.investment-cell input {
  width: 100%;
  height: 34px;
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 6px 8px;
}

.market-tab {
  display: grid;
  gap: 5px;
  text-align: left;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: #ffffff;
  padding: 13px 14px;
}

.market-tab.active {
  border-color: var(--accent);
  background: var(--accent-soft);
}

.market-tab span {
  color: var(--muted);
  font-size: 12px;
}

.market-layout {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  gap: 14px;
  align-items: start;
}

.panel-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  padding: 14px 16px;
  border-bottom: 1px solid var(--line);
  background: #f7f9fc;
}

.panel-head.compact {
  display: block;
}

.panel-head strong {
  display: block;
  margin-bottom: 4px;
  font-size: 16px;
}

.investment-box,
.sequence-box {
  display: grid;
  gap: 12px;
  padding: 16px;
}

.investment-box label {
  display: grid;
  gap: 6px;
}

.investment-box label span {
  color: var(--muted);
  font-size: 13px;
  font-weight: 700;
}

.investment-box input {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 10px 12px;
}

.btn {
  border-radius: 12px;
  padding: 10px 14px;
}

.btn.primary {
  background: var(--accent);
  color: #ffffff;
  border-color: var(--accent);
}

.btn.danger {
  background: #fff5f5;
  color: var(--danger);
  border-color: #efc4c4;
}

.btn.full {
  width: 100%;
}

.sequence-box {
  border-top: 1px solid var(--line);
}

.segment-row {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #fbfcfe;
  padding: 10px 12px;
  text-align: left;
}

.segment-row.current {
  border-color: var(--success);
  background: #eefaf2;
}

.segment-row em {
  color: var(--muted);
  font-style: normal;
  font-size: 12px;
  white-space: nowrap;
}

.segment-summary,
.action-line {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  align-items: center;
  padding: 14px 16px;
}

.segment-summary span {
  border: 1px solid var(--line);
  border-radius: 999px;
  background: #ffffff;
  padding: 7px 10px;
  font-size: 13px;
}

.sequence-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 10px;
  padding: 0 16px 16px;
}

.sequence-card {
  display: grid;
  gap: 4px;
  min-height: 92px;
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 10px 12px;
  background: #ffffff;
}

.sequence-card.current {
  border-color: var(--accent);
  background: var(--accent-soft);
}

.sequence-card.self {
  box-shadow: inset 3px 0 0 var(--accent);
}

.sequence-card span,
.sequence-card em,
.sequence-card small {
  color: var(--muted);
  font-size: 12px;
  font-style: normal;
}

.selected-order {
  display: grid;
  gap: 10px;
  margin: 0 16px 16px;
  border: 1px solid #b7e2c5;
  border-radius: 14px;
  background: #edfdf3;
  padding: 12px 14px;
}

.selected-order span {
  display: block;
  margin-top: 4px;
}

.delivery-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}

.delivery-actions select {
  min-width: 128px;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #ffffff;
  padding: 9px 10px;
}

.delivery-note {
  color: var(--success);
  font-style: normal;
  font-size: 13px;
}

.delivery-panel {
  grid-column: 1 / -1;
  border: 1px solid var(--line);
  border-radius: 16px;
  background: #ffffff;
}

.delivery-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
  padding: 14px 16px;
}

.delivery-toolbar select {
  min-width: 128px;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #ffffff;
  padding: 9px 10px;
}

.delivery-list {
  display: grid;
  gap: 10px;
  padding: 0 16px 16px;
}

.delivery-row {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto auto;
  gap: 10px;
  align-items: center;
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 10px 12px;
}

.delivery-row em {
  color: var(--success);
  font-style: normal;
  font-weight: 700;
}

.order-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 12px;
  padding: 0 16px 16px;
}

.order-card {
  display: grid;
  gap: 12px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: #ffffff;
  padding: 14px;
}

.order-card.locked {
  opacity: 0.58;
  background: #f4f6f9;
}

.order-title {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: center;
}

.order-title span {
  border-radius: 999px;
  background: #edfdf3;
  color: var(--success);
  padding: 4px 8px;
  font-size: 12px;
}

.locked .order-title span {
  background: #e2e8f0;
  color: var(--muted);
}

.order-card dl {
  display: grid;
  gap: 7px;
  margin: 0;
}

.order-card dl div {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.order-card dt {
  color: var(--muted);
}

.order-card dd {
  margin: 0;
  font-weight: 700;
}

.empty-state,
.loading-card {
  display: grid;
  gap: 8px;
  padding: 28px;
  text-align: center;
  color: var(--muted);
}

.empty-state strong {
  color: var(--text);
}

.empty-state.nested {
  margin: 16px;
}

@media (max-width: 1280px) {
  .market-layout {
    grid-template-columns: 1fr;
  }

  .status-grid,
  .market-strip,
  .forecast-market-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .investment-grid {
    overflow-x: auto;
    grid-template-columns: 110px repeat(4, 130px);
  }
}

@media (max-width: 900px) {
  .player-page {
    padding: 12px;
  }

  .shell {
    width: 100%;
  }

  .page-header,
  .toolbar-top {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-pills {
    justify-content: flex-start;
  }

  .status-grid,
  .market-strip,
  .forecast-market-grid {
    grid-template-columns: 1fr;
  }
}
</style>
