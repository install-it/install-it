<script setup lang="ts">
import ModalFrame from '@/components/modals/ModalFrame.vue'
import { storage } from '@/wailsjs/go/models'
import { ref, useTemplateRef } from 'vue'
import TaggedInput from './TaggedInput.vue'

defineExpose({
  show: (rule?: { _id: number | undefined } & storage.Rule) => {
    input.value = rule || {
      _id: undefined,
      source: storage.RuleSource.CPU,
      type: storage.RuleType.CONTAIN,
      is_case_sensitive: false,
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
  type: storage.RuleType.CONTAIN,
  is_case_sensitive: false,
  values: []
})
</script>

<template>
  <ModalFrame :on-demand="true" :immediate="false" ref="frame">
    <div class="w-4/5">
      <div class="bg-white rounded-lg shadow-sm">
        <div class="flex items-center justify-between h-12 px-4 border-b rounded-t">
          <h3 class="font-semibold">配對規則</h3>

          <button
            type="button"
            class="p-3 text-sm text-gray-400 hover:text-gray-900 bg-transparent rounded-lg"
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
              $emit('submit', input)
              frame?.hide()
            }
          "
        >
          <div class="flex flex-col gap-y-2 max-h-[75vh] overflow-auto p-4" ref="modalBody">
            <div class="flex gap-1">
              <fieldset class="fieldset flex-1">
                <legend class="fieldset-legend text-sm">規則</legend>
                <select v-model="input.source" class="select select-accent" required>
                  <option v-for="s in storage.RuleSource" :value="s" :key="s">{{ s }}</option>
                </select>
              </fieldset>

              <fieldset class="fieldset flex-1">
                <legend class="fieldset-legend text-sm">規則</legend>
                <select v-model="input.type" class="select select-accent" required>
                  <option v-for="t in storage.RuleType" :value="t" :key="t">{{ t }}</option>
                </select>
              </fieldset>
            </div>

            <fieldset class="fieldset flex-1">
              <legend class="fieldset-legend text-sm">區分大小寫</legend>
              <label class="flex items-center select-none cursor-pointer">
                <input
                  type="checkbox"
                  value="set_password"
                  v-model="input.is_case_sensitive"
                  class="checkbox checkbox-sm checkbox-primary me-1.5"
                />
                啟用
              </label>
            </fieldset>

            <!-- <fieldset class="fieldset flex-1">
              <legend class="fieldset-legend text-sm">區分大小寫</legend>
              <input type="checkbox" class="checkbox" v-model="input.is_case_sensitive" />
            </fieldset> -->

            <fieldset class="border input-bordered rounded w-full py-1">
              <legend class="fieldset-legend text-sm">值</legend>
              <TaggedInput v-model="input.values"></TaggedInput>
            </fieldset>
          </div>

          <div class="flex gap-x-2 px-4 py-2 border-t">
            <button type="submit" class="btn btn-sm btn-secondary w-full">
              {{ $t('common.save') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </ModalFrame>
</template>
