import { request } from '@/api/http'
import type {
  AdminDictionaryApplySchemeRequest,
  AdminDictionaryChangeLogPageResult,
  AdminDictionaryCurrentResult,
  AdminDictionaryRevisionResult,
  AdminDictionarySaveSchemeRequest,
  AdminDictionarySchemeDetailResult,
  AdminDictionarySchemeListResult,
  AdminDictionaryUpdateCurrentRequest,
} from '@/types/sandbox-game-admin'

export function getAdminDictionaryCurrent(editionCode?: string | null) {
  const query = editionCode ? `?editionCode=${encodeURIComponent(editionCode)}` : ''
  return request<AdminDictionaryCurrentResult>(`/api/v1/sandbox-game/admin-dictionary/get-current${query}`)
}

export function listAdminDictionarySchemes(editionCode?: string | null) {
  const query = editionCode ? `?editionCode=${encodeURIComponent(editionCode)}` : ''
  return request<AdminDictionarySchemeListResult>(`/api/v1/sandbox-game/admin-dictionary/list-schemes${query}`)
}

export function getAdminDictionarySchemeDetail(id: number) {
  return request<AdminDictionarySchemeDetailResult>(
    `/api/v1/sandbox-game/admin-dictionary/get-scheme-detail?id=${encodeURIComponent(String(id))}`,
  )
}

export function saveAdminDictionaryScheme(payload: AdminDictionarySaveSchemeRequest) {
  return request<AdminDictionarySchemeDetailResult>('/api/v1/sandbox-game/admin-dictionary/save-scheme', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function deleteAdminDictionaryScheme(id: number) {
  return request<{ deleted: boolean }>(`/api/v1/sandbox-game/admin-dictionary/delete-scheme?id=${encodeURIComponent(String(id))}`, {
    method: 'DELETE',
  })
}

export function updateAdminDictionaryCurrent(payload: AdminDictionaryUpdateCurrentRequest) {
  return request<AdminDictionaryCurrentResult>('/api/v1/sandbox-game/admin-dictionary/update-current', {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}

export function applyAdminDictionarySchemeToCurrent(payload: AdminDictionaryApplySchemeRequest) {
  return request<AdminDictionaryCurrentResult>('/api/v1/sandbox-game/admin-dictionary/apply-scheme-to-current', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function restoreAdminDictionaryCurrentDefault() {
  return request<AdminDictionaryCurrentResult>('/api/v1/sandbox-game/admin-dictionary/restore-current-default', {
    method: 'POST',
  })
}

export function pageAdminDictionaryChangeLogs(pageNo = 1, pageSize = 20) {
  return request<AdminDictionaryChangeLogPageResult>(
    `/api/v1/sandbox-game/admin-dictionary/page-change-logs?pageNo=${encodeURIComponent(String(pageNo))}&pageSize=${encodeURIComponent(String(pageSize))}`,
  )
}

export function getAdminDictionaryRevision() {
  return request<AdminDictionaryRevisionResult>('/api/v1/sandbox-game/admin-dictionary/get-revision')
}
