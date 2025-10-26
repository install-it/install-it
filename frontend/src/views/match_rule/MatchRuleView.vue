<script setup lang="ts">
import { useDriverGroupStore, useMatchRuleStore } from '@/store'
// import { storage } from '@/wailsjs/go/models'
import * as matchRuleStorage from '@/wailsjs/go/storage/MatchRuleStorage'

const [ruleStore, driverStore] = [useMatchRuleStore(), useDriverGroupStore()]
</script>

<template>
  <div class="flex flex-col h-full gap-y-2">
    <div class="flex flex-col grow p-1.5 min-h-48 overflow-y-scroll shadow-md rounded-md">
      <div
        v-for="rs in ruleStore.ruleSets"
        :key="rs.id"
        class="driver-card m-1 px-2 py-1 border border-gray-200 rounded-lg shadow-sm"
      >
        <div class="flex justify-between">
          <h2 class="my-1 truncate oveflow-x-hidden align-middle text-sm font-bold">
            {{ rs.name }}
          </h2>

          <div class="flex gap-x-1.5 py-1">
            <RouterLink
              :to="`/match-rules/edit/${rs.id}`"
              class="px-1 bg-gray-200 hover:bg-gray-300 transition-all rounded-sm"
            >
              <font-awesome-icon icon="fa-solid fa-pen-to-square" class="text-gray-500" />
            </RouterLink>

            <button
              class="px-1 bg-gray-200 hover:bg-gray-300 transition-all rounded-sm"
              @click="
                matchRuleStorage.Add(rs).then(() =>
                  matchRuleStorage
                    .All()
                    .then(rs => (ruleStore.ruleSets = rs))
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
                matchRuleStorage.Remove(rs.id).then(() =>
                  matchRuleStorage
                    .All()
                    .then(rs => (ruleStore.ruleSets = rs))
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

        <div class="grid grid-cols-5 gap-1 py-1 text-xs bg-gray-100">
          <div class="col-span-1 font-semibold">{{ $t('matchRule.source') }}</div>
          <div class="col-span-1 font-semibold">{{ $t('matchRule.operator') }}</div>
          <div class="col-span-3 font-semibold">{{ $t('matchRule.pattern') }}</div>
        </div>

        <div v-for="(r, ri) in rs.rules" :key="ri" class="grid grid-cols-5 gap-1 py-1 text-xs">
          <div class="col-span-1">
            {{ $t(`common.${r.source}`) }}
          </div>
          <div class="col-span-1">
            <span class="bg-gray-200 rounded-sm">
              <font-awesome-icon v-if="r.is_case_sensitive" icon="fa-solid fa-a" />
            </span>
            {{ $t(`matchRule.${r.operator}`) }}
          </div>
          <div class="col-span-3 space-x-1 space-y-0.5 line-clamp-3">
            <span v-for="(v, i) in r.values" :key="i" class="badge badge-neutral badge-sm px-0.5">
              {{ v }}
            </span>
          </div>
        </div>

        <hr class="my-1" />

        <div class="flex gap-2 text-xs">
          <p class="font-semibold">{{ $t('matchRule.matchTo') }}</p>
          <div class="flex-1 line-clamp-2 space-x-1 space-y-1">
            <span
              v-for="(group, i) in driverStore.groups.filter(g =>
                rs.driver_group_ids.includes(g.id)
              )"
              :key="i"
              class="badge badge-sm px-0.5"
              :style="`--badge-color: var(--color-${group.type})`"
            >
              {{ group.name }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-end gap-x-3">
      <button class="btn btn-primary">
        <RouterLink :to="{ path: '/match-rules/create' }">
          {{ $t('matchRule.createRule') }}
        </RouterLink>
      </button>
    </div>
  </div>
</template>
