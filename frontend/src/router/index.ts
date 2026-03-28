import { createRouter, createWebHistory } from 'vue-router'
import { useGameStore } from '@/stores/game'

const routes = [
  {
    path: '/',
    name: 'start',
    component: () => import('@/features/start/StartView.vue'),
  },
  {
    path: '/character',
    name: 'character',
    component: () => import('@/features/character/CharacterSelectionView.vue'),
  },
  {
    path: '/game',
    name: 'gameplay',
    component: () => import('@/features/gameplay/GameplayView.vue'),
    meta: { requiresGame: true },
  },
  {
    path: '/final',
    name: 'final',
    component: () => import('@/features/final/FinalView.vue'),
    meta: { requiresGame: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const gameStore = useGameStore()
  
  if (to.meta.requiresGame && !gameStore.gameId) {
    next({ name: 'start' })
  } else {
    next()
  }
})

export default router