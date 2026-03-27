import { request } from '@/api/http'
import type {
  PlayerOperatingView,
  SaveDraftRequest,
  SaveDraftResponse,
  SubmitStageRequest,
  SubmitStageResponse,
} from '@/types/sandbox-game'

export function getPlayerOperatingYearView(yearNo: number) {
  return request<PlayerOperatingView>(`/api/v1/sandbox-game/player-operating/get-year-view?yearNo=${yearNo}`)
}

export function savePlayerOperatingDraft(payload: SaveDraftRequest) {
  return request<SaveDraftResponse>('/api/v1/sandbox-game/player-operating/save-draft', {
    method: 'PUT',
    body: JSON.stringify(payload),
  })
}

export function submitPlayerOperatingStage(payload: SubmitStageRequest) {
  return request<SubmitStageResponse>('/api/v1/sandbox-game/player-operating/submit-stage', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}