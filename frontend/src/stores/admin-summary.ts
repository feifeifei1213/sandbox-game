import { defineStore } from 'pinia'
import { ref } from 'vue'

import { getAdminFinalRanking, getAdminYearSummary } from '@/api/sandbox-game/admin-summary'
import type {
  AdminFinalRankingResult,
  AdminYearSummaryResult,
} from '@/types/sandbox-game-admin'
import type { PageMessage } from '@/stores/admin-shell'

export const useAdminSummaryStore = defineStore('sandbox-admin-summary', () => {
  const yearSummaries = ref<AdminYearSummaryResult[]>([])
  const finalRanking = ref<AdminFinalRankingResult | null>(null)
  const rankingUnavailableReason = ref('')
  const loading = ref(false)
  const pageMessage = ref<PageMessage | null>(null)

  async function bootstrap(finalYear: number) {
    loading.value = true
    pageMessage.value = null
    try {
      if (finalYear < 1) {
        yearSummaries.value = []
        finalRanking.value = null
        rankingUnavailableReason.value = '当前尚未设置正式年份，年度汇总区暂不可展示。'
        return
      }

      const formalYears = Array.from({ length: finalYear }, (_, index) => index + 1)
      yearSummaries.value = await Promise.all(formalYears.map((yearNo) => getAdminYearSummary(yearNo)))

      try {
        finalRanking.value = await getAdminFinalRanking()
        rankingUnavailableReason.value = ''
      } catch (error) {
        finalRanking.value = null
        rankingUnavailableReason.value = toErrorText(error, '最终排名暂未开放。')
      }
    } catch (error) {
      pageMessage.value = toErrorMessage(error, '读取管理员汇总失败')
      throw error
    } finally {
      loading.value = false
    }
  }

  return {
    yearSummaries,
    finalRanking,
    rankingUnavailableReason,
    loading,
    pageMessage,
    bootstrap,
  }
})

function toErrorMessage(error: unknown, fallback: string): PageMessage {
  return {
    type: 'error',
    text: toErrorText(error, fallback),
  }
}

function toErrorText(error: unknown, fallback: string) {
  if (error instanceof Error && error.message) {
    return error.message
  }
  return fallback
}
