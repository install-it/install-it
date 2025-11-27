import type { storage } from '@/wailsjs/go/models'
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

/**
 * Tests whether the given input string satisfies the specified match rule.
 */
export function testMatchRule(rule: storage.Rule, input: string) {
  input = rule.is_case_sensitive ? input : input.toLowerCase()
  const values = rule.is_case_sensitive ? rule.values : rule.values.map(v => v.toLowerCase())
  const hits = values.map((v: string): boolean => {
    switch (rule.operator) {
      case 'contain':
        return input.includes(v)
      case 'not_contain':
        return !input.includes(v)
      case 'equal':
        return input === v
      case 'not_equal':
        return input !== v
      case 'regex': {
        try {
          return new RegExp(v, rule.is_case_sensitive ? '' : 'i').test(input)
        } catch {
          return false
        }
      }
      default:
        return false
    }
  })

  return rule.should_hit_all ? hits.every(Boolean) : hits.some(Boolean)
}
