# Frontend Implementation Plan — Narrative Survival Game «Армейка»

**Версия:** 1.0  
**Дата:** 2026-03-26  
**Статус:** Production-Ready Implementation Plan  
**Framework:** Vue 3 (Composition API) + TypeScript + Vite + Pinia + Apollo Client  

---

## 1. Project Overview

### 1.1 Technology Stack

| Component | Technology | Version | Justification |
|-----------|------------|---------|----------------|
| Framework | Vue 3 | 3.4+ | Composition API for better logic reuse and TypeScript support |
| Language | TypeScript | 5.0+ | Strict mode for type safety |
| Build Tool | Vite | 5.0+ | Fast HMR, optimized builds |
| State Management | Pinia | 2.1+ | Vue 3 recommended state management, better TypeScript support than Vuex |
| GraphQL Client | @apollo/client | 3.8+ | Full-featured client with cache management |
| GraphQL Code Generator | @graphql-codegen/cli | 6.0+ | Type-safe client generation |
| Testing | Vitest | 1.0+ | Vite-native testing, fast execution |
| Component Testing | @vue/test-utils | 2.4+ | Official Vue testing utilities |
| Styling | Tailwind CSS | 3.4+ | Utility-first, scalable, easy to implement design system |
| Router | Vue Router | 4.2+ | Official Vue router |

### 1.2 Architecture Pattern

**Feature-Based Architecture** with the following principles:

1. **Feature-First Organization** — Code grouped by business feature, not by technical layer
2. **Shared Components** — Reusable UI components in `/components`
3. **Composables** — Business logic extracted into composable functions
4. **API Layer** — GraphQL operations isolated from components

### 1.3 Design System Alignment

The implementation follows the design system from `docs/DESIGN_SPECIFICATION.md`:

| Design Element | Implementation |
|----------------|----------------|
| Colors | CSS variables matching military palette (#1A1F16, #2D3527, #4A5D3F, etc.) |
| Typography | Roboto Condensed for headers, Roboto/Inter for body |
| Spacing | 8px grid system (xs=4, sm=8, md=16, lg=24, xl=32) |
| Components | Match Hoosegow-style layout with military theme |

---

## 2. Project Structure

### 2.1 Directory Structure

```
frontend/
├── public/
│   └── favicon.svg
├── src/
│   ├── assets/
│   │   └── styles/
│   │       ├── main.css          # Global styles, CSS variables
│   │       └── tailwind.css     # Tailwind imports
│   │
│   ├── components/               # Shared components (reusable)
│   │   ├── ui/
│   │   │   ├── Button.vue        # Primary, secondary, choice buttons
│   │   │   ├── Card.vue          # Base card container
│   │   │   ├── Modal.vue         # Modal overlay component
│   │   │   ├── Spinner.vue       # Loading indicator
│   │   │   └── Toast.vue         # Notification component
│   │   ├── stat/
│   │   │   ├── StatBar.vue       # Compact stat indicator (name + bar + value)
│   │   │   └── StatCard.vue      # Detailed stat card for modal
│   │   └── event/
│   │       ├── EventCard.vue     # Event container with illustration
│   │       ├── EventIllustration.vue  # Illustration placeholder
│   │       └── ChoiceButton.vue  # Individual choice button
│   │
│   ├── composables/              # Business logic hooks
│   │   ├── usePlayer.ts          # Player state, stats, progression
│   │   ├── useGame.ts            # Game flow, current event, history
│   │   ├── useEvents.ts          # Event generation, choice handling
│   │   └── useApi.ts             # GraphQL operations wrapper
│   │
│   ├── stores/                   # Pinia stores
│   │   ├── player.ts             # Player stats, rank, status
│   │   ├── game.ts               # Current event, history, game state
│   │   └── ui.ts                 # UI state: modals, loading, errors
│   │
│   ├── api/                      # GraphQL API layer
│   │   ├── client.ts             # Apollo Client setup
│   │   ├── queries.ts            # GraphQL queries
│   │   ├── mutations.ts          # GraphQL mutations
│   │   └── generated/            # Auto-generated types
│   │       ├── types.ts
│   │       └── operations.ts
│   │
│   ├── features/                 # Feature-based modules
│   │   ├── start/
│   │   │   ├── StartScreen.vue   # Main menu
│   │   │   └── StartView.vue      # View wrapper
│   │   ├── gameplay/
│   │   │   ├── GameplayScreen.vue
│   │   │   ├── StatusBar.vue     # Top bar with day + icons
│   │   │   └── StatsRow.vue      # Compact stat display
│   │   ├── stats/
│   │   │   ├── StatsPanel.vue    # Modal with detailed stats
│   │   │   └── StatsPanelView.vue
│   │   ├── history/
│   │   │   ├── HistoryPanel.vue  # Slide-up panel with event list
│   │   │   └── HistoryPanelView.vue
│   │   └── final/
│   │       ├── FinalScreen.vue   # Game over screen
│   │       └── FinalView.vue
│   │
│   ├── router/
│   │   └── index.ts              # Vue Router configuration
│   │
│   ├── types/                    # TypeScript type definitions
│   │   ├── game.ts               # Game-related types
│   │   └── api.ts                # API-related types
│   │
│   ├── App.vue                   # Root component
│   └── main.ts                   # Entry point
│
├── index.html
├── package.json
├── vite.config.ts
├── tsconfig.json
├── tailwind.config.js
├── postcss.config.js
└── env.d.ts
```

### 2.2 File Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Components | PascalCase + .vue | `StartScreen.vue`, `StatBar.vue` |
| Composables | camelCase + .ts | `usePlayer.ts`, `useGame.ts` |
| Stores | camelCase + .ts | `player.ts`, `game.ts` |
| Types | PascalCase + .ts | `game.ts`, `api.ts` |
| Constants | SCREAMING_SNAKE_CASE | `STAT_RANGES`, `FORMAL_RANKS` |

---

## 3. Step-by-Step Implementation Plan

### Phase 1: Foundation Setup (Steps 1-3)

#### Step 1: Project Initialization and Configuration

**Task:** Set up Vue 3 project with TypeScript, Vite, Tailwind, and core dependencies.

**TDD Workflow:**

1. **Write Tests (before implementation):**
   ```typescript
   // tests/unit/setup.test.ts
   import { describe, it, expect } from 'vitest';
   import { render, screen } from '@testing-library/vue';
   import App from '@/App.vue';

   describe('App', () => {
     it('should render without errors', () => {
       const { container } = render(App);
       expect(container).toBeTruthy();
     });

     it('should have correct title', () => {
       const { getByText } = render(App);
       expect(getByText('АРМЕЙКА')).toBeTruthy();
     });
   });
   ```

2. **Run Tests:** `npm test -- tests/unit/setup.test.ts` → Expected to FAIL (App doesn't exist)

3. **Implement Minimal Logic:**
   - Initialize Vite project: `npm create vite@latest frontend -- --template vue-ts`
   - Install dependencies: `npm install pinia @apollo/client graphql vue-router@4`
   - Install dev dependencies: `npm install -D tailwindcss postcss autoprefixer vitest @vue/test-utils @testing-library/vue jsdom`
   - Configure Tailwind with design system colors
   - Create basic App.vue with title

4. **Refactor:** Organize imports, ensure clean component structure

**Verification:** Run tests again → PASS

---

#### Step 2: GraphQL Client Setup and Type Generation

**Task:** Configure Apollo Client with auto-generated types.

**TDD Workflow:**

1. **Write Tests:**
   ```typescript
   // tests/unit/api/client.test.ts
   import { describe, it, expect, vi, beforeEach } from 'vitest';
   import { ApolloClient, InMemoryCache } from '@apollo/client/core';

   describe('Apollo Client', () => {
     let client: ApolloClient<any>;

     beforeEach(() => {
       client = new ApolloClient({
         uri: 'http://localhost:8080/graphql',
         cache: new InMemoryCache(),
       });
     });

     it('should be defined', () => {
       expect(client).toBeDefined();
     });

     it('should have correct uri', () => {
       expect(client.link.uri).toBe('http://localhost:8080/graphql');
     });
   });
   ```

2. **Run Tests:** FAIL (client not implemented)

3. **Implement:**
   ```typescript
   // src/api/client.ts
   import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core';
   import { onError } from '@apollo/client/link/error';

   const httpLink = createHttpLink({
     uri: import.meta.env.VITE_GRAPHQL_URL || 'http://localhost:8080/graphql',
   });

   const errorLink = onError(({ graphQLErrors }) => {
     if (graphQLErrors) {
       graphQLErrors.forEach(({ message, locations, path }) => {
         console.error(`[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`);
       });
     }
   });

   export const apolloClient = new ApolloClient({
     link: errorLink.concat(httpLink),
     cache: new InMemoryCache(),
     defaultOptions: {
       watchQuery: { fetchPolicy: 'network-only' },
       query: { fetchPolicy: 'network-only' },
     },
   });
   ```

4. **Run Tests:** PASS

**Continue:** Set up GraphQL code generation with `@graphql-codegen/cli`

---

#### Step 3: Pinia Store Setup

**Task:** Create stores for player, game, and UI state.

**TDD Workflow:**

1. **Write Tests:**
   ```typescript
   // tests/unit/stores/player.test.ts
   import { describe, it, expect, beforeEach } from 'vitest';
   import { setActivePinia, createPinia } from 'pinia';
   import { usePlayerStore } from '@/stores/player';

   describe('Player Store', () => {
     beforeEach(() => {
       setActivePinia(createPinia());
     });

     it('should have default stats', () => {
       const store = usePlayerStore();
       expect(store.stats.str).toBe(50);
       expect(store.stats.end).toBe(50);
       expect(store.stats.agi).toBe(50);
       expect(store.stats.mor).toBe(50);
       expect(store.stats.disc).toBe(0);
     });

     it('should clamp stats to valid ranges', () => {
       const store = usePlayerStore();
       store.setStat('mor', -10);
       expect(store.stats.mor).toBe(0); // Clamped to min
      
       store.setStat('str', 150);
       expect(store.stats.str).toBe(100); // Clamped to max
     });

     it('should calculate formal rank correctly', () => {
       const store = usePlayerStore();
       store.stats.disc = 60;
       expect(store.formalRank).toBe('МЛ_СЕРЖАНТ');
     });
   });
   ```

2. **Run Tests:** FAIL (store doesn't exist)

3. **Implement:**
   ```typescript
   // src/stores/player.ts
   import { defineStore } from 'pinia';
   import { ref, computed } from 'vue';
   import type { PlayerStats, FormalRank, InformalStatus } from '@/types/game';

   const STAT_RANGES = {
     str: { min: 1, max: 100 },
     end: { min: 1, max: 100 },
     agi: { min: 1, max: 100 },
     mor: { min: 0, max: 100 },
     disc: { min: -100, max: 100 },
   };

   const FORMAL_RANKS: Record<FormalRank, number> = {
     'РЯДОВОЙ': -Infinity,
     'ЕФРЕЙТОР': 25,
     'МЛ_СЕРЖАНТ': 50,
     'СЕРЖАНТ': 75,
   };

   export const usePlayerStore = defineStore('player', () => {
     const stats = ref<PlayerStats>({
       str: 50,
       end: 50,
       agi: 50,
       mor: 50,
       disc: 0,
     });

     const formalRank = computed<FormalRank>(() => {
       const disc = stats.value.disc;
       if (disc >= 75) return 'СЕРЖАНТ';
       if (disc >= 50) return 'МЛ_СЕРЖАНТ';
       if (disc >= 25) return 'ЕФРЕЙТОР';
       return 'РЯДОВОЙ';
     });

     const informalStatus = computed< InformalStatus>(() => {
       const disc = stats.value.disc;
       if (disc <= -90) return 'ДЕМБЕЛЬ';
       if (disc <= -75) return 'ДЕД';
       if (disc <= -50) return 'ЧЕРПАК';
       if (disc <= -25) return 'СЛОН';
       if (disc < 0) return 'ДУХ';
       return 'ЗАПАХ';
     });

     function setStat(stat: keyof PlayerStats, value: number) {
       const range = STAT_RANGES[stat];
       stats.value[stat] = Math.max(range.min, Math.min(range.max, value));
     }

     function updateStats(newStats: PlayerStats) {
       Object.keys(newStats).forEach(key => {
         setStat(key as keyof PlayerStats, newStats[key as keyof PlayerStats]);
       });
     }

     function reset() {
       stats.value = { str: 50, end: 50, agi: 50, mor: 50, disc: 0 };
     }

     return {
       stats,
       formalRank,
       informalStatus,
       setStat,
       updateStats,
       reset,
     };
   });
   ```

4. **Run Tests:** PASS

**Verification:** Run full test suite → All PASS

---

### Phase 2: Core Components (Steps 4-6)

#### Step 4: UI Components Library

**Task:** Create reusable UI components matching design system.

**Components to Implement (TDD for each):**

**4.1 Button Component**

```typescript
// tests/unit/components/Button.test.ts
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/vue';
import Button from '@/components/ui/Button.vue';

describe('Button', () => {
  it('should render text', () => {
    render(Button, { props: { text: 'Начать игру' } });
    expect(screen.getByText('Начать игру')).toBeTruthy();
  });

  it('should emit click event', async () => {
    const { emitted } = render(Button, { props: { text: 'Test' } });
    await fireEvent.click(screen.getByText('Test'));
    expect(emitted().click).toBeTruthy();
  });

  it('should be disabled when disabled prop is true', () => {
    render(Button, { props: { text: 'Test', disabled: true } });
    const button = screen.getByText('Test') as HTMLButtonElement;
    expect(button.disabled).toBe(true);
  });

  it('should apply variant classes', () => {
    const { container } = render(Button, { props: { text: 'Primary', variant: 'primary' } });
    expect(container.querySelector('.btn-primary')).toBeTruthy();
  });
});
```

**Implementation:**
```vue
<!-- src/components/ui/Button.vue -->
<script setup lang="ts">
interface Props {
  text: string;
  variant?: 'primary' | 'secondary' | 'choice';
  disabled?: boolean;
  loading?: boolean;
}

withDefaults(defineProps<Props>(), {
  variant: 'primary',
  disabled: false,
  loading: false,
});

const emit = defineEmits<{
  click: [];
}>();
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
  @apply font-roboto-condensed font-bold uppercase tracking-wider;
  @apply h-11 px-6 rounded-md transition-all duration-100;
  @apply min-w-[44px] touch-target;
}

.btn-primary {
  @apply bg-military-green text-light-gray;
  @apply hover:bg-military-green-light active:bg-military-green-dark;
  @apply active:scale-[0.98];
}

.btn-secondary {
  @apply bg-surface text-secondary border border-border;
  @apply hover:bg-surface-elevated;
}

.btn-choice {
  @apply w-full text-left bg-surface border border-border;
  @apply hover:bg-surface-elevated active:scale-[0.98];
  @apply h-12 px-4;
}

.btn-disabled {
  @apply opacity-50 cursor-not-allowed;
}
</style>
```

**4.2 StatBar Component**

```typescript
// tests/unit/components/StatBar.test.ts
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/vue';
import StatBar from '@/components/stat/StatBar.vue';

describe('StatBar', () => {
  it('should display stat name and value', () => {
    render(StatBar, { props: { name: 'STR', value: 75, max: 100 } });
    expect(screen.getByText('STR')).toBeTruthy();
    expect(screen.getByText('75')).toBeTruthy();
  });

  it('should show warning color when value < 50', () => {
    const { container } = render(StatBar, { props: { name: 'MOR', value: 30, max: 100 } });
    const bar = container.querySelector('.stat-bar-fill');
    expect(bar?.classList.contains('bg-warning')).toBe(true);
  });

  it('should show danger color when value < 30', () => {
    const { container } = render(StatBar, { props: { name: 'MOR', value: 20, max: 100 } });
    const bar = container.querySelector('.stat-bar-fill');
    expect(bar?.classList.contains('bg-danger')).toBe(true);
  });
});
```

**4.3 Modal Component**

```typescript
// tests/unit/components/Modal.test.ts
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/vue';
import Modal from '@/components/ui/Modal.vue';

describe('Modal', () => {
  it('should render slot content', () => {
    render(Modal, { slots: { default: 'Modal Content' } });
    expect(screen.getByText('Modal Content')).toBeTruthy();
  });

  it('should emit close event on backdrop click', async () => {
    const { emitted } = render(Modal, { props: { show: true } });
    await fireEvent.click(document.querySelector('.modal-overlay')!);
    expect(emitted().close).toBeTruthy();
  });
});
```

**Verify:** Run all component tests → PASS

---

#### Step 5: Composables Implementation

**Task:** Create composables for business logic.

**5.1 usePlayer Composable**

```typescript
// tests/unit/composables/usePlayer.test.ts
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { usePlayer } from '@/composables/usePlayer';

describe('usePlayer', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('should initialize with default stats', () => {
    const { stats } = usePlayer();
    expect(stats.value.str).toBe(50);
    expect(stats.value.mor).toBe(50);
  });

  it('should apply stat changes with clamping', () => {
    const { applyEffect } = usePlayer();
    applyEffect({ stat: 'MOR', delta: -60, previousValue: 50, newValue: -10 });
    expect(usePlayer().stats.value.mor).toBe(0);
  });

  it('should correctly calculate formal rank', () => {
    const { stats, formalRank } = usePlayer();
    stats.value.disc = 80;
    expect(formalRank.value).toBe('СЕРЖАНТ');
  });
});
```

**Implementation:**
```typescript
// src/composables/usePlayer.ts
import { computed } from 'vue';
import { usePlayerStore } from '@/stores/player';
import type { PlayerStats, FormalRank, InformalStatus, Effect } from '@/types/game';

export function usePlayer() {
  const store = usePlayerStore();

  const stats = computed(() => store.stats);
  const formalRank = computed(() => store.formalRank);
  const informalStatus = computed(() => store.informalStatus);

  function applyEffect(effect: Effect) {
    const clampedValue = Math.max(
      getStatMin(effect.stat),
      Math.min(getStatMax(effect.stat), effect.newValue)
    );
    store.setStat(effect.stat as keyof PlayerStats, clampedValue);
  }

  function applyEffects(effects: Effect[]) {
    effects.forEach(applyEffect);
  }

  function reset() {
    store.reset();
  }

  return {
    stats,
    formalRank,
    informalStatus,
    applyEffect,
    applyEffects,
    reset,
  };
}

function getStatMin(stat: string): number {
  const ranges: Record<string, number> = { STR: 1, END: 1, AGI: 1, MOR: 0, DISC: -100 };
  return ranges[stat] ?? 0;
}

function getStatMax(stat: string): number {
  const ranges: Record<string, number> = { STR: 100, END: 100, AGI: 100, MOR: 100, DISC: 100 };
  return ranges[stat] ?? 100;
}
```

**5.2 useGame Composable**

```typescript
// tests/unit/composables/useGame.test.ts
import { describe, it, expect, beforeEach } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useGameStore } from '@/stores/game';
import { useGame } from '@/composables/useGame';

describe('useGame', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('should return current turn', () => {
    const { turn } = useGame();
    expect(turn.value).toBe(1);
  });

  it('should determine game over when mor is 0', () => {
    const store = useGameStore();
    store.$patch({ playerStats: { mor: 0 } });
    const { isGameOver } = useGame();
    expect(isGameOver.value).toBe(true);
  });
});
```

**Verify:** Run all composable tests → PASS

---

#### Step 6: GraphQL Queries and Mutations

**Task:** Implement GraphQL operations with type generation.

**6.1 Query Tests**

```typescript
// tests/unit/api/queries.test.ts
import { describe, it, expect } from 'vitest';
import { gql } from '@apollo/client';

describe('GraphQL Queries', () => {
  it('should have player query defined', () => {
    expect(gql`query GetPlayer($id: ID!) { player(id: $id) { id } }`).toBeTruthy();
  });

  it('should have currentEvent query defined', () => {
    expect(gql`query GetCurrentEvent($playerId: ID!) { currentEvent(playerId: $playerId) { id } }`).toBeTruthy();
  });
});

describe('GraphQL Mutations', () => {
  it('should have startGame mutation defined', () => {
    expect(gql`mutation StartGame { startGame { gameId } }`).toBeTruthy();
  });

  it('should have choose mutation with all fields', () => {
    expect(gql`mutation Choose($playerId: ID!, $choiceId: ID!, $expectedVersion: Int!) { choose(playerId: $playerId, choiceId: $choiceId, expectedVersion: $expectedVersion) { success } }`).toBeTruthy();
  });
});
```

**Implementation:**
```typescript
// src/api/queries.ts
import { gql } from '@apollo/client';

export const GET_PLAYER = gql`
  query GetPlayer($id: ID!) {
    player(id: $id) {
      id
      stats {
        str
        end
        agi
        mor
        disc
      }
      formalRank
      informalStatus
      turn
      flags
      isFinished
      version
    }
  }
`;

export const GET_CURRENT_EVENT = gql`
  query GetCurrentEvent($playerId: ID!) {
    currentEvent(playerId: $playerId) {
      id
      templateId
      description
      resolvedVariables
      choices {
        id
        text
        available
      }
      context {
        time
        location
        urgency
      }
    }
  }
`;

export const GET_EVENT_HISTORY = gql`
  query GetEventHistory($playerId: ID!, $limit: Int = 10) {
    eventHistory(playerId: $playerId, limit: $limit) {
      id
      turn
      eventDescription
      choiceText
      checkResult {
        success
        outcome
        description
      }
      effects {
        stat
        delta
      }
    }
  }
`;
```

```typescript
// src/api/mutations.ts
import { gql } from '@apollo/client';

export const START_GAME = gql`
  mutation StartGame {
    startGame {
      gameId
      player {
        id
        stats {
          str
          end
          agi
          mor
          disc
        }
        formalRank
        informalStatus
        turn
      }
      currentEvent {
        id
        description
        choices {
          id
          text
          available
        }
      }
      isGameOver
    }
  }
`;

export const CHOOSE = gql`
  mutation Choose($playerId: ID!, $choiceId: ID!, $expectedVersion: Int!) {
    choose(playerId: $playerId, choiceId: $choiceId, expectedVersion: $expectedVersion) {
      success
      checkResult {
        success
        outcome
        description
      }
      effects {
        stat
        delta
        previousValue
        newValue
      }
      updatedPlayer {
        id
        stats {
          str
          end
          agi
          mor
          disc
        }
        formalRank
        informalStatus
        turn
        version
      }
      nextEvent {
        id
        description
        choices {
          id
          text
          available
        }
      }
      gameOver
      final {
        type
        title
        subtitle
        description
        finalStats {
          str
          end
          agi
          mor
          disc
        }
        achievedOnTurn
      }
      newVersion
    }
  }
`;

export const RESTART_GAME = gql`
  mutation RestartGame($playerId: ID!) {
    restartGame(playerId: $playerId) {
      gameId
      player {
        id
        stats {
          str
          end
          agi
          mor
          disc
        }
        turn
      }
      currentEvent {
        id
        description
        choices {
          id
          text
          available
        }
      }
      isGameOver
    }
  }
`;
```

**Verify:** Run API tests → PASS

---

### Phase 3: Feature Screens (Steps 7-10)

#### Step 7: Start Screen Implementation

**TDD Tests:**

```typescript
// tests/unit/features/start/StartScreen.test.ts
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/vue';
import { createRouter, createWebHistory } from 'vue-router';
import StartScreen from '@/features/start/StartScreen.vue';
import { setActivePinia, createPinia } from 'pinia';

const router = createRouter({
  history: createWebHistory(),
  routes: [{ path: '/', name: 'start' }],
});

describe('StartScreen', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('should display game title', () => {
    render(StartScreen, { global: { plugins: [router] } });
    expect(screen.getByText('АРМЕЙКА')).toBeTruthy();
  });

  it('should show "Начать игру" button', () => {
    render(StartScreen, { global: { plugins: [router] } });
    expect(screen.getByText('Начать игру')).toBeTruthy();
  });

  it('should not show "Продолжить" when no saved game', () => {
    render(StartScreen, { global: { plugins: [router] } });
    expect(screen.queryByText('Продолжить')).toBeNull();
  });

  it('should emit startGame event on button click', async () => {
    const { emitted } = render(StartScreen, { global: { plugins: [router] } });
    await fireEvent.click(screen.getByText('Начать игру'));
    expect(emitted().startGame).toBeTruthy();
  });
});
```

**Implementation:**
```vue
<!-- src/features/start/StartScreen.vue -->
<script setup lang="ts">
import { computed } from 'vue';
import { usePlayerStore } from '@/stores/player';
import { useGameStore } from '@/stores/game';
import Button from '@/components/ui/Button.vue';

const emit = defineEmits<{
  startGame: [];
  continueGame: [];
}>();

const playerStore = usePlayerStore();
const gameStore = useGameStore();

const hasSavedGame = computed(() => !!gameStore.gameId);

function handleStart() {
  emit('startGame');
}

function handleContinue() {
  emit('continueGame');
}
</script>

<template>
  <div class="start-screen">
    <div class="start-screen__content">
      <h1 class="start-screen__title">АРМЕЙКА</h1>
      <p class="start-screen__subtitle">Narrative Survival Game</p>
      
      <div class="start-screen__actions">
        <Button
          text="Начать игру"
          variant="primary"
          class="start-screen__btn-primary"
          @click="handleStart"
        />
        
        <Button
          v-if="hasSavedGame"
          text="Продолжить"
          variant="secondary"
          class="start-screen__btn-secondary"
          @click="handleContinue"
        />
      </div>
    </div>
    
    <div class="start-screen__background"></div>
  </div>
</template>

<style scoped>
.start-screen {
  @apply min-h-screen bg-background-primary flex items-center justify-center;
  @apply relative overflow-hidden;
}

.start-screen__content {
  @apply relative z-10 flex flex-col items-center;
  @apply px-6 py-12;
}

.start-screen__title {
  @apply font-roboto-condensed font-bold text-5xl tracking-widest;
  @apply text-text-primary mb-4;
  text-transform: uppercase;
}

.start-screen__subtitle {
  @apply font-roboto text-text-secondary text-sm mb-12;
}

.start-screen__actions {
  @apply flex flex-col gap-4 w-full max-w-xs;
}

.start-screen__background {
  @apply absolute inset-0;
  background: linear-gradient(135deg, #1A1F16 0%, #252B1E 50%, #1A1F16 100%);
}
</style>
```

---

#### Step 8: Gameplay Screen Implementation

**TDD Tests:**

```typescript
// tests/unit/features/gameplay/GameplayScreen.test.ts
import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/vue';
import { setActivePinia, createPinia } from 'pinia';
import { useGameStore } from '@/stores/game';
import GameplayScreen from '@/features/gameplay/GameplayScreen.vue';

describe('GameplayScreen', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    const gameStore = useGameStore();
    gameStore.$patch({
      currentEvent: {
        id: 'test-event',
        description: 'Сержант проверяет построение',
        choices: [
          { id: 'choice-1', text: 'Стоять смирно', available: true },
          { id: 'choice-2', text: 'Поправить ремень', available: true },
        ],
      },
    });
  });

  it('should display current day', () => {
    const { getByText } = render(GameplayScreen);
    expect(getByText('ДЕНЬ 1')).toBeTruthy();
  });

  it('should display event description', () => {
    const { getByText } = render(GameplayScreen);
    expect(getByText('Сержант проверяет построение')).toBeTruthy();
  });

  it('should display all choices', () => {
    const { getByText } = render(GameplayScreen);
    expect(getByText('Стоять смирно')).toBeTruthy();
    expect(getByText('Поправить ремень')).toBeTruthy();
  });

  it('should emit selectChoice event when choice clicked', async () => {
    const { emitted } = render(GameplayScreen);
    await fireEvent.click(screen.getByText('Стоять смирно'));
    expect(emitted().selectChoice).toBeTruthy();
  });
});
```

**Implementation:**
```vue
<!-- src/features/gameplay/GameplayScreen.vue -->
<script setup lang="ts">
import { computed } from 'vue';
import { useGameStore } from '@/stores/game';
import { usePlayerStore } from '@/stores/player';
import StatusBar from '@/features/gameplay/StatusBar.vue';
import StatsRow from '@/features/gameplay/StatsRow.vue';
import EventCard from '@/components/event/EventCard.vue';
import ChoiceButton from '@/components/event/ChoiceButton.vue';

const emit = defineEmits<{
  selectChoice: [choiceId: string];
  openStats: [];
  openHistory: [];
}>();

const gameStore = useGameStore();
const playerStore = usePlayerStore();

const currentEvent = computed(() => gameStore.currentEvent);
const turn = computed(() => gameStore.turn);
const stats = computed(() => playerStore.stats);

function handleChoiceSelect(choiceId: string) {
  if (!gameStore.loading) {
    emit('selectChoice', choiceId);
  }
}
</script>

<template>
  <div class="gameplay-screen">
    <StatusBar
      :turn="turn"
      @open-stats="emit('openStats')"
      @open-history="emit('openHistory')"
    />
    
    <StatsRow :stats="stats" />
    
    <div class="gameplay-screen__content">
      <EventCard
        v-if="currentEvent"
        :description="currentEvent.description"
        :illustration="currentEvent.templateId"
      />
      
      <div class="gameplay-screen__choices">
        <ChoiceButton
          v-for="choice in currentEvent?.choices"
          :key="choice.id"
          :text="choice.text"
          :disabled="!choice.available || gameStore.loading"
          @click="handleChoiceSelect(choice.id)"
        />
      </div>
    </div>
    
    <div v-if="gameStore.loading" class="gameplay-screen__loading">
      <Spinner />
    </div>
  </div>
</template>

<style scoped>
.gameplay-screen {
  @apply min-h-screen bg-background-primary flex flex-col;
}

.gameplay-screen__content {
  @apply flex-1 flex flex-col p-4 gap-4;
}

.gameplay-screen__choices {
  @apply flex flex-col gap-3 mt-auto;
}

.gameplay-screen__loading {
  @apply fixed inset-0 bg-black/50 flex items-center justify-center;
  @apply z-50;
}
</style>
```

---

#### Step 9: Stats and History Panels

**TDD Tests for StatsPanel:**

```typescript
// tests/unit/features/stats/StatsPanel.test.ts
import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/vue';
import { setActivePinia, createPinia } from 'pinia';
import { usePlayerStore } from '@/stores/player';
import { useGameStore } from '@/stores/game';
import StatsPanel from '@/features/stats/StatsPanel.vue';

describe('StatsPanel', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('should display all 5 stats', () => {
    render(StatsPanel);
    expect(screen.getByText('STR')).toBeTruthy();
    expect(screen.getByText('END')).toBeTruthy();
    expect(screen.getByText('AGI')).toBeTruthy();
    expect(screen.getByText('MOR')).toBeTruthy();
    expect(screen.getByText('DISC')).toBeTruthy();
  });

  it('should display formal rank', () => {
    const store = usePlayerStore();
    store.stats.disc = 60;
    const { getByText } = render(StatsPanel);
    expect(getByText('МЛ.СЕРЖАНТ')).toBeTruthy();
  });

  it('should display turn progress', () => {
    const gameStore = useGameStore();
    gameStore.$patch({ turn: 15 });
    const { getByText } = render(StatsPanel);
    expect(getByText('День 15 из 30')).toBeTruthy();
  });
});
```

**TDD Tests for HistoryPanel:**

```typescript
// tests/unit/features/history/HistoryPanel.test.ts
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/vue';
import HistoryPanel from '@/features/history/HistoryPanel.vue';

describe('HistoryPanel', () => {
  it('should display event history', () => {
    render(HistoryPanel, {
      props: {
        history: [
          {
            turn: 5,
            eventDescription: 'Проверка',
            choiceText: 'Принять позу',
            checkResult: { success: true, outcome: 'SUCCESS', description: 'Успех' },
            effects: [{ stat: 'MOR', delta: -2 }],
          },
        ],
      },
    });
    expect(screen.getByText('День 5')).toBeTruthy();
    expect(screen.getByText('Проверка')).toBeTruthy();
  });
});
```

---

#### Step 10: Final Screen

**TDD Tests:**

```typescript
// tests/unit/features/final/FinalScreen.test.ts
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/vue';
import FinalScreen from '@/features/final/FinalScreen.vue';

describe('FinalScreen', () => {
  it('should display тихий дембель final', () => {
    render(FinalScreen, {
      props: {
        final: {
          type: 'ТИХИЙ_ДЕМБЕЛЬ',
          title: 'Тихий дембель',
          subtitle: 'Приспособленец',
          description: 'Ваша служба прошла неспеша...',
          finalStats: { str: 50, end: 50, agi: 50, mor: 30, disc: 10 },
          achievedOnTurn: 30,
        },
      },
    });
    expect(screen.getByText('Тихий дембель')).toBeTruthy();
    expect(screen.getByText('Приспособленец')).toBeTruthy();
  });

  it('should display уважаемый дембель final', () => {
    render(FinalScreen, {
      props: {
        final: {
          type: 'УВАЖАЕМЫЙ_ДЕМБЕЛЬ',
          title: 'Уважаемый дембель',
          description: 'Вы ушли со службы в золотом кашне...',
          finalStats: { str: 80, end: 75, agi: 60, mor: 55, disc: -60 },
          achievedOnTurn: 28,
        },
      },
    });
    expect(screen.getByText('Уважаемый дембель')).toBeTruthy();
  });

  it('should display сломанный дембель final', () => {
    render(FinalScreen, {
      props: {
        final: {
          type: 'СЛОМАННЫЙ_ДЕМБЕЛЬ',
          title: 'Сломанный дембель',
          subtitle: 'Косячник',
          description: 'Служба в армии принесла вам лишь горе...',
          finalStats: { str: 20, end: 15, agi: 30, mor: 0, disc: -80 },
          achievedOnTurn: 12,
        },
      },
    });
    expect(screen.getByText('Сломанный дембель')).toBeTruthy();
    expect(screen.getByText('Косячник')).toBeTruthy();
  });

  it('should emit restart event on button click', async () => {
    const { emitted } = render(FinalScreen, {
      props: {
        final: {
          type: 'ТИХИЙ_ДЕМБЕЛЬ',
          title: 'Test',
          description: 'Test',
          finalStats: { str: 50, end: 50, agi: 50, mor: 50, disc: 0 },
          achievedOnTurn: 30,
        },
      },
    });
    const { fireEvent } = await import('@testing-library/vue');
    await fireEvent.click(screen.getByText('ИГРАТЬ СНОВА'));
    expect(emitted().restart).toBeTruthy();
  });
});
```

---

## 4. Component Architecture

### 4.1 Component Hierarchy

```
App
├── RouterView
│   ├── StartView
│   │   └── StartScreen
│   │       └── Button (x1-2)
│   │
│   ├── GameplayView
│   │   └── GameplayScreen
│   │       ├── StatusBar
│   │       │   ├── IconButton (x2)
│   │       │   └── DayIndicator
│   │       ├── StatsRow
│   │       │   └── StatBar (x5)
│   │       ├── EventCard
│   │       │   ├── EventIllustration
│   │       │   └── EventDescription
│   │       └── ChoiceButton (x2-4)
│   │
│   ├── StatsPanelView (Modal)
│   │   └── StatsPanel
│   │       ├── Modal
│   │       ├── FormalRankCard
│   │       ├── StatCard (x5)
│   │       ├── InformalStatusCard
│   │       └── ProgressBar
│   │
│   ├── HistoryPanelView (Bottom Sheet)
│   │   └── HistoryPanel
│   │       └── HistoryItem (x10 max)
│   │
│   └── FinalView
│       └── FinalScreen
│           ├── FinalTitle
│           ├── FinalDescription
│           ├── FinalStats
│           └── Button
│
└── Toast (Global)
```

### 4.2 Component Communication

| Parent | Child | Communication |
|--------|-------|---------------|
| GameplayScreen | ChoiceButton | `@click` → emits `selectChoice` |
| GameplayScreen | StatusBar | `@openStats` → opens StatsPanel |
| GameplayScreen | StatsRow | Props: `stats` object |
| StatsPanel | Modal | `@close` → closes panel |
| FinalScreen | Button | `@click` → emits `restart` |

---

## 5. API Integration Strategy

### 5.1 GraphQL Flow

```
Component (GameplayScreen)
    │
    ├── useGame composable
    │       │
    │       └── useApi composable
    │               │
    │               └── apolloClient (mutate/query)
    │
    └── Pinia Store (game.ts)
            │
            └── localStorage (gameId persistence)
```

### 5.2 Optimistic Updates

```typescript
// In useGame composable - handle optimistic updates
async function makeChoice(choiceId: string) {
  loading.value = true;
  const expectedVersion = playerStore.version;
  
  try {
    const result = await chooseMutation({
      playerId: gameStore.gameId!,
      choiceId,
      expectedVersion,
    });
    
    // Update stores from server response
    playerStore.updateStats(result.updatedPlayer.stats);
    playerStore.setVersion(result.newVersion);
    gameStore.setTurn(result.updatedPlayer.turn);
    gameStore.setCurrentEvent(result.nextEvent);
    
    if (result.gameOver) {
      gameStore.setGameOver(result.final);
    }
  } catch (error) {
    if (error.extensions?.code === 'CONCURRENT_MODIFICATION') {
      // Reload from server and retry
      await reloadGame();
    }
    handleError(error);
  } finally {
    loading.value = false;
  }
}
```

---

## 6. State Management Approach

### 6.1 Store Responsibilities

| Store | State | Actions |
|-------|-------|---------|
| `player.ts` | stats, formalRank, informalStatus, version | setStat, updateStats, reset |
| `game.ts` | gameId, currentEvent, eventHistory, turn, isGameOver, final | setEvent, addToHistory, setGameOver |
| `ui.ts` | statsPanelOpen, historyPanelOpen, loading, error | openPanel, closePanel, setLoading |

### 6.2 State Persistence

```typescript
// Persist gameId to localStorage
watch(
  () => gameStore.gameId,
  (newId) => {
    if (newId) {
      localStorage.setItem('army_game_gameId', newId);
    } else {
      localStorage.removeItem('army_game_gameId');
    }
  },
  { immediate: true }
);

// Load saved game on init
const savedGameId = localStorage.getItem('army_game_gameId');
if (savedGameId) {
  await loadGame(savedGameId);
}
```

---

## 7. Routing Structure

### 7.1 Routes

```typescript
// src/router/index.ts
const routes = [
  {
    path: '/',
    name: 'start',
    component: () => import('@/features/start/StartView.vue'),
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
];
```

### 7.2 Route Guards

```typescript
// Redirect to start if no game
router.beforeEach((to, from, next) => {
  const gameStore = useGameStore();
  
  if (to.meta.requiresGame && !gameStore.gameId) {
    next({ name: 'start' });
  } else {
    next();
  }
});
```

---

## 8. Error Handling and Loading States

### 8.1 Error Types

| Error Type | Source | Handling |
|------------|--------|----------|
| `CHOICE_UNAVAILABLE` | GraphQL | Show toast, disable button |
| `CONCURRENT_MODIFICATION` | GraphQL | Reload state, retry |
| `PLAYER_NOT_FOUND` | GraphQL | Clear saved game, redirect to start |
| `NETWORK_ERROR` | Apollo | Show retry dialog |

### 8.2 Loading States

```typescript
// Global loading state in UI store
const uiStore = useUiStore();

// Show spinner during any GraphQL operation
// Use Apollo's useLazyQuery with loading ref
```

### 8.3 Error Boundary

```typescript
// App-level error handling
window.addEventListener('unhandledrejection', (event) => {
  console.error('Unhandled promise rejection:', event.reason);
  showErrorToast('Произошла ошибка. Попробуйте перезагрузить страницу.');
});
```

---

## 9. Testing Strategy

### 9.1 Test Structure

```
tests/
├── unit/
│   ├── setup.test.ts           # App initialization
│   ├── stores/
│   │   ├── player.test.ts      # Player store tests
│   │   ├── game.test.ts        # Game store tests
│   │   └── ui.test.ts          # UI store tests
│   ├── composables/
│   │   ├── usePlayer.test.ts   # Player composable tests
│   │   ├── useGame.test.ts     # Game composable tests
│   │   └── useApi.test.ts      # API wrapper tests
│   ├── components/
│   │   ├── ui/
│   │   │   ├── Button.test.ts
│   │   │   ├── Modal.test.ts
│   │   │   └── Spinner.test.ts
│   │   ├── stat/
│   │   │   ├── StatBar.test.ts
│   │   │   └── StatCard.test.ts
│   │   └── event/
│   │       ├── EventCard.test.ts
│   │       └── ChoiceButton.test.ts
│   ├── features/
│   │   ├── start/
│   │   │   └── StartScreen.test.ts
│   │   ├── gameplay/
│   │   │   └── GameplayScreen.test.ts
│   │   ├── stats/
│   │   │   └── StatsPanel.test.ts
│   │   ├── history/
│   │   │   └── HistoryPanel.test.ts
│   │   └── final/
│   │       └── FinalScreen.test.ts
│   └── api/
│       ├── client.test.ts     # Apollo client tests
│       ├── queries.test.ts    # GraphQL query tests
│       └── mutations.test.ts  # GraphQL mutation tests
│
├── integration/
│   └── game-flow.test.ts      # Full game flow tests
│
└── e2e/
    └── app.test.ts            # (Optional) Playwright tests
```

### 9.2 Test Coverage Targets

| Category | Target Coverage |
|----------|-----------------|
| Stores | 90%+ |
| Composables | 85%+ |
| Components | 80%+ (render, props, events) |
| Integration | Core flows |

### 9.3 Running Tests

```bash
# Run all tests
npm test

# Run with coverage
npm test -- --coverage

# Run specific test file
npm test -- tests/unit/stores/player.test.ts

# Run in watch mode
npm test -- --watch
```

---

## 10. Implementation Checklist

### Phase 1: Foundation (Steps 1-3)
- [x] Step 1: Initialize Vue 3 + Vite + TypeScript + Tailwind
- [x] Step 2: Configure Apollo Client + GraphQL code generation
- [x] Step 3: Implement Pinia stores (player, game, ui)

### Phase 2: Core Components (Steps 4-6)
- [x] Step 4: Build UI components (Button, Modal, Card, StatBar)
- [x] Step 5: Create composables (usePlayer, useGame, useApi)
- [x] Step 6: Implement GraphQL queries and mutations

### Phase 3: Feature Screens (Steps 7-10)
- [x] Step 7: Start Screen with new game / continue logic
- [x] Step 8: Gameplay Screen with event display and choices
- [x] Step 9: Stats Panel and History Panel modals
- [x] Step 10: Final Screen with all three final types

### Phase 4: Character Selection (v1.1)
- [x] Character Selection Screen with slider
- [x] 5 character types (Толстяк, Очкарик, Худой, Коротышка, Блатной)
- [x] Touch/swipe navigation between characters
- [x] Character stats display (STR, END, AGI, MOR, DISC)
- [x] Character images copied to public folder
- [x] Router updated: start → character → gameplay
- [x] Game store updated with selectedCharacter

### Final Integration
- [x] Connect all screens via router
- [ ] Implement optimistic locking with version
- [ ] Add loading states and error handling
- [ ] Test full game flow (start → character → play → finish → restart)
- [x] Verify responsive design (mobile + desktop)

---

**Документ обновлён:** 2026-03-26  
**Следующий шаг:** Подключение бэкенда