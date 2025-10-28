<script setup lang="ts">
import { useAppSettingStore, useDriverGroupStore, useMatchRuleStore } from '@/store'
import { AppVersion } from '@/wailsjs/go/main/App'
import * as appSettingStorage from '@/wailsjs/go/storage/AppSettingStorage'
import * as driverGroupStorage from '@/wailsjs/go/storage/DriverGroupStorage'
import * as matchRuleStorage from '@/wailsjs/go/storage/MatchRuleStorage'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { RouteLocationRaw } from 'vue-router'
import { useToast } from 'vue-toast-notification'
import { latestRelease } from './utils'

const { t, locale } = useI18n()

const $toast = useToast({ position: 'top-right' })

const [settingsStore, groupStore, matchStore] = [
  useAppSettingStore(),
  useDriverGroupStore(),
  useMatchRuleStore()
]

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
    }),
  matchRuleStorage
    .All()
    .then(rs => {
      matchStore.ruleSets = rs
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
  { to: '/match-rules', icon: 'fa-solid fa-location-crosshairs' },
  { to: '/settings', icon: 'fa-solid fa-gear' },
  { to: '/porter', icon: 'fa-solid fa-people-arrows' },
  { to: '/app-info', icon: 'fa-solid fa-info' }
]
</script>

<template>
  <Transition name="fade" mode="out-in">
    <template v-if="!initilisating">
      <div class="flex h-screen w-screen">
        <aside class="w-12">
          <div class="flex h-full justify-center bg-gray-50">
            <ul class="mt-6 space-y-3 font-medium">
              <li v-for="(link, i) in routes" :key="i">
                <RouterLink
                  :to="link.to"
                  class="flex rounded-lg p-2 hover:bg-gray-200"
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

        <main class="flex-1 overflow-hidden p-3">
          <RouterView></RouterView>
        </main>
      </div>
    </template>

    <template v-else>
      <div class="flex h-screen w-screen justify-center">
        <span class="loading loading-xl loading-dots"></span>
      </div>
    </template>
  </Transition>
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
