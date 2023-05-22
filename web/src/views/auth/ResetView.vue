<template>
    <div>
        <Loading :loading="loading" />

        <h1>Password Reset</h1>
        <p>
            Please note that your email must be verified in order to request a
            password reset.
        </p>

        <form>
            <ValidEmailItem v-model="email" />

            <input
                :disabled="!valid"
                type="submit"
                value="Submit"
                @click.prevent="sendVerification"
            />
        </form>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref } from "vue";
import { useRequest } from "@/use/useAPI";
import { useValidEmail } from "@/use/useValidate";
import type { PasswordReset } from "@/types/api/auth";

import ValidEmailItem from "@/components/valid/ValidEmailItem.vue";
import Loading from "@/components/LoadingItem.vue";
import Swal from "sweetalert2";

const email = ref("");
const loading = ref(false);
const valid = computed(() => useValidEmail(email));

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
