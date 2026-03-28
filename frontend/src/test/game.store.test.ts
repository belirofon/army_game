import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useGameStore } from '@/stores/game'
import { usePlayerStore } from '@/stores/player'
import { mockPlayer, mockEvent } from './setup'

const localStorageMock = (() => {
  let store: Record<string, string> = {}
  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => { store[key] = value }),
    removeItem: vi.fn((key: string) => { delete store[key] }),
    clear: vi.fn(() => { store = {} }),
  }
})()

Object.defineProperty(global, 'localStorage', { value: localStorageMock })

describe('Game Store', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.clear()
    setActivePinia(createPinia())
  })

  it('initializes with default values', () => {
    const store = useGameStore()
    expect(store.gameId).toBeNull()
    expect(store.selectedCharacter).toBeNull()
    expect(store.currentEvent).toBeNull()
    expect(store.eventHistory).toEqual([])
    expect(store.turn).toBe(1)
    expect(store.isGameOver).toBe(false)
    expect(store.final).toBeNull()
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
  })

  it('setGameId updates gameId and localStorage', () => {
    const store = useGameStore()
    
    store.setGameId('game-123')
    
    expect(store.gameId).toBe('game-123')
    expect(localStorageMock.setItem).toHaveBeenCalledWith('army_game_gameId', 'game-123')
  })

  it('setGameId with null removes from localStorage', () => {
    const store = useGameStore()
    store.setGameId('game-123')
    
    store.setGameId(null)
    
    expect(store.gameId).toBeNull()
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('army_game_gameId')
  })

  it('loads gameId from localStorage on initialization', () => {
    localStorageMock.getItem.mockReturnValue('saved-game-id')
    
    setActivePinia(createPinia())
    const store = useGameStore()
    
    expect(store.gameId).toBe('saved-game-id')
  })

  it('setSelectedCharacter updates character and player stats', () => {
    const store = useGameStore()
    const playerStore = usePlayerStore()
    
    const character = {
      id: 'char-1',
      name: 'Иван',
      stats: { str: 70, end: 60, agi: 80, mor: 50, disc: 20 },
    }
    
    store.setSelectedCharacter(character)
    
    expect(store.selectedCharacter).toEqual(character)
    expect(playerStore.stats.str).toBe(70)
    expect(playerStore.stats.end).toBe(60)
    expect(playerStore.stats.agi).toBe(80)
    expect(playerStore.stats.mor).toBe(50)
    expect(playerStore.stats.disc).toBe(20)
  })

  it('setCurrentEvent updates current event', () => {
    const store = useGameStore()
    
    store.setCurrentEvent(mockEvent)
    
    expect(store.currentEvent).toEqual(mockEvent)
  })

  it('setTurn updates turn number', () => {
    const store = useGameStore()
    
    store.setTurn(10)
    
    expect(store.turn).toBe(10)
  })

  it('addToHistory adds entry to history', () => {
    const store = useGameStore()
    
    const entry = {
      id: 'entry-1',
      playerId: 'player-1',
      turn: 1,
      eventDescription: 'Test event',
      choiceText: 'Choice 1',
      checkResult: {
        success: true,
        outcome: 'SUCCESS' as const,
        description: 'Success',
      },
      effects: [],
      createdAt: '2024-01-01',
    }
    
    store.addToHistory(entry)
    
    expect(store.eventHistory).toHaveLength(1)
    expect(store.eventHistory[0]).toEqual(entry)
  })

  it('addToHistory keeps only last 10 entries', () => {
    const store = useGameStore()
    
    for (let i = 0; i < 15; i++) {
      store.addToHistory({
        id: `entry-${i}`,
        playerId: 'player-1',
        turn: i,
        eventDescription: `Event ${i}`,
        choiceText: `Choice ${i}`,
        checkResult: { success: true, outcome: 'SUCCESS' as const, description: '' },
        effects: [],
        createdAt: '',
      })
    }
    
    expect(store.eventHistory).toHaveLength(10)
    expect(store.eventHistory[0].id).toBe('entry-14')
  })

  it('setGameOver sets game over state', () => {
    const store = useGameStore()
    const final = {
      type: 'ТИХИЙ_ДЕМБЕЛЬ' as const,
      title: 'Тихо дембель',
      description: 'Description',
      finalStats: { str: 50, end: 50, agi: 50, mor: 50, disc: -50 },
      achievedOnTurn: 30,
    }
    
    store.setGameOver(final)
    
    expect(store.isGameOver).toBe(true)
    expect(store.final).toEqual(final)
  })

  it('setLoading updates loading state', () => {
    const store = useGameStore()
    
    store.setLoading(true)
    expect(store.loading).toBe(true)
    
    store.setLoading(false)
    expect(store.loading).toBe(false)
  })

  it('setError updates error state', () => {
    const store = useGameStore()
    
    store.setError('Something went wrong')
    expect(store.error).toBe('Something went wrong')
    
    store.setError(null)
    expect(store.error).toBeNull()
  })

  it('loadGameState loads full game state', () => {
    const store = useGameStore()
    
    const history = [
      {
        id: 'entry-1',
        playerId: 'player-1',
        turn: 1,
        eventDescription: 'Event 1',
        choiceText: 'Choice 1',
        checkResult: { success: true, outcome: 'SUCCESS' as const, description: '' },
        effects: [],
        createdAt: '',
      },
    ]
    
    store.loadGameState(mockPlayer, mockEvent, history, true, null)
    
    expect(usePlayerStore().stats.str).toBe(60)
    expect(store.turn).toBe(5)
    expect(store.currentEvent).toEqual(mockEvent)
    expect(store.eventHistory).toEqual(history)
    expect(store.isGameOver).toBe(true)
  })

  it('reset clears all state', () => {
    const store = useGameStore()
    const playerStore = usePlayerStore()
    
    store.setGameId('game-123')
    store.setSelectedCharacter({
      id: 'char-1',
      name: 'Test',
      stats: { str: 80, end: 80, agi: 80, mor: 80, disc: 50 },
    })
    store.setTurn(15)
    playerStore.setStat('str', 90)
    
    store.reset()
    
    expect(store.gameId).toBeNull()
    expect(store.selectedCharacter).toBeNull()
    expect(store.currentEvent).toBeNull()
    expect(store.eventHistory).toEqual([])
    expect(store.turn).toBe(1)
    expect(store.isGameOver).toBe(false)
    expect(store.final).toBeNull()
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
    expect(playerStore.stats.str).toBe(50)
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('army_game_gameId')
  })

  it('computedIsGameOver is true when mor <= 0', () => {
    const store = useGameStore()
    const playerStore = usePlayerStore()
    
    playerStore.setStat('mor', 0)
    
    expect(store.isGameOver).toBe(true)
  })

  it('computedIsGameOver is true when turn >= 30', () => {
    const store = useGameStore()
    
    store.setTurn(30)
    
    expect(store.isGameOver).toBe(true)
  })

  it('computedIsGameOver is true when isGameOver flag is set', () => {
    const store = useGameStore()
    
    store.setGameOver(null)
    
    expect(store.isGameOver).toBe(true)
  })
})
