import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import useUnsavedFormStore from '@/stores/useUnsavedFormStore'

describe('useUnsavedFormStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('show is initially false', () => {
    const store = useUnsavedFormStore()
    expect(store.show).toBe(false)
  })

  it('confirmLeave(true) hides the modal and calls the handler with true', () => {
    const store = useUnsavedFormStore()
    const handler = vi.fn()
    store.setAnswerHandler(handler)
    store.show = true

    store.confirmLeave(true)

    expect(store.show).toBe(false)
    expect(handler).toHaveBeenCalledOnce()
    expect(handler).toHaveBeenCalledWith(true)
  })

  it('confirmLeave(false) hides the modal and calls the handler with false', () => {
    const store = useUnsavedFormStore()
    const handler = vi.fn()
    store.setAnswerHandler(handler)
    store.show = true

    store.confirmLeave(false)

    expect(store.show).toBe(false)
    expect(handler).toHaveBeenCalledOnce()
    expect(handler).toHaveBeenCalledWith(false)
  })

  it('confirmLeave without a registered handler does not throw', () => {
    const store = useUnsavedFormStore()
    expect(() => store.confirmLeave(true)).not.toThrow()
    expect(store.show).toBe(false)
  })

  it('handler is cleared after the first confirmLeave call', () => {
    const store = useUnsavedFormStore()
    const handler = vi.fn()
    store.setAnswerHandler(handler)

    store.confirmLeave(true)
    store.confirmLeave(true) // second call — handler should already be null

    expect(handler).toHaveBeenCalledTimes(1)
  })

  it('setAnswerHandler replaces a previously registered handler', () => {
    const store = useUnsavedFormStore()
    const first = vi.fn()
    const second = vi.fn()

    store.setAnswerHandler(first)
    store.setAnswerHandler(second)
    store.confirmLeave(true)

    expect(first).not.toHaveBeenCalled()
    expect(second).toHaveBeenCalledOnce()
  })

  it('show can be set to true externally', () => {
    const store = useUnsavedFormStore()
    store.show = true
    expect(store.show).toBe(true)
  })
})
