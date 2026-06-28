<script setup lang="ts">
import { storage } from '@/wailsjs/go/models'
import * as driverGroupStorage from '@/wailsjs/go/storage/DriverGroupStorage'
import * as groupStorage from '@/wailsjs/go/storage/DriverGroupStorage'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

function categoryKey(type: string): string {
  return `category${type.charAt(0).toUpperCase() + type.slice(1)}`
}

const toast = useToast()

const router = useRouter()

const groupStore = useDriverGroupStore()

const reordering = ref(false)

const { scrollContainer } = useScrollPosition('driverGroup', () =>
  ['/drivers/create', '/drivers/:id/edit'].some(
    v =>
      (router.options.history.state.forward ?? router.options.history.state.back)
        ?.toString()
        .includes(v) ?? false
  )
)
</script>

<template>
  <div class="flex h-full flex-col gap-y-2">
    <PageHeader
      variant="compact"
      :title="$t('fieldInstallOption')"
      :description="$t('descInstallOption')"
    >
      <div class="flex w-1/3 flex-wrap justify-center gap-0.5 px-0.5 text-xs select-none">
        <router-link
          :to="{ path: '/drivers' }"
          class="flex-1/3 truncate rounded-sm text-center font-bold uppercase shadow-xs"
          :class="{
            'bg-white text-half-baked-500': $route.query.type != undefined,
            'bg-half-baked-500 text-white': $route.query.type == undefined
          }"
          draggable="false"
        >
          {{ $t('all') }}
        </router-link>

        <router-link
          v-for="type in storage.DriverType"
          :key="type"
          :to="{ path: '/drivers', query: { type: type } }"
          class="flex-1/3 truncate rounded-sm text-center font-bold uppercase shadow-xs"
          :class="{
            'bg-white text-half-baked-500': $route.query.type !== type,
            'bg-half-baked-500 text-white': $route.query.type === type
          }"
          draggable="false"
        >
          {{ $t(categoryKey(type)) }}
        </router-link>
      </div>
    </PageHeader>

    <div
      ref="scrollContainer"
      class="flex min-h-48 grow flex-col overflow-y-scroll rounded-md p-1.5 shadow-md"
    >
      <div
        v-for="(g, i) in groupStore.groups.filter(
          g => $route.query.type == undefined || g.type == $route.query.type
        )"
        :key="g.id"
        class="driver-card m-1 rounded-lg border border-gray-200 px-2 py-1 shadow-sm"
        :class="reordering ? 'cursor-pointer select-none' : ''"
        :draggable="reordering"
        @dragstart="
          event => {
            if (!reordering) {
              return event.preventDefault()
            }

            event.dataTransfer!.setData('id', g.id.toString())
            const fullIdx = groupStore.groups.findIndex(g2 => g2.id === g.id)
            event.dataTransfer!.setData('index', fullIdx.toString())
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

            // async function will cause event.dataTransfer to lose data
            const sourceId = parseInt(event.dataTransfer!.getData('id'))
            const sourceIdx = parseInt(event.dataTransfer!.getData('index'))
            const targetIdx = groupStore.groups.findIndex(g2 => g2.id === g.id)
            let moveBehindIdx = targetIdx
            if (sourceIdx <= targetIdx) {
              moveBehindIdx -= 1
            }
            groupStorage
              .MoveBehind(sourceId, moveBehindIdx)
              .then(() => groupStorage.All().then(gs => (groupStore.groups = gs)))
          }
        "
      >
        <div class="flex justify-between">
          <div class="flex min-w-0 items-center gap-1">
            <UBadge size="sm" :style="`background-color: var(--color-${g.type})`"> &nbsp; </UBadge>

            <h2 class="truncate">
              {{ g.name }}
            </h2>
          </div>

          <div class="flex gap-x-1.5 py-1">
            <RouterLink :to="`/drivers/${g.id}/edit`" :title="$t('edit')">
              <UButton color="neutral" variant="outline" size="xs" class="h-6">
                <Icon icon="mdi:pencil" class="text-gray-500" />
              </UButton>
            </RouterLink>

            <UButton
              color="neutral"
              variant="outline"
              size="xs"
              class="h-6"
              :title="$t('clone')"
              @click="
                groupStorage.Clone(g.id).then(() =>
                  driverGroupStorage
                    .All()
                    .then(gs => (groupStore.groups = gs))
                    .catch(() => {
                      toast.add({ title: $t('toastReadDriversFailed'), color: 'error' })
                    })
                )
              "
            >
              <Icon icon="mdi:content-duplicate" class="text-gray-500" />
            </UButton>

            <UButton
              color="neutral"
              variant="outline"
              size="xs"
              class="h-6"
              :title="$t('delete')"
              @click="
                groupStorage.Remove(g.id).then(() =>
                  driverGroupStorage
                    .All()
                    .then(gs => (groupStore.groups = gs))
                    .catch(() => {
                      toast.add({ title: $t('toastReadDriversFailed'), color: 'error' })
                    })
                )
              "
            >
              <Icon icon="mdi:trash-can" class="text-gray-500" />
            </UButton>
          </div>
        </div>

        <div class="grid grid-cols-12 gap-1 bg-gray-100 py-1 text-xs">
          <div class="col-span-2 font-medium lg:col-span-3">{{ $t('name') }}</div>

          <div class="col-span-5 font-medium lg:col-span-5">{{ $t('path') }}</div>

          <div class="col-span-3 font-medium lg:col-span-3">{{ $t('fieldArgument') }}</div>

          <div class="col-span-2 font-medium lg:col-span-1">
            {{ $t('fieldOther') }}
          </div>
        </div>

        <div v-for="d in g.drivers" :key="d.id" class="grid grid-cols-12 gap-1 py-1 text-xs">
          <div class="col-span-2 line-clamp-2 break-all lg:col-span-3">
            {{ d.name }}
          </div>

          <div
            class="col-span-5 line-clamp-2 break-all lg:col-span-5"
            :class="{ 'text-red-600': groupStore.notFoundDrivers.includes(d.id) }"
          >
            {{ d.path }}
          </div>

          <div class="col-span-3 line-clamp-2 break-all lg:col-span-3">
            {{ d.flags }}
          </div>

          <div class="col-span-2 flex gap-x-1 text-sm lg:col-span-1">
            <span
              v-show="g.mutuallyExclusive"
              class="inline-block max-h-5 rounded-xs bg-orange-300 p-0.5"
              :title="$t('fieldMutuallyExclusive')"
            >
              <Icon icon="mdi:chart-timeline" />
            </span>

            <span
              v-show="d.incompatibles.length > 0"
              class="inline-block max-h-5 rounded-xs bg-yellow-300 p-0.5"
              :title="$t('labelIncompatibleWith')"
            >
              <Icon icon="mdi:source-merge" />
            </span>

            <span
              v-show="d.allowRtCodes.length > 0"
              class="inline-block max-h-5 rounded-xs bg-blue-300 p-0.5"
              :title="$t('fieldAllowedExitCode')"
            >
              <Icon icon="mdi:numeric-1-box-outline" />
            </span>
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-end gap-x-3">
      <div v-show="groupStore.groups?.filter(d => d.type == $route.query.type).length > 1">
        <UButton
          type="button"
          size="sm"
          class="text-white"
          :style="
            reordering
              ? '--btn-color: var(--color-apple-green-800); animation: var(--animate-blink-75);'
              : '--btn-color: #D9BD68'
          "
          @click="reordering = !reordering"
        >
          {{ reordering ? $t('view') : $t('fieldOrder') }}
        </UButton>
      </div>

      <RouterLink :to="{ path: '/drivers/create', query: { type: $route.query.type } }">
        <UButton color="primary" size="sm">
          {{ $t('create') }}
        </UButton>
      </RouterLink>
    </div>
  </div>
</template>
