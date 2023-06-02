<template>
    <label class="input-group">
        <input
            :value="modelValue"
            :disabled="disabled"
            :class="{ 'input-error': !valid && !disabled }"
            type="text"
            placeholder="Email"
            class="input input-bordered"
            @input="updateText"
        />

        <button
            @click.prevent="why"
            :disabled="valid || disabled"
            class="btn btn-error"
        >
            ?
        </button>
    </label>
</template>

<script lang="ts" setup>
import { computed, toRef } from "vue";
import { useValidEmail } from "@/use/useValidate";
import Swal from "sweetalert2";

const props = defineProps<{
    modelValue: string;
    disabled?: boolean;
}>();

const emit = defineEmits<{
    (e: "update:modelValue", evt: String): void;
}>();

const updateText = (evt: Event) => {
    emit("update:modelValue", (evt.target as HTMLInputElement).value);
};

const valid = computed(() => useValidEmail(toRef(props.modelValue)));
const why = () => {
    Swal.fire({
        title: "Invalid Email!",
        text: "Must be a valid address between 7 and 30 in length.",
        icon: "warning",
        confirmButtonText: "Okay",
    });
};
</script>
