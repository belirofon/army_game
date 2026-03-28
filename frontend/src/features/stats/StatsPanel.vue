<script setup lang="ts">
import { computed } from 'vue'
import { usePlayerStore } from '@/stores/player'
import { useGameStore } from '@/stores/game'

const emit = defineEmits<{
  close: []
}>()

const playerStore = usePlayerStore()
const gameStore = useGameStore()

const stats = computed(() => playerStore.stats)
const formalRank = computed(() => playerStore.formalRank)
const informalStatus = computed(() => playerStore.informalStatus)
const turn = computed(() => gameStore.turn)

const progressPercent = computed(() => Math.round((turn.value / 24) * 100))

const monthsDisplay = computed(() => {
  const years = Math.floor(turn.value / 12)
  const months = turn.value % 12
  
  if (years === 0) {
    return `${months} мес.`
  } else if (months === 0) {
    return `${years} ${years === 1 ? 'год' : 'лет'}`
  } else {
    return `${years} ${years === 1 ? 'год' : 'лет'}, ${months} мес.`
  }
})

const rankName = computed(() => {
  const names: Record<string, string> = {
    'РЯДОВОЙ': 'Рядовой',
    'ЕФРЕЙТОР': 'Ефрейтор',
    'МЛ_СЕРЖАНТ': 'Мл. сержант',
    'СЕРЖАНТ': 'Сержант',
  }
  return names[formalRank.value] || formalRank.value
})

const statusName = computed(() => {
  const names: Record<string, string> = {
    'ЗАПАХ': 'Запах',
    'ДУХ': 'Дух',
    'СЛОН': 'Слон',
    'ЧЕРПАК': 'Черпак',
    'ДЕД': 'Дед',
    'ДЕМБЕЛЬ': 'Дембель',
  }
  return names[informalStatus.value] || informalStatus.value
})

function handleOverlayClick(e: MouseEvent) {
  if ((e.target as HTMLElement).classList.contains('modal-overlay')) {
    emit('close')
  }
}
</script>

<template>
  <div class="modal-overlay" @click="handleOverlayClick">
    <div class="stats-panel">
      <div class="stats-panel__header">
        <button class="stats-panel__close" @click="emit('close')">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
          </svg>
        </button>
        <h2 class="stats-panel__title">ХАРАКТЕРИСТИКИ</h2>
      </div>

      <div class="stats-panel__section">
        <h3 class="stats-panel__section-title">ФОРМАЛЬНЫЙ СТАТУС</h3>
        <div class="stats-panel__rank-card">{{ rankName }}</div>
      </div>

      <div class="stats-panel__stats-grid">
        <div class="stat-card">
          <span class="stat-card__label">STR</span>
          <div class="stat-card__bar">
            <div class="stat-card__fill" :style="{ width: `${stats.str}%` }"></div>
          </div>
          <span class="stat-card__value">{{ stats.str }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-card__label">END</span>
          <div class="stat-card__bar">
            <div class="stat-card__fill" :style="{ width: `${stats.end}%` }"></div>
          </div>
          <span class="stat-card__value">{{ stats.end }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-card__label">AGI</span>
          <div class="stat-card__bar">
            <div class="stat-card__fill" :style="{ width: `${stats.agi}%` }"></div>
          </div>
          <span class="stat-card__value">{{ stats.agi }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-card__label">MOR</span>
          <div class="stat-card__bar" :class="{ 'stat-card__bar--warning': stats.mor < 30 }">
            <div class="stat-card__fill" :style="{ width: `${stats.mor}%` }"></div>
          </div>
          <span class="stat-card__value">{{ stats.mor }}</span>
        </div>
        <div class="stat-card">
          <span class="stat-card__label">DSC</span>
          <div class="stat-card__bar">
            <div class="stat-card__fill" :style="{ width: `${(stats.disc + 100) / 2}%` }"></div>
          </div>
          <span class="stat-card__value">{{ stats.disc > 0 ? `+${stats.disc}` : stats.disc }}</span>
        </div>
      </div>

      <div class="stats-panel__divider"></div>

      <div class="stats-panel__section">
        <h3 class="stats-panel__section-title">НЕФОРМАЛЬНЫЙ СТАТУС</h3>
        <div class="stats-panel__status-card">{{ statusName }}</div>
      </div>

      <div class="stats-panel__section">
        <h3 class="stats-panel__section-title">ПРОГРЕСС</h3>
        <div class="stats-panel__progress">
          <span>Служба: {{ monthsDisplay }} из 2 лет</span>
          <div class="stats-panel__progress-bar">
            <div class="stats-panel__progress-fill" :style="{ width: `${progressPercent}%` }"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  animation: fadeIn 0.2s;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.stats-panel {
  background-color: #2D3527;
  border-radius: 12px;
  width: 90%;
  max-width: 360px;
  max-height: 90vh;
  overflow-y: auto;
  animation: slideUp 0.3s;
}

@keyframes slideUp {
  from { transform: translateY(20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.stats-panel__header {
  display: flex;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #3D4533;
}

.stats-panel__close {
  background: none;
  border: none;
  color: #9A968E;
  cursor: pointer;
  padding: 0.5rem;
  margin-right: 0.5rem;
}
.stats-panel__close:hover {
  color: #E8E4DC;
}
.stats-panel__close svg {
  width: 1.5rem;
  height: 1.5rem;
}

.stats-panel__title {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 1.125rem;
  color: #E8E4DC;
  letter-spacing: 0.05em;
}

.stats-panel__section {
  padding: 1rem;
}

.stats-panel__section-title {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 0.75rem;
  color: #9A968E;
  letter-spacing: 0.1em;
  margin-bottom: 0.75rem;
}

.stats-panel__rank-card,
.stats-panel__status-card {
  background-color: #1A1F16;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 1.25rem;
  color: #E8E4DC;
  text-align: center;
}

.stats-panel__stats-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 0.5rem;
  padding: 0 1rem;
}

.stat-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
}

.stat-card__label {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 0.625rem;
  color: #9A968E;
}

.stat-card__bar {
  width: 100%;
  height: 0.5rem;
  background-color: #3D4533;
  border-radius: 4px;
  overflow: hidden;
}

.stat-card__fill {
  height: 100%;
  background-color: #4A5D3F;
  border-radius: 4px;
  transition: width 0.3s;
}

.stat-card__bar--warning .stat-card__fill {
  background-color: #8B3A3A;
}

.stat-card__value {
  font-family: 'Roboto Mono', monospace;
  font-size: 0.75rem;
  color: #E8E4DC;
}

.stats-panel__divider {
  height: 1px;
  background-color: #3D4533;
  margin: 0.5rem 0;
}

.stats-panel__progress {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.stats-panel__progress span {
  font-family: 'Roboto', sans-serif;
  font-size: 0.875rem;
  color: #9A968E;
}

.stats-panel__progress-bar {
  height: 0.5rem;
  background-color: #3D4533;
  border-radius: 4px;
  overflow: hidden;
}

.stats-panel__progress-fill {
  height: 100%;
  background-color: #4A5D3F;
  border-radius: 4px;
  transition: width 0.3s;
}
</style>