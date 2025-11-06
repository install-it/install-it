<script setup lang="ts">
import { useDriverGroupStore, useMatchRuleStore } from '@/store'
// import { storage } from '@/wailsjs/go/models'
import * as matchRuleStorage from '@/wailsjs/go/storage/MatchRuleStorage'

const [ruleStore, driverStore] = [useMatchRuleStore(), useDriverGroupStore()]
</script>

<template>
  <div class="flex h-full flex-col gap-y-2">
    <div class="flex min-h-48 grow flex-col overflow-y-scroll rounded-md p-1.5 shadow-md">
      <div
        v-for="rs in ruleStore.ruleSets"
        :key="rs.id"
        class="driver-card m-1 rounded-lg border border-gray-200 px-2 py-1 shadow-sm"
      >
        <div class="flex justify-between">
          <div class="flex min-w-0 items-center gap-1">
            <span
              v-if="rs.should_hit_all"
              class="badge badge-sm px-1 text-nowrap text-white"
              style="--badge-color: var(--color-rose-400)"
            >
              {{ $t('matchRule.hitAll') }}
            </span>

            <h2 class="truncate">{{ rs.name }}</h2>
          </div>

          <div class="flex gap-x-1.5 py-1">
            <RouterLink
              :to="`/match-rules/edit/${rs.id}`"
              class="btn size-6 btn-xs"
              :title="$t('common.edit')"
            >
              <font-awesome-icon icon="fa-solid fa-pen-to-square" class="text-gray-500" />
            </RouterLink>

            <button
              class="btn size-6 btn-xs"
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
              class="btn size-6 btn-xs"
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

        <div class="grid grid-cols-10 gap-1 bg-gray-100 py-1 text-xs">
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
                class="me-0.5 badge badge-sm px-1 text-white md:me-1"
                style="--badge-color: var(--color-rose-400)"
              >
                {{ $t('matchRule.hitAll') }}
              </span>

              <span
                v-if="r.is_case_sensitive"
                class="me-0.5 badge badge-sm px-1 md:me-1"
                style="--badge-color: var(--color-orange-300)"
              >
                Aa
              </span>

              <span
                v-for="(v, i) in r.values"
                :key="i"
                class="me-0.5 badge badge-sm px-1 badge-neutral md:me-1"
              >
                {{ v }}
              </span>
            </div>
          </div>
        </div>

        <hr class="my-1" />

        <div class="flex gap-2 text-xs">
          <p class="content-center font-semibold">{{ $t('matchRule.matchTo') }}</p>

          <div class="line-clamp-2 flex-1">
            <span
              v-for="(group, i) in driverStore.groups.filter(g =>
                rs.driver_group_ids.includes(g.id)
              )"
              :key="i"
              class="me-0.5 badge badge-xs px-1 md:badge-sm md:px-1.5"
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
          {{ $t('common.create') }}
        </RouterLink>
      </button>
    </div>
  </div>
</template>
