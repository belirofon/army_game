-- ============================================================================
-- PLAYERS — Единственный источник правды о состоянии игрока
-- v2: Stats as separate columns with CHECK constraints
-- ============================================================================
CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    str INT NOT NULL DEFAULT 50 CHECK (str BETWEEN 1 AND 100),
    end_ INT NOT NULL DEFAULT 50 CHECK (end_ BETWEEN 1 AND 100),
    agi INT NOT NULL DEFAULT 50 CHECK (agi BETWEEN 1 AND 100),
    mor INT NOT NULL DEFAULT 50 CHECK (mor BETWEEN 0 AND 100),
    disc INT NOT NULL DEFAULT 0 CHECK (disc BETWEEN -100 AND 100),
    formal_rank VARCHAR(20) NOT NULL DEFAULT 'РЯДОВОЙ' 
        CHECK (formal_rank IN ('РЯДОВОЙ', 'ЕФРЕЙТОР', 'МЛ_СЕРЖАНТ', 'СЕРЖАНТ')),
    informal_status VARCHAR(20) NOT NULL DEFAULT 'ЗАПАХ'
        CHECK (informal_status IN ('ЗАПАХ', 'ДУХ', 'СЛОН', 'ЧЕРПАК', 'ДЕД', 'ДЕМБЕЛЬ')),
    turn INT NOT NULL DEFAULT 1 CHECK (turn BETWEEN 1 AND 24),
    flags JSONB DEFAULT '[]',
    version INT NOT NULL DEFAULT 1,
    is_finished BOOLEAN NOT NULL DEFAULT FALSE,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_players_id ON players(id);
CREATE INDEX IF NOT EXISTS idx_players_turn ON players(turn);
CREATE INDEX IF NOT EXISTS idx_players_finished ON players(is_finished) WHERE is_finished = FALSE;
CREATE INDEX IF NOT EXISTS idx_players_updated_at ON players(updated_at DESC);

-- ============================================================================
-- GAME_LOGS — Неизменяемая история всех выборов
-- ============================================================================
CREATE TABLE IF NOT EXISTS game_logs (
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

CREATE INDEX IF NOT EXISTS idx_game_logs_player_id ON game_logs(player_id);
CREATE INDEX IF NOT EXISTS idx_game_logs_player_turn ON game_logs(player_id, turn DESC);
CREATE INDEX IF NOT EXISTS idx_game_logs_created_at ON game_logs(created_at DESC);
