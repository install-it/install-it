<script setup lang="ts">
import { AppVersion } from '@/wailsjs/go/main/App'
import * as app_manager from '@/wailsjs/go/store/AppSettingStore.js'
import { onBeforeMount } from 'vue'
import { useI18n } from 'vue-i18n'
import type { RouteLocationRaw } from 'vue-router'
import { useToast } from 'vue-toast-notification'
import { latestRelease } from './utils'

const { t, locale } = useI18n()

const routes: Array<{ to: RouteLocationRaw; icon: string }> = [
  { to: '/', icon: 'fa-regular fa-house' },
  { to: '/drivers', icon: 'fa-regular fa-file-code' },
  { to: '/settings', icon: 'fa-solid fa-gear' },
  { to: '/porter', icon: 'fa-solid fa-people-arrows' },
  { to: '/app-info', icon: 'fa-solid fa-info' }
]

onBeforeMount(() => {
  app_manager.All().then(s => {
    locale.value = s.language

    if (s.auto_check_update) {
      AppVersion().then(version =>
        latestRelease(version).then(release => {
          if (release.hasUpdate) {
            useToast({ position: 'top-right' }).info(t('toast.updateAvailable'))
          }
        })
      )
    }
  })
})
</script>

<template>
  <div class="flex h-screen">
    <aside class="w-12">
      <div class="flex justify-center h-full bg-gray-50">
        <ul class="mt-6 space-y-3 font-medium">
          <li v-for="(link, i) in routes" :key="i">
            <RouterLink
              :to="link.to"
              class="flex p-2 rounded-lg hover:bg-gray-200"
              activeClass="text-apple-green-900 bg-powder-blue-400"
            >
              <font-awesome-icon :icon="link.icon" />
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
