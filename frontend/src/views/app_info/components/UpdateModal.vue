<!-- eslint-disable vue/no-v-html -->

<script setup lang="ts">
import ModalFrame from '@/components/modals/ModalFrame.vue'
import { latestRelease } from '@/utils'
import { Update } from '@/wailsjs/go/main/App'
import { Quit } from '@/wailsjs/runtime/runtime'
import { ref, useTemplateRef } from 'vue'
import { useLoading } from 'vue-loading-overlay'

const frame = useTemplateRef('frame')

defineProps<{ app: { version: string; binaryType: string } }>()

defineExpose({
  show: (releaseInfo_: typeof releaseInfo.value, isWebview: boolean) => {
    releaseInfo.value = releaseInfo_
    webviewVersion.value = isWebview
    frame.value?.show()
  },
  hide: frame.value?.hide || (() => {})
})

const $loading = useLoading({ lockScroll: true })

const releaseInfo = ref<Awaited<ReturnType<typeof latestRelease>>>()

const webviewVersion = ref(false)
</script>

<template>
  <ModalFrame ref="frame" :on-demand="true" :immediate="false">
    <div class="w-4/5 max-w-4xl">
      <!-- Modal content -->
      <div class="rounded-lg bg-white shadow-sm">
        <!-- Modal header -->
        <div class="flex h-12 items-center justify-between rounded-t border-b px-4">
          <h3 class="font-semibold">
            {{ $t('info.updateInfoTitle') }}
          </h3>

          <button
            type="button"
            class="rounded-lg bg-transparent p-3 text-sm text-gray-400 hover:bg-gray-100 hover:text-gray-900"
            @click="
              () => {
                frame?.hide()
              }
            "
          >
            <font-awesome-icon icon="fa-solid fa-xmark" />
          </button>
        </div>

        <!-- Modal body -->
        <div class="flex max-h-96 min-h-40 flex-col gap-y-3 overflow-y-auto px-4 py-2">
          <div class="flex grow flex-col gap-y-2">
            <div class="flex">
              <h1 class="min-w-34 font-medium">
                {{ $t('info.currentVersion') }}
              </h1>

              <p>{{ $props.app.version }}</p>
            </div>

            <div class="flex">
              <h1 class="min-w-34 font-medium">
                {{ $t('info.latestVersion') }}
              </h1>

              <p>
                {{ `${releaseInfo?.version} (${releaseInfo?.releaseAt.toLocaleDateString()})` }}
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
                v-html="releaseInfo?.releaseNotes || `<i>${$t('info.noUpdateInfo')}</i>`"
              ></div>
            </div>

            <hr />

            <div class="flex flex-col">
              <h1 class="font-medium">
                {{ $t('info.updateOption') }}
              </h1>

              <label class="flex w-full cursor-pointer items-center select-none">
                <input
                  v-model="webviewVersion"
                  type="checkbox"
                  name="create_partition"
                  class="checkbox me-1.5 checkbox-primary"
                />
                {{ $t('info.downloadBuiltInWebView2Version') }}
              </label>
            </div>
          </div>

          <button
            class="btn w-full btn-secondary"
            @click="
              () => {
                if (!releaseInfo) {
                  return
                }

                $toast.info($t('toast.downloadingUpdater'), { duration: 60 * 1000 })
                const loader = $loading.show()

                Update($props.app.version, releaseInfo.version, webviewVersion)
                  .then(() => Quit())
                  .catch(reason => $toast.error(reason))
                  .finally(() => loader.hide())
              }
            "
          >
            {{ $t('info.update') }}
          </button>
        </div>
      </div>
    </div>
  </ModalFrame>
</template>

<style scoped>
label:has(+ input:required, + select:required):after,
label:has(+ div > input:required):after {
  content: ' *';
  color: red;
}

#release-notes * {
  all: revert;
}
</style>
