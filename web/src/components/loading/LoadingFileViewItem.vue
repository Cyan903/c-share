<template>
    <img
        v-if="props.loading"
        src="/svg/loading/loading-dots.svg"
        width="30"
        class="m-auto my-4"
    />
</template>

<script lang="ts" setup>
const props = defineProps<{
    loading: boolean;
}>();
</script>
