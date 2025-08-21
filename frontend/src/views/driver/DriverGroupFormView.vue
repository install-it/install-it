<script setup lang="ts">
import UnsaveConfirmModal from '@/components/modals/UnsaveConfirmModal.vue'
import { getNotExistDrivers } from '@/utils/index'
import { store } from '@/wailsjs/go/models'
import * as groupStore from '@/wailsjs/go/store/DriverGroupStore'
import { onBeforeMount, ref, toRaw, useTemplateRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'
import { useToast } from 'vue-toast-notification'
import DriverInputModal from './components/DriverInputModal.vue'

const { t } = useI18n()

const $route = useRoute()

const $router = useRouter()

const $toast = useToast({ position: 'top-right' })

const questionModal = useTemplateRef('questionModal')

const inputModal = useTemplateRef('inputModal')

const notExistDrivers = ref<Array<string>>([])

const group = ref<store.DriverGroup>(
  new store.DriverGroup({
    type:
      store.DriverType[
        $route.query.type?.toString().toUpperCase() as keyof typeof store.DriverType
      ] ?? undefined,
    name: $route.query.name ?? '',
    drivers: []
  })
)

/** A clone of the `group` variable */
let groupOriginal: store.DriverGroup = structuredClone(toRaw(group.value))

onBeforeMount(() => {
  groupStore
    .Get($route.params.id as string)
    .then(g => {
      group.value = g
      groupOriginal = structuredClone(toRaw(g))

      getNotExistDrivers(g.drivers).then(result => {
        notExistDrivers.value = result
      })
    })
    .catch(() => undefined)
    .finally(() => {
      // let the async call fails when $route.params.id is undefinded to avoid duplicate watcher setup logic
      //
      // setup watchers after the async call to avoid the triggering due to group.value replacement
      watch(
        () => group.value.drivers,
        newDrivers => {
          getNotExistDrivers(newDrivers).then(result => {
            notExistDrivers.value = result
          })
        },
        { deep: true }
      )
    })
})

onBeforeRouteLeave((to, from, next) => {
  if (JSON.stringify(groupOriginal) != JSON.stringify(group.value)) {
    questionModal.value?.show(answer => {
      next(answer == 'yes')
    })
  } else {
    next(true)
  }
})

function handleSubmit(event: SubmitEvent) {
  if (group.value.drivers.length == 0) {
    $toast.warning(t('toast.addAtLeastOneDriver'))
    return
  }

  const action =
    group.value.id == undefined
      ? groupStore.Add(group.value).then(gid => {
          group.value.id = gid
          // no need to update the URL,
          // as users is not able to refresh the page in production build
        })
      : groupStore.Update({
          ...group.value,
          drivers: group.value.drivers.map(d => {
            if (d.id.includes('new:')) {
              d.id = ''
            }
            return d
          })
        })

  action
    .then(() => {
      $toast.success(t('toast.updated'))

      groupStore.Get(group.value.id).then(g => {
        group.value = g
        groupOriginal = structuredClone(g)

        if (event.submitter?.id != 'driver-submit-btn') {
          $router.back()
        } else {
          getNotExistDrivers(g.drivers).then(result => {
            notExistDrivers.value = result
          })
        }
      })
    })
    .catch(reason => {
      $toast.error(reason)
    })
}
</script>

<template>
  <form
    class="flex flex-col justify-center h-full max-w-full lg:max-w-2xl xl:max-w-4xl mx-auto gap-y-8 overflow-y-auto"
    autocomplete="off"
    @submit.prevent="event => handleSubmit(event as SubmitEvent)"
  >
    <div class="flex gap-x-3">
      <div class="w-32">
        <label class="block mb-2 text-sm font-medium text-gray-900">
          {{ $t('driverForm.type') }}
        </label>
        <select name="type" v-model="group.type" class="w-full p-1.5 text-sm shadow-xs" required>
          <option v-for="type in store.DriverType" :key="type" :value="type">
            {{ $t(`driverCatetory.${type}`) }}
          </option>
        </select>
      </div>

      <div class="flex-1">
        <label class="block mb-2 text-sm font-medium text-gray-900">
          {{ $t('driverForm.name') }}
        </label>
        <input
          type="text"
          name="name"
          v-model="group.name"
          class="w-full p-1.5 text-sm shadow-xs"
          required
        />
      </div>
    </div>

    <div class="flex flex-col gap-y-4">
      <label class="block text-sm font-medium text-gray-900">
        {{ $t('driverForm.driver') }}
      </label>

      <div class="max-h-[40vh] text-sm overflow-y-auto">
        <div class="grid grid-rows">
          <div class="grid grid-cols-10 gap-2 py-1.5 border-y">
            <div class="col-span-2">{{ $t('driverForm.name') }}</div>
            <div class="col-span-3">{{ $t('driverForm.path') }}</div>
            <div class="col-span-2">{{ $t('driverForm.argument') }}</div>
            <div class="col-span-2">{{ $t('driverForm.otherSetting') }}</div>
          </div>

          <div v-if="group.drivers.length == 0" class="py-1 text-center last:border-b">N/A</div>

          <div
            v-else
            v-for="(d, i) in group.drivers"
            :key="d.id"
            class="grid grid-cols-10 items-center gap-2 py-1.5 text-xs border-b"
            :class="{ 'bg-lime-50': d.id.includes('new:') }"
          >
            <div class="col-span-2">
              <p class="break-all line-clamp-2">
                {{ d.name }}
              </p>
            </div>

            <div class="col-span-3">
              <p
                class="font-mono break-all line-clamp-2"
                :class="{ 'text-red-600': notExistDrivers.includes(d.id) }"
              >
                {{ d.path }}
              </p>
            </div>

            <div class="col-span-2">
              <p class="break-all line-clamp-2">
                {{ d.flags.join(', ') }}
              </p>
            </div>

            <div class="flex col-span-2 gap-x-1">
              <span
                v-show="d.incompatibles.length > 0"
                class="inline-block p-0.5 max-h-5 bg-yellow-300 rounded-xs"
                :title="$t('driverForm.incompatibleWith')"
              >
                <font-awesome-icon icon="fa-solid fa-code-merge" />
              </span>

              <span
                v-show="d.allowRtCodes.length > 0"
                class="inline-block p-0.5 max-h-5 bg-blue-300 rounded-xs"
                :title="$t('driverForm.allowedExitCode')"
              >
                <font-awesome-icon icon="fa-solid fa-0" />
              </span>
            </div>

            <div>
              <div class="flex gap-x-2">
                <button type="button" @click="inputModal?.show(d)">
                  <font-awesome-icon icon="fa-solid fa-pen-to-square" />
                </button>
                <button type="button" @click="group.drivers.splice(i, 1)">
                  <font-awesome-icon icon="fa-solid fa-trash" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <p class="text-hint">
          {{ $t('driverForm.incompatibleForNewHelp') }}
        </p>
      </div>

      <div class="flex justify-end gap-x-3">
        <button
          v-show="JSON.stringify(group.drivers) != JSON.stringify(groupOriginal.drivers)"
          type="submit"
          id="driver-submit-btn"
          class="h-8 px-2 text-sm font-medium text-white bg-half-baked-600 hover:bg-half-baked-500 rounded-lg"
        >
          <font-awesome-icon icon="fa-solid fa-floppy-disk" />
        </button>

        <button
          type="button"
          class="h-8 px-2 text-sm font-medium text-white bg-powder-blue-800 hover:bg-powder-blue-600 rounded-lg"
          @click="inputModal?.show()"
        >
          <font-awesome-icon icon="fa-regular fa-square-plus" />
        </button>
      </div>
    </div>

    <div class="flex h-8 gap-x-5">
      <button
        type="button"
        class="w-full text-sm font-medium text-white bg-gray-400 hover:bg-gray-300 rounded-lg"
        @click="$router.back()"
      >
        {{ $t('common.back') }}
      </button>

      <button
        type="submit"
        class="w-full text-sm font-medium text-white bg-half-baked-600 hover:bg-half-baked-500 rounded-lg"
      >
        {{ $t('common.save') }}
      </button>
    </div>
  </form>

  <DriverInputModal
    @submit="
      newDriver => {
        if (newDriver.id) {
          group.drivers = group.drivers.map(d => (d.id == newDriver.id ? newDriver : d))
        } else {
          group.drivers.push({
            ...newDriver,
            id: `new:${group.drivers.length + 1}` // assign a temporary ID for editing
          })
        }
        inputModal?.hide()
      }
    "
    ref="inputModal"
  ></DriverInputModal>

  <UnsaveConfirmModal ref="questionModal"></UnsaveConfirmModal>
</template>

<style scoped>
label:has(+ input:required, + select:required):after,
label:has(+ div > input:required):after {
  content: ' *';
  color: red;
}
</style>
