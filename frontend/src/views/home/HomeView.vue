<script setup lang="ts">
import { useAppSettingStore, useDriverGroupStore, useMatchRuleStore } from '@/store'
import * as utils from '@/utils'
import CommandStatueModal from '@/views/home/components/CommandStatusModal.vue'
import * as executor from '@/wailsjs/go/execute/CommandExecutor'
import { storage } from '@/wailsjs/go/models'
import { computed, onBeforeMount, ref, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'vue-toast-notification'
import type { Command } from './types'

const { t } = useI18n()

const $toast = useToast({ position: 'top-right' })

const [statusModal, form] = [useTemplateRef('statusModal'), useTemplateRef('form')]

const [groupStore, settingStore, ruleStore] = [
  useDriverGroupStore(),
  useAppSettingStore(),
  useMatchRuleStore()
]

/** Driver groups with conditional exclusion based on app configuration */
const groups = computed(() =>
  settingStore.settings.hide_not_found
    ? groupStore.groups.filter(g => groupStore.isAllDriversExist(g))
    : groupStore.groups
)

const hwinfos = ref<{
  cpu: Array<string>
  gpu: Array<string>
  motherboard: Array<string>
  memory: Array<string>
  nic: Array<string>
  storage: Array<string>
} | null>(null)

onBeforeMount(() => {
  utils.getHardware().then(v => (hwinfos.value = v))
})

function selectMatchedOptions() {
  if (hwinfos.value === null) {
    $toast.error('沒有資訊')
    return
  }

  /** Determines if there is any hardware name matches the provided rule */
  const nameTest = (rule: storage.Rule): boolean => {
    return hwinfos.value![rule.source].some(src => utils.testMatchRule(rule, src))
  }

  ruleStore.ruleSets.forEach(rs => {
    const matched = rs.should_hit_all
      ? rs.rules.map(nameTest).every(Boolean)
      : rs.rules.map(nameTest).some(Boolean)

    if (form.value && matched) {
      rs.driver_group_ids.forEach(gid => {
        const el = form.value!.querySelector(`input[value="${gid}"], option[value="${gid}"]`)
        if (el instanceof HTMLInputElement) {
          el.checked = true
        } else if (el instanceof HTMLOptionElement) {
          el.selected = true
        }
      })
    }
  })
}

async function handleSubmit() {
  if (!form.value) {
    $toast.error(t('toast.readInputFailed'))
    return
  }

  const inputs = new FormData(form.value)
  const commands: Array<Command> = []

  if (settingStore.settings.set_password) {
    commands.push({
      id: 'set_password',
      groupName: t('task.setPassword'),
      config: {
        program: 'powershell',
        options: [
          '-WindowStyle',
          'Hidden',
          '-Command',
          `Set-LocalUser -Name $Env:UserName -Password ${
            settingStore.settings.password == ''
              ? '(new-object System.Security.SecureString)'
              : `(ConvertTo-SecureString ${settingStore.settings.password} -AsPlainText -Force)`
          }`
        ],
        minExeTime: 0.5,
        allowRtCodes: [0],
        incompatibles: []
      }
    })
  }

  if (settingStore.settings.create_partition) {
    commands.push({
      id: 'create_partition',
      groupName: t('task.createPartitions'),
      config: {
        program: 'powershell',
        options: [
          '-WindowStyle',
          'Hidden',
          '-Command',
          'Get-Disk | Where-Object PartitionStyle -Eq "RAW" | Initialize-Disk -PassThru | New-Partition -AssignDriveLetter -UseMaximumSize | Format-Volume'
        ],
        minExeTime: 1,
        allowRtCodes: [0],
        incompatibles: []
      }
    })
  }

  groupStore.groups
    .filter(group =>
      [inputs.get('network'), inputs.get('display'), ...inputs.getAll('miscellaneous')].includes(
        group.id
      )
    )
    .forEach(group => {
      group.drivers.forEach(driver => {
        commands.push({
          id: driver.id,
          name: driver.name,
          groupName: group.name,
          config: {
            program: driver.path,
            options: driver.flags,
            minExeTime: driver.minExeTime,
            allowRtCodes: driver.allowRtCodes,
            incompatibles: driver.incompatibles
          }
        })
      })
    })

  if (commands.length == 0) {
    $toast.warning(t('toast.noInputWarning'))
    return
  }

  statusModal.value?.show(settingStore.settings.parallel_install, commands)
}
</script>

<template>
  <div class="flex h-full flex-col">
    <div id="sysinfo" class="flex flex-1 flex-col gap-y-1 overflow-y-auto rounded-sm border p-1">
      <template v-if="hwinfos !== null">
        <div v-for="[part, names] in Object.entries(hwinfos)" :key="part">
          <h2 class="text-sm font-bold">{{ $t(`common.${part}`) }}</h2>

          <p
            v-for="(name, i) in names.filter(
              n =>
                part !== 'nic' ||
                ((!settingStore.settings.filter_miniport_nic || !n.includes('Miniport')) &&
                  (!settingStore.settings.filter_microsoft_nic || !n.includes('Microsoft')))
            )"
            :key="i"
            class="text-sm"
          >
            {{ name }}
          </p>
        </div>
      </template>

      <template v-else>
        <div v-for="i in 6" :key="i">
          <h2
            class="mb-1 h-5 skeleton"
            :style="{ width: `${Math.random() * (25 - 15) + 15}%` }"
          ></h2>
          <p class="h-5 skeleton" :style="{ width: `${Math.random() * (85 - 30) + 30}%` }"></p>
        </div>
      </template>
    </div>

    <form ref="form" class="mt-3 flex h-28 gap-x-3">
      <div class="flex flex-1 flex-col justify-between">
        <div class="relative w-full">
          <label
            class="pointer-events-none absolute start-4 top-0 h-full translate-y-1 text-xs text-gray-500"
          >
            {{ $t('driverCatetory.network') }}
          </label>

          <select name="network" class="w-full rounded-lg ps-3 pe-9 pt-5 pb-1">
            <option>{{ $t('common.pleaseSelect') }}</option>
            <option v-for="g in groups.filter(g => g.type == 'network')" :key="g.id" :value="g.id">
              {{ g.name }}
            </option>
          </select>
        </div>

        <div class="relative w-full">
          <label
            class="pointer-events-none absolute start-4 top-0 h-full translate-y-1 text-xs text-gray-500"
          >
            {{ $t('driverCatetory.display') }}
          </label>

          <select name="display" class="w-full rounded-lg ps-3 pe-9 pt-5 pb-1">
            <option>{{ $t('common.pleaseSelect') }}</option>
            <option v-for="g in groups.filter(g => g.type == 'display')" :key="g.id" :value="g.id">
              {{ g.name }}
            </option>
          </select>
        </div>
      </div>

      <div class="flex flex-1">
        <div class="relative mb-3 h-full w-full">
          <label
            class="pointer-events-none absolute top-1 left-3 origin-[0_0] -translate-y-[0.55rem] scale-[0.9] bg-white px-2 text-xs text-gray-500"
          >
            {{ $t('driverCatetory.miscellaneous') }}
          </label>

          <div class="h-full overflow-y-scroll rounded-lg border border-apple-green-600 px-2 pt-3">
            <template v-for="g in groups.filter(g => g.type == 'miscellaneous')" :key="g.id">
              <label class="flex w-full cursor-pointer items-center select-none">
                <input
                  type="checkbox"
                  name="miscellaneous"
                  class="checkbox me-1.5 checkbox-sm checkbox-primary"
                  :value="g.id"
                />
                {{ g.name }}
              </label>
            </template>
          </div>
        </div>
      </div>
    </form>

    <hr class="my-3" />

    <div class="flex gap-x-6">
      <div class="flex flex-col">
        <p class="font-semibold">{{ $t('installSetting.title') }}</p>

        <div class="flex flex-1 flex-col justify-around">
          <div class="flex gap-x-4">
            <label class="flex cursor-pointer items-center gap-x-1.5 select-none">
              <input
                v-model="settingStore.settings.create_partition"
                type="checkbox"
                name="create_partition"
                class="checkbox checkbox-sm checkbox-primary"
              />
              {{ $t('installSetting.createPartition') }}
            </label>

            <label class="flex cursor-pointer items-center gap-x-1.5 select-none">
              <input
                v-model="settingStore.settings.parallel_install"
                type="checkbox"
                name="parallel_install"
                class="checkbox checkbox-sm checkbox-primary"
              />
              {{ $t('installSetting.parallelInstall') }}
            </label>
          </div>

          <div class="flex gap-x-2">
            <label class="flex cursor-pointer items-center gap-x-1.5 select-none">
              <input
                v-model="settingStore.settings.set_password"
                type="checkbox"
                name="set_password"
                class="checkbox checkbox-sm checkbox-primary"
              />
              {{ $t('installSetting.setPassword') }}
            </label>

            <input
              v-model="settingStore.settings.password"
              type="text"
              name="password"
              class="input input-sm max-w-28 input-accent"
              :disabled="!settingStore.settings.set_password"
            />
          </div>
        </div>
      </div>

      <div class="flex grow flex-col justify-between">
        <div>
          <label class="mb-1 block text-sm text-gray-900">
            {{ $t('installSetting.successAction') }}
          </label>

          <select
            v-model="settingStore.settings.success_action"
            name="success_action"
            class="select w-full select-accent"
          >
            <option v-for="action in storage.SuccessAction" :key="action" :value="action">
              {{ $t(`successAction.${action}`) }}
            </option>
          </select>
        </div>

        <div class="mt-2 flex h-8 flex-row items-center justify-end gap-x-3">
          <button
            type="button"
            class="btn border-2 btn-outline btn-neutral"
            @click="selectMatchedOptions"
          >
            {{ $t('matchRule.match') }}
          </button>

          <button
            type="button"
            class="btn border-2 btn-outline btn-secondary"
            @click="
              () => {
                form?.reset()
                // settingStore.restore()
              }
            "
          >
            {{ $t('installSetting.reset') }}
          </button>

          <button class="btn btn-secondary" @click="handleSubmit">
            {{ $t('installSetting.execute') }}
          </button>
        </div>
      </div>
    </div>
  </div>

  <CommandStatueModal
    ref="statusModal"
    @completed="
      () => {
        switch (settingStore.settings.success_action) {
          case 'shutdown':
            executor.RunAndOutput(
              'cmd',
              ['/C', `shutdown /s /t ${settingStore.settings.success_action_delay}`],
              true
            )
            break
          case 'reboot':
            executor.RunAndOutput(
              'cmd',
              ['/C', `shutdown /r /t ${settingStore.settings.success_action_delay}`],
              true
            )
            break
          case 'firmware':
            executor
              .RunAndOutput(
                'cmd',
                ['/C', `shutdown /r /fw /t ${settingStore.settings.success_action_delay}`],
                true
              )
              .then(result => {
                if (result.exitCode != 0) {
                  // sometimes, /fw would resulted in an error: 'The system could not find the environment option that was entered. (203)'
                  // execute again normally solve the error
                  executor.RunAndOutput(
                    'cmd',
                    ['/C', `shutdown /r /fw /t ${settingStore.settings.success_action_delay}`],
                    true
                  )
                }
              })
            break
        }
      }
    "
  ></CommandStatueModal>
</template>
