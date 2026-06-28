<script setup lang="ts">
import DriverInputModal from '@/components/DriverInputModal.vue'
import { ExecutableExists } from '@/wailsjs/go/main/App'
import { storage } from '@/wailsjs/go/models'
import * as groupStorage from '@/wailsjs/go/storage/DriverGroupStorage'
import { computed, ref, toRaw, useTemplateRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

const props = defineProps<{ id?: number }>()

const { t } = useI18n()

const $route = useRoute()

function driverCategoryKey(type: string): string {
  return `driverCategory${type.charAt(0).toUpperCase() + type.slice(1)}`
}
const $router = useRouter()
const toast = useToast()

const inputModal = useTemplateRef('inputModal')

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
    toast.add({ title: t('toastAddAtLeastOneDriver'), color: 'warning' })
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
    const updateGroup = toRaw(group.value)
    groupStorage
      .Update(updateGroup)
      .then(handleSuccess)
      .catch(reason => toast.add({ title: reason.toString(), color: 'error' }))
  }
}
</script>

<template>
  <form
    class="mx-auto flex h-full max-w-full flex-col justify-center gap-y-8 overflow-y-auto lg:max-w-2xl xl:max-w-4xl"
    autocomplete="off"
    @submit.prevent="handleSubmit"
  >
    <div class="flex gap-x-3 px-1">
      <div class="w-32">
        <fieldset class="fieldset">
          <legend class="fieldset-legend text-sm">{{ $t('driverFormType') }}</legend>

          <USelect
            v-model="group.type"
            name="type"
            class="w-full"
            :items="
              Object.values(storage.DriverType).map(type => ({
                label: $t(driverCategoryKey(type)),
                value: type
              }))
            "
            required
          />
        </fieldset>
      </div>

      <div class="grow">
        <fieldset class="fieldset">
          <legend class="fieldset-legend text-sm">{{ $t('driverFormName') }}</legend>

          <UInput v-model="group.name" type="text" class="w-full" required />
        </fieldset>
      </div>
    </div>

    <div>
      <label class="flex w-full cursor-pointer items-center select-none">
        <UCheckbox v-model="group.mutuallyExclusive" class="me-1.5" />
        {{ $t('driverFormMutuallyExclusive') }}
      </label>

      <p class="text-hint text-xs">{{ $t('driverFormMutuallyExclusiveHelp') }}</p>
    </div>

    <fieldset class="fieldset">
      <legend class="fieldset-legend text-sm">{{ $t('driverFormDriver') }}</legend>

      <div>
        <div class="max-h-[40vh] overflow-y-auto">
          <div class="grid-rows grid text-sm">
            <div class="grid grid-cols-10 gap-2 border-y py-1.5">
              <div class="col-span-2">{{ $t('driverFormName') }}</div>

              <div class="col-span-3">{{ $t('driverFormPath') }}</div>

              <div class="col-span-2">{{ $t('driverFormArgument') }}</div>

              <div class="col-span-2">{{ $t('driverFormOtherSetting') }}</div>
            </div>

            <div v-if="group.drivers.length == 0" class="py-1 text-center last:border-b">N/A</div>

            <div
              v-for="(d, i) in group.drivers"
              v-else
              :key="d.id"
              class="grid grid-cols-10 items-center gap-2 border-b py-1.5 text-xs"
              :class="{ 'bg-lime-50': d.id === 0 }"
            >
              <div class="col-span-2">
                <p class="line-clamp-2 break-all">
                  {{ d.name }}
                </p>
              </div>

              <div class="col-span-3">
                <p
                  class="line-clamp-2 font-mono break-all"
                  :class="{ 'text-red-600': notFoundDrivers.includes(d.id) }"
                >
                  {{ d.path }}
                </p>
              </div>

              <div class="col-span-2">
                <p class="line-clamp-2 break-all">
                  {{ d.flags.join(', ') }}
                </p>
              </div>

              <div class="col-span-2 flex gap-x-1">
                <span
                  v-show="d.incompatibles.length > 0"
                  class="inline-block max-h-5 rounded-xs bg-yellow-300 p-0.5"
                  :title="$t('driverFormIncompatibleWith')"
                >
                  <Icon icon="mdi:source-merge" />
                </span>

                <span
                  v-show="d.allowRtCodes.length > 0"
                  class="inline-block max-h-5 rounded-xs bg-blue-300 p-0.5"
                  :title="$t('driverFormAllowedExitCode')"
                >
                  <Icon icon="mdi:numeric-0-box-outline" />
                </span>
              </div>

              <div>
                <div class="flex gap-x-2">
                  <button type="button" :title="$t('commonEdit')" @click="inputModal?.show(d)">
                    <Icon icon="mdi:pencil" class="size-4" />
                  </button>

                  <button
                    type="button"
                    :title="$t('commonDelete')"
                    @click="group.drivers.splice(i, 1)"
                  >
                    <Icon icon="mdi:trash-can" class="size-4" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <p class="text-hint">
            {{ $t('driverFormDriverGroupHelp') }}
          </p>
        </div>

        <div class="flex justify-end gap-x-3">
          <UButton type="button" class="px-2" color="primary" @click="inputModal?.show()">
            <Icon icon="mdi:plus-box" />
          </UButton>
        </div>
      </div>
    </fieldset>

    <div class="flex h-8 gap-x-5">
      <UButton
        type="button"
        class="grow justify-center"
        color="neutral"
        variant="outline"
        style="--btn-color: var(--color-gray-100)"
        @click="$router.back()"
      >
        {{ $t('commonBack') }}
      </UButton>

      <UButton type="submit" class="grow justify-center" color="secondary">
        {{ $t('commonSave') }}
      </UButton>
    </div>
  </form>

  <DriverInputModal
    ref="inputModal"
    @submit="
      newDriver => {
        console.log(newDriver)
        if (newDriver.id) {
          group.drivers = group.drivers.map(d => (d.id == newDriver.id ? newDriver : d))
        } else {
          group.drivers.push({
            ...newDriver,
            id: 0
          })

          console.table(toRaw(group))
        }
        inputModal?.hide()
      }
    "
  ></DriverInputModal>
</template>

<style scoped>
legend:has(+ input:required, + select:required):after,
legend:has(+ div > input:required):after {
  content: ' *';
  color: red;
}
</style>
