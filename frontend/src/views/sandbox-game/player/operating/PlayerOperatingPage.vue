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
          <span class="pill">组别：{{ currentView?.groupId ?? '--' }}</span>
          <span class="pill">开放年份：{{ currentConfig?.currentOpenYear ?? '--' }}</span>
          <span class="pill">最终年份：{{ currentConfig?.finalYear ?? '--' }}</span>
        </div>
      </header>

      <section class="toolbar-card">
        <div class="toolbar-top">
          <YearTabs :tabs="yearTabs" :active-year="selectedYear" @select="handleYearSelect" />
          <PageModeSwitch :report-enabled="reportEnabled" @report="goReport" />
        </div>
      </section>

      <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
        {{ pageMessage.text }}
      </section>

      <div class="workspace">
        <main class="main-panel">
          <section class="status-banner">
            <div class="status-item">
              <span>当前年份</span>
              <strong>{{ selectedYear }} 年</strong>
            </div>
            <div class="status-item">
              <span>年度状态</span>
              <strong>{{ currentView?.yearStatus || '--' }}</strong>
            </div>
            <div class="status-item">
              <span>经营阶段</span>
              <strong>{{ currentView?.currentStageCode || '--' }}</strong>
            </div>
            <div class="status-item">
              <span>财报状态</span>
              <strong>{{ currentView?.reportStatus || '--' }}</strong>
            </div>
          </section>

          <section v-if="loading || yearViewLoading" class="loading-card">正在加载经营页数据...</section>

          <OperatingSheet
            v-else-if="currentView"
            :model-value="draftPayload"
            :editable-scopes="currentView.editableScopes"
            :quarter-cash-checks="currentView.quarterCashChecks"
            :derived-values="currentView.derivedValues"
            :period-end-cash="currentView.periodEndCash"
            @update:model-value="updatePayload"
          />

          <section v-else class="loading-card">当前没有可展示的经营页数据。</section>
        </main>

        <OperatingSidebar
          class="side-panel"
          :view="currentView"
          :year-label="`${selectedYear} 年经营`"
          :dirty="dirty"
          :saving="saving"
          :submitting="submitting"
          @save="handleSave"
          @submit="handleSubmit"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'

import YearTabs from '@/components/sandbox-game/common/YearTabs.vue'
import OperatingSheet from '@/components/sandbox-game/player/OperatingSheet.vue'
import OperatingSidebar from '@/components/sandbox-game/player/OperatingSidebar.vue'
import PageModeSwitch from '@/components/sandbox-game/player/PageModeSwitch.vue'
import { usePlayerOperatingStore } from '@/stores/player-operating'
import type { OperatingPayload } from '@/types/sandbox-game'

const AUTO_SAVE_INTERVAL = 5 * 60 * 1000

const route = useRoute()
const router = useRouter()
const store = usePlayerOperatingStore()
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
} = storeToRefs(store)

let initialized = false
let autoSaveTimer = 0

onMounted(async () => {
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
  () => route.query.yearNo,
  async () => {
    const targetYear = readRouteYear()
    if (!initialized || targetYear === undefined || targetYear === selectedYear.value) {
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
  if (yearNo === selectedYear.value) {
    return
  }
  syncRouteYear(yearNo)
}

function updatePayload(nextPayload: OperatingPayload) {
  store.updateDraft(nextPayload)
}

async function handleSave() {
  await store.saveDraft()
}

async function handleSubmit() {
  await store.submitCurrentStage()
}

function goReport() {
  if (!reportEnabled.value) {
    return
  }
  router.push({
    path: '/sandbox-game/player/report',
    query: { yearNo: String(selectedYear.value) },
  })
}
</script>

<style scoped>
.player-page {
  min-height: 100vh;
  padding: 20px;
}

.shell {
  width: min(1680px, calc(100vw - 40px));
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
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 18px;
  padding: 18px;
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