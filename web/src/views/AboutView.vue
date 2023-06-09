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
                    You can with the API. You can find more information down
                    <a class="link" href="#fileUpload">below</a>.
                </div>
            </div>

            <div
                class="collapse collapse-arrow bg-base-300 mb-4"
                @click="toggleOpen"
                v-if="source"
            >
                <div class="collapse-title text-xl font-medium">
                    Is this open source?
                </div>
                <div class="collapse-content">
                    <p>
                        Yes, you can view the source code
                        <a :href="source" class="link">here</a>.
                    </p>
                </div>
            </div>
        </div>

        <div class="divider"></div>
        <h1 class="sm:text-4xl text-3xl font-bold text-center">
            API Reference
        </h1>

        <h4 class="text-sm text-center font-bold opacity-60 my-5">
            File Upload/Delete API
        </h4>

        <div class="w-full md:w-2/4 m-auto">
            <p>
                An API key is required to interact with the upload API. You can
                request one in
                <router-link to="/@me/profile/api" class="link"
                    >settings</router-link
                >. Please note that your email must be verified in order to use
                API tokens.
            </p>

            <div v-if="auth.isLoggedIn" class="text-center my-4 font-bold">
                <p v-if="!auth.userData.emailVerified">
                    Your email is <span class="text-error">not verified</span>.
                    You can verify it

                    <router-link to="/@me/profile/email" class="link">
                        here.
                    </router-link>
                </p>
                <p v-else>
                    Your email is
                    <span class="text-success">verified!</span> You may interact
                    with the upload API.
                </p>
            </div>

            <div>
                <div class="divider"></div>
                <h1 class="text-2xl font-bold">
                    <a class="opacity-20 text-lg font-bold">#</a>
                    Upload API
                </h1>

                <p class="my-4" id="fileUpload">
                    <code class="rounded bg-base-300 p-1"
                        >POST - /api/upload</code
                    >
                    Upload files endpoint. Here you can automate file uploading
                    with a token.
                </p>

                <div class="overflow-x-auto my-5">
                    <table class="w-full" bordered>
                        <thead>
                            <th>Query</th>
                            <th>Value</th>
                        </thead>
                        <tbody>
                            <tr>
                                <td>?token</td>
                                <td>"api_token"</td>
                            </tr>

                            <tr>
                                <td>&perm</td>
                                <td>"public" | "private" | "unlisted"</td>
                            </tr>

                            <tr>
                                <td>&password (ONLY if perm == unlisted)</td>
                                <td>"your_password"</td>
                            </tr>

                            <tr>
                                <td>&comment</td>
                                <td>"your_comment'</td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div>
                    <h4 class="font-bold text-lg">Request Body</h4>
                    <p>
                        The request body only contains the file being uploaded.
                        All other data is in the URL parameters.
                    </p>

                    <div class="alert my-4">
                        <span>
                            <code class="rounded bg-base-300 p-1">upload</code>
                            - multipart/form-data
                        </span>
                    </div>
                </div>

                <RequestViewerItem
                    defaultCode="Public"
                    :response="Upload.response"
                    :code="Upload.code"
                />
            </div>

            <div>
                <div class="divider"></div>
                <h1 class="text-2xl font-bold">
                    <a class="opacity-20 text-lg font-bold">#</a>
                    Delete API
                </h1>

                <p class="my-4" id="fileUpload">
                    <code class="rounded bg-base-300 p-1"
                        >DELETE - /api/upload</code
                    >
                    Delete files endpoint. Here you can remove multiple files
                    with a token.
                </p>

                <div class="overflow-x-auto my-5">
                    <table class="w-full" bordered>
                        <thead>
                            <th>Query</th>
                            <th>Value</th>
                        </thead>
                        <tbody>
                            <tr>
                                <td>?token</td>
                                <td>"api_token"</td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div>
                    <h4 class="font-bold text-lg">Request Body</h4>
                    <p>
                        The request body only contains a JSON array of file IDs
                        to be removed.
                    </p>

                    <div class="alert my-4">
                        <span>
                            <code class="rounded bg-base-300 p-1"
                                >["abc123", "def456", ...]</code
                            >
                            - application/json
                        </span>
                    </div>
                </div>

                <RequestViewerItem
                    defaultCode="One File"
                    :response="Delete.response"
                    :code="Delete.code"
                />
            </div>
        </div>

        <div v-if="source" class="text-center my-5 mt-10">
            <a class="link hover:opacity-50 text-sm" :href="source">
                Further documentation here
            </a>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useAuthStore } from "@/stores/auth";
import { computed } from "vue";

import Delete from "@/docs/delete";
import Upload from "@/docs/upload";

import RequestViewerItem from "@/components/about/RequestViewerItem.vue";

const auth = useAuthStore();
const source = computed(() => import.meta.env.VITE_SOURCE || "");

const toggleOpen = (e: Event) =>
    (e.target as HTMLElement).parentElement!.classList.toggle("collapse-open");
</script>

<style scoped>
table * {
    border: 1px solid rgba(255, 255, 255, 0.3);
    padding: 0.5rem;
}
</style>
