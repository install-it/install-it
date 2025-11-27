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
    <div class="mb-1 line-clamp-1 text-xs">
      <span class="inline">
        {{ $t('driverForm.selectedWithCount', { count: model.length }) }}
      </span>
    </div>

    <div class="mb-2 flex gap-x-2">
      <input
        v-model="searchPhrase"
        :placeholder="$t('driverForm.search')"
        class="input ms-1 grow border-none bg-gray-100 focus:outline-gray-200"
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

    <ul class="h-44 overflow-auto rounded-lg border p-1.5">
      <li
        v-if="!excludeBuiltin"
        v-show="
          searchPhrase === '' ||
          'set password'.includes(searchPhrase) ||
          $t('installSetting.setPassword').includes(searchPhrase)
        "
        class="px-4 py-2 text-sm"
      >
        <label class="flex w-full cursor-pointer items-center select-none">
          <input
            v-model="model"
            type="checkbox"
            value="set_password"
            class="checkbox me-1.5 checkbox-sm checkbox-primary"
          />
          <span class="me-1 badge px-1" :style="`--badge-color: var(--color-builtin)`">
            &nbsp;
          </span>
          <span class="line-clamp-2">
            {{ $t('installSetting.setPassword') }}
          </span>
        </label>
      </li>

      <li
        v-if="!excludeBuiltin"
        v-show="
          searchPhrase === '' ||
          'create partition'.includes(searchPhrase) ||
          $t('installSetting.createPartition').includes(searchPhrase)
        "
        class="px-4 py-2 text-sm"
      >
        <label class="flex w-full cursor-pointer items-center select-none">
          <input
            v-model="model"
            type="checkbox"
            value="create_partition"
            class="checkbox me-1.5 checkbox-sm checkbox-primary"
          />
          <span class="me-1 badge px-1" :style="`--badge-color: var(--color-builtin)`">
            &nbsp;
          </span>
          <span class="line-clamp-2">
            {{ $t('installSetting.createPartition') }}
          </span>
        </label>
      </li>

      <template v-for="g in filteredGroups" :key="g.id">
        <template v-if="$props.groupBy === 'group'">
          <li class="px-4 py-2 text-sm">
            <label class="flex w-full cursor-pointer items-center select-none">
              <input
                v-model="model"
                type="checkbox"
                :value="g.id"
                class="checkbox me-1.5 checkbox-sm checkbox-primary"
              />
              <span
                class="me-1 badge px-1"
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
            <li class="px-4 py-2 text-sm">
              <label class="flex w-full cursor-pointer items-center select-none">
                <input
                  v-model="model"
                  type="checkbox"
                  :value="d.id"
                  class="checkbox me-1.5 checkbox-sm checkbox-primary"
                />
                <span
                  class="me-1 badge px-1"
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
