<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>汇总</h2>
        <p>同页展示年度汇总区与最终排名区。首版不做单独排名页面，也不做导出。</p>
      </div>
      <div class="hero-actions">
        <button type="button" class="btn" :disabled="loading" @click="handleRefresh">刷新汇总</button>
      </div>
    </header>

    <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
      {{ pageMessage.text }}
    </section>

    <div class="cards-grid">
      <article class="card">
        <span class="card-label">正式年份范围</span>
        <strong>{{ formalYears.length > 0 ? `1 - ${config?.finalYear}` : '尚未配置' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">当前开放年份</span>
        <strong>{{ config?.currentOpenYear ?? '--' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">共享基线状态</span>
        <strong>{{ config?.initialBaselineSubmitted ? '已提交' : '未提交' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">最近管理员动作</span>
        <strong>{{ latestAdminActionText }}</strong>
      </article>
    </div>

    <section class="sheet-card">
      <div class="sheet-head">
        <div>
          <strong>年度汇总区</strong>
          <span>按小组为行、按正式年份为列组。0 年不进入正式汇总。</span>
        </div>
      </div>

      <div v-if="loading" class="empty-state">正在读取年度汇总数据...</div>
      <div v-else-if="formalYears.length === 0" class="empty-state">当前最终年份仍为 0，正式年度汇总暂不可展示。</div>
      <div v-else class="table-scroll">
        <table class="summary-table">
          <thead>
            <tr>
              <th rowspan="2">小组</th>
              <th v-for="yearNo in formalYears" :key="`year-${yearNo}`" colspan="3" class="year-head">{{ yearNo }}年</th>
              <th rowspan="2">经营状态</th>
            </tr>
            <tr>
              <template v-for="yearNo in formalYears" :key="`year-sub-${yearNo}`">
                <th>收入</th>
                <th>利润</th>
                <th>权益</th>
              </template>
            </tr>
          </thead>
          <tbody>
            <tr v-if="summaryRows.length === 0">
              <td :colspan="formalYears.length * 3 + 2" class="empty-row">当前没有可展示的正式汇总数据。</td>
            </tr>
            <tr v-for="row in summaryRows" :key="row.groupId">
              <td class="row-title">第{{ row.groupNo }}组</td>
              <template v-for="yearNo in formalYears" :key="`${row.groupId}-${yearNo}`">
                <td class="number-cell">{{ formatMetric(row.byYear[yearNo]?.revenue) }}</td>
                <td class="number-cell">{{ formatMetric(row.byYear[yearNo]?.profit) }}</td>
                <td class="number-cell">{{ formatMetric(row.byYear[yearNo]?.equity) }}</td>
              </template>
              <td class="status-cell">
                <span class="status-tag" :class="row.businessStatus === 'BANKRUPT' ? 'danger' : 'ok'">
                  {{ formatBusinessStatus(row.businessStatus) }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section class="sheet-card ranking-card">
      <div class="sheet-head">
        <div>
          <strong>最终排名区</strong>
          <span>仅在最终年份满足开放条件后显示有效排名。</span>
        </div>
      </div>

      <div v-if="finalRanking?.list?.length" class="table-scroll narrow-scroll">
        <table class="ranking-table">
          <thead>
            <tr>
              <th>排名</th>
              <th>小组</th>
              <th>权益</th>
              <th>收入</th>
              <th>利润</th>
              <th>经营状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in finalRanking.list" :key="item.groupId">
              <td class="center-cell">{{ item.ranking }}</td>
              <td>第{{ item.groupNo }}组</td>
              <td class="number-cell">{{ formatMetric(item.equity) }}</td>
              <td class="number-cell">{{ formatMetric(item.revenue) }}</td>
              <td class="number-cell">{{ formatMetric(item.profit) }}</td>
              <td class="status-cell">
                <span class="status-tag" :class="item.businessStatus === 'BANKRUPT' ? 'danger' : 'ok'">
                  {{ formatBusinessStatus(item.businessStatus) }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else class="placeholder">
        {{ rankingUnavailableReason || '最终排名暂未开放，请等待所有未破产小组完成最终年度财报后查看。' }}
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'

import { useAdminShellStore } from '@/stores/admin-shell'
import { useAdminSummaryStore } from '@/stores/admin-summary'
import type { AdminYearSummaryItem, AdminYearSummaryResult } from '@/types/sandbox-game-admin'
import { formatBusinessStatus } from '@/utils/sandbox-game-display'

interface SummaryRow {
  groupId: number
  groupNo: number
  groupName: string
  businessStatus: string
  byYear: Partial<Record<number, AdminYearSummaryItem>>
}

const shellStore = useAdminShellStore()
const summaryStore = useAdminSummaryStore()
const { config } = storeToRefs(shellStore)
const { yearSummaries, finalRanking, rankingUnavailableReason, loading, pageMessage } = storeToRefs(summaryStore)

const initialized = ref(false)

const formalYears = computed(() => {
  const finalYear = Math.max(config.value?.finalYear ?? 0, 0)
  return Array.from({ length: finalYear }, (_, index) => index + 1)
})

const summaryRows = computed(() => buildSummaryRows(yearSummaries.value))

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
    await summaryStore.bootstrap(config.value?.finalYear ?? 0)
  } catch {
    // 页面消息由 store 统一处理。
  } finally {
    initialized.value = true
  }
})

watch(
  () => config.value?.finalYear,
  async (next, previous) => {
    if (!initialized.value || next === undefined || next === previous) {
      return
    }
    try {
      await summaryStore.bootstrap(next)
    } catch {
      // 页面消息由 store 统一处理。
    }
  },
)

async function handleRefresh() {
  try {
    await shellStore.refreshConfig({ silent: true })
    await summaryStore.bootstrap(config.value?.finalYear ?? 0)
  } catch {
    // 页面消息由 store 统一处理。
  }
}

function buildSummaryRows(results: AdminYearSummaryResult[]) {
  const rowMap = new Map<number, SummaryRow>()
  for (const result of results) {
    for (const item of result.list) {
      const existing = rowMap.get(item.groupId) ?? {
        groupId: item.groupId,
        groupNo: item.groupNo,
        groupName: item.groupName,
        businessStatus: item.businessStatus,
        byYear: {},
      }
      existing.businessStatus = item.businessStatus
      existing.byYear[result.yearNo] = item
      rowMap.set(item.groupId, existing)
    }
  }
  return Array.from(rowMap.values()).sort((left, right) => left.groupNo - right.groupNo)
}

function formatMetric(value?: number) {
  if (typeof value !== 'number') {
    return '--'
  }
  return value.toLocaleString('zh-CN', {
    minimumFractionDigits: Number.isInteger(value) ? 0 : 2,
    maximumFractionDigits: 2,
  })
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

.cards-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.card,
.sheet-card {
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

.sheet-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  padding: 14px 16px;
  background: #f7f9fc;
  border-bottom: 1px solid var(--line);
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

.table-scroll {
  overflow: auto;
}

.summary-table,
.ranking-table {
  width: 100%;
  min-width: 1080px;
  border-collapse: collapse;
}

.ranking-table {
  min-width: 720px;
}

.summary-table th,
.summary-table td,
.ranking-table th,
.ranking-table td {
  border: 1px solid var(--line);
  padding: 10px 12px;
  font-size: 14px;
}

.summary-table th,
.ranking-table th {
  background: #f4f6f9;
  font-weight: 700;
  text-align: center;
}

.year-head {
  background: #f6decb !important;
}

.row-title {
  background: #fbfbfc;
  font-weight: 700;
  white-space: nowrap;
}

.number-cell {
  text-align: right;
}

.center-cell,
.status-cell {
  text-align: center;
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

.status-tag.danger {
  background: #fee4e2;
  color: var(--danger);
}

.empty-state,
.placeholder,
.empty-row {
  padding: 24px;
  color: var(--muted);
  text-align: center;
}

.ranking-card {
  overflow: hidden;
}

.narrow-scroll {
  border-top: 1px solid var(--line);
}

@media (max-width: 1240px) {
  .cards-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .cards-grid {
    grid-template-columns: 1fr;
  }
}
</style>
