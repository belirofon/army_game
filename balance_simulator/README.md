# Balance Simulator

Monte Carlo симулятор для проверки баланса игры «Армейка» перед implementation.

## Быстрый старт

```bash
cd balance_simulator

# Запустить 1000 симуляций
python simulator.py --runs 1000

# С воспроизводимым seed
python simulator.py --runs 1000 --seed 42

# Подробный вывод
python simulator.py --runs 100 --verbose

# Сохранить результаты
python simulator.py --runs 1000 --save results/run1.csv
```

## Опции

| Опция | Описание | Default |
|-------|----------|---------|
| `--runs`, `-n` | Количество симуляций | 1000 |
| `--seed`, `-s` | Random seed для reproducibility | None |
| `--events`, `-e` | Путь к events.json | events.json |
| `--output`, `-o` | Формат вывода: terminal/csv/json | terminal |
| `--save` | Путь для сохранения CSV | None |
| `--verbose`, `-v` | Печатать каждый playthrough | False |

## Форматы вывода

### Terminal (default)

```
╔══════════════════════════════════════════════════════════════════════╗
║                    BALANCE SIMULATOR REPORT                            ║
╠══════════════════════════════════════════════════════════════════════╣
║  Runs: 1000    Deaths: 243    Victories: 757                        ║
╠══════════════════════════════════════════════════════════════════════╣
║  KEY METRICS                                                          ║
║  ──────────────────────────────────────────────────────────────────── ║
║  Death Rate:         24.3%  ████████░░░░░░░░░░░░░  target: 20-30%  ✓   ║
║  Victory Rate:        75.7%  ██████████████████░░░  target: 70-80%  ✓   ║
║  Perfect Runs:         8.2%  █████░░░░░░░░░░░░░░░░░  target: 5-10%   ✓   ║
╠══════════════════════════════════════════════════════════════════════╣
║  FINAL STATUS: BALANCED ✓                                             ║
╚══════════════════════════════════════════════════════════════════════╝
```

### CSV

```bash
python simulator.py --runs 1000 --output csv > results.csv
```

### JSON

```bash
python simulator.py --runs 1000 --output json
```

## Структура событий (events.json)

```json
{
  "id": "e001",
  "type": "ROUTINE",           // ROUTINE, SOCIAL, INSPECTION, INFORMAL, EMERGENCY, SAFE
  "difficulty": 1,            // 1-5 (влияет на probability)
  "choices": [
    {
      "id": "e001_c1",
      "probability": 0.7,       // Базовый шанс успеха
      "success_effects": {"mor": 2, "disc": 2},
      "partial_effects": {"mor": 0, "disc": 1},
      "failure_effects": {"mor": -3, "disc": -2}
    }
  ]
}
```

## Целевые метрики

| Метрика | Target | Описание |
|---------|--------|---------|
| Death Rate | 20-35% | % игроков погибших до victory |
| Victory Rate | 65-80% | % игроков дошедших до конца |
| Perfect Runs | 0-5% | % игроков с MOR > 50 и DISC ∈ [-20, +20] |
| Avg MOR/turn | -2.0 to -1.2 | Среднее изменение MOR за ход |

**Примечание:** Perfect runs крайне редки (~0.1%) при текущем балансе. Это ожидаемо для challenging gameplay.

## Как использовать

1. **Настройка событий** — отредактируй `events.json`
2. **Запуск симуляций** — `python simulator.py --runs 1000`
3. **Проверка метрик** — смотри output на ✓/✗
4. **Корректировка** — меняй effects/probability в events.json
5. **Повтор** — запусти снова для проверки

## Правила симуляции

- **Difficulty Curve**: ходы 1-10 легче, 21-30 сложнее
- **MOR Recovery**: после 3 MOR- событий — гарантированно SAFE событие
- **Weighted Selection**: события выбираются по весам, не pure random
- **Stat Clamping**: все статы ограничены 0-100 (для MOR/disc: -100 до 100)

## Расширение

### Добавление новых событий

Добавь в `events.json`:

```json
{
  "id": "e021",
  "type": "INFORMAL",
  "difficulty": 3,
  "choices": [
    {
      "id": "e021_c1",
      "probability": 0.5,
      "success_effects": {"mor": 3, "disc": -10},
      "partial_effects": {"mor": 0, "disc": -5},
      "failure_effects": {"mor": -4, "disc": -8}
    }
  ]
}
```

### Изменение целевых метрик

Измени `StatisticsAnalyzer.TARGETS` в `simulator.py`:

```python
TARGETS = {
    'death_rate': (0.15, 0.25),  # Изменено
    ...
}
```

## Troubleshooting

**Death rate слишком высокий (>40%)?**
→ Увеличь `probability` в событиях или уменьши негативные effects

**Victory rate слишком низкий (<60%)?**
→ Добавь больше SAFE событий или увеличь MOR gains

**Perfect runs слишком много (>20%)?**
→ Игра слишком легкая — увеличь сложность поздних событий

## История результатов

Результаты сохраняются в папку `results/`:

| Файл | Симуляций | Death Rate | Victory Rate | MOR/turn |
|------|-----------|------------|--------------|----------|
| `run_baseline.csv` | 1000 | 30.3% | 69.7% | -1.40 |
| `run_5000.csv` | 5000 | 31.4% | 68.6% | -1.41 |

Используй `--save results/filename.csv` для сохранения новых результатов.
