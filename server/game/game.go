package game

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type game struct {
	gameId         string
	gameType       string
	word           string
	guessedLetters []string
	livesRemaining int
	turn           string
	full           bool
	lastGuessTime  time.Time
	connPlayerA    *websocket.Conn
	connPlayerB    *websocket.Conn
	abandoned      bool
}

const (
	gameTypeLocal  = "local"
	gameTypeRemote = "remote"
)

const (
	playerA = "a"
	playerB = "b"
)

const gameIdLetters = "ABCDEF0123456789"
const gameIdLength = 18

const gameLives = 10
const gameGuessWordLockDuration = 5

var remoteGames map[string]*game = make(map[string]*game)

var (
	errUnknownGameType = errors.New("unknown game type")
	errGameNotFound    = errors.New("game not found")
	errGameOver        = errors.New("game is already over")
	errLetterUnknown   = errors.New("unknown letter")
	errLetterGuessed   = errors.New("already guessed this letter")
	errGameNotOver     = errors.New("game is not over")
	errGameAbandoned   = errors.New("game is abandoned")
)

func newGame(gameType string, conn *websocket.Conn) (*game, error) {
	if gameType != gameTypeLocal && gameType != gameTypeRemote {
		return nil, errUnknownGameType
	}

	gameId := ""
	for i := 0; i < gameIdLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(gameIdLetters))))
		if err != nil {
			panic(err)
		}

		gameId += string(gameIdLetters[num.Int64()])
	}

	wordIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(wordList))))
	if err != nil {
		panic(err)
	}

	game := &game{
		gameId:         gameId,
		gameType:       gameType,
		word:           wordList[wordIndex.Int64()],
		guessedLetters: []string{},
		livesRemaining: gameLives,
		turn:           playerA,
		full:           gameType == gameTypeLocal,
		lastGuessTime:  time.UnixMilli(0),
		connPlayerA:    conn,
		abandoned:      false,
	}

	if game.gameType == gameTypeRemote {
		remoteGames[game.gameId] = game
	}

	return game, nil
}

func findGame(gameId string) (*game, error) {
	game, ok := remoteGames[gameId]
	if !ok {
		return nil, errGameNotFound
	}

	return game, nil
}

func (game *game) guessLetter(letter string) error {
	if game.abandoned {
		return errGameAbandoned
	}

	if game.gameOver() {
		return errGameOver
	}

	if !strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZ", letter) {
		return errLetterUnknown
	}

	for _, guessedLetter := range game.guessedLetters {
		if string(guessedLetter) == letter {
			return errLetterGuessed
		}
	}

	if !strings.Contains(game.word, letter) {
		game.livesRemaining--
	}

	if game.livesRemaining > 0 {
		game.guessedLetters = append(game.guessedLetters, letter)
	} else {
		game.guessedLetters = strings.Split(game.word, "")
	}

	game.swapTurn()

	return nil
}

func (game *game) guessWord(word string) error {
	if game.abandoned {
		return errGameAbandoned
	}

	if game.guessWordLocked() {
		return nil
	}

	if game.gameOver() {
		return errGameOver
	}

	if game.word == word {
		game.guessedLetters = strings.Split(game.word, "")
	}

	game.lastGuessTime = time.Now()

	game.swapTurn()

	return nil
}

func (game *game) rematch() error {
	if game.abandoned {
		return errGameAbandoned
	}

	if !game.gameOver() {
		return errGameNotOver
	}

	wordIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(wordList))))
	if err != nil {
		panic(err)
	}

	game.word = wordList[wordIndex.Int64()]
	game.guessedLetters = []string{}
	game.livesRemaining = gameLives
	game.turn = playerA
	game.lastGuessTime = time.UnixMilli(0)

	return nil
}

func (game *game) hiddenWord() string {
	hiddenWord := ""

	for _, letter := range game.word {
		letterWasGuessed := false

		for _, guessedLetter := range game.guessedLetters {
			if string(letter) == guessedLetter {
				letterWasGuessed = true
				break
			}
		}

		if letterWasGuessed {
			hiddenWord += string(letter)
		} else {
			hiddenWord += "_"
		}
	}

	return hiddenWord
}

func (game *game) swapTurn() {
	if game.turn == playerA {
		game.turn = playerB
	} else if game.turn == playerB {
		game.turn = playerA
	}
}

func (game *game) guessWordLocked() bool {
	return time.Since(game.lastGuessTime).Seconds() < gameGuessWordLockDuration
}

func (game *game) gameOver() bool {
	return !strings.Contains(game.hiddenWord(), "_")
}
