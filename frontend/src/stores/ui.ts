import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUiStore = defineStore('ui', () => {
  const statsPanelOpen = ref(false)
  const historyPanelOpen = ref(false)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const toast = ref<{ message: string; type: 'success' | 'error' | 'warning' } | null>(null)

  function openStatsPanel() {
    statsPanelOpen.value = true
  }

  function closeStatsPanel() {
    statsPanelOpen.value = false
  }

  function openHistoryPanel() {
    historyPanelOpen.value = true
  }

  function closeHistoryPanel() {
    historyPanelOpen.value = false
  }

  function setLoading(l: boolean) {
    loading.value = l
  }

  function setError(e: string | null) {
    error.value = e
  }

  function showToast(message: string, type: 'success' | 'error' | 'warning' = 'success') {
    toast.value = { message, type }
    setTimeout(() => {
      toast.value = null
    }, 3000)
  }

  function clearToast() {
    toast.value = null
  }

  return {
    statsPanelOpen,
    historyPanelOpen,
    loading,
    error,
    toast,
    openStatsPanel,
    closeStatsPanel,
    openHistoryPanel,
    closeHistoryPanel,
    setLoading,
    setError,
    showToast,
    clearToast,
  }
})