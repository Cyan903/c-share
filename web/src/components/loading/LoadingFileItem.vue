<template>
    <tr v-if="props.loading" class="hidden lg:block mb-20 mt-3">
        <img
            class="animate-spin text-center spinner-table"
            src="/svg/loading/loading-spinner.svg"
            width="80"
        />
    </tr>
    <tr v-if="props.loading" class="lg:hidden block">
        <img class="my-4" src="/svg/loading/loading-dots.svg" width="40" />
    </tr>
</template>

<script lang="ts" setup>
const props = defineProps<{
    loading: boolean;
}>();
</script>

<style scoped>
.spinner-table {
    position: absolute;
    margin: auto;
    left: 0;
    right: 0;
}
</style>
