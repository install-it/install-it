import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ref } from 'vue'
import { setActivePinia, createPinia } from 'pinia'
import { useEditor } from '@/composables/useEditor'

// useEditor uses onBeforeRouteLeave from vue-router inside a guarded branch
// (warnOnUnsavedLeave). We never set that flag to true in these tests, so
// vue-router's inject context is never required.

describe('useEditor', () => {
  beforeEach(() => {
    // Some tests instantiate useUnsavedFormStore indirectly; keep pinia ready.
    setActivePinia(createPinia())
  })

  // ─── initial state ────────────────────────────────────────────────────────

  it('data starts as a deep clone of the source', () => {
    const source = { name: 'test', count: 1 }
    const editor = useEditor({ source })
    expect(editor.data.value).toEqual(source)
    // Must be a clone, not the same reference
    expect(editor.data.value).not.toBe(source)
  })

  it('data initialised from a ref source is also a deep clone', () => {
    const sourceRef = ref({ value: 42 })
    const editor = useEditor({ source: sourceRef })
    expect(editor.data.value).toEqual({ value: 42 })
    expect(editor.data.value).not.toBe(sourceRef.value)
  })

  it('modified is false immediately after creation', () => {
    const editor = useEditor({ source: { name: 'original' } })
    expect(editor.modified.value).toBe(false)
  })

  // ─── tracking ─────────────────────────────────────────────────────────────

  it('modified becomes true when data diverges from source', () => {
    const source = { name: 'original' }
    const editor = useEditor({ source })
    editor.data.value.name = 'changed'
    expect(editor.modified.value).toBe(true)
  })

  it('modified returns to false when data is manually restored to match source', () => {
    const source = { name: 'original' }
    const editor = useEditor({ source })
    editor.data.value.name = 'changed'
    expect(editor.modified.value).toBe(true)
    editor.data.value.name = 'original'
    expect(editor.modified.value).toBe(false)
  })

  it('modified tracks against the current value of a ref source', () => {
    const sourceRef = ref({ value: 42 })
    const editor = useEditor({ source: sourceRef })

    editor.data.value.value = 99
    expect(editor.modified.value).toBe(true)

    // Mutate the source ref to match the edited data
    sourceRef.value.value = 99
    expect(editor.modified.value).toBe(false)
  })

  // ─── reset ────────────────────────────────────────────────────────────────

  it('reset() restores data to the current source value', async () => {
    const source = { name: 'original' }
    const editor = useEditor({ source })
    editor.data.value.name = 'changed'

    await editor.reset()

    expect(editor.data.value).toEqual(source)
    expect(editor.modified.value).toBe(false)
  })

  it('reset() invokes onBeforeReset before restoring', async () => {
    const onBeforeReset = vi.fn()
    const editor = useEditor({ source: { x: 1 }, onBeforeReset })
    editor.data.value.x = 99

    await editor.reset()

    expect(onBeforeReset).toHaveBeenCalledOnce()
  })

  it('reset() invokes onAfterReset after restoring', async () => {
    const order: string[] = []
    const onBeforeReset = vi.fn(() => {
      order.push('before')
    })
    const onAfterReset = vi.fn(() => {
      order.push('after')
    })
    const editor = useEditor({ source: { x: 1 }, onBeforeReset, onAfterReset })

    await editor.reset()

    expect(order).toEqual(['before', 'after'])
  })

  it('reset() waits for an async onBeforeReset before restoring', async () => {
    const calls: string[] = []
    const onBeforeReset = vi.fn(async () => {
      await new Promise<void>(resolve => setTimeout(resolve, 0))
      calls.push('before')
    })
    const onAfterReset = vi.fn(() => {
      calls.push('after')
    })
    const editor = useEditor({ source: { x: 1 }, onBeforeReset, onAfterReset })

    await editor.reset()

    expect(calls).toEqual(['before', 'after'])
  })

  it('reset() against a ref source uses the current ref value', async () => {
    const sourceRef = ref({ label: 'v1' })
    const editor = useEditor({ source: sourceRef })
    editor.data.value.label = 'edited'

    // Advance the source before resetting
    sourceRef.value.label = 'v2'
    await editor.reset()

    expect(editor.data.value.label).toBe('v2')
  })

  // ─── updateSource ─────────────────────────────────────────────────────────

  it('updateSource() replaces data with a deep clone of the new value', () => {
    const source = { name: 'original' }
    const editor = useEditor({ source })

    editor.updateSource({ name: 'replaced' })

    expect(editor.data.value.name).toBe('replaced')
    // The new data must be a fresh clone, not the object passed in
    const newSource = { name: 'replaced' }
    editor.updateSource(newSource)
    expect(editor.data.value).not.toBe(newSource)
  })

  it('updateSource() does not alter the original source closure', () => {
    const source = { name: 'original' }
    const editor = useEditor({ source })
    editor.updateSource({ name: 'new-data' })

    // toValue(source) is still 'original', data.value is 'new-data' → modified
    expect(editor.modified.value).toBe(true)
  })

  it('updateSource() accepts a ref as the new source value', () => {
    const source = { count: 0 }
    const editor = useEditor({ source })
    const newRef = ref({ count: 5 })

    editor.updateSource(newRef)

    expect(editor.data.value.count).toBe(5)
  })

  // ─── custom equals ────────────────────────────────────────────────────────

  it('custom equals function is used instead of JSON.stringify comparison', () => {
    type Item = { id: number; name: string }
    const source: Item = { id: 1, name: 'original' }
    // Only compare by id — ignore the name field
    const idOnlyEquals = (a: Item, b: Item) => a.id === b.id

    const editor = useEditor({ source, equals: idOnlyEquals })
    editor.data.value.name = 'completely different'

    // Despite the name change, the custom equals still considers them equal
    expect(editor.modified.value).toBe(false)
  })

  it('custom equals reflecting id change marks the editor as modified', () => {
    type Item = { id: number; name: string }
    const source: Item = { id: 1, name: 'original' }
    const idOnlyEquals = (a: Item, b: Item) => a.id === b.id

    const editor = useEditor({ source, equals: idOnlyEquals })
    editor.data.value.id = 99

    expect(editor.modified.value).toBe(true)
  })

  // ─── getter source ────────────────────────────────────────────────────────

  it('accepts a getter function as source', () => {
    let raw = { value: 10 }
    const editor = useEditor({ source: () => raw })

    expect(editor.data.value).toEqual({ value: 10 })
    expect(editor.modified.value).toBe(false)
  })

  it('modified reacts when the reactive ref read inside a getter changes', () => {
    // The getter must read reactive data for computed() to track it.
    // A plain variable reassignment is invisible to Vue's reactivity system.
    const inner = ref({ value: 10 })
    const editor = useEditor({ source: () => inner.value })

    editor.data.value.value = 20
    expect(editor.modified.value).toBe(true)

    // Mutate the reactive value that the getter reads — computed() re-evaluates
    inner.value.value = 20
    expect(editor.modified.value).toBe(false)
  })
})
