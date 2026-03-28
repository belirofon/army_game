<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '@/stores/game'
import Button from '@/components/ui/Button.vue'

const router = useRouter()
const gameStore = useGameStore()

const final = computed(() => gameStore.final)

const titleColor = computed(() => {
  switch (final.value?.type) {
    case 'СЛОМАННЫЙ_ДЕМБЕЛЬ':
      return 'text-danger'
    case 'УВАЖАЕМЫЙ_ДЕМБЕЛЬ':
      return 'text-military-green'
    default:
      return 'text-accent-tertiary'
  }
})

const serviceTime = computed(() => {
  if (!final.value) return ''
  const months = final.value.achievedOnTurn
  const years = Math.floor(months / 12)
  const remainingMonths = months % 12
  
  if (years === 0) {
    return `${remainingMonths} мес.`
  } else if (remainingMonths === 0) {
    return `${years} ${years === 1 ? 'год' : 'лет'}`
  } else {
    return `${years} ${years === 1 ? 'год' : 'лет'}, ${remainingMonths} мес.`
  }
})

function handleRestart() {
  gameStore.reset()
  router.push('/')
}

function goHome() {
  router.push('/')
}
</script>

<template>
  <div v-if="final" class="final-screen">
    <div class="final-screen__content">
      <h1 class="final-screen__header">РЕЗУЛЬТАТ</h1>
      
      <div class="final-screen__divider"></div>
      
      <h2 :class="['final-screen__title', titleColor]">{{ final.title }}</h2>
      <p v-if="final.subtitle" class="final-screen__subtitle">{{ final.subtitle }}</p>
      
      <div class="final-screen__divider"></div>
      
      <p class="final-screen__description">{{ final.description }}</p>
      
      <div class="final-screen__stats">
        <h3 class="final-screen__stats-title">ИТОГИ:</h3>
        <div class="final-screen__stats-grid">
          <div class="final-screen__stat">
            <span>Служба:</span>
            <span>{{ serviceTime }}</span>
          </div>
          <div class="final-screen__stat">
            <span>MOR:</span>
            <span>{{ final.finalStats.mor }}</span>
          </div>
          <div class="final-screen__stat">
            <span>DISC:</span>
            <span>{{ final.finalStats.disc > 0 ? `+${final.finalStats.disc}` : final.finalStats.disc }}</span>
          </div>
        </div>
      </div>
      
      <Button
        text="ИГРАТЬ СНОВА"
        variant="primary"
        class="final-screen__btn"
        @click="handleRestart"
      />
      
      <Button
        text="НА ГЛАВНУЮ"
        variant="secondary"
        class="final-screen__btn"
        @click="goHome"
      />
    </div>
  </div>
</template>

<style scoped>
.final-screen {
  min-height: 100vh;
  background-color: #1A1F16;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
}

.final-screen__content {
  max-width: 28rem;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.final-screen__header {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 1.5rem;
  letter-spacing: 0.1em;
  color: #E8E4DC;
  margin-bottom: 1rem;
}

.final-screen__divider {
  width: 100%;
  height: 1px;
  background-color: #3D4533;
  margin: 1rem 0;
}

.final-screen__title {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 1.875rem;
  margin-bottom: 0.5rem;
}

.final-screen__title.text-danger {
  color: #8B3A3A;
}

.final-screen__title.text-military-green {
  color: #4A5D3F;
}

.final-screen__title.text-accent-tertiary {
  color: #5D7A4A;
}

.final-screen__subtitle {
  font-family: 'Roboto', sans-serif;
  font-size: 1.125rem;
  color: #9A968E;
  margin-bottom: 1rem;
}

.final-screen__description {
  font-family: 'Roboto', sans-serif;
  color: #E8E4DC;
  line-height: 1.625;
  margin-bottom: 2rem;
}

.final-screen__stats {
  width: 100%;
  margin-bottom: 2rem;
}

.final-screen__stats-title {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 1.125rem;
  color: #E8E4DC;
  margin-bottom: 1rem;
}

.final-screen__stats-grid {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.final-screen__stat {
  display: flex;
  justify-content: space-between;
  color: #9A968E;
}

.final-screen__btn {
  width: 100%;
}

.final-screen__btn + .final-screen__btn {
  margin-top: 0.75rem;
}
</style>