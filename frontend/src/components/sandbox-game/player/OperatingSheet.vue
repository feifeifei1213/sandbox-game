<template>
  <div class="sheet-frame">
    <div class="sheet-scroll">
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
            <td class="note-cell center" colspan="4">一般企业 25%，高新企业 15%</td>
            <td class="result-cell">{{ formatNumber(taxPaymentDisplay) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">2</th>
            <td class="group-title" colspan="2">2. 召开年度经营会议</td>
            <td class="note-cell center" colspan="3">调整公司战略，制定当年经营和工作计划</td>
            <td class="note-cell center">计划收入</td>
            <td class="result-cell">{{ formatNumber(planRevenueDisplay) }}</td>
            <td class="note-cell center">综合费用</td>
            <td class="result-cell">{{ formatNumber(comprehensiveCostDisplay) }}</td>
          </tr>
          <tr>
            <th class="row-head">3</th>
            <td class="group-title center" rowspan="5">3. 市场竞标</td>
            <td class="note-cell center">区域</td>
            <td v-for="field in marketProductFields" :key="`market-head-${field.key}`" class="note-cell center">{{ field.label }}</td>
            <td class="note-cell center">区域市场订单总价</td>
            <td class="note-cell center">订单总额</td>
            <td class="note-cell center">市场投入</td>
            <td class="empty-cell"></td>
          </tr>
          <tr v-for="(region, index) in marketRegions" :key="region.key">
            <th class="row-head">{{ 4 + index }}</th>
            <td class="note-cell center">{{ region.label }}</td>
            <td v-for="field in marketProductFields" :key="`${region.key}-${field.key}`" :class="editableCellClass('YEAR_START')">
              <input
                :value="displayCell(marketBidRows[index][field.key])"
                :disabled="!isScopeEditable('YEAR_START')"
                inputmode="decimal"
                @input="updateMarketBidField(index, field.key, $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getMarketRowTotal(index)) }}</td>
            <td v-if="index === 0" class="result-cell" :rowspan="marketRegions.length">{{ formatNumber(marketOrderTotal) }}</td>
            <td v-if="index === 0" :class="editableCellClass('YEAR_START')" :rowspan="marketRegions.length">
              <input
                :value="displayCell(marketInvestmentTotal)"
                :disabled="!isScopeEditable('YEAR_START')"
                inputmode="decimal"
                @input="updateMarketInvestmentTotal($event)"
              />
            </td>
            <td v-if="index === 0" class="empty-cell" :rowspan="marketRegions.length"></td>
          </tr>

          <tr>
            <th class="row-head">8</th>
            <td class="section-band band-quarter" rowspan="32">每季度工作</td>
            <td class="group-title" rowspan="4">1. 短期贷款更新账期</td>
            <td class="note-cell center"></td>
            <td v-for="quarter in quarterList" :key="`loan-head-${quarter.key}`" class="note-cell center">{{ quarter.label }}</td>
            <td class="note-cell center">总计</td>
            <td class="note-cell center">增减</td>
            <td class="empty-cell"></td>
          </tr>
          <tr v-for="(field, index) in shortTermLoanFields" :key="`loan-${field.key}`">
            <th class="row-head">{{ 9 + index }}</th>
            <td class="note-cell">{{ field.label }}</td>
            <td v-for="quarter in quarterList" :key="`loan-${field.key}-${quarter.key}`" :class="editableCellClass(quarter.scope)">
              <input
                :value="displayCell(getQuarterFieldValue('shortTermLoan', quarter.key, field.key))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="decimal"
                @input="updateQuarterField('shortTermLoan', quarter.key, field.key, $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('shortTermLoan', field.key)) }}</td>
            <td v-if="index === 0" class="orange-cell" :rowspan="shortTermLoanFields.length">{{ formatNumber(shortTermLoanDelta) }}</td>
            <td v-if="index === 0" class="empty-cell" :rowspan="shortTermLoanFields.length"></td>
          </tr>

          <tr>
            <th class="row-head">12</th>
            <td class="group-title" rowspan="6">2. 支付材料费用</td>
            <td class="note-cell center"></td>
            <td v-for="quarter in quarterList" :key="`material-head-${quarter.key}`" class="note-cell center">{{ quarter.label }}</td>
            <td class="note-cell center">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr v-for="(field, index) in materialFields" :key="`material-${field.key}`">
            <th class="row-head">{{ 13 + index }}</th>
            <td class="note-cell">{{ field.label }}</td>
            <td v-for="quarter in quarterList" :key="`material-${field.key}-${quarter.key}`" :class="editableCellClass(quarter.scope)">
              <input
                :value="displayCell(getQuarterFieldValue('materialPayment', quarter.key, field.key))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="decimal"
                @input="updateQuarterField('materialPayment', quarter.key, field.key, $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('materialPayment', field.key)) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">17</th>
            <td class="result-title">合计</td>
            <td v-for="quarter in quarterList" :key="`material-total-${quarter.key}`" class="result-cell">{{ formatNumber(getQuarterGroupQuarterTotal('materialPayment', materialFieldKeys, quarter.key)) }}</td>
            <td class="result-cell">{{ formatNumber(getQuarterGroupGrandTotal('materialPayment', materialFieldKeys)) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">18</th>
            <td class="group-title" rowspan="7">3. 生产线调整</td>
            <td class="note-cell center"></td>
            <td v-for="quarter in quarterList" :key="`line-head-${quarter.key}`" class="note-cell center">{{ quarter.label }}</td>
            <td class="note-cell center">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr v-for="(field, index) in productionLineFields" :key="`line-${field.key}`">
            <th class="row-head">{{ 19 + index }}</th>
            <td class="note-cell">{{ field.label }}</td>
            <td v-for="quarter in quarterList" :key="`line-${field.key}-${quarter.key}`" :class="editableCellClass(quarter.scope)">
              <input
                :value="displayCell(getQuarterFieldValue('productionLineAdjustment', quarter.key, field.key))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="decimal"
                @input="updateQuarterField('productionLineAdjustment', quarter.key, field.key, $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('productionLineAdjustment', field.key)) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">25</th>
            <td class="group-title" rowspan="2">4. 人力资源</td>
            <td class="note-cell center"></td>
            <td v-for="quarter in quarterList" :key="`hr-head-${quarter.key}`" class="note-cell center">{{ quarter.label }}</td>
            <td class="note-cell center">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">26</th>
            <td class="note-cell">人力资源费用</td>
            <td v-for="quarter in quarterList" :key="`hr-${quarter.key}`" :class="editableCellClass(quarter.scope)">
              <input
                :value="displayCell(getQuarterFieldValue('humanResource', quarter.key, 'staffCost'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="decimal"
                @input="updateQuarterField('humanResource', quarter.key, 'staffCost', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('humanResource', 'staffCost')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">27</th>
            <td class="group-title" rowspan="2">5. 工资与生产</td>
            <td class="note-cell center"></td>
            <td v-for="quarter in quarterList" :key="`salary-head-${quarter.key}`" class="note-cell center">{{ quarter.label }}</td>
            <td class="note-cell center">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">28</th>
            <td class="note-cell">工资与生产费用</td>
            <td v-for="quarter in quarterList" :key="`salary-${quarter.key}`" :class="editableCellClass(quarter.scope)">
              <input
                :value="displayCell(getQuarterFieldValue('salaryAndProduction', quarter.key, 'salaryCost'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="decimal"
                @input="updateQuarterField('salaryAndProduction', quarter.key, 'salaryCost', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('salaryAndProduction', 'salaryCost')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">29</th>
            <td class="group-title" rowspan="3">6. 研发与管理</td>
            <td class="note-cell center"></td>
            <td v-for="quarter in quarterList" :key="`research-head-${quarter.key}`" class="note-cell center">{{ quarter.label }}</td>
            <td class="note-cell center">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr v-for="(field, index) in researchFields" :key="`research-${field.key}`">
            <th class="row-head">{{ 30 + index }}</th>
            <td class="note-cell">{{ field.label }}</td>
            <td v-for="quarter in quarterList" :key="`research-${field.key}-${quarter.key}`" :class="editableCellClass(quarter.scope)">
              <input
                :value="displayCell(getQuarterFieldValue('researchAndManagement', quarter.key, field.key))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="decimal"
                @input="updateQuarterField('researchAndManagement', quarter.key, field.key, $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('researchAndManagement', field.key)) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">32</th>
            <td class="group-title" rowspan="2">7. 应收更新</td>
            <td class="note-cell center"></td>
            <td v-for="quarter in quarterList" :key="`receivable-head-${quarter.key}`" class="note-cell center">{{ quarter.label }}</td>
            <td class="note-cell center">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">33</th>
            <td class="note-cell">应收回款</td>
            <td v-for="quarter in quarterList" :key="`receivable-${quarter.key}`" :class="editableCellClass(quarter.scope)">
              <input
                :value="displayCell(getQuarterFieldValue('receivableUpdate', quarter.key, 'receivableCollection'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="decimal"
                @input="updateQuarterField('receivableUpdate', quarter.key, 'receivableCollection', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('receivableUpdate', 'receivableCollection')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">34</th>
            <td class="group-title" rowspan="3">8. 产品验收交付</td>
            <td class="note-cell center"></td>
            <td v-for="quarter in quarterList" :key="`delivery-head-${quarter.key}`" class="note-cell center">{{ quarter.label }}</td>
            <td class="note-cell center">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">35</th>
            <td class="note-cell">销售收入</td>
            <td v-for="quarter in quarterList" :key="`sales-${quarter.key}`" :class="editableCellClass(quarter.scope)">
              <input
                :value="displayCell(getQuarterFieldValue('deliverySettlement', quarter.key, 'salesRevenue'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="decimal"
                @input="updateQuarterField('deliverySettlement', quarter.key, 'salesRevenue', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('deliverySettlement', 'salesRevenue')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">36</th>
            <td class="note-cell">直接成本</td>
            <td v-for="quarter in quarterList" :key="`cost-${quarter.key}`" :class="editableCellClass(quarter.scope)">
              <input
                :value="displayCell(getQuarterFieldValue('deliverySettlement', quarter.key, 'directCost'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="decimal"
                @input="updateQuarterField('deliverySettlement', quarter.key, 'directCost', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('deliverySettlement', 'directCost')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">37</th>
            <td class="group-title" rowspan="2">9. 支付管理人员费用</td>
            <td class="note-cell center">每季 1M</td>
            <td v-for="quarter in quarterList" :key="`management-head-${quarter.key}`" class="note-cell center">{{ quarter.label }}</td>
            <td class="note-cell center">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">38</th>
            <td class="note-cell">管理人员费用</td>
            <td v-for="quarter in quarterList" :key="`management-${quarter.key}`" :class="editableCellClass(quarter.scope)">
              <input
                :value="displayCell(getQuarterFieldValue('deliverySettlement', quarter.key, 'managementStaffCost'))"
                :disabled="!isScopeEditable(quarter.scope)"
                inputmode="decimal"
                @input="updateQuarterField('deliverySettlement', quarter.key, 'managementStaffCost', $event)"
              />
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('deliverySettlement', 'managementStaffCost')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">41</th>
            <td class="group-title" colspan="2">核对季末现金</td>
            <td class="quarter-cash-cell">{{ displayQuarterCash('Q1') }}</td>
            <td class="quarter-cash-cell">{{ displayQuarterCash('Q2') }}</td>
            <td class="quarter-cash-cell">{{ displayQuarterCash('Q3') }}</td>
            <td class="quarter-cash-cell">{{ displayQuarterCash('Q4') }}</td>
            <td class="result-cell">{{ formatNumber(periodEndCash) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">42</th>
            <td class="section-band band-year-end" rowspan="12">年末工作</td>
            <td class="group-title" rowspan="3">1. 办理长期贷款账期更新</td>
            <td class="note-cell">付利息</td>
            <td class="note-cell center" colspan="5">年利率 5%</td>
            <td :class="editableCellClass('YEAR_END')">
              <input
                :value="displayCell(getYearEndFieldValue('longTermLoan', 'interest'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="decimal"
                @input="updateYearEndField('longTermLoan', 'interest', $event)"
              />
            </td>
            <td class="orange-cell" rowspan="3">{{ formatNumber(longTermLoanDelta) }}</td>
            <td class="empty-cell" rowspan="3"></td>
          </tr>
          <tr>
            <th class="row-head">43</th>
            <td class="note-cell">到期还款</td>
            <td class="note-cell center" colspan="5">------------------------------</td>
            <td :class="editableCellClass('YEAR_END')">
              <input
                :value="displayCell(getYearEndFieldValue('longTermLoan', 'repayment'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="decimal"
                @input="updateYearEndField('longTermLoan', 'repayment', $event)"
              />
            </td>
          </tr>
          <tr>
            <th class="row-head">44</th>
            <td class="note-cell">办理新贷款</td>
            <td class="note-cell center" colspan="5">------------------------------</td>
            <td :class="editableCellClass('YEAR_END')">
              <input
                :value="displayCell(getYearEndFieldValue('longTermLoan', 'newLoan'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="decimal"
                @input="updateYearEndField('longTermLoan', 'newLoan', $event)"
              />
            </td>
          </tr>

          <tr>
            <th class="row-head">45</th>
            <td class="group-title">2. 支付生产线年度维护费</td>
            <td class="note-cell center">1M / 条</td>
            <td class="empty-cell" colspan="5"></td>
            <td :class="editableCellClass('YEAR_END')">
              <input
                :value="displayCell(getYearEndFieldValue('assetAdjustment', 'lineMaintenance'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="decimal"
                @input="updateYearEndField('assetAdjustment', 'lineMaintenance', $event)"
              />
            </td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">46</th>
            <td class="group-title" rowspan="2">3. 数据中心资产</td>
            <td class="note-cell center">购买</td>
            <td class="note-cell center" colspan="5">数据中心 A / B / C 价位（40M、32M、16M）</td>
            <td :class="editableCellClass('YEAR_END')">
              <input
                :value="displayCell(getYearEndFieldValue('assetAdjustment', 'purchase'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="decimal"
                @input="updateYearEndField('assetAdjustment', 'purchase', $event)"
              />
            </td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">47</th>
            <td class="note-cell center">出售</td>
            <td class="note-cell center" colspan="5">数据中心 A / B / C 价位（40M、32M、16M）</td>
            <td :class="editableCellClass('YEAR_END')">
              <input
                :value="displayCell(getYearEndFieldValue('assetAdjustment', 'sale'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="decimal"
                @input="updateYearEndField('assetAdjustment', 'sale', $event)"
              />
            </td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">48</th>
            <td class="group-title">4. 数据中心租金</td>
            <td class="note-cell center">付租金</td>
            <td class="note-cell center" colspan="5">数据中心 A / B / C 租金（4M、3M、2M）</td>
            <td :class="editableCellClass('YEAR_END')">
              <input
                :value="displayCell(getYearEndFieldValue('assetAdjustment', 'rent'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="decimal"
                @input="updateYearEndField('assetAdjustment', 'rent', $event)"
              />
            </td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">49</th>
            <td class="group-title">5. 生产线残值</td>
            <td class="empty-cell" colspan="6"></td>
            <td class="result-cell">{{ formatNumber(derivedMetric('lineResidual')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">50</th>
            <td class="group-title" rowspan="3">6. 生产线折旧</td>
            <td class="note-cell">折旧前待折资产总价值</td>
            <td class="empty-cell" colspan="5"></td>
            <td class="result-cell">{{ formatNumber(derivedMetric('depreciableAssetTotal')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">51</th>
            <td class="note-cell">折旧费</td>
            <td class="note-cell center" colspan="5">按待折资产的 1 / 3 取整</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('depreciation')) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr>
            <th class="row-head">52</th>
            <td class="note-cell">未完工生产线价值</td>
            <td class="empty-cell" colspan="5"></td>
            <td :class="editableCellClass('YEAR_END')">
              <input
                :value="displayCell(getYearEndFieldValue('assetAdjustment', 'workInConstruction'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="decimal"
                @input="updateYearEndField('assetAdjustment', 'workInConstruction', $event)"
              />
            </td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">53</th>
            <td class="group-title">7. 新市场培育</td>
            <td class="empty-cell"></td>
            <td class="note-cell center" colspan="5">每年可向区域、全国、全球各投 1M</td>
            <td :class="editableCellClass('YEAR_END')">
              <input
                :value="displayCell(getYearEndFieldValue('assetAdjustment', 'marketCultivation'))"
                :disabled="!isScopeEditable('YEAR_END')"
                inputmode="decimal"
                @input="updateYearEndField('assetAdjustment', 'marketCultivation', $event)"
              />
            </td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">54</th>
            <td class="section-band band-misc" rowspan="4">其他收支</td>
            <td class="group-title" rowspan="4">额外收入 / 罚款</td>
            <td class="note-cell center"></td>
            <td v-for="quarter in quarterList" :key="`extra-head-${quarter.key}`" class="note-cell center">{{ quarter.label }}</td>
            <td class="note-cell center">总计</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>
          <tr v-for="(field, index) in extraFields" :key="`extra-${field.key}`">
            <th class="row-head">{{ 55 + index }}</th>
            <td class="note-cell">{{ field.label }}</td>
            <td v-for="quarter in quarterList" :key="`extra-${field.key}-${quarter.key}`" :class="extraCellClass(field.key, quarter.scope)">
              <template v-if="isAdminIssuedExtraField(field.key)">
                <span class="cell-readonly-value">{{ displayExtraReadonlyCell(quarter.key, field.key) }}</span>
              </template>
              <template v-else>
                <input
                  :value="displayCell(getQuarterFieldValue('incomeAndPenalty', quarter.key, field.key))"
                  :disabled="!isScopeEditable(quarter.scope)"
                  inputmode="decimal"
                  @input="updateQuarterField('incomeAndPenalty', quarter.key, field.key, $event)"
                />
              </template>
            </td>
            <td class="result-cell">{{ formatNumber(getQuarterRowTotal('incomeAndPenalty', field.key)) }}</td>
            <td class="empty-cell" colspan="2"></td>
          </tr>

          <tr>
            <th class="row-head">59</th>
            <td class="metric-label" colspan="2">期末现金</td>
            <td class="result-cell">{{ formatNumber(periodEndCash) }}</td>
            <td class="metric-label">现金收入</td>
            <td class="result-cell">{{ formatNumber(cashInflow) }}</td>
            <td class="metric-label">现金支出</td>
            <td class="result-cell">{{ formatNumber(cashOutflow) }}</td>
            <td class="metric-label">市场回报比</td>
            <td class="result-cell">{{ formatPercent(marketReturnRatio) }}</td>
            <td class="result-cell">{{ formatNumber(laborProductivity) }}</td>
          </tr>
          <tr>
            <th class="row-head">60</th>
            <td class="metric-label" colspan="2">财务收入 / 支出</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('financeIncomeExpense')) }}</td>
            <td class="metric-label">不动产增减值</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('factoryAssetChange')) }}</td>
            <td class="metric-label">待折资产增减</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('depreciableAssetChange')) }}</td>
            <td class="metric-label">净资产收益</td>
            <td class="result-cell">{{ formatPercent(netAssetYield) }}</td>
            <td class="result-cell">{{ formatPercent(totalAssetYield) }}</td>
          </tr>
          <tr>
            <th class="row-head">61</th>
            <td class="metric-label" colspan="2">额外收入 / 支出</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('extraIncomeExpense')) }}</td>
            <td class="metric-label">残值增减</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('lineResidualChange')) }}</td>
            <td class="metric-label">应收增减</td>
            <td class="result-cell">{{ formatNumber(derivedMetric('receivableChange')) }}</td>
            <td class="metric-label">净利润率</td>
            <td class="result-cell">{{ formatPercent(netProfitRate) }}</td>
            <td class="result-cell">{{ formatPercent(grossMarginRate) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import { cloneOperatingPayload, type NumericCellValue, type OperatingPayload, type QuarterValueMap } from '@/types/sandbox-game'

type MarketBidKey = 'basicProductTotal' | 'standardProductTotal' | 'precisionProductTotal' | 'intelligentProductTotal'
type MarketBidRow = Record<MarketBidKey | 'orderAmount', NumericCellValue>

const props = defineProps<{
  modelValue: OperatingPayload
  editableScopes: string[]
  invalidScopes: string[]
  quarterCashChecks: Record<string, number>
  currentStageCode: string
  derivedValues: Record<string, number>
  periodEndCash: number
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

const marketProductFields = [
  { key: 'basicProductTotal', label: '基础产品总价' },
  { key: 'standardProductTotal', label: '标准产品总价' },
  { key: 'precisionProductTotal', label: '精密产品总价' },
  { key: 'intelligentProductTotal', label: '智能产品总价' },
] as const satisfies ReadonlyArray<{ key: MarketBidKey; label: string }>

const quarterList = [
  { key: 'q1', label: '第一季度', scope: 'Q1' },
  { key: 'q2', label: '第二季度', scope: 'Q2' },
  { key: 'q3', label: '第三季度', scope: 'Q3' },
  { key: 'q4', label: '第四季度', scope: 'Q4' },
] as const

const shortTermLoanFields = [
  { key: 'dueRepayment', label: '到期还贷' },
  { key: 'interest', label: '付利息' },
  { key: 'newLoan', label: '新增贷款' },
] as const

const materialFields = [
  { key: 'basicProduct', label: '基础产品' },
  { key: 'standardProduct', label: '标准产品' },
  { key: 'precisionProduct', label: '精密产品' },
  { key: 'intelligentProduct', label: '智能产品' },
] as const

const productionLineFields = [
  { key: 'changeProduct', label: '变更产品' },
  { key: 'dismantleCost', label: '生产线拆除 1M / 条' },
  { key: 'lineSale', label: '生产线出售' },
  { key: 'newLineInstall', label: '新生产线安装' },
  { key: 'constructionToFixed', label: '转入固定资产' },
  { key: 'newDepreciableAsset', label: '新增待折资产' },
] as const

const researchFields = [
  { key: 'technologyResearch', label: '技术研发' },
  { key: 'managementSystem', label: '管理体系' },
] as const

const extraFields = [
  { key: 'discountExpense', label: '折现费用' },
  { key: 'extraExpensePenalty', label: '额外支出 / 罚款' },
  { key: 'extraIncomeReward', label: '额外收入 / 奖励' },
] as const
const materialFieldKeys = materialFields.map((item) => item.key)

const marketBidRows = computed(() => buildFixedMarketBidRows(props.modelValue.beginning.marketBid))
const marketOrderTotal = computed(() => marketBidRows.value.reduce((total, _row, index) => total + getMarketRowTotal(index), 0))
const marketInvestmentTotal = computed(() => getStoredMarketInvestmentValue())
const taxPaymentDisplay = computed(() => resolveMetric(props.modelValue.beginning.taxAndPlanning.taxPayment))
const planRevenueDisplay = computed(() => resolveMetric(props.modelValue.beginning.taxAndPlanning.planRevenue, props.derivedValues.orderTotal, marketOrderTotal.value))
const comprehensiveCostDisplay = computed(() => resolveMetric(props.modelValue.beginning.taxAndPlanning.comprehensiveCostPlan, props.derivedValues.comprehensiveCostTotal))
const shortTermLoanDelta = computed(() => getQuarterRowTotal('shortTermLoan', 'newLoan') - getQuarterRowTotal('shortTermLoan', 'dueRepayment'))
const longTermLoanDelta = computed(() => toNumber(getYearEndFieldValue('longTermLoan', 'newLoan')) - toNumber(getYearEndFieldValue('longTermLoan', 'repayment')))

const cashInflow = computed(() =>
  derivedMetric('salesRevenue') +
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
  toNumber(getYearEndFieldValue('assetAdjustment', 'marketCultivation')) +
  derivedMetric('discountExpense') +
  derivedMetric('extraExpensePenalty'),
)

const marketReturnRatio = computed(() => ratio(marketOrderTotal.value, toNumber(marketInvestmentTotal.value)))
const laborProductivity = computed(() => safeValue(derivedMetric('salesRevenue'), derivedMetric('salaryAndProductionCost') + derivedMetric('humanResourceCost')))
const netAssetYield = computed(() => ratio(computeNetProfit(), derivedMetric('lineResidual') + derivedMetric('depreciableAssetTotal')))
const totalAssetYield = computed(() => ratio(computeNetProfit(), props.periodEndCash + derivedMetric('lineResidual') + derivedMetric('depreciableAssetTotal')))
const netProfitRate = computed(() => ratio(computeNetProfit(), derivedMetric('salesRevenue')))
const grossMarginRate = computed(() => ratio(derivedMetric('salesRevenue') - derivedMetric('directCost'), derivedMetric('salesRevenue')))

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

function computeNetProfit() {
  return derivedMetric('salesRevenue') - derivedMetric('directCost') - derivedMetric('comprehensiveCostTotal') - derivedMetric('financeIncomeExpense') + derivedMetric('extraIncomeExpense')
}

function safeValue(numerator: number, denominator: number) {
  if (denominator === 0) {
    return 0
  }
  return numerator / denominator
}

function ratio(numerator: number, denominator: number) {
  return safeValue(numerator, denominator)
}

function isScopeEditable(scope: string) {
  return props.editableScopes.includes(scope)
}

function isScopeInvalidDraft(scope: string) {
  return props.invalidScopes.includes(scope)
}

function editableCellClass(scope: string) {
  if (isScopeEditable(scope)) {
    return 'input-cell'
  }
  if (isScopeInvalidDraft(scope)) {
    return 'invalid-draft-cell'
  }
  return 'locked-cell'
}

function isAdminIssuedExtraField(fieldKey: string) {
  return fieldKey === 'extraExpensePenalty' || fieldKey === 'extraIncomeReward'
}

function extraCellClass(fieldKey: string, scope: string) {
  if (isAdminIssuedExtraField(fieldKey)) {
    return 'admin-issued-cell'
  }
  return editableCellClass(scope)
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
  return Number.isNaN(parsed) ? '' : parsed
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

function createMarketBidRow(): MarketBidRow {
  return {
    basicProductTotal: '',
    standardProductTotal: '',
    precisionProductTotal: '',
    intelligentProductTotal: '',
    orderAmount: '',
  }
}

function buildFixedMarketBidRows(source: Array<Record<string, NumericCellValue>>) {
  return marketRegions.map((_, index) => ({
    ...createMarketBidRow(),
    ...(source[index] ?? {}),
  })) as MarketBidRow[]
}

function ensureMarketBidRows(next: OperatingPayload) {
  const normalized = buildFixedMarketBidRows(next.beginning.marketBid)
  next.beginning.marketBid = normalized.map((item) => ({ ...item })) as Array<Record<string, NumericCellValue>>
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

function getMarketRowTotal(index: number) {
  const row = marketBidRows.value[index]
  return marketProductFields.reduce((total, field) => total + toNumber(row[field.key]), 0)
}

function syncMarketBidDerived(next: OperatingPayload) {
  ensureMarketBidRows(next)

  let orderTotal = 0
  next.beginning.marketBid.forEach((row) => {
    const typedRow = row as Record<string, NumericCellValue>
    const rowTotal = marketProductFields.reduce((total, field) => total + toNumber(typedRow[field.key]), 0)
    typedRow.orderAmount = rowTotal
    orderTotal += rowTotal
  })

  const marketInvestment = normalizeNumericValue(String(next.beginning.taxAndPlanning.marketInvestmentTotal ?? ''))
  next.beginning.taxAndPlanning.marketInvestmentTotal = marketInvestment
  next.beginning.taxAndPlanning.marketBidCost = marketInvestment
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
  letter-spacing: 2px;
  font-weight: 700;
  font-size: 14px;
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
.empty-cell {
  padding: 8px 10px;
  font-size: 13px;
  line-height: 1.45;
}

.group-title,
.result-title {
  background: #fbfbfc;
  font-weight: 700;
}

.note-cell,
.metric-label {
  background: #f8fafc;
}

.metric-label {
  text-align: center;
  font-weight: 700;
}

.empty-cell {
  background: #ffffff;
}

.center {
  text-align: center;
}

.input-cell,
.locked-cell,
.invalid-draft-cell,
.result-cell,
.orange-cell,
.quarter-cash-cell {
  text-align: center;
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

.input-cell input,
.locked-cell input,
.invalid-draft-cell input {
  width: 100%;
  border: none;
  background: transparent;
  padding: 8px 10px;
  outline: none;
  font: inherit;
  color: inherit;
  text-align: center;
}

.admin-issued-cell {
  background: #eef3fb;
  color: #526175;
  font-weight: 700;
  padding: 8px 10px;
}

.cell-readonly-value {
  display: block;
  min-height: 20px;
}
.locked-cell input {
  color: #8491a6;
}

.invalid-draft-cell input {
  color: #8a5a17;
}

.result-cell,
.quarter-cash-cell {
  background: var(--calc-bg);
  padding: 8px 10px;
  font-size: 13px;
  font-weight: 700;
}

.quarter-cash-cell {
  background: #f8c15d;
}

.orange-cell {
  background: #f8c15d;
  padding: 8px 10px;
  font-size: 13px;
  font-weight: 700;
}

@media (max-width: 1400px) {
  .sheet-table {
    min-width: 1180px;
  }
}
</style>







