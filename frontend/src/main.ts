import './assets/main.css'

import { i18n } from '@/i18n'
import { Icon } from '@iconify/vue'
import ui from '@nuxt/ui/vue-plugin'
import { createPinia } from 'pinia'
import { createApp } from 'vue'
import { LoadingPlugin } from 'vue-loading-overlay'
import 'vue-loading-overlay/dist/css/index.css'
import App from './App.vue'
import router from './router'

const app = createApp(App)
  .use(router)
  .use(createPinia())
  .use(i18n)
  .use(LoadingPlugin)
  .use(ui)
  .component('Icon', Icon)

app.config.globalProperties.$window = window

app.mount('#app')
