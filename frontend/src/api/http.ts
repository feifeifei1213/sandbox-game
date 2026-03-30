import type { CommonResult } from '@/types/http'
import { clearStoredAuthSession, getStoredAuthSession } from '@/utils/auth-session'

const API_BASE = import.meta.env.VITE_API_BASE ?? ''

function buildDefaultHeaders(): HeadersInit {
  const headers: Record<string, string> = {
    Accept: 'application/json',
  }

  const session = getStoredAuthSession()
  if (session) {
    headers.Authorization = `${session.tokenType} ${session.accessToken}`
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
  const json = tryParseCommonResult<T>(raw)

  if (response.status === 401) {
    clearStoredAuthSession()
  }

  if (!response.ok) {
    throw new Error(resolveHttpErrorMessage(response.status, json, raw))
  }

  if (!json) {
    throw new Error(resolveNonJSONMessage(raw))
  }

  if (json.code === 401) {
    clearStoredAuthSession()
  }

  if (json.code !== 0) {
    throw new Error(json.msg || '接口返回失败')
  }

  return json.data
}

function tryParseCommonResult<T>(raw: string): CommonResult<T> | null {
  if (!raw.trim()) {
    return null
  }

  try {
    return JSON.parse(raw) as CommonResult<T>
  } catch {
    return null
  }
}

function resolveHttpErrorMessage<T>(status: number, json: CommonResult<T> | null, raw: string): string {
  if (json?.msg) {
    return json.msg
  }
  if (raw.trim()) {
    return `接口返回异常：${truncateRawText(raw)}`
  }
  return `请求失败（${status}）`
}

function resolveNonJSONMessage(raw: string): string {
  if (raw.trim()) {
    return `接口返回非 JSON 内容：${truncateRawText(raw)}`
  }
  return '接口返回为空'
}

function truncateRawText(raw: string): string {
  const normalized = raw.replace(/\s+/g, ' ').trim()
  if (normalized.length <= 120) {
    return normalized
  }
  return `${normalized.slice(0, 117)}...`
}
