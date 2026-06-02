import { request } from '@/api/http'
import type {
  CreateSnapshotRequest,
  RestoreGroupSnapshotRequest,
  RestoreGroupSnapshotResult,
  SnapshotDetailResult,
  SnapshotListResult,
} from '@/types/sandbox-game-admin'

export interface ListSnapshotsParams {
  snapshotScope?: string
  snapshotType?: string
  groupId?: number | null
  yearNo?: number | null
  stageCode?: string
  pageNo?: number
  pageSize?: number
}

export function listAdminSnapshots(params: ListSnapshotsParams) {
  return request<SnapshotListResult>(`/api/v1/sandbox-game/admin-rollback/list-snapshots${buildQuery(params)}`)
}

export function getAdminSnapshotDetail(snapshotId: number) {
  return request<SnapshotDetailResult>(`/api/v1/sandbox-game/admin-rollback/get-snapshot-detail?snapshotId=${snapshotId}`)
}

export function createAdminSnapshot(payload: CreateSnapshotRequest) {
  return request('/api/v1/sandbox-game/admin-rollback/create-snapshot', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function restoreAdminGroupSnapshot(payload: RestoreGroupSnapshotRequest) {
  return request<RestoreGroupSnapshotResult>('/api/v1/sandbox-game/admin-rollback/restore-group-snapshot', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

function buildQuery(params: ListSnapshotsParams) {
  const query = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null || value === '') {
      continue
    }
    query.set(key, String(value))
  }
  const value = query.toString()
  return value ? `?${value}` : ''
}
