import { ref } from 'vue'

const isLoading = ref(false)

export function useLoading() {
  return {
    show: () => {
      isLoading.value = true
      return {
        hide: () => {
          isLoading.value = false
        }
      }
    },
    hide: () => {
      isLoading.value = false
    },
    isLoading
  }
}
