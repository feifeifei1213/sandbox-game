import { request } from '@/api/http'
import type { CurrentGameConfigResult, YearTabsResult } from '@/types/sandbox-game'
import type { GameEdition } from '@/types/sandbox-game-admin'

export function getCurrentGameConfig() {
  return request<CurrentGameConfigResult>('/api/v1/sandbox-game/game-config/get-current')
}

export function getYearTabs() {
  return request<YearTabsResult>('/api/v1/sandbox-game/game-config/get-year-tabs')
}

export function listGameEditions() {
  return request<GameEdition[]>('/api/v1/sandbox-game/game-config/list-game-editions')
}
