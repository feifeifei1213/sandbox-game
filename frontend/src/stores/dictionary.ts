import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import {
  getAdminDictionaryCurrent,
  getAdminDictionaryRevision,
  listAdminDictionarySchemes,
  pageAdminDictionaryChangeLogs,
} from '@/api/sandbox-game/admin-dictionary'
import type {
  AdminDictionaryChangeLogPageResult,
  AdminDictionaryCurrentResult,
  AdminDictionaryItem,
  AdminDictionarySchemeListResult,
} from '@/types/sandbox-game-admin'

const SYNC_INTERVAL_MS = 10_000

export const useDictionaryStore = defineStore('sandbox-dictionary', () => {
  const current = ref<AdminDictionaryCurrentResult | null>(null)
  const schemes = ref<AdminDictionarySchemeListResult | null>(null)
  const changeLogs = ref<AdminDictionaryChangeLogPageResult | null>(null)
  const loading = ref(false)
  const syncing = ref(false)
  const lastEditionCode = ref<string | null>(null)
  const timerId = ref<number | null>(null)

  const itemMap = computed(() => {
    const map = new Map<string, AdminDictionaryItem>()
    for (const item of current.value?.items ?? []) {
      map.set(item.itemCode, item)
    }
    return map
  })

  function displayName(itemCode: string, fallback: string) {
    return itemMap.value.get(itemCode)?.displayName || fallback
  }

  function marketName(code: string, fallback: string) {
    return displayName(`market.${code}`, fallback)
  }

  function orderTypeName(code: string, fallback: string) {
    return displayName(`orderType.${code}`, fallback)
  }

  async function loadCurrent(editionCode?: string | null, options?: { silent?: boolean }) {
    if (!options?.silent) {
      loading.value = true
    }
    try {
      const result = await getAdminDictionaryCurrent(editionCode)
      current.value = result
      lastEditionCode.value = result.editionCode
      return result
    } finally {
      if (!options?.silent) {
        loading.value = false
      }
    }
  }

  async function loadSchemes(editionCode?: string | null) {
    schemes.value = await listAdminDictionarySchemes(editionCode ?? lastEditionCode.value)
    return schemes.value
  }

  async function loadChangeLogs(pageNo = 1, pageSize = 20) {
    changeLogs.value = await pageAdminDictionaryChangeLogs(pageNo, pageSize)
    return changeLogs.value
  }

  function applyCurrent(result: AdminDictionaryCurrentResult) {
    current.value = result
    lastEditionCode.value = result.editionCode
  }

  async function checkRevisionAndSync() {
    if (!current.value || syncing.value) {
      return
    }
    syncing.value = true
    try {
      const revision = await getAdminDictionaryRevision()
      if (
        revision.dictionaryRevision !== current.value.dictionaryRevision
        || revision.editionCode !== current.value.editionCode
      ) {
        await loadCurrent(revision.editionCode, { silent: true })
      }
    } finally {
      syncing.value = false
    }
  }

  function startSilentSync() {
    if (timerId.value !== null || typeof window === 'undefined') {
      return
    }
    timerId.value = window.setInterval(() => {
      void checkRevisionAndSync()
    }, SYNC_INTERVAL_MS)
  }

  function stopSilentSync() {
    if (timerId.value !== null && typeof window !== 'undefined') {
      window.clearInterval(timerId.value)
    }
    timerId.value = null
  }

  return {
    current,
    schemes,
    changeLogs,
    loading,
    syncing,
    itemMap,
    displayName,
    marketName,
    orderTypeName,
    loadCurrent,
    loadSchemes,
    loadChangeLogs,
    applyCurrent,
    checkRevisionAndSync,
    startSilentSync,
    stopSilentSync,
  }
})
