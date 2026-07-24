<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>年度控制</h2>
        
      </div>
      <div class="hero-actions">
        <button type="button" class="btn" :disabled="shellLoading || updatingFinalYear || openingNextYear" @click="handleRefresh">刷新状态</button>
        <button type="button" class="btn primary" :disabled="!config?.canOpenNextYear || openingNextYear" @click="handleOpenNextYear">
          开放下一年
        </button>
      </div>
    </header>

    <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
      {{ pageMessage.text }}
    </section>

    <div class="cards-grid">
      <article class="card">
        <span class="card-label">最终年份</span>
        <strong>{{ config?.finalYear ?? '--' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">当前开放年份</span>
        <strong>{{ config?.currentOpenYear ?? '--' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">下一可开放年份</span>
        <strong>{{ config?.nextOpenableYear ?? '--' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">共享基线</span>
        <strong>{{ config?.initialBaselineSubmitted ? '已提交' : '未提交' }}</strong>
      </article>
    </div>

    <div class="split-layout">
      <section class="panel-card">
        <div class="panel-head">
          <div>
            <strong>控制区</strong>
            
          </div>
        </div>

        <div class="form-grid">
          <label class="field">
            <span>最终年份</span>
            <input
              v-model.number="finalYearDraft"
              type="number"
              min="0"
              step="1"
              inputmode="numeric"
              :class="{ invalid: hasFractionInput(finalYearDraft) }"
              :disabled="updatingFinalYear || shellLoading"
              data-enter-confirm
              @keydown.enter="confirmInputOnEnter"
            >
          </label>
          <label class="field">
            <span>准备开放到</span>
            <input :value="config ? `${config.nextOpenableYear}年` : '--'" readonly>
          </label>
        </div>

        <div class="state-line">
          <span class="status-tag" :class="config?.canOpenNextYear ? 'ok' : 'warn'">
            {{ config?.canOpenNextYear ? '可以开放下一年' : '当前不可开放' }}
          </span>
          
        </div>

        <div class="action-row">
          <button type="button" class="btn" :disabled="updatingFinalYear || shellLoading" @click="handleUpdateFinalYear">保存最终年份</button>
          <button type="button" class="btn primary" :disabled="!config?.canOpenNextYear || openingNextYear" @click="handleOpenNextYear">
            确认开放 {{ config?.nextOpenableYear ?? '--' }} 年
          </button>
        </div>
      </section>

      <section class="panel-card">
        <div class="panel-head">
          <div>
            <strong>阻断原因与最近动作</strong>
            
          </div>
        </div>

        <div class="timeline">
          <div class="timeline-item">
            <strong>当前阻断原因</strong>
            <span>{{ config?.openNextYearBlockedReason || '无阻断，可推进下一年。' }}</span>
          </div>
          <div class="timeline-item">
            <strong>规则版本</strong>
            <span>{{ config?.ruleVersion ?? '--' }}</span>
          </div>
          <div class="timeline-item">
            <strong>模板版本</strong>
            <span>{{ config?.templateVersion ?? '--' }}</span>
          </div>
          <div class="timeline-item">
            <strong>最近管理员动作</strong>
            <span>{{ latestAdminActionText }}</span>
          </div>
        </div>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'

import { openAdminNextYear, updateAdminFinalYear } from '@/api/sandbox-game/admin-control'
import { useAdminShellStore, type PageMessage } from '@/stores/admin-shell'
import { confirmInputOnEnter } from '@/utils/input-navigation'
import { hasFractionInput } from '@/utils/manual-integer'

const shellStore = useAdminShellStore()
const { config, loading: shellLoading } = storeToRefs(shellStore)

const finalYearDraft = ref(0)
const updatingFinalYear = ref(false)
const openingNextYear = ref(false)
const pageMessage = ref<PageMessage | null>(null)

const latestAdminActionText = computed(() => {
  const action = config.value?.latestAdminAction
  if (!action) {
    return '暂无'
  }
  return `${action.operatorName} · ${formatDateTime(action.operateTime)}`
})

onMounted(async () => {
  try {
    if (!config.value) {
      await shellStore.bootstrap()
    }
    syncFinalYearDraft()
  } catch {
    // 错误消息由 shell store 统一处理。
  }
})

watch(
  () => config.value?.finalYear,
  () => {
    syncFinalYearDraft()
  },
)

async function handleRefresh() {
  pageMessage.value = null
  try {
    await shellStore.refreshConfig({ silent: true })
    syncFinalYearDraft()
  } catch {
    pageMessage.value = { type: 'error', text: '刷新年度控制状态失败。' }
  }
}

async function handleUpdateFinalYear() {
  if (!config.value) {
    return
  }
  if (hasFractionInput(finalYearDraft.value)) {
    pageMessage.value = { type: 'error', text: '最终年份必须填写整数。' }
    return
  }
  updatingFinalYear.value = true
  pageMessage.value = null
  try {
    const result = await updateAdminFinalYear(finalYearDraft.value)
    await shellStore.refreshConfig({ silent: true })
    syncFinalYearDraft()
    pageMessage.value = {
      type: 'success',
      text: buildUpdateSuccessText(result.initializedYearCount, result.finalYear),
    }
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '更新最终年份失败')
  } finally {
    updatingFinalYear.value = false
  }
}

async function handleOpenNextYear() {
  if (!config.value) {
    return
  }
  openingNextYear.value = true
  pageMessage.value = null
  try {
    const result = await openAdminNextYear(config.value.nextOpenableYear)
    await shellStore.refreshConfig({ silent: true })
    syncFinalYearDraft()
    pageMessage.value = {
      type: 'success',
      text: `${result.openedYearNo} 年已开放。当前开放年份更新为 ${result.currentOpenYear} 年。`,
    }
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '开放下一年失败')
  } finally {
    openingNextYear.value = false
  }
}

function syncFinalYearDraft() {
  finalYearDraft.value = config.value?.finalYear ?? 0
}

function buildUpdateSuccessText(initializedYearCount: number, finalYear: number) {
  if (initializedYearCount > 0) {
    return `最终年份已更新为 ${finalYear} 年，并补齐 ${initializedYearCount} 个未来年份状态。`
  }
  return `最终年份已更新为 ${finalYear} 年。`
}

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
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
.action-row {
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

.cards-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.card,
.panel-card {
  border: 1px solid var(--line);
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.05);
}

.card {
  padding: 16px;
}

.card-label {
  display: block;
  color: var(--muted);
  font-size: 12px;
  margin-bottom: 8px;
}

.card strong {
  font-size: 20px;
}

.panel-head {
  padding: 14px 16px;
  border-bottom: 1px solid var(--line);
  background: #f7f9fc;
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

.split-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(320px, 0.9fr);
  gap: 14px;
}

.panel-card {
  padding-bottom: 16px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  padding: 16px;
}

.field {
  display: grid;
  gap: 6px;
}

.field span {
  font-size: 13px;
  font-weight: 700;
  color: var(--muted);
}

.field input {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #ffffff;
  padding: 10px 12px;
}

.field input.invalid {
  border-color: #e17979;
  background: #fff6f6;
  color: var(--danger);
}

.state-line {
  display: grid;
  gap: 10px;
  padding: 0 16px 16px;
}

.inline-tip {
  color: var(--muted);
  font-size: 13px;
}

.status-tag {
  display: inline-flex;
  align-items: center;
  width: fit-content;
  padding: 6px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
}

.status-tag.ok {
  background: #e8f7ee;
  color: var(--success);
}

.status-tag.warn {
  background: #fff7ed;
  color: var(--warning);
}

.timeline {
  display: grid;
  gap: 12px;
  padding: 16px;
}

.timeline-item {
  border-left: 3px solid #c9d7f2;
  padding-left: 12px;
}

.timeline-item strong {
  display: block;
  margin-bottom: 4px;
}

.timeline-item span {
  color: var(--muted);
  line-height: 1.6;
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

@media (max-width: 1240px) {
  .cards-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .split-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 900px) {
  .hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .cards-grid,
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>

