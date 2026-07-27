<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>订单管理</h2>
        
      </div>
      <div class="hero-actions">
        <select v-model.number="store.selectedYearNo" class="year-select" :disabled="loading || savingConfig || generatingPool" @change="handleYearChange">
          <option v-for="item in yearOptions" :key="item.value" :value="item.value">{{ item.label }}</option>
        </select>
        <button type="button" class="btn" :disabled="loading" @click="handleRefresh">刷新</button>
      </div>
    </header>

    <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
      {{ pageMessage.text }}
    </section>

    <section class="stats-grid">
      <article class="stat-card">
        <span>当前年份</span>
        <strong>{{ selectedYearNo }}年</strong>
      </article>
      <article class="stat-card">
        <span>生成状态</span>
        <strong>{{ formatGenerationStatus(config?.generationStatus) }}</strong>
      </article>
      <article class="stat-card">
        <span>配置订单数</span>
        <strong>{{ totalOrderCount }}</strong>
      </article>
      <article class="stat-card">
        <span>已开启市场</span>
        <strong>{{ enabledMarketCount }}</strong>
      </article>
      <article class="stat-card">
        <span>已生成订单</span>
        <strong>{{ totalGeneratedCount }}</strong>
      </article>
    </section>

    <nav class="flow-tabs" aria-label="订单管理流程">
      <button type="button" class="flow-tab" :class="{ active: activeOrderTab === 'forecast' }" @click="activeOrderTab = 'forecast'">
        <strong>1. 数量控制</strong>
        <span>{{ lockedForecastYearCount > 0 ? `${lockedForecastYearCount} 个年份已锁定` : '可编辑' }}</span>
      </button>
      <button type="button" class="flow-tab" :class="{ active: activeOrderTab === 'market' }" @click="activeOrderTab = 'market'">
        <strong>2. 市场设置</strong>
        <span>{{ enabledMarketCount }} 个市场开启</span>
      </button>
      <button type="button" class="flow-tab" :class="{ active: activeOrderTab === 'sequence' }" @click="activeOrderTab = 'sequence'">
        <strong>3. 标段顺序</strong>
        <span>{{ totalOrderCount }} 单配置</span>
      </button>
      <button type="button" class="flow-tab" :class="{ active: activeOrderTab === 'pool' }" @click="activeOrderTab = 'pool'">
        <strong>4. 订单池</strong>
        <span>{{ formatGenerationStatus(config?.generationStatus) }}</span>
      </button>
      <button type="button" class="flow-tab" :class="{ active: activeOrderTab === 'bidding' }" @click="activeOrderTab = 'bidding'">
        <strong>5. 竞标控制</strong>
        <span>{{ formatSegmentStatus(marketSelectionStatus?.marketBidStatus) }}</span>
      </button>
    </nav>

    <template v-if="activeOrderTab === 'forecast'">
    <section class="panel-card forecast-control-panel">
      <div class="panel-head">
        <div>
          <strong>多年订单数量控制台</strong>
          <span>保存后同步刷新未确认年份的市场预测和预览订单池；已确认年份整列锁定。</span>
        </div>
        <button type="button" class="btn primary" :disabled="savingForecastControl || generatingPool" @click="handleSaveForecastControl">
          {{ savingForecastControl ? '保存中...' : '保存数量并刷新预览' }}
        </button>
      </div>
      <div v-for="stage in forecastStages" :key="stage.forecastStageCode" class="forecast-stage-block">
        <div class="forecast-stage-head">
          <strong>{{ stage.forecastStageName }}</strong>
          <span>公式版本 {{ forecastControl?.forecast.formulaVersion ?? config?.forecast.formulaVersion ?? '--' }}</span>
        </div>
        <div class="table-scroll">
          <table class="forecast-control-table">
            <thead>
              <tr>
                <th>市场</th>
                <th>订单类型</th>
                <th v-for="yearNo in stage.years" :key="yearNo" :class="{ 'locked-year-cell': isForecastYearLocked(yearNo) }">
                  {{ yearNo }}年数量
                  <span v-if="isForecastYearLocked(yearNo)" class="locked-year-label">{{ forecastYearLockReason(yearNo) }}</span>
                </th>
              </tr>
            </thead>
            <tbody>
              <template v-for="market in marketOptions" :key="`${stage.forecastStageCode}-${market.code}`">
                <tr v-for="(orderType, orderIndex) in orderTypeOptions" :key="`${market.code}-${orderType.code}`">
                  <td v-if="orderIndex === 0" :rowspan="orderTypeOptions.length">{{ market.name }}</td>
                  <td>{{ orderType.name }}</td>
                  <td
                    v-for="yearNo in stage.years"
                    :key="yearNo"
                    class="forecast-input-cell"
                    :class="{ 'locked-year-cell': isForecastYearLocked(yearNo) }"
                  >
                    <input
                      :ref="(el) => setForecastInputRef(stage.forecastStageCode, yearNo, market.code, orderType.code, el)"
                      :value="getForecastItem(yearNo, market.code, orderType.code)?.orderCount ?? 0"
                      type="number"
                      min="0"
                      :max="maxCardCount > 0 ? maxCardCount : undefined"
                      step="1"
                      class="compact-input"
                      :class="{ invalid: hasFractionInput(getForecastItem(yearNo, market.code, orderType.code)?.orderCount ?? 0) }"
                      :disabled="savingForecastControl || generatingPool || isForecastYearLocked(yearNo)"
                      data-admin-forecast-nav
                      @keydown="handleForecastCountNavigation"
                      @input="handleForecastCountInput(yearNo, market.code, orderType.code, $event)"
                    >
                  </td>
                </tr>
                <tr class="forecast-total-row">
                  <td colspan="2">{{ market.name }}预测金额</td>
                  <td v-for="yearNo in stage.years" :key="yearNo" class="number-cell" :class="{ 'locked-year-cell': isForecastYearLocked(yearNo) }">{{ formatIntegerAmount(getForecastYear(stage.forecastStageCode, market.code, yearNo)?.totalForecastAmount ?? 0) }}</td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
        <div class="forecast-chart-grid">
          <section v-for="market in stage.markets" :key="`${stage.forecastStageCode}-${market.marketCode}-chart`" class="forecast-market-card">
            <div class="forecast-market-head">
              <strong>{{ market.marketName }}</strong>
              <span>{{ formatIntegerAmount(stageMarketAmount(market)) }}</span>
            </div>
            <p v-if="market.narrative">{{ market.narrative }}</p>
            <div class="forecast-chart" :style="forecastChartColumns(market)">
              <div v-for="year in market.years" :key="year.yearNo" class="forecast-year-group">
                <div class="forecast-bars">
                  <span
                    v-for="product in year.products"
                    :key="product.orderType"
                    class="forecast-bar"
                    :class="`type-${product.orderType}`"
                    :style="{ height: `${forecastBarHeight(product.forecastAmount, market)}%` }"
                    :title="`${year.yearNo}年 ${product.orderTypeName}: ${product.orderCount}单 / ${formatIntegerAmount(product.forecastAmount)}`"
                  />
                </div>
                <span class="forecast-year-label">{{ year.yearNo }}年</span>
                <em>{{ year.totalOrderCount }}单 / {{ formatIntegerAmount(year.totalForecastAmount) }}</em>
              </div>
            </div>
            <div class="forecast-legend">
              <span v-for="orderType in orderTypeOptions" :key="orderType.code" :class="`type-${orderType.code}`">
                {{ orderType.name }}
              </span>
            </div>
          </section>
        </div>
        <div class="forecast-narratives">
          <label v-for="market in marketOptions" :key="`${stage.forecastStageCode}-${market.code}`" class="field">
            <span>{{ market.name }}说明</span>
            <textarea
              :value="getForecastNarrative(stage.forecastStageCode, market.code)?.content ?? ''"
              rows="2"
              :disabled="savingForecastControl || generatingPool"
              @input="handleForecastNarrativeInput(stage.forecastStageCode, market.code, $event)"
            />
          </label>
        </div>
      </div>
    </section>
    </template>

    <template v-if="activeOrderTab === 'market'">
    <section class="panel-card">
      <div class="panel-head">
        <div>
          <strong>市场开启与投入上限</strong>
          
        </div>
        <button type="button" class="btn primary" :disabled="savingMarketConfig || !config?.canUpdateConfig" @click="handleSaveMarketConfig">
          {{ savingMarketConfig ? '保存中...' : '保存市场' }}
        </button>
      </div>
      <div class="market-config-grid">
        <div v-for="market in editableMarketConfigs" :key="market.marketCode" class="market-config-item" :class="{ disabled: !market.enabled }">
          <label class="market-toggle">
            <input v-model="market.enabled" type="checkbox" :disabled="savingMarketConfig || !config?.canUpdateConfig" @change="syncMarketDraftToItems">
            <span>{{ market.marketName }}</span>
            <em>{{ market.enabled ? '已开启' : '未开启' }}</em>
          </label>
          <label class="limit-field">
            <span>单市场上限（M）</span>
            <input
              :value="market.marketInvestmentLimit ?? ''"
              type="number"
              min="0"
              step="1"
              inputmode="numeric"
              placeholder="无上限"
              :class="{ invalid: hasFractionInput(market.marketInvestmentLimit ?? '') }"
              :disabled="savingMarketConfig || !config?.canUpdateConfig || !market.enabled"
              data-enter-confirm
              @keydown.enter="confirmInputOnEnter"
              @input="handleMarketLimitInput(market.marketCode, $event)"
            >
          </label>
        </div>
      </div>
    </section>
    </template>

    <template v-if="activeOrderTab === 'sequence'">
    <section class="panel-card">
      <div class="panel-head">
        <div>
          <strong>标段释放顺序</strong>
        </div>
        <button type="button" class="btn primary" :disabled="savingConfig || !config?.canUpdateConfig" @click="handleSaveConfig">
          {{ savingConfig ? '保存中...' : '保存顺序' }}
        </button>
      </div>

      <div class="table-scroll" data-enter-nav-scope @keydown="handleSequentialInputNavigation">
        <table class="config-table release-sequence-table">
          <colgroup>
            <col class="release-order-col">
            <col>
            <col>
            <col class="generated-col">
            <col class="status-col">
          </colgroup>
          <thead>
            <tr>
              <th>释放顺序</th>
              <th>市场</th>
              <th>订单类型</th>
              <th>已生成</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="sortedItems.length === 0">
              <td colspan="5" class="empty-row">暂无订单配置。</td>
            </tr>
            <tr v-for="item in sortedItems" :key="`${item.marketCode}-${item.orderType}`" :class="{ 'row-disabled': !item.marketEnabled }">
              <td class="release-sequence-cell">
                <input v-model.number="item.releaseSequenceNo" type="number" min="1" step="1" inputmode="numeric" data-enter-nav :disabled="savingConfig || generatingPool || !config?.canUpdateConfig || !item.marketEnabled" class="compact-input release-sequence-input" :class="{ invalid: hasFractionInput(item.releaseSequenceNo) }">
              </td>
              <td>
                {{ item.marketName }}
                <span v-if="!item.marketEnabled" class="muted-inline">市场未开启</span>
              </td>
              <td>{{ item.orderTypeName }}</td>
              <td class="number-cell">{{ item.generatedCount }}</td>
              <td>
                <span v-if="!item.marketEnabled" class="status-tag muted">市场未开启</span>
                <span v-else class="status-tag" :class="item.configStatus === 'LOCKED' ? 'locked' : 'draft'">
                  {{ item.configStatus === 'LOCKED' ? '已锁定' : '草稿' }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
    </template>

    <template v-if="activeOrderTab === 'pool'">
    <section class="panel-card generation-panel">
      <div class="panel-head">
        <div>
          <strong>订单池状态</strong>
          <span>{{ poolStatusTip }}</span>
        </div>
        <button type="button" class="btn primary" :disabled="confirmingPool || !config?.canConfirmPool" @click="handleConfirmPool">
          {{ confirmingPool ? '确认中...' : '确认订单池' }}
        </button>
      </div>
      <div class="batch-grid">
        <article class="batch-card">
          <span>预览批次</span>
          <strong>{{ config?.latestPreviewBatch ? `#${config.latestPreviewBatch.batchId}` : '--' }}</strong>
          <em>{{ config?.latestPreviewBatch ? `${config.latestPreviewBatch.generatedCount} 单 · ${formatDateTime(config.latestPreviewBatch.generatedAt)}` : '尚未生成预览' }}</em>
        </article>
        <article class="batch-card">
          <span>确认批次</span>
          <strong>{{ config?.confirmedBatch ? `#${config.confirmedBatch.batchId}` : '--' }}</strong>
          <em>{{ config?.confirmedBatch?.confirmedAt ? `${config.confirmedBatch.generatedCount} 单 · ${formatDateTime(config.confirmedBatch.confirmedAt)}` : '尚未确认' }}</em>
        </article>
        <article class="batch-card">
          <span>公式版本</span>
          <strong>{{ config?.latestPreviewBatch?.formulaVersion ?? config?.confirmedBatch?.formulaVersion ?? '--' }}</strong>
          <em>{{ config?.latestPreviewBatch?.randomSeed ? `Seed ${config.latestPreviewBatch.randomSeed}` : '确认后固化随机种子' }}</em>
        </article>
      </div>
    </section>

    <section v-if="config?.warnings.length" class="warning-list">
      <strong>数量风险提示</strong>
      <span v-for="warning in config.warnings" :key="`${warning.level}-${warning.message}`">{{ warning.message }}</span>
    </section>
    </template>

    <template v-if="activeOrderTab === 'bidding'">
    <section class="panel-card control-panel">
      <div class="panel-head">
        <div>
          <strong>市场竞标控制</strong>
          
        </div>
        <button type="button" class="btn" :disabled="loadingSelectionStatus" @click="handleLoadSelectionStatus">
          {{ loadingSelectionStatus ? '加载中...' : '刷新状态' }}
        </button>
      </div>

      <div class="sequence-action-bar">
        <div>
          <strong>选单顺序</strong>
          <span>{{ marketInvestmentSubmitSummary }}</span>
        </div>
        <button type="button" class="btn primary" :disabled="generatingSequence" @click="handleGenerateSelectionSequence">
          {{ generatingSequence ? '生成中...' : '生成选单顺序' }}
        </button>
      </div>

      <div class="control-grid">
        <label class="field">
          <span>控制市场</span>
          <select v-model="store.controlForm.marketCode" @change="handleControlMarketChange">
            <option v-for="item in marketOptions" :key="item.code" :value="item.code">{{ item.name }}</option>
          </select>
        </label>
        <div class="control-actions">
          <button type="button" class="btn primary" :disabled="releasingSegment" @click="handleReleaseNextSegment">
            {{ releasingSegment ? '释放中...' : '释放下一个标段' }}
          </button>
          <button
            type="button"
            class="btn primary"
            :disabled="currentSegment?.segmentStatus !== 'ROUND_READY' || !currentSegment.nextRoundNo || openingNextRound"
            @click="handleOpenNextRound"
          >
            {{ openingNextRound ? '开启中...' : '开启下一轮' }}
          </button>
        </div>
      </div>

      <div class="selection-overview">
        <article class="status-mini">
          <span>市场状态</span>
          <strong>{{ formatSegmentStatus(marketSelectionStatus?.marketBidStatus) }}</strong>
        </article>
        <article class="status-mini">
          <span>市场龙头</span>
          <strong>{{ marketLeaderText }}</strong>
          <em v-if="marketLeaderAmountText">{{ marketLeaderAmountText }}</em>
        </article>
        <article class="status-mini">
          <span>当前选单标段</span>
          <strong>{{ currentSegment ? `${currentSegment.marketName} ${currentSegment.orderTypeName}` : '--' }}</strong>
        </article>
        <article class="status-mini">
          <span>当前小组</span>
          <strong>{{ currentGroupName }}</strong>
        </article>
        <article class="status-mini">
          <span>当前 / 下一轮</span>
          <strong>{{ currentSegment?.currentRoundNo ? `第 ${currentSegment.currentRoundNo} 轮` : '--' }}</strong>
          <em v-if="currentSegment?.nextRoundNo">下一轮：第 {{ currentSegment.nextRoundNo }} 轮</em>
        </article>
        <article class="status-mini">
          <span>理论机会 / 剩余订单</span>
          <strong>{{ selectedSequenceSegment?.theoreticalMaxSelections ?? 0 }} / {{ selectedSequenceSegment?.availableCount ?? 0 }}</strong>
        </article>
      </div>

      <div v-if="selectedSequenceSegment?.warnings?.length" class="warning-list">
        <span v-for="warning in selectedSequenceSegment.warnings" :key="warning">{{ warning }}</span>
      </div>

      <div class="skip-row">
        <input v-model="store.controlForm.skipReason" type="text" placeholder="代跳过原因，必填">
        <button type="button" class="btn danger" :disabled="!currentSegment?.currentGroupId || skippingGroup" @click="handleSkipCurrentGroup">
          {{ skippingGroup ? '跳过中...' : '跳过当前小组' }}
        </button>
      </div>

      <div class="table-scroll">
        <table class="selection-table">
          <thead>
            <tr>
              <th>释放顺序</th>
              <th>市场</th>
              <th>订单类型</th>
              <th>标段状态</th>
              <th>当前组</th>
              <th>可选</th>
              <th>已选</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!marketSelectionStatus || marketSelectionStatus.segments.length === 0">
              <td colspan="7" class="empty-row">当前市场暂无标段状态，请先生成订单池。</td>
            </tr>
            <tr
              v-for="segment in marketSelectionStatus?.segments ?? []"
              :key="`${segment.marketCode}-${segment.orderType}`"
              class="clickable-row"
              :class="{ selected: selectedSequenceSegmentKey === segmentKey(segment) }"
              @click="selectSequenceSegment(segment)"
            >
              <td>#{{ segment.releaseSequenceNo }}</td>
              <td>{{ segment.marketName }}</td>
              <td>{{ segment.orderTypeName }}</td>
              <td>{{ formatSegmentStatus(segment.segmentStatus) }}</td>
              <td>{{ formatGroupName(segment.currentGroupId) }}</td>
              <td class="number-cell">{{ segment.availableCount }}</td>
              <td class="number-cell">{{ segment.selectedCount }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-scroll sequence-scroll">
        <table class="selection-table">
          <thead>
            <tr>
              <th colspan="7" class="sequence-title">
                {{ selectedSequenceSegment ? `${selectedSequenceSegment.marketName} · ${selectedSequenceSegment.orderTypeName} 选单顺序` : '选单顺序' }}
              </th>
            </tr>
            <tr>
              <th>顺序</th>
              <th>小组</th>
              <th>本标段投入</th>
              <th>上年该市场订单额</th>
              <th>市场龙头</th>
              <th>状态</th>
              <th>已选订单</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!selectedSequenceSegment">
              <td colspan="7" class="empty-row">当前选中标段暂无选单顺序，请先生成选单顺序。</td>
            </tr>
            <template v-for="roundNo in selectionRoundNumbers" :key="roundNo">
              <tr class="round-title-row">
                <td colspan="7">第 {{ roundNo }} 轮选单顺序</td>
              </tr>
              <tr v-if="roundSelectionItems(roundNo).length === 0">
                <td colspan="7" class="empty-row">本轮无参与小组</td>
              </tr>
              <tr v-for="item in roundSelectionItems(roundNo)" :key="`${item.roundNo}-${item.groupId}`">
                <td>#{{ item.sequenceNo }}</td>
                <td>{{ item.groupName }}</td>
                <td class="number-cell">{{ formatAmount(item.marketInvestment) }}</td>
                <td class="number-cell">{{ formatIntegerAmount(item.previousMarketOrderAmount) }}</td>
                <td>{{ item.isMarketLeader ? '是' : '否' }}</td>
                <td>{{ formatSelectionStatus(item.selectionStatus) }}</td>
                <td>{{ item.selectedOrderNo || (item.selectedOrderId ? `#${item.selectedOrderId}` : '--') }}</td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </section>
    </template>

    <template v-if="activeOrderTab === 'pool'">
    <section class="panel-card">
      <div class="panel-head">
        <div>
          <strong>订单池查看</strong>
          <span>默认查看全部订单，可按年份、市场和订单类型筛选。</span>
        </div>
        <div class="hero-actions">
          <select v-model.number="store.selectedYearNo" class="year-select" :disabled="loading || loadingPool" @change="handleYearChange">
            <option v-for="item in yearOptions" :key="item.value" :value="item.value">{{ item.label }}</option>
          </select>
          <button type="button" class="btn" :disabled="loadingPool" @click="handleLoadPool">
            {{ loadingPool ? '加载中...' : '刷新' }}
          </button>
        </div>
      </div>

      <div class="pool-filter">
        <label class="field">
          <span>市场</span>
          <select v-model="store.poolFilter.marketCode" @change="handleLoadPool">
            <option value="ALL">全部市场</option>
            <option v-for="item in marketOptions" :key="item.code" :value="item.code">{{ item.name }}</option>
          </select>
        </label>
        <label class="field">
          <span>订单类型</span>
          <select v-model="store.poolFilter.orderType" @change="handleLoadPool">
            <option value="ALL">全部类型</option>
            <option v-for="item in orderTypeOptions" :key="item.code" :value="item.code">{{ item.name }}</option>
          </select>
        </label>
      </div>

      <div class="table-scroll">
        <table class="pool-table">
          <thead>
            <tr>
              <th>订单编号</th>
              <th>市场</th>
              <th>订单类型</th>
              <th>金额</th>
              <th>数量</th>
              <th>单价</th>
              <th>账期</th>
              <th v-for="field in extraPoolFields" :key="field.code">{{ field.name }}</th>
              <th>状态</th>
              <th>选中组</th>
              <th>来源</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!orderPool || orderPool.list.length === 0">
              <td :colspan="10 + extraPoolFields.length" class="empty-row">当前筛选下没有订单池记录。</td>
            </tr>
            <tr v-for="item in orderPool?.list ?? []" :key="item.orderId">
              <td>{{ item.businessOrderNo || `#${item.orderId}` }}</td>
              <td>{{ item.marketName }}</td>
              <td>{{ item.orderTypeName }}</td>
              <td class="number-cell">{{ formatIntegerAmount(item.orderAmount) }}</td>
              <td class="number-cell">{{ formatQuantity(item.orderQuantity) }}</td>
              <td class="number-cell">{{ formatUnitPrice(item.unitPrice) }}</td>
              <td>{{ item.accountTerm }}季度</td>
              <td v-for="field in extraPoolFields" :key="field.code">{{ formatOrderField(item, field) }}</td>
              <td>{{ formatPoolStatus(item.poolStatus) }}</td>
              <td>{{ item.selectedGroupId ? `组ID ${item.selectedGroupId}` : '--' }}</td>
              <td>{{ item.sourceSheetName }} {{ item.sourceCell }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
    </template>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'

import { useAdminShellStore } from '@/stores/admin-shell'
import { useAdminOrderStore } from '@/stores/admin-order'
import { useDictionaryStore } from '@/stores/dictionary'
import type { OrderPoolItem, OrderPoolStatus } from '@/types/sandbox-game-admin'
import type { AdminOrderSegmentStatus, OrderMarketForecastMarket, OrderTemplateField } from '@/types/sandbox-game-order'
import {
  confirmInputOnEnter,
  focusAdjacentInputFromList,
  handleSequentialInputNavigation,
  resolveInputNavigationDirection,
} from '@/utils/input-navigation'
import { hasFractionInput } from '@/utils/manual-integer'

const shellStore = useAdminShellStore()
const store = useAdminOrderStore()
const dictionaryStore = useDictionaryStore()

const { config: shellConfig } = storeToRefs(shellStore)
const {
  selectedYearNo,
  config,
  orderPool,
  loading,
  savingConfig,
  savingForecastControl,
  savingMarketConfig,
  generatingPool,
  confirmingPool,
  generatingSequence,
  loadingPool,
  loadingSelectionStatus,
  releasingSegment,
  openingNextRound,
  skippingGroup,
  pageMessage,
  marketSelectionStatus,
  currentSegment,
  sortedItems,
  forecastControl,
  forecastStages,
  forecastYearLocks,
  editableForecastItems,
  editableForecastNarratives,
  editableMarketConfigs,
  totalOrderCount,
  totalGeneratedCount,
  enabledMarketCount,
  marketOptions: templateMarketOptions,
  orderTypeOptions: templateOrderTypeOptions,
  maxCardCount,
  segmentCount,
} = storeToRefs(store)

const marketOptions = computed(() =>
  templateMarketOptions.value.map((item) => ({
    ...item,
    name: dictionaryStore.marketName(item.code, item.name),
  })),
)
const orderTypeOptions = computed(() =>
  templateOrderTypeOptions.value.map((item) => ({
    ...item,
    name: dictionaryStore.orderTypeName(item.code, item.name),
  })),
)
const extraPoolFields = computed(() =>
  (config.value?.orderTemplate.fields ?? orderPool.value?.orderTemplate.fields ?? [])
    .filter((field) => !['orderAmount', 'orderQuantity', 'unitPrice', 'accountTerm'].includes(field.code))
    .sort((a, b) => a.order - b.order),
)
const selectedSequenceSegmentKey = ref('')
const activeOrderTab = ref<'forecast' | 'market' | 'sequence' | 'pool' | 'bidding'>('market')
const forecastInputRefs = new Map<string, HTMLInputElement>()
let selectionStatusTimer: number | null = null
const lockedForecastYearCount = computed(() => forecastYearLocks.value.filter((item) => item.locked).length)

const yearOptions = computed(() => {
  const finalYear = Math.max(shellConfig.value?.finalYear ?? config.value?.finalYear ?? 1, 1)
  return Array.from({ length: finalYear }, (_, index) => ({
    value: index + 1,
    label: `${index + 1}年`,
  }))
})
const currentGroupName = computed(() => formatGroupName(currentSegment.value?.currentGroupId ?? null))
const selectedSequenceSegment = computed(() => {
  const segments = marketSelectionStatus.value?.segments ?? []
  if (selectedSequenceSegmentKey.value) {
    const matched = segments.find((item) => segmentKey(item) === selectedSequenceSegmentKey.value)
    if (matched) {
      return matched
    }
  }
  return currentSegment.value
    ?? segments.find((item) => item.selectionOrder.length > 0 && item.segmentStatus === 'SEQUENCE_READY')
    ?? segments.find((item) => item.selectionOrder.length > 0)
    ?? null
})
const selectionRoundNumbers = computed(() => {
  const maxRounds = Math.max(1, Number(selectedSequenceSegment.value ? store.orderTemplate.maxSelectionRounds : 1) || 1)
  return Array.from({ length: maxRounds }, (_, index) => index + 1)
})
const marketLeaderOrder = computed(() => selectedSequenceSegment.value?.selectionOrder.find((item) => item.isMarketLeader) ?? null)
const marketLeaderText = computed(() => {
  const leader = marketLeaderOrder.value
  if (leader) {
    return leader.groupName || `组ID ${leader.groupId}`
  }
  return marketSelectionStatus.value?.leaderGroupId ? formatGroupName(marketSelectionStatus.value.leaderGroupId) : '--'
})
const marketLeaderAmountText = computed(() => {
  const leader = marketLeaderOrder.value
  if (!leader) {
    return ''
  }
  return `上年该市场订单额 ${formatIntegerAmount(leader.previousMarketOrderAmount)}`
})
const poolStatusTip = computed(() => {
  if (config.value?.confirmedBatch) {
    return '订单池已确认。'
  }
  if (config.value?.latestPreviewBatch) {
    return '当前为预览订单，尚未保存成正式订单；确认后才会开放玩家市场投入。'
  }
  return '请先在数量控制台保存订单数量，系统会自动生成预览订单池。'
})
const marketInvestmentSubmitSummary = computed(() => {
  const bids = marketSelectionStatus.value?.bids ?? []
  if (bids.length === 0) {
    return '暂无小组提交市场投入。'
  }
  const activeBids = bids.filter((item) => item.businessStatus !== 'BANKRUPT')
  const submitted = activeBids.filter((item) => item.submitted).length
  const missing = activeBids.filter((item) => !item.submitted).map((item) => item.groupName)
  if (missing.length === 0) {
    return `${submitted}/${activeBids.length} 组已提交，可生成选单顺序。`
  }
  return `${submitted}/${activeBids.length} 组已提交，未提交：${missing.join('、')}`
})

watch(
  () => marketSelectionStatus.value,
  () => {
    const segments = marketSelectionStatus.value?.segments ?? []
    const selectedStillExists = segments.some((item) => segmentKey(item) === selectedSequenceSegmentKey.value)
    if (selectedStillExists) {
      return
    }
    const fallback = currentSegment.value
      ?? segments.find((item) => item.selectionOrder.length > 0 && item.segmentStatus === 'SEQUENCE_READY')
      ?? segments.find((item) => item.selectionOrder.length > 0)
    selectedSequenceSegmentKey.value = fallback ? segmentKey(fallback) : ''
  },
)

watch(
  () => activeOrderTab.value,
  (tab) => {
    if (tab === 'bidding') {
      void refreshSelectionStatusSilently()
      startSelectionStatusPolling()
      return
    }
    stopSelectionStatusPolling()
  },
)

onMounted(async () => {
  try {
    if (!shellStore.config) {
      await shellStore.bootstrap()
    }
    const defaultYear = Math.max(shellConfig.value?.currentOpenYear ?? 1, 1)
    await Promise.all([
      store.bootstrap(defaultYear),
      dictionaryStore.loadCurrent(shellStore.config?.editionCode ?? shellStore.setupStatus?.editionCode, { silent: true }),
    ])
    dictionaryStore.startSilentSync()
    await store.loadPool({ silent: true })
    await store.loadSelectionStatus({ silent: true })
    if (activeOrderTab.value === 'bidding') {
      startSelectionStatusPolling()
    }
  } catch {
    // 页面消息由 store 统一处理。
  }
})

onUnmounted(() => {
  stopSelectionStatusPolling()
  dictionaryStore.stopSilentSync()
})

async function handleYearChange() {
  try {
    await store.loadConfig()
    await store.loadPool({ silent: true })
    await store.loadSelectionStatus({ silent: true })
    if (activeOrderTab.value === 'bidding') {
      restartSelectionStatusPolling()
    }
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleRefresh() {
  try {
    await store.loadForecastControl({ silent: true })
    await store.loadConfig()
    await store.loadPool({ silent: true })
    await store.loadSelectionStatus({ silent: true })
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleSaveConfig() {
  try {
    await store.saveConfig()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleSaveForecastControl() {
  try {
    await store.saveForecastControl()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleSaveMarketConfig() {
  try {
    await store.saveMarketConfig()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

function syncMarketDraftToItems() {
  const enabledMap = new Map(editableMarketConfigs.value.map((item) => [item.marketCode, item.enabled]))
  for (const item of sortedItems.value) {
    item.marketEnabled = enabledMap.get(item.marketCode) ?? item.marketEnabled
  }
}

function handleMarketLimitInput(marketCode: string, event: Event) {
  const input = event.target as HTMLInputElement
  const market = editableMarketConfigs.value.find((item) => item.marketCode === marketCode)
  if (!market) {
    return
  }
  if (input.value === '') {
    market.marketInvestmentLimit = null
    return
  }
  const raw = Number(input.value)
  market.marketInvestmentLimit = Number.isFinite(raw) ? raw : null
}

function getForecastItem(yearNo: number, marketCode: string, orderType: string) {
  return editableForecastItems.value.find((item) => item.yearNo === yearNo && item.marketCode === marketCode && item.orderType === orderType) ?? null
}

function forecastInputKey(stageCode: string, yearNo: number, marketCode: string, orderType: string) {
  return `${stageCode}|${yearNo}|${marketCode}|${orderType}`
}

function setForecastInputRef(stageCode: string, yearNo: number, marketCode: string, orderType: string, el: unknown) {
  const key = forecastInputKey(stageCode, yearNo, marketCode, orderType)
  if (el instanceof HTMLInputElement) {
    forecastInputRefs.set(key, el)
    return
  }

  forecastInputRefs.delete(key)
}

function orderedForecastInputKeys() {
  const keys: string[] = []
  for (const stage of forecastStages.value) {
    for (const yearNo of stage.years) {
      if (isForecastYearLocked(yearNo)) {
        continue
      }
      for (const market of marketOptions.value) {
        for (const orderType of orderTypeOptions.value) {
          keys.push(forecastInputKey(stage.forecastStageCode, yearNo, market.code, orderType.code))
        }
      }
    }
  }
  return keys
}

function handleForecastCountNavigation(event: KeyboardEvent) {
  const direction = resolveInputNavigationDirection(event)
  if (!direction) {
    return
  }

  const keys = orderedForecastInputKeys()
  focusAdjacentInputFromList(event, keys.map((key) => forecastInputRefs.get(key)), direction)
}

function getForecastNarrative(stageCode: string, marketCode: string) {
  return editableForecastNarratives.value.find((item) => item.forecastStageCode === stageCode && item.marketCode === marketCode) ?? null
}

function getForecastYear(stageCode: string, marketCode: string, yearNo: number) {
  return forecastStages.value
    .find((stage) => stage.forecastStageCode === stageCode)
    ?.markets.find((market) => market.marketCode === marketCode)
    ?.years.find((year) => year.yearNo === yearNo) ?? null
}

function forecastYearLock(yearNo: number) {
  return forecastYearLocks.value.find((item) => item.yearNo === yearNo) ?? null
}

function isForecastYearLocked(yearNo: number) {
  return forecastYearLock(yearNo)?.locked ?? false
}

function forecastYearLockReason(yearNo: number) {
  return forecastYearLock(yearNo)?.reason ?? '已锁定'
}

function handleForecastCountInput(yearNo: number, marketCode: string, orderType: string, event: Event) {
  const input = event.target as HTMLInputElement
  const raw = Number(input.value)
  const nonNegative = Number.isFinite(raw) ? Math.max(raw, 0) : 0
  const value = maxCardCount.value > 0 ? Math.min(nonNegative, maxCardCount.value) : nonNegative
  const item = getForecastItem(yearNo, marketCode, orderType)
  if (item) {
    item.orderCount = value
  }
}

function handleForecastNarrativeInput(stageCode: string, marketCode: string, event: Event) {
  const input = event.target as HTMLTextAreaElement
  const item = getForecastNarrative(stageCode, marketCode)
  if (item) {
    item.content = input.value
  }
}

async function handleGeneratePool() {
  try {
    await store.generatePool(true)
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleConfirmPool() {
  try {
    await store.confirmPool()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleGenerateSelectionSequence() {
  try {
    await store.generateSelectionSequence()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleLoadPool() {
  try {
    await store.loadPool()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleLoadSelectionStatus() {
  try {
    await store.loadSelectionStatus()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleControlMarketChange() {
  try {
    selectedSequenceSegmentKey.value = ''
    await store.loadSelectionStatus()
    restartSelectionStatusPolling()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleReleaseNextSegment() {
  try {
    await store.releaseNextSegment()
    if (currentSegment.value) {
      selectedSequenceSegmentKey.value = segmentKey(currentSegment.value)
    }
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function handleOpenNextRound() {
  try {
    await store.openNextRound()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

function roundSelectionItems(roundNo: number) {
  return selectedSequenceSegment.value?.selectionOrder.filter((item) => item.roundNo === roundNo) ?? []
}

async function handleSkipCurrentGroup() {
  try {
    await store.skipCurrentGroup()
  } catch {
    // 页面消息由 store 统一处理。
  }
}

async function refreshSelectionStatusSilently() {
  try {
    await store.loadSelectionStatus({ silent: true })
  } catch {
    // 静默轮询不打断管理员当前操作。
  }
}

function startSelectionStatusPolling() {
  if (selectionStatusTimer !== null) {
    return
  }
  selectionStatusTimer = window.setInterval(() => {
    if (activeOrderTab.value !== 'bidding') {
      stopSelectionStatusPolling()
      return
    }
    void refreshSelectionStatusSilently()
  }, 3000)
}

function stopSelectionStatusPolling() {
  if (selectionStatusTimer === null) {
    return
  }
  window.clearInterval(selectionStatusTimer)
  selectionStatusTimer = null
}

function restartSelectionStatusPolling() {
  if (activeOrderTab.value !== 'bidding') {
    return
  }
  stopSelectionStatusPolling()
  startSelectionStatusPolling()
}

function segmentKey(segment: Pick<AdminOrderSegmentStatus, 'marketCode' | 'orderType'>) {
  return `${segment.marketCode}|${segment.orderType}`
}

function selectSequenceSegment(segment: AdminOrderSegmentStatus) {
  selectedSequenceSegmentKey.value = segmentKey(segment)
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

function formatAmount(value: number) {
  return Number(value || 0).toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })
}

function formatIntegerAmount(value: number) {
  return Math.round(Number(value || 0)).toLocaleString('zh-CN', { maximumFractionDigits: 0 })
}

function formatQuantity(value: number) {
  return Math.round(Number(value || 0)).toLocaleString('zh-CN', { maximumFractionDigits: 0 })
}

function formatUnitPrice(value: number) {
  return Number(value || 0).toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 4 })
}

function formatOrderField(order: OrderPoolItem, field: OrderTemplateField) {
  const value = resolveOrderFieldValue(order, field.code)
  if (value === null || value === undefined || value === '') {
    return '--'
  }
  if (field.valueType === 'money') {
    if (field.code === 'unitPrice') {
      return formatUnitPrice(Number(value))
    }
    return `${formatIntegerAmount(Number(value))}${field.unit ? field.unit : ''}`
  }
  if (field.valueType === 'percent') {
    return `${formatQuantity(Number(value))}${field.unit || '%'}`
  }
  if (field.valueType === 'number') {
    const unit = field.unit ? ` ${field.unit}` : ''
    return `${formatQuantity(Number(value))}${unit}`
  }
  return String(value)
}

function resolveOrderFieldValue(order: OrderPoolItem, code: string) {
  if (code === 'orderAmount') {
    return order.orderAmount
  }
  if (code === 'orderQuantity' || code === 'flightCount') {
    return order.orderPayload?.[code] ?? order.orderQuantity
  }
  if (code === 'unitPrice') {
    return order.unitPrice
  }
  if (code === 'accountTerm') {
    return order.accountTerm
  }
  return order.orderPayload?.[code]
}

function stageMarketAmount(market: OrderMarketForecastMarket) {
  return market.years.reduce((sum, year) => sum + Number(year.totalForecastAmount || 0), 0)
}

function stageMarketMaxAmount(market: OrderMarketForecastMarket) {
  return Math.max(
    ...market.years.flatMap((year) => year.products.map((product) => Number(product.forecastAmount || 0))),
    0,
  )
}

function forecastBarHeight(amount: number, market: OrderMarketForecastMarket) {
  const max = stageMarketMaxAmount(market)
  if (max <= 0) {
    return 0
  }
  return Math.max(4, Math.round((Number(amount || 0) / max) * 100))
}

function forecastChartColumns(market: OrderMarketForecastMarket) {
  return {
    gridTemplateColumns: `repeat(${Math.max(market.years.length, 1)}, minmax(72px, 1fr))`,
  }
}

function formatPoolStatus(value: OrderPoolStatus) {
  if (value === 'AVAILABLE') {
    return '可选'
  }
  if (value === 'SELECTED') {
    return '已选'
  }
  if (value === 'VOID') {
    return '作废'
  }
  return value
}

function formatGenerationStatus(value?: string) {
  const map: Record<string, string> = {
    NOT_GENERATED: '未生成',
    PREVIEW_GENERATED: '预览已生成',
    POOL_CONFIRMED: '订单池已确认',
    SELECTING: '选单中',
    ROUND_READY: '下一轮待开启',
    COMPLETED: '已完成',
  }
  return value ? map[value] ?? value : '--'
}

function formatSegmentStatus(value?: string) {
  const map: Record<string, string> = {
    WAITING_INVESTMENT: '等待投入',
    MARKET_DISABLED: '市场未开启',
    NO_ORDER_CONFIG: '未配置订单',
    BID_OPEN: '投入开放',
    BID_CLOSED: '投入关闭',
    SEQUENCE_READY: '顺序已生成',
    WAITING_RELEASE: '等待释放',
    SELECTING: '选单中',
    COMPLETED: '已完成',
    SKIPPED: '已跳过',
  }
  return value ? map[value] ?? value : '--'
}

function formatSelectionStatus(value: string) {
  const map: Record<string, string> = {
    INELIGIBLE: '无资格',
    WAITING: '待选择',
    CURRENT: '当前选择',
    SELECTED: '已选择',
    PASSED: '已放弃',
    ADMIN_SKIPPED: '管理员跳过',
    INELIGIBLE_BANKRUPT: '已破产，本轮无资格',
  }
  return map[value] ?? value
}

function formatGroupName(groupId?: number | null) {
  if (!groupId) {
    return '--'
  }
  const bid = marketSelectionStatus.value?.bids.find((item) => item.groupId === groupId)
  return bid?.groupName ?? `组ID ${groupId}`
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
.generation-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  align-items: center;
}

.btn,
.year-select {
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

.btn.danger {
  color: var(--danger);
  border-color: #efc4c4;
  background: #fff5f5;
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 14px;
}

.flow-tabs {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px;
  padding: 8px;
  border: 1px solid var(--line);
  border-radius: 16px;
  background: #ffffff;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.05);
}

.flow-tab {
  display: grid;
  gap: 5px;
  min-height: 72px;
  border: 1px solid transparent;
  border-radius: 12px;
  background: #f8fafc;
  padding: 11px 12px;
  color: var(--text);
  text-align: left;
  cursor: pointer;
}

.flow-tab strong {
  font-size: 14px;
}

.flow-tab span {
  overflow: hidden;
  color: var(--muted);
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.flow-tab.active {
  border-color: var(--accent);
  background: var(--accent-soft);
  box-shadow: inset 0 0 0 1px rgba(37, 99, 235, 0.12);
}

.flow-tab.active strong {
  color: var(--accent);
}

.stat-card,
.panel-card {
  border: 1px solid var(--line);
  border-radius: 18px;
  background: #ffffff;
  box-shadow: 0 8px 22px rgba(15, 23, 42, 0.05);
}

.stat-card {
  display: grid;
  gap: 8px;
  padding: 16px;
}

.stat-card span {
  color: var(--muted);
  font-size: 12px;
}

.stat-card strong {
  font-size: 20px;
}

.panel-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  padding: 14px 16px;
  border-bottom: 1px solid var(--line);
  background: #f7f9fc;
}

.panel-head strong {
  display: block;
  margin-bottom: 4px;
  font-size: 16px;
}

.panel-head span,
.inline-tip {
  color: var(--muted);
  font-size: 13px;
}

.generation-actions,
.batch-grid,
.pool-filter {
  padding: 16px;
}

.batch-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  border-top: 1px solid var(--line);
}

.batch-card {
  display: grid;
  gap: 6px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: #fbfcfe;
  padding: 12px;
}

.batch-card span,
.batch-card em {
  color: var(--muted);
  font-size: 13px;
  font-style: normal;
}

.warning-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.warning-list span {
  padding: 7px 10px;
  border-radius: 999px;
  border: 1px solid var(--line);
  background: #f8fafc;
  font-size: 12px;
}

.warning-list {
  padding: 12px 14px;
  border: 1px solid #efd4aa;
  border-radius: 14px;
  background: #fffaf0;
}

.warning-list strong {
  width: 100%;
  color: var(--warning);
}

.warning-list.compact {
  margin-top: 12px;
  padding: 0;
  border: 0;
  background: transparent;
}

.table-scroll {
  overflow: auto;
}

.config-table,
.pool-table,
.forecast-control-table {
  width: 100%;
  min-width: 900px;
  border-collapse: collapse;
}

.release-sequence-table {
  min-width: 720px;
  table-layout: fixed;
}

.config-table th,
.config-table td,
.pool-table th,
.pool-table td,
.forecast-control-table th,
.forecast-control-table td {
  border: 1px solid var(--line);
  padding: 10px 12px;
  font-size: 14px;
}

.config-table th,
.pool-table th,
.forecast-control-table th {
  background: #f4f6f9;
  text-align: center;
}

.release-sequence-table th,
.release-sequence-table td {
  padding: 8px 10px;
}

.release-sequence-table .release-order-col {
  width: 108px;
}

.release-sequence-table .generated-col {
  width: 92px;
}

.release-sequence-table .status-col {
  width: 128px;
}

.release-sequence-input {
  width: 68px;
  padding: 6px 8px;
  border-radius: 8px;
  text-align: center;
}

.forecast-input-cell:focus-within,
.release-sequence-cell:focus-within {
  box-shadow: inset 0 0 0 2px #2563eb, 0 0 0 2px rgba(37, 99, 235, 0.12);
}

.forecast-control-table .locked-year-cell {
  background: #eef2f6;
  color: #64748b;
}

.locked-year-label {
  display: block;
  margin-top: 4px;
  color: #64748b;
  font-size: 11px;
  font-weight: 500;
}

.locked-year-cell .compact-input {
  background: #e5e7eb;
  color: #64748b;
}

.compact-input {
  width: 96px;
  border: 1px solid var(--line);
  border-radius: 10px;
  padding: 8px 10px;
}

.compact-input.invalid,
.limit-field input.invalid {
  border-color: var(--danger);
  background: #fff5f5;
  color: var(--danger);
}

.number-cell {
  text-align: right;
}

.readonly-number {
  display: inline-flex;
  min-width: 56px;
  justify-content: flex-end;
  border: 1px solid var(--line);
  border-radius: 10px;
  background: #f8fafc;
  padding: 8px 10px;
}

.forecast-control-panel {
  overflow: hidden;
}

.forecast-stage-block {
  border-top: 1px solid var(--line);
}

.forecast-stage-block:first-of-type {
  border-top: 0;
}

.forecast-stage-head {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  padding: 12px 16px;
  background: #fbfcfe;
}

.forecast-stage-head span {
  color: var(--muted);
  font-size: 13px;
}

.forecast-total-row {
  background: #f8fafc;
  color: var(--muted);
}

.forecast-chart-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  padding: 14px 16px 0;
}

.forecast-market-card {
  display: grid;
  gap: 10px;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #ffffff;
  padding: 11px 12px;
}

.forecast-market-head {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.forecast-market-head span {
  color: var(--success);
  font-weight: 700;
}

.forecast-market-card p {
  margin: 0;
  color: var(--muted);
  font-size: 12px;
  line-height: 1.5;
}

.forecast-chart {
  display: grid;
  gap: 8px;
  align-items: end;
  min-height: 172px;
  border: 1px solid var(--line);
  border-radius: 10px;
  background:
    linear-gradient(to top, rgba(148, 163, 184, 0.18) 1px, transparent 1px) 0 0 / 100% 25%,
    #fbfcfe;
  padding: 12px 10px 10px;
}

.forecast-year-group {
  display: grid;
  grid-template-rows: 112px auto auto;
  gap: 5px;
  min-width: 0;
  text-align: center;
}

.forecast-bars {
  display: flex;
  gap: 4px;
  align-items: end;
  justify-content: center;
  height: 112px;
  border-bottom: 1px solid #cbd5e1;
}

.forecast-bar {
  width: 10px;
  min-height: 0;
  border-radius: 3px 3px 0 0;
  background: #3b82f6;
}

.forecast-year-label,
.forecast-year-group em {
  overflow: hidden;
  color: var(--muted);
  font-size: 11px;
  font-style: normal;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.forecast-year-label {
  color: var(--text);
  font-weight: 700;
}

.forecast-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 10px;
}

.forecast-legend span {
  display: inline-flex;
  gap: 5px;
  align-items: center;
  color: var(--muted);
  font-size: 11px;
}

.forecast-legend span::before {
  width: 8px;
  height: 8px;
  border-radius: 2px;
  background: currentColor;
  content: "";
}

.type-AGENCY_INSPECTION {
  background-color: #2563eb;
  color: #2563eb;
}

.type-TWO_CABIN_VIP {
  background-color: #16a34a;
  color: #16a34a;
}

.type-BUSINESS_VIP {
  background-color: #d97706;
  color: #d97706;
}

.type-MEMBER_CUSTOM {
  background-color: #7c3aed;
  color: #7c3aed;
}

.forecast-narratives {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  padding: 14px 16px 16px;
}

.empty-row {
  text-align: center;
  color: var(--muted);
}

.status-tag {
  display: inline-flex;
  padding: 5px 9px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
}

.status-tag.draft {
  background: var(--accent-soft);
  color: var(--accent);
}

.status-tag.locked {
  background: #e8f7ee;
  color: var(--success);
}

.status-tag.muted {
  background: #eef2f6;
  color: var(--muted);
}

.muted-inline {
  display: block;
  margin-top: 4px;
  color: var(--muted);
  font-size: 12px;
}

.row-disabled {
  background: #f8fafc;
  color: var(--muted);
}

.market-config-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  padding: 16px;
}

.market-config-item {
  display: grid;
  gap: 10px;
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 12px;
  background: #ffffff;
}

.market-config-item.disabled {
  background: #f8fafc;
  color: var(--muted);
}

.market-config-item em {
  color: var(--muted);
  font-size: 12px;
  font-style: normal;
}

.market-toggle {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: 8px;
  align-items: center;
}

.limit-field {
  display: grid;
  gap: 5px;
}

.limit-field span {
  color: var(--muted);
  font-size: 12px;
}

.limit-field input {
  width: 100%;
  height: 34px;
  border: 1px solid var(--line);
  border-radius: 8px;
  padding: 6px 8px;
}

.pool-filter {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 240px));
  gap: 12px;
}

.sequence-action-bar {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  margin: 16px;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #fbfcfe;
  padding: 12px 14px;
}

.sequence-action-bar strong {
  display: block;
  margin-bottom: 4px;
}

.sequence-action-bar span {
  color: var(--muted);
  font-size: 13px;
}

.control-grid {
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr);
  gap: 12px;
  align-items: end;
  padding: 16px;
}

.control-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.selection-overview {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  padding: 0 16px 16px;
}

.status-mini {
  display: grid;
  gap: 6px;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #fbfcfe;
  padding: 12px;
}

.status-mini span {
  color: var(--muted);
  font-size: 12px;
}

.status-mini strong {
  font-size: 16px;
}

.status-mini em {
  color: var(--muted);
  font-size: 12px;
  font-style: normal;
}

.skip-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 10px;
  padding: 0 16px 16px;
}

.skip-row input {
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 10px 12px;
}

.selection-table {
  width: 100%;
  min-width: 760px;
  border-collapse: collapse;
}

.round-title-row td {
  background: #eef4fb;
  color: var(--text);
  font-weight: 700;
  text-align: left;
}

.selection-table th,
.selection-table td {
  border: 1px solid var(--line);
  padding: 10px 12px;
  font-size: 14px;
}

.selection-table th {
  background: #f4f6f9;
  text-align: center;
}

.clickable-row {
  cursor: pointer;
}

.clickable-row.selected {
  background: var(--accent-soft);
  box-shadow: inset 3px 0 0 var(--accent);
}

.sequence-title {
  color: var(--text);
  text-align: left !important;
}

.sequence-scroll {
  border-top: 1px solid var(--line);
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

.field select {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 12px;
  background: #ffffff;
  padding: 10px 12px;
}

.field textarea {
  width: 100%;
  resize: vertical;
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 9px 10px;
  font: inherit;
}

@media (max-width: 1240px) {
  .stats-grid,
  .flow-tabs {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .hero,
  .panel-head {
    flex-direction: column;
    align-items: flex-start;
  }

  .stats-grid,
  .flow-tabs,
  .batch-grid,
  .sequence-action-bar,
  .pool-filter,
  .forecast-chart-grid,
  .forecast-narratives,
  .market-config-grid,
  .control-grid,
  .selection-overview,
  .skip-row {
    grid-template-columns: 1fr;
  }
}
</style>
