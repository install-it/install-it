import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import useAppSettingStore from '@/stores/useAppSettingStore'
import { storage } from '@/wailsjs/go/models'

describe('useAppSettingStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('initializes with a default AppSetting instance', () => {
    const store = useAppSettingStore()
    expect(store.settings).toBeInstanceOf(storage.AppSetting)
  })

  it('default AppSetting has the expected property shape', () => {
    const store = useAppSettingStore()
    expect(store.settings).toHaveProperty('create_partition')
    expect(store.settings).toHaveProperty('set_password')
    expect(store.settings).toHaveProperty('language')
    expect(store.settings).toHaveProperty('success_action')
    expect(store.settings).toHaveProperty('parallel_install')
  })

  it('settings ref is writable with a new AppSetting', () => {
    const store = useAppSettingStore()
    const updated = new storage.AppSetting({ language: 'zh-hk' })
    store.settings = updated
    expect(store.settings.language).toBe('zh-hk')
  })

  it('settings ref is reactive — field mutations are tracked', () => {
    const store = useAppSettingStore()
    store.settings.language = 'en'
    expect(store.settings.language).toBe('en')
    store.settings.language = 'zh-hk'
    expect(store.settings.language).toBe('zh-hk')
  })

  it('parallel_install can be toggled', () => {
    const store = useAppSettingStore()
    store.settings.parallel_install = true
    expect(store.settings.parallel_install).toBe(true)
    store.settings.parallel_install = false
    expect(store.settings.parallel_install).toBe(false)
  })

  it('each store instance is independent when pinia is recreated', () => {
    const store1 = useAppSettingStore()
    store1.settings.language = 'en'

    setActivePinia(createPinia())
    const store2 = useAppSettingStore()
    // A fresh pinia means a fresh store — language should be reset
    expect(store2.settings.language).toBeUndefined()
  })
})
