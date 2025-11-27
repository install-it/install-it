<script setup lang="ts">
import { porter } from '@/wailsjs/go/models'
import { computed } from 'vue'

const props = defineProps<{ progress?: porter.Progress }>()

const inProgress = computed(
  () => props.progress?.status != 'pending' && props.progress?.status.includes('ing')
)
</script>

<template>
  <li class="flex items-center" :class="{ grow: progress !== undefined }">
    <span
      class="flex aspect-square h-7 items-center justify-center rounded-full border-4 border-gray-100 md:h-9 lg:h-11 lg:border-6"
      :style="[
        progress != undefined
          ? {
              'background-color': `var(--color-${progress.status})`,
              color: 'white'
            }
          : undefined
      ]"
    >
      <slot></slot>
    </span>

    <div v-if="progress !== undefined" class="relative flex w-full flex-col text-center">
      <span
        class="absolute -top-4.5 w-full truncate px-1 text-xs lg:-top-5.5 lg:text-sm"
        :class="{ 'text-gray-400': !inProgress }"
      >
        {{ $t(`porter.${progress.name}`) }}
      </span>

      <div class="h-1.5 w-full bg-gray-100 lg:h-2">
        <div
          class="h-full transition-all"
          :class="[{ 'animate-pulse': inProgress }]"
          :style="{
            width: `${progress.total === 0 ? 0 : Math.floor((progress.current / progress.total) * 100)}%`,
            'background-color': `var(--color-${progress.status})`
          }"
        ></div>
      </div>

      <span
        v-if="inProgress"
        class="absolute -bottom-4 w-full truncate px-1 text-xs text-gray-400 lg:-bottom-5"
      >
        {{ `${Math.floor((progress.current / progress.total) * 100)}%` }}
      </span>
    </div>
  </li>
</template>
