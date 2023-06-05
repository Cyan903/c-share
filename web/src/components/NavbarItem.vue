<template>
    <div class="navbar bg-base-100 rounded-box my-4 shadow-xl">
        <div class="flex-1">
            <div class="dropdown">
                <label tabindex="0" class="btn btn-ghost btn-circle">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-5 w-5"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M4 6h16M4 12h16M4 18h7"
                        />
                    </svg>
                </label>
                <ul
                    tabindex="0"
                    class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-52"
                >
                    <li>
                        <router-link to="/">Home</router-link>
                    </li>

                    <li>
                        <router-link to="/about">About</router-link>
                    </li>

                    <div v-if="auth.isLoggedIn">
                        <li class="lg:hidden">
                            <router-link to="/@me">Dashboard</router-link>
                        </li>

                        <li class="lg:hidden">
                            <router-link to="/@me/profile/email">
                                Settings
                            </router-link>
                        </li>

                        <li class="lg:hidden">
                            <a href="#" @click.prevent="logout">Logout</a>
                        </li>
                    </div>
                    <div v-else>
                        <li class="lg:hidden">
                            <router-link to="/auth/login">Login</router-link>
                        </li>

                        <li class="lg:hidden">
                            <router-link to="/auth/register">
                                Register
                            </router-link>
                        </li>
                    </div>
                </ul>
            </div>

            <router-link
                class="logo btn btn-ghost normal-case text-xl"
                :to="auth.isLoggedIn ? '/@me' : '/'"
            >
                {{ appName }}
            </router-link>
        </div>
        <div class="flex-none">
            <ul
                v-if="auth.isLoggedIn"
                class="menu menu-horizontal px-1 hidden lg:flex"
            >
                <li>
                    <router-link to="/@me">Dashboard</router-link>
                </li>

                <li tabindex="0">
                    <a>
                        {{ auth.userData.nickname }}
                        <svg
                            class="fill-current"
                            xmlns="http://www.w3.org/2000/svg"
                            width="20"
                            height="20"
                            viewBox="0 0 24 24"
                        >
                            <path
                                d="M7.41,8.58L12,13.17L16.59,8.58L18,10L12,16L6,10L7.41,8.58Z"
                            />
                        </svg>
                    </a>
                    <ul class="p-2 bg-base-100">
                        <li>
                            <router-link to="/@me/profile/email">
                                Settings
                            </router-link>
                        </li>

                        <li>
                            <a href="#" @click.prevent="logout">Logout</a>
                        </li>
                    </ul>
                </li>
            </ul>
            <ul v-else class="menu menu-horizontal px-1 hidden lg:flex">
                <li>
                    <router-link to="/auth/login">Login</router-link>
                </li>

                <li>
                    <router-link to="/auth/register">Register</router-link>
                </li>
            </ul>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useAuthStore } from "@/stores/auth";
import { RouterLink, useRouter } from "vue-router";
import { computed } from "vue";

import Swal from "sweetalert2";

const router = useRouter();
const auth = useAuthStore();

const appName = computed(() => import.meta.env.VITE_APP || "c-share");
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
.router-link-exact-active:not(.logo) {
    color: rgba(255, 255, 255, 0.1);
}

.navbar * {
    z-index: 1;
}
</style>
