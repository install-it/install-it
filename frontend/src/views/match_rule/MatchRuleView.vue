<script setup lang="ts">
import { useDriverGroupStore, useMatchRuleStore } from '@/store'
// import { storage } from '@/wailsjs/go/models'
import * as matchRuleStorage from '@/wailsjs/go/storage/MatchRuleStorage'

const ruleStore = useMatchRuleStore()
const driverStore = useDriverGroupStore()
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
          <p class="my-1 truncate oveflow-x-hidden align-middle">
            <span class="text-sm">
              {{
                rs.driver_group_ids
                  .map(id => driverStore.groups.find(g => g.id == id)?.name ?? '')
                  .join(', ')
              }}
            </span>
          </p>

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

        <div class="grid grid-cols-7 gap-1 py-1 text-xs bg-gray-100">
          <div class="col-span-2 font-medium">目標</div>
          <div class="col-span-2 font-medium">規則類型</div>
          <div class="col-span-3 font-medium">規則</div>
        </div>

        <div v-for="(r, ri) in rs.rules" :key="ri" class="grid grid-cols-7 gap-1 py-1 text-xs">
          <div class="col-span-2 font-medium">{{ r.source }}</div>
          <div class="col-span-2 font-medium">
            <span class="bg-gray-200 rounded-sm">
              <font-awesome-icon v-if="r.is_case_sensitive" icon="fa-solid fa-a" />
            </span>
            {{ r.type }}
          </div>
          <div class="col-span-3 font-medium">{{ r.values }}</div>
        </div>
      </div>
    </div>

    <div class="flex justify-end gap-x-3">
      <button class="btn btn-primary">
        <RouterLink :to="{ path: '/match-rules/create' }"> 建立配對規則 </RouterLink>
      </button>
    </div>
  </div>
</template>
