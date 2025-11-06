<script setup lang="ts">
import ModalFrame from '@/components/modals/ModalFrame.vue'
import { storage } from '@/wailsjs/go/models'
import { ref, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import TaggedInput from './TaggedInput.vue'

const { t } = useI18n()

defineExpose({
  show: (rule?: { _id: number | undefined } & storage.Rule) => {
    input.value = rule || {
      _id: undefined,
      source: storage.RuleSource.CPU,
      operator: storage.RuleOperator.CONTAIN,
      is_case_sensitive: false,
      should_hit_all: false,
      values: []
    }

    frame.value?.show()
  },
  hide: () => frame.value?.hide()
})

defineEmits<{ submit: [rules: { _id: number | undefined } & storage.Rule] }>()

const [frame, modalBody] = [useTemplateRef('frame'), useTemplateRef('modalBody')]

const input = ref<{ _id: number | undefined } & storage.Rule>({
  _id: undefined,
  source: storage.RuleSource.CPU,
  operator: storage.RuleOperator.CONTAIN,
  is_case_sensitive: false,
  should_hit_all: false,
  values: []
})
</script>

<template>
  <ModalFrame :on-demand="true" :immediate="false" ref="frame">
    <div class="w-4/5">
      <div class="rounded-lg bg-white shadow-sm">
        <div class="flex h-12 items-center justify-between rounded-t border-b px-4">
          <h3 class="font-semibold">
            {{ $t('matchRule.matchRule') }}
          </h3>

          <button
            type="button"
            class="rounded-lg bg-transparent p-3 text-sm text-gray-400 hover:text-gray-900"
            @click="frame?.hide()"
          >
            <font-awesome-icon icon="fa-solid fa-xmark" />
          </button>
        </div>

        <form
          class=""
          autocomplete="off"
          @submit.prevent="
            () => {
              if (input.values.length == 0) {
                $toast.warning(t('toast.addAtLeastOnePattern'))
              } else {
                $emit('submit', input)
                frame?.hide()
              }
            }
          "
        >
          <div class="flex max-h-[75vh] flex-col gap-y-2 overflow-auto p-4" ref="modalBody">
            <div class="flex gap-1">
              <fieldset class="fieldset flex-1">
                <legend class="fieldset-legend text-sm">
                  {{ $t('matchRule.source') }}
                </legend>
                <select v-model="input.source" class="select select-accent" required>
                  <option v-for="s in storage.RuleSource" :value="s" :key="s">
                    {{ $t(`common.${s}`) }}
                  </option>
                </select>
              </fieldset>

              <fieldset class="fieldset flex-1">
                <legend class="fieldset-legend text-sm">
                  {{ $t('matchRule.operator') }}
                </legend>
                <select v-model="input.operator" class="select select-accent" required>
                  <option v-for="t in storage.RuleOperator" :value="t" :key="t">
                    {{ $t(`matchRule.${t}`) }}
                  </option>
                </select>
              </fieldset>
            </div>

            <div class="flex">
              <fieldset class="fieldset flex-1">
                <legend class="fieldset-legend text-sm">
                  {{ $t('matchRule.caseSensitive') }}
                </legend>
                <label class="flex cursor-pointer items-center select-none">
                  <input
                    type="checkbox"
                    v-model="input.is_case_sensitive"
                    class="checkbox me-1.5 checkbox-sm checkbox-primary"
                    :disabled="input.operator === 'regex'"
                  />
                  {{ $t('common.enable') }}
                </label>
              </fieldset>

              <fieldset class="fieldset flex-1">
                <legend class="fieldset-legend text-sm">
                  {{ $t('matchRule.multiPatternMatching') }}
                </legend>
                <label class="flex cursor-pointer items-center select-none">
                  <input
                    type="checkbox"
                    v-model="input.should_hit_all"
                    class="checkbox me-1.5 checkbox-sm checkbox-primary"
                  />
                  {{ $t('matchRule.hitAllPatterns') }}
                </label>

                <p class="text-hint">{{ $t('matchRule.multiPatternMatchingHelp') }}</p>
              </fieldset>
            </div>

            <fieldset class="input-bordered w-full rounded border py-1">
              <legend class="text-required fieldset-legend text-sm">
                {{ $t('matchRule.pattern') }}
              </legend>
              <TaggedInput v-model="input.values"></TaggedInput>
            </fieldset>
          </div>

          <div class="flex gap-x-2 border-t px-4 py-2">
            <button type="submit" class="btn w-full btn-sm btn-secondary">
              {{ $t('common.save') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </ModalFrame>
</template>

<style scoped>
legend:has(+ input:required, + select:required):after,
legend:has(+ div > input:required):after {
  content: ' *';
  color: red;
}
</style>
