<template>
    <div>
        <input
            type="password"
            :value="modelValue"
            :disabled="disabled"
            @input="updateText"
            :placeholder="placehold || 'Password'"
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
