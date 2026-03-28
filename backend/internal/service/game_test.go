package service

import (
	"context"
	"testing"

	"army-game-backend/internal/domain"
)

type mockPlayerRepo struct {
	players map[string]*domain.Player
}

func (m *mockPlayerRepo) Create(ctx context.Context, player *domain.Player) error {
	m.players[player.ID] = player
	return nil
}

func (m *mockPlayerRepo) GetByID(ctx context.Context, id string) (*domain.Player, error) {
	if p, ok := m.players[id]; ok {
		return p, nil
	}
	return nil, ErrPlayerNotFound
}

func (m *mockPlayerRepo) Update(ctx context.Context, player *domain.Player) error {
	m.players[player.ID] = player
	return nil
}

func (m *mockPlayerRepo) Delete(ctx context.Context, id string) error {
	delete(m.players, id)
	return nil
}

type mockHistoryRepo struct {
	entries map[string][]domain.GameLogEntry
}

func (m *mockHistoryRepo) CreateEntry(ctx context.Context, entry *domain.GameLogEntry) error {
	m.entries[entry.PlayerID] = append([]domain.GameLogEntry{*entry}, m.entries[entry.PlayerID]...)
	if len(m.entries[entry.PlayerID]) > 10 {
		m.entries[entry.PlayerID] = m.entries[entry.PlayerID][:10]
	}
	return nil
}

func (m *mockHistoryRepo) GetByPlayerID(ctx context.Context, playerID string, limit int) ([]domain.GameLogEntry, error) {
	if entries, ok := m.entries[playerID]; ok {
		if limit > 0 && len(entries) > limit {
			return entries[:limit], nil
		}
		return entries, nil
	}
	return []domain.GameLogEntry{}, nil
}

func (m *mockHistoryRepo) DeleteByPlayerID(ctx context.Context, playerID string) error {
	delete(m.entries, playerID)
	return nil
}

func NewTestService() *GameService {
	playerRepo := &mockPlayerRepo{players: make(map[string]*domain.Player)}
	historyRepo := &mockHistoryRepo{entries: make(map[string][]domain.GameLogEntry)}
	return NewGameService(playerRepo, historyRepo)
}

func TestStartGame(t *testing.T) {
	svc := NewTestService()

	result, err := svc.StartGame(context.Background())
	if err != nil {
		t.Fatalf("StartGame failed: %v", err)
	}

	if result.GameID == "" {
		t.Error("GameID should not be empty")
	}

	if result.Player == nil {
		t.Error("Player should not be nil")
	}

	if result.Player.Turn != 1 {
		t.Errorf("Expected turn 1, got %d", result.Player.Turn)
	}

	if result.Player.Stats.Str != 50 {
		t.Errorf("Expected str 50, got %d", result.Player.Stats.Str)
	}

	if result.Player.Stats.Mor != 50 {
		t.Errorf("Expected mor 50, got %d", result.Player.Stats.Mor)
	}

	if result.CurrentEvent == nil {
		t.Error("CurrentEvent should not be nil")
	}

	if result.IsGameOver {
		t.Error("Game should not be over on start")
	}
}

func TestMakeChoice(t *testing.T) {
	svc := NewTestService()

	startResult, _ := svc.StartGame(context.Background())
	gameID := startResult.GameID

	choiceResult, err := svc.MakeChoice(context.Background(), gameID, "choice_1", 1)
	if err != nil {
		t.Fatalf("MakeChoice failed: %v", err)
	}

	if !choiceResult.Success {
		t.Error("Choice should succeed")
	}

	if choiceResult.UpdatedPlayer.Turn != 2 {
		t.Errorf("Expected turn 2, got %d", choiceResult.UpdatedPlayer.Turn)
	}

	if choiceResult.UpdatedPlayer.Version != 2 {
		t.Errorf("Expected version 2, got %d", choiceResult.UpdatedPlayer.Version)
	}

	if choiceResult.NewVersion != 2 {
		t.Errorf("Expected newVersion 2, got %d", choiceResult.NewVersion)
	}

	if choiceResult.GameOver {
		t.Error("Game should not be over yet")
	}

	if choiceResult.NextEvent == nil {
		t.Error("NextEvent should be provided")
	}
}

func TestMakeChoiceInvalidChoice(t *testing.T) {
	svc := NewTestService()

	startResult, _ := svc.StartGame(context.Background())
	gameID := startResult.GameID

	_, err := svc.MakeChoice(context.Background(), gameID, "invalid_choice", 1)
	if err != ErrInvalidChoice {
		t.Errorf("Expected ErrInvalidChoice, got %v", err)
	}
}

func TestMakeChoiceConcurrentModification(t *testing.T) {
	svc := NewTestService()

	startResult, _ := svc.StartGame(context.Background())
	gameID := startResult.GameID

	_, err := svc.MakeChoice(context.Background(), gameID, "choice_1", 999)
	if err != ErrConcurrentModify {
		t.Errorf("Expected ErrConcurrentModify, got %v", err)
	}
}

func TestGameOverAtMaxTurn(t *testing.T) {
	svc := NewTestService()

	startResult, _ := svc.StartGame(context.Background())
	gameID := startResult.GameID

	for i := 0; i < 29; i++ {
		_, err := svc.MakeChoice(context.Background(), gameID, "choice_1", i+1)
		if err != nil {
			t.Fatalf("MakeChoice failed: %v", err)
		}
	}

	result, err := svc.LoadGame(context.Background(), gameID)
	if err != nil {
		t.Fatalf("LoadGame failed: %v", err)
	}

	if !result.IsGameOver {
		t.Error("Game should be over at turn 30")
	}

	if result.Final == nil {
		t.Error("Final should be provided")
	}
}

func TestLoadGame(t *testing.T) {
	svc := NewTestService()

	startResult, _ := svc.StartGame(context.Background())
	gameID := startResult.GameID

	svc.MakeChoice(context.Background(), gameID, "choice_1", 1)
	svc.MakeChoice(context.Background(), gameID, "choice_1", 2)

	result, err := svc.LoadGame(context.Background(), gameID)
	if err != nil {
		t.Fatalf("LoadGame failed: %v", err)
	}

	if result.Player.Turn != 3 {
		t.Errorf("Expected turn 3, got %d", result.Player.Turn)
	}

	if len(result.EventHistory) != 2 {
		t.Errorf("Expected 2 history entries, got %d", len(result.EventHistory))
	}
}

func TestRestartGame(t *testing.T) {
	svc := NewTestService()

	startResult, _ := svc.StartGame(context.Background())
	gameID := startResult.GameID

	svc.MakeChoice(context.Background(), gameID, "choice_1", 1)

	restartResult, err := svc.RestartGame(context.Background(), gameID)
	if err != nil {
		t.Fatalf("RestartGame failed: %v", err)
	}

	if restartResult.Player.Turn != 1 {
		t.Errorf("Expected turn 1 after restart, got %d", restartResult.Player.Turn)
	}

	if restartResult.Player.Stats.Mor != 50 {
		t.Errorf("Expected mor 50 after restart, got %d", restartResult.Player.Stats.Mor)
	}
}

func TestFormalRankCalculation(t *testing.T) {
	svc := NewTestService()

	tests := []struct {
		disc     int
		wantRank domain.FormalRank
	}{
		{0, domain.FormalRankRyadovoy},
		{24, domain.FormalRankRyadovoy},
		{25, domain.FormalRankEfreitor},
		{49, domain.FormalRankEfreitor},
		{50, domain.FormalRankMlSerzhant},
		{74, domain.FormalRankMlSerzhant},
		{75, domain.FormalRankSerzhant},
		{100, domain.FormalRankSerzhant},
	}

	for _, tt := range tests {
		player := &domain.Player{
			Stats: domain.PlayerStats{Disc: tt.disc},
		}
		svc.updateRanks(player)
		if player.FormalRank != tt.wantRank {
			t.Errorf("disc=%d: expected rank %s, got %s", tt.disc, tt.wantRank, player.FormalRank)
		}
	}
}

func TestInformalStatusCalculation(t *testing.T) {
	svc := NewTestService()

	tests := []struct {
		disc       int
		wantStatus domain.InformalStatus
	}{
		{0, domain.InformalStatusZapah},
		{-1, domain.InformalStatusDukh},
		{-49, domain.InformalStatusDukh},
		{-50, domain.InformalStatusSlon},
		{-74, domain.InformalStatusSlon},
		{-75, domain.InformalStatusCherpak},
		{-89, domain.InformalStatusCherpak},
		{-90, domain.InformalStatusDed},
		{-99, domain.InformalStatusDed},
		{-100, domain.InformalStatusDembel},
		{-101, domain.InformalStatusDembel},
	}

	for _, tt := range tests {
		player := &domain.Player{
			Stats: domain.PlayerStats{Disc: tt.disc},
		}
		svc.updateRanks(player)
		if player.InformalStatus != tt.wantStatus {
			t.Errorf("disc=%d: expected status %s, got %s", tt.disc, tt.wantStatus, player.InformalStatus)
		}
	}
}

func TestStatClamping(t *testing.T) {
	player := &domain.Player{
		Stats: domain.PlayerStats{Str: 50, End: 50, Agi: 50, Mor: 50, Disc: 0},
	}

	effects := []domain.Effect{
		{Stat: "str", Delta: 100, PreviousValue: 50, NewValue: 150},
		{Stat: "mor", Delta: -100, PreviousValue: 50, NewValue: -50},
		{Stat: "disc", Delta: -200, PreviousValue: 0, NewValue: -200},
	}

	svc := NewTestService()

	for _, effect := range effects {
		svc.applyEffect(&player.Stats, effect)
	}

	if player.Stats.Str != 100 {
		t.Errorf("Expected str 100 (clamped), got %d", player.Stats.Str)
	}

	if player.Stats.Mor != 0 {
		t.Errorf("Expected mor 0 (clamped), got %d", player.Stats.Mor)
	}

	if player.Stats.Disc != -100 {
		t.Errorf("Expected disc -100 (clamped), got %d", player.Stats.Disc)
	}
}
