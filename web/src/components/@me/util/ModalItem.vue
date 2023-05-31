<template>
    <teleport to="#modal">
        <div v-if="show">
            <div class="modal-background" @click="$emit('hide')"></div>
            <div class="modal-content">
                <button class="modal-close" @click="$emit('hide')">x</button>
                <slot></slot>
            </div>
        </div>
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
.modal-background {
    background-color: rgba(0, 0, 0, 0.5);
    width: 100vw;
    height: 100vh;
    position: fixed;
    top: 0;
}

.modal-content {
    background-color: #fff;
    max-width: 40vw;
    max-height: 80vh;
    padding: 1em;
    position: fixed;
    margin-left: auto;
    margin-right: auto;
    left: 0;
    right: 0;
    top: 10vh;
    overflow-y: scroll;
}

.modal-close {
    float: right;
}
</style>
