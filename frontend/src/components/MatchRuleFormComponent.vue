<script setup lang="ts">
import DriverSelector from '@/components/DriverSelector.vue'
import MatchRuleInputModal from '@/components/MatchRuleInputModal.vue'
import { storage } from '@/wailsjs/go/models'
import * as ruleSetStorage from '@/wailsjs/go/storage/RuleSetStorage'
import { computed, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const props = defineProps<{ id?: number }>()

const inputModal = useTemplateRef('inputModal')

const { t } = useI18n()

const $router = useRouter()

function commonKey(suffix: string): string {
  return `common${suffix.charAt(0).toUpperCase() + suffix.slice(1)}`
}

function matchRuleKey(suffix: string): string {
  return `matchRule${suffix.charAt(0).toUpperCase() + suffix.slice(1)}`
}

const toast = useToast()

const [ruleStore, groupStore] = [useMatchRuleStore(), useDriverGroupStore()]

const sourceRuleSet = computed(
  () =>
    ruleStore.ruleSets.find(r => r.id === props.id) ??
    new storage.RuleSet({ rules: [], driver_group_ids: [] })
)

const { data: ruleSet } = useEditor({
  source: sourceRuleSet,
  warnOnUnsavedLeave: true
})

function handleSubmit() {
  if (ruleSet.value.rules.length == 0) {
    toast.add({ title: t('toastAddAtLeastOneRule'), color: 'warning' })
    return
  }

  const handleSuccess = () => {
    toast.add({ title: t('toastUpdated'), color: 'success' })

    ruleSetStorage.All().then(newMatchRule => {
      ruleStore.ruleSets = newMatchRule
      $router.push({ path: '/match-rules/' })
    })
  }

  if (!ruleSet.value.id) {
    ruleSetStorage
      .Add(ruleSet.value)
      .then(handleSuccess)
      .catch(reason => toast.add({ title: reason.toString(), color: 'error' }))
  } else {
    ruleSetStorage
      .Update(ruleSet.value)
      .then(handleSuccess)
      .catch(reason => toast.add({ title: reason.toString(), color: 'error' }))
  }
}
</script>

<template>
  <form class="size-full overflow-y-auto" autocomplete="off" @submit.prevent="handleSubmit">
    <div class="mx-auto h-full max-w-full content-center space-y-3 lg:max-w-2xl xl:max-w-4xl">
      <h1 class="mb-2 text-xl font-bold">
        {{ $t('matchRuleCreateRule') }}
      </h1>

      <fieldset class="fieldset">
        <legend class="fieldset-legend text-sm">
          {{ $t('matchRuleName') }}
        </legend>

        <UInput v-model="ruleSet.name" type="text" class="ms-1" />
      </fieldset>

      <fieldset class="fieldset">
        <legend class="text-required fieldset-legend text-sm">
          {{ $t('matchRuleRule') }}
        </legend>

        <div>
          <div class="grid-rows grid text-sm">
            <div class="grid grid-cols-6 gap-2 border-y py-1.5">
              <div class="col-span-1">{{ $t('matchRuleSource') }}</div>

              <div class="col-span-1">{{ $t('matchRuleOperator') }}</div>

              <div class="col-span-3">{{ $t('matchRulePattern') }}</div>

              <div class="col-span-1">{{ $t('commonAction') }}</div>
            </div>

            <div class="max-h-38 overflow-y-auto">
              <div
                v-if="ruleSet.rules.length == 0"
                class="h-full content-center border-b py-1 text-center"
              >
                N/A
              </div>

              <div
                v-for="(r, i) in ruleSet.rules"
                v-else
                :key="i"
                class="grid grid-cols-6 items-center gap-2 border-b py-1.5 text-xs lg:text-sm"
              >
                <div class="col-span-1">
                  <p class="line-clamp-2 break-all">
                    {{ $t(commonKey(r.source)) }}
                  </p>
                </div>

                <div class="col-span-1">
                  <span class="font-mono">
                    {{ $t(matchRuleKey(r.operator)) }}
                  </span>
                </div>

                <div class="col-span-3">
                  <div>
                    <UBadge
                      v-if="r.should_hit_all"
                      size="sm"
                      class="me-0.5 bg-rose-400 text-white md:me-1"
                    >
                      {{ $t('matchRuleHitAll') }}
                    </UBadge>

                    <UBadge
                      v-if="r.is_case_sensitive"
                      size="sm"
                      class="me-0.5 bg-orange-300 md:me-1"
                    >
                      Aa
                    </UBadge>

                    <UBadge
                      v-for="(v, vi) in r.values"
                      :key="vi"
                      size="sm"
                      color="tertiary"
                      class="me-0.5 md:me-1"
                    >
                      {{ v }}
                    </UBadge>
                  </div>
                </div>

                <div class="col-span-1 flex gap-x-1">
                  <div class="flex gap-x-2">
                    <button
                      type="button"
                      :title="$t('commonEdit')"
                      @click="inputModal?.show(JSON.parse(JSON.stringify({ ...r, _id: i })))"
                    >
                      <Icon icon="mdi:pencil" class="size-4" />
                    </button>

                    <button
                      type="button"
                      :title="$t('commonDelete')"
                      @click="ruleSet.rules.splice(i, 1)"
                    >
                      <Icon icon="mdi:trash-can" class="size-4" />
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <p class="text-hint"></p>
          </div>

          <div class="flex justify-end gap-x-3">
            <UButton type="button" class="px-2" color="primary" @click="inputModal?.show()">
              <Icon icon="mdi:plus-box" />
            </UButton>
          </div>
        </div>
      </fieldset>

      <fieldset class="fieldset">
        <legend class="fieldset-legend text-sm">{{ $t('matchRuleMultiRuleMatching') }}</legend>

        <label class="flex w-full cursor-pointer items-center select-none">
          <UCheckbox v-model="ruleSet.should_hit_all" class="me-1.5" />
          {{ $t('matchRuleHitAllRule') }}
        </label>

        <p class="text-hint">{{ $t('matchRuleMultiRuleMatchingHelp') }}</p>
      </fieldset>

      <hr />

      <fieldset class="fieldset">
        <legend class="text-required fieldset-legend text-sm">
          {{ $t('matchRuleMatchTo') }}
        </legend>

        <DriverSelector
          v-model="ruleSet.driver_group_ids"
          group-by="group"
          :driver-groups="groupStore.groups"
        ></DriverSelector>
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
          {{ $t('commonBack') }}
        </UButton>

        <UButton type="submit" class="grow justify-center" color="secondary">
          {{ $t('commonSave') }}
        </UButton>
      </div>
    </div>
  </form>

  <MatchRuleInputModal
    ref="inputModal"
    @submit="
      rule => {
        const { _id, ...input } = rule
        if (_id !== undefined) {
          ruleSet.rules[_id] = input
        } else {
          ruleSet.rules.push(rule)
        }
      }
    "
  ></MatchRuleInputModal>
</template>

<style scoped>
legend:has(+ input:required, + select:required):after,
legend:has(+ div > input:required):after {
  content: ' *';
  color: red;
}
</style>
