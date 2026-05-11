import { storage } from '@/wailsjs/go/models'
import { defineStore } from 'pinia'

export default defineStore('matchRuleGroup', () => {
  const ruleSets = ref<storage.RuleSet[]>([])

  return {
    ruleSets
  }
})
