<template>
    <teleport to="#modal">
        <div v-if="props.loading" class="modal-background"></div>
        <dialog
            class="modal w-11/12 max-w-2xl"
            :class="{ 'modal-open': props.loading }"
        >
            <img
                class="animate-spin"
                src="/svg/loading/loading-spinner.svg"
                width="100"
            />
        </dialog>
    </teleport>
</template>

<script lang="ts" setup>
const props = defineProps<{
    loading: boolean;
}>();
</script>
