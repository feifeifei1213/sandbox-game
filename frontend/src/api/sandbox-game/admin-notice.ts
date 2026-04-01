import { request } from '@/api/http'
import type {
  AdminNoticeRecordsResult,
  SendAdminAdjustmentRequest,
  SendAdminAdjustmentResult,
  SendAdminGeneralNoticeRequest,
  SendAdminGeneralNoticeResult,
} from '@/types/sandbox-game-admin'

export function getAdminNoticeRecords(limit = 20) {
  return request<AdminNoticeRecordsResult>(`/api/v1/sandbox-game/admin-notice/get-records?limit=${limit}`)
}

export function sendAdminGeneralNotice(payload: SendAdminGeneralNoticeRequest) {
  return request<SendAdminGeneralNoticeResult>('/api/v1/sandbox-game/admin-notice/send-general', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function sendAdminAdjustment(payload: SendAdminAdjustmentRequest) {
  return request<SendAdminAdjustmentResult>('/api/v1/sandbox-game/admin-notice/send-adjustment', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}
