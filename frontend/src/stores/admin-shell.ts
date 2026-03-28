import { defineStore } from 'pinia'
import { ref } from 'vue'

import { getAdminControlConfig } from '@/api/sandbox-game/admin-control'
import type { AdminControlConfigResult } from '@/types/sandbox-game-admin'

type MessageType = 'success' | 'error' | 'info'

export interface PageMessage {
  type: MessageType
  text: string
}

export const useAdminShellStore = defineStore('sandbox-admin-shell', () => {
  const config = ref<AdminControlConfigResult | null>(null)
  const loading = ref(false)
  const pageMessage = ref<PageMessage | null>(null)

  async function bootstrap() {
    loading.value = true
    pageMessage.value = null
    try {
      config.value = await getAdminControlConfig()
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '初始化管理员页面失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  async function refreshConfig(options?: { silent?: boolean }) {
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      config.value = await getAdminControlConfig()
    } catch (error) {
      if (!options?.silent) {
        pageMessage.value = toErrorMessage(error, '刷新管理员配置失败')
      }
      throw error
    }
  }

  function setMessage(message: PageMessage | null) {
    pageMessage.value = message
  }

  return {
    config,
    loading,
    pageMessage,
    bootstrap,
    refreshConfig,
    setMessage,
  }
})

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}
