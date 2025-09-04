import { ExecutableExists } from '@/wailsjs/go/main/App'
import { storage } from '@/wailsjs/go/models'
import * as appSettingStorage from '@/wailsjs/go/storage/AppSettingStorage'
import * as driverGroupStorage from '@/wailsjs/go/storage/DriverGroupStorage'
import { defineStore } from 'pinia'
import { computed, ref, toRaw, watch } from 'vue'

export const useAppSettingStore = defineStore('appSetting', () => {
  const loading = ref(false)

  const settings = ref<storage.AppSetting>(new storage.AppSetting())

  return {
    loading,
    settings,
    read: async () => {
      loading.value = true
      return appSettingStorage
        .All()
        .then(s => (settings.value = s))
        .finally(() => (loading.value = false))
    },
    editor: () => {
      const settingsClone = ref(structuredClone(toRaw(settings.value)))
      return {
        settings: settingsClone,
        modified: computed(
          () => JSON.stringify(settingsClone.value) != JSON.stringify(settings.value)
        ),
        restore: () => (settingsClone.value = structuredClone(toRaw(settings.value))),
        save: () => {
          settings.value = structuredClone(toRaw(settingsClone.value))
          settingsClone.value = structuredClone(toRaw(settings.value))
        }
      }
    }
  }
})

export const useDriverGroupStore = defineStore('driverGroup', () => {
  const loading = ref(false)

  const groups = ref<storage.DriverGroup[]>([])
  const notFoundDrivers = ref<Array<string>>([])

  watch(
    groups,
    newGroups =>
      Promise.all(
        newGroups
          .flatMap(g => g.drivers)
          .flatMap(d => ExecutableExists(d.path).then(exist => ({ id: d.id, exist: exist })))
      )
        .then(results => {
          return results
            .map(result => (result.exist ? undefined : result.id))
            .filter(v => v !== undefined)
        })
        .then(ids => (notFoundDrivers.value = ids)),
    { immediate: true }
  )

  return {
    loading,
    groups,
    notFoundDrivers,
    read: async () => {
      loading.value = true
      return driverGroupStorage
        .All()
        .then(g => (groups.value = g))
        .finally(() => (loading.value = false))
    }
  }
})
