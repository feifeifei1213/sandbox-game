<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>通知与奖惩</h2>
        
      </div>
      <div class="hero-actions">
        <button type="button" class="btn" :disabled="loading" @click="handleRefresh">刷新记录</button>
      </div>
    </header>

    <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
      {{ pageMessage.text }}
    </section>

    <div class="page-grid">
      <section class="panel-card form-card">
        <div class="panel-head">
          <strong>普通通知发送区</strong>
          
        </div>
        <div class="form-body">
          <label class="field">
            <span>目标范围</span>
            <select v-model="store.generalForm.targetScope">
              <option value="ALL">全体小组</option>
              <option value="GROUP">单个小组</option>
            </select>
          </label>

          <label v-if="store.generalForm.targetScope === 'GROUP'" class="field">
            <span>目标小组</span>
            <select v-model.number="store.generalForm.targetGroupId">
              <option :value="null" disabled>请选择小组</option>
              <option v-for="item in groups" :key="item.groupId" :value="item.groupId">第{{ item.groupNo }}组</option>
            </select>
          </label>

          <label class="field checkbox-field">
            <input v-model="store.generalForm.pinned" type="checkbox" />
            <span>设为置顶通知</span>
          </label>

          <label class="field">
            <span>通知内容</span>
            <textarea v-model="store.generalForm.content" rows="6" placeholder="例如：准备开启下一年，请所有小组尽快完成本年经营与财报。"></textarea>
          </label>

          <div class="form-actions">
            <button type="button" class="btn primary" :disabled="sendingGeneral" @click="handleSendGeneral">
              {{ sendingGeneral ? '发送中...' : '发送普通通知' }}
            </button>
          </div>
        </div>
      </section>

      <section class="panel-card form-card">
        <div class="panel-head">
          <strong>奖惩下发区</strong>
          
        </div>
        <div class="form-body">
          <label class="field">
            <span>目标小组</span>
            <select v-model.number="store.adjustmentForm.groupId">
              <option :value="null" disabled>请选择小组</option>
              <option v-for="item in groups" :key="item.groupId" :value="item.groupId">第{{ item.groupNo }}组</option>
            </select>
          </label>

          <div class="two-col-grid">
            <label class="field">
              <span>年份</span>
              <select v-model.number="store.adjustmentForm.yearNo">
                <option v-for="item in yearOptions" :key="item.value" :value="item.value">{{ item.label }}</option>
              </select>
            </label>

            <label class="field">
              <span>季度</span>
              <select v-model="store.adjustmentForm.stageCode">
                <option value="Q1">Q1</option>
                <option value="Q2">Q2</option>
                <option value="Q3">Q3</option>
                <option value="Q4">Q4</option>
              </select>
            </label>
          </div>

          <div class="two-col-grid">
            <label class="field">
              <span>类型</span>
              <select v-model="store.adjustmentForm.adjustmentType">
                <option value="REWARD">奖励</option>
                <option value="PENALTY">罚款</option>
              </select>
            </label>

            <label class="field">
              <span>金额</span>
              <input
                v-model="store.adjustmentForm.amount"
                type="number"
                min="0"
                step="1"
                inputmode="numeric"
                placeholder="请输入整数金额"
                :class="{ invalid: hasFractionInput(store.adjustmentForm.amount) }"
              />
            </label>
          </div>

          <label class="field">
            <span>原因说明</span>
            <textarea v-model="store.adjustmentForm.reason" rows="6" placeholder="例如：本季度市场竞标表现优异，奖励 5。"></textarea>
          </label>

          <div class="form-actions">
            <button type="button" class="btn primary" :disabled="sendingAdjustment" @click="handleSendAdjustment">
              {{ sendingAdjustment ? '下发中...' : '下发奖惩' }}
            </button>
          </div>
        </div>
      </section>
    </div>

    <div class="record-grid">
      <section class="panel-card">
        <div class="panel-head">
          <strong>最近普通通知</strong>
          
        </div>
        <div class="table-scroll">
          <table class="record-table">
            <thead>
              <tr>
                <th>时间</th>
                <th>目标</th>
                <th>内容</th>
                <th>置顶</th>
                <th>发送人</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="generalNotices.length === 0">
                <td colspan="5" class="empty-row">当前还没有普通通知记录。</td>
              </tr>
              <tr v-for="item in generalNotices" :key="item.id">
                <td>{{ formatDateTime(item.publishedAt) }}</td>
                <td>{{ formatNoticeTarget(item.targetScope, item.targetGroupName) }}</td>
                <td class="content-cell">{{ item.content }}</td>
                <td>{{ item.pinned ? '是' : '否' }}</td>
                <td>{{ item.operatorName }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="panel-card">
        <div class="panel-head">
          <strong>最近奖惩记录</strong>
          
        </div>
        <div class="table-scroll">
          <table class="record-table">
            <thead>
              <tr>
                <th>时间</th>
                <th>小组</th>
                <th>年份</th>
                <th>季度</th>
                <th>类型</th>
                <th>金额</th>
                <th>原因</th>
                <th>发送人</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="adjustments.length === 0">
                <td colspan="8" class="empty-row">当前还没有奖惩下发记录。</td>
              </tr>
              <tr v-for="item in adjustments" :key="item.id">
                <td>{{ formatDateTime(item.publishedAt) }}</td>
                <td>第{{ item.groupNo }}组</td>
                <td>{{ item.yearNo }}年</td>
                <td>{{ item.stageCode }}</td>
                <td>{{ item.adjustmentType === 'REWARD' ? '奖励' : '罚款' }}</td>
                <td class="number-cell">{{ formatAmount(item.amount) }}</td>
                <td class="content-cell">{{ item.reason }}</td>
                <td>{{ item.operatorName }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'

import { useAdminShellStore } from '@/stores/admin-shell'
import { useAdminNoticeStore } from '@/stores/admin-notice'
import { hasFractionInput } from '@/utils/manual-integer'

const shellStore = useAdminShellStore()
const store = useAdminNoticeStore()
const { config } = storeToRefs(shellStore)
const { groups, generalNotices, adjustments, loading, sendingGeneral, sendingAdjustment, pageMessage } = storeToRefs(store)

const yearOptions = computed(() => {
  const currentOpenYear = Math.max(config.value?.currentOpenYear ?? 0, 0)
  return Array.from({ length: currentOpenYear + 1 }, (_, index) => ({
    value: index,
    label: `${index}年`,
  }))
})

onMounted(async () => {
  try {
    if (!shellStore.config) {
      await shellStore.bootstrap()
    }
    await store.bootstrap(config.value?.currentOpenYear ?? 0)
  } catch {
    // 页面消息由 store 统一处理。
  }
})

watch(
  () => config.value?.currentOpenYear,
  (value) => {
    if (typeof value === 'number') {
      store.adjustmentForm.yearNo = value
    }
  },
)

async function handleRefresh() {
  try {
    await store.refreshRecords()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleSendGeneral() {
  try {
    await store.sendGeneral()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleSendAdjustment() {
  try {
    await store.sendAdjustmentNotice()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

function formatDateTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString('zh-CN', { hour12: false })
}

function formatNoticeTarget(scope: string, groupName: string | null) {
  if (scope === 'GROUP') {
    return groupName ? groupName : '单组'
  }
  return '全体小组'
}

function formatAmount(value: number) {
  return value.toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
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

.page-grid,
.record-grid {
  display: grid;
  gap: 14px;
}

.page-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.record-grid {
  grid-template-columns: 1fr;
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

.form-body {
  display: grid;
  gap: 14px;
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
.field textarea,
.field input {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #ffffff;
  padding: 10px 12px;
  font: inherit;
}

.field input.invalid {
  border-color: var(--danger);
  background: #fff5f5;
  color: var(--danger);
}

.checkbox-field {
  display: flex;
  align-items: center;
  gap: 10px;
}

.checkbox-field input {
  width: auto;
}

.two-col-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.table-scroll {
  overflow: auto;
}

.record-table {
  width: 100%;
  min-width: 840px;
  border-collapse: collapse;
}

.record-table th,
.record-table td {
  border: 1px solid var(--line);
  padding: 10px 12px;
  font-size: 14px;
}

.record-table th {
  background: #f4f6f9;
  text-align: center;
}

.content-cell {
  min-width: 220px;
  line-height: 1.6;
}

.number-cell {
  text-align: right;
}

.empty-row {
  text-align: center;
  color: var(--muted);
}

@media (max-width: 1200px) {
  .page-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 900px) {
  .hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .two-col-grid {
    grid-template-columns: 1fr;
  }
}
</style>
