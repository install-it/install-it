<script setup lang="ts">
import UnsaveConfirmModal from '@/components/modals/UnsaveConfirmModal.vue'
import { store } from '@/wailsjs/go/models'
import * as appSettingStore from '@/wailsjs/go/store/AppSettingStore'
import { onBeforeMount, ref, toRaw, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'vue-toast-notification'

const { t, locale } = useI18n()

const $toast = useToast({ position: 'top-right' })

const questionModal = useTemplateRef('questionModal')

const tabKeys = ['softwareSetting', 'defaultInstallSetting', 'displaySetting'] as const

const currentTab = ref<(typeof tabKeys)[number]>(tabKeys[0])

const settings = ref<store.AppSetting>(new store.AppSetting())

let settingsOriginal: store.AppSetting

onBeforeMount(() => {
  appSettingStore
    .All()
    .then(s => {
      settings.value = s
      settingsOriginal = structuredClone(s)
    })
    .catch(() => {
      $toast.error(t('toast.readAppSettingFailed'))
    })
})

function handleTabClick(key: (typeof tabKeys)[number]) {
  if (JSON.stringify(settings.value) == JSON.stringify(settingsOriginal)) {
    currentTab.value = key
    return
  }

  questionModal.value?.show(answer => {
    if (answer == 'yes') {
      settings.value = structuredClone(settingsOriginal)
      currentTab.value = key
    }
    questionModal.value?.hide()
  })
}

function handleSubmit() {
  console.log(settings.value)
  appSettingStore.Update(settings.value).then(() => {
    locale.value = settings.value.language
    settingsOriginal = structuredClone(toRaw(settings.value))

    $toast.success(t('toast.updated'), { duration: 1500, position: 'top-right' })
  })
}
</script>

<template>
  <form class="flex flex-col h-full gap-y-3" @submit.prevent="handleSubmit()">
    <div class="flex items-center border-b-2">
      <button
        v-for="key in tabKeys"
        :key="key"
        type="button"
        class="px-4 py-2"
        :class="
          currentTab == key ? 'font-semibold border-b-2 border-b-kashmir-blue-500 -mb-[2px]' : ''
        "
        @click="() => handleTabClick(key)"
      >
        {{ $t(`setting.${key}`) }}
      </button>
    </div>

    <div v-show="currentTab == 'softwareSetting'" class="flex flex-col gap-y-3">
      <section>
        <p class="font-bold mb-2">
          {{ $t('setting.generalSetting') }}
        </p>

        <div class="flex flex-col gap-y-3">
          <div>
            <p class="block mb-2 text-gray-900">
              {{ $t('setting.autoCheckUpdate') }}
            </p>

            <label class="flex item-center w-full select-none cursor-pointer">
              <input
                type="checkbox"
                name="create_partition"
                v-model="settings.auto_check_update"
                class="me-1.5 text-blue-600 bg-gray-100 border-gray-300 rounded-sm focus:ring-blue-500"
              />
              {{ $t('common.enable') }}
            </label>
          </div>

          <div>
            <label class="block mb-2 text-gray-900">
              {{ $t('setting.successActionDelay') }}
            </label>

            <input
              type="number"
              name="success_action_delay"
              min="0"
              step="0"
              v-model="settings.success_action_delay"
              class="w-20 p-1.5 text-sm shadow-xs"
              required
            />
            &nbsp; {{ $t('setting.second') }}
          </div>
        </div>
      </section>

      <section>
        <p class="font-bold mb-2">{{ $t('setting.porter') }}</p>

        <div class="flex flex-col gap-y-3">
          <div>
            <label class="block mb-2 text-gray-900">{{ $t('setting.importUrl') }}</label>

            <input
              type="url"
              name="driver_download_url"
              v-model="settings.driver_download_url"
              class="w-full p-1.5 text-sm shadow-xs"
            />
          </div>
        </div>
      </section>
    </div>

    <div v-show="currentTab == 'defaultInstallSetting'" class="flex flex-col gap-y-3">
      <section>
        <p class="font-bold mb-2">
          {{ $t('setting.task') }}
        </p>

        <div class="flex flex-col gap-y-3">
          <div class="flex items-center">
            <label class="flex item-center w-full select-none cursor-pointer">
              <input
                type="checkbox"
                name="create_partition"
                v-model="settings.create_partition"
                class="me-1.5 text-blue-600 bg-gray-100 border-gray-300 rounded-sm focus:ring-blue-500"
              />
              {{ $t('installOption.createPartition') }}
            </label>
          </div>

          <div class="flex gap-3">
            <div class="flex items-center">
              <label class="flex item-center w-full select-none cursor-pointer">
                <input
                  type="checkbox"
                  name="set_password"
                  v-model="settings.set_password"
                  class="me-1.5 text-blue-600 bg-gray-100 border-gray-300 rounded-sm focus:ring-blue-500"
                />
                {{ $t('installOption.setPassword') }}
              </label>
            </div>

            <div class="flex shrink">
              <input
                type="text"
                name="password"
                v-model="settings.password"
                class="w-full min-w-20 p-1.5 text-sm shadow-xs"
                :disabled="!settings.set_password"
              />
            </div>
          </div>
        </div>
      </section>

      <section>
        <p class="font-bold mb-2">
          {{ $t('setting.installOption') }}
        </p>

        <div class="flex flex-col gap-y-3">
          <div class="flex items-center">
            <label class="flex item-center w-full select-none cursor-pointer">
              <input
                type="checkbox"
                name="parallel_install"
                v-model="settings.parallel_install"
                class="me-1.5 text-blue-600 bg-gray-100 border-gray-300 rounded-sm focus:ring-blue-500"
              />
              {{ $t('installOption.parallelInstall') }}
            </label>
          </div>

          <div>
            <label class="block mb-2 text-gray-900">
              {{ $t('installOption.successAction') }}
            </label>
            <select
              name="success_action"
              v-model="settings.success_action"
              class="w-full max-w-72 min-w-24 p-1.5 text-sm shadow-xs"
            >
              <option v-for="action in store.SuccessAction" :key="action" :value="action">
                {{ $t(`successAction.${action}`) }}
              </option>
            </select>
          </div>
        </div>
      </section>
    </div>

    <div v-show="currentTab == 'displaySetting'" class="flex flex-col gap-y-3">
      <section>
        <p class="font-bold mb-2">
          {{ $t('setting.language') }}
        </p>

        <div>
          <select
            name="language"
            v-model="settings.language"
            class="w-full max-w-72 min-w-24 p-1.5 text-sm shadow-xs"
          >
            <option value="en">English</option>
            <option value="zh_Hant_HK">繁體中文</option>
          </select>
        </div>
      </section>

      <section>
        <p class="font-bold mb-2">
          {{ $t('setting.hardwareInfo') }}
        </p>

        <div class="flex flex-col gap-y-3">
          <div class="flex items-center">
            <label class="flex item-center w-full select-none cursor-pointer">
              <input
                type="checkbox"
                name="filter_miniport_nic"
                v-model="settings.filter_miniport_nic"
                class="me-1.5 text-blue-600 bg-gray-100 border-gray-300 rounded-sm focus:ring-blue-500"
              />
              {{ $t('setting.filterMiniportNic') }}
            </label>
          </div>
        </div>

        <div class="flex flex-col gap-y-3">
          <div class="flex items-center">
            <label class="flex item-center w-full select-none cursor-pointer">
              <input
                type="checkbox"
                name="filter_microsoft_nic"
                v-model="settings.filter_microsoft_nic"
                class="me-1.5 text-blue-600 bg-gray-100 border-gray-300 rounded-sm focus:ring-blue-500"
              />
              {{ $t('setting.filterMicorsoftNic') }}
            </label>
          </div>
        </div>
      </section>
    </div>

    <div>
      <button
        type="submit"
        class="h-8 mt-3 px-3 text-white text-sm bg-half-baked-600 hover:bg-half-baked-500 rounded-sm"
      >
        {{ $t('common.save') }}
      </button>
    </div>
  </form>

  <UnsaveConfirmModal ref="questionModal"></UnsaveConfirmModal>
</template>
