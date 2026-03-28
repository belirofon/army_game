package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"army-game-backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// tableName constants
const (
	tablePlayers  = "players"
	tableGameLogs = "game_logs"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, player *domain.Player) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO players (id, str, end_, agi, mor, disc, formal_rank, informal_status, turn, flags, is_finished, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, player.ID, player.Stats.Str, player.Stats.End, player.Stats.Agi, player.Stats.Mor, player.Stats.Disc,
		player.FormalRank, player.InformalStatus, player.Turn, player.Flags, player.IsFinished, player.Version, player.CreatedAt, player.UpdatedAt)

	if err != nil {
		return fmt.Errorf("insert player: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Player, error) {
	var player domain.Player
	var formalRank string
	var informalStatus string

	err := r.pool.QueryRow(ctx, `
		SELECT id, str, end_, agi, mor, disc, formal_rank, informal_status, turn, flags, is_finished, finished_at, version, created_at, updated_at
		FROM players
		WHERE id = $1
	`, id).Scan(
		&player.ID,
		&player.Stats.Str,
		&player.Stats.End,
		&player.Stats.Agi,
		&player.Stats.Mor,
		&player.Stats.Disc,
		&formalRank,
		&informalStatus,
		&player.Turn,
		&player.Flags,
		&player.IsFinished,
		&player.FinishedAt,
		&player.Version,
		&player.CreatedAt,
		&player.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("player not found")
		}
		return nil, fmt.Errorf("get player: %w", err)
	}

	player.FormalRank = domain.FormalRank(formalRank)
	player.InformalStatus = domain.InformalStatus(informalStatus)

	return &player, nil
}

func (r *PostgresRepository) Update(ctx context.Context, player *domain.Player) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE players
		SET str = $2, end_ = $3, agi = $4, mor = $5, disc = $6, 
			formal_rank = $7, informal_status = $8, turn = $9, flags = $10, 
			is_finished = $11, finished_at = $12, version = $13, updated_at = $14
		WHERE id = $1
	`, player.ID, player.Stats.Str, player.Stats.End, player.Stats.Agi, player.Stats.Mor, player.Stats.Disc,
		player.FormalRank, player.InformalStatus, player.Turn, player.Flags,
		player.IsFinished, player.FinishedAt, player.Version, player.UpdatedAt)

	if err != nil {
		return fmt.Errorf("update player: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM players WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete player: %w", err)
	}
	return nil
}

func (r *PostgresRepository) CreateEntry(ctx context.Context, entry *domain.GameLogEntry) error {
	checkResultJSON, err := json.Marshal(entry.CheckResult)
	if err != nil {
		return fmt.Errorf("marshal check result: %w", err)
	}

	effectsJSON, err := json.Marshal(entry.Effects)
	if err != nil {
		return fmt.Errorf("marshal effects: %w", err)
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO game_logs (id, player_id, turn, event_template_id, event_description, choice_id, choice_text, check_result, effects, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, entry.ID, entry.PlayerID, entry.Turn, entry.EventTemplateID, entry.EventDescription, "", entry.ChoiceText, checkResultJSON, effectsJSON, entry.CreatedAt)

	if err != nil {
		return fmt.Errorf("insert log entry: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByPlayerID(ctx context.Context, playerID string, limit int) ([]domain.GameLogEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, player_id, turn, event_template_id, event_description, choice_id, choice_text, check_result, effects, created_at
		FROM game_logs
		WHERE player_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, playerID, limit)

	if err != nil {
		return nil, fmt.Errorf("query history: %w", err)
	}
	defer rows.Close()

	var entries []domain.GameLogEntry

	for rows.Next() {
		var entry domain.GameLogEntry
		var checkResultJSON []byte
		var effectsJSON []byte
		var eventTemplateID, choiceID string

		err := rows.Scan(
			&entry.ID,
			&entry.PlayerID,
			&entry.Turn,
			&eventTemplateID,
			&entry.EventDescription,
			&choiceID,
			&entry.ChoiceText,
			&checkResultJSON,
			&effectsJSON,
			&entry.CreatedAt,
		)

		entry.EventTemplateID = eventTemplateID
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		if err := json.Unmarshal(checkResultJSON, &entry.CheckResult); err != nil {
			return nil, fmt.Errorf("unmarshal check result: %w", err)
		}

		if err := json.Unmarshal(effectsJSON, &entry.Effects); err != nil {
			return nil, fmt.Errorf("unmarshal effects: %w", err)
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return entries, nil
}

func (r *PostgresRepository) DeleteByPlayerID(ctx context.Context, playerID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM game_logs WHERE player_id = $1`, playerID)
	if err != nil {
		return fmt.Errorf("delete history: %w", err)
	}
	return nil
}
