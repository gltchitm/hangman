package main

import (
	"net/http"

	"github.com/gltchitm/hangman/server/game"
)

func main() {
	http.HandleFunc("/ws", game.SocketHandler)

	fileServer := http.FileServer(http.Dir("../static"))

	http.Handle("/", fileServer)

	http.ListenAndServe(":5522", nil)
}
