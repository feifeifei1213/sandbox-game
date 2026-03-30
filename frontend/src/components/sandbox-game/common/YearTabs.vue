<template>
  <div class="year-tabs">
    <button
      v-for="item in tabs"
      :key="item.yearNo"
      type="button"
      class="year-tab"
      :class="{
        active: item.yearNo === activeYear,
        locked: !item.canEnter,
        completed: item.tabStatus === 'COMPLETED',
        bankrupt: item.tabStatus === 'BANKRUPT_READONLY',
      }"
      :disabled="!item.canEnter"
      @click="$emit('select', item.yearNo)"
    >
      <span class="label">{{ item.label }}</span>
      <span class="status">{{ formatYearTabStatus(item.tabStatus) }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import type { YearTabItem } from '@/types/sandbox-game'
import { formatYearTabStatus } from '@/utils/sandbox-game-display'

defineProps<{
  tabs: YearTabItem[]
  activeYear: number
}>()

defineEmits<{
  (event: 'select', yearNo: number): void
}>()
</script>

<style scoped>
.year-tabs {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding: 2px 0;
}

.year-tab {
  min-width: 102px;
  border: 1px solid var(--line-strong);
  border-radius: 12px;
  background: #ffffff;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.year-tab.active {
  background: var(--accent);
  color: #ffffff;
  border-color: var(--accent);
}

.year-tab.locked {
  background: var(--readonly-bg);
  color: var(--muted);
}

.year-tab.completed {
  border-color: #7fbf9b;
}

.year-tab.bankrupt {
  border-color: #e29a9a;
  color: var(--danger);
}

.label {
  font-weight: 700;
}

.status {
  font-size: 12px;
  opacity: 0.82;
}
</style>
