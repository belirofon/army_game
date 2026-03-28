import { gql } from '@apollo/client/core'

export const LOAD_GAME = gql`
  query LoadGame($gameId: ID!) {
    loadGame(gameId: $gameId) {
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
      eventHistory {
        id
        playerId
        turn
        eventDescription
        choiceText
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
        createdAt
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
      gameId
    }
  }
`