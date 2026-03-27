<template>
  <aside class="sidebar">
    <section class="panel">
      <h3>页面状态</h3>
      <dl class="meta-list">
        <div>
          <dt>年份</dt>
          <dd>{{ yearLabel }}</dd>
        </div>
        <div>
          <dt>年度状态</dt>
          <dd>{{ view?.yearStatus || '--' }}</dd>
        </div>
        <div>
          <dt>当前阶段</dt>
          <dd>{{ view?.currentStageCode || '--' }}</dd>
        </div>
        <div>
          <dt>财报状态</dt>
          <dd>{{ view?.reportStatus || '--' }}</dd>
        </div>
        <div>
          <dt>经营状态</dt>
          <dd :class="{ danger: view?.businessStatus === 'BANKRUPT' }">{{ view?.businessStatus || '--' }}</dd>
        </div>
        <div>
          <dt>最近草稿</dt>
          <dd>{{ lastDraftSavedAt }}</dd>
        </div>
      </dl>
    </section>

    <section class="panel">
      <h3>提交控制</h3>
      <div class="action-list">
        <button type="button" class="btn" :disabled="saving || submitting || !view?.canEdit || !dirty" @click="$emit('save')">
          {{ saving ? '保存中...' : '保存草稿' }}
        </button>
        <button type="button" class="btn primary" :disabled="saving || submitting || !view?.canSubmit" @click="$emit('submit')">
          {{ submitting ? '提交中...' : `提交 ${view?.currentStageCode || ''}` }}
        </button>
      </div>
      <p class="hint">阶段提交成功后，会提醒玩家关注贷款更新。</p>
    </section>

    <section class="panel">
      <h3>提交历史</h3>
      <ul class="history-list">
        <li v-if="!view?.stageSubmitHistory?.length" class="empty">当前年份还没有正式提交记录</li>
        <li v-for="item in view?.stageSubmitHistory" :key="`${item.stageCode}-${item.submitVersion}`">
          <strong>{{ item.stageCode }}</strong>
          <span>v{{ item.submitVersion }}</span>
          <span>期末现金 {{ item.periodEndCash.toLocaleString('zh-CN') }}</span>
          <span>{{ formatTime(item.submitTime) }}</span>
        </li>
      </ul>
    </section>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import type { PlayerOperatingView } from '@/types/sandbox-game'

const props = defineProps<{
  view: PlayerOperatingView | null
  yearLabel: string
  dirty: boolean
  saving: boolean
  submitting: boolean
}>()

defineEmits<{
  (event: 'save'): void
  (event: 'submit'): void
}>()

const lastDraftSavedAt = computed(() => {
  const value = props.view?.lastDraftSavedAt
  if (!value) {
    return '尚未保存'
  }
  return formatTime(value)
})

function formatTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString('zh-CN', { hour12: false })
}
</script>

<style scoped>
.sidebar {
  display: grid;
  gap: 14px;
}

.panel {
  border: 1px solid var(--line);
  border-radius: 16px;
  background: #ffffff;
  padding: 16px;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.04);
}

.panel h3 {
  margin: 0 0 14px;
  font-size: 16px;
}

.meta-list {
  display: grid;
  gap: 10px;
  margin: 0;
}

.meta-list div {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.meta-list dt {
  color: var(--muted);
}

.meta-list dd {
  margin: 0;
  text-align: right;
  font-weight: 700;
}

.meta-list dd.danger {
  color: var(--danger);
}

.action-list {
  display: grid;
  gap: 10px;
}

.btn {
  border: 1px solid var(--line-strong);
  border-radius: 12px;
  padding: 12px 14px;
  background: #ffffff;
}

.btn.primary {
  background: var(--accent);
  border-color: var(--accent);
  color: #ffffff;
}

.hint {
  margin: 12px 0 0;
  font-size: 12px;
  color: var(--muted);
  line-height: 1.6;
}

.history-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 10px;
}

.history-list li {
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 10px 12px;
  display: grid;
  gap: 4px;
}

.history-list .empty {
  color: var(--muted);
}
</style>