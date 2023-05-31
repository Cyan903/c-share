<template>
    <NavbarItem />

    <main>
        <RouterView />
    </main>
</template>

<script lang="ts" setup>
import { RouterView } from "vue-router";
import NavbarItem from "@/components/NavbarItem.vue";
</script>
