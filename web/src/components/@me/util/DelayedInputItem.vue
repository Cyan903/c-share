<template>
    <input
        type="text"
        :placeholder="placehold"
        :disabled="disabled"
        @input="updateInput"
    />
</template>

<script lang="ts" setup>
defineProps<{
    modelValue: string;
    placehold: string;
    disabled?: boolean;
}>();

let timer = 0;

const emit = defineEmits<{
    (e: "update:modelValue", evt: String): void;
}>();

const updateInput = (evt: Event) => {
    const text = (evt.target as HTMLInputElement).value;

    clearInterval(timer);
    timer = setTimeout(() => {
        emit("update:modelValue", text);
    }, 500);
};
</script>
