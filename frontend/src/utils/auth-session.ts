import type {
  AuthCurrentUserResult,
  AuthLoginResult,
  AuthSession,
  AuthUser,
} from '@/types/auth'

const AUTH_SESSION_STORAGE_KEY = 'sandbox-game-auth-session'

export function normalizeDefaultRoute(route: string): string {
  const trimmed = route.trim()
  if (!trimmed) {
    return '/sandbox-game/login'
  }
  if (trimmed.startsWith('/')) {
    return trimmed
  }
  if (trimmed.startsWith('sandbox-game/')) {
    return `/${trimmed}`
  }
  return `/sandbox-game/${trimmed.replace(/^\/+/, '')}`
}

export function buildAuthSessionFromLoginResult(result: AuthLoginResult): AuthSession {
  const expiresAt = Date.now() + result.expiresIn * 1000
  return {
    accessToken: result.accessToken,
    tokenType: result.tokenType,
    expiresAt,
    user: {
      userId: result.user.userId,
      username: result.user.username,
      roleType: result.user.roleType,
      groupId: result.user.groupId,
      defaultRoute: normalizeDefaultRoute(result.defaultRoute),
    },
  }
}

export function buildAuthUserFromCurrentUser(result: AuthCurrentUserResult): AuthUser {
  return {
    userId: result.userId,
    username: result.username,
    roleType: result.roleType,
    groupId: result.groupId,
    defaultRoute: normalizeDefaultRoute(result.defaultRoute),
  }
}

export function getStoredAuthSession(): AuthSession | null {
  if (typeof window === 'undefined') {
    return null
  }

  const raw = window.localStorage.getItem(AUTH_SESSION_STORAGE_KEY)
  if (!raw) {
    return null
  }

  try {
    const parsed = JSON.parse(raw) as unknown
    if (!isValidAuthSession(parsed)) {
      clearStoredAuthSession()
      return null
    }
    if (parsed.expiresAt <= Date.now()) {
      clearStoredAuthSession()
      return null
    }
    return parsed
  } catch {
    clearStoredAuthSession()
    return null
  }
}

export function setStoredAuthSession(session: AuthSession) {
  if (typeof window === 'undefined') {
    return
  }
  window.localStorage.setItem(AUTH_SESSION_STORAGE_KEY, JSON.stringify(session))
}

export function clearStoredAuthSession() {
  if (typeof window === 'undefined') {
    return
  }
  window.localStorage.removeItem(AUTH_SESSION_STORAGE_KEY)
}

function isValidAuthSession(value: unknown): value is AuthSession {
  if (!value || typeof value !== 'object') {
    return false
  }

  const session = value as Partial<AuthSession>
  const user = session.user as Partial<AuthUser> | undefined
  return (
    typeof session.accessToken === 'string' &&
    typeof session.tokenType === 'string' &&
    typeof session.expiresAt === 'number' &&
    !!user &&
    typeof user.userId === 'number' &&
    typeof user.username === 'string' &&
    (user.roleType === 'ADMIN' || user.roleType === 'GROUP') &&
    (user.groupId === null || typeof user.groupId === 'number') &&
    typeof user.defaultRoute === 'string'
  )
}
