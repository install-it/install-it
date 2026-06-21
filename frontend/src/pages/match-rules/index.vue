<script setup lang="ts">
// import { storage } from '@/wailsjs/go/models'
import * as ruleSetStorage from '@/wailsjs/go/storage/RuleSetStorage'
import { useRouter } from 'vue-router'

const toast = useToast()

const router = useRouter()

const [ruleStore, driverStore] = [useMatchRuleStore(), useDriverGroupStore()]

const { scrollContainer } = useScrollPosition('matchRule', () =>
  ['/match-rules/create', '/match-rules/:id/edit'].some(
    v =>
      (router.options.history.state.forward ?? router.options.history.state.back)
        ?.toString()
        .includes(v) ?? false
  )
)
</script>

<template>
  <div class="flex h-full flex-col gap-y-2">
    <div>
      <h1 class="font-bold">{{ $t('matchRule.matchRule') }}</h1>

      <p class="text-xs">{{ $t('matchRule.matchRuleHelp') }}</p>
    </div>

    <div
      ref="scrollContainer"
      class="flex min-h-48 grow flex-col overflow-y-scroll rounded-md p-1.5 shadow-md"
    >
      <div
        v-for="rs in ruleStore.ruleSets"
        :key="rs.id"
        class="driver-card m-1 rounded-lg border border-gray-200 px-2 py-1 shadow-sm"
      >
        <div class="flex justify-between">
          <div class="flex min-w-0 items-center gap-1">
            <UBadge v-if="rs.should_hit_all" size="sm" class="bg-rose-400 text-nowrap text-white">
              {{ $t('matchRule.hitAll') }}
            </UBadge>

            <h2 class="truncate">{{ rs.name }}</h2>
          </div>

          <div class="flex gap-x-1.5 py-1">
            <RouterLink :to="`/match-rules/${rs.id}/edit`" :title="$t('common.edit')">
              <UButton color="neutral" variant="outline" size="xs" class="h-6">
                <Icon icon="mdi:pencil" class="text-gray-500" />
              </UButton>
            </RouterLink>

            <UButton
              color="neutral"
              variant="outline"
              size="xs"
              class="h-6"
              :title="$t('common.clone')"
              @click="
                ruleSetStorage.Clone(rs.id).then(() =>
                  ruleSetStorage
                    .All()
                    .then(rs => (ruleStore.ruleSets = rs))
                    .catch(() => {
                      toast.add({ title: $t('toast.readDriverFailed'), color: 'error' })
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
              :title="$t('common.delete')"
              @click="
                ruleSetStorage.Remove(rs.id).then(() =>
                  ruleSetStorage
                    .All()
                    .then(rs => (ruleStore.ruleSets = rs))
                    .catch(() => {
                      toast.add({ title: $t('toast.readDriverFailed'), color: 'error' })
                    })
                )
              "
            >
              <Icon icon="mdi:trash-can" class="text-gray-500" />
            </UButton>
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
              <UBadge
                v-if="r.should_hit_all"
                size="sm"
                class="me-0.5 bg-rose-400 text-white md:me-1"
              >
                {{ $t('matchRule.hitAll') }}
              </UBadge>

              <UBadge
                v-if="r.is_case_sensitive"
                size="sm"
                class="me-0.5 bg-orange-300 text-white md:me-1"
              >
                Aa
              </UBadge>

              <UBadge
                v-for="(v, i) in r.values"
                :key="i"
                size="sm"
                color="tertiary"
                class="me-0.5 md:me-1"
              >
                {{ v }}
              </UBadge>
            </div>
          </div>
        </div>

        <hr class="my-1" />

        <div class="flex gap-2 text-xs">
          <p class="content-center font-semibold">{{ $t('matchRule.matchTo') }}</p>

          <div class="line-clamp-2 flex-1">
            <UBadge
              v-for="(group, i) in driverStore.groups.filter(g =>
                rs.driver_group_ids?.includes(g.id)
              )"
              :key="i"
              size="sm"
              class="me-0.5 text-zinc-600"
              :style="`background-color: var(--color-${group.type})`"
            >
              {{ group.name }}
            </UBadge>
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-end gap-x-3">
      <RouterLink :to="{ path: '/match-rules/create' }">
        <UButton color="primary" size="sm">
          {{ $t('common.create') }}
        </UButton>
      </RouterLink>
    </div>
  </div>
</template>
