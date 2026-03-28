<script setup lang="ts">
interface Props {
  size?: 'sm' | 'md' | 'lg'
}

withDefaults(defineProps<Props>(), {
  size: 'md',
})
</script>

<template>
  <div class="spinner-wrapper" :class="`spinner-${size}`">
    <div class="spinner-circle"></div>
  </div>
</template>

<style scoped>
.spinner-wrapper {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.spinner-circle {
  border-radius: 50%;
  border-style: solid;
  border-color: var(--color-military-green);
  border-top-color: transparent;
  animation: spin 0.8s linear infinite;
}

.spinner-sm .spinner-circle {
  width: 1rem;
  height: 1rem;
  border-width: 2px;
}

.spinner-md .spinner-circle {
  width: 1.5rem;
  height: 1.5rem;
  border-width: 3px;
}

.spinner-lg .spinner-circle {
  width: 2.5rem;
  height: 2.5rem;
  border-width: 4px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
