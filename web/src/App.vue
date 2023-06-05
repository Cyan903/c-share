<template>
    <NavbarItem />

    <main>
        <RouterView />
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
