import { request } from '@/api/http'
import type { PlayerOperatingView, PlayerReportView } from '@/types/sandbox-game'
import type { AdminGroupOperationContextResult, ListAdminGroupsResult } from '@/types/sandbox-game-admin'

export function listAdminGroups() {
  return request<ListAdminGroupsResult>('/api/v1/sandbox-game/admin-group-data/list-groups')
}

export function getAdminGroupOperationContext(groupId: number) {
  return request<AdminGroupOperationContextResult>(`/api/v1/sandbox-game/admin-group-data/get-operation-context?groupId=${groupId}`)
}

export function getAdminGroupOperatingView(groupId: number, yearNo: number) {
  return request<PlayerOperatingView>(`/api/v1/sandbox-game/admin-group-data/get-operating-view?groupId=${groupId}&yearNo=${yearNo}`)
}

export function getAdminGroupReportView(groupId: number, yearNo: number) {
  return request<PlayerReportView>(`/api/v1/sandbox-game/admin-group-data/get-report-view?groupId=${groupId}&yearNo=${yearNo}`)
}
