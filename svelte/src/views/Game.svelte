<script lang="ts">
    import WaitingForPartner from '../components/WaitingForPartner.svelte'
    import Loading from '../components/Loading.svelte'
    import Modal from '../components/Modal.svelte'

    import { GameType, Player } from '../enum'
    import { socket } from '../store'

    export let gameType: GameType
    export let createGame: boolean
    export let joinGame: string

    let player = Player.A
    let me = Player.A

    let loading = true
    let guessWordLocked = false
    let waitingForPartner = false
    let confirmEndGame = false
    let gameOver = false
    let abandoned = false
    let canMove = true
    let wordGuess = ''

    let word: string
    let lives: number
    let gameId: string
    let letters: string[]
    let guessWordTimeout: ReturnType<typeof setTimeout>

    const send = (message: object) => {
        if ($socket.readyState === WebSocket.OPEN) {
            $socket.send(JSON.stringify(message))
        }
    }

    const guess = (character: string) => send({ action: 'guess_letter', letter: character })

    const lockGuessWord = () => {
        guessWordLocked = true
        if (guessWordTimeout) {
            clearTimeout(guessWordTimeout)
        }
        guessWordTimeout = setTimeout(() => {
            guessWordLocked = false
        }, 5000)
    }
    const guessWord = () => {
        send({
            action: 'guess_word',
            word: wordGuess.toUpperCase().trim()
        })
        wordGuess = ''
        lockGuessWord()
    }

    const newGame = () => send({ action: 'new_game' })

    $socket.addEventListener('message', ({ data }) => {
        const payload = JSON.parse(data)
        switch (payload.message) {
            case 'update':
                word = payload.word
                lives = Number(payload.lives)
                letters = payload.letters
                player = payload.player === 'a' ? Player.A : Player.B
                gameOver = lives <= 0 || !word.includes('_')
                waitingForPartner = !payload.full
                if (gameType === GameType.RemoteMultiplayer) {
                    canMove = player === me
                }
                guessWordLocked = payload.guess_word_locked
                break
            case 'wait_for_partner':
                gameId = payload.game_id
                waitingForPartner = true
                break
            case 'lock_guess_word':
                lockGuessWord()
                break
            case 'abandoned':
                abandoned = true
                break
            default:
                return
        }
        loading = false
    })

    if (gameType === GameType.RemoteMultiplayer && !createGame) {
            me = Player.B
            canMove = false
            send({
                action: 'join_game',
                game_id: joinGame
            })
    } else {
        send({
            action: 'create_game',
            game_type: gameType !== GameType.RemoteMultiplayer ? 'local' : 'remote'
        })
    }

</script>
<style>
    .word {
        letter-spacing: 3px;
        font-family: monospace;
        font-size: 64px;
        color: white;
    }
    @media screen and (max-width: 500px) {
        .word {
            font-size: 30px;
        }
    }
    .game {
        min-height: 100vh;
        display: flex;
        align-items: center;
        flex-direction: column;
        justify-content: center;
        padding: 20px;
    }
    .hud, .controls {
        display: flex;
        flex-direction: column;
        align-items: center;
    }
    .letters {
        display: flex;
        flex-wrap: wrap;
        justify-content: center;
        align-items: center;
        margin-top: 15px;
        margin-bottom: 15px;
        max-width: 530px;
        gap: 10px;
        padding: 10px;
    }
    .letters * {
        padding: 5px;
        width: 30px;
        height: 30px;
    }
    .guess-word {
        display: flex;
        gap: 5px;
        margin-bottom: 15px;
    }
    .game-controls {
        display: flex;
        gap: 5px;
    }
    .game-status {
        font-size: 30px;
    }
    .lives {
        font-size: 14px;
        color: white;
    }
</style>

{#if confirmEndGame}
    <Modal>
        <svelte:fragment slot="title">
            {#if gameType === GameType.RemoteMultiplayer}
                Leave Game
            {:else}
                End Game
            {/if}
        </svelte:fragment>
        <svelte:fragment slot="body">
            Are you sure you want to do this?
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                on:click={() => window.location.reload()}
            >Yes</button>
            <button
                class="btn btn-secondary"
                on:click={() => confirmEndGame = false}
            >No</button>
        </svelte:fragment>
    </Modal>
{/if}
{#if abandoned}
    <Modal>
        <svelte:fragment slot="title">Abandoned By Partner</svelte:fragment>
        <svelte:fragment slot="body">Your partner has abandoned the game.</svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                on:click={() => window.location.reload()}
            >Main Menu</button>
        </svelte:fragment>
    </Modal>
{/if}
{#if loading}
    <Loading />
{:else if waitingForPartner}
    <WaitingForPartner {gameId} />
{:else}
    <div class="game">
        <div class="hud">
            {#if gameType !== GameType.LocalSolo}
                <span class="game-status text-dark">
                    {#if gameOver}
                        Game Over
                    {:else}
                        {#if gameType === GameType.LocalMultiplayer}
                            {#if player === Player.A}
                                Player 1's turn
                            {:else}
                                Player 2's turn
                            {/if}
                        {:else}
                            {#if player === me}
                                Your turn
                            {:else}
                                Partner's turn
                            {/if}
                        {/if}
                    {/if}
                </span>
            {/if}
            <span class="word">{word}</span>
            <span class="lives">
                {#if gameOver && lives}
                    You win!
                {:else if gameOver}
                    You lose!
                {:else}
                    Lives Remaining: {lives}
                {/if}
            </span>
            <div class="letters bg-light rounded">
                {#each 'ABCDEFGHIJKLMNOPQRSTUVWXYZ' as letter}
                    <button
                        class="btn btn-secondary"
                        on:click={() => guess(letter)}
                        disabled={letters.includes(letter) || gameOver || !canMove}
                    >{letter}</button>
                {/each}
            </div>
        </div>
        <div class="controls">
            <div class="guess-word">
                <input
                    class="form-control"
                    placeholder="Guess Word"
                    bind:value={wordGuess}
                    disabled={guessWordLocked || gameOver || !canMove}
                />
                <button
                    class="btn btn-secondary"
                    disabled={guessWordLocked || gameOver || !canMove}
                    on:click={guessWord}
                >Guess</button>
            </div>
            <div class="game-controls">
                {#if gameOver}
                    <button
                        class="btn btn-secondary"
                        on:click={newGame}
                    >New Game</button>
                {/if}
                <button class="btn btn-secondary" on:click={() => confirmEndGame = true}>
                    {#if gameType === GameType.RemoteMultiplayer}
                        Leave Game
                    {:else}
                        End Game
                    {/if}
                </button>
            </div>
        </div>
    </div>
{/if}
