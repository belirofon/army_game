import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { PlayerStats, FormalRank, InformalStatus, Player } from '@/types/game'
import { STAT_RANGES, FORMAL_RANKS, INFORMAL_STATUSES } from '@/types/game'

export const usePlayerStore = defineStore('player', () => {
  const stats = ref<PlayerStats>({
    str: 50,
    end: 50,
    agi: 50,
    mor: 50,
    disc: 0,
  })

  const formalRank = ref<FormalRank>('РЯДОВОЙ')
  const informalStatus = ref<InformalStatus>('ЗАПАХ')
  const version = ref(1)
  const isLoaded = ref(false)

  const computedFormalRank = computed<FormalRank>(() => {
    const disc = stats.value.disc
    if (disc >= 75) return 'СЕРЖАНТ'
    if (disc >= 50) return 'МЛ_СЕРЖАНТ'
    if (disc >= 25) return 'ЕФРЕЙТОР'
    return 'РЯДОВОЙ'
  })

  const computedInformalStatus = computed<InformalStatus>(() => {
    const disc = stats.value.disc
    if (disc <= -90) return 'ДЕМБЕЛЬ'
    if (disc <= -75) return 'ДЕД'
    if (disc <= -50) return 'ЧЕРПАК'
    if (disc <= -25) return 'СЛОН'
    if (disc < 0) return 'ДУХ'
    return 'ЗАПАХ'
  })

  function setStat(stat: keyof PlayerStats, value: number) {
    const range = STAT_RANGES[stat]
    stats.value[stat] = Math.max(range.min, Math.min(range.max, value))
    updateRanks()
  }

  function updateStats(newStats: PlayerStats) {
    Object.keys(newStats).forEach(key => {
      setStat(key as keyof PlayerStats, newStats[key as keyof PlayerStats])
    })
  }

  function setVersion(newVersion: number) {
    version.value = newVersion
  }

  function updateRanks() {
    formalRank.value = computedFormalRank.value
    informalStatus.value = computedInformalStatus.value
  }

  function loadFromPlayer(player: Player) {
    stats.value = { ...player.stats }
    formalRank.value = player.formalRank
    informalStatus.value = player.informalStatus
    version.value = player.version
    isLoaded.value = true
  }

  function reset() {
    stats.value = { str: 50, end: 50, agi: 50, mor: 50, disc: 0 }
    formalRank.value = 'РЯДОВОЙ'
    informalStatus.value = 'ЗАПАХ'
    version.value = 1
    isLoaded.value = false
  }

  return {
    stats,
    formalRank,
    informalStatus,
    version,
    isLoaded,
    setStat,
    updateStats,
    setVersion,
    loadFromPlayer,
    reset,
  }
})