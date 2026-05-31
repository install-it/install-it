<script setup lang="ts">
import { TriggerNativeUpdate } from '@/wailsjs/go/update/Updater'
import { Quit } from '@/wailsjs/runtime/runtime'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { computed, ref, watch } from 'vue'

defineProps<{ currentVersion: string }>()

interface UpdateCheckResult {
  hasUpdate: boolean
  latestVersion: string
  downloadUrl: string
  downloadUrlBundled: string
  releaseNotes: string
  releaseAt: string
}

const isOpen = ref(false)
const updateResult = ref<UpdateCheckResult>()
const parsedNotes = ref('')
const webviewVersion = ref(false)

const releaseAt = ref('')

defineExpose({
  show: (result: UpdateCheckResult) => {
    updateResult.value = result
    webviewVersion.value = !!result.downloadUrlBundled
    isOpen.value = true
  },
  hide: () => {
    isOpen.value = false
  }
})

const toast = useToast()
const $loading = useLoading()

const selectedUrl = computed(() => {
  if (!updateResult.value) return ''
  return webviewVersion.value && updateResult.value.downloadUrlBundled
    ? updateResult.value.downloadUrlBundled
    : updateResult.value.downloadUrl
})

watch(
  updateResult,
  async result => {
    if (!result) {
      parsedNotes.value = ''
      releaseAt.value = ''
      return
    }

    releaseAt.value = new Date(result.releaseAt).toLocaleDateString()

    if (result.releaseNotes) {
      const html = await marked.parse(result.releaseNotes)
      parsedNotes.value = DOMPurify.sanitize(html)
    } else {
      parsedNotes.value = ''
    }
  },
  { deep: true }
)
</script>

<template>
  <UModal v-model:open="isOpen" :title="$t('info.updateInfoTitle')">
    <template #body>
      <div class="flex flex-col gap-y-3">
        <div class="flex grow flex-col gap-y-2">
          <div class="flex">
            <h1 class="min-w-34 font-medium">
              {{ $t('info.currentVersion') }}
            </h1>

            <p>{{ currentVersion }}</p>
          </div>

          <div class="flex">
            <h1 class="min-w-34 font-medium">
              {{ $t('info.latestVersion') }}
            </h1>

            <p>
              {{ `${updateResult?.latestVersion} (${releaseAt})` }}
            </p>
          </div>

          <hr />

          <div class="flex grow flex-col">
            <h1 class="mb-1 min-w-32 font-medium">
              {{ $t('info.updateInfo') }}
            </h1>

            <div
              id="release-notes"
              class="rounded-lg border px-1"
              v-html="parsedNotes || `<i>${$t('info.noUpdateInfo')}</i>`"
            ></div>
          </div>

          <hr />

          <div class="flex flex-col">
            <h1 class="font-medium">
              {{ $t('info.updateOption') }}
            </h1>

            <label class="flex w-full cursor-pointer items-center select-none">
              <UCheckbox v-model="webviewVersion" name="webview_version" color="primary" />

              <span class="ms-1.5">{{ $t('info.downloadBuiltInWebView2Version') }}</span>
            </label>
          </div>
        </div>

        <UButton
          color="secondary"
          block
          class="justify-center"
          @click="
            () => {
              if (!selectedUrl) {
                toast.add({ title: $t('toast.noAssetUrl'), color: 'error' })
                return
              }

              toast.add({
                title: $t('toast.downloadingUpdater'),
                color: 'info',
                duration: 60 * 1000
              })
              const loader = $loading.show()

              TriggerNativeUpdate(selectedUrl)
                .then(() => Quit())
                .catch(reason => toast.add({ title: reason, color: 'error' }))
                .finally(() => loader.hide())
            }
          "
        >
          {{ $t('info.update') }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>

<style scoped>
#release-notes * {
  all: revert;
}
</style>
