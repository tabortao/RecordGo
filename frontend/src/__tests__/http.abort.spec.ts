import { describe, it, expect } from 'vitest'
import { isAbortError } from '@/services/http'

describe('isAbortError', () => {
  it('returns true for ERR_CANCELED code', () => {
    const e = { code: 'ERR_CANCELED', message: 'canceled' }
    expect(isAbortError(e)).toBe(true)
  })

  it('returns true for aborted message', () => {
    const e = { code: '', message: 'Request aborted by the browser' }
    expect(isAbortError(e)).toBe(true)
  })

  it('returns false for generic network error', () => {
    const e = { code: 'ECONNRESET', message: 'Network error' }
    expect(isAbortError(e)).toBe(false)
  })
})

