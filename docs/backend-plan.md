# Backend Implementation Plan — Narrative Survival Game «Армейка»

**Версия:** 1.0  
**Дата:** 2026-03-26  
**Статус:** Production-Ready Implementation Plan  
**Language:** Go 1.21+  
**Framework:** gqlgen (GraphQL) + PostgreSQL (pgx)  

---

## 1. Project Overview

### 1.1 Technology Stack

| Component | Technology | Version | Justification |
|-----------|------------|---------|----------------|
| Language | Go | 1.21+ | Performance, concurrency, type safety |
| GraphQL | gqlgen | 0.14+ | Code-first GraphQL, type-safe |
| Database | PostgreSQL | 15+ | ACID compliance, JSONB, constraints |
| DB Driver | pgx | 5.5+ | High-performance PostgreSQL driver |
| Configuration | Viper | 1.18+ | YAML/config management |
| Logging | Zap | 1.24+ | Structured JSON logging |
| Validation | go-playground/validator | 10.15+ | Input validation |
| Testing | testify | 1.9+ | Assertions, mocks, test helpers |
| Migrations | golang-migrate | 4.16+ | Database migrations |

### 1.2 Architecture Pattern

**Clean Architecture** with three distinct layers:

1. **Application Layer** (`internal/app/`) — HTTP handlers, GraphQL resolvers, middleware
2. **Domain Layer** (`internal/domain/`) — Business logic, entities, repository interfaces
3. **Infrastructure Layer** (`internal/infrastructure/`) — Database, cache, external services

### 1.3 Database Design Alignment

Following the corrected schema from `docs/DB_ANALYSIS.md`:

- Stats stored as individual columns (NOT JSONB)
- CHECK constraints for stat ranges
- Version field for optimistic locking
- Proper indexes for query patterns

---

## 2. Project Structure

### 2.1 Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
│
├── internal/
│   ├── app/
│   │   ├── handlers/               # GraphQL resolvers
│   │   │   ├── player.go           # Player queries/mutations
│   │   │   ├── event.go            # Event queries/mutations
│   │   │   └── game.go             # Game flow resolvers
│   │   ├── middleware/             # HTTP middleware
│   │   │   ├── logging.go          # Request logging
│   │   │   ├── cors.go             # CORS handling
│   │   │   └── recovery.go         # Panic recovery
│   │   └── services/               # Use cases / application services
│   │       ├── player_service.go   # Player operations
│   │       ├── game_service.go    # Game flow orchestration
│   │       └── event_service.go   # Event operations
│   │
│   ├── domain/
│   │   ├── entities/              # Domain entities
│   │   │   ├── player.go           # Player entity
│   │   │   ├── event.go            # Event template/entity
│   │   │   ├── game_log.go         # Game log entry
│   │   │   └── final.go            # Final type entity
│   │   ├── services/              # Business logic
│   │   │   ├── check_engine.go     # Check resolution (Threshold/Probability/Catastrophic)
│   │   │   ├── effect_engine.go   # Stat change application with clamping
│   │   │   ├── progression.go     # Rank/status calculation
│   │   │   ├── event_selector.go  # Weighted event selection
│   │   │   └── final_determiner.go # Final type determination
│   │   └── repositories/          # Repository interfaces
│   │       ├── player_repository.go
│   │       ├── event_repository.go
│   │       └── game_log_repository.go
│   │
│   └── infrastructure/
│       ├── database/
│       │   ├── connection.go       # PostgreSQL connection
│       │   ├── migrations/         # SQL migrations
│       │   └── seed/               # Event templates seeding
│       ├── repositories/          # Repository implementations
│       │   ├── postgres_player.go
│       │   ├── postgres_event.go
│       │   └── postgres_game_log.go
│       └── config/                 # Configuration
│           └── config.go
│
├── pkg/
│   ├── graphql/
│   │   ├── schema.graphql          # GraphQL schema definition
│   │   ├── generated.go            # Generated code (DO NOT EDIT)
│   │   ├── models.go               # GraphQL models
│   │   └── resolvers.go            # Resolver interfaces
│   ├── errors/
│   │   ├── errors.go               # Custom error types
│   │   └── codes.go                # Error codes
│   └── logger/
│       └── logger.go               # Logger setup
│
├── migrations/
│   ├── 001_create_tables.up.sql
│   └── 001_create_tables.down.sql
│
├── go.mod
├── go.sum
├── .gqlgen.yml                     # gqlgen configuration
├── docker-compose.yml              # Local development
└── Makefile                        # Build commands
```

### 2.2 File Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Packages | lowercase, short | `player`, `event`, `services` |
| Files | snake_case.go | `player_service.go`, `check_engine.go` |
| Interfaces | PascalCase + er suffix | `PlayerRepository`, `EventSelector` |
| Variables | camelCase | `playerID`, `currentTurn` |
| Constants | PascalCase or SCREAMING_SNAKE | `FormalRank`, `MAX_TURN` |

---

## 3. Step-by-Step Implementation Plan

### Phase 1: Foundation Setup (Steps 1-3)

#### Step 1: Project Initialization and Configuration

**TDD Workflow:**

1. **Write Tests (before implementation):**
   ```go
   // internal/infrastructure/config/config_test.go
   package config

   import (
       "testing"
       "github.com/stretchr/testify/assert"
   )

   func TestConfig_Load(t *testing.T) {
       // Set environment variables
       os.Setenv("DB_HOST", "localhost")
       os.Setenv("DB_PORT", "5432")
       defer os.Unsetenv("DB_HOST")
       defer os.Unsetenv("DB_PORT")

       cfg, err := Load()
       assert.NoError(t, err)
       assert.Equal(t, "localhost", cfg.Database.Host)
       assert.Equal(t, 5432, cfg.Database.Port)
   }

   func TestConfig_Validate(t *testing.T) {
       cfg := &Config{
           Server: ServerConfig{Port: 8080},
           Database: DatabaseConfig{Host: "localhost", Port: 5432, Name: "test", User: "test", Password: "test"},
       }
       err := cfg.Validate()
       assert.NoError(t, err)
   }

   func TestConfig_Validate_MissingFields(t *testing.T) {
       cfg := &Config{}
       err := cfg.Validate()
       assert.Error(t, err)
       assert.Contains(t, err.Error(), "server port")
   }
   ```

2. **Run Tests:** `go test ./internal/infrastructure/config/...` → FAIL (config package doesn't exist)

3. **Implement Minimal Logic:**
   - Create go.mod: `go mod init army-game-backend`
   - Install dependencies: `go get github.com/99designs/gqlgen`, `github.com/jackc/pgx/v5`, etc.
   - Create config struct and Load function
   - Add validation logic

4. **Refactor:** Organize into proper package structure

**Verification:** Run tests again → PASS

---

#### Step 2: Database Schema and Migrations

**TDD Tests:**

```go
// internal/infrastructure/database/connection_test.go
package database

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestConnection_Ping(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    conn, err := NewConnection(ConnectionConfig{
        Host:     "localhost",
        Port:     5432,
        Database: "army_game_test",
        User:     "postgres",
        Password: "postgres",
    })
    require.NoError(t, err)
    defer conn.Close()
    
    err = conn.Ping()
    assert.NoError(t, err)
}

func TestPlayersTable_Schema(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Test that CHECK constraints exist
    // This would be verified through migration scripts
}
```

**SQL Migration:**
```sql
-- migrations/001_create_tables.up.sql

-- ============================================================================
-- PLAYERS TABLE — Single source of truth for player state
-- ============================================================================
CREATE TABLE players (
    -- Identification
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Stats as individual columns (NOT JSONB) with CHECK constraints
    str INT NOT NULL DEFAULT 50 CHECK (str BETWEEN 1 AND 100),
    end_ INT NOT NULL DEFAULT 50 CHECK (end_ BETWEEN 1 AND 100),
    agi INT NOT NULL DEFAULT 50 CHECK (agi BETWEEN 1 AND 100),
    mor INT NOT NULL DEFAULT 50 CHECK (mor BETWEEN 0 AND 100),
    disc INT NOT NULL DEFAULT 0 CHECK (disc BETWEEN -100 AND 100),
    
    -- Progression
    formal_rank VARCHAR(50) NOT NULL DEFAULT 'РЯДОВОЙ',
    informal_status VARCHAR(50) NOT NULL DEFAULT 'ЗАПАХ',
    
    -- Game State
    turn INT NOT NULL DEFAULT 1 CHECK (turn BETWEEN 1 AND 30),
    flags JSONB DEFAULT '[]',
    current_event_template_id VARCHAR(100),
    
    -- Optimistic Locking
    version INT NOT NULL DEFAULT 1,
    
    -- Metadata
    is_finished BOOLEAN NOT NULL DEFAULT FALSE,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_players_updated_at ON players(updated_at DESC);
CREATE INDEX idx_players_turn ON players(turn);
CREATE INDEX idx_players_finished ON players(is_finished) WHERE is_finished = FALSE;
CREATE INDEX idx_players_finished_at ON players(is_finished, finished_at) WHERE is_finished = TRUE;

-- ============================================================================
-- GAME_LOGS TABLE — Immutable history of all choices
-- ============================================================================
CREATE TABLE game_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    turn INT NOT NULL,
    event_template_id VARCHAR(100) NOT NULL,
    event_description TEXT NOT NULL,
    choice_id VARCHAR(100) NOT NULL,
    choice_text TEXT NOT NULL,
    check_result JSONB NOT NULL,
    effects JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT uq_game_logs_player_turn UNIQUE (player_id, turn)
);

-- Indexes
CREATE INDEX idx_game_logs_player_id ON game_logs(player_id);
CREATE INDEX idx_game_logs_player_turn ON game_logs(player_id, turn DESC);
CREATE INDEX idx_game_logs_created_at ON game_logs(created_at);

-- ============================================================================
-- EVENT_TEMPLATES TABLE — Read-only game content
-- ============================================================================
CREATE TABLE event_templates (
    id VARCHAR(100) PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    tags JSONB DEFAULT '[]',
    context JSONB NOT NULL,
    template JSONB NOT NULL,
    choices JSONB NOT NULL,
    used_count INT DEFAULT 0,
    last_used_at TIMESTAMP WITH TIME ZONE,
    version INT DEFAULT 1,
    is_active BOOLEAN DEFAULT TRUE
);

-- Indexes
CREATE INDEX idx_event_templates_type ON event_templates(type);
CREATE INDEX idx_event_templates_type_active ON event_templates(type) WHERE is_active = TRUE;
CREATE INDEX idx_event_templates_used_count ON event_templates(used_count);
```

---

#### Step 3: Domain Entities

**TDD Tests:**

```go
// internal/domain/entities/player_test.go
package entities

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestPlayer_New(t *testing.T) {
    player := NewPlayer()
    
    assert.Equal(t, 50, player.Str)
    assert.Equal(t, 50, player.End)
    assert.Equal(t, 50, player.Agi)
    assert.Equal(t, 50, player.Mor)
    assert.Equal(t, 0, player.Disc)
    assert.Equal(t, 1, player.Turn)
    assert.Equal(t, "РЯДОВОЙ", player.FormalRank)
    assert.Equal(t, "ЗАПАХ", player.InformalStatus)
}

func TestPlayer_SetStat(t *testing.T) {
    player := NewPlayer()
    
    player.SetStat("mor", -10)
    assert.Equal(t, 0, player.Mor) // Clamped to min
    
    player.SetStat("str", 150)
    assert.Equal(t, 100, player.Str) // Clamped to max
    
    player.SetStat("disc", 200)
    assert.Equal(t, 100, player.Disc)
    
    player.SetStat("disc", -200)
    assert.Equal(t, -100, player.Disc)
}

func TestPlayer_ApplyEffect(t *testing.T) {
    player := NewPlayer()
    
    effect := Effect{
        Stat:          "mor",
        Delta:         -20,
        PreviousValue: 50,
        NewValue:      30,
    }
    
    player.ApplyEffect(effect)
    assert.Equal(t, 30, player.Mor)
}

func TestPlayer_IsGameOver(t *testing.T) {
    player := NewPlayer()
    
    player.Mor = 0
    assert.True(t, player.IsGameOver())
    
    player.Mor = 50
    player.Turn = 30
    assert.True(t, player.IsGameOver())
    
    player.Turn = 29
    player.Mor = 50
    assert.False(t, player.IsGameOver())
}
```

**Implementation:**
```go
// internal/domain/entities/player.go
package entities

type Player struct {
    ID        string `json:"id"`
    Str       int    `json:"str"`
    End       int    `json:"end"`
    Agi       int    `json:"agi"`
    Mor       int    `json:"mor"`
    Disc      int    `json:"disc"`
    FormalRank   string   `json:"formalRank"`
    InformalStatus string `json:"informalStatus"`
    Turn      int       `json:"turn"`
    Flags     []string  `json:"flags"`
    CurrentEventTemplateID string `json:"currentEventTemplateId,omitempty"`
    Version   int       `json:"version"`
    IsFinished bool      `json:"isFinished"`
    FinishedAt *string   `json:"finishedAt,omitempty"`
    CreatedAt string    `json:"createdAt"`
    UpdatedAt string    `json:"updatedAt"`
}

const (
    MinStr = 1
    MaxStr = 100
    MinEnd = 1
    MaxEnd = 100
    MinAgi = 1
    MaxAgi = 100
    MinMor = 0
    MaxMor = 100
    MinDisc = -100
    MaxDisc = 100
    MaxTurn = 30
)

func NewPlayer() *Player {
    return &Player{
        Str:   50,
        End:   50,
        Agi:   50,
        Mor:   50,
        Disc:  0,
        FormalRank:    "РЯДОВОЙ",
        InformalStatus: "ЗАПАХ",
        Turn:    1,
        Flags:   []string{},
        Version: 1,
    }
}

func (p *Player) SetStat(stat string, value int) {
    switch stat {
    case "str":
        p.Str = clamp(value, MinStr, MaxStr)
    case "end":
        p.End = clamp(value, MinEnd, MaxEnd)
    case "agi":
        p.Agi = clamp(value, MinAgi, MaxAgi)
    case "mor":
        p.Mor = clamp(value, MinMor, MaxMor)
    case "disc":
        p.Disc = clamp(value, MinDisc, MaxDisc)
    }
}

func (p *Player) ApplyEffect(effect Effect) {
    p.SetStat(effect.Stat, effect.NewValue)
}

func (p *Player) IsGameOver() bool {
    return p.Mor <= 0 || p.Turn >= MaxTurn
}

func clamp(value, min, max int) int {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}
```

---

### Phase 2: Core Services (Steps 4-7)

#### Step 4: Check Engine Implementation

**TDD Tests:**

```go
// internal/domain/services/check_engine_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "army-game-backend/internal/domain/entities"
)

func TestCheckEngine_ThresholdCheck_Success(t *testing.T) {
    engine := NewCheckEngine()
    
    player := &entities.Player{Str: 70, End: 50, Agi: 60, Mor: 50, Disc: 0}
    check := &Check{
        Type:      CheckTypeThreshold,
        Stat:      "str",
        Threshold: 60,
        Operator:  ">=",
    }
    
    result := engine.ExecuteCheck(check, player)
    
    assert.True(t, result.Success)
    assert.Equal(t, entities.OutcomeSuccess, result.Outcome)
    assert.Contains(t, result.Description, "успех")
}

func TestCheckEngine_ThresholdCheck_Failure(t *testing.T) {
    engine := NewCheckEngine()
    
    player := &entities.Player{Str: 40, End: 50, Agi: 60, Mor: 50, Disc: 0}
    check := &Check{
        Type:      CheckTypeThreshold,
        Stat:      "str",
        Threshold: 60,
        Operator:  ">=",
    }
    
    result := engine.ExecuteCheck(check, player)
    
    assert.False(t, result.Success)
    assert.Equal(t, entities.OutcomeFailure, result.Outcome)
    assert.Contains(t, result.Description, "провал")
}

func TestCheckEngine_ProbabilityCheck(t *testing.T) {
    engine := NewCheckEngine()
    engine.rand = &mockRand{fixedValue: 50} // 50% chance
    
    player := &entities.Player{Str: 50, End: 50, Agi: 50, Mor: 50, Disc: 0}
    check := &Check{
        Type:       CheckTypeProbability,
        Stat:       "str",
        BaseChance: 60, // 60% base
    }
    
    result := engine.ExecuteCheck(check, player)
    
    // With difficulty modifier 1.0 (mid game), 60 / 1.0 = 60%, rand 50 < 60 = success
    assert.True(t, result.Success)
}

func TestCheckEngine_CatastrophicCheck_Noticed(t *testing.T) {
    engine := NewCheckEngine()
    engine.rand = &mockRand{fixedValue: 30} // 30% notice chance, will notice
    
    player := &entities.Player{Str: 30, End: 50, Agi: 50, Mor: 50, Disc: 0}
    check := &Check{
        Type:          CheckTypeCatastrophic,
        NoticeChance:  50,
        PowerStat:      "str",
        PowerThreshold: 40,
    }
    
    result := engine.ExecuteCheck(check, player)
    
    assert.Equal(t, entities.OutcomeNoticedFailure, result.Outcome)
    assert.Contains(t, result.Description, "заметили")
}

func TestCheckEngine_CatastrophicCheck_NotNoticed(t *testing.T) {
    engine := NewCheckEngine()
    engine.rand = &mockRand{fixedValue: 80} // 80% > 50% notice chance = not noticed
    
    player := &entities.Player{Str: 30, End: 50, Agi: 50, Mor: 50, Disc: 0}
    check := &Check{
        Type:          CheckTypeCatastrophic,
        NoticeChance:  50,
        PowerStat:      "str",
        PowerThreshold: 40,
    }
    
    result := engine.ExecuteCheck(check, player)
    
    assert.Equal(t, entities.OutcomeIgnored, result.Outcome)
}
```

**Implementation:**
```go
// internal/domain/services/check_engine.go
package services

import (
    "math/rand"
    "army-game-backend/internal/domain/entities"
)

type CheckType string

const (
    CheckTypeThreshold    CheckType = "threshold"
    CheckTypeProbability  CheckType = "probability"
    CheckTypeCatastrophic CheckType = "catastrophic"
)

type Check struct {
    Type          CheckType `json:"type"`
    Stat          string    `json:"stat"`
    Threshold     int       `json:"threshold,omitempty"`
    Operator      string    `json:"operator,omitempty"`
    BaseChance    int       `json:"baseChance,omitempty"`
    NoticeChance  int       `json:"noticeChance,omitempty"`
    PowerStat     string    `json:"powerStat,omitempty"`
    PowerThreshold int      `json:"powerThreshold,omitempty"`
    Difficulty    float64   `json:"difficulty,omitempty"`
}

type CheckResult struct {
    Success    bool                 `json:"success"`
    Outcome    entities.OutcomeType `json:"outcome"`
    Description string             `json:"description"`
}

type CheckEngine struct {
    rand rand.Rand
}

func NewCheckEngine() *CheckEngine {
    return &CheckEngine{
        rand: *rand.New(rand.NewSource(rand.Int63())),
    }
}

func (e *CheckEngine) ExecuteCheck(check *Check, player *entities.Player) *CheckResult {
    switch check.Type {
    case CheckTypeThreshold:
        return e.executeThresholdCheck(check, player)
    case CheckTypeProbability:
        return e.executeProbabilityCheck(check, player)
    case CheckTypeCatastrophic:
        return e.executeCatastrophicCheck(check, player)
    default:
        return &CheckResult{
            Success:    false,
            Outcome:    entities.OutcomeFailure,
            Description: "Неизвестный тип проверки",
        }
    }
}

func (e *CheckEngine) executeThresholdCheck(check *Check, player *entities.Player) *CheckResult {
    statValue := e.getStatValue(check.Stat, player)
    success := evaluateCondition(statValue, check.Operator, check.Threshold)
    
    return &CheckResult{
        Success:    success,
        Outcome:    e.determineOutcome(success),
        Description: e.generateDescription(check.Stat, success),
    }
}

func (e *CheckEngine) executeProbabilityCheck(check *Check, player *entities.Player) *CheckResult {
    baseChance := check.BaseChance
    if baseChance == 0 {
        baseChance = 50
    }
    
    // Apply difficulty modifier
    difficulty := check.Difficulty
    if difficulty == 0 {
        difficulty = 1.0
    }
    effectiveChance := int(float64(baseChance) / difficulty)
    
    roll := e.rand.Intn(100)
    success := roll < effectiveChance
    
    // Determine partial success (within 10% of threshold)
    partialSuccess := roll >= effectiveChance && roll < effectiveChance+10
    
    outcome := entities.OutcomeFailure
    if success {
        outcome = entities.OutcomeSuccess
    } else if partialSuccess {
        outcome = entities.OutcomePartial
    }
    
    return &CheckResult{
        Success:    success || partialSuccess,
        Outcome:    outcome,
        Description: e.generateDescription(check.Stat, success),
    }
}

func (e *CheckEngine) executeCatastrophicCheck(check *Check, player *entities.Player) *CheckResult {
    noticeRoll := e.rand.Intn(100)
    
    if noticeRoll < check.NoticeChance {
        // Caught! Now check power
        powerStat := e.getStatValue(check.PowerStat, player)
        success := powerStat >= check.PowerThreshold
        
        return &CheckResult{
            Success:    success,
            Outcome:    e.determineNoticedOutcome(success),
            Description: e.generateNoticedDescription(success),
        }
    }
    
    // Not noticed - ignore the action
    return &CheckResult{
        Success:    true,
        Outcome:    entities.OutcomeIgnored,
        Description: "Вас не заметили",
    }
}

func (e *CheckEngine) getStatValue(stat string, player *entities.Player) int {
    switch stat {
    case "str": return player.Str
    case "end": return player.End
    case "agi": return player.Agi
    case "mor": return player.Mor
    case "disc": return player.Disc
    default: return 0
    }
}

func evaluateCondition(value int, operator string, threshold int) bool {
    switch operator {
    case ">=": return value >= threshold
    case ">":  return value > threshold
    case "<=": return value <= threshold
    case "<":  return value < threshold
    case "==": return value == threshold
    default:   return value >= threshold
    }
}

func (e *CheckEngine) determineOutcome(success bool) entities.OutcomeType {
    if success {
        return entities.OutcomeSuccess
    }
    return entities.OutcomeFailure
}

func (e *CheckEngine) determineNoticedOutcome(success bool) entities.OutcomeType {
    if success {
        return entities.OutcomeNoticedSuccess
    }
    return entities.OutcomeNoticedFailure
}

func (e *CheckEngine) generateDescription(stat string, success bool) string {
    if success {
        return "Успех! Вы справились."
    }
    return "Провал. Не получилось."
}

func (e *CheckEngine) generateNoticedDescription(success bool) string {
    if success {
        return "Вас заметили, но вы сумели справиться!"
    }
    return "Вас заметили. Это плохо закончилось."
}
```

---

#### Step 5: Effect Engine Implementation

**TDD Tests:**

```go
// internal/domain/services/effect_engine_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "army-game-backend/internal/domain/entities"
)

func TestEffectEngine_ApplyEffect_ClampToMax(t *testing.T) {
    engine := NewEffectEngine()
    player := &entities.Player{Mor: 50}
    
    effect := Effect{
        Stat:          "mor",
        Delta:         100,
        PreviousValue: 50,
        NewValue:      150,
    }
    
    engine.ApplyEffect(player, effect)
    assert.Equal(t, 100, player.Mor) // Clamped to max
}

func TestEffectEngine_ApplyEffect_ClampToMin(t *testing.T) {
    engine := NewEffectEngine()
    player := &entities.Player{Mor: 50}
    
    effect := Effect{
        Stat:          "mor",
        Delta:         -100,
        PreviousValue: 50,
        NewValue:      -50,
    }
    
    engine.ApplyEffect(player, effect)
    assert.Equal(t, 0, player.Mor) // Clamped to min
}

func TestEffectEngine_ApplyEffects_Multiple(t *testing.T) {
    engine := NewEffectEngine()
    player := &entities.Player{Str: 50, End: 50, Mor: 50}
    
    effects := []Effect{
        {Stat: "str", Delta: 10, PreviousValue: 50, NewValue: 60},
        {Stat: "mor", Delta: -20, PreviousValue: 50, NewValue: 30},
    }
    
    engine.ApplyEffects(player, effects)
    
    assert.Equal(t, 60, player.Str)
    assert.Equal(t, 30, player.Mor)
}

func TestEffectEngine_CalculateNewValue(t *testing.T) {
    engine := NewEffectEngine()
    
    newValue := engine.CalculateNewValue(50, 10)
    assert.Equal(t, 60, newValue)
    
    newValue = engine.CalculateNewValue(50, -60)
    assert.Equal(t, -10, newValue)
}
```

**Implementation:**
```go
// internal/domain/services/effect_engine.go
package services

import "army-game-backend/internal/domain/entities"

type Effect struct {
    Stat          string `json:"stat"`
    Delta         int    `json:"delta"`
    PreviousValue int    `json:"previousValue"`
    NewValue      int    `json:"newValue"`
}

type EffectEngine struct{}

func NewEffectEngine() *EffectEngine {
    return &EffectEngine{}
}

func (e *EffectEngine) ApplyEffect(player *entities.Player, effect Effect) {
    player.ApplyEffect(effect)
}

func (e *EffectEngine) ApplyEffects(player *entities.Player, effects []Effect) {
    for _, effect := range effects {
        e.ApplyEffect(player, effect)
    }
}

func (e *EffectEngine) CalculateNewValue(currentValue, delta int) int {
    return currentValue + delta
}

func (e *EffectEngine) GetStatMin(stat string) int {
    mins := map[string]int{
        "str":  entities.MinStr,
        "end":  entities.MinEnd,
        "agi":  entities.MinAgi,
        "mor":  entities.MinMor,
        "disc": entities.MinDisc,
    }
    return mins[stat]
}

func (e *EffectEngine) GetStatMax(stat string) int {
    maxs := map[string]int{
        "str":  entities.MaxStr,
        "end":  entities.MaxEnd,
        "agi":  entities.MaxAgi,
        "mor":  entities.MaxMor,
        "disc": entities.MaxDisc,
    }
    return maxs[stat]
}
```

---

#### Step 6: Progression Service

**TDD Tests:**

```go
// internal/domain/services/progression_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "army-game-backend/internal/domain/entities"
)

func TestProgressionService_CalculateFormalRank(t *testing.T) {
    svc := NewProgressionService()
    
    player := &entities.Player{Disc: 0}
    rank := svc.CalculateFormalRank(player)
    assert.Equal(t, "РЯДОВОЙ", rank)
    
    player.Disc = 30
    rank = svc.CalculateFormalRank(player)
    assert.Equal(t, "ЕФРЕЙТОР", rank)
    
    player.Disc = 60
    rank = svc.CalculateFormalRank(player)
    assert.Equal(t, "МЛ_СЕРЖАНТ", rank)
    
    player.Disc = 80
    rank = svc.CalculateFormalRank(player)
    assert.Equal(t, "СЕРЖАНТ", rank)
}

func TestProgressionService_CalculateInformalStatus(t *testing.T) {
    svc := NewProgressionService()
    
    player := &entities.Player{Disc: 0}
    status := svc.CalculateInformalStatus(player)
    assert.Equal(t, "ЗАПАХ", status)
    
    player.Disc = -30
    status = svc.CalculateInformalStatus(player)
    assert.Equal(t, "ДУХ", status)
    
    player.Disc = -60
    status = svc.CalculateInformalStatus(player)
    assert.Equal(t, "СЛОН", status)
    
    player.Disc = -80
    status = svc.CalculateInformalStatus(player)
    assert.Equal(t, "ЧЕРПАК", status)
    
    player.Disc = -95
    status = svc.CalculateInformalStatus(player)
    assert.Equal(t, "ДЕД", status)
}
```

**Implementation:**
```go
// internal/domain/services/progression.go
package services

import "army-game-backend/internal/domain/entities"

type ProgressionService struct{}

func NewProgressionService() *ProgressionService {
    return &ProgressionService{}
}

func (s *ProgressionService) CalculateFormalRank(player *entities.Player) string {
    disc := player.Disc
    if disc >= 75 {
        return "СЕРЖАНТ"
    }
    if disc >= 50 {
        return "МЛ_СЕРЖАНТ"
    }
    if disc >= 25 {
        return "ЕФРЕЙТОР"
    }
    return "РЯДОВОЙ"
}

func (s *ProgressionService) CalculateInformalStatus(player *entities.Player) string {
    disc := player.Disc
    if disc <= -90 {
        return "ДЕМБЕЛЬ"
    }
    if disc <= -75 {
        return "ДЕД"
    }
    if disc <= -50 {
        return "ЧЕРПАК"
    }
    if disc <= -25 {
        return "СЛОН"
    }
    if disc < 0 {
        return "ДУХ"
    }
    return "ЗАПАХ"
}

func (s *ProgressionService) UpdateProgression(player *entities.Player) {
    player.FormalRank = s.CalculateFormalRank(player)
    player.InformalStatus = s.CalculateInformalStatus(player)
}
```

---

#### Step 7: Event Selector and Generator

**TDD Tests:**

```go
// internal/domain/services/event_selector_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "army-game-backend/internal/domain/entities"
)

func TestEventSelector_Select_ExcludesRecentEvents(t *testing.T) {
    selector := NewEventSelector()
    
    templates := []*entities.EventTemplate{
        {ID: "event1", Type: "ROUTINE"},
        {ID: "event2", Type: "ROUTINE"},
        {ID: "event3", Type: "ROUTINE"},
        {ID: "event4", Type: "SAFE"},
    }
    
    // Select first event
    selected := selector.Select(templates, []string{}, &entities.Player{Turn: 5})
    assert.NotNil(t, selected)
    
    // Second selection should exclude event1 (just used)
    selected = selector.Select(templates, []string{"event1"}, &entities.Player{Turn: 5})
    require.NotNil(t, selected)
    assert.NotEqual(t, "event1", selected.ID)
}

func TestEventSelector_Select_AppliesWeights(t *testing.T) {
    selector := NewEventSelector()
    
    templates := []*entities.EventTemplate{
        {ID: "routine1", Type: "ROUTINE"},
        {ID: "safe1", Type: "SAFE"},
    }
    
    // Player with low morale should get priority for SAFE events
    player := &entities.Player{Turn: 5, Mor: 20}
    selected := selector.Select(templates, []string{}, player)
    
    // In low morale, SAFE should have higher probability
    assert.NotNil(t, selected)
}

func TestEventSelector_GetDifficultyModifier(t *testing.T) {
    selector := NewEventSelector()
    
    modifier := selector.GetDifficultyModifier(5)  // Early game
    assert.Equal(t, 0.6, modifier)
    
    modifier = selector.GetDifficultyModifier(15) // Mid game  
    assert.Equal(t, 0.85, modifier)
    
    modifier = selector.GetDifficultyModifier(25) // Late game
    assert.Equal(t, 1.05, modifier)
}
```

**Implementation:**
```go
// internal/domain/services/event_selector.go
package services

import (
    "math/rand"
    "army-game-backend/internal/domain/entities"
)

type EventSelector struct {
    rand *rand.Rand
}

func NewEventSelector() *EventSelector {
    return &EventSelector{
        rand: *rand.New(rand.Int63()),
    }
}

func (s *EventSelector) Select(templates []*entities.EventTemplate, recentHistory []string, player *entities.Player) *entities.EventTemplate {
    if len(templates) == 0 {
        return nil
    }
    
    // Filter out recent events
    available := make([]*entities.EventTemplate, 0)
    for _, t := range templates {
        if !s.isRecent(t.ID, recentHistory) {
            available = append(available, t)
        }
    }
    
    if len(available) == 0 {
        // If all excluded, reset history
        available = templates
    }
    
    // Calculate weights
    weights := make([]float64, len(available))
    totalWeight := 0.0
    
    for i, t := range available {
        weight := s.calculateWeight(t, player)
        weights[i] = weight
        totalWeight += weight
    }
    
    // Select by weighted random
    roll := s.rand.Float64() * totalWeight
    cumulative := 0.0
    
    for i, w := range weights {
        cumulative += w
        if cumulative >= roll {
            return available[i]
        }
    }
    
    return available[0]
}

func (s *EventSelector) calculateWeight(template *entities.EventTemplate, player *entities.Player) float64 {
    weight := 1.0
    
    // Disc-based weighting
    if player.Disc > 50 && template.Type == "INSPECTION" {
        weight += 0.3
    }
    if player.Disc < -50 && template.Type == "INFORMAL" {
        weight += 0.3
    }
    
    // Morale-based weighting (SAFE events priority when low)
    if player.Mor < 30 && template.Type == "SAFE" {
        weight += 0.3
    }
    
    // Turn-based weighting
    if player.Turn <= 10 && template.Type == "ROUTINE" {
        weight += 0.1
    }
    if player.Turn >= 20 && template.Type == "EMERGENCY" {
        weight += 0.2
    }
    
    return weight
}

func (s *EventSelector) isRecent(templateID string, recentHistory []string) bool {
    for _, id := range recentHistory {
        if id == templateID {
            return true
        }
    }
    return false
}

func (s *EventSelector) GetDifficultyModifier(turn int) float64 {
    switch {
    case turn <= 10:
        return 0.6  // Easy early game
    case turn <= 20:
        return 0.85 // Normal mid game
    default:
        return 1.05 // Hard late game
    }
}
```

---

### Phase 3: Application Layer (Steps 8-10)

#### Step 8: Final Determiner Service

**TDD Tests:**

```go
// internal/domain/services/final_determiner_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "army-game-backend/internal/domain/entities"
)

func TestFinalDeterminer_Determine_СломанныйДембель(t *testing.T) {
    determiner := NewFinalDeterminer()
    
    player := &entities.Player{Mor: 0, Turn: 15}
    final := determiner.Determine(player)
    
    assert.NotNil(t, final)
    assert.Equal(t, "СЛОМАННЫЙ_ДЕМБЕЛЬ", final.Type)
    assert.Equal(t, "Сломанный дембель", final.Title)
    assert.Equal(t, "Косячник", final.Subtitle)
}

func TestFinalDeterminer_Determine_ТихийДембель(t *testing.T) {
    determiner := NewFinalDeterminer()
    
    player := &entities.Player{Mor: 30, Turn: 28, Disc: 10, InformalStatus: "ЗАПАХ"}
    final := determiner.Determine(player)
    
    assert.NotNil(t, final)
    assert.Equal(t, "ТИХИЙ_ДЕМБЕЛЬ", final.Type)
    assert.Equal(t, "Тихий дембель", final.Title)
    assert.Equal(t, "Приспособленец", final.Subtitle)
}

func TestFinalDeterminer_Determine_УважаемыйДембель(t *testing.T) {
    determiner := NewFinalDeterminer()
    
    player := &entities.Player{Mor: 55, Turn: 28, Disc: -60, InformalStatus: "ДЕД"}
    final := determiner.Determine(player)
    
    assert.NotNil(t, final)
    assert.Equal(t, "УВАЖАЕМЫЙ_ДЕМБЕЛЬ", final.Type)
    assert.Equal(t, "Уважаемый дембель", final.Title)
    assert.Nil(t, final.Subtitle)
}

func TestFinalDeterminer_Determine_NoFinal(t *testing.T) {
    determiner := NewFinalDeterminer()
    
    player := &entities.Player{Mor: 50, Turn: 15, Disc: 0}
    final := determiner.Determine(player)
    
    assert.Nil(t, final)
}
```

---

#### Step 9: Game Service (Orchestration)

**TDD Tests:**

```go
// internal/app/services/game_service_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "army-game-backend/internal/domain/entities"
)

type MockPlayerRepository struct {
    mock.Mock
}

func (m *MockPlayerRepository) Create(ctx context.Context, player *entities.Player) error {
    args := m.Called(ctx, player)
    return args.Error(0)
}

func (m *MockPlayerRepository) GetByID(ctx context.Context, id string) (*entities.Player, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entities.Player), args.Error(1)
}

func (m *MockPlayerRepository) Update(ctx context.Context, player *entities.Player) error {
    args := m.Called(ctx, player)
    return args.Error(0)
}

func (m *MockPlayerRepository) UpdateWithVersion(ctx context.Context, player *entities.Player, expectedVersion int) error {
    args := m.Called(ctx, player, expectedVersion)
    return args.Error(0)
}

func TestGameService_StartGame(t *testing.T) {
    mockRepo := new(MockPlayerRepository)
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Player")).Return(nil)
    
    svc := NewGameService(mockRepo, nil, nil, nil, nil)
    gameState, err := svc.StartGame()
    
    assert.NoError(t, err)
    assert.NotNil(t, gameState)
    assert.NotEmpty(t, gameState.GameID)
    assert.NotNil(t, gameState.Player)
    assert.Equal(t, 1, gameState.Player.Turn)
    assert.NotNil(t, gameState.CurrentEvent)
}
```

---

#### Step 10: GraphQL Handlers

**TDD Tests:**

```go
// internal/app/handlers/player_test.go
package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "army-game-backend/internal/domain/entities"
)

func TestResolver_Player(t *testing.T) {
    // Test that player query returns correct structure
    player := &entities.Player{
        ID:    "test-id",
        Str:   50,
        End:   50,
        Agi:   50,
        Mor:   50,
        Disc:  0,
        FormalRank:    "РЯДОВОЙ",
        InformalStatus: "ЗАПАХ",
        Turn:    5,
        Flags:   []string{"test_flag"},
        Version: 3,
    }
    
    // Verify field mapping to GraphQL type
    assert.Equal(t, "test-id", player.ID)
    assert.Equal(t, 50, player.Str)
    assert.Equal(t, "РЯДОВОЙ", player.FormalRank)
}

func TestResolver_StartGame(t *testing.T) {
    // Test startGame mutation returns expected structure
    // Full integration test would use test database
}
```

**GraphQL Schema:**
```graphql
# pkg/graphql/schema.graphql

scalar Time

type Query {
  player(id: ID!): Player
  currentEvent(playerId: ID!): EventInstance
  eventHistory(playerId: ID!, limit: Int = 10): [GameLogEntry!]!
  loadGame(gameId: ID!): GameState
}

type Mutation {
  startGame: GameState!
  choose(playerId: ID!, choiceId: ID!, expectedVersion: Int!): ChooseResult!
  restartGame(playerId: ID!): GameState!
}

type Player {
  id: ID!
  stats: PlayerStats!
  formalRank: FormalRank!
  informalStatus: InformalStatus!
  turn: Int!
  flags: [String!]!
  isFinished: Boolean!
  version: Int!
  createdAt: Time!
  updatedAt: Time!
}

type PlayerStats {
  str: Int!
  end: Int!
  agi: Int!
  mor: Int!
  disc: Int!
}

enum FormalRank {
  РЯДОВОЙ
  ЕФРЕЙТОР
  МЛ_СЕРЖАНТ
  СЕРЖАНТ
}

enum InformalStatus {
  ЗАПАХ
  ДУХ
  СЛОН
  ЧЕРПАК
  ДЕД
  ДЕМБЕЛЬ
}

type EventInstance {
  id: ID!
  templateId: ID!
  description: String!
  resolvedVariables: JSON
  choices: [Choice!]!
  context: EventContext!
}

type Choice {
  id: ID!
  text: String!
  available: Boolean!
}

type EventContext {
  time: String!
  location: String!
  urgency: String!
}

type GameLogEntry {
  id: ID!
  playerId: ID!
  turn: Int!
  eventDescription: String!
  choiceText: String!
  checkResult: CheckResult!
  effects: [Effect!]!
  createdAt: Time!
}

type CheckResult {
  success: Boolean!
  outcome: OutcomeType!
  description: String!
}

enum OutcomeType {
  SUCCESS
  PARTIAL
  FAILURE
  IGNORED
  NOTICED_SUCCESS
  NOTICED_FAILURE
}

type Effect {
  stat: String!
  delta: Int!
  previousValue: Int!
  newValue: Int!
}

type GameState {
  player: Player!
  currentEvent: EventInstance
  eventHistory: [GameLogEntry!]!
  isGameOver: Boolean!
  final: Final
  gameId: ID!
}

type Final {
  type: FinalType!
  title: String!
  subtitle: String
  description: String!
  finalStats: PlayerStats!
  achievedOnTurn: Int!
}

enum FinalType {
  ТИХИЙ_ДЕМБЕЛЬ
  УВАЖАЕМЫЙ_ДЕМБЕЛЬ
  СЛОМАННЫЙ_ДЕМБЕЛЬ
}

type ChooseResult {
  success: Boolean!
  checkResult: CheckResult!
  effects: [Effect!]!
  updatedPlayer: Player!
  nextEvent: EventInstance
  gameOver: Boolean!
  final: Final
  newVersion: Int!
}
```

---

## 4. Service Module Breakdown

### 4.1 Service Responsibilities

| Service | Responsibility | Public API |
|---------|----------------|------------|
| `PlayerService` | CRUD operations, state management | `Create()`, `GetByID()`, `UpdateWithVersion()` |
| `GameService` | Game flow orchestration | `StartGame()`, `MakeChoice()`, `RestartGame()` |
| `EventService` | Event template management | `GetTemplates()`, `GenerateEvent()`, `SelectEvent()` |
| `CheckEngine` | Check resolution | `ExecuteCheck()` |
| `EffectEngine` | Stat changes | `ApplyEffect()`, `ApplyEffects()` |
| `ProgressionService` | Rank/status calculation | `CalculateFormalRank()`, `CalculateInformalStatus()`, `UpdateProgression()` |
| `EventSelector` | Weighted event selection | `Select()`, `GetDifficultyModifier()` |
| `FinalDeterminer` | Final type determination | `Determine()` |

### 4.2 Dependency Graph

```
GameService
    ├── PlayerService (interface)
    ├── EventService (interface)
    ├── CheckEngine
    ├── EffectEngine
    ├── ProgressionService
    ├── EventSelector
    └── FinalDeterminer
```

---

## 5. Error Handling Strategy

### 5.1 Custom Error Types

```go
// pkg/errors/errors.go
package errors

import "fmt"

type AppError struct {
    Code    string
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// Error constructors
func NewPlayerNotFound(id string) *AppError {
    return &AppError{Code: "PLAYER_NOT_FOUND", Message: fmt.Sprintf("Player not found: %s", id)}
}

func NewChoiceUnavailable(choiceID, playerID string) *AppError {
    return &AppError{Code: "CHOICE_UNAVAILABLE", Message: fmt.Sprintf("Choice %s unavailable for player %s", choiceID, playerID)}
}

func NewVersionMismatch(expected, actual int) *AppError {
    return &AppError{Code: "VERSION_MISMATCH", Message: fmt.Sprintf("Expected version %d, got %d", expected, actual)}
}

func NewConcurrentModification(playerID string) *AppError {
    return &AppError{Code: "CONCURRENT_MODIFICATION", Message: fmt.Sprintf("Player %s was modified concurrently", playerID)}
}
```

### 5.2 Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `PLAYER_NOT_FOUND` | 404 | Player ID not found |
| `CHOICE_UNAVAILABLE` | 400 | Choice not available for player |
| `CHOICE_NOT_IN_CURRENT_EVENT` | 400 | Choice doesn't match current event |
| `GAME_ALREADY_FINISHED` | 400 | Game is over, no more moves |
| `INVALID_STAT_VALUE` | 400 | Stat value out of range |
| `CONCURRENT_MODIFICATION` | 409 | Version mismatch, retry needed |
| `VERSION_MISMATCH` | 409 | Expected version doesn't match |

---

## 6. Logging and Observability

### 6.1 Structured Logging

```go
// pkg/logger/logger.go
package logger

import (
    "go.uber.org/zap"
)

var log *zap.SugaredLogger

func Init(debug bool) {
    var cfg zap.Config
    if debug {
        cfg = zap.NewDevelopmentConfig()
    } else {
        cfg = zap.NewProductionConfig()
    }
    
    l, _ := cfg.Build()
    log = l.Sugar()
}

func Info(msg string, fields ...zap.Field) {
    log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
    log.Error(msg, fields...)
}

func With(field string, value interface{}) zap.Field {
    return zap.Any(field, value)
}
```

### 6.2 Request Logging Middleware

```go
// internal/app/middleware/logging.go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        logger.Info("Request started",
            logger.With("method", r.Method),
            logger.With("path", r.URL.Path),
            logger.With("remote_addr", r.RemoteAddr),
        )
        
        next.ServeHTTP(w, r)
        
        logger.Info("Request completed",
            logger.With("method", r.Method),
            logger.With("path", r.URL.Path),
            logger.With("duration", time.Since(start)),
        )
    })
}
```

---

## 7. Concurrency Handling

### 7.1 Optimistic Locking

```go
// Internal domain/repository pattern for optimistic locking
func (r *PostgresPlayerRepository) UpdateWithVersion(ctx context.Context, player *entities.Player, expectedVersion int) error {
    query := `
        UPDATE players 
        SET str = $1, end_ = $2, agi = $3, mor = $4, disc = $5,
            formal_rank = $6, informal_status = $7, turn = $8,
            flags = $9, version = version + 1, updated_at = NOW(),
            is_finished = $10, finished_at = $11
        WHERE id = $12 AND version = $13
    `
    
    result, err := r.pool.Exec(ctx, query,
        player.Str, player.End, player.Agi, player.Mor, player.Disc,
        player.FormalRank, player.InformalStatus, player.Turn,
        player.Flags, player.IsFinished, player.FinishedAt,
        player.ID, expectedVersion,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update player: %w", err)
    }
    
    if result.RowsAffected() == 0 {
        return errors.ConcurrentModification(player.ID)
    }
    
    player.Version = expectedVersion + 1
    return nil
}
```

### 7.2 Context for Request Cancellation

```go
func (s *GameService) StartGame() (*GameState, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // All DB operations use context
    player, err := s.playerRepo.Create(ctx, entities.NewPlayer())
    // ...
}
```

---

## 8. Testing Strategy

### 8.1 Test Structure

```
backend/
├── internal/
│   ├── domain/
│   │   ├── entities/
│   │   │   └── player_test.go
│   │   │   └── event_test.go
│   │   ├── services/
│   │   │   ├── check_engine_test.go
│   │   │   ├── effect_engine_test.go
│   │   │   ├── progression_test.go
│   │   │   ├── event_selector_test.go
│   │   │   └── final_determiner_test.go
│   │   └── repositories/
│   │       └── player_repository_test.go
│   ├── app/
│   │   ├── handlers/
│   │   │   └── player_test.go
│   │   └── services/
│   │       └── game_service_test.go
│   └── infrastructure/
│       ├── config/
│       │   └── config_test.go
│       └── database/
│           └── connection_test.go
└── pkg/
    ├── errors/
    │   └── errors_test.go
    └── logger/
        └── logger_test.go
```

### 8.2 Test Types

| Type | Tool | Target | Coverage Goal |
|------|------|--------|----------------|
| Unit | `testing` + `testify` | Domain services, entities | 90%+ |
| Integration | `testcontainers` | Repository implementations | 80%+ |
| Handler | `httptest` | GraphQL resolvers | 70%+ |

### 8.3 Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/domain/services/... -v

# Run short tests (skip integration)
go test -short ./...

# Run with race detector
go test -race ./...
```

### 8.4 Example Unit Test Patterns

```go
// Table-driven test pattern
func TestCheckEngine_ThresholdCheck(t *testing.T) {
    tests := []struct {
        name       string
        statValue  int
        threshold  int
        operator   string
        wantSuccess bool
    }{
        {"success_equal", 50, 50, ">=", true},
        {"success_greater", 51, 50, ">=", true},
        {"failure_less", 49, 50, ">=", false},
        {"success_less_or_equal", 50, 50, "<=", true},
    }
    
    engine := NewCheckEngine()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            player := &entities.Player{}
            // Set stat value
            player.SetStat("str", tt.statValue)
            
            check := &Check{
                Type:      CheckTypeThreshold,
                Stat:      "str",
                Threshold: tt.threshold,
                Operator:  tt.operator,
            }
            
            result := engine.ExecuteCheck(check, player)
            assert.Equal(t, tt.wantSuccess, result.Success)
        })
    }
}
```

---

## 9. Scalability Considerations

### 9.1 Connection Pooling

```go
// internal/infrastructure/database/connection.go
func NewConnectionPool(cfg ConnectionConfig) (*pgxpool.Pool, error) {
    poolConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?pool_max_conns=%d",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.PoolMaxConns,
    ))
    if err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create pool: %w", err)
    }
    
    return pool, nil
}
```

### 9.2 Caching Strategy

```go
// Event templates caching
func (s *EventService) GetTemplates() ([]*entities.EventTemplate, error) {
    // Check cache first
    if cached, ok := s.cache.Get("event_templates"); ok {
        return cached.([]*entities.EventTemplate), nil
    }
    
    // Load from DB
    templates, err := s.repo.GetActiveTemplates()
    if err != nil {
        return nil, err
    }
    
    // Cache for 1 hour
    s.cache.Set("event_templates", templates, 1*time.Hour)
    
    return templates, nil
}
```

### 9.3 Future Scaling Path

| Current | Future | Trigger |
|---------|--------|---------|
| Single PostgreSQL | Read replicas | >100 concurrent users |
| In-memory cache | Redis cluster | >1000 concurrent users |
| Monolith | Microservices | Team grows beyond 5 |

---

## 10. Implementation Checklist

### Phase 1: Foundation (Steps 1-3)
- [ ] Step 1: Initialize Go project, config, logging
- [ ] Step 2: Create database migrations with correct schema
- [ ] Step 3: Implement domain entities (Player, Event, GameLog)

### Phase 2: Core Services (Steps 4-7)
- [ ] Step 4: Implement CheckEngine (Threshold, Probability, Catastrophic)
- [ ] Step 5: Implement EffectEngine with clamping
- [ ] Step 6: Implement ProgressionService (rank/status calculation)
- [ ] Step 7: Implement EventSelector (weighted selection, difficulty curve)

### Phase 3: Application Layer (Steps 8-10)
- [ ] Step 8: Implement FinalDeterminer
- [ ] Step 9: Implement GameService (orchestration)
- [ ] Step 10: Implement GraphQL handlers and schema

### Final Integration
- [ ] Connect all services in main.go
- [ ] Add error handling middleware
- [ ] Add CORS and security headers
- [ ] Write integration tests with testcontainers
- [ ] Run full test suite with coverage
- [ ] Verify API contracts match frontend expectations

---

**Документ завершён:** 2026-03-26  
**Следующий шаг:** Начать реализацию Phase 1 (Project Setup)