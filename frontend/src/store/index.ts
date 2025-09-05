import { ExecutableExists } from '@/wailsjs/go/main/App'
import { storage } from '@/wailsjs/go/models'

import { defineStore } from 'pinia'
import { computed, ref, toRaw, watch } from 'vue'

export const useAppSettingStore = defineStore('appSetting', () => {
  const settings = ref<storage.AppSetting>(new storage.AppSetting())

  return {
    settings,
    editor: () => {
      const settingsClone = ref(structuredClone(toRaw(settings.value)))
      return {
        settings: settingsClone,
        modified: computed(
          () => JSON.stringify(settingsClone.value) != JSON.stringify(settings.value)
        ),
        reset: () => (settingsClone.value = structuredClone(toRaw(settings.value)))
      }
    }
  }
})

export const useDriverGroupStore = defineStore('driverGroup', () => {
  const groups = ref<storage.DriverGroup[]>([])
  const notFoundDrivers = ref<Array<string>>([])

  const findNotExists = (drivers: Array<storage.Driver>) =>
    Promise.all(
      drivers.map(d => ExecutableExists(d.path).then(exist => ({ id: d.id, exist: exist })))
    ).then(results => {
      return results
        .map(result => (result.exist ? undefined : result.id))
        .filter(v => v !== undefined)
    })

  watch(
    groups,
    newGroups =>
      findNotExists(newGroups.flatMap(g => g.drivers)).then(ids => (notFoundDrivers.value = ids)),
    { immediate: true }
  )

  return {
    groups,
    notFoundDrivers,
    editor: (id: string | null | undefined, defaultType?: storage.DriverType) => {
      const groupClone = ref<storage.DriverGroup>(
        structuredClone(
          toRaw(
            groups.value.find(g => g.id == id) ??
              new storage.DriverGroup({ type: defaultType, name: '', drivers: [] })
          )
        )
      )
      const notFoundDrivers = ref<Array<string>>([])

      watch(
        groupClone.value.drivers,
        newDrivers => findNotExists(newDrivers).then(ids => (notFoundDrivers.value = ids)),
        { immediate: true }
      )

      return {
        group: groupClone,
        notFoundDrivers,
        modified: computed(
          () =>
            JSON.stringify(groupClone.value) !=
            JSON.stringify(
              groups.value.find(g => g.id == groupClone.value.id) ||
                new storage.DriverGroup({ type: defaultType, name: '', drivers: [] })
            )
        ),
        modifiedDrivers: computed(
          () =>
            JSON.stringify(groupClone.value.drivers) !=
            JSON.stringify(groups.value.find(g => g.id == groupClone.value.id)?.drivers || [])
        ),
        reset: () => {
          groupClone.value = structuredClone(
            toRaw(
              groups.value.find(g => g.id == groupClone.value.id) ||
                new storage.DriverGroup({ type: defaultType, name: '', drivers: [] })
            )
          )
        }
      }
    }
  }
})
