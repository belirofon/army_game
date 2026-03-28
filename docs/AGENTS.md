# AGENTS.md — Narrative Survival Game "Армейка"

## Project Overview

Army Game — narrative survival / choice-based RPG simulating military conscription experience. 
The project is a monorepo with backend (Go), frontend (Vue), and balance simulator (Python).

**Stack:**
- Backend: Go + GraphQL (gqlgen) + PostgreSQL + Redis
- Frontend: Vue 3 + TypeScript + Vite + Apollo Client + Pinia
- Balance Simulator: Python 3.10+ (stdlib only)
- Infrastructure: Docker Compose

---

## Project Structure

```
army_game/
├── backend/              # Go + GraphQL API (git submodule)
├── frontend/             # Vue SPA (git submodule)  
├── balance_simulator/    # Python Monte Carlo simulator
├── docs/                 # Technical specifications
└── docker-compose.yml    # PostgreSQL + Redis
```

---

## Build, Lint, and Test Commands

### Backend (Go)

```bash
cd backend

# Build
go build ./...

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run single test
go test ./internal/domain/services/... -run TestPlayerService -v

# Lint (requires golangci-lint)
golangci-lint run ./...

# Format code
go fmt ./...

# Generate GraphQL code
go generate ./...
```

### Frontend (Vue + TypeScript)

```bash
cd frontend

# Install dependencies
npm install

# Development server
npm run dev

# Build for production
npm run build

# Lint
npm run lint

# Type check
npm run typecheck

# Run single test file
npm test -- src/hooks/usePlayer.test.ts

# Run tests with coverage
npm run test -- --coverage

# Format code (Prettier)
npm run format
```

### Balance Simulator (Python)

```bash
cd balance_simulator

# Run simulations
python simulator.py --runs 1000
python simulator.py --runs 5000 --seed 42

# Verbose output
python simulator.py --runs 100 --verbose

# Save results
python simulator.py --runs 1000 --save results/run1.csv

# Different output formats
python simulator.py --runs 1000 --output csv
python simulator.py --runs 1000 --output json
```

### Docker Compose

```bash
# Start all services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f

# Restart a specific service
docker-compose restart postgres
```

---

## Code Style Guidelines

### General Principles

1. **TypeScript strict mode** — Always enable `strict: true` in tsconfig
2. **Single Responsibility Principle** — Each component/module handles one task
3. **Readable code** — Avoid magic numbers and obscure logic
4. **Constants and enums** — Use for fixed values (statuses, event types, NPC roles)
5. **Minimal nesting** — Use guard clauses and early returns
6. **Document complex logic** — Use JSDoc/TSDoc comments

### Vue 3 + TypeScript

**Imports order:**
1. Vue and external libraries
2. Internal imports (from `@/` alias)
3. Relative imports
4. Type imports first, then regular imports

```typescript
// Good
import { computed, ref } from 'vue';
import { useQuery } from '@apollo/client';
import type { PlayerStats } from '@/types';
import { usePlayer } from '@/composables/usePlayer';
import { StatsPanel } from './StatsPanel.vue';

// Bad
import { usePlayer } from '@/composables/usePlayer';
import { computed, ref } from 'vue';
import { useQuery } from '@apollo/client';
```

**Naming conventions:**
- Variables/functions: `camelCase`
- Components: `PascalCase` with `.vue` extension
- Composables: `camelCase.ts` (e.g., `usePlayer.ts`)
- Constants: `SCREAMING_SNAKE_CASE` or `PascalCase` for public

**Types:**
- Always type function parameters and return values
- Use interfaces for object shapes, types for unions/intersections
- Avoid `any` — use `unknown` when type is truly unknown

```typescript
// Good - Props interface
interface EventCardProps {
  title: string;
  description: string;
  disabled?: boolean;
}

// Vue 3 Component
defineComponent({
  props: {
    title: { type: String, required: true },
    description: { type: String, required: true },
    disabled: { type: Boolean, default: false },
  },
  emits: ['select'],
  setup(props, { emit }) {
    const handleClick = (id: string) => emit('select', id);
    return { handleClick };
  },
});
```

### Go

**Package structure:**
```
cmd/server/main.go
internal/
  app/handlers/       # GraphQL resolvers
  app/middleware/     # HTTP middleware
  app/services/       # Use cases
  domain/entities/    # Domain entities
  domain/services/    # Business logic
  domain/repositories/# Repository interfaces
  infrastructure/
    database/         # PostgreSQL connection
    repositories/     # Repository implementations
    config/           # Configuration
pkg/
  graphql/            # GraphQL types
  errors/             # Custom error types
```

**Naming:**
- Packages: lowercase, short (`player`, `event`, `game`)
- Files: `snake_case.go` (`player_service.go`)
- Variables: `camelCase` for local, `PascalCase` for exported
- Interfaces: `PascalCase` with `er` suffix (`PlayerRepository`)

**Error handling:**
- Custom error types in `pkg/errors/`
- Error wrapping with `fmt.Errorf` and `%w`
- Structured logging with `zap`
- Always handle errors explicitly

```go
// Good
func (s *PlayerService) GetPlayer(ctx context.Context, id string) (*Player, error) {
    player, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrPlayerNotFound
        }
        return nil, fmt.Errorf("failed to get player: %w", err)
    }
    return player, nil
}

// Bad
func (s *PlayerService) GetPlayer(id string) (*Player, error) {
    player, _ := s.repo.GetByID(context.Background(), id)
    return player, nil
}
```

### Python (Balance Simulator)

**Style:**
- Type hints for function signatures
- dataclasses for data structures
- Clear docstrings for complex functions
- Avoid star imports

```python
# Good
@dataclass
class PlayerState:
    str_stat: int = 50
    end: int = 50
    agi: int = 50
    mor: int = 50
    disc: int = 0
    turn: int = 1

    def apply_effect(self, stat: str, delta: int) -> None:
        """Apply stat change with clamping to valid ranges."""
        if stat == 'str':
            self.str_stat = max(0, min(100, self.str_stat + delta))
        # ...
```

### Vue 3 Components

- Composition API with `<script setup>`
- Small, reusable components (EventCard, ChoiceButton)
- Complex logic in composables (custom hooks)
- Props typed via interfaces with `defineProps`
- State via Pinia store, avoid global mutable state
- Keep components under 50-70 lines, extract logic to composables

```vue
<!-- Good - component under 50 lines -->
<script setup lang="ts">
interface Props {
  event: Event;
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
});

const emit = defineEmits<{
  select: [choiceId: string];
}>();

const handleChoice = (choiceId: string) => {
  emit('select', choiceId);
};
</script>

<template>
  <div class="event-card">
    <p>{{ event.description }}</p>
    <div class="choices">
      <button
        v-for="choice in event.choices"
        :key="choice.id"
        :disabled="disabled || !choice.available"
        @click="handleChoice(choice.id)"
      >
        {{ choice.text }}
      </button>
    </div>
  </div>
</template>
```

### GraphQL

- Use **gqlgen** for code-first approach
- Schema in `pkg/graphql/schema.graphql`
- Generated types in `pkg/graphql/models/`
- Resolvers in `internal/app/handlers/`
- Strict typing for all operations via graphql-codegen on frontend

---

## Game Logic Guidelines

### Stats Constraints (MUST enforce)

| Stat | Range | Description |
|------|-------|-------------|
| STR | [1, 100] | Physical strength |
| END | [1, 100] | Endurance |
| AGI | [1, 100] | Agility |
| MOR | [0, 100] | Morale (MOR=0 = game over) |
| DISC | [-100, +100] | Discipline (positive = formal, negative = informal) |

### Event System

- Events use LEGO-template pattern with variable substitution
- Always clamp stat effects to valid ranges
- Track flags for conditional events
- Anti-repetition: exclude last 5 events from selection

### Balance Rules

- Death rate target: 20-35%
- Victory rate target: 65-80%
- Average MOR loss: ~1.5 per turn
- Recovery trigger: SAFE event after 2 consecutive negative MOR events

---

## Testing Requirements

### Backend (Go)
- Unit tests for stat validation, effects, checks
- Integration tests with testcontainers
- Table-driven tests for parameterized scenarios
- Mock interfaces via testify/mock

### Frontend (Vue)
- Component tests via Vue Testing Library + Vitest
- Hook tests with renderHook
- Store tests for state management
- E2E optional (Cypress/Playwright)

### Balance Simulator (Python)
- Run with `--runs 1000 --seed 42` for reproducible results
- Check metrics against targets after any events.json change

---

## Documentation Standards

- All public functions/modules require docstrings
- Complex business logic requires inline comments explaining "why"
- Update TZ.md (docs/tech.md) when changing game mechanics
- Keep GraphQL schema as source of truth for API contracts

---

## Workflow Requirements

### Before Submitting Changes

1. Run linting: `npm run lint` (frontend) / `golangci-lint run` (backend)
2. Run type checks: `npm run typecheck` (frontend) / `go vet` (backend)
3. Run tests: ensure all pass, especially for changed modules
4. Test balance changes with simulator: `python simulator.py --runs 1000`

### Code Review Checklist

- [ ] Types are explicit, no `any`
- [ ] Errors are handled explicitly
- [ ] Magic numbers replaced with constants
- [ ] Functions are small and single-purpose
- [ ] Tests cover edge cases
- [ ] Game logic changes validated with simulator
