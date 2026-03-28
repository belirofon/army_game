<script setup lang="ts">
interface Props {
  message: string
  type?: 'success' | 'error' | 'warning'
}

withDefaults(defineProps<Props>(), {
  type: 'success',
})

const emit = defineEmits<{
  close: []
}>()

function handleClose() {
  emit('close')
}
</script>

<template>
  <div class="toast" :class="`toast-${type}`">
    <div class="toast-icon">
      <svg v-if="type === 'success'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="20 6 9 17 4 12"></polyline>
      </svg>
      <svg v-else-if="type === 'error'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="15" y1="9" x2="9" y2="15"></line>
        <line x1="9" y1="9" x2="15" y2="15"></line>
      </svg>
      <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
        <line x1="12" y1="9" x2="12" y2="13"></line>
        <line x1="12" y1="17" x2="12.01" y2="17"></line>
      </svg>
    </div>
    <span class="toast-message">{{ message }}</span>
    <button class="toast-close" @click="handleClose">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
    </button>
  </div>
</template>

<style scoped>
.toast {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-radius: 0.375rem;
  border: 1px solid;
  min-width: 280px;
  max-width: 400px;
}

.toast-success {
  background-color: rgba(93, 122, 74, 0.2);
  border-color: var(--color-success);
  color: var(--color-text-primary);
}

.toast-error {
  background-color: rgba(139, 58, 58, 0.2);
  border-color: var(--color-danger);
  color: var(--color-text-primary);
}

.toast-warning {
  background-color: rgba(184, 134, 11, 0.2);
  border-color: var(--color-warning);
  color: var(--color-text-primary);
}

.toast-icon {
  flex-shrink: 0;
  display: flex;
}

.toast-message {
  flex: 1;
  font-size: 0.875rem;
}

.toast-close {
  flex-shrink: 0;
  background: none;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding: 0.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.15s;
}

.toast-close:hover {
  color: var(--color-text-primary);
}
</style>
