<script lang="ts">
    import { fade, fly } from 'svelte/transition'
    import { quintOut } from 'svelte/easing'
</script>

<div
    class="modal d-block"
    in:fly={{ y: -50, duration: 300 }}
    out:fly={{ y: -50, duration: 300, easing: quintOut }}
>
<div class="modal-dialog">
    <div class="modal-content">
        <div class="modal-header">
            <h5 class="modal-title">
                <slot name="title" />
            </h5>
        </div>
        <div class="modal-body">
            <slot name="body" />
        </div>
        <div class="modal-footer">
            <slot name="footer" />
        </div>
    </div>
</div>
</div>
<div class="modal-backdrop show" transition:fade={{ duration: 150 }}></div>
