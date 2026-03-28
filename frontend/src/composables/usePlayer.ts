import { computed } from 'vue'
import { usePlayerStore } from '@/stores/player'
import type { PlayerStats, Effect } from '@/types/game'
import { STAT_RANGES } from '@/types/game'

export function usePlayer() {
  const store = usePlayerStore()

  const stats = computed(() => store.stats)
  const formalRank = computed(() => store.formalRank)
  const informalStatus = computed(() => store.informalStatus)
  const version = computed(() => store.version)

  function applyEffect(effect: Effect) {
    const clampedValue = Math.max(
      STAT_RANGES[effect.stat].min,
      Math.min(STAT_RANGES[effect.stat].max, effect.newValue)
    )
    store.setStat(effect.stat, clampedValue)
  }

  function applyEffects(effects: Effect[]) {
    effects.forEach(applyEffect)
  }

  function reset() {
    store.reset()
  }

  return {
    stats,
    formalRank,
    informalStatus,
    version,
    applyEffect,
    applyEffects,
    reset,
  }
}