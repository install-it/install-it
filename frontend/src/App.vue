<script setup lang="ts">
import { useAppSettingStore, useDriverGroupStore } from '@/store'
import { AppVersion } from '@/wailsjs/go/main/App'
import * as appSettingStorage from '@/wailsjs/go/storage/AppSettingStorage'
import * as driverGroupStorage from '@/wailsjs/go/storage/DriverGroupStorage'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { RouteLocationRaw } from 'vue-router'
import { useToast } from 'vue-toast-notification'
import { latestRelease } from './utils'

const { t, locale } = useI18n()

const $toast = useToast({ position: 'top-right' })

const settingsStore = useAppSettingStore()

const groupStore = useDriverGroupStore()

const initilisating = ref(true)

const hasUpdate = ref(false)

Promise.all([
  driverGroupStorage
    .All()
    .then(gs => (groupStore.groups = gs))
    .catch(() => {
      $toast.error(t('toast.readDriverFailed'))
    }),
  appSettingStorage
    .All()
    .then(s => {
      settingsStore.settings = s
      locale.value = s.language
    })
    .catch(() => {
      $toast.error(t('toast.readAppSettingFailed'))
    })
])
  .then(() => {
    setTimeout(() => {
      if (settingsStore.settings.auto_check_update) {
        return AppVersion().then(version =>
          latestRelease(version).then(release => {
            hasUpdate.value = release.hasUpdate
            if (release.hasUpdate) {
              $toast.info(t('toast.updateAvailable'))
            }
          })
        )
      }
    }, 1000)
  })
  .finally(() => (initilisating.value = false))

const routes: Array<{ to: RouteLocationRaw; icon: string }> = [
  { to: '/', icon: 'fa-regular fa-house' },
  { to: '/drivers', icon: 'fa-regular fa-file-code' },
  { to: '/settings', icon: 'fa-solid fa-gear' },
  { to: '/porter', icon: 'fa-solid fa-people-arrows' },
  { to: '/app-info', icon: 'fa-solid fa-info' }
]
</script>

<template>
  <div class="flex h-screen">
    <Transition name="fade" mode="out-in">
      <template v-if="!initilisating">
        <div class="flex w-full">
          <aside class="w-12">
            <div class="flex justify-center h-full bg-gray-50">
              <ul class="mt-6 space-y-3 font-medium">
                <li v-for="(link, i) in routes" :key="i">
                  <RouterLink
                    :to="link.to"
                    class="flex p-2 rounded-lg hover:bg-gray-200"
                    activeClass="text-apple-green-900 bg-powder-blue-400"
                    draggable="false"
                  >
                    <div class="indicator">
                      <span
                        class="indicator-item status status-neutral"
                        style="background-image: unset"
                        v-if="link.to == '/app-info' && hasUpdate"
                      ></span>
                      <font-awesome-icon :icon="link.icon" />
                    </div>
                  </RouterLink>
                </li>
              </ul>
            </div>
          </aside>

          <main class="w-full h-full p-3">
            <RouterView></RouterView>
          </main>
        </div>
      </template>

      <template v-else>
        <div class="m-auto">
          <span class="loading loading-dots loading-xl"></span>
        </div>
      </template>
    </Transition>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease-in;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
