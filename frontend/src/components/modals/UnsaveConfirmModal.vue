<script setup lang="ts">
import { ref } from 'vue'

const isOpen = ref(false)
let callback: (answer: 'yes' | 'no') => void

defineExpose({
  show: (cb: typeof callback) => {
    callback = cb
    isOpen.value = true
  },
  hide: () => {
    isOpen.value = false
  }
})
</script>

<template>
  <UModal v-model:open="isOpen" title="$t('common.unsaveConfirmTitle')">
    <template #body>
      <p>{{ $t('common.unsaveConfirmMessage') }}</p>
    </template>

    <template #footer>
      <div class="flex gap-x-2">
        <UButton
          color="primary"
          class="flex-1"
          @click="
            () => {
              isOpen = false
              callback('yes')
            }
          "
        >
          {{ $t('common.confirm') }}
        </UButton>

        <UButton
          variant="outline"
          color="neutral"
          class="flex-1"
          @click="
            () => {
              isOpen = false
              callback('no')
            }
          "
        >
          {{ $t('common.cancel') }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>
