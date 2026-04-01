<template>
  <section class="panel notice-panel">
    <div class="notice-head">
      <h3>通知区</h3>
      <span>{{ hasNotice ? '管理员最新消息' : '当前暂无消息' }}</span>
    </div>

    <article v-if="noticeBoard?.pinnedNotice" class="pinned-card">
      <div class="notice-label-row">
        <strong class="notice-label notice-label-pinned">置顶通知</strong>
        <span>{{ formatTime(noticeBoard.pinnedNotice.publishedAt) }}</span>
      </div>
      <strong class="notice-title">{{ noticeBoard.pinnedNotice.title }}</strong>
      <p class="notice-content">{{ noticeBoard.pinnedNotice.content }}</p>
    </article>

    <ul v-if="noticeBoard?.recentList?.length" class="notice-list">
      <li v-for="item in noticeBoard.recentList" :key="`${item.kind}-${item.id}`">
        <div class="notice-label-row">
          <strong class="notice-label" :class="resolveNoticeClass(item.kind)">{{ resolveNoticeLabel(item.kind) }}</strong>
          <span>{{ formatTime(item.publishedAt) }}</span>
        </div>
        <strong class="notice-title">{{ resolveNoticeTitle(item) }}</strong>
        <p class="notice-content">{{ item.content || '无补充说明' }}</p>
        <p v-if="item.amount !== null && item.amount !== undefined" class="notice-meta">
          金额：{{ formatAmount(item.amount) }}
        </p>
      </li>
    </ul>

    <p v-else class="empty-text">主持人暂未发送新的通知或奖惩消息。</p>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import type { PlayerNoticeBoard, PlayerNoticeItem } from '@/types/sandbox-game'

const props = defineProps<{
  noticeBoard: PlayerNoticeBoard | null | undefined
}>()

const hasNotice = computed(() => Boolean(props.noticeBoard?.pinnedNotice || props.noticeBoard?.recentList?.length))

function resolveNoticeLabel(kind: string) {
  if (kind === 'REWARD') {
    return '奖励'
  }
  if (kind === 'PENALTY') {
    return '罚款'
  }
  return '通知'
}

function resolveNoticeClass(kind: string) {
  if (kind === 'REWARD') {
    return 'notice-label-reward'
  }
  if (kind === 'PENALTY') {
    return 'notice-label-penalty'
  }
  return 'notice-label-general'
}

function resolveNoticeTitle(item: PlayerNoticeItem) {
  if (item.kind === 'GENERAL') {
    return item.title
  }
  const yearText = typeof item.yearNo === 'number' ? `${item.yearNo}年` : ''
  const stageText = item.stageCode ? ` ${item.stageCode}` : ''
  return `${yearText}${stageText} ${resolveNoticeLabel(item.kind)}`.trim()
}

function formatTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString('zh-CN', { hour12: false })
}

function formatAmount(value: number) {
  return value.toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}
</script>

<style scoped>
.notice-panel {
  border-color: #d7e1ef;
  background: linear-gradient(180deg, #fbfdff 0%, #f4f8ff 100%);
}

.notice-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: baseline;
  margin-bottom: 14px;
}

.notice-head span {
  color: var(--muted);
  font-size: 12px;
}

.pinned-card,
.notice-list li {
  border: 1px solid var(--line);
  border-radius: 14px;
  background: #ffffff;
  padding: 12px;
}

.pinned-card {
  margin-bottom: 12px;
  border-color: #d4bb78;
  background: #fffaf0;
}

.notice-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 10px;
}

.notice-label-row {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.notice-label-row span {
  color: var(--muted);
  font-size: 12px;
}

.notice-label {
  display: inline-flex;
  align-items: center;
  padding: 4px 8px;
  border-radius: 999px;
  font-size: 11px;
  line-height: 1;
}

.notice-label-general {
  background: #edf4ff;
  color: #2458aa;
}

.notice-label-pinned {
  background: #f7e4ac;
  color: #875c00;
}

.notice-label-reward {
  background: #e9f7ee;
  color: #1f6b40;
}

.notice-label-penalty {
  background: #fff1f1;
  color: #b24040;
}

.notice-title {
  display: block;
  font-size: 13px;
  margin-bottom: 6px;
}

.notice-content,
.notice-meta,
.empty-text {
  margin: 0;
  color: var(--muted);
  font-size: 12px;
  line-height: 1.7;
}

.notice-meta {
  margin-top: 6px;
}
</style>
