# Army Game — AI Prompts (Characters & Environment)

## 1. Environment Prompt (Base Style)

### Prompt

```
mobile game environment, russian military mess hall cafeteria, stylized semi-realistic illustration, clean vector-like shapes, soft gradient shading, high contrast lighting, cinematic rim light, slightly exaggerated perspective, polished game art, high detail, consistent game style
```

### Explanation

* **mobile game environment** — оптимизация под мобильный UX
* **stylized semi-realistic** — баланс между реализмом и читаемостью
* **vector-like shapes** — чистые формы для читаемости
* **cinematic rim light** — добавляет глубину и драму
* **consistent game style** — важно для пайплайна

---

## 2. Composition Prompt

### Prompt

```
wide interior composition, centered perspective with strong depth lines, camera slightly above table height, clear UI space at bottom, readable layout for gameplay
```

### Explanation

* **centered perspective** — ключевой паттерн Hoosegow
* **strong depth lines** — создает ощущение глубины
* **camera slightly above** — улучшает обзор сцены
* **UI space at bottom** — место под кнопки

---

## 3. Mess Hall Environment Details

### Prompt

```
large military dining hall, long rows of metal tables and benches, trays with simple army food, industrial ceiling lights casting cold and warm mixed lighting, steam rising from food, slightly crowded atmosphere with recruits eating, institutional and repetitive environment, clean but rough military aesthetic, muted green and gray tones
```

### Explanation

* **long rows of tables** — повторяемость (важно для атмосферы)
* **mixed lighting** — визуальный интерес
* **steam** — добавляет «живость» сцене
* **muted colors** — военная эстетика

---

## 4. MASTER PROMPT (Full Pipeline)

### Prompt

```
You are a professional game UI/UX designer and concept artist.

Your task is to generate a complete set of mobile game screens in the style of Hoosegow: Prison Survival, but for a military survival game (army setting).

STEP 1 — DEFINE GLOBAL STYLE
- mobile game UI
- stylized semi-realistic illustration
- clean vector-like shapes
- soft gradient shading
- high contrast lighting
- cinematic rim light
- muted military palette (green, gray, brown)

STEP 2 — DEFINE CAMERA
- centered composition
- strong depth
- camera slightly above
- bottom space for UI

STEP 3 — SCREENS
1. Mess hall
2. Event screen
3. Barracks
4. Inventory
5. Character
6. Daily routine

For each screen:
- generate image prompt
- describe UI layout

Ensure consistency and readability.
```

### Explanation

* Используется как **system prompt**
* Генерирует сразу все экраны
* Контролирует стиль

---

## 5. UI Screen Prompt (Gameplay)

### Prompt

```
Design a mobile game UI screen for a survival choice-based game similar to Hoosegow: Prison Survival, but in a military setting.

SCREEN TYPE: Main gameplay screen (mess hall)

LAYOUT:
- Top bar with player stats (health, hunger, morale, authority)
- Center: large environment illustration
- Bottom: 3 decision buttons

STYLE:
- stylized semi-realistic
- clean UI
- dark military palette

UX:
- clear focus on choices
- large tap areas
- minimal clutter

OUTPUT:
- structured UI
- auto-layout
- reusable components
```

### Explanation

* Генерация UI в Figma
* Подходит для Fluvara / Figma AI

---

## 6. Example Final Image Prompt (Mess Hall)

### Prompt

```
mobile game environment, military mess hall, long metal tables, soldiers eating, trays with food, steam rising, industrial lighting, green-gray palette, stylized semi-realistic, cinematic lighting, centered composition, strong depth perspective, camera slightly above table level, clean vector-like shapes, high detail, UI space at bottom
```

### Explanation

* Финальный prompt для генерации изображения
* Объединяет стиль + композицию + детали

---

## 7. Architecture Insight

```
BDD → Scene → Prompt → Image → UI → Game
```

### Explanation

* Единый pipeline
* Можно автоматизировать генерацию контента
* Масштабируется под новые сцены

---

## Notes

* Все промты можно комбинировать
* Лучше разбивать на блоки (style / composition / details)
* Использовать как DSL для генерации
