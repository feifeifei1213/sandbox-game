<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>赛前配置</h2>
        <p>先确认本场比赛的小组数量，再由系统一次性生成小组、账号和年份状态。</p>
      </div>
      <div class="hero-actions">
        <button type="button" class="btn" :disabled="shellLoading || initializing" @click="handleRefresh">刷新状态</button>
        <button
          v-if="setupStatus?.initialized"
          type="button"
          class="btn primary"
          @click="goToSummary"
        >
          进入汇总页
        </button>
      </div>
    </header>

    <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
      {{ pageMessage.text }}
    </section>

    <div class="cards-grid">
      <article class="card">
        <span class="card-label">初始化状态</span>
        <strong>{{ setupStatus?.initialized ? '已完成' : '未完成' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">当前小组数</span>
        <strong>{{ setupStatus?.groupCount ?? '--' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">最终年份</span>
        <strong>{{ setupStatus?.finalYear ?? '--' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">当前开放年份</span>
        <strong>{{ setupStatus?.currentOpenYear ?? '--' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">沙盘版本</span>
        <strong>{{ setupStatus?.editionName ?? '--' }}</strong>
      </article>
    </div>

    <div class="split-layout">
      <section class="panel-card">
        <div class="panel-head">
          <div>
            <strong>初始化比赛环境</strong>
            <span>首版只支持赛前初始化配置小组数量，比赛开始后不支持动态增减。</span>
          </div>
        </div>

        <div class="form-grid">
          <label class="field">
            <span>沙盘版本</span>
            <select
              v-model="editionCodeDraft"
              :disabled="setupStatus?.initialized || initializing || availableEditions.length === 0"
            >
              <option v-for="item in availableEditions" :key="item.editionCode" :value="item.editionCode">
                {{ item.editionName }}
              </option>
            </select>
          </label>
          <label class="field">
            <span>小组数量</span>
            <input
              v-model.number="groupCountDraft"
              type="number"
              min="1"
              max="10"
              step="1"
              :disabled="setupStatus?.initialized || initializing"
            >
          </label>
          <label class="field">
            <span>默认玩家账号</span>
            <input :value="accountPreviewText" readonly>
          </label>
        </div>

        <div class="note-list">
          <span>管理员账号继续保留为 `admin`。</span>
          <span>玩家账号将按 `group01 ~ groupNN` 自动生成，默认密码沿用 `123456`。</span>
          <span>沙盘版本只允许赛前选择，初始化完成后锁定。</span>
          <span>初始化完成后默认打开 `0年`，正式年份预置但保持锁定。</span>
        </div>

        <div class="action-row">
          <button
            type="button"
            class="btn primary"
            :disabled="setupStatus?.initialized || initializing"
            @click="handleInitialize"
          >
            {{ initializing ? '初始化中...' : '确认初始化比赛' }}
          </button>
        </div>
      </section>

      <section class="panel-card">
        <div class="panel-head">
          <div>
            <strong>当前说明</strong>
            <span>初始化完成后，汇总页、组数据页、通知与奖惩范围都会按实际小组数量动态适配。</span>
          </div>
        </div>

        <div class="timeline">
          <div class="timeline-item">
            <strong>默认入口</strong>
            <span>{{ setupStatus?.defaultRoute ?? '/sandbox-game/admin/setup' }}</span>
          </div>
          <div class="timeline-item">
            <strong>共享基线</strong>
            <span>{{ setupStatus?.initialBaselineSubmitted ? '已提交' : '未提交' }}</span>
          </div>
          <div class="timeline-item">
            <strong>当前版本</strong>
            <span>{{ currentEditionDescription }}</span>
          </div>
          <div class="timeline-item">
            <strong>字段模板</strong>
            <span>{{ currentEditionTemplateText }}</span>
          </div>
          <div class="timeline-item">
            <strong>初始化后</strong>
            <span>系统会生成小组主数据、玩家账号与全部年份主状态数据。</span>
          </div>
          <div class="timeline-item">
            <strong>风险边界</strong>
            <span>若当前环境已初始化，再次提交会被服务端拒绝，不会覆盖现有比赛数据。</span>
          </div>
        </div>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'

import { initializeAdminGame } from '@/api/sandbox-game/admin-control'
import { useAdminShellStore, type PageMessage } from '@/stores/admin-shell'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const shellStore = useAdminShellStore()
const authStore = useAuthStore()
const { config, setupStatus, loading: shellLoading } = storeToRefs(shellStore)

const groupCountDraft = ref(10)
const editionCodeDraft = ref('VIP_SERVICE_V1')
const initializing = ref(false)
const pageMessage = ref<PageMessage | null>(null)

const availableEditions = computed(() => setupStatus.value?.availableEditions ?? [])

const accountPreviewText = computed(() => {
  const count = normalizeGroupCount(groupCountDraft.value)
  if (count <= 1) {
    return 'group01'
  }
  return `group01 ~ group${String(count).padStart(2, '0')}`
})

const selectedEdition = computed(() =>
  availableEditions.value.find((item) => item.editionCode === editionCodeDraft.value)
  ?? availableEditions.value.find((item) => item.defaultEdition)
  ?? null,
)

const currentEditionDescription = computed(() => {
  const edition = selectedEdition.value
  if (edition) {
    return `${edition.editionName} · ${edition.description}`
  }
  return setupStatus.value?.editionName ?? '--'
})

const currentEditionTemplateText = computed(() => {
  const edition = selectedEdition.value
  if (!edition) {
    return setupStatus.value?.templateVersion ?? '--'
  }
  return `${edition.operatingTemplateVersion} / ${edition.reportTemplateVersion} / ${edition.orderTemplateVersion}`
})

onMounted(async () => {
  try {
    if (!setupStatus.value || !config.value) {
      await shellStore.bootstrap()
    }
    syncGroupCountDraft()
    syncEditionCodeDraft()
  } catch {
    // 错误消息由 shell store 统一展示。
  }
})

watch(
  () => setupStatus.value?.groupCount,
  () => {
    syncGroupCountDraft()
  },
)

watch(
  () => [setupStatus.value?.editionCode, setupStatus.value?.availableEditions?.length],
  () => {
    syncEditionCodeDraft()
  },
)

async function handleRefresh() {
  pageMessage.value = null
  try {
    await shellStore.refreshAll({ silent: true })
    syncGroupCountDraft()
    syncEditionCodeDraft()
  } catch {
    pageMessage.value = { type: 'error', text: '刷新赛前配置状态失败。' }
  }
}

async function handleInitialize() {
  initializing.value = true
  pageMessage.value = null
  try {
    const result = await initializeAdminGame(normalizeGroupCount(groupCountDraft.value), editionCodeDraft.value)
    await shellStore.refreshAll({ silent: true })
    await authStore.refreshCurrentUser()
    pageMessage.value = {
      type: 'success',
      text: `比赛初始化完成，已生成 ${result.createdGroupCount} 个小组、${result.createdAccountCount} 个玩家账号和 ${result.createdYearStateCount} 条年份状态。`,
    }
    await router.replace(authStore.resolveDefaultRoute())
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '初始化比赛失败')
  } finally {
    initializing.value = false
  }
}

async function goToSummary() {
  await router.replace('/sandbox-game/admin/summary')
}

function syncGroupCountDraft() {
  const nextGroupCount = setupStatus.value?.groupCount
  groupCountDraft.value = nextGroupCount && nextGroupCount > 0 ? nextGroupCount : 10
}

function syncEditionCodeDraft() {
  const statusEditionCode = setupStatus.value?.editionCode
  if (statusEditionCode) {
    editionCodeDraft.value = statusEditionCode
    return
  }
  const defaultEdition = availableEditions.value.find((item) => item.defaultEdition) ?? availableEditions.value[0]
  editionCodeDraft.value = defaultEdition?.editionCode ?? 'VIP_SERVICE_V1'
}

function normalizeGroupCount(value: number) {
  const rounded = Number.isFinite(value) ? Math.round(value) : 10
  return Math.min(10, Math.max(1, rounded))
}

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
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
.action-row {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.btn {
  border: 1px solid var(--line);
  background: #ffffff;
  border-radius: 12px;
  padding: 10px 16px;
  cursor: pointer;
}

.btn.primary {
  background: var(--accent);
  color: #ffffff;
  border-color: var(--accent);
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
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

.cards-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 12px;
}

.card,
.panel-card {
  border: 1px solid var(--line);
  border-radius: 18px;
  background: #ffffff;
  padding: 18px;
}

.card {
  display: grid;
  gap: 8px;
}

.card-label {
  color: var(--muted);
  font-size: 13px;
}

.card strong {
  font-size: 24px;
}

.split-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(320px, 0.8fr);
  gap: 16px;
}

.panel-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
}

.panel-head strong {
  display: block;
  margin-bottom: 6px;
  font-size: 18px;
}

.panel-head span {
  color: var(--muted);
  line-height: 1.6;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.field {
  display: grid;
  gap: 8px;
}

.field span {
  font-size: 13px;
  color: var(--muted);
}

.field input,
.field select {
  width: 100%;
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid var(--line);
  background: #fbfcfe;
}

.note-list,
.timeline {
  display: grid;
  gap: 10px;
  margin-top: 18px;
}

.note-list span,
.timeline-item span {
  color: var(--muted);
  line-height: 1.6;
}

.timeline-item {
  display: grid;
  gap: 6px;
  padding: 12px 14px;
  border-radius: 14px;
  background: #f7f9fc;
}

.timeline-item strong {
  font-size: 14px;
}

@media (max-width: 1180px) {
  .cards-grid,
  .split-layout,
  .form-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 860px) {
  .hero {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
