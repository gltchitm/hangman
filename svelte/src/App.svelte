<script lang="ts">
    import { FormGroup, Input, Label, Modal, ModalBody, ModalFooter, ModalHeader } from 'sveltestrap'

    import Welcome from './views/Welcome.svelte'
    import Game from './views/Game.svelte'

    import SocketClosed from './components/SocketClosed.svelte'
    import Button from './components/Button.svelte'

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

<Modal isOpen={showJoinOrCreateGameModal} transitionOptions={{ duration: 100 }}>
    <ModalHeader
        toggle={() => showJoinOrCreateGameModal = false}
    >Remote Multiplayer</ModalHeader>
    <ModalFooter>
        <Button
            on:click={joinOrCreateGameSelection(true)}
            color="secondary"
        >Create Game</Button>
        <Button
            on:click={joinOrCreateGameSelection(false)}
            color="secondary"
        >Join Game</Button>
    </ModalFooter>
</Modal>
<Modal isOpen={showJoinGameModal} transitionOptions={{ duration: 100 }}>
    <ModalHeader
        toggle={() => showJoinGameModal = false}
    >Remote Multiplayer</ModalHeader>
    <ModalBody>
        <FormGroup>
            <Label>Game ID</Label>
            <Input bind:value={joinGame} type="text" />
        </FormGroup>
    </ModalBody>
    <ModalFooter>
        <Button
            on:click={() => showJoinGameModal = false}
            color="secondary"
        >Cancel</Button>
        <Button
            on:click={join}
            color="secondary"
        >Join</Button>
    </ModalFooter>
</Modal>
<Modal isOpen={showCannotJoinModal} transitionOptions={{ duration: 100 }}>
    <ModalHeader>Cannot Join Game</ModalHeader>
    <ModalBody>You may have entered an invalid ID.</ModalBody>
    <ModalFooter>
        <Button
            on:click={() => showCannotJoinModal = false}
            color="secondary"
        >Close</Button>
    </ModalFooter>
</Modal>
{#if socketClosed}
    <SocketClosed />
{:else if view === View.Welcome}
    <Welcome on:select={gametypeSelect} />
{:else if view === View.Game}
    <Game {createGame} {gameType} {joinGame} />
{/if}
