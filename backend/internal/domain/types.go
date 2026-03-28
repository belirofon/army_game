package domain

import (
	"time"
)

type FormalRank string

const (
	FormalRankRyadovoy   FormalRank = "РЯДОВОЙ"
	FormalRankEfreitor   FormalRank = "ЕФРЕЙТОР"
	FormalRankMlSerzhant FormalRank = "МЛ_СЕРЖАНТ"
	FormalRankSerzhant   FormalRank = "СЕРЖАНТ"
)

type InformalStatus string

const (
	InformalStatusZapah   InformalStatus = "ЗАПАХ"
	InformalStatusDukh    InformalStatus = "ДУХ"
	InformalStatusSlon    InformalStatus = "СЛОН"
	InformalStatusCherpak InformalStatus = "ЧЕРПАК"
	InformalStatusDed     InformalStatus = "ДЕД"
	InformalStatusDembel  InformalStatus = "ДЕМБЕЛЬ"
)

type OutcomeType string

const (
	OutcomeSuccess        OutcomeType = "SUCCESS"
	OutcomePartial        OutcomeType = "PARTIAL"
	OutcomeFailure        OutcomeType = "FAILURE"
	OutcomeIgnored        OutcomeType = "IGNORED"
	OutcomeNoticedSuccess OutcomeType = "NOTICED_SUCCESS"
	OutcomeNoticedFailure OutcomeType = "NOTICED_FAILURE"
)

type FinalType string

type EventType string

const (
	EventTypeRoutine    EventType = "ROUTINE"
	EventTypeSafe       EventType = "SAFE"
	EventTypeSocial     EventType = "SOCIAL"
	EventTypeInformal   EventType = "INFORMAL"
	EventTypeInspection EventType = "INSPECTION"
	EventTypeEmergency  EventType = "EMERGENCY"
)

type CheckType string

const (
	CheckTypeThreshold    CheckType = "THRESHOLD"
	CheckTypeProbability  CheckType = "PROBABILITY"
	CheckTypeCatastrophic CheckType = "CATASTROPHIC"
)

type EventTemplate struct {
	ID          string        `json:"id"`
	Type        EventType     `json:"type"`
	CheckType   CheckType     `json:"checkType"`
	Location    EventLocation `json:"location"`
	Description string        `json:"description"`
	Choices     []Choice      `json:"choices"`
}

type EventsData struct {
	Events []EventTemplate `json:"events"`
}

const (
	FinalTypeTihiyDembel      FinalType = "ТИХИЙ_ДЕМБЕЛЬ"
	FinalTypeUvażajemyiDembel FinalType = "УВАЖАЕМЫЙ"
	FinalTypeSlomannyiDembel  FinalType = "СЛОМАННЫЙ"
)

type PlayerStats struct {
	Str  int `json:"str"`
	End  int `json:"end"`
	Agi  int `json:"agi"`
	Mor  int `json:"mor"`
	Disc int `json:"disc"`
}

type Player struct {
	ID             string         `json:"id"`
	Stats          PlayerStats    `json:"stats"`
	FormalRank     FormalRank     `json:"formalRank"`
	InformalStatus InformalStatus `json:"informalStatus"`
	Turn           int            `json:"turn"`
	Flags          []string       `json:"flags"`
	IsFinished     bool           `json:"isFinished"`
	FinishedAt     *time.Time     `json:"finishedAt,omitempty"`
	Version        int            `json:"version"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

type EventLocation string

const (
	LocationDiningRoom EventLocation = "DINING_ROOM"
	LocationBarracks   EventLocation = "BARRACKS"
	LocationTraining   EventLocation = "TRAINING_GROUND"
	LocationStorage    EventLocation = "STORAGE"
	LocationGuardDuty  EventLocation = "GUARD_DUTY"
	LocationDefault    EventLocation = "DEFAULT"
)

type EventContext struct {
	Time     string        `json:"time"`
	Location EventLocation `json:"location"`
	Urgency  string        `json:"urgency"`
}

type Choice struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Available bool      `json:"available"`
	CheckType CheckType `json:"checkType,omitempty"`
	Threshold int       `json:"threshold,omitempty"`
	Chance    int       `json:"Chance,omitempty"`
	Effects   []Effect  `json:"effects,omitempty"`
}

type EventInstance struct {
	ID                string            `json:"id"`
	TemplateID        string            `json:"templateId"`
	Description       string            `json:"description"`
	ResolvedVariables map[string]string `json:"resolvedVariables"`
	Choices           []Choice          `json:"choices"`
	Context           EventContext      `json:"context"`
}

type CheckResult struct {
	Success     bool        `json:"success"`
	Outcome     OutcomeType `json:"outcome"`
	Description string      `json:"description"`
}

type Effect struct {
	Stat          string `json:"stat"`
	Delta         int    `json:"delta"`
	PreviousValue int    `json:"previousValue"`
	NewValue      int    `json:"newValue"`
}

type GameLogEntry struct {
	ID               string      `json:"id"`
	PlayerID         string      `json:"playerId"`
	Turn             int         `json:"turn"`
	EventTemplateID  string      `json:"eventTemplateId"`
	EventDescription string      `json:"eventDescription"`
	ChoiceText       string      `json:"choiceText"`
	CheckResult      CheckResult `json:"checkResult"`
	Effects          []Effect    `json:"effects"`
	CreatedAt        time.Time   `json:"createdAt"`
}

type Final struct {
	Type           FinalType   `json:"type"`
	Title          string      `json:"title"`
	Subtitle       string      `json:"subtitle,omitempty"`
	Description    string      `json:"description"`
	FinalStats     PlayerStats `json:"finalStats"`
	AchievedOnTurn int         `json:"achievedOnTurn"`
}

type ChooseResult struct {
	Success       bool           `json:"success"`
	CheckResult   CheckResult    `json:"checkResult"`
	Effects       []Effect       `json:"effects"`
	UpdatedPlayer Player         `json:"updatedPlayer"`
	NextEvent     *EventInstance `json:"nextEvent,omitempty"`
	GameOver      bool           `json:"gameOver"`
	Final         *Final         `json:"final,omitempty"`
	NewVersion    int            `json:"newVersion"`
}

type GameState struct {
	Player       Player         `json:"player"`
	CurrentEvent *EventInstance `json:"currentEvent"`
	EventHistory []GameLogEntry `json:"eventHistory"`
	IsGameOver   bool           `json:"isGameOver"`
	Final        *Final         `json:"final,omitempty"`
	GameID       string         `json:"gameId"`
}
