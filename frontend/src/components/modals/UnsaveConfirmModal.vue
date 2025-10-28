<script setup lang="ts">
import { useTemplateRef } from 'vue'
import ModalFrame from './ModalFrame.vue'

const frame = useTemplateRef('frame')

defineExpose({
  show: (cb: typeof callback) => {
    callback = cb
    frame.value?.show()
  },
  hide: frame.value?.hide || (() => {})
})

let callback: (answer: 'yes' | 'no') => void
</script>

<template>
  <ModalFrame :on-demand="false" :immediate="false" ref="frame">
    <div class="max-w-[60vw]">
      <!-- Modal content -->
      <div class="rounded-lg bg-white shadow-sm">
        <!-- Modal header -->
        <div class="flex h-12 items-center justify-between rounded-t border-b px-4">
          <h3 class="font-semibold">
            {{ $t('common.unsaveConfirmTitle') }}
          </h3>

          <button
            type="button"
            class="rounded-lg bg-transparent p-3 text-sm text-gray-400 hover:text-gray-900"
            @click="
              () => {
                frame?.hide()
                callback('no')
              }
            "
          >
            <font-awesome-icon icon="fa-solid fa-xmark" />
          </button>
        </div>

        <!-- Modal body -->
        <div class="px-3 py-5">
          <p>
            {{ $t('common.unsaveConfirmMessage') }}
          </p>
        </div>

        <div class="flex gap-x-2 border-t px-4 py-2">
          <button
            type="button"
            class="btn flex-1 text-gray-700 btn-accent"
            @click="
              () => {
                frame?.hide()
                callback('yes')
              }
            "
          >
            {{ $t('common.confirm') }}
          </button>

          <button
            type="button"
            class="btn flex-1"
            @click="
              () => {
                frame?.hide()
                callback('no')
              }
            "
          >
            {{ $t('common.cancel') }}
          </button>
        </div>
      </div>
    </div>
  </ModalFrame>
</template>
