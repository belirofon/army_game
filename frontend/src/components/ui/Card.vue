<script setup lang="ts">
interface Props {
  variant?: 'default' | 'elevated'
}

withDefaults(defineProps<Props>(), {
  variant: 'default',
})
</script>

<template>
  <div class="card" :class="{ 'card-elevated': variant === 'elevated' }">
    <slot />
  </div>
</template>

<style scoped>
.card {
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  padding: 1rem;
}

.card-elevated {
  background-color: var(--color-surface-elevated);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3);
}
</style>
