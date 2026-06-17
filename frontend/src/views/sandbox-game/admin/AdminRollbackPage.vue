<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>回退与修正</h2>
        </div>
      <div class="hero-actions">
        <button type="button" class="btn" :disabled="loading" @click="handleRefresh">刷新</button>
      </div>
    </header>

    <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
      {{ pageMessage.text }}
    </section>

    <nav class="tabs">
      <button type="button" :class="{ active: activeTab === 'unlock' }" @click="activeTab = 'unlock'">退回重提</button>
      <button type="button" :class="{ active: activeTab === 'snapshot' }" @click="activeTab = 'snapshot'">恢复快照</button>
    </nav>

    <section v-if="activeTab === 'unlock'" class="sheet-card">
      <div class="sheet-head">
        <div>
          <strong>退回重提</strong>
          
        </div>
      </div>
      <div class="form-grid">
        <label class="field">
          <span>目标小组</span>
          <select v-model.number="store.unlockForm.groupId">
            <option v-for="group in groups" :key="group.groupId" :value="group.groupId">第{{ group.groupNo }}组 · {{ group.groupName }}</option>
          </select>
        </label>
        <label class="field">
          <span>年份</span>
          <input v-model.number="store.unlockForm.yearNo" type="number" min="0" step="1">
        </label>
        <label class="field">
          <span>目标类型</span>
          <select v-model="store.unlockForm.unlockTargetType" @change="handleUnlockTargetChange">
            <option value="OPERATING">经营页</option>
            <option value="REPORT">财报页</option>
          </select>
        </label>
        <label v-if="store.unlockForm.unlockTargetType === 'OPERATING'" class="field">
          <span>经营阶段</span>
          <select v-model="store.unlockForm.targetStageCode">
            <option value="Q1">Q1</option>
            <option value="Q2">Q2</option>
            <option value="Q3">Q3</option>
            <option value="Q4">Q4</option>
            <option value="YEAR_END">年末</option>
          </select>
        </label>
      </div>
      <div class="action-row">
        <button type="button" class="primary-btn" :disabled="operating" @click="handleUnlockRetry">提交退回重提</button>
      </div>
    </section>

    <section v-else class="snapshot-layout">
      <article class="sheet-card">
        <div class="sheet-head">
          <div>
            <strong>快照筛选</strong>
            <span>首版只支持恢复单组快照；全局快照仅用于审计。</span>
          </div>
        </div>
        <div class="filter-grid">
          <label class="field">
            <span>范围</span>
            <select v-model="store.filters.snapshotScope">
              <option value="">全部</option>
              <option value="GROUP">单组</option>
              <option value="GLOBAL">全局</option>
            </select>
          </label>
          <label class="field">
            <span>类型</span>
            <select v-model="store.filters.snapshotType">
              <option value="">全部</option>
              <option value="AUTO">自动</option>
              <option value="MANUAL">手动</option>
              <option value="SAFETY">安全</option>
            </select>
          </label>
          <label class="field">
            <span>小组</span>
            <select :value="store.filters.groupId ?? ''" @change="handleFilterGroupChange">
              <option value="">全部</option>
              <option v-for="group in groups" :key="group.groupId" :value="group.groupId">第{{ group.groupNo }}组</option>
            </select>
          </label>
          <label class="field">
            <span>年份</span>
            <input v-model.number="store.filters.yearNo" type="number" min="0" step="1">
          </label>
          <label class="field">
            <span>阶段</span>
            <select v-model="store.filters.stageCode">
              <option value="">全部</option>
              <option value="Q1">Q1</option>
              <option value="Q2">Q2</option>
              <option value="Q3">Q3</option>
              <option value="Q4">Q4</option>
              <option value="YEAR_END">年末</option>
              <option value="REPORT">财报</option>
            </select>
          </label>
          <div class="field action-field">
            <span>&nbsp;</span>
            <button type="button" class="btn" :disabled="loading" @click="handleLoadSnapshots">查询快照</button>
          </div>
        </div>

        <div class="table-scroll">
          <table class="snapshot-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>范围</th>
                <th>类型</th>
                <th>节点</th>
                <th>说明</th>
                <th>小组</th>
                <th>年份</th>
                <th>阶段</th>
                <th>创建人</th>
                <th>时间</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="snapshots.length === 0">
                <td colspan="10" class="empty-row">暂无快照。</td>
              </tr>
              <tr
                v-for="snapshot in snapshots"
                :key="snapshot.id"
                :class="{ selected: selectedSnapshotId === snapshot.id }"
                @click="handleSelectSnapshot(snapshot.id)"
              >
                <td>#{{ snapshot.id }}</td>
                <td>{{ formatSnapshotScope(snapshot.snapshotScope) }}</td>
                <td>{{ formatSnapshotType(snapshot.snapshotType) }}</td>
                <td>{{ formatTrigger(snapshot.triggerCode) }}</td>
                <td>{{ snapshot.description || '--' }}</td>
                <td>{{ snapshot.groupName ?? '--' }}</td>
                <td>{{ snapshot.yearNo ?? '--' }}</td>
                <td>{{ formatStage(snapshot.stageCode) }}</td>
                <td>{{ snapshot.createdByName }}</td>
                <td>{{ formatDateTime(snapshot.createdAt) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </article>

      <aside class="side-stack">
        <article class="sheet-card">
          <div class="sheet-head compact">
            <div>
              <strong>手动快照</strong>
              <span>不改变业务状态。</span>
            </div>
          </div>
          <div class="small-form">
            <label class="field">
              <span>范围</span>
              <select v-model="store.manualSnapshotForm.snapshotScope">
                <option value="GROUP">单组</option>
                <option value="GLOBAL">全局</option>
              </select>
            </label>
            <label v-if="store.manualSnapshotForm.snapshotScope === 'GROUP'" class="field">
              <span>小组</span>
              <select v-model.number="store.manualSnapshotForm.groupId">
                <option v-for="group in groups" :key="group.groupId" :value="group.groupId">第{{ group.groupNo }}组</option>
              </select>
            </label>
            <label class="field">
              <span>年份</span>
              <input v-model.number="store.manualSnapshotForm.yearNo" type="number" min="0" step="1">
            </label>
            <label class="field">
              <span>阶段</span>
              <select v-model="store.manualSnapshotForm.stageCode">
                <option value="">按当前状态</option>
                <option value="Q1">Q1</option>
                <option value="Q2">Q2</option>
                <option value="Q3">Q3</option>
                <option value="Q4">Q4</option>
                <option value="YEAR_END">年末</option>
                <option value="REPORT">财报</option>
              </select>
            </label>
            <label class="field">
              <span>说明</span>
              <textarea v-model="store.manualSnapshotForm.description" rows="3"></textarea>
            </label>
            <button type="button" class="btn full" :disabled="operating" @click="handleCreateSnapshot">创建手动快照</button>
          </div>
        </article>

        <article class="sheet-card">
          <div class="sheet-head compact">
            <div>
              <strong>快照详情</strong>
              <span>{{ selectedSnapshot ? `#${selectedSnapshot.id}` : '请选择快照' }}</span>
            </div>
          </div>
          <div v-if="snapshotDetail" class="detail-panel">
            <dl>
              <div v-for="item in detailRows" :key="item.label">
                <dt>{{ item.label }}</dt>
                <dd>{{ item.value }}</dd>
              </div>
            </dl>
            <label class="field">
              <span>恢复原因</span>
              <textarea v-model="store.restoreForm.reason" rows="3"></textarea>
            </label>
            <label class="field">
              <span>确认文本</span>
              <input v-model="store.restoreForm.confirmText" placeholder="确认恢复">
            </label>
            <button
              type="button"
              class="danger-btn full"
              :disabled="operating || snapshotDetail.snapshot.snapshotScope !== 'GROUP'"
              @click="handleRestoreSnapshot"
            >
              恢复单组快照
            </button>
          </div>
          <div v-else class="empty-row">点击左侧快照查看详情。</div>
        </article>
      </aside>
    </section>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'

import { useAdminRollbackStore } from '@/stores/admin-rollback'
import { useAdminShellStore } from '@/stores/admin-shell'

type TabKey = 'unlock' | 'snapshot'

const shellStore = useAdminShellStore()
const store = useAdminRollbackStore()
const activeTab = ref<TabKey>('unlock')
const {
  groups,
  snapshots,
  loading,
  operating,
  pageMessage,
  selectedSnapshotId,
  selectedSnapshot,
  snapshotDetail,
} = storeToRefs(store)

const detailRows = computed(() => {
  if (!snapshotDetail.value) {
    return []
  }
  const summary = snapshotDetail.value.stateSummary
  const preview = snapshotDetail.value.payloadPreview
  return [
    { label: '生成节点', value: formatTrigger(snapshotDetail.value.snapshot.triggerCode) },
    { label: '快照说明', value: snapshotDetail.value.snapshot.description || '--' },
    { label: '目标小组', value: String(summary.groupName ?? snapshotDetail.value.snapshot.groupName ?? '--') },
    { label: '年份', value: String(summary.yearNo ?? snapshotDetail.value.snapshot.yearNo ?? '--') },
    { label: '阶段', value: formatStage(String(summary.targetStageCode ?? snapshotDetail.value.snapshot.stageCode ?? '')) },
    { label: '年度状态', value: String(summary.yearStatus ?? '--') },
    { label: '财报状态', value: String(summary.reportStatus ?? '--') },
    { label: '业务状态', value: String(summary.businessStatus ?? '--') },
    { label: '载荷版本', value: snapshotDetail.value.payloadVersion },
    { label: '载荷大小', value: `${snapshotDetail.value.payloadSize} 字节` },
    { label: '经营草稿数', value: String(preview.operatingDraftCount ?? '--') },
    { label: '订单记录数', value: String(preview.orderSelectionCount ?? preview.orderPoolCount ?? '--') },
  ]
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
    await shellStore.refreshConfig({ silent: true })
    await store.loadSnapshots()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleLoadSnapshots() {
  await store.loadSnapshots()
}

async function handleSelectSnapshot(snapshotId: number) {
  await store.selectSnapshot(snapshotId)
}

async function handleUnlockRetry() {
  await store.submitUnlockRetry()
  await shellStore.refreshConfig({ silent: true })
}

async function handleCreateSnapshot() {
  await store.submitManualSnapshot()
}

async function handleRestoreSnapshot() {
  await store.submitRestoreSnapshot()
  await shellStore.refreshConfig({ silent: true })
}

function handleUnlockTargetChange() {
  if (store.unlockForm.unlockTargetType === 'REPORT') {
    store.unlockForm.targetStageCode = null
  } else if (!store.unlockForm.targetStageCode) {
    store.unlockForm.targetStageCode = 'YEAR_END'
  }
}

function handleFilterGroupChange(event: Event) {
  const value = (event.target as HTMLSelectElement).value
  store.setFilterGroup(value ? Number(value) : null)
}

function formatSnapshotScope(value?: string | null) {
  if (value === 'GROUP') {
    return '单组'
  }
  if (value === 'GLOBAL') {
    return '全局'
  }
  return value || '--'
}

function formatSnapshotType(value?: string | null) {
  if (value === 'AUTO') {
    return '自动'
  }
  if (value === 'MANUAL') {
    return '手动'
  }
  if (value === 'SAFETY') {
    return '安全'
  }
  return value || '--'
}

function formatTrigger(value?: string | null) {
  const map: Record<string, string> = {
    STAGE_SUBMITTED: '经营提交',
    REPORT_SUBMITTED: '财报提交',
    ROLLBACK_STAGE_RESUBMITTED: '回退补提经营',
    ROLLBACK_REPORT_RESUBMITTED: '回退补提财报',
    ORDER_POOL_CONFIRMED: '订单池确认',
    SEGMENT_COMPLETED: '标段完成',
    OPEN_NEXT_YEAR: '开放下一年',
    BEFORE_ROLLBACK: '回退前',
    MANUAL: '手动',
  }
  return value ? map[value] ?? value : '--'
}

function formatStage(value?: string | null) {
  const map: Record<string, string> = {
    Q1: 'Q1',
    Q2: 'Q2',
    Q3: 'Q3',
    Q4: 'Q4',
    YEAR_END: '年末',
    REPORT: '财报',
  }
  return value ? map[value] ?? value : '--'
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

.btn,
.primary-btn,
.danger-btn {
  border: 1px solid var(--line);
  background: #ffffff;
  border-radius: 12px;
  padding: 10px 14px;
}

.primary-btn {
  background: var(--accent);
  border-color: var(--accent);
  color: #ffffff;
}

.danger-btn {
  background: #b42318;
  border-color: #b42318;
  color: #ffffff;
}

.full {
  width: 100%;
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

.tabs {
  display: inline-flex;
  gap: 8px;
  padding: 6px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: #ffffff;
  width: fit-content;
}

.tabs button {
  border: 0;
  border-radius: 10px;
  background: transparent;
  padding: 9px 14px;
}

.tabs button.active {
  background: #eef4ff;
  color: var(--accent);
  font-weight: 700;
}

.sheet-card {
  border: 1px solid var(--line);
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.05);
  overflow: hidden;
}

.sheet-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  padding: 14px 16px;
  background: #f7f9fc;
  border-bottom: 1px solid var(--line);
}

.sheet-head.compact {
  padding: 12px 14px;
}

.sheet-head strong {
  display: block;
  margin-bottom: 4px;
  font-size: 16px;
}

.sheet-head span {
  color: var(--muted);
  font-size: 13px;
}

.form-grid,
.filter-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
  padding: 16px;
}

.filter-grid {
  grid-template-columns: repeat(6, minmax(0, 1fr));
}

.small-form {
  display: grid;
  gap: 12px;
  padding: 14px;
}

.field {
  display: grid;
  gap: 6px;
}

.field.wide {
  grid-column: span 4;
}

.field span {
  color: var(--muted);
  font-size: 13px;
}

.field input,
.field select,
.field textarea {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 10px;
  background: #ffffff;
  padding: 9px 10px;
  font: inherit;
}

.field textarea {
  resize: vertical;
}

.action-field {
  align-self: end;
}

.action-row {
  padding: 0 16px 16px;
}

.snapshot-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 16px;
  align-items: start;
}

.side-stack {
  display: grid;
  gap: 16px;
}

.table-scroll {
  overflow: auto;
}

.snapshot-table {
  width: 100%;
  min-width: 1120px;
  border-collapse: collapse;
}

.snapshot-table th,
.snapshot-table td {
  border: 1px solid var(--line);
  padding: 10px 12px;
  font-size: 13px;
}

.snapshot-table th {
  background: #f4f6f9;
  text-align: left;
}

.snapshot-table tr {
  cursor: pointer;
}

.snapshot-table tr.selected td {
  background: #eef4ff;
}

.empty-row {
  padding: 20px;
  color: var(--muted);
  text-align: center;
}

.detail-panel {
  display: grid;
  gap: 14px;
  padding: 14px;
}

.detail-panel dl {
  display: grid;
  gap: 8px;
  margin: 0;
}

.detail-panel dl div {
  display: grid;
  grid-template-columns: 88px minmax(0, 1fr);
  gap: 8px;
}

.detail-panel dt {
  color: var(--muted);
}

.detail-panel dd {
  margin: 0;
  word-break: break-word;
}

@media (max-width: 1280px) {
  .snapshot-layout {
    grid-template-columns: 1fr;
  }

  .filter-grid,
  .form-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .field.wide {
    grid-column: span 2;
  }
}

@media (max-width: 820px) {
  .hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .filter-grid,
  .form-grid {
    grid-template-columns: 1fr;
  }

  .field.wide {
    grid-column: span 1;
  }
}
</style>
