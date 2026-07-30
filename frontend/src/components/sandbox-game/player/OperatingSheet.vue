<template>
  <div class="sheet-frame">
    <div class="sheet-scroll" data-enter-nav-scope @keydown="handleSequentialInputNavigation">
      <table class="sheet-table">
        <colgroup>
          <col class="col-index" />
          <col class="col-band" />
          <col class="col-main" />
          <col class="col-sub" />
          <col class="col-stage" />
          <col class="col-stage" />
          <col class="col-stage" />
          <col class="col-stage" />
          <col class="col-total" />
          <col class="col-side" />
          <col class="col-side" />
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
            <th class="col-head">I</th>
            <th class="col-head">J</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <th class="row-head">1</th>
            <td class="section-band band-year-start" rowspan="7">年初工作</td>
            <td class="group-title" colspan="2">1. 支付上年度所得税</td>
            <td class="note-cell rule-note rule-note-short" colspan="4">{{ labels.taxPolicyNote }}</td>
            <td class="result-cell">{{ formatNumber(taxPaymentDisplay) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">2</th>
            <td class="group-title" colspan="2">2. 召开年度经营会议</td>
            <td class="note-cell instruction-note" colspan="3">调整公司战略，制定当年经营和工作计划</td>
            <td class="note-cell item-label">计划收入</td>
            <td class="result-cell">{{ formatNumber(planRevenueDisplay) }}</td>
            <td class="note-cell item-label">综合费用</td>
            <td class="result-cell">{{ formatNumber(comprehensiveCostDisplay) }}</td>
          </tr>
          <tr>
            <th class="row-head">3</th>
            <td class="group-title" rowspan="5">3. 市场竞标</td>
            <td class="note-cell period-label">区域</td>
            <td v-for="field in marketProductFields" :key="`market-head-${field.key}`" class="note-cell period-label">{{ field.label }}</td>
            <td class="note-cell period-label">区域市场订单总价</td>
            <td class="note-cell period-label">订单总额</td>
            <td class="note-cell period-label">市场投入</td>
            <td class="empty-cell"></td>
          </tr>
          <tr v-for="(region, index) in marketRegions" :key="region.key">
            <th class="row-head">{{ 4 + index }}</th>
            <td class="note-cell item-label">{{ region.label }}</td>
            <td
              v-for="field in marketProductFields"
              :key="`${region.key}-${field.key}`"
              :class="marketBidReadonly ? 'linked-order-cell' : editableCellClass('YEAR_START', marketBidRows[index][field.key])"
            >
              <input
                :value="displayCell(marketBidRows[index][field.key])"
                :disabled="marketBidReadonly || !isScopeEditable('YEAR_START')"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateMarketBidField(index, field.key, $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getMarketRowTotal(index)) }}</td>
            <td v-if="index === 0" class="result-cell" :rowspan="marketRegions.length">{{ formatNumber(marketOrderTotal) }}</td>
            <td
              v-if="index === 0"
              :class="marketBidReadonly ? 'linked-order-cell' : editableCellClass('YEAR_START', marketInvestmentTotal)"
              :rowspan="marketRegions.length"
            >
              <input
                :value="displayCell(marketInvestmentTotal)"
                :disabled="marketBidReadonly || !isScopeEditable('YEAR_START')"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateMarketInvestmentTotal($event)"
              />
            </td>
            <td v-if="index === 0" class="empty-cell" :rowspan="marketRegions.length"></td>
          </tr>

          <tr>
            <th class="row-head">8</th>
            <td class="section-band band-quarter" rowspan="40">每季度工作</td>
            <td class="group-title" rowspan="4">1. 短期贷款更新账期</td>
            <td class="note-cell item-label"></td>
            <td v-for="quarter in quarterList" :key="`loan-head-${quarter.key}`" class="note-cell period-label">{{ quarter.label }}</td>
            <td class="note-cell period-label">总计</td>
            <td class="note-cell period-label">增减</td>
            <td class="empty-cell"></td>
          </tr>
          <tr v-for="(field, index) in shortTermLoanFields" :key="`loan-${field.key}`">
            <th class="row-head">{{ 9 + index }}</th>
            <td class="note-cell">{{ field.label }}</td>
            <td v-for="quarter in quarterList" :key="`loan-${field.key}-${quarter.key}`" :class="editableCellClass(quarter.scope, getQuarterFieldValue('shortTermLoan', quarter.key, field.key))">
              <input
                :value="displayCell(getQuarterFieldValue('shortTermLoan', quarter.key, field.key))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateQuarterField('shortTermLoan', quarter.key, field.key, $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('shortTermLoan', field.key)) }}</td>
            <td v-if="index === 0" class="orange-cell" :rowspan="shortTermLoanFields.length">{{ formatNumber(shortTermLoanDelta) }}</td>
            <td v-if="index === 0" class="empty-cell" :rowspan="shortTermLoanFields.length"></td>
          </tr>
          <tr>
            <th class="row-head">12</th>
            <td class="reminder-cell" colspan="9">请将短期贷款右移一格</td>
          </tr>

          <tr>
            <th class="row-head">13</th>
            <td class="group-title" rowspan="6">{{ labels.materialGroupTitle }}</td>
            <td class="note-cell item-label"></td>
            <td v-for="quarter in quarterList" :key="`material-head-${quarter.key}`" class="note-cell period-label">{{ quarter.label }}</td>
            <td class="note-cell period-label">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr v-for="(field, index) in materialFields" :key="`material-${field.key}`">
            <th class="row-head">{{ 14 + index }}</th>
            <td class="note-cell item-label">{{ field.label }}</td>
            <td v-for="quarter in quarterList" :key="`material-${field.key}-${quarter.key}`" :class="editableCellClass(quarter.scope, getQuarterFieldValue('materialPayment', quarter.key, field.key))">
              <input
                :value="displayCell(getQuarterFieldValue('materialPayment', quarter.key, field.key))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateQuarterField('materialPayment', quarter.key, field.key, $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('materialPayment', field.key)) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">18</th>
            <td class="result-title">合计</td>
            <td v-for="quarter in quarterList" :key="`material-total-${quarter.key}`" class="result-cell">{{ formatNumber(getQuarterGroupQuarterTotal('materialPayment', materialFieldKeys, quarter.key)) }}</td>
            <td class="result-cell">{{ formatNumber(getQuarterGroupGrandTotal('materialPayment', materialFieldKeys)) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">19</th>
            <td class="excel-nested-cell project-progress-cell" colspan="9">
              <table class="nested-excel-table project-progress-table">
                <tbody>
                  <tr>
                    <td class="nested-title-cell project-progress-title-cell" rowspan="5">
                      <div>3. 项目进度更新</div>
                      <span>（每个在建交付项目向前移动一格，满格的移动到项目交付库）</span>
                    </td>
                    <td class="nested-head-cell project-factory-head">生产厂房</td>
                    <td v-for="quarter in quarterList" :key="`project-quarter-${quarter.key}`" class="nested-head-cell project-quarter-head" colspan="2">
                      {{ quarter.label }}
                    </td>
                  </tr>
                  <tr>
                    <td class="nested-head-cell project-factory-subhead"></td>
                    <template v-for="quarter in quarterList" :key="`project-subhead-${quarter.key}`">
                      <td class="nested-subhead-cell">产线</td>
                      <td class="nested-subhead-cell">进度</td>
                    </template>
                  </tr>
                  <tr v-for="factory in projectFactoryRows" :key="factory.key">
                    <td class="nested-head-cell project-factory-cell">{{ factory.label }}</td>
                    <template v-for="quarter in quarterList" :key="`${factory.key}-${quarter.key}`">
                      <td :class="projectProgressCellClass(quarter.scope, getProjectProgressCell(factory.key, quarter.key, 'lineType'))">
                        <select
                          :value="getProjectProgressCell(factory.key, quarter.key, 'lineType')"
                          :disabled="!isScopeEditable(quarter.scope)"
                          @change="updateProjectProgressCell(factory.key, quarter.key, 'lineType', $event)"
                        >
                          <option value=""></option>
                          <option v-for="option in projectLineTypeOptions" :key="option" :value="option">{{ option }}</option>
                        </select>
                      </td>
                      <td :class="projectProgressCellClass(quarter.scope, getProjectProgressCell(factory.key, quarter.key, 'progress'))">
                        <input
                          :value="displayCell(getProjectProgressCell(factory.key, quarter.key, 'progress'))"
                          :disabled="!isScopeEditable(quarter.scope)"
                          inputmode="numeric"
                          pattern="[0-9-]*"
                          data-enter-nav
                          @input="updateProjectProgressCell(factory.key, quarter.key, 'progress', $event)"
                        />
                      </td>
                    </template>
                  </tr>
                </tbody>
              </table>
            </td>
          </tr>

          <tr>
            <th class="row-head">20</th>
            <td class="group-title" colspan="2">{{ labels.productionAdjustmentGroupTitle }}</td>
            <td v-for="quarter in quarterList" :key="`line-head-${quarter.key}`" class="note-cell period-label">{{ quarter.label }}</td>
            <td class="note-cell period-label">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr v-for="(field, index) in productionLineRows" :key="`line-${field.key}`">
            <th class="row-head">{{ 21 + index }}</th>
            <template v-if="field.groupLabel">
              <td class="nested-group-cell" :rowspan="field.groupRowspan ?? 1">{{ field.groupLabel }}</td>
              <td class="note-cell item-label">{{ field.label }}</td>
            </template>
            <template v-else-if="field.labelColspan === 2">
              <td class="note-cell item-label" colspan="2">{{ field.label }}</td>
            </template>
            <template v-else-if="field.useExistingGroupCell">
              <td class="note-cell item-label">{{ field.label }}</td>
            </template>
            <template v-else>
              <td class="note-cell item-label" colspan="2">{{ field.label }}</td>
            </template>
            <td v-for="quarter in quarterList" :key="`line-${field.key}-${quarter.key}`" :class="editableCellClass(quarter.scope, getQuarterFieldValue('productionLineAdjustment', quarter.key, field.key))">
              <input
                :value="displayCell(getQuarterFieldValue('productionLineAdjustment', quarter.key, field.key))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateQuarterField('productionLineAdjustment', quarter.key, field.key, $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('productionLineAdjustment', field.key)) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">27</th>
            <td class="group-title" rowspan="2">{{ labels.humanResourceGroupTitle }}</td>
            <td class="note-cell item-label"></td>
            <td v-for="quarter in quarterList" :key="`hr-head-${quarter.key}`" class="note-cell period-label">{{ quarter.label }}</td>
            <td class="note-cell period-label">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">28</th>
            <td class="note-cell item-label">{{ labels.humanResourceLabel }}</td>
            <td v-for="quarter in quarterList" :key="`hr-${quarter.key}`" :class="editableCellClass(quarter.scope, getQuarterFieldValue('humanResource', quarter.key, 'staffCost'))">
              <input
                :value="displayCell(getQuarterFieldValue('humanResource', quarter.key, 'staffCost'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateQuarterField('humanResource', quarter.key, 'staffCost', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('humanResource', 'staffCost')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">29</th>
            <td class="group-title" rowspan="2">{{ labels.salaryGroupTitle }}</td>
            <td class="note-cell item-label"></td>
            <td v-for="quarter in quarterList" :key="`salary-head-${quarter.key}`" class="note-cell period-label">{{ quarter.label }}</td>
            <td class="note-cell period-label">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">30</th>
            <td class="note-cell item-label">{{ labels.salaryLabel }}</td>
            <td v-for="quarter in quarterList" :key="`salary-${quarter.key}`" :class="editableCellClass(quarter.scope, getQuarterFieldValue('salaryAndProduction', quarter.key, 'salaryCost'))">
              <input
                :value="displayCell(getQuarterFieldValue('salaryAndProduction', quarter.key, 'salaryCost'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateQuarterField('salaryAndProduction', quarter.key, 'salaryCost', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('salaryAndProduction', 'salaryCost')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">31</th>
            <td class="group-title" rowspan="3">{{ labels.researchAndManagementGroupTitle }}</td>
            <td class="note-cell item-label"></td>
            <td v-for="quarter in quarterList" :key="`research-head-${quarter.key}`" class="note-cell period-label">{{ quarter.label }}</td>
            <td class="note-cell period-label">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr v-for="(field, index) in researchFields" :key="`research-${field.key}`">
            <th class="row-head">{{ 32 + index }}</th>
            <td class="note-cell item-label">{{ field.label }}</td>
            <td v-for="quarter in quarterList" :key="`research-${field.key}-${quarter.key}`" :class="editableCellClass(quarter.scope, getQuarterFieldValue('researchAndManagement', quarter.key, field.key))">
              <input
                :value="displayCell(getQuarterFieldValue('researchAndManagement', quarter.key, field.key))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateQuarterField('researchAndManagement', quarter.key, field.key, $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('researchAndManagement', field.key)) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">34</th>
            <td class="group-title" rowspan="5">{{ labels.newOrderReminder }}</td>
            <td class="note-cell rule-note rule-note-short">单位：个</td>
            <td v-for="quarter in quarterList" :key="`supply-order-head-${quarter.key}`" class="note-cell period-label">{{ quarter.label }}</td>
            <td class="note-cell period-label">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr v-for="(field, index) in supplyChainOrderFields" :key="`supply-order-${field.key}`">
            <th class="row-head">{{ 35 + index }}</th>
            <td class="note-cell item-label">{{ field.label }}</td>
            <td
              v-for="quarter in quarterList"
              :key="`supply-order-${field.key}-${quarter.key}`"
              :class="supplyChainOrderCellClass(quarter.scope, getQuarterFieldValue('supplyChainOrderRecord', quarter.key, field.key))"
            >
              <input
                :value="displayCell(getQuarterFieldValue('supplyChainOrderRecord', quarter.key, field.key))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="numeric"
                min="0"
                pattern="[0-9]*"
                data-enter-nav
                @input="updateQuarterField('supplyChainOrderRecord', quarter.key, field.key, $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('supplyChainOrderRecord', field.key)) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">39</th>
            <td class="group-title" rowspan="2">{{ labels.receivableGroupTitle }}</td>
            <td class="note-cell item-label"></td>
            <td v-for="quarter in quarterList" :key="`receivable-head-${quarter.key}`" class="note-cell period-label">{{ quarter.label }}</td>
            <td class="note-cell period-label">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">40</th>
            <td class="note-cell item-label">{{ labels.receivableLabel }}</td>
            <td v-for="quarter in quarterList" :key="`receivable-${quarter.key}`" :class="editableCellClass(quarter.scope, getQuarterFieldValue('receivableUpdate', quarter.key, 'receivableCollection'))">
              <input
                :value="displayCell(getQuarterFieldValue('receivableUpdate', quarter.key, 'receivableCollection'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateQuarterField('receivableUpdate', quarter.key, 'receivableCollection', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('receivableUpdate', 'receivableCollection')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">41</th>
            <td class="reminder-cell" colspan="9">{{ labels.receivableReminder }}</td>
          </tr>

          <tr>
            <th class="row-head">42</th>
            <td class="group-title" rowspan="3">{{ labels.deliveryGroupTitle }}</td>
            <td class="note-cell item-label"></td>
            <td v-for="quarter in quarterList" :key="`delivery-head-${quarter.key}`" class="note-cell period-label">{{ quarter.label }}</td>
            <td class="note-cell period-label">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">43</th>
            <td class="note-cell item-label">{{ labels.deliveryRevenueLabel }}</td>
            <td
              v-for="quarter in quarterList"
              :key="`sales-${quarter.key}`"
              :class="deliverySalesRevenueCellClass(quarter.scope, getQuarterFieldValue('deliverySettlement', quarter.key, 'salesRevenue'))"
            >
              <span v-if="isOrderSalesRevenueLinked()" class="system-linked-value">
                {{ displayLinkedSalesRevenue(quarter.key) }}
              </span>
              <input
                v-else
                :value="displayCell(getQuarterFieldValue('deliverySettlement', quarter.key, 'salesRevenue'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateQuarterField('deliverySettlement', quarter.key, 'salesRevenue', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('deliverySettlement', 'salesRevenue')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">44</th>
            <td class="note-cell item-label">{{ labels.deliveryCostLabel }}</td>
            <td v-for="quarter in quarterList" :key="`cost-${quarter.key}`" :class="editableCellClass(quarter.scope, getQuarterFieldValue('deliverySettlement', quarter.key, 'directCost'))">
              <input
                :value="displayCell(getQuarterFieldValue('deliverySettlement', quarter.key, 'directCost'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateQuarterField('deliverySettlement', quarter.key, 'directCost', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('deliverySettlement', 'directCost')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">45</th>
            <td class="group-title" rowspan="2">{{ labels.managementStaffGroupTitle }}</td>
            <td class="note-cell rule-note rule-note-short">每季 1M</td>
            <td v-for="quarter in quarterList" :key="`management-head-${quarter.key}`" class="note-cell period-label">{{ quarter.label }}</td>
            <td class="note-cell period-label">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">46</th>
            <td class="note-cell item-label">管理人员费用</td>
            <td v-for="quarter in quarterList" :key="`management-${quarter.key}`" :class="editableCellClass(quarter.scope, getQuarterFieldValue('deliverySettlement', quarter.key, 'managementStaffCost'))">
              <input
                :value="displayCell(getQuarterFieldValue('deliverySettlement', quarter.key, 'managementStaffCost'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateQuarterField('deliverySettlement', quarter.key, 'managementStaffCost', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('deliverySettlement', 'managementStaffCost')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">47</th>
            <td class="group-title" colspan="2">核对季末现金</td>
            <td class="quarter-cash-cell">{{ displayQuarterCash('Q1') }}</td>
            <td class="quarter-cash-cell">{{ displayQuarterCash('Q2') }}</td>
            <td class="quarter-cash-cell">{{ displayQuarterCash('Q3') }}</td>
            <td class="quarter-cash-cell">{{ displayQuarterCash('Q4') }}</td>
            <td class="result-cell">{{ formatNumber(periodEndCash) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">48</th>
            <td class="section-band band-year-end" rowspan="13">年末工作</td>
            <td class="group-title" rowspan="3">1. 办理长期贷款账期更新</td>
            <td class="note-cell item-label">付利息</td>
            <td class="note-cell rule-note rule-note-short" colspan="5">年利率 5%</td>
            <td :class="editableCellClass('YEAR_END', getYearEndFieldValue('longTermLoan', 'interest'))">
              <input
                :value="displayCell(getYearEndFieldValue('longTermLoan', 'interest'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateYearEndField('longTermLoan', 'interest', $event)"
              />
            </td>
            <td class="orange-cell" rowspan="3">{{ formatNumber(longTermLoanDelta) }}</td>
            <td class="empty-cell" rowspan="3"></td>
          </tr>
          <tr>
            <th class="row-head">49</th>
            <td class="note-cell item-label">到期还款</td>
            <td class="note-cell rule-note rule-note-short" colspan="5">------------------------------</td>
            <td :class="editableCellClass('YEAR_END', getYearEndFieldValue('longTermLoan', 'repayment'))">
              <input
                :value="displayCell(getYearEndFieldValue('longTermLoan', 'repayment'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateYearEndField('longTermLoan', 'repayment', $event)"
              />
            </td>
          </tr>
          <tr>
            <th class="row-head">50</th>
            <td class="note-cell item-label">办理新贷款</td>
            <td class="note-cell rule-note rule-note-short" colspan="5">------------------------------</td>
            <td :class="editableCellClass('YEAR_END', getYearEndFieldValue('longTermLoan', 'newLoan'))">
              <input
                :value="displayCell(getYearEndFieldValue('longTermLoan', 'newLoan'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateYearEndField('longTermLoan', 'newLoan', $event)"
              />
            </td>
          </tr>

          <tr>
            <th class="row-head">51</th>
            <td class="group-title">{{ labels.lineMaintenanceGroupTitle }}</td>
            <td class="note-cell rule-note rule-note-short">{{ labels.lineMaintenanceNote }}</td>
            <td class="empty-cell" colspan="5"></td>
            <td :class="editableCellClass('YEAR_END', getYearEndFieldValue('assetAdjustment', 'lineMaintenance'))">
              <input
                :value="displayCell(getYearEndFieldValue('assetAdjustment', 'lineMaintenance'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateYearEndField('assetAdjustment', 'lineMaintenance', $event)"
              />
            </td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">52</th>
            <td class="group-title" rowspan="2">{{ labels.assetGroupTitle }}</td>
            <td class="note-cell item-label">购买</td>
            <td class="note-cell rule-note rule-note-short" colspan="5">{{ labels.assetValueNote }}</td>
            <td :class="editableCellClass('YEAR_END', getYearEndFieldValue('assetAdjustment', 'purchase'))">
              <input
                :value="displayCell(getYearEndFieldValue('assetAdjustment', 'purchase'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateYearEndField('assetAdjustment', 'purchase', $event)"
              />
            </td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">53</th>
            <td class="note-cell item-label">出售</td>
            <td class="note-cell rule-note rule-note-short" colspan="5">{{ labels.assetValueNote }}</td>
            <td :class="editableCellClass('YEAR_END', getYearEndFieldValue('assetAdjustment', 'sale'))">
              <input
                :value="displayCell(getYearEndFieldValue('assetAdjustment', 'sale'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateYearEndField('assetAdjustment', 'sale', $event)"
              />
            </td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">54</th>
            <td class="group-title">{{ labels.rentGroupTitle }}</td>
            <td class="note-cell item-label">付租金</td>
            <td class="note-cell rule-note rule-note-short" colspan="5">{{ labels.rentValueNote }}</td>
            <td :class="editableCellClass('YEAR_END', getYearEndFieldValue('assetAdjustment', 'rent'))">
              <input
                :value="displayCell(getYearEndFieldValue('assetAdjustment', 'rent'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateYearEndField('assetAdjustment', 'rent', $event)"
              />
            </td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">55</th>
            <td class="group-title">{{ labels.residualGroupTitle }}</td>
            <td class="empty-cell" colspan="6"></td>
            <td class="result-cell">{{ formatNumber(derivedMetric('lineResidual')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">56</th>
            <td class="group-title" rowspan="3">{{ labels.depreciationGroupTitle }}</td>
            <td class="note-cell item-label">折旧前待折资产总价值</td>
            <td class="empty-cell" colspan="5"></td>
            <td class="result-cell">{{ formatNumber(derivedMetric('depreciableAssetTotal')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">57</th>
            <td class="note-cell item-label">折旧费</td>
            <td class="note-cell rule-note rule-note-short" colspan="5">按待折资产的 1 / 3 取整</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('depreciation')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">58</th>
            <td class="note-cell item-label">{{ labels.workInConstructionLabel }}</td>
            <td class="empty-cell" colspan="5"></td>
            <td :class="editableCellClass('YEAR_END', getYearEndFieldValue('assetAdjustment', 'workInConstruction'))">
              <input
                :value="displayCell(getYearEndFieldValue('assetAdjustment', 'workInConstruction'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="numeric"
                pattern="[0-9-]*"
                data-enter-nav
                @input="updateYearEndField('assetAdjustment', 'workInConstruction', $event)"
              />
            </td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">59</th>
            <td class="group-title">7. 新市场培育</td>
            <td class="excel-nested-cell market-cultivation-cell" colspan="9">
              <table class="nested-excel-table market-cultivation-table">
                <tbody>
                  <tr v-for="(market, index) in marketCultivationFields" :key="market.key">
                    <td v-if="index === 0" class="nested-title-cell market-rule-cell" :rowspan="marketCultivationFields.length">
                      每年可向区域、全国、全球各投1M
                    </td>
                    <td class="nested-head-cell market-name-cell">{{ market.label }}</td>
                    <td :class="marketCultivationInputClass(market.key)">
                    <input
                      :value="displayCell(getMarketCultivationItem(market.key).annualInvestment)"
                      :disabled="!isScopeEditable('YEAR_END') || getMarketCultivationItem(market.key).locked || getMarketCultivationItem(market.key).lockedByPrevious"
                      inputmode="numeric"
                      pattern="[0-1]*"
                      data-enter-nav
                      @input="updateMarketCultivationField(market.key, $event)"
                    />
                    </td>
                    <td class="nested-head-cell">总投入</td>
                    <td class="nested-result-cell">{{ formatNumber(getMarketCultivationDisplayCumulative(market.key)) }}M</td>
                    <td class="nested-head-cell">状态</td>
                    <td class="nested-status-cell">{{ marketCultivationStatusText(market.key) }}</td>
                    <td v-if="index === 0" class="nested-result-cell market-total-cell" :rowspan="marketCultivationFields.length">
                      本年合计<br />
                      {{ formatNumber(marketCultivationAnnualTotal) }}M
                    </td>
                  </tr>
                </tbody>
              </table>
            </td>
          </tr>

          <tr>
            <th class="row-head">60</th>
            <td class="group-title">8. 资质认证</td>
            <template v-for="item in qualificationFields" :key="item.key">
              <td class="qualification-label-cell">
                {{ item.label }}
              </td>
              <td :class="qualificationCellClass(item.key)">
                <select
                  :value="getQualificationItem(item.key).status"
                  :disabled="!isScopeEditable('YEAR_END') || getQualificationItem(item.key).locked || getQualificationItem(item.key).lockedByPrevious"
                  @change="updateQualificationStatus(item.key, $event)"
                >
                  <option value="未解锁">未解锁</option>
                  <option value="解锁">解锁</option>
                </select>
              </td>
            </template>
            <td class="empty-cell"></td>
          </tr>

          <tr>
            <th class="row-head">61</th>
            <td class="section-band band-misc" rowspan="4">其他收支</td>
            <td class="group-title" rowspan="4">额外收入 / 罚款</td>
            <td class="note-cell item-label"></td>
            <td v-for="period in extraPeriodList" :key="`extra-head-${period.key}`" class="note-cell period-label">{{ period.label }}</td>
            <td class="note-cell period-label">总计</td>
            <td class="empty-cell"></td>
          </tr>
          <tr v-for="(field, index) in extraFields" :key="`extra-${field.key}`">
            <th class="row-head">{{ 62 + index }}</th>
            <td class="note-cell item-label">{{ field.label }}</td>
            <td v-for="period in extraPeriodList" :key="`extra-${field.key}-${period.key}`" :class="extraCellClass(field.key, period.scope)">
              <template v-if="field.key === 'discountExpense' && period.key === 'year_end'">
                <span class="cell-readonly-value">--</span>
              </template>
              <template v-else-if="isAdminIssuedExtraField(field.key)">
                <span class="cell-readonly-value">{{ displayExtraReadonlyCell(period.key, field.key) }}</span>
              </template>
              <template v-else>
                <input
                  :value="displayCell(getQuarterFieldValue('incomeAndPenalty', period.key, field.key))"
                  :disabled="!isScopeEditable(period.scope)"
                  inputmode="numeric"
                  pattern="[0-9-]*"
                  data-enter-nav
                  @input="updateQuarterField('incomeAndPenalty', period.key, field.key, $event)"
                />
              </template>
            </td>
            <td class="result-cell">{{ formatNumber(getExtraRowTotal(field.key)) }}</td>
            <td class="empty-cell"></td>
          </tr>

          <tr>
            <th class="row-head">65</th>
            <td class="metric-label">期末现金</td>
            <td class="result-cell">{{ formatNumber(periodEndCash) }}</td>
            <td class="metric-label">现金收入</td>
            <td class="result-cell">{{ formatNumber(cashInflow) }}</td>
            <td class="metric-label">现金支出</td>
            <td class="result-cell">{{ formatNumber(cashOutflow) }}</td>
            <td class="metric-label">市场回报比</td>
            <td class="result-cell">{{ formatDisplayNumber(optionalDerivedMetric('marketReturnRatio')) }}</td>
            <td class="dual-metric-cell" colspan="2">
              <div class="metric-stack-line">
                <span>研发投入强度</span>
                <strong>{{ formatDisplayPercent(optionalDerivedMetric('researchIntensity')) }}</strong>
              </div>
              <div class="metric-stack-line">
                <span>劳动生产率</span>
                <strong>{{ formatDisplayNumber(optionalDerivedMetric('laborProductivity')) }}</strong>
              </div>
            </td>
          </tr>
          <tr>
            <th class="row-head">66</th>
            <td class="metric-label">财务收入/支出</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('financeIncomeExpense')) }}</td>
            <td class="metric-label">不动产增减值</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('factoryAssetChange')) }}</td>
            <td class="metric-label">待折资产增减</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('depreciableAssetChange')) }}</td>
            <td class="metric-label">净产收益率</td>
            <td class="result-cell">{{ formatDisplayPercent(optionalDerivedMetric('netAssetYield')) }}</td>
            <td class="metric-label">总产收益</td>
            <td class="result-cell">{{ formatDisplayPercent(optionalDerivedMetric('totalAssetYield')) }}</td>
          </tr>
          <tr>
            <th class="row-head">67</th>
            <td class="metric-label">额外收入/支出</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('extraIncomeExpense')) }}</td>
            <td class="metric-label">残值增减</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('lineResidualChange')) }}</td>
            <td class="metric-label">应收增减</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('receivableChange')) }}</td>
            <td class="metric-label">净利润率</td>
            <td class="result-cell">{{ formatDisplayPercent(optionalDerivedMetric('netProfitRate')) }}</td>
            <td class="metric-label">毛率润率</td>
            <td class="result-cell">{{ formatDisplayPercent(optionalDerivedMetric('grossMarginRate')) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import { serviceOperatingLabels, type SandboxGameOperatingLabels } from '@/configs/sandbox-game-service-labels'
import {
  cloneOperatingPayload,
  type CellValue,
  type MarketCultivationItem,
  type NumericCellValue,
  type OperatingCarryForward,
  type OperatingPayload,
  type ProjectProgressItem,
  type QualificationItem,
  type QuarterValueMap,
} from '@/types/sandbox-game'
import { handleSequentialInputNavigation } from '@/utils/input-navigation'
import { hasFractionInput } from '@/utils/manual-integer'

type MarketBidKey =
  | 'basicProductTotal'
  | 'standardProductTotal'
  | 'precisionProductTotal'
  | 'intelligentProductTotal'
  | 'agencyInspectionTotal'
  | 'twoCabinVipTotal'
  | 'businessVipTotal'
  | 'memberCustomTotal'
type MarketBidRow = Record<MarketBidKey | 'orderAmount', NumericCellValue>
type ProjectFactoryKey = 'factoryA' | 'factoryB' | 'factoryC'
type ProjectProgressQuarterKey = 'q1' | 'q2' | 'q3' | 'q4'
type ProjectProgressCellKey = 'lineType' | 'progress'
type MarketCultivationKey = 'regional' | 'national' | 'global'
type QualificationKey = 'qualityEnvironmentalHealth' | 'highTechEnterprise' | 'specializedInnovation' | 'listedCompany'

const props = defineProps<{
  modelValue: OperatingPayload
  editableScopes: string[]
  invalidScopes: string[]
  quarterCashChecks: Record<string, number>
  currentStageCode: string
  derivedValues: Record<string, number>
  periodEndCash: number
  carryForward?: OperatingCarryForward | null
  labels?: SandboxGameOperatingLabels
}>()

const emit = defineEmits<{
  (event: 'update:modelValue', value: OperatingPayload): void
}>()

const marketRegions = [
  { key: 'local', label: '本地' },
  { key: 'regional', label: '区域' },
  { key: 'national', label: '全国' },
  { key: 'global', label: '全球' },
] as const

const labels = computed(() => props.labels ?? serviceOperatingLabels)

const marketProductFields = computed(() => labels.value.marketProductFields as ReadonlyArray<{ key: MarketBidKey; label: string }>)

const legacyMarketProductKeys = ['basicProductTotal', 'standardProductTotal', 'precisionProductTotal', 'intelligentProductTotal'] as const satisfies ReadonlyArray<MarketBidKey>

const quarterList = [
  { key: 'q1', label: '第一季度', scope: 'Q1' },
  { key: 'q2', label: '第二季度', scope: 'Q2' },
  { key: 'q3', label: '第三季度', scope: 'Q3' },
  { key: 'q4', label: '第四季度', scope: 'Q4' },
] as const
const extraPeriodList = [
  ...quarterList,
  { key: 'year_end', label: '年末', scope: 'YEAR_END' },
] as const

const projectFactoryRows = [
  { key: 'factoryA', label: '生产厂房 A' },
  { key: 'factoryB', label: '生产厂房 B' },
  { key: 'factoryC', label: '生产厂房 C' },
] as const satisfies ReadonlyArray<{ key: ProjectFactoryKey; label: string }>

const defaultProjectProgressRows: ProjectProgressItem[] = projectFactoryRows.flatMap((factory) =>
  quarterList.map((quarter) => ({
    projectName: `${factory.label}-${quarter.label}`,
    lineType: '',
    progress: '',
  })),
)
const projectLineTypeOptions = ['人工', '半自动', '自动', '智能'] as const

const marketCultivationFields = [
  { key: 'regional', label: '区域', threshold: 1 },
  { key: 'national', label: '全国', threshold: 2 },
  { key: 'global', label: '全球', threshold: 3 },
] as const satisfies ReadonlyArray<{ key: MarketCultivationKey; label: string; threshold: number }>

const qualificationFields = [
  { key: 'qualityEnvironmentalHealth', label: '质量、环境健康体系认证企业' },
  { key: 'highTechEnterprise', label: '高新技术企业' },
  { key: 'specializedInnovation', label: '专精特新小巨人' },
  { key: 'listedCompany', label: '上市企业' },
] as const satisfies ReadonlyArray<{ key: QualificationKey; label: string }>

const shortTermLoanFields = [
  { key: 'dueRepayment', label: '到期还贷' },
  { key: 'interest', label: '付利息' },
  { key: 'newLoan', label: '新增贷款' },
] as const

const materialFields = computed(() => labels.value.materialFields)
const supplyChainOrderFields = computed(() => labels.value.materialFields)

interface ProductionLineRow {
  key: string
  label: string
  labelColspan?: number
  groupLabel?: string
  groupRowspan?: number
  useExistingGroupCell?: boolean
}

const productionLineRows = computed<ReadonlyArray<ProductionLineRow>>(() => labels.value.productionLineRows)

const researchFields = computed(() => labels.value.researchFields)

const extraFields = [
  { key: 'discountExpense', label: '折现费用' },
  { key: 'extraExpensePenalty', label: '额外支出 / 罚款' },
  { key: 'extraIncomeReward', label: '额外收入 / 奖励' },
] as const
const materialFieldKeys = computed(() => materialFields.value.map((item) => item.key))

const marketBidRows = computed(() => buildFixedMarketBidRows(props.modelValue.beginning.marketBid))
const projectProgressRows = computed<ProjectProgressItem[]>(() => {
  const items = props.modelValue.yearEnd.projectProgressUpdate?.items ?? []
  return defaultProjectProgressRows.map((fallback, index) => ({
    ...fallback,
    ...(items[index] ?? {}),
  }))
})
const marketOrderTotal = computed(() => resolveMetric(props.modelValue.beginning.taxAndPlanning.orderTotal, props.derivedValues.orderTotal, marketBidRows.value.reduce((total, _row, index) => total + getMarketRowTotal(index), 0)))
const marketInvestmentTotal = computed(() => getStoredMarketInvestmentValue())
const marketBidReadonly = computed(() => isOrderLinkedMarketBid())
const taxPaymentDisplay = computed(() => resolveMetric(props.carryForward?.previousIncomeTax, props.modelValue.beginning.taxAndPlanning.taxPayment))
const planRevenueDisplay = computed(() => resolveMetric(props.modelValue.beginning.taxAndPlanning.planRevenue, props.derivedValues.orderTotal, marketOrderTotal.value))
const comprehensiveCostDisplay = computed(() => resolveMetric(props.modelValue.beginning.taxAndPlanning.comprehensiveCostPlan, props.derivedValues.comprehensiveCostTotal))
const shortTermLoanDelta = computed(() => getQuarterRowTotal('shortTermLoan', 'newLoan') - getQuarterRowTotal('shortTermLoan', 'dueRepayment'))
const longTermLoanDelta = computed(() => toNumber(getYearEndFieldValue('longTermLoan', 'newLoan')) - toNumber(getYearEndFieldValue('longTermLoan', 'repayment')))
const marketCultivationAnnualTotal = computed(() => getMarketCultivationAnnualAmount())

const cashInflow = computed(() =>
  derivedMetric('receivableRecovered') +
  derivedMetric('newShortTermLoan') +
  toNumber(getYearEndFieldValue('longTermLoan', 'newLoan')) +
  derivedMetric('lineSaleValue') +
  derivedMetric('factorySale') +
  derivedMetric('extraIncomeReward'),
)

const cashOutflow = computed(() =>
  taxPaymentDisplay.value +
  toNumber(marketInvestmentTotal.value) +
  derivedMetric('shortTermRepayment') +
  derivedMetric('shortTermInterest') +
  derivedMetric('materialPayment') +
  derivedMetric('changeProductCost') +
  derivedMetric('lineDismantleCost') +
  derivedMetric('newLineInstall') +
  derivedMetric('humanResourceCost') +
  derivedMetric('salaryAndProductionCost') +
  derivedMetric('researchCost') +
  derivedMetric('managementSystemCost') +
  derivedMetric('managementSalary') +
  toNumber(getYearEndFieldValue('assetAdjustment', 'lineMaintenance')) +
  toNumber(getYearEndFieldValue('assetAdjustment', 'purchase')) +
  toNumber(getYearEndFieldValue('longTermLoan', 'interest')) +
  toNumber(getYearEndFieldValue('longTermLoan', 'repayment')) +
  toNumber(getYearEndFieldValue('assetAdjustment', 'rent')) +
  marketCultivationAnnualTotal.value +
  derivedMetric('discountExpense') +
  derivedMetric('extraExpensePenalty'),
)

function visibleQuarterCount() {
  switch (props.currentStageCode) {
    case 'Q1':
      return 1
    case 'Q2':
      return 2
    case 'Q3':
      return 3
    case 'Q4':
    case 'YEAR_END':
      return 4
    default:
      return 4
  }
}

function displayQuarterCash(stageCode: 'Q1' | 'Q2' | 'Q3' | 'Q4') {
  const order = ['Q1', 'Q2', 'Q3', 'Q4']
  const index = order.indexOf(stageCode)
  if (index === -1 || index >= visibleQuarterCount()) {
    return ''
  }
  return formatNumber(props.quarterCashChecks[stageCode])
}

function derivedMetric(key: string) {
  return toNumber(props.derivedValues[key])
}

function optionalDerivedMetric(key: string) {
  const value = props.derivedValues[key]
  if (typeof value !== 'number' || Number.isNaN(value)) {
    return null
  }
  return value
}

function isScopeEditable(scope: string) {
  return props.editableScopes.includes(scope)
}

function isScopeInvalidDraft(scope: string) {
  return props.invalidScopes.includes(scope)
}

function editableCellClass(scope: string, value?: unknown) {
  if (isScopeEditable(scope)) {
    return hasFractionInput(value) ? 'input-cell manual-integer-invalid-cell' : 'input-cell'
  }
  if (isScopeInvalidDraft(scope)) {
    return 'invalid-draft-cell'
  }
  return 'locked-cell'
}

function isOrderSalesRevenueLinked() {
  return Number(props.modelValue.derived?.values?.orderSalesRevenueLinked ?? 0) === 1
}

function deliverySalesRevenueCellClass(scope: string, value?: unknown) {
  if (isOrderSalesRevenueLinked()) {
    return isScopeInvalidDraft(scope) ? 'linked-order-cell invalid-draft-cell' : 'linked-order-cell'
  }
  return editableCellClass(scope, value)
}

function displayLinkedSalesRevenue(quarterKey: string) {
  return formatNumber(getQuarterFieldValue('deliverySettlement', quarterKey, 'salesRevenue'))
}

function isAdminIssuedExtraField(fieldKey: string) {
  return fieldKey === 'extraExpensePenalty' || fieldKey === 'extraIncomeReward'
}

function extraCellClass(fieldKey: string, scope: string) {
  if (fieldKey === 'discountExpense' && scope === 'YEAR_END') {
    return 'not-applicable-cell'
  }
  if (isAdminIssuedExtraField(fieldKey)) {
    return 'admin-issued-cell'
  }
  return editableCellClass(scope, getQuarterFieldValue('incomeAndPenalty', scope.toLowerCase(), fieldKey))
}

function supplyChainOrderCellClass(scope: string, value?: unknown) {
  const baseClass = editableCellClass(scope, value)
  const parsed = parseNumber(value)
  if (value !== '' && value !== undefined && value !== null && (parsed === null || !Number.isInteger(parsed) || parsed < 0)) {
    return `${baseClass} manual-integer-invalid-cell`
  }
  return baseClass
}

function projectProgressIndex(factoryKey: ProjectFactoryKey, quarterKey: ProjectProgressQuarterKey) {
  const factoryIndex = projectFactoryRows.findIndex((item) => item.key === factoryKey)
  const quarterIndex = quarterList.findIndex((item) => item.key === quarterKey)
  if (factoryIndex === -1 || quarterIndex === -1) {
    return -1
  }
  return factoryIndex * quarterList.length + quarterIndex
}

function getProjectProgressCell(factoryKey: ProjectFactoryKey, quarterKey: ProjectProgressQuarterKey, key: ProjectProgressCellKey) {
  const index = projectProgressIndex(factoryKey, quarterKey)
  if (index < 0) {
    return ''
  }
  return projectProgressRows.value[index]?.[key] ?? ''
}

function projectProgressCellClass(scope: string, value?: unknown) {
  return editableCellClass(scope, value)
}

function updateProjectProgressCell(factoryKey: ProjectFactoryKey, quarterKey: ProjectProgressQuarterKey, key: ProjectProgressCellKey, event: Event) {
  const quarter = quarterList.find((item) => item.key === quarterKey)
  if (!quarter || !isScopeEditable(quarter.scope)) {
    return
  }
  const index = projectProgressIndex(factoryKey, quarterKey)
  if (index < 0) {
    return
  }
  const next = cloneOperatingPayload(props.modelValue)
  const sourceItems = next.yearEnd.projectProgressUpdate.items ?? []
  next.yearEnd.projectProgressUpdate.items = defaultProjectProgressRows.map((fallback, itemIndex) => ({
    ...fallback,
    ...(sourceItems[itemIndex] ?? {}),
  }))
  const target = next.yearEnd.projectProgressUpdate.items[index]
  if (!target) {
    return
  }
  const raw = (event.target as HTMLInputElement | HTMLSelectElement).value
  if (key === 'progress') {
    target.progress = normalizeNumericValue(raw)
  } else {
    target.lineType = raw
  }
  emit('update:modelValue', next)
}

function getMarketCultivationItem(key: MarketCultivationKey): MarketCultivationItem {
  return props.modelValue.yearEnd.marketCultivation?.[key] ?? createEmptyMarketCultivationItem()
}

function getMarketCultivationAnnualAmount() {
  const marketCultivation = props.modelValue.yearEnd.marketCultivation
  const items = marketCultivationFields.map((field) => marketCultivation?.[field.key] ?? createEmptyMarketCultivationItem())
  if (marketCultivation?.stateApplied && items.some(hasMarketCultivationItemSignal)) {
    return items.reduce((total, item) => total + getMarketCultivationEffectiveAnnual(item), 0)
  }
  if (items.some((item) => item.annualInvestment !== '' && item.annualInvestment !== undefined && item.annualInvestment !== null)) {
    return items.reduce((total, item) => total + getMarketCultivationEffectiveAnnual(item), 0)
  }
  return toNumber(getYearEndFieldValue('assetAdjustment', 'marketCultivation'))
}

function getMarketCultivationEffectiveAnnual(item: MarketCultivationItem) {
  if (item.lockedByPrevious) {
    return 0
  }
  if (item.annualInvestment !== '' && item.annualInvestment !== undefined && item.annualInvestment !== null) {
    return toNumber(item.annualInvestment)
  }
  return toNumber(item.effectiveAnnualInvestment)
}

function getMarketCultivationDisplayCumulative(key: MarketCultivationKey) {
  const item = getMarketCultivationItem(key)
  if (item.lockedByPrevious || item.locked || item.unlocked) {
    return toNumber(item.cumulativeInvestment)
  }
  if (item.annualInvestment !== '' && item.annualInvestment !== undefined && item.annualInvestment !== null) {
    return toNumber(item.previousCumulative) + getMarketCultivationEffectiveAnnual(item)
  }
  return toNumber(item.cumulativeInvestment)
}

function hasMarketCultivationItemSignal(item: MarketCultivationItem) {
  return (
    (item.annualInvestment !== '' && item.annualInvestment !== undefined && item.annualInvestment !== null) ||
    toNumber(item.effectiveAnnualInvestment) !== 0 ||
    toNumber(item.previousCumulative) !== 0 ||
    toNumber(item.cumulativeInvestment) !== 0 ||
    item.unlocked ||
    item.locked ||
    item.lockedByPrevious ||
    item.willUnlock
  )
}

function marketCultivationStatusText(key: MarketCultivationKey) {
  const item = getMarketCultivationItem(key)
  const field = marketCultivationFields.find((current) => current.key === key)
  if (item.lockedByPrevious) {
    return '已解锁（继承）'
  }
  if (item.locked || item.unlocked) {
    return '已解锁'
  }
  if (item.willUnlock || (getMarketCultivationEffectiveAnnual(item) > 0 && getMarketCultivationDisplayCumulative(key) >= (field?.threshold ?? 0))) {
    return '提交后解锁'
  }
  return `未解锁 / 门槛 ${field?.threshold ?? 0}M`
}

function marketCultivationInputClass(key: MarketCultivationKey) {
  const item = getMarketCultivationItem(key)
  const disabled = !isScopeEditable('YEAR_END') || item.locked || item.lockedByPrevious
  const invalid = isMarketCultivationAnnualInvalid(item.annualInvestment)
  const base = disabled ? 'market-cultivation-input locked-cell' : 'market-cultivation-input input-cell'
  return invalid ? `${base} manual-integer-invalid-cell` : base
}

function isMarketCultivationAnnualInvalid(value: unknown) {
  if (value === '' || value === undefined || value === null) {
    return false
  }
  const parsed = parseNumber(value)
  return parsed === null || !Number.isInteger(parsed) || (parsed !== 0 && parsed !== 1)
}

function updateMarketCultivationField(key: MarketCultivationKey, event: Event) {
  const next = cloneOperatingPayload(props.modelValue)
  const item = next.yearEnd.marketCultivation[key]
  if (!isScopeEditable('YEAR_END') || item.locked || item.lockedByPrevious) {
    return
  }
  item.annualInvestment = normalizeNumericValue((event.target as HTMLInputElement).value)
  emit('update:modelValue', next)
}

function getQualificationItem(key: QualificationKey): QualificationItem {
  return props.modelValue.yearEnd.qualificationCertification?.[key] ?? createEmptyQualificationItem()
}

function updateQualificationStatus(key: QualificationKey, event: Event) {
  const next = cloneOperatingPayload(props.modelValue)
  const item = next.yearEnd.qualificationCertification[key]
  if (!isScopeEditable('YEAR_END') || item.locked || item.lockedByPrevious) {
    return
  }
  const status = (event.target as HTMLSelectElement).value === '解锁' ? '解锁' : '未解锁'
  item.status = status
  item.unlocked = status === '解锁'
  emit('update:modelValue', next)
}

function qualificationCellClass(key: QualificationKey) {
  const item = getQualificationItem(key)
  return !isScopeEditable('YEAR_END') || item.locked || item.lockedByPrevious ? 'qualification-input-cell locked-cell' : 'qualification-input-cell input-cell'
}

function displayExtraReadonlyCell(quarterKey: string, fieldKey: string) {
  return formatNumber(getQuarterFieldValue('incomeAndPenalty', quarterKey, fieldKey))
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
  if (Number.isNaN(parsed)) {
    return ''
  }
  return Number.isInteger(parsed) ? parsed : value
}

function parseNumber(value: unknown): number | null {
  if (typeof value === 'number' && !Number.isNaN(value)) {
    return value
  }
  if (typeof value === 'string') {
    const trimmed = value.trim()
    if (trimmed === '') {
      return null
    }
    const parsed = Number(trimmed)
    return Number.isNaN(parsed) ? null : parsed
  }
  return null
}

function resolveMetric(...values: unknown[]) {
  for (const value of values) {
    const parsed = parseNumber(value)
    if (parsed !== null) {
      return parsed
    }
  }
  return 0
}

function toNumber(value: unknown) {
  const parsed = parseNumber(value)
  return parsed ?? 0
}

function formatNumber(value: unknown) {
  return toNumber(value).toLocaleString('zh-CN', { maximumFractionDigits: 2 })
}

function formatPercent(value: number) {
  return `${(value * 100).toFixed(2)}%`
}

function formatDisplayNumber(value: number | null) {
  if (value === null) {
    return '--'
  }
  return formatNumber(value)
}

function formatDisplayPercent(value: number | null) {
  if (value === null) {
    return '--'
  }
  return formatPercent(value)
}

function createMarketBidRow(): MarketBidRow {
  return {
    basicProductTotal: '',
    standardProductTotal: '',
    precisionProductTotal: '',
    intelligentProductTotal: '',
    agencyInspectionTotal: '',
    twoCabinVipTotal: '',
    businessVipTotal: '',
    memberCustomTotal: '',
    orderAmount: '',
  }
}

function createEmptyMarketCultivationItem(): MarketCultivationItem {
  return {
    annualInvestment: '',
    effectiveAnnualInvestment: 0,
    previousCumulative: 0,
    cumulativeInvestment: 0,
    unlocked: false,
    locked: false,
    lockedByPrevious: false,
    willUnlock: false,
  }
}

function createEmptyQualificationItem(): QualificationItem {
  return {
    status: '未解锁',
    unlocked: false,
    locked: false,
    lockedByPrevious: false,
  }
}

function buildFixedMarketBidRows(source: Array<Record<string, CellValue>>) {
  return marketRegions.map((_, index) => ({
    ...createMarketBidRow(),
    ...(source[index] ?? {}),
  })) as MarketBidRow[]
}

function ensureMarketBidRows(next: OperatingPayload) {
  const normalized = buildFixedMarketBidRows(next.beginning.marketBid)
  next.beginning.marketBid = normalized.map((item) => ({ ...item })) as Array<Record<string, CellValue>>
}

function getStoredMarketInvestmentValue(): NumericCellValue {
  const direct = props.modelValue.beginning.taxAndPlanning.marketInvestmentTotal
  if (direct !== undefined) {
    return direct as NumericCellValue
  }
  const fallback = props.modelValue.beginning.taxAndPlanning.marketBidCost
  if (fallback !== undefined) {
    return fallback as NumericCellValue
  }
  return ''
}

function isOrderLinkedMarketBid() {
  return props.modelValue.beginning.taxAndPlanning.orderLinked === true
}

function getMarketRowTotal(index: number) {
  const row = marketBidRows.value[index]
  const explicit = parseNumber(row.orderAmount)
  if (explicit !== null && isOrderLinkedMarketBid()) {
    return explicit
  }
  const currentTotal = marketProductFields.value.reduce((total, field) => total + toNumber(row[field.key]), 0)
  if (currentTotal !== 0) {
    return currentTotal
  }
  return legacyMarketProductKeys.reduce((total, key) => total + toNumber(row[key]), 0)
}

function syncMarketBidDerived(next: OperatingPayload) {
  ensureMarketBidRows(next)

  let orderTotal = 0
  next.beginning.marketBid.forEach((row) => {
    const typedRow = row as Record<string, CellValue>
    const rowTotal = marketProductFields.value.reduce((total, field) => total + toNumber(typedRow[field.key]), 0)
    typedRow.orderAmount = rowTotal
    orderTotal += rowTotal
  })

  const marketInvestment = normalizeNumericValue(String(next.beginning.taxAndPlanning.marketInvestmentTotal ?? ''))
  next.beginning.taxAndPlanning.marketInvestmentTotal = marketInvestment
  delete next.beginning.taxAndPlanning.marketBidCost
  next.beginning.taxAndPlanning.orderTotal = orderTotal
}

function updateMarketBidField(index: number, key: MarketBidKey, event: Event) {
  const next = cloneOperatingPayload(props.modelValue)
  ensureMarketBidRows(next)
  next.beginning.marketBid[index][key] = normalizeNumericValue((event.target as HTMLInputElement).value)
  syncMarketBidDerived(next)
  emit('update:modelValue', next)
}

function updateMarketInvestmentTotal(event: Event) {
  const next = cloneOperatingPayload(props.modelValue)
  next.beginning.taxAndPlanning.marketInvestmentTotal = normalizeNumericValue((event.target as HTMLInputElement).value)
  syncMarketBidDerived(next)
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
    case 'supplyChainOrderRecord':
      return props.modelValue.quarter.supplyChainOrderRecord
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
  return (entry?.[fieldKey] as NumericCellValue | undefined) ?? ''
}

function getQuarterRowTotal(source: string, fieldKey: string) {
  return quarterList.reduce((total, quarter) => total + toNumber(getQuarterFieldValue(source, quarter.key, fieldKey)), 0)
}

function getExtraRowTotal(fieldKey: string) {
  const quarterTotal = getQuarterRowTotal('incomeAndPenalty', fieldKey)
  if (!isAdminIssuedExtraField(fieldKey)) return quarterTotal
  return quarterTotal + toNumber(getQuarterFieldValue('incomeAndPenalty', 'year_end', fieldKey))
}

function getQuarterGroupQuarterTotal(source: string, fieldKeys: readonly string[], quarterKey: string) {
  return fieldKeys.reduce((total, fieldKey) => total + toNumber(getQuarterFieldValue(source, quarterKey, fieldKey)), 0)
}

function getQuarterGroupGrandTotal(source: string, fieldKeys: readonly string[]) {
  return quarterList.reduce((total, quarter) => total + getQuarterGroupQuarterTotal(source, fieldKeys, quarter.key), 0)
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
    case 'supplyChainOrderRecord':
      apply(next.quarter.supplyChainOrderRecord)
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
    return (props.modelValue.yearEnd.longTermLoan[fieldKey] as NumericCellValue | undefined) ?? ''
  }
  return (props.modelValue.yearEnd.assetAdjustment[fieldKey] as NumericCellValue | undefined) ?? ''
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
</script>

<style scoped>
.sheet-frame {
  border: 1px solid var(--line);
  border-radius: 18px;
  overflow: hidden;
  background: #ffffff;
  --operating-font-family: "Microsoft YaHei", "Segoe UI", Arial, sans-serif;
  --operating-section-title-font-size: 14px;
  --operating-group-title-font-size: 16px;
  --operating-body-font-size: 14px;
  --operating-cell-padding-y: 9px;
  --operating-cell-padding-x: 11px;
  --operating-title-padding-y: 10px;
  --operating-title-padding-x: 12px;
  --operating-input-min-height: 34px;
}

.sheet-scroll {
  overflow: auto;
  background: var(--sheet-bg);
}

.sheet-table {
  width: 100%;
  table-layout: fixed;
  border-collapse: collapse;
  background: #ffffff;
  color: #243447;
  font-family: var(--operating-font-family);
  font-size: var(--operating-body-font-size);
}

.col-index {
  width: 3.5%;
}

.col-band {
  width: 3.5%;
}

.col-main {
  width: 13%;
}

.col-sub {
  width: 13%;
}

.col-stage {
  width: 9.5%;
}

.col-total {
  width: 8%;
}

.col-side {
  width: 7%;
}

.sheet-table th,
.sheet-table td {
  border: 1px solid var(--line);
  padding: 0;
  vertical-align: middle;
  box-sizing: border-box;
  word-break: break-word;
}

.corner,
.col-head,
.row-head {
  background: #f1f4f8;
  color: #5d6878;
  text-align: center;
  font-weight: 700;
}

.col-head,
.corner {
  height: 32px;
  font-size: 12px;
}

.row-head {
  font-size: 12px;
}

.section-band {
  writing-mode: vertical-rl;
  text-orientation: upright;
  text-align: center;
  padding: 10px 6px;
  letter-spacing: 1.8px;
  font-weight: 700;
  font-size: var(--operating-section-title-font-size);
  line-height: 1.25;
}

.band-year-start {
  background: var(--year-start-bg);
}

.band-quarter {
  background: var(--quarter-bg);
}

.band-year-end {
  background: var(--year-end-bg);
}

.band-misc {
  background: #f7ead4;
}

.group-title,
.result-title,
.note-cell,
.metric-label,
.dual-metric-cell,
.empty-cell,
.nested-group-cell,
.reminder-cell {
  padding: var(--operating-cell-padding-y) var(--operating-cell-padding-x);
  font-size: var(--operating-body-font-size);
  line-height: 1.45;
}

.group-title,
.result-title {
  background: #fbfbfc;
  font-weight: 700;
  text-align: center;
}

.group-title {
  height: 40px;
  padding: var(--operating-title-padding-y) var(--operating-title-padding-x);
  font-size: var(--operating-group-title-font-size);
  line-height: 1.45;
}

.result-title {
  height: 36px;
  font-size: var(--operating-body-font-size);
  line-height: 1.45;
}

.note-cell,
.metric-label {
  background: #f8fafc;
}

.note-cell {
  height: 36px;
  padding: var(--operating-cell-padding-y) var(--operating-cell-padding-x);
  text-align: center;
  font-size: var(--operating-body-font-size);
  font-weight: 500;
  line-height: 1.5;
}

.item-label {
  text-align: center;
  font-weight: 500;
}

.period-label {
  height: 36px;
  text-align: center;
  font-weight: 700;
  line-height: 1.4;
}

.instruction-note {
  height: 40px;
  padding: 10px 12px;
  text-align: left;
  font-weight: 600;
  line-height: 1.55;
}

.rule-note {
  padding: 10px 12px;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.rule-note-short {
  text-align: center;
  line-height: 1.5;
}

.rule-note-long {
  text-align: left;
  line-height: 1.55;
}

.nested-group-cell {
  background: #f3f6fa;
  text-align: center;
  font-weight: 700;
  font-size: var(--operating-body-font-size);
  line-height: 1.45;
  color: #314458;
}

.reminder-cell {
  background: linear-gradient(90deg, #fff1b3 0%, #ffe39a 100%);
  border-left: 4px solid #df7c0a;
  color: #8f3f00;
  font-size: 15px;
  font-weight: 800;
  line-height: 1.5;
  padding: 10px 12px;
  letter-spacing: 0.02em;
}

.excel-nested-cell {
  background: #ffffff;
  padding: 0;
}

.nested-excel-table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

.nested-excel-table td {
  height: 32px;
  border: 1px solid var(--line);
  padding: 0;
  box-sizing: border-box;
  color: #243447;
  font-size: 13px;
  line-height: 1.35;
  text-align: center;
  vertical-align: middle;
}

.nested-title-cell {
  background: #fff1b3;
  color: #8f3f00;
  font-weight: 800;
}

.nested-title-cell span {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  font-weight: 700;
  line-height: 1.35;
}

.nested-head-cell,
.nested-subhead-cell,
.qualification-label-cell {
  background: #f3f6fa;
  color: #314458;
  font-weight: 800;
}

.nested-subhead-cell {
  font-size: 12px;
}

.nested-result-cell,
.nested-status-cell {
  background: var(--calc-bg);
  color: #1f4f82;
  font-weight: 800;
  font-variant-numeric: tabular-nums;
}

.project-progress-title-cell {
  width: 25%;
  padding: 8px 10px !important;
  line-height: 1.45 !important;
}

.project-factory-head,
.project-factory-subhead,
.project-factory-cell {
  width: 11%;
}

.project-quarter-head {
  width: 16%;
}

.market-rule-cell {
  width: 25%;
  padding: 8px 10px !important;
  line-height: 1.45 !important;
}

.market-name-cell {
  width: 9%;
}

.market-total-cell {
  width: 13%;
  line-height: 1.5 !important;
}

.market-cultivation-input,
.qualification-input-cell {
  text-align: center;
  font-variant-numeric: tabular-nums;
}

.qualification-label-cell {
  padding: 8px 10px;
  text-align: center;
  font-size: 12px;
  line-height: 1.35;
}

.nested-excel-table input,
.nested-excel-table select,
.qualification-input-cell select {
  width: 100%;
  min-height: 32px;
  border: none;
  background: transparent;
  padding: 5px 6px;
  box-sizing: border-box;
  color: inherit;
  font: inherit;
  font-size: 13px;
  font-weight: 600;
  line-height: 1.35;
  text-align: center;
  outline: none;
}

.nested-excel-table input:focus,
.nested-excel-table select:focus,
.qualification-input-cell select:focus {
  box-shadow: inset 0 0 0 2px #2563eb, 0 0 0 2px rgba(37, 99, 235, 0.14);
}

.nested-excel-table input:disabled,
.nested-excel-table select:disabled,
.qualification-input-cell select:disabled {
  color: #8491a6;
  cursor: not-allowed;
}

.dual-metric-cell {
  background: #f8fafc;
  padding: 8px 10px;
}

.metric-label {
  text-align: center;
  font-weight: 700;
  font-size: var(--operating-body-font-size);
  line-height: 1.45;
}

.metric-stack-line {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  align-items: center;
  font-size: 13px;
  color: #33475b;
  font-variant-numeric: tabular-nums;
}

.metric-stack-line + .metric-stack-line {
  margin-top: 6px;
  padding-top: 6px;
  border-top: 1px dashed #d6dfea;
}

.metric-stack-line strong {
  color: #173753;
}

.empty-cell {
  background: #ffffff;
}

.input-cell,
.locked-cell,
.invalid-draft-cell,
.result-cell,
.orange-cell,
.quarter-cash-cell {
  text-align: center;
  font-variant-numeric: tabular-nums;
}

.input-cell,
.locked-cell,
.invalid-draft-cell {
  background: var(--input-bg);
}

.locked-cell {
  background: var(--readonly-bg);
}

.invalid-draft-cell {
  background: #fff1dc;
}

.linked-order-cell {
  background: var(--calc-bg);
}

.input-cell:focus-within,
.manual-integer-invalid-cell:focus-within {
  position: relative;
  box-shadow: inset 0 0 0 2px #2563eb, 0 0 0 2px rgba(37, 99, 235, 0.14);
}

.input-cell input,
.locked-cell input,
.invalid-draft-cell input,
.linked-order-cell input {
  width: 100%;
  border: none;
  background: transparent;
  padding: 8px 10px;
  outline: none;
  min-height: var(--operating-input-min-height);
  box-sizing: border-box;
  font: inherit;
  font-size: var(--operating-body-font-size);
  font-weight: 500;
  line-height: 1.4;
  font-variant-numeric: tabular-nums;
  color: inherit;
  text-align: center;
}

.admin-issued-cell {
  background: #eef3fb;
  color: #526175;
  font-weight: 700;
  padding: 8px 10px;
  font-size: var(--operating-body-font-size);
  line-height: 1.4;
  text-align: center;
  font-variant-numeric: tabular-nums;
}

.not-applicable-cell {
  background: #f4f6f8;
  color: #8491a6;
  font-weight: 700;
  padding: 8px 10px;
  font-size: var(--operating-body-font-size);
  line-height: 1.4;
  text-align: center;
  font-variant-numeric: tabular-nums;
}

.cell-readonly-value {
  display: block;
  min-height: var(--operating-input-min-height);
  line-height: var(--operating-input-min-height);
  font-variant-numeric: tabular-nums;
}
.locked-cell input {
  color: #8491a6;
}

.invalid-draft-cell input {
  color: #8a5a17;
}

.linked-order-cell input {
  color: inherit;
}

.system-linked-value {
  display: block;
  min-height: var(--operating-input-min-height);
  line-height: var(--operating-input-min-height);
  padding: 8px 10px;
  box-sizing: border-box;
  color: #1f4f82;
  font-size: var(--operating-body-font-size);
  font-weight: 800;
  font-variant-numeric: tabular-nums;
  text-align: center;
}

.manual-integer-invalid-cell {
  background: #fff5f5;
}

.manual-integer-invalid-cell input {
  color: var(--danger);
}

.result-cell,
.quarter-cash-cell {
  background: var(--calc-bg);
  padding: 8px 10px;
  font-size: var(--operating-body-font-size);
  font-weight: 700;
  line-height: 1.4;
  font-variant-numeric: tabular-nums;
}

.quarter-cash-cell {
  background: #f8c15d;
}

.orange-cell {
  background: #f8c15d;
  padding: 8px 10px;
  font-size: var(--operating-body-font-size);
  font-weight: 700;
  line-height: 1.4;
  font-variant-numeric: tabular-nums;
}

@media (max-width: 1400px) {
  .sheet-table {
    min-width: 1180px;
  }
}
</style>










