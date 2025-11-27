<script setup lang="ts">
import ModalFrame from '@/components/modals/ModalFrame.vue'
import * as executor from '@/wailsjs/go/execute/CommandExecutor'
import { status } from '@/wailsjs/go/models'
import * as runtime from '@/wailsjs/runtime/runtime'
import AsyncLock from 'async-lock'
import { ref, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'vue-toast-notification'
import type { Command, Process } from '../types'
import TaskStatus from './TaskStatus.vue'

const frame = useTemplateRef('frame')

defineExpose({
  show: async (parallel: boolean, cmds: Array<Command>) => {
    frame.value?.show()

    isParallel = parallel

    processes.value = cmds.map(vals => ({ command: { ...vals }, status: status.Status.PENDING }))
    dispatchCommand()
  },
  hide: frame.value?.hide || (() => {})
})

const emit = defineEmits<{ completed: [] }>()

const { t } = useI18n()

const $toast = useToast({ position: 'top-left', duration: 7000 })

const lock = new AsyncLock()

let isParallel = false

const processes = ref<Array<Process>>([])

runtime.EventsOn('execute:exited', async (id: string, result: NonNullable<Process['result']>) => {
  const process = processes.value.find(c => c.procId === id)!
  process.result = result

  if (result.aborted) {
    process.status = status.Status.ABORTED
  } else if (![0, ...process.command.config.allowRtCodes].includes(result.exitCode)) {
    process.status = status.Status.FAILED
  } else if (result.lapse < process.command.config.minExeTime) {
    process.status = status.Status.SPEEDED
  } else {
    process.status = status.Status.COMPLETED
  }

  dispatchCommand().then(() => {
    if (processes.value.every(c => c.status === 'completed')) {
      emit('completed')
      $toast.success(t('toast.finished'), { position: 'bottom-right' })
    } else if (processes.value.every(c => !c.status.includes('ing'))) {
      $toast.info(t('toast.finished'), { position: 'bottom-right' })
    }
  })
})

function getProcessName(process: Process) {
  return process.command.name
    ? `${process.command.groupName} - ${process.command.name}`
    : process.command.groupName
}

async function dispatchCommand() {
  lock.acquire('executor', async () => {
    const pendings = processes.value
      .filter(c => c.status === 'pending')
      .slice(0, isParallel ? undefined : 1)

    for (const process of pendings) {
      if (
        !process.command.config.incompatibles.every(id =>
          processes.value.filter(p => p.status === 'running').every(p => p.command.id != id)
        )
      ) {
        continue
      }

      await executor
        .Run(process.command.config.program, process.command.config.options)
        .then(processId => {
          process.status = status.Status.RUNNING
          process.procId = processId
        })
        .catch(error => {
          process.status = status.Status.ERRORED
          process.result = {
            lapse: -1,
            exitCode: -1,
            stdout: '',
            stderr: '',
            error: (error as Error).toString(),
            aborted: false
          }
        })
    }
  })
}

async function handleAbort(process: Process) {
  return lock
    .acquire('executor', () => {
      if (process.status == 'pending' || process.status == 'running') {
        process.status =
          process.procId == undefined || process.procId == ''
            ? status.Status.ABORTED
            : status.Status.ABORTING
      }
    })
    .then(() => {
      if (process.status != 'aborting') {
        return
      }

      // `aborted` status will be updated at `execute:exited` event handler
      executor.Abort(process.procId!).catch(error => {
        if (error.includes('process does not exist')) {
          $toast.warning(
            t('toast.cancelCompletedFailed', {
              name: getProcessName(process)
            })
          )
          return
        }

        error
          .toString()
          .split('\n')
          .forEach((err: string) => {
            if (err.includes('abort failed')) {
              $toast.warning(
                t('toast.cancelFailed', {
                  name: getProcessName(process)
                })
              )
            } else {
              $toast.error(`[${getProcessName(process)}] ${err}`)
            }
          })

        process.status = status.Status.ERRORED
        process.result = {
          lapse: -1,
          exitCode: -1,
          stdout: '',
          stderr: '',
          error: error.toString(),
          aborted: false
        }
      })
    })
}
</script>

<template>
  <ModalFrame ref="frame" :on-demand="true" :immediate="false">
    <div class="w-[65vw] max-w-3xl">
      <!-- Modal content -->
      <div class="rounded-sm bg-white shadow-sm">
        <!-- Modal header -->
        <div class="flex items-center justify-between rounded-t border-b px-3 py-1.5">
          <h3 class="font-semibold">
            {{ $t('execute.title') }}
          </h3>
          <button
            type="button"
            class="ms-auto inline-flex h-8 w-8 items-center justify-center rounded-lg bg-transparent text-sm text-gray-400 enabled:hover:bg-gray-200 enabled:hover:text-gray-900"
            :disabled="
              processes.some(cmd => ['pending', 'running', 'aborting'].includes(cmd.status))
            "
            @click="frame?.hide()"
          >
            <font-awesome-icon icon="fa-solid fa-xmark" />
          </button>
        </div>

        <!-- Modal body -->
        <div class="max-h-[70vh] overflow-y-auto px-4 py-2">
          <template v-for="(process, i) in processes" :key="i">
            <TaskStatus :process="process" @abort="handleAbort(process)"></TaskStatus>
          </template>
        </div>

        <div
          v-show="
            processes.every(p => p.status.includes('ed')) &&
            processes.some(p => p.status != 'completed')
          "
          class="flex justify-end px-4 pb-2"
        >
          <button
            class="btn font-normal btn-sm btn-secondary"
            @click="
              event => {
                $emit('completed')
                $toast.success(t('toast.finished'), { position: 'bottom-right' })

                // @ts-ignore
                event.currentTarget?.remove()
              }
            "
          >
            <font-awesome-icon icon="fa-solid fa-forward" />
            {{ $t('execute.forceComplete') }}
          </button>
        </div>
      </div>
    </div>
  </ModalFrame>
</template>
