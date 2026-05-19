import { request } from '@/api/http'
import type {
  AdminOrderType,
  GenerateOrderPoolRequest,
  GenerateOrderPoolResult,
  OrderControlConfigResult,
  OrderPoolResult,
  UpdateOrderControlConfigRequest,
  UpdateOrderControlConfigResult,
  UploadOrderExcelResult,
} from '@/types/sandbox-game-admin'
import type { AdminMarketSelectionStatus, AdminOrderControlResult, OrderMarketCode, OrderTypeCode } from '@/types/sandbox-game-order'
import type { CommonResult } from '@/types/http'
import { clearStoredAuthSession, getStoredAuthSession } from '@/utils/auth-session'

const API_BASE = import.meta.env.VITE_API_BASE ?? ''

export function getAdminOrderControlConfig(yearNo: number) {
  return request<OrderControlConfigResult>(`/api/v1/sandbox-game/admin-order/get-control-config?yearNo=${yearNo}`)
}

export async function uploadAdminOrderExcel(file: File) {
  const form = new FormData()
  form.append('file', file)
  return multipartRequest<UploadOrderExcelResult>('/api/v1/sandbox-game/admin-order/upload-excel', form)
}

export function updateAdminOrderControlConfig(payload: UpdateOrderControlConfigRequest) {
  return request<UpdateOrderControlConfigResult>('/api/v1/sandbox-game/admin-order/update-control-config', {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}

export function generateAdminOrderPool(payload: GenerateOrderPoolRequest) {
  return request<GenerateOrderPoolResult>('/api/v1/sandbox-game/admin-order/generate-order-pool', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function getAdminOrderPool(yearNo: number, marketCode: OrderMarketCode, orderType: AdminOrderType) {
  const params = new URLSearchParams({
    yearNo: String(yearNo),
    marketCode,
    orderType,
  })
  return request<OrderPoolResult>(`/api/v1/sandbox-game/admin-order/get-order-pool?${params.toString()}`)
}

export function openAdminMarketBidding(payload: { yearNo: number; marketCode: OrderMarketCode }) {
  return request<AdminOrderControlResult>('/api/v1/sandbox-game/admin-order/open-market-bidding', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function closeAdminMarketBidding(payload: { yearNo: number; marketCode: OrderMarketCode }) {
  return request<AdminOrderControlResult>('/api/v1/sandbox-game/admin-order/close-market-bidding', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function getAdminMarketSelectionStatus(yearNo: number, marketCode: OrderMarketCode) {
  const params = new URLSearchParams({
    yearNo: String(yearNo),
    marketCode,
  })
  return request<AdminMarketSelectionStatus>(`/api/v1/sandbox-game/admin-order/get-market-selection-status?${params.toString()}`)
}

export function releaseNextAdminOrderSegment(payload: { yearNo: number }) {
  return request<AdminOrderControlResult>('/api/v1/sandbox-game/admin-order/release-next-segment', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function adminSkipCurrentOrderGroup(payload: {
  yearNo: number
  marketCode: OrderMarketCode
  orderType: OrderTypeCode
  groupId: number
  reason: string
}) {
  return request<AdminOrderControlResult>('/api/v1/sandbox-game/admin-order/admin-skip-current-group', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

async function multipartRequest<T>(url: string, body: FormData): Promise<T> {
  const headers: HeadersInit = {
    Accept: 'application/json',
  }
  const session = getStoredAuthSession()
  if (session) {
    headers.Authorization = `${session.tokenType} ${session.accessToken}`
  }

  const response = await fetch(`${API_BASE}${url}`, {
    method: 'POST',
    headers,
    body,
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
