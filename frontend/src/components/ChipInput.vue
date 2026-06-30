<script setup lang="ts">
import { ref } from 'vue'

const props = withDefaults(
  defineProps<{
    placeholder?: string
    commitKeys?: Array<'enter' | 'space' | 'tab' | 'comma'>
    parse?: (raw: string) => string
    accept?: (parsed: string) => boolean
    popOnBackspace?: boolean
    minInputWidth?: string
  }>(),
  {
    placeholder: '',
    commitKeys: () => ['enter'],
    parse: undefined,
    accept: undefined,
    popOnBackspace: true,
    minInputWidth: '120px'
  }
)

const model = defineModel<Array<string>>({ required: true })

const input = ref('')

function commit() {
  const parsed = (props.parse ?? ((s: string) => s))(input.value)
  if (!parsed) return
  if (!(props.accept?.(parsed) ?? true)) return
  if (!model.value.includes(parsed)) model.value.push(parsed)
  input.value = ''
}

function remove(i: number) {
  model.value.splice(i, 1)
}

function onBackspace() {
  if (props.popOnBackspace && input.value === '' && model.value.length > 0) {
    model.value.pop()
  }
}

function onKeydown(event: KeyboardEvent) {
  const key = event.key.toLowerCase()
  if ((props.commitKeys ?? ['enter']).includes(key as 'enter' | 'space' | 'tab' | 'comma')) {
    event.preventDefault()
    commit()
  }
}
</script>

<template>
  <div
    class="flex flex-wrap items-center gap-1.5 rounded-lg border border-gray-200 bg-gray-50/50 px-2 py-1.5 shadow-inner transition-all focus-within:bg-white"
    :style="{
      '--chip-bg': 'var(--color-half-baked-50)',
      '--chip-border': 'var(--color-half-baked-200)',
      '--chip-text': 'var(--color-half-baked-700)',
      '--chip-close-text': 'var(--color-half-baked-400)',
      '--chip-close-bg-hover': 'var(--color-half-baked-100)',
      '--focus-border': 'var(--color-half-baked-500)',
      '--focus-ring': 'color-mix(in oklab, var(--color-half-baked-500) 30%, transparent)'
    }"
  >
    <span
      v-for="(item, i) in model"
      :key="`${item}-${i}`"
      class="inline-flex items-center gap-1 rounded-md border px-2 py-0.5 text-xs font-bold shadow-sm"
      :style="{
        backgroundColor: 'var(--chip-bg)',
        borderColor: 'var(--chip-border)',
        color: 'var(--chip-text)'
      }"
    >
      <span class="font-mono">{{ item }}</span>

      <button
        type="button"
        class="rounded-sm p-0.5 transition-colors"
        :style="{
          color: 'var(--chip-close-text)',
          '--hover-bg': 'var(--chip-close-bg-hover)',
          '--hover-text': 'var(--chip-text)'
        }"
        @click="remove(i)"
      >
        <Icon icon="mdi:close" class="h-3 w-3" />
      </button>
    </span>

    <input
      v-model="input"
      type="text"
      :placeholder="placeholder"
      class="border-none bg-transparent px-2 py-1 font-mono text-xs text-gray-800 placeholder:text-gray-400 focus:ring-0 focus:outline-none xl:text-sm"
      :style="{ minWidth: minInputWidth, flex: '1 1 0%' }"
      @keydown="onKeydown"
      @keydown.backspace="onBackspace"
    />

    <slot />
  </div>
</template>

<style scoped>
div:focus-within {
  border-color: var(--focus-border);
  box-shadow: 0 0 0 2px var(--focus-ring);
}

button:hover {
  background-color: var(--hover-bg);
  color: var(--hover-text);
}
</style>
