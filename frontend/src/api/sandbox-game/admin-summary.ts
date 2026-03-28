import { buildBypassHeaders, request } from '@/api/http'
import type { AdminFinalRankingResult, AdminYearSummaryResult } from '@/types/sandbox-game-admin'

const adminHeaders = buildBypassHeaders({
  roleType: 'ADMIN',
  userId: import.meta.env.VITE_ADMIN_BYPASS_USER_ID ?? '1',
  username: import.meta.env.VITE_ADMIN_BYPASS_USERNAME ?? 'admin',
})

export function getAdminYearSummary(yearNo: number) {
  return request<AdminYearSummaryResult>(`/api/v1/sandbox-game/admin-summary/get-year-summary?yearNo=${yearNo}`, {
    headers: adminHeaders,
  })
}

export function getAdminFinalRanking() {
  return request<AdminFinalRankingResult>('/api/v1/sandbox-game/admin-summary/get-final-ranking', {
    headers: adminHeaders,
  })
}
