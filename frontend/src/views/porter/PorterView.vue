<script setup lang="ts">
import { onBeforeMount, ref, useTemplateRef } from 'vue'

import { Cwd, SelectFile, SelectFolder } from '@/wailsjs/go/main/App'
import * as appStorage from '@/wailsjs/go/storage/AppSettingStorage'
import DownloadProgressModal from './components/ProgressModal.vue'

const progressModal = useTemplateRef('progressModal')

const exportDirectory = ref('')

const importInput = ref<{
  from: 'file' | 'url'
  filePath: string
  url: string
}>({
  from: 'file',
  filePath: '',
  url: ''
})

onBeforeMount(() => {
  appStorage.All().then(s => (importInput.value.url = s.driver_download_url))

  Cwd().then(cwd => {
    exportDirectory.value = cwd
  })
})
</script>

<template>
  <div class="flex flex-col h-full gap-y-6">
    <div>
      <h1 class="text-xl font-bold">{{ $t('porter.title') }}</h1>
      <p class="text-gray-400">{{ $t('porter.titleHint') }}</p>

      <hr class="mt-2 -mb-3" />
    </div>

    <div class="flex flex-col gap-y-3">
      <h2 class="mb-1 text-lg font-medium">{{ $t('porter.exportToFile') }}</h2>

      <div class="flex gap-x-6">
        <label class="w-24 content-center text-gray-900">
          {{ $t('porter.exportDestination') }}
        </label>

        <div class="flex gap-x-2 w-full">
          <input
            type="url"
            name="export_directory"
            v-model="exportDirectory"
            class="grow input input-accent"
          />

          <button
            type="button"
            class="btn btn-primary"
            @click="
              () => {
                SelectFolder(false).then(path => {
                  if (path != '') {
                    exportDirectory = path
                  }
                })
              }
            "
          >
            {{ $t('common.select') }}
          </button>
        </div>
      </div>

      <div class="flex justify-end">
        <button
          type="button"
          class="btn btn-secondary mt-3 w-28"
          @click="
            () => {
              if (!exportDirectory) {
                $toast.warning($t('toast.enterExportPath'), { position: 'bottom-right' })
              } else {
                progressModal?.export(exportDirectory)
              }
            }
          "
        >
          {{ $t('porter.export') }}
        </button>
      </div>
    </div>

    <div class="flex flex-col gap-y-3">
      <div class="flex gap-x-4">
        <h2 class="mb-1 text-lg font-medium">
          {{ $t('porter.import') }}
        </h2>

        <div class="relative inline-flex p-0.5 border rounded-3xl">
          <button
            class="z-10 px-3 text-center text-xs rounded-3xl select-none"
            @click="importInput.from = 'file'"
          >
            {{ $t('porter.importFromFile') }}
          </button>

          <button
            class="z-10 px-3 text-center text-xs rounded-3xl select-none"
            @click="importInput.from = 'url'"
          >
            {{ $t('porter.importFromNetwork') }}
          </button>

          <span
            class="absolute top-1 bg-gray-300 rounded-3xl transition duration-200"
            :class="{ 'translate-x-full': importInput.from == 'url' }"
            style="width: calc(50% - 2px); height: calc(100% - 8px)"
          ></span>
        </div>
      </div>

      <form
        @submit.prevent="
          progressModal?.import(
            importInput.from,
            importInput.from == 'file' ? importInput.filePath : importInput.url
          )
        "
      >
        <!-- from file -->
        <div v-if="importInput.from == 'file'" class="flex gap-x-6">
          <label class="w-24 content-center text-gray-900">
            {{ $t('porter.file') }}
          </label>

          <div class="flex gap-x-2 w-full">
            <input
              type="text"
              name="driver_download_url"
              placeholder="install-it.zip"
              v-model="importInput.filePath"
              class="grow input input-accent pointer-events-none"
              :required="importInput.from == 'file'"
            />

            <button
              type="button"
              class="btn btn-primary"
              @click="
                () => {
                  SelectFile(false).then(path => {
                    if (path != '') {
                      importInput.filePath = path
                    }
                  })
                }
              "
            >
              {{ $t('common.select') }}
            </button>
          </div>
        </div>

        <!-- from url -->
        <div v-else class="flex gap-x-6">
          <label class="w-24 content-center text-gray-900">
            {{ $t('porter.url') }}
          </label>

          <div class="flex gap-x-2 w-full">
            <input
              type="url"
              placeholder="https://..."
              v-model="importInput.url"
              class="grow input input-accent"
              :required="importInput.from == 'url'"
            />
          </div>
        </div>

        <div class="flex justify-end">
          <button type="submit" class="btn btn-secondary mt-3 w-28">
            {{ $t('porter.import') }}
          </button>
        </div>
      </form>
    </div>
  </div>

  <DownloadProgressModal ref="progressModal"></DownloadProgressModal>
</template>
