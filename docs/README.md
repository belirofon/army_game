# Documentation

## Tech Specifications

| File | Description |
|------|-------------|
| `TZ.md` | Technical specification (TZ) — game mechanics, API contracts, DB schema, balance requirements |
| `architecture.md` | System architecture, DDD layers, module breakdown |
| `DB_ANALYSIS.md` | Database schema analysis with recommendations |
| `AGENTS.md` | Developer guidelines, code style, build commands |

## UI/UX Design

| File | Description |
|------|-------------|
| `DESIGN_SPECIFICATION.md` | Complete UI/UX design system, screens, components, Figma prompts |
| `Army_Game_Prompts.md` | AI prompts for generating game illustrations |

## Balance Simulator

| File | Description |
|------|-------------|
| `balance_simulator/simulator.py` | Monte Carlo simulation engine |
| `balance_simulator/events.json` | Game events pool (42 narrative events) |
| `balance_simulator/results/` | Simulation run results |

### Running Balance Tests

```bash
cd balance_simulator
python simulator.py --runs 1000 --seed 42
```

Target metrics: Death 20-35%, Victory 65-80%, Avg MOR/turn -2.0 to -1.5
