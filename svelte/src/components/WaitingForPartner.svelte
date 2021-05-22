<script lang="ts">
    export let gameId: string
</script>
<style>
    div {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        text-align: center;
    }
    .unselectable {
        user-select: none;
    }
</style>

<div>
    <h1 class="unselectable h3">Waiting for Partner</h1>
    <p>
        <strong class="unselectable">Game ID: </strong>{gameId}
    </p>
</div>
