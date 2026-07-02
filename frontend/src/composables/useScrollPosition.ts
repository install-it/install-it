import { nextTick, onMounted, ref, watch } from 'vue'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'

export type ScrollRecord = {
  from: string
  to: string
  position: number
}

let previousRoute = ''
let hookInstalled = false

export function useScrollPosition(
  name: string,
  shouldApply?: (record: ScrollRecord | null) => boolean,
  source?: () => unknown
) {
  const key = `scroll_${name}`
  const route = useRoute()
  const router = useRouter()
  const scrollContainer = ref<HTMLDivElement | null>(null)

  if (!hookInstalled) {
    hookInstalled = true
    router.afterEach((_to, from) => {
      previousRoute = from.fullPath
    })
  }

  function restore(position: number) {
    if (scrollContainer.value) {
      scrollContainer.value.scrollTop = position
    }
  }

  onMounted(() => {
    const item = sessionStorage.getItem(key)
    if (!item) {
      return
    }

    sessionStorage.removeItem(key)

    let record: ScrollRecord | null = null
    try {
      record = JSON.parse(item)
    } catch {
      return
    }

    if (previousRoute !== record?.to || shouldApply?.(record) === false) {
      return
    }

    if (source) {
      watch(
        source,
        val => {
          if (Array.isArray(val) ? val.length > 0 : val) {
            restore(record!.position)
          }
        },
        { immediate: true }
      )
    } else {
      nextTick(() => restore(record!.position))
    }
  })

  onBeforeRouteLeave(to => {
    sessionStorage.setItem(
      key,
      JSON.stringify({
        from: route.fullPath,
        to: to.fullPath,
        position: scrollContainer.value?.scrollTop ?? 0
      })
    )
  })

  return { scrollContainer }
}
