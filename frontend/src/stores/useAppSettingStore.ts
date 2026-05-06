import { storage } from '@/wailsjs/go/models'
import { defineStore } from 'pinia'

export default defineStore('appSetting', () => {
  const settings = ref<storage.AppSetting>(new storage.AppSetting())

  return {
    settings
  }
})
