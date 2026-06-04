<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>初始基线</h2>
        <p>页面语义是 1 份共享模板。提交后按组应用到全部小组，并整体锁定。</p>
      </div>
      <div class="hero-actions">
        <button type="button" class="btn" :disabled="loading || submitting" @click="handleRefresh">刷新</button>
        <button type="button" class="btn primary" :disabled="!view?.editable || submitting" @click="handleSubmit">
          {{ submitting ? '提交中...' : '提交并应用到全部小组' }}
        </button>
      </div>
    </header>

    <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
      {{ pageMessage.text }}
    </section>

    <div class="split-layout">
      <section class="sheet-card">
        <div class="sheet-head">
          <div>
            <strong>共享初始基线模板</strong>
            <span>首版使用单页高信息密度表格，方便和 Excel 对照录入。</span>
          </div>
          <span class="status-tag" :class="view?.submitted ? 'ok' : 'warn'">
            {{ view?.submitted ? '已提交并锁定' : '待提交' }}
          </span>
        </div>

        <div v-if="loading" class="empty-state">正在读取初始基线...</div>
        <div v-else class="table-scroll">
          <table class="baseline-table">
            <thead>
              <tr>
                <th>类别</th>
                <th>项目</th>
                <th>值</th>
                <th>说明</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="field in fieldDefs" :key="field.key">
                <td class="section-cell">{{ field.section }}</td>
                <td class="row-title">{{ field.label }}</td>
                <td class="input-cell" :class="{ invalid: hasFractionValue(draftPayload[field.key]) }">
                  <input
                    :value="draftPayload[field.key]"
                    type="number"
                    step="1"
                    inputmode="numeric"
                    :disabled="!view?.editable"
                    @input="handleFieldInput(field.key, $event)"
                  >
                </td>
                <td class="note-cell">{{ field.note }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <aside class="side-stack">
        <section class="panel-card">
          <div class="panel-head">
            <strong>提交状态</strong>
            <span>当前页只有一份共享模板，不做逐组编辑。</span>
          </div>
          <div class="meta-list">
            <div class="meta-item">
              <span>状态</span>
              <strong>{{ view?.submitted ? '已提交' : '未提交' }}</strong>
            </div>
            <div class="meta-item">
              <span>应用组数</span>
              <strong>{{ view?.appliedGroupCount ?? 0 }}</strong>
            </div>
            <div class="meta-item">
              <span>提交人</span>
              <strong>{{ view?.submitterName ?? '--' }}</strong>
            </div>
            <div class="meta-item">
              <span>提交时间</span>
              <strong>{{ formatDateTime(view?.submittedAt) }}</strong>
            </div>
          </div>
        </section>

        <section class="panel-card">
          <div class="panel-head">
            <strong>首版约束</strong>
            <span>当前按赛前锁定的沙盘版本展示字段名。</span>
          </div>
          <ul class="note-list">
            <li>提交前可编辑，提交后锁定只读。</li>
            <li>共享模板提交后由后端按组扇出应用。</li>
            <li>首版不支持提交后反复编辑或回滚。</li>
          </ul>
        </section>
      </aside>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { storeToRefs } from 'pinia'

import { applyDictionaryToBaselineLabels, resolveBaselineLabels } from '@/configs/sandbox-game-service-labels'
import { useAdminBaselineStore } from '@/stores/admin-baseline'
import { useAdminShellStore } from '@/stores/admin-shell'
import { useDictionaryStore } from '@/stores/dictionary'
import type { BaselinePayload } from '@/types/sandbox-game-admin'
import { hasFractionInput } from '@/utils/manual-integer'

const shellStore = useAdminShellStore()
const baselineStore = useAdminBaselineStore()
const dictionaryStore = useDictionaryStore()
const { view, draftPayload, loading, submitting, pageMessage } = storeToRefs(baselineStore)

const baselineLabels = computed(() =>
  applyDictionaryToBaselineLabels(
    resolveBaselineLabels(shellStore.config?.editionCode ?? shellStore.setupStatus?.editionCode),
    dictionaryStore.displayName,
  ),
)
const fieldDefs = computed<Array<{
  key: keyof BaselinePayload
  section: string
  label: string
  note: string
}>>(() => [
  { key: 'baselineSalesRevenue', section: '损益', label: baselineLabels.value.salesRevenue, note: '共享模板基线输入项' },
  { key: 'baselineDirectCost', section: '损益', label: baselineLabels.value.directCost, note: '共享模板基线输入项' },
  { key: 'baselineComprehensiveCost', section: '损益', label: baselineLabels.value.comprehensiveCost, note: '共享模板基线输入项' },
  { key: 'baselineDepreciation', section: '损益', label: baselineLabels.value.depreciation, note: '共享模板基线输入项' },
  { key: 'baselineFinanceIncomeExpense', section: '损益', label: baselineLabels.value.financeIncomeExpense, note: '共享模板基线输入项' },
  { key: 'baselineExtraIncomeExpense', section: '损益', label: baselineLabels.value.extraIncomeExpense, note: '共享模板基线输入项' },
  { key: 'baselineIncomeTax', section: '损益', label: baselineLabels.value.incomeTax, note: '共享模板基线输入项' },
  { key: 'baselineWorkInConstruction', section: '资产', label: baselineLabels.value.workInConstruction, note: '共享模板基线输入项' },
  { key: 'baselineFactoryAsset', section: '资产', label: baselineLabels.value.factoryAsset, note: '共享模板基线输入项' },
  { key: 'baselineLineResidual', section: '资产', label: baselineLabels.value.lineResidual, note: '共享模板基线输入项' },
  { key: 'baselineDepreciableAsset', section: '资产', label: baselineLabels.value.depreciableAsset, note: '共享模板基线输入项' },
  { key: 'baselineCash', section: '资产', label: baselineLabels.value.cash, note: '共享模板基线输入项' },
  { key: 'baselineReceivable', section: '资产', label: baselineLabels.value.receivable, note: '共享模板基线输入项' },
  { key: 'baselineWorkInProgress', section: '存货', label: baselineLabels.value.workInProgress, note: '共享模板基线输入项' },
  { key: 'baselineFinishedGoods', section: '存货', label: baselineLabels.value.finishedGoods, note: '共享模板基线输入项' },
  { key: 'baselineRawMaterials', section: '存货', label: baselineLabels.value.rawMaterials, note: '共享模板基线输入项' },
  { key: 'baselineShortTermLoan', section: '负债', label: baselineLabels.value.shortTermLoan, note: '共享模板基线输入项' },
  { key: 'baselineLongTermLoan', section: '负债', label: baselineLabels.value.longTermLoan, note: '共享模板基线输入项' },
  { key: 'baselineShareCapital', section: '权益', label: baselineLabels.value.shareCapital, note: '共享模板基线输入项' },
  { key: 'baselineRetainedEarnings', section: '权益', label: baselineLabels.value.retainedEarnings, note: '共享模板基线输入项' },
])

onMounted(async () => {
  try {
    if (!shellStore.config) {
      await shellStore.bootstrap()
    }
    await Promise.all([
      baselineStore.bootstrap(),
      dictionaryStore.loadCurrent(shellStore.config?.editionCode ?? shellStore.setupStatus?.editionCode, { silent: true }),
    ])
    dictionaryStore.startSilentSync()
  } catch {
    // 页面消息由 store 统一处理。
  }
})

onUnmounted(() => {
  dictionaryStore.stopSilentSync()
})

async function handleRefresh() {
  try {
    await Promise.all([shellStore.refreshConfig({ silent: true }), baselineStore.bootstrap()])
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleSubmit() {
  try {
    await baselineStore.submit()
    await shellStore.refreshConfig({ silent: true })
  } catch {
    // 页面消息由 store 统一处理。
  }
}

function handleFieldInput(key: keyof BaselinePayload, event: Event) {
  const target = event.target as HTMLInputElement
  const nextValue = Number(target.value)
  baselineStore.updateField(key, Number.isFinite(nextValue) ? nextValue : 0)
}

function hasFractionValue(value: number) {
  return hasFractionInput(value)
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

.split-layout {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) 340px;
  gap: 14px;
}

.sheet-card,
.panel-card {
  border: 1px solid var(--line);
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.05);
}

.sheet-head,
.panel-head {
  padding: 14px 16px;
  background: #f7f9fc;
  border-bottom: 1px solid var(--line);
}

.sheet-head strong,
.panel-head strong {
  display: block;
  margin-bottom: 4px;
  font-size: 16px;
}

.sheet-head span,
.panel-head span {
  color: var(--muted);
  font-size: 13px;
}

.status-tag {
  display: inline-flex;
  align-items: center;
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

.table-scroll {
  overflow: auto;
}

.baseline-table {
  width: 100%;
  min-width: 860px;
  border-collapse: collapse;
}

.baseline-table th,
.baseline-table td {
  border: 1px solid var(--line);
  padding: 10px 12px;
  font-size: 14px;
}

.baseline-table th {
  background: #f4f6f9;
  text-align: center;
}

.section-cell {
  background: #fbfbfc;
  font-weight: 700;
  white-space: nowrap;
}

.row-title {
  white-space: nowrap;
}

.input-cell input {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 10px;
  background: #f5fff0;
  padding: 8px 10px;
  text-align: right;
}

.input-cell input:disabled {
  background: var(--readonly-bg);
}

.input-cell.invalid input {
  border-color: var(--danger);
  background: #fff5f5;
  color: var(--danger);
}

.note-cell {
  color: var(--muted);
}

.side-stack {
  display: grid;
  gap: 14px;
}

.meta-list {
  display: grid;
  gap: 12px;
  padding: 16px;
}

.meta-item {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.meta-item span {
  color: var(--muted);
  font-size: 13px;
}

.meta-item strong {
  font-size: 14px;
}

.note-list {
  margin: 0;
  padding: 16px 16px 16px 32px;
  color: var(--muted);
  line-height: 1.8;
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

.empty-state {
  padding: 24px;
  color: var(--muted);
  text-align: center;
}

@media (max-width: 1180px) {
  .split-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 900px) {
  .hero {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>

