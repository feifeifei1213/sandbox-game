import { request } from '@/api/http'
import type {
  PlayerReportView,
  SavePlayerReportDraftRequest,
  SavePlayerReportDraftResponse,
  SubmitPlayerReportRequest,
  SubmitPlayerReportResponse,
} from '@/types/sandbox-game'

export function getPlayerReportView(yearNo: number) {
  return request<PlayerReportView>(`/api/v1/sandbox-game/player-report/get-view?yearNo=${yearNo}`)
}

export function savePlayerReportDraft(payload: SavePlayerReportDraftRequest) {
  return request<SavePlayerReportDraftResponse>('/api/v1/sandbox-game/player-report/save-draft', {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}

export function submitPlayerReport(payload: SubmitPlayerReportRequest) {
  return request<SubmitPlayerReportResponse>('/api/v1/sandbox-game/player-report/submit', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}