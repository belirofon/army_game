<script setup lang="ts">
interface Props {
  templateId?: string
}

defineProps<Props>()
</script>

<template>
  <div class="illustration">
    <div class="illustration-placeholder">
      <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
        <circle cx="8.5" cy="8.5" r="1.5"></circle>
        <polyline points="21 15 16 10 5 21"></polyline>
      </svg>
    </div>
  </div>
</template>

<style scoped>
.illustration {
  width: 100%;
  aspect-ratio: 16 / 9;
  background-color: var(--color-surface-elevated);
  border-radius: 0.375rem;
  overflow: hidden;
}

.illustration-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-muted);
  background: linear-gradient(
    135deg,
    var(--color-surface) 0%,
    var(--color-surface-elevated) 50%,
    var(--color-surface) 100%
  );
}
</style>
