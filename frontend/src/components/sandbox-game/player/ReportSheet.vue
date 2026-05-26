<template>
  <div class="sheet-frame">
    <div class="sheet-scroll">
      <table class="sheet-table">
        <colgroup>
          <col class="col-index" />
          <col class="col-label" />
          <col class="col-value" />
          <col class="col-spacer" />
          <col class="col-label" />
          <col class="col-value" />
          <col class="col-spacer" />
          <col class="col-label" />
          <col class="col-value" />
        </colgroup>
        <thead>
          <tr>
            <th class="corner"></th>
            <th class="col-head">A</th>
            <th class="col-head">B</th>
            <th class="col-head">C</th>
            <th class="col-head">D</th>
            <th class="col-head">E</th>
            <th class="col-head">F</th>
            <th class="col-head">G</th>
            <th class="col-head">H</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <th class="row-head">1</th>
            <td class="sheet-title" colspan="8">
              <div class="title-main">年度财报表</div>
              <div class="title-sub">黄色单元格为系统计算结果，绿色单元格为玩家手工填写项</div>
            </td>
          </tr>
          <tr>
            <th class="row-head">2</th>
            <td class="section-head" colspan="2">损益表</td>
            <td class="spacer-cell"></td>
            <td class="section-head" colspan="2">资产</td>
            <td class="spacer-cell"></td>
            <td class="section-head" colspan="2">负债与权益</td>
          </tr>
          <tr>
            <th class="row-head">3</th>
            <td class="sub-head">项目</td>
            <td class="sub-head amount-head">金额</td>
            <td class="spacer-cell"></td>
            <td class="sub-head">项目</td>
            <td class="sub-head amount-head">金额</td>
            <td class="spacer-cell"></td>
            <td class="sub-head">项目</td>
            <td class="sub-head amount-head">金额</td>
          </tr>

          <tr v-for="row in sheetRows" :key="row.rowNo">
            <th class="row-head">{{ row.rowNo }}</th>
            <td :class="labelCellClass(row.profit)">{{ row.profit.label || '' }}</td>
            <td :class="valueCellClass(row.profit)">
              <template v-if="row.profit.kind === 'manual-select'">
                <select :value="manualSelectValue(row.profit.manualKey)" :disabled="!canEdit" @change="updateTaxRate($event)">
                  <option value="">请选择</option>
                  <option v-for="item in taxRateOptions" :key="item" :value="String(item)">{{ formatTaxRate(item) }}</option>
                </select>
              </template>
              <template v-else-if="row.profit.kind === 'manual-number'">
                <input
                  :value="manualInputValue(row.profit.manualKey)"
                  :disabled="!canEdit"
                  inputmode="decimal"
                  @input="updateManualNumber(row.profit.manualKey, $event)"
                />
              </template>
              <template v-else-if="row.profit.computedKey">
                {{ formatComputedCell(row.profit) }}
              </template>
            </td>
            <td class="spacer-cell"></td>
            <td :class="labelCellClass(row.asset)">{{ row.asset.label || '' }}</td>
            <td :class="valueCellClass(row.asset)">
              <template v-if="row.asset.kind === 'manual-number'">
                <input
                  :value="manualInputValue(row.asset.manualKey)"
                  :disabled="!canEdit"
                  inputmode="decimal"
                  @input="updateManualNumber(row.asset.manualKey, $event)"
                />
              </template>
              <template v-else-if="row.asset.kind === 'manual-select'">
                <select :value="manualSelectValue(row.asset.manualKey)" :disabled="!canEdit" @change="updateTaxRate($event)">
                  <option value="">请选择</option>
                  <option v-for="item in taxRateOptions" :key="item" :value="String(item)">{{ formatTaxRate(item) }}</option>
                </select>
              </template>
              <template v-else-if="row.asset.computedKey">
                {{ formatComputedCell(row.asset) }}
              </template>
            </td>
            <td class="spacer-cell"></td>
            <td :class="labelCellClass(row.liability)">{{ row.liability.label || '' }}</td>
            <td :class="valueCellClass(row.liability)">
              <template v-if="row.liability.kind === 'manual-number'">
                <input
                  :value="manualInputValue(row.liability.manualKey)"
                  :disabled="!canEdit"
                  inputmode="decimal"
                  @input="updateManualNumber(row.liability.manualKey, $event)"
                />
              </template>
              <template v-else-if="row.liability.kind === 'manual-select'">
                <select :value="manualSelectValue(row.liability.manualKey)" :disabled="!canEdit" @change="updateTaxRate($event)">
                  <option value="">请选择</option>
                  <option v-for="item in taxRateOptions" :key="item" :value="String(item)">{{ formatTaxRate(item) }}</option>
                </select>
              </template>
              <template v-else-if="row.liability.computedKey">
                {{ formatComputedCell(row.liability) }}
              </template>
            </td>
          </tr>

          <tr>
            <th class="row-head">30</th>
            <td class="note-cell" colspan="2">资产负债平衡校验</td>
            <td class="spacer-cell"></td>
            <td class="label-cell footer-label">总资产</td>
            <td class="summary-cell">{{ formatNumber(computedPayload.reportTotalAssets) }}</td>
            <td class="spacer-cell"></td>
            <td class="label-cell footer-label">总负债和权益</td>
            <td class="summary-cell">{{ formatNumber(computedPayload.reportTotalLiabilityEquity) }}</td>
          </tr>
          <tr>
            <th class="row-head">31</th>
            <td class="note-cell" colspan="2">提交前需满足：总资产 = 总负债 + 总权益</td>
            <td class="spacer-cell"></td>
            <td class="label-cell footer-label">差额</td>
            <td :class="balanceValueCellClass">{{ formatNumber(balanceGap) }}</td>
            <td class="spacer-cell"></td>
            <td class="note-cell" colspan="2">税率仅允许选择 0.25 / 0.15 / 0，绿色得分项需显式填写</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import { serviceReportLabels } from '@/configs/sandbox-game-service-labels'
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
  suffix?: string
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
  hasInvalidDraft: boolean
  taxRateOptions: number[]
}>()

const emit = defineEmits<{
  (event: 'update:modelValue', value: ReportManualPayload): void
}>()

const blankCell: SheetCell = { kind: 'blank' }

const sheetRows: SheetRow[] = [
  { rowNo: 4, profit: { label: '销售收入', kind: 'computed', computedKey: 'reportSalesRevenue' }, asset: { label: '非流动资产', kind: 'blank' }, liability: { label: '负债', kind: 'blank' } },
  { rowNo: 5, profit: { label: '直接成本', kind: 'computed', computedKey: 'reportDirectCost' }, asset: { label: serviceReportLabels.workInConstruction, kind: 'computed', computedKey: 'reportWorkInConstruction' }, liability: blankCell },
  { rowNo: 6, profit: blankCell, asset: { label: serviceReportLabels.factoryAsset, kind: 'computed', computedKey: 'reportFactoryAsset' }, liability: blankCell },
  { rowNo: 7, profit: { label: '毛利', kind: 'computed', computedKey: 'reportGrossProfit' }, asset: { label: serviceReportLabels.lineResidual, kind: 'computed', computedKey: 'reportLineResidual' }, liability: { label: '短期负债', kind: 'computed', computedKey: 'reportShortTermLiability' } },
  { rowNo: 8, profit: blankCell, asset: { label: '待折资产', kind: 'computed', computedKey: 'reportDepreciableAsset' }, liability: { label: '长期负债', kind: 'computed', computedKey: 'reportLongTermLiability' } },
  { rowNo: 9, profit: blankCell, asset: { label: '总非流动资产', kind: 'computed', computedKey: 'reportTotalNonCurrentAssets' }, liability: blankCell },
  { rowNo: 10, profit: { label: '综合费用', kind: 'computed', computedKey: 'reportComprehensiveCost' }, asset: blankCell, liability: { label: '总负债', kind: 'computed', computedKey: 'reportTotalLiability' } },
  { rowNo: 11, profit: { label: '折旧', kind: 'computed', computedKey: 'reportDepreciation' }, asset: { label: '流动资产', kind: 'blank' }, liability: blankCell },
  { rowNo: 12, profit: blankCell, asset: blankCell, liability: blankCell },
  { rowNo: 13, profit: { label: '营业利润', kind: 'computed', computedKey: 'reportOperatingProfit' }, asset: { label: '现金', kind: 'computed', computedKey: 'reportCash' }, liability: { label: '总权益', kind: 'blank' } },
  { rowNo: 14, profit: blankCell, asset: { label: '应收款', kind: 'computed', computedKey: 'reportReceivable' }, liability: blankCell },
  { rowNo: 15, profit: blankCell, asset: { label: serviceReportLabels.workInProgress, kind: 'manual-number', manualKey: 'workInProgress' }, liability: { label: '股东资本', kind: 'computed', computedKey: 'reportShareCapital' } },
  { rowNo: 16, profit: { label: '财务收入/支出', kind: 'computed', computedKey: 'reportFinanceIncomeExpense' }, asset: { label: serviceReportLabels.finishedGoods, kind: 'manual-number', manualKey: 'finishedGoods' }, liability: { label: '利润留存', kind: 'computed', computedKey: 'reportRetainedEarnings' } },
  { rowNo: 17, profit: { label: '额外收入/支出', kind: 'computed', computedKey: 'reportExtraIncomeExpense' }, asset: { label: serviceReportLabels.rawMaterials, kind: 'manual-number', manualKey: 'rawMaterials' }, liability: { label: '年度净利润', kind: 'computed', computedKey: 'reportNetProfit' } },
  { rowNo: 18, profit: blankCell, asset: { label: '所得税后现金', kind: 'computed', computedKey: 'reportPostTaxCash' }, liability: blankCell },
  { rowNo: 19, profit: { label: serviceReportLabels.row19Label, kind: 'computed', computedKey: 'reportPreTaxProfit' }, asset: { label: '总流动资产', kind: 'computed', computedKey: 'reportTotalCurrentAssets' }, liability: { label: '总股东权益', kind: 'computed', computedKey: 'reportTotalEquity' } },
  { rowNo: 20, profit: blankCell, asset: blankCell, liability: blankCell },
  { rowNo: 21, profit: { label: '所得税', kind: 'computed', computedKey: 'reportIncomeTax' }, asset: blankCell, liability: blankCell },
  { rowNo: 22, profit: { label: '年度净利润', kind: 'computed', computedKey: 'reportNetProfit' }, asset: { label: '总资产', kind: 'computed', computedKey: 'reportTotalAssets' }, liability: { label: serviceReportLabels.totalLiabilityEquity, kind: 'computed', computedKey: 'reportTotalLiabilityEquity' } },
  { rowNo: 23, profit: blankCell, asset: blankCell, liability: blankCell },
  { rowNo: 24, profit: { label: '所得税税率', kind: 'manual-select', manualKey: 'incomeTaxRate' }, asset: { label: '最佳市场经营总监得分', kind: 'computed', computedKey: 'reportBestMarketDirectorBaseScore', suffix: '+' }, liability: { label: '企业认证得分', kind: 'manual-number', manualKey: 'enterpriseCertificationScore' } },
  { rowNo: 25, profit: blankCell, asset: { label: '最佳科技创新总监得分', kind: 'computed', computedKey: 'reportBestTechnologyDirectorScore' }, liability: blankCell },
  { rowNo: 26, profit: blankCell, asset: { label: serviceReportLabels.bestProductionHumanDirector, kind: 'manual-number', manualKey: 'productionHumanScore' }, liability: blankCell },
  { rowNo: 27, profit: blankCell, asset: { label: '最佳销售总监得分', kind: 'computed', computedKey: 'reportBestSalesDirectorScore' }, liability: blankCell },
  { rowNo: 28, profit: blankCell, asset: { label: '最佳 CFO（财务总监）得分', kind: 'computed', computedKey: 'reportBestCfoBaseScore', suffix: '+' }, liability: { label: '关账速度得分', kind: 'manual-number', manualKey: 'closingSpeedScore' } },
  { rowNo: 29, profit: blankCell, asset: { label: '最佳 CEO（总经理）得分', kind: 'computed', computedKey: 'reportBestCeoScore' }, liability: blankCell },
]

const balanceGap = computed(() => {
  const assets = props.computedPayload.reportTotalAssets ?? 0
  const liabilitiesEquity = props.computedPayload.reportTotalLiabilityEquity ?? 0
  return assets - liabilitiesEquity
})

const balanceValueCellClass = computed(() => ({
  'value-cell': true,
  'summary-cell': true,
  'balance-pass': Math.abs(balanceGap.value) <= 0.000001,
  'balance-fail': Math.abs(balanceGap.value) > 0.000001,
}))

function labelCellClass(cell: SheetCell) {
  return {
    'label-cell': true,
    'blank-label-cell': cell.kind === 'blank',
    'manual-label-cell': cell.kind === 'manual-number' || cell.kind === 'manual-select',
  }
}

function valueCellClass(cell: SheetCell) {
  const isManualCell = cell.kind === 'manual-number' || cell.kind === 'manual-select'
  return {
    'value-cell': true,
    'auto-cell': cell.kind === 'computed',
    'manual-cell': isManualCell,
    'manual-invalid-cell': isManualCell && props.hasInvalidDraft,
    'blank-cell': cell.kind === 'blank',
  }
}

function formatNumber(value: number | null | undefined) {
  if (typeof value !== 'number' || Number.isNaN(value)) {
    return ''
  }
  return value.toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

function formatComputedCell(cell: SheetCell) {
  if (!cell.computedKey) {
    return ''
  }
  const value = formatNumber(props.computedPayload[cell.computedKey])
  if (!value) {
    return ''
  }
  return cell.suffix ? `${value} ${cell.suffix}` : value
}

function manualInputValue(key?: ManualKey) {
  if (!key) {
    return ''
  }
  const value = props.modelValue[key]
  return value === null ? '' : String(value)
}

function manualSelectValue(key?: ManualKey) {
  if (!key) {
    return ''
  }
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
.sheet-frame {
  border: 1px solid var(--line);
  border-radius: 18px;
  background: #ffffff;
  overflow: hidden;
}

.sheet-scroll {
  overflow-x: auto;
}

.sheet-table {
  width: 100%;
  min-width: 1080px;
  border-collapse: collapse;
  table-layout: fixed;
}

.sheet-table th,
.sheet-table td {
  border: 1px solid #d6dfea;
  padding: 8px 10px;
  font-size: 13px;
}

.col-index {
  width: 52px;
}

.col-label {
  width: 210px;
}

.col-value {
  width: 128px;
}

.col-spacer {
  width: 22px;
}

.corner,
.row-head {
  background: #eef3f8;
  color: #5b6978;
  text-align: center;
  font-weight: 700;
}

.col-head {
  background: #eef3f8;
  color: #5b6978;
  text-align: center;
  font-weight: 700;
}

.sheet-title {
  background: linear-gradient(180deg, #f9fbfe 0%, #edf3f9 100%);
  text-align: center;
  padding: 14px 16px;
}

.title-main {
  font-size: 18px;
  font-weight: 700;
  color: #18324a;
}

.title-sub {
  margin-top: 4px;
  font-size: 12px;
  color: #66788b;
}

.section-head {
  background: #dbe8f4;
  color: #173753;
  text-align: center;
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.04em;
}

.sub-head {
  background: #f4f7fb;
  color: #566779;
  font-weight: 700;
}

.amount-head {
  text-align: right;
}

.label-cell {
  background: #fbfcfe;
  color: #33475b;
  font-weight: 600;
}

.manual-label-cell {
  background: #eef7e5;
}

.blank-label-cell {
  background: #ffffff;
}

.value-cell {
  text-align: right;
  font-variant-numeric: tabular-nums;
}

.auto-cell {
  background: #fff1a9;
  color: #584300;
  font-weight: 700;
}

.manual-cell {
  background: #d7efc3;
}

.manual-invalid-cell {
  background: #f7e3c1;
}

.blank-cell {
  background: #ffffff;
}

.summary-cell {
  background: #e8f2fd;
  color: #173753;
  font-weight: 700;
}

.note-cell {
  background: #f7f9fc;
  color: #607284;
}

.footer-label {
  background: #f3f6fa;
}

.balance-pass {
  background: #eaf8ee;
  color: #1f6b40;
}

.balance-fail {
  background: #fff4f4;
  color: #b24040;
}

.spacer-cell {
  background: #edf2f7;
  padding: 0;
}

.manual-cell input,
.manual-cell select,
.manual-invalid-cell input,
.manual-invalid-cell select {
  width: 100%;
  border: none;
  background: transparent;
  padding: 0;
  font: inherit;
  color: #1c3e28;
  text-align: right;
  outline: none;
}

.manual-cell input:disabled,
.manual-cell select:disabled,
.manual-invalid-cell input:disabled,
.manual-invalid-cell select:disabled {
  color: #3f5d49;
  opacity: 1;
}

.manual-invalid-cell input,
.manual-invalid-cell select,
.manual-invalid-cell input:disabled,
.manual-invalid-cell select:disabled {
  color: #7a5419;
}

@media (max-width: 1024px) {
  .sheet-table th,
  .sheet-table td {
    padding: 7px 8px;
    font-size: 12px;
  }

  .title-main {
    font-size: 16px;
  }
}
</style>

