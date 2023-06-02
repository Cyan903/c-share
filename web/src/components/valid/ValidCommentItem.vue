<template>
    <label class="input-group">
        <input
            :value="modelValue"
            :class="{ 'input-error': !valid }"
            type="text"
            placeholder="File Comment"
            class="input input-bordered"
            @input="updateText"
        />

        <button
            @click.prevent="why"
            :disabled="valid"
            class="btn btn-error"
        >
            ?
        </button>
    </label>
</template>

<script lang="ts" setup>
import { computed, toRef } from "vue";
import { useValidComment } from "@/use/useValidate";
import Swal from "sweetalert2";

const props = defineProps<{
    modelValue: string;
}>();

const emit = defineEmits<{
    (e: "update:modelValue", evt: String): void;
}>();

const updateText = (evt: Event) => {
    emit("update:modelValue", (evt.target as HTMLInputElement).value);
};

const valid = computed(() => useValidComment(toRef(props.modelValue)));
const why = () => {
    Swal.fire({
        title: "Invalid comment!",
        text: "Must be between 1 and 100 characters.",
        icon: "warning",
        confirmButtonText: "Okay",
    });
};
</script>
