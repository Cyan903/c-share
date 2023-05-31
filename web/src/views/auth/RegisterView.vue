<template>
    <div>
        <Loading :loading="loading" />

        <h1>Register</h1>
        <form>
            <ValidNicknameItem v-model="nick" />
            <ValidEmailItem v-model="email" />
            <ValidPasswordItem v-model="password" />
            <ValidPasswordConfirmItem
                v-model="passwordConfirm"
                :check="password"
            />

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

import ValidNicknameItem from "@/components/valid/ValidNicknameItem.vue";
import ValidEmailItem from "@/components/valid/ValidEmailItem.vue";
import ValidPasswordItem from "@/components/valid/ValidPasswordItem.vue";
import ValidPasswordConfirmItem from "@/components/valid/ValidPasswordConfirmItem.vue";

import { useRequest } from "@/use/useAPI";
import { useAuthStore } from "@/stores/auth";
import type { Register } from "@/types/api/auth";

import Loading from "@/components/LoadingItem.vue";
import Swal from "sweetalert2";

const router = useRouter();
const auth = useAuthStore();

const nick = ref("");
const email = ref("");
const password = ref("");
const passwordConfirm = ref("");

const loading = ref(false);
const valid = computed(
    () =>
        useValidNickname(nick) &&
        useValidEmail(email) &&
        useValidPassword(password) &&
        passwordConfirm.value == password.value
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
