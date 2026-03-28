<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '@/stores/game'
import { usePlayerStore } from '@/stores/player'
import { useApi } from '@/composables/useApi'
import StatusBar from '@/features/gameplay/StatusBar.vue'
import StatsRow from '@/features/gameplay/StatsRow.vue'
import EventCard from '@/components/event/EventCard.vue'
import ChoiceButton from '@/components/event/ChoiceButton.vue'
import StatsPanel from '@/features/stats/StatsPanel.vue'
import HistoryPanel from '@/features/history/HistoryPanel.vue'

const router = useRouter()
const gameStore = useGameStore()
const playerStore = usePlayerStore()
const api = useApi()

const currentEvent = computed(() => gameStore.currentEvent)
const turn = computed(() => gameStore.turn)
const stats = computed(() => playerStore.stats)
const isGameOver = computed(() => gameStore.isGameOver)
const final = computed(() => gameStore.final)

const eventLocation = computed(() => currentEvent.value?.context?.location || '')
const characterImage = computed(() => {
  const selectedChar = gameStore.selectedCharacter
  if (!selectedChar) return null
  
  const characterMap: Record<string, string> = {
    'fatty': 'Fatty_without_bg.png',
    'nerd': 'Nerd_without_bg.png',
    'scrawny': 'Scrawny_without_bg.png',
    'shorty': 'Shorty_without_bg.png',
    'jailbird': 'Jailbird_without_bg.png',
  }
  return characterMap[selectedChar.id] || null
})

const showStats = ref(false)
const showHistory = ref(false)

onMounted(async () => {
  if (!gameStore.gameId) {
    router.push('/')
    return
  }

  const loaded = await api.loadGame()
  if (!loaded) {
    gameStore.setError('Failed to load game')
    return
  }

  if (isGameOver.value || final.value) {
    router.push('/final')
  }
})

async function handleChoiceSelect(choiceId: string) {
  if (api.loading.value) return

  gameStore.setLoading(true)
  const success = await api.makeChoice(choiceId)
  gameStore.setLoading(false)
  
  if (!success) {
    gameStore.setError(api.error.value || 'Failed to make choice')
    return
  }

  if (isGameOver.value || final.value) {
    router.push('/final')
  }
}

function openStats() {
  showStats.value = true
}

function closeStats() {
  showStats.value = false
}

function openHistory() {
  showHistory.value = true
}

function closeHistory() {
  showHistory.value = false
}
</script>

<template>
  <div class="gameplay-screen">
    <StatusBar
      :turn="turn"
      @open-stats="openStats"
      @open-history="openHistory"
    />
    
    <StatsRow :stats="stats" />
    
    <div class="gameplay-screen__content">
      <EventCard
        v-if="currentEvent"
        :description="currentEvent.description"
        :template-id="currentEvent.templateId"
        :location="eventLocation"
        :show-character="!!characterImage"
        :character-image="characterImage || undefined"
      />
      
      <div v-else class="gameplay-screen__no-event">
        <p>Загрузка события...</p>
      </div>
      
      <div class="gameplay-screen__choices">
        <ChoiceButton
          v-for="choice in currentEvent?.choices"
          :key="choice.id"
          :text="choice.text"
          :disabled="!choice.available || api.loading.value"
          @click="handleChoiceSelect(choice.id)"
        />
      </div>
    </div>
    
    <div v-if="api.loading.value" class="gameplay-screen__loading">
      <div class="spinner-lg"></div>
    </div>

    <StatsPanel v-if="showStats" @close="closeStats" />
    <HistoryPanel v-if="showHistory" @close="closeHistory" />
  </div>
</template>

<style scoped>
.gameplay-screen {
  min-height: 100vh;
  background-color: #1A1F16;
  display: flex;
  flex-direction: column;
}

.gameplay-screen__content {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 1rem;
  gap: 1rem;
}

.gameplay-screen__choices {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-top: auto;
}

.gameplay-screen__no-event {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9A968E;
}

.gameplay-screen__loading {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
}

.spinner-lg {
  width: 3rem;
  height: 3rem;
  border: 4px solid #4A5D3F;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>