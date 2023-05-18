<template>
    <div>
        <Loading :loading="loading" />

        <h1>Login</h1>
        <form class="auth-form">
            <input type="text" v-model="email" placeholder="Email" v-focus />
            <input type="password" v-model="password" placeholder="Password" />
            <input
                :disabled="!invalid"
                type="submit"
                value="Submit"
                @click.prevent="login"
            />
        </form>

        <router-link to="/auth/pwreset">Forgot Password?</router-link>
        <pre>{{ auth.userData }}</pre>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue";
import type { Login } from "@/types/api/auth";

import { useValidEmail, useValidPassword } from "@/use/useValidate";
import { useRequest } from "@/use/useAPI";
import { useAuthStore } from "@/stores/auth";
import { vFocus } from "@/directives/vFocus";
import { useRouter } from "vue-router";

import Loading from "@/components/LoadingItem.vue";
import Swal from "sweetalert2";

const router = useRouter();
const auth = useAuthStore();

const email = ref("");
const password = ref("");
const loading = ref(false);

const invalid = computed(
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
