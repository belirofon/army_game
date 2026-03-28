<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  name: string
  value: number
  max?: number
  description?: string
}

const props = withDefaults(defineProps<Props>(), {
  max: 100,
})

const percentage = computed(() => {
  return Math.round((props.value / props.max) * 100)
})

const statusClass = computed(() => {
  if (percentage.value < 30) return 'stat-critical'
  if (percentage.value < 50) return 'stat-warning'
  return 'stat-normal'
})
</script>

<template>
  <div class="stat-card">
    <div class="stat-header">
      <span class="stat-name">{{ name }}</span>
      <span class="stat-value" :class="statusClass">{{ value }}/{{ max }}</span>
    </div>
    <div class="stat-bar-bg">
      <div 
        class="stat-bar-fill" 
        :class="statusClass"
        :style="{ width: `${percentage}%` }"
      ></div>
    </div>
    <p v-if="description" class="stat-description">{{ description }}</p>
  </div>
</template>

<style scoped>
.stat-card {
  padding: 0.75rem;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 0.375rem;
}

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.stat-name {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-secondary);
}

.stat-value {
  font-family: 'Roboto Mono', monospace;
  font-size: 0.875rem;
  font-weight: 500;
}

.stat-normal {
  color: var(--color-text-primary);
}

.stat-warning {
  color: var(--color-warning);
}

.stat-critical {
  color: var(--color-danger);
}

.stat-bar-bg {
  height: 0.5rem;
  background-color: var(--color-background-secondary);
  border-radius: 0.25rem;
  overflow: hidden;
}

.stat-bar-fill {
  height: 100%;
  border-radius: 0.25rem;
  transition: width 0.3s ease;
}

.stat-bar-fill.stat-normal {
  background-color: var(--color-military-green);
}

.stat-bar-fill.stat-warning {
  background-color: var(--color-warning);
}

.stat-bar-fill.stat-critical {
  background-color: var(--color-danger);
}

.stat-description {
  margin-top: 0.5rem;
  font-size: 0.75rem;
  color: var(--color-text-muted);
}
</style>
