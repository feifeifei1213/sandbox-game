import { request } from '@/api/http'
import type { AdminFinalRankingResult, AdminYearSummaryResult } from '@/types/sandbox-game-admin'

export function getAdminYearSummary(yearNo: number) {
  return request<AdminYearSummaryResult>(`/api/v1/sandbox-game/admin-summary/get-year-summary?yearNo=${yearNo}`)
}

export function getAdminFinalRanking() {
  return request<AdminFinalRankingResult>('/api/v1/sandbox-game/admin-summary/get-final-ranking')
}
