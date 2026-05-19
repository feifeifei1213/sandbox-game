<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>订单管理</h2>
        <p>管理员在开标前维护订单 Excel、标段数量、释放顺序，并生成固定订单池。</p>
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
        <span>最新 Excel 批次</span>
        <strong>{{ config?.latestBatchId ? `#${config.latestBatchId}` : '未上传' }}</strong>
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

    <section class="panel-card upload-panel">
      <div class="panel-head">
        <div>
          <strong>订单 Excel 上传</strong>
          <span>上传后只解析订单卡片保存值，不直接覆盖已生成订单池。</span>
        </div>
      </div>
      <div class="upload-body">
        <label class="file-picker">
          <input type="file" accept=".xlsx" :disabled="uploading" @change="handleFileChange">
          <span>{{ selectedFileName || '选择订单推算 Excel' }}</span>
        </label>
        <button type="button" class="btn primary" :disabled="!selectedFile || uploading" @click="handleUpload">
          {{ uploading ? '上传解析中...' : '上传并解析' }}
        </button>
        <span class="inline-tip">{{ latestBatchText }}</span>
      </div>
      <div v-if="uploadResult" class="parse-preview">
        <div class="preview-head">
          <strong>解析预览</strong>
          <span>{{ uploadResult.originalFileName }} · {{ formatDateTime(uploadResult.uploadedAt) }}</span>
        </div>
        <div class="summary-chips">
          <span v-for="item in uploadResult.summary.slice(0, 12)" :key="`${item.yearNo}-${item.marketCode}-${item.orderType}`">
            {{ item.yearNo }}年 {{ item.marketName }} {{ item.orderTypeName }}：{{ item.availableCount }}
          </span>
        </div>
        <div v-if="uploadResult.warnings.length > 0" class="warning-list compact">
          <span v-for="warning in uploadResult.warnings" :key="warning">{{ warning }}</span>
        </div>
      </div>
    </section>

    <section v-if="config?.warnings.length" class="warning-list">
      <strong>数量风险提示</strong>
      <span v-for="warning in config.warnings" :key="`${warning.level}-${warning.message}`">{{ warning.message }}</span>
    </section>

    <section class="panel-card">
      <div class="panel-head">
        <div>
          <strong>标段数量与释放顺序</strong>
          <span>每个市场 + 订单类型单独开标；开标过程中释放顺序锁定。</span>
        </div>
        <button type="button" class="btn" :disabled="generatingPool || totalOrderCount <= 0" @click="handleGeneratePool">
          {{ generatingPool ? '生成中...' : '生成/覆盖订单池' }}
        </button>
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
                <input v-model.number="item.releaseSequenceNo" type="number" min="1" step="1" :disabled="savingConfig || generatingPool" class="compact-input">
              </td>
              <td>{{ item.marketName }}</td>
              <td>{{ item.orderTypeName }}</td>
              <td>
                <input v-model.number="item.orderCount" type="number" min="0" step="1" :disabled="savingConfig || generatingPool" class="compact-input">
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
import { computed, onMounted, ref } from 'vue'
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
  uploadResult,
  orderPool,
  loading,
  uploading,
  savingConfig,
  generatingPool,
  loadingPool,
  pageMessage,
  sortedItems,
  totalOrderCount,
  totalGeneratedCount,
} = storeToRefs(store)

const selectedFile = ref<File | null>(null)
const marketOptions = MARKET_OPTIONS
const orderTypeOptions = ORDER_TYPE_OPTIONS

const selectedFileName = computed(() => selectedFile.value?.name ?? '')
const yearOptions = computed(() => {
  const finalYear = Math.max(shellConfig.value?.finalYear ?? config.value?.finalYear ?? 1, 1)
  return Array.from({ length: finalYear }, (_, index) => ({
    value: index + 1,
    label: `${index + 1}年`,
  }))
})
const latestBatchText = computed(() => {
  if (!config.value?.latestBatchId) {
    return '当前还没有成功解析的订单 Excel。'
  }
  return `最新批次 #${config.value.latestBatchId}，上传时间 ${formatDateTime(config.value.latestBatchUploadedAt)}。`
})

onMounted(async () => {
  try {
    if (!shellStore.config) {
      await shellStore.bootstrap()
    }
    const defaultYear = Math.max(shellConfig.value?.currentOpenYear ?? 1, 1)
    await store.bootstrap(defaultYear)
    await store.loadPool({ silent: true })
  } catch {
    // 页面消息由 store 统一处理。
  }
})

async function handleYearChange() {
  try {
    await store.loadConfig()
    await store.loadPool({ silent: true })
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleRefresh() {
  try {
    await store.loadConfig()
    await store.loadPool({ silent: true })
  } catch {
    // 页面消息由 store 统一处理。
  }
}

function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  selectedFile.value = input.files?.[0] ?? null
}

async function handleUpload() {
  if (!selectedFile.value) {
    return
  }
  try {
    await store.uploadExcel(selectedFile.value)
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

async function handleLoadPool() {
  try {
    await store.loadPool()
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
.upload-body {
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

.upload-body,
.parse-preview,
.pool-filter {
  padding: 16px;
}

.file-picker {
  position: relative;
  display: inline-flex;
  min-width: 260px;
}

.file-picker input {
  position: absolute;
  inset: 0;
  opacity: 0;
  cursor: pointer;
}

.file-picker span {
  width: 100%;
  border: 1px dashed var(--line-strong);
  border-radius: 12px;
  padding: 10px 14px;
  background: #fbfcfe;
  color: var(--muted);
}

.parse-preview {
  border-top: 1px solid var(--line);
}

.preview-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.preview-head span {
  color: var(--muted);
  font-size: 13px;
}

.summary-chips,
.warning-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.summary-chips span,
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
  .panel-head,
  .preview-head {
    flex-direction: column;
    align-items: flex-start;
  }

  .stats-grid,
  .pool-filter {
    grid-template-columns: 1fr;
  }
}
</style>
