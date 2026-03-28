<script setup lang="ts">
import type { PlayerStats } from '@/types/game'
import StatBar from '@/components/stat/StatBar.vue'

defineProps<{
  stats: PlayerStats
}>()
</script>

<template>
  <div class="stats-row">
    <StatBar name="STR" :value="stats.str" :max="100" />
    <StatBar name="END" :value="stats.end" :max="100" />
    <StatBar name="AGI" :value="stats.agi" :max="100" />
    <StatBar 
      name="MOR" 
      :value="stats.mor" 
      :max="100"
      :warning="stats.mor < 30"
    />
    <StatBar name="DSC" :value="stats.disc" :max="100" show-sign />
  </div>
</template>

<style scoped>
.stats-row {
  background-color: #2D3527;
  padding: 0.5rem 1rem;
  display: flex;
  gap: 0.75rem;
  overflow-x: auto;
  border-bottom: 1px solid #3D4533;
}
</style>