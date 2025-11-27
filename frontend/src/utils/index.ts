import * as libsysi from '@/wailsjs/go/sysinfo/SysInfo'
import { marked } from 'marked'
import * as semver from 'semver'

export async function latestRelease(currentVersion: string) {
  return fetch('https://api.github.com/repos/install-it/install-it/releases/latest')
    .then(response => response.json())
    .then(async body => {
      const version = semver.clean(body.tag_name) || '0.0.0'

      return {
        hasUpdate: semver.gt(version, currentVersion),
        name: body.name as string,
        releaseAt: new Date(Date.parse(body.published_at)),
        releaseNotes: await marked.parse(body.body),
        tag: body.tag_name as string,
        url: body.html_url as string,
        version: version
      }
    })
}

/**
 * Retrieves detailed hardware information from the system.
 */
export async function getHardware() {
  return Promise.all([
    libsysi.CpuInfo().then(vs => vs.map(v => v.Name)),
    libsysi
      .GpuInfo()
      .then(vs => vs.map(v => `${v.Name} (${Math.round(v.AdapterRAM / Math.pow(1024, 3))}GB)`)),
    libsysi
      .MemoryInfo()
      .then(vs =>
        vs.map(
          v =>
            `${v.Manufacturer} ${v.PartNumber.trim()} ${v.Capacity / Math.pow(1024, 3)}GB ${v.Speed}MHz`
        )
      ),
    libsysi.MotherboardInfo().then(vs => vs.map(v => `${v.Manufacturer} ${v.Product}`)),
    libsysi.NicInfo().then(vs => vs.map(v => v.Name)),
    libsysi
      .DiskInfo()
      .then(vs => vs.map(v => `${v.Model} (${Math.round(v.Size / Math.pow(1024, 3))}GB)`))
  ]).then(parts => ({
    cpu: parts[0],
    gpu: parts[1],
    memory: parts[2],
    motherboard: parts[3],
    nic: parts[4],
    storage: parts[5]
  }))
}
