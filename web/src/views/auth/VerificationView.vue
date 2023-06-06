<template>
    <div class="auth-form flex flex-wrap items-center justify-center">
        <LoadingAuthItem :loading="loading" />

        <div class="card card-normal bg-base-300 shadow-xl">
            <div class="card-body items-center text-center">
                <h2 class="card-title text-3xl">Password Reset</h2>

                <ValidPasswordItem v-model="password" />
                <ValidPasswordConfirmItem
                    v-model="passwordConfirm"
                    :check="password"
                />

                <input
                    class="btn btn-primary w-5/6 my-4"
                    :disabled="!valid"
                    type="submit"
                    value="Reset Password"
                    @click.prevent="updatePassword"
                />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref } from "vue";
import { useRoute, useRouter } from "vue-router";

import { useRequest } from "@/use/useAPI";
import { useValidPassword } from "@/use/useValidate";
import type { ResetToken } from "@/types/api/auth";

import ValidPasswordItem from "@/components/valid/ValidPasswordItem.vue";
import ValidPasswordConfirmItem from "@/components/valid/ValidPasswordConfirmItem.vue";
import LoadingAuthItem from "@/components/loading/LoadingAuthItem.vue";

import Swal from "sweetalert2";

const password = ref("");
const passwordConfirm = ref("");
const loading = ref(false);

const [route, router] = [useRoute(), useRouter()];
const valid = computed(
    () => useValidPassword(password) && password.value == passwordConfirm.value
);

const updatePassword = async () => {
    const req = await useRequest<ResetToken>(
        "/auth/" + route.params.id,
        {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ password: password.value }),
        },
        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not reset password!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        return;
    }

    Swal.fire({
        title: "Password has been reset!",
        icon: "success",
        confirmButtonText: "Okay",
    }).then(() => router.push("/auth/login"));
};
</script>
