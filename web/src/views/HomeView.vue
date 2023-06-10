<template>
    <div>
        <div class="hero bg-base-200 relative bottom-2">
            <div class="hero-content flex-col lg:flex-row">
                <img
                    src="/assets/landing.png"
                    class="max-w-sm rounded-lg hidden lg:block"
                />
                <div>
                    <h1 class="sm:text-5xl text-4xl font-bold">
                        File Sharing Made Easy
                    </h1>
                    <p class="py-6 sm:text-base text-sm">
                        {{ appName }} is a secure password protected file
                        hosting service. Upload and transfer files quickly and
                        efficiently. {{ appName }} is open source and
                        self-hostable to anyone.
                    </p>

                    <div
                        v-if="!auth.isLoggedIn"
                        class="flex flex-wrap flex-column"
                    >
                        <router-link
                            class="btn btn-primary sm:w-auto w-full mb-3"
                            to="/auth/register"
                        >
                            Get Started
                        </router-link>

                        <router-link
                            class="btn btn-secondary mx-0 sm:mx-4 sm:w-auto w-full"
                            to="/about"
                        >
                            Documentation
                        </router-link>
                    </div>
                    <div v-else>
                        <router-link
                            class="btn btn-primary sm:w-auto w-full mb-3"
                            to="/@me"
                        >
                            Dashboard
                        </router-link>

                        <router-link
                            class="btn btn-secondary mx-0 sm:mx-4 sm:w-auto w-full"
                            to="/@me/profile/email"
                        >
                            Settings
                        </router-link>
                    </div>
                </div>
            </div>
        </div>

        <div class="divider"></div>
        <div class="text-center my-10">
            <div class="stats stats-vertical lg:stats-horizontal shadow">
                <div class="stat">
                    <div class="stat-figure text-secondary">
                        <img src="/svg/storage.svg" width="24" />
                    </div>
                    <div class="stat-title">Used Storage</div>
                    <div class="stat-value">{{ data.storage }}</div>
                </div>

                <div class="stat">
                    <div class="stat-figure text-secondary">
                        <img src="/svg/users.svg" width="24" />
                    </div>
                    <div class="stat-title">New Users</div>
                    <div class="stat-value">{{ data.users }}</div>
                </div>

                <div class="stat">
                    <div class="stat-figure text-secondary">
                        <img src="/svg/files.svg" width="24" />
                    </div>
                    <div class="stat-title">Total Files</div>
                    <div class="stat-value">{{ data.files }}</div>
                </div>
            </div>

            <h4 class="text-base-content font-bold my-4 text-sm opacity-70">
                <a :href="appSource" class="hover:opacity-60"
                    >Powered by {{ appName }}</a
                >
            </h4>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useRequest } from "@/use/useAPI";
import { useStorage } from "@/use/useStorage";
import type { ServerInfo } from "@/types/api/general";

import Swal from "sweetalert2";

const auth = useAuthStore();
const appName = computed(() => import.meta.env.VITE_APP || "c-share");
const appSource = computed(() => import.meta.env.VITE_SOURCE || "#");

const data = reactive({
    storage: "0",
    users: 0,
    files: 0,
});

onMounted(async () => {
    const req = await useRequest<ServerInfo>(
        "/",
        {
            method: "GET",
            headers: { "Content-Type": "application/json" },
        },
        ref(false)
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not server info!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        return;
    }

    data.users = req.json.data.users;
    data.storage = useStorage(parseInt(req.json.data.storage));
    data.files = req.json.data.total_files;
});
</script>

<style>
.hero {
    min-height: 80vh !important;
}
</style>
