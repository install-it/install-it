<script setup lang="ts">
import { storage } from '@/wailsjs/go/models'
import * as appSettingStorage from '@/wailsjs/go/storage/AppSettingStorage'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const toast = useToast()

const tabs = ref({ general: true, installDefaults: false, filters: false })

const settingStore = useAppSettingStore()
const {
  data: settings,
  modified,
  reset
} = useEditor({ source: () => settingStore.settings, warnOnUnsavedLeave: true })

function askLeave() {
  const formStore = useUnsavedFormStore()
  if (!modified.value) {
    return Promise.resolve(true)
  }

  formStore.show = true
  return new Promise<boolean>(resolve => {
    formStore.setAnswerHandler(resolve)
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
      toast.add({ title: t('toast.saved'), color: 'success' })
    })
    .catch(() => {
      toast.add({ title: t('toast.failedToSave'), color: 'error' })
    })
}
</script>

<template>
  <form
    class="flex h-full flex-col bg-transparent text-gray-900 dark:text-gray-100"
    @submit.prevent="handleSubmit"
  >
    <header
      class="shrink-0 border-b border-gray-200 bg-gray-50/50 px-6 py-2.5 dark:border-gray-800 dark:bg-gray-900/20"
    >
      <nav class="flex gap-x-1.5">
        <UButton
          v-for="key in Object.keys(tabs)"
          :key="key"
          type="button"
          :color="tabs[key as keyof typeof tabs] ? 'secondary' : 'neutral'"
          :variant="tabs[key as keyof typeof tabs] ? 'solid' : 'ghost'"
          class="font-semibold"
          @click="
            askLeave().then(leave => {
              if (leave) {
                reset()
                Object.keys(tabs).forEach(k => (tabs[k as keyof typeof tabs] = k === key))
              }
            })
          "
        >
          {{ $t(`setting.${key}`) }}
        </UButton>
      </nav>
    </header>

    <main class="flex-1 overflow-y-auto px-6 py-6">
      <div class="max-w-3xl">
        <div v-show="tabs.general" class="flex flex-col gap-y-6">
          <section>
            <h3 class="mb-3 text-lg font-bold">{{ $t('setting.generalSetting') }}</h3>

            <div class="flex flex-col gap-y-1">
              <label>{{ $t('setting.language') }}</label>
              <USelect
                v-model="settings.language"
                name="language"
                color="primary"
                variant="outline"
                class="w-56"
                :items="[
                  { label: 'English', value: 'en' },
                  { label: '繁體中文', value: 'zh_Hant_HK' }
                ]"
              />
            </div>
          </section>

          <hr class="border-gray-100 dark:border-gray-800" />

          <section>
            <h3 class="mb-3 text-lg font-bold">{{ $t('info.updateOption') }}</h3>

            <div class="flex flex-col gap-y-3">
              <label class="flex w-fit cursor-pointer items-center select-none">
                <UCheckbox
                  v-model="settings.auto_check_update"
                  name="auto_check_update"
                  color="primary"
                  class="me-2"
                />
                <span>{{ $t('setting.autoCheckUpdate') }}</span>
              </label>

              <label class="flex w-fit cursor-pointer items-center select-none">
                <UCheckbox
                  v-model="settings.allow_pre_release"
                  name="allow_pre_release"
                  color="primary"
                  class="me-2"
                />
                <span>{{ $t('setting.preferPreRelease') }}</span>
              </label>
            </div>
          </section>

          <hr class="border-gray-100 dark:border-gray-800" />

          <section>
            <h3 class="mb-3 text-lg font-bold">{{ $t('setting.porter') }}</h3>

            <div class="flex flex-col gap-y-1">
              <label>{{ $t('setting.importUrl') }}</label>
              <UTextarea
                v-model="settings.driver_download_url"
                name="driver_download_url"
                color="primary"
                placeholder="https://"
                class="w-full max-w-lg"
                :rows="2"
              />
            </div>
          </section>
        </div>

        <div v-show="tabs.installDefaults" class="flex flex-col gap-y-6">
          <section>
            <h3 class="mb-3 text-lg font-bold">{{ $t('setting.task') }}</h3>

            <div class="flex flex-col gap-y-3">
              <label class="flex w-fit cursor-pointer items-center select-none">
                <UCheckbox
                  v-model="settings.create_partition"
                  name="create_partition"
                  color="primary"
                  class="me-2"
                />
                <span>{{ $t('installSetting.createPartition') }}</span>
              </label>

              <div class="flex flex-col gap-y-2">
                <label class="flex w-fit cursor-pointer items-center select-none">
                  <UCheckbox
                    v-model="settings.set_password"
                    name="set_password"
                    color="primary"
                    class="me-2"
                  />
                  <span>{{ $t('installSetting.setPassword') }}</span>
                </label>

                <div
                  class="ml-6 transition-opacity duration-200"
                  :class="{ 'opacity-50': !settings.set_password }"
                >
                  <UInput
                    v-model="settings.password"
                    type="password"
                    name="password"
                    placeholder="••••••••"
                    color="primary"
                    class="w-56"
                    :disabled="!settings.set_password"
                  />
                </div>
              </div>
            </div>
          </section>

          <hr class="border-gray-100 dark:border-gray-800" />

          <section>
            <h3 class="mb-3 text-lg font-bold">{{ $t('setting.installSetting') }}</h3>

            <div class="flex flex-col gap-y-4">
              <label class="flex w-fit cursor-pointer items-center select-none">
                <UCheckbox
                  v-model="settings.parallel_install"
                  name="parallel_install"
                  color="primary"
                  class="me-2"
                />
                <span>{{ $t('installSetting.parallelInstall') }}</span>
              </label>

              <div class="flex flex-wrap items-start gap-4">
                <div class="flex flex-col gap-y-1">
                  <label>{{ $t('installSetting.successAction') }}</label>
                  <USelect
                    v-model="settings.success_action"
                    name="success_action"
                    color="primary"
                    class="w-56"
                    :items="
                      Object.values(storage.SuccessAction).map((action: string) => ({
                        label: $t(`successAction.${action}`),
                        value: action
                      }))
                    "
                  />
                </div>

                <div
                  class="flex flex-col gap-y-1 transition-opacity duration-200"
                  :class="{
                    'pointer-events-none opacity-40': settings.success_action === 'nothing'
                  }"
                >
                  <label>{{ $t('setting.successActionDelay') }}</label>
                  <div class="flex items-center gap-x-2">
                    <UInput
                      v-model="settings.success_action_delay"
                      type="number"
                      name="success_action_delay"
                      min="0"
                      color="primary"
                      class="w-24"
                      :disabled="settings.success_action === 'nothing'"
                      required
                    />
                    <span class="text-sm text-gray-500">{{ $t('setting.second') }}</span>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>

        <div v-show="tabs.filters" class="flex flex-col gap-y-6">
          <section>
            <h3 class="mb-3 text-lg font-bold">{{ $t('setting.hardwareInfo') }}</h3>

            <div class="flex flex-col gap-y-3">
              <label class="flex w-fit cursor-pointer items-center select-none">
                <UCheckbox
                  v-model="settings.filter_miniport_nic"
                  name="filter_miniport_nic"
                  color="primary"
                  class="me-2"
                />
                <span>{{ $t('setting.filterMiniportNic') }}</span>
              </label>

              <label class="flex w-fit cursor-pointer items-center select-none">
                <UCheckbox
                  v-model="settings.filter_microsoft_nic"
                  name="filter_microsoft_nic"
                  color="primary"
                  class="me-2"
                />
                <span>{{ $t('setting.filterMicorsoftNic') }}</span>
              </label>
            </div>
          </section>

          <hr class="border-gray-100 dark:border-gray-800" />

          <section>
            <h3 class="mb-3 text-lg font-bold">{{ $t('setting.installOption') }}</h3>

            <div class="flex flex-col gap-y-3">
              <label class="flex w-fit cursor-pointer items-center select-none">
                <UCheckbox
                  v-model="settings.hide_not_found"
                  name="hide_not_found"
                  color="primary"
                  class="me-2"
                />
                <span>{{ $t('setting.hideNotFound') }}</span>
              </label>
            </div>
          </section>
        </div>
      </div>
    </main>

    <footer
      class="flex shrink-0 items-center justify-end border-t border-gray-200 bg-gray-50/50 px-6 py-3 dark:border-gray-800 dark:bg-gray-900/20"
    >
      <UButton type="submit" color="secondary" class="px-6 font-medium">
        {{ $t('common.save') }}
      </UButton>
    </footer>
  </form>
</template>
