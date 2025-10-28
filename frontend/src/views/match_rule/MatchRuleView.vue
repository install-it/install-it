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
          <div class="flex items-center gap-1">
            <span
              v-if="rs.should_hit_all"
              class="badge px-1.5 text-xs text-white"
              style="--badge-color: var(--color-rose-400)"
            >
              {{ $t('matchRule.hitAll') }}
            </span>

            <h2 class="truncate oveflow-x-hidden text-sm font-bold">
              {{ rs.name }}
            </h2>
          </div>

          <div class="flex gap-x-1.5 py-1">
            <RouterLink
              :to="`/match-rules/edit/${rs.id}`"
              class="btn btn-xs size-6"
              :title="$t('common.edit')"
            >
              <font-awesome-icon icon="fa-solid fa-pen-to-square" class="text-gray-500" />
            </RouterLink>

            <button
              class="btn btn-xs size-6"
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
              :title="$t('common.clone')"
            >
              <font-awesome-icon icon="fa-solid fa-clone" class="text-gray-500" />
            </button>

            <button
              class="btn btn-xs size-6"
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
              :title="$t('common.delete')"
            >
              <font-awesome-icon icon="fa-solid fa-trash" class="text-gray-500" />
            </button>
          </div>
        </div>

        <div class="grid grid-cols-10 gap-1 py-1 text-xs bg-gray-100">
          <div class="col-span-2 font-semibold">{{ $t('matchRule.source') }}</div>
          <div class="col-span-2 font-semibold">{{ $t('matchRule.operator') }}</div>
          <div class="col-span-6 font-semibold">{{ $t('matchRule.pattern') }}</div>
        </div>

        <div
          v-for="(r, ri) in rs.rules"
          :key="ri"
          class="grid grid-cols-10 items-center gap-1 py-1 text-xs"
        >
          <div class="col-span-2">
            {{ $t(`common.${r.source}`) }}
          </div>

          <div class="col-span-2">
            <span class="font-mono">
              {{ $t(`matchRule.${r.operator}`) }}
            </span>
          </div>

          <div class="col-span-6">
            <div class="line-clamp-2">
              <span
                v-if="r.should_hit_all"
                class="badge px-0.5 me-0.5 h-4 text-white text-xs"
                style="--badge-color: var(--color-rose-400)"
              >
                {{ $t('matchRule.hitAll') }}
              </span>

              <span
                v-if="r.is_case_sensitive"
                class="badge px-0.5 me-0.5 h-4 text-xs"
                style="--badge-color: var(--color-orange-300)"
              >
                Aa
              </span>

              <span
                v-for="(v, i) in r.values"
                :key="i"
                class="badge badge-neutral px-0.5 me-0.5 h-4 text-xs"
              >
                {{ v }}
              </span>
            </div>
          </div>
        </div>

        <hr class="my-1" />

        <div class="flex gap-2 text-xs">
          <p class="content-center font-semibold">{{ $t('matchRule.matchTo') }}</p>

          <div class="flex-1 line-clamp-2">
            <span
              v-for="(group, i) in driverStore.groups.filter(g =>
                rs.driver_group_ids.includes(g.id)
              )"
              :key="i"
              class="badge px-0.5 me-0.5 h-4 text-xs"
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
