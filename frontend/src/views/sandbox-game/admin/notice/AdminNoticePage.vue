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
            <select v-model.number="store.adjustmentForm.groupId" @change="handleAdjustmentGroupChange">
              <option :value="null" disabled>请选择小组</option>
              <option v-for="item in groups" :key="item.groupId" :value="item.groupId">第{{ item.groupNo }}组</option>
            </select>
          </label>

          <div class="two-col-grid">
            <div class="field readonly-field">
              <span>系统归属年份</span>
              <div class="readonly-value">{{ operationContext ? `${operationContext.operationYearNo}年` : '--' }}</div>
            </div>

            <div class="field readonly-field">
              <span>系统归属阶段</span>
              <div class="readonly-value">{{ operationContext?.adjustmentStageCode ? formatStage(operationContext.adjustmentStageCode) : '--' }}</div>
            </div>
          </div>

          <div v-if="operationContext && !operationContext.canAdjust" class="context-note">
            {{ operationContext.blockedReason || '目标小组当前状态无法下发奖惩。' }}
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
                data-enter-confirm
                @keydown.enter="confirmInputOnEnter"
              />
            </label>
          </div>

          <label class="field">
            <span>原因说明</span>
            <textarea v-model="store.adjustmentForm.reason" rows="6" placeholder="例如：本季度市场竞标表现优异，奖励 5。"></textarea>
          </label>

          <div class="form-actions">
            <button type="button" class="btn primary" :disabled="previewingAdjustment || sendingAdjustment || !operationContext?.canAdjust" @click="handlePreviewAdjustment">
              {{ previewingAdjustment ? '计算中...' : '预览影响' }}
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
                <th>状态</th>
                <th>发送人</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="adjustments.length === 0">
                <td colspan="10" class="empty-row">当前还没有奖惩下发记录。</td>
              </tr>
              <tr v-for="item in adjustments" :key="item.id">
                <td>{{ formatDateTime(item.publishedAt) }}</td>
                <td>第{{ item.groupNo }}组</td>
                <td>{{ item.yearNo }}年</td>
                <td>{{ item.stageCode }}</td>
                <td>{{ item.adjustmentType === 'REWARD' ? '奖励' : '罚款' }}</td>
                <td class="number-cell">{{ formatAmount(item.amount) }}</td>
                <td class="content-cell">{{ item.reason }}</td>
                <td><span class="status-pill" :class="item.status.toLowerCase()">{{ formatAdjustmentStatus(item.status) }}</span></td>
                <td>{{ item.operatorName }}</td>
                <td>
                  <button v-if="item.canVoid" type="button" class="btn small danger" :disabled="previewingAdjustment" @click="handlePreviewVoid(item.id)">作废</button>
                  <span v-else class="muted">--</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <div v-if="impactPreview" class="modal-backdrop" @click.self="closeImpactPreview">
      <section class="impact-dialog" role="dialog" aria-modal="true" aria-labelledby="impact-dialog-title">
        <div class="panel-head dialog-head">
          <div>
            <strong id="impact-dialog-title">{{ previewOperation === 'CREATE' ? '确认下发奖惩' : '确认作废奖惩' }}</strong>
            <span>系统已按目标组最新保存数据重新计算</span>
          </div>
          <button type="button" class="btn small" @click="closeImpactPreview">关闭</button>
        </div>

        <div class="impact-body">
          <div v-if="impactPreview.willBankrupt" class="bankruptcy-warning">
            本次操作将使所得税后现金小于 0。确认后会立即生成只读破产快照并将该组永久标记为破产。
          </div>

          <dl class="impact-grid">
            <div><dt>系统归属阶段</dt><dd>{{ formatStage(impactPreview.resolvedStageCode) }}</dd></div>
            <div><dt>计算依据时间</dt><dd>{{ impactPreview.calculationBasisSavedAt ? formatDateTime(impactPreview.calculationBasisSavedAt) : '尚无草稿保存时间' }}</dd></div>
            <div><dt>税后现金（调整前）</dt><dd>{{ formatAmount(impactPreview.cashBefore) }}</dd></div>
            <div><dt>税后现金（调整后）</dt><dd :class="{ negative: impactPreview.cashAfter < 0 }">{{ formatAmount(impactPreview.cashAfter) }}</dd></div>
            <div><dt>税前利润（调整后）</dt><dd>{{ formatAmount(impactPreview.preTaxProfitAfter) }}</dd></div>
            <div><dt>所得税（调整后）</dt><dd>{{ formatAmount(impactPreview.incomeTaxAfter) }}</dd></div>
            <div><dt>净利润（调整后）</dt><dd>{{ formatAmount(impactPreview.netProfitAfter) }}</dd></div>
            <div><dt>所有者权益（调整后）</dt><dd>{{ formatAmount(impactPreview.totalEquityAfter) }}</dd></div>
          </dl>

          <label v-if="previewOperation === 'VOID'" class="field">
            <span>作废原因</span>
            <textarea v-model="voidReason" rows="3" placeholder="请填写本次作废原因"></textarea>
          </label>

          <div class="dialog-actions">
            <button type="button" class="btn" @click="closeImpactPreview">取消</button>
            <button
              type="button"
              class="btn primary"
              :class="{ danger: impactPreview.willBankrupt }"
              :disabled="sendingAdjustment || voidingAdjustment || (previewOperation === 'VOID' && !voidReason.trim())"
              @click="handleConfirmImpact"
            >
              {{ confirmButtonText }}
            </button>
          </div>
        </div>
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'

import { useAdminShellStore } from '@/stores/admin-shell'
import { useAdminNoticeStore } from '@/stores/admin-notice'
import { confirmInputOnEnter } from '@/utils/input-navigation'
import { hasFractionInput } from '@/utils/manual-integer'
import type { AdjustmentImpactResult, AdjustmentRecordStatus } from '@/types/sandbox-game-admin'

const shellStore = useAdminShellStore()
const store = useAdminNoticeStore()
const { groups, generalNotices, adjustments, loading, sendingGeneral, sendingAdjustment, previewingAdjustment, voidingAdjustment, pageMessage, operationContext } = storeToRefs(store)
const impactPreview = ref<AdjustmentImpactResult | null>(null)
const previewOperation = ref<'CREATE' | 'VOID'>('CREATE')
const previewAdjustmentId = ref<number | null>(null)
const voidReason = ref('')

const confirmButtonText = computed(() => {
  const action = previewOperation.value === 'CREATE' ? '下发' : '作废'
  return impactPreview.value?.willBankrupt ? `确认${action}并执行破产判定` : `确认${action}`
})

onMounted(async () => {
  try {
    if (!shellStore.config) {
      await shellStore.bootstrap()
    }
    await store.bootstrap()
  } catch {
    // 页面消息由 store 统一处理。
  }
})

async function handleRefresh() {
  try {
    await store.loadOperationContext()
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

async function handlePreviewAdjustment() {
  try {
    impactPreview.value = await store.previewAdjustmentNotice()
    previewOperation.value = 'CREATE'
    previewAdjustmentId.value = null
    voidReason.value = ''
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleAdjustmentGroupChange() {
  try {
    await store.loadOperationContext()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handlePreviewVoid(adjustmentId: number) {
  try {
    impactPreview.value = await store.previewVoidAdjustment(adjustmentId)
    previewOperation.value = 'VOID'
    previewAdjustmentId.value = adjustmentId
    voidReason.value = ''
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleConfirmImpact() {
  try {
    if (previewOperation.value === 'CREATE') {
      await store.sendAdjustmentNotice()
    } else if (previewAdjustmentId.value) {
      await store.voidAdjustmentNotice(previewAdjustmentId.value, voidReason.value.trim())
    }
    closeImpactPreview()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

function closeImpactPreview() {
  impactPreview.value = null
  previewAdjustmentId.value = null
  voidReason.value = ''
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

function formatAdjustmentStatus(value: AdjustmentRecordStatus) {
  if (value === 'EFFECTIVE') return '有效'
  if (value === 'VOIDED') return '已作废'
  return '快照失效'
}

function formatStage(value?: string | null) {
  const map: Record<string, string> = {
    Q1: 'Q1',
    Q2: 'Q2',
    Q3: 'Q3',
    Q4: 'Q4',
    YEAR_END: '年末',
  }
  return value ? map[value] ?? value : '--'
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

.readonly-value {
  min-height: 42px;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #f8fafc;
  padding: 10px 12px;
  color: var(--text);
  font-weight: 700;
}

.context-note {
  border: 1px solid #fdb022;
  border-radius: 12px;
  background: #fffaeb;
  padding: 10px 12px;
  color: #93370d;
  line-height: 1.5;
  font-size: 13px;
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

.btn.small {
  padding: 6px 10px;
  border-radius: 9px;
  font-size: 13px;
}

.btn.danger,
.btn.primary.danger {
  color: #ffffff;
  border-color: #b42318;
  background: #b42318;
}

.muted {
  color: var(--muted);
}

.status-pill {
  display: inline-flex;
  white-space: nowrap;
  border-radius: 999px;
  padding: 3px 8px;
  font-size: 12px;
  font-weight: 700;
}

.status-pill.effective {
  color: #067647;
  background: #ecfdf3;
}

.status-pill.voided,
.status-pill.snapshot_inactive {
  color: #667085;
  background: #f2f4f7;
}

.modal-backdrop {
  position: fixed;
  z-index: 1000;
  inset: 0;
  display: grid;
  place-items: center;
  padding: 24px;
  background: rgba(15, 23, 42, 0.52);
}

.impact-dialog {
  width: min(760px, 100%);
  max-height: calc(100vh - 48px);
  overflow: auto;
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 24px 64px rgba(15, 23, 42, 0.24);
}

.dialog-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.impact-body {
  display: grid;
  gap: 16px;
  padding: 18px;
}

.bankruptcy-warning {
  border: 1px solid #f04438;
  border-radius: 12px;
  padding: 12px 14px;
  color: #b42318;
  background: #fef3f2;
  line-height: 1.6;
}

.impact-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1px;
  margin: 0;
  border: 1px solid var(--line);
  background: var(--line);
}

.impact-grid > div {
  display: grid;
  gap: 6px;
  padding: 12px;
  background: #ffffff;
}

.impact-grid dt {
  color: var(--muted);
  font-size: 12px;
}

.impact-grid dd {
  margin: 0;
  font-weight: 700;
}

.impact-grid .negative {
  color: #b42318;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
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

  .impact-grid {
    grid-template-columns: 1fr;
  }
}
</style>
