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
      <div class="bg-white rounded-lg shadow-sm">
        <!-- Modal header -->
        <div class="flex items-center justify-between h-12 px-4 border-b rounded-t">
          <h3 class="font-semibold">
            {{ $t('common.unsaveConfirmTitle') }}
          </h3>

          <button
            type="button"
            class="p-3 text-sm text-gray-400 hover:text-gray-900 bg-transparent rounded-lg"
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

        <div class="flex gap-x-2 h-12 py-2 px-4 border-t">
          <button
            type="button"
            class="flex-1 text-white bg-apple-green-800 hover:bg-apple-green-600 rounded-sm"
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
            class="flex-1 text-white bg-gray-400 hover:bg-gray-300 rounded-sm"
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
