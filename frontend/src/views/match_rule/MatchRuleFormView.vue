<script setup lang="ts">
import DriverSelector from '@/components/DriverSelector.vue'
import { useDriverGroupStore, useMatchRuleStore } from '@/store'
import * as matchRuleStorage from '@/wailsjs/go/storage/MatchRuleStorage'
import { useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'vue-toast-notification'
import MatchRuleInputModal from './components/MatchRuleInputModal.vue'

const inputModal = useTemplateRef('inputModal')

const { t } = useI18n()

const [$route, $router] = [useRoute(), useRouter()]

const $toast = useToast({ position: 'top-right' })

const [ruleStore, groupStore] = [useMatchRuleStore(), useDriverGroupStore()]

const ruleEditor = ruleStore.editor($route.params.id as string | undefined)

const ruleSet = ruleEditor.ruleSet // alias

function handleSubmit(event: SubmitEvent) {
  if (ruleSet.value.rules.length == 0) {
    $toast.warning(t('toast.addAtLeastOneRule'))
    return
  }

  const handleSuccess = () => {
    $toast.success(t('toast.updated'))

    matchRuleStorage.All().then(newMatchRule => {
      ruleStore.ruleSets = newMatchRule
      ruleEditor.reset()

      if (event.submitter?.id !== 'driver-submit-btn') {
        $router.back()
      } else {
        $router.replace({ path: `/` })
      }
    })
  }

  if (ruleSet.value.id == undefined) {
    matchRuleStorage
      .Add(ruleSet.value)
      .then(rid => (ruleSet.value.id = rid))
      .then(handleSuccess)
      .catch(reason => $toast.error(reason.toString()))
  } else {
    matchRuleStorage
      .Update(ruleSet.value)
      .then(handleSuccess)
      .catch(reason => $toast.error(reason.toString()))
  }
}
</script>

<template>
  <form
    class="size-full overflow-y-auto"
    autocomplete="off"
    @submit.prevent="event => handleSubmit(event as SubmitEvent)"
  >
    <div class="mx-auto h-full max-w-full content-center space-y-3 lg:max-w-2xl xl:max-w-4xl">
      <h1 class="mb-2 text-xl font-bold">
        {{ $t('matchRule.createRule') }}
      </h1>

      <fieldset class="fieldset">
        <legend class="fieldset-legend text-sm">
          {{ $t('matchRule.name') }}
        </legend>
        <input type="text" class="input ms-1 input-accent" v-model="ruleSet.name" />
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
                v-else
                v-for="(r, i) in ruleSet.rules"
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

                <div class="col-span-1 flex gap-x-1">
                  <div class="flex gap-x-2">
                    <button
                      type="button"
                      @click="inputModal?.show(JSON.parse(JSON.stringify({ ...r, _id: i })))"
                      :title="$t('common.edit')"
                    >
                      <font-awesome-icon icon="fa-solid fa-pen-to-square" />
                    </button>
                    <button
                      type="button"
                      @click="ruleSet.rules.splice(i, 1)"
                      :title="$t('common.delete')"
                    >
                      <font-awesome-icon icon="fa-solid fa-trash" />
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <p class="text-hint"></p>
          </div>

          <div class="flex justify-end gap-x-3">
            <button type="button" class="btn px-2 btn-primary" @click="inputModal?.show()">
              <font-awesome-icon icon="fa-regular fa-square-plus" />
            </button>
          </div>
        </div>
      </fieldset>

      <fieldset class="fieldset">
        <legend class="fieldset-legend text-sm">{{ $t('matchRule.multiRuleMatching') }}</legend>

        <label class="flex w-full cursor-pointer items-center select-none">
          <input
            type="checkbox"
            class="checkbox me-1.5 checkbox-primary"
            v-model="ruleSet.should_hit_all"
          />
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
          group-by="group"
          :driver-groups="groupStore.groups"
          :exclude-builtin="true"
          v-model="ruleSet.driver_group_ids"
        ></DriverSelector>
      </fieldset>

      <div class="flex h-8 gap-x-5">
        <button
          type="button"
          class="btn grow"
          style="--btn-color: var(--color-gray-100)"
          @click="$router.back()"
        >
          {{ $t('common.back') }}
        </button>

        <button type="submit" class="btn grow btn-secondary">
          {{ $t('common.save') }}
        </button>
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
