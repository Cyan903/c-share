<template>
    <div>
        <button @click="updatePage(-1)" :disabled="disabled || page <= 0">
            &lt;
        </button>

        <span>{{ page + 1 }}</span>

        <button @click="updatePage(1)" :disabled="disabled">&gt;</button>
    </div>
</template>

<script lang="ts" setup>
const props = defineProps<{
    page: number;
    disabled: boolean;
}>();

const emits = defineEmits<{
    (e: "clicked", evt: number): void;
}>();

const updatePage = (n: number) => {
    if (props.page + n < 0) {
        return;
    }

    emits("clicked", props.page + n);
};
</script>
