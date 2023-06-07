<template>
    <label class="input-group">
        <input
            :value="modelValue"
            :disabled="disabled"
            :placeholder="placehold || 'Password'"
            :class="{ 'input-error': !valid && !disabled }"
            type="password"
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
import { useValidPassword } from "@/use/useValidate";
import Swal from "sweetalert2";

const props = defineProps<{
    modelValue: string;
    disabled?: boolean;
    placehold?: string;
}>();

const emit = defineEmits<{
    (e: "update:modelValue", evt: String): void;
}>();

const updateText = (evt: Event) => {
    emit("update:modelValue", (evt.target as HTMLInputElement).value);
};

const valid = computed(() => useValidPassword(toRef(props.modelValue)));
const why = () => {
    Swal.fire({
        title: "Invalid Password!",
        text: "Must be between 7 and 30 characters.",
        icon: "warning",
        confirmButtonText: "Okay",
    });
};
</script>
