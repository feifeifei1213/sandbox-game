<template>
  <section class="dictionary-editor">
    <div class="dictionary-summary">
      <div class="summary-metrics">
        <span class="summary-pill changed">已修改 {{ changedCount }} 项</span>
        <span class="summary-pill" :class="{ invalid: invalidCount > 0 }">空值 {{ invalidCount }} 项</span>
        <span class="summary-pill muted">共 {{ editableCount }} 项可维护</span>
      </div>
      <div class="summary-tools">
        <label class="dictionary-search">
          <span>搜索</span>
          <input v-model.trim="searchKeyword" type="search" placeholder="字段名、当前名或修改后名称">
        </label>
        <button
          type="button"
          class="filter-toggle"
          :class="{ active: changedOnly }"
          @click="changedOnly = !changedOnly"
        >
          只看已修改
        </button>
      </div>
    </div>

    <div class="dictionary-accordion">
      <section
        v-for="group in visibleGroups"
        :key="group.category"
        class="dictionary-category"
        :class="{ changed: group.changedCount > 0, invalid: group.invalidCount > 0 }"
      >
        <button type="button" class="category-header" @click="toggleCategory(group.category)">
          <span class="expand-icon">{{ expandedCategories[group.category] ? '收起' : '展开' }}</span>
          <strong>{{ group.label }}</strong>
          <span class="category-meta">{{ group.items.length }} 项</span>
          <span v-if="group.changedCount > 0" class="category-badge changed">{{ group.changedCount }} 项修改</span>
          <span v-if="group.invalidCount > 0" class="category-badge invalid">{{ group.invalidCount }} 项空值</span>
        </button>

        <div v-if="expandedCategories[group.category]" class="dictionary-table">
          <div class="table-head">
            <span>默认名称</span>
            <span>当前名称</span>
            <span>修改后名称</span>
            <span>操作</span>
          </div>

          <div
            v-for="row in group.items"
            :key="row.item.itemCode"
            class="dictionary-row"
            :class="{ changed: row.changed, invalid: row.invalid, readonly: readonly || !row.item.editable }"
          >
            <span class="row-name">{{ row.item.defaultName }}</span>
            <span class="row-current">{{ row.currentName }}</span>
            <label class="row-input">
              <input
                :value="row.item.displayName"
                :readonly="readonly || !row.item.editable"
                :aria-label="`${row.item.defaultName} 修改后名称`"
                @input="handleInput(row.item.itemCode, $event)"
              >
            </label>
            <span class="row-actions">
              <button
                v-if="row.changed"
                type="button"
                class="text-btn"
                :disabled="readonly || !row.item.editable"
                @click="emit('update', row.item.itemCode, row.currentName)"
              >
                撤销
              </button>
              <button
                type="button"
                class="text-btn"
                :disabled="readonly || !row.item.editable || normalizeName(row.item.displayName) === normalizeName(row.item.defaultName)"
                @click="emit('update', row.item.itemCode, row.item.defaultName)"
              >
                默认
              </button>
            </span>
          </div>
        </div>

        <div v-if="expandedCategories[group.category] && group.items.length > 0" class="category-footer">
          <span>本类修改只影响页面展示名称，不改变字段编码和计算规则。</span>
          <button type="button" class="text-btn" :disabled="readonly" @click="restoreCategoryDefault(group.category)">
            恢复本类默认
          </button>
        </div>
      </section>

      <div v-if="visibleGroups.length === 0" class="dictionary-empty">
        没有匹配的字段。
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'

import type { AdminDictionaryItem } from '@/types/sandbox-game-admin'

type DictionaryRow = {
  item: AdminDictionaryItem
  currentName: string
  changed: boolean
  invalid: boolean
}

type DictionaryGroup = {
  category: string
  label: string
  items: DictionaryRow[]
  changedCount: number
  invalidCount: number
}

const props = withDefaults(defineProps<{
  items: AdminDictionaryItem[]
  baseItems?: AdminDictionaryItem[]
  readonly?: boolean
}>(), {
  baseItems: () => [],
  readonly: false,
})

const emit = defineEmits<{
  update: [itemCode: string, displayName: string]
}>()

const categoryLabels: Record<string, string> = {
  MARKET: '市场',
  ORDER_TYPE: '订单类型',
  OPERATING: '经营页',
  REPORT: '财报页',
  BASELINE: '初始基线',
}

const categoryOrder = ['MARKET', 'ORDER_TYPE', 'OPERATING', 'REPORT', 'BASELINE']
const searchKeyword = ref('')
const changedOnly = ref(false)
const expandedCategories = ref<Record<string, boolean>>({})

const baseItemMap = computed(() => new Map(props.baseItems.map((item) => [item.itemCode, item])))

const groupedItems = computed<DictionaryGroup[]>(() => {
  const groups = new Map<string, DictionaryRow[]>()
  for (const item of props.items) {
    const baseItem = baseItemMap.value.get(item.itemCode)
    const currentName = baseItem?.displayName ?? item.defaultName
    const row: DictionaryRow = {
      item,
      currentName,
      changed: normalizeName(item.displayName) !== normalizeName(currentName),
      invalid: !normalizeName(item.displayName),
    }
    const key = item.itemCategory || 'OTHER'
    groups.set(key, [...(groups.get(key) ?? []), row])
  }

  return Array.from(groups.entries())
    .map(([category, items]) => {
      const sortedItems = [...items].sort((a, b) => a.item.displayOrder - b.item.displayOrder)
      return {
        category,
        label: categoryLabels[category] ?? category,
        items: sortedItems,
        changedCount: sortedItems.filter((row) => row.changed).length,
        invalidCount: sortedItems.filter((row) => row.invalid).length,
      }
    })
    .sort((a, b) => categorySortIndex(a.category) - categorySortIndex(b.category))
})

const visibleGroups = computed<DictionaryGroup[]>(() => {
  const keyword = searchKeyword.value.toLowerCase()
  return groupedItems.value
    .map((group) => {
      const filteredItems = group.items.filter((row) => {
        const matchedKeyword = !keyword
          || row.item.itemCode.toLowerCase().includes(keyword)
          || row.item.defaultName.toLowerCase().includes(keyword)
          || row.currentName.toLowerCase().includes(keyword)
          || row.item.displayName.toLowerCase().includes(keyword)
        const matchedChanged = !changedOnly.value || row.changed
        return matchedKeyword && matchedChanged
      })
      return {
        ...group,
        items: filteredItems,
        changedCount: filteredItems.filter((row) => row.changed).length,
        invalidCount: filteredItems.filter((row) => row.invalid).length,
      }
    })
    .filter((group) => group.items.length > 0)
})

const changedCount = computed(() => groupedItems.value.reduce((sum, group) => sum + group.changedCount, 0))
const invalidCount = computed(() => groupedItems.value.reduce((sum, group) => sum + group.invalidCount, 0))
const editableCount = computed(() => props.items.filter((item) => item.editable).length)

watch(
  groupedItems,
  (groups) => {
    const next: Record<string, boolean> = { ...expandedCategories.value }
    const hasExpanded = Object.values(next).some(Boolean)
    for (const group of groups) {
      if (next[group.category] === undefined) {
        next[group.category] = group.changedCount > 0 || group.invalidCount > 0 || (!hasExpanded && group.category === groups[0]?.category)
      }
    }
    expandedCategories.value = next
  },
  { immediate: true },
)

function toggleCategory(category: string) {
  expandedCategories.value = {
    ...expandedCategories.value,
    [category]: !expandedCategories.value[category],
  }
}

function handleInput(itemCode: string, event: Event) {
  emit('update', itemCode, (event.target as HTMLInputElement).value)
}

function restoreCategoryDefault(category: string) {
  const group = groupedItems.value.find((item) => item.category === category)
  if (!group) {
    return
  }
  for (const row of group.items) {
    if (row.item.editable) {
      emit('update', row.item.itemCode, row.item.defaultName)
    }
  }
}

function categorySortIndex(category: string) {
  const index = categoryOrder.indexOf(category)
  return index >= 0 ? index : categoryOrder.length + 1
}

function normalizeName(value?: string | null) {
  return (value ?? '').trim()
}
</script>

<style scoped>
.dictionary-editor {
  display: grid;
  gap: 12px;
}

.dictionary-summary {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  padding: 12px;
  border: 1px solid #d9e3f0;
  border-radius: 8px;
  background: #f8fbff;
}

.summary-metrics,
.summary-tools,
.row-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: center;
}

.summary-pill,
.category-badge {
  display: inline-flex;
  align-items: center;
  min-height: 28px;
  padding: 4px 9px;
  border-radius: 999px;
  border: 1px solid #d9e3f0;
  background: #ffffff;
  color: #475569;
  font-size: 12px;
  font-weight: 700;
}

.summary-pill.changed,
.category-badge.changed {
  background: #fff7db;
  border-color: #f2d475;
  color: #8a5a00;
}

.summary-pill.invalid,
.category-badge.invalid {
  background: #fff1f1;
  border-color: #efb8b8;
  color: #b42318;
}

.summary-pill.muted {
  color: #64748b;
}

.dictionary-search {
  display: grid;
  gap: 4px;
}

.dictionary-search span {
  color: #64748b;
  font-size: 12px;
}

.dictionary-search input {
  width: 260px;
  max-width: 100%;
  border: 1px solid #ccd6e5;
  border-radius: 8px;
  padding: 8px 10px;
  background: #ffffff;
}

.filter-toggle,
.text-btn {
  border: 1px solid #ccd6e5;
  border-radius: 8px;
  background: #ffffff;
  color: #334155;
  cursor: pointer;
}

.filter-toggle {
  padding: 9px 12px;
}

.filter-toggle.active {
  background: #163b68;
  border-color: #163b68;
  color: #ffffff;
}

.dictionary-accordion {
  display: grid;
  gap: 10px;
}

.dictionary-category {
  overflow: hidden;
  border: 1px solid #d9e3f0;
  border-radius: 8px;
  background: #ffffff;
}

.dictionary-category.changed {
  border-color: #efd37d;
}

.dictionary-category.invalid {
  border-color: #efb8b8;
}

.category-header {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 10px;
  border: 0;
  border-left: 4px solid transparent;
  background: #ffffff;
  padding: 12px;
  text-align: left;
  cursor: pointer;
}

.dictionary-category.changed .category-header {
  border-left-color: #d79b00;
  background: #fffaf0;
}

.dictionary-category.invalid .category-header {
  border-left-color: #d92d20;
  background: #fff7f7;
}

.category-header strong {
  font-size: 15px;
}

.expand-icon,
.category-meta {
  color: #64748b;
  font-size: 12px;
}

.expand-icon {
  width: 34px;
}

.dictionary-table {
  display: grid;
  border-top: 1px solid #e6edf5;
}

.table-head,
.dictionary-row {
  display: grid;
  grid-template-columns: minmax(150px, 1fr) minmax(150px, 1fr) minmax(220px, 1.5fr) minmax(96px, 0.5fr);
  gap: 10px;
  align-items: center;
}

.table-head {
  padding: 10px 12px;
  background: #f4f7fb;
  color: #64748b;
  font-size: 12px;
  font-weight: 700;
}

.dictionary-row {
  padding: 10px 12px;
  border-top: 1px solid #edf2f7;
}

.dictionary-row.changed {
  background: #fffaf0;
  box-shadow: inset 4px 0 0 #d79b00;
}

.dictionary-row.invalid {
  background: #fff7f7;
  box-shadow: inset 4px 0 0 #d92d20;
}

.dictionary-row.readonly {
  color: #64748b;
}

.row-name,
.row-current {
  color: #334155;
}

.row-input input {
  width: 100%;
  min-width: 0;
  border: 1px solid #ccd6e5;
  border-radius: 8px;
  padding: 8px 10px;
  background: #ffffff;
}

.dictionary-row.changed .row-input input {
  border-color: #d79b00;
}

.dictionary-row.invalid .row-input input {
  border-color: #d92d20;
  color: #b42318;
}

.row-input input[readonly] {
  background: #f1f5f9;
  color: #64748b;
}

.text-btn {
  padding: 7px 9px;
  font-size: 12px;
}

.text-btn:disabled,
.filter-toggle:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.category-footer {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  padding: 10px 12px;
  border-top: 1px solid #edf2f7;
  color: #64748b;
  font-size: 12px;
}

.dictionary-empty {
  padding: 18px;
  border: 1px dashed #ccd6e5;
  border-radius: 8px;
  color: #64748b;
  text-align: center;
}

@media (max-width: 1120px) {
  .dictionary-summary,
  .category-footer {
    align-items: flex-start;
    flex-direction: column;
  }

  .dictionary-search input {
    width: 100%;
  }

  .summary-tools {
    width: 100%;
  }

  .dictionary-search {
    width: 100%;
  }

  .table-head {
    display: none;
  }

  .dictionary-row {
    grid-template-columns: 1fr;
    gap: 7px;
  }

  .row-name::before {
    content: '默认名称：';
    color: #64748b;
  }

  .row-current::before {
    content: '当前名称：';
    color: #64748b;
  }
}
</style>
