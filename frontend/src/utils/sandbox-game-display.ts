const FALLBACK_TEXT = '--'

const yearTabStatusLabels: Record<string, string> = {
  ENTERABLE: '可进入',
  LOCKED: '未开放',
  COMPLETED: '已完成',
  BANKRUPT_READONLY: '破产只读',
  ROLLBACK_PENDING: '待重提',
}

const yearStatusLabels: Record<string, string> = {
  LOCKED: '已锁定',
  OPERATING: '经营中',
  REPORT_PENDING: '待填财报',
  REPORTING: '财报填写中',
  COMPLETED: '本年完成',
}

const reportStatusLabels: Record<string, string> = {
  REPORT_LOCKED: '未开放',
  REPORT_OPEN: '可填写',
  REPORT_SUBMITTED: '已提交',
}

const businessStatusLabels: Record<string, string> = {
  NORMAL: '正常',
  BANKRUPT: '已破产',
}

const stageCodeLabels: Record<string, string> = {
  YEAR_START: '年初',
  Q1: '第一季度',
  Q2: '第二季度',
  Q3: '第三季度',
  Q4: '第四季度',
  YEAR_END: '年末',
}

const stageStatusLabels: Record<string, string> = {
  Q1_OPEN: '第一季度',
  Q2_OPEN: '第二季度',
  Q3_OPEN: '第三季度',
  Q4_OPEN: '第四季度',
  YEAR_END_OPEN: '年末',
}

function formatMappedValue(value: string | null | undefined, labels: Record<string, string>) {
  if (!value) {
    return FALLBACK_TEXT
  }
  return labels[value] ?? value
}

export function formatYearTabStatus(value: string | null | undefined) {
  return formatMappedValue(value, yearTabStatusLabels)
}

export function formatYearStatus(value: string | null | undefined) {
  return formatMappedValue(value, yearStatusLabels)
}

export function formatReportStatus(value: string | null | undefined) {
  return formatMappedValue(value, reportStatusLabels)
}

export function formatBusinessStatus(value: string | null | undefined) {
  return formatMappedValue(value, businessStatusLabels)
}

export function formatStageCode(value: string | null | undefined) {
  return formatMappedValue(value, stageCodeLabels)
}

export function formatStageStatus(value: string | null | undefined) {
  return formatMappedValue(value, stageStatusLabels)
}
