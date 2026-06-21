<script setup lang="ts">
import { porter } from '@/wailsjs/go/models'
import * as programPorter from '@/wailsjs/go/porter/Porter'
import * as runtime from '@/wailsjs/runtime'
import { nextTick, ref, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'

const isOpen = ref(false)

const mode = ref<'export' | 'import'>('export')

defineExpose({
  export: (destination: string) => {
    isOpen.value = true
    mode.value = 'export'
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
    isOpen.value = true
    mode.value = 'import'
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

const messageBox = useTemplateRef('message-box')

const { t } = useI18n()

const toast = useToast()

let interval: ReturnType<typeof setInterval> | number = -1

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
    toast.add({ title: t('toast.pathNotFind'), color: 'error' })
  } else if (err.includes('unsupported protocol scheme')) {
    toast.add({ title: t('toast.unsupportUrlProtocal'), color: 'error' })
  } else if (err.includes('no such host')) {
    toast.add({ title: t('toast.noSuchHost'), color: 'error' })
  } else if (err == 'zip: not a valid zip file') {
    toast.add({ title: t('toast.invalidZipFile'), color: 'error' })
  } else {
    toast.add({ title: err, color: 'error' })
  }
}
</script>

<template>
  <UModal
    v-model:open="isOpen"
    :title="t('porter.progress')"
    :close="progress?.status.includes('ed') && !(mode == 'import' && progress?.status == 'completed')"
    :ui="{
      content: 'h-[80vh] max-h-150',
      body: 'flex flex-col overflow-hidden'
    }"
  >
    <template #body>
      <div class="flex min-h-0 flex-1 flex-col gap-y-2">
        <div class="flex items-center gap-x-3">
          <h2 class="text-lg font-bold">{{ title }}</h2>

          <p class="proc-badge h-6" :class="[`proc-badge-${progress?.status}`]">
            <span class="truncate capitalize">{{ $t(`status.${progress?.status}`) }}</span>
          </p>
        </div>

        <ol class="flex w-full items-center">
          <ProgressNode v-for="(task, i) in progress?.tasks ?? []" :key="i" :progress="task">
            <i class="text-xs lg:text-base">
              <Icon v-if="task.status == 'pending'" icon="mdi:hourglass" />

              <Icon
                v-else-if="task.status.includes('ing')"
                icon="mdi:loading"
                class="animate-spin"
              />

              <Icon v-else-if="task.status == 'completed'" icon="mdi:check" />

              <Icon v-else icon="mdi:alert" />
            </i>
          </ProgressNode>

          <ProgressNode>
            <i class="text-xs lg:text-base">
              <Icon icon="mdi:goal" />
            </i>
          </ProgressNode>
        </ol>

        <div
          ref="message-box"
          class="flex flex-1 flex-col gap-y-2 overflow-y-auto rounded-sm border p-1"
        >
          <p v-for="(m, i) in messages" :key="i" class="text-xs break-all text-gray-400">
            {{ m }}
          </p>
        </div>

        <div class="flex justify-end gap-x-2">
          <div v-show="progress?.status == 'pending' || progress?.status == 'running'">
            <UButton
              type="button"
              color="error"
              @click="
                () => {
                  programPorter.Abort().catch(err => toast.add({ title: err, color: 'error' }))
                }
              "
            >
              {{ $t('common.cancel') }}
            </UButton>
          </div>

          <div v-show="mode == 'import' && progress?.status == 'completed'">
            <UButton
              type="button"
              color="primary"
              @click="
                () => {
                  isOpen = false
                  runtime.WindowReloadApp()
                }
              "
            >
              {{ $t('common.refresh') }}
            </UButton>
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>
