import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { usePlayer } from '@/composables/usePlayer'
import { usePlayerStore } from '@/stores/player'

describe('usePlayer Composable', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('returns stats from player store', () => {
    const { stats } = usePlayer()
    const store = usePlayerStore()
    
    store.setStat('str', 75)
    
    expect(stats.value.str).toBe(75)
  })

  it('returns formalRank from player store', () => {
    const { formalRank } = usePlayer()
    const store = usePlayerStore()
    
    store.setStat('disc', 60)
    
    expect(formalRank.value).toBe('МЛ_СЕРЖАНТ')
  })

  it('returns informalStatus from player store', () => {
    const { informalStatus } = usePlayer()
    const store = usePlayerStore()
    
    store.setStat('disc', -60)
    
    expect(informalStatus.value).toBe('ЧЕРПАК')
  })

  it('returns version from player store', () => {
    const { version } = usePlayer()
    const store = usePlayerStore()
    
    store.setVersion(5)
    
    expect(version.value).toBe(5)
  })

  it('applyEffect updates player stat', () => {
    const { stats, applyEffect } = usePlayer()
    const store = usePlayerStore()
    
    store.setStat('str', 50)
    
    applyEffect({
      stat: 'str',
      delta: 10,
      previousValue: 50,
      newValue: 60,
    })
    
    expect(stats.value.str).toBe(60)
  })

  it('applyEffect clamps value to valid range', () => {
    const { stats, applyEffect } = usePlayer()
    const store = usePlayerStore()
    
    store.setStat('str', 50)
    
    applyEffect({
      stat: 'str',
      delta: 100,
      previousValue: 50,
      newValue: 150,
    })
    
    expect(stats.value.str).toBe(100)
  })

  it('applyEffect clamps negative values to minimum', () => {
    const { stats, applyEffect } = usePlayer()
    const store = usePlayerStore()
    
    store.setStat('str', 50)
    
    applyEffect({
      stat: 'str',
      delta: -100,
      previousValue: 50,
      newValue: -50,
    })
    
    expect(stats.value.str).toBe(1)
  })

  it('applyEffects applies multiple effects', () => {
    const { stats, applyEffects } = usePlayer()
    const store = usePlayerStore()
    
    store.setStat('str', 50)
    store.setStat('end', 50)
    store.setStat('mor', 50)
    
    applyEffects([
      { stat: 'str', delta: 10, previousValue: 50, newValue: 60 },
      { stat: 'end', delta: 20, previousValue: 50, newValue: 70 },
      { stat: 'mor', delta: -30, previousValue: 50, newValue: 20 },
    ])
    
    expect(stats.value.str).toBe(60)
    expect(stats.value.end).toBe(70)
    expect(stats.value.mor).toBe(20)
  })

  it('reset calls store reset', () => {
    const { stats, reset } = usePlayer()
    const store = usePlayerStore()
    
    store.setStat('str', 90)
    store.setStat('disc', 50)
    store.setVersion(10)
    
    reset()
    
    expect(stats.value.str).toBe(50)
    expect(stats.value.disc).toBe(0)
    expect(usePlayerStore().version).toBe(1)
  })

  it('works with mor stat correctly (min is 0)', () => {
    const { stats, applyEffect } = usePlayer()
    const store = usePlayerStore()
    
    store.setStat('mor', 10)
    
    applyEffect({
      stat: 'mor',
      delta: -20,
      previousValue: 10,
      newValue: -10,
    })
    
    expect(stats.value.mor).toBe(0)
  })

  it('works with disc stat with negative range', () => {
    const { stats, applyEffect } = usePlayer()
    const store = usePlayerStore()
    
    store.setStat('disc', 0)
    
    applyEffect({
      stat: 'disc',
      delta: -50,
      previousValue: 0,
      newValue: -50,
    })
    
    expect(stats.value.disc).toBe(-50)
  })

  it('handles disc going below -100', () => {
    const { stats, applyEffect } = usePlayer()
    const store = usePlayerStore()
    
    store.setStat('disc', -80)
    
    applyEffect({
      stat: 'disc',
      delta: -50,
      previousValue: -80,
      newValue: -130,
    })
    
    expect(stats.value.disc).toBe(-100)
  })
})
