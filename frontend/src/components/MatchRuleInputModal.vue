<script setup lang="ts">
import { storage } from '@/wailsjs/go/models'
import { computed, ref, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'

defineEmits<{ submit: [rules: { _id: number | undefined } & storage.Rule] }>()

const { t } = useI18n()

const toast = useToast()

const isOpen = ref(false)

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

    isOpen.value = true
  },
  hide: () => {
    isOpen.value = false
  }
})

const modalBody = useTemplateRef('modalBody')

const input = ref<{ _id: number | undefined } & storage.Rule>({
  _id: undefined,
  source: storage.RuleSource.CPU,
  operator: storage.RuleOperator.CONTAIN,
  is_case_sensitive: false,
  should_hit_all: false,
  values: []
})

const sourceItems = computed(() =>
  Object.entries(storage.RuleSource).map(([, value]) => ({     label: t(`common${value.charAt(0).toUpperCase() + value.slice(1)}`), value }))
)

const operatorItems = computed(() =>
  Object.entries(storage.RuleOperator).map(([, value]) => ({
    label: t(`matchRule${value.charAt(0).toUpperCase() + value.slice(1)}`),
    value
  }))
)
</script>

<template>
  <UModal v-model:open="isOpen" :title="$t('matchRuleMatchRule')">
    <template #body>
      <form
        autocomplete="off"
        @submit.prevent="
          () => {
            if (input.values.length == 0) {
              toast.add({ title: t('toastAddAtLeastOnePattern'), color: 'warning' })
            } else {
              $emit('submit', input)
              isOpen = false
            }
          }
        "
      >
        <div ref="modalBody" class="flex flex-col gap-y-2">
          <div class="flex gap-1">
            <fieldset class="fieldset flex-1">
              <legend class="fieldset-legend text-sm">
                {{ $t('matchRuleSource') }}
              </legend>

              <USelect
                v-model="input.source"
                color="primary"
                class="w-full"
                :items="sourceItems"
                required
              />
            </fieldset>

            <fieldset class="fieldset flex-1">
              <legend class="fieldset-legend text-sm">
                {{ $t('matchRuleOperator') }}
              </legend>

              <USelect
                v-model="input.operator"
                color="primary"
                class="w-full"
                :items="operatorItems"
                required
              />
            </fieldset>
          </div>

          <div class="flex">
            <fieldset class="fieldset flex-1">
              <legend class="fieldset-legend text-sm">
                {{ $t('matchRuleCaseSensitive') }}
              </legend>

              <label class="flex cursor-pointer items-center select-none">
                <UCheckbox
                  v-model="input.is_case_sensitive"
                  color="primary"
                  size="sm"
                  :disabled="input.operator === 'regex'"
                />

                <span class="ms-1.5">{{ $t('commonEnable') }}</span>
              </label>
            </fieldset>

            <fieldset class="fieldset flex-1">
              <legend class="fieldset-legend text-sm">
                {{ $t('matchRuleMultiPatternMatching') }}
              </legend>

              <label class="flex cursor-pointer items-center select-none">
                <UCheckbox v-model="input.should_hit_all" color="primary" size="sm" />

                <span class="ms-1.5">{{ $t('matchRuleHitAllPatterns') }}</span>
              </label>

              <p class="text-hint">{{ $t('matchRuleMultiPatternMatchingHelp') }}</p>
            </fieldset>
          </div>

          <fieldset class="input-bordered w-full rounded border py-1">
            <legend class="text-required fieldset-legend text-sm">
              {{ $t('matchRulePattern') }}
            </legend>

            <TaggedInput v-model="input.values"></TaggedInput>
          </fieldset>
        </div>

        <div class="flex gap-x-2 border-t pt-2">
          <UButton type="submit" color="secondary" size="sm" block class="justify-center">
            {{ $t('commonSave') }}
          </UButton>
        </div>
      </form>
    </template>
  </UModal>
</template>

<style scoped>
legend:has(+ input:required, + select:required):after,
legend:has(+ div > input:required):after {
  content: ' *';
  color: red;
}
</style>
