<template>
  <div class="page-mode-switch">
    <button
      type="button"
      class="mode-btn"
      :class="{ active: activeMode === 'operating' }"
      :disabled="activeMode === 'operating'"
      @click="$emit('operating')"
    >
      经营页
    </button>
    <button
      type="button"
      class="mode-btn"
      :class="{ active: activeMode === 'report' }"
      :disabled="activeMode === 'report' || !reportEnabled"
      @click="$emit('report')"
    >
      财报页
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  activeMode: 'operating' | 'report'
  reportEnabled: boolean
}>()

defineEmits<{
  (event: 'operating'): void
  (event: 'report'): void
}>()
</script>

<style scoped>
.page-mode-switch {
  display: inline-flex;
  gap: 8px;
  padding: 6px;
  border-radius: 14px;
  background: #eef2f6;
}

.mode-btn {
  border: 1px solid transparent;
  border-radius: 10px;
  background: transparent;
  padding: 10px 16px;
  color: var(--muted);
}

.mode-btn.active {
  background: #ffffff;
  color: var(--text);
  border-color: var(--line-strong);
  font-weight: 700;
}

.mode-btn:disabled {
  opacity: 0.6;
}
</style>