import { request } from '@/api/http'
import type { CurrentGameConfigResult, YearTabsResult } from '@/types/sandbox-game'

export function getCurrentGameConfig() {
  return request<CurrentGameConfigResult>('/api/v1/sandbox-game/game-config/get-current')
}

export function getYearTabs() {
  return request<YearTabsResult>('/api/v1/sandbox-game/game-config/get-year-tabs')
}