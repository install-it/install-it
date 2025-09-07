import { marked } from 'marked'
import * as semver from 'semver'

export async function latestRelease(currentVersion: string) {
  return fetch('https://api.github.com/repos/driverbox/driver-box/releases/latest')
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
