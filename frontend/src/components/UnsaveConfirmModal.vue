<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{ callback?: (confirm: boolean) => void }>()

const isOpen = ref(false)

let callback_: typeof props.callback = props.callback

defineExpose({
  show: (cb: typeof callback_) => {
    callback_ = cb
    isOpen.value = true
  },
  hide: () => {
    isOpen.value = false
  }
})
</script>

<template>
  <UModal v-model:open="isOpen" :title="$t('common.unsaveConfirmTitle')">
    <template #body>
      <p>{{ $t('common.unsaveConfirmMessage') }}</p>
    </template>

    <template #footer>
      <div class="flex gap-x-2">
        <UButton
          color="primary"
          class="flex-1 justify-center"
          @click="
            () => {
              isOpen = false
              if (callback_) {
                callback_(true)
              }
            }
          "
        >
          {{ $t('common.confirm') }}
        </UButton>

        <UButton
          variant="outline"
          color="neutral"
          class="flex-1 justify-center"
          @click="
            () => {
              isOpen = false
              if (callback_) {
                callback_(false)
              }
            }
          "
        >
          {{ $t('common.cancel') }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>
