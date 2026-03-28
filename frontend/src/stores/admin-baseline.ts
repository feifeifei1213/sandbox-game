import { defineStore } from 'pinia'
import { ref } from 'vue'

import {
  getAdminInitialBaseline,
  submitAdminInitialBaseline,
} from '@/api/sandbox-game/admin-control'
import type {
  BaselinePayload,
  InitialBaselineViewResult,
} from '@/types/sandbox-game-admin'
import { cloneBaselinePayload, createEmptyBaselinePayload } from '@/types/sandbox-game-admin'
import type { PageMessage } from '@/stores/admin-shell'

export const useAdminBaselineStore = defineStore('sandbox-admin-baseline', () => {
  const view = ref<InitialBaselineViewResult | null>(null)
  const draftPayload = ref<BaselinePayload>(createEmptyBaselinePayload())
  const loading = ref(false)
  const submitting = ref(false)
  const dirty = ref(false)
  const pageMessage = ref<PageMessage | null>(null)

  async function bootstrap() {
    loading.value = true
    pageMessage.value = null
    try {
      const result = await getAdminInitialBaseline()
      view.value = result
      draftPayload.value = cloneBaselinePayload(result.baselinePayload)
      dirty.value = false
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '读取初始基线失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  function updateField(key: keyof BaselinePayload, value: number) {
    draftPayload.value = {
      ...draftPayload.value,
      [key]: Number.isFinite(value) ? value : 0,
    }
    dirty.value = true
  }

  async function submit() {
    if (!view.value?.editable) {
      return
    }

    submitting.value = true
    try {
      const result = await submitAdminInitialBaseline({
        baselinePayload: cloneBaselinePayload(draftPayload.value),
      })
      const latestView = await getAdminInitialBaseline()
      view.value = latestView
      draftPayload.value = cloneBaselinePayload(latestView.baselinePayload)
      dirty.value = false
      pageMessage.value = {
        type: 'success',
        text: `初始基线已提交，已应用到 ${result.appliedGroupCount} 个小组。`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '提交初始基线失败')
      throw error
    } finally {
      submitting.value = false
    }
  }

  return {
    view,
    draftPayload,
    loading,
    submitting,
    dirty,
    pageMessage,
    bootstrap,
    updateField,
    submit,
  }
})

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}
