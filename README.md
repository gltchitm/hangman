# Hangman
Game of Hangman written in Go with remote multiplayer support.

## Running
Execute `start.sh` to start the server on port 5522.

## Go Server
The Hangman server was originally written in Elixir but has since been rewritten in Go. Rewriting it has brought many new improvements such as better input validation. The rewritten server was designed to be very backwards compatible with the old server so that old clients can still connect to it.

## License
[MIT](LICENSE)
