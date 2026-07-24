<template>
  <section class="page-content">
    <header class="hero">
      <div>
        <h2>赛前配置</h2>
      </div>
      <div class="hero-actions">
        <button type="button" class="btn" :disabled="shellLoading || dictionaryLoading || initializing" @click="handleRefresh">刷新状态</button>
        <button v-if="setupStatus?.initialized" type="button" class="btn primary" @click="goToSummary">进入汇总页</button>
      </div>
    </header>

    <section v-if="pageMessage" class="message-bar" :class="pageMessage.type">
      {{ pageMessage.text }}
    </section>

    <div class="cards-grid">
      <article class="card">
        <span class="card-label">初始化状态</span>
        <strong>{{ setupStatus?.initialized ? '已完成' : '未完成' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">当前小组数</span>
        <strong>{{ setupStatus?.groupCount ?? '--' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">最终年份</span>
        <strong>{{ setupStatus?.finalYear ?? '--' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">沙盘版本</span>
        <strong>{{ setupStatus?.editionName ?? selectedEdition?.editionName ?? '--' }}</strong>
      </article>
      <article class="card">
        <span class="card-label">字典版本</span>
        <strong>{{ dictionaryCurrent?.dictionaryRevision ?? setupStatus?.dictionaryRevision ?? 0 }}</strong>
      </article>
    </div>

    <section v-if="!setupStatus?.initialized" class="panel-card setup-flow">
      <div class="steps-row">
        <span v-for="item in preSetupSteps" :key="item" class="step-pill">{{ item }}</span>
      </div>

      <div class="setup-grid">
        <section class="setup-section">
          <div class="panel-head">
            <strong>比赛基础信息</strong>
          </div>
          <div class="form-grid">
            <label class="field">
              <span>小组数量</span>
              <input
                v-model.number="groupCountDraft"
                type="number"
                min="1"
                max="10"
                step="1"
                inputmode="numeric"
                :class="{ invalid: hasFractionInput(groupCountDraft) }"
                :disabled="initializing"
                data-enter-confirm
                @keydown.enter="confirmInputOnEnter"
              >
            </label>
            <label class="field">
              <span>默认玩家账号</span>
              <input :value="accountPreviewText" readonly>
            </label>
            <label class="field">
              <span>最终年份</span>
              <input :value="setupStatus?.finalYear ?? '--'" readonly>
            </label>
            <label class="field">
              <span>共享基线</span>
              <input :value="setupStatus?.initialBaselineSubmitted ? '已提交' : '未提交'" readonly>
            </label>
          </div>
          <button type="button" class="btn link-btn" :disabled="!setupStatus?.initialized" @click="goToBaseline">进入初始基线页</button>
        </section>

        <section class="setup-section">
          <div class="panel-head">
            <strong>版本包</strong>
            <span>版本包决定字段结构、公式版本和流程规则，初始化完成后不可切换。</span>
          </div>
          <label class="field">
            <span>沙盘版本</span>
            <select v-model="editionCodeDraft" :disabled="initializing || availableEditions.length === 0">
              <option v-for="item in availableEditions" :key="item.editionCode" :value="item.editionCode">
                {{ item.editionName }}
              </option>
            </select>
          </label>
          <div class="meta-panel">
            <span>{{ currentEditionDescription }}</span>
            <span>{{ currentEditionTemplateText }}</span>
          </div>
        </section>
      </div>

      <section class="setup-section">
        <div class="panel-head">
          <div>
            <strong>业务显示字典</strong>
            <span>只修改页面显示名，不改变字段编码、公式、流程或数据库含义。</span>
          </div>
          <div class="toolbar-actions">
            <select v-model.number="selectedSchemeId" :disabled="dictionaryLoading || initializing">
              <option :value="0">版本默认名称</option>
              <option v-for="item in customSchemes" :key="item.id" :value="item.id">
                {{ item.schemeName }}
              </option>
            </select>
            <button type="button" class="btn" :disabled="dictionaryLoading || initializing" @click="restoreDraftDefault">恢复默认</button>
            <button type="button" class="btn" :disabled="dictionaryLoading || initializing || selectedSchemeId <= 0 || dictionaryDraftInvalidCount > 0" @click="updateSelectedScheme">更新方案</button>
            <button type="button" class="btn danger" :disabled="dictionaryLoading || initializing || selectedSchemeId <= 0" @click="deleteSelectedScheme">删除方案</button>
            <button type="button" class="btn" :disabled="dictionaryLoading || initializing || dictionaryDraftInvalidCount > 0" @click="saveDraftAsScheme">另存为方案</button>
          </div>
        </div>

        <DictionaryEditor
          :items="dictionaryDraftItems"
          :base-items="dictionaryCurrent?.items ?? []"
          :readonly="initializing"
          @update="updateDictionaryDraftItem"
        />
      </section>

      <section class="confirm-panel">
        <div>
          <strong>确认初始化</strong>
          <span>将生成小组、账号、年份状态，并保存当前显示名称快照。初始基线可在初始化后从本页入口进入提交。</span>
        </div>
        <button type="button" class="btn primary" :disabled="initializing || dictionaryDraftInvalidCount > 0" @click="handleInitialize">
          {{ initializing ? '初始化中...' : '确认初始化比赛' }}
        </button>
      </section>
    </section>

    <section v-else class="panel-card initialized-view">
      <div class="panel-head">
        <div>
          <strong>当前比赛配置</strong>
          
        </div>
        <div class="toolbar-actions">
          <button type="button" class="btn" @click="goToBaseline">初始基线</button>
          <button type="button" class="btn" @click="unlockCurrentDictionary">
            {{ currentDictionaryUnlocked ? '收起修改' : '解锁修改显示名称' }}
          </button>
        </div>
      </div>

      <div class="meta-grid">
        <div><span>版本包</span><strong>{{ setupStatus?.editionName }}</strong></div>
        <div><span>公式版本</span><strong>{{ setupStatus?.formulaVersion }}</strong></div>
        <div><span>流程规则</span><strong>{{ setupStatus?.processRuleVersion }}</strong></div>
        <div><span>字段模板</span><strong>{{ currentEditionTemplateText }}</strong></div>
      </div>

      <section class="setup-section">
        <div class="panel-head">
          <div>
            <strong>当前比赛显示名称</strong>
            
          </div>
          <div class="toolbar-actions">
            <select v-model.number="selectedSchemeId" :disabled="!currentDictionaryUnlocked || dictionaryLoading">
              <option :value="0">版本默认名称</option>
              <option v-for="item in customSchemes" :key="item.id" :value="item.id">
                {{ item.schemeName }}
              </option>
            </select>
            <button type="button" class="btn" :disabled="!currentDictionaryUnlocked || selectedSchemeId <= 0 || dictionaryDraftInvalidCount > 0" @click="applySchemeToCurrent">应用方案</button>
            <button type="button" class="btn" :disabled="!currentDictionaryUnlocked" @click="restoreCurrentDefault">恢复默认</button>
            <button type="button" class="btn" :disabled="!currentDictionaryUnlocked || selectedSchemeId <= 0 || dictionaryDraftInvalidCount > 0" @click="updateSelectedScheme">更新方案</button>
            <button type="button" class="btn danger" :disabled="!currentDictionaryUnlocked || selectedSchemeId <= 0" @click="deleteSelectedScheme">删除方案</button>
            <button type="button" class="btn" :disabled="!currentDictionaryUnlocked || dictionaryDraftInvalidCount > 0" @click="saveDraftAsScheme">另存为方案</button>
            <button type="button" class="btn primary" :disabled="!currentDictionaryUnlocked || savingDictionary || dictionaryDraftInvalidCount > 0 || dictionaryPendingChangeCount === 0" @click="saveCurrentDictionary">
              {{ currentDictionarySaveButtonText }}
            </button>
          </div>
        </div>
        <DictionaryEditor
          :items="dictionaryDraftItems"
          :base-items="dictionaryCurrent?.items ?? []"
          :readonly="!currentDictionaryUnlocked || savingDictionary"
          @update="updateDictionaryDraftItem"
        />
      </section>

      <section class="setup-section">
        <div class="panel-head">
          <strong>字段名称修改记录</strong>
          <span>只展示本场比赛和方案维护相关记录。</span>
        </div>
        <div class="log-list">
          <div v-for="item in changeLogItems" :key="item.id" class="log-item">
            <strong>{{ formatChangeType(item.changeType) }} · revision {{ item.revision }}</strong>
            <span>{{ item.changedSummary }} · {{ item.operatorName }} · {{ formatDateTime(item.operateTime) }}</span>
            <small v-if="item.reason">{{ item.reason }}</small>
          </div>
          <div v-if="changeLogItems.length === 0" class="empty-state">暂无修改记录。</div>
        </div>
      </section>
    </section>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'

import DictionaryEditor from '@/components/sandbox-game/admin/DictionaryEditor.vue'
import {
  applyAdminDictionarySchemeToCurrent,
  deleteAdminDictionaryScheme,
  getAdminDictionarySchemeDetail,
  restoreAdminDictionaryCurrentDefault,
  saveAdminDictionaryScheme,
  updateAdminDictionaryCurrent,
} from '@/api/sandbox-game/admin-dictionary'
import { initializeAdminGameWithDictionary } from '@/api/sandbox-game/admin-control'
import { useAdminShellStore, type PageMessage } from '@/stores/admin-shell'
import { useAuthStore } from '@/stores/auth'
import { useDictionaryStore } from '@/stores/dictionary'
import type { AdminDictionaryItem } from '@/types/sandbox-game-admin'
import { confirmInputOnEnter } from '@/utils/input-navigation'
import { hasFractionInput } from '@/utils/manual-integer'

const router = useRouter()
const shellStore = useAdminShellStore()
const authStore = useAuthStore()
const dictionaryStore = useDictionaryStore()
const { config, setupStatus, loading: shellLoading } = storeToRefs(shellStore)
const { current: dictionaryCurrent, schemes, changeLogs, loading: dictionaryLoading } = storeToRefs(dictionaryStore)

const groupCountDraft = ref(10)
const editionCodeDraft = ref('VIP_SERVICE_V1')
const selectedSchemeId = ref(0)
const dictionaryDraftItems = ref<AdminDictionaryItem[]>([])
const initializing = ref(false)
const savingDictionary = ref(false)
const currentDictionaryUnlocked = ref(false)
const pageMessage = ref<PageMessage | null>(null)

const preSetupSteps = ['基础信息', '版本包', '显示名称', '初始基线', '确认初始化']
const availableEditions = computed(() => setupStatus.value?.availableEditions ?? [])
const customSchemes = computed(() => (schemes.value?.list ?? []).filter((item) => item.id > 0))
const changeLogItems = computed(() => changeLogs.value?.list ?? [])

const accountPreviewText = computed(() => {
  const count = normalizeGroupCount(groupCountDraft.value)
  return count <= 1 ? 'group01' : `group01 ~ group${String(count).padStart(2, '0')}`
})

const selectedEdition = computed(() =>
  availableEditions.value.find((item) => item.editionCode === editionCodeDraft.value)
  ?? availableEditions.value.find((item) => item.defaultEdition)
  ?? null,
)

const currentEditionDescription = computed(() => {
  const edition = selectedEdition.value
  return edition ? `${edition.editionName} · ${edition.description}` : setupStatus.value?.editionName ?? '--'
})

const currentEditionTemplateText = computed(() => {
  const edition = selectedEdition.value
  if (edition) {
    return `${edition.operatingTemplateVersion} / ${edition.reportTemplateVersion} / ${edition.orderTemplateVersion}`
  }
  return [
    setupStatus.value?.operatingTemplateVersion,
    setupStatus.value?.reportTemplateVersion,
    setupStatus.value?.orderTemplateVersion,
  ].filter(Boolean).join(' / ') || '--'
})

const dictionaryDraftInvalidCount = computed(() => dictionaryDraftItems.value.filter((item) => !item.displayName.trim()).length)
const dictionaryPendingChangeCount = computed(() => countChangedDisplayNames(dictionaryCurrent.value?.items ?? [], dictionaryDraftItems.value))
const currentDictionarySaveButtonText = computed(() => {
  if (savingDictionary.value) {
    return '保存中...'
  }
  if (dictionaryDraftInvalidCount.value > 0) {
    return `还有 ${dictionaryDraftInvalidCount.value} 项为空`
  }
  if (dictionaryPendingChangeCount.value === 0) {
    return '暂无修改'
  }
  return `保存 ${dictionaryPendingChangeCount.value} 项修改`
})

onMounted(async () => {
  await bootstrapPage()
})

watch(
  () => setupStatus.value?.groupCount,
  () => syncGroupCountDraft(),
)

watch(
  () => [setupStatus.value?.editionCode, setupStatus.value?.availableEditions?.length],
  () => syncEditionCodeDraft(),
)

watch(
  () => editionCodeDraft.value,
  async (nextEdition) => {
    if (setupStatus.value?.initialized) {
      return
    }
    await loadDictionaryForEdition(nextEdition)
  },
)

watch(
  () => selectedSchemeId.value,
  async (nextSchemeId) => {
    if (!schemes.value || nextSchemeId === undefined) {
      return
    }
    if (nextSchemeId === 0) {
      restoreDraftDefault()
      return
    }
    await loadSelectedSchemeDetail(nextSchemeId)
  },
)

async function loadSelectedSchemeDetail(schemeId: number) {
  pageMessage.value = null
  try {
    const detail = await getAdminDictionarySchemeDetail(schemeId)
    if (detail.editionCode !== editionCodeDraft.value) {
      pageMessage.value = { type: 'error', text: '该字典方案不属于当前沙盘版本。' }
      selectedSchemeId.value = 0
      restoreDraftDefault()
      return
    }
    dictionaryDraftItems.value = mergeDictionaryItems(dictionaryCurrent.value?.items ?? dictionaryDraftItems.value, detail.items)
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '读取字典方案详情失败')
  }
}

async function bootstrapPage() {
  try {
    await shellStore.bootstrap()
    syncGroupCountDraft()
    syncEditionCodeDraft()
    await Promise.all([
      loadDictionaryForEdition(setupStatus.value?.initialized ? setupStatus.value.editionCode : editionCodeDraft.value),
      dictionaryStore.loadChangeLogs(1, 20),
    ])
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '读取赛前配置失败')
  }
}

async function loadDictionaryForEdition(editionCode?: string | null) {
  const [current] = await Promise.all([
    dictionaryStore.loadCurrent(editionCode, { silent: true }),
    dictionaryStore.loadSchemes(editionCode),
  ])
  selectedSchemeId.value = 0
  dictionaryDraftItems.value = cloneDictionaryItems(current.items)
}

async function handleRefresh() {
  pageMessage.value = null
  await bootstrapPage()
}

async function handleInitialize() {
  if (hasFractionInput(groupCountDraft.value)) {
    pageMessage.value = { type: 'error', text: '小组数量必须填写整数。' }
    return
  }
  if (dictionaryDraftItems.value.some((item) => !item.displayName.trim())) {
    pageMessage.value = { type: 'error', text: '显示名称不能为空。' }
    return
  }
  initializing.value = true
  pageMessage.value = null
  try {
    const result = await initializeAdminGameWithDictionary(
      normalizeGroupCount(groupCountDraft.value),
      editionCodeDraft.value,
      toDictionaryInputs(dictionaryDraftItems.value),
      selectedSchemeId.value > 0 ? selectedSchemeId.value : null,
    )
    await shellStore.refreshAll({ silent: true })
    await authStore.refreshCurrentUser()
    await loadDictionaryForEdition(result.editionCode)
    pageMessage.value = {
      type: 'success',
      text: `比赛初始化完成，已生成 ${result.createdGroupCount} 个小组、${result.createdAccountCount} 个玩家账号和 ${result.createdYearStateCount} 条年份状态。`,
    }
    await router.replace(authStore.resolveDefaultRoute())
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '初始化比赛失败')
  } finally {
    initializing.value = false
  }
}

function updateDictionaryDraftItem(itemCode: string, displayName: string) {
  const matched = dictionaryDraftItems.value.find((item) => item.itemCode === itemCode)
  if (matched) {
    matched.displayName = displayName
  }
}

function restoreDraftDefault() {
  dictionaryDraftItems.value = dictionaryDraftItems.value.map((item) => ({
    ...item,
    displayName: item.defaultName,
  }))
}

async function saveDraftAsScheme() {
  if (dictionaryDraftInvalidCount.value > 0) {
    pageMessage.value = { type: 'error', text: '显示名称不能为空，请先处理红色标记字段。' }
    return
  }
  const schemeName = window.prompt('请输入字典方案名称', `${selectedEdition.value?.editionName ?? '本场比赛'}显示名称`)
  if (!schemeName?.trim()) {
    return
  }
  try {
    await saveAdminDictionaryScheme({
      editionCode: editionCodeDraft.value,
      schemeName: schemeName.trim(),
      description: '管理员在赛前配置页保存的显示名称方案',
      items: toDictionaryInputs(dictionaryDraftItems.value),
    })
    await dictionaryStore.loadSchemes(editionCodeDraft.value)
    await dictionaryStore.loadChangeLogs(1, 20)
    pageMessage.value = { type: 'success', text: '字典方案已保存。' }
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '保存字典方案失败')
  }
}

async function updateSelectedScheme() {
  const scheme = customSchemes.value.find((item) => item.id === selectedSchemeId.value)
  if (!scheme) {
    return
  }
  if (dictionaryDraftInvalidCount.value > 0) {
    pageMessage.value = { type: 'error', text: '显示名称不能为空，请先处理红色标记字段。' }
    return
  }
  if (!window.confirm(`确认用当前编辑内容覆盖字典方案「${scheme.schemeName}」吗？`)) {
    return
  }
  savingDictionary.value = true
  pageMessage.value = null
  try {
    await saveAdminDictionaryScheme({
      schemeId: scheme.id,
      editionCode: scheme.editionCode,
      schemeName: scheme.schemeName,
      description: scheme.description,
      items: toDictionaryInputs(dictionaryDraftItems.value),
    })
    await dictionaryStore.loadSchemes(editionCodeDraft.value)
    await dictionaryStore.loadChangeLogs(1, 20)
    pageMessage.value = { type: 'success', text: '字典方案已更新。' }
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '更新字典方案失败')
  } finally {
    savingDictionary.value = false
  }
}

async function deleteSelectedScheme() {
  const scheme = customSchemes.value.find((item) => item.id === selectedSchemeId.value)
  if (!scheme) {
    return
  }
  if (!window.confirm(`确认删除字典方案「${scheme.schemeName}」吗？删除方案不会影响当前比赛显示名称。`)) {
    return
  }
  savingDictionary.value = true
  pageMessage.value = null
  try {
    await deleteAdminDictionaryScheme(scheme.id)
    selectedSchemeId.value = 0
    await dictionaryStore.loadSchemes(editionCodeDraft.value)
    await dictionaryStore.loadChangeLogs(1, 20)
    restoreDraftDefault()
    pageMessage.value = { type: 'success', text: '字典方案已删除，当前编辑内容已恢复为版本默认名称。' }
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '删除字典方案失败')
  } finally {
    savingDictionary.value = false
  }
}

function unlockCurrentDictionary() {
  currentDictionaryUnlocked.value = !currentDictionaryUnlocked.value
  if (currentDictionaryUnlocked.value) {
    dictionaryDraftItems.value = cloneDictionaryItems(dictionaryCurrent.value?.items ?? [])
  }
}

async function saveCurrentDictionary() {
  if (dictionaryDraftInvalidCount.value > 0) {
    pageMessage.value = { type: 'error', text: '显示名称不能为空，请先处理红色标记字段。' }
    return
  }
  if (dictionaryPendingChangeCount.value === 0) {
    pageMessage.value = { type: 'success', text: '当前没有需要保存的显示名称修改。' }
    return
  }
  const preview = buildChangedDisplayNamePreview(dictionaryCurrent.value?.items ?? [], dictionaryDraftItems.value)
  if (!window.confirm(`确认保存当前比赛显示名称修改吗？${preview}`)) {
    return
  }
  savingDictionary.value = true
  pageMessage.value = null
  try {
    const result = await updateAdminDictionaryCurrent({
      items: toDictionaryInputs(dictionaryDraftItems.value),
      reason: '管理员在赛前配置页修改当前比赛显示名称',
    })
    dictionaryStore.applyCurrent(result)
    await dictionaryStore.loadChangeLogs(1, 20)
    currentDictionaryUnlocked.value = false
    pageMessage.value = { type: 'success', text: '当前比赛显示名称已更新。' }
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '保存当前比赛显示名称失败')
  } finally {
    savingDictionary.value = false
  }
}

async function applySchemeToCurrent() {
  if (selectedSchemeId.value <= 0) {
    return
  }
  if (dictionaryDraftInvalidCount.value > 0) {
    pageMessage.value = { type: 'error', text: '显示名称不能为空，请先处理红色标记字段。' }
    return
  }
  const scheme = customSchemes.value.find((item) => item.id === selectedSchemeId.value)
  let schemeItems = dictionaryDraftItems.value
  try {
    const detail = await getAdminDictionarySchemeDetail(selectedSchemeId.value)
    schemeItems = mergeDictionaryItems(dictionaryCurrent.value?.items ?? dictionaryDraftItems.value, detail.items)
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '读取字典方案详情失败')
    return
  }
  const preview = buildChangedDisplayNamePreview(dictionaryCurrent.value?.items ?? [], schemeItems)
  if (!window.confirm(`确认应用字典方案「${scheme?.schemeName ?? selectedSchemeId.value}」到当前比赛吗？${preview}`)) {
    return
  }
  savingDictionary.value = true
  try {
    const result = await applyAdminDictionarySchemeToCurrent({ schemeId: selectedSchemeId.value })
    dictionaryStore.applyCurrent(result)
    dictionaryDraftItems.value = cloneDictionaryItems(result.items)
    await dictionaryStore.loadChangeLogs(1, 20)
    pageMessage.value = { type: 'success', text: '已应用字典方案到当前比赛。' }
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '应用字典方案失败')
  } finally {
    savingDictionary.value = false
  }
}

async function restoreCurrentDefault() {
  const currentItems = dictionaryCurrent.value?.items ?? []
  const defaultItems = currentItems.map((item) => ({
    ...item,
    displayName: item.defaultName,
  }))
  const preview = buildChangedDisplayNamePreview(currentItems, defaultItems)
  if (!window.confirm(`确认恢复当前版本默认显示名称吗？${preview}`)) {
    return
  }
  savingDictionary.value = true
  try {
    const result = await restoreAdminDictionaryCurrentDefault()
    dictionaryStore.applyCurrent(result)
    dictionaryDraftItems.value = cloneDictionaryItems(result.items)
    await dictionaryStore.loadChangeLogs(1, 20)
    pageMessage.value = { type: 'success', text: '已恢复当前比赛默认显示名称。' }
  } catch (error) {
    pageMessage.value = toErrorMessage(error, '恢复默认显示名称失败')
  } finally {
    savingDictionary.value = false
  }
}

async function goToSummary() {
  await router.replace('/sandbox-game/admin/summary')
}

async function goToBaseline() {
  await router.push('/sandbox-game/admin/baseline')
}

function syncGroupCountDraft() {
  const nextGroupCount = setupStatus.value?.groupCount
  groupCountDraft.value = nextGroupCount && nextGroupCount > 0 ? nextGroupCount : 10
}

function syncEditionCodeDraft() {
  const statusEditionCode = setupStatus.value?.editionCode
  if (statusEditionCode) {
    editionCodeDraft.value = statusEditionCode
    return
  }
  const defaultEdition = availableEditions.value.find((item) => item.defaultEdition) ?? availableEditions.value[0]
  editionCodeDraft.value = defaultEdition?.editionCode ?? 'VIP_SERVICE_V1'
}

function normalizeGroupCount(value: number) {
  const rounded = Number.isFinite(value) ? Math.round(value) : 10
  return Math.min(10, Math.max(1, rounded))
}

function cloneDictionaryItems(items: AdminDictionaryItem[]) {
  return items.map((item) => ({ ...item }))
}

function mergeDictionaryItems(baseItems: AdminDictionaryItem[], overrideItems: AdminDictionaryItem[]) {
  const overrideMap = new Map(overrideItems.map((item) => [item.itemCode, item.displayName]))
  return cloneDictionaryItems(baseItems).map((item) => ({
    ...item,
    displayName: overrideMap.get(item.itemCode) ?? item.displayName,
  }))
}

function buildChangedDisplayNamePreview(beforeItems: AdminDictionaryItem[], afterItems: AdminDictionaryItem[]) {
  const beforeMap = new Map(beforeItems.map((item) => [item.itemCode, item]))
  const changedItems = afterItems
    .map((item) => ({
      item,
      beforeName: beforeMap.get(item.itemCode)?.displayName ?? item.defaultName,
    }))
    .filter(({ item, beforeName }) => beforeName.trim() !== item.displayName.trim())
  const categoryLabels: Record<string, string> = {
    MARKET: '市场',
    ORDER_TYPE: '订单类型',
    OPERATING: '经营页',
    REPORT: '财报页',
    BASELINE: '初始基线',
  }
  const grouped = new Map<string, typeof changedItems>()
  for (const changedItem of changedItems) {
    const key = changedItem.item.itemCategory || 'OTHER'
    grouped.set(key, [...(grouped.get(key) ?? []), changedItem])
  }
  const detailLines = Array.from(grouped.entries()).flatMap(([category, items]) => {
    const header = `\n【${categoryLabels[category] ?? category}】`
    const lines = items.slice(0, 6).map(({ item, beforeName }) => `\n- ${item.defaultName}: ${beforeName} -> ${item.displayName.trim()}`)
    const moreText = items.length > 6 ? `\n- 其余 ${items.length - 6} 项略` : ''
    return [header, ...lines, moreText]
  })
  return `预计会覆盖 ${changedItems.length} 个显示名称。${detailLines.join('')}`
}

function countChangedDisplayNames(beforeItems: AdminDictionaryItem[], afterItems: AdminDictionaryItem[]) {
  const beforeMap = new Map(beforeItems.map((item) => [item.itemCode, item]))
  return afterItems.filter((item) => {
    const beforeName = beforeMap.get(item.itemCode)?.displayName ?? item.defaultName
    return beforeName.trim() !== item.displayName.trim()
  }).length
}

function toDictionaryInputs(items: AdminDictionaryItem[]) {
  return items.map((item) => ({
    itemCode: item.itemCode,
    displayName: item.displayName.trim(),
  }))
}

function formatDateTime(value?: string | null) {
  if (!value) {
    return '--'
  }
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString('zh-CN', { hour12: false })
}

function formatChangeType(value: string) {
  const map: Record<string, string> = {
    INITIALIZE: '初始化',
    UPDATE_CURRENT: '修改当前名称',
    APPLY_SCHEME: '应用方案',
    RESTORE_DEFAULT: '恢复默认',
    SAVE_SCHEME: '保存方案',
    DELETE_SCHEME: '删除方案',
  }
  return map[value] ?? value
}

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
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

.hero p,
.panel-head span,
.confirm-panel span,
.meta-panel span,
.log-item span,
.log-item small {
  color: var(--muted);
  line-height: 1.6;
}

.hero-actions,
.toolbar-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  align-items: center;
}

.btn {
  border: 1px solid var(--line);
  background: #ffffff;
  border-radius: 12px;
  padding: 10px 14px;
  cursor: pointer;
}

.btn.primary {
  background: var(--accent);
  color: #ffffff;
  border-color: var(--accent);
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.link-btn {
  margin-top: 14px;
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

.cards-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 12px;
}

.card,
.panel-card,
.setup-section,
.confirm-panel {
  border: 1px solid var(--line);
  border-radius: 16px;
  background: #ffffff;
}

.card {
  display: grid;
  gap: 8px;
  padding: 18px;
}

.card-label {
  color: var(--muted);
  font-size: 13px;
}

.card strong {
  font-size: 22px;
}

.panel-card,
.setup-section,
.confirm-panel {
  padding: 16px;
}

.setup-flow,
.initialized-view {
  display: grid;
  gap: 16px;
}

.steps-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.step-pill {
  padding: 7px 10px;
  border-radius: 999px;
  background: #eef4ff;
  color: var(--accent);
  border: 1px solid #ccd9f3;
  font-size: 12px;
  font-weight: 700;
}

.setup-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(320px, 0.9fr);
  gap: 16px;
}

.panel-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.panel-head strong {
  display: block;
  margin-bottom: 4px;
  font-size: 18px;
}

.form-grid,
.meta-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.field {
  display: grid;
  gap: 8px;
}

.field span {
  color: var(--muted);
  font-size: 13px;
}

.field input,
.field select,
.toolbar-actions select {
  width: 100%;
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid var(--line);
  background: #fbfcfe;
}

.field input.invalid {
  border-color: var(--danger);
  background: #fff6f6;
  color: var(--danger);
}

.meta-panel {
  display: grid;
  gap: 8px;
  margin-top: 14px;
  padding: 12px;
  border-radius: 12px;
  background: #f7f9fc;
}

.meta-grid div {
  display: grid;
  gap: 6px;
  padding: 12px;
  border-radius: 12px;
  background: #f7f9fc;
}

.meta-grid span {
  color: var(--muted);
  font-size: 12px;
}

.meta-grid strong {
  font-size: 14px;
}

.confirm-panel {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.log-list {
  display: grid;
  gap: 10px;
}

.log-item {
  display: grid;
  gap: 4px;
  padding: 12px;
  border-radius: 12px;
  background: #f7f9fc;
}

.log-item small {
  font-size: 12px;
}

.empty-state {
  color: var(--muted);
  text-align: center;
  padding: 16px;
}

@media (max-width: 1180px) {
  .cards-grid,
  .setup-grid,
  .form-grid,
  .meta-grid {
    grid-template-columns: 1fr;
  }

  .panel-head,
  .confirm-panel,
  .hero {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
