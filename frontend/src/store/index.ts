import { storage } from '@/wailsjs/go/models'
import * as appManager from '@/wailsjs/go/storage/AppSettingManager'
import { defineStore } from 'pinia'
import { readonly, ref, toRaw } from 'vue'

export const useAppSettingStore = defineStore('appSetting', () => {
  const loading = ref(false)
  const settings = ref<storage.AppSetting>(new storage.AppSetting())
  const settingOriginal = ref(settings.value)

  return {
    loading,
    settings,
    settingOriginal: readonly(settingOriginal),
    read: async () => {
      loading.value = true
      return appManager
        .Read()
        .then(s => {
          settings.value = s
          settingOriginal.value = structuredClone(s)
        })
        .finally(() => (loading.value = false))
    },
    write: () =>
      appManager
        .Update(settings.value)
        .then(() => (settingOriginal.value = structuredClone(toRaw(settings.value))))
  }
})
