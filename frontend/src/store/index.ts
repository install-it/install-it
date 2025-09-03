import { storage } from '@/wailsjs/go/models'
import * as appManager from '@/wailsjs/go/storage/AppSettingManager'
import { defineStore } from 'pinia'
import { computed, ref, toRaw } from 'vue'

export const useAppSettingStore = defineStore('appSetting', () => {
  const loading = ref(false)

  const settings = ref<storage.AppSetting>(new storage.AppSetting())
  const original = ref(structuredClone(toRaw(settings.value)))

  return {
    loading,
    settings,
    modified: computed(() => JSON.stringify(original.value) !== JSON.stringify(settings.value)),
    restore: () => (settings.value = structuredClone(toRaw(original.value))),
    read: async () => {
      loading.value = true
      return appManager
        .Read()
        .then(s => {
          settings.value = s
          original.value = structuredClone(s)
        })
        .finally(() => (loading.value = false))
    },
    write: () =>
      appManager
        .Update(settings.value)
        .then(() => (original.value = structuredClone(toRaw(settings.value))))
  }
})
