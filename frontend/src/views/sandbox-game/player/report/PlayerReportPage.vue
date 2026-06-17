<template>
  <div class="player-page">
    <div class="shell">
      <header class="page-header">
        <div>
          <p class="eyebrow">Sandbox Game / Player</p>
          <h1>玩家财报页</h1>
          <p class="subtext">本年经营结束后开放，黄色区域自动计算，绿色区域由玩家填写或确认。</p>
        </div>
        <div class="header-pills">
          <span class="pill">组别：{{ activeView?.groupId ?? '--' }}</span>
          <span class="pill">开放年份：{{ activeConfig?.currentOpenYear ?? '--' }}</span>
          <span class="pill">最终年份：{{ activeConfig?.finalYear ?? '--' }}</span>
          <span v-if="!previewMode && currentUser" class="pill">账号：{{ currentUser.username }}</span>
          <button v-if="!previewMode" class="logout-button" type="button" @click="handleLogout">退出登录</button>
          <span v-if="previewMode" class="pill preview-pill">开发预览</span>
        </div>
      </header>

      <section class="toolbar-card">
        <div class="toolbar-top">
          <YearTabs :tabs="activeYearTabs" :active-year="activeSelectedYear" @select="handleYearSelect" />
          <PageModeSwitch active-mode="report" :report-enabled="true" @operating="goOperating" @order="goOrder" />
        </div>
      </section>

      <section v-if="activePageMessage" class="message-bar" :class="activePageMessage.type">
        {{ activePageMessage.text }}
      </section>

      <div class="workspace">
        <main class="main-panel">
          <section class="status-banner">
            <div class="status-item">
              <span>当前年份</span>
              <strong>{{ activeSelectedYear }} 年</strong>
            </div>
            <div class="status-item">
              <span>年度状态</span>
              <strong>{{ formatYearStatus(activeView?.yearStatus) }}</strong>
            </div>
            <div class="status-item">
              <span>财报状态</span>
              <strong>{{ formatReportStatus(activeView?.reportStatus) }}</strong>
            </div>
            <div class="status-item">
              <span>经营状态</span>
              <strong>{{ formatBusinessStatus(activeView?.businessStatus) }}</strong>
            </div>
            <div class="status-item" :class="activeBalancePassed ? 'status-pass' : 'status-fail'">
              <span>平衡校验</span>
              <strong>{{ activeBalancePassed ? '通过' : '未通过' }}</strong>
            </div>
          </section>

          <section v-if="activeView?.rollbackPending && activeView.rollbackNotice" class="message-bar info">
            {{ activeView.rollbackNotice }}
          </section>

          <section v-if="!previewMode && (loading || yearViewLoading)" class="loading-card">正在加载财报页数据...</section>

          <section v-else-if="airportPendingMode" class="loading-card pending-template-card">
            <strong>机场沙盘版财报页待接入</strong>
            <span>当前仅支持订单模块联调，财报字段和公式将在机场版财报 Excel 给到后接入。</span>
          </section>

          <ReportSheet
            v-else-if="activeView"
            :model-value="activeDraftManualPayload"
            :computed-payload="activeComputedPayload"
            :can-edit="activeView.canEdit"
            :has-invalid-draft="activeView.hasInvalidDraft"
            :tax-rate-options="activeView.manualFieldOptions.incomeTaxRateOptions"
            :labels="activeReportLabels"
            @update:model-value="updatePayload"
          />

          <section v-else class="loading-card">当前没有可展示的财报页数据。</section>
        </main>

        <ReportSidebar
          v-if="!airportPendingMode"
          class="side-panel"
          :view="activeView"
          :year-label="`${activeSelectedYear} 年财报`"
          :computed-payload="activeComputedPayload"
          :balance-gap="activeBalanceGap"
          :balance-passed="activeBalancePassed"
          :tax-rate-options="activeView?.manualFieldOptions.incomeTaxRateOptions ?? []"
          :missing-fields="activeMissingFields"
          :integer-issues="activeManualIntegerIssues"
          :dirty="activeDirty"
          :saving="activeSaving"
          :submitting="activeSubmitting"
          :submit-ready="activeSubmitReady"
          :labels="activeReportLabels"
          @save="handleSave"
          @submit="handleSubmit"
        />
        <section v-else class="side-panel pending-side">
          <strong>财报页未开放</strong>
          <span>机场版经营页和财报页未接入前，不允许填写或提交财报。</span>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'

import YearTabs from '@/components/sandbox-game/common/YearTabs.vue'
import PageModeSwitch from '@/components/sandbox-game/player/PageModeSwitch.vue'
import ReportSheet from '@/components/sandbox-game/player/ReportSheet.vue'
import ReportSidebar from '@/components/sandbox-game/player/ReportSidebar.vue'
import {
  applyDictionaryToReportLabels,
  applyDictionaryToReportRequiredFieldLabels,
  resolveReportLabels,
  resolveReportRequiredFieldLabels,
} from '@/configs/sandbox-game-service-labels'
import { useAuthStore } from '@/stores/auth'
import { useDictionaryStore } from '@/stores/dictionary'
import { formatBusinessStatus, formatReportStatus, formatYearStatus } from '@/utils/sandbox-game-display'
import { usePlayerReportStore } from '@/stores/player-report'
import type { PageMessage } from '@/stores/player-report'
import {
  buildReportComputedPreview,
  cloneReportManualPayload,
  createEmptyReportComputedPayload,
  type CurrentGameConfigResult,
  type PlayerReportView,
  type ReportComputedPayload,
  type ReportManualPayload,
  type YearTabItem,
} from '@/types/sandbox-game'

const AUTO_SAVE_INTERVAL = 5 * 60 * 1000
const balanceTolerance = 0.000001

const route = useRoute()
const router = useRouter()
const store = usePlayerReportStore()
const authStore = useAuthStore()
const dictionaryStore = useDictionaryStore()
const {
  currentConfig,
  yearTabs,
  currentView,
  draftManualPayload,
  selectedYear,
  loading,
  yearViewLoading,
  saving,
  submitting,
  dirty,
  pageMessage,
  previewComputedPayload: storePreviewComputedPayload,
  balanceGap,
  balancePassed,
  missingFields,
  manualIntegerIssues,
  submitReady,
} = storeToRefs(store)
const { currentUser } = storeToRefs(authStore)

const previewDraftManualPayload = ref<ReportManualPayload>({
  workInProgress: 6,
  finishedGoods: 4,
  rawMaterials: 2,
  incomeTaxRate: 0.25,
  enterpriseCertificationScore: 3,
  productionHumanScore: 4,
  closingSpeedScore: 2,
})
const previewDirty = ref(false)
const previewSaving = ref(false)
const previewSubmitting = ref(false)
const previewLastDraftSavedAt = ref<string | null>(null)
const previewPageMessage = ref<PageMessage | null>(null)

let initialized = false
let autoSaveTimer = 0

const previewMode = computed(() => {
  const raw = Array.isArray(route.query.preview) ? route.query.preview[0] : route.query.preview
  return raw === '1' || raw === 'true'
})

const previewYear = computed(() => readRouteYear() ?? 0)
const previewConfig = computed<CurrentGameConfigResult>(() => ({
  currentOpenYear: previewYear.value,
  finalYear: 8,
  editionCode: 'VIP_SERVICE_V1',
  editionName: '贵宾服务版 V1',
  ruleVersion: 'COMMON_FORMULA_V1',
  formulaVersion: 'COMMON_FORMULA_V1',
  templateVersion: 'VIP_SERVICE_V1',
  operatingTemplateVersion: 'VIP_OPERATING_TEMPLATE_V1',
  reportTemplateVersion: 'VIP_REPORT_TEMPLATE_V1',
  orderTemplateVersion: 'VIP_ORDER_TEMPLATE_V1',
  processRuleVersion: 'COMMON_PROCESS_V1',
  dictionaryRevision: 0,
  demoYearEnabled: true,
}))
const previewYearTabs = computed<YearTabItem[]>(() => {
  const tabs: YearTabItem[] = []
  for (let yearNo = 0; yearNo <= previewConfig.value.finalYear; yearNo += 1) {
    tabs.push({
      yearNo,
      label: `${yearNo}年`,
      tabStatus: yearNo === previewYear.value ? 'ENTERABLE' : yearNo < previewYear.value ? 'COMPLETED' : 'ENTERABLE',
      canEnter: true,
      isCurrentOpenYear: yearNo === previewYear.value,
      isFormalYear: yearNo > 0,
    })
  }
  return tabs
})
const previewBaseComputedPayload = computed<ReportComputedPayload>(() => {
  const yearFactor = previewYear.value
  return {
    ...createEmptyReportComputedPayload(),
    reportSalesRevenue: 86 + yearFactor * 4,
    reportDirectCost: 42 + yearFactor * 2,
    reportGrossProfit: 44 + yearFactor * 2,
    reportComprehensiveCost: 18 + yearFactor,
    reportDepreciation: 6,
    reportOperatingProfit: 20 + yearFactor,
    reportFinanceIncomeExpense: 3,
    reportExtraIncomeExpense: 1,
    reportPreTaxProfit: 18 + yearFactor,
    reportIncomeTax: 0,
    reportNetProfit: 0,
    reportWorkInProgress: 0,
    reportFinishedGoods: 0,
    reportRawMaterials: 0,
    reportWorkInConstruction: 5,
    reportFactoryAsset: 40,
    reportLineResidual: 9,
    reportDepreciableAsset: 12,
    reportTotalNonCurrentAssets: 66,
    reportCash: 28 + yearFactor * 3,
    reportReceivable: 14,
    reportPostTaxCash: 0,
    reportTotalCurrentAssets: 0,
    reportTotalAssets: 0,
    reportShortTermLiability: 22,
    reportLongTermLiability: 10,
    reportTotalLiability: 32,
    reportShareCapital: 50,
    reportRetainedEarnings: 12 + yearFactor,
    reportTotalEquity: 0,
    reportTotalLiabilityEquity: 0,
    reportBestMarketDirectorBaseScore: 2 + yearFactor,
    reportBestMarketDirectorScore: 0,
    reportBestTechnologyDirectorScore: 4 + yearFactor * 2,
    reportBestSalesDirectorScore: 36 + yearFactor * 8,
    reportBestCfoBaseScore: 5 + yearFactor,
    reportBestCfoScore: 0,
    reportBestCeoScore: 0,
  }
})
const previewComputedPayloadLocal = computed(() =>
  buildReportComputedPreview(previewBaseComputedPayload.value, previewDraftManualPayload.value),
)
const previewBalanceGap = computed(
  () => previewComputedPayloadLocal.value.reportTotalAssets - previewComputedPayloadLocal.value.reportTotalLiabilityEquity,
)
const previewBalancePassed = computed(() => Math.abs(previewBalanceGap.value) <= balanceTolerance)
const previewRequiredFieldLabels = computed(() =>
  applyDictionaryToReportRequiredFieldLabels(resolveReportRequiredFieldLabels(previewConfig.value.editionCode), dictionaryStore.displayName),
)
const previewMissingFields = computed(() => {
  const missing: string[] = []
  if (previewDraftManualPayload.value.workInProgress === null) {
    missing.push(previewRequiredFieldLabels.value.workInProgress)
  }
  if (previewDraftManualPayload.value.finishedGoods === null) {
    missing.push(previewRequiredFieldLabels.value.finishedGoods)
  }
  if (previewDraftManualPayload.value.rawMaterials === null) {
    missing.push(previewRequiredFieldLabels.value.rawMaterials)
  }
  if (previewDraftManualPayload.value.incomeTaxRate === null) {
    missing.push(previewRequiredFieldLabels.value.incomeTaxRate)
  }
  if (previewDraftManualPayload.value.enterpriseCertificationScore === null) {
    missing.push(previewRequiredFieldLabels.value.enterpriseCertificationScore)
  }
  if (previewDraftManualPayload.value.productionHumanScore === null) {
    missing.push(previewRequiredFieldLabels.value.productionHumanScore)
  }
  if (previewDraftManualPayload.value.closingSpeedScore === null) {
    missing.push(previewRequiredFieldLabels.value.closingSpeedScore)
  }
  return missing
})
const previewSubmitReady = computed(() => previewMissingFields.value.length === 0 && previewBalancePassed.value)
const previewView = computed<PlayerReportView>(() => ({
  groupId: 1,
  yearNo: previewYear.value,
  yearStatus: 'REPORTING',
  reportStatus: 'REPORT_OPEN',
  businessStatus: 'NORMAL',
  canView: true,
  canEdit: true,
  canSubmit: true,
  hasInvalidDraft: false,
  rollbackPending: false,
  rollbackTargetYearNo: null,
  rollbackTargetStageCode: null,
  rollbackNotice: '',
  reportComputedPayload: previewBaseComputedPayload.value,
  reportManualPayload: cloneReportManualPayload(previewDraftManualPayload.value),
  manualFieldOptions: {
    incomeTaxRateOptions: [0.25, 0.15, 0],
  },
  lastDraftSavedAt: previewLastDraftSavedAt.value,
  noticeBoard: buildPreviewNoticeBoard(previewYear.value),
}))

const activeConfig = computed(() => (previewMode.value ? previewConfig.value : currentConfig.value))
const activeYearTabs = computed(() => (previewMode.value ? previewYearTabs.value : yearTabs.value))
const activeView = computed(() => (previewMode.value ? previewView.value : currentView.value))
const activeSelectedYear = computed(() => (previewMode.value ? previewYear.value : selectedYear.value))
const activeDraftManualPayload = computed(() => (previewMode.value ? previewDraftManualPayload.value : draftManualPayload.value))
const activeComputedPayload = computed(() =>
  previewMode.value
    ? previewComputedPayloadLocal.value
    : activeView.value?.rollbackPending && activeView.value.hasInvalidDraft && !activeView.value.canEdit
      ? activeView.value.reportComputedPayload
      : storePreviewComputedPayload.value,
)
const activeBalanceGap = computed(() => (previewMode.value ? previewBalanceGap.value : balanceGap.value))
const activeBalancePassed = computed(() => (previewMode.value ? previewBalancePassed.value : balancePassed.value))
const activeMissingFields = computed(() => (previewMode.value ? previewMissingFields.value : missingFields.value))
const activeManualIntegerIssues = computed(() => (previewMode.value ? [] : manualIntegerIssues.value))
const activeSubmitReady = computed(() => (previewMode.value ? previewSubmitReady.value : submitReady.value))
const activePageMessage = computed(() => (previewMode.value ? previewPageMessage.value : pageMessage.value))
const activeDirty = computed(() => (previewMode.value ? previewDirty.value : dirty.value))
const activeSaving = computed(() => (previewMode.value ? previewSaving.value : saving.value))
const activeSubmitting = computed(() => (previewMode.value ? previewSubmitting.value : submitting.value))
const activeReportLabels = computed(() =>
  applyDictionaryToReportLabels(resolveReportLabels(activeConfig.value?.editionCode), dictionaryStore.displayName),
)
const airportPendingMode = computed(() =>
  !previewMode.value
  && (activeConfig.value?.editionCode === 'AIRPORT_V1' || activeConfig.value?.reportTemplateVersion === 'PENDING_REPORT_TEMPLATE'),
)

onMounted(async () => {
  if (previewMode.value) {
    previewPageMessage.value = {
      type: 'info',
      text: '当前为开发预览模式，页面样式可查看，但不会读取或提交真实业务数据。',
    }
    return
  }

  try {
    await Promise.all([
      store.bootstrap(readRouteYear()),
      dictionaryStore.loadCurrent(null, { silent: true }),
    ])
    dictionaryStore.startSilentSync()
    initialized = true
    if (selectedYear.value !== readRouteYear()) {
      syncRouteYear(selectedYear.value)
    }
  } catch {
    // 页面消息由 store 统一处理。
  }

  autoSaveTimer = window.setInterval(async () => {
    if (!dirty.value || !currentView.value?.canEdit || saving.value || submitting.value) {
      return
    }
    try {
      await store.saveDraft({ silent: true })
    } catch {
      // 自动保存失败时，store 内部会保留页面消息。
    }
  }, AUTO_SAVE_INTERVAL)
})

watch(
  () => [route.query.yearNo, route.query.preview],
  async () => {
    if (previewMode.value) {
      previewPageMessage.value = {
        type: 'info',
        text: '当前为开发预览模式，页面样式可查看，但不会读取或提交真实业务数据。',
      }
      return
    }

    const targetYear = readRouteYear()
    if (!initialized || targetYear === undefined || targetYear === selectedYear.value) {
      return
    }
    try {
      await store.loadYearView(targetYear)
    } catch {
      syncRouteYear(selectedYear.value)
    }
  },
)

onBeforeUnmount(() => {
  if (autoSaveTimer) {
    window.clearInterval(autoSaveTimer)
  }
  dictionaryStore.stopSilentSync()
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
  if (yearNo === activeSelectedYear.value) {
    return
  }
  syncRouteYear(yearNo)
}

function updatePayload(nextPayload: ReportManualPayload) {
  if (previewMode.value) {
    previewDraftManualPayload.value = cloneReportManualPayload(nextPayload)
    previewDirty.value = true
    return
  }
  store.updateDraft(nextPayload)
}

async function handleSave() {
  if (previewMode.value) {
    previewSaving.value = true
    previewLastDraftSavedAt.value = new Date().toISOString()
    previewDirty.value = false
    previewPageMessage.value = {
      type: 'success',
      text: `预览模式草稿已本地保存，时间 ${new Date(previewLastDraftSavedAt.value).toLocaleString('zh-CN', { hour12: false })}`,
    }
    previewSaving.value = false
    return
  }
  await store.saveDraft()
}

async function handleSubmit() {
  if (previewMode.value) {
    previewSubmitting.value = true
    previewPageMessage.value = {
      type: 'info',
      text: '当前是开发预览模式，提交按钮只用于查看交互，不会提交到后端。',
    }
    previewSubmitting.value = false
    return
  }
  await store.submitCurrentReport()
}

function buildPreviewNoticeBoard(yearNo: number) {
  return {
    pinnedNotice: {
      id: 11,
      kind: 'GENERAL',
      title: '系统通知',
      content: '主持人提示：本年经营结束后，请尽快完成财报绿色手工项填写。',
      pinned: true,
      publishedAt: new Date(2026, 2, 21, 13, 30).toISOString(),
      yearNo,
      stageCode: null,
      amount: null,
    },
    recentList: [
      {
        id: 12,
        kind: 'PENALTY',
        title: `${yearNo}年 Q4 罚款`,
        content: '示例：因现场判罚扣减 2。',
        pinned: false,
        publishedAt: new Date(2026, 2, 21, 13, 45).toISOString(),
        yearNo,
        stageCode: 'Q4',
        amount: 2,
      },
    ],
  }
}
function goOperating() {
  router.push({
    path: '/sandbox-game/player/operating',
    query: previewMode.value ? { yearNo: String(activeSelectedYear.value), preview: '1' } : { yearNo: String(activeSelectedYear.value) },
  })
}

function goOrder() {
  router.push({
    path: '/sandbox-game/player/orders',
    query: previewMode.value ? { yearNo: String(activeSelectedYear.value), preview: '1' } : { yearNo: String(activeSelectedYear.value) },
  })
}

async function handleLogout() {
  try {
    await authStore.logout()
  } finally {
    await router.replace('/sandbox-game/login')
  }
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

.preview-pill {
  border-color: #e0b24c;
  background: #fff5d8;
  color: #946200;
}

.logout-button {
  height: 36px;
  padding: 0 14px;
  border-radius: 999px;
  border: 1px solid #d4dae4;
  background: #ffffff;
  color: var(--text);
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

.workspace {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 280px;
  gap: 14px;
  padding: 14px;
  align-items: start;
}

.main-panel {
  min-width: 0;
}

.status-banner {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.status-item {
  border: 1px solid var(--line);
  border-radius: 16px;
  background: #ffffff;
  padding: 14px 16px;
}

.status-item span {
  display: block;
  color: var(--muted);
  font-size: 12px;
  margin-bottom: 6px;
}

.status-item strong {
  font-size: 18px;
}

.status-pass {
  background: #eefaf2;
  border-color: #b7dec6;
}

.status-fail {
  background: #fff6f6;
  border-color: #efc6c6;
}

.loading-card {
  border: 1px dashed var(--line-strong);
  border-radius: 16px;
  background: #ffffff;
  padding: 28px;
  text-align: center;
  color: var(--muted);
}

.pending-template-card {
  display: grid;
  gap: 8px;
}

.pending-template-card strong {
  color: var(--text);
  font-size: 18px;
}

.side-panel {
  position: sticky;
  top: 18px;
}

.pending-side {
  display: grid;
  gap: 8px;
  border: 1px solid var(--line);
  border-radius: 16px;
  background: #ffffff;
  padding: 18px;
}

.pending-side span {
  color: var(--muted);
  line-height: 1.6;
}

@media (max-width: 1360px) {
  .workspace {
    grid-template-columns: 1fr;
  }

  .side-panel {
    position: static;
  }
}

@media (max-width: 1024px) {
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

  .status-banner {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>











