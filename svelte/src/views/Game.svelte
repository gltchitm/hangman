<script lang="ts">
    import Loading from '../components/Loading.svelte'
    import Button from '../components/Button.svelte'
    import WaitingForPartner from '../components/WaitingForPartner.svelte'

    import { Input, Modal, ModalBody, ModalFooter, ModalHeader } from 'sveltestrap'

    import { GameType, Player } from '../enum'
    import { socket } from '../store'

    export let gameType: GameType
    export let createGame: boolean
    export let joinGame: string

    const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.split('')

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

    const guess = (character: string) => {
        send({
            action: 'guess_letter',
            letter: character
        })
    }
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
    const newGame = () => {
        send({ action: 'new_game' })
    }

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
    }
    .game {
        display: flex;
        flex-direction: column;
        align-items: center;
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
    }
    .characters {
        display: flex;
        flex-wrap: wrap;
        justify-content: center;
        align-items: center;
        margin: 15px;
    }
</style>

{#if confirmEndGame}
    <Modal isOpen={true} transitionOptions={{ duration: 100 }}>
        <ModalHeader>
            {#if gameType === GameType.RemoteMultiplayer}
                Leave Game
            {:else}
                End Game
            {/if}
        </ModalHeader>
        <ModalBody>
            Are you sure you want to do this?
        </ModalBody>
        <ModalFooter>
            <Button
                on:click={() => window.location.reload()}
                color="secondary"
            >Yes</Button>
            <Button
                on:click={() => confirmEndGame = false}
                color="secondary"
            >No</Button>
        </ModalFooter>
    </Modal>
{/if}
{#if abandoned}
<Modal isOpen={true} transitionOptions={{ duration: 100 }}>
    <ModalHeader>Abandoned By Partner</ModalHeader>
    <ModalBody>Your partner has abandoned the game.</ModalBody>
    <ModalFooter>
        <Button
            on:click={() => window.location.reload()}
            color="secondary"
        >Main Menu</Button>
    </ModalFooter>
</Modal>
{/if}
{#if loading}
    <Loading />
{:else if waitingForPartner}
    <WaitingForPartner {gameId} />
{:else}
    <div class="game">
        {#if gameType !== GameType.LocalSolo}
            <span class="h4 text-dark">
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
        <span class="h1 word">{word}</span>
        <span class="h6">
            {#if gameOver && lives > 0}
                    You win!
            {:else if gameOver}
                You lose!
            {:else}
                Lives Remaining: {lives}
            {/if}
        </span>
        <div class="characters bg-light rounded">
            {#each alphabet as letter}
                <Button
                    style="margin: 5px; padding: 5px; width: 30px; height: 30px;"
                    on:click={() => guess(letter)}
                    disabled={letters.includes(letter) || gameOver || !canMove}
                >{letter}</Button>
            {/each}
        </div>
        <div class="form-inline" style="margin-bottom: 15px;">
            <Input
                style="margin-right: 5px;"
                placeholder="Guess Word"
                bind:value={wordGuess}
                disabled={guessWordLocked || gameOver || !canMove}
            />
            <Button
                disabled={guessWordLocked || gameOver || !canMove}
                on:click={guessWord}
            >Guess</Button>
        </div>
        <div class="d-flex">
            {#if gameOver}
                <Button
                    style="margin-right: 5px;"
                    on:click={newGame}
                >New Game</Button>
            {/if}
            <Button on:click={() => confirmEndGame = true}>
                {#if gameType === GameType.RemoteMultiplayer}
                    Leave Game
                {:else}
                    End Game
                {/if}
            </Button>
        </div>
    </div>
{/if}
