<script setup lang="ts">
import ModalFrame from '@/components/modals/ModalFrame.vue'
import { useDriverGroupStore } from '@/store'
import { SelectFile } from '@/wailsjs/go/main/App'
import { storage } from '@/wailsjs/go/models'
import { computed, nextTick, ref, toRaw, useTemplateRef } from 'vue'

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

defineEmits<{ submit: [dri: storage.Driver] }>()

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

const searchPhrase = ref('')

const driver = ref<
  Partial<Omit<storage.Driver, 'allowRtCodes' | 'flags'> & { allowRtCodes: string; flags: string }>
>({})

const filterGroups = computed(() => {
  return searchPhrase.value === ''
    ? groupStore.groups
    : groupStore.groups.filter(
        g =>
          g.name.includes(searchPhrase.value) ||
          g.drivers.some(d => d.name.includes(searchPhrase.value))
      )
})
</script>

<template>
  <ModalFrame :on-demand="true" :immediate="false" ref="frame">
    <div class="w-[75vw] max-w-[650px]">
      <!-- Modal content -->
      <div class="bg-white rounded-lg shadow-sm">
        <!-- Modal header -->
        <div class="flex items-center justify-between h-12 px-4 border-b rounded-t">
          <h3 class="font-semibold">
            {{ driver ? $t('driverForm.editDriver') : $t('driverForm.createDriver') }}
          </h3>

          <button
            type="button"
            class="p-3 text-sm text-gray-400 hover:text-gray-900 bg-transparent hover:bg-gray-100 rounded-lg"
            @click="frame?.hide()"
          >
            <font-awesome-icon icon="fa-solid fa-xmark" />
          </button>
        </div>

        <!-- Modal body -->
        <div class="max-h-[70vh] overflow-auto py-2 px-4" ref="modalBody">
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
                type="text"
                name="name"
                v-model="driver.name"
                class="input input-accent w-full"
              />
            </fieldset>

            <fieldset class="fieldset">
              <legend class="fieldset-legend text-sm">{{ $t('driverForm.path') }}</legend>

              <div class="join">
                <button
                  type="button"
                  class="w-32 btn join-item"
                  @click="
                    SelectFile(true).then(path => {
                      driver.path = path
                    })
                  "
                >
                  {{ $t('driverForm.selectFile') }}
                </button>

                <input
                  type="text"
                  name="path"
                  v-model="driver.path"
                  class="input input-accent w-full join-item"
                  ref="pathInput"
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
                  <div tabindex="0" role="button" class="join-item btn m-1 w-30">
                    {{ $t('common.select') }}
                  </div>
                  <ul
                    tabindex="0"
                    class="dropdown-content menu bg-base-100 rounded-box z-1 w-52 p-2 shadow-sm"
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
                  type="text"
                  name="flags"
                  v-model="driver.flags"
                  class="input input-accent w-full join-item"
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
                  type="number"
                  name="minExeTime"
                  v-model="driver.minExeTime"
                  step="0.1"
                  class="input input-accent w-full"
                  required
                />

                <p class="label text-apple-green-800 text-wrap">
                  {{ $t('driverForm.minExecuteTimeHelp') }}
                </p>
              </fieldset>

              <fieldset class="fieldset flex-1">
                <legend class="fieldset-legend text-sm">
                  {{ $t('driverForm.allowedExitCode') }}
                </legend>

                <input
                  type="text"
                  name="allowRtCodes"
                  v-model="driver.allowRtCodes"
                  class="input input-accent"
                />

                <p class="label text-apple-green-800 text-wrap">
                  {{ $t('driverForm.commaSeparated') }}
                </p>
              </fieldset>
            </div>

            <fieldset class="fieldset flex-1">
              <legend class="fieldset-legend text-sm">
                {{ $t('driverForm.incompatibleWith') }}
              </legend>

              <div class="mb-1 text-xs line-clamp-1">
                <span class="inline">
                  {{ $t('driverForm.selectedWithCount', { count: driver.incompatibles?.length }) }}
                </span>
              </div>

              <div class="flex mb-2 gap-x-2">
                <input
                  v-model="searchPhrase"
                  :placeholder="$t('driverForm.search')"
                  class="input border-none focus:outline-gray-200 bg-gray-100 grow"
                />

                <button
                  type="button"
                  class="btn px-2 text-white"
                  style="--btn-color: var(--color-powder-blue-800)"
                  :title="$t('driverForm.selectAll')"
                  @click="
                    () => {
                      driver.incompatibles = [
                        ...groupStore.groups.flatMap(g => g.drivers.flatMap(d => d.id)),
                        'set_password',
                        'create_partition'
                      ]
                    }
                  "
                >
                  <font-awesome-icon icon="fa-regular fa-square-check" />
                </button>

                <button
                  type="button"
                  class="btn px-2 text-white"
                  style="--btn-color: var(--color-rose-400)"
                  :title="$t('driverForm.selectNone')"
                  @click="
                    () => {
                      driver.incompatibles = []
                    }
                  "
                >
                  <font-awesome-icon icon="fa-regular fa-square" />
                </button>
              </div>

              <ul class="h-44 p-1.5 overflow-auto border rounded-lg">
                <li
                  class="py-2.5 px-4 text-sm"
                  v-show="
                    searchPhrase === '' ||
                    'set password'.includes(searchPhrase) ||
                    $t('installOption.setPassword').includes(searchPhrase)
                  "
                >
                  <label class="flex items-center w-full select-none cursor-pointer">
                    <input
                      type="checkbox"
                      value="set_password"
                      v-model="driver.incompatibles"
                      class="checkbox checkbox-sm checkbox-primary me-1.5"
                    />
                    <span class="badge px-1 me-1" :style="`--badge-color: var(--color-builtin)`">
                      &nbsp;
                    </span>
                    <span class="line-clamp-2">
                      {{ $t('installOption.setPassword') }}
                    </span>
                  </label>
                </li>

                <li
                  class="py-2.5 px-4 text-sm"
                  v-show="
                    searchPhrase === '' ||
                    'create partition'.includes(searchPhrase) ||
                    $t('installOption.createPartition').includes(searchPhrase)
                  "
                >
                  <label class="flex items-center w-full select-none cursor-pointer">
                    <input
                      type="checkbox"
                      value="create_partition"
                      v-model="driver.incompatibles"
                      class="checkbox checkbox-sm checkbox-primary me-1.5"
                    />
                    <span class="badge px-1 me-1" :style="`--badge-color: var(--color-builtin)`">
                      &nbsp;
                    </span>
                    <span class="line-clamp-2">
                      {{ $t('installOption.createPartition') }}
                    </span>
                  </label>
                </li>

                <template v-for="g in filterGroups" :key="g.id">
                  <template v-for="d in g.drivers.filter(d => d.id != driver.id)" :key="d.id">
                    <li class="py-2.5 px-4 text-sm">
                      <label class="flex items-center w-full select-none cursor-pointer">
                        <input
                          type="checkbox"
                          :value="d.id"
                          v-model="driver.incompatibles"
                          class="checkbox checkbox-sm checkbox-primary me-1.5"
                        />
                        <span
                          class="badge px-1 me-1"
                          :class="[`badge-${g.type}`]"
                          :style="`--badge-color: var(--color-${g.type})`"
                        >
                          &nbsp;
                        </span>
                        <span class="line-clamp-2">
                          {{ `[${g.name}] ${d.name}` }}
                        </span>
                      </label>
                    </li>
                  </template>
                </template>
              </ul>
            </fieldset>

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
