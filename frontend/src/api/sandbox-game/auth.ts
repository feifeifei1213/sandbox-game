import { request } from '@/api/http'
import type {
  AuthCurrentUserResult,
  AuthLoginPayload,
  AuthLoginResult,
  AuthLogoutResult,
} from '@/types/auth'

export function login(payload: AuthLoginPayload) {
  return request<AuthLoginResult>('/api/v1/sandbox-game/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function getCurrentUser() {
  return request<AuthCurrentUserResult>('/api/v1/sandbox-game/auth/get-current-user')
}

export function logout() {
  return request<AuthLogoutResult>('/api/v1/sandbox-game/auth/logout', {
    method: 'POST',
  })
}
