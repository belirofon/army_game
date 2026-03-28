<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '@/stores/game'
import { usePlayerStore } from '@/stores/player'
import { useApi } from '@/composables/useApi'
import Button from '@/components/ui/Button.vue'

const router = useRouter()
const gameStore = useGameStore()
const playerStore = usePlayerStore()
const api = useApi()

const hasSavedGame = computed(() => !!gameStore.gameId)

async function handleStart() {
  gameStore.setLoading(true)
  try {
    const success = await api.startGame()
    if (success) {
      router.push('/character')
    } else {
      gameStore.setError(api.error.value || 'Failed to start game')
    }
  } catch (e: any) {
    gameStore.setError(e.message || 'Failed to start game')
  } finally {
    gameStore.setLoading(false)
  }
}

async function handleContinue() {
  gameStore.setLoading(true)
  try {
    const success = await api.loadGame()
    if (success) {
      router.push('/game')
    } else {
      gameStore.setError(api.error.value || 'Failed to load game')
    }
  } catch (e: any) {
    gameStore.setError(e.message || 'Failed to load game')
  } finally {
    gameStore.setLoading(false)
  }
}
</script>

<template>
  <div class="start-screen">
    <div class="start-screen__content">
      <h1 class="start-screen__title">АРМЕЙКА</h1>
      <p class="start-screen__subtitle">Narrative Survival Game</p>
      
      <div class="start-screen__actions">
        <Button
          text="Начать игру"
          variant="primary"
          :loading="gameStore.loading"
          @click="handleStart"
        />
        
        <Button
          v-if="hasSavedGame"
          text="Продолжить"
          variant="secondary"
          @click="handleContinue"
        />
      </div>
    </div>
    
    <div class="start-screen__background"></div>
  </div>
</template>

<style scoped>
.start-screen {
  min-height: 100vh;
  background-color: #1A1F16;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.start-screen__content {
  position: relative;
  z-index: 10;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 3rem 1.5rem;
}

.start-screen__title {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 3rem;
  letter-spacing: 0.2em;
  color: #E8E4DC;
  margin-bottom: 1rem;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.8);
}

.start-screen__subtitle {
  font-family: 'Roboto', sans-serif;
  font-size: 0.875rem;
  color: #9A968E;
  margin-bottom: 3rem;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.8);
}

.start-screen__actions {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  width: 100%;
  max-width: 20rem;
}

.start-screen__background {
  position: absolute;
  inset: 0;
  background-image: url('/start_screen_main.png');
  background-size: cover;
  background-position: center;
  opacity: 0.4;
}

@media (max-width: 730px) {
  .start-screen__title {
    font-size: 2rem;
    letter-spacing: 0.15em;
  }

  .start-screen__subtitle {
    font-size: 0.75rem;
    margin-bottom: 2rem;
  }

  .start-screen__background {
    opacity: 0.3;
  }
}

@media (max-width: 400px) {
  .start-screen__title {
    font-size: 1.5rem;
  }

  .start-screen__content {
    padding: 2rem 1rem;
  }
}
</style>