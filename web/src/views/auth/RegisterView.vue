<template>
    <div>
        <Loading :loading="loading" />

        <h1>Register</h1>
        <form class="auth-form">
            <input type="text" v-model="nick" placeholder="Nickname" v-focus />
            <input type="text" v-model="email" placeholder="Email" />
            <input type="password" v-model="password" placeholder="Password" />
            <input
                :disabled="!valid"
                type="submit"
                value="Submit"
                @click.prevent="register"
            />
        </form>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import {
    useValidEmail,
    useValidPassword,
    useValidNickname,
} from "@/use/useValidate";

import { useRequest } from "@/use/useAPI";
import { useAuthStore } from "@/stores/auth";
import { vFocus } from "@/directives/vFocus";
import type { Register } from "@/types/api/auth";

import Loading from "@/components/LoadingItem.vue";
import Swal from "sweetalert2";

// TODO: Should probably indicate why a username/password is invalid.
// TODO: Confirm password

const router = useRouter();
const auth = useAuthStore();

const nick = ref("");
const email = ref("");
const password = ref("");

const loading = ref(false);
const valid = computed(
    () =>
        useValidNickname(nick) &&
        useValidEmail(email) &&
        useValidPassword(password)
);

const register = async () => {
    const req = await useRequest<Register>(
        "/auth/register",
        {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                email: email.value,
                nickname: nick.value,
                password: password.value,
            }),
        },
        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not register!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        email.value = "";
        password.value = "";
        return;
    }

    if (await auth.login(req.json.data)) {
        router.push("/@me");
        return;
    }
};
</script>
