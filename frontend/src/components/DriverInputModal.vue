<script setup lang="ts">
import DriverSelector from '@/components/DriverSelector.vue'
import ModalFrame from '@/components/modals/ModalFrame.vue'
import { useDriverGroupStore } from '@/store'
import { SelectFile } from '@/wailsjs/go/main/App'
import { storage } from '@/wailsjs/go/models'
import { nextTick, ref, toRaw, useTemplateRef } from 'vue'

defineEmits<{ submit: [dri: storage.Driver] }>()

const frame = useTemplateRef('frame')

defineExpose({
  show: (data?: Partial<storage.Driver>) => {
    frame.value?.show()

    if (data) {
      driver.value = {
        ...data,
        flags: data.flags?.join(','),
        allowRtCodes: data.allowRtCodes?.join(',')
      }
    } else {
      driver.value = { minExeTime: 5, incompatibles: [] }
    }

    nextTick(() => {
      // wait for the modal to open
      modalBody.value?.scrollTo({ top: 0, behavior: 'smooth' })
    })
  },
  hide: frame.value?.hide || (() => {})
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
  Partial<Omit<storage.Driver, 'allowRtCodes' | 'flags'> & { allowRtCodes: string; flags: string }>
>({})
</script>

<template>
  <ModalFrame ref="frame" :on-demand="true" :immediate="false">
    <div class="w-[75vw] max-w-[650px]">
      <!-- Modal content -->
      <div class="rounded-lg bg-white shadow-sm">
        <!-- Modal header -->
        <div class="flex h-12 items-center justify-between rounded-t border-b px-4">
          <h3 class="font-semibold">
            {{ driver ? $t('driverForm.editDriver') : $t('driverForm.createDriver') }}
          </h3>

          <button
            type="button"
            class="rounded-lg bg-transparent p-3 text-sm text-gray-400 hover:bg-gray-100 hover:text-gray-900"
            @click="frame?.hide()"
          >
            <font-awesome-icon icon="fa-solid fa-xmark" />
          </button>
        </div>

        <!-- Modal body -->
        <div ref="modalBody" class="max-h-[70vh] overflow-auto px-4 py-2">
          <form
            class="flex flex-col gap-y-2"
            autocomplete="off"
            @submit.prevent="
              _ => {
                $emit(
                  'submit',
                  new storage.Driver({
                    ...driver,
                    flags: driver.flags ? driver.flags.split(',') : [],
                    allowRtCodes: driver.allowRtCodes
                      ? driver.allowRtCodes
                          ?.split(',')
                          .map(c => parseInt(c))
                          .filter(c => !Number.isNaN(c))
                      : [],
                    incompatibles: toRaw(driver.incompatibles) ?? []
                  })
                )

                frame?.hide()
              }
            "
          >
            <fieldset class="fieldset">
              <legend class="fieldset-legend text-sm">{{ $t('driverForm.name') }}</legend>

              <input
                v-model="driver.name"
                type="text"
                name="name"
                class="input w-full input-accent"
              />
            </fieldset>

            <fieldset class="fieldset">
              <legend class="fieldset-legend text-sm">{{ $t('driverForm.path') }}</legend>

              <div class="join">
                <button
                  type="button"
                  class="btn join-item w-32"
                  @click="
                    SelectFile(true).then(path => {
                      driver.path = path
                    })
                  "
                >
                  {{ $t('driverForm.selectFile') }}
                </button>

                <input
                  ref="pathInput"
                  v-model="driver.path"
                  type="text"
                  name="path"
                  class="input join-item w-full input-accent"
                  required
                />
              </div>
            </fieldset>

            <fieldset class="fieldset">
              <legend class="fieldset-legend text-sm">{{ $t('driverForm.argument') }}</legend>

              <div class="join items-center">
                <!-- <select
                  name="flags"
                  class="w-32 select select-accent join-item ps-1"
                  @change="
                    event => {
                      driver.flags = (event.target as HTMLSelectElement).value
                    }
                  "
                >
                  <option value="">{{ $t('driverForm.manualInput') }}</option>
                  <option
                    v-for="(flag, name) in FLAGS"
                    :key="name"
                    :value="flag.join(',')"
                    :selected="driver.flags === flag.join()"
                  >
                    {{ name }}
                  </option>
                </select> -->
                <div class="dropdown">
                  <div tabindex="0" role="button" class="btn m-1 join-item w-30">
                    {{ $t('common.select') }}
                  </div>

                  <ul
                    tabindex="0"
                    class="dropdown-content menu z-1 w-52 rounded-box bg-base-100 p-2 shadow-sm"
                  >
                    <div class="h-36 overflow-y-auto">
                      <li
                        v-for="(flag, name) in FLAGS"
                        :key="name"
                        @click="
                          event => {
                            driver.flags = flag.join()
                            ;(
                              event.currentTarget as HTMLLIElement
                            ).parentElement?.parentElement?.blur()
                          }
                        "
                      >
                        <a>
                          {{ name }}
                        </a>
                      </li>
                    </div>
                  </ul>
                </div>

                <input
                  v-model="driver.flags"
                  type="text"
                  name="flags"
                  class="input join-item w-full input-accent"
                />
              </div>

              <p class="label text-apple-green-800">
                {{ $t('driverForm.commaSeparated') }}
              </p>
            </fieldset>

            <div class="flex gap-x-3">
              <fieldset class="fieldset flex-1">
                <legend class="fieldset-legend text-sm">
                  {{ $t('driverForm.minExecuteTime') }}
                </legend>

                <input
                  v-model="driver.minExeTime"
                  type="number"
                  name="minExeTime"
                  step="0.1"
                  class="input w-full input-accent"
                  required
                />

                <p class="label text-wrap text-apple-green-800">
                  {{ $t('driverForm.minExecuteTimeHelp') }}
                </p>
              </fieldset>

              <fieldset class="fieldset flex-1">
                <legend class="fieldset-legend text-sm">
                  {{ $t('driverForm.allowedExitCode') }}
                </legend>

                <input
                  v-model="driver.allowRtCodes"
                  type="text"
                  name="allowRtCodes"
                  class="input input-accent"
                />

                <p class="label text-wrap text-apple-green-800">
                  {{ $t('driverForm.commaSeparated') }}
                </p>
              </fieldset>
            </div>

            <DriverSelector
              v-model="driver.incompatibles"
              group-by="driver"
              :driver-groups="groupStore.groups"
            ></DriverSelector>

            <button type="submit" class="btn btn-secondary">
              {{ $t('common.save') }}
            </button>
          </form>
        </div>
      </div>
    </div>
  </ModalFrame>
</template>

<style scoped>
legend:has(+ input:required, + select:required):after,
legend:has(+ div > input:required):after {
  content: ' *';
  color: red;
}
</style>
