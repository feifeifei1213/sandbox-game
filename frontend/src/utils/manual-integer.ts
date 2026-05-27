export function isIntegerLikeInput(value: unknown) {
  if (value === '' || value === null || value === undefined) {
    return true
  }
  const text = String(value).trim()
  if (text === '') {
    return true
  }
  const parsed = Number(text)
  return Number.isFinite(parsed) && Number.isInteger(parsed)
}

export function parseIntegerInput(value: string, fallback: number | null = null) {
  const text = value.trim()
  if (text === '') {
    return fallback
  }
  const parsed = Number(text)
  if (!Number.isFinite(parsed)) {
    return fallback
  }
  return Number.isInteger(parsed) ? parsed : parsed
}

export const hasFractionInput = (value: unknown) => !isIntegerLikeInput(value)
