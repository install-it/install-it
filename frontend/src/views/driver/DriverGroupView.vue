<script setup lang="ts">
import { useDriverGroupStore } from '@/store'
import { storage } from '@/wailsjs/go/models'
import * as driverGroupStorage from '@/wailsjs/go/storage/DriverGroupStorage'
import * as groupStorage from '@/wailsjs/go/storage/DriverGroupStorage'
import { ref } from 'vue'

const groupStore = useDriverGroupStore()
const reordering = ref(false)
</script>

<template>
  <div class="flex flex-col h-full gap-y-2">
    <div class="flex flex-row gap-x-3 list-none text-center">
      <router-link
        :to="{ path: '/drivers' }"
        class="w-full py-3 text-xs font-bold uppercase shadow-lg rounded-sm"
        :class="{
          'text-half-baked-600 bg-white': $route.query.type != undefined,
          'text-white bg-half-baked-600': $route.query.type == undefined
        }"
        draggable="false"
      >
        {{ $t(`common.all`) }}
      </router-link>

      <router-link
        v-for="type in storage.DriverType"
        :key="type"
        :to="{ path: '/drivers', query: { type: type } }"
        class="w-full py-3 text-xs font-bold uppercase shadow-lg rounded-sm"
        :class="{
          'text-half-baked-600 bg-white': $route.query.type !== type,
          'text-white bg-half-baked-600': $route.query.type === type
        }"
        draggable="false"
      >
        {{ $t(`driverCatetory.${type}`) }}
      </router-link>
    </div>

    <div class="flex flex-col grow p-1.5 min-h-48 overflow-y-scroll shadow-md rounded-md">
      <div
        v-for="(g, i) in groupStore.groups.filter(
          g => $route.query.type == undefined || g.type == $route.query.type
        )"
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

            groupStorage.IndexOf(g.id).then(targetIndex => {
              if (parseInt(sourcePosition) <= i) {
                // aligning MoveBehind's logic and UI draging's logic
                targetIndex -= 1
              }

              groupStorage.MoveBehind(sourceId, targetIndex).then(result => {
                groupStore.groups = result
              })
            })
          }
        "
        :draggable="reordering"
      >
        <div class="flex justify-between">
          <p class="my-1 truncate oveflow-x-hidden align-middle">
            <span class="badge h-4 px-1 me-1" :style="`--badge-color: var(--color-${g.type})`">
              &nbsp;
            </span>
            <span>{{ g.name }}</span>
          </p>

          <div class="flex gap-x-1.5 py-1">
            <RouterLink
              :to="`/drivers/edit/${g.id}`"
              class="px-1 bg-gray-200 hover:bg-gray-300 transition-all rounded-sm"
            >
              <font-awesome-icon icon="fa-solid fa-pen-to-square" class="text-gray-500" />
            </RouterLink>

            <button
              class="px-1 bg-gray-200 hover:bg-gray-300 transition-all rounded-sm"
              @click="
                groupStorage.Add(g).then(() =>
                  driverGroupStorage
                    .All()
                    .then(gs => (groupStore.groups = gs))
                    .catch(() => {
                      $toast.error($t('toast.readDriverFailed'))
                    })
                )
              "
            >
              <font-awesome-icon icon="fa-solid fa-clone" class="text-gray-500" />
            </button>

            <button
              class="px-1 bg-gray-200 hover:bg-gray-300 transition-all rounded-sm"
              @click="
                groupStorage.Remove(g.id).then(() =>
                  driverGroupStorage
                    .All()
                    .then(gs => (groupStore.groups = gs))
                    .catch(() => {
                      $toast.error($t('toast.readDriverFailed'))
                    })
                )
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
            :class="{ 'text-red-600': groupStore.notFoundDrivers.includes(d.id) }"
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
        v-show="groupStore.groups?.filter(d => d.type == $route.query.type).length > 1"
        type="button"
        class="btn text-white"
        :style="
          reordering
            ? '--btn-color: var(--color-apple-green-800); animation: var(--animate-blink-75);'
            : '--btn-color: #D9BD68'
        "
        @click="reordering = !reordering"
      >
        {{ reordering ? $t('driverForm.view') : $t('driverForm.order') }}
      </button>

      <button class="btn btn-primary">
        <RouterLink :to="{ path: '/drivers/create', query: { type: $route.query.type } }">
          {{ $t('driverForm.create') }}
        </RouterLink>
      </button>
    </div>
  </div>
</template>
