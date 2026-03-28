<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  name: string
  value: number
  max: number
  warning?: boolean
  showSign?: boolean
}>()

const percentage = computed(() => Math.round((props.value / props.max) * 100))

const barColor = computed(() => {
  if (props.warning || percentage.value < 30) return 'bg-danger'
  if (percentage.value < 50) return 'bg-warning'
  return 'bg-military-green'
})

const displayValue = computed(() => {
  if (props.showSign && props.value > 0) return `+${props.value}`
  return String(props.value)
})
</script>

<template>
  <div class="stat-bar">
    <span class="stat-bar__name">{{ name }}</span>
    <div class="stat-bar__track">
      <div 
        class="stat-bar__fill" 
        :class="barColor"
        :style="{ width: `${percentage}%` }"
      ></div>
    </div>
    <span class="stat-bar__value">{{ displayValue }}</span>
  </div>
</template>

<style scoped>
.stat-bar {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 80px;
}

.stat-bar__name {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 0.75rem;
  color: #9A968E;
  width: 2rem;
}

.stat-bar__track {
  flex: 1;
  height: 0.5rem;
  background-color: #3D4533;
  border-radius: 9999px;
  overflow: hidden;
}

.stat-bar__fill {
  height: 100%;
  border-radius: 9999px;
  transition: width 0.3s;
}

.stat-bar__fill.bg-danger {
  background-color: #8B3A3A;
}

.stat-bar__fill.bg-warning {
  background-color: #B8860B;
}

.stat-bar__fill.bg-military-green {
  background-color: #4A5D3F;
}

.stat-bar__value {
  font-family: 'Roboto Mono', monospace;
  font-size: 0.75rem;
  color: #E8E4DC;
  width: 1.5rem;
  text-align: right;
}
</style>