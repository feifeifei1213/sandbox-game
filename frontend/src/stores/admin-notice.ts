import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'

import { listAdminGroups } from '@/api/sandbox-game/admin-group-data'
import {
  getAdminNoticeRecords,
  previewAdminAdjustment,
  sendAdminAdjustment,
  sendAdminGeneralNotice,
  voidAdminAdjustment,
} from '@/api/sandbox-game/admin-notice'
import type {
  AdminAdjustmentRecord,
  AdminGeneralNoticeRecord,
  AdminGroupOption,
  AdjustmentType,
  NoticeTargetScope,
} from '@/types/sandbox-game-admin'
import { hasFractionInput } from '@/utils/manual-integer'

type MessageType = 'success' | 'error' | 'info'

export interface PageMessage {
  type: MessageType
  text: string
}

export const useAdminNoticeStore = defineStore('sandbox-admin-notice', () => {
  const groups = ref<AdminGroupOption[]>([])
  const generalNotices = ref<AdminGeneralNoticeRecord[]>([])
  const adjustments = ref<AdminAdjustmentRecord[]>([])
  const loading = ref(false)
  const sendingGeneral = ref(false)
  const sendingAdjustment = ref(false)
  const previewingAdjustment = ref(false)
  const voidingAdjustment = ref(false)
  const pageMessage = ref<PageMessage | null>(null)

  const generalForm = reactive({
    targetScope: 'ALL' as NoticeTargetScope,
    targetGroupId: null as number | null,
    content: '',
    pinned: false,
  })

  const adjustmentForm = reactive({
    groupId: null as number | null,
    yearNo: 0,
    adjustmentType: 'REWARD' as AdjustmentType,
    amount: '',
    reason: '',
  })

  async function bootstrap(defaultYear = 0) {
    loading.value = true
    pageMessage.value = null
    try {
      const [groupResult, records] = await Promise.all([listAdminGroups(), getAdminNoticeRecords()])
      groups.value = groupResult.list
      generalNotices.value = records.generalNotices
      adjustments.value = records.adjustments
      if (!adjustmentForm.groupId && groups.value.length > 0) {
        adjustmentForm.groupId = groups.value[0].groupId
      }
      adjustmentForm.yearNo = defaultYear
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '初始化通知与奖惩页面失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  async function refreshRecords(options?: { silent?: boolean }) {
    if (!options?.silent) {
      pageMessage.value = null
    }
    const records = await getAdminNoticeRecords()
    generalNotices.value = records.generalNotices
    adjustments.value = records.adjustments
  }

  async function sendGeneral() {
    sendingGeneral.value = true
    pageMessage.value = null
    try {
      const result = await sendAdminGeneralNotice({
        targetScope: generalForm.targetScope,
        targetGroupId: generalForm.targetScope === 'GROUP' ? generalForm.targetGroupId : null,
        content: generalForm.content,
        pinned: generalForm.pinned,
      })
      generalForm.content = ''
      generalForm.pinned = false
      await refreshRecords({ silent: true })
      pageMessage.value = {
        type: 'success',
        text: `普通通知已发送，时间 ${formatDateTime(result.publishedAt)}`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '发送普通通知失败')
      throw error
    } finally {
      sendingGeneral.value = false
    }
  }

  async function sendAdjustmentNotice() {
    validateAdjustmentAmount()
    sendingAdjustment.value = true
    pageMessage.value = null
    try {
      const amount = Number(adjustmentForm.amount)
      const result = await sendAdminAdjustment({
        groupId: adjustmentForm.groupId ?? 0,
        yearNo: adjustmentForm.yearNo,
        adjustmentType: adjustmentForm.adjustmentType,
        amount,
        reason: adjustmentForm.reason,
      })
      adjustmentForm.amount = ''
      adjustmentForm.reason = ''
      await refreshRecords({ silent: true })
      pageMessage.value = {
        type: 'success',
        text: `${formatAdjustmentType(result.adjustmentType)}已下发，时间 ${formatDateTime(result.publishedAt)}`,
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '下发奖惩失败')
      throw error
    } finally {
      sendingAdjustment.value = false
    }
  }

  async function previewAdjustmentNotice() {
    validateAdjustmentAmount()
    previewingAdjustment.value = true
    pageMessage.value = null
    try {
      return await previewAdminAdjustment({
        operation: 'CREATE',
        groupId: adjustmentForm.groupId ?? 0,
        yearNo: adjustmentForm.yearNo,
        adjustmentType: adjustmentForm.adjustmentType,
        amount: Number(adjustmentForm.amount),
        reason: adjustmentForm.reason,
      })
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '预览奖惩影响失败')
      throw error
    } finally {
      previewingAdjustment.value = false
    }
  }

  async function previewVoidAdjustment(adjustmentId: number) {
    previewingAdjustment.value = true
    pageMessage.value = null
    try {
      return await previewAdminAdjustment({ operation: 'VOID', adjustmentId })
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '预览作废影响失败')
      throw error
    } finally {
      previewingAdjustment.value = false
    }
  }

  async function voidAdjustmentNotice(adjustmentId: number, reason: string) {
    voidingAdjustment.value = true
    pageMessage.value = null
    try {
      const result = await voidAdminAdjustment({ adjustmentId, reason })
      await refreshRecords({ silent: true })
      pageMessage.value = { type: 'success', text: `奖惩已作废，时间 ${formatDateTime(result.voidedAt)}` }
      return result
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '作废奖惩失败')
      throw error
    } finally {
      voidingAdjustment.value = false
    }
  }

  function validateAdjustmentAmount() {
    const raw = adjustmentForm.amount
    if (raw === '' || !Number.isFinite(Number(raw)) || Number(raw) <= 0) {
      const error = new Error('奖惩金额必须填写大于 0 的整数')
      pageMessage.value = toErrorMessage(error, error.message)
      throw error
    }
    if (hasFractionInput(raw)) {
      const error = new Error('奖惩金额必须填写整数')
      pageMessage.value = toErrorMessage(error, error.message)
      throw error
    }
  }

  return {
    groups,
    generalNotices,
    adjustments,
    loading,
    sendingGeneral,
    sendingAdjustment,
    previewingAdjustment,
    voidingAdjustment,
    pageMessage,
    generalForm,
    adjustmentForm,
    bootstrap,
    refreshRecords,
    sendGeneral,
    sendAdjustmentNotice,
    previewAdjustmentNotice,
    previewVoidAdjustment,
    voidAdjustmentNotice,
  }
})

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  if (error instanceof Error && error.message) {
    return { type: 'error', text: error.message }
  }
  return { type: 'error', text: fallback }
}

function formatDateTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString('zh-CN', { hour12: false })
}

function formatAdjustmentType(value: AdjustmentType) {
  return value === 'REWARD' ? '奖励' : '罚款'
}
