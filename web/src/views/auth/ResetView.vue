<template>
    <div class="auth-form flex flex-wrap items-center justify-center">
        <LoadingAuthItem :loading="loading" />

        <div class="card card-normal bg-base-300 shadow-xl">
            <div class="card-body items-center text-center">
                <h2 class="card-title text-3xl">Password Reset</h2>
                <p class="mb-4">
                    Please note that your email must be verified in order to
                    request a password reset.
                </p>

                <ValidEmailItem v-model="email" />

                <input
                    class="btn btn-primary w-5/6 my-4"
                    :disabled="!valid"
                    type="submit"
                    value="Send Reset"
                    @click.prevent="sendVerification"
                />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref } from "vue";
import { useRequest } from "@/use/useAPI";
import { useValidEmail } from "@/use/useValidate";
import type { PasswordReset } from "@/types/api/auth";

import ValidEmailItem from "@/components/valid/ValidEmailItem.vue";
import LoadingAuthItem from "@/components/loading/LoadingAuthItem.vue";

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
