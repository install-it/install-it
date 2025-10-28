<script setup lang="ts">
import type { storage } from '@/wailsjs/go/models'
import { computed, ref } from 'vue'

const props = defineProps<{
  driverGroups: Array<storage.DriverGroup>
  excludes?: Array<string>
  excludeBuiltin?: boolean
  groupBy: 'group' | 'driver'
}>()

const model = defineModel<Array<string>>({ default: [] })

const searchPhrase = ref('')

const filteredGroups = computed(() => {
  return searchPhrase.value === ''
    ? props.driverGroups
    : props.driverGroups.filter(
        g =>
          g.name.includes(searchPhrase.value) ||
          g.drivers.some(d => d.name.includes(searchPhrase.value))
      )
})
</script>

<template>
  <div>
    <div class="mb-1 text-xs line-clamp-1">
      <span class="inline">
        {{ $t('driverForm.selectedWithCount', { count: model.length }) }}
      </span>
    </div>

    <div class="flex mb-2 gap-x-2">
      <input
        v-model="searchPhrase"
        :placeholder="$t('driverForm.search')"
        class="input border-none focus:outline-gray-200 bg-gray-100 grow ms-1"
      />

      <button
        type="button"
        class="btn px-2 text-white"
        style="--btn-color: var(--color-powder-blue-800)"
        :title="$t('driverForm.selectAll')"
        @click="
          () => {
            model = [
              ...($props.groupBy === 'group'
                ? driverGroups.map(g => g.id)
                : driverGroups.flatMap(g => g.drivers.flatMap(d => d.id))),
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
        @click="model = []"
      >
        <font-awesome-icon icon="fa-regular fa-square" />
      </button>
    </div>

    <ul class="h-44 p-1.5 overflow-auto border rounded-lg">
      <li
        v-if="!excludeBuiltin"
        class="py-2 px-4 text-sm"
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
            v-model="model"
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
        v-if="!excludeBuiltin"
        class="py-2 px-4 text-sm"
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
            v-model="model"
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

      <template v-for="g in filteredGroups" :key="g.id">
        <template v-if="$props.groupBy === 'group'">
          <li class="py-2 px-4 text-sm">
            <label class="flex items-center w-full select-none cursor-pointer">
              <input
                type="checkbox"
                :value="g.id"
                v-model="model"
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
                {{ g.name }}
              </span>
            </label>
          </li>
        </template>

        <template v-else>
          <template v-for="d in g.drivers.filter(d => !excludes?.includes(d.id))" :key="d.id">
            <li class="py-2 px-4 text-sm">
              <label class="flex items-center w-full select-none cursor-pointer">
                <input
                  type="checkbox"
                  :value="d.id"
                  v-model="model"
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
      </template>
    </ul>
  </div>
</template>
