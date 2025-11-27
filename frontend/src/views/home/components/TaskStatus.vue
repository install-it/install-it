<script setup lang="ts">
import type { Process } from '../types'

const props = defineProps<{ process: Process }>()

defineEmits<{ abort: [] }>()
</script>

<template>
  <div class="flex min-h-9 border-t border-kashmir-blue-100 last:border-b">
    <div class="w-2/6 content-center truncate pe-1 text-xs">
      <p class="truncate font-medium">{{ props.process.command.groupName }}</p>

      <p v-if="props.process.command.name" class="truncate">
        &nbsp;&nbsp;{{ `⤷ ${props.process.command.name}` }}
      </p>
    </div>

    <div class="flex w-4/6 items-center py-1 ps-1">
      <!-- status badget -->
      <div class="w-[4.1rem] shrink-0">
        <p class="proc-badge h-6" :class="[`proc-badge-${props.process.status}`]">
          <span class="truncate text-sm">
            {{ $t(`status.@${props.process.status}`).toUpperCase() }}
          </span>
        </p>
      </div>

      <!-- messages -->
      <template v-if="props.process.status == 'speeded' || props.process.status == 'failed'">
        <div class="line-clamp-3 text-sm break-all">
          {{ $t('execute.exitCode', { code: props.process.result?.exitCode }) }}

          <p v-if="props.process.status == 'speeded'" class="text-xs text-orange-300">
            {{
              $t('execute.earlyExit', {
                second: `${(props.process.result?.lapse ?? -1).toFixed(1)}/${props.process.command.config.minExeTime}`
              })
            }}
          </p>

          <p
            v-else-if="
              props.process.result &&
              props.process.result.error !== '' &&
              !props.process.result.error.includes('exit status')
            "
            class="font-mono text-xs text-red-400"
          >
            {{
              props.process.result.error.includes('file does not exist') ||
              props.process.result.error.includes('The system cannot find the file specified.') ||
              props.process.result.error.includes('The system cannot find the path specified.')
                ? $t('execute.fileNotExist')
                : props.process.result.error.split(':').slice(1).join(':').trim()
            }}
          </p>

          <p v-else class="font-mono text-xs text-red-400">
            {{ props.process.result?.stderr || props.process.result?.stdout }}
          </p>
        </div>
      </template>

      <template v-else-if="props.process.status == 'errored'">
        <div class="line-clamp-2 font-mono text-sm break-all">
          {{
            props.process.result?.error?.split(':').slice(1).join(':').trim() ??
            $t('execute.startFailed')
          }}
        </div>
      </template>

      <template v-else-if="props.process.status == 'completed'">
        <div class="text-xs text-gray-300">
          <p class="truncate">
            {{ $t('execute.exitCode', { code: props.process.result?.exitCode }) }}
          </p>

          <p class="truncate">
            {{
              $t('execute.executeTime', { second: Math.round(props.process.result?.lapse ?? -1) })
            }}
          </p>
        </div>
      </template>

      <!-- abort button -->
      <button
        v-show="props.process.status == 'pending' || props.process.status == 'running'"
        class="btn ms-auto font-normal btn-xs"
        @click="$emit('abort')"
      >
        {{ $t('execute.abort') }}
      </button>
    </div>
  </div>
</template>
