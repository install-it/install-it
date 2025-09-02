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
      class="flex items-center justify-center h-7 md:h-9 lg:h-11 aspect-square border-4 lg:border-6 border-gray-100 rounded-full"
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

    <div class="relative flex flex-col w-full text-center" v-if="progress !== undefined">
      <span
        class="text-xs lg:text-sm absolute -top-4.5 lg:-top-5.5 truncate w-full px-1"
        :class="{ 'text-gray-400': !inProgress }"
      >
        {{ $t(`porter.${progress.name}`) }}
      </span>

      <div class="w-full h-1.5 lg:h-2 bg-gray-100">
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
        class="text-xs absolute -bottom-4 lg:-bottom-5 truncate w-full px-1 text-gray-400"
        v-if="inProgress"
      >
        {{ `${Math.floor((progress.current / progress.total) * 100)}%` }}
      </span>
    </div>
  </li>
</template>
