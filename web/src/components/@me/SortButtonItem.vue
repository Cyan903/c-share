<template>
    <button :disabled="disabled" @click="updateMode">{{ mode }}</button>
</template>

<script lang="ts" setup>
const props = defineProps<{
    mode: string;
    disabled: boolean;
}>();

const emits = defineEmits<{
    (e: "clicked", evt: string): void;
}>();

const updateMode = () => emits("clicked", props.mode == "asc" ? "desc" : "asc");
</script>
