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
    $toast.warning(t('toast.addAtLeastOneDriver'))
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
      <h1 class="mb-2 text-xl font-bold">建立配對規則</h1>

      <fieldset class="fieldset">
        <legend class="fieldset-legend text-sm">規則</legend>

        <div>
          <div class="grid grid-rows text-sm">
            <div class="grid grid-cols-10 gap-2 py-1.5 border-y">
              <div class="col-span-2">來源</div>
              <div class="col-span-3">類別</div>
              <div class="col-span-2">值</div>
              <div class="col-span-2">動作</div>
            </div>

            <div class="h-16 max-h-[23vh] overflow-y-auto border-b">
              <div v-if="ruleSet.rules.length == 0" class="h-full content-center py-1 text-center">
                N/A
              </div>

              <div
                v-else
                v-for="(r, i) in ruleSet.rules"
                :key="i"
                class="grid grid-cols-10 items-center gap-2 py-1.5 text-xs border-b"
              >
                <div class="col-span-2">
                  <p class="break-all line-clamp-2">{{ r.source }}</p>
                </div>

                <div class="col-span-3">
                  <p class="font-mono break-all line-clamp-2">
                    <span class="bg-gray-200 rounded-sm">
                      <font-awesome-icon v-if="r.is_case_sensitive" icon="fa-solid fa-a" />
                    </span>
                    {{ r.type }}
                  </p>
                </div>

                <div class="col-span-2">
                  <p class="break-all line-clamp-2">{{ r.values }}</p>
                </div>

                <div class="flex col-span-2 gap-x-1">
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

            <p class="text-hint">規則幫助</p>
          </div>

          <div class="flex justify-end gap-x-3">
            <button type="button" class="btn btn-primary px-2" @click="inputModal?.show()">
              <font-awesome-icon icon="fa-regular fa-square-plus" />
            </button>
          </div>
        </div>
      </fieldset>

      <hr />

      <fieldset class="fieldset">
        <legend class="fieldset-legend text-sm">目標</legend>

        <DriverSelector
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
legend:has(+ input:required, + select:required):after,
legend:has(+ div > input:required):after {
  content: ' *';
  color: red;
}
</style>
