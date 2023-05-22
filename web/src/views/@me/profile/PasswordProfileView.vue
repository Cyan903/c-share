<template>
    <div>
        <Loading :loading="loading" />
        <EmailStatusItem />

        <h2>Password Change</h2>
        <p v-if="!auth.userData.emailVerified">
            Your email must be verified in order to change your password. You
            can send a verification code
            <router-link to="/@me/profile/email">here</router-link>.
        </p>

        <form>
            <ValidPasswordItem
                v-model="pw.old"
                :disabled="!auth.userData.emailVerified"
                placehold="Old Password"
            />

            <ValidPasswordItem
                v-model="pw.new"
                :disabled="!auth.userData.emailVerified"
                placehold="New Password"
            />

            <ValidPasswordConfirmItem
                v-model="pw.confirm"
                :check="pw.new"
                :disabled="!auth.userData.emailVerified"
            />

            <input
                type="submit"
                value="Update"
                :disabled="!valid"
                @click.prevent="updatePassword"
            />
        </form>
    </div>
</template>

<script lang="ts" setup>
import { computed, reactive, ref, toRef } from "vue";

import { useValidPassword } from "@/use/useValidate";
import { useAuthStore } from "@/stores/auth";
import { useRequest } from "@/use/useAPI";
import type { PasswordUpdate } from "@/types/api/@me/profile";

import EmailStatusItem from "@/components/profile/EmailStatusItem.vue";
import ValidPasswordItem from "@/components/valid/ValidPasswordItem.vue";
import ValidPasswordConfirmItem from "@/components/valid/ValidPasswordConfirmItem.vue";

import Loading from "@/components/LoadingItem.vue";
import Swal from "sweetalert2";

const loading = ref(false);
const auth = useAuthStore();
const pw = reactive({
    old: "",
    new: "",
    confirm: "",
});

const valid = computed(
    () =>
        useValidPassword(toRef(pw.new)) &&
        useValidPassword(toRef(pw.old)) &&
        pw.confirm == pw.new &&
        auth.userData.emailVerified
);

const updatePassword = async () => {
    const req = await useRequest<PasswordUpdate>(
        "/@me/profile/password",
        {
            method: "POST",
            headers: { "Content-Type": "application/json", Token: auth.token },
            body: JSON.stringify({
                old_password: pw.old,
                new_password: pw.new,
            }),
        },
        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not update password!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        return;
    }

    Swal.fire({
        title: "Success",
        text: req.json.message,
        icon: "success",
        confirmButtonText: "Okay",
    });

    pw.old = "";
    pw.new = "";
};
</script>
