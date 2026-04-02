import { defineStore } from 'pinia'
import { ref } from 'vue'

import { getAdminControlConfig, getAdminControlSetupStatus } from '@/api/sandbox-game/admin-control'
import type { AdminControlConfigResult, AdminControlSetupStatusResult } from '@/types/sandbox-game-admin'

type MessageType = 'success' | 'error' | 'info'

export interface PageMessage {
  type: MessageType
  text: string
}

export const useAdminShellStore = defineStore('sandbox-admin-shell', () => {
  const config = ref<AdminControlConfigResult | null>(null)
  const setupStatus = ref<AdminControlSetupStatusResult | null>(null)
  const loading = ref(false)
  const pageMessage = ref<PageMessage | null>(null)

  async function bootstrap() {
    loading.value = true
    pageMessage.value = null
    try {
      const [nextSetupStatus, nextConfig] = await Promise.all([
        getAdminControlSetupStatus(),
        getAdminControlConfig(),
      ])
      setupStatus.value = nextSetupStatus
      config.value = nextConfig
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

  async function refreshSetupStatus(options?: { silent?: boolean }) {
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      setupStatus.value = await getAdminControlSetupStatus()
    } catch (error) {
      if (!options?.silent) {
        pageMessage.value = toErrorMessage(error, '刷新赛前初始化状态失败')
      }
      throw error
    }
  }

  async function refreshAll(options?: { silent?: boolean }) {
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      const [nextSetupStatus, nextConfig] = await Promise.all([
        getAdminControlSetupStatus(),
        getAdminControlConfig(),
      ])
      setupStatus.value = nextSetupStatus
      config.value = nextConfig
    } catch (error) {
      if (!options?.silent) {
        pageMessage.value = toErrorMessage(error, '刷新管理员状态失败')
      }
      throw error
    }
  }

  function setMessage(message: PageMessage | null) {
    pageMessage.value = message
  }

  return {
    config,
    setupStatus,
    loading,
    pageMessage,
    bootstrap,
    refreshConfig,
    refreshSetupStatus,
    refreshAll,
    setMessage,
  }
})

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}
