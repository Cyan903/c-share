<template>
    <footer class="footer items-center p-4 bg-base-300">
        <div class="items-center grid-flow-col">
            <img src="/assets/footer.png" width="35" height="35" />
            <p>{{ appName }} - fast & efficient file hosting</p>
        </div>
        <div
            v-if="appSource"
            class="grid-flow-col gap-4 hidden md:flex md:place-self-center md:justify-self-end"
        >
            <a :href="appSource" target="_blank">
                <img
                    src="/ico/source.png"
                    height="24"
                    width="24"
                    alt="source"
                />
            </a>
        </div>
    </footer>
</template>

<script lang="ts" setup>
import { computed } from "vue";

const appName = computed(() => import.meta.env.VITE_APP);
const appSource = computed(() => import.meta.env.VITE_SOURCE);
</script>
