import { readdirSync, readFileSync, writeFileSync } from 'fs'
import { join } from 'path'

const pathEnTranslation = join(import.meta.dirname, 'en.json')
const template = sortObject(JSON.parse(readFileSync(pathEnTranslation, 'utf8')))

/**
 * Sort the given `obj` by key
 * @param {object} obj
 * @returns {object}
 */
function sortObject(obj) {
  return Object.keys(obj)
    .sort()
    .reduce((r, k) => ((r[k] = obj[k]), r), {})
}

// Save the sorted JSON back to template.json
writeFileSync(pathEnTranslation, JSON.stringify(template, null, 2))

// Add missing keys from template to fileJson
readdirSync(import.meta.dirname).forEach(file => {
  const pathCurrent = join(import.meta.dirname, file)

  if (pathCurrent === pathEnTranslation || !pathCurrent.includes('.json')) {
    return
  }

  const translation = JSON.parse(readFileSync(pathCurrent, 'utf8'))

  Object.keys(template).forEach(key => {
    if (!(key in translation)) {
      console.log(`[${file}]\t + "${key}"`)
      translation[key] = template[key]
    }
  })

  Object.keys(translation).forEach(key => {
    if (!(key in template)) {
      console.log(`[${file}]\t - "${key}"`)
      delete translation[key]
    }
  })

  writeFileSync(pathCurrent, JSON.stringify(sortObject(translation), null, 2))
})
