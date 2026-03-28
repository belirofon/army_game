import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUiStore } from '@/stores/ui'

describe('UI Store', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    setActivePinia(createPinia())
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('initializes with default values', () => {
    const store = useUiStore()
    expect(store.statsPanelOpen).toBe(false)
    expect(store.historyPanelOpen).toBe(false)
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
    expect(store.toast).toBeNull()
  })

  it('openStatsPanel sets statsPanelOpen to true', () => {
    const store = useUiStore()
    store.openStatsPanel()
    expect(store.statsPanelOpen).toBe(true)
  })

  it('closeStatsPanel sets statsPanelOpen to false', () => {
    const store = useUiStore()
    store.openStatsPanel()
    store.closeStatsPanel()
    expect(store.statsPanelOpen).toBe(false)
  })

  it('openHistoryPanel sets historyPanelOpen to true', () => {
    const store = useUiStore()
    store.openHistoryPanel()
    expect(store.historyPanelOpen).toBe(true)
  })

  it('closeHistoryPanel sets historyPanelOpen to false', () => {
    const store = useUiStore()
    store.openHistoryPanel()
    store.closeHistoryPanel()
    expect(store.historyPanelOpen).toBe(false)
  })

  it('setLoading updates loading state', () => {
    const store = useUiStore()
    store.setLoading(true)
    expect(store.loading).toBe(true)
    store.setLoading(false)
    expect(store.loading).toBe(false)
  })

  it('setError updates error state', () => {
    const store = useUiStore()
    store.setError('Test error')
    expect(store.error).toBe('Test error')
    store.setError(null)
    expect(store.error).toBeNull()
  })

  it('showToast displays toast and clears after timeout', () => {
    const store = useUiStore()
    
    store.showToast('Success message', 'success')
    
    expect(store.toast).toEqual({
      message: 'Success message',
      type: 'success',
    })
    
    vi.advanceTimersByTime(3000)
    
    expect(store.toast).toBeNull()
  })

  it('showToast accepts different types', () => {
    const store = useUiStore()
    
    store.showToast('Error message', 'error')
    expect(store.toast?.type).toBe('error')
    
    store.showToast('Warning message', 'warning')
    expect(store.toast?.type).toBe('warning')
  })

  it('showToast defaults to success type', () => {
    const store = useUiStore()
    
    store.showToast('Default message')
    
    expect(store.toast?.type).toBe('success')
  })

  it('clearToast clears toast immediately', () => {
    const store = useUiStore()
    
    store.showToast('Message', 'error')
    expect(store.toast).not.toBeNull()
    
    store.clearToast()
    expect(store.toast).toBeNull()
    
    vi.advanceTimersByTime(3000)
    expect(store.toast).toBeNull()
  })

  it('multiple showToast calls override previous toast', () => {
    const store = useUiStore()
    
    store.showToast('First message', 'success')
    store.showToast('Second message', 'error')
    
    expect(store.toast?.message).toBe('Second message')
  })

  it('toast auto-clears even if already cleared manually', () => {
    const store = useUiStore()
    
    store.showToast('Message', 'success')
    store.clearToast()
    
    vi.advanceTimersByTime(3000)
    
    expect(store.toast).toBeNull()
  })
})
