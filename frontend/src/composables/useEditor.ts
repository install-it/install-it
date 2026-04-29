import { computed, ref, toRaw, type Ref } from 'vue'

/**
 * Creates an editable copy of data with change tracking
 *
 * @param options - Configuration options
 * @param options.source - Original data source. The editor will clone this value for editing
 * @param options.equals - Custom equality check function (defaults to JSON.stringify comparison)
 * @param options.onBeforeReset - Callback invoked before reset. Useful for cleanup or confirmation logic
 * @param options.onAfterReset - Callback invoked after reset
 *
 * @returns Editor object with data, modified state, reset and updateSource functions
 *
 * @example
 * ```ts
 * const editor = useEditor({ source: { theme: 'dark', fontSize: 14 } })
 *
 * editor.data.value.theme = 'light'
 * console.log(editor.modified.value) // true
 *
 * editor.reset()
 * console.log(editor.data.value.theme) // 'dark'
 * ```
 */
export function useEditor<T>(options: {
  /** Original data source */
  source: T
  /** Custom equality check function */
  equals?: (a: T, b: T) => boolean
  /** Callback invoked before reset */
  onBeforeReset?: () => void | Promise<void>
  /** Callback invoked after reset */
  onAfterReset?: () => void
}) {
  const { source, equals, onBeforeReset, onAfterReset } = options

  // Create a deep clone of the source
  const data = ref<T>(structuredClone(toRaw(source))) as Ref<T>

  const equalityCheck =
    equals ??
    ((a: T, b: T): boolean => {
      try {
        return JSON.stringify(a) === JSON.stringify(b)
      } catch (error) {
        console.warn('Editor equality check failed, falling back to false:', error)
        return false
      }
    })

  // Compute modified state
  const modified = computed(() => {
    return !equalityCheck(data.value, source)
  })

  // Reset function
  const reset = async () => {
    if (onBeforeReset) {
      await onBeforeReset()
    }

    data.value = structuredClone(toRaw(source))

    if (onAfterReset) {
      onAfterReset()
    }
  }

  // Allow updating the source (useful when ID changes)
  const updateSource = (newSource: T) => {
    data.value = structuredClone(toRaw(newSource))
  }

  return {
    data,
    modified,
    reset,
    updateSource
  }
}
