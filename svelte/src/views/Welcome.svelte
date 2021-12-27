<script>
    import { createEventDispatcher } from 'svelte'

    import { GameType } from '../enum'

    const dispatcher = createEventDispatcher()

    const select = gamemode => () => dispatcher('select', gamemode)
</script>
<style>
    .main-container {
        height: 100vh;
        width: 100vw;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    .welcome {
        width: 200px;
    }
    .hangman-label {
        font-size: 40px;
        text-align: center;
    }
</style>

<div class="main-container">
    <div class="welcome">
        <h1 class="hangman-label">Hangman</h1>
        <button
            class="btn btn-secondary btn-block"
            on:click={select(GameType.LocalSolo)}
        >Local Solo</button>
        <button
            class="btn btn-secondary btn-block"
            on:click={select(GameType.LocalMultiplayer)}
        >Local Multiplayer</button>
        <button
            class="btn btn-secondary btn-block"
            on:click={select(GameType.RemoteMultiplayer)}
        >Remote Multiplayer</button>
    </div>
</div>
