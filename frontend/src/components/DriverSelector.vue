<script setup lang="ts">
import type { storage } from '@/wailsjs/go/models'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  driverGroups: Array<storage.DriverGroup>
  excludes?: Array<string>
  excludeBuiltin?: boolean
  groupBy: 'group' | 'driver'
}>()

const model = defineModel<Array<string> | undefined>({ default: [] })

const { t: $t } = useI18n()

const searchPhrase = ref('')

// ensure model is always an array
watch(
  () => model.value,
  val => {
    if (val === null || val === undefined) {
      model.value = []
    }
  },
  { immediate: true }
)

const filteredGroups = computed(() => {
  return searchPhrase.value === ''
    ? props.driverGroups
    : props.driverGroups.filter(
        g =>
          g.name.includes(searchPhrase.value) ||
          g.drivers.some(d => d.name.includes(searchPhrase.value))
      )
})

const builtinItems = computed(() => [
  { label: $t('installSetting.setPassword'), value: 'set_password', type: 'builtin' },
  { label: $t('installSetting.createPartition'), value: 'create_partition', type: 'builtin' }
])

const groupItems = computed(() =>
  props.groupBy === 'group'
    ? filteredGroups.value.map(g => ({ label: g.name, value: g.id, type: g.type }))
    : filteredGroups.value.flatMap(g =>
        g.drivers
          .filter(d => !props.excludes?.includes(d.id))
          .map(d => ({ label: `[${g.name}] ${d.name}`, value: d.id, groupType: g.type }))
      )
)
</script>

<template>
  <div>
    <div class="mb-1 line-clamp-1 text-xs">
      <span class="inline">
        {{ $t('driverForm.selectedWithCount', { count: model?.length ?? 0 }) }}
      </span>
    </div>

    <div class="mb-2 flex gap-x-2">
      <UInput
        v-model="searchPhrase"
        :placeholder="$t('driverForm.search')"
        class="ms-1 grow"
        variant="none"
        :ui="{ base: 'border-none bg-gray-100 focus:outline-gray-200' }"
      />

      <UButton
        type="button"
        class="px-2 text-white"
        style="--btn-color: var(--color-powder-blue-800)"
        :title="$t('driverForm.selectAll')"
        @click="
          () => {
            model = [
              ...($props.groupBy === 'group'
                ? props.driverGroups.map(g => g.id)
                : props.driverGroups.flatMap(g => g.drivers.flatMap(d => d.id))),
              'set_password',
              'create_partition'
            ]
          }
        "
      >
        <Icon icon="mdi:checkbox-marked" />
      </UButton>

      <UButton
        type="button"
        class="px-2 text-white"
        style="--btn-color: var(--color-rose-400)"
        :title="$t('driverForm.selectNone')"
        @click="model = []"
      >
        <Icon icon="mdi:checkbox-blank-outline" />
      </UButton>
    </div>

    <div class="overflow-auto rounded-lg border p-1.5">
      <div class="max-h-44 space-y-1.5">
        <!-- Builtin items -->
        <template v-if="!excludeBuiltin">
          <label
            v-for="item in builtinItems"
            :key="item.value"
            class="flex cursor-pointer items-center select-none"
          >
            <UCheckbox
              size="sm"
              :model-value="model?.includes(item.value) ?? false"
              @update:model-value="
                (checked: boolean | 'indeterminate') => {
                  if (checked === true) {
                    model = [...(model || []), item.value]
                  } else {
                    model = model?.filter(v => v !== item.value) ?? []
                  }
                }
              "
            />

            <UBadge
              v-if="item.type"
              size="sm"
              class="ms-1.5"
              :style="`background-color: var(--color-${item.type})`"
            >
              &nbsp;
            </UBadge>

            <span class="ms-1.5">{{ item.label }}</span>
          </label>
        </template>

        <!-- Group/Driver items with color blocks -->
        <template v-for="item in groupItems" :key="item.value">
          <label class="flex cursor-pointer items-center select-none">
            <UCheckbox
              size="sm"
              :model-value="model?.includes(item.value) ?? false"
              @update:model-value="
                (checked: boolean | 'indeterminate') => {
                  if (checked === true) {
                    model = [...(model || []), item.value]
                  } else {
                    model = model?.filter(v => v !== item.value) ?? []
                  }
                }
              "
            />

            <UBadge
              v-if="'type' in item ? item.type : 'groupType' in item ? item.groupType : undefined"
              size="sm"
              class="ms-1.5"
              :style="`background-color: var(--color-${'type' in item ? item.type : item.groupType})`"
            >
              &nbsp;
            </UBadge>

            <span class="ms-1.5">{{ item.label }}</span>
          </label>
        </template>
      </div>
    </div>
  </div>
</template>
