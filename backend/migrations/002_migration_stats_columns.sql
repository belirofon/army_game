-- ============================================================================
-- Миграция: Stats из JSONB в отдельные колонки
-- Выполнить ПОСЛЕ применения 001_schema.sql
-- ============================================================================

-- Если таблицы уже существуют в старом формате, выполнить миграцию:

BEGIN;

-- Добавить колонки статов в существующую таблицу players
ALTER TABLE players ADD COLUMN IF NOT EXISTS str INT;
ALTER TABLE players ADD COLUMN IF NOT EXISTS end_ INT;
ALTER TABLE players ADD COLUMN IF NOT EXISTS agi INT;
ALTER TABLE players ADD COLUMN IF NOT EXISTS mor INT;
ALTER TABLE players ADD COLUMN IF NOT EXISTS disc INT;

-- Мигрировать данные из JSONB
UPDATE players SET 
    str = COALESCE((stats->>'str')::int, 50),
    end_ = COALESCE((stats->>'end')::int, 50),
    agi = COALESCE((stats->>'agi')::int, 50),
    mor = COALESCE((stats->>'mor')::int, 50),
    disc = COALESCE((stats->>'disc')::int, 0)
WHERE stats IS NOT NULL;

-- Установить значения по умолчанию для новых записей
ALTER TABLE players ALTER COLUMN str SET DEFAULT 50;
ALTER TABLE players ALTER COLUMN end_ SET DEFAULT 50;
ALTER TABLE players ALTER COLUMN agi SET DEFAULT 50;
ALTER TABLE players ALTER COLUMN mor SET DEFAULT 50;
ALTER TABLE players ALTER COLUMN disc SET DEFAULT 0;

-- Добавить NOT NULL после миграции
ALTER TABLE players ALTER COLUMN str SET NOT NULL;
ALTER TABLE players ALTER COLUMN end_ SET NOT NULL;
ALTER TABLE players ALTER COLUMN agi SET NOT NULL;
ALTER TABLE players ALTER COLUMN mor SET NOT NULL;
ALTER TABLE players ALTER COLUMN disc SET NOT NULL;

-- Добавить CHECK constraints
ALTER TABLE players ADD CHECK (str BETWEEN 1 AND 100);
ALTER TABLE players ADD CHECK (end_ BETWEEN 1 AND 100);
ALTER TABLE players ADD CHECK (agi BETWEEN 1 AND 100);
ALTER TABLE players ADD CHECK (mor BETWEEN 0 AND 100);
ALTER TABLE players ADD CHECK (disc BETWEEN -100 AND 100);

-- Переименовать event_history в game_logs если существует
ALTER TABLE IF EXISTS event_history RENAME TO game_logs;

-- Обновить turn constraint если уже есть данные
ALTER TABLE players DROP CONSTRAINT IF EXISTS players_turn_check;
ALTER TABLE players ADD CHECK (turn BETWEEN 1 AND 30);

COMMIT;
