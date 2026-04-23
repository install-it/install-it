<script setup lang="ts">
import DriverInputModal from '@/components/DriverInputModal.vue'
import UnsaveConfirmModal from '@/components/UnsaveConfirmModal.vue'
import { useDriverGroupStore } from '@/store'
import { storage } from '@/wailsjs/go/models'
import * as groupStorage from '@/wailsjs/go/storage/DriverGroupStorage'
import { toRaw, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'

const props = defineProps<{ id?: string }>()

const { t } = useI18n()

const $route = useRoute()
const $router = useRouter()
const toast = useToast()

const questionModal = useTemplateRef('questionModal')
const inputModal = useTemplateRef('inputModal')

const groupStore = useDriverGroupStore()

const groupEditor = groupStore.editor(
  props.id,
  storage.DriverType[
    ($route.query.type as string | undefined)?.toUpperCase() as keyof typeof storage.DriverType
  ] ?? storage.DriverType.NETWORK
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
    toast.add({ title: t('toast.addAtLeastOneDriver'), color: 'warning' })
    return
  }

  const handleSuccess = () => {
    toast.add({ title: t('toast.updated'), color: 'success' })
    groupStorage.All().then(newDriverGroups => {
      groupStore.groups = newDriverGroups
      groupEditor.reset()

      if (event.submitter?.id !== 'driver-submit-btn') {
        $router.back()
      } else {
        $router.replace({ path: `/drivers/${group.value.id}/edit` })
      }
    })
  }

  if (group.value.id == undefined) {
    groupStorage
      .Add(group.value)
      .then(gid => (group.value.id = gid))
      .then(handleSuccess)
      .catch(reason => toast.add({ title: reason.toString(), color: 'error' }))
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
          <legend class="fieldset-legend text-sm">{{ $t('driverForm.type') }}</legend>

          <USelect
            v-model="group.type"
            name="type"
            class="w-full"
            :items="
              Object.values(storage.DriverType).map(type => ({
                label: $t(`driverCatetory.${type}`),
                value: type
              }))
            "
            required
          />
        </fieldset>
      </div>

      <div class="grow">
        <fieldset class="fieldset">
          <legend class="fieldset-legend text-sm">{{ $t('driverForm.name') }}</legend>

          <UInput v-model="group.name" type="text" class="w-full" required />
        </fieldset>
      </div>
    </div>

    <fieldset class="fieldset">
      <legend class="fieldset-legend text-sm">{{ $t('driverForm.driver') }}</legend>

      <div>
        <div class="max-h-[40vh] overflow-y-auto">
          <div class="grid-rows grid text-sm">
            <div class="grid grid-cols-10 gap-2 border-y py-1.5">
              <div class="col-span-2">{{ $t('driverForm.name') }}</div>

              <div class="col-span-3">{{ $t('driverForm.path') }}</div>

              <div class="col-span-2">{{ $t('driverForm.argument') }}</div>

              <div class="col-span-2">{{ $t('driverForm.otherSetting') }}</div>
            </div>

            <div v-if="group.drivers.length == 0" class="py-1 text-center last:border-b">N/A</div>

            <div
              v-for="(d, i) in group.drivers"
              v-else
              :key="d.id"
              class="grid grid-cols-10 items-center gap-2 border-b py-1.5 text-xs"
              :class="{ 'bg-lime-50': d.id.includes('new:') }"
            >
              <div class="col-span-2">
                <p class="line-clamp-2 break-all">
                  {{ d.name }}
                </p>
              </div>

              <div class="col-span-3">
                <p
                  class="line-clamp-2 font-mono break-all"
                  :class="{ 'text-red-600': groupEditor.notFoundDrivers.value.includes(d.id) }"
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
                  :title="$t('driverForm.incompatibleWith')"
                >
                  <Icon icon="mdi:source-merge" />
                </span>

                <span
                  v-show="d.allowRtCodes.length > 0"
                  class="inline-block max-h-5 rounded-xs bg-blue-300 p-0.5"
                  :title="$t('driverForm.allowedExitCode')"
                >
                  <Icon icon="mdi:numeric-0-box-outline" />
                </span>
              </div>

              <div>
                <div class="flex gap-x-2">
                  <button type="button" :title="$t('common.edit')" @click="inputModal?.show(d)">
                    <Icon icon="mdi:pencil" class="size-4" />
                  </button>

                  <button
                    type="button"
                    :title="$t('common.delete')"
                    @click="group.drivers.splice(i, 1)"
                  >
                    <Icon icon="mdi:trash-can" class="size-4" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <p class="text-hint">
            {{ $t('driverForm.driverGroupHelp') }}
          </p>

          <p v-show="groupEditor.modifiedDrivers.value" class="text-hint">
            {{ $t('driverForm.incompatibleForNewHelp') }}
          </p>
        </div>

        <div class="flex justify-end gap-x-3">
          <div v-show="groupEditor.modifiedDrivers.value">
            <UButton id="driver-submit-btn" type="submit" class="px-2" color="secondary">
              <Icon icon="mdi:content-save" />
            </UButton>
          </div>

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
        {{ $t('common.back') }}
      </UButton>

      <UButton type="submit" class="grow justify-center" color="secondary">
        {{ $t('common.save') }}
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
            id: `new:${group.drivers.length + 1}` // assign a temporary ID for editing
          })

          console.table(toRaw(group))
        }
        inputModal?.hide()
      }
    "
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
