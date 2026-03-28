export type FormalRank = 'РЯДОВОЙ' | 'ЕФРЕЙТОР' | 'МЛ_СЕРЖАНТ' | 'СЕРЖАНТ'

export type InformalStatus = 'ЗАПАХ' | 'ДУХ' | 'СЛОН' | 'ЧЕРПАК' | 'ДЕД' | 'ДЕМБЕЛЬ'

export type OutcomeType = 
  | 'SUCCESS' 
  | 'PARTIAL' 
  | 'FAILURE' 
  | 'IGNORED' 
  | 'NOTICED_SUCCESS' 
  | 'NOTICED_FAILURE'

export type FinalType = 'ТИХИЙ_ДЕМБЕЛЬ' | 'УВАЖАЕМЫЙ_ДЕМБЕЛЬ' | 'СЛОМАННЫЙ_ДЕМБЕЛЬ'

export interface PlayerStats {
  str: number
  end: number
  agi: number
  mor: number
  disc: number
}

export interface Player {
  id: string
  stats: PlayerStats
  formalRank: FormalRank
  informalStatus: InformalStatus
  turn: number
  flags: string[]
  isFinished: boolean
  version: number
  createdAt: string
  updatedAt: string
}

export interface EventContext {
  time: string
  location: string
  urgency: string
}

export interface Choice {
  id: string
  text: string
  available: boolean
}

export interface EventInstance {
  id: string
  templateId: string
  description: string
  resolvedVariables: Record<string, string>
  choices: Choice[]
  context: EventContext
}

export interface CheckResult {
  success: boolean
  outcome: OutcomeType
  description: string
}

export interface Effect {
  stat: keyof PlayerStats
  delta: number
  previousValue: number
  newValue: number
}

export interface GameLogEntry {
  id: string
  playerId: string
  turn: number
  eventDescription: string
  choiceText: string
  checkResult: CheckResult
  effects: Effect[]
  createdAt: string
}

export interface Final {
  type: FinalType
  title: string
  subtitle?: string
  description: string
  finalStats: PlayerStats
  achievedOnTurn: number
}

export interface GameState {
  player: Player
  currentEvent: EventInstance | null
  eventHistory: GameLogEntry[]
  isGameOver: boolean
  final: Final | null
  gameId: string
}

export interface ChooseResult {
  success: boolean
  checkResult: CheckResult
  effects: Effect[]
  updatedPlayer: Player
  nextEvent: EventInstance | null
  gameOver: boolean
  final: Final | null
  newVersion: number
}

export const STAT_RANGES: Record<keyof PlayerStats, { min: number; max: number }> = {
  str: { min: 1, max: 100 },
  end: { min: 1, max: 100 },
  agi: { min: 1, max: 100 },
  mor: { min: 0, max: 100 },
  disc: { min: -100, max: 100 },
}

export const FORMAL_RANKS: Record<FormalRank, number> = {
  'РЯДОВОЙ': -Infinity,
  'ЕФРЕЙТОР': 25,
  'МЛ_СЕРЖАНТ': 50,
  'СЕРЖАНТ': 75,
}

export const INFORMAL_STATUSES: Record<InformalStatus, number> = {
  'ЗАПАХ': 0,
  'ДУХ': -25,
  'СЛОН': -50,
  'ЧЕРПАК': -75,
  'ДЕД': -90,
  'ДЕМБЕЛЬ': -100,
}

export const MAX_TURN = 30