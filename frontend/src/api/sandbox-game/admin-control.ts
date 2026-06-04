import { request } from '@/api/http'
import type {
  AdminControlConfigResult,
  AdminControlSetupStatusResult,
  InitializeGameResult,
  InitialBaselineViewResult,
  OpenNextYearResult,
  SubmitInitialBaselineRequest,
  SubmitInitialBaselineResult,
  UnlockYearRequest,
  UnlockYearResult,
  UpdateFinalYearResult,
} from '@/types/sandbox-game-admin'

export function getAdminControlSetupStatus() {
  return request<AdminControlSetupStatusResult>('/api/v1/sandbox-game/admin-control/get-setup-status')
}

export function initializeAdminGame(groupCount: number, editionCode: string) {
  return request<InitializeGameResult>('/api/v1/sandbox-game/admin-control/initialize-game', {
    method: 'POST',
    body: JSON.stringify({ groupCount, editionCode }),
  })
}

export function initializeAdminGameWithDictionary(
  groupCount: number,
  editionCode: string,
  dictionaryItems: Array<{ itemCode: string; displayName: string }>,
  dictionarySchemeId?: number | null,
) {
  return request<InitializeGameResult>('/api/v1/sandbox-game/admin-control/initialize-game', {
    method: 'POST',
    body: JSON.stringify({
      groupCount,
      editionCode,
      dictionaryItems,
      dictionarySchemeId: dictionarySchemeId ?? null,
    }),
  })
}

export function getAdminControlConfig() {
  return request<AdminControlConfigResult>('/api/v1/sandbox-game/admin-control/get-config')
}

export function updateAdminFinalYear(finalYear: number) {
  return request<UpdateFinalYearResult>('/api/v1/sandbox-game/admin-control/update-final-year', {
    method: 'PUT',
    body: JSON.stringify({ finalYear }),
  })
}

export function openAdminNextYear(targetYearNo: number) {
  return request<OpenNextYearResult>('/api/v1/sandbox-game/admin-control/open-next-year', {
    method: 'POST',
    body: JSON.stringify({ targetYearNo }),
  })
}

export function getAdminInitialBaseline() {
  return request<InitialBaselineViewResult>('/api/v1/sandbox-game/admin-control/get-initial-baseline')
}

export function submitAdminInitialBaseline(payload: SubmitInitialBaselineRequest) {
  return request<SubmitInitialBaselineResult>('/api/v1/sandbox-game/admin-control/submit-initial-baseline', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function unlockAdminYear(payload: UnlockYearRequest) {
  return request<UnlockYearResult>('/api/v1/sandbox-game/admin-control/unlock-year', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}
