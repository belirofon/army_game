import { useApolloClient } from '@vue/apollo-composable'
import { ref } from 'vue'
import { useGameStore } from '@/stores/game'
import { usePlayerStore } from '@/stores/player'
import { START_GAME, CHOOSE, SELECT_CHARACTER, RESTART_GAME } from '@/api/mutations'
import { LOAD_GAME } from '@/api/queries'
import type { ChooseResult, EventInstance, GameLogEntry, Final } from '@/types/game'

export function useApi() {
  const client = useApolloClient()
  const gameStore = useGameStore()
  const playerStore = usePlayerStore()

  const loading = ref(false)
  const error = ref<string | null>(null)

  async function startGame(): Promise<boolean> {
    loading.value = true
    error.value = null
    
    try {
      const result = await client.client.mutate({
        mutation: START_GAME,
      })
      
      const data = result.data?.startGame
      if (data) {
        gameStore.setGameId(data.gameId)
        playerStore.loadFromPlayer(data.player)
        gameStore.setTurn(data.player.turn)
        gameStore.setCurrentEvent(data.currentEvent)
        gameStore.setGameOver(data.isGameOver ? data.final : null)
        return true
      }
      return false
    } catch (e: any) {
      error.value = e.message || 'Failed to start game'
      return false
    } finally {
      loading.value = false
    }
  }

  async function makeChoice(choiceId: string): Promise<boolean> {
    const gameId = gameStore.gameId
    const expectedVersion = playerStore.version
    
    if (!gameId) {
      error.value = 'No active game'
      return false
    }

    loading.value = true
    error.value = null

    try {
      const result = await client.client.mutate({
        mutation: CHOOSE,
        variables: {
          playerId: gameId,
          choiceId,
          expectedVersion,
        },
      })
      
      const data = result.data?.choose
      if (data) {
        // Reload game state to ensure consistency
        return await loadGame()
      }
      return false
    } catch (e: any) {
      if (e.graphQLErrors?.[0]?.extensions?.code === 'CONCURRENT_MODIFICATION') {
        error.value = 'Game state changed, please retry'
        await loadGame()
      } else {
        error.value = e.message || 'Failed to make choice'
      }
      return false
    } finally {
      loading.value = false
    }
  }

  async function loadGame(): Promise<boolean> {
    const gameId = gameStore.gameId
    if (!gameId) return false

    loading.value = true
    error.value = null

    try {
      const result = await client.client.query({
        query: LOAD_GAME,
        variables: { gameId },
      })
      
      const data = result.data?.loadGame
      if (data) {
        gameStore.loadGameState(
          data.player,
          data.currentEvent,
          data.eventHistory,
          data.isGameOver,
          data.final
        )
        return true
      }
      return false
    } catch (e: any) {
      error.value = e.message || 'Failed to load game'
      return false
    } finally {
      loading.value = false
    }
  }

  async function selectCharacter(characterType: string, stats: {
    str: number
    end: number
    agi: number
    mor: number
    disc: number
  }): Promise<boolean> {
    const gameId = gameStore.gameId
    if (!gameId) {
      error.value = 'No active game'
      return false
    }

    loading.value = true
    error.value = null

    try {
      const result = await client.client.mutate({
        mutation: SELECT_CHARACTER,
        variables: {
          playerId: gameId,
          characterType,
          stats,
        },
      })

      const data = result.data?.selectCharacter
      if (data) {
        playerStore.loadFromPlayer(data)
        return true
      }
      return false
    } catch (e: any) {
      error.value = e.message || 'Failed to select character'
      return false
    } finally {
      loading.value = false
    }
  }

  async function restartGame(): Promise<boolean> {
    const gameId = gameStore.gameId
    if (!gameId) return startGame()

    loading.value = true
    error.value = null

    try {
      const result = await client.client.mutate({
        mutation: RESTART_GAME,
        variables: { playerId: gameId },
      })
      
      const data = result.data?.restartGame
      if (data) {
        gameStore.setGameId(data.gameId)
        playerStore.loadFromPlayer(data.player)
        gameStore.setTurn(data.player.turn)
        gameStore.setCurrentEvent(data.currentEvent)
        gameStore.setGameOver(null)
        return true
      }
      return false
    } catch (e: any) {
      error.value = e.message || 'Failed to restart game'
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    startGame,
    makeChoice,
    selectCharacter,
    loadGame,
    restartGame,
  }
}