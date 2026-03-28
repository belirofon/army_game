<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '@/stores/game'
import { usePlayerStore } from '@/stores/player'
import { useApi } from '@/composables/useApi'

const router = useRouter()
const gameStore = useGameStore()
const playerStore = usePlayerStore()
const api = useApi()

interface CharacterStats {
  str: number
  end: number
  agi: number
  mor: number
  disc: number
}

interface Character {
  id: string
  name: string
  image: string
  description: string
  stats: CharacterStats
}

const BASE_STATS = {
  mor: 50,
  disc: 0
}

const characters: Character[] = [
  {
    id: 'fatty',
    name: 'Толстяк',
    image: '/Fatty.png',
    description: 'Бывший студент спортивного факультета. За год до армии набрал 30 кг лишнего веса. В душе добрый, но внешне грозный. Физически сильный, но медленный. Подойдёт тем, кто любит держать удар.',
    stats: { str: 70, end: 60, agi: 30, mor: 50, disc: 0 }
  },
  {
    id: 'nerd',
    name: 'Очкарик',
    image: '/Nerd.png',
    description: 'Ходил в очках с первого класса, программировал с десяти лет. В армии никогда не был, зато читал Генри Форда. Умный, но физически слабый. Найти выход — его главный талант.',
    stats: { str: 30, end: 40, agi: 50, mor: 50, disc: 0 }
  },
  {
    id: 'scrawny',
    name: 'Худой',
    image: '/Scrawny.png',
    description: 'Вечно худой, вечно голодный. Ест в три раза больше, чем ест, но всё равно худой. Быстрый как ртуть, ловкий как обезьяна. Выживает за счёт скорости и смекалки.',
    stats: { str: 35, end: 45, agi: 80, mor: 50, disc: 0 }
  },
  {
    id: 'shorty',
    name: 'Коротышка',
    image: '/Shorty.png',
    description: 'Маленький рост — не недостаток, а преимущество. Все недооценивают маленьких, а зря. Юркий, незаметный, везде пролезет. Компенсирует рост характером.',
    stats: { str: 45, end: 55, agi: 75, mor: 50, disc: 0 }
  },
  {
    id: 'jailbird',
    name: 'Блатной',
    image: '/Jailbird.png',
    description: 'Знает жизнь не понаслышке. Два года отсидел, прежде чем попал в армию. Тяжёлый взгляд, железные нервы. Не боится никого, кроме сержанта с жёстким голосом.',
    stats: { str: 60, end: 70, agi: 45, mor: 50, disc: 0 }
  }
]

const currentIndex = ref(0)
const startX = ref(0)
const isDragging = ref(false)

const currentCharacter = computed(() => characters[currentIndex.value])

function nextCharacter() {
  currentIndex.value = (currentIndex.value + 1) % characters.length
}

function prevCharacter() {
  currentIndex.value = (currentIndex.value - 1 + characters.length) % characters.length
}

function onTouchStart(e: TouchEvent) {
  startX.value = e.touches[0].clientX
  isDragging.value = true
}

function onTouchEnd(e: TouchEvent) {
  if (!isDragging.value) return
  const diff = startX.value - e.changedTouches[0].clientX
  if (Math.abs(diff) > 50) {
    if (diff > 0) nextCharacter()
    else prevCharacter()
  }
  isDragging.value = false
}

const isLoading = ref(false)

async function onSelect() {
  const stats = {
    str: currentCharacter.value.stats.str,
    end: currentCharacter.value.stats.end,
    agi: currentCharacter.value.stats.agi,
    mor: BASE_STATS.mor,
    disc: BASE_STATS.disc
  }

  gameStore.setSelectedCharacter({
    id: currentCharacter.value.id,
    name: currentCharacter.value.name,
    stats
  })
  gameStore.setTurn(1)

  // Update backend with selected character stats
  isLoading.value = true
  const success = await api.selectCharacter(currentCharacter.value.id, stats)
  isLoading.value = false

  if (!success) {
    // Still navigate even if API fails - local state is saved
    console.error('Failed to update character on backend:', api.error.value)
  }

  router.push('/game')
}
</script>

<template>
  <div class="character-select">
    <div class="character-select__header">
      <h1 class="character-select__title">ВЫБЕРИ СОЛДАТА</h1>
      <p class="character-select__subtitle">Кем ты был до армии?</p>
    </div>

    <div 
      class="character-select__slider"
      @touchstart="onTouchStart"
      @touchend="onTouchEnd"
    >
      <button class="character-select__nav character-select__nav--prev" @click="prevCharacter">
        <svg viewBox="0 0 24 24" fill="currentColor">
          <path d="M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12z"/>
        </svg>
      </button>

      <div class="character-select__card">
        <div class="character-select__image-wrapper">
          <img 
            :src="currentCharacter.image" 
            :alt="currentCharacter.name"
            class="character-select__image"
          />
        </div>
        
        <div class="character-select__info">
          <h2 class="character-select__name">{{ currentCharacter.name }}</h2>
          <p class="character-select__desc">{{ currentCharacter.description }}</p>
          
          <div class="character-select__stats">
            <div class="character-select__stat">
              <span class="character-select__stat-label">СИЛА</span>
              <div class="character-select__stat-bar">
                <div class="character-select__stat-fill" :style="{ width: `${currentCharacter.stats.str}%` }"></div>
              </div>
            </div>
            <div class="character-select__stat">
              <span class="character-select__stat-label">ВЫНОС</span>
              <div class="character-select__stat-bar">
                <div class="character-select__stat-fill" :style="{ width: `${currentCharacter.stats.end}%` }"></div>
              </div>
            </div>
            <div class="character-select__stat">
              <span class="character-select__stat-label">ЛОВКОСТЬ</span>
              <div class="character-select__stat-bar">
                <div class="character-select__stat-fill" :style="{ width: `${currentCharacter.stats.agi}%` }"></div>
              </div>
            </div>
            <div class="character-select__stat">
              <span class="character-select__stat-label">ДУХ</span>
              <div class="character-select__stat-bar">
                <div class="character-select__stat-fill" :style="{ width: `${currentCharacter.stats.mor}%` }"></div>
              </div>
            </div>
            <div class="character-select__stat">
              <span class="character-select__stat-label">Дисциплина</span>
              <div class="character-select__stat-bar">
                <div class="character-select__stat-fill" :style="{ width: `${(currentCharacter.stats.disc + 100) / 2}%` }"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <button class="character-select__nav character-select__nav--next" @click="nextCharacter">
        <svg viewBox="0 0 24 24" fill="currentColor">
          <path d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z"/>
        </svg>
      </button>
    </div>

    <div class="character-select__dots">
      <button 
        v-for="(char, idx) in characters" 
        :key="char.id"
        class="character-select__dot"
        :class="{ 'character-select__dot--active': idx === currentIndex }"
        @click="currentIndex = idx"
      ></button>
    </div>

    <button class="character-select__confirm" @click="onSelect">
      ВЫБРАТЬ
    </button>
  </div>
</template>

<style scoped>
.character-select {
  min-height: 100vh;
  background-color: #1A1F16;
  display: flex;
  flex-direction: column;
  padding: 1rem;
}

.character-select__header {
  text-align: center;
  margin-bottom: 1rem;
}

.character-select__title {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 1.5rem;
  color: #E8E4DC;
  letter-spacing: 0.1em;
  margin-bottom: 0.5rem;
}

.character-select__subtitle {
  font-family: 'Roboto', sans-serif;
  font-size: 0.875rem;
  color: #9A968E;
}

.character-select__slider {
  flex: 1;
  display: flex;
  align-items: center;
  position: relative;
  max-width: 400px;
  width: 100%;
  margin: 0 auto;
}

.character-select__card {
  width: 100%;
  display: flex;
  flex-direction: column;
  background-color: #2D3527;
  border-radius: 12px;
  overflow: hidden;
}

.character-select__image-wrapper {
  height: 60vh;
  max-height: 400px;
  background-color: #1A1F16;
  display: flex;
  align-items: center;
  justify-content: center;
}

.character-select__image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.character-select__info {
  padding: 1rem;
  flex: 1;
}

.character-select__name {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 1.5rem;
  color: #E8E4DC;
  text-align: center;
  margin-bottom: 0.75rem;
}

.character-select__desc {
  font-family: 'Roboto', sans-serif;
  font-size: 0.875rem;
  color: #9A968E;
  line-height: 1.5;
  text-align: center;
  margin-bottom: 1rem;
}

.character-select__stats {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.character-select__stat {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.character-select__stat-label {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 0.625rem;
  color: #9A968E;
  width: 3.5rem;
}

.character-select__stat-bar {
  flex: 1;
  height: 0.375rem;
  background-color: #3D4533;
  border-radius: 2px;
}

.character-select__stat-fill {
  height: 100%;
  background-color: #4A5D3F;
  border-radius: 2px;
}

.character-select__nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background-color: rgba(45, 53, 39, 0.9);
  border: none;
  border-radius: 50%;
  width: 3rem;
  height: 3rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 10;
  color: #E8E4DC;
}
.character-select__nav--prev { left: -1rem; }
.character-select__nav--next { right: -1rem; }
.character-select__nav svg { width: 1.5rem; height: 1.5rem; }

.character-select__dots {
  display: flex;
  justify-content: center;
  gap: 0.5rem;
  margin: 1rem 0;
}

.character-select__dot {
  width: 0.5rem;
  height: 0.5rem;
  border-radius: 50%;
  background-color: #3D4533;
  border: none;
  cursor: pointer;
}
.character-select__dot--active {
  background-color: #4A5D3F;
}

.character-select__confirm {
  background-color: #4A5D3F;
  color: #E8E4DC;
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 1rem;
  letter-spacing: 0.1em;
  padding: 1rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  max-width: 400px;
  width: 100%;
  margin: 0 auto;
}
.character-select__confirm:hover {
  background-color: #5A6D4F;
}

@media (max-width: 400px) {
  .character-select__nav { display: none; }
  .character-select__name { font-size: 1.25rem; }
  .character-select__desc { font-size: 0.75rem; }
}
</style>