<script setup lang="ts">
defineProps<{
  text: string
  disabled?: boolean
}>()

const emit = defineEmits<{
  click: []
}>()
</script>

<template>
  <button
    class="choice-btn"
    :class="{ 'choice-btn--disabled': disabled }"
    :disabled="disabled"
    @click="emit('click')"
  >
    {{ text }}
  </button>
</template>

<style scoped>
.choice-btn {
  width: 100%;
  text-align: left;
  background-color: #2D3527;
  border: 1px solid #3D4533;
  border-radius: 0.5rem;
  padding: 0.75rem 1rem;
  font-family: 'Roboto', sans-serif;
  color: #E8E4DC;
  transition: all 0.1s;
  min-height: 44px;
}
.choice-btn:hover {
  background-color: #353D2E;
}
.choice-btn:active {
  transform: scale(0.98);
}

.choice-btn--disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.choice-btn--disabled:hover {
  background-color: #2D3527;
}
.choice-btn--disabled:active {
  transform: none;
}
</style>