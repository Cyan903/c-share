<script setup lang="ts">
import { RouterLink, RouterView } from "vue-router";
</script>

<template>
    <nav>
        <div>
            <RouterLink to="/">Home</RouterLink> |
            <RouterLink to="/about">About</RouterLink>
        </div>

        <div>
            <RouterLink to="/auth/login">Login</RouterLink> |
            <RouterLink to="/auth/register">Register</RouterLink>
        </div>
    </nav>

    <main>
        <RouterView />
    </main>
</template>

<style scoped>
nav a.router-link-exact-active {
    color: red;
}
</style>
