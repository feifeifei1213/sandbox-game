<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>组数据</h2>
        
      </div>
      <div class="hero-actions">
        <button type="button" class="btn" :disabled="loading" @click="handleRefresh">刷新</button>
      </div>
    </header>

    <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
      {{ pageMessage.text }}
    </section>

    <section class="panel-card query-card">
      <div class="panel-head">
        <div>
          <strong>查询条件</strong>
          
        </div>
      </div>
      <div class="query-grid">
        <label class="field">
          <span>小组</span>
          <select :value="selectedGroupId ?? ''" :disabled="loading || groups.length === 0" @change="handleGroupChange">
            <option value="" disabled>请选择小组</option>
            <option v-for="item in groups" :key="item.groupId" :value="item.groupId">
              第{{ item.groupNo }}组 · {{ formatBusinessStatus(item.businessStatus) }}
            </option>
          </select>
        </label>

        <label class="field">
          <span>年份</span>
          <select :value="selectedYear" :disabled="loading" @change="handleYearChange">
            <option v-for="item in yearOptions" :key="item.value" :value="item.value">{{ item.label }}</option>
          </select>
        </label>

        <div class="field">
          <span>页面类型</span>
          <div class="mode-switch">
            <button
              type="button"
              class="mode-btn"
              :class="{ active: selectedPageType === 'operating' }"
              :disabled="loading"
              @click="switchPageType('operating')"
            >
              经营页
            </button>
            <button
              type="button"
              class="mode-btn"
              :class="{ active: selectedPageType === 'report' }"
              :disabled="loading"
              @click="switchPageType('report')"
            >
              财报页
            </button>
          </div>
        </div>
      </div>
    </section>

    <div class="content-layout">
      <main class="preview-panel">
        <section v-if="loading" class="loading-card">正在读取组数据...</section>
        <section v-else-if="!selectedGroup" class="loading-card">当前没有可选小组数据。</section>
        <OperatingSheet
          v-else-if="selectedPageType === 'operating' && operatingView"
          :model-value="operatingView.operatingPayload"
          :editable-scopes="[]"
          :invalid-scopes="operatingView.invalidScopes"
          :quarter-cash-checks="operatingView.quarterCashChecks"
          :current-stage-code="operatingView.currentStageCode"
          :derived-values="operatingView.derivedValues"
          :period-end-cash="operatingView.periodEndCash"
          :carry-forward="operatingView.carryForward"
          :labels="activeOperatingLabels"
          @update:model-value="noopOperatingUpdate"
        />
        <ReportSheet
          v-else-if="selectedPageType === 'report' && reportView"
          :model-value="reportView.reportManualPayload"
          :computed-payload="reportView.reportComputedPayload"
          :can-edit="false"
          :has-invalid-draft="reportView.hasInvalidDraft"
          :tax-rate-options="reportView.manualFieldOptions.incomeTaxRateOptions"
          :labels="activeReportLabels"
          @update:model-value="noopReportUpdate"
        />
        <section v-else class="loading-card">当前暂无可展示的数据。</section>
      </main>

      <aside class="side-stack">
        <section class="panel-card">
          <div class="panel-head">
            <strong>当前目标</strong>
            
          </div>
          <div class="meta-list">
            <div class="meta-item">
              <span>小组</span>
              <strong>{{ selectedGroup ? `第${selectedGroup.groupNo}组` : '--' }}</strong>
            </div>
            <div class="meta-item">
              <span>年份</span>
              <strong>{{ selectedYear }} 年</strong>
            </div>
            <div class="meta-item">
              <span>页面</span>
              <strong>{{ selectedPageType === 'operating' ? '经营页' : '财报页' }}</strong>
            </div>
            <div class="meta-item">
              <span>经营状态</span>
              <strong>{{ activeBusinessStatusText }}</strong>
            </div>
          </div>
        </section>

        <section class="panel-card">
          <div class="panel-head">
            <strong>当前状态</strong>
            
          </div>
          <div class="status-grid">
            <div class="status-item">
              <span>年度状态</span>
              <strong>{{ activeYearStatusText }}</strong>
            </div>
            <div class="status-item" v-if="selectedPageType === 'operating' && operatingView">
              <span>阶段状态</span>
              <strong>{{ formatStageStatus(operatingView.stageStatus) }}</strong>
            </div>
            <div class="status-item" v-if="selectedPageType === 'operating' && operatingView">
              <span>当前阶段</span>
              <strong>{{ formatStageCode(operatingView.currentStageCode) }}</strong>
            </div>
            <div class="status-item" v-if="selectedPageType === 'operating' && operatingView">
              <span>财报状态</span>
              <strong>{{ formatReportStatus(operatingView.reportStatus) }}</strong>
            </div>
            <div class="status-item" v-if="selectedPageType === 'operating' && operatingView?.hasInvalidDraft">
              <span>失效草稿</span>
              <strong>{{ operatingView.hasRetainedReportDraft ? '经营待重提，且含财报失效草稿' : '经营结果待重新提交' }}</strong>
            </div>
            <div class="status-item" v-if="selectedPageType === 'report' && reportView">
              <span>财报状态</span>
              <strong>{{ formatReportStatus(reportView.reportStatus) }}</strong>
            </div>
            <div class="status-item" v-if="selectedPageType === 'report' && reportView">
              <span>平衡差额</span>
              <strong>{{ formatMetric(reportBalanceGap(reportView.reportComputedPayload)) }}</strong>
            </div>
            <div class="status-item" v-if="selectedPageType === 'report' && reportView">
              <span>草稿保存</span>
              <strong>{{ formatDateTime(reportView.lastDraftSavedAt) }}</strong>
            </div>
            <div class="status-item" v-if="selectedPageType === 'report' && reportView?.hasInvalidDraft">
              <span>失效草稿</span>
              <strong>财报草稿待重新提交</strong>
            </div>
          </div>
        </section>
      </aside>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted } from 'vue'
import { storeToRefs } from 'pinia'

import OperatingSheet from '@/components/sandbox-game/player/OperatingSheet.vue'
import ReportSheet from '@/components/sandbox-game/player/ReportSheet.vue'
import {
  applyDictionaryToOperatingLabels,
  applyDictionaryToReportLabels,
  resolveOperatingLabels,
  resolveReportLabels,
} from '@/configs/sandbox-game-service-labels'
import { useAdminGroupDataStore } from '@/stores/admin-group-data'
import { useAdminShellStore } from '@/stores/admin-shell'
import { useDictionaryStore } from '@/stores/dictionary'
import { reportBalanceGap, type OperatingPayload, type ReportManualPayload } from '@/types/sandbox-game'
import type { AdminGroupDataPageType } from '@/types/sandbox-game-admin'
import {
  formatBusinessStatus,
  formatReportStatus,
  formatStageCode,
  formatStageStatus,
  formatYearStatus,
} from '@/utils/sandbox-game-display'

const shellStore = useAdminShellStore()
const groupDataStore = useAdminGroupDataStore()
const dictionaryStore = useDictionaryStore()
const { config } = storeToRefs(shellStore)
const {
  groups,
  selectedGroupId,
  selectedYear,
  selectedPageType,
  operatingView,
  reportView,
  loading,
  pageMessage,
  selectedGroup,
} = storeToRefs(groupDataStore)

const yearOptions = computed(() => {
  const finalYear = Math.max(config.value?.finalYear ?? 0, 0)
  return Array.from({ length: finalYear + 1 }, (_, index) => ({
    value: index,
    label: `${index}年`,
  }))
})

const activeBusinessStatusText = computed(() => {
  const status = selectedPageType.value === 'operating' ? operatingView.value?.businessStatus : reportView.value?.businessStatus
  return formatBusinessStatus(status)
})

const activeYearStatusText = computed(() => {
  const status = selectedPageType.value === 'operating' ? operatingView.value?.yearStatus : reportView.value?.yearStatus
  return formatYearStatus(status)
})
const activeOperatingLabels = computed(() =>
  applyDictionaryToOperatingLabels(resolveOperatingLabels(config.value?.editionCode), dictionaryStore.displayName),
)
const activeReportLabels = computed(() =>
  applyDictionaryToReportLabels(resolveReportLabels(config.value?.editionCode), dictionaryStore.displayName),
)

onMounted(async () => {
  try {
    if (!shellStore.config) {
      await shellStore.bootstrap()
    }
    await Promise.all([
      dictionaryStore.loadCurrent(config.value?.editionCode, { silent: true }),
      groupDataStore.bootstrap(config.value?.finalYear ?? 0, config.value?.currentOpenYear ?? 0),
    ])
    dictionaryStore.startSilentSync()
  } catch {
    // 页面消息由 store 统一处理。
  }
})

onBeforeUnmount(() => {
  dictionaryStore.stopSilentSync()
})

async function handleRefresh() {
  try {
    await shellStore.refreshConfig({ silent: true })
    await groupDataStore.bootstrap(config.value?.finalYear ?? 0, config.value?.currentOpenYear ?? 0)
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleGroupChange(event: Event) {
  const target = event.target as HTMLSelectElement
  const nextGroupId = Number(target.value)
  if (!Number.isFinite(nextGroupId)) {
    return
  }
  groupDataStore.setSelectedGroupId(nextGroupId)
  await groupDataStore.loadCurrentView()
}

async function handleYearChange(event: Event) {
  const target = event.target as HTMLSelectElement
  const nextYear = Number(target.value)
  if (!Number.isFinite(nextYear)) {
    return
  }
  groupDataStore.setSelectedYear(nextYear)
  await groupDataStore.loadCurrentView()
}

async function switchPageType(pageType: AdminGroupDataPageType) {
  if (pageType === selectedPageType.value) {
    return
  }
  groupDataStore.setSelectedPageType(pageType)
  await groupDataStore.loadCurrentView()
}

function noopOperatingUpdate(_: OperatingPayload) {
  // 管理员组数据页复用玩家视图，但始终只读。
}

function noopReportUpdate(_: ReportManualPayload) {
  // 管理员组数据页复用玩家视图，但始终只读。
}

function formatMetric(value?: number) {
  if (typeof value !== 'number') {
    return '--'
  }
  return value.toLocaleString('zh-CN', {
    minimumFractionDigits: Number.isInteger(value) ? 0 : 2,
    maximumFractionDigits: 2,
  })
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

.hero-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.btn {
  border: 1px solid var(--line);
  background: #ffffff;
  border-radius: 12px;
  padding: 10px 14px;
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

.panel-card {
  border: 1px solid var(--line);
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.05);
}

.panel-head {
  padding: 14px 16px;
  background: #f7f9fc;
  border-bottom: 1px solid var(--line);
}

.panel-head strong {
  display: block;
  margin-bottom: 4px;
  font-size: 16px;
}

.panel-head span {
  color: var(--muted);
  font-size: 13px;
}

.query-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  padding: 16px;
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

.field select,
.field textarea {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #ffffff;
  padding: 10px 12px;
  font: inherit;
}

.mode-switch {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.mode-btn {
  border: 1px solid var(--line);
  background: #ffffff;
  border-radius: 12px;
  padding: 10px 12px;
}

.mode-btn.active {
  color: #ffffff;
  border-color: var(--accent);
  background: var(--accent);
}

.content-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 14px;
  align-items: start;
}

.preview-panel {
  min-width: 0;
}

.side-stack {
  display: grid;
  gap: 14px;
}

.loading-card {
  border: 1px dashed var(--line-strong);
  border-radius: 16px;
  background: #ffffff;
  padding: 28px;
  text-align: center;
  color: var(--muted);
}

.meta-list,
.status-grid {
  display: grid;
  gap: 12px;
  padding: 16px;
}

.meta-item,
.status-item {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.meta-item span,
.status-item span {
  color: var(--muted);
  font-size: 13px;
}

.meta-item strong,
.status-item strong {
  font-size: 14px;
}

@media (max-width: 1280px) {
  .content-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 960px) {
  .hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .query-grid {
    grid-template-columns: 1fr;
  }
}
</style>



