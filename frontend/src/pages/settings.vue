<script setup lang="ts">
import UnsaveConfirmModal from '@/components/UnsaveConfirmModal.vue'
import { useEditor } from '@/composables/useEditor'
import { useAppSettingStore } from '@/store'
import { storage } from '@/wailsjs/go/models'
import * as appSettingStorage from '@/wailsjs/go/storage/AppSettingStorage'
import { ref, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave } from 'vue-router'

const { t, locale } = useI18n()

const toast = useToast()

const questionModal = useTemplateRef('questionModal')

const tabs = ref({ softwareSetting: true, defaultInstallSetting: false, displaySetting: false })

const settingStore = useAppSettingStore()
const { data: settings, modified, reset } = useEditor({ source: () => settingStore.settings })

onBeforeRouteLeave(async (to, from, next) => {
  next(await handleMoveAway())
})

async function handleMoveAway() {
  return new Promise<boolean>(resolve => {
    if (modified.value) {
      questionModal.value?.show(answer => {
        if (answer == 'yes') {
          reset()
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
      return reset()
    })
    .then(() => {
      locale.value = settings.value.language
      toast.add({ title: t('toast.saved'), color: 'success', duration: 1500 })
    })
    .catch(() => {
      toast.add({ title: t('toast.failedToSave'), color: 'error', duration: 1500 })
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
              <UCheckbox
                v-model="settings.auto_check_update"
                name="auto_check_update"
                color="primary"
                class="me-1.5"
              />
              {{ $t('common.enable') }}
            </label>
          </div>

          <div>
            <label class="mb-2 block text-gray-900">
              {{ $t('setting.successActionDelay') }}
            </label>

            <UInput
              v-model="settings.success_action_delay"
              type="number"
              name="success_action_delay"
              min="0"
              step="0"
              color="primary"
              class="w-20"
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

            <UInput
              v-model="settings.driver_download_url"
              type="url"
              name="driver_download_url"
              color="primary"
              class="w-full"
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
              <UCheckbox
                v-model="settings.create_partition"
                name="create_partition"
                color="primary"
                class="me-1.5"
              />
              {{ $t('installSetting.createPartition') }}
            </label>
          </div>

          <div class="flex gap-3">
            <div class="flex">
              <label class="flex w-full cursor-pointer items-center select-none">
                <UCheckbox
                  v-model="settings.set_password"
                  name="set_password"
                  color="primary"
                  class="me-1.5"
                />
                {{ $t('installSetting.setPassword') }}
              </label>
            </div>

            <div class="flex shrink">
              <UInput
                v-model="settings.password"
                type="password"
                name="password"
                color="primary"
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
              <UCheckbox
                v-model="settings.parallel_install"
                name="parallel_install"
                color="primary"
                class="me-1.5"
              />
              {{ $t('installSetting.parallelInstall') }}
            </label>
          </div>

          <div>
            <label class="mb-2 block text-gray-900">
              {{ $t('installSetting.successAction') }}
            </label>

            <USelect
              v-model="settings.success_action"
              name="success_action"
              color="primary"
              :items="
                Object.values(storage.SuccessAction).map((action: string) => ({
                  label: $t(`successAction.${action}`),
                  value: action
                }))
              "
            />
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
          <USelect
            v-model="settings.language"
            name="language"
            color="primary"
            :items="[
              { label: 'English', value: 'en' },
              { label: '繁體中文', value: 'zh_Hant_HK' }
            ]"
          />
        </div>
      </section>

      <section>
        <p class="mb-2 font-bold">
          {{ $t('setting.hardwareInfo') }}
        </p>

        <div class="flex flex-col gap-y-3">
          <div class="flex">
            <label class="flex w-full cursor-pointer items-center select-none">
              <UCheckbox
                v-model="settings.filter_miniport_nic"
                name="filter_miniport_nic"
                color="primary"
                class="me-1.5"
              />
              {{ $t('setting.filterMiniportNic') }}
            </label>
          </div>
        </div>

        <div class="flex flex-col gap-y-3">
          <div class="flex">
            <label class="flex w-full cursor-pointer items-center select-none">
              <UCheckbox
                v-model="settings.filter_microsoft_nic"
                name="filter_microsoft_nic"
                color="primary"
                class="me-1.5"
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
              <UCheckbox
                v-model="settings.hide_not_found"
                name="hide_not_found"
                color="primary"
                class="me-1.5"
              />
              {{ $t('setting.hideNotFound') }}
            </label>
          </div>
        </div>
      </section>
    </div>

    <div class="mt-6">
      <UButton type="submit" color="secondary">
        {{ $t('common.save') }}
      </UButton>
    </div>
  </form>

  <UnsaveConfirmModal ref="questionModal"></UnsaveConfirmModal>
</template>
