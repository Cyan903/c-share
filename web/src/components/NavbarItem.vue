<template>
    <nav>
        <div>
            <router-link to="/">Home</router-link> |
            <router-link to="/about">About</router-link>
        </div>

        <div v-if="auth.isLoggedIn">
            <router-link to="/@me">Dashboard</router-link> |
            <router-link to="/@me/settings">Settings</router-link> |
            <a href="#" @click.prevent="logout">Logout</a>
        </div>
        <div v-else>
            <router-link to="/auth/login">Login</router-link> |
            <router-link to="/auth/register">Register</router-link>
        </div>
    </nav>
</template>

<script lang="ts" setup>
import { useAuthStore } from "@/stores/auth";
import { RouterLink, useRouter } from "vue-router";
import Swal from "sweetalert2";

const router = useRouter();
const auth = useAuthStore();

const logout = () => {
    Swal.fire({
        title: "Are you sure you want to logout?",
        showCancelButton: true,
        confirmButtonText: "Logout",
    }).then((result) => {
        if (result.isConfirmed) {
            auth.logout();
            router.push("/");
        }
    });
};
</script>

<style scoped>
nav a.router-link-exact-active {
    color: red;
}
</style>
