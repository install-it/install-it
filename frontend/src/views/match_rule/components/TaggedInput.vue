<script setup lang="ts">
import { ref } from 'vue'

defineProps<{ title?: string }>()

const input = ref<string>('')

const tags = defineModel<Array<string>>({ default: [] })

function pushTag() {
  if (input.value.trim()) {
    tags.value.push(input.value.trim())
    input.value = ''
  }
}

function removeTag(index: number) {
  tags.value.splice(index, 1)
}
</script>

<template>
  <div class="flex flex-wrap gap-1 p-1">
    <button
      v-for="(tag, i) in tags"
      :key="tag"
      class="btn btn-sm"
      type="button"
      @click="removeTag(i)"
    >
      {{ tag }}
      <font-awesome-icon v-if="i < tags.length - 1" icon="fa-solid fa-xmark" class="h-6 w-6" />
      <font-awesome-icon v-else icon="fa-solid fa-delete-left" class="h-6 w-6" />
    </button>

    <input
      v-model="input"
      type="text"
      class="input input-sm grow input-accent focus:outline-0"
      @keydown.backspace="
        () => {
          if (input.length == 0) {
            removeTag(tags.length - 1)
          }
        }
      "
      @keydown.enter="
        event => {
          if (input != '') {
            event.preventDefault()
          }
          pushTag()
        }
      "
      @focusout="pushTag"
    />
  </div>
</template>
