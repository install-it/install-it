<script setup lang="ts">
import { ExecutableExists } from '@/wailsjs/go/main/App'
import { store } from '@/wailsjs/go/models'
import * as groupStore from '@/wailsjs/go/store/DriverGroupStore'
import { onBeforeMount, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'vue-toast-notification'

const { t } = useI18n()

const $router = useRouter()

const $route = useRoute()

const $toast = useToast({ position: 'top-right' })

const driverType = ref(
  store.DriverType[$route.query.type?.toString().toUpperCase() as keyof typeof store.DriverType] ??
    store.DriverType.NETWORK
)

const groups = ref<Array<store.DriverGroup>>([])

/** driver ID of drivers that the executable cannot be found */
const notExistDrivers = ref<Array<string>>([])

const reordering = ref(false)

onBeforeMount(() => {
  groupStore
    .All()
    .then(g => (groups.value = g))
    .catch(() => {
      $toast.error(t('toast.readDriverFailed'))
    })
})

watch(groups, newGroups => {
  Promise.all(
    newGroups.flatMap(g =>
      g.drivers.flatMap(d => ExecutableExists(d.path).then(exist => ({ id: d.id, exist: exist })))
    )
  ).then(results => {
    notExistDrivers.value = results
      .map(result => (result.exist ? undefined : result.id))
      .filter(v => v !== undefined)
  })
})

watch(driverType, newType => {
  $router.replace({ path: '/drivers', query: { type: newType } })
})
</script>

<template>
  <div class="flex flex-col h-full gap-y-2">
    <div class="flex flex-row gap-x-3 list-none text-center">
      <button
        v-for="type in store.DriverType"
        :key="type"
        class="w-full py-3 text-xs font-bold uppercase shadow-lg rounded-sm"
        :class="{
          'text-half-baked-600 bg-white': driverType !== type,
          'text-white bg-half-baked-600': driverType === type
        }"
        @click="driverType = type"
      >
        {{ t(`driverCatetory.${type}`) }}
      </button>
    </div>

    <div class="flex flex-col grow py-2 min-h-48 overflow-y-scroll shadow-lg rounded-sm">
      <div
        v-for="(g, i) in groups.filter(g => g.type == driverType)"
        :key="g.id"
        class="driver-card m-1 px-2 py-1 border border-gray-200 rounded-lg shadow-sm"
        :class="reordering ? 'select-none cursor-pointer' : ''"
        @dragstart="
          event => {
            if (!reordering) {
              return event.preventDefault()
            }

            event.dataTransfer!.setData('id', g.id)
            event.dataTransfer!.setData('position', i.toString())
          }
        "
        @dragover.prevent="
          event => {
            ;(event.target as HTMLDivElement)
              .closest('.driver-card')!
              .classList.add('border-b-2', 'border-b-half-baked-700')
          }
        "
        @dragenter.prevent
        @dragleave="
          event => {
            ;(event.target as HTMLDivElement)
              .closest('.driver-card')!
              .classList.remove('border-b-2', 'border-b-half-baked-700')
          }
        "
        @drop="
          async event => {
            ;(event.target as HTMLDivElement)
              .closest('.driver-card')!
              .classList.remove('border-b-2', 'border-b-half-baked-700')

            // async functuion will cause event.dataTransfer lost data
            const sourceId = event.dataTransfer!.getData('id')
            const sourcePosition = event.dataTransfer!.getData('position')

            groupStore.IndexOf(g.id).then(targetIndex => {
              if (parseInt(sourcePosition) <= i) {
                // aligning MoveBehind's logic and UI draging's logic
                targetIndex -= 1
              }

              groupStore.MoveBehind(sourceId, targetIndex).then(result => {
                groups = result
              })
            })
          }
        "
        :draggable="reordering"
      >
        <div class="flex justify-between">
          <p class="my-1 truncate oveflow-x-hidden align-middle">
            <span class="badge h-4 px-1 me-1" :class="[`badge-${g.type}`]">&nbsp;</span>
            <span>{{ g.name }}</span>
          </p>

          <div class="flex gap-x-1.5 py-1">
            <RouterLink :to="`/drivers/edit/${g.id}`" class="px-0.5 bg-gray-200 rounded-sm">
              <font-awesome-icon icon="fa-solid fa-pen-to-square" class="text-gray-500" />
            </RouterLink>

            <button
              class="px-0.5 bg-gray-200 rounded-sm"
              @click="
                () => {
                  groupStore.Add(g).then(() => {
                    groupStore.All().then(g => {
                      groups = g
                    })
                  })
                }
              "
            >
              <font-awesome-icon icon="fa-solid fa-clone" class="text-gray-500" />
            </button>

            <button
              class="px-0.5 bg-gray-200 rounded-sm"
              @click="
                () => {
                  groupStore.Remove(g.id).then(() => {
                    groupStore.All().then(g => {
                      groups = g
                    })
                  })
                }
              "
            >
              <font-awesome-icon icon="fa-solid fa-trash" class="text-gray-500" />
            </button>
          </div>
        </div>

        <div class="grid grid-cols-12 gap-1 py-1 text-xs bg-gray-100">
          <div class="col-span-2 lg:col-span-3 font-medium">{{ $t('driverForm.name') }}</div>
          <div class="col-span-5 lg:col-span-5 font-medium">{{ $t('driverForm.path') }}</div>
          <div class="col-span-3 lg:col-span-3 font-medium">{{ $t('driverForm.argument') }}</div>
          <div class="col-span-2 lg:col-span-1 font-medium">
            {{ $t('driverForm.otherSetting') }}
          </div>
        </div>

        <div v-for="d in g.drivers" :key="d.id" class="grid grid-cols-12 gap-1 py-1 text-xs">
          <div class="col-span-2 lg:col-span-3 break-all line-clamp-2">
            {{ d.name }}
          </div>

          <div
            class="col-span-5 lg:col-span-5 break-all line-clamp-2"
            :class="{ 'text-red-600': notExistDrivers.includes(d.id) }"
          >
            {{ d.path }}
          </div>

          <div class="col-span-3 lg:col-span-3 break-all line-clamp-2">
            {{ d.flags }}
          </div>

          <div class="flex col-span-2 lg:col-span-1 gap-x-1">
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
        </div>
      </div>
    </div>

    <div class="flex justify-end gap-x-3">
      <button
        v-show="groups?.filter(d => d.type == driverType).length > 1"
        type="button"
        class="h-8 px-3 text-white text-sm focus:outline-hidden rounded-sm"
        :class="[
          reordering
            ? 'bg-apple-green-800 hover:bg-apple-green-700 animate-blink-75'
            : 'bg-[#D9BD68] hover:bg-[#E5D195]'
        ]"
        @click="reordering = !reordering"
      >
        {{ reordering ? $t('driverForm.view') : $t('driverForm.order') }}
      </button>

      <RouterLink
        :to="{ path: '/drivers/create', query: { type: driverType } }"
        class="flex items-center h-8 px-3 text-white text-sm focus:outline-hidden bg-powder-blue-800 hover:bg-powder-blue-700 rounded-sm"
      >
        {{ $t('driverForm.create') }}
      </RouterLink>
    </div>
  </div>
</template>
