<style>
    .socket-closed {
        height: 100vh;
        width: 100vw;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }
    .reload-btn {
        width: 200px;
    }
</style>

<div class="socket-closed">
    <h4>Socket Closed</h4>
    <button
        class="btn btn-secondary btn-block reload-btn"
        on:click={() => window.location.reload()}
    >Reload</button>
</div>
