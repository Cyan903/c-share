<template>
    <div>
        <Loading :loading="loading" />

        <h1>New Password</h1>
        <form>
            <ValidPasswordItem v-model="password" />
            <ValidPasswordConfirmItem
                v-model="passwordConfirm"
                :check="password"
            />

            <input
                :disabled="!valid"
                type="submit"
                value="Submit"
                @click.prevent="updatePassword"
            />
        </form>
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

import Loading from "@/components/LoadingItem.vue";
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
