# Анализ архитектуры базы данных — «Армейка»

**Дата:** 2026-03-22  
**Анализ:** TZ.md, раздел 10.2 — Схема базы данных  
**Статус:** Требует исправлений

---

## 1. Выбор типа базы данных

**Рекомендация:** PostgreSQL 15+ (как указано в ТЗ)

**Обоснование:**
- Жёсткое соответствие ACID для целостности состояния игры
- Поддержка JSONB для гибкого хранения шаблонов событий (допустимо для read-only данных)
- Драйвер pgx обеспечивает отличную интеграцию с Go
- Connection pooling через pgxpool

**Не рекомендуется для MVP:**
- MongoDB — избыточная сложность, нет документально-ориентированного юзкейса
- Redis — эфемерное хранилище, не подходит для персистентного состояния игры

---

## 2. КРИТИЧЕСКИЕ ПРОБЛЕМЫ (Исправить до реализации)

### 2.1 Проблема #1: Статы хранятся как JSONB — нет ограничений

**Текущий дизайн:**
```sql
stats JSONB NOT NULL DEFAULT '{"str":50,"end":50,"agi":50,"mor":50,"disc":0}'
```

**Проблемы:**
1. ❌ Невозможно enforce `STR/END/AGI ∈ [1,100]`, `MOR ∈ [0,100]`, `DISC ∈ [-100,+100]` на уровне БД
2. ❌ Невозможно индексировать отдельные статы для аналитики
3. ❌ Невозможно эффективно искать «все игроки с MOR < 20»
4. ❌ Нет type safety — невалидный JSON тихо принимается

**Влияние:** Критическое нарушение целостности данных. Баг в приложении может установить `MOR = -50`.

**Рекомендация:**
```sql
CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    str INT NOT NULL DEFAULT 50 CHECK (str BETWEEN 1 AND 100),
    end_ INT NOT NULL DEFAULT 50 CHECK (end_ BETWEEN 1 AND 100),  -- end - зарезервированное слово
    agi INT NOT NULL DEFAULT 50 CHECK (agi BETWEEN 1 AND 100),
    mor INT NOT NULL DEFAULT 50 CHECK (mor BETWEEN 0 AND 100),
    disc INT NOT NULL DEFAULT 0 CHECK (disc BETWEEN -100 AND 100),
    formal_rank VARCHAR(50) NOT NULL DEFAULT 'РЯДОВОЙ',
    informal_status VARCHAR(50) NOT NULL DEFAULT 'ЗАПАХ',
    turn INT NOT NULL DEFAULT 1 CHECK (turn BETWEEN 1 AND 30),
    flags JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

**Обоснование:** Статы — атомарные целые числа с фиксированными диапазонами — идеальны для колонок с CHECK constraints.

---

### 2.2 Проблема #2: Дублирование данных — event_history в players

**Текущий дизайн:**
```sql
CREATE TABLE players (
    ...
    event_history JSONB DEFAULT '[]',  -- ❌ ДУБЛИРОВАНИЕ ДАННЫХ
);

CREATE TABLE game_logs (
    ...
    -- Эта таблица уже содержит те же данные
);
```

**Проблемы:**
1. ❌ Одни и те же данные хранятся в ДВУХ местах
2. ❌ Кошмар синхронизации — что является источником правды?
3. ❌ Потраченное место
4. ❌ Потенциальная несогласованность после сбоев

**Влияние:** Возможна порча данных. Сложная логика синхронизации.

**Рекомендация:**
- УДАЛИТЬ `event_history` из таблицы `players`
- `game_logs` — единственный источник правды
- Для производительности создать materialized view `current_player_state` или использовать Redis

---

### 2.3 Проблема #3: Избыточная таблица game_saves

**Текущий дизайн:**
```sql
CREATE TABLE game_saves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id),
    save_data JSONB NOT NULL,  -- Содержит всё состояние игры?
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(player_id)
);
```

**Проблемы:**
1. ❌ Что содержит `save_data`? Если статы/turn/flags → ДУБЛИРУЕТ `players`
2. ❌ Если game_logs — источник правды → зачем сохранять отдельно?
3. ❌ `UNIQUE(player_id)` = только ОДНО сохранение на игрока = нет бэкапа перед рискованным выбором
4. ❌ Неясно когда это пишется vs когда обновляется `players`

**Влияние:** Архитектура неясна. Возможны сценарии потери данных.

**Рекомендация:**

Вариант А — Если нужны сохранения для «отката»:
```sql
CREATE TABLE game_saves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id),
    player_state JSONB NOT NULL,  -- Снимок на момент сохранения
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX idx_game_saves_player_created ON game_saves(player_id, created_at DESC);
```
Хранить несколько сохранений на игрока.

Вариант Б — Если достаточно автосохранения:
- УДАЛИТЬ таблицу `game_saves`
- `players.updated_at` достаточно для «последняя игра»
- Состояние игры всегда восстановимо из `players` + `game_logs`

**Рекомендация:** Использовать Вариант Б для MVP. Добавить сохранения позже при необходимости.

---

### 2.4 Проблема #4: Отсутствует механизм оптимистичной блокировки

**Текущий дизайн:**
```sql
players (
    ...
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

**Проблемы:**
1. ❌ `updatedAt` меняется при каждом обновлении — нельзя обнаружить конкурентные изменения
2. ❌ Нет колонки `version`
3. ❌ Race condition: два параллельных `choose()` читают состояние игрока, оба инкрементят turn, оба сохраняют → turn инкрементится только ОДИН раз вместо двух

**Влияние:** Возможна порча состояния игры. Потерянные обновления.

**Рекомендация:**
```sql
CREATE TABLE players (
    ...
    version INT NOT NULL DEFAULT 1,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

Логика обновления:
```sql
UPDATE players 
SET 
    mor = $new_mor,
    turn = turn + 1,
    version = version + 1,
    updated_at = NOW()
WHERE id = $player_id AND version = $expected_version;
-- Если rows_affected = 0 → конкурентное изменение, повторить
```

---

### 2.5 Проблема #5: current_event_id ссылается на эфемерные данные

**Текущий дизайн:**
```sql
CREATE TABLE players (
    ...
    current_event_id UUID,  -- ❌ Ссылается на что?
);
```

**Проблемы:**
1. ❌ EventInstance генерируется на лету — не хранится в БД
2. ❌ Эта колонка ВСЕГДА будет NULL или ссылаться на мусор
3. ❌ GraphQL возвращает `currentEvent`, но нет соответствующей записи в БД

**Влияние:** Эта колонка не нужна. Вызовет путаницу.

**Рекомендация:**
- УДАЛИТЬ `current_event_id`
- Хранить `current_event_template_id VARCHAR(100)` если нужно отслеживать «в каком событии игрок сейчас»
- EventInstance генерируется динамически из шаблона + состояния игрока

---

## 3. ПРОБЛЕМЫ ВЫСОКОГО ПРИОРИТЕТА (Следует исправить)

### 3.1 Проблема #6: Недостаточно индексов

**Текущие индексы:**
```sql
CREATE INDEX idx_game_logs_player_id ON game_logs(player_id);
CREATE INDEX idx_game_logs_turn ON game_logs(player_id, turn DESC);
CREATE INDEX idx_event_templates_type ON event_templates(type);
CREATE INDEX idx_players_updated_at ON players(updated_at DESC);
```

**Недостающие индексы:**

| Паттерн запроса | Недостающий индекс |
|----------------|--------------------|
| Найти игроков по turn (например, «все игры на ходу 27+») | `idx_players_turn ON players(turn)` |
| Получить последнюю запись лога игрока | `idx_game_logs_player_turn` (уже есть) |
| Найти завершённые игры (turn = 30 ИЛИ mor = 0) | Нужен partial index |
| Шаблоны событий по использованию (взвешенная случайная выборка) | `idx_event_templates_used ON event_templates(used_count)` |
| Игры для очистки (старше 90 дней, завершены) | Partial index на `created_at` |

**Рекомендация:**
```sql
-- Для аналитических запросов
CREATE INDEX idx_players_turn ON players(turn);

-- Для поиска завершённых игр
CREATE INDEX idx_players_finished ON players(turn) WHERE turn = 30;

-- Для оптимизации выбора событий (взвешенный случайный)
CREATE INDEX idx_event_templates_used_count ON event_templates(used_count);
```

---

### 3.2 Проблема #7: Нет стратегии архивации / очистки

**Текущий дизайн:** Нет механизма архивировать или очищать завершённые игры.

**Проблемы:**
1. ❌ Таблица `players` растёт бесконечно
2. ❌ Аналитические запросы замедляются со временем
3. ❌ Проблемы GDPR/privacy (хранение данных игроков вечно)

**Рекомендация:**
```sql
-- Добавить колонки для архивации
ALTER TABLE players ADD COLUMN is_finished BOOLEAN DEFAULT FALSE;
ALTER TABLE players ADD COLUMN finished_at TIMESTAMP WITH TIME ZONE;

-- Индекс для job архивации
CREATE INDEX idx_players_finished_at ON players(is_finished, finished_at) 
WHERE is_finished = TRUE;

-- Ежемесячный cron job:
-- DELETE FROM players WHERE is_finished = TRUE AND finished_at < NOW() - INTERVAL '90 days';
```

---

### 3.3 Проблема #8: Race condition в event_templates.used_count

**Текущий дизайн:**
```sql
event_templates (
    used_count INT DEFAULT 0,
    last_used_at TIMESTAMP WITH TIME ZONE
);
```

**Проблема:**
Конкурентные игровые сессии одновременно инкрементят `used_count` → потерянные обновления.

**Рекомендация:**
```sql
-- Использовать атомарное обновление
UPDATE event_templates 
SET used_count = used_count + 1, 
    last_used_at = NOW()
WHERE id = $template_id;
-- Это атомарно в PostgreSQL
```

---

### 3.4 Проблема #9: Нет рекомендаций по размеру Connection Pool

**Текущее:** "Connection pool (PostgreSQL) | 25 connections"

**Проблема:**
Произвольное число. Может быть слишком мало или слишком много.

**Рекомендация:**
- Правило: `pool_size = (number_of_CPU_cores * 2) + effective_spindle_count`
- Для выделенного DB сервера с SSD: 50-100 соединений
- Мониторить `pg_stat_activity` для настройки
- Использовать PgBouncer для connection pooling перед PostgreSQL если нужно

---

## 4. ПРОБЛЕМЫ СРЕДНЕГО ПРИОРИТЕТА

### 4.1 Проблема #10: event_templates загружается из JSON файлов

**Текущее:** Таблица `event_templates` заполняется из JSON файлов при старте.

**Проблемы:**
1. ❌ Каждый рестарт сервера перезагружает шаблоны (или полагается на предсуществующие данные)
2. ❌ Нет версионирования шаблонов
3. ❌ Нет отката если плохой шаблон задеплоен

**Рекомендация:**
```sql
-- Добавить версионирование шаблонов
ALTER TABLE event_templates ADD COLUMN version INT DEFAULT 1;
ALTER TABLE event_templates ADD COLUMN is_active BOOLEAN DEFAULT TRUE;

-- Миграционный скрипт загружает новые шаблоны, старые помечает как неактивные
-- Хранить историю для игр в процессе
```

---

### 4.2 Проблема #11: Нет стратегии обработки NULL

**Текущее:** Некоторые колонки допускают NULL (например, `current_event_id`)

**Проблема:**
Семантика NULL различается, усложняет запросы.

**Рекомендация:**
- Заменить NULLable UUID на пустую строку `''` или использовать `COALESCE` везде
- Быть явным в том, что означает NULL

---

### 4.3 Проблема #12: Недостаёт Unique Constraint в game_logs

**Текущее:** Нет уникального ограничения на `(player_id, turn)`

**Проблема:**
Можно вставить дублирующиеся записи лога для одного хода.

**Рекомендация:**
```sql
ALTER TABLE game_logs 
ADD CONSTRAINT uq_game_logs_player_turn UNIQUE (player_id, turn);
```

---

## 5. ПРОБЛЕМЫ НИЗКОГО ПРИОРИТЕТА (Желательно иметь)

### 5.1 Проблема #13: Не определены Composite Types

**Проблема:**
GraphQL типы (PlayerStats, Effect и т.д.) не имеют эквивалентов в БД.

**Рекомендация:**
```sql
CREATE TYPE player_stats AS (
    str INT,
    end INT,
    agi INT,
    mor INT,
    disc INT
);

-- Теперь можно использовать:
-- player_stats DEFAULT ROW(50, 50, 50, 50, 0)
-- Но для MVP отдельные колонки подходят
```

---

### 5.2 Проблема #14: Нет стратегии партиционирования для масштабирования

**Текущее:** Одна таблица `players`.

**Будущая проблема:**
При масштабировании за пределы MVP:

```sql
-- Партиционирование по месяцам (при публичном запуске)
CREATE TABLE players (
    ...
) PARTITION BY RANGE (created_at);
```

---

## 6. АНАЛИЗ ПАТТЕРНОВ ДОСТУПА

### 6.1 Частота операций

| Операция | Частота | Паттерн |
|----------|---------|---------|
| `startGame()` | Низкая | Вставить нового игрока |
| `choose()` | Высокая | Прочитать игрока → Обновить игрока → Вставить лог |
| `loadGame()` | Средняя | Прочитать игрока + недавние логи |
| `eventHistory()` | Средняя | Прочитать логи с пагинацией |

### 6.2 Горячий путь: Мутация `choose()`

```sql
-- Текущий (проблемный):
BEGIN;
SELECT * FROM players WHERE id = $1;  -- ❌ Без блокировки
-- Обработать выбор, рассчитать эффекты
UPDATE players SET ... WHERE id = $1;  -- ❌ Race condition
INSERT INTO game_logs ...;
COMMIT;
```

**Рекомендуемый (с оптимистичной блокировкой):**
```sql
BEGIN;
SELECT * FROM players WHERE id = $1 FOR UPDATE;  -- Row lock
-- Обработать выбор
UPDATE players SET version = version + 1, ... WHERE id = $1 AND version = $expected;
-- Проверить rows_affected
INSERT INTO game_logs ...;
COMMIT;
```

---

## 7. ИСПРАВЛЕННАЯ СХЕМА БАЗЫ ДАННЫХ

```sql
-- ============================================================================
-- ТАБЛИЦА PLAYERS — Единственный источник правды о состоянии игрока
-- ============================================================================
CREATE TABLE players (
    -- Идентификация
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Статы с ограничениями (НЕ JSONB)
    str INT NOT NULL DEFAULT 50 CHECK (str BETWEEN 1 AND 100),
    end_ INT NOT NULL DEFAULT 50 CHECK (end_ BETWEEN 1 AND 100),  -- end - зарезервированное слово
    agi INT NOT NULL DEFAULT 50 CHECK (agi BETWEEN 1 AND 100),
    mor INT NOT NULL DEFAULT 50 CHECK (mor BETWEEN 0 AND 100),
    disc INT NOT NULL DEFAULT 0 CHECK (disc BETWEEN -100 AND 100),
    
    -- Прогрессия
    formal_rank VARCHAR(50) NOT NULL DEFAULT 'РЯДОВОЙ',
    informal_status VARCHAR(50) NOT NULL DEFAULT 'ЗАПАХ',
    
    -- Состояние игры
    turn INT NOT NULL DEFAULT 1 CHECK (turn BETWEEN 1 AND 30),
    flags JSONB DEFAULT '[]',
    current_event_template_id VARCHAR(100),
    
    -- Метаданные
    is_finished BOOLEAN NOT NULL DEFAULT FALSE,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    version INT NOT NULL DEFAULT 1  -- Оптимистичная блокировка
);

-- Индексы
CREATE INDEX idx_players_updated_at ON players(updated_at DESC);
CREATE INDEX idx_players_turn ON players(turn);
CREATE INDEX idx_players_finished ON players(is_finished) WHERE is_finished = FALSE;
CREATE INDEX idx_players_finished_at ON players(is_finished, finished_at) 
    WHERE is_finished = TRUE;

-- ============================================================================
-- ТАБЛИЦА GAME_LOGS — Неизменяемая история всех выборов
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

-- Индексы
CREATE INDEX idx_game_logs_player_id ON game_logs(player_id);
CREATE INDEX idx_game_logs_player_turn ON game_logs(player_id, turn DESC);
CREATE INDEX idx_game_logs_created_at ON game_logs(created_at);

-- ============================================================================
-- ТАБЛИЦА EVENT_TEMPLATES — Контент игры (только для чтения)
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

-- Индексы
CREATE INDEX idx_event_templates_type ON event_templates(type);
CREATE INDEX idx_event_templates_type_active ON event_templates(type) 
    WHERE is_active = TRUE;
CREATE INDEX idx_event_templates_used_count ON event_templates(used_count);

-- ============================================================================
-- MATERIALIZED VIEW — Текущее резюме игроков (опциональная оптимизация)
-- ============================================================================
CREATE MATERIALIZED VIEW player_summary AS
SELECT 
    p.id,
    p.turn,
    p.mor,
    p.disc,
    p.formal_rank,
    p.informal_status,
    p.is_finished,
    COUNT(gl.id) AS total_choices,
    MAX(gl.created_at) AS last_action_at
FROM players p
LEFT JOIN game_logs gl ON p.id = gl.player_id
GROUP BY p.id;

CREATE UNIQUE INDEX idx_player_summary_id ON player_summary(id);
```

---

## 8. ПЛАН МИГРАЦИИ

### Этап 1: Быстрые исправления (Перед первым деплоем)
1. Добавить CHECK constraints к таблице `players`
2. Добавить колонку `version` для оптимистичной блокировки
3. Удалить `event_history` из таблицы players
4. Решить судьбу `game_saves` (удалить или переделать с множественными сохранениями)
5. Удалить/обновить `current_event_id`

### Этап 2: Стабилизация (После запуска MVP)
1. Добавить колонки и job для архивации
2. Создать partial indexes для частых запросов
3. Добавить версионирование шаблонов

### Этап 3: Масштабирование (Когда потребуется)
1. Добавить connection pooling (PgBouncer)
2. Рассмотреть партиционирование по дате
3. Добавить read replicas

---

## 9. СВОДКА РЕКОМЕНДАЦИЙ

| Приоритет | Проблема | Действие |
|-----------|----------|----------|
| 🔴 КРИТИЧЕСКАЯ | Статы как JSONB | Изменить на отдельные колонки с CHECK |
| 🔴 КРИТИЧЕСКАЯ | Дублирование event_history | Удалить из таблицы players |
| 🔴 КРИТИЧЕСКАЯ | Избыточная game_saves | Удалить (или переделать с множественными сохранениями) |
| 🔴 КРИТИЧЕСКАЯ | Нет оптимистичной блокировки | Добавить колонку version |
| 🟡 ВЫСОКИЙ | Недостаточно индексов | Добавить turn, partial indexes |
| 🟡 ВЫСОКИЙ | Нет архивации | Добавить is_finished, cleanup job |
| 🟡 ВЫСОКИЙ | Race condition used_count | Использовать атомарное UPDATE |
| 🟢 СРЕДНИЙ | Версионирование шаблонов | Добавить колонку version |
| 🟢 СРЕДНИЙ | Обработка NULL | Быть явным |

---

**Анализ завершён:** 2026-03-22  
**Следующий шаг:** Обсудить с backend-командой, обновить TZ.md раздел 10.2
