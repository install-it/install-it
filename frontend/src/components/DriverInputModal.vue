<script setup lang="ts">
import DriverSelector from '@/components/DriverSelector.vue'
import TaggedInput from '@/components/TaggedInput.vue'
import { ExecutableExists, SelectFile } from '@/wailsjs/go/main/App'
import { storage } from '@/wailsjs/go/models'
import { nextTick, onUnmounted, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const toast = useToast()

const props = defineProps<{
  open: boolean
  editData?: Partial<storage.Driver>
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  submit: [dri: storage.Driver]
}>()

const FLAGS = {
  'Intel LAN': ['/s'],
  'Realtek LAN': ['-s'],
  'Nvidia Display': ['-s', '-noreboot', 'Display.Driver'],
  'AMD Display': ['-install'],
  'Intel Display': ['-s', '--noExtras'],
  'Intel Wifi': ['-q'],
  'Intel BT': ['/quiet', '/norestart'],
  'Intel Chipset': ['-s', '-norestart'],
  'AMD Chipset': ['/S']
}

const groupStore = useDriverGroupStore()

const modalBody = ref<HTMLDivElement | null>(null)

const driver = ref<{
  id?: number
  type?: storage.DriverType
  name?: string
  path?: string
  flags: string[]
  minExeTime: number
  allowRtCodes: string[]
  incompatibles: number[]
}>({ flags: [], minExeTime: 5, allowRtCodes: [], incompatibles: [] })

const pathExists = ref<boolean | null>(null)
let pathCheckTimeout: ReturnType<typeof setTimeout> | null = null

const flagItems = Object.entries(FLAGS).map(([name, flags]) => ({
  label: name,
  onSelect: () => {
    driver.value.flags = [...flags]
  }
}))

watch(
  () => props.open,
  val => {
    if (!val) return

    if (props.editData) {
      driver.value = {
        id: props.editData.id,
        type: props.editData.type,
        name: props.editData.name,
        path: props.editData.path,
        flags: props.editData.flags ?? [],
        minExeTime: props.editData.minExeTime ?? 5,
        allowRtCodes: props.editData.allowRtCodes?.map(c => String(c)) ?? [],
        incompatibles: Array.isArray(props.editData.incompatibles)
          ? props.editData.incompatibles
          : []
      }
    } else {
      driver.value = { flags: [], minExeTime: 5, allowRtCodes: [], incompatibles: [] }
    }

    pathExists.value = null

    nextTick(() => {
      modalBody.value?.scrollTo({ top: 0, behavior: 'smooth' })
    })
  }
)

watch(
  () => driver.value.path,
  path => {
    pathExists.value = null

    if (pathCheckTimeout) {
      clearTimeout(pathCheckTimeout)
    }

    if (!path) return

    pathCheckTimeout = setTimeout(() => {
      ExecutableExists(path)
        .then(exists => {
          pathExists.value = exists
        })
        .catch(() => {
          pathExists.value = false
        })
    }, 300)
  }
)

onUnmounted(() => {
  if (pathCheckTimeout) clearTimeout(pathCheckTimeout)
})

function handleSubmit() {
  if (!driver.value.name?.trim() || !driver.value.path?.trim()) {
    toast.add({ title: t('toastNameAndPathRequired'), color: 'warning' })
    return
  }

  emit(
    'submit',
    new storage.Driver({
      ...driver.value,
      minExeTime: Number(driver.value.minExeTime) || 5,
      allowRtCodes: driver.value.allowRtCodes
        .map(c => parseInt(c))
        .filter(c => !Number.isNaN(c)),
      incompatibles: toRaw(driver.value.incompatibles)
    })
  )
  emit('update:open', false)
}
</script>

<template>
  <UModal
    :open="props.open"
    :title="props.editData ? t('titleEditDriver') : t('titleCreateDriver')"
    class="max-w-2xl"
    @update:open="emit('update:open', $event)"
  >
    <template #body>
      <div ref="modalBody">
        <form
          id="driver-form"
          class="flex flex-col gap-y-5"
          autocomplete="off"
          @submit.prevent="handleSubmit"
        >
          <fieldset>
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-1">
              {{ t('name') }} <span class="text-red-500">*</span>
            </label>
            <UInput
              v-model="driver.name"
              type="text"
              name="name"
              color="primary"
              class="w-full"
              required
            />
          </fieldset>

          <fieldset>
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-1">
              {{ t('path') }} <span class="text-red-500">*</span>
            </label>
            <div class="flex gap-2">
              <UButton
                type="button"
                color="neutral"
                variant="outline"
                icon="mdi:folder-open"
                @click="
                  SelectFile(true).then(path => {
                    if (path) driver.path = path
                  })
                "
              >
                {{ t('labelSelectFile') }}
              </UButton>
              <UInput
                v-model="driver.path"
                type="text"
                name="path"
                class="flex-1 font-mono"
                required
              />
            </div>
          </fieldset>

          <fieldset>
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-1">
              {{ t('fieldArgument') }}
            </label>
            <div class="flex items-center gap-2">
              <UDropdownMenu :items="[flagItems]" :ui="{ content: 'max-h-58 overflow-y-auto' }">
                <UButton color="neutral" variant="outline">
                  {{ t('select') }}
                </UButton>
              </UDropdownMenu>
            </div>
            <TaggedInput v-model="driver.flags" :title="t('fieldArgument')" />
            <p class="text-hint">
              {{ t('descCommaSeparated') }}
            </p>
          </fieldset>

          <div class="flex gap-x-3">
            <fieldset class="flex-1">
              <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-1">
                {{ t('fieldMinExecuteTime') }}
              </label>
              <UInput
                v-model.number="driver.minExeTime"
                type="number"
                name="minExeTime"
                step="0.1"
                color="primary"
                required
              />
              <p class="text-hint">
                {{ t('descMinExecuteTime') }}
              </p>
            </fieldset>

            <fieldset class="flex-1">
              <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-1">
                {{ t('fieldAllowedExitCode') }}
              </label>
              <TaggedInput v-model="driver.allowRtCodes" :title="t('fieldAllowedExitCode')" />
              <p class="text-hint">
                {{ t('descCommaSeparated') }}
              </p>
            </fieldset>
          </div>

          <details class="group border border-gray-200 rounded-lg bg-gray-50 shadow-sm overflow-hidden">
            <summary class="flex justify-between items-center p-3 cursor-pointer text-sm font-bold text-gray-700 select-none list-none">
              <div class="flex items-center gap-2">
                <Icon icon="mdi:alert-octagram-outline" class="text-amber-500 text-lg" />
                {{ t('labelIncompatibleWith') }}
                <span v-if="driver.incompatibles.length"
                  class="bg-amber-100 text-amber-700 px-1.5 py-0.5 rounded text-xs border border-amber-200">
                  {{ driver.incompatibles.length }} set
                </span>
              </div>
              <Icon icon="mdi:chevron-down" class="transition-transform group-open:rotate-180" />
            </summary>
            <div class="p-3 border-t border-gray-200 bg-white">
              <p class="text-xs text-gray-500 mb-2 font-medium">{{ t('descIncompatible') }}</p>
              <DriverSelector
                v-model="driver.incompatibles"
                group-by="driver"
                :driver-groups="groupStore.groups"
              />
            </div>
          </details>
        </form>
      </div>
    </template>

    <template #footer>
      <div class="flex justify-end gap-2">
        <UButton color="neutral" variant="ghost" @click="emit('update:open', false)">
          {{ t('cancel') }}
        </UButton>
        <UButton color="primary" type="submit" form="driver-form">
          {{ t('save') }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>
