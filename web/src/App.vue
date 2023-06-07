<template>
    <NavbarItem />

    <main>
        <router-view v-slot="{ Component }">
            <transition name="router" mode="out-in">
                <component :is="Component"></component>
            </transition>
        </router-view>
    </main>

    <teleport to="body">
        <FooterItem />
    </teleport>
</template>

<script lang="ts" setup>
import { RouterView } from "vue-router";

import NavbarItem from "@/components/NavbarItem.vue";
import FooterItem from "@/components/FooterItem.vue";
</script>

<style scoped>
.router-enter-from {
    opacity: 0;
}

.router-enter-active {
    transition: all 0.1s ease-in-out;
}
</style>
