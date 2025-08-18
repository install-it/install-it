import { i18n } from '@/i18n'
import { library } from '@fortawesome/fontawesome-svg-core'
import { fab } from '@fortawesome/free-brands-svg-icons'
import { far } from '@fortawesome/free-regular-svg-icons'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { createApp } from 'vue'
import { LoadingPlugin } from 'vue-loading-overlay'
import 'vue-loading-overlay/dist/css/index.css'
import ToastPlugin from 'vue-toast-notification'
import 'vue-toast-notification/dist/theme-bootstrap.css'
import App from './App.vue'
import './assets/main.css'
import router from './router'

library.add(fas, far, fab)

const app = createApp(App)
  .use(router)
  .use(ToastPlugin, { position: 'top-right' })
  .use(i18n)
  .use(LoadingPlugin)
  .component('font-awesome-icon', FontAwesomeIcon)

app.config.globalProperties.$window = window

app.mount('#app')
