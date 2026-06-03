import { describe, it, expect } from 'vitest'
import { testMatchRule } from '@/utils/index'
import { storage } from '@/wailsjs/go/models'

/**
 * Helper: build a storage.Rule with safe defaults, then apply overrides.
 */
function makeRule(overrides: {
  operator?: storage.RuleOperator
  is_case_sensitive?: boolean
  should_hit_all?: boolean
  values?: string[]
}): storage.Rule {
  return new storage.Rule({
    source: storage.RuleSource.CPU,
    operator: storage.RuleOperator.CONTAIN,
    is_case_sensitive: false,
    should_hit_all: false,
    values: [],
    ...overrides
  })
}

describe('testMatchRule', () => {
  // ─── CONTAIN ──────────────────────────────────────────────────────────────

  describe('CONTAIN operator', () => {
    it('returns true when value is found in input', () => {
      const rule = makeRule({ operator: storage.RuleOperator.CONTAIN, values: ['Intel'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(true)
    })

    it('returns false when value is not found in input', () => {
      const rule = makeRule({ operator: storage.RuleOperator.CONTAIN, values: ['AMD'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })

    it('is case-insensitive when is_case_sensitive is false', () => {
      const rule = makeRule({ operator: storage.RuleOperator.CONTAIN, values: ['INTEL'], is_case_sensitive: false })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(true)
    })

    it('is case-sensitive when is_case_sensitive is true', () => {
      const rule = makeRule({ operator: storage.RuleOperator.CONTAIN, values: ['INTEL'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })
  })

  // ─── NOT_CONTAIN ──────────────────────────────────────────────────────────

  describe('NOT_CONTAIN operator', () => {
    it('returns true when value is absent from input', () => {
      const rule = makeRule({ operator: storage.RuleOperator.NOT_CONTAIN, values: ['AMD'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(true)
    })

    it('returns false when value is present in input', () => {
      const rule = makeRule({ operator: storage.RuleOperator.NOT_CONTAIN, values: ['Intel'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })

    it('is case-insensitive when is_case_sensitive is false', () => {
      const rule = makeRule({ operator: storage.RuleOperator.NOT_CONTAIN, values: ['intel'], is_case_sensitive: false })
      // lowercased input 'intel core i9' contains 'intel', so NOT_CONTAIN → false
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })
  })

  // ─── EQUAL ────────────────────────────────────────────────────────────────

  describe('EQUAL operator', () => {
    it('returns true for exact match', () => {
      const rule = makeRule({ operator: storage.RuleOperator.EQUAL, values: ['Intel Core i9'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(true)
    })

    it('returns false for partial match', () => {
      const rule = makeRule({ operator: storage.RuleOperator.EQUAL, values: ['Intel'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })

    it('is case-insensitive when is_case_sensitive is false', () => {
      const rule = makeRule({ operator: storage.RuleOperator.EQUAL, values: ['intel core i9'], is_case_sensitive: false })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(true)
    })

    it('is case-sensitive when is_case_sensitive is true', () => {
      const rule = makeRule({ operator: storage.RuleOperator.EQUAL, values: ['intel core i9'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })
  })

  // ─── NOT_EQUAL ────────────────────────────────────────────────────────────

  describe('NOT_EQUAL operator', () => {
    it('returns true when input does not equal value', () => {
      const rule = makeRule({ operator: storage.RuleOperator.NOT_EQUAL, values: ['AMD Ryzen 9'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(true)
    })

    it('returns false when input equals value', () => {
      const rule = makeRule({ operator: storage.RuleOperator.NOT_EQUAL, values: ['Intel Core i9'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })

    it('is case-insensitive when is_case_sensitive is false', () => {
      const rule = makeRule({ operator: storage.RuleOperator.NOT_EQUAL, values: ['intel core i9'], is_case_sensitive: false })
      // lowercased input matches, so NOT_EQUAL → false
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })
  })

  // ─── REGEX ────────────────────────────────────────────────────────────────

  describe('REGEX operator', () => {
    it('returns true when pattern matches', () => {
      const rule = makeRule({ operator: storage.RuleOperator.REGEX, values: ['Intel.+i\\d'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(true)
    })

    it('returns false when pattern does not match', () => {
      const rule = makeRule({ operator: storage.RuleOperator.REGEX, values: ['^AMD'], is_case_sensitive: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })

    it('does not throw and returns false for an invalid regex pattern', () => {
      const rule = makeRule({ operator: storage.RuleOperator.REGEX, values: ['[invalid'], is_case_sensitive: false })
      expect(() => testMatchRule(rule, 'Intel Core i9')).not.toThrow()
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })

    it('applies case-insensitive flag when is_case_sensitive is false', () => {
      // Pattern is lowercase; input starts with uppercase 'I'. With 'i' flag it still matches.
      const rule = makeRule({ operator: storage.RuleOperator.REGEX, values: ['^intel'], is_case_sensitive: false })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(true)
    })

    it('handles complex patterns (anchors, quantifiers, character classes)', () => {
      const rule = makeRule({
        operator: storage.RuleOperator.REGEX,
        values: ['\\bCore\\s+i[0-9]+\\b'],
        is_case_sensitive: true
      })
      expect(testMatchRule(rule, 'Intel Core i9 Gen 13')).toBe(true)
      expect(testMatchRule(rule, 'Intel Xeon E5-2690')).toBe(false)
    })
  })

  // ─── should_hit_all aggregation ───────────────────────────────────────────

  describe('should_hit_all aggregation', () => {
    it('should_hit_all=true requires ALL values to match', () => {
      const rule = makeRule({
        operator: storage.RuleOperator.CONTAIN,
        values: ['Intel', 'i9'],
        is_case_sensitive: true,
        should_hit_all: true
      })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(true)
      expect(testMatchRule(rule, 'Intel Core i7')).toBe(false)
    })

    it('should_hit_all=false returns true when ANY value matches', () => {
      const rule = makeRule({
        operator: storage.RuleOperator.CONTAIN,
        values: ['AMD', 'Intel'],
        is_case_sensitive: true,
        should_hit_all: false
      })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(true)
      expect(testMatchRule(rule, 'Nvidia RTX 4090')).toBe(false)
    })

    it('should_hit_all=false returns false when no values match', () => {
      const rule = makeRule({
        operator: storage.RuleOperator.CONTAIN,
        values: ['AMD', 'Ryzen'],
        is_case_sensitive: true,
        should_hit_all: false
      })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })
  })

  // ─── Edge cases ───────────────────────────────────────────────────────────

  describe('edge cases', () => {
    it('empty values with should_hit_all=true returns true (vacuous truth)', () => {
      const rule = makeRule({ operator: storage.RuleOperator.CONTAIN, values: [], should_hit_all: true })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(true)
    })

    it('empty values with should_hit_all=false returns false', () => {
      const rule = makeRule({ operator: storage.RuleOperator.CONTAIN, values: [], should_hit_all: false })
      expect(testMatchRule(rule, 'Intel Core i9')).toBe(false)
    })

    it('handles unicode hardware names', () => {
      const rule = makeRule({
        operator: storage.RuleOperator.CONTAIN,
        values: ['英特尔'],
        is_case_sensitive: false,
        should_hit_all: false
      })
      expect(testMatchRule(rule, '英特尔 Core i9')).toBe(true)
      expect(testMatchRule(rule, 'AMD Ryzen')).toBe(false)
    })

    it('empty input string with CONTAIN returns false', () => {
      const rule = makeRule({ operator: storage.RuleOperator.CONTAIN, values: ['Intel'], is_case_sensitive: false })
      expect(testMatchRule(rule, '')).toBe(false)
    })

    it('empty input string with NOT_CONTAIN returns true', () => {
      const rule = makeRule({ operator: storage.RuleOperator.NOT_CONTAIN, values: ['Intel'], is_case_sensitive: false })
      expect(testMatchRule(rule, '')).toBe(true)
    })
  })
})
