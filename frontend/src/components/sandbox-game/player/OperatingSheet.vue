<template>
  <div class="sheet-wrap">
    <section class="sheet-card year-start-card">
      <div class="sheet-card-head">年初区</div>
      <table class="sheet-table simple-table">
        <thead>
          <tr>
            <th>项目</th>
            <th>填写值</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="field in yearStartFields" :key="field.key">
            <th>{{ field.label }}</th>
            <td :class="cellClass(isScopeEditable('YEAR_START'))">
              <input
                :value="displayCell(modelValue.beginning.taxAndPlanning[field.key])"
                :disabled="!isScopeEditable('YEAR_START')"
                placeholder="请输入"
                @input="updateYearStartField(field.key, $event)"
              />
            </td>
          </tr>
        </tbody>
      </table>

      <div class="sub-block">
        <div class="sub-head">
          <strong>市场投入 / 订单总额</strong>
          <button type="button" class="mini-btn" :disabled="!isScopeEditable('YEAR_START')" @click="addMarketBidRow">新增一行</button>
        </div>
        <table class="sheet-table market-table">
          <thead>
            <tr>
              <th>条目</th>
              <th>市场投入</th>
              <th>订单总额</th>
              <th class="action-col">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in marketBidRows" :key="`market-${index}`">
              <th>条目 {{ index + 1 }}</th>
              <td :class="cellClass(isScopeEditable('YEAR_START'))">
                <input
                  :value="displayCell(item.marketInvestment)"
                  :disabled="!isScopeEditable('YEAR_START')"
                  placeholder="请输入"
                  @input="updateMarketBidField(index, 'marketInvestment', $event)"
                />
              </td>
              <td :class="cellClass(isScopeEditable('YEAR_START'))">
                <input
                  :value="displayCell(item.orderAmount)"
                  :disabled="!isScopeEditable('YEAR_START')"
                  placeholder="请输入"
                  @input="updateMarketBidField(index, 'orderAmount', $event)"
                />
              </td>
              <td class="action-cell">
                <button type="button" class="mini-btn danger" :disabled="!isScopeEditable('YEAR_START') || marketBidRows.length === 1" @click="removeMarketBidRow(index)">删除</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section class="sheet-card quarter-card">
      <div class="sheet-card-head">季度经营区</div>
      <table class="sheet-table quarter-table">
        <thead>
          <tr>
            <th class="sticky-col">项目</th>
            <th>Q1</th>
            <th>Q2</th>
            <th>Q3</th>
            <th>Q4</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="group in quarterFieldGroups" :key="group.title">
            <tr class="group-row">
              <th colspan="5">{{ group.title }}</th>
            </tr>
            <tr v-for="field in group.fields" :key="`${group.title}-${field.key}`">
              <th class="sticky-col">{{ field.label }}</th>
              <td v-for="quarter in quarterList" :key="`${group.title}-${field.key}-${quarter.key}`" :class="cellClass(isScopeEditable(quarter.scope))">
                <input
                  :value="displayCell(getQuarterFieldValue(group.source, quarter.key, field.key))"
                  :disabled="!isScopeEditable(quarter.scope)"
                  placeholder="请输入"
                  @input="updateQuarterField(group.source, quarter.key, field.key, $event)"
                />
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </section>

    <section class="sheet-card quarter-cash-card">
      <div class="sheet-card-head">核对季末现金</div>
      <table class="sheet-table cash-check-table">
        <tbody>
          <tr>
            <th>Q1</th>
            <td class="calc-cell">{{ formatNumber(quarterCashChecks.Q1) }}</td>
            <th>Q2</th>
            <td class="calc-cell">{{ formatNumber(quarterCashChecks.Q2) }}</td>
            <th>Q3</th>
            <td class="calc-cell">{{ formatNumber(quarterCashChecks.Q3) }}</td>
            <th>Q4</th>
            <td class="calc-cell">{{ formatNumber(quarterCashChecks.Q4) }}</td>
          </tr>
        </tbody>
      </table>
    </section>

    <section class="sheet-card year-end-card">
      <div class="sheet-card-head">年末区</div>
      <table class="sheet-table simple-table">
        <thead>
          <tr>
            <th>项目</th>
            <th>填写值</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="group in yearEndGroups" :key="group.title">
            <tr class="group-row">
              <th colspan="2">{{ group.title }}</th>
            </tr>
            <tr v-for="field in group.fields" :key="`${group.title}-${field.key}`">
              <th>{{ field.label }}</th>
              <td :class="cellClass(isScopeEditable('YEAR_END'))">
                <input
                  :value="displayCell(getYearEndFieldValue(group.source, field.key))"
                  :disabled="!isScopeEditable('YEAR_END')"
                  placeholder="请输入"
                  @input="updateYearEndField(group.source, field.key, $event)"
                />
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </section>

    <section class="sheet-card derived-card">
      <div class="sheet-card-head">关键派生指标</div>
      <table class="sheet-table derived-table">
        <thead>
          <tr>
            <th>指标</th>
            <th>系统结果</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="field in derivedFields" :key="field.key">
            <th>{{ field.label }}</th>
            <td class="calc-cell">{{ formatNumber(derivedValues[field.key]) }}</td>
          </tr>
          <tr>
            <th>当前期末现金</th>
            <td class="calc-cell strong">{{ formatNumber(periodEndCash) }}</td>
          </tr>
        </tbody>
      </table>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import { cloneOperatingPayload, type NumericCellValue, type OperatingPayload, type QuarterValueMap } from '@/types/sandbox-game'

const props = defineProps<{
  modelValue: OperatingPayload
  editableScopes: string[]
  quarterCashChecks: Record<string, number>
  derivedValues: Record<string, number>
  periodEndCash: number
}>()

const emit = defineEmits<{
  (event: 'update:modelValue', value: OperatingPayload): void
}>()

const yearStartFields = [
  { key: 'taxPayment', label: '年初缴税' },
  { key: 'planRevenue', label: '计划收入' },
  { key: 'comprehensiveCostPlan', label: '综合费用计划' },
]

const quarterList = [
  { key: 'q1', scope: 'Q1' },
  { key: 'q2', scope: 'Q2' },
  { key: 'q3', scope: 'Q3' },
  { key: 'q4', scope: 'Q4' },
]

const quarterFieldGroups = [
  {
    title: '短期贷款',
    source: 'shortTermLoan',
    fields: [
      { key: 'dueRepayment', label: '到期还贷' },
      { key: 'interest', label: '付利息' },
      { key: 'newLoan', label: '新增贷款' },
    ],
  },
  {
    title: '材料费',
    source: 'materialPayment',
    fields: [{ key: 'materialCost', label: '材料费' }],
  },
  {
    title: '生产线调整',
    source: 'productionLineAdjustment',
    fields: [
      { key: 'changeProduct', label: '变更产品' },
      { key: 'dismantleCost', label: '拆除费用' },
      { key: 'lineSale', label: '出售生产线' },
      { key: 'newLineInstall', label: '新安装生产线' },
      { key: 'constructionToFixed', label: '转入固定资产' },
      { key: 'newDepreciableAsset', label: '新增待折资产' },
    ],
  },
  {
    title: '人力资源',
    source: 'humanResource',
    fields: [{ key: 'staffCost', label: '人力资源费用' }],
  },
  {
    title: '工资与生产',
    source: 'salaryAndProduction',
    fields: [{ key: 'salaryCost', label: '工资与生产费用' }],
  },
  {
    title: '研发与管理',
    source: 'researchAndManagement',
    fields: [
      { key: 'technologyResearch', label: '技术研发' },
      { key: 'managementSystem', label: '管理体系' },
    ],
  },
  {
    title: '应收更新',
    source: 'receivableUpdate',
    fields: [{ key: 'receivableCollection', label: '应收回款' }],
  },
  {
    title: '交货结算',
    source: 'deliverySettlement',
    fields: [
      { key: 'salesRevenue', label: '销售收入' },
      { key: 'directCost', label: '直接成本' },
      { key: 'managementStaffCost', label: '管理人员费用' },
    ],
  },
  {
    title: '额外收入 / 罚款',
    source: 'incomeAndPenalty',
    fields: [
      { key: 'discountExpense', label: '折现费用' },
      { key: 'extraExpensePenalty', label: '额外支出 / 罚款' },
      { key: 'extraIncomeReward', label: '额外收入 / 奖励' },
    ],
  },
] as const

const yearEndGroups = [
  {
    title: '长期贷款',
    source: 'longTermLoan',
    fields: [
      { key: 'interest', label: '长期贷款利息' },
      { key: 'repayment', label: '长期贷款还款' },
      { key: 'newLoan', label: '新增长期贷款' },
    ],
  },
  {
    title: '资产调整',
    source: 'assetAdjustment',
    fields: [
      { key: 'lineMaintenance', label: '维护费' },
      { key: 'purchase', label: '厂房购置' },
      { key: 'sale', label: '厂房出售' },
      { key: 'rent', label: '厂房租金' },
      { key: 'workInConstruction', label: '在建工程' },
      { key: 'marketCultivation', label: '新市场培育' },
    ],
  },
] as const

const derivedFields = [
  { key: 'marketBidCost', label: '市场投入合计' },
  { key: 'orderTotal', label: '订单总额' },
  { key: 'comprehensiveCostTotal', label: '综合费用合计' },
  { key: 'shortTermLoanBalance', label: '短期贷款余额' },
  { key: 'longTermLoanBalance', label: '长期贷款余额' },
  { key: 'lineResidual', label: '生产线残值' },
  { key: 'depreciableAssetTotal', label: '待折资产合计' },
  { key: 'depreciation', label: '折旧' },
  { key: 'financeIncomeExpense', label: '财务收入 / 支出' },
  { key: 'receivableChange', label: '应收变动' },
  { key: 'extraIncomeExpense', label: '额外收入 / 支出' },
]

const marketBidRows = computed(() => {
  if (props.modelValue.beginning.marketBid.length > 0) {
    return props.modelValue.beginning.marketBid
  }
  return [{ marketInvestment: '', orderAmount: '' }]
})

function isScopeEditable(scope: string) {
  return props.editableScopes.includes(scope)
}

function cellClass(editable: boolean) {
  return editable ? 'input-cell' : 'readonly-cell'
}

function displayCell(value: string | number | undefined | null) {
  if (value === '' || value === undefined || value === null) {
    return ''
  }
  return String(value)
}

function normalizeNumericValue(raw: string): NumericCellValue {
  const value = raw.trim()
  if (value === '') {
    return ''
  }
  const parsed = Number(value)
  return Number.isNaN(parsed) ? '' : parsed
}

function updateYearStartField(key: string, event: Event) {
  const next = cloneOperatingPayload(props.modelValue)
  next.beginning.taxAndPlanning[key] = normalizeNumericValue((event.target as HTMLInputElement).value)
  emit('update:modelValue', next)
}

function addMarketBidRow() {
  const next = cloneOperatingPayload(props.modelValue)
  next.beginning.marketBid.push({ marketInvestment: '', orderAmount: '' })
  emit('update:modelValue', next)
}

function removeMarketBidRow(index: number) {
  const next = cloneOperatingPayload(props.modelValue)
  if (next.beginning.marketBid.length <= 1) {
    return
  }
  next.beginning.marketBid.splice(index, 1)
  emit('update:modelValue', next)
}

function updateMarketBidField(index: number, key: string, event: Event) {
  const next = cloneOperatingPayload(props.modelValue)
  while (next.beginning.marketBid.length <= index) {
    next.beginning.marketBid.push({ marketInvestment: '', orderAmount: '' })
  }
  next.beginning.marketBid[index][key] = normalizeNumericValue((event.target as HTMLInputElement).value)
  emit('update:modelValue', next)
}

function getQuarterMap(source: string): QuarterValueMap {
  switch (source) {
    case 'shortTermLoan':
      return props.modelValue.quarter.shortTermLoan
    case 'materialPayment':
      return props.modelValue.quarter.materialPayment
    case 'productionLineAdjustment':
      return props.modelValue.quarter.productionLineAdjustment
    case 'humanResource':
      return props.modelValue.quarter.humanResource
    case 'salaryAndProduction':
      return props.modelValue.quarter.salaryAndProduction
    case 'researchAndManagement':
      return props.modelValue.quarter.researchAndManagement
    case 'receivableUpdate':
      return props.modelValue.quarter.receivableUpdate
    case 'deliverySettlement':
      return props.modelValue.quarter.deliverySettlement
    case 'incomeAndPenalty':
      return props.modelValue.extra.incomeAndPenalty
    default:
      return {}
  }
}

function getQuarterFieldValue(source: string, quarterKey: string, fieldKey: string): NumericCellValue {
  const entry = getQuarterMap(source)[quarterKey]
  return entry?.[fieldKey] ?? ''
}

function updateQuarterField(source: string, quarterKey: string, fieldKey: string, event: Event) {
  const next = cloneOperatingPayload(props.modelValue)
  const nextValue = normalizeNumericValue((event.target as HTMLInputElement).value)

  const apply = (target: QuarterValueMap) => {
    target[quarterKey] = target[quarterKey] ?? {}
    target[quarterKey][fieldKey] = nextValue
  }

  switch (source) {
    case 'shortTermLoan':
      apply(next.quarter.shortTermLoan)
      break
    case 'materialPayment':
      apply(next.quarter.materialPayment)
      break
    case 'productionLineAdjustment':
      apply(next.quarter.productionLineAdjustment)
      break
    case 'humanResource':
      apply(next.quarter.humanResource)
      break
    case 'salaryAndProduction':
      apply(next.quarter.salaryAndProduction)
      break
    case 'researchAndManagement':
      apply(next.quarter.researchAndManagement)
      break
    case 'receivableUpdate':
      apply(next.quarter.receivableUpdate)
      break
    case 'deliverySettlement':
      apply(next.quarter.deliverySettlement)
      break
    case 'incomeAndPenalty':
      apply(next.extra.incomeAndPenalty)
      break
    default:
      break
  }

  emit('update:modelValue', next)
}

function getYearEndFieldValue(source: string, fieldKey: string): NumericCellValue {
  if (source === 'longTermLoan') {
    return props.modelValue.yearEnd.longTermLoan[fieldKey] ?? ''
  }
  return props.modelValue.yearEnd.assetAdjustment[fieldKey] ?? ''
}

function updateYearEndField(source: string, fieldKey: string, event: Event) {
  const next = cloneOperatingPayload(props.modelValue)
  const nextValue = normalizeNumericValue((event.target as HTMLInputElement).value)
  if (source === 'longTermLoan') {
    next.yearEnd.longTermLoan[fieldKey] = nextValue
  } else {
    next.yearEnd.assetAdjustment[fieldKey] = nextValue
  }
  emit('update:modelValue', next)
}

function formatNumber(value?: number) {
  if (typeof value !== 'number' || Number.isNaN(value)) {
    return '--'
  }
  return value.toLocaleString('zh-CN', { maximumFractionDigits: 2 })
}
</script>

<style scoped>
.sheet-wrap {
  display: grid;
  gap: 16px;
}

.sheet-card {
  background: #ffffff;
  border: 1px solid var(--line);
  border-radius: 18px;
  overflow: hidden;
}

.sheet-card-head {
  padding: 14px 18px;
  font-size: 16px;
  font-weight: 700;
  border-bottom: 1px solid var(--line);
}

.year-start-card .sheet-card-head {
  background: var(--year-start-bg);
}

.quarter-card .sheet-card-head {
  background: var(--quarter-bg);
}

.quarter-cash-card .sheet-card-head,
.derived-card .sheet-card-head {
  background: var(--derived-bg);
}

.year-end-card .sheet-card-head {
  background: var(--year-end-bg);
}

.sheet-table {
  width: 100%;
  border-collapse: collapse;
}

.sheet-table th,
.sheet-table td {
  border: 1px solid var(--line);
  padding: 0;
  vertical-align: middle;
}

.sheet-table thead th,
.group-row th {
  background: #f3f5f8;
  padding: 10px 12px;
  text-align: left;
}

.simple-table tbody th,
.derived-table tbody th,
.cash-check-table tbody th,
.quarter-table tbody th {
  width: 220px;
  min-width: 220px;
  background: #fbfcfd;
  text-align: left;
  padding: 10px 12px;
}

.group-row th {
  font-weight: 700;
}

.input-cell,
.readonly-cell,
.calc-cell {
  min-width: 140px;
}

.input-cell input,
.readonly-cell input {
  width: 100%;
  border: none;
  padding: 12px;
  background: transparent;
  outline: none;
}

.input-cell {
  background: var(--input-bg);
}

.readonly-cell {
  background: var(--readonly-bg);
}

.calc-cell {
  background: var(--calc-bg);
  padding: 12px;
  font-weight: 700;
}

.calc-cell.strong {
  font-size: 16px;
}

.sub-block {
  padding: 16px 18px 18px;
}

.sub-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  margin-bottom: 12px;
}

.market-table .action-col,
.action-cell {
  width: 88px;
  text-align: center;
}

.mini-btn {
  border: 1px solid var(--line-strong);
  background: #ffffff;
  border-radius: 10px;
  padding: 8px 10px;
}

.mini-btn.danger {
  color: var(--danger);
}

.quarter-table {
  table-layout: fixed;
}

.quarter-table td,
.quarter-table th {
  min-width: 140px;
}

.sticky-col {
  position: sticky;
  left: 0;
  z-index: 1;
}

.quarter-table .sticky-col {
  background: #fbfcfd;
}

.quarter-table .group-row .sticky-col,
.quarter-table .group-row th {
  background: #f3f5f8;
}

@media (max-width: 1280px) {
  .sheet-wrap {
    overflow-x: auto;
  }

  .sheet-card {
    min-width: 920px;
  }
}
</style>