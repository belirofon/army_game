import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { EventInstance, GameLogEntry, Final, Player, PlayerStats } from '@/types/game'
import { usePlayerStore } from './player'

const STORAGE_KEY = 'army_game_gameId'

interface SelectedCharacter {
  id: string
  name: string
  stats: PlayerStats
}

export const useGameStore = defineStore('game', () => {
  const playerStore = usePlayerStore()
  
  const gameId = ref<string | null>(null)
  const selectedCharacter = ref<SelectedCharacter | null>(null)
  const currentEvent = ref<EventInstance | null>(null)
  const eventHistory = ref<GameLogEntry[]>([])
  const turn = ref(1)
  const isGameOver = ref(false)
  const final = ref<Final | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const computedIsGameOver = computed(() => {
    return isGameOver.value || playerStore.stats.mor <= 0 || turn.value >= 24
  })

  function setGameId(id: string | null) {
    gameId.value = id
    if (id) {
      localStorage.setItem(STORAGE_KEY, id)
    } else {
      localStorage.removeItem(STORAGE_KEY)
    }
  }

  function setSelectedCharacter(char: SelectedCharacter | null) {
    selectedCharacter.value = char
    if (char) {
      playerStore.updateStats(char.stats)
    }
  }

  function loadFromStorage() {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      gameId.value = saved
    }
  }

  function setCurrentEvent(event: EventInstance | null) {
    currentEvent.value = event
  }

  function setTurn(newTurn: number) {
    turn.value = newTurn
  }

  function addToHistory(entry: GameLogEntry) {
    eventHistory.value.unshift(entry)
    if (eventHistory.value.length > 10) {
      eventHistory.value.pop()
    }
  }

  function setGameOver(f: Final | null) {
    isGameOver.value = true
    final.value = f
  }

  function setLoading(l: boolean) {
    loading.value = l
  }

  function setError(e: string | null) {
    error.value = e
  }

  function loadGameState(player: Player, event: EventInstance | null, history: GameLogEntry[], gameFinished: boolean = false, f: Final | null = null) {
    playerStore.loadFromPlayer(player)
    turn.value = player.turn
    currentEvent.value = event
    eventHistory.value = history
    isGameOver.value = gameFinished
    final.value = f
  }

  function reset() {
    gameId.value = null
    selectedCharacter.value = null
    currentEvent.value = null
    eventHistory.value = []
    turn.value = 1
    isGameOver.value = false
    final.value = null
    loading.value = false
    error.value = null
    playerStore.reset()
    localStorage.removeItem(STORAGE_KEY)
  }

  loadFromStorage()

  return {
    gameId,
    selectedCharacter,
    currentEvent,
    eventHistory,
    turn,
    isGameOver: computedIsGameOver,
    final,
    loading,
    error,
    setGameId,
    setSelectedCharacter,
    setCurrentEvent,
    setTurn,
    addToHistory,
    setGameOver,
    setLoading,
    setError,
    loadGameState,
    reset,
  }
})