import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import VueRouter from 'vue-router/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [tailwindcss(), VueRouter(), vue()],
  resolve: {
    alias: {
      '@/wailsjs': fileURLToPath(new URL('./wailsjs', import.meta.url)),
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
