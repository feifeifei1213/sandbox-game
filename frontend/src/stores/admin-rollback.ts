import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { unlockAdminYear } from '@/api/sandbox-game/admin-control'
import {
  createAdminSnapshot,
  getAdminSnapshotDetail,
  listAdminSnapshots,
  restoreAdminGroupSnapshot,
} from '@/api/sandbox-game/admin-rollback'
import { getAdminGroupOperationContext, listAdminGroups } from '@/api/sandbox-game/admin-group-data'
import type { PageMessage } from '@/stores/admin-shell'
import type {
  AdminGroupOperationContextResult,
  AdminGroupOption,
  SnapshotDetailResult,
  SnapshotScope,
  SnapshotSummary,
  SnapshotType,
  UnlockStageCode,
  UnlockTargetType,
} from '@/types/sandbox-game-admin'

const DEFAULT_PAGE_SIZE = 20

export const useAdminRollbackStore = defineStore('sandbox-admin-rollback', () => {
  const groups = ref<AdminGroupOption[]>([])
  const snapshots = ref<SnapshotSummary[]>([])
  const total = ref(0)
  const pageNo = ref(1)
  const pageSize = ref(DEFAULT_PAGE_SIZE)
  const loading = ref(false)
  const operating = ref(false)
  const pageMessage = ref<PageMessage | null>(null)
  const selectedSnapshotId = ref<number | null>(null)
  const snapshotDetail = ref<SnapshotDetailResult | null>(null)
  const operationContext = ref<AdminGroupOperationContextResult | null>(null)

  const filters = ref({
    snapshotScope: 'GROUP' as SnapshotScope | '',
    snapshotType: '' as SnapshotType | '',
    groupId: null as number | null,
    yearNo: null as number | null,
    stageCode: '',
  })

  const unlockForm = ref({
    groupId: null as number | null,
    yearNo: 0,
    unlockTargetType: 'OPERATING' as UnlockTargetType,
    targetStageCode: 'YEAR_END' as UnlockStageCode | null,
    reason: '',
  })

  const manualSnapshotForm = ref({
    snapshotScope: 'GROUP' as SnapshotScope,
    groupId: null as number | null,
    yearNo: 0,
    stageCode: '',
    description: '',
  })

  const selectedGroup = computed(() => groups.value.find((item) => item.groupId === unlockForm.value.groupId) ?? null)
  const selectedSnapshot = computed(() => snapshots.value.find((item) => item.id === selectedSnapshotId.value) ?? null)

  async function bootstrap() {
    loading.value = true
    pageMessage.value = null
    try {
      const result = await listAdminGroups()
      groups.value = result.list
      if (!unlockForm.value.groupId) {
        unlockForm.value.groupId = groups.value[0]?.groupId ?? null
      }
      if (!manualSnapshotForm.value.groupId) {
        manualSnapshotForm.value.groupId = groups.value[0]?.groupId ?? null
      }
      if (!filters.value.groupId) {
        filters.value.groupId = groups.value[0]?.groupId ?? null
      }
      await loadOperationContext(unlockForm.value.groupId)
      await loadSnapshots({ silent: true })
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '初始化回退与修正页面失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  async function loadSnapshots(options?: { silent?: boolean }) {
    loading.value = true
    if (!options?.silent) {
      pageMessage.value = null
    }
    try {
      const result = await listAdminSnapshots({
        ...filters.value,
        pageNo: pageNo.value,
        pageSize: pageSize.value,
      })
      snapshots.value = result.list
      total.value = result.total
      pageNo.value = result.pageNo
      pageSize.value = result.pageSize
      if (selectedSnapshotId.value && !snapshots.value.some((item) => item.id === selectedSnapshotId.value)) {
        selectedSnapshotId.value = null
        snapshotDetail.value = null
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '读取状态快照失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  async function selectSnapshot(snapshotId: number) {
    selectedSnapshotId.value = snapshotId
    snapshotDetail.value = null
    pageMessage.value = null
    try {
      snapshotDetail.value = await getAdminSnapshotDetail(snapshotId)
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '读取快照详情失败')
      throw error
    }
  }

  async function loadOperationContext(groupId = unlockForm.value.groupId) {
    if (!groupId) {
      operationContext.value = null
      return null
    }
    try {
      const result = await getAdminGroupOperationContext(groupId)
      operationContext.value = result
      unlockForm.value.groupId = result.groupId
      unlockForm.value.yearNo = result.operationYearNo
      return result
    } catch (error) {
      operationContext.value = null
      pageMessage.value = toErrorMessage(error, '读取目标小组操作年份失败')
      throw error
    }
  }

  async function submitUnlockRetry() {
    if (!unlockForm.value.groupId) {
      pageMessage.value = { type: 'error', text: '请选择目标小组' }
      throw new Error('请选择目标小组')
    }
    if (unlockForm.value.unlockTargetType === 'OPERATING' && !unlockForm.value.targetStageCode) {
      pageMessage.value = { type: 'error', text: '请选择经营阶段' }
      throw new Error('请选择经营阶段')
    }
    const context =
      operationContext.value?.groupId === unlockForm.value.groupId
        ? operationContext.value
        : await loadOperationContext(unlockForm.value.groupId)
    if (!context) {
      pageMessage.value = { type: 'error', text: '请先选择目标小组' }
      throw new Error('请先选择目标小组')
    }
    if (!context.canUnlockRetry) {
      const message = buildUnlockBlockedMessage(context)
      pageMessage.value = { type: 'error', text: message }
      throw new Error(message)
    }

    operating.value = true
    pageMessage.value = null
    try {
      const result = await unlockAdminYear({
        groupId: unlockForm.value.groupId,
        yearNo: context.operationYearNo,
        unlockTargetType: unlockForm.value.unlockTargetType,
        targetStageCode: unlockForm.value.unlockTargetType === 'OPERATING' ? unlockForm.value.targetStageCode : null,
        reason: unlockForm.value.reason.trim(),
      })
      pageMessage.value = {
        type: 'success',
        text: buildUnlockSuccessMessage(result, selectedGroup.value),
      }
      unlockForm.value.reason = ''
      await loadOperationContext(result.groupId)
      await loadSnapshots({ silent: true })
      return result
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '提交退回重提失败')
      throw error
    } finally {
      operating.value = false
    }
  }

  async function submitManualSnapshot() {
    if (manualSnapshotForm.value.snapshotScope === 'GROUP' && !manualSnapshotForm.value.groupId) {
      throw new Error('单组快照必须选择小组')
    }
    if (!manualSnapshotForm.value.description.trim()) {
      throw new Error('快照说明不能为空')
    }
    operating.value = true
    try {
      await createAdminSnapshot({
        snapshotScope: manualSnapshotForm.value.snapshotScope,
        groupId: manualSnapshotForm.value.snapshotScope === 'GROUP' ? manualSnapshotForm.value.groupId : null,
        yearNo: manualSnapshotForm.value.yearNo,
        stageCode: manualSnapshotForm.value.stageCode,
        description: manualSnapshotForm.value.description.trim(),
      })
      pageMessage.value = { type: 'success', text: '手动快照已创建。' }
      manualSnapshotForm.value.description = ''
      await loadSnapshots({ silent: true })
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '创建手动快照失败')
      throw error
    } finally {
      operating.value = false
    }
  }

  async function submitRestoreSnapshot() {
    if (!selectedSnapshotId.value) {
      throw new Error('请选择要恢复的单组快照')
    }
    operating.value = true
    try {
      const result = await restoreAdminGroupSnapshot({
        snapshotId: selectedSnapshotId.value,
      })
      pageMessage.value = {
        type: 'success',
        text: `快照恢复已完成，回退日志 #${result.rollbackLogId}，安全快照 #${result.safetySnapshotId}。`,
      }
      await loadOperationContext(result.groupId)
      await loadSnapshots({ silent: true })
      if (selectedSnapshotId.value) {
        await selectSnapshot(selectedSnapshotId.value)
      }
      return result
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '恢复快照失败')
      throw error
    } finally {
      operating.value = false
    }
  }

  function setFilterGroup(groupId: number | null) {
    filters.value.groupId = groupId
    pageNo.value = 1
  }

  return {
    groups,
    snapshots,
    total,
    pageNo,
    pageSize,
    loading,
    operating,
    pageMessage,
    filters,
    unlockForm,
    manualSnapshotForm,
    operationContext,
    selectedSnapshotId,
    selectedSnapshot,
    snapshotDetail,
    selectedGroup,
    bootstrap,
    loadOperationContext,
    loadSnapshots,
    selectSnapshot,
    submitUnlockRetry,
    submitManualSnapshot,
    submitRestoreSnapshot,
    setFilterGroup,
  }
})

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}

function buildUnlockSuccessMessage(result: { yearNo: number; unlockTargetType: string; targetStageCode?: string | null; editableStageCode?: string | null }, group: AdminGroupOption | null) {
  const groupLabel = group ? `第${group.groupNo}组` : '目标小组'
  if (result.unlockTargetType === 'REPORT') {
    return `已退回${groupLabel} ${result.yearNo} 年财报页，可重新提交财报。`
  }
  return `已退回${groupLabel} ${result.yearNo} 年经营页到 ${formatUnlockStage(result.editableStageCode ?? result.targetStageCode)}，可重新提交该阶段。`
}

function buildUnlockBlockedMessage(context: AdminGroupOperationContextResult) {
  if (context.businessStatus === 'BANKRUPT') {
    return '目标小组已破产，不能退回重提'
  }
  if (context.currentOpenYear > context.operationYearNo && !context.rollbackPending) {
    return '普通退回重提只处理当前操作年份；跨年修正请使用恢复快照'
  }
  return '目标小组当前状态不能退回重提'
}

function formatUnlockStage(value?: string | null) {
  const map: Record<string, string> = {
    Q1: 'Q1',
    Q2: 'Q2',
    Q3: 'Q3',
    Q4: 'Q4',
    YEAR_END: '年末',
  }
  return value ? map[value] ?? value : '--'
}
