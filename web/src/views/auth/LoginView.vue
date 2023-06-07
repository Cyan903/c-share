<template>
    <div class="auth-form flex flex-wrap items-center justify-center">
        <LoadingAuthItem :loading="loading" />

        <div class="card card-normal bg-base-300 shadow-xl">
            <div class="card-body items-center text-center">
                <h2 class="card-title text-3xl">Login</h2>

                <ValidEmailItem v-model="email" />
                <ValidPasswordItem v-model="password" />

                <router-link class="link link-info" to="/auth/pwreset">
                    Forgot Password?
                </router-link>

                <input
                    class="btn btn-primary w-5/6 my-4"
                    :disabled="!valid"
                    type="submit"
                    value="Login"
                    @click.prevent="login"
                />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue";
import type { Login } from "@/types/api/auth";

import { useValidEmail, useValidPassword } from "@/use/useValidate";
import { useRequest } from "@/use/useAPI";
import { useAuthStore } from "@/stores/auth";
import { useRouter } from "vue-router";

import ValidEmailItem from "@/components/valid/ValidEmailItem.vue";
import ValidPasswordItem from "@/components/valid/ValidPasswordItem.vue";
import LoadingAuthItem from "@/components/loading/LoadingAuthItem.vue";

import Swal from "sweetalert2";

const router = useRouter();
const auth = useAuthStore();

const email = ref("");
const password = ref("");
const loading = ref(false);

const valid = computed(
    () => useValidEmail(email) && useValidPassword(password)
);

const login = async () => {
    const req = await useRequest<Login>(
        "/auth/login",
        {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                email: email.value,
                password: password.value,
            }),
        },

        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not login!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        password.value = "";
        return;
    }

    if (await auth.login(req.json.data)) {
        router.push("/@me");
        return;
    }

    Swal.fire({
        title: "Could not obtain user information.",
        text: "Please report this bug.",
        icon: "error",
        confirmButtonText: "Okay",
    });
};
</script>
