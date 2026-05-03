import { fileURLToPath, URL } from 'node:url'

import ui from '@nuxt/ui/vite'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import VueRouter from 'vue-router/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    tailwindcss(),
    VueRouter(),
    vue(),
    ui({
      theme: {
        colors: [
          'primary',
          'secondary',
          'tertiary',
          'accent',
          'info',
          'success',
          'warning',
          'error'
        ]
      },
      ui: {
        colors: {
          primary: 'powder-blue',
          secondary: 'half-baked',
          tertiary: 'kashmir-blue',
          accent: 'apple-green'
        },
        modal: {
          slots: {
            header: 'flex items-center gap-1.5 px-4 py-2 sm:px-6 min-h-12',
            close: 'absolute top-2 inset-e-4',
            content: 'lg:max-w-3xl'
          }
        }
      }
    })
  ],
  resolve: {
    alias: {
      '@/wailsjs': fileURLToPath(new URL('./wailsjs', import.meta.url)),
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
