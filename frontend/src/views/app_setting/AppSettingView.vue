<script setup lang="ts">
import UnsaveConfirmModal from '@/components/modals/UnsaveConfirmModal.vue'
import { useAppSettingStore } from '@/store'
import { storage } from '@/wailsjs/go/models'
import * as appSettingStorage from '@/wailsjs/go/storage/AppSettingStorage'
import { ref, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave } from 'vue-router'
import { useToast } from 'vue-toast-notification'

const { t, locale } = useI18n()

const $toast = useToast()

const questionModal = useTemplateRef('questionModal')

const tabs = ref({ softwareSetting: true, defaultInstallSetting: false, displaySetting: false })

const settingEditor = useAppSettingStore().editor()

const settings = settingEditor.settings // alias

onBeforeRouteLeave(async (to, from, next) => {
  next(await handleMoveAway())
})

async function handleMoveAway() {
  return new Promise<boolean>(resolve => {
    if (settingEditor.modified.value) {
      questionModal.value?.show(answer => {
        if (answer == 'yes') {
          settingEditor.reset()
        }
        resolve(answer == 'yes')
      })
    } else {
      resolve(true)
    }
  })
}

function handleSubmit() {
  appSettingStorage
    .Update(settings.value)
    .then(newAppSettings => {
      useAppSettingStore().settings = newAppSettings
      settingEditor.reset()

      locale.value = settings.value.language
      $toast.success(t('toast.saved'), { duration: 1500, position: 'top-right' })
    })
    .catch(() => {
      $toast.error(t('toast.failedToSave'), { duration: 1500, position: 'top-right' })
    })
}
</script>

<template>
  <form class="flex flex-col h-full gap-y-3" @submit.prevent="handleSubmit">
    <div class="flex items-center border-b-2">
      <button
        v-for="key in Object.keys(tabs)"
        :key="key"
        type="button"
        class="px-4 py-2"
        :class="
          tabs[key as keyof typeof tabs]
            ? 'font-semibold border-b-2 border-b-kashmir-blue-500 -mb-[2px]'
            : ''
        "
        @click="
          handleMoveAway().then(leave => {
            if (leave) {
              Object.keys(tabs).forEach(k => (tabs[k as keyof typeof tabs] = k == key))
            }
          })
        "
      >
        {{ $t(`setting.${key}`) }}
      </button>
    </div>

    <div v-show="tabs.softwareSetting" class="flex flex-col gap-y-3">
      <section>
        <p class="font-bold mb-2">
          {{ $t('setting.generalSetting') }}
        </p>

        <div class="flex flex-col gap-y-3">
          <div>
            <p class="block mb-2 text-gray-900">
              {{ $t('setting.autoCheckUpdate') }}
            </p>

            <label class="flex items-center w-full select-none cursor-pointer">
              <input
                type="checkbox"
                name="auto_check_update"
                v-model="settings.auto_check_update"
                class="checkbox checkbox-primary me-1.5"
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
              class="w-20 input input-accent shadow-xs"
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
              class="w-full input input-accent shadow-xs"
            />
          </div>
        </div>
      </section>
    </div>

    <div v-show="tabs.defaultInstallSetting" class="flex flex-col gap-y-3">
      <section>
        <p class="font-bold mb-2">
          {{ $t('setting.task') }}
        </p>

        <div class="flex flex-col gap-y-3">
          <div class="flex">
            <label class="flex items-center w-full select-none cursor-pointer">
              <input
                type="checkbox"
                name="create_partition"
                v-model="settings.create_partition"
                class="checkbox checkbox-primary me-1.5"
              />
              {{ $t('installOption.createPartition') }}
            </label>
          </div>

          <div class="flex gap-3">
            <div class="flex">
              <label class="flex items-center w-full select-none cursor-pointer">
                <input
                  type="checkbox"
                  name="set_password"
                  v-model="settings.set_password"
                  class="checkbox checkbox-primary me-1.5"
                />
                {{ $t('installOption.setPassword') }}
              </label>
            </div>

            <div class="flex shrink">
              <input
                type="text"
                name="password"
                v-model="settings.password"
                class="input input-accent"
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
          <div class="flex">
            <label class="flex items-center w-full select-none cursor-pointer">
              <input
                type="checkbox"
                name="parallel_install"
                v-model="settings.parallel_install"
                class="checkbox checkbox-primary me-1.5"
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
              class="select select-accent"
            >
              <option v-for="action in storage.SuccessAction" :key="action" :value="action">
                {{ $t(`successAction.${action}`) }}
              </option>
            </select>
          </div>
        </div>
      </section>
    </div>

    <div v-show="tabs.displaySetting" class="flex flex-col gap-y-3">
      <section>
        <p class="font-bold mb-2">
          {{ $t('setting.language') }}
        </p>

        <div>
          <select name="language" v-model="settings.language" class="select select-accent">
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
          <div class="flex">
            <label class="flex items-center w-full select-none cursor-pointer">
              <input
                type="checkbox"
                name="filter_miniport_nic"
                v-model="settings.filter_miniport_nic"
                class="checkbox checkbox-primary me-1.5"
              />
              {{ $t('setting.filterMiniportNic') }}
            </label>
          </div>
        </div>

        <div class="flex flex-col gap-y-3">
          <div class="flex">
            <label class="flex items-center w-full select-none cursor-pointer">
              <input
                type="checkbox"
                name="filter_microsoft_nic"
                v-model="settings.filter_microsoft_nic"
                class="checkbox checkbox-primary me-1.5"
              />
              {{ $t('setting.filterMicorsoftNic') }}
            </label>
          </div>
        </div>
      </section>
    </div>

    <div class="mt-6">
      <button type="submit" class="btn btn-secondary">
        {{ $t('common.save') }}
      </button>
    </div>
  </form>

  <UnsaveConfirmModal ref="questionModal"></UnsaveConfirmModal>
</template>
