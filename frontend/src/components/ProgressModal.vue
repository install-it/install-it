<script setup lang="ts">
import { porter } from '@/wailsjs/go/models'
import * as programPorter from '@/wailsjs/go/porter/Porter'
import * as runtime from '@/wailsjs/runtime'
import { nextTick, ref, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'

const isOpen = ref(false)

const mode = ref<'export' | 'import'>('export')

const snapshot = ref<porter.JobSnapshot | null>(null)

const messages = ref<string[]>([])

const title = ref<string>('')

const messageBox = useTemplateRef('message-box')

const { t } = useI18n()

const toast = useToast()

let interval: ReturnType<typeof setInterval> | number = -1

function resetInterval() {
  clearInterval(interval)
  interval = -1
}

function startPolling() {
  updateProgress()
  interval = setInterval(updateProgress, 300)
}

defineExpose({
  export: (destination: string) => {
    isOpen.value = true
    mode.value = 'export'
    title.value = t('porter.export')
    snapshot.value = null
    messages.value = []

    programPorter
      .Export(destination)
      .catch(toastErrMsg)
      .finally(() => {
        resetInterval()
        updateProgress()
      })

    startPolling()
  },
  download: (url: string): Promise<porter.ImportPreview> => {
    return new Promise((resolve, reject) => {
      isOpen.value = true
      mode.value = 'import'
      title.value = t('porter.download')
      snapshot.value = null
      messages.value = []

      programPorter
        .DownloadAndValidate(url)
        .then(preview => {
          resolve(preview)
          isOpen.value = false
        })
        .catch(err => {
          toastErrMsg(err)
          reject(err)
        })
        .finally(() => {
          resetInterval()
          updateProgress()
        })

      startPolling()
    })
  },
  import: (from: 'url' | 'file', source: string, opts: porter.ImportOptions) => {
    isOpen.value = true
    mode.value = 'import'
    title.value = `${t('porter.import')} (${t(`porter.${from}`)})`
    snapshot.value = null
    messages.value = []

    if (from === 'url') {
      programPorter
        .ImportFromURL(opts)
        .catch(toastErrMsg)
        .finally(() => {
          resetInterval()
          updateProgress()
        })
    } else {
      programPorter
        .ImportFromFile(source, opts)
        .catch(toastErrMsg)
        .finally(() => {
          resetInterval()
          updateProgress()
        })
    }

    startPolling()
  }
})

function updateProgress() {
  return programPorter.Progress().then(p => {
    let scroll = false
    if (
      messageBox.value &&
      messageBox.value.scrollTop + messageBox.value.clientHeight >=
        messageBox.value.scrollHeight * 0.99
    ) {
      scroll = true
    }

    snapshot.value = p
    messages.value.push(...p.messages.slice(messages.value.length).filter(m => m !== ''))

    if (scroll && messageBox.value) {
      nextTick(() => {
        messageBox.value!.scrollTop = messageBox.value!.scrollHeight
      })
    }
  })
}

function toastErrMsg(err: string) {
  if (err.includes('context canceled')) return
  if (err.includes('The system cannot find the path specified.'))
    toast.add({ title: t('toast.pathNotFind'), color: 'error' })
  else if (err.includes('unsupported protocol scheme'))
    toast.add({ title: t('toast.unsupportUrlProtocal'), color: 'error' })
  else if (err.includes('no such host'))
    toast.add({ title: t('toast.noSuchHost'), color: 'error' })
  else if (err == 'zip: not a valid zip file')
    toast.add({ title: t('toast.invalidZipFile'), color: 'error' })
  else toast.add({ title: err, color: 'error' })
}
</script>

<template>
  <UModal
    v-model:open="isOpen"
    :title="t('porter.progress')"
    :close="snapshot?.status.includes('ed') && !(mode == 'import' && snapshot?.status == 'completed')"
    :ui="{
      content: 'h-[80vh] max-h-150',
      body: 'flex flex-col overflow-hidden'
    }"
  >
    <template #body>
      <div class="flex min-h-0 flex-1 flex-col gap-y-4">
        <div class="flex items-center gap-x-3">
          <h2 class="text-lg font-bold">{{ title }}</h2>

          <p class="proc-badge h-6" :class="[`proc-badge-${snapshot?.status}`]">
            <span class="truncate capitalize">{{ $t(`status.${snapshot?.status}`) }}</span>
          </p>
        </div>

        <UProgress :model-value="(snapshot?.progress ?? 0) * 100" color="primary" />

        <p v-if="snapshot?.step" class="text-sm text-gray-600">
          {{ $t(`porter.step.${snapshot.step}`) }}
        </p>

        <div
          ref="message-box"
          class="flex flex-1 flex-col gap-y-2 overflow-y-auto rounded-sm border p-2"
        >
          <p v-for="(m, i) in messages" :key="i" class="text-xs break-all text-gray-400">
            {{ m }}
          </p>
        </div>

        <div class="flex justify-end gap-x-2">
          <div v-show="snapshot?.status == 'pending' || snapshot?.status == 'running'">
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

          <div v-show="mode == 'import' && snapshot?.status == 'completed'">
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
