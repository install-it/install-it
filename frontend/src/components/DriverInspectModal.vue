<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  groupId: number | null
}>()

const emit = defineEmits<{
  close: []
  edit: [id: number]
}>()

const groupStore = useDriverGroupStore()

const isOpen = computed({
  get: () => props.groupId !== null,
  set: val => {
    if (!val) emit('close')
  }
})

const group = computed(() => groupStore.groups.find(g => g.id === props.groupId))

function categoryKey(type: string): string {
  return `category${type.charAt(0).toUpperCase() + type.slice(1)}`
}
</script>

<template>
  <UModal v-model:open="isOpen" :title="$t('titleInspectGroup')">
    <template #body>
      <div v-if="group" class="flex flex-col gap-y-4">
        <!-- Metrics strip: 3-col grid in white card -->
        <div
          class="grid grid-cols-3 gap-2 rounded-lg border border-gray-200 bg-white p-3 text-center shadow-sm"
        >
          <div>
            <span class="block text-[10px] font-bold text-gray-400 uppercase">Type</span>

            <span class="text-xs font-bold" :style="`color: var(--color-${group.type})`">{{
              $t(categoryKey(group.type))
            }}</span>
          </div>

          <div class="border-x border-gray-100">
            <span class="block text-[10px] font-bold text-gray-400 uppercase">Exclusive Flow</span>

            <span class="text-xs font-bold text-gray-800">{{
              group.mutuallyExclusive ? 'Yes' : 'No'
            }}</span>
          </div>

          <div>
            <span class="block text-[10px] font-bold text-gray-400 uppercase">Drivers</span>

            <span class="text-xs font-bold text-gray-800">{{ group.drivers.length }}</span>
          </div>
        </div>

        <!-- Per-driver cards -->
        <div
          v-for="(d, i) in group.drivers"
          :key="d.id"
          class="relative space-y-3 overflow-hidden rounded-lg border border-gray-200 bg-white p-4 shadow-sm"
        >
          <!-- Left edge accent -->
          <div
            class="absolute top-0 left-0 h-full w-1"
            :class="groupStore.notFoundDrivers.includes(d.id) ? 'bg-red-400' : 'bg-gray-200'"
          />

          <!-- Header -->
          <div class="flex items-center justify-between pl-2">
            <h4 class="flex items-center gap-2 text-sm font-bold text-gray-900">
              <span class="text-gray-400">#{{ i + 1 }}</span>
              {{ d.name }}
            </h4>
            <!-- Only show badge when path is MISSING, not when valid -->
            <span
              v-if="groupStore.notFoundDrivers.includes(d.id)"
              class="rounded border border-red-200 bg-red-100 px-2 py-0.5 text-[10px] font-bold text-red-700"
            >
              Missing Exe
            </span>
          </div>

          <!-- Path -->
          <div class="space-y-2 pl-2">
            <div>
              <span class="text-[9px] font-bold text-gray-400 uppercase">{{ $t('path') }}</span>

              <div
                class="rounded border border-gray-100 bg-gray-50 p-2 font-mono text-[11px] break-all text-gray-700"
              >
                {{ d.path }}
              </div>
            </div>

            <!-- Details grid -->
            <div class="grid grid-cols-2 gap-3 rounded border border-gray-100 bg-gray-50/50 p-2">
              <div>
                <span class="block text-[9px] font-bold text-gray-400 uppercase">{{
                  $t('fieldArgument')
                }}</span>

                <span class="mt-0.5 block font-mono text-[11px] font-bold text-gray-700">
                  {{ d.flags.length > 0 ? d.flags.join(' ') : 'None' }}
                </span>
              </div>

              <div>
                <span class="block text-[9px] font-bold text-gray-400 uppercase">{{
                  $t('fieldAllowedExitCode')
                }}</span>

                <span class="mt-0.5 block font-mono text-[11px] font-bold text-gray-700">
                  {{ d.allowRtCodes.length > 0 ? d.allowRtCodes.join(', ') : 'Any' }}
                </span>
              </div>

              <div>
                <span class="block text-[9px] font-bold text-gray-400 uppercase">{{
                  $t('fieldMinExecuteTime')
                }}</span>

                <span class="mt-0.5 block font-mono text-[11px] font-bold text-gray-700"
                  >{{ d.minExeTime }}s</span
                >
              </div>

              <div>
                <span class="block text-[9px] font-bold text-gray-400 uppercase">{{
                  $t('labelIncompatibleWith')
                }}</span>

                <span class="mt-0.5 block font-mono text-[11px] font-bold text-gray-700">
                  {{ d.incompatibles.length > 0 ? d.incompatibles.length + ' set' : 'None' }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="py-6 text-center text-sm text-gray-500">
        {{ $t('msgGroupNotFound') }}
      </div>
    </template>

    <template #footer>
      <div class="flex justify-end gap-2">
        <UButton color="neutral" variant="ghost" size="sm" @click="emit('close')">
          {{ $t('cancel') }}
        </UButton>

        <UButton v-if="group" color="primary" size="sm" @click="emit('edit', group.id)">
          {{ $t('edit') }}
        </UButton>
      </div>
    </template>
  </UModal>
</template>
