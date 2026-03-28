package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"army-game-backend/internal/domain"

	"github.com/google/uuid"
)

var (
	ErrPlayerNotFound   = errors.New("player not found")
	ErrGameNotFound     = errors.New("game not found")
	ErrInvalidChoice    = errors.New("invalid choice")
	ErrConcurrentModify = errors.New("concurrent modification")
	ErrNoEvents         = errors.New("no events loaded")
)

const MaxTurn = 24 // 24 months = 24 events

type GameService struct {
	playerRepo     domain.PlayerRepository
	historyRepo    domain.EventHistoryRepository
	rng            *rand.Rand
	eventTemplates []domain.EventTemplate
}

type playerState struct {
	recentEvents      []string
	consecutiveNegMOR int
}

func (s *GameService) getDifficultyModifier(turn int) float64 {
	switch {
	case turn <= 10:
		return 0.6
	case turn <= 20:
		return 0.85
	default:
		return 1.05
	}
}

func (s *GameService) calculateEventWeights(player *domain.Player, recentEvents []string) map[string]float64 {
	weights := make(map[string]float64)

	for _, tmpl := range s.eventTemplates {
		weight := 1.0

		isRecent := false
		for _, recentID := range recentEvents {
			if recentID == tmpl.ID {
				isRecent = true
				break
			}
		}
		if isRecent {
			weight = 0
		}

		if tmpl.Type == "SAFE" && s.hasRecentNegativeMOR(player) {
			weight += 0.3
		}
		if tmpl.Type == "INSPECTION" && player.Stats.Disc > 50 {
			weight += 0.2
		}
		if tmpl.Type == "INFORMAL" && player.Stats.Disc < -50 {
			weight += 0.2
		}
		if tmpl.Type == "EMERGENCY" && player.Turn >= 20 {
			weight += 0.2
		}

		weights[tmpl.ID] = weight
	}

	return weights
}

func (s *GameService) hasRecentNegativeMOR(player *domain.Player) bool {
	return false
}

func (s *GameService) selectWeightedEvent(weights map[string]float64) *domain.EventTemplate {
	var totalWeight float64
	for _, w := range weights {
		totalWeight += w
	}

	if totalWeight <= 0 {
		return nil
	}

	roll := s.rng.Float64() * totalWeight
	var cumulative float64

	for _, tmpl := range s.eventTemplates {
		cumulative += weights[tmpl.ID]
		if roll <= cumulative {
			return &tmpl
		}
	}

	return nil
}

func NewGameService(playerRepo domain.PlayerRepository, historyRepo domain.EventHistoryRepository) *GameService {
	svc := &GameService{
		playerRepo:  playerRepo,
		historyRepo: historyRepo,
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	svc.loadEvents()
	return svc
}

func (s *GameService) loadEvents() {
	// Try multiple paths for events.json
	paths := []string{
		"data/events.json",
		"../data/events.json",
		"backend/data/events.json",
	}

	var data []byte
	var err error

	for _, path := range paths {
		data, err = os.ReadFile(path)
		if err == nil {
			fmt.Printf("Loaded events from: %s\n", path)
			break
		}
	}

	if data == nil {
		fmt.Printf("Warning: could not load events file from any path: %v\n", err)
		s.eventTemplates = s.getDefaultEvents()
		return
	}

	var eventsData domain.EventsData
	if err := json.Unmarshal(data, &eventsData); err != nil {
		fmt.Printf("Warning: could not parse events JSON: %v\n", err)
		s.eventTemplates = s.getDefaultEvents()
		return
	}

	s.eventTemplates = eventsData.Events
	fmt.Printf("Loaded %d events from file\n", len(s.eventTemplates))
}

func (s *GameService) getDefaultEvents() []domain.EventTemplate {
	return []domain.EventTemplate{
		{
			ID:          "morning_formation",
			Description: "Утренняя построение. Сержант проверяет солдат.",
			Choices: []domain.Choice{
				{ID: "choice_1", Text: "Стоять смирно", Available: true},
				{ID: "choice_2", Text: "Поправить ремень", Available: true},
				{ID: "choice_3", Text: "Зевнуть", Available: true},
			},
		},
	}
}

type StartGameResult struct {
	GameID       string
	Player       *domain.Player
	CurrentEvent *domain.EventInstance
	IsGameOver   bool
	Final        *domain.Final
}

func (s *GameService) StartGame(ctx context.Context) (*StartGameResult, error) {
	gameID := uuid.New().String()
	player := &domain.Player{
		ID:             gameID,
		Stats:          domain.PlayerStats{Str: 50, End: 50, Agi: 50, Mor: 50, Disc: 0},
		FormalRank:     domain.FormalRankRyadovoy,
		InformalStatus: domain.InformalStatusZapah,
		Turn:           1,
		Flags:          []string{},
		IsFinished:     false,
		Version:        1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.playerRepo.Create(ctx, player); err != nil {
		return nil, fmt.Errorf("create player: %w", err)
	}

	event, err := s.generateEvent(ctx, player)
	if err != nil {
		return nil, fmt.Errorf("generate event: %w", err)
	}

	return &StartGameResult{
		GameID:       gameID,
		Player:       player,
		CurrentEvent: event,
		IsGameOver:   false,
		Final:        nil,
	}, nil
}

func (s *GameService) MakeChoice(ctx context.Context, playerID, choiceID string, expectedVersion int) (*domain.ChooseResult, error) {
	player, err := s.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return nil, ErrPlayerNotFound
	}

	if player.Version != expectedVersion {
		return nil, ErrConcurrentModify
	}

	currentEvent, err := s.getCurrentEventForPlayer(ctx, player)
	if err != nil {
		return nil, err
	}

	choice := s.findChoice(currentEvent, choiceID)
	if choice == nil || !choice.Available {
		return nil, ErrInvalidChoice
	}

	difficultyModifier := s.getDifficultyModifier(player.Turn)
	checkResult := s.performCheck(ctx, player, choice, difficultyModifier)
	effects := s.calculateEffects(ctx, player, checkResult, choice)

	for _, effect := range effects {
		s.applyEffect(&player.Stats, effect)
	}

	player.Turn++
	player.Version++
	player.UpdatedAt = time.Now()

	s.updateRanks(player)

	gameOver := s.checkGameOver(ctx, player)

	var final *domain.Final
	if gameOver {
		player.IsFinished = true
		now := time.Now()
		player.FinishedAt = &now
		final = s.calculateFinal(ctx, player)
	}

	if err := s.playerRepo.Update(ctx, player); err != nil {
		return nil, fmt.Errorf("update player: %w", err)
	}

	logEntry := &domain.GameLogEntry{
		ID:               uuid.New().String(),
		PlayerID:         playerID,
		Turn:             player.Turn - 1,
		EventTemplateID:  currentEvent.TemplateID,
		EventDescription: currentEvent.Description,
		ChoiceText:       choice.Text,
		CheckResult:      *checkResult,
		Effects:          effects,
		CreatedAt:        time.Now(),
	}

	if err := s.historyRepo.CreateEntry(ctx, logEntry); err != nil {
		return nil, fmt.Errorf("create log entry: %w", err)
	}

	var nextEvent *domain.EventInstance
	if !gameOver {
		nextEvent, _ = s.generateEvent(ctx, player)
	}

	return &domain.ChooseResult{
		Success:       true,
		CheckResult:   *checkResult,
		Effects:       effects,
		UpdatedPlayer: *player,
		NextEvent:     nextEvent,
		GameOver:      gameOver,
		Final:         final,
		NewVersion:    player.Version,
	}, nil
}

func (s *GameService) LoadGame(ctx context.Context, gameID string) (*domain.GameState, error) {
	player, err := s.playerRepo.GetByID(ctx, gameID)
	if err != nil {
		return nil, ErrGameNotFound
	}

	history, err := s.historyRepo.GetByPlayerID(ctx, gameID, 10)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}

	var currentEvent *domain.EventInstance
	if !player.IsFinished {
		currentEvent, _ = s.generateEvent(ctx, player)
	}

	gameOver := player.IsFinished || player.Turn >= MaxTurn || player.Stats.Mor <= 0

	var final *domain.Final
	if gameOver {
		final = s.calculateFinal(ctx, player)
	}

	return &domain.GameState{
		Player:       *player,
		CurrentEvent: currentEvent,
		EventHistory: history,
		IsGameOver:   gameOver,
		Final:        final,
		GameID:       gameID,
	}, nil
}

func (s *GameService) SelectCharacter(ctx context.Context, gameID string, characterType string, stats domain.PlayerStats) (*domain.Player, error) {
	player, err := s.playerRepo.GetByID(ctx, gameID)
	if err != nil {
		return nil, ErrPlayerNotFound
	}

	player.Stats = stats
	player.Version++
	player.UpdatedAt = time.Now()

	if err := s.playerRepo.Update(ctx, player); err != nil {
		return nil, fmt.Errorf("update player: %w", err)
	}

	return player, nil
}

func (s *GameService) RestartGame(ctx context.Context, gameID string) (*StartGameResult, error) {
	player, err := s.playerRepo.GetByID(ctx, gameID)
	if err != nil {
		return nil, ErrGameNotFound
	}

	if err := s.historyRepo.DeleteByPlayerID(ctx, gameID); err != nil {
		return nil, fmt.Errorf("delete history: %w", err)
	}

	player.Stats = domain.PlayerStats{Str: 50, End: 50, Agi: 50, Mor: 50, Disc: 0}
	player.FormalRank = domain.FormalRankRyadovoy
	player.InformalStatus = domain.InformalStatusZapah
	player.Turn = 1
	player.Flags = []string{}
	player.IsFinished = false
	player.Version = 1
	player.UpdatedAt = time.Now()

	if err := s.playerRepo.Update(ctx, player); err != nil {
		return nil, fmt.Errorf("update player: %w", err)
	}

	event, err := s.generateEvent(ctx, player)
	if err != nil {
		return nil, fmt.Errorf("generate event: %w", err)
	}

	return &StartGameResult{
		GameID:       gameID,
		Player:       player,
		CurrentEvent: event,
		IsGameOver:   false,
		Final:        nil,
	}, nil
}

func (s *GameService) generateEvent(ctx context.Context, player *domain.Player) (*domain.EventInstance, error) {
	if len(s.eventTemplates) == 0 {
		return nil, ErrNoEvents
	}

	recentEvents := s.getRecentEventIDs(ctx, player.ID)
	weights := s.calculateEventWeights(player, recentEvents)

	var template *domain.EventTemplate
	attempts := 0
	for attempts < 10 {
		template = s.selectWeightedEvent(weights)
		if template != nil {
			break
		}
		weights = s.calculateEventWeights(player, []string{})
		attempts++
	}

	if template == nil {
		template = &s.eventTemplates[s.rng.Intn(len(s.eventTemplates))]
	}

	choices := make([]domain.Choice, len(template.Choices))
	for i, c := range template.Choices {
		choices[i] = domain.Choice{
			ID:        c.ID,
			Text:      c.Text,
			Available: true,
			CheckType: c.CheckType,
			Threshold: c.Threshold,
			Chance:    c.Chance,
			Effects:   c.Effects,
		}
	}

	return &domain.EventInstance{
		ID:                uuid.New().String(),
		TemplateID:        template.ID,
		Description:       template.Description,
		ResolvedVariables: map[string]string{},
		Choices:           choices,
		Context: domain.EventContext{
			Time:     "день",
			Location: template.Location,
			Urgency:  "обычный",
		},
	}, nil
}

func (s *GameService) getRecentEventIDs(ctx context.Context, playerID string) []string {
	history, err := s.historyRepo.GetByPlayerID(ctx, playerID, 5)
	if err != nil || len(history) == 0 {
		return []string{}
	}

	ids := make([]string, len(history))
	for i, entry := range history {
		ids[i] = entry.EventTemplateID
	}
	return ids
}

func (s *GameService) getCurrentEventForPlayer(ctx context.Context, player *domain.Player) (*domain.EventInstance, error) {
	return s.generateEvent(ctx, player)
}

func (s *GameService) findChoice(event *domain.EventInstance, choiceID string) *domain.Choice {
	for i := range event.Choices {
		if event.Choices[i].ID == choiceID {
			return &event.Choices[i]
		}
	}
	return nil
}

func (s *GameService) performCheck(ctx context.Context, player *domain.Player, choice *domain.Choice, difficultyModifier float64) *domain.CheckResult {
	roll := s.rng.Intn(100) + 1

	checkType := choice.CheckType
	if checkType == "" {
		checkType = domain.CheckTypeThreshold
	}

	switch checkType {
	case domain.CheckTypeThreshold:
		return s.thresholdCheck(player, choice, roll, difficultyModifier)
	case domain.CheckTypeProbability:
		return s.probabilityCheck(player, choice, roll, difficultyModifier)
	case domain.CheckTypeCatastrophic:
		return s.catastrophicCheck(player, choice, roll, difficultyModifier)
	default:
		return s.thresholdCheck(player, choice, roll, difficultyModifier)
	}
}

func (s *GameService) thresholdCheck(player *domain.Player, choice *domain.Choice, roll int, modifier float64) *domain.CheckResult {
	effectiveThreshold := float64(50) * modifier

	success := roll <= int(effectiveThreshold)

	if success {
		outcome := domain.OutcomeSuccess
		description := "Успех!"
		if roll > int(effectiveThreshold*0.7) {
			description = "Успех! Но было близко."
		}
		return &domain.CheckResult{
			Success:     true,
			Outcome:     outcome,
			Description: description,
		}
	}

	return &domain.CheckResult{
		Success:     false,
		Outcome:     domain.OutcomeFailure,
		Description: "Неудача.",
	}
}

func (s *GameService) probabilityCheck(player *domain.Player, choice *domain.Choice, roll int, modifier float64) *domain.CheckResult {
	chance := choice.Chance
	if chance == 0 {
		chance = 50
	}

	effectiveChance := float64(chance) / modifier

	success := roll <= int(effectiveChance)

	if success {
		return &domain.CheckResult{
			Success:     true,
			Outcome:     domain.OutcomeSuccess,
			Description: "Повезло!",
		}
	}

	return &domain.CheckResult{
		Success:     false,
		Outcome:     domain.OutcomeFailure,
		Description: "Не повезло.",
	}
}

func (s *GameService) catastrophicCheck(player *domain.Player, choice *domain.Choice, roll int, modifier float64) *domain.CheckResult {
	noticeChance := 30

	noticed := roll <= noticeChance

	if !noticed {
		return &domain.CheckResult{
			Success:     true,
			Outcome:     domain.OutcomeIgnored,
			Description: "Никто не заметил.",
		}
	}

	powerRoll := s.rng.Intn(100) + 1
	effectiveThreshold := float64(50) * modifier

	if powerRoll <= int(effectiveThreshold) {
		return &domain.CheckResult{
			Success:     true,
			Outcome:     domain.OutcomeNoticedSuccess,
			Description: "Заметили, но удалось отмазаться!",
		}
	}

	return &domain.CheckResult{
		Success:     false,
		Outcome:     domain.OutcomeNoticedFailure,
		Description: "Заметили и наказали!",
	}
}

func (s *GameService) calculateEffects(ctx context.Context, player *domain.Player, result *domain.CheckResult, choice *domain.Choice) []domain.Effect {
	effects := []domain.Effect{}

	if len(choice.Effects) > 0 {
		for _, eff := range choice.Effects {
			prevValue := s.getStatValue(&player.Stats, eff.Stat)
			newValue := prevValue + eff.Delta
			effects = append(effects, domain.Effect{
				Stat:          eff.Stat,
				Delta:         eff.Delta,
				PreviousValue: prevValue,
				NewValue:      newValue,
			})
		}
		return effects
	}

	effectMap := map[string]int{
		"str":  0,
		"end":  0,
		"agi":  0,
		"mor":  0,
		"disc": 0,
	}

	switch result.Outcome {
	case domain.OutcomeSuccess, domain.OutcomeNoticedSuccess:
		effectMap["mor"] = s.rng.Intn(3) + 1
		effectMap["disc"] = s.rng.Intn(2)
	case domain.OutcomePartial:
		effectMap["mor"] = -s.rng.Intn(2) - 1
	case domain.OutcomeFailure, domain.OutcomeNoticedFailure:
		effectMap["mor"] = -s.rng.Intn(4) - 3
		effectMap["disc"] = -s.rng.Intn(2) - 1
	case domain.OutcomeIgnored:
	}

	for stat, delta := range effectMap {
		if delta != 0 {
			prevValue := s.getStatValue(&player.Stats, stat)
			effects = append(effects, domain.Effect{
				Stat:          stat,
				Delta:         delta,
				PreviousValue: prevValue,
				NewValue:      prevValue + delta,
			})
		}
	}

	return effects
}

func (s *GameService) getStatValue(stats *domain.PlayerStats, stat string) int {
	switch stat {
	case "str":
		return stats.Str
	case "end":
		return stats.End
	case "agi":
		return stats.Agi
	case "mor":
		return stats.Mor
	case "disc":
		return stats.Disc
	}
	return 0
}

func (s *GameService) applyEffect(stats *domain.PlayerStats, effect domain.Effect) {
	newValue := effect.NewValue

	switch effect.Stat {
	case "str":
		stats.Str = clamp(newValue, 1, 100)
	case "end":
		stats.End = clamp(newValue, 1, 100)
	case "agi":
		stats.Agi = clamp(newValue, 1, 100)
	case "mor":
		stats.Mor = clamp(newValue, 0, 100)
	case "disc":
		stats.Disc = clamp(newValue, -100, 100)
	}
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (s *GameService) updateRanks(player *domain.Player) {
	disc := player.Stats.Disc

	if disc >= 75 {
		player.FormalRank = domain.FormalRankSerzhant
	} else if disc >= 50 {
		player.FormalRank = domain.FormalRankMlSerzhant
	} else if disc >= 25 {
		player.FormalRank = domain.FormalRankEfreitor
	} else {
		player.FormalRank = domain.FormalRankRyadovoy
	}

	if disc <= -100 {
		player.InformalStatus = domain.InformalStatusDembel
	} else if disc <= -90 {
		player.InformalStatus = domain.InformalStatusDed
	} else if disc <= -75 {
		player.InformalStatus = domain.InformalStatusCherpak
	} else if disc <= -50 {
		player.InformalStatus = domain.InformalStatusSlon
	} else if disc < 0 {
		player.InformalStatus = domain.InformalStatusDukh
	} else {
		player.InformalStatus = domain.InformalStatusZapah
	}
}

func (s *GameService) checkGameOver(ctx context.Context, player *domain.Player) bool {
	return player.Turn >= MaxTurn || player.Stats.Mor <= 0
}

func (s *GameService) calculateFinal(ctx context.Context, player *domain.Player) *domain.Final {
	disc := player.Stats.Disc
	mor := player.Stats.Mor

	if mor <= 0 {
		return &domain.Final{
			Type:           domain.FinalTypeSlomannyiDembel,
			Title:          "Сломанный",
			Subtitle:       "Косячник",
			Description:    "Служба в армии принесла вам лишь горе, затрещины и сломанный нос, когда вы ночью упали с кровати, исполняя 'летучую мышь' по команде сержанта Златопузова. Впрочем, вы выжили — а это главное.",
			FinalStats:     player.Stats,
			AchievedOnTurn: player.Turn,
		}
	}

	if disc > 0 {
		return &domain.Final{
			Type:           domain.FinalTypeTihiyDembel,
			Title:          "Тихий дембель",
			Subtitle:       "Приспособленец",
			Description:    "Ваша служба прошла неспеша, вы не выделялись, не косячили, выполняли приказы 'шакалов' и офицеров. В общем — вы просто жили.",
			FinalStats:     player.Stats,
			AchievedOnTurn: player.Turn,
		}
	}

	return &domain.Final{
		Type:           domain.FinalTypeUvażajemyiDembel,
		Title:          "Уважаемый",
		Subtitle:       "Армани",
		Description:    "Вы ушли со службы в золотом кашне. Вся рота провожала вас со слезами на глазах, и обсуждала, как вы смогли так 'подняться' за такой короткий срок. А ваш кореш — старший прапорщик Залупенко плакал от горя, ведь вы смогли 'загнать' 300 тонн гнутых гвоздей...",
		FinalStats:     player.Stats,
		AchievedOnTurn: player.Turn,
	}
}
