<template>
    <div>
        <input
            type="password"
            :value="modelValue"
            :disabled="disabled"
            @input="updateText"
            placeholder="Confirm Password"
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
import { computed } from "vue";
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

const valid = computed(() => props.modelValue == props.check);
const why = () => {
    Swal.fire({
        title: "Passwords must match!",
        text: "Passwords do not match.",
        icon: "warning",
        confirmButtonText: "Okay",
    });
};
</script>
