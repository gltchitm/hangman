package game

const (
	actionCreateGame  = "create_game"
	actionJoinGame    = "join_game"
	actionGuessLetter = "guess_letter"
	actionGuessWord   = "guess_word"
	actionRematch     = "new_game"
	actionPing        = "ping"
)

const (
	messageUpdate         = "update"
	messageWaitForPartner = "wait_for_partner"
	messageCannotJoinGame = "cannot_join_game"
	messageAbandoned      = "abandoned"
)

type serverboundPacket struct {
	Action string `json:"action"`

	// create_game
	GameType string `json:"game_type,omitempty"`

	// join_game
	GameId string `json:"game_id,omitempty"`

	// guess_letter
	Letter string `json:"letter,omitempty"`

	// guess_word
	Word string `json:"word,omitempty"`

	// rematch

	// ping
}

type clientboundUpdatePacket struct {
	Message string `json:"message"`

	Word            string   `json:"word"`
	GuessedLetters  []string `json:"letters"`
	LivesRemaining  int      `json:"lives"`
	Full            bool     `json:"full"`
	Turn            string   `json:"player"`
	GuessWordLocked bool     `json:"guess_word_locked"`
}

type clientboundWaitForPartnerPacket struct {
	Message string `json:"message"`

	GameId string `json:"game_id"`
}

type clientboundCannotJoinGamePacket struct {
	Message string `json:"message"`
}

type clientboundAbandonedPacket struct {
	Message string `json:"message"`
}
