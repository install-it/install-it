<script setup lang="ts">
import DriverInputModal from '@/components/DriverInputModal.vue'
import { ExecutableExists } from '@/wailsjs/go/main/App'
import { storage } from '@/wailsjs/go/models'
import * as groupStorage from '@/wailsjs/go/storage/DriverGroupStorage'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

const props = defineProps<{ id?: number }>()

const { t } = useI18n()

const $route = useRoute()

function categoryKey(type: string): string {
  return `category${type.charAt(0).toUpperCase() + type.slice(1)}`
}
const $router = useRouter()
const toast = useToast()

const isModalOpen = ref(false)
const editingDriver = ref<Partial<storage.Driver> | undefined>(undefined)

const groupStore = useDriverGroupStore()

// Create computed source for dynamic group lookup
const sourceGroup = computed(
  () =>
    groupStore.groups.find(g => g.id === props.id) ??
    new storage.DriverGroup({
      type:
        storage.DriverType[
          (
            $route.query.type as string | undefined
          )?.toUpperCase() as keyof typeof storage.DriverType
        ] ?? storage.DriverType.NETWORK,
      name: '',
      drivers: []
    })
)

const { data: group, reset } = useEditor({
  source: sourceGroup,
  warnOnUnsavedLeave: true
})

// Track drivers that don't exist on system
const notFoundDrivers = ref<number[]>([])

const findNotExists = (drivers: Array<storage.Driver>) =>
  Promise.all(
    drivers.map(d => ExecutableExists(d.path).then(exist => ({ id: d.id, exist: exist })))
  ).then(results => {
    return results
      .map(result => (result.exist ? undefined : result.id))
      .filter(v => v !== undefined)
  })

watch(
  () => group.value.drivers,
  newDrivers => findNotExists(newDrivers).then(ids => (notFoundDrivers.value = ids)),
  { immediate: true }
)

function handleSubmit() {
  if (group.value.drivers.length == 0) {
    toast.add({ title: t('toastAddDriverRequired'), color: 'warning' })
    return
  }

  const handleSuccess = () => {
    toast.add({ title: t('toastUpdated'), color: 'success' })
    groupStorage
      .All()
      .then(newDriverGroups => {
        groupStore.groups = newDriverGroups
        return reset()
      })
      .then(() => {
        $router.back()
      })
  }

  if (!group.value.id) {
    groupStorage
      .Add(group.value)
      .then(handleSuccess)
      .catch(reason => toast.add({ title: reason.toString(), color: 'error' }))
  } else {
    const updateGroup = group.value
    groupStorage
      .Update(updateGroup)
      .then(handleSuccess)
      .catch(reason => toast.add({ title: reason.toString(), color: 'error' }))
  }
}
</script>

<template>
  <form
    class="mx-auto flex h-full max-w-full flex-col gap-y-5 overflow-y-auto lg:max-w-2xl xl:max-w-4xl"
    autocomplete="off"
    @submit.prevent="handleSubmit"
  >
    <!-- Basic Info Row -->
    <div class="flex gap-4">
      <div class="w-48 shrink-0">
        <label class="mb-1 block text-xs font-bold tracking-wider text-gray-500 uppercase">
          {{ $t('fieldDriverType') }} <span class="text-red-500">*</span>
        </label>

        <USelect
          v-model="group.type"
          name="type"
          class="w-full"
          :items="
            Object.values(storage.DriverType).map(type => ({
              label: $t(categoryKey(type)),
              value: type
            }))
          "
          required
        />
      </div>

      <div class="flex-1">
        <label class="mb-1 block text-xs font-bold tracking-wider text-gray-500 uppercase">
          {{ $t('name') }} <span class="text-red-500">*</span>
        </label>

        <UInput v-model="group.name" type="text" class="w-full" required />
      </div>
    </div>

    <!-- Mutually Exclusive Card -->
    <div class="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
      <label class="flex cursor-pointer items-start gap-3">
        <UCheckbox v-model="group.mutuallyExclusive" class="mt-1" />

        <div>
          <span class="block text-sm font-bold text-gray-800">{{
            $t('fieldMutuallyExclusive')
          }}</span>

          <p class="mt-0.5 text-xs text-gray-500">{{ $t('descMutuallyExclusive') }}</p>
        </div>
      </label>
    </div>

    <!-- Drivers Stack -->
    <div class="flex flex-col">
      <label class="mb-2 block text-xs font-bold tracking-wider text-gray-500 uppercase">
        {{ $t('fieldDriver') }}
      </label>

      <div
        class="flex flex-col overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm"
      >
        <!-- List area -->
        <div class="max-h-64 divide-y divide-gray-100 overflow-y-auto">
          <!-- Empty state -->
          <div
            v-if="group.drivers.length === 0"
            class="flex flex-col items-center py-10 text-center text-gray-400"
          >
            <Icon icon="mdi:package-variant" class="mb-2 text-4xl opacity-50" />

            <span class="text-sm font-medium">{{ $t('msgNoDriversInGroup') }}</span>
          </div>

          <!-- Driver rows -->
          <div
            v-for="(d, i) in group.drivers"
            v-else
            :key="d.id"
            class="group/row flex items-start gap-4 p-4 transition-colors hover:bg-gray-50"
            :class="{ 'bg-lime-50': d.id <= 0 }"
          >
            <span
              class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded bg-gray-100 text-sm font-bold text-gray-500"
            >
              {{ i + 1 }}
            </span>

            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span class="text-base font-bold text-gray-800">{{ d.name }}</span>

                <span
                  v-if="notFoundDrivers.includes(d.id)"
                  class="rounded border border-red-200 bg-red-100 px-1.5 py-0.5 text-xs font-bold text-red-600"
                >
                  File Not Found
                </span>
              </div>

              <p class="mt-1 font-mono text-sm break-all text-gray-500">
                {{ $t('labelPathPrefix') }}: {{ d.path }}
              </p>

              <div class="mt-2 flex flex-wrap gap-2">
                <span
                  v-if="d.flags.length"
                  class="rounded bg-gray-100 px-1.5 py-0.5 font-mono text-xs text-gray-600"
                >
                  {{ $t('labelFlagsPrefix') }}: {{ d.flags.join(', ') }}
                </span>

                <span
                  v-if="d.allowRtCodes?.length"
                  class="flex items-center gap-1 rounded bg-purple-50 px-1.5 py-0.5 text-xs font-semibold text-purple-600"
                >
                  <Icon icon="mdi:shield-check-outline" /> {{ $t('fieldAllowedExitCode') }}:
                  {{ d.allowRtCodes.join(', ') }}
                </span>

                <span
                  v-if="d.incompatibles.length"
                  class="flex items-center gap-1 rounded bg-yellow-100 px-1.5 py-0.5 text-xs font-semibold text-yellow-700"
                >
                  <Icon icon="mdi:source-merge" /> {{ $t('labelIncompatibleWith') }}:
                  {{ d.incompatibles.length }}
                </span>

                <span
                  class="flex items-center gap-1 rounded bg-blue-50 px-1.5 py-0.5 text-xs font-semibold text-blue-600"
                >
                  <Icon icon="mdi:timer-outline" /> {{ $t('labelMinPrefix') }}: {{ d.minExeTime }}s
                </span>
              </div>
            </div>

            <div class="flex gap-1 opacity-0 transition-opacity group-hover/row:opacity-100">
              <button
                type="button"
                :title="$t('edit')"
                class="flex h-8 w-8 items-center justify-center rounded text-gray-600 hover:bg-gray-200"
                @click="
                  editingDriver = d
                  isModalOpen = true
                "
              >
                <Icon icon="mdi:pencil" />
              </button>

              <button
                type="button"
                :title="$t('delete')"
                class="flex h-8 w-8 items-center justify-center rounded text-red-500 hover:bg-red-100"
                @click="group.drivers.splice(i, 1)"
              >
                <Icon icon="mdi:trash-can" />
              </button>
            </div>
          </div>
        </div>

        <!-- Footer bar -->
        <div class="flex items-center justify-between border-t border-gray-200 bg-gray-50 p-4">
          <span class="text-sm text-gray-500 italic">{{ $t('descDriverGroup') }}</span>

          <UButton
            type="button"
            color="primary"
            size="sm"
            @click="
              editingDriver = undefined
              isModalOpen = true
            "
          >
            <Icon icon="mdi:plus-circle-outline" class="mr-1 text-sm" />
            {{ $t('labelAddDriver') }}
          </UButton>
        </div>
      </div>
    </div>

    <!-- Master Form Actions -->
    <div class="mt-4 flex shrink-0 gap-4 border-t border-gray-200 pt-4">
      <UButton
        type="button"
        color="neutral"
        variant="outline"
        class="flex-1 justify-center text-sm"
        @click="$router.back()"
      >
        {{ $t('back') }}
      </UButton>

      <UButton type="submit" color="secondary" class="flex-1 justify-center text-sm">
        {{ $t('save') }}
      </UButton>
    </div>
  </form>

  <DriverInputModal
    v-model:open="isModalOpen"
    :edit-data="editingDriver"
    @submit="
      newDriver => {
        if (newDriver.id) {
          group.drivers = group.drivers.map(d => (d.id === newDriver.id ? newDriver : d))
        } else {
          group.drivers.push({ ...newDriver, id: -Date.now() })
        }
      }
    "
  />
</template>
