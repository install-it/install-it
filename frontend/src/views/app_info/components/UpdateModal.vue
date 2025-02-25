<script setup lang="ts">
import CrossIcon from '@/components/icons/CrossIcon.vue'
import ModalFrame from '@/components/modals/ModalFrame.vue'
import { latestRelease } from '@/utils'
import { Update } from '@/wailsjs/go/main/App'
import { Quit } from '@/wailsjs/runtime/runtime'
import { ref, useTemplateRef } from 'vue'
import { useLoading } from 'vue-loading-overlay'

const frame = useTemplateRef('frame')

const props = defineProps<{
  app: { version: string; binaryType: string; builtinWebview: boolean }
}>()

defineExpose({
  show: (releaseInfo_: typeof releaseInfo.value) => {
    releaseInfo.value = releaseInfo_
    frame.value?.show()
  },
  hide: frame.value?.hide || (() => {})
})

const $loading = useLoading({ lockScroll: true })

const releaseInfo = ref<Awaited<ReturnType<typeof latestRelease>>>()

const webviewVersion = ref(!props.app.builtinWebview)
</script>

<template>
  <ModalFrame :on-demand="true" :immediate="false" ref="frame">
    <div class="w-4/5 max-w-4xl">
      <!-- Modal content -->
      <div class="bg-white rounded-lg shadow-sm">
        <!-- Modal header -->
        <div class="flex items-center justify-between h-12 px-4 border-b rounded-t">
          <h3 class="font-semibold">
            {{ $t('info.updateInfoTitle') }}
          </h3>

          <button
            type="button"
            class="p-3 text-sm text-gray-400 hover:text-gray-900 bg-transparent hover:bg-gray-100 rounded-lg"
            @click="
              () => {
                $refs.frame?.hide()
              }
            "
          >
            <CrossIcon></CrossIcon>
          </button>
        </div>

        <!-- Modal body -->
        <div class="flex flex-col gap-y-3 min-h-40 max-h-96 overflow-y-auto py-2 px-4">
          <div class="flex flex-col gap-y-2 grow">
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

            <div class="flex flex-col grow">
              <h1 class="min-w-32 mb-1 font-medium">
                {{ $t('info.updateInfo') }}
              </h1>

              <div
                v-html="releaseInfo?.releaseNotes || `<i>${$t('info.noUpdateInfo')}</i>`"
                id="release-notes"
                class="px-1 border rounded-lg"
              ></div>
            </div>

            <hr />

            <div class="flex flex-col">
              <h1 class="font-medium">
                {{ $t('info.updateOption') }}
              </h1>

              <label class="flex item-center w-full select-none cursor-pointer">
                <input
                  type="checkbox"
                  name="create_partition"
                  v-model="webviewVersion"
                  class="me-1.5 text-blue-600 bg-gray-100 border-gray-300 rounded-sm focus:ring-blue-500"
                />
                {{ $t('info.downloadBuiltInWebView2Version') }}
              </label>
            </div>
          </div>

          <button
            class="w-full py-1 text-white bg-half-baked-600 hover:bg-half-baked-500 rounded-sm"
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
