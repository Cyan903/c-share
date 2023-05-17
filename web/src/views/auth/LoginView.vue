<template>
    <div>
        <h1>Login</h1>
        <div v-show="loading" class="loading">Loading...</div>
        <form>
            <input type="text" v-model="email" placeholder="Email" />
            <input type="password" v-model="password" placeholder="Password" />
            <input
                :disabled="!invalid"
                type="submit"
                value="Submit"
                @click.prevent="login"
            />
        </form>

        <pre>{{ auth.userData }}</pre>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue";
import { type Login } from "@/types/api/auth";

import { useRequest } from "@/use/useAPI";
import { useAuthStore } from "@/stores/auth";
import { useRouter } from "vue-router";

import Swal from "sweetalert2";

const router = useRouter();
const auth = useAuthStore();

const email = ref("");
const password = ref("");
const loading = ref(false);

const invalid = computed(
    () =>
        email.value.length > 6 &&
        email.value.length < 30 &&
        password.value.length > 6 &&
        password.value.length < 30
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

<style scoped>
form input {
    display: block;
}
</style>
