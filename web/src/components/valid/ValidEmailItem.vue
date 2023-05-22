<template>
    <div>
        <input
            type="text"
            :value="modelValue"
            :disabled="disabled"
            @input="updateText"
            placeholder="Email"
        />

        <button
            :class="{ invalid: !valid && !disabled }"
            @click.prevent="why"
            :disabled="valid || disabled"
        >
            ?
        </button>
    </div>
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
