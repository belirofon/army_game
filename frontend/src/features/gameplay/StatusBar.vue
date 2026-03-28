<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  turn: number
}>()

const emit = defineEmits<{
  openStats: []
  openHistory: []
}>()

const remainingTime = computed(() => {
  const remaining = 24 - props.turn
  const years = Math.floor(remaining / 12)
  const months = remaining % 12
  
  if (years === 0) {
    return `${months} мес.`
  } else if (months === 0) {
    return `${years} ${years === 1 ? 'год' : 'лет'}`
  } else {
    return `${years} ${years === 1 ? 'год' : 'лет'}, ${months} мес.`
  }
})
</script>

<template>
  <div class="status-bar">
    <button class="status-bar__btn" @click="emit('openStats')">
      <svg class="icon" viewBox="0 0 24 24" fill="currentColor">
        <path d="M3 13h2v-2H3v2zm0 4h2v-2H3v2zm0-8h2V7H3v2zm4 4h14v-2H7v2zm0 4h14v-2H7v2zM7 7v2h14V7H7z"/>
      </svg>
    </button>
    
    <span class="status-bar__title">МЕСЯЦ {{ turn }}</span>
    <span class="status-bar__remaining">Осталось: {{ remainingTime }}</span>
    
    <button class="status-bar__btn" @click="emit('openHistory')">
      <svg class="icon" viewBox="0 0 24 24" fill="currentColor">
        <path d="M13 3c-4.97 0-9 4.03-9 9H1l3.89 3.89.07.14L9 12H6c0-3.87 3.13-7 7-7s7 3.13 7 7-3.13 7-7 7c-1.93 0-3.68-.79-4.94-2.06l-1.42 1.42C8.27 19.99 10.51 21 13 21c4.97 0 9-4.03 9-9s-4.03-9-9-9zm-1 5v5l4.28 2.54.72-1.21-3.5-2.08V8H12z"/>
      </svg>
    </button>
  </div>
</template>

<style scoped>
.status-bar {
  height: 3rem;
  background-color: #252B1E;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1rem;
  border-bottom: 1px solid #3D4533;
}

.status-bar__btn {
  padding: 0.5rem;
  border-radius: 0.375rem;
  background: transparent;
  border: none;
  color: #9A968E;
  cursor: pointer;
  transition: background-color 0.2s, color 0.2s;
}
.status-bar__btn:hover {
  background-color: #2D3527;
  color: #E8E4DC;
}

.status-bar__title {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 1.125rem;
  letter-spacing: 0.05em;
  color: #E8E4DC;
}

.status-bar__remaining {
  font-family: 'Roboto', sans-serif;
  font-size: 0.75rem;
  color: #9A968E;
  margin-left: 0.5rem;
}

.icon {
  width: 1.5rem;
  height: 1.5rem;
}
</style>