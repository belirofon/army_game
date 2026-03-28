package domain

import "context"

type PlayerRepository interface {
	Create(ctx context.Context, player *Player) error
	GetByID(ctx context.Context, id string) (*Player, error)
	Update(ctx context.Context, player *Player) error
	Delete(ctx context.Context, id string) error
}

type EventHistoryRepository interface {
	CreateEntry(ctx context.Context, entry *GameLogEntry) error
	GetByPlayerID(ctx context.Context, playerID string, limit int) ([]GameLogEntry, error)
	DeleteByPlayerID(ctx context.Context, playerID string) error
}
