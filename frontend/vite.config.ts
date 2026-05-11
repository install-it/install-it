import { fileURLToPath, URL } from 'node:url'

import ui from '@nuxt/ui/vite'
import vue from '@vitejs/plugin-vue'
import { normalize, resolve } from 'pathe'
import { convertPathToPattern, glob } from 'tinyglobby'
import { scanExports } from 'unimport'
import { defineConfig } from 'vite'
import VueRouter from 'vue-router/vite'

async function scanAutoImports(...dirs: string[]): Promise<
  {
    from: string
    imports: (string | [string, string])[]
  }[]
> {
  const fileToSpecifier = new Map<string, string>()
  const patterns = dirs.map(
    dir =>
      `${convertPathToPattern(
        normalize(resolve(new URL('.', import.meta.url).pathname, dir))
      )}/**/*.{ts,js,mjs,cjs,mts,cts}`
  )

  const files = await glob(patterns, { absolute: true, ignore: ['**/*.d.ts'] })

  const scanResults = await Promise.all(
    files.sort().map(async file => {
      const abs = normalize(file)
      // the 'from' specifier for the import (no extension, forward slashes)
      fileToSpecifier.set(abs, abs.replace(/\.[^.]+$/, ''))
      return { abs, exports: await scanExports(abs, false) }
    })
  )

  const groupedImports = new Map<string, (string | [string, string])[]>()
  const seenImports = new Set<string>()

  for (const { abs, exports } of scanResults) {
    const from = fileToSpecifier.get(abs)
    if (!from) {
      continue
    }

    for (const entry of exports) {
      if (entry.type) {
        continue
      }

      const alias = entry.as ?? entry.name
      const dedupeKey = `${from}:${entry.name}:${alias}`

      if (seenImports.has(dedupeKey)) {
        continue
      }
      seenImports.add(dedupeKey)

      const fileImports = groupedImports.get(from) ?? []
      fileImports.push(alias === entry.name ? entry.name : [entry.name, alias])
      groupedImports.set(from, fileImports)
    }
  }

  return [...groupedImports.entries()]
    .sort(([l], [r]) => l.localeCompare(r))
    .map(([from, fileImports]) => ({ from, imports: fileImports }))
}

// https://vite.dev/config/
export default defineConfig(async () => {
  return {
    plugins: [
      VueRouter(),
      vue(),
      ui({
        autoImport: {
          imports: await scanAutoImports('./src/composables', './src/stores')
        },
        components: {
          directoryAsNamespace: true
        },
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
  }
})
