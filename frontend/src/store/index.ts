import { ExecutableExists } from '@/wailsjs/go/main/App'
import { storage } from '@/wailsjs/go/models'

import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useAppSettingStore = defineStore('appSetting', () => {
  const settings = ref<storage.AppSetting>(new storage.AppSetting())

  return {
    settings
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
    isAllDriversExist: (g: storage.DriverGroup) =>
      g.drivers.flatMap(d => d.id).every(id => !notFoundDrivers.value.includes(id))
  }
})

export const useMatchRuleStore = defineStore('matchRuleGroup', () => {
  const ruleSets = ref<storage.RuleSet[]>([])

  return {
    ruleSets
  }
})

export const useUnsavedFormStore = defineStore('unsavedForm', () => {
  const show = ref(false)

  let answerHandler: ((allow: boolean) => void) | null = null

  return {
    show,
    confirmLeave: (answer: boolean) => {
      show.value = false
      if (answerHandler) {
        answerHandler(answer)
        answerHandler = null
      }
    },
    setAnswerHandler: (handler: (allow: boolean) => void) => {
      answerHandler = handler
    }
  }
})
