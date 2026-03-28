# UI/UX Design Specification — Narrative Survival Game «Армейка»

**Версия:** 1.0  
**Дата:** 2026-03-25  
**Статус:** Черновик для разработки  
**Подход:** Mobile-First  

---

## 1. Overview

### 1.1 Продукт

Веб-игра «Армейка» — narrative survival / choice-based RPG, симулирующая опыт срочной военной службы. Игрок — солдат, вынужденный выживать в системе постоянного давления между формальной уставной иерархией и неформальной солдатской средой.

### 1.2 Целевая аудитория

| Сегмент | Описание |
|---------|----------|
| Основной | Казуальные геймеры, 18-35 лет |
| Сессия | 15-30 минут на игру |
| Платформа | Мобильные устройства (приоритет), десктоп |
| Язык | Русский |

### 1.3 Референсы

#### Hoosegow: Prison Survival (Основной референс)

**Описание игры:**
Hoosegow: Prison Survival — мобильная choice-based RPG от разработчика Demian Credit. Игрок попадает в тюрьму и должен выживать, принимая решения, которые влияют на его характеристики (Hunger, Health, Sanity, Respect).

**Ссылка:** App Store / Google Play (мобильная игра)

**Ключевые особенности UI/UX:**

| Аспект | Hoosegow | Применение к Армейка |
|--------|----------|---------------------|
| **Layout** | Stats bar вверху, иллюстрация по центру, 3 кнопки внизу | Адаптировать для military тематики |
| **Иллюстрации** | Стилизованные векторные сцены тюрьмы | Использовать military environments (казарма, столовая, плац) |
| **Цветовая палитра** | Тёмная, коричнево-серая (тюремная) | Тёмная olive-green (военная) |
| **Типографика** | Roboto, чистая и читаемая | Roboto Condensed для заголовков |
| **Анимации** | Плавные переходы между событиями | Адаптировать crossfade для иллюстраций |
| **Текст** | Крупный, хорошо читаемый на мобильных | 16px body, min 44px tap targets |
| **Тон** | Тёмный юмор, реалистичная атмосфера | Реалистичный military narrative |

**Элементы Hoosegow для адаптации:**

1. **Stats Display (Верхняя панель)**
   - Hoosegow: 4 характеристики с иконками и числовыми значениями
   - Армейка: 5 статов (STR, END, AGI, MOR, DISC) — адаптировать layout

2. **Event Area (Центральная часть)**
   - Hoosegow: Крупная иллюстрация сцены (70% экрана)
   - Армейка: Иллюстрация 40-50% + описание события

3. **Choice Buttons (Нижняя часть)**
   - Hoosegow: 3 кнопки на всю ширину, вертикальный стек
   - Армейка: 2-4 кнопки, тот же паттерн

4. **Modal Panels**
   - Hoosegow: Slide-up панели для детальной информации
   - Армейка: Stats Panel, History Panel — тот же паттерн

5. **Game Flow**
   - Hoosegow: Дни в тюрьме → финал
   - Армейка: Дни в армии → дембель (3 типа финала)

**Визуальный стиль Hoosegow:**

```
- Semi-realistic illustration style
- Clean vector-like shapes
- Soft gradient shading on characters/environments
- High contrast lighting with rim light
- Dark, desaturated color palette
- Muted browns, grays, and warm accent colors
- Cinematic framing with strong depth
- Camera slightly above eye level (centered perspective)
- Space reserved at bottom for UI/buttons
```

**Почему Hoosegow — хороший референс:**

1. ✅ Похожий жанр: narrative survival / choice-based RPG
2. ✅ Mobile-first дизайн с thumb-friendly layout
3. ✅ Clear visual hierarchy
4. ✅ Успешная монетизация (free-to-play)
5. ✅ Сильный narrative focus
6. ✅ Атмосферный визуал без перегруза

#### Дополнительные референсы

| Игра / Продукт | Релевантность |
|----------------|---------------|
| **This War of Mine** | Тёмная атмосфера, survival narrative, decision consequences |
| **Papers, Please** | Текстовый нарратив, бюрократия, моральные дилеммы |
| **Reigns** | Карточный UI, тапы вместо кнопок, последствия решений |
| ** Fallout Shelter** | Stats management, progress tracking |

---

## 2. User Flow

### 2.1 Основной поток

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  START      │────▶│  GAMEPLAY   │────▶│   FINAL     │
│  SCREEN     │     │  (Events)   │     │   SCREEN    │
└─────────────┘     └──────┬──────┘     └─────────────┘
                          │
                    ┌─────┴─────┐
                    │  STATS    │
                    │  PANEL    │
                    └───────────┘
```

### 2.2 Детализация шагов

1. **Открытие приложения** → Проверка сохранённой игры
   - Есть сохранение → Загрузка Gameplay Screen
   - Нет сохранения → Показ Start Screen

2. **Start Screen** → Пользователь нажимает "Начать игру"
   - Создание нового Player со стартовыми статами
   - Генерация первого события
   - Переход на Gameplay Screen

3. **Gameplay Screen** → Пользователь читает событие и делает выбор
   - Выбор варианта → API запрос
   - Применение эффектов → Сохранение
   - Генерация следующего события или финал

4. **Финал** → Отображение результата
   - **Тихий дембель** (Приспособленец): turn ≥ 27, MOR > 0, DISC ∈ [-20,+20]
   - **Уважаемый дембель**: turn ≥ 27, MOR > 0, DISC < -50, informalStatus ≥ «Дед»
   - **Сломанный дембель** (Косячник): MOR = 0

---

## 3. Screen List

### 3.1 Mobile Screens (Приоритет)

| # | Экран | Описание | Приоритет |
|---|-------|----------|-----------|
| 1 | Start Screen | Главное меню, кнопка старта | 🔴 Высокий |
| 2 | Gameplay Screen | Основной игровой экран | 🔴 Высокий |
| 3 | Stats Panel | Модальная панель статов | 🔴 Высокий |
| 4 | History Panel | Панель истории событий | 🟡 Средний |
| 5 | Final Screen | Экран окончания игры | 🔴 Высокий |
| 6 | Loading State | Состояние загрузки | 🟡 Средний |

### 3.2 Desktop Screens

| # | Экран | Описание | Приоритет |
|---|-------|----------|-----------|
| 1 | Desktop Gameplay | Расширенная версия мобильного | 🟢 Низкий |
| 2 | Desktop Start | Адаптированное меню | 🟢 Низкий |

---

## 4. Screen Specifications

### 4.1 Start Screen (Мобильный)

```
┌─────────────────────────┐
│                         │
│      ЛОГОТИП            │
│                         │
│      "АРМЕЙКА"          │
│                         │
│   [Начать игру]         │
│                         │
│   Продолжить (если       │
│   есть сохранение)      │
│                         │
└─────────────────────────┘
```

**Компоненты:**
- Логотип игры (текстовый или графический)
- Заголовок «АРМЕЙКА» — крупный, жирный, upper-case
- Кнопка «Начать игру» — основная CTA, яркий акцент
- Кнопка «Продолжить» — вторичная (показывается при наличии сохранения)

**Поведение:**
- При наличии gameId в localStorage → показать «Продолжить»
- Тап по «Начать игру» → мутация startGame() → переход на Gameplay

**Визуал:**
- Фон: военная казарма (иллюстрация)
- Центрированный layout
- Кнопки в нижней трети экрана

---

### 4.2 Gameplay Screen (Мобильный)

```
┌─────────────────────────┐
│ [📊]  ДЕНЬ 5    [📜]  │  ← Status Bar (48px)
├─────────────────────────┤
│  STR ████░░  60        │
│  END █████░  70        │  ← Stats Row (compact)
│  AGI ███░░░  50        │
│  MOR ██░░░░  30 ▼     │
│  DSC ██████ +15       │
├─────────────────────────┤
│                         │
│   ┌─────────────────┐   │
│   │                 │   │
│   │   ИЛЛЮСТРАЦИЯ  │   │  ← Event Illustration
│   │   (40-50%)     │   │     (стилизованное
│   │                 │   │      изображение)
│   └─────────────────┘   │
│                         │
├─────────────────────────┤
│                         │
│  "Сержант Кривонос     │  ← Event Description
│  проверяет построение. │  (scrollable)
│  Ты стоишь не по        │
│  уставу..."             │
│                         │
├─────────────────────────┤
│                         │
│ [Выбор 1: "Принять     │  ← Choice Buttons
│  позу"            ]    │     (2-4 кнопки)
│                         │
│ [Выбор 2: "Стоять     │  ← Min height: 44px
│  как стоял"        ]   │     (tap target)
│                         │
│ [Выбор 3: "Спросить   │
│  у старших"       ]    │
│                         │
└─────────────────────────┘
```

**Компоненты:**

#### Status Bar (верхняя панель)
- Иконка статов [📊] — открывает Stats Panel
- Заголовок «ДЕНЬ X» — текущий ход
- Иконка истории [📜] — открывает History Panel

#### Stats Row (строка статов)
- 5 компактных индикаторов: STR, END, AGI, MOR, DISC
- Визуальная полоска (bar) + числовое значение
- MOR показан с треугольником-индикатором (▼ при низком)
- DISC показан с + или - префиксом

#### Event Illustration (область иллюстрации)
- Стилизованное изображение (40-50% высоты экрана)
- Соотношение сторон: 16:9 или 4:3
- Плавная анимация смены при переходе к следующему событию

#### Event Description (текст события)
- Описание ситуации на русском языке
- Шрифт: Inter/Roboto, 16px
- Скролл при длинном тексте
- Подстановка переменных: {soldier_name}, {time}, {location}

#### Choice Buttons (кнопки выбора)
- 2-4 кнопки (варианты действий)
- Минимальная высота: 44px (touch target)
- Полная ширина с padding 16px по бокам
- Disabled состояние для недоступных выборов (серый, неактивный)
- Анимация нажатия (scale 0.98)

**Поведение:**
- При выборе → блокировка кнопок → API запрос → анимация результата → следующее событие
- При MOR < 30 → индикатор MOR мигает красным
- При недоступном выборе → кнопка серая + tooltip "Недоступно"

---

### 4.3 Stats Panel (Модальное окно)

```
┌─────────────────────────┐
│  ←  ХАРАКТЕРИСТИКИ     │  ← Header (закрыть)
├─────────────────────────┤
│                         │
│  ФОРМАЛЬНЫЙ СТАТУС     │
│  ┌───────────────────┐ │
│  │    Рядовой        │ │  ← Formal Rank
│  └───────────────────┘ │
│                         │
│  ┌───┐ ┌───┐ ┌───┐     │
│  │STR│ │END│ │AGI│     │  ← Physical Stats
│  │ 60│ │ 70│ │ 50│     │
│  └───┘ └───┘ └───┘     │
│                         │
│  ┌───┐ ┌───┐           │
│  │MOR│ │DSC│           │  ← Mental/Social
│  │ 30│ │+15│           │
│  └───┘ └───┘           │
│                         │
│  ─────────────────────│
│  НЕФОРМАЛЬНЫЙ СТАТУС   │
│  ┌───────────────────┐ │
│  │    Запах         │ │  ← Informal Status
│  └───────────────────┘ │
│                         │
│  ПРОГРЕСС              │
│  День 5 из 30          │  ← Turn Progress
│  ████████░░░░░░░  17%  │
│                         │
└─────────────────────────┘
```

**Компоненты:**
- Header с кнопкой закрытия [←]
- Блок «Формальный статус» — текущее звание
- Карточки статов (5 штук) с иконками и значениями
- Линия-разделитель
- Блок «Неформальный статус» — текущий статус
- Прогресс-бар «День X из 30»

**Визуал:**
- Modal overlay с затемнением фона (rgba(0,0,0,0.7))
- Surface color: #2D3527
- Анимация появления: slide-up 300ms ease-out

---

### 4.4 History Panel (Панель истории)

```
┌─────────────────────────┐
│     ИСТОРИЯ СОБЫТИЙ    │  ← Header
├─────────────────────────┤
│                         │
│  День 5: "Проверка"    │  ← Event Entry
│  > "Принять позу" ✓   │
│  MOR: 50 → 48 (-2)    │
│  ─────────────────    │
│                         │
│  День 4: "Построение"  │
│  > "Стоять как стоял"  │
│  DSC: +5 → +10 (+5)   │
│  ─────────────────    │
│                         │
│  День 3: "Столовая"    │
│  > "Сесть с дедами"    │
│  MOR: 52 → 55 (+3)    │
│  ─────────────────    │
│                         │
│  ... (скролл, макс 10) │
│                         │
└─────────────────────────┘
```

**Компоненты:**
- Header с названием
- Список последних 10 событий (List Items)
- Каждый item: день, название события, выбор, изменение статов
- Индикатор успеха/неудачи (✓ / ✗)
- Скролл внутри панели

**Визуал:**
- Slide-up panel (bottom sheet)
- Высота: 70% экрана
- Анимация: slide-up 300ms

---

### 4.5 Final Screen (Экран финала)

```
┌─────────────────────────┐
│                         │
│      РЕЗУЛЬТАТ         │
│                         │
│  ═══════════════════   │
│                         │
│  "ТИХИЙ ДЕМБЕЛЬ"       │  ← Приспособленец
│  "УВАЖАЕМЫЙ ДЕМБЕЛЬ"   │  или
│  "СЛОМАННЫЙ ДЕМБЕЛЬ"   │  ← Косячник
│                         │
│  ═══════════════════   │
│                         │
│  "Ваша служба прошла  │
│   неспеша, вы не      │
│   выделялись..."       │  ← Final Description
│                         │
│  ─────────────────────  │
│                         │
│  ИТОГИ:                 │
│  Дней прослужено: 30   │
│  MOR: 12              │
│  DSC: +5 (Рядовой)    │
│                         │
│  ─────────────────────  │
│                         │
│  [ИГРАТЬ СНОВА]        │  ← Primary CTA
│                         │
└─────────────────────────┘
```

**Компоненты:**
- Заголовок «РЕЗУЛЬТАТ»
- Название финала (крупно, акцент)
- Описание финала (текст 2-3 предложения)
- Блок «Итоги» — статистика игры
- Кнопка «Играть снова»

**Типы финалов:**

| Финал | Условие | Описание | Цвет акцента |
|-------|---------|----------|--------------|
| **Тихий дембель**<br>(Приспособленец) | turn ≥ 27, MOR > 0, DISC ∈ [-20,+20] | «Ваша служба прошла неспеша, вы не выделялись, не косячили, выполняли приказы "шакалов" и офицеров. В общем — вы просто жили.» | #5D7A4A (зеленый) |
| **Уважаемый дембель** | turn ≥ 27, MOR > 0, DISC < -50, status ≥ «Дед» | «Вы ушли со службы в золотом кашне. Вся рота провожала вас со слезами на глазах, и обсуждала, как вы смогли так "подняться" за такой короткий срок. А ваш кореш — старший прапорщик Залупенко плакал от горя, ведь вы смогли "загнать" 300 тонн гнутых гвоздей...» | #4A5D3F (военный зеленый) |
| **Сломанный дембель**<br>(Косячник) | MOR = 0 | «Служба в армии принесла вам лишь горе, затрещины и сломанный нос, когда вы ночью упали с кровати, исполняя "летучую мышь" по команде сержанта Златопузова. Впрочем, вы выжили — а это главное.» | #8B3A3A (красный) |

---

### 4.6 Loading State

```
┌─────────────────────────┐
│                         │
│         ⏳              │  ←Spinner
│                         │
│    Загрузка...         │  ← Text
│                         │
└─────────────────────────┘
```

- Центрированный спиннер
- Текст «Загрузка...»
- Используется при: startGame, choose, loadGame

---

## 5. Design System

### 5.1 Цветовая палитра

#### Основные цвета

| Название | HEX | Применение |
|----------|-----|------------|
| Background Primary | `#1A1F16` | Основной фон |
| Background Secondary | `#252B1E` | Фон карточек, панелей |
| Surface | `#2D3527` | Модальные окна, панели |
| Surface Elevated | `#353D2E` | Поднятые элементы |

#### Акцентные цвета

| Название | HEX | Применение |
|----------|-----|------------|
| Accent Primary | `#4A5D3F` | Основные кнопки, активные элементы |
| Accent Secondary | `#8B7355` | Вторичные элементы |
| Accent Tertiary | `#5D7A4A` | Успех, положительные изменения |

#### Текстовые цвета

| Название | HEX | Применение |
|----------|-----|------------|
| Text Primary | `#E8E4DC` | Основной текст |
| Text Secondary | `#9A968E` | Вторичный текст, описания |
| Text Muted | `#6B665E` | Плейсхолдеры, подсказки |

#### Статусные цвета

| Название | HEX | Применение |
|----------|-----|------------|
| Success | `#5D7A4A` | Успешные исходы |
| Warning | `#B8860B` | Низкие статы (MOR < 30) |
| Danger | `#8B3A3A` | Критические статы, неудачи |
| Info | `#4A5D7F` | Информационные элементы |

#### Граница

| Название | HEX | Применение |
|----------|-----|------------|
| Border | `#3D4533` | Границы элементов |
| Border Light | `#4A5540` | Границы в светлых местах |

---

### 5.2 Типографика

#### Шрифты

| Применение | Шрифт | Вес | Размер |
|------------|-------|-----|--------|
| Заголовок (H1) | Roboto Condensed | Bold | 24px |
| Заголовок (H2) | Roboto Condensed | Bold | 20px |
| Подзаголовок | Roboto | Medium | 18px |
| Body | Inter / Roboto | Regular | 16px |
| Body Small | Inter / Roboto | Regular | 14px |
| Caption | Inter / Roboto | Regular | 12px |
| Stats Numbers | Roboto Mono | Medium | 14px |
| Buttons | Roboto Condensed | Bold | 14px |

#### Стили текста

```
H1 — Заголовок экрана
- Roboto Condensed, Bold
- Uppercase
- Letter-spacing: 2px
- Color: Text Primary

H2 — Заголовок секции
- Roboto Condensed, Bold
- Letter-spacing: 1px
- Color: Text Primary

Body — Основной текст
- Inter/Roboto, Regular
- Line-height: 1.5
- Color: Text Primary

Caption — Подписи
- Inter/Roboto, Regular
- Line-height: 1.4
- Color: Text Secondary
```

---

### 5.3 Spacing (Сетка 8px)

| Название | Значение | Применение |
|----------|----------|------------|
| xs | 4px | Tight spacing |
| sm | 8px | Icon padding |
| md | 16px | Standard padding |
| lg | 24px | Section spacing |
| xl | 32px | Large gaps |
| xxl | 48px | Screen margins |

#### Padding стандартных элементов

| Элемент | Padding |
|---------|---------|
| Card | 16px |
| Button | 12px 24px |
| Modal | 24px |
| Screen | 16px horizontal |

---

### 5.4 Border Radius

| Название | Значение | Применение |
|----------|----------|------------|
| Small | 4px | Кнопки, маленькие элементы |
| Medium | 8px | Карточки, инпуты |
| Large | 12px | Модальные окна, панели |
| Full | 50% | Круглые элементы |

---

### 5.5 Тени

| Название | Значение | Применение |
|----------|----------|------------|
| Card | `0 2px 8px rgba(0,0,0,0.3)` | Карточки |
| Modal | `0 8px 24px rgba(0,0,0,0.5)` | Модальные окна |
| Button | `0 2px 4px rgba(0,0,0,0.2)` | Кнопки |
| Elevated | `0 4px 12px rgba(0,0,0,0.4)` | Поднятые элементы |

---

### 5.6 Компоненты

#### Button (Кнопка)

```
States:
- Default: Background Accent Primary, Text Primary
- Hover: Background #5A6D4F (lighten 10%)
- Active/Pressed: Scale 0.98, Background #3A4D2F (darken 10%)
- Disabled: Background Surface, Text Text Muted, opacity 0.5

Sizes:
- Large: height 56px, font-size 16px
- Medium: height 44px, font-size 14px (default for choices)
- Small: height 36px, font-size 12px
```

#### Stat Bar (Индикатор стата)

```
Structure:
┌─ Иконка ─┐ ┌─ Bar ──────────┐ ┌─ Value ─┐
│   [ICON] │ │ ████████░░░░░░ │ │   60   │
└──────────┘ └────────────────┘ └─────────┘
   24px           flex-grow         32px

Bar Colors:
- Normal (>50%): Accent Primary
- Warning (30-50%): Warning
- Critical (<30%): Danger

Bar Animation: width transition 300ms ease
```

#### Card (Карточка)

```
Background: Surface (#2D3527)
Border: 1px solid Border (#3D4533)
Border-radius: Medium (8px)
Padding: md (16px)
Shadow: Card
```

#### Modal / Panel (Модальное окно)

```
Overlay: rgba(0,0,0,0.7)
Background: Surface (#2D3527)
Border-radius: Large (12px)
Padding: lg (24px)
Animation: fade-in 200ms + slide-up 300ms
```

#### Toast (Уведомление)

```
Position: bottom-center, 80px from bottom
Background: Surface Elevated (#353D2E)
Border-radius: Medium (8px)
Padding: sm md (8px 16px)
Animation: fade-in 200ms, auto-dismiss 3s
```

---

## 6. Responsive Breakpoints

### 6.1 Mobile-First подход

| Breakpoint | Ширина | Описание |
|------------|--------|----------|
| xs | 320px | iPhone SE, маленькие Android |
| sm | 375px | iPhone 12/13/14 |
| md | 428px | iPhone Pro Max, большие Android |
| lg | 768px | Tablet portrait, маленькие десктопы |
| xl | 1024px | Tablet landscape, десктоп |
| xxl | 1440px+ | Большие десктопы |

### 6.2 Mobile (320px — 428px)

- Стандартный layout, описанный выше
- Кнопки на всю ширину
- Single column

### 6.3 Tablet (768px — 1024px)

- Увеличенные отступы
- Max-width контента: 600px (центрировано)
- Иллюстрация может быть меньше (30-40%)

### 6.4 Desktop (1024px+)

```
┌────────────────────────────────────────────┐
│  [Stats]          TITLE          [History]│
├────────────────────────────────────────────┤
│                                            │
│   ┌──────────────┐   ┌──────────────────┐ │
│   │              │   │                  │ │
│   │ ИЛЛЮСТРАЦИЯ │   │  EVENT TEXT       │ │
│   │   (50%)     │   │   (50%)          │ │
│   │              │   │                  │ │
│   └──────────────┘   └──────────────────┘ │
│                                            │
├────────────────────────────────────────────┤
│   [Choice 1]  [Choice 2]  [Choice 3]       │
└────────────────────────────────────────────┘
```

- Layout: 2 колонки (иллюстрация + текст)
- Max-width контента: 900px
- Padding: xl (32px)

---

## 7. Figma Prompts

### 7.1 Главный промт для Figma (Mobile Gameplay Screen)

```
Mobile game UI, main gameplay screen for military survival game "Армейка",
layout with status bar at top showing 5 stats (STR, END, AGI, MOR, DISC),
center area featuring stylized semi-realistic military mess hall illustration,
bottom section with 3 decision buttons in dark military palette,

STYLE:
- dark military theme (olive green #1A1F16, brown-gray #2D3527)
- clean vector-like shapes
- soft gradient shading
- high contrast lighting
- muted military palette (green, gray, brown)

TYPOGRAPHY:
- Roboto Condensed for headers (bold, uppercase)
- Inter/Roboto for body text
- Monospace for stats numbers

COMPONENTS:
- Compact stat bars with icons and values
- Large illustrated event area (50% screen height)
- Full-width choice buttons (min 44px tap targets)
- Clean card-based layout

consistency with military survival game aesthetic,
clean grid layout, professional mobile UI,
high usability, minimalistic, pixel-perfect
```

### 7.2 Промт для Start Screen

```
Mobile game UI, start screen / main menu for military survival game "АРМЕЙКА",
centered layout with game logo/title at top,
dark military background with subtle texture,
primary CTA button "Начать игру" in accent green,
secondary button "Продолжить" shown below when save exists,

STYLE:
- dark military theme
- military green accent color (#4A5D3F)
- cinematic lighting
- stylized semi-realistic background illustration

clean layout, mobile-first design,
high contrast text, readable typography,
professional mobile game UI, pixel-perfect
```

### 7.3 Промт для Stats Panel

```
Mobile game UI, stats panel modal overlay for military survival game,
dark surface color (#2D3527) with subtle shadow,
header with close button,
section showing formal military rank "Рядовой",
grid of 5 stat cards (STR, END, AGI, MOR, DISC) with icons and progress bars,
section showing informal status "Запах",
progress bar showing "День X из 30",

STYLE:
- dark military theme
- stat bars with color coding (green/yellow/red based on value)
- clean card components with rounded corners (8px)

TYPOGRAPHY:
- Roboto Condensed for headers
- Roboto Mono for numbers

modal overlay style, slide-up animation ready,
professional mobile UI, clean layout, pixel-perfect
```

### 7.4 Промт для History Panel

```
Mobile game UI, event history slide-up panel for military survival game,
dark surface background (#2D3527),
header "ИСТОРИЯ СОБЫТИЙ",
list of 10 recent events showing:
- day number
- event name
- player choice (marked with >)
- stat changes with +/- indicators
- success/failure icons

STYLE:
- dark military theme
- list items with subtle borders
- green indicators for positive changes
- red indicators for negative changes

scrollable list ready,
clean typography, professional mobile UI,
pixel-perfect, consistent with game aesthetic
```

### 7.5 Промт для Final Screen

```
Mobile game UI, game over / final results screen for military survival game "Армейка",
centered layout with "РЕЗУЛЬТАТ" header,

Three possible final types with Russian descriptions:
1. "ТИХИЙ ДЕМБЕЛЬ" (Приспособленец) - green accent color
   Description: "Ваша служба прошла неспеша, вы не выделялись, не косячили, выполняли приказы 'шакалов' и офицеров. В общем — вы просто жили."

2. "УВАЖАЕМЫЙ ДЕМБЕЛЬ" - military green accent
   Description: "Вы ушли со службы в золотом кашне. Вся рота провожала вас со слезами на глазах, и обсуждала, как вы смогли так 'подняться' за такой короткий срок. А ваш кореш — старший прапорщик Залупенко плакал от горя, ведь вы смогли 'загнать' 300 тонн гнутых гвоздей..."

3. "СЛОМАННЫЙ ДЕМБЕЛЬ" (Косячник) - red accent color
   Description: "Служба в армии принесла вам лишь горе, затрещины и сломанный нос, когда вы ночью упали с кровати, исполняя 'летучую мышь' по команде сержанта Златопузова. Впрочем, вы выжили — а это главное."

Statistics section showing:
- days served (Дней прослужено)
- final MOR value
- final DISC value and formal/informal rank

STYLE:
- dark military theme
- color-coded final type (green for quiet, dark green for respected, red for broken)
- atmospheric description text

primary CTA button "ИГРАТЬ СНОВА" at bottom,
professional mobile game UI, pixel-perfect
```

### 7.6 Промт для Desktop Gameplay

```
Desktop web game UI, gameplay screen for military survival game "Армейка",
two-column layout: illustration on left (50%), event text and choices on right (50%),
top bar with all stats visible, full-width below,

ILLUSTRATION:
- large stylized military environment (mess hall, barracks, training ground)
- clean vector art style with cinematic lighting
- muted military palette

RIGHT PANEL:
- event title and description (scrollable)
- 3-4 choice buttons in horizontal row or 2x2 grid

STYLE:
- dark military theme (#1A1F16 background)
- consistent with mobile version
- professional game UI aesthetic

responsive desktop layout, max-width 1200px centered,
professional web game UI, pixel-perfect
```

---

## 8. Анимации и переходы

### 8.1 Переходы между событиями

| Действие | Анимация | Длительность |
|----------|----------|--------------|
| Смена события | Fade-out → Fade-in | 300ms |
| Иллюстрация | Crossfade | 400ms |
| Текст | Slide-up + Fade | 300ms |

### 8.2 Взаимодействие

| Действие | Анимация | Длительность |
|----------|----------|--------------|
| Нажатие кнопки | Scale 0.98 | 100ms |
| Disabled кнопка | Opacity 0.5 | 200ms |
| Stat change | Number count-up + bar width | 400ms |

### 8.3 Модальные окна

| Действие | Анимация | Длительность |
|----------|----------|--------------|
| Открытие | Fade-in + Slide-up | 300ms |
| Закрытие | Fade-out + Slide-down | 200ms |

---

## 9. Технические требования

### 9.1 Accessibility

- Minimum touch target: 44x44px
- Color contrast: minimum 4.5:1 for text
- Semantic HTML для скринридеров
- Focus states для клавиатурной навигации

### 9.2 Performance

- First Contentful Paint: < 1.5s
- Time to Interactive: < 3s
- Images: lazy loading для иллюстраций
- Bundle size target: < 500KB gzipped

### 9.3 Безопасность

- XSS protection для пользовательского текста
- Sanitize input в GraphQL
- CSP headers

---

## 10. Глоссарий UI терминов

| Термин | Определение |
|--------|-------------|
| Stats Bar | Верхняя панель с иконками статов |
| Event Card | Область с иллюстрацией и описанием события |
| Choice Button | Кнопка выбора действия |
| Modal / Overlay | Модальное окно поверх основного экрана |
| Toast | Временное уведомление |
| Stat Bar | Визуальный индикатор значения (полоска) |
| Progress Bar | Индикатор прогресса (день X из 30) |
| Surface | Основной цвет фона карточек/панелей |

---

## 12. MASTER PROMPT для Figma

Ниже представлен готовый к использованию промт для генерации всех UI макетов в Figma. Промт основан на **реальных скриншотах Hoosegow: Prison Survival**:

---

```
You are a professional game UI/UX designer and Figma expert.

Your task is to design a complete mobile game UI for "Армейка" (Army Game) — a narrative survival choice-based RPG about military conscription in Russia.

## EXACT REFERENCE: Hoosegow: Prison Survival
Study the screenshots carefully and replicate this EXACT layout:

**MAIN GAMEPLAY SCREEN LAYOUT:**
┌─────────────────────────────────────┐
│ [=] [DAY 3] [📜]                   │  ← Top bar with menu, day, history
├─────────────────────────────────────┤
│ HEALTH  ████████░░  80             │  ← Stats row: icon + name + bar + value
│ HUNGER  ████░░░░░░  40             │
│ SANITY  ████████░░  78             │
│ RESPECT ██░░░░░░░░  23             │
├─────────────────────────────────────┤
│                                     │
│     ┌───────────────────────┐       │
│     │                       │       │
│     │    ILLUSTRATION      │       │  ← ~40% screen height
│     │    (cell scene)      │       │
│     │                       │       │
│     └───────────────────────┘       │
│                                     │
│  "You wake up in your cell.        │  ← Event description
│   The guard walks by..."            │
│                                     │
├─────────────────────────────────────┤
│ ┌─────────────────────────────────┐   │
│ │ SLEEP                           │   │
│ └─────────────────────────────────┘   │
│ ┌─────────────────────────────────┐   │
│ │ WAIT                            │   │  ← Choice buttons (3)
│ └─────────────────────────────────┘   │
│ ┌─────────────────────────────────┐   │
│ │ WORK                            │   │
│ └─────────────────────────────────┘   │
└─────────────────────────────────────┘

## KEY UI PATTERNS FROM HOOSEGOW

### 1. Stats Bar (TOP)
- 4 stats in vertical stack OR horizontal row
- Each stat: NAME + progress bar + numeric value
- Progress bar colors:
  - GREEN (#4CAF50): >60%
  - YELLOW (#FFC107): 30-60%
  - RED (#F44336): <30%
- Dark background (#1A1A1A - #2A2A2A)
- White text for values

### 2. Top Navigation Bar
- Left: Menu/hamburger icon or stats icon
- Center: "DAY X" (current turn)
- Right: History icon

### 3. Illustration Area
- ~40% of screen height
- Clean vector illustration
- Rounded corners (8px)
- Subtle border or shadow

### 4. Event Description
- Below illustration
- 16px font, white/light text
- Scrollable if long

### 5. Choice Buttons
- 3 buttons stacked vertically
- Full-width with padding (16px sides)
- Height: 48-56px
- Uppercase text
- Rounded corners (8px)
- Hover/pressed states

### 6. Stats Panel (Modal)
- Title: "CHARACTER"
- Grid of stat cards
- Each card: icon, name, progress bar, value
- Dark surface color

### 7. History Panel (Bottom Sheet)
- Title: "HISTORY"
- List format:
  - Day number
  - Event description
  - Choice made (prefixed with >)
  - Result with icon (✓ or ✗)
- Scrollable

### 8. Start Screen
- Game title (centered, bold)
- Dark gradient background
- 3 menu options (vertical stack):
  - CONTINUE (if save exists)
  - NEW GAME
  - SETTINGS (optional)

## GLOBAL STYLE
- Mobile game UI, optimized for thumb interaction
- Dark theme: background #1A1A1A to #2A2A2A
- Stats colors: #4CAF50 (green), #FFC107 (yellow), #F44336 (red)
- Clean, minimalist aesthetic
- No gradients on backgrounds (flat dark)
- Subtle shadows only

## TYPOGRAPHY
- Headers: Roboto Bold, 18-24px, uppercase
- Body: Roboto Regular, 14-16px
- Stats: Roboto Bold for numbers
- Buttons: Roboto Bold, uppercase, letter-spacing 1px

## COLOR PALETTE
| Element | Color | Hex |
|---------|-------|-----|
| Background | Dark gray | #1A1A1A |
| Surface | Lighter gray | #2A2A2A |
| Card | Gray | #333333 |
| Text Primary | White | #FFFFFF |
| Text Secondary | Gray | #9E9E9E |
| Stat Green | Green | #4CAF50 |
| Stat Yellow | Yellow | #FFC107 |
| Stat Red | Red | #F44336 |
| Button | Military green | #5D7A4A |
| Button Pressed | Darker green | #4A6339 |

## COMPONENTS TO CREATE
1. StatRow (name + progress bar + value)
2. StatCard (icon + name + bar + value for modal)
3. ChoiceButton (full-width, uppercase)
4. TopBar (menu + day + history icons)
5. IllustrationFrame (rounded, with shadow)
6. HistoryItem (day + event + choice + result)
7. ModalContainer (overlay + content)
8. BottomSheet (slide-up panel)

## SCREENS TO DESIGN

### 1. START SCREEN
- Game title centered
- 3 buttons: CONTINUE, NEW GAME, SETTINGS
- Dark gradient background

### 2. GAMEPLAY SCREEN
- Top bar with DAY X
- Stats: STR, END, AGI, MOR, DISC (5 stats)
- Illustration area
- Event description
- 3 choice buttons

### 3. STATS PANEL
- CHARACTER title
- 5 stat cards in grid
- Day progress

### 4. HISTORY PANEL
- HISTORY title
- List of 10 events

### 5. FINAL SCREEN
- Result title
- Final type (ТИХИЙ/УВАЖАЕМЫЙ/СЛОМАННЫЙ ДЕМБЕЛЬ)
- Stats summary
- PLAY AGAIN button

## OUTPUT REQUIREMENTS
- Use Figma auto-layout
- Create component variants (default, hover, pressed, disabled)
- Name components properly (e.g., "Button/Choice", "Stat/Row")
- Set up Figma variables for colors and typography
- Design for 375px width (iPhone standard)
- Minimum touch target: 44x44px

Generate the complete UI kit in Figma matching Hoosegow's exact layout and style.
```

---

## 13. Следующие шаги

1. **Создать Figma макеты** используя MASTER PROMPT из раздела 12
2. **Согласовать дизайн** с stakeholders
3. **Экспортировать assets** (если нужны иллюстрации)
4. **Начать разработку** Frontend компонентов

---

**Документ подготовлен:** 2026-03-25  
**Автор:** AI Assistant (ui-analyst, pm, analyst skills)
