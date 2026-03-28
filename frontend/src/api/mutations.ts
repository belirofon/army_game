import { gql } from '@apollo/client/core'

export const START_GAME = gql`
  mutation StartGame {
    startGame {
      gameId
      player {
        id
        stats {
          str
          end
          agi
          mor
          disc
        }
        formalRank
        informalStatus
        turn
        flags
        isFinished
        version
        createdAt
        updatedAt
      }
      currentEvent {
        id
        templateId
        description
        choices {
          id
          text
          available
        }
        context {
          time
          location
          urgency
        }
      }
      isGameOver
      final {
        type
        title
        subtitle
        description
        finalStats {
          str
          end
          agi
          mor
          disc
        }
        achievedOnTurn
      }
    }
  }
`

export const CHOOSE = gql`
  mutation Choose($playerId: ID!, $choiceId: String!, $expectedVersion: Int!) {
    choose(playerId: $playerId, choiceId: $choiceId, expectedVersion: $expectedVersion) {
      success
      checkResult {
        success
        outcome
        description
      }
      effects {
        stat
        delta
        previousValue
        newValue
      }
      updatedPlayer {
        id
        stats {
          str
          end
          agi
          mor
          disc
        }
        formalRank
        informalStatus
        turn
        flags
        isFinished
        version
        createdAt
        updatedAt
      }
      nextEvent {
        id
        templateId
        description
        choices {
          id
          text
          available
        }
        context {
          time
          location
          urgency
        }
      }
      gameOver
      final {
        type
        title
        subtitle
        description
        finalStats {
          str
          end
          agi
          mor
          disc
        }
        achievedOnTurn
      }
      newVersion
    }
  }
`

export const SELECT_CHARACTER = gql`
  mutation SelectCharacter($playerId: ID!, $characterType: String!, $stats: PlayerStatsInput!) {
    selectCharacter(playerId: $playerId, characterType: $characterType, stats: $stats) {
      id
      stats {
        str
        end
        agi
        mor
        disc
      }
      formalRank
      informalStatus
      turn
      flags
      isFinished
      version
      createdAt
      updatedAt
    }
  }
`

export const RESTART_GAME = gql`
  mutation RestartGame($playerId: ID!) {
    restartGame(playerId: $playerId) {
      gameId
      player {
        id
        stats {
          str
          end
          agi
          mor
          disc
        }
        formalRank
        informalStatus
        turn
        flags
        isFinished
        version
        createdAt
        updatedAt
      }
      currentEvent {
        id
        templateId
        description
        choices {
          id
          text
          available
        }
        context {
          time
          location
          urgency
        }
      }
      isGameOver
      final {
        type
        title
        subtitle
        description
        finalStats {
          str
          end
          agi
          mor
          disc
        }
        achievedOnTurn
      }
    }
  }
`