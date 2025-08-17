<script setup lang="ts">
import { statusBadget } from '@/definitions/styles'
import { porter } from '@/wailsjs/go/models'
import ProgressNode from './ProgressNode.vue'

const props = defineProps<{ progresses: porter.Progresses | null }>()

// const tasks = computed(() => {
//   if (props.progresses === null) {
//     return []
//   }

//   if (props.progresses.tasks.length <= 3) {
//     return props.progresses.tasks
//   }

//   return []
// })
</script>

<template>
  <div class="flex flex-col gap-y-2">
    <div class="flex justify-between">
      <div class="shrink-0 w-[4.1rem]">
        <p
          class="inline-flex justify-center items-center max-w-[96%] h-6 px-1 rounded-sm"
          :class="[
            { 'animate-pulse': progresses?.status.includes('ing') },
            statusBadget[progresses?.status as keyof typeof statusBadget]
          ]"
        >
          <span class="text-sm truncate">{{ $t(`executeStatus.${progresses?.status}`) }}</span>
        </p>
      </div>
    </div>

    <ol class="flex items-center w-full">
      <ProgressNode
        v-for="(progress, i) in props.progresses?.tasks ?? []"
        :progress
        :key="i"
      ></ProgressNode>

      <ProgressNode>
        <svg
          class="w-4 h-4 text-gray-500 lg:w-5 lg:h-5 dark:text-gray-100"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          fill="currentColor"
          viewBox="0 0 18 20"
        >
          <path
            d="M16 1h-3.278A1.992 1.992 0 0 0 11 0H7a1.993 1.993 0 0 0-1.722 1H2a2 2 0 0 0-2 2v15a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V3a2 2 0 0 0-2-2ZM7 2h4v3H7V2Zm5.7 8.289-3.975 3.857a1 1 0 0 1-1.393 0L5.3 12.182a1.002 1.002 0 1 1 1.4-1.436l1.328 1.289 3.28-3.181a1 1 0 1 1 1.392 1.435Z"
          />
        </svg>
      </ProgressNode>
    </ol>
  </div>
</template>
