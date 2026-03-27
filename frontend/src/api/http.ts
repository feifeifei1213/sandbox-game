import type { CommonResult } from '@/types/http'

const API_BASE = import.meta.env.VITE_API_BASE ?? ''

function buildDefaultHeaders(): HeadersInit {
  const headers: Record<string, string> = {
    Accept: 'application/json',
  }

  if (import.meta.env.DEV) {
    headers['X-Role-Type'] = import.meta.env.VITE_BYPASS_ROLE_TYPE ?? 'GROUP'
    headers['X-User-Id'] = import.meta.env.VITE_BYPASS_USER_ID ?? '101'
    headers['X-Username'] = import.meta.env.VITE_BYPASS_USERNAME ?? 'group01'
    headers['X-Group-Id'] = import.meta.env.VITE_BYPASS_GROUP_ID ?? '1'
  }

  return headers
}

export async function request<T>(url: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(`${API_BASE}${url}`, {
    ...init,
    headers: {
      ...buildDefaultHeaders(),
      ...(init.body ? { 'Content-Type': 'application/json' } : {}),
      ...(init.headers ?? {}),
    },
  })

  const raw = await response.text()
  const json = raw ? (JSON.parse(raw) as CommonResult<T>) : null

  if (!response.ok) {
    throw new Error(json?.msg || `请求失败（${response.status}）`)
  }

  if (!json) {
    throw new Error('接口返回为空')
  }

  if (json.code !== 0) {
    throw new Error(json.msg || '接口返回失败')
  }

  return json.data
}