import { createI18n } from 'vue-i18n'
import en from './en.json'
import zh_Hant_HK from './zh_Hant_HK.json'

const i18n = createI18n({
  legacy: false,
  locale: 'zh_Hant_HK',
  messages: {
    en,
    zh_Hant_HK
  },
  fallbackLocale: 'en'
})

export { i18n }
