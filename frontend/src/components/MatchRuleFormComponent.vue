<script setup lang="ts">
import DriverSelector from '@/components/DriverSelector.vue'
import MatchRuleInputModal from '@/components/MatchRuleInputModal.vue'
import { useEditor } from '@/composables/useEditor'
import { useDriverGroupStore, useMatchRuleStore } from '@/store'
import { storage } from '@/wailsjs/go/models'
import * as matchRuleStorage from '@/wailsjs/go/storage/MatchRuleStorage'
import { computed, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const props = defineProps<{
  id?: string
}>()

const inputModal = useTemplateRef('inputModal')

const { t } = useI18n()

const $router = useRouter()

const toast = useToast()

const [ruleStore, groupStore] = [useMatchRuleStore(), useDriverGroupStore()]

const sourceRuleSet = computed(
  () =>
    ruleStore.ruleSets.find(r => r.id === props.id) ??
    new storage.RuleSet({ rules: [], driver_group_ids: [] })
)

const { data: ruleSet } = useEditor({ source: sourceRuleSet.value })

function handleSubmit() {
  if (ruleSet.value.rules.length == 0) {
    toast.add({ title: t('toast.addAtLeastOneRule'), color: 'warning' })
    return
  }

  const handleSuccess = () => {
    toast.add({ title: t('toast.updated'), color: 'success' })

    matchRuleStorage.All().then(newMatchRule => {
      ruleStore.ruleSets = newMatchRule
      $router.push({ path: '/match-rules/' })
    })
  }

  if (ruleSet.value.id == undefined) {
    matchRuleStorage
      .Add(ruleSet.value)
      .then(rid => (ruleSet.value.id = rid))
      .then(handleSuccess)
      .catch(reason => toast.add({ title: reason.toString(), color: 'error' }))
  } else {
    matchRuleStorage
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
        {{ $t('matchRule.createRule') }}
      </h1>

      <fieldset class="fieldset">
        <legend class="fieldset-legend text-sm">
          {{ $t('matchRule.name') }}
        </legend>

        <UInput v-model="ruleSet.name" type="text" class="ms-1" />
      </fieldset>

      <fieldset class="fieldset">
        <legend class="text-required fieldset-legend text-sm">
          {{ $t('matchRule.rule') }}
        </legend>

        <div>
          <div class="grid-rows grid text-sm">
            <div class="grid grid-cols-6 gap-2 border-y py-1.5">
              <div class="col-span-1">{{ $t('matchRule.source') }}</div>

              <div class="col-span-1">{{ $t('matchRule.operator') }}</div>

              <div class="col-span-3">{{ $t('matchRule.pattern') }}</div>

              <div class="col-span-1">{{ $t('common.action') }}</div>
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
                    {{ $t(`common.${r.source}`) }}
                  </p>
                </div>

                <div class="col-span-1">
                  <span class="font-mono">
                    {{ $t(`matchRule.${r.operator}`) }}
                  </span>
                </div>

                <div class="col-span-3">
                  <div>
                    <span
                      v-if="r.should_hit_all"
                      class="badge badge-sm me-0.5 px-1 text-white md:me-1"
                      style="--badge-color: var(--color-rose-400)"
                    >
                      {{ $t('matchRule.hitAll') }}
                    </span>

                    <span
                      v-if="r.is_case_sensitive"
                      class="badge badge-sm me-0.5 px-1 md:me-1"
                      style="--badge-color: var(--color-orange-300)"
                    >
                      Aa
                    </span>

                    <span
                      v-for="(v, vi) in r.values"
                      :key="vi"
                      class="badge badge-sm badge-neutral me-0.5 px-1 md:me-1"
                    >
                      {{ v }}
                    </span>
                  </div>
                </div>

                <div class="col-span-1 flex gap-x-1">
                  <div class="flex gap-x-2">
                    <button
                      type="button"
                      :title="$t('common.edit')"
                      @click="inputModal?.show(JSON.parse(JSON.stringify({ ...r, _id: i })))"
                    >
                      <Icon icon="mdi:pencil" class="size-4" />
                    </button>

                    <button
                      type="button"
                      :title="$t('common.delete')"
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
        <legend class="fieldset-legend text-sm">{{ $t('matchRule.multiRuleMatching') }}</legend>

        <label class="flex w-full cursor-pointer items-center select-none">
          <UCheckbox v-model="ruleSet.should_hit_all" class="me-1.5" />
          {{ $t('matchRule.hitAllRule') }}
        </label>

        <p class="text-hint">{{ $t('matchRule.multiRuleMatchingHelp') }}</p>
      </fieldset>

      <hr />

      <fieldset class="fieldset">
        <legend class="text-required fieldset-legend text-sm">
          {{ $t('matchRule.matchTo') }}
        </legend>

        <DriverSelector
          v-model="ruleSet.driver_group_ids"
          group-by="group"
          :driver-groups="groupStore.groups"
          :exclude-builtin="true"
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
          {{ $t('common.back') }}
        </UButton>

        <UButton type="submit" class="grow justify-center" color="secondary">
          {{ $t('common.save') }}
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
