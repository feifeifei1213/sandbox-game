import { buildBypassHeaders, request } from '@/api/http'
import type {
  AdminControlConfigResult,
  InitialBaselineViewResult,
  OpenNextYearResult,
  SubmitInitialBaselineRequest,
  SubmitInitialBaselineResult,
  UpdateFinalYearResult,
} from '@/types/sandbox-game-admin'

const adminHeaders = buildBypassHeaders({
  roleType: 'ADMIN',
  userId: import.meta.env.VITE_ADMIN_BYPASS_USER_ID ?? '1',
  username: import.meta.env.VITE_ADMIN_BYPASS_USERNAME ?? 'admin',
})

export function getAdminControlConfig() {
  return request<AdminControlConfigResult>('/api/v1/sandbox-game/admin-control/get-config', {
    headers: adminHeaders,
  })
}

export function updateAdminFinalYear(finalYear: number) {
  return request<UpdateFinalYearResult>('/api/v1/sandbox-game/admin-control/update-final-year', {
    method: 'PUT',
    headers: adminHeaders,
    body: JSON.stringify({ finalYear }),
  })
}

export function openAdminNextYear(targetYearNo: number) {
  return request<OpenNextYearResult>('/api/v1/sandbox-game/admin-control/open-next-year', {
    method: 'POST',
    headers: adminHeaders,
    body: JSON.stringify({ targetYearNo }),
  })
}

export function getAdminInitialBaseline() {
  return request<InitialBaselineViewResult>('/api/v1/sandbox-game/admin-control/get-initial-baseline', {
    headers: adminHeaders,
  })
}

export function submitAdminInitialBaseline(payload: SubmitInitialBaselineRequest) {
  return request<SubmitInitialBaselineResult>('/api/v1/sandbox-game/admin-control/submit-initial-baseline', {
    method: 'POST',
    headers: adminHeaders,
    body: JSON.stringify(payload),
  })
}
