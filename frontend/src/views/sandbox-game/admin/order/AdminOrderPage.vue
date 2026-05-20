<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>订单管理</h2>
        <p>管理员按订单生成规则配置标段数量，生成预览订单池，确认后再生成选单顺序并逐段释放。</p>
      </div>
      <div class="hero-actions">
        <select v-model.number="store.selectedYearNo" class="year-select" :disabled="loading || savingConfig || generatingPool" @change="handleYearChange">
          <option v-for="item in yearOptions" :key="item.value" :value="item.value">{{ item.label }}</option>
        </select>
        <button type="button" class="btn" :disabled="loading" @click="handleRefresh">刷新</button>
        <button type="button" class="btn primary" :disabled="savingConfig" @click="handleSaveConfig">
          {{ savingConfig ? '保存中...' : '保存配置' }}
        </button>
      </div>
    </header>

    <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
      {{ pageMessage.text }}
    </section>

    <section class="stats-grid">
      <article class="stat-card">
        <span>当前年份</span>
        <strong>{{ selectedYearNo }}年</strong>
      </article>
      <article class="stat-card">
        <span>生成状态</span>
        <strong>{{ formatGenerationStatus(config?.generationStatus) }}</strong>
      </article>
      <article class="stat-card">
        <span>配置订单数</span>
        <strong>{{ totalOrderCount }}</strong>
      </article>
      <article class="stat-card">
        <span>已生成订单</span>
        <strong>{{ totalGeneratedCount }}</strong>
      </article>
    </section>

    <section class="panel-card generation-panel">
      <div class="panel-head">
        <div>
          <strong>订单生成控制台</strong>
          <span>先保存 16 个标段数量与释放顺序，再生成预览；确认后订单池和数量配置锁定。</span>
        </div>
      </div>
      <div class="generation-actions">
        <button type="button" class="btn" :disabled="generatingPool || !config?.canGeneratePreview" @click="handleGeneratePool">
          {{ generatingPool ? '生成中...' : '生成/覆盖预览订单池' }}
        </button>
        <button type="button" class="btn primary" :disabled="confirmingPool || !config?.canConfirmPool" @click="handleConfirmPool">
          {{ confirmingPool ? '确认中...' : '确认订单池' }}
        </button>
        <button type="button" class="btn primary" :disabled="generatingSequence" @click="handleGenerateSelectionSequence">
          {{ generatingSequence ? '生成中...' : '生成选单顺序' }}
        </button>
      </div>
      <div class="batch-grid">
        <article class="batch-card">
          <span>预览批次</span>
          <strong>{{ config?.latestPreviewBatch ? `#${config.latestPreviewBatch.batchId}` : '--' }}</strong>
          <em>{{ config?.latestPreviewBatch ? `${config.latestPreviewBatch.generatedCount} 单 · ${formatDateTime(config.latestPreviewBatch.generatedAt)}` : '尚未生成预览' }}</em>
        </article>
        <article class="batch-card">
          <span>确认批次</span>
          <strong>{{ config?.confirmedBatch ? `#${config.confirmedBatch.batchId}` : '--' }}</strong>
          <em>{{ config?.confirmedBatch?.confirmedAt ? `${config.confirmedBatch.generatedCount} 单 · ${formatDateTime(config.confirmedBatch.confirmedAt)}` : '尚未确认' }}</em>
        </article>
        <article class="batch-card">
          <span>公式版本</span>
          <strong>{{ config?.latestPreviewBatch?.formulaVersion ?? config?.confirmedBatch?.formulaVersion ?? '--' }}</strong>
          <em>{{ config?.latestPreviewBatch?.randomSeed ? `Seed ${config.latestPreviewBatch.randomSeed}` : '确认后固化随机种子' }}</em>
        </article>
      </div>
    </section>

    <section v-if="config?.warnings.length" class="warning-list">
      <strong>数量风险提示</strong>
      <span v-for="warning in config.warnings" :key="`${warning.level}-${warning.message}`">{{ warning.message }}</span>
    </section>

    <section class="panel-card control-panel">
      <div class="panel-head">
        <div>
          <strong>市场竞标控制</strong>
          <span>玩家提交完整 16 项投入后，管理员生成全部标段选单顺序，再按释放顺序逐个释放。</span>
        </div>
        <button type="button" class="btn" :disabled="loadingSelectionStatus" @click="handleLoadSelectionStatus">
          {{ loadingSelectionStatus ? '加载中...' : '刷新状态' }}
        </button>
      </div>

      <div class="control-grid">
        <label class="field">
          <span>控制市场</span>
          <select v-model="store.controlForm.marketCode" @change="handleControlMarketChange">
            <option v-for="item in marketOptions" :key="item.code" :value="item.code">{{ item.name }}</option>
          </select>
        </label>
        <div class="control-actions">
          <button type="button" class="btn primary" :disabled="releasingSegment" @click="handleReleaseNextSegment">
            {{ releasingSegment ? '释放中...' : '释放下一个标段' }}
          </button>
        </div>
      </div>

      <div class="selection-overview">
        <article class="status-mini">
          <span>市场状态</span>
          <strong>{{ formatSegmentStatus(marketSelectionStatus?.marketBidStatus) }}</strong>
        </article>
        <article class="status-mini">
          <span>市场龙头</span>
          <strong>{{ marketSelectionStatus?.leaderGroupId ? `组ID ${marketSelectionStatus.leaderGroupId}` : '--' }}</strong>
        </article>
        <article class="status-mini">
          <span>当前标段</span>
          <strong>{{ currentSegment ? `${currentSegment.marketName} ${currentSegment.orderTypeName}` : '--' }}</strong>
        </article>
        <article class="status-mini">
          <span>当前小组</span>
          <strong>{{ currentGroupName }}</strong>
        </article>
      </div>

      <div class="skip-row">
        <input v-model="store.controlForm.skipReason" type="text" placeholder="代跳过原因，必填">
        <button type="button" class="btn danger" :disabled="!currentSegment?.currentGroupId || skippingGroup" @click="handleSkipCurrentGroup">
          {{ skippingGroup ? '跳过中...' : '跳过当前小组' }}
        </button>
      </div>

      <div class="table-scroll">
        <table class="selection-table">
          <thead>
            <tr>
              <th>释放顺序</th>
              <th>市场</th>
              <th>订单类型</th>
              <th>标段状态</th>
              <th>当前组</th>
              <th>可选</th>
              <th>已选</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!marketSelectionStatus || marketSelectionStatus.segments.length === 0">
              <td colspan="7" class="empty-row">当前市场暂无标段状态，请先生成订单池。</td>
            </tr>
            <tr v-for="segment in marketSelectionStatus?.segments ?? []" :key="`${segment.marketCode}-${segment.orderType}`">
              <td>#{{ segment.releaseSequenceNo }}</td>
              <td>{{ segment.marketName }}</td>
              <td>{{ segment.orderTypeName }}</td>
              <td>{{ formatSegmentStatus(segment.segmentStatus) }}</td>
              <td>{{ formatGroupName(segment.currentGroupId) }}</td>
              <td class="number-cell">{{ segment.availableCount }}</td>
              <td class="number-cell">{{ segment.selectedCount }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-scroll sequence-scroll">
        <table class="selection-table">
          <thead>
            <tr>
              <th>顺序</th>
              <th>小组</th>
              <th>市场投入</th>
              <th>市场龙头</th>
              <th>状态</th>
              <th>已选订单</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!currentSegment || currentSegment.selectionOrder.length === 0">
              <td colspan="6" class="empty-row">当前没有正在选单的标段。</td>
            </tr>
            <tr v-for="item in currentSegment?.selectionOrder ?? []" :key="item.groupId">
              <td>#{{ item.sequenceNo }}</td>
              <td>{{ item.groupName }}</td>
              <td class="number-cell">{{ formatAmount(item.marketInvestment) }}</td>
              <td>{{ item.isMarketLeader ? '是' : '否' }}</td>
              <td>{{ formatSelectionStatus(item.selectionStatus) }}</td>
              <td>{{ item.selectedOrderId ? `#${item.selectedOrderId}` : '--' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section class="panel-card">
      <div class="panel-head">
        <div>
          <strong>标段数量与释放顺序</strong>
          <span>每个市场 + 订单类型单独开标；开标过程中释放顺序锁定。</span>
        </div>
      </div>

      <div class="table-scroll">
        <table class="config-table">
          <thead>
            <tr>
              <th>释放顺序</th>
              <th>市场</th>
              <th>订单类型</th>
              <th>订单数量</th>
              <th>Excel 可用</th>
              <th>已生成</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="sortedItems.length === 0">
              <td colspan="7" class="empty-row">暂无订单配置。</td>
            </tr>
            <tr v-for="item in sortedItems" :key="`${item.marketCode}-${item.orderType}`">
              <td>
                <input v-model.number="item.releaseSequenceNo" type="number" min="1" step="1" :disabled="savingConfig || generatingPool || !config?.canUpdateConfig" class="compact-input">
              </td>
              <td>{{ item.marketName }}</td>
              <td>{{ item.orderTypeName }}</td>
              <td>
                <input v-model.number="item.orderCount" type="number" min="0" max="15" step="1" :disabled="savingConfig || generatingPool || !config?.canUpdateConfig" class="compact-input">
              </td>
              <td class="number-cell">{{ item.availableCount }}</td>
              <td class="number-cell">{{ item.generatedCount }}</td>
              <td>
                <span class="status-tag" :class="item.configStatus === 'LOCKED' ? 'locked' : 'draft'">
                  {{ item.configStatus === 'LOCKED' ? '已锁定' : '草稿' }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section class="panel-card">
      <div class="panel-head">
        <div>
          <strong>订单池查看</strong>
          <span>按年份、市场和订单类型查看固定订单池状态。</span>
        </div>
        <button type="button" class="btn" :disabled="loadingPool" @click="handleLoadPool">
          {{ loadingPool ? '加载中...' : '查看订单池' }}
        </button>
      </div>

      <div class="pool-filter">
        <label class="field">
          <span>市场</span>
          <select v-model="store.poolFilter.marketCode" @change="handleLoadPool">
            <option v-for="item in marketOptions" :key="item.code" :value="item.code">{{ item.name }}</option>
          </select>
        </label>
        <label class="field">
          <span>订单类型</span>
          <select v-model="store.poolFilter.orderType" @change="handleLoadPool">
            <option v-for="item in orderTypeOptions" :key="item.code" :value="item.code">{{ item.name }}</option>
          </select>
        </label>
      </div>

      <div class="table-scroll">
        <table class="pool-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>金额</th>
              <th>数量</th>
              <th>单价</th>
              <th>账期</th>
              <th>状态</th>
              <th>选中组</th>
              <th>来源</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!orderPool || orderPool.list.length === 0">
              <td colspan="8" class="empty-row">当前筛选下没有订单池记录。</td>
            </tr>
            <tr v-for="item in orderPool?.list ?? []" :key="item.orderId">
              <td>#{{ item.orderId }}</td>
              <td class="number-cell">{{ formatAmount(item.orderAmount) }}</td>
              <td class="number-cell">{{ formatAmount(item.orderQuantity) }}</td>
              <td class="number-cell">{{ formatAmount(item.unitPrice) }}</td>
              <td>{{ item.accountTerm }}季度</td>
              <td>{{ formatPoolStatus(item.poolStatus) }}</td>
              <td>{{ item.selectedGroupId ? `组ID ${item.selectedGroupId}` : '--' }}</td>
              <td>{{ item.sourceSheetName }} {{ item.sourceCell }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'

import { useAdminShellStore } from '@/stores/admin-shell'
import { MARKET_OPTIONS, ORDER_TYPE_OPTIONS, useAdminOrderStore } from '@/stores/admin-order'
import type { OrderPoolStatus } from '@/types/sandbox-game-admin'

const shellStore = useAdminShellStore()
const store = useAdminOrderStore()

const { config: shellConfig } = storeToRefs(shellStore)
const {
  selectedYearNo,
  config,
  orderPool,
  loading,
  savingConfig,
  generatingPool,
  confirmingPool,
  generatingSequence,
  loadingPool,
  loadingSelectionStatus,
  releasingSegment,
  skippingGroup,
  pageMessage,
  marketSelectionStatus,
  currentSegment,
  sortedItems,
  totalOrderCount,
  totalGeneratedCount,
} = storeToRefs(store)

const marketOptions = MARKET_OPTIONS
const orderTypeOptions = ORDER_TYPE_OPTIONS

const yearOptions = computed(() => {
  const finalYear = Math.max(shellConfig.value?.finalYear ?? config.value?.finalYear ?? 1, 1)
  return Array.from({ length: finalYear }, (_, index) => ({
    value: index + 1,
    label: `${index + 1}年`,
  }))
})
const currentGroupName = computed(() => formatGroupName(currentSegment.value?.currentGroupId ?? null))

onMounted(async () => {
  try {
    if (!shellStore.config) {
      await shellStore.bootstrap()
    }
    const defaultYear = Math.max(shellConfig.value?.currentOpenYear ?? 1, 1)
    await store.bootstrap(defaultYear)
    await store.loadPool({ silent: true })
    await store.loadSelectionStatus({ silent: true })
  } catch {
    // 页面消息由 store 统一处理。
  }
})

async function handleYearChange() {
  try {
    await store.loadConfig()
    await store.loadPool({ silent: true })
    await store.loadSelectionStatus({ silent: true })
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleRefresh() {
  try {
    await store.loadConfig()
    await store.loadPool({ silent: true })
    await store.loadSelectionStatus({ silent: true })
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleSaveConfig() {
  try {
    await store.saveConfig()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleGeneratePool() {
  try {
    await store.generatePool(true)
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleConfirmPool() {
  try {
    await store.confirmPool()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleGenerateSelectionSequence() {
  try {
    await store.generateSelectionSequence()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleLoadPool() {
  try {
    await store.loadPool()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleLoadSelectionStatus() {
  try {
    await store.loadSelectionStatus()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleControlMarketChange() {
  try {
    await store.loadSelectionStatus()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleReleaseNextSegment() {
  try {
    await store.releaseNextSegment()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleSkipCurrentGroup() {
  try {
    await store.skipCurrentGroup()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

function formatDateTime(value?: string | null) {
  if (!value) {
    return '--'
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString('zh-CN', { hour12: false })
}

function formatAmount(value: number) {
  return Number(value || 0).toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

function formatPoolStatus(value: OrderPoolStatus) {
  if (value === 'AVAILABLE') {
    return '可选'
  }
  if (value === 'SELECTED') {
    return '已选'
  }
  if (value === 'VOID') {
    return '作废'
  }
  return value
}

function formatGenerationStatus(value?: string) {
  const map: Record<string, string> = {
    NOT_GENERATED: '未生成',
    PREVIEW_GENERATED: '预览已生成',
    POOL_CONFIRMED: '订单池已确认',
    SELECTING: '选单中',
    COMPLETED: '已完成',
  }
  return value ? map[value] ?? value : '--'
}

function formatSegmentStatus(value?: string) {
  const map: Record<string, string> = {
    WAITING_INVESTMENT: '等待投入',
    BID_OPEN: '投入开放',
    BID_CLOSED: '投入关闭',
    SEQUENCE_READY: '顺序已生成',
    WAITING_RELEASE: '等待释放',
    SELECTING: '选单中',
    COMPLETED: '已完成',
    SKIPPED: '已跳过',
  }
  return value ? map[value] ?? value : '--'
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

function formatGroupName(groupId?: number | null) {
  if (!groupId) {
    return '--'
  }
  const bid = marketSelectionStatus.value?.bids.find((item) => item.groupId === groupId)
  return bid?.groupName ?? `组ID ${groupId}`
}
</script>

<style scoped>
.page-content {
  display: grid;
  gap: 16px;
}

.hero {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
}

.hero h2 {
  margin: 0 0 6px;
  font-size: 24px;
}

.hero p {
  margin: 0;
  color: var(--muted);
  line-height: 1.6;
}

.hero-actions,
.generation-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  align-items: center;
}

.btn,
.year-select {
  border: 1px solid var(--line);
  background: #ffffff;
  border-radius: 12px;
  padding: 10px 14px;
}

.btn.primary {
  color: #ffffff;
  border-color: var(--accent);
  background: var(--accent);
}

.btn.danger {
  color: var(--danger);
  border-color: #efc4c4;
  background: #fff5f5;
}

.message-bar {
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.stat-card,
.panel-card {
  border: 1px solid var(--line);
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.05);
}

.stat-card {
  display: grid;
  gap: 8px;
  padding: 16px;
}

.stat-card span {
  color: var(--muted);
  font-size: 12px;
}

.stat-card strong {
  font-size: 20px;
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

.panel-head strong {
  display: block;
  margin-bottom: 4px;
  font-size: 16px;
}

.panel-head span,
.inline-tip {
  color: var(--muted);
  font-size: 13px;
}

.generation-actions,
.batch-grid,
.pool-filter {
  padding: 16px;
}

.batch-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  border-top: 1px solid var(--line);
}

.batch-card {
  display: grid;
  gap: 6px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: #fbfcfe;
  padding: 12px;
}

.batch-card span,
.batch-card em {
  color: var(--muted);
  font-size: 13px;
  font-style: normal;
}

.warning-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.warning-list span {
  padding: 7px 10px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: #f8fafc;
  font-size: 12px;
}

.warning-list {
  padding: 12px 14px;
  border: 1px solid #efd4aa;
  border-radius: 14px;
  background: #fffaf0;
}

.warning-list strong {
  width: 100%;
  color: var(--warning);
}

.warning-list.compact {
  margin-top: 12px;
  padding: 0;
  border: 0;
  background: transparent;
}

.table-scroll {
  overflow: auto;
}

.config-table,
.pool-table {
  width: 100%;
  min-width: 900px;
  border-collapse: collapse;
}

.config-table th,
.config-table td,
.pool-table th,
.pool-table td {
  border: 1px solid var(--line);
  padding: 10px 12px;
  font-size: 14px;
}

.config-table th,
.pool-table th {
  background: #f4f6f9;
  text-align: center;
}

.compact-input {
  width: 96px;
  border: 1px solid var(--line);
  border-radius: 10px;
  padding: 8px 10px;
}

.number-cell {
  text-align: right;
}

.empty-row {
  text-align: center;
  color: var(--muted);
}

.status-tag {
  display: inline-flex;
  padding: 5px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
}

.status-tag.draft {
  background: var(--accent-soft);
  color: var(--accent);
}

.status-tag.locked {
  background: #e8f7ee;
  color: var(--success);
}

.pool-filter {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 240px));
  gap: 12px;
}

.control-grid {
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr);
  gap: 12px;
  align-items: end;
  padding: 16px;
}

.control-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.selection-overview {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  padding: 0 16px 16px;
}

.status-mini {
  display: grid;
  gap: 6px;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #fbfcfe;
  padding: 12px;
}

.status-mini span {
  color: var(--muted);
  font-size: 12px;
}

.status-mini strong {
  font-size: 16px;
}

.skip-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 10px;
  padding: 0 16px 16px;
}

.skip-row input {
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 10px 12px;
}

.selection-table {
  width: 100%;
  min-width: 760px;
  border-collapse: collapse;
}

.selection-table th,
.selection-table td {
  border: 1px solid var(--line);
  padding: 10px 12px;
  font-size: 14px;
}

.selection-table th {
  background: #f4f6f9;
  text-align: center;
}

.sequence-scroll {
  border-top: 1px solid var(--line);
}

.field {
  display: grid;
  gap: 6px;
}

.field span {
  color: var(--muted);
  font-size: 13px;
  font-weight: 700;
}

.field select {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #ffffff;
  padding: 10px 12px;
}

@media (max-width: 1240px) {
  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .hero,
  .panel-head {
    flex-direction: column;
    align-items: flex-start;
  }

  .stats-grid,
  .batch-grid,
  .pool-filter,
  .control-grid,
  .selection-overview,
  .skip-row {
    grid-template-columns: 1fr;
  }
}
</style>
