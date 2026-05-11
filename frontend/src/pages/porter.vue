<script setup lang="ts">
import { Cwd, SelectFile, SelectFolder } from '@/wailsjs/go/main/App'
import * as appStorage from '@/wailsjs/go/storage/AppSettingStorage'
import { onBeforeMount, ref, useTemplateRef } from 'vue'

const toast = useToast()

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
  <div class="flex h-full flex-col gap-y-6 p-2">
    <div>
      <h1 class="text-lg font-bold">{{ $t('porter.title') }}</h1>

      <p class="text-gray-400">{{ $t('porter.titleHint') }}</p>

      <hr class="mt-2 -mb-3" />
    </div>

    <div class="flex flex-col gap-y-3">
      <h2 class="mb-1 font-medium">{{ $t('porter.exportToFile') }}</h2>

      <div class="flex gap-x-6">
        <label class="w-24 content-center text-gray-900">
          {{ $t('porter.exportDestination') }}
        </label>

        <div class="flex w-full gap-x-2">
          <UInput
            v-model="exportDirectory"
            type="url"
            name="export_directory"
            color="primary"
            class="grow"
          />

          <UButton
            type="button"
            color="primary"
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
          </UButton>
        </div>
      </div>

      <div class="flex justify-end">
        <UButton
          type="button"
          color="secondary"
          class="mt-3 w-28 justify-center"
          @click="
            () => {
              if (!exportDirectory) {
                toast.add({ title: $t('toast.enterExportPath'), color: 'warning' })
              } else {
                progressModal?.export(exportDirectory)
              }
            }
          "
        >
          {{ $t('porter.export') }}
        </UButton>
      </div>
    </div>

    <div class="flex flex-col gap-y-3">
      <div class="flex gap-x-4">
        <h2 class="mb-1 font-medium">
          {{ $t('porter.import') }}
        </h2>

        <div class="relative inline-flex rounded-3xl border p-0.5">
          <button
            class="z-10 rounded-3xl px-3 text-center text-xs select-none"
            @click="importInput.from = 'file'"
          >
            {{ $t('porter.importFromFile') }}
          </button>

          <button
            class="z-10 rounded-3xl px-3 text-center text-xs select-none"
            @click="importInput.from = 'url'"
          >
            {{ $t('porter.importFromNetwork') }}
          </button>

          <span
            class="absolute top-1 rounded-3xl bg-gray-300 transition duration-200"
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

          <div class="flex w-full gap-x-2">
            <UInput
              v-model="importInput.filePath"
              type="text"
              name="driver_download_url"
              placeholder="install-it.zip"
              color="primary"
              class="pointer-events-none grow"
              :required="importInput.from == 'file'"
            />

            <UButton
              type="button"
              color="primary"
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
            </UButton>
          </div>
        </div>

        <!-- from url -->
        <div v-else class="flex gap-x-6">
          <label class="w-24 content-center text-gray-900">
            {{ $t('porter.url') }}
          </label>

          <div class="flex w-full gap-x-2">
            <UInput
              v-model="importInput.url"
              type="url"
              placeholder="https://..."
              color="primary"
              class="grow"
              :required="importInput.from == 'url'"
            />
          </div>
        </div>

        <div class="flex justify-end">
          <UButton type="submit" color="secondary" class="mt-3 w-28 justify-center">
            {{ $t('porter.import') }}
          </UButton>
        </div>
      </form>
    </div>
  </div>

  <ProgressModal ref="progressModal"></ProgressModal>
</template>
