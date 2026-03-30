<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>组数据</h2>
        <p>管理员按“组 → 年 → 页面类型”查看数据，并在同页发起异常解锁。</p>
      </div>
      <div class="hero-actions">
        <button type="button" class="btn" :disabled="loading" @click="handleRefresh">刷新</button>
        <button type="button" class="btn primary" :disabled="!selectedGroup || unlocking" @click="openUnlockDialog">
          异常解锁
        </button>
      </div>
    </header>

    <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
      {{ pageMessage.text }}
    </section>

    <section class="panel-card query-card">
      <div class="panel-head">
        <div>
          <strong>查询条件</strong>
          <span>默认展示经营页，只读复用玩家端页面结构。</span>
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
          :quarter-cash-checks="operatingView.quarterCashChecks"
          :current-stage-code="operatingView.currentStageCode"
          :derived-values="operatingView.derivedValues"
          :period-end-cash="operatingView.periodEndCash"
          @update:model-value="noopOperatingUpdate"
        />
        <ReportSheet
          v-else-if="selectedPageType === 'report' && reportView"
          :model-value="reportView.reportManualPayload"
          :computed-payload="reportView.reportComputedPayload"
          :can-edit="false"
          :tax-rate-options="reportView.manualFieldOptions.incomeTaxRateOptions"
          @update:model-value="noopReportUpdate"
        />
        <section v-else class="loading-card">当前暂无可展示的数据。</section>
      </main>

      <aside class="side-stack">
        <section class="panel-card">
          <div class="panel-head">
            <strong>当前目标</strong>
            <span>管理员只读查看，不在本页直接编辑业务数据。</span>
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
            <span>根据当前查看页面展示年度 / 阶段 / 财报状态。</span>
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
          </div>
        </section>

        <section class="panel-card">
          <div class="panel-head">
            <strong>异常解锁</strong>
            <span>是否解锁由管理员现场判断，服务端只做最小硬校验。</span>
          </div>
          <div class="unlock-box">
            <p>若发现该组该年数据需要重新提交，可发起异常解锁，恢复当前年份为可编辑状态。</p>
            <button type="button" class="btn primary full" :disabled="!selectedGroup || unlocking" @click="openUnlockDialog">
              {{ unlocking ? '提交中...' : '打开异常解锁弹窗' }}
            </button>
          </div>
        </section>

        <section class="panel-card" v-if="latestUnlockResult">
          <div class="panel-head">
            <strong>最近一次解锁记录</strong>
            <span>显示当前页面最近一次成功提交的异常解锁结果。</span>
          </div>
          <div class="meta-list">
            <div class="meta-item">
              <span>解锁日志 ID</span>
              <strong>{{ latestUnlockResult.unlockLogId }}</strong>
            </div>
            <div class="meta-item">
              <span>年份状态</span>
              <strong>{{ formatYearStatus(latestUnlockResult.yearStatus) }}</strong>
            </div>
            <div class="meta-item">
              <span>阶段状态</span>
              <strong>{{ formatStageStatus(latestUnlockResult.stageStatus) }}</strong>
            </div>
            <div class="meta-item">
              <span>财报状态</span>
              <strong>{{ formatReportStatus(latestUnlockResult.reportStatus) }}</strong>
            </div>
            <div class="meta-item wide-item">
              <span>解锁原因</span>
              <strong>{{ latestUnlockReason }}</strong>
            </div>
          </div>
        </section>
      </aside>
    </div>

    <div v-if="unlockDialogVisible" class="dialog-mask" @click.self="closeUnlockDialog">
      <div class="dialog-card">
        <div class="dialog-head">
          <div>
            <strong>异常解锁确认</strong>
            <span>提交前请再次确认目标组、年份和原因。</span>
          </div>
        </div>
        <div class="dialog-body">
          <div class="dialog-target">
            <div><span>目标小组</span><strong>{{ selectedGroup ? `第${selectedGroup.groupNo}组` : '--' }}</strong></div>
            <div><span>目标年份</span><strong>{{ selectedYear }} 年</strong></div>
            <div><span>当前页面</span><strong>{{ selectedPageType === 'operating' ? '经营页' : '财报页' }}</strong></div>
          </div>
          <label class="field">
            <span>解锁原因</span>
            <textarea :value="unlockReason" rows="5" placeholder="请输入管理员现场确认后的解锁原因" @input="handleUnlockReasonInput" />
          </label>
        </div>
        <div class="dialog-actions">
          <button type="button" class="btn" :disabled="unlocking" @click="closeUnlockDialog">取消</button>
          <button type="button" class="btn primary" :disabled="unlocking" @click="handleUnlockSubmit">
            {{ unlocking ? '提交中...' : '确认异常解锁' }}
          </button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'

import OperatingSheet from '@/components/sandbox-game/player/OperatingSheet.vue'
import ReportSheet from '@/components/sandbox-game/player/ReportSheet.vue'
import { useAdminGroupDataStore } from '@/stores/admin-group-data'
import { useAdminShellStore } from '@/stores/admin-shell'
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
const { config } = storeToRefs(shellStore)
const {
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
  latestUnlockResult,
  latestUnlockReason,
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

onMounted(async () => {
  try {
    if (!shellStore.config) {
      await shellStore.bootstrap()
    }
    await groupDataStore.bootstrap(config.value?.finalYear ?? 0, config.value?.currentOpenYear ?? 0)
  } catch {
    // 页面消息由 store 统一处理。
  }
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

function openUnlockDialog() {
  groupDataStore.openUnlockDialog()
}

function closeUnlockDialog() {
  groupDataStore.closeUnlockDialog()
}

function handleUnlockReasonInput(event: Event) {
  const target = event.target as HTMLTextAreaElement
  groupDataStore.setUnlockReason(target.value)
}

async function handleUnlockSubmit() {
  try {
    await groupDataStore.submitUnlock()
    await shellStore.refreshConfig({ silent: true })
  } catch {
    // 页面消息由 store 统一处理。
  }
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

.hero-actions,
.dialog-actions {
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

.btn.primary {
  color: #ffffff;
  border-color: var(--accent);
  background: var(--accent);
}

.btn.full {
  width: 100%;
  justify-content: center;
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
.status-item span,
.dialog-target span {
  color: var(--muted);
  font-size: 13px;
}

.meta-item strong,
.status-item strong,
.dialog-target strong {
  font-size: 14px;
}

.wide-item {
  align-items: flex-start;
}

.wide-item strong {
  text-align: right;
  line-height: 1.6;
}

.unlock-box {
  display: grid;
  gap: 12px;
  padding: 16px;
}

.unlock-box p {
  margin: 0;
  color: var(--muted);
  line-height: 1.7;
}

.dialog-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.36);
  display: grid;
  place-items: center;
  padding: 20px;
  z-index: 20;
}

.dialog-card {
  width: min(560px, 100%);
  border-radius: 20px;
  background: #ffffff;
  box-shadow: 0 24px 80px rgba(15, 23, 42, 0.28);
  overflow: hidden;
}

.dialog-head {
  padding: 18px 20px 14px;
  border-bottom: 1px solid var(--line);
  background: #f8fafc;
}

.dialog-head strong {
  display: block;
  margin-bottom: 4px;
  font-size: 18px;
}

.dialog-head span {
  color: var(--muted);
  font-size: 13px;
}

.dialog-body {
  display: grid;
  gap: 16px;
  padding: 20px;
}

.dialog-target {
  display: grid;
  gap: 10px;
}

.dialog-target div {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.dialog-actions {
  justify-content: flex-end;
  padding: 0 20px 20px;
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



