<template>
    <input
        v-if="!label"
        type="text"
        :placeholder="placehold"
        :disabled="disabled"
        :class="classes || ''"
        :value="modelValue"
        @input="updateInput"
    />

    <label v-else class="input-group" :class="classes || ''">
        <span>{{ label }}</span>
        <input
            type="text"
            class="input input-bordered"
            :placeholder="placehold"
            :disabled="disabled"
            :value="modelValue"
            @input="updateInput"
        />
    </label>
</template>

<script lang="ts" setup>
let timer = 0;

defineProps<{
    modelValue: string;
    placehold: string;
    classes?: string;
    label?: string;
    disabled?: boolean;
}>();

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
