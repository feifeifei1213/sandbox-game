<template>
  <aside class="sidebar">
    <PlayerNoticePanel :notice-board="view?.noticeBoard" />

    <section class="panel">
      <h3>当前状态</h3>
      <dl class="meta-list">
        <div>
          <dt>年份</dt>
          <dd>{{ yearLabel }}</dd>
        </div>
        <div>
          <dt>年度状态</dt>
          <dd>{{ formatYearStatus(view?.yearStatus) }}</dd>
        </div>
        <div>
          <dt>财报状态</dt>
          <dd>{{ formatReportStatus(view?.reportStatus) }}</dd>
        </div>
        <div>
          <dt>经营状态</dt>
          <dd :class="{ danger: view?.businessStatus === 'BANKRUPT' }">{{ formatBusinessStatus(view?.businessStatus) }}</dd>
        </div>
        <div>
          <dt>最近草稿</dt>
          <dd>{{ lastDraftSavedAt }}</dd>
        </div>
      </dl>
    </section>

    <section v-if="view?.hasInvalidDraft" class="panel warning-panel">
      <h3>失效草稿</h3>
      <p>当前保留的是上次已失效的财报草稿，需要玩家重新核对并再次提交。</p>
      <p class="hint warning-hint">经营页重新提交后，请回到本页确认绿色手工项，再完成财报提交。</p>
    </section>

    <section class="panel">
      <h3>平衡校验</h3>
      <div class="balance-card" :class="{ pass: balancePassed, fail: !balancePassed }">
        <strong>{{ balancePassed ? '校验通过' : '校验未通过' }}</strong>
        <span>总资产 {{ formatNumber(computedPayload.reportTotalAssets) }}</span>
        <span>总负债和权益 {{ formatNumber(computedPayload.reportTotalLiabilityEquity) }}</span>
        <span>差额 {{ formatNumber(balanceGap) }}</span>
      </div>
    </section>

    <section class="panel">
      <h3>税率说明</h3>
      <ul class="text-list">
        <li v-for="item in taxRateOptions" :key="item">{{ formatTaxRate(item) }}</li>
      </ul>
      <p class="hint">首版财报页中，所得税税率只能通过下拉框选择，不允许自由输入。</p>
    </section>

    <section class="panel">
      <h3>提交控制</h3>
      <div class="action-list">
        <button type="button" class="btn" :disabled="saving || submitting || !view?.canEdit || !dirty" @click="$emit('save')">
          {{ saving ? '保存中...' : '保存草稿' }}
        </button>
        <button type="button" class="btn primary" :disabled="saving || submitting || !submitReady" @click="$emit('submit')">
          {{ submitting ? '提交中...' : '提交财报' }}
        </button>
      </div>
      <p class="hint">绿色单元格为手工项，黄色单元格为系统计算结果。</p>
      <p v-if="view?.hasInvalidDraft" class="warning">当前为失效草稿状态，必须重新提交后才会恢复正式结果。</p>
      <p v-if="missingFields.length" class="warning">待填写：{{ missingFields.join('、') }}</p>
      <p v-else-if="!balancePassed" class="warning">当前资产负债尚未平衡，不能提交。</p>
      <p v-else-if="!view?.canEdit" class="hint">当前年份财报未开放或已提交，页面为只读状态。</p>
    </section>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import PlayerNoticePanel from '@/components/sandbox-game/player/PlayerNoticePanel.vue'
import { serviceReportLabels, type SandboxGameReportLabels } from '@/configs/sandbox-game-service-labels'
import type { PlayerReportView, ReportComputedPayload } from '@/types/sandbox-game'
import {
  formatBusinessStatus,
  formatReportStatus,
  formatYearStatus,
} from '@/utils/sandbox-game-display'

const props = defineProps<{
  view: PlayerReportView | null
  yearLabel: string
  computedPayload: ReportComputedPayload
  balanceGap: number
  balancePassed: boolean
  taxRateOptions: number[]
  missingFields: string[]
  dirty: boolean
  saving: boolean
  submitting: boolean
  submitReady: boolean
  labels?: SandboxGameReportLabels
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
const labels = computed(() => props.labels ?? serviceReportLabels)

function formatTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString('zh-CN', { hour12: false })
}

function formatNumber(value: number | null | undefined) {
  if (typeof value !== 'number' || Number.isNaN(value)) {
    return '--'
  }
  return value.toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

function formatTaxRate(value: number) {
  if (value === 0.25) {
    return '0.25 一般企业'
  }
  if (value === 0.15) {
    return labels.value.taxRateEnterprise15
  }
  return '0 全额弥补亏损'
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

.warning-panel {
  border-color: #efd4aa;
  background: #fffaf2;
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

.balance-card {
  display: grid;
  gap: 6px;
  border-radius: 14px;
  padding: 14px;
  border: 1px solid #d8e1ea;
  background: #f8fbfe;
}

.balance-card.pass {
  border-color: #b9dec9;
  background: #eefaf2;
  color: #1f6b40;
}

.balance-card.fail {
  border-color: #efc5c5;
  background: #fff6f6;
  color: #b24040;
}

.text-list {
  margin: 0;
  padding-left: 18px;
  display: grid;
  gap: 6px;
  color: #334155;
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

.warning {
  margin: 8px 0 0;
  color: #b24040;
  font-size: 12px;
  line-height: 1.6;
}

.warning-hint {
  color: #8a5a17;
}
</style>
