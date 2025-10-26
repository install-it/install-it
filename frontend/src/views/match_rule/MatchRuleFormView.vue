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
    <div class="h-full content-center space-y-3 max-w-full lg:max-w-2xl xl:max-w-4xl mx-auto">
      <h1 class="mb-2 text-xl font-bold">
        {{ $t('matchRule.createRule') }}
      </h1>

      <fieldset class="fieldset">
        <legend class="fieldset-legend text-sm">
          {{ $t('matchRule.name') }}
        </legend>
        <input type="text" class="input input-accent ms-1" v-model="ruleSet.name" />
      </fieldset>

      <fieldset class="fieldset">
        <legend class="fieldset-legend text-sm text-required">
          {{ $t('matchRule.rule') }}
        </legend>

        <div>
          <div class="grid grid-rows text-sm">
            <div class="grid grid-cols-6 gap-2 py-1.5 border-y">
              <div class="col-span-1">{{ $t('matchRule.source') }}</div>
              <div class="col-span-1">{{ $t('matchRule.operator') }}</div>
              <div class="col-span-3">{{ $t('matchRule.pattern') }}</div>
              <div class="col-span-1">{{ $t('common.action') }}</div>
            </div>

            <div class="max-h-38 overflow-y-auto">
              <div
                v-if="ruleSet.rules.length == 0"
                class="h-full content-center py-1 text-center border-b"
              >
                N/A
              </div>

              <div
                v-else
                v-for="(r, i) in ruleSet.rules"
                :key="i"
                class="grid grid-cols-6 items-center gap-2 py-1.5 text-xs border-b"
              >
                <div class="col-span-1">
                  <p class="break-all line-clamp-2">
                    {{ $t(`common.${r.source}`) }}
                  </p>
                </div>

                <div class="col-span-1">
                  <div class="flex gap-0.5">
                    <span v-if="r.is_case_sensitive">
                      <font-awesome-icon
                        icon="fa-solid fa-a"
                        class="bg-gray-200 rounded-sm p-0.5"
                      />
                    </span>
                    <span class="font-mono">
                      {{ $t(`matchRule.${r.operator}`) }}
                    </span>
                  </div>
                </div>

                <div class="col-span-3 space-x-0.5 md:space-x-1.5 space-y-0.5">
                  <span
                    v-for="(v, i) in r.values"
                    :key="i"
                    class="badge badge-neutral badge-sm px-0.5"
                  >
                    {{ v }}
                  </span>
                </div>

                <div class="flex col-span-1 gap-x-1">
                  <div class="flex gap-x-2">
                    <button
                      type="button"
                      @click="inputModal?.show(JSON.parse(JSON.stringify({ ...r, _id: i })))"
                    >
                      <font-awesome-icon icon="fa-solid fa-pen-to-square" />
                    </button>
                    <button type="button" @click="ruleSet.rules.splice(i, 1)">
                      <font-awesome-icon icon="fa-solid fa-trash" />
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <p class="text-hint"></p>
          </div>

          <div class="flex justify-end gap-x-3">
            <button type="button" class="btn btn-primary px-2" @click="inputModal?.show()">
              <font-awesome-icon icon="fa-regular fa-square-plus" />
            </button>
          </div>
        </div>
      </fieldset>

      <fieldset class="fieldset">
        <legend class="fieldset-legend text-sm">{{ $t('matchRule.multiRuleMatching') }}</legend>

        <label class="flex items-center w-full select-none cursor-pointer">
          <input
            type="checkbox"
            class="me-1.5 checkbox checkbox-primary"
            v-model="ruleSet.should_hit_all"
          />
          {{ $t('matchRule.hitAllRule') }}
        </label>

        <p class="text-hint">{{ $t('matchRule.multiRuleMatchingHelp') }}</p>
      </fieldset>

      <hr />

      <fieldset class="fieldset">
        <legend class="fieldset-legend text-sm">
          {{ $t('matchRule.matchTo') }}
        </legend>

        <DriverSelector
          group-by="group"
          :driver-groups="groupStore.groups"
          v-model="ruleSet.driver_group_ids"
        ></DriverSelector>
      </fieldset>

      <div class="flex h-8 gap-x-5">
        <button
          type="button"
          class="grow btn"
          style="--btn-color: var(--color-gray-100)"
          @click="$router.back()"
        >
          {{ $t('common.back') }}
        </button>

        <button type="submit" class="grow btn btn-secondary">
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
.text-required:after,
legend:has(+ input:required, + select:required):after,
legend:has(+ div > input:required):after {
  content: ' *';
  color: red;
}
</style>
