<script setup lang="ts">
const loading = useLoading()
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="loading.isLoading.value"
        class="fixed inset-0 z-100 flex items-center justify-center bg-gray-700/50 backdrop-blur-xs"
      >
        <div class="flex flex-col items-center gap-4">
          <div
            class="size-10 animate-spin rounded-full border-4 border-gray-300 border-t-white"
          ></div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
