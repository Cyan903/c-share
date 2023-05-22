<template>
    <div>
        <input
            type="text"
            :value="modelValue"
            @input="updateText"
            placeholder="Nickname"
        />

        <button
            :class="{ invalid: !valid }"
            @click.prevent="why"
            :disabled="valid"
        >
            ?
        </button>
    </div>
</template>

<script lang="ts" setup>
import { computed, toRef } from "vue";
import { useValidNickname } from "@/use/useValidate";
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

const valid = computed(() => useValidNickname(toRef(props.modelValue)));
const why = () => {
    Swal.fire({
        title: "Invalid Nickname!",
        text: "Must be between 4 and 10 characters. Cannot use numbers or special characters.",
        icon: "warning",
        confirmButtonText: "Okay",
    });
};
</script>
