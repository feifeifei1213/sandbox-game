import { request } from '@/api/http'
import type { PlayerAdjustmentSyncResult } from '@/types/sandbox-game'

export function getPlayerAdjustmentSync(yearNo: number, knownRevision: number) {
  const query = new URLSearchParams({ yearNo: String(yearNo), knownRevision: String(knownRevision) })
  return request<PlayerAdjustmentSyncResult>(`/api/v1/sandbox-game/player-notice/get-adjustment-sync?${query.toString()}`)
}
