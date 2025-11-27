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
  <form class="flex h-full flex-col gap-y-3" @submit.prevent="handleSubmit">
    <div class="flex items-center border-b-2">
      <button
        v-for="key in Object.keys(tabs)"
        :key="key"
        type="button"
        class="px-4 py-2"
        :class="
          tabs[key as keyof typeof tabs]
            ? '-mb-0.5 border-b-2 border-b-kashmir-blue-500 font-semibold'
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
        <p class="mb-2 font-bold">
          {{ $t('setting.generalSetting') }}
        </p>

        <div class="flex flex-col gap-y-3">
          <div>
            <p class="mb-2 block text-gray-900">
              {{ $t('setting.autoCheckUpdate') }}
            </p>

            <label class="flex w-full cursor-pointer items-center select-none">
              <input
                v-model="settings.auto_check_update"
                type="checkbox"
                name="auto_check_update"
                class="checkbox me-1.5 checkbox-primary"
              />
              {{ $t('common.enable') }}
            </label>
          </div>

          <div>
            <label class="mb-2 block text-gray-900">
              {{ $t('setting.successActionDelay') }}
            </label>

            <input
              v-model="settings.success_action_delay"
              type="number"
              name="success_action_delay"
              min="0"
              step="0"
              class="input w-20 shadow-xs input-accent"
              required
            />
            &nbsp; {{ $t('setting.second') }}
          </div>
        </div>
      </section>

      <section>
        <p class="mb-2 font-bold">{{ $t('setting.porter') }}</p>

        <div class="flex flex-col gap-y-3">
          <div>
            <label class="mb-2 block text-gray-900">{{ $t('setting.importUrl') }}</label>

            <input
              v-model="settings.driver_download_url"
              type="url"
              name="driver_download_url"
              class="input w-full shadow-xs input-accent"
            />
          </div>
        </div>
      </section>
    </div>

    <div v-show="tabs.defaultInstallSetting" class="flex flex-col gap-y-3">
      <section>
        <p class="mb-2 font-bold">
          {{ $t('setting.task') }}
        </p>

        <div class="flex flex-col gap-y-3">
          <div class="flex">
            <label class="flex w-full cursor-pointer items-center select-none">
              <input
                v-model="settings.create_partition"
                type="checkbox"
                name="create_partition"
                class="checkbox me-1.5 checkbox-primary"
              />
              {{ $t('installSetting.createPartition') }}
            </label>
          </div>

          <div class="flex gap-3">
            <div class="flex">
              <label class="flex w-full cursor-pointer items-center select-none">
                <input
                  v-model="settings.set_password"
                  type="checkbox"
                  name="set_password"
                  class="checkbox me-1.5 checkbox-primary"
                />
                {{ $t('installSetting.setPassword') }}
              </label>
            </div>

            <div class="flex shrink">
              <input
                v-model="settings.password"
                type="text"
                name="password"
                class="input input-accent"
                :disabled="!settings.set_password"
              />
            </div>
          </div>
        </div>
      </section>

      <section>
        <p class="mb-2 font-bold">
          {{ $t('setting.installSetting') }}
        </p>

        <div class="flex flex-col gap-y-3">
          <div class="flex">
            <label class="flex w-full cursor-pointer items-center select-none">
              <input
                v-model="settings.parallel_install"
                type="checkbox"
                name="parallel_install"
                class="checkbox me-1.5 checkbox-primary"
              />
              {{ $t('installSetting.parallelInstall') }}
            </label>
          </div>

          <div>
            <label class="mb-2 block text-gray-900">
              {{ $t('installSetting.successAction') }}
            </label>

            <select
              v-model="settings.success_action"
              name="success_action"
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
        <p class="mb-2 font-bold">
          {{ $t('setting.language') }}
        </p>

        <div>
          <select v-model="settings.language" name="language" class="select select-accent">
            <option value="en">English</option>

            <option value="zh_Hant_HK">繁體中文</option>
          </select>
        </div>
      </section>

      <section>
        <p class="mb-2 font-bold">
          {{ $t('setting.hardwareInfo') }}
        </p>

        <div class="flex flex-col gap-y-3">
          <div class="flex">
            <label class="flex w-full cursor-pointer items-center select-none">
              <input
                v-model="settings.filter_miniport_nic"
                type="checkbox"
                name="filter_miniport_nic"
                class="checkbox me-1.5 checkbox-primary"
              />
              {{ $t('setting.filterMiniportNic') }}
            </label>
          </div>
        </div>

        <div class="flex flex-col gap-y-3">
          <div class="flex">
            <label class="flex w-full cursor-pointer items-center select-none">
              <input
                v-model="settings.filter_microsoft_nic"
                type="checkbox"
                name="filter_microsoft_nic"
                class="checkbox me-1.5 checkbox-primary"
              />
              {{ $t('setting.filterMicorsoftNic') }}
            </label>
          </div>
        </div>
      </section>

      <section>
        <p class="mb-2 font-bold">{{ $t('setting.installOption') }}</p>

        <div class="flex flex-col gap-y-3">
          <div class="flex">
            <label class="flex w-full cursor-pointer items-center select-none">
              <input
                v-model="settings.hide_not_found"
                type="checkbox"
                name="hide_not_found"
                class="checkbox me-1.5 checkbox-primary"
              />
              {{ $t('setting.hideNotFound') }}
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
