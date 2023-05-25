<template>
    <button @click="updateMode">{{ mode }}</button>
</template>

<script lang="ts" setup>
const props = defineProps<{
    mode: string;
}>();

const emits = defineEmits<{
    (e: "clicked", evt: string): void;
}>();

const updateMode = () => emits("clicked", props.mode == "asc" ? "desc" : "asc");
</script>
