<script setup lang="ts">
import ModalFrame from '@/components/modals/ModalFrame.vue'
import { porter } from '@/wailsjs/go/models'
import * as programPorter from '@/wailsjs/go/porter/Porter'
import * as runtime from '@/wailsjs/runtime'
import { nextTick, ref, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'vue-toast-notification'
import ProgressNode from './ProgressNode.vue'

defineExpose({
  export: (destination: string) => {
    frame.value?.show()
    title.value = t('porter.export')
    progress.value = null
    messages.value = []

    programPorter
      .Export(destination)
      .catch(toastErrMsg)
      .finally(() => {
        clearInterval(interval)
        updateProgress()
      })

    updateProgress()
    interval = setInterval(updateProgress, 300)
  },
  import: (from: 'url' | 'file', source: string) => {
    frame.value?.show()
    title.value = `${t('porter.import')} (${t(`porter.${from}`)})`
    progress.value = null
    messages.value = []

    if (from == 'url') {
      programPorter
        .ImportFromURL(source)
        .catch(toastErrMsg)
        .finally(() => {
          clearInterval(interval)
          updateProgress()
        })
    } else {
      programPorter
        .ImportFromFile(source)
        .catch(toastErrMsg)
        .finally(() => {
          clearInterval(interval)
          updateProgress()
        })
    }

    updateProgress()
    interval = setInterval(updateProgress, 300)
  }
})

const frame = useTemplateRef('frame')

const messageBox = useTemplateRef('message-box')

const { t } = useI18n()

const $toast = useToast({ position: 'top-right' })

let interval = -1

const messages = ref<Array<string>>([])

const progress = ref<porter.Progresses | null>(null)

const title = ref<string>()

function updateProgress() {
  return programPorter.Progress().then(p => {
    let scroll = false
    if (
      messageBox.value!.scrollTop + messageBox.value!.clientHeight >=
      messageBox.value!.scrollHeight * 0.99
    ) {
      scroll = true
    }

    progress.value = p
    messages.value.push(...p.messages.filter(m => m !== ''))

    if (scroll) {
      nextTick(() => {
        messageBox.value!.scrollTop = messageBox.value!.scrollHeight
      })
    }
  })
}

function toastErrMsg(err: string) {
  if (err.includes('context canceled')) {
    return
  } else if (err.includes('The system cannot find the path specified.')) {
    $toast.error(t('toast.pathNotFind'))
  } else if (err.includes('unsupported protocol scheme')) {
    $toast.error(t('toast.unsupportUrlProtocal'))
  } else if (err.includes('no such host')) {
    $toast.error(t('toast.noSuchHost'))
  } else if (err == 'zip: not a valid zip file') {
    $toast.error(t('toast.invalidZipFile'))
  } else {
    $toast.error(err)
  }
}
</script>

<template>
  <ModalFrame ref="frame" :on-demand="true" :immediate="false">
    <div>
      <!-- Modal content -->
      <div class="rounded-lg bg-white shadow-sm">
        <!-- Modal header -->
        <div class="flex h-12 items-center justify-between rounded-t border-b px-4">
          <h3 class="font-semibold">
            {{ t('porter.progress') }}
          </h3>

          <button
            v-show="progress?.status.includes('ed')"
            type="button"
            class="rounded-lg bg-transparent p-3 text-sm text-gray-400 hover:bg-gray-100 hover:text-gray-900"
            @click="
              () => {
                if (progress?.status == 'completed') {
                  runtime.WindowReloadApp()
                } else {
                  frame?.hide()
                }
              }
            "
          >
            <font-awesome-icon icon="fa-solid fa-xmark" />
          </button>
        </div>

        <!-- Modal body -->
        <div class="h-[70vh] w-[70vw] overflow-auto px-4 py-2">
          <div class="flex h-full flex-col gap-y-2">
            <div class="flex items-center gap-x-3">
              <h2 class="text-lg font-bold">{{ title }}</h2>

              <p class="proc-badge h-6" :class="[`proc-badge-${progress?.status}`]">
                <span class="truncate capitalize">{{ $t(`status.${progress?.status}`) }}</span>
              </p>
            </div>

            <ol class="flex w-full items-center">
              <ProgressNode v-for="(task, i) in progress?.tasks ?? []" :key="i" :task>
                <i class="text-xs lg:text-base">
                  <font-awesome-icon
                    v-if="task.status == 'pending'"
                    icon="fa-solid fa-hourglass-start"
                  />

                  <font-awesome-icon
                    v-else-if="task.status.includes('ing')"
                    icon="fa-solid fa-spinner"
                    spin
                  />

                  <font-awesome-icon
                    v-else-if="task.status == 'completed'"
                    icon="fa-solid fa-check"
                  />

                  <font-awesome-icon v-else icon="fa-solid fa-exclamation" />
                </i>
              </ProgressNode>

              <ProgressNode></ProgressNode>
            </ol>

            <div
              ref="message-box"
              class="flex min-h-48 flex-1 flex-col gap-y-2 overflow-y-auto rounded-sm border p-1"
            >
              <p v-for="(m, i) in messages" :key="i" class="text-xs break-all text-gray-400">
                {{ m }}
              </p>
            </div>

            <div class="flex justify-end">
              <button
                v-show="progress?.status == 'pending' || progress?.status == 'running'"
                type="button"
                class="btn btn-error"
                @click="
                  () => {
                    programPorter.Abort().catch(err => $toast.error(err))
                  }
                "
              >
                {{ $t('common.cancel') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </ModalFrame>
</template>
