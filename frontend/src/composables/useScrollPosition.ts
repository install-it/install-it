import { onMounted, ref } from 'vue'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'

type ScrollRecord = {
  from: string
  to: string
  position: number
}

export function useScrollPosition(
  name: string,
  shouldApply?: (scroll: ScrollRecord | null) => boolean
) {
  const key = `scroll_${name}`
  const [route, router] = [useRoute(), useRouter()]
  const scrollContainer = ref<HTMLDivElement | null>(null)

  onMounted(() => {
    try {
      const item = sessionStorage.getItem(key)
      const record: ScrollRecord | null = item ? JSON.parse(item) : null

      if (
        (shouldApply && shouldApply(record)) ||
        (!shouldApply &&
          record?.to ===
            (router.options.history.state.forward ?? router.options.history.state.back) &&
          record?.from === route.fullPath)
      ) {
        scrollContainer.value!.scrollTop = record?.position ?? 0
      }
    } finally {
      sessionStorage.removeItem(key)
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
