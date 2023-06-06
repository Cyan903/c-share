<template>
    <div class="bg-base-100 rounded-box shadow-xl p-6">
        <h1 class="sm:text-4xl text-3xl font-bold text-center">
            Documentation
        </h1>

        <h4 class="text-sm text-center font-bold opacity-60 my-5">
            Frequently Asked Questions
        </h4>

        <div class="w-full md:w-2/4 m-auto">
            <div
                class="collapse collapse-arrow bg-base-300 mb-4"
                @click="toggleOpen"
            >
                <div class="collapse-title text-xl font-medium">
                    How do I upload and share files?
                </div>
                <div class="collapse-content">
                    <p>
                        To upload files, simply visit your
                        <router-link to="/@me" class="link"
                            >dashboard</router-link
                        >. You can share a file by copying the file's ID. Please
                        note that private files cannot be shared, and unlisted
                        files require a password.
                    </p>
                </div>
            </div>

            <div
                class="collapse collapse-arrow bg-base-300 mb-4"
                @click="toggleOpen"
            >
                <div class="collapse-title text-xl font-medium">
                    What is the difference between a private file and a unlisted
                    file?
                </div>
                <div class="collapse-content">
                    <p>
                        A private file may only be view by you. You must be
                        logged in to view your private file. A unlisted file may
                        be shared but requires a password. You can change a
                        file's permissions in your
                        <router-link to="/@me" class="link">
                            dashboard.
                        </router-link>
                    </p>
                </div>
            </div>

            <div
                class="collapse collapse-arrow bg-base-300 mb-4"
                @click="toggleOpen"
            >
                <div class="collapse-title text-xl font-medium">
                    How can I automate file uploading?
                </div>
                <div class="collapse-content">
                    <p>Not possible at the moment.</p>
                </div>
            </div>

            <div
                class="collapse collapse-arrow bg-base-300 mb-4"
                @click="toggleOpen"
            >
                <div class="collapse-title text-xl font-medium">
                    Is this open source?
                </div>
                <div class="collapse-content">
                    <p>
                        Yes! You can view the source code
                        <a :href="source" class="link">here</a>.
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed } from "vue";

const source = computed(() => import.meta.env.VITE_SOURCE || "");

const toggleOpen = (e: Event) =>
    (e.target as HTMLElement).parentElement!.classList.toggle("collapse-open");
</script>
