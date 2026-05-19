<template>
  <div class="player-page">
    <div class="shell">
      <header class="page-header">
        <div>
          <p class="eyebrow">Sandbox Game / Player</p>
          <h1>玩家经营页</h1>
          <p class="subtext">正式前端工程版本，当前已接入真实年份标签与经营页查询 / 保存 / 提交接口。</p>
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
          <PageModeSwitch active-mode="operating" :report-enabled="activeReportEnabled" @order="goOrder" @report="goReport" />
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
              <span>经营阶段</span>
              <strong>{{ formatStageCode(activeView?.currentStageCode) }}</strong>
            </div>
            <div class="status-item">
              <span>财报状态</span>
              <strong>{{ formatReportStatus(activeView?.reportStatus) }}</strong>
            </div>
          </section>

          <section v-if="!previewMode && (loading || yearViewLoading)" class="loading-card">正在加载经营页数据...</section>

          <OperatingSheet
            v-else-if="activeView"
            :model-value="activeDraftPayload"
            :editable-scopes="activeView.editableScopes"
            :invalid-scopes="activeView.invalidScopes"
            :quarter-cash-checks="activeQuarterCashChecks"
            :current-stage-code="activeView.currentStageCode"
            :derived-values="activeDerivedValues"
            :period-end-cash="activePeriodEndCash"
            :carry-forward="activeView.carryForward"
            @update:model-value="updatePayload"
          />

          <section v-else class="loading-card">当前没有可展示的经营页数据。</section>
        </main>

        <OperatingSidebar
          class="side-panel"
          :view="activeView"
          :year-label="`${activeSelectedYear} 年经营`"
          :dirty="activeDirty"
          :saving="activeSaving"
          :submitting="activeSubmitting"
          @save="handleSave"
          @submit="handleSubmit"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'

import YearTabs from '@/components/sandbox-game/common/YearTabs.vue'
import OperatingSheet from '@/components/sandbox-game/player/OperatingSheet.vue'
import OperatingSidebar from '@/components/sandbox-game/player/OperatingSidebar.vue'
import PageModeSwitch from '@/components/sandbox-game/player/PageModeSwitch.vue'
import { useAuthStore } from '@/stores/auth'
import { formatReportStatus, formatStageCode, formatYearStatus } from '@/utils/sandbox-game-display'
import { buildOperatingPreviewCalculation } from '@/utils/sandbox-game-operating-preview'
import type { PageMessage } from '@/stores/player-operating'
import { usePlayerOperatingStore } from '@/stores/player-operating'
import {
  cloneOperatingPayload,
  createEmptyOperatingPayload,
  type CurrentGameConfigResult,
  type OperatingPayload,
  type PlayerOperatingView,
  type YearTabItem,
} from '@/types/sandbox-game'

const AUTO_SAVE_INTERVAL = 5 * 60 * 1000

const route = useRoute()
const router = useRouter()
const store = usePlayerOperatingStore()
const authStore = useAuthStore()
const {
  currentConfig,
  yearTabs,
  currentView,
  draftPayload,
  selectedYear,
  loading,
  yearViewLoading,
  saving,
  submitting,
  dirty,
  pageMessage,
  reportEnabled,
  previewCalculation: storePreviewCalculation,
} = storeToRefs(store)
const { currentUser } = storeToRefs(authStore)

const previewDraftPayload = ref<OperatingPayload>(buildPreviewOperatingPayload(0))
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
  ruleVersion: 'preview',
  templateVersion: 'preview',
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
const previewView = computed<PlayerOperatingView>(() => buildPreviewView(previewYear.value, previewLastDraftSavedAt.value))
const previewCalculation = computed(() =>
  buildOperatingPreviewCalculation({
    payload: previewDraftPayload.value,
    carryForward: previewView.value.carryForward,
    currentStageCode: previewView.value.currentStageCode,
    fallback: {
      quarterCashChecks: previewView.value.quarterCashChecks,
      derivedValues: previewView.value.derivedValues,
      periodEndCash: previewView.value.periodEndCash,
    },
  }),
)

const activeConfig = computed(() => (previewMode.value ? previewConfig.value : currentConfig.value))
const activeYearTabs = computed(() => (previewMode.value ? previewYearTabs.value : yearTabs.value))
const activeView = computed(() => (previewMode.value ? previewView.value : currentView.value))
const activeSelectedYear = computed(() => (previewMode.value ? previewYear.value : selectedYear.value))
const activeDraftPayload = computed(() => (previewMode.value ? previewDraftPayload.value : draftPayload.value))
const activeQuarterCashChecks = computed(() => (previewMode.value ? previewCalculation.value.quarterCashChecks : storePreviewCalculation.value.quarterCashChecks))
const activeDerivedValues = computed(() => (previewMode.value ? previewCalculation.value.derivedValues : storePreviewCalculation.value.derivedValues))
const activePeriodEndCash = computed(() => (previewMode.value ? previewCalculation.value.periodEndCash : storePreviewCalculation.value.periodEndCash))
const activePageMessage = computed(() => (previewMode.value ? previewPageMessage.value : pageMessage.value))
const activeDirty = computed(() => (previewMode.value ? previewDirty.value : dirty.value))
const activeSaving = computed(() => (previewMode.value ? previewSaving.value : saving.value))
const activeSubmitting = computed(() => (previewMode.value ? previewSubmitting.value : submitting.value))
const activeReportEnabled = computed(() => (previewMode.value ? true : reportEnabled.value))

onMounted(async () => {
  if (previewMode.value) {
    resetPreviewState(previewYear.value)
    return
  }

  await store.bootstrap(readRouteYear())
  initialized = true
  if (selectedYear.value !== readRouteYear()) {
    syncRouteYear(selectedYear.value)
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
      resetPreviewState(previewYear.value)
      return
    }

    const targetYear = readRouteYear()
    if (!initialized) {
      await store.bootstrap(targetYear)
      initialized = true
      return
    }
    if (targetYear === undefined || targetYear === selectedYear.value) {
      return
    }
    try {
      await store.loadYearView(targetYear)
    } catch {
      // 错误消息由 store 统一处理。
    }
  },
)

onBeforeUnmount(() => {
  if (autoSaveTimer) {
    window.clearInterval(autoSaveTimer)
  }
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

function updatePayload(nextPayload: OperatingPayload) {
  if (previewMode.value) {
    previewDraftPayload.value = cloneOperatingPayload(nextPayload)
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
  await store.submitCurrentStage()
}

function goReport() {
  if (!activeReportEnabled.value) {
    return
  }
  router.push({
    path: '/sandbox-game/player/report',
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

function resetPreviewState(yearNo: number) {
  previewDraftPayload.value = buildPreviewOperatingPayload(yearNo)
  previewDirty.value = false
  previewSaving.value = false
  previewSubmitting.value = false
  previewPageMessage.value = {
    type: 'info',
    text: '当前为开发预览模式，页面样式可查看，但不会读取或提交真实业务数据。',
  }
}

function buildPreviewView(yearNo: number, lastDraftSavedAt: string | null): PlayerOperatingView {
  const calculation = buildOperatingPreviewCalculation({
    payload: previewDraftPayload.value,
    carryForward: buildPreviewCarryForward(yearNo),
    currentStageCode: 'YEAR_END',
  })

  return {
    groupId: 1,
    yearNo,
    yearStatus: 'OPERATING',
    stageStatus: 'YEAR_END_OPEN',
    reportStatus: 'REPORT_LOCKED',
    businessStatus: 'NORMAL',
    currentStageCode: 'YEAR_END',
    canView: true,
    canEdit: true,
    canSubmit: true,
    hasInvalidDraft: false,
    invalidScopes: [],
    hasRetainedReportDraft: false,
    operatingPayload: cloneOperatingPayload(previewDraftPayload.value),
    editableScopes: ['YEAR_START', 'Q1', 'Q2', 'Q3', 'Q4', 'YEAR_END'],
    readonlyScopes: [],
    stageSubmitHistory: [
      { stageCode: 'Q1', submitVersion: 1, periodEndCash: 38 + yearNo, submitTime: new Date(2026, 2, 21, 9, 30).toISOString() },
      { stageCode: 'Q2', submitVersion: 1, periodEndCash: 43 + yearNo, submitTime: new Date(2026, 2, 21, 10, 30).toISOString() },
      { stageCode: 'Q3', submitVersion: 1, periodEndCash: 41 + yearNo, submitTime: new Date(2026, 2, 21, 11, 30).toISOString() },
    ],
    lastDraftSavedAt,
    quarterCashChecks: calculation.quarterCashChecks,
    derivedValues: calculation.derivedValues,
    periodEndCash: calculation.periodEndCash,
    carryForward: buildPreviewCarryForward(yearNo),
    noticeBoard: buildPreviewNoticeBoard(yearNo),
  }
}

function buildPreviewCarryForward(yearNo: number) {
  return {
    previousIncomeTax: 1 + yearNo,
    previousShortTermLoan: 8 + yearNo,
    previousLongTermLoan: 4 + yearNo,
    previousEquipmentResidual: 6 + yearNo,
    previousDepreciableAsset: 12 + yearNo,
    previousCash: 36 + yearNo,
    previousReceivable: 5 + yearNo,
    shareholderCapital: 50,
    retainedEarnings: 18 + yearNo,
  }
}

function buildPreviewNoticeBoard(yearNo: number) {
  return {
    pinnedNotice: {
      id: 1,
      kind: 'GENERAL',
      title: '系统通知',
      content: '主持人提示：请各小组按现场节奏推进经营，季度提交后关注贷款更新。',
      pinned: true,
      publishedAt: new Date(2026, 2, 21, 8, 45).toISOString(),
      yearNo,
      stageCode: null,
      amount: null,
    },
    recentList: [
      {
        id: 2,
        kind: 'REWARD',
        title: `${yearNo}年 Q2 奖励`,
        content: '示例：本季度经营表现优秀，奖励 3。',
        pinned: false,
        publishedAt: new Date(2026, 2, 21, 10, 0).toISOString(),
        yearNo,
        stageCode: 'Q2',
        amount: 3,
      },
      {
        id: 3,
        kind: 'GENERAL',
        title: '系统通知',
        content: '示例：准备进入下一阶段前，请再次核对经营数据。',
        pinned: false,
        publishedAt: new Date(2026, 2, 21, 9, 20).toISOString(),
        yearNo,
        stageCode: null,
        amount: null,
      },
    ],
  }
}
function buildPreviewOperatingPayload(yearNo: number): OperatingPayload {
  const payload = createEmptyOperatingPayload()
  payload.beginning.taxAndPlanning = {
    taxPayment: 5 + yearNo,
    planRevenue: 108 + yearNo * 8,
    comprehensiveCostPlan: 32 + yearNo * 2,
  }
  payload.beginning.marketBid = [
    { basicProductTotal: 18, standardProductTotal: 6, precisionProductTotal: 0, intelligentProductTotal: 0, marketInvestment: 2, orderAmount: 24 },
    { basicProductTotal: 12, standardProductTotal: 8, precisionProductTotal: 4, intelligentProductTotal: 0, marketInvestment: 3, orderAmount: 24 },
    { basicProductTotal: 0, standardProductTotal: 10, precisionProductTotal: 8, intelligentProductTotal: 4, marketInvestment: 4, orderAmount: 22 },
    { basicProductTotal: 0, standardProductTotal: 0, precisionProductTotal: 12, intelligentProductTotal: 14, marketInvestment: 5, orderAmount: 26 },
  ]
  payload.quarter.shortTermLoan = {
    q1: { dueRepayment: 2, interest: 1, newLoan: 5 },
    q2: { dueRepayment: 2, interest: 1, newLoan: 2 },
    q3: { dueRepayment: 2, interest: 1, newLoan: 3 },
    q4: { dueRepayment: 2, interest: 1, newLoan: 2 },
  }
  payload.quarter.materialPayment = {
    q1: { basicProduct: 5, standardProduct: 3, precisionProduct: 1, intelligentProduct: 0 },
    q2: { basicProduct: 4, standardProduct: 4, precisionProduct: 2, intelligentProduct: 1 },
    q3: { basicProduct: 3, standardProduct: 4, precisionProduct: 3, intelligentProduct: 1 },
    q4: { basicProduct: 2, standardProduct: 3, precisionProduct: 3, intelligentProduct: 2 },
  }
  payload.quarter.productionLineAdjustment = {
    q1: { changeProduct: 1, dismantleCost: 0, lineSale: 0, newLineInstall: 2, constructionToFixed: 1, newDepreciableAsset: 2 },
    q2: { changeProduct: 1, dismantleCost: 1, lineSale: 0, newLineInstall: 1, constructionToFixed: 1, newDepreciableAsset: 1 },
    q3: { changeProduct: 2, dismantleCost: 0, lineSale: 1, newLineInstall: 1, constructionToFixed: 1, newDepreciableAsset: 2 },
    q4: { changeProduct: 1, dismantleCost: 1, lineSale: 0, newLineInstall: 1, constructionToFixed: 0, newDepreciableAsset: 1 },
  }
  payload.quarter.humanResource = {
    q1: { staffCost: 2 },
    q2: { staffCost: 2 },
    q3: { staffCost: 3 },
    q4: { staffCost: 2 },
  }
  payload.quarter.salaryAndProduction = {
    q1: { salaryCost: 4 },
    q2: { salaryCost: 4 },
    q3: { salaryCost: 5 },
    q4: { salaryCost: 5 },
  }
  payload.quarter.researchAndManagement = {
    q1: { technologyResearch: 2, managementSystem: 1 },
    q2: { technologyResearch: 2, managementSystem: 1 },
    q3: { technologyResearch: 3, managementSystem: 1 },
    q4: { technologyResearch: 2, managementSystem: 1 },
  }
  payload.quarter.receivableUpdate = {
    q1: { receivableCollection: 6 },
    q2: { receivableCollection: 7 },
    q3: { receivableCollection: 5 },
    q4: { receivableCollection: 6 },
  }
  payload.quarter.deliverySettlement = {
    q1: { salesRevenue: 24, directCost: 14, managementStaffCost: 1 },
    q2: { salesRevenue: 27, directCost: 16, managementStaffCost: 1 },
    q3: { salesRevenue: 30, directCost: 17, managementStaffCost: 1 },
    q4: { salesRevenue: 35 + yearNo * 2, directCost: 20 + yearNo, managementStaffCost: 1 },
  }
  payload.yearEnd.longTermLoan = {
    interest: 2,
    repayment: 3,
    newLoan: 6,
  }
  payload.yearEnd.assetAdjustment = {
    lineMaintenance: 2,
    purchase: 16,
    sale: 4,
    rent: 3,
    workInConstruction: 5,
    marketCultivation: 2,
  }
  payload.extra.incomeAndPenalty = {
    q1: { discountExpense: 0, extraExpensePenalty: 0, extraIncomeReward: 1 },
    q2: { discountExpense: 0, extraExpensePenalty: 1, extraIncomeReward: 0 },
    q3: { discountExpense: 1, extraExpensePenalty: 0, extraIncomeReward: 1 },
    q4: { discountExpense: 0, extraExpensePenalty: 0, extraIncomeReward: 1 },
  }
  payload.derived.values = buildPreviewDerivedValues(yearNo)
  return payload
}

function buildPreviewDerivedValues(yearNo: number) {
  return {
    orderTotal: 96 + yearNo * 8,
    comprehensiveCostTotal: 34 + yearNo * 2,
    shortTermRepayment: 8,
    shortTermInterest: 4,
    newShortTermLoan: 12,
    materialPayment: 38,
    changeProductCost: 5,
    lineDismantleCost: 2,
    lineSaleValue: 5,
    newLineInstall: 5,
    humanResourceCost: 9,
    salaryAndProductionCost: 18,
    researchCost: 9,
    managementSystemCost: 4,
    managementSalary: 4,
    receivableRecovered: 24,
    salesRevenue: 116 + yearNo * 8,
    directCost: 67 + yearNo * 3,
    discountExpense: 1,
    extraExpensePenalty: 1,
    extraIncomeReward: 3,
    lineResidual: 16,
    depreciableAssetTotal: 21,
    depreciation: 7,
    financeIncomeExpense: 5,
    extraIncomeExpense: 2,
    factoryAssetChange: 12,
    factorySale: 4,
    receivableChange: 10,
    lineResidualChange: 3,
    depreciableAssetChange: 4,
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
  grid-template-columns: repeat(4, minmax(0, 1fr));
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

.loading-card {
  border: 1px dashed var(--line-strong);
  border-radius: 16px;
  background: #ffffff;
  padding: 28px;
  text-align: center;
  color: var(--muted);
}

.side-panel {
  position: sticky;
  top: 18px;
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













