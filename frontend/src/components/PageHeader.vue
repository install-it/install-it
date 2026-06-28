<script setup lang="ts">
defineProps<{
  title: string
  description?: string
  variant?: 'compact' | 'normal'
}>()
</script>

<template>
  <div v-if="variant === 'compact'" class="flex">
    <div class="grow">
      <h1 class="font-bold">{{ title }}</h1>

      <p v-if="description" class="text-xs">{{ description }}</p>
    </div>

    <slot />
  </div>

  <div v-else>
    <h1 class="text-lg font-bold">{{ title }}</h1>

    <p v-if="description" class="text-gray-400">{{ description }}</p>

    <hr class="mt-2 -mb-3" />
  </div>
</template>
