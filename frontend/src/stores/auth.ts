import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { getCurrentUser, login as loginApi, logout as logoutApi } from '@/api/sandbox-game/auth'
import type { AuthSession } from '@/types/auth'
import {
  buildAuthSessionFromLoginResult,
  buildAuthUserFromCurrentUser,
  clearStoredAuthSession,
  getStoredAuthSession,
  normalizeDefaultRoute,
  setStoredAuthSession,
} from '@/utils/auth-session'

export const useAuthStore = defineStore('sandbox-auth', () => {
  const session = ref<AuthSession | null>(null)
  const hydrated = ref(false)
  const validated = ref(false)
  const loadingCurrentUser = ref(false)

  let validationPromise: Promise<boolean> | null = null

  const currentUser = computed(() => session.value?.user ?? null)
  const isAuthenticated = computed(() => currentUser.value !== null)

  function restoreSession() {
    if (hydrated.value) {
      return
    }
    session.value = getStoredAuthSession()
    hydrated.value = true
    validated.value = false
  }

  async function login(username: string, password: string) {
    const result = await loginApi({ username, password })
    const nextSession = buildAuthSessionFromLoginResult(result)
    persistSession(nextSession, true)
    return resolveDefaultRoute()
  }

  async function ensureAuthenticated() {
    restoreSession()
    if (!session.value) {
      return false
    }
    if (validated.value) {
      return true
    }
    if (validationPromise) {
      return validationPromise
    }

    validationPromise = (async () => {
      loadingCurrentUser.value = true
      try {
        const result = await getCurrentUser()
        if (!session.value) {
          return false
        }
        persistSession(
          {
            ...session.value,
            user: buildAuthUserFromCurrentUser(result),
          },
          true,
        )
        return true
      } catch (error) {
        if (!getStoredAuthSession()) {
          clearSession()
          return false
        }
        throw error
      } finally {
        loadingCurrentUser.value = false
        validationPromise = null
      }
    })()

    return validationPromise
  }

  async function logout() {
    try {
      if (session.value) {
        await logoutApi()
      }
    } finally {
      clearSession()
    }
  }

  function clearSession() {
    session.value = null
    hydrated.value = true
    validated.value = false
    clearStoredAuthSession()
  }

  function resolveDefaultRoute() {
    return normalizeDefaultRoute(currentUser.value?.defaultRoute ?? '/sandbox-game/login')
  }

  function persistSession(nextSession: AuthSession, nextValidated: boolean) {
    session.value = nextSession
    hydrated.value = true
    validated.value = nextValidated
    setStoredAuthSession(nextSession)
  }

  return {
    session,
    currentUser,
    isAuthenticated,
    hydrated,
    validated,
    loadingCurrentUser,
    restoreSession,
    login,
    ensureAuthenticated,
    logout,
    clearSession,
    resolveDefaultRoute,
  }
})
