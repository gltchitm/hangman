<script>
    import { Container } from 'sveltestrap'
    import { createEventDispatcher } from 'svelte'

    import Button from '../components/Button.svelte'
    
    import { GameType } from '../enum'

    const dispatcher = createEventDispatcher()

    const select = gamemode => () => dispatcher('select', gamemode)
</script>
<style>
    div {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
    }
</style>

<div>
    <Container class="welcome">
        <h1 class="h3">Hangman</h1>
        <Button
            on:click={select(GameType.LocalSolo)}
            block
        >Local Solo</Button>
        <Button
            on:click={select(GameType.LocalMultiplayer)}
            block
        >Local Multiplayer</Button>
        <Button
            on:click={select(GameType.RemoteMultiplayer)}
            block
        >Remote Multiplayer</Button>
    </Container>
</div>
