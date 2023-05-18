<template>
    <div>
        <Loading :loading="loading" />

        <h1>Password Reset</h1>
        <p>
            Please note that your email must be verified in order to request a
            password reset.
        </p>

        <form>
            <input type="text" v-model="email" placeholder="Email" />

            <input
                :disabled="!invalid"
                type="submit"
                value="Submit"
                @click.prevent="sendVerification"
            />
        </form>
    </div>
</template>

<script lang="ts" setup>
import { useRequest } from "@/use/useAPI";
import { computed, ref } from "vue";
import { useValidEmail } from "@/use/useValidate";
import type { PasswordReset } from "@/types/api/auth";

import Loading from "@/components/LoadingItem.vue";
import Swal from "sweetalert2";

// TODO: v-autofocus

const email = ref("");
const loading = ref(false);
const invalid = computed(() => useValidEmail(email));

const sendVerification = async () => {
    const req = await useRequest<PasswordReset>(
        "/auth/pwreset",
        {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                email: email.value,
            }),
        },
        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not send password reset!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        return;
    }

    Swal.fire({
        title: "Check your inbox.",
        text: req.json.message,
        icon: "success",
        confirmButtonText: "Okay",
    });
};
</script>
