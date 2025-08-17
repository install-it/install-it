<script setup lang="ts">
import { statusBadget } from '@/definitions/styles'
import { porter } from '@/wailsjs/go/models'

defineProps<{ progress?: porter.Progress }>()
</script>

<template>
  <li class="flex items-center" :class="{ grow: progress !== undefined }">
    <span
      class="flex items-center justify-center h-7 md:h-9 lg:h-11 aspect-square rounded-full text-sm text-white"
      :class="[
        { 'animate-pulse': progress?.status.includes('ing') },
        statusBadget[progress?.status as keyof typeof statusBadget] || 'bg-gray-300'
      ]"
    >
      <slot>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          fill="currentColor"
          class="bi bi-list-check"
          viewBox="0 0 16 16"
        >
          <path
            fill-rule="evenodd"
            d="M5 11.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5m0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5m0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5M3.854 2.146a.5.5 0 0 1 0 .708l-1.5 1.5a.5.5 0 0 1-.708 0l-.5-.5a.5.5 0 1 1 .708-.708L2 3.293l1.146-1.147a.5.5 0 0 1 .708 0m0 4a.5.5 0 0 1 0 .708l-1.5 1.5a.5.5 0 0 1-.708 0l-.5-.5a.5.5 0 1 1 .708-.708L2 7.293l1.146-1.147a.5.5 0 0 1 .708 0m0 4a.5.5 0 0 1 0 .708l-1.5 1.5a.5.5 0 0 1-.708 0l-.5-.5a.5.5 0 0 1 .708-.708l.146.147 1.146-1.147a.5.5 0 0 1 .708 0"
          />
        </svg>
      </slot>
    </span>

    <div class="relative flex flex-col w-full text-center" v-if="progress !== undefined">
      <span class="text-xs lg:text-sm absolute -top-5 lg:-top-6 truncate w-full px-1">{{
        progress.name
      }}</span>

      <div class="w-full h-1.5 lg:h-2 bg-gray-200 rounded-full">
        <div
          class="h-full transition-all"
          :class="[statusBadget[progress.status as keyof typeof statusBadget]]"
          :style="{
            width: `${progress.total === 0 ? 0 : Math.floor((progress.current / progress.total) * 100)}%`
          }"
        ></div>
      </div>
    </div>
  </li>
</template>
