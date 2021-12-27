<script lang="ts">
    export let gameId: string
</script>
<style>
    .waiting-for-partner {
        height: 100vh;
        width: 100vw;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }
    .unselectable {
        user-select: none;
    }
</style>

<div class="waiting-for-partner">
    <h1 class="unselectable h3">Waiting for Partner</h1>
    <p>
        <strong class="unselectable">Game ID: </strong>{gameId}
    </p>
    <button
        class="btn btn-secondary"
        on:click={() => navigator.clipboard.writeText(gameId)}
    >Copy</button>
</div>
