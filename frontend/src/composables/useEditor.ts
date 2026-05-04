import { useUnsavedFormStore } from '@/store'
import { computed, ref, toRaw, toValue, type MaybeRefOrGetter, type Ref } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'

/**
 * Creates an editable copy of data with change tracking
 *
 * @param options - Configuration options
 * @param options.source - Original data source. Can be a plain value, Ref, ComputedRef, or getter function.
 *                         The editor will always track against the CURRENT source value.
 * @param options.equals - Custom equality check function (defaults to JSON.stringify comparison)
 * @param options.onBeforeReset - Callback invoked before reset. Useful for cleanup or confirmation logic
 * @param options.onAfterReset - Callback invoked after reset
 *
 * @returns Editor object with data, modified state, reset and updateSource functions
 *
 * @example
 * ```ts
 * // Static value (backward compatible)
 * const editor1 = useEditor({ source: { theme: 'dark', fontSize: 14 } })
 *
 * // Reactive ref
 * const settingsRef = ref({ theme: 'dark' })
 * const editor2 = useEditor({ source: settingsRef })
 *
 * // Computed ref
 * const computedSettings = computed(() => store.settings)
 * const editor3 = useEditor({ source: computedSettings })
 *
 * // Getter function (always fresh)
 * const editor4 = useEditor({ source: () => store.settings })
 *
 * editor1.data.value.theme = 'light'
 * console.log(editor1.modified.value) // true
 *
 * editor1.reset()
 * console.log(editor1.data.value.theme) // 'dark'
 * ```
 */
export function useEditor<T>(options: {
  /** Original data source - can be a value, Ref, ComputedRef, or getter function */
  source: MaybeRefOrGetter<T>
  /** Custom equality check function */
  equals?: (a: T, b: T) => boolean
  /** Callback invoked before reset */
  onBeforeReset?: () => void | Promise<void>
  /** Callback invoked after reset */
  onAfterReset?: () => void
  /** Enable unsaved changes warning on route leave */
  warnOnUnsavedLeave?: boolean
}) {
  const { source, equals, onBeforeReset, onAfterReset, warnOnUnsavedLeave } = options
  const data = ref<T>(structuredClone(toRaw(toValue(source)))) as Ref<T>

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

  const modified = computed(() => {
    return !equalityCheck(data.value, toValue(source))
  })

  if (warnOnUnsavedLeave) {
    const formStore = useUnsavedFormStore()

    onBeforeRouteLeave(() => {
      if (modified.value) {
        formStore.show = true
        return new Promise<boolean>(resolve => {
          formStore.setAnswerHandler(resolve)
        })
      }
      return true
    })
  }

  return {
    data,
    modified,
    reset: async () => {
      if (onBeforeReset) {
        await onBeforeReset()
      }

      data.value = structuredClone(toRaw(toValue(source)))

      if (onAfterReset) {
        onAfterReset()
      }
    },
    updateSource: (newSource: MaybeRefOrGetter<T>) => {
      data.value = structuredClone(toRaw(toValue(newSource)))
    }
  }
}
