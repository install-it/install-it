<!-- eslint-disable vue/multi-word-component-names -->
<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{ onDemand: boolean; immediate: boolean }>()

defineExpose({
  show: () => {
    render.value = true

    setTimeout(() => (show.value = true), 50)
  },
  hide: () => {
    show.value = false

    if (props.onDemand) {
      setTimeout(() => (render.value = false), 500)
    }
  }
})

const show = ref(props.immediate || false)

const render = ref(props.onDemand)
</script>

<template>
  <TransitionGroup name="modal" v-if="render">
    <div class="fixed inset-0 z-40 bg-gray-900/50" v-show="show" key="0"></div>

    <div
      class="fixed top-0 right-0 left-0 z-50 flex h-full w-full items-center justify-center"
      v-show="show"
      key="1"
    >
      <slot></slot>
    </div>
  </TransitionGroup>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.1s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
