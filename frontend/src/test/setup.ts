import { setActivePinia, createPinia } from 'pinia'
import { vi } from 'vitest'

export function setupPiniaForTesting() {
  const pinia = createPinia()
  setActivePinia(pinia)
  return pinia
}

export function cleanupPinia() {
  // Pinia automatically cleans up when a new one is set
}

export const mockPlayer = {
  id: 'test-player-id',
  stats: {
    str: 60,
    end: 55,
    agi: 45,
    mor: 50,
    disc: 10,
  },
  formalRank: 'РЯДОВОЙ' as const,
  informalStatus: 'ЗАПАХ' as const,
  turn: 5,
  flags: ['started'],
  isFinished: false,
  version: 1,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
}

export const mockEvent = {
  id: 'event-1',
  templateId: 'training',
  description: 'Утренняя построение',
  resolvedVariables: {},
  choices: [
    { id: 'choice-1', text: 'Стоять смирно', available: true },
    { id: 'choice-2', text: 'Отвлечься', available: true },
  ],
  context: {
    time: 'утро',
    location: 'плац',
    urgency: 'обычный',
  },
}
