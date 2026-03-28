import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { usePlayerStore } from '@/stores/player'
import { STAT_RANGES } from '@/types/game'

describe('Player Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initializes with default stats', () => {
    const store = usePlayerStore()
    expect(store.stats.str).toBe(50)
    expect(store.stats.end).toBe(50)
    expect(store.stats.agi).toBe(50)
    expect(store.stats.mor).toBe(50)
    expect(store.stats.disc).toBe(0)
  })

  it('initializes with default ranks', () => {
    const store = usePlayerStore()
    expect(store.formalRank).toBe('РЯДОВОЙ')
    expect(store.informalStatus).toBe('ЗАПАХ')
  })

  it('setStat clamps values to valid range', () => {
    const store = usePlayerStore()
    
    store.setStat('str', 150)
    expect(store.stats.str).toBe(STAT_RANGES.str.max)
    
    store.setStat('str', -50)
    expect(store.stats.str).toBe(STAT_RANGES.str.min)
    
    store.setStat('str', 75)
    expect(store.stats.str).toBe(75)
  })

  it('setStat updates rank when disc changes', () => {
    const store = usePlayerStore()
    
    store.setStat('disc', 80)
    expect(store.formalRank).toBe('СЕРЖАНТ')
    
    store.setStat('disc', 30)
    expect(store.formalRank).toBe('ЕФРЕЙТОР')
  })

  it('computedFormalRank returns correct rank based on disc', () => {
    const store = usePlayerStore()
    
    store.setStat('disc', 0)
    expect(store.formalRank).toBe('РЯДОВОЙ')
    
    store.setStat('disc', 25)
    expect(store.formalRank).toBe('ЕФРЕЙТОР')
    
    store.setStat('disc', 50)
    expect(store.formalRank).toBe('МЛ_СЕРЖАНТ')
    
    store.setStat('disc', 75)
    expect(store.formalRank).toBe('СЕРЖАНТ')
  })

  it('computedInformalStatus returns correct status based on disc', () => {
    const store = usePlayerStore()
    
    store.setStat('disc', 0)
    expect(store.informalStatus).toBe('ЗАПАХ')
    
    store.setStat('disc', -30)
    expect(store.informalStatus).toBe('СЛОН')
    
    store.setStat('disc', -55)
    expect(store.informalStatus).toBe('ЧЕРПАК')
    
    store.setStat('disc', -80)
    expect(store.informalStatus).toBe('ДЕД')
    
    store.setStat('disc', -85)
    expect(store.informalStatus).toBe('ДЕД')
    
    store.setStat('disc', -95)
    expect(store.informalStatus).toBe('ДЕМБЕЛЬ')
  })

  it('updateStats updates multiple stats at once', () => {
    const store = usePlayerStore()
    
    store.updateStats({
      str: 80,
      end: 70,
      agi: 60,
      mor: 40,
      disc: 20,
    })
    
    expect(store.stats.str).toBe(80)
    expect(store.stats.end).toBe(70)
    expect(store.stats.agi).toBe(60)
    expect(store.stats.mor).toBe(40)
    expect(store.stats.disc).toBe(20)
  })

  it('setVersion updates version', () => {
    const store = usePlayerStore()
    store.setVersion(5)
    expect(store.version).toBe(5)
  })

  it('loadFromPlayer loads player data correctly', () => {
    const store = usePlayerStore()
    
    const player = {
      id: 'test-id',
      stats: { str: 70, end: 60, agi: 80, mor: 50, disc: 30 },
      formalRank: 'ЕФРЕЙТОР' as const,
      informalStatus: 'ЗАПАХ' as const,
      turn: 10,
      flags: [],
      isFinished: false,
      version: 2,
      createdAt: '',
      updatedAt: '',
    }
    
    store.loadFromPlayer(player)
    
    expect(store.stats.str).toBe(70)
    expect(store.stats.end).toBe(60)
    expect(store.stats.agi).toBe(80)
    expect(store.stats.mor).toBe(50)
    expect(store.stats.disc).toBe(30)
    expect(store.formalRank).toBe('ЕФРЕЙТОР')
    expect(store.informalStatus).toBe('ЗАПАХ')
    expect(store.version).toBe(2)
    expect(store.isLoaded).toBe(true)
  })

  it('reset restores default values', () => {
    const store = usePlayerStore()
    
    store.setStat('str', 90)
    store.setStat('disc', 50)
    store.setVersion(10)
    
    store.reset()
    
    expect(store.stats.str).toBe(50)
    expect(store.stats.end).toBe(50)
    expect(store.stats.agi).toBe(50)
    expect(store.stats.mor).toBe(50)
    expect(store.stats.disc).toBe(0)
    expect(store.formalRank).toBe('РЯДОВОЙ')
    expect(store.informalStatus).toBe('ЗАПАХ')
    expect(store.version).toBe(1)
    expect(store.isLoaded).toBe(false)
  })
})
