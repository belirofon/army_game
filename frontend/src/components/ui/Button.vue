<script setup lang="ts">
interface Props {
  text: string
  variant?: 'primary' | 'secondary' | 'choice'
  disabled?: boolean
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  variant: 'primary',
  disabled: false,
  loading: false,
})

const emit = defineEmits<{
  click: []
}>()
</script>

<template>
  <button
    class="btn"
    :class="[
      `btn-${variant}`,
      { 'btn-disabled': disabled || loading }
    ]"
    :disabled="disabled || loading"
    @click="emit('click')"
  >
    <span v-if="loading" class="spinner"></span>
    <span v-else>{{ text }}</span>
  </button>
</template>

<style scoped>
.btn {
  font-family: 'Roboto Condensed', sans-serif;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  height: 2.75rem;
  padding: 0 1.5rem;
  border-radius: 0.375rem;
  transition: all 0.1s;
  min-width: 44px;
  min-height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-primary {
  background-color: #4A5D3F;
  color: #E8E4DC;
}
.btn-primary:hover {
  background-color: #5A6D4F;
}
.btn-primary:active {
  background-color: #3A4D2F;
  transform: scale(0.98);
}

.btn-secondary {
  background-color: #2D3527;
  color: #9A968E;
  border: 1px solid #3D4533;
}
.btn-secondary:hover {
  background-color: #353D2E;
}

.btn-choice {
  width: 100%;
  text-align: left;
  background-color: #2D3527;
  border: 1px solid #3D4533;
  padding: 0.75rem 1rem;
  text-align: left;
  height: 3rem;
}
.btn-choice:hover {
  background-color: #353D2E;
}
.btn-choice:active {
  transform: scale(0.98);
}

.btn-disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinner {
  width: 1.25rem;
  height: 1.25rem;
  border: 2px solid #E8E4DC;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>