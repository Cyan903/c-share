<template>
    <h6 v-if="props.loading">Loading...</h6>
</template>

<script lang="ts" setup>
const props = defineProps<{
    loading: boolean;
}>();
</script>
