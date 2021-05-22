<script>
    import Button from './Button.svelte'
</script>
<style>
    div {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        text-align: center;
    }
</style>

<div>
    <h4>Socket Closed</h4>
    <Button block on:click={() => window.location.reload()}>Reload</Button>
</div>
