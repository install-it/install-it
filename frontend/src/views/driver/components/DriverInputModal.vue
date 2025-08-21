<script setup lang="ts">
import ModalFrame from '@/components/modals/ModalFrame.vue'
import { flags } from '@/definitions/flags'
import { SelectFile } from '@/wailsjs/go/main/App'
import { store } from '@/wailsjs/go/models'
import * as groupStore from '@/wailsjs/go/store/DriverGroupStore'
import { computed, nextTick, ref, useTemplateRef } from 'vue'

const frame = useTemplateRef('frame')

defineExpose({
  show: (data?: Partial<store.Driver>) => {
    frame.value?.show()

    groupStore.All().then(g => (groups.value = g))

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

defineEmits<{ submit: [dri: store.Driver] }>()

const groups = ref<Array<store.DriverGroup>>([])

const modalBody = useTemplateRef<HTMLDivElement>('modalBody')

const searchPhrase = ref('')

const driver = ref<
  Partial<Omit<store.Driver, 'allowRtCodes' | 'flags'> & { allowRtCodes: string; flags: string }>
>({})

const filterGroups = computed(() => {
  return searchPhrase.value === ''
    ? groups.value
    : groups.value.filter(
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
            class="flex flex-col gap-y-3"
            autocomplete="off"
            @submit.prevent="
              _ => {
                $emit(
                  'submit',
                  new store.Driver({
                    ...driver,
                    flags: driver.flags ? driver.flags.split(',') : [],
                    allowRtCodes: driver.allowRtCodes
                      ? driver.allowRtCodes
                          ?.split(',')
                          .map(c => parseInt(c))
                          .filter(c => !Number.isNaN(c))
                      : [],
                    incompatibles: driver.incompatibles ?? []
                  })
                )

                frame?.hide()
              }
            "
          >
            <div>
              <label class="block mb-2 text-sm font-medium text-gray-900">
                {{ $t('driverForm.name') }}
              </label>
              <input
                type="text"
                name="name"
                v-model="driver.name"
                class="w-full p-1.5 text-sm shadow-xs"
              />
            </div>

            <div>
              <label class="block mb-2 text-sm font-medium text-gray-900">
                {{ $t('driverForm.path') }}
              </label>

              <div class="flex">
                <button
                  type="button"
                  class="w-28 px-3 text-sm text-gray-900 bg-gray-200 border border-e-0 rounded-s-md rounded-e-0 border-gray-300"
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
                  class="block flex-1 min-w-0 w-full p-1.5 text-sm rounded-none rounded-e shadow-xs"
                  ref="pathInput"
                  required
                />
              </div>
            </div>

            <div>
              <label class="block mb-2 text-sm font-medium text-gray-900">
                {{ $t('driverForm.argument') }}
              </label>

              <div class="flex">
                <select
                  name="flags"
                  class="w-28 text-sm border border-e-0 rounded-e-none rounded-s"
                  @change="
                    event => {
                      driver.flags = (event.target as HTMLSelectElement).value
                    }
                  "
                >
                  <option value="">
                    {{ $t('driverForm.manualInput') }}
                  </option>
                  <option
                    v-for="(flag, name) in flags"
                    :key="name"
                    :value="flag.join(',')"
                    :selected="driver.flags === flag.join()"
                  >
                    {{ name }}
                  </option>
                </select>
                <input
                  type="text"
                  name="flags"
                  v-model="driver.flags"
                  class="flex-1 p-1.5 text-sm rounded-none rounded-e shadow-xs"
                />
              </div>

              <p class="text-hint">
                {{ $t('driverForm.commaSeparated') }}
              </p>
            </div>

            <div class="flex gap-x-3">
              <div>
                <label class="block mb-2 text-sm font-medium text-gray-900">
                  {{ $t('driverForm.minExecuteTime') }}
                </label>
                <input
                  type="number"
                  name="minExeTime"
                  v-model="driver.minExeTime"
                  step="0.1"
                  class="w-full p-1.5 text-sm shadow-xs"
                  required
                />
                <p class="text-hint">
                  {{ $t('driverForm.minExecuteTimeHelp') }}
                </p>
              </div>

              <div>
                <label class="block mb-2 text-sm font-medium text-gray-900">
                  {{ $t('driverForm.allowedExitCode') }}
                </label>
                <input
                  type="text"
                  name="allowRtCodes"
                  v-model="driver.allowRtCodes"
                  class="w-full p-1.5 text-sm shadow-xs"
                />
                <p class="text-hint">
                  {{ $t('driverForm.allowedExitCodeHelp') }}
                </p>
                <p class="text-hint">
                  {{ $t('driverForm.commaSeparated') }}
                </p>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-900">
                {{ $t('driverForm.incompatibleWith') }}
              </label>

              <div class="mb-1 text-xs line-clamp-1">
                <span class="inline">
                  {{ $t('driverForm.selectedWithCount', { count: driver.incompatibles?.length }) }}
                </span>
              </div>

              <div class="flex mb-2 gap-x-2">
                <input
                  v-model="searchPhrase"
                  :placeholder="$t('driverForm.search')"
                  class="px-3 py-2 w-full text-black text-sm border-none rounded-sm bg-gray-100"
                />

                <button
                  type="button"
                  class="px-3 text-sm font-medium text-white bg-powder-blue-800 hover:bg-powder-blue-600 rounded-sm"
                  :title="$t('driverForm.selectAll')"
                  @click="
                    () => {
                      driver.incompatibles = [
                        ...groups.flatMap(g => g.drivers.flatMap(d => d.id)),
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
                  class="px-3 text-sm font-medium text-white bg-rose-400 hover:bg-rose-300 rounded-sm"
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
                  <label class="flex item-center w-full select-none cursor-pointer">
                    <input
                      type="checkbox"
                      value="set_password"
                      v-model="driver.incompatibles"
                      class="me-1.5"
                    />
                    <span class="badge badge-builtin me-1">&nbsp;</span>
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
                  <label class="flex item-center w-full select-none cursor-pointer">
                    <input
                      type="checkbox"
                      value="create_partition"
                      v-model="driver.incompatibles"
                      class="me-1.5"
                    />
                    <span class="badge badge-builtin me-1">&nbsp;</span>
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
                          class="me-1.5"
                        />
                        <span class="badge me-1" :class="[`badge-${g.type}`]">&nbsp;</span>
                        <span class="line-clamp-2">
                          {{ `[${g.name}] ${d.name}` }}
                        </span>
                      </label>
                    </li>
                  </template>
                </template>
              </ul>
            </div>

            <button
              type="submit"
              class="w-full my-1 py-2 text-sm font-medium text-white bg-half-baked-600 hover:bg-half-baked-500 rounded-lg"
            >
              {{ $t('common.save') }}
            </button>
          </form>
        </div>
      </div>
    </div>
  </ModalFrame>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.5s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

label:has(+ input:required, + select:required):after,
label:has(+ div > input:required):after {
  content: ' *';
  color: red;
}
</style>
