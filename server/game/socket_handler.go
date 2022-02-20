package game

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func createUpdatePacket(game *game) clientboundUpdatePacket {
	return clientboundUpdatePacket{
		Message:         messageUpdate,
		Word:            game.hiddenWord(),
		Turn:            game.turn,
		GuessedLetters:  game.guessedLetters,
		LivesRemaining:  game.livesRemaining,
		Full:            game.full,
		GuessWordLocked: game.guessWordLocked(),
	}
}

func SocketHandler(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		panic(err)
	}

	var game *game
	var me string

	conn.SetCloseHandler(func(code int, text string) error {
		if game != nil {
			abandonedPacket := clientboundAbandonedPacket{Message: messageAbandoned}

			if game.gameType == gameTypeRemote {
				if me == playerA {
					err := game.connPlayerB.WriteJSON(abandonedPacket)
					if err != nil {
						panic(err)
					}
				} else if me == playerB {
					err := game.connPlayerA.WriteJSON(abandonedPacket)
					if err != nil {
						panic(err)
					}
				}

				game.abandoned = true
			}

			delete(remoteGames, game.gameId)
		}

		return &websocket.CloseError{
			Code: code,
			Text: text,
		}
	})

	for {
		var packet serverboundPacket

		err := conn.ReadJSON(&packet)
		if err != nil {
			panic(err)
		}

		switch packet.Action {
		case actionCreateGame:
			if game != nil {
				panic("already in a game!")
			}

			game, err = newGame(packet.GameType, conn)
			if err != nil {
				panic(err)
			}

			me = playerA

			updatePacket := createUpdatePacket(game)

			err = game.connPlayerA.WriteJSON(updatePacket)
			if err != nil {
				panic(err)
			}

			if game.gameType == gameTypeRemote {
				err = conn.WriteJSON(clientboundWaitForPartnerPacket{
					Message: messageWaitForPartner,
					GameId:  game.gameId,
				})
				if err != nil {
					panic(err)
				}
			}
		case actionJoinGame:
			if game != nil {
				panic("already in a game!")
			}

			game, err = findGame(packet.GameId)
			if err != nil && err == errGameNotFound {
				err = conn.WriteJSON(clientboundCannotJoinGamePacket{
					Message: messageCannotJoinGame,
				})
				if err != nil {
					panic(err)
				}
				return
			} else if err != nil {
				panic(err)
			}

			game.connPlayerB = conn
			game.full = true

			delete(remoteGames, game.gameId)

			me = playerB

			updatePacket := createUpdatePacket(game)

			err = game.connPlayerA.WriteJSON(updatePacket)
			if err != nil {
				panic(err)
			}

			if game.connPlayerB != nil {
				err = game.connPlayerB.WriteJSON(updatePacket)
				if err != nil {
					panic(err)
				}
			}
		case actionGuessLetter:
			if game.gameType == gameTypeRemote && game.turn != me {
				panic("partner's turn")
			}

			err = game.guessLetter(packet.Letter)
			if err != nil {
				panic(err)
			}

			updatePacket := createUpdatePacket(game)

			err = game.connPlayerA.WriteJSON(updatePacket)
			if err != nil {
				panic(err)
			}

			if game.connPlayerB != nil {
				err = game.connPlayerB.WriteJSON(updatePacket)
				if err != nil {
					panic(err)
				}
			}
		case actionGuessWord:
			if game.gameType == gameTypeRemote && game.turn != me {
				panic("partner's turn")
			}

			err = game.guessWord(packet.Word)
			if err != nil {
				panic(err)
			}

			updatePacket := createUpdatePacket(game)

			err = game.connPlayerA.WriteJSON(updatePacket)
			if err != nil {
				panic(err)
			}

			if game.connPlayerB != nil {
				err = game.connPlayerB.WriteJSON(updatePacket)
				if err != nil {
					panic(err)
				}
			}
		case actionRematch:
			err = game.rematch()
			if err != nil {
				panic(err)
			}

			updatePacket := createUpdatePacket(game)

			err = game.connPlayerA.WriteJSON(updatePacket)
			if err != nil {
				panic(err)
			}

			if game.connPlayerB != nil {
				err = game.connPlayerB.WriteJSON(updatePacket)
				if err != nil {
					panic(err)
				}
			}
		case actionPing:
			// ignore ping packets
		default:
			panic("unknown packet")
		}
	}
}
