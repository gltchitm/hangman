<script lang="ts">
    import { Button } from 'sveltestrap'
</script>

<Button
    on:click
    onmouseup="this.blur()"
    onmouseleave="this.blur()"
    {...$$props}
>
    <slot />
</Button>
