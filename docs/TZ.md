# Техническая спецификация — Narrative Survival Game «Армейка»

**Версия:** 1.5  
**Дата:** 2026-03-25  
**Статус:** Черновик  
**История изменений:** 
- v1.5 — Обновлены описания финалов: «Тихий дембель» (Приспособленец), «Уважаемый дембель», «Сломанный дембель» (Косячник). Добавлено поле subtitle в GraphQL schema. Добавлена ссылка на docs/DESIGN_SPECIFICATION.md в разделы 1 и 14. Добавлен раздел 10.2 Frontend Architecture (FSD). Добавлен раздел 10.3 Backend Architecture (Go + GraphQL). Добавлен раздел 15.2 Testing Strategy (TDD). Добавлен раздел 15.3 CI/CD Pipeline.
- v1.4 — Добавлены 42 narrative события (баланс проверен: Death 34.6%, Victory 65.4%, Avg MOR/turn -1.48)
- v1.3 — Section 7.6 обновлён на основе Monte Carlo симуляции (1000 runs): параметры difficulty_modifier, event distribution, MOR economy, метрики baseline
- v1.2 — добавлены quantitative balance requirements, difficulty curve, state synchronization, missing edge cases
- v1.1 — исправлена схема БД (stats как колонки, удалён game_saves, добавлен optimistic locking)

---

## 1. Overview

### 1.1 Что строится

Веб-игра «Армейка» — narrative survival / choice-based RPG, симулирующая опыт срочной военной службы. Игрок — солдат, вынужденный выживать в системе постоянного давления между формальной уставной иерархией и неформальной солдатской средой.

> **Документация:** Полная спецификация UI/UX доступна в `docs/DESIGN_SPECIFICATION.md`. Дизайн-система и промты для Figma находятся там же.

### 1.2 Зачем существует

- Проверить механику «конфликта двух иерархий» как основу геймплея
- Собрать фидбек от 10–20 тестировщиков
- Создать foundation для будущего масштабирования (мобильные платформы, Steam)

### 1.3 Кто будет использовать

| Роль | Описание |
|------|----------|
| Игрок | Анонимный пользователь, проходящий игру сессиями 15–30 минут |
| Разработчик | Fullstack-разработчик, поддерживающий и расширяющий систему |

---

## 2. Goals / Non-Goals

### 2.1 Goals (Что МUST быть реализовано)

- [ ] Полный игровой цикл: событие → описание → выбор → проверка → последствия → следующий ход
- [ ] 5 игровых параметров с валидными диапазонами: STR/END/AGI ∈ [1,100], MOR ∈ [0,100], DISC ∈ [-100,+100]
- [ ] Система LEGO-шаблонов событий с подстановкой переменных
- [ ] 3 типа проверок: Threshold, Probability, Catastrophic
- [ ] Две параллельные иерархии: формальная (4 звания) и неформальная (6 статусов)
- [ ] 3 финала: «Тихий дембель» (Приспособленец), «Уважаемый дембель», «Сломанный дембель» (Косячник)
- [ ] Автосохранение после каждого хода
- [ ] История последних 10 событий
- [ ] GraphQL API для взаимодействия frontend ↔ backend
- [ ] RESTART (новая игра) без перезагрузки страницы
- [ ] Weighted event selection (не pure random)
- [ ] Difficulty curve (легче в начале, сложнее к концу)
- [ ] MOR recovery guarantee (после 3 негативных — позитивное событие)
- [ ] Оптимистичная блокировка (concurrency protection)
- [ ] Баланс протестирован через симулятор до implementation

### 2.2 Non-Goals (Что запрещено делать в MVP)

- [ ] Пользовательская авторизация (Google OAuth, JWT)
- [ ] Админка для создания событий
- [ ] AI-арты и иллюстрации (только placeholder UI)
- [ ] Flutter / мобильная версия
- [ ] Steam achievements и таблицы лидеров
- [ ] Монетизация
- [ ] Локализация на английский язык
- [ ] Полноценный саундтрек (только базовые SFX)
- [ ] P2W механики (pay-to-win ЗАПРЕЩЕНЫ)

---

## 3. Definitions

### 3.1 Игровые термины

| Термин | Определение |
|--------|-------------|
| **Turn (Ход)** | Один день службы. Максимум 30 ходов в игре. |
| **Player (Игрок)** | Сущность, представляющая солдата. Содержит статы, прогрессию, флаги. |
| **Stat (Стат)** | Игровой параметр, влияющий на проверки и доступность выборов. |
| **Event (Событие)** | Ситуация, с которой сталкивается игрок. Состоит из описания и вариантов выбора. |
| **EventTemplate (Шаблон события)** | LEGO-блок для генерации EventInstance. Хранится в JSON. |
| **EventInstance (Экземпляр события)** | Конкретное событие с подставленными переменными. |
| **Choice (Выбор)** | Вариант действия игрока. Имеет условия доступности и результаты. |
| **Check (Проверка)** | Механизм разрешения исхода выбора на основе статов. |
| **Effect (Эффект)** | Изменение статов после выбора. |
| **Flag (Флаг)** | Булевый маркер, влияющий на будущие события. |
| **FormalRank (Формальное звание)** | Позиция в уставной иерархии. Зависит от DISC. |
| **InformalStatus (Неформальный статус)** | Позиция в солдатской иерархии. Зависит от DISC. |
| **Final (Финал)** | Конечный результат игры, определяемый совокупностью параметров. |

### 3.2 Технические термины

| Термин | Определение |
|--------|-------------|
| **MVP** | Minimum Viable Product — первая играбельная версия |
| **LEGO-template** | Шаблон события с переменными, заменяемыми при генерации |
| **Clamp** | Ограничение значения в пределах допустимого диапазона |
| **Autonomous Save** | Автоматическое сохранение состояния игры на сервере |

---

## 4. User Roles

### 4.1 Guest (Гость)

**Capabilities:**
- Начать новую игру
- Играть в рамках сессии
- Получить случайный финал
- Начать новую игру после завершения

**Restrictions:**
- Нет авторизации
- Нет доступа к сохранениям других игроков
- Нет админских функций
- Все игры анонимны (один активный gameId на браузер)

---

## 5. User Scenarios

### 5.1 Основной игровой поток

```
1. Пользователь открывает приложение
   → Система отображает главный экран с кнопкой "Начать игру"

2. Пользователь нажимает "Начать игру"
   → Система вызывает startGame mutation
   → Сервер создаёт Player со стартовыми статами
   → Сервер генерирует первое EventInstance из пула шаблонов
   → Клиент получает Player + EventInstance
   → UI отображает событие и варианты выбора

3. Пользователь читает описание события
   → UI отображает текст с подставленными переменными
   → UI отображает доступные варианты выбора (2–4 кнопки)

4. Пользователь нажимает на вариант выбора
   → Система вызывает choose mutation
   → Сервер выполняет проверку (check)
   → Сервер применяет эффекты (effects) с clamp
   → Сервер обновляет прогрессию (если изменились пороги)
   → Сервер проверяет условия финала
   → Сервер автосохраняет состояние
   → Сервер генерирует следующее событие (или определяет финал)
   → Клиент получает обновлённый Player + результат + следующее событие

5. Цикл повторяется (шаги 3–4) до:
   - Turn = 30 → финал «Дембель»
   - MOR = 0 → финал «Сломанный»
   - Дембель (turn ≥ 27 и DISC ∈ [-20,+20]) → финал «Тихий»

6. Игра окончена
   → UI отображает экран финала с описанием и статистикой
   → UI отображает кнопку "Играть снова"
```

### 5.2 Сценарий прерванной игры

```
1. Пользователь открывает приложение
   → Система проверяет localStorage на наличие gameId
   → Если gameId существует:
     → Система вызывает loadGame query
     → Сервер возвращает сохранённое состояние
     → UI отображает текущий прогресс
   → Если gameId не существует:
     → UI отображает главный экран
```

### 5.3 Сценарий недоступного выбора

```
1. UI отображает событие с вариантами
2. Система вычисляет availability для каждого choice
   → Если availability не satisfied: кнопка disabled
   → Если availability satisfied: кнопка enabled
3. Пользователь видит disabled кнопки и tooltip "Недоступно"
4. Пользователь выбирает только доступные варианты
```

---

## 6. Functional Requirements

### 6.1 FR-001: Система статов

| ID | Требование | Валидация |
|----|------------|-----------|
| FR-001.1 | Система MUST хранить 5 статов в БД как отдельные колонки: STR, END, AGI ∈ [1,100], MOR ∈ [0,100], DISC ∈ [-100,+100] | Unit-тест + CHECK constraints |
| FR-001.2 | БД MUST enforce CHECK constraints на диапазоны статов | PostgreSQL constraint |
| FR-001.3 | Все эффекты MUST быть clamped к допустимым диапазонам после применения (на уровне приложения + БД) | Unit-тест |
| FR-001.4 | Начальные статы MUST быть: STR=50, END=50, AGI=50, MOR=50, DISC=0 | Интеграционный тест |
| FR-001.5 | DISC < 0 интерпретируется как «неформальная ориентация» | Документация |
| FR-001.6 | DISC > 0 интерпретируется как «формальная ориентация» | Документация |

### 6.2 FR-002: Система событий

| ID | Требование | Валидация |
|----|------------|-----------|
| FR-002.1 | Система MUST поддерживать минимум 20 EventTemplate в пуле | Интеграционный тест |
| FR-002.2 | EventTemplate MUST содержать: id, type, tags, context, template, choices | JSON-схема |
| FR-002.3 | Система MUST выбирать события с анти-рандомизацией (не повторять последние 5) | Unit-тест |
| FR-002.4 | Генерация EventInstance MUST подставлять случайные значения переменных | Unit-тест |
| FR-002.5 | EventInstance MUST фильтровать choices по availability перед отправкой клиенту | Unit-тест |
| FR-002.6 | Пул событий MUST содержать: 40% Routine, 25% Social, 15% Inspection, 15% Informal, 5% Emergency | JSON-контент |

### 6.3 FR-003: Система проверок

| ID | Требование | Валидация |
|----|------------|-----------|
| FR-003.1 | Threshold check: `formula >= threshold` → success, иначе failure | Unit-тест |
| FR-003.2 | Probability check: `rand(100) < calculatedChance` → success/partial/failure | Unit-тест (mock rand) |
| FR-003.3 | Catastrophic check: noticeChance → если noticed → powerCheck → исход | Unit-тест |
| FR-003.4 | Формулы MUST поддерживать: STAT, коэффициенты (*0.3), константы | JSON-парсер |
| FR-003.5 | Формулы MUST NOT выполнять произвольный код | Без eval() |
| FR-003.6 | Проверка availability MUST производиться до выполнения check | Unit-тест |

### 6.4 FR-004: Система прогрессии

| ID | Требование | Валидация |
|----|------------|-----------|
| FR-004.1 | Формальная иерархия MUST содержать 4 звания: Рядовой, Ефрейтор, Мл.сержант, Сержант | Константы |
| FR-004.2 | Неформальная иерархия MUST содержать 6 статусов: Запах, Дух, Слон, Черпак, Дед, Дембель | Константы |
| FR-004.3 | Формальное звание MUST обновляться при DISC > +25/+50/+75 | Unit-тест |
| FR-004.4 | Неформальный статус MUST обновляться при DISC < -25/-50/-75/-90 | Unit-тест |
| FR-004.5 | Повышение по формальной иерархии MUST давать бонусы к офицерским проверкам | Документация |
| FR-004.6 | Повышение по неформальной иерархии MUST давать бонусы к солдатским проверкам | Документация |

### 6.5 FR-005: Система финалов

| ID | Требование | Валидация |
|----|------------|-----------|
| FR-005.1 | Финал «Тихий дембель» (Приспособленец): turn ≥ 27 AND MOR > 0 AND DISC ∈ [-20,+20] | Unit-тест |
| FR-005.2 | Финал «Уважаемый дембель»: turn ≥ 27 AND MOR > 0 AND DISC < -50 AND informalStatus ≥ «Дед» | Unit-тест |
| FR-005.3 | Финал «Сломанный дембель» (Косячник): MOR = 0 (в любой момент) | Unit-тест |
| FR-005.4 | При достижении финала MUST прекращаться генерация событий | Интеграционный тест |
| FR-005.5 | Финал MUST содержать: title, subtitle (тип), description, finalStats | API-контракт |

### 6.6 FR-006: Сохранение и загрузка

| ID | Требование | Валидация |
|----|------------|-----------|
| FR-006.1 | Автосохранение MUST происходить после каждого choose | Интеграционный тест |
| FR-006.2 | Состояние игры MUST храниться в двух таблицах: players (текущее состояние) и game_logs (история) | Unit-тест |
| FR-006.3 | Таблица players MUST содержать: id, str, end_, agi, mor, disc, formal_rank, informal_status, turn, flags, version, is_finished, finished_at | Schema verification |
| FR-006.4 | Загрузка MUST восстанавливать состояние из players + последние 10 записей из game_logs | Интеграционный тест |
| FR-006.5 | При отсутствии сохранения MUST создаваться новый Player | Unit-тест |
| FR-006.6 | Каждый UPDATE players MUST инкрементировать version | Unit-тест |
| FR-006.7 | mutation choose MUST проверять expectedVersion против текущей версии | Интеграционный тест |

### 6.7 FR-007: API-контракт

| ID | Требование | Валидация |
|----|------------|-----------|
| FR-007.1 | Query player(id) MUST возвращать Player со всеми статами и текущей version | GraphQL schema |
| FR-007.2 | Query currentEvent(playerId) MUST возвращать динамически сгенерированный EventInstance | GraphQL schema |
| FR-007.3 | Query eventHistory(playerId, limit) MUST возвращать массив последних событий из game_logs | GraphQL schema |
| FR-007.4 | Mutation startGame() MUST создавать нового Player и возвращать EventInstance | GraphQL schema |
| FR-007.5 | Mutation choose(playerId, choiceId, expectedVersion) MUST применять эффекты с проверкой версии и возвращать результат + следующее событие + newVersion | GraphQL schema |
| FR-007.6 | Все mutations MUST валидировать входные данные | Integration тест |
| FR-007.7 | При ошибке валидации MUST возвращаться structured error | GraphQL errors |
| FR-007.8 | При CONCURRENT_MODIFICATION MUST возвращаться error с code VERSION_MISMATCH | Integration тест |

---

## 7. Non-Functional Requirements

### 7.1 Производительность

| Параметр | Значение | Метрика |
|----------|---------|---------|
| Время отклика API (p95) | < 200ms | Grafana |
| Время генерации события | < 50ms | Benchmark |
| Время применения эффектов | < 10ms | Benchmark |
| Размер бандла frontend | < 500KB (gzipped) | Lighthouse |
| TTFB | < 100ms | Web Vitals |

### 7.2 Надёжность

| Параметр | Значение | Метрика |
|----------|---------|---------|
| Uptime | 99.5% | SLA |
| Recovery Time Objective (RTO) | < 5 min | Runbook |
| Recovery Point Objective (RPO) | < 1 min | Backup schedule |

### 7.3 Масштабируемость

| Параметр | Значение |
|----------|---------|
| Concurrent users (MVP) | 100 |
| Connection pool (PostgreSQL) | 25 connections |
| Redis cache TTL | 1 hour |

### 7.4 Безопасность

| Требование | Реализация |
|------------|------------|
| Нет eval() или runtime code execution | Формулы парсятся через lexer |
| CORS настроен на разрешённые origins | Middleware |
| Rate limiting | 100 req/min per IP |
| Нет секретов в коде | Environment variables |
| Input validation | Все GraphQL inputs валидируются |

### 7.5 Мониторинг

| Компонент | Инструмент |
|-----------|------------|
| Логи | Structured JSON → stdout → Loki |
| Метрики | Prometheus metrics endpoint |
| Трейсинг | OpenTelemetry (опционально) |

### 7.6 Баланс и игровая экономика (Quantitative Requirements)

**Важно:** Эти требования критичны для success игры. Нарушение приводит к P2W-подобному опыту или фрустрации (см. Hoosegow failures в ANALYSIS.md).

**История изменений:**
- v1.3 — Обновлено на основе Monte Carlo симуляции (1000+ runs). Параметры проверены симулятором до implementation.

#### 7.6.1 Ограничения на статы

| Параметр | Min | Max | Примечание |
|----------|-----|-----|------------|
| MOR decay per event | - | -10 | EMERGENCY failure может быть до -10 |
| MOR gain per event | - | +4 | SAFE success обычно +2-4 |
| DISC change per event | -10 | +10 | Ограничение на изменение |

**Обоснование:** Симуляция показала, что max MOR decay = -10 (EMERGENCY) создаёт оптимальный баланс: игроки умирают на turn 26-27 (3-4 хода до финиша), создавая "почти сделал" ощущение без hopelessness.

#### 7.6.2 Распределение событий по сложности

| Фаза игры | Turns | Easy (Routine/Safe) | Medium (Social/Informal) | Hard (Inspection/Emergency) |
|-----------|-------|---------------------|------------------------|----------------------------|
| Early | 1-10 | 35% | 35% | 30% |
| Mid | 11-20 | 35% | 35% | 30% |
| Late | 21-30 | 35% | 35% | 30% |

**Изменение v1.3:** Распределение одинаковое по всем фазам. Difficulty curve реализуется через difficulty_modifier (Section 7.6.4), а не через event distribution.

**Обоснование:** Тестовая симуляция с 60/30/10 early, 30/50/20 mid, 10/40/50 late показала Too Hard в early game. Uniform distribution + probability modifier = лучший баланс.

#### 7.6.3 Recovery events (MOR-positive)

| Тараметр | Значение | Обоснование |
|----------|---------|-------------|
| MIN событий с MOR+ в пуле | 15-20% | SAFE events составляют 15% пула |
| MAX MOR- событий подряд | 2 | После 2 негативных — гарантированно SAFE event |
| Recovery boost | +2 to +4 MOR | SAFE success даёт +2-4 MOR |

**Обоснование:** recovery_trigger = 2 создаёт "reset loop" каждые 3-4 хода, предотвращая death spiral. recovery_trigger = 3 рекомендуется для production (больше challenge).

#### 7.6.4 Difficulty Curve (Probability Modifier)

```
difficulty_modifier = {
  "turn_1_10": 0.60,   // base_prob / 0.60 = base_prob * 1.67 (67% easier)
  "turn_11_20": 0.85,  // base_prob / 0.85 = base_prob * 1.18 (18% easier)
  "turn_21_30": 1.05   // base_prob / 1.05 = base_prob * 0.95 (5% harder)
}

Эффективная probability = base_prob / difficulty_modifier[current_phase]
```

**Эффект:** Игрок имеет ~67% higher success probability в early game, создавая "honeymoon period".

#### 7.6.5 Weighted Event Selection

```
Алгоритм выбора события:
1. Получить все активные EventTemplate
2. Для каждого template вычислить weight:
   - base_weight = 1.0
   - Если type = "SAFE": +0.2 если consecutive_negative >= 1
   - Если type = "INSPECTION": +0.3 если DISC > +50
   - Если type = "INFORMAL": +0.3 если DISC < -50
   - Если type = "ROUTINE": +0.1 если turn <= 10
   - Если type = "EMERGENCY": +0.2 если turn >= 20
   - Если template_id ∈ recent_history[последние 5]: weight = 0
   - Если template.used_count > 20: weight *= 0.5
3. Выбрать случайный template с probability proportional to weight
4. Не выбирать если все weights = 0 (reset history)
```

#### 7.6.6 Целевые метрики баланса (v1.4 Narrative Events)

| Метрика | Baseline (v1.4) | Target | Status | Notes |
|---------|----------|--------|--------|-------|
| Death Rate | 34.6% | 20-35% | ✅ | 42 narrative events, 1000 runs |
| Victory Rate | 65.4% | 65-80% | ✅ | Turn 30 с MOR > 0 |
| Perfect Runs | 0.0% | 0-5% | ✅ | Turn 30, MOR > 50 |
| Avg MOR/turn | -1.48 | -2.5 to -1.5 | ✅ | Balance approved |
| Avg Death Turn | 26.7 | 12-28 | ✅ | 3-4 хода до финиша |
| Avg Success Rate | ~40% | 40-60% | ✅ | Per choice statistics |

**Примечание:** 42 narrative events покрывают все 6 типов (ROUTINE, SAFE, SOCIAL, INFORMAL, INSPECTION, EMERGENCY). События имеют иммерсивный русский military narrative.

#### 7.6.7 MOR Economy (Подробно)

| Event Type | Success MOR | Partial MOR | Failure MOR | % пула |
|------------|-------------|-------------|-------------|--------|
| ROUTINE | +0 | -2 to -3 | -5 to -6 | 20% |
| SAFE | +2 to +4 | 0 to -1 | -2 to -3 | 15% |
| SOCIAL | +0 to +1 | -2 to -3 | -5 to -7 | 20% |
| INFORMAL | -1 to +0 | -3 to -4 | -6 to -8 | 15% |
| INSPECTION | -1 to +0 | -3 to -5 | -6 to -9 | 15% |
| EMERGENCY | +0 | -4 to -5 | -7 to -10 | 15% |

**Expected MOR drain:** ~50 MOR за 30 turns. Starting MOR = 50. Финиш с MOR ≈ 0 = "тихий дембель".

---

## 8. System Behavior

### 8.1 Диаграмма состояний игры

```
                    ┌─────────────┐
                    │   START     │
                    └──────┬──────┘
                           │ startGame()
                           ▼
                    ┌─────────────┐
              ┌───▶│   PLAYING   │
              │     └──────┬──────┘
              │            │
              │    choose()│
              │            ▼
              │     ┌─────────────┐    turn = 30 OR
              │     │   CHECKING  │──── MOR = 0 OR
              │     └──────┬──────┘    turn ≥ 27 AND conditions
              │            │
              │            ▼
              │     ┌─────────────┐
              │     │   FINISHED  │
              │     └─────────────┘
              │
              │ restart()
              │
              └─────────────▶ (back to START)
```

### 8.2 Алгоритм выбора следующего события

```
1. Получить пул доступных EventTemplate
2. Исключить: templateId ∈ последних 5 использованных
3. Если DISC > 50: приоритет Inspection событий (+20% weight)
4. Если DISC < -50: приоритет Informal событий (+20% weight)
5. Если MOR < 30: приоритет Safe событий (+30% weight)
6. Выбрать случайный template по весам
7. Сгенерировать EventInstance:
   a. Подставить случайные значения для каждой variable
   b. Фильтровать choices по availability
   c. Возвращать только доступные choices
8. Сохранить templateId в history
9. Вернуть EventInstance
```

### 8.3 Алгоритм применения эффектов

```
1. Для каждого effect в outcome.effects:
   a. stat = player.stats[effect.stat]
   b. newValue = stat + effect.delta
   c. clampedValue = clamp(newValue, min, max)
   d. player.stats[effect.stat] = clampedValue
2. Обновить флаги (flags.push(...effect.flags))
3. Пересчитать формальную прогрессию:
   a. Если DISC >= 25 AND rank != "Ефрейтор": rank = "Ефрейтор"
   b. Если DISC >= 50 AND rank != "Мл.сержант": rank = "Мл.сержант"
   c. Если DISC >= 75 AND rank != "Сержант": rank = "Сержант"
4. Пересчитать неформальную прогрессию:
   a. Если DISC <= -25 AND status != "Дух": status = "Дух"
   b. Если DISC <= -50 AND status != "Слон": status = "Слон"
   c. Если DISC <= -75 AND status != "Черпак": status = "Черпак"
   d. Если DISC <= -90 AND status != "Дед": status = "Дед"
5. Сохранить player
```

### 8.4 Алгоритм проверки финала

```
1. Если MOR = 0: 
   return "Сломанный дембель" (Косячник)
   Description: "Служба в армии принесла вам лишь горе, затрещины и сломанный нос, когда вы ночью упали с кровати, исполняя 'летучую мышь' по команде сержанта Златопузова. Впрочем, вы выжили — а это главное."

2. Если turn ≥ 27 AND MOR > 0 AND DISC ∈ [-20,+20]: 
   return "Тихий дембель" (Приспособленец)
   Description: "Ваша служба прошла неспеша, вы не выделялись, не косячили, выполняли приказы 'шакалов' и офицеров. В общем — вы просто жили."

3. Если turn ≥ 27 AND MOR > 0 AND DISC < -50 AND informalStatus ≥ "Дед": 
   return "Уважаемый дембель"
   Description: "Вы ушли со службы в золотом кашне. Вся рота провожала вас со слезами на глазах, и обсуждала, как вы смогли так 'подняться' за такой короткий срок. А ваш кореш — старший прапорщик Залупенко плакал от горя, ведь вы смогли 'загнать' 300 тонн гнутых гвоздей..."

4. Если turn < 30 AND MOR > 0: return null (продолчать игру)
```

### 8.5 Difficulty Curve Algorithm

```
Входные данные: player, turn
Выходные данные: difficulty modifier для проверок

1. Определить phase:
   - Early: turn <= 10
   - Mid: 11 <= turn <= 20
   - Late: turn > 20

2. Вычислить base difficulty:
   - Early: 0.9 (легче)
   - Mid: 1.0 (нормально)
   - Late: 1.1 (сложнее)

3. Применить modifiers:
   - Если formalRank = "Сержант": difficulty *= 1.1 (выше ожидания)
   - Если informalStatus = "Дед": difficulty *= 1.1 (давление)
   - Если flags.contains("hospital"): difficulty *= 0.9 (восстановление)

4. Return clamped(difficulty, 0.7, 1.3)
```

### 8.6 State Synchronization Rules

**Принцип:** Server является source of truth. Клиент НЕ доверяет своему локальному состоянию.

```
При choose():
1. Клиент отправляет: { playerId, choiceId, expectedVersion }
2. Сервер выполняет:
   a. Валидирует choiceId принадлежит currentEvent игрока
   b. Применяет effects
   c. Генерирует nextEvent
   d. Сохраняет в БД
3. Сервер возвращает: { updatedPlayer, nextEvent, newVersion, checkResult }
4. Клиент MUST перезаписать своё состояние из ответа сервера

При ошибке CONCURRENT_MODIFICATION:
1. Сервер возвращает текущее состояние игрока
2. Клиент MUST перезагрузить состояние из ответа
3. Клиент повторяет запрос с новой версией

При network timeout:
1. Клиент определяет timeout как возможную успешную операцию
2. Клиент вызывает loadGame() для проверки статуса
3. Если turn incremented → применить ответ, показать результат
4. Если turn не изменился → retry choose()
```

---

## 9. API / Data Contracts

### 9.1 GraphQL Schema

```graphql
scalar Time
scalar JSON

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

type EventTemplate {
  id: ID!
  type: EventType!
  tags: [String!]!
  context: EventContext!
}

type EventContext {
  time: TimeOfDay!
  location: Location!
  urgency: Urgency!
}

enum TimeOfDay {
  MORNING
  DAY
  NIGHT
}

enum Location {
  BARRACKS
  TRAINING_GROUND
  CAFETERIA
  GUARD_DUTY
  STORAGE
}

enum Urgency {
  LOW
  MEDIUM
  HIGH
}

enum EventType {
  ROUTINE
  SOCIAL
  INSPECTION
  INFORMAL
  EMERGENCY
  SAFE
}

type EventInstance {
  id: ID!
  templateId: ID!
  description: String!
  resolvedVariables: JSON!
  choices: [Choice!]!
  context: EventContext!
}

type Choice {
  id: ID!
  text: String!
  available: Boolean!
}

type GameLogEntry {
  id: ID!
  playerId: ID!
  turn: Int!
  eventTemplateId: ID!
  eventDescription: String!
  choiceId: ID!
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
  stat: StatName!
  delta: Int!
  previousValue: Int!
  newValue: Int!
}

enum StatName {
  STR
  END
  AGI
  MOR
  DISC
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
  title: String!           # "Тихий дембель", "Уважаемый дембель", "Сломанный дембель"
  subtitle: String          # "Приспособленец", null, "Косячник"
  description: String!      # Полное описание финала (см. раздел 8.4)
  finalStats: PlayerStats!
  achievedOnTurn: Int!
}

enum FinalType {
  ТИХИЙ_ДЕМБЕЛЬ    # Условие: turn ≥ 27, MOR > 0, DISC ∈ [-20,+20]
  УВАЖАЕМЫЙ_ДЕМБЕЛЬ # Условие: turn ≥ 27, MOR > 0, DISC < -50, informalStatus ≥ "Дед"
  СЛОМАННЫЙ_ДЕМБЕЛЬ # Условие: MOR = 0 (в любой момент)
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

### 9.2 Запросы и мутации

#### Query: player

```graphql
query GetPlayer($id: ID!) {
  player(id: $id) {
    id
    stats { str end agi mor disc }
    formalRank
    informalStatus
    turn
    flags
  }
}
```

**Response:**
```json
{
  "data": {
    "player": {
      "id": "abc123",
      "stats": { "str": 52, "end": 48, "agi": 55, "mor": 45, "disc": -10 },
      "formalRank": "РЯДОВОЙ",
      "informalStatus": "ЗАПАХ",
      "turn": 5,
      "flags": ["risk_report"]
    }
  }
}
```

#### Query: currentEvent

```graphql
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
```

#### Query: eventHistory

```graphql
query GetEventHistory($playerId: ID!, $limit: Int = 10) {
  eventHistory(playerId: $playerId, limit: $limit) {
    id
    turn
    eventDescription
    choiceText
    checkResult { success outcome }
    effects { stat delta }
  }
}
```

#### Mutation: startGame

```graphql
mutation StartGame {
  startGame {
    gameId
    player { id stats { str end agi mor disc } formalRank informalStatus turn }
    currentEvent { id description choices { id text available } }
    isGameOver
  }
}
```

**Response:**
```json
{
  "data": {
    "startGame": {
      "gameId": "new-game-uuid",
      "player": {
        "id": "player-uuid",
        "stats": { "str": 50, "end": 50, "agi": 50, "mor": 50, "disc": 0 },
        "formalRank": "РЯДОВОЙ",
        "informalStatus": "ЗАПАХ",
        "turn": 1
      },
      "currentEvent": {
        "id": "event-instance-uuid",
        "description": "Внезапно сержант закричал: 'Смирно! Рота, построение!'. Чёрт, я не успел забрать берцы из сушилки.",
        "choices": [
          { "id": "choice_001", "text": "Напрячь салагу", "available": true },
          { "id": "choice_002", "text": "Быстро сгонять самому", "available": true },
          { "id": "choice_003", "text": "Поискать под кроватью", "available": false },
          { "id": "choice_004", "text": "Выйти в тапочках", "available": false }
        ]
      },
      "isGameOver": false
    }
  }
}
```

#### Mutation: choose

```graphql
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
      stats { str end agi mor disc }
      formalRank
      informalStatus
      turn
      version
    }
    nextEvent {
      id
      description
      choices { id text available }
    }
    gameOver
    final {
      type
      title
      description
      finalStats { str end agi mor disc }
      achievedOnTurn
    }
    newVersion
  }
}
```

**Response:**
```json
{
  "data": {
    "choose": {
      "success": true,
      "checkResult": {
        "success": true,
        "outcome": "SUCCESS",
        "description": "Салага молча приносит берцы."
      },
      "effects": [
        { "stat": "DISC", "delta": -8, "previousValue": 0, "newValue": -8 },
        { "stat": "MOR", "delta": -2, "previousValue": 50, "newValue": 48 }
      ],
      "updatedPlayer": {
        "id": "player-uuid",
        "stats": { "str": 50, "end": 50, "agi": 50, "mor": 48, "disc": -8 },
        "formalRank": "РЯДОВОЙ",
        "informalStatus": "ЗАПАХ",
        "turn": 2
      },
      "nextEvent": {
        "id": "event-instance-uuid-2",
        "description": "Командир ночного наряда сержант обнаружил, что ты заснул на посту.",
        "choices": [
          { "id": "choice_001", "text": "Притвориться больным", "available": true },
          { "id": "choice_002", "text": "Поправить на месте", "available": false }
        ]
      },
      "gameOver": false,
      "final": null
    }
  }
}
```

### 9.3 Обработка ошибок GraphQL

```json
{
  "errors": [
    {
      "message": "Choice is not available",
      "extensions": {
        "code": "CHOICE_UNAVAILABLE",
        "choiceId": "choice_003",
        "playerId": "player-uuid"
      }
    }
  ],
  "data": null
}
```

**Error codes:**
| Code | HTTP Status | Description |
|------|-------------|-------------|
| PLAYER_NOT_FOUND | 404 | Игрок с указанным ID не найден |
| CHOICE_UNAVAILABLE | 400 | Выбор недоступен по условиям availability |
| CHOICE_NOT_IN_CURRENT_EVENT | 400 | Выбор не принадлежит текущему событию |
| GAME_ALREADY_FINISHED | 400 | Игра уже завершена, новые ходы невозможны |
| INVALID_STAT_VALUE | 400 | Некорректное значение стата |
| CONCURRENT_MODIFICATION | 409 | Версия игрока изменилась, повторите запрос |
| VERSION_MISMATCH | 409 | Ожидаемая версия не совпадает с текущей |
| TEMPLATE_NOT_FOUND | 500 | Шаблон события не найден (internal error) |

---

## 10. Data Flow

### 10.1 Архитектура данных

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │   Pinia      │    │   Vue 3      │    │   Apollo     │  │
│  │   Store      │◄───│   Components  │───▶│   Client     │  │
│  │  (game.ts)   │    │              │    │              │  │
│  └──────┬───────┘    └──────────────┘    └──────┬───────┘  │
│         │                                        │          │
│         │ version: 1                             │          │
└─────────┼────────────────────────────────────────┼──────────┘
          │ GraphQL + version                     │
          ▼                                        ▼
┌─────────────────────────────────────────────────────────────┐
│                        Backend                               │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              GraphQL Handler (gqlgen)                   │ │
│  └─────────────────────────────────────────────────────────┘ │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              Application Layer (Services)               │ │
│  │  ┌────────────┐ ┌────────────┐ ┌────────────┐          │ │
│  │  │  GameService│ │PlayerService│ │EventService│          │ │
│  │  └────────────┘ └────────────┘ └────────────┘          │ │
│  └─────────────────────────────────────────────────────────┘ │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              Domain Layer (Business Logic)               │ │
│  │  ┌────────────┐ ┌────────────┐ ┌────────────┐          │ │
│  │  │  CheckEngine│ │EffectEngine│ │Progression │          │ │
│  │  └────────────┘ └────────────┘ └────────────┘          │ │
│  └─────────────────────────────────────────────────────────┘ │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              Infrastructure Layer                       │ │
│  │  ┌────────────────┐  ┌────────────────┐                  │ │
│  │  │ PlayerRepo     │  │ EventTemplateRepo│                │ │
│  │  │ (PostgreSQL)   │  │ (PostgreSQL)    │                  │ │
│  │  └────────────────┘  └────────────────┘                  │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    PostgreSQL                                 │
│  ┌──────────────┐  ┌──────────────┐                        │
│  │   players    │  │   game_logs  │                        │
│  │  (stats as   │  │  (immutable) │                        │
│  │   columns)   │  │              │                        │
│  └──────────────┘  └──────────────┘                        │
│  ┌──────────────┐                                          │
│  │event_templates│                                         │
│  │  (readonly)   │                                         │
│  └──────────────┘                                          │
└─────────────────────────────────────────────────────────────┘
```

### 10.2 Схема базы данных

```sql
-- ============================================================================
-- PLAYERS — Единственный источник правды о состоянии игрока
-- ============================================================================
CREATE TABLE players (
    -- Идентификация
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Статы с ограничениями (НЕ JSONB)
    str INT NOT NULL DEFAULT 50 CHECK (str BETWEEN 1 AND 100),
    end_ INT NOT NULL DEFAULT 50 CHECK (end_ BETWEEN 1 AND 100),
    agi INT NOT NULL DEFAULT 50 CHECK (agi BETWEEN 1 AND 100),
    mor INT NOT NULL DEFAULT 50 CHECK (mor BETWEEN 0 AND 100),
    disc INT NOT NULL DEFAULT 0 CHECK (disc BETWEEN -100 AND 100),
    
    -- Прогрессия
    formal_rank VARCHAR(50) NOT NULL DEFAULT 'РЯДОВОЙ',
    informal_status VARCHAR(50) NOT NULL DEFAULT 'ЗАПАХ',
    
    -- Состояние игры
    turn INT NOT NULL DEFAULT 1 CHECK (turn BETWEEN 1 AND 30),
    flags JSONB DEFAULT '[]',
    
    -- Версионирование для optimistic locking
    version INT NOT NULL DEFAULT 1,
    
    -- Метаданные
    is_finished BOOLEAN NOT NULL DEFAULT FALSE,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для players
CREATE INDEX idx_players_updated_at ON players(updated_at DESC);
CREATE INDEX idx_players_turn ON players(turn);
CREATE INDEX idx_players_finished ON players(is_finished) WHERE is_finished = FALSE;
CREATE INDEX idx_players_finished_at ON players(is_finished, finished_at) 
    WHERE is_finished = TRUE;

-- ============================================================================
-- GAME_LOGS — Неизменяемая история всех выборов
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
    
    -- Уникальный constraint: один лог на игрока на ход
    CONSTRAINT uq_game_logs_player_turn UNIQUE (player_id, turn)
);

-- Индексы для game_logs
CREATE INDEX idx_game_logs_player_id ON game_logs(player_id);
CREATE INDEX idx_game_logs_player_turn ON game_logs(player_id, turn DESC);
CREATE INDEX idx_game_logs_created_at ON game_logs(created_at);

-- ============================================================================
-- EVENT_TEMPLATES — Контент игры (только для чтения)
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

-- Индексы для event_templates
CREATE INDEX idx_event_templates_type ON event_templates(type);
CREATE INDEX idx_event_templates_type_active ON event_templates(type) 
    WHERE is_active = TRUE;
CREATE INDEX idx_event_templates_used_count ON event_templates(used_count);
```

### 10.3 Алгоритмы работы с данными

#### 10.3.1 Автосохранение (после каждого choose)

```
1. BEGIN TRANSACTION
2. SELECT player FROM players WHERE id = $playerId FOR UPDATE
3. Проверить version == expectedVersion
   → Если нет: ROLLBACK, вернуть CONCURRENT_MODIFICATION
4. Применить effects: UPDATE stats с CHECK constraints
5. UPDATE players SET 
     turn = turn + 1,
     flags = flags || new_flags,
     version = version + 1,
     updated_at = NOW(),
     is_finished = CASE WHEN финал THEN TRUE ELSE FALSE END,
     finished_at = CASE WHEN финал THEN NOW() ELSE finished_at END
   WHERE id = $playerId AND version = $expectedVersion
6. Проверить rows_affected == 1
   → Если нет: ROLLBACK, вернуть CONCURRENT_MODIFICATION
7. INSERT INTO game_logs (turn, event, choice, result, effects)
8. COMMIT
9. Сгенерировать следующее EventInstance из шаблона
10. RETURN updatedPlayer + nextEvent + newVersion
```

#### 10.3.2 Загрузка игры (loadGame)

```
1. SELECT * FROM players WHERE id = $playerId
   → Если не найден: вернуть PLAYER_NOT_FOUND
2. SELECT * FROM game_logs 
   WHERE player_id = $playerId 
   ORDER BY turn DESC 
   LIMIT 10
3. Сгенерировать currentEvent из последнего template_id
4. RETURN { player, currentEvent, history }
```

#### 10.3.3 Генерация EventInstance (dynamic, not stored)

```
1. SELECT * FROM event_templates WHERE id = $templateId
2. Подставить переменные из template.variables
   → Для каждой variable: random.choice(variable.values)
3. Фильтровать choices по availability
   → Проверить stats conditions (min/max)
   → Проверить flags conditions
4. Сгенерировать уникальный instance_id (UUID)
5. RETURN EventInstance (never stored in DB)
```

### 10.3 Структура EventTemplate (JSON)

```typescript
interface EventTemplate {
  id: string;
  type: 'ROUTINE' | 'SOCIAL' | 'INSPECTION' | 'INFORMAL' | 'EMERGENCY' | 'SAFE';
  tags: string[];
  context: {
    time: 'MORNING' | 'DAY' | 'NIGHT';
    location: 'BARRACKS' | 'TRAINING_GROUND' | 'CAFETERIA' | 'GUARD_DUTY' | 'STORAGE';
    urgency: 'LOW' | 'MEDIUM' | 'HIGH';
  };
  template: {
    description: string; // e.g. "Внезапно {{npc_role}} закричал: '{{command}}'."
    variables: Record<string, string[]>; // e.g. { "npc_role": ["сержант", "офицер"] }
  };
  choices: Choice[];
}

interface Choice {
  id: string;
  text: string;
  availability?: AvailabilityCondition;
  check: CheckConfig;
  outcomes: {
    success?: Outcome;
    partial?: Outcome;
    failure?: Outcome;
    ignored?: Outcome;
    noticed_success?: Outcome;
    noticed_failure?: Outcome;
  };
  nextEvent?: string; // templateId for chained events
}

interface AvailabilityCondition {
  stats?: Partial<Record<StatName, { min?: number; max?: number }>>;
  flags?: { required?: string[]; forbidden?: string[] };
}

interface CheckConfig {
  type: 'THRESHOLD' | 'PROBABILITY' | 'CATASTROPHIC';
  formula?: string; // e.g. "STR + (DISC * -0.3) + (MOR * 0.2)"
  chanceFormula?: string; // for probability
  threshold?: number;
  noticeChanceFormula?: string; // for catastrophic
  powerCheck?: {
    formula: string;
    threshold: number;
  };
}

interface Outcome {
  text: string;
  effects?: Effect[];
  flags?: string[];
  nextEvent?: string;
}

interface Effect {
  stat: StatName;
  delta: number;
}
```

---

## 10.2 Frontend Архитектура (Feature-Sliced Design)

### 10.2.1 Обзор

Frontend реализован на Vue 3 с использованием **Feature-Sliced Design (FSD)** методологии — модульной архитектуры, которая разделяет код на функциональные слои.

**Технологический стек:**
- Vue 3 + Composition API (`<script setup>`)
- TypeScript (strict mode)
- Pinia (state management)
- Apollo Client (GraphQL)
- Vite (bundler)

### 10.2.2 Принципы FSD

| Принцип | Описание |
|---------|----------|
| Слои | app → pages → widgets → features → entities → shared |
| Изоляция | Каждый слой зависит только от нижележащих |
| Инкапсуляция | Features экспортируют только публичный API |
| Переиспользуемость | shared содержит только универсальный код |

### 10.2.3 Структура папок

```
frontend/
├── src/
│   ├── app/                    # Корневая конфигурация
│   │   ├── App.vue             # Корневой компонент
│   │   ├── router.ts           # Vue Router
│   │   └── index.ts            # Инициализация приложения
│   │
│   ├── pages/                  # Страницы (роуты)
│   │   ├── StartPage/          # Главный экран
│   │   │   ├── StartPage.vue
│   │   │   └── index.ts
│   │   ├── GamePage/           # Игровой экран
│   │   │   ├── GamePage.vue
│   │   │   └── index.ts
│   │   └── FinalPage/          # Экран финала
│   │       ├── FinalPage.vue
│   │       └── index.ts
│   │
│   ├── widgets/                # Композитные компоненты
│   │   ├── GameScreen/         # Основной игровой контейнер
│   │   │   ├── GameScreen.vue
│   │   │   └── components/
│   │   │       ├── TopBar.vue         # Верхняя панель
│   │   │       ├── StatsRow.vue       # Строка статов
│   │   │       ├── EventCard.vue      # Карточка события
│   │   │       └── ChoiceButtons.vue  # Кнопки выбора
│   │   ├── StatsPanel/         # Модальная панель статов
│   │   │   ├── StatsPanel.vue
│   │   │   └── index.ts
│   │   ├── HistoryPanel/       # Панель истории
│   │   │   ├── HistoryPanel.vue
│   │   │   └── index.ts
│   │   └── LoadingState/       # Состояние загрузки
│   │       └── LoadingState.vue
│   │
│   ├── features/               # Бизнес-логика
│   │   ├── start-game/         # Начало игры
│   │   │   ├── components/
│   │   │   │   └── StartButton.vue
│   │   │   ├── useStartGame.ts  # Composable
│   │   │   └── index.ts
│   │   ├── make-choice/        # Обработка выбора
│   │   │   ├── components/
│   │   │   │   └── ChoiceButton.vue
│   │   │   ├── useMakeChoice.ts
│   │   │   └── index.ts
│   │   ├── view-stats/         # Отображение статов
│   │   │   ├── components/
│   │   │   │   ├── StatCard.vue
│   │   │   │   └── ProgressBar.vue
│   │   │   ├── useStats.ts
│   │   │   └── index.ts
│   │   ├── view-history/       # Просмотр истории
│   │   │   ├── components/
│   │   │   │   └── HistoryItem.vue
│   │   │   ├── useHistory.ts
│   │   │   └── index.ts
│   │   └── game-progress/     # Прогрессия игры
│   │       ├── useGameProgress.ts
│   │       └── index.ts
│   │
│   ├── entities/              # Доменные модели
│   │   ├── Player/            # Модель игрока
│   │   │   ├── model.ts       # Интерфейсы Player, PlayerStats
│   │   │   └── index.ts
│   │   ├── Event/             # Модель события
│   │   │   ├── model.ts       # EventTemplate, EventInstance, Choice
│   │   │   └── index.ts
│   │   ├── Game/              # Модель игры
│   │   │   ├── model.ts       # GameState, Final types
│   │   │   └── index.ts
│   │   └── Final/            # Типы финалов
│   │       ├── model.ts
│   │       └── index.ts
│   │
│   ├── shared/                # Переиспользуемый код
│   │   ├── api/               # GraphQL клиент
│   │   │   ├── client.ts     # Apollo Client настройка
│   │   │   ├── queries/      # GraphQL запросы
│   │   │   │   ├── player.ts
│   │   │   │   ├── event.ts
│   │   │   │   └── history.ts
│   │   │   └── mutations/    # GraphQL мутации
│   │   │       ├── startGame.ts
│   │   │       ├── choose.ts
│   │   │       └── restartGame.ts
│   │   ├── ui/               # Базовые UI компоненты
│   │   │   ├── Button/
│   │   │   ├── Modal/
│   │   │   ├── ProgressBar/
│   │   │   └── Toast/
│   │   ├── lib/              # Утилиты
│   │   │   ├── format.ts     # Форматирование (дат, чисел)
│   │   │   ├── constants.ts  # Константы
│   │   │   └── validation.ts
│   │   └── config/
│   │       └── index.ts
│   │
│   ├── stores/                # Pinia сторы
│   │   ├── game.ts           # Основной стор игры
│   │   └── index.ts
│   │
│   ├── assets/                # Статические ресурсы
│   │   ├── styles/
│   │   │   └── main.css
│   │   └── images/
│   │
│   └── main.ts               # Точка входа
│
├── public/
│   └── index.html
│
├── package.json
├── vite.config.ts
└── tsconfig.json
```

### 10.2.4 Зависимости между слоями

```
┌─────────────────────────────────────┐
│              app/                   │  → Маршрутизация, инициализация
└──────────────┬────────────────────┘
               │
┌──────────────▼────────────────────┐
│             pages/                  │  → Страницы (StartPage, GamePage, FinalPage)
└──────────────┬────────────────────┘
               │
┌──────────────▼────────────────────┐
│            widgets/                 │  → GameScreen, StatsPanel, HistoryPanel
└──────────────┬────────────────────┘
               │
┌──────────────▼────────────────────┐
│           features/                 │  → start-game, make-choice, view-stats
└──────────────┬────────────────────┘
               │
┌──────────────▼────────────────────┐
│           entities/                 │  → Player, Event, Game, Final
└──────────────┬────────────────────┘
               │
┌──────────────▼────────────────────┐
│            shared/                 │  → api, ui, lib, config
└─────────────────────────────────────┘
```

### 10.2.5 Pinia Store (game.ts)

```typescript
// stores/game.ts
interface GameState {
  player: Player | null;
  currentEvent: EventInstance | null;
  eventHistory: GameLogEntry[];
  isLoading: boolean;
  error: string | null;
  isGameOver: boolean;
  final: Final | null;
  version: number;
}

interface GameActions {
  startGame(): Promise<void>;
  makeChoice(choiceId: string): Promise<void>;
  loadGame(): Promise<void>;
  restartGame(): Promise<void>;
}
```

### 10.2.6 Composables

| Composable | Слой | Ответственность |
|------------|-------|----------------|
| `useStartGame` | features/start-game | Начало новой игры |
| `useMakeChoice` | features/make-choice | Обработка выбора игрока |
| `useStats` | features/view-stats | Форматирование и отображение статов |
| `useHistory` | features/view-history | Управление историей событий |
| `useGameProgress` | features/game-progress | Отслеживание прогресса |

### 10.2.7 GraphQL Integration

```typescript
// shared/api/client.ts
import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core';

const httpLink = createHttpLink({
  uri: import.meta.env.VITE_GRAPHQL_URL || 'http://localhost:4000/graphql',
});

export const apolloClient = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache(),
  defaultOptions: {
    watchQuery: { fetchPolicy: 'network-only' },
    query: { fetchPolicy: 'network-only' },
  },
});
```

### 10.2.8 GraphQL Queries/Mutations

| Операция | Тип | Назначение |
|----------|-----|------------|
| `GetPlayer` | Query | Получить данные игрока |
| `GetCurrentEvent` | Query | Получить текущее событие |
| `GetEventHistory` | Query | Получить историю событий |
| `StartGame` | Mutation | Начать новую игру |
| `Choose` | Mutation | Сделать выбор |
| `RestartGame` | Mutation | Начать заново |

### 10.2.9 Компоненты и их использование

**Pages (используют Widgets):**
- `StartPage` → `StartButton` (feature)
- `GamePage` → `GameScreen` (widget)
- `FinalPage` → (inline)

**Widgets (используют Features + Entities):**
- `GameScreen` → `TopBar`, `StatsRow`, `EventCard`, `ChoiceButtons` + `useMakeChoice`
- `StatsPanel` → `StatCard`, `ProgressBar` + `useStats`
- `HistoryPanel` → `HistoryItem` + `useHistory`

**Features (используют Entities + Shared):**
- `useStartGame` → `Player` entity + `startGame` mutation
- `useMakeChoice` → `Event`, `Player` entities + `choose` mutation

---

## 10.3 Backend Архитектура (Go + GraphQL)

### 10.3.1 Обзор

Backend реализован на **Go** с использованием **GraphQL** (gqlgen) для API. Архитектура следует принципам **DDD** (Domain-Driven Design) с чётким разделением на слои.

**Технологический стек:**
- Go 1.21+
- GraphQL (gqlgen)
- PostgreSQL 15+ (pgx driver)
- Redis 7+ (go-redis)
- Docker

### 10.3.2 Структура проекта

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа, запуск сервера
│
├── internal/                    # Приватные пакеты (enforced by Go)
│   ├── domain/                 # Доменный слой
│   │   ├── entities/          # Сущности
│   │   │   ├── player.go      # Player entity
│   │   │   ├── event.go       # Event entity
│   │   │   ├── game.go        # Game state
│   │   │   └── final.go       # Final types
│   │   ├── valueobjects/       # Value objects
│   │   │   ├── stats.go       # PlayerStats
│   │   │   └── rank.go        # Formal/Informal rank
│   │   └── errors/            # Доменные ошибки
│   │       ├── errors.go      # Sentinel errors
│   │       └── game.go        # Game-specific errors
│   │
│   ├── service/                # Бизнес-логика
│   │   ├── player_service.go  # Управление игроками
│   │   ├── game_service.go    # Игровая логика
│   │   ├── event_service.go   # Генерация событий
│   │   ├── progression_service.go # Прогрессия, финалы
│   │   └── analytics_service.go  # Аналитика
│   │
│   ├── handler/               # GraphQL handlers (gqlgen)
│   │   ├── resolver.go       # Resolver struct
│   │   ├── query.go          # Query resolvers
│   │   ├── mutation.go        # Mutation resolvers
│   │   └── generated/         # autogenerated by gqlgen
│   │
│   ├── repository/            # Data access layer (interfaces)
│   │   ├── player_repo.go    # Player repository interface
│   │   ├── event_repo.go     # Event template interface
│   │   └── game_log_repo.go  # Game log interface
│   │
│   └── middleware/            # HTTP/GraphQL middleware
│       ├── logging.go        # Structured logging
│       ├── recovery.go       # Panic recovery
│       └── cors.go          # CORS
│
├── pkg/                      # Публичные пакеты
│   ├── graphql/              # GraphQL схемы и типы
│   │   ├── schema.graphql   # GraphQL schema
│   │   └── models/           # Generated models
│   ├── config/               # Конфигурация
│   │   └── config.go
│   └── logger/              # Логирование
│       └── logger.go
│
├── infrastructure/            # Инфраструктура
│   ├── database/            # PostgreSQL
│   │   ├── pool.go          # Connection pool
│   │   ├── migrations/     # SQL migrations
│   │   └── migrations.go  # Migration runner
│   ├── cache/              # Redis
│   │   └── cache.go
│   └── repositories/       # Repository implementations
│       ├── postgres_player_repo.go
│       ├── postgres_event_repo.go
│       └── postgres_game_log_repo.go
│
├── migrations/              # SQL миграции
│   ├── 001_create_tables.sql
│   └── ...
│
├── go.mod
├── go.sum
└── generate.go             # go generate для gqlgen
```

### 10.3.3 Зависимости между слоями

```
┌─────────────────────────────────────────────────┐
│                   cmd/server                     │  → main.go, запуск
└─────────────────────┬───────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────┐
│                 handler/                         │  → GraphQL resolvers
│                 (incoming requests)              │
└─────────────────────┬───────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────┐
│                 service/                         │  → Business logic
│  player_service, game_service, event_service    │
└─────────────────────┬───────────────────────────┘
                      │
        ┌─────────────┴─────────────┐
        │                           │
┌───────▼───────┐         ┌────────▼────────┐
│  repository/  │         │   domain/        │
│ (interfaces)  │         │   (entities)     │
└───────┬───────┘         └──────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────┐
│           infrastructure/repositories/            │  → PostgreSQL, Redis
│           (implementations)                     │
└─────────────────────────────────────────────────┘
```

### 10.3.4 Domain Entities

```go
// internal/domain/entities/player.go
type Player struct {
    ID             uuid.UUID
    Stats          PlayerStats
    FormalRank     FormalRank
    InformalStatus InformalStatus
    Turn           int
    Flags          []string
    Version        int
    IsFinished     bool
    FinishedAt     *time.Time
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

type PlayerStats struct {
    STR  int // [1, 100]
    END  int // [1, 100]
    AGI  int // [1, 100]
    MOR  int // [0, 100]
    DISC int // [-100, +100]
}

type FormalRank string

const (
    RankPrivate    FormalRank = "РЯДОВОЙ"
    RankCorporal   FormalRank = "ЕФРЕЙТОР"
    RankJuniorSergeant FormalRank = "МЛ_СЕРЖАНТ"
    RankSergeant   FormalRank = "СЕРЖАНТ"
)

type InformalStatus string

const (
    StatusSmell    InformalStatus = "ЗАПАХ"
    StatusSpirit   InformalStatus = "ДУХ"
    StatusElephant InformalStatus = "СЛОН"
    StatusLadler   InformalStatus = "ЧЕРПАК"
    StatusGrandfather InformalStatus = "ДЕД"
    StatusDembel  InformalStatus = "ДЕМБЕЛЬ"
)
```

### 10.3.5 Service Layer

```go
// internal/service/game_service.go
type GameService interface {
    StartGame(ctx context.Context) (*GameState, error)
    MakeChoice(ctx context.Context, playerID uuid.UUID, choiceID string, expectedVersion int) (*ChooseResult, error)
    LoadGame(ctx context.Context, playerID uuid.UUID) (*GameState, error)
    RestartGame(ctx context.Context, playerID uuid.UUID) (*GameState, error)
}

type gameService struct {
    playerRepo   repository.PlayerRepository
    eventRepo   repository.EventRepository
    gameLogRepo repository.GameLogRepository
    logger      *logger.Logger
}

func (s *gameService) MakeChoice(ctx context.Context, playerID uuid.UUID, choiceID string, expectedVersion int) (*ChooseResult, error) {
    // 1. Получить игрока с lock
    player, err := s.playerRepo.GetForUpdate(ctx, playerID)
    if err != nil {
        return nil, fmt.Errorf("get player: %w", err)
    }
    
    // 2. Проверить версию (optimistic locking)
    if player.Version != expectedVersion {
        return nil, errors.New("VERSION_MISMATCH")
    }
    
    // 3. Обработать выбор
    result, err := s.processChoice(ctx, player, choiceID)
    if err != nil {
        return nil, fmt.Errorf("process choice: %w", err)
    }
    
    // 4. Сохранить изменения
    if err := s.playerRepo.Update(ctx, player); err != nil {
        return nil, fmt.Errorf("update player: %w", err)
    }
    
    // 5. Записать в лог
    if err := s.gameLogRepo.Create(ctx, logEntry); err != nil {
        return nil, fmt.Errorf("create log: %w", err)
    }
    
    return result, nil
}
```

### 10.3.6 Repository Pattern

```go
// internal/repository/player_repo.go
type PlayerRepository interface {
    Create(ctx context.Context, player *entities.Player) error
    GetByID(ctx context.Context, id uuid.UUID) (*entities.Player, error)
    GetForUpdate(ctx context.Context, id uuid.UUID) (*entities.Player, error) // SELECT FOR UPDATE
    Update(ctx context.Context, player *entities.Player) error
    Delete(ctx context.Context, id uuid.UUID) error
}

// internal/infrastructure/repositories/postgres_player_repo.go
type postgresPlayerRepository struct {
    db *pgxpool.Pool
}

func (r *postgresPlayerRepository) GetForUpdate(ctx context.Context, id uuid.UUID) (*entities.Player, error) {
    row := r.db.QueryRow(ctx, `
        SELECT id, str, end_, agi, mor, disc, formal_rank, informal_status, 
               turn, flags, version, is_finished, finished_at, created_at, updated_at
        FROM players 
        WHERE id = $1 
        FOR UPDATE`, id)
    
    var player entities.Player
    err := row.Scan(...)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrPlayerNotFound
        }
        return nil, fmt.Errorf("scan player: %w", err)
    }
    return &player, nil
}
```

### 10.3.7 GraphQL Handler (gqlgen)

```go
// internal/handler/mutation.go
type MutationResolver struct {
    gameService   service.GameService
    playerService service.PlayerService
}

func (r *MutationResolver) StartGame(ctx context.Context) (*GameState, error) {
    state, err := r.gameService.StartGame(ctx)
    if err != nil {
        return nil, handleError(err)
    }
    return toGameState(state), nil
}

func (r *MutationResolver) Choose(ctx context.Context, playerID string, choiceID string, expectedVersion int) (*ChooseResult, error) {
    pid, err := uuid.Parse(playerID)
    if err != nil {
        return nil, gqlerror.Errorf("invalid player id")
    }
    
    result, err := r.gameService.MakeChoice(ctx, pid, choiceID, expectedVersion)
    if err != nil {
        return nil, handleError(err)
    }
    return toChooseResult(result), nil
}
```

### 10.3.8 Error Handling

```go
// internal/domain/errors/game.go
var (
    ErrPlayerNotFound     = errors.New("player not found")
    ErrGameAlreadyFinished = errors.New("game already finished")
    ErrChoiceUnavailable  = errors.New("choice unavailable")
    ErrVersionMismatch    = errors.New("version mismatch")
    ErrEventNotFound      = errors.New("event not found")
)

type GameError struct {
    Code    string
    Message string
    Err     error
}

func (e *GameError) Error() string {
    return e.Message
}

// internal/handler/errors.go
func handleError(err error) error {
    if errors.Is(err, ErrPlayerNotFound) {
        return gqlerror.Errorf("player not found")
    }
    if errors.Is(err, ErrVersionMismatch) {
        return gqlerror.Errorf("version mismatch")
    }
    // Логируем оригинальную ошибку
    log.Error(err)
    return gqlerror.Errorf("internal error")
}
```

### 10.3.9 Конфигурация

```go
// pkg/config/config.go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Logger   LoggerConfig
}

type ServerConfig struct {
    Port string
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    PoolSize int32
}

func Load() (*Config, error) {
    return &Config{
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "4000"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvInt("DB_PORT", 5432),
            User:     getEnv("DB_USER", "army_game"),
            Password: getEnv("DB_PASSWORD", ""),
            DBName:   getEnv("DB_NAME", "army_game"),
            PoolSize: 25,
        },
    }, nil
}
```

### 10.3.10 GraphQL Schema

```graphql
# pkg/graphql/schema.graphql

type Query {
    player(id: ID!): Player
    currentEvent(playerId: ID!): EventInstance
    eventHistory(playerId: ID!, limit: Int): [GameLogEntry!]!
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
```

### 10.3.11 Database Migrations

```sql
-- migrations/001_create_tables.sql

-- Players table
CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    str INT NOT NULL DEFAULT 50 CHECK (str BETWEEN 1 AND 100),
    end_ INT NOT NULL DEFAULT 50 CHECK (end_ BETWEEN 1 AND 100),
    agi INT NOT NULL DEFAULT 50 CHECK (agi BETWEEN 1 AND 100),
    mor INT NOT NULL DEFAULT 50 CHECK (mor BETWEEN 0 AND 100),
    disc INT NOT NULL DEFAULT 0 CHECK (disc BETWEEN -100 AND 100),
    formal_rank VARCHAR(50) NOT NULL DEFAULT 'РЯДОВОЙ',
    informal_status VARCHAR(50) NOT NULL DEFAULT 'ЗАПАХ',
    turn INT NOT NULL DEFAULT 1 CHECK (turn BETWEEN 1 AND 30),
    flags JSONB DEFAULT '[]',
    version INT NOT NULL DEFAULT 1,
    is_finished BOOLEAN NOT NULL DEFAULT FALSE,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Game logs table
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

-- Event templates table
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
```

---

## 11. Edge Cases

### 11.1 Пустые и граничные значения

| Edge Case | Ожидаемое поведение |
|------------|---------------------|
| MOR = 0 | Немедленный финал «Сломанный дембель» (Косячник) — описание: "Служба в армии принесла вам лишь горе, затрещины и сломанный нос..." |
| MOR > 100 после эффекта | Clamp к 100 |
| MOR < 0 после эффекта | Clamp к 0 (→ финал) |
| DISC < -100 после эффекта | Clamp к -100 |
| DISC > +100 после эффекта | Clamp к +100 |
| turn > 30 | Не допускается, игра завершается на turn = 30 |
| turn = 0 | Не допускается, минимум 1 |
| Все choices unavailable | Генерировать новое событие (без decrement turn) |
| Нет доступных событий в пуле | Использовать любой template (сбросить recent_history) |
| Player не найден | Вернуть ошибку PLAYER_NOT_FOUND |
| choiceId не принадлежит событию | Вернуть ошибку CHOICE_NOT_IN_CURRENT_EVENT |
| MOR delta > 10 за событие | Log warning, cap на +10 |
| MOR delta < -5 за событие | Log warning, cap на -5 |

### 11.2 Специальные игровые сценарии

| Edge Case | Ожидаемое поведение |
|------------|---------------------|
| **Perfect run** (turn = 30, MOR > 50) | Если DISC ∈ [+21,+50] и formalRank ≥ "Ефрейтор" → «Дисциплинированный дембель» (SPECIAL_ENDING_1) |
| **DISC exploitation** (только формальная линия) | Inspection events становятся +50% сложнее при formalRank ≥ "Сержант" |
| **3 MOR- подряд** | Следующее событие FORCED SAFE type (не выбирается, а гарантируется) |
| **Flags overflow** (> 20 flags) | Truncate oldest flags (FIFO) |
| **Negative DISC spiral** (< -80) | Informal events дают diminishing returns |

### 11.3 Состояние гонки (Race Conditions)

| Edge Case | Решение |
|-----------|---------|
| Двойной choose (быстрые клики) | Optimistic locking через version field. Клиент отправляет expectedVersion. Сервер проверяет совпадение. |
| Concurrent modification detected | Вернуть CONCURRENT_MODIFICATION. Клиент перезагружает состояние и повторяет запрос с новой версией. |
| Конкурентные запросы loadGame | Нет проблемы — SELECT без блокировки, данные eventually consistent. |
| Race condition в used_count шаблонов | Атомарное UPDATE: `SET used_count = used_count + 1` — PostgreSQL гарантирует атомарность. |

### 11.3 Сетевые ошибки

| Edge Case | Решение |
|-----------|---------|
| Таймаут запроса (> 5s) | Retry до 3 раз с exponential backoff |
| Сервер недоступен | UI: «Ошибка соединения. Попробуйте позже.» |
| GraphQL ошибка | Отобразить user-friendly message |
| Потеря связи во время choose | Сохранить pending state в localStorage, retry при восстановлении |

### 11.4 Повреждение данных

| Edge Case | Решение |
|-----------|---------|
| Невалидный JSON в effects | Откатить транзакцию, вернуть ошибку |
| Некорректный stat name | Игнорировать эффект, залогировать warning |
| Отсутствует обязательное поле | Вернуть ошибку валидации |

---

## 12. Error Handling

### 12.1 Иерархия ошибок (Backend)

```go
// pkg/errors/errors.go
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Cause   error  `json:"-"`
}

func (e *AppError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Конкретные ошибки
var (
    ErrPlayerNotFound = &AppError{Code: "PLAYER_NOT_FOUND", Message: "Player not found"}
    ErrChoiceUnavailable = &AppError{Code: "CHOICE_UNAVAILABLE", Message: "Choice is not available"}
    ErrGameFinished = &AppError{Code: "GAME_ALREADY_FINISHED", Message: "Game is already finished"}
    ErrInvalidStat = &AppError{Code: "INVALID_STAT_VALUE", Message: "Invalid stat value"}
    ErrConcurrentModification = &AppError{Code: "CONCURRENT_MODIFICATION", Message: "Player state was modified by another request. Please retry."}
    ErrVersionMismatch = &AppError{Code: "VERSION_MISMATCH", Message: "Expected version does not match current version."}
    ErrTemplateNotFound = &AppError{Code: "TEMPLATE_NOT_FOUND", Message: "Event template not found"}
)
```

### 12.2 Обработка ошибок (Frontend)

```typescript
// utils/errorHandler.ts
interface GraphQLError {
  message: string;
  extensions?: {
    code: string;
    [key: string]: unknown;
  };
}

function handleGraphQLError(error: GraphQLError): string {
  const errorMessages: Record<string, string> = {
    PLAYER_NOT_FOUND: 'Игра не найдена. Начните новую игру.',
    CHOICE_UNAVAILABLE: 'Этот выбор сейчас недоступен.',
    GAME_ALREADY_FINISHED: 'Игра уже завершена. Начните новую игру.',
    CONCURRENT_MODIFICATION: 'Игра обновилась. Перезагружаю...',
    VERSION_MISMATCH: 'Версия игры не совпадает. Перезагружаю...',
    NETWORK_ERROR: 'Ошибка сети. Проверьте подключение.',
    INTERNAL_ERROR: 'Внутренняя ошибка сервера. Попробуйте позже.',
  };

  const code = error.extensions?.code ?? 'INTERNAL_ERROR';
  return errorMessages[code] ?? errorMessages.INTERNAL_ERROR;
}

// Пример обработки concurrent modification в store
async function makeChoice(choiceId: string) {
  try {
    const result = await chooseMutation({
      playerId: player.value.id,
      choiceId: choiceId,
      expectedVersion: player.value.version
    });
    // Обновить локальное состояние
    player.value = result.updatedPlayer;
    // Сохранить новую версию
    localStorage.setItem('gameVersion', result.newVersion.toString());
  } catch (error) {
    if (error.extensions?.code === 'CONCURRENT_MODIFICATION') {
      // Перезагрузить состояние игры
      await loadGame(player.value.id);
      // Повторить выбор с новой версией
      await makeChoice(choiceId);
    } else {
      throw error;
    }
  }
}
```
```

### 12.3 Retry Policy

| Ситуация | Policy |
|----------|--------|
| Network timeout | Retry 3x, delay: 1s, 2s, 4s (exponential) |
| 5xx server error | Retry 3x, delay: 2s, 4s, 8s |
| 4xx client error | No retry, show error message |
| CONCURRENT_MODIFICATION | Reload game state, retry once with new version |
| choose mutation | Retry with updated expectedVersion |

---

## 13. Permissions & Access Control

### 13.1 Ролевая модель (MVP)

| Роль | Capabilities |
|------|--------------|
| guest | startGame, choose, loadGame (своя), restartGame |

### 13.2 Ресурсные ограничения

| Ресурс | Ограничение | Enforcement |
|--------|-------------|-------------|
| max games per browser | 1 active + 1 history | localStorage gameId |
| max event history | 30 entries | Truncate on save |
| max flags per player | 20 flags | Cap at limit |
| request rate | 100 req/min per IP | Middleware rate limiter |

### 13.3 Запрещённые действия

| Действие | Причина |
|----------|---------|
| Загрузка чужого сохранения | Нет авторизации в MVP |
| Изменение статов напрямую | Только через choose |
| Доступ к admin queries | Не реализовано в MVP |

---

## 14. UI/UX Requirements

### 14.1 Экраны

#### Screen 1: Главный экран (Home)

```
┌────────────────────────────────────────────────────────────┐
│                        АРМЕЙКА                              │
│              Narrative Survival Game                        │
├────────────────────────────────────────────────────────────┤
│                                                            │
│    [Найти берцы] — начать новую игру                       │
│                                                            │
│    ───────────────────────────────────────                 │
│                                                            │
│    Если есть сохранение:                                   │
│    [Продолжить игру]                                       │
│                                                            │
├────────────────────────────────────────────────────────────┤
│    📊 Статистика: X игр сыграно, Y финалов достигнуто      │
└────────────────────────────────────────────────────────────┘
```

#### Screen 2: Игровой экран (Game)

```
┌────────────────────────────────────────────────────────────┐
│  [≡ Меню]                              Ход: 5 / 30        │
├───────────────────────────────┬────────────────────────────┤
│                               │  СТАТЫ                     │
│  ┌─────────────────────────┐  │  ━━━━━━━━━━━━━━━━━━━━━━━  │
│  │                         │  │  💪 Сила:     52 [████░] │
│  │   [ИЛЛЮСТРАЦИЯ]         │  │  🫀 Выносливость: 48 [███░]│
│  │                         │  │  🏃 Ловкость:  55 [████░] │
│  └─────────────────────────┘  │  😊 Мораль:    45 [███░]  │
│                               │  ⚖️ Дисциплина: -10 [░░░░] │
│  Ситуация:                    │                            │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━   │  ЗВАНИЕ: Рядовой           │
│                               │  СТАТУС: Запах             │
│  Внезапно сержант закричал:   │                            │
│  'Смирно! Рота, построение!'  │  ФЛАГИ: риск_рапорта       │
│                               │                            │
│  Чёрт, я не успел забрать     │                            │
│  берцы из сушилки.            │                            │
│                               │                            │
├───────────────────────────────┴────────────────────────────┤
│  ВЫБОР:                                                   │
│  ┌──────────────────────────────────────────────────────┐ │
│  │ 1. Напрячь салагу                    [ВЫБРАТЬ]      │ │
│  └──────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────┐ │
│  │ 2. Быстро сгонять самому                [ВЫБРАТЬ]    │ │
│  └──────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────┐ │
│  │ 3. Поискать под кроватью              [ недоступно]│ │
│  └──────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────┐ │
│  │ 4. Выйти в тапочках                    [ недоступно]│ │
│  └──────────────────────────────────────────────────────┘ │
├────────────────────────────────────────────────────────────┤
│  [Последние события ▼]                                   │
└────────────────────────────────────────────────────────────┘
```

#### Screen 3: Экран результата проверки

```
┌────────────────────────────────────────────────────────────┐
│                                                            │
│                     ✓ УСПЕХ                               │
│                                                            │
│  Салага молча приносит берцы.                             │
│                                                            │
│  ─────────────────────────────────────────────────────    │
│                                                            │
│  Эффекты:                                                 │
│  • Дисциплина: -8 (0 → -8)                               │
│  • Мораль: -2 (50 → 48)                                   │
│                                                            │
│  ┌──────────────────────────────────────────────────────┐ │
│  │                  СЛЕДУЮЩИЙ ХОД →                     │ │
│  └──────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────┘
```

#### Screen 4: Экран финала

```
┌────────────────────────────────────────────────────────────┐
│                                                            │
│               ФИНАЛ: ТИХИЙ ДЕМБЕЛЬ                         │
│                                                            │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━   │
│                                                            │
│  Ты прошёл службу без лишнего шума. Не герой, не изгой —   │
│  просто ещё один мужик, который отдал долг Родине.          │
│                                                            │
│  30 дней. 15 событий. 47 выборов.                          │
│                                                            │
│  ─────────────────────────────────────────────────────    │
│                                                            │
│  Финальные статы:                                         │
│  💪 Сила: 52  🫀 Выносливость: 48                         │
│  🏃 Ловкость: 55  😊 Мораль: 35  ⚖️ Дисциплина: +5        │
│                                                            │
│  Звание: Ефрейтор  |  Статус: Дух                          │
│                                                            │
│  ┌──────────────────────────────────────────────────────┐ │
│  │               ИГРАТЬ СНОВА                           │ │
│  └──────────────────────────────────────────────────────┘ │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

### 14.2 UI Components

| Компонент | Описание | States |
|-----------|---------|--------|
| `StatBar` | Прогресс-бар для стата | default, low (< 30), critical (< 15), buffed |
| `ChoiceButton` | Кнопка выбора | enabled, disabled, selected, loading |
| `EventCard` | Карточка события | loading, loaded, error |
| `ResultBadge` | Бейдж результата проверки | success, partial, failure, catastrophic |
| `FinalScreen` | Экран финала | по типу финала |
| `GameSidebar` | Боковая панель статов | collapsible на mobile |
| `HistoryDrawer` | Выдвижная панель истории | collapsed, expanded |

### 14.3 Loading States

| Состояние | UI |
|-----------|-----|
| Загрузка события | Skeleton: текст + 4 кнопки |
| Обработка выбора | Spinner на кнопке + disabled |
| Загрузка сохранения | Full-screen loader |

### 14.4 Error States

| Ошибка | UI |
|--------|-----|
| Network error | Toast: «Нет соединения» + Retry button |
| Validation error | Inline message под полем |
| Server error | Modal: «Что-то пошло не так» + Restart button |

### 14.5 Responsive Breakpoints

| Breakpoint | Layout |
|------------|--------|
| < 768px (mobile) | Stack: sidebar above, event below; collapsible |
| ≥ 768px (tablet) | Split: 2/3 event, 1/3 sidebar |
| ≥ 1024px (desktop) | Split: 3/4 event, 1/4 sidebar + history |

### 14.6 Animation & Transition Timings

| Animation | Duration | Easing | Trigger |
|-----------|----------|--------|---------|
| Stat change (bar) | 300ms | ease-out | After choose |
| Stat increase | 300ms + green flash | ease-out | MOR/STR/END/AGI increase |
| Stat decrease | 300ms + red flash | ease-out | Any stat decrease |
| Choice button hover | 150ms | ease-in-out | Mouse enter |
| Choice button disabled | 200ms | ease-out | Availability change |
| Screen transition | 200ms | ease-in-out | Between screens |
| Modal open | 200ms | ease-out | Error/success modal |
| Modal close | 150ms | ease-in | Dismiss |
| Toast notification | 300ms slide-in, 5s visible | ease-out | Error/success |
| Result badge appear | 200ms | ease-out | After check |

### 14.7 Sound Design Requirements

| Sound | Type | Duration | Trigger |
|-------|------|----------|---------|
| Choice select | SFX | < 100ms | Click on choice |
| Check success | SFX (positive) | < 200ms | Success outcome |
| Check failure | SFX (negative) | < 200ms | Failure/partial outcome |
| Stat increase | SFX (subtle positive) | < 100ms | Any stat going up |
| Stat decrease | SFX (subtle negative) | < 100ms | Any stat going down |
| New event | SFX (whoosh/transition) | < 300ms | Next event appears |
| Rank up | SFX (fanfare) | < 500ms | Formal/informal promotion |
| Game over | Music sting | < 3s | Final screen |
| Final victory | Music | Loop | Final screen |

**Примечание:** Звуковое сопровождение НЕ является blocker для MVP. Но базовые SFX должны быть реализованы для feedback.

### 14.8 Детальная спецификация дизайна

> **Важно:** Полная спецификация UI/UX дизайна вынесена в отдельный документ.

| Документ | Описание |
|----------|----------|
| `docs/DESIGN_SPECIFICATION.md` | Детальная спецификация UI/UX: экраны, компоненты, дизайн-система, промты для Figma |
| `docs/Army_Game_Prompts.md` | AI-промты для генерации игровых иллюстраций |

**Краткое содержание DESIGN_SPECIFICATION.md:**
- Mobile-first подход (375px базовая ширина)
- 6 основных экранов: Start, Gameplay, Stats Panel, History Panel, Final, Loading
- Дизайн-система: цвета, типографика, spacing (8px grid), компоненты
- Figma промты для генерации макетов
- Референс: Hoosegow: Prison Survival

**Референсные скриншоты:** `reference_images/` — реальные скриншоты Hoosegow для анализа UI паттернов.

---

## 15. Constraints & Assumptions

### 15.1 Технические ограничения

| Ограничение | Значение |
|-------------|----------|
| Браузеры | Chrome 90+, Firefox 88+, Safari 14+, Edge 90+ |
| Node.js | ^20.19.0 |
| Go | 1.21+ |
| PostgreSQL | 15+ |
| Redis | 7+ |
| Docker | 20.10+ |

### 15.2 Внешние зависимости

| Зависимость | Версия | Fallback |
|-------------|--------|----------|
| GraphQL endpoint | localhost:3001 | Hardcoded env |
| PostgreSQL | localhost:5432 | docker-compose |
| Redis | localhost:6379 | In-memory (dev) |
| Event templates | `/content/events/*.json` | Bundled JSON |

### 15.3 Rate Limits

| Endpoint | Limit |
|----------|-------|
| `/graphql` | 100 req/min per IP |
| `startGame` | 5 req/min per IP |
| `choose` | 30 req/min per IP |

### 15.4 Assumptions

| # | Assumption | Mitigation |
|---|-----------|------------|
| A1 | Игрок не будет пытаться взломать клиент | Backend валидирует все операции |
| A2 | Максимум 100 одновременных игроков | Capacity planning на это значение |
| A3 | 20 шаблонов достаточно для MVP | Можно добавить до 90 без изменения кода |
| A4 | Игроки не будут спамить запросы | Rate limiting + debounce на клиенте |
| A5 | Архитектура фронтенда — Vue 3 (не Vue как в docs) | Используется Vue 3 согласно коду |
| A6 | Двойной клик на кнопку выбора — редкий кейс | Optimistic locking с retry решает эту проблему |
| A7 | Баланс будет протестирован до implementation | Симулятор в Excel создаётся до кода |

### 15.5 Content Creation Pipeline (Out of Scope)

**Важно:** Создание игрового контента (событий, диалогов, баланса) выходит за scope данного ТЗ и описано в отдельном документе.

**Требуемые роли для создания контента:**
| Роль | Responsibility |
|------|---------------|
| Narrative Designer | Написание событий, диалогов, lore |
| Game Designer | Баланс формул, метрики, difficulty curve |
| Domain Expert | Валидация реализма (армейский сленг, традиции) |

**Content deliverables:**
- 20-30 событий для MVP
- Balance simulator (Excel/Sheets)
- JSON event templates с validation schema

---

## 15.2 Testing Strategy (TDD)

> **Важно:** Проект использует **TDD (Test-Driven Development)** подход. Сначала пишутся тесты, затем код.

### 15.2.1 Общие принципы TDD

| Принцип | Описание |
|---------|----------|
| Red-Green-Refactor | Red (тест падает) → Green (тест проходит) → Refactor (рефакторинг) |
| Тесты first | Код пишется только после написания теста |
| One test at a time | Один тест — одна проверка |
| Fast tests | Все unit тесты выполняются < 1 минуты |

### 15.2.2 Backend Testing (Go)

**Уровни тестирования:**

| Уровень | Инструменты | Что тестируется |
|---------|------------|-----------------|
| Unit | `testing` + `testify` | Business logic, services |
| Integration | `testcontainers` + `sqlx` | Database, repositories |
| E2E | `httptest` | GraphQL API |

**Тестируемые компоненты:**

```go
// internal/service/player_service_test.go
func TestPlayerService_Create(t *testing.T) {
    // Arrange
    repo := &mockPlayerRepository{}
    svc := NewPlayerService(repo)
    
    // Act
    player, err := svc.CreatePlayer(context.Background())
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, 50, player.Stats.STR)
    assert.Equal(t, 50, player.Stats.END)
    assert.Equal(t, 1, player.Turn)
}
```

**Test coverage требования:**

| Компонент | Min Coverage |
|----------|-------------|
| domain/entities | 90% |
| service/ | 80% |
| handler/ | 60% |
| **Общий** | **70%** |

**Команды:**

```bash
# Запустить все тесты
go test ./...

# Запустить с покрытием
go test -cover ./...

# Запустить конкретный тест
go test ./internal/service/... -run TestPlayerService -v

# Запустить интеграционные тесты
go test ./internal/integration/... -tags=integration
```

### 15.2.3 Frontend Testing (Vue 3)

**Уровни тестирования:**

| Уровень | Инструменты | Что тестируется |
|---------|------------|-----------------|
| Unit | Vitest | Composables, utilities |
| Component | Vue Testing Library | Components |
| Integration | Vitest + Vue Test Utils | User flows |
| E2E | Playwright | Critical paths |

**Тестируемые компоненты:**

```typescript
// features/start-game/useStartGame.test.ts
import { describe, it, expect, vi } from 'vitest'
import { useStartGame } from './useStartGame'

vi.mock('@/shared/api/mutations/startGame', () => ({
  startGameMutation: vi.fn()
}))

describe('useStartGame', () => {
  it('should start new game and return player', async () => {
    const { startGame, isLoading } = useStartGame()
    
    await startGame()
    
    expect(isLoading.value).toBe(false)
    // ...
  })
})
```

**Test coverage требования:**

| Компонент | Min Coverage |
|----------|-------------|
| composables / use* | 80% |
| components | 60% |
| **Общий** | **60%** |

**Команды:**

```bash
# Запустить все тесты
npm test

# Запустить с покрытием
npm run test -- --coverage

# Запустить тесты в watch режиме
npm run test -- --watch

# Запустить E2E тесты
npx playwright test
```

### 15.2.4 Game Balance Testing

**Баланс валидируется через Python симулятор:**

```bash
cd balance_simulator
python simulator.py --runs 1000 --seed 42
```

**Требования к балансу:**

| Метрика | Target | Tolerance |
|---------|--------|-----------|
| Death Rate | 20-35% | ±5% |
| Victory Rate | 65-80% | ±5% |
| Avg MOR/turn | -1.5 to -2.5 | ±0.5 |
| Avg Death Turn | 12-28 | ±3 |

**Правило:** Любое изменение в `events.json` должно быть проверено через симулятор.

### 15.2.5 TDD Workflow

```
1. Написать тест (RED)
   ↓
2. Запустить тест → Падает
   ↓
3. Написать минимальный код (GREEN)
   ↓
4. Запустить тест → Проходит
   ↓
5. Рефакторить код
   ↓
6. Запустить все тесты → Проходят
   ↓
7. Commit
```

### 15.2.6 CI Integration

**Тесты запускаются автоматически:**

- При push в любую ветку
- При Pull Request
- Тесты должны проходить перед merge

---

## 15.3 CI/CD Pipeline

### 15.3.1 GitHub Actions Workflow

```yaml
# .github/workflows/ci.yml
name: CI/CD

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  # Backend
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Lint
        run: golangci-lint run ./...
      
      - name: Test
        run: go test -v -cover ./...
      
      - name: Build
        run: go build -o bin/server ./cmd/server

  # Frontend
  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Node
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      
      - name: Install
        run: npm ci
      
      - name: Lint
        run: npm run lint
      
      - name: Type Check
        run: npm run typecheck
      
      - name: Test
        run: npm run test -- --coverage
      
      - name: Build
        run: npm run build

  # Balance Simulator
  balance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Balance Tests
        run: |
          cd balance_simulator
          python simulator.py --runs 1000 --seed 42
      
      - name: Validate Metrics
        run: python validate_results.py

  # Docker Build
  docker:
    runs-on: ubuntu-latest
    needs: [backend, frontend]
    steps:
      - uses: actions/checkout@v4
      
      - name: Build Docker
        run: docker-compose build
      
      - name: Push to Registry
        if: github.ref == 'refs/heads/main'
        run: docker-compose push
```

### 15.3.2 Сборочные команды

**Backend:**

```bash
# Линтинг
golangci-lint run ./...

# Типы
go vet ./...

# Тесты
go test -v ./...

# Покрытие
go test -cover ./...

# Билд
go build -o bin/server ./cmd/server
```

**Frontend:**

```bash
# Установка зависимостей
npm ci

# Линтинг
npm run lint

# Проверка типов
npm run typecheck

# Тесты
npm test

# Билд
npm run build
```

**Balance:**

```bash
cd balance_simulator
python simulator.py --runs 1000 --seed 42
```

### 15.3.3 Деплой (Future)

| Stage | Environment | Trigger |
|-------|-------------|---------|
| Dev | Docker Compose (localhost) | manual |
| Staging | Cloud (GCP/AWS) | tag `v*.*.*-rc*` |
| Prod | Cloud (GCP/AWS) | tag `v*.*.*` |

---

## 16. Acceptance Criteria

### 16.1 Критерии готовности MVP

| # | Критерий | Definition of Done |
|---|---------|-------------------|
| AC-01 | Игрок может начать новую игру | startGame возвращает Player + EventInstance |
| AC-02 | Игрок может сделать выбор | choose применяет эффекты, обновляет статы |
| AC-03 | Статы корректно обновляются | Clamp работает на границах диапазонов |
| AC-04 | Проверки работают детерминированно | Unit-тесты проходят с mock rand |
| AC-05 | Генерируются разные описания событий | Подстановка переменных работает |
| AC-06 | Прогрессия обновляется | Звания и статусы меняются по порогам |
| AC-07 | Финалы достигаются | 3 типа финалов реализованы |
| AC-08 | Сохранение работает | Данные сохраняются между сессиями |
| AC-09 | UI отображает игровой поток | Все экраны реализованы |
| AC-10 | GraphQL API полностью функционален | Все queries и mutations работают |
| AC-11 | Optimistic locking работает | При CONCURRENT_MODIFICATION клиент перезагружает состояние и повторяет |
| AC-12 | Weighted event selection | События выбираются по весам, не pure random |
| AC-13 | Difficulty curve | Ранние ходы легче, поздние — сложнее |
| AC-14 | MOR recovery guarantee | После 3 MOR- событий — гарантированно MOR+ событие |

### 16.2 Definition of Done для задачи

Каждая задача считается выполненной, когда:
1. Код написан и соответствует code style
2. Unit-тесты написаны и проходят
3. Integration тесты написаны и проходят
4. Типовая проверка (type-check) проходит
5. Линтинг проходит без warnings
6. PR создан и reviewed

---

## 17. Appendix

### 17.1 Глоссарий кодовых имён

| Кодовое имя | Значение |
|------------|----------|
| LEGO-template | Паттерн проектирования событий с переменными |
| Clamp | Ограничение значения в диапазоне |
| Demob | Сокращение от «демобилизация» = конец службы |
| Black eye | Синяк под глазом = поражение в драке |

### 17.2 Формулы расчёта

**Формула шанса (Probability):**
```
chance = clamp(STAT * coeff1 + DISC * coeff2 + MOR * coeff3, 0, 100)
success = rand(100) < chance
partial = rand(100) < chance + partialThreshold
```

**Формула threshold:**
```
calculated = STAT + DISC * coeff1 + MOR * coeff2
success = calculated >= threshold
```

**Формула catastrophic notice:**
```
noticeChance = baseChance + DISC * coeff
noticed = rand(100) < noticeChance
```

---

**Документ утверждён:** _________________ / _________________ / 2026-03-22
