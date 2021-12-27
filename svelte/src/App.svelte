<script lang="ts">
    import Welcome from './views/Welcome.svelte'
    import Game from './views/Game.svelte'

    import SocketClosed from './components/SocketClosed.svelte'
    import Modal from './components/Modal.svelte'

    import { GameType } from './enum'
    import { socket } from './store'

    enum View {
        Welcome,
        JoinOrCreateGame,
        Game
    }

    let view = View.Welcome
    let showJoinOrCreateGameModal = false
    let showJoinGameModal = false
    let showCannotJoinModal = false
    let createGame = false
    let socketClosed = false
    let gameType: GameType
    let joinGame: string

    $socket.addEventListener('message', ({ data }) => {
        const payload = JSON.parse(data)
        if (payload.message === 'cannot_join_game') {
            view = View.Welcome
            showCannotJoinModal = true
        }
    })
    $socket.addEventListener('close', () => {
        socketClosed = true
    })

    const gametypeSelect = ({ detail }) => {
        gameType = detail
        if (gameType === GameType.RemoteMultiplayer) {
            showJoinOrCreateGameModal = true
        } else {
            view = View.Game
        }
    }
    const joinOrCreateGameSelection = (result: boolean) => {
        return () => {
            showJoinOrCreateGameModal = false
            createGame = result
            if (result) {
                view = View.Game
            } else {
                showJoinGameModal = true
            }
        }
    }
    const join = () => {
        showJoinGameModal = false
        view = View.Game
    }

    setInterval(() => {
        $socket.send(JSON.stringify({ action: 'ping' }))
    }, 5000)
</script>

{#if showJoinOrCreateGameModal}
    <Modal>
        <svelte:fragment slot="title">Remote Multiplayer</svelte:fragment>
        <svelte:fragment slot="body">
            Create or join game?
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                on:click={joinOrCreateGameSelection(true)}
            >Create Game</button>
            <button
                class="btn btn-secondary"
                on:click={joinOrCreateGameSelection(false)}
            >Join Game</button>
        </svelte:fragment>
    </Modal>
{/if}
{#if showJoinGameModal}
    <Modal>
        <svelte:fragment slot="title">Remote Multiplayer</svelte:fragment>
        <svelte:fragment slot="body">
            <label for="gameId">Game ID</label>
            <input class="form-control" id="gameId" type="text" bind:value={joinGame} />
        </svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                class="btn btn-secondary"
                on:click={() => showJoinGameModal = false}
            >Cancel</button>
            <button
                class="btn btn-secondary"
                on:click={join}
            >Join</button>
        </svelte:fragment>
    </Modal>
{/if}
{#if showCannotJoinModal}
    <Modal>
        <svelte:fragment slot="title">Cannot Join Game</svelte:fragment>
        <svelte:fragment slot="body">You may have entered an invalid ID.</svelte:fragment>
        <svelte:fragment slot="footer">
            <button
                on:click={() => showCannotJoinModal = false}
                class="btn btn-secondary"
            >Close</button>
        </svelte:fragment>
    </Modal>
{/if}
{#if socketClosed}
    <SocketClosed />
{:else if view === View.Welcome}
    <Welcome on:select={gametypeSelect} />
{:else if view === View.Game}
    <Game {createGame} {gameType} {joinGame} />
{/if}
