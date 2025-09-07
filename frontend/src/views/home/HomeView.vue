<script setup lang="ts">
import { useAppSettingStore, useDriverGroupStore } from '@/store'
import CommandStatueModal from '@/views/home/components/CommandStatusModal.vue'
import * as executor from '@/wailsjs/go/execute/CommandExecutor'
import { storage, sysinfo } from '@/wailsjs/go/models'
import * as sysinfoqy from '@/wailsjs/go/sysinfo/SysInfo'
import { onBeforeMount, ref, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'vue-toast-notification'
import type { Command } from './types'

const { t } = useI18n()

const $toast = useToast({ position: 'top-right' })

const statusModal = useTemplateRef('statusModal')

const form = useTemplateRef('form')

const settingStore = useAppSettingStore()

const groupStore = useDriverGroupStore()

const hwinfos = ref<{
  motherboard: Array<sysinfo.Win32_BaseBoard>
  cpu: Array<sysinfo.Win32_Processor>
  gpu: Array<sysinfo.Win32_VideoController>
  memory: Array<sysinfo.Win32_PhysicalMemory>
  nic: Array<sysinfo.Win32_NetworkAdapter>
  disk: Array<sysinfo.Win32_DiskDrive>
} | null>(null)

onBeforeMount(() => {
  Promise.all([
    sysinfoqy.MotherboardInfo(),
    sysinfoqy.CpuInfo(),
    sysinfoqy.GpuInfo(),
    sysinfoqy.MemoryInfo(),
    sysinfoqy.NicInfo(),
    sysinfoqy.DiskInfo()
  ]).then(infos => {
    hwinfos.value = ['motherboard', 'cpu', 'gpu', 'memory', 'nic', 'disk'].reduce(
      (obj, key, index) => {
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        obj[key] = infos[index]
        return obj
      },
      {} as typeof hwinfos.value
    )
  })
})

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
  <div class="flex flex-col h-full">
    <div id="sysinfo" class="flex flex-col flex-1 gap-y-1 overflow-y-auto p-1 border rounded-sm">
      <template v-if="hwinfos !== null">
        <div>
          <h2 class="text-sm font-bold">{{ $t('common.motherboard') }}</h2>

          <p v-for="(mb, i) in hwinfos.motherboard" :key="i" class="text-sm">
            {{ `${mb.Manufacturer} ${mb.Product}` }}
          </p>
        </div>

        <div>
          <h2 class="text-sm font-bold">{{ $t('common.cpu') }}</h2>

          <p v-for="(cpu, i) in hwinfos.cpu" :key="i" class="text-sm">
            {{ cpu.Name }}
          </p>
        </div>

        <div>
          <h2 class="text-sm font-bold">{{ $t('common.ram') }}</h2>

          <p v-for="(mem, i) in hwinfos.memory" :key="i" class="text-sm">
            {{
              `${mem.Manufacturer} ${mem.PartNumber.trim()} ${mem.Capacity / Math.pow(1024, 3)}GB ${mem.Speed}MHz`
            }}
          </p>
        </div>

        <div>
          <h2 class="text-sm font-bold">{{ $t('common.gpu') }}</h2>

          <p v-for="(dp, i) in hwinfos.gpu" :key="i" class="text-sm">
            {{ `${dp.Name} (${dp.AdapterRAM / Math.pow(1024, 3)}GB)` }}
          </p>
        </div>

        <div>
          <h2 class="text-sm font-bold">{{ $t('common.nic') }}</h2>

          <p
            v-for="(dp, i) in hwinfos.nic
              .filter(
                n =>
                  !settingStore.settings.filter_miniport_nic ||
                  (settingStore.settings.filter_miniport_nic && !n.Name.includes('Miniport'))
              )
              .filter(
                n =>
                  !settingStore.settings.filter_microsoft_nic ||
                  (settingStore.settings.filter_microsoft_nic && !n.Name.includes('Microsoft'))
              )"
            :key="i"
            class="text-sm"
          >
            {{ dp.Name }}
          </p>
        </div>

        <div>
          <h2 class="text-sm font-bold">{{ $t('common.storage') }}</h2>

          <p v-for="(dp, i) in hwinfos.disk" :key="i" class="text-sm">
            {{ `${dp.Model} (${Math.round(dp.Size / Math.pow(1024, 3))}GB)` }}
          </p>
        </div>
      </template>

      <template v-else>
        <div v-for="i in 6" :key="i">
          <h2
            class="skeleton h-5 mb-1"
            :style="{ width: `${Math.random() * (25 - 15) + 15}%` }"
          ></h2>
          <p class="skeleton h-5" :style="{ width: `${Math.random() * (85 - 30) + 30}%` }"></p>
        </div>
      </template>
    </div>

    <form class="flex gap-x-3 h-28 mt-3" ref="form">
      <div class="flex flex-col flex-1 justify-between">
        <div class="relative w-full">
          <label
            class="absolute top-0 start-4 h-full translate-y-1 text-xs text-gray-500 pointer-events-none"
          >
            {{ $t('driverCatetory.network') }}
          </label>

          <select name="network" class="w-full ps-3 pe-9 pt-5 pb-1 rounded-lg">
            <option>{{ $t('common.pleaseSelect') }}</option>
            <option
              v-for="d in groupStore.groups.filter(d => d.type == 'network')"
              :key="d.id"
              :value="d.id"
            >
              {{ `${d.name}${groupStore.notFoundDrivers.includes(d.id) ? ' ⚠' : ''}` }}
            </option>
          </select>
        </div>

        <div class="relative w-full">
          <label
            class="absolute top-0 start-4 h-full translate-y-1 text-xs text-gray-500 pointer-events-none"
          >
            {{ $t('driverCatetory.display') }}
          </label>

          <select name="display" class="w-full ps-3 pe-9 pt-5 pb-1 rounded-lg">
            <option>{{ $t('common.pleaseSelect') }}</option>
            <option
              v-for="d in groupStore.groups.filter(d => d.type == 'display')"
              :key="d.id"
              :value="d.id"
            >
              {{ `${d.name}${groupStore.notFoundDrivers.includes(d.id) ? ' ⚠' : ''}` }}
            </option>
          </select>
        </div>
      </div>

      <div class="flex flex-1">
        <div class="relative w-full h-full mb-3">
          <label
            class="absolute left-3 top-1 origin-[0_0] -translate-y-[0.55rem] px-2 bg-white text-xs scale-[0.9] text-gray-500 pointer-events-none"
          >
            {{ $t('driverCatetory.miscellaneous') }}
          </label>

          <div class="h-full overflow-y-scroll px-2 pt-3 rounded-lg border border-apple-green-600">
            <template
              v-for="d in groupStore.groups.filter(d => d.type == 'miscellaneous')"
              :key="d.id"
            >
              <label class="flex items-center w-full select-none cursor-pointer">
                <input
                  type="checkbox"
                  name="miscellaneous"
                  class="checkbox checkbox-sm checkbox-primary me-1.5"
                  :value="d.id"
                />
                {{ `${d.name}${groupStore.notFoundDrivers.includes(d.id) ? ' ⚠' : ''}` }}
              </label>
            </template>
          </div>
        </div>
      </div>
    </form>

    <hr class="my-3" />

    <div class="flex gap-x-6">
      <div class="flex flex-col">
        <p class="font-semibold">{{ $t('installOption.title') }}</p>

        <div class="flex flex-col flex-1 justify-around">
          <div class="flex gap-x-4">
            <label class="flex items-center gap-x-1.5 select-none cursor-pointer">
              <input
                type="checkbox"
                name="create_partition"
                v-model="settingStore.settings.create_partition"
                class="checkbox checkbox-sm checkbox-primary"
              />
              {{ $t('installOption.createPartition') }}
            </label>

            <label class="flex items-center gap-x-1.5 select-none cursor-pointer">
              <input
                type="checkbox"
                name="parallel_install"
                v-model="settingStore.settings.parallel_install"
                class="checkbox checkbox-sm checkbox-primary"
              />
              {{ $t('installOption.parallelInstall') }}
            </label>
          </div>

          <div class="flex gap-x-2">
            <label class="flex items-center gap-x-1.5 select-none cursor-pointer">
              <input
                type="checkbox"
                name="set_password"
                v-model="settingStore.settings.set_password"
                class="checkbox checkbox-sm checkbox-primary"
              />
              {{ $t('installOption.setPassword') }}
            </label>

            <input
              type="text"
              name="password"
              v-model="settingStore.settings.password"
              class="max-w-28 input input-sm input-accent"
              :disabled="!settingStore.settings.set_password"
            />
          </div>
        </div>
      </div>

      <div class="flex flex-col grow justify-between">
        <div>
          <label class="block mb-1 text-sm text-gray-900">
            {{ $t('installOption.successAction') }}
          </label>

          <select
            name="success_action"
            v-model="settingStore.settings.success_action"
            class="select select-accent w-full"
          >
            <option v-for="action in storage.SuccessAction" :key="action" :value="action">
              {{ $t(`successAction.${action}`) }}
            </option>
          </select>
        </div>

        <div class="flex flex-row gap-x-3 justify-end items-center mt-2 h-8">
          <button
            type="button"
            class="btn btn-outline btn-secondary border-2"
            @click="
              () => {
                form?.reset()
                // settingStore.restore()
              }
            "
          >
            {{ $t('installOption.reset') }}
          </button>
          <button class="btn btn-secondary" @click="handleSubmit">
            {{ $t('installOption.execute') }}
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
