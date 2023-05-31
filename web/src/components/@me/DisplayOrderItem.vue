<template>
    <div>
        <span>{{ total }}</span> | <span>Order: {{ order }}</span> |
        <span>Listing: {{ listing }}</span> |
        <span>Type: {{ type || "Any" }}</span>
    </div>
</template>

<script lang="ts" setup>
defineProps<{
    total: number;
    order: string;
    listing: string;
    type: string;
}>();
</script>

<style scoped>
span {
    text-transform: capitalize;
}
</style>
