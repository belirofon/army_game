package graphql

import (
	"context"
	"log"
	"net/http"

	"army-game-backend/internal/domain"
	"army-game-backend/internal/service"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type Resolver struct {
	gameService *service.GameService
}

func NewResolver(gameService *service.GameService) *Resolver {
	return &Resolver{
		gameService: gameService,
	}
}

func (r *Resolver) Mutation() *MutationResolver {
	return &MutationResolver{r}
}

func (r *Resolver) Query() *QueryResolver {
	return &QueryResolver{r}
}

type MutationResolver struct {
	*Resolver
}

type QueryResolver struct {
	*Resolver
}

func (r *MutationResolver) StartGame(ctx context.Context) *StartGamePayloadResolver {
	result, err := r.gameService.StartGame(ctx)
	if err != nil {
		log.Printf("StartGame error: %v", err)
		return nil
	}

	return &StartGamePayloadResolver{
		gameID:       result.GameID,
		player:       result.Player,
		currentEvent: result.CurrentEvent,
		isGameOver:   result.IsGameOver,
		final:        result.Final,
	}
}

func (r *MutationResolver) Choose(ctx context.Context, args struct {
	PlayerID        graphql.ID
	ChoiceID        string
	ExpectedVersion int32
}) *ChooseResultResolver {
	result, err := r.gameService.MakeChoice(
		ctx,
		string(args.PlayerID),
		args.ChoiceID,
		int(args.ExpectedVersion),
	)
	if err != nil {
		log.Printf("Choose error: %v", err)
		return nil
	}

	return &ChooseResultResolver{
		result: result,
	}
}

func (r *MutationResolver) SelectCharacter(ctx context.Context, args struct {
	PlayerID      graphql.ID
	CharacterType string
	Stats         struct {
		Str  int32
		End  int32
		Agi  int32
		Mor  int32
		Disc int32
	}
}) *PlayerResolver {
	result, err := r.gameService.SelectCharacter(
		ctx,
		string(args.PlayerID),
		args.CharacterType,
		domain.PlayerStats{
			Str:  int(args.Stats.Str),
			End:  int(args.Stats.End),
			Agi:  int(args.Stats.Agi),
			Mor:  int(args.Stats.Mor),
			Disc: int(args.Stats.Disc),
		},
	)
	if err != nil {
		log.Printf("SelectCharacter error: %v", err)
		return nil
	}

	return &PlayerResolver{player: result}
}

func (r *MutationResolver) RestartGame(ctx context.Context, args struct {
	PlayerID graphql.ID
}) *StartGamePayloadResolver {
	result, err := r.gameService.RestartGame(ctx, string(args.PlayerID))
	if err != nil {
		log.Printf("RestartGame error: %v", err)
		return nil
	}

	return &StartGamePayloadResolver{
		gameID:       result.GameID,
		player:       result.Player,
		currentEvent: result.CurrentEvent,
		isGameOver:   result.IsGameOver,
		final:        result.Final,
	}
}

func (r *QueryResolver) LoadGame(ctx context.Context, args struct {
	GameID graphql.ID
}) *GameStateResolver {
	result, err := r.gameService.LoadGame(ctx, string(args.GameID))
	if err != nil {
		log.Printf("LoadGame error: %v", err)
		return nil
	}

	return &GameStateResolver{
		state: result,
	}
}

type StartGamePayloadResolver struct {
	gameID       string
	player       *domain.Player
	currentEvent *domain.EventInstance
	isGameOver   bool
	final        *domain.Final
}

func (r *StartGamePayloadResolver) GameID() graphql.ID {
	return graphql.ID(r.gameID)
}

func (r *StartGamePayloadResolver) Player() *PlayerResolver {
	if r.player == nil {
		return nil
	}
	return &PlayerResolver{player: r.player}
}

func (r *StartGamePayloadResolver) CurrentEvent() *EventInstanceResolver {
	if r.currentEvent == nil {
		return nil
	}
	return &EventInstanceResolver{event: r.currentEvent}
}

func (r *StartGamePayloadResolver) IsGameOver() bool {
	return r.isGameOver
}

func (r *StartGamePayloadResolver) Final() *FinalResolver {
	if r.final == nil {
		return nil
	}
	return &FinalResolver{final: r.final}
}

type ChooseResultResolver struct {
	result *domain.ChooseResult
}

func (r *ChooseResultResolver) Success() bool {
	return r.result.Success
}

func (r *ChooseResultResolver) CheckResult() *CheckResultResolver {
	return &CheckResultResolver{result: &r.result.CheckResult}
}

func (r *ChooseResultResolver) Effects() []*EffectResolver {
	effects := make([]*EffectResolver, len(r.result.Effects))
	for i, e := range r.result.Effects {
		effects[i] = &EffectResolver{effect: &e}
	}
	return effects
}

func (r *ChooseResultResolver) UpdatedPlayer() *PlayerResolver {
	return &PlayerResolver{player: &r.result.UpdatedPlayer}
}

func (r *ChooseResultResolver) NextEvent() *EventInstanceResolver {
	if r.result.NextEvent == nil {
		return nil
	}
	return &EventInstanceResolver{event: r.result.NextEvent}
}

func (r *ChooseResultResolver) GameOver() bool {
	return r.result.GameOver
}

func (r *ChooseResultResolver) Final() *FinalResolver {
	if r.result.Final == nil {
		return nil
	}
	return &FinalResolver{final: r.result.Final}
}

func (r *ChooseResultResolver) NewVersion() int32 {
	return int32(r.result.NewVersion)
}

type GameStateResolver struct {
	state *domain.GameState
}

func (r *GameStateResolver) Player() *PlayerResolver {
	return &PlayerResolver{player: &r.state.Player}
}

func (r *GameStateResolver) CurrentEvent() *EventInstanceResolver {
	if r.state.CurrentEvent == nil {
		return nil
	}
	return &EventInstanceResolver{event: r.state.CurrentEvent}
}

func (r *GameStateResolver) EventHistory() []*GameLogEntryResolver {
	history := make([]*GameLogEntryResolver, len(r.state.EventHistory))
	for i, h := range r.state.EventHistory {
		history[i] = &GameLogEntryResolver{entry: &h}
	}
	return history
}

func (r *GameStateResolver) IsGameOver() bool {
	return r.state.IsGameOver
}

func (r *GameStateResolver) Final() *FinalResolver {
	if r.state.Final == nil {
		return nil
	}
	return &FinalResolver{final: r.state.Final}
}

func (r *GameStateResolver) GameID() graphql.ID {
	return graphql.ID(r.state.GameID)
}

type PlayerResolver struct {
	player *domain.Player
}

func (r *PlayerResolver) ID() graphql.ID {
	return graphql.ID(r.player.ID)
}

func (r *PlayerResolver) Stats() *PlayerStatsResolver {
	return &PlayerStatsResolver{stats: &r.player.Stats}
}

func (r *PlayerResolver) FormalRank() string {
	return string(r.player.FormalRank)
}

func (r *PlayerResolver) InformalStatus() string {
	return string(r.player.InformalStatus)
}

func (r *PlayerResolver) Turn() int32 {
	return int32(r.player.Turn)
}

func (r *PlayerResolver) Flags() []string {
	return r.player.Flags
}

func (r *PlayerResolver) IsFinished() bool {
	return r.player.IsFinished
}

func (r *PlayerResolver) Version() int32 {
	return int32(r.player.Version)
}

func (r *PlayerResolver) CreatedAt() string {
	return r.player.CreatedAt.Format("2006-01-02T15:04:05Z")
}

func (r *PlayerResolver) UpdatedAt() string {
	return r.player.UpdatedAt.Format("2006-01-02T15:04:05Z")
}

type PlayerStatsResolver struct {
	stats *domain.PlayerStats
}

func (r *PlayerStatsResolver) Str() int32 {
	return int32(r.stats.Str)
}

func (r *PlayerStatsResolver) End() int32 {
	return int32(r.stats.End)
}

func (r *PlayerStatsResolver) Agi() int32 {
	return int32(r.stats.Agi)
}

func (r *PlayerStatsResolver) Mor() int32 {
	return int32(r.stats.Mor)
}

func (r *PlayerStatsResolver) Disc() int32 {
	return int32(r.stats.Disc)
}

type EventInstanceResolver struct {
	event *domain.EventInstance
}

func (r *EventInstanceResolver) ID() graphql.ID {
	return graphql.ID(r.event.ID)
}

func (r *EventInstanceResolver) TemplateID() string {
	return r.event.TemplateID
}

func (r *EventInstanceResolver) Description() string {
	return r.event.Description
}

func (r *EventInstanceResolver) Choices() []*ChoiceResolver {
	choices := make([]*ChoiceResolver, len(r.event.Choices))
	for i, c := range r.event.Choices {
		choices[i] = &ChoiceResolver{choice: &c}
	}
	return choices
}

func (r *EventInstanceResolver) Context() *EventContextResolver {
	return &EventContextResolver{context: r.event.Context}
}

type ChoiceResolver struct {
	choice *domain.Choice
}

func (r *ChoiceResolver) ID() graphql.ID {
	return graphql.ID(r.choice.ID)
}

func (r *ChoiceResolver) Text() string {
	return r.choice.Text
}

func (r *ChoiceResolver) Available() bool {
	return r.choice.Available
}

type EventContextResolver struct {
	context domain.EventContext
}

func (r *EventContextResolver) Time() string {
	return r.context.Time
}

func (r *EventContextResolver) Location() string {
	return string(r.context.Location)
}

func (r *EventContextResolver) Urgency() string {
	return r.context.Urgency
}

type CheckResultResolver struct {
	result *domain.CheckResult
}

func (r *CheckResultResolver) Success() bool {
	return r.result.Success
}

func (r *CheckResultResolver) Outcome() string {
	return string(r.result.Outcome)
}

func (r *CheckResultResolver) Description() string {
	return r.result.Description
}

type EffectResolver struct {
	effect *domain.Effect
}

func (r *EffectResolver) Stat() string {
	return r.effect.Stat
}

func (r *EffectResolver) Delta() int32 {
	return int32(r.effect.Delta)
}

func (r *EffectResolver) PreviousValue() int32 {
	return int32(r.effect.PreviousValue)
}

func (r *EffectResolver) NewValue() int32 {
	return int32(r.effect.NewValue)
}

type GameLogEntryResolver struct {
	entry *domain.GameLogEntry
}

func (r *GameLogEntryResolver) ID() graphql.ID {
	return graphql.ID(r.entry.ID)
}

func (r *GameLogEntryResolver) PlayerID() graphql.ID {
	return graphql.ID(r.entry.PlayerID)
}

func (r *GameLogEntryResolver) Turn() int32 {
	return int32(r.entry.Turn)
}

func (r *GameLogEntryResolver) EventDescription() string {
	return r.entry.EventDescription
}

func (r *GameLogEntryResolver) ChoiceText() string {
	return r.entry.ChoiceText
}

func (r *GameLogEntryResolver) CheckResult() *CheckResultResolver {
	return &CheckResultResolver{result: &r.entry.CheckResult}
}

func (r *GameLogEntryResolver) Effects() []*EffectResolver {
	effects := make([]*EffectResolver, len(r.entry.Effects))
	for i, e := range r.entry.Effects {
		effects[i] = &EffectResolver{effect: &e}
	}
	return effects
}

func (r *GameLogEntryResolver) CreatedAt() string {
	return r.entry.CreatedAt.Format("2006-01-02T15:04:05Z")
}

type FinalResolver struct {
	final *domain.Final
}

func (r *FinalResolver) Type() string {
	return string(r.final.Type)
}

func (r *FinalResolver) Title() string {
	return r.final.Title
}

func (r *FinalResolver) Subtitle() *string {
	if r.final.Subtitle == "" {
		return nil
	}
	return &r.final.Subtitle
}

func (r *FinalResolver) Description() string {
	return r.final.Description
}

func (r *FinalResolver) FinalStats() *PlayerStatsResolver {
	return &PlayerStatsResolver{stats: &r.final.FinalStats}
}

func (r *FinalResolver) AchievedOnTurn() int32 {
	return int32(r.final.AchievedOnTurn)
}

func AddGraphQLHandler(mux *http.ServeMux, resolver *Resolver) {
	schema := `
		type Query {
			loadGame(gameId: ID!): GameState
		}
		
		type Mutation {
			startGame: StartGamePayload!
			choose(playerId: ID!, choiceId: String!, expectedVersion: Int!): ChooseResult!
			selectCharacter(playerId: ID!, characterType: String!, stats: PlayerStatsInput!): Player!
			restartGame(playerId: ID!): StartGamePayload!
		}
		
		type StartGamePayload {
			gameId: ID!
			player: Player!
			currentEvent: EventInstance
			isGameOver: Boolean!
			final: Final
		}
		
		type Player {
			id: ID!
			stats: PlayerStats!
			formalRank: String!
			informalStatus: String!
			turn: Int!
			flags: [String!]!
			isFinished: Boolean!
			version: Int!
			createdAt: String!
			updatedAt: String!
		}
		
		type PlayerStats {
			str: Int!
			end: Int!
			agi: Int!
			mor: Int!
			disc: Int!
		}
		
		input PlayerStatsInput {
			str: Int!
			end: Int!
			agi: Int!
			mor: Int!
			disc: Int!
		}
		
		type EventInstance {
			id: ID!
			templateId: String!
			description: String!
			choices: [Choice!]!
			context: EventContext!
		}
		
		type Choice {
			id: ID!
			text: String!
			available: Boolean!
		}
		
		type EventContext {
			time: String!
			location: String!
			urgency: String!
		}
		
		type CheckResult {
			success: Boolean!
			outcome: String!
			description: String!
		}
		
		type Effect {
			stat: String!
			delta: Int!
			previousValue: Int!
			newValue: Int!
		}
		
		type ChooseResult {
			success: Boolean!
			checkResult: CheckResult!
			effects: [Effect!]!
			updatedPlayer: Player!
			nextEvent: EventInstance
			gameOver: Boolean!
			final: Final
			newVersion: Int!
		}
		
		type GameLogEntry {
			id: ID!
			playerId: ID!
			turn: Int!
			eventDescription: String!
			choiceText: String!
			checkResult: CheckResult!
			effects: [Effect!]!
			createdAt: String!
		}
		
		type Final {
			type: String!
			title: String!
			subtitle: String
			description: String!
			finalStats: PlayerStats!
			achievedOnTurn: Int!
		}
		
		type GameState {
			player: Player!
			currentEvent: EventInstance
			eventHistory: [GameLogEntry!]!
			isGameOver: Boolean!
			final: Final
			gameId: ID!
		}
	`

	s, err := graphql.ParseSchema(schema, resolver)
	if err != nil {
		log.Fatalf("Failed to parse schema: %v", err)
	}

	h := relay.Handler{Schema: s}
	mux.Handle("/graphql", &h)
}
