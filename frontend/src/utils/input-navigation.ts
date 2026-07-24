function isKeyboardComposing(event: KeyboardEvent) {
  return event.isComposing || event.key === 'Process'
}

function isVisibleInput(input: HTMLInputElement) {
  if (!input.isConnected || input.type === 'hidden' || input.closest('[hidden]')) {
    return false
  }

  const view = input.ownerDocument.defaultView
  if (view) {
    const style = view.getComputedStyle(input)
    if (style.display === 'none' || style.visibility === 'hidden') {
      return false
    }
  }

  return input.getClientRects().length > 0
}

function canFocusInput(input: HTMLInputElement) {
  return !input.disabled && !input.readOnly && isVisibleInput(input)
}

export function shouldHandleEnterNavigation(event: KeyboardEvent) {
  return !isKeyboardComposing(event)
}

export function canFocusEnterInput(input: HTMLInputElement | null | undefined): input is HTMLInputElement {
  return !!input && canFocusInput(input)
}

function isEnterNavigableInput(input: HTMLInputElement) {
  return input.matches('input[data-enter-nav]') && canFocusInput(input)
}

function isConfirmableNumericInput(input: HTMLInputElement) {
  return (
    input.matches('input[data-enter-confirm]')
    || input.type === 'number'
    || input.inputMode === 'numeric'
  ) && canFocusInput(input)
}

function selectInputText(input: HTMLInputElement) {
  try {
    input.select()
  } catch {
    // 部分输入类型不支持 select，聚焦即可。
  }
}

export function focusInputElement(input: HTMLInputElement | null | undefined) {
  if (!canFocusEnterInput(input)) {
    return false
  }

  input.focus()
  selectInputText(input)
  return true
}

export function confirmInputOnEnter(event: KeyboardEvent) {
  if (!shouldHandleEnterNavigation(event)) {
    return
  }

  const current = event.target
  if (!(current instanceof HTMLInputElement) || !isConfirmableNumericInput(current)) {
    return
  }

  event.preventDefault()
  current.blur()
}

export function focusNextInputOnEnter(event: KeyboardEvent) {
  if (!shouldHandleEnterNavigation(event)) {
    return
  }

  const current = event.target
  if (!(current instanceof HTMLInputElement) || !isEnterNavigableInput(current)) {
    return
  }

  event.preventDefault()

  const scope = current.closest('[data-enter-nav-scope]')
  if (!scope) {
    current.blur()
    return
  }

  const inputs = Array.from(scope.querySelectorAll<HTMLInputElement>('input[data-enter-nav]'))

  focusNextInputFromList(event, inputs)
}

export function focusNextInputFromList(event: KeyboardEvent, inputs: Array<HTMLInputElement | null | undefined>) {
  if (!shouldHandleEnterNavigation(event)) {
    return
  }

  const current = event.target
  if (!(current instanceof HTMLInputElement) || !canFocusEnterInput(current)) {
    return
  }

  event.preventDefault()

  const focusableInputs = inputs.filter(canFocusEnterInput)
  const currentIndex = focusableInputs.indexOf(current)
  const candidates = currentIndex >= 0 ? focusableInputs.slice(currentIndex + 1) : focusableInputs
  const next = candidates.find((input) => input !== current) ?? null

  if (!focusInputElement(next)) {
    current.blur()
  }
}
