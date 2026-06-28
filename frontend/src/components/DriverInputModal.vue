<script setup lang="ts">
import DriverSelector from '@/components/DriverSelector.vue'
import { SelectFile } from '@/wailsjs/go/main/App'
import { storage } from '@/wailsjs/go/models'
import { nextTick, ref, toRaw, useTemplateRef } from 'vue'

const emit = defineEmits<{ submit: [dri: storage.Driver] }>()

const isOpen = ref(false)

defineExpose({
  show: (data?: Partial<storage.Driver>) => {
    isOpen.value = true

    if (data) {
      driver.value = {
        ...data,
        flags: data.flags?.join(','),
        allowRtCodes: data.allowRtCodes?.join(','),
        incompatibles: Array.isArray(data.incompatibles) ? data.incompatibles : []
      }
    } else {
      driver.value = { minExeTime: 5, incompatibles: [] }
    }

    nextTick(() => {
      // wait for the modal to open
      modalBody.value?.scrollTo({ top: 0, behavior: 'smooth' })
    })
  },
  hide: () => {
    isOpen.value = false
  }
})

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

const modalBody = useTemplateRef<HTMLDivElement>('modalBody')

const driver = ref<
  Partial<
    Omit<storage.Driver, 'allowRtCodes' | 'flags' | 'incompatibles'> & {
      allowRtCodes: string
      flags: string
      incompatibles: string[]
    }
  >
>({ incompatibles: [] })

const flagItems = Object.entries(FLAGS).map(([name, flags]) => ({
  label: name,
  onSelect: () => {
    driver.value.flags = flags.join(',')
  }
}))

function handleSubmit() {
  emit(
    'submit',
    new storage.Driver({
      ...driver.value,
      flags: driver.value.flags ? driver.value.flags.split(',') : [],
      allowRtCodes: driver.value.allowRtCodes
        ? driver.value.allowRtCodes
            ?.split(',')
            .map((c: string) => parseInt(c))
            .filter((c: number) => !Number.isNaN(c))
        : [],
      incompatibles: toRaw(driver.value.incompatibles) ?? []
    })
  )

  isOpen.value = false
}
</script>

<template>
  <UModal
    v-model:open="isOpen"
    :title="driver.name ? $t('driverFormEditDriver') : $t('driverFormCreateDriver')"
  >
    <template #body>
      <div ref="modalBody">
        <form class="flex flex-col gap-y-2" autocomplete="off" @submit.prevent="handleSubmit">
          <fieldset class="fieldset">
            <legend class="fieldset-legend text-sm">{{ $t('driverFormName') }}</legend>

            <UInput v-model="driver.name" type="text" name="name" color="primary" class="w-full" />
          </fieldset>

          <fieldset class="fieldset">
            <legend class="fieldset-legend text-sm">{{ $t('driverFormPath') }}</legend>

            <div class="flex gap-2">
              <UButton
                type="button"
                color="neutral"
                variant="outline"
                @click="
                  SelectFile(true).then(path => {
                    driver.path = path
                  })
                "
              >
                {{ $t('driverFormSelectFile') }}
              </UButton>

              <UInput
                ref="pathInput"
                v-model="driver.path"
                type="text"
                name="path"
                color="primary"
                class="flex-1"
                required
              />
            </div>
          </fieldset>

          <fieldset class="fieldset">
            <legend class="fieldset-legend text-sm">{{ $t('driverFormArgument') }}</legend>

            <div class="flex items-center gap-2">
              <UDropdownMenu :items="[flagItems]" :ui="{ content: 'max-h-58 overflow-y-auto' }">
                <UButton color="neutral" variant="outline">
                  {{ $t('commonSelect') }}
                </UButton>
              </UDropdownMenu>

              <UInput
                v-model="driver.flags"
                type="text"
                name="flags"
                color="primary"
                class="flex-1"
              />
            </div>

            <p class="text-hint">
                {{ $t('driverFormCommaSeparated') }}
            </p>
          </fieldset>

          <div class="flex gap-x-3">
            <fieldset class="fieldset flex-1">
              <legend class="fieldset-legend text-sm">
                {{ $t('driverFormMinExecuteTime') }}
              </legend>

              <UInput
                v-model="driver.minExeTime"
                type="number"
                name="minExeTime"
                step="0.1"
                color="primary"
                required
              />

              <p class="text-hint">
                {{ $t('driverFormMinExecuteTimeHelp') }}
              </p>
            </fieldset>

            <fieldset class="fieldset flex-1">
              <legend class="fieldset-legend text-sm">
                {{ $t('driverFormAllowedExitCode') }}
              </legend>

              <UInput
                v-model="driver.allowRtCodes"
                type="text"
                name="allowRtCodes"
                color="primary"
              />

              <p class="text-hint">
              {{ $t('driverFormCommaSeparated') }}
              </p>
            </fieldset>
          </div>

          <DriverSelector
            v-model="driver.incompatibles"
            group-by="driver"
            :driver-groups="groupStore.groups"
          ></DriverSelector>

          <UButton type="submit" color="secondary" block class="justify-center">
            {{ $t('commonSave') }}
          </UButton>
        </form>
      </div>
    </template>
  </UModal>
</template>

<style scoped>
legend:has(+ input:required, + select:required):after,
legend:has(+ div > input:required):after {
  content: ' *';
  color: red;
}
</style>
