<script setup lang="ts">
import UnsaveConfirmModal from '@/components/modals/UnsaveConfirmModal.vue'
import { useDriverGroupStore } from '@/store'
import { storage } from '@/wailsjs/go/models'
import * as groupStorage from '@/wailsjs/go/storage/DriverGroupStorage'
import { toRaw, useTemplateRef } from 'vue'
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

const groupStore = useDriverGroupStore()

const groupEditor = groupStore.editor(
  $route.params.id as string | undefined,
  storage.DriverType[$route.query.type?.toString().toUpperCase() as keyof typeof storage.DriverType]
)

const group = groupEditor.group // alias

onBeforeRouteLeave((to, from, next) => {
  if (groupEditor.modified.value) {
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

  const handleSuccess = () => {
    $toast.success(t('toast.updated'))
    groupStorage.All().then(newDriverGroups => {
      groupStore.groups = newDriverGroups
      groupEditor.reset()
      if (event.submitter?.id !== 'driver-submit-btn') {
        $router.back()
      } else {
        $router.replace({ path: `/drivers/edit/${group.value.id}` })
      }
    })
  }

  if (group.value.id == undefined) {
    groupStorage
      .Add(group.value)
      .then(gid => (group.value.id = gid))
      .then(handleSuccess)
      .catch(reason => $toast.error(reason.toString()))
  } else {
    groupStorage
      .Update({
        ...group.value,
        drivers: group.value.drivers.map(d => {
          if (d.id.includes('new:')) {
            d.id = ''
          }
          return d
        })
      })
      .then(handleSuccess)
      .catch(reason => $toast.error(reason.toString()))
  }
}
</script>

<template>
  <form
    class="flex flex-col justify-center h-full max-w-full lg:max-w-2xl xl:max-w-4xl mx-auto gap-y-8 overflow-y-auto"
    autocomplete="off"
    @submit.prevent="event => handleSubmit(event as SubmitEvent)"
  >
    <div class="flex gap-x-3 px-1">
      <div class="w-32">
        <fieldset class="fieldset">
          <legend class="fieldset-legend text-sm">{{ $t('driverForm.type') }}</legend>

          <select name="type" v-model="group.type" class="w-full select select-accent" required>
            <option v-for="type in storage.DriverType" :key="type" :value="type">
              {{ $t(`driverCatetory.${type}`) }}
            </option>
          </select>
        </fieldset>
      </div>

      <div class="grow">
        <fieldset class="fieldset">
          <legend class="fieldset-legend text-sm">{{ $t('driverForm.name') }}</legend>
          <input type="text" v-model="group.name" class="input input-accent w-full" required />
        </fieldset>
      </div>
    </div>

    <fieldset class="fieldset">
      <legend class="fieldset-legend text-sm">{{ $t('driverForm.driver') }}</legend>

      <div>
        <div class="max-h-[40vh] overflow-y-auto">
          <div class="grid grid-rows text-sm">
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
                  :class="{ 'text-red-600': groupEditor.notFoundDrivers.value.includes(d.id) }"
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
            v-show="groupEditor.modifiedDrivers.value"
            type="submit"
            id="driver-submit-btn"
            class="btn btn-secondary px-2"
          >
            <font-awesome-icon icon="fa-solid fa-floppy-disk" />
          </button>

          <button type="button" class="btn btn-primary px-2" @click="inputModal?.show()">
            <font-awesome-icon icon="fa-regular fa-square-plus" />
          </button>
        </div>
      </div>
    </fieldset>

    <div class="flex h-8 gap-x-5">
      <button
        type="button"
        class="grow btn"
        style="--btn-color: var(--color-gray-100)"
        @click="$router.back()"
      >
        {{ $t('common.back') }}
      </button>

      <button type="submit" class="grow btn btn-secondary">
        {{ $t('common.save') }}
      </button>
    </div>
  </form>

  <DriverInputModal
    @submit="
      newDriver => {
        console.log(newDriver)
        if (newDriver.id) {
          group.drivers = group.drivers.map(d => (d.id == newDriver.id ? newDriver : d))
        } else {
          group.drivers.push({
            ...newDriver,
            id: `new:${group.drivers.length + 1}` // assign a temporary ID for editing
          })

          console.table(toRaw(group))
        }
        inputModal?.hide()
      }
    "
    ref="inputModal"
  ></DriverInputModal>

  <UnsaveConfirmModal ref="questionModal"></UnsaveConfirmModal>
</template>

<style scoped>
legend:has(+ input:required, + select:required):after,
legend:has(+ div > input:required):after {
  content: ' *';
  color: red;
}
</style>
