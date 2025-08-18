<script setup lang="ts">
import CrossIcon from '@/components/icons/CrossIcon.vue'
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
  <ModalFrame :on-demand="true" :immediate="false" ref="frame">
    <div>
      <!-- Modal content -->
      <div class="bg-white rounded-lg shadow-sm">
        <!-- Modal header -->
        <div class="flex items-center justify-between h-12 px-4 border-b rounded-t">
          <h3 class="font-semibold">
            {{ t('porter.progress') }}
          </h3>

          <button
            v-show="progress?.status.includes('ed')"
            type="button"
            class="p-3 text-sm text-gray-400 hover:text-gray-900 bg-transparent hover:bg-gray-100 rounded-lg"
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
            <CrossIcon></CrossIcon>
          </button>
        </div>

        <!-- Modal body -->
        <div class="h-[70vh] w-[70vw] overflow-auto py-2 px-4">
          <div class="flex flex-col gap-y-2 h-full">
            <div class="flex items-center gap-x-3">
              <h2 class="text-lg font-bold">{{ title }}</h2>

              <p class="h-6 badge" :class="[`badge-${progress?.status}`]">
                <span class="truncate capitalize">{{ $t(`status.${progress?.status}`) }}</span>
              </p>
            </div>

            <ProgressBar :progresses="progress"></ProgressBar>
            <ol class="flex items-center w-full">
              <ProgressNode v-for="(progress, i) in progress?.tasks ?? []" :progress :key="i">
                <i class="text-xs lg:text-base">
                  <font-awesome-icon
                    v-if="progress.status == 'pending'"
                    icon="fa-solid fa-hourglass-start"
                  />
                  <font-awesome-icon
                    v-else-if="progress.status.includes('ing')"
                    icon="fa-solid fa-spinner"
                    spin
                  />
                  <font-awesome-icon
                    v-else-if="progress.status == 'completed'"
                    icon="fa-solid fa-check"
                  />
                  <font-awesome-icon v-else icon="fa-solid fa-exclamation" />
                </i>
              </ProgressNode>

              <ProgressNode></ProgressNode>
            </ol>

            <div
              class="flex flex-col flex-1 gap-y-2 overflow-y-auto min-h-48 p-1 border rounded-sm"
              ref="message-box"
            >
              <p v-for="(m, i) in messages" :key="i" class="text-xs text-gray-400 break-all">
                {{ m }}
              </p>
            </div>

            <div class="flex justify-end">
              <button
                v-show="progress?.status == 'pending' || progress?.status == 'running'"
                type="button"
                class="px-2 py-1 text-white bg-rose-600 rounded-sm"
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
