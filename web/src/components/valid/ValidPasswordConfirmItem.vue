<template>
    <label class="input-group">
        <input
            :value="modelValue"
            :disabled="disabled"
            :class="{ 'input-error': !valid && !disabled }"
            type="password"
            placeholder="Confirm Password"
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
    check: string;
    disabled?: boolean;
}>();

const emit = defineEmits<{
    (e: "update:modelValue", evt: String): void;
}>();

const updateText = (evt: Event) => {
    emit("update:modelValue", (evt.target as HTMLInputElement).value);
};

const valid = computed(() => useValidPassword(toRef(props.modelValue)) && props.modelValue == props.check);
const why = () => {
    Swal.fire({
        title: "Passwords must match!",
        text: "Passwords do not match.",
        icon: "warning",
        confirmButtonText: "Okay",
    });
};
</script>
