<template>
  <div class="sheet-shell">
    <div class="sheet-scroll">
      <table class="report-sheet">
        <thead>
          <tr class="title-row">
            <th class="corner-head">#</th>
            <th class="section-head" colspan="2">损益表</th>
            <th class="spacer-head"></th>
            <th class="section-head" colspan="2">资产</th>
            <th class="spacer-head"></th>
            <th class="section-head" colspan="2">负债与权益</th>
          </tr>
          <tr class="sub-head-row">
            <th class="row-head muted">行</th>
            <th class="sub-head">项目</th>
            <th class="sub-head">金额</th>
            <th class="spacer-head"></th>
            <th class="sub-head">项目</th>
            <th class="sub-head">金额</th>
            <th class="spacer-head"></th>
            <th class="sub-head">项目</th>
            <th class="sub-head">金额</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in sheetRows" :key="row.rowNo">
            <th class="row-head">{{ row.rowNo }}</th>
            <td class="label-cell">{{ row.profit.label || '' }}</td>
            <td :class="valueCellClass(row.profit)">
              <template v-if="row.profit.kind === 'manual-select'">
                <select :value="manualSelectValue('incomeTaxRate')" :disabled="!canEdit" @change="updateTaxRate($event)">
                  <option value="">请选择</option>
                  <option v-for="item in taxRateOptions" :key="item" :value="String(item)">{{ formatTaxRate(item) }}</option>
                </select>
              </template>
              <template v-else-if="row.profit.computedKey">{{ formatNumber(computedPayload[row.profit.computedKey]) }}</template>
            </td>
            <td class="spacer-cell"></td>
            <td class="label-cell">{{ row.asset.label || '' }}</td>
            <td :class="valueCellClass(row.asset)">
              <template v-if="row.asset.kind === 'manual-number'">
                <input
                  :value="manualInputValue(row.asset.manualKey)"
                  :disabled="!canEdit"
                  inputmode="decimal"
                  @input="updateManualNumber(row.asset.manualKey, $event)"
                />
              </template>
              <template v-else-if="row.asset.computedKey">{{ formatNumber(computedPayload[row.asset.computedKey]) }}</template>
            </td>
            <td class="spacer-cell"></td>
            <td class="label-cell">{{ row.liability.label || '' }}</td>
            <td :class="valueCellClass(row.liability)">
              <template v-if="row.liability.computedKey">{{ formatNumber(computedPayload[row.liability.computedKey]) }}</template>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ReportComputedPayload, ReportManualPayload } from '@/types/sandbox-game'
import { cloneReportManualPayload } from '@/types/sandbox-game'

type ManualKey = keyof ReportManualPayload
type ComputedKey = keyof ReportComputedPayload
type CellKind = 'blank' | 'computed' | 'manual-number' | 'manual-select'

interface SheetCell {
  label?: string
  kind: CellKind
  computedKey?: ComputedKey
  manualKey?: ManualKey
}

interface SheetRow {
  rowNo: number
  profit: SheetCell
  asset: SheetCell
  liability: SheetCell
}

const props = defineProps<{
  modelValue: ReportManualPayload
  computedPayload: ReportComputedPayload
  canEdit: boolean
  taxRateOptions: number[]
}>()

const emit = defineEmits<{
  (event: 'update:modelValue', value: ReportManualPayload): void
}>()

const blankCell: SheetCell = { kind: 'blank' }

const sheetRows: SheetRow[] = [
  { rowNo: 4, profit: { label: '销售收入', kind: 'computed', computedKey: 'reportSalesRevenue' }, asset: blankCell, liability: blankCell },
  { rowNo: 5, profit: { label: '直接成本', kind: 'computed', computedKey: 'reportDirectCost' }, asset: blankCell, liability: blankCell },
  { rowNo: 6, profit: blankCell, asset: { label: '厂房', kind: 'computed', computedKey: 'reportFactoryAsset' }, liability: blankCell },
  { rowNo: 7, profit: { label: '毛利', kind: 'computed', computedKey: 'reportGrossProfit' }, asset: { label: '生产线残值', kind: 'computed', computedKey: 'reportLineResidual' }, liability: { label: '短期负债', kind: 'computed', computedKey: 'reportShortTermLiability' } },
  { rowNo: 8, profit: blankCell, asset: { label: '待折资产', kind: 'computed', computedKey: 'reportDepreciableAsset' }, liability: { label: '长期负债', kind: 'computed', computedKey: 'reportLongTermLiability' } },
  { rowNo: 9, profit: blankCell, asset: { label: '非流动资产合计', kind: 'computed', computedKey: 'reportTotalNonCurrentAssets' }, liability: { label: '负债合计', kind: 'computed', computedKey: 'reportTotalLiability' } },
  { rowNo: 10, profit: { label: '综合费用', kind: 'computed', computedKey: 'reportComprehensiveCost' }, asset: blankCell, liability: blankCell },
  { rowNo: 11, profit: { label: '折旧', kind: 'computed', computedKey: 'reportDepreciation' }, asset: blankCell, liability: blankCell },
  { rowNo: 12, profit: blankCell, asset: blankCell, liability: blankCell },
  { rowNo: 13, profit: { label: '营业利润', kind: 'computed', computedKey: 'reportOperatingProfit' }, asset: { label: '现金', kind: 'computed', computedKey: 'reportCash' }, liability: blankCell },
  { rowNo: 14, profit: blankCell, asset: { label: '应收款', kind: 'computed', computedKey: 'reportReceivable' }, liability: blankCell },
  { rowNo: 15, profit: blankCell, asset: { label: '在制品', kind: 'manual-number', manualKey: 'workInProgress' }, liability: { label: '股东资本', kind: 'computed', computedKey: 'reportShareCapital' } },
  { rowNo: 16, profit: { label: '财务收入/支出', kind: 'computed', computedKey: 'reportFinanceIncomeExpense' }, asset: { label: '成品', kind: 'manual-number', manualKey: 'finishedGoods' }, liability: { label: '利润留存', kind: 'computed', computedKey: 'reportRetainedEarnings' } },
  { rowNo: 17, profit: { label: '额外收入/支出', kind: 'computed', computedKey: 'reportExtraIncomeExpense' }, asset: { label: '材料', kind: 'manual-number', manualKey: 'rawMaterials' }, liability: blankCell },
  { rowNo: 18, profit: blankCell, asset: { label: '税后现金', kind: 'computed', computedKey: 'reportPostTaxCash' }, liability: blankCell },
  { rowNo: 19, profit: { label: '税前利润', kind: 'computed', computedKey: 'reportPreTaxProfit' }, asset: { label: '流动资产合计', kind: 'computed', computedKey: 'reportTotalCurrentAssets' }, liability: { label: '权益合计', kind: 'computed', computedKey: 'reportTotalEquity' } },
  { rowNo: 20, profit: blankCell, asset: blankCell, liability: blankCell },
  { rowNo: 21, profit: { label: '所得税', kind: 'computed', computedKey: 'reportIncomeTax' }, asset: blankCell, liability: blankCell },
  { rowNo: 22, profit: { label: '年度净利润', kind: 'computed', computedKey: 'reportNetProfit' }, asset: { label: '总资产', kind: 'computed', computedKey: 'reportTotalAssets' }, liability: { label: '总负债和权益', kind: 'computed', computedKey: 'reportTotalLiabilityEquity' } },
  { rowNo: 23, profit: blankCell, asset: blankCell, liability: blankCell },
  { rowNo: 24, profit: { label: '所得税税率', kind: 'manual-select', manualKey: 'incomeTaxRate' }, asset: blankCell, liability: blankCell },
]

function valueCellClass(cell: SheetCell) {
  return {
    'value-cell': true,
    'auto-cell': cell.kind === 'computed',
    'manual-cell': cell.kind === 'manual-number' || cell.kind === 'manual-select',
    'blank-cell': cell.kind === 'blank',
  }
}

function formatNumber(value: number | null | undefined) {
  if (typeof value !== 'number' || Number.isNaN(value)) {
    return ''
  }
  return value.toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

function manualInputValue(key?: ManualKey) {
  if (!key) {
    return ''
  }
  const value = props.modelValue[key]
  return value === null ? '' : String(value)
}

function manualSelectValue(key: ManualKey) {
  const value = props.modelValue[key]
  return value === null ? '' : String(value)
}

function updateManualNumber(key: ManualKey | undefined, event: Event) {
  if (!key) {
    return
  }
  const next = cloneReportManualPayload(props.modelValue)
  const raw = (event.target as HTMLInputElement).value.trim()
  next[key] = raw === '' ? null : Number(raw)
  emit('update:modelValue', next)
}

function updateTaxRate(event: Event) {
  const next = cloneReportManualPayload(props.modelValue)
  const raw = (event.target as HTMLSelectElement).value
  next.incomeTaxRate = raw === '' ? null : Number(raw)
  emit('update:modelValue', next)
}

function formatTaxRate(value: number) {
  if (value === 0) {
    return '0'
  }
  return String(value)
}
</script>

<style scoped>
.sheet-shell {
  border: 1px solid var(--line);
  border-radius: 18px;
  background: #ffffff;
  overflow: hidden;
}

.sheet-scroll {
  overflow-x: auto;
}

.report-sheet {
  width: 100%;
  min-width: 1120px;
  border-collapse: collapse;
  table-layout: fixed;
}

.report-sheet th,
.report-sheet td {
  border: 1px solid #d7dee7;
  padding: 8px 10px;
  font-size: 13px;
}

.corner-head,
.row-head {
  width: 56px;
  background: #eff3f7;
  color: #526171;
  text-align: center;
  font-weight: 700;
}

.section-head {
  background: #d9e6f2;
  color: #15324b;
  font-size: 15px;
  letter-spacing: 0.04em;
}

.sub-head-row th {
  background: #f3f6fa;
  color: #5c6978;
}

.sub-head {
  font-weight: 700;
  text-align: left;
}

.muted {
  color: #7c8794;
}

.label-cell {
  background: #f9fbfd;
  color: #314457;
  font-weight: 600;
}

.value-cell {
  background: #fff9d8;
  text-align: right;
  font-variant-numeric: tabular-nums;
}

.auto-cell {
  background: #fff1a8;
  font-weight: 700;
}

.manual-cell {
  background: #d9f0c4;
}

.blank-cell {
  background: #ffffff;
}

.spacer-cell,
.spacer-head {
  width: 22px;
  min-width: 22px;
  background: #eef3f8;
  padding: 0;
}

.manual-cell input,
.manual-cell select {
  width: 100%;
  border: none;
  background: transparent;
  padding: 0;
  font: inherit;
  color: #163724;
  text-align: right;
  outline: none;
}

.manual-cell input:disabled,
.manual-cell select:disabled {
  color: #375240;
  opacity: 1;
}

@media (max-width: 1024px) {
  .report-sheet th,
  .report-sheet td {
    padding: 7px 8px;
    font-size: 12px;
  }
}
</style>