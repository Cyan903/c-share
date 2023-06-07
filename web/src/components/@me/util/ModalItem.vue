<template>
    <teleport to="#modal">
        <div v-if="show" class="modal-background" @click="$emit('hide')"></div>
        <dialog class="modal w-11/12 max-w-2xl" :class="{ 'modal-open': show }">
            <div v-if="show" class="modal-box">
                <button
                    @click="$emit('hide')"
                    class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
                >
                    ✕
                </button>

                <slot></slot>

                <div class="modal-action">
                    <button @click="$emit('hide')" class="btn">Close</button>
                </div>
            </div>
        </dialog>
    </teleport>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted } from "vue";

defineProps<{
    show: boolean;
}>();

const emit = defineEmits<{
    (e: "hide"): void;
}>();

const escape = (e: KeyboardEvent) => (e.key == "Escape" ? emit("hide") : 0);

onMounted(() => window.addEventListener("keydown", escape));
onUnmounted(() => window.removeEventListener("keydown", escape));
</script>

<style scoped>
.c-modal {
    background-color: rgba(0, 0, 0, 0.3);
}
</style>
