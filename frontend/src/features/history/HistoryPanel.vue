<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '@/stores/game'

const emit = defineEmits<{
  close: []
}>()

const gameStore = useGameStore()
const history = computed(() => gameStore.eventHistory)

function handleOverlayClick(e: MouseEvent) {
  if ((e.target as HTMLElement).classList.contains('history-overlay')) {
    emit('close')
  }
}

function getOutcomeIcon(outcome: string): string {
  if (outcome === 'SUCCESS' || outcome === 'NOTICED_SUCCESS') return '✓'
  if (outcome === 'PARTIAL') return '~'
  return '✗'
}

function getOutcomeClass(outcome: string): string {
  if (outcome === 'SUCCESS' || outcome === 'NOTICED_SUCCESS') return 'history-item__icon--success'
  if (outcome === 'PARTIAL') return 'history-item__icon--partial'
  return 'history-item__icon--failure'
}

function formatMonth(turn: number): string {
  const years = Math.floor(turn / 12)
  const months = turn % 12
  
  if (years === 0) {
    return `${months} мес.`
  } else if (months === 0) {
    return `${years} ${years === 1 ? 'год' : 'лет'}`
  } else {
    return `${years} ${years === 1 ? 'год' : 'лет'}, ${months} мес.`
  }
}
</script>

<template>
  <div class="history-overlay" @click="handleOverlayClick">
    <div class="history-panel">
      <div class="history-panel__header">
        <h2 class="history-panel__title">ИСТОРИЯ СОБЫТИЙ</h2>
        <button class="history-panel__close" @click="emit('close')">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
          </svg>
        </button>
      </div>

      <div class="history-panel__content">
        <div v-if="history.length === 0" class="history-panel__empty">
          Пока нет событий
        </div>

        <div 
          v-for="entry in history" 
          :key="entry.id" 
          class="history-item"
        >
          <div class="history-item__header">
            <span class="history-item__turn">{{ formatMonth(entry.turn) }}</span>
            <span 
              class="history-item__icon"
              :class="getOutcomeClass(entry.checkResult.outcome)"
            >
              {{ getOutcomeIcon(entry.checkResult.outcome) }}
            </span>
          </div>
          <div class="history-item__event">{{ entry.eventDescription }}</div>
          <div class="history-item__choice">&gt; {{ entry.choiceText }}</div>
          <div class="history-item__effects">
            <span 
              v-for="effect in entry.effects" 
              :key="effect.stat"
              class="history-item__effect"
              :class="{ 'history-item__effect--negative': effect.delta < 0 }"
            >
              {{ effect.stat }}: {{ effect.delta > 0 ? '+' : '' }}{{ effect.delta }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.history-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: flex-end;
  z-index: 100;
  animation: fadeIn 0.2s;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.history-panel {
  background-color: #2D3527;
  border-radius: 12px 12px 0 0;
  width: 100%;
  max-height: 70vh;
  display: flex;
  flex-direction: column;
  animation: slideUp 0.3s;
}

@keyframes slideUp {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}

.history-panel__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem;
  border-bottom: 1px solid #3D4533;
}

.history-panel__title {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 1rem;
  color: #E8E4DC;
  letter-spacing: 0.05em;
}

.history-panel__close {
  background: none;
  border: none;
  color: #9A968E;
  cursor: pointer;
  padding: 0.25rem;
}
.history-panel__close:hover {
  color: #E8E4DC;
}
.history-panel__close svg {
  width: 1.5rem;
  height: 1.5rem;
}

.history-panel__content {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
}

.history-panel__empty {
  text-align: center;
  color: #9A968E;
  padding: 2rem;
}

.history-item {
  padding: 0.75rem 0;
  border-bottom: 1px solid #3D4533;
}

.history-item:last-child {
  border-bottom: none;
}

.history-item__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.25rem;
}

.history-item__turn {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 0.875rem;
  color: #E8E4DC;
}

.history-item__icon {
  font-size: 1rem;
  font-weight: bold;
}

.history-item__icon--success {
  color: #5D7A4A;
}

.history-item__icon--partial {
  color: #B8860B;
}

.history-item__icon--failure {
  color: #8B3A3A;
}

.history-item__event {
  font-family: 'Roboto', sans-serif;
  font-size: 0.875rem;
  color: #E8E4DC;
  margin-bottom: 0.25rem;
}

.history-item__choice {
  font-family: 'Roboto', sans-serif;
  font-size: 0.75rem;
  color: #9A968E;
  margin-bottom: 0.5rem;
}

.history-item__effects {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.history-item__effect {
  font-family: 'Roboto Mono', monospace;
  font-size: 0.625rem;
  color: #5D7A4A;
  background-color: rgba(93, 122, 74, 0.2);
  padding: 0.125rem 0.375rem;
  border-radius: 4px;
}

.history-item__effect--negative {
  color: #8B3A3A;
  background-color: rgba(139, 58, 58, 0.2);
}
</style>