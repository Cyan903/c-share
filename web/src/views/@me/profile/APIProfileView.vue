<template>
    <div>
        <LoadingAuthItem :loading="loading" />
        <div>
            <h2 class="font-semibold text-3xl my-4">API Tokens</h2>
            <EmailStatusItem />

            <p v-if="!auth.userData.emailVerified" class="link mb-4">
                <router-link to="/@me/profile/email">
                    Your email must be verified in order to request an API key.
                    You can send a verification code here.
                </router-link>
            </p>
            <p v-else class="m-4">
                Your email is verified! You may request an API key. Be sure to
                read the
                <router-link class="link" to="/about">documentation</router-link
                >.
            </p>
        </div>

        <div>
            <div class="input-group">
                <input
                    v-model="token"
                    class="input"
                    type="text"
                    placeholder="No token..."
                    readonly
                />

                <button
                    class="btn hidden md:inline-block"
                    :disabled="!auth.userData.emailVerified || token == ''"
                    @click="copyID"
                >
                    Copy
                </button>
            </div>

            <div class="my-4">
                <button
                    class="btn btn-outline w-full md:w-auto btn-success md:mx-2 my-2"
                    :disabled="!auth.userData.emailVerified"
                    @click="generateToken"
                >
                    New Token
                </button>
                <button
                    class="btn btn-outline w-full md:w-auto btn-error keep-error"
                    :disabled="!auth.userData.emailVerified || token == ''"
                    @click="deleteToken"
                >
                    Delete Token
                </button>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useRequest } from "@/use/useAPI";
import { useToast } from "vue-toastification";
import type { UserTokenData } from "@/types/api/@me/api";

import EmailStatusItem from "@/components/profile/EmailStatusItem.vue";
import LoadingAuthItem from "@/components/loading/LoadingAuthItem.vue";

import Swal from "sweetalert2";

const auth = useAuthStore();
const toast = useToast();

const loading = ref(false);
const token = ref("");
const date = ref("");

const copyID = () => {
    navigator.clipboard.writeText(token.value).then(
        () => {
            toast.success(`Copied ${token.value} to clipboard`);
        },
        (err) => {
            console.warn(`[copyID] could not copy id to clipboard, ${err}`);
            toast.error(`Could not copy ${token.value} to clipboard`);
        }
    );
};

const generateToken = () => {
    Swal.fire({
        title: "Are you sure you want to generate a new token?",
        text: token.value ? "Anything using this token will break!" : "",
        showCancelButton: true,
        confirmButtonText: "Yes",
    }).then(async (result) => {
        if (!result.isConfirmed) return;
        const req = await useRequest<UserTokenData>(
            "/@me/api/token",
            {
                method: "POST",
                headers: { Token: auth.token },
            },
            loading
        );

        if (req.error) return;
        if (req.response.status != 200) {
            Swal.fire({
                title: "Could generate a new token!",
                text: req.json.message,
                icon: "warning",
                confirmButtonText: "Okay",
            });

            return;
        }

        token.value = req.json.message;
        toast.success("New token has been generated!");
    });
};

const deleteToken = () => {
    Swal.fire({
        title: "Are you sure you want to delete this token?",
        text: token.value ? "Anything using this token will break!" : "",
        showCancelButton: true,
        confirmButtonText: "Yes",
    }).then(async (result) => {
        if (!result.isConfirmed) return;
        const req = await useRequest<UserTokenData>(
            "/@me/api/token",
            {
                method: "DELETE",
                headers: { Token: auth.token },
            },
            loading
        );

        if (req.error) return;
        if (req.response.status != 200) {
            Swal.fire({
                title: "Could delete your token!",
                text: req.json.message,
                icon: "warning",
                confirmButtonText: "Okay",
            });

            return;
        }

        token.value = "";
        toast.success("Token has been deleted!");
    });
};

onMounted(async () => {
    const wt = localStorage.getItem("token");
    if (!wt) return;

    const req = await useRequest<UserTokenData>(
        "/@me/api",
        {
            method: "GET",
            headers: { Token: wt },
        },
        loading
    );

    if (req.error) return;
    if (req.response.status != 200 && req.response.status != 401) {
        Swal.fire({
            title: "Could get API token information!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        return;
    }

    if (req.response.status != 401) {
        token.value = req.json.data.token;
        date.value = req.json.data.created_at;
    }
});
</script>
