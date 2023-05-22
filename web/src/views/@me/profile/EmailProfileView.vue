<template>
    <div>
        <Loading :loading="loading" />
        <EmailStatusItem />

        <h1>Email</h1>
        <form>
            <ValidEmailItem
                v-model="email"
                :disabled="!auth.userData.emailVerified"
            />

            <input
                :disabled="!valid"
                type="submit"
                value="Update Email"
                @click.prevent="updateEmail"
            />
        </form>

        <hr />

        <div>
            <h3>Email Verification</h3>
            <p>
                You cannot change your email address or reset your password if
                your email is not verified. It is highly recommended that you
                verify your email.
            </p>

            <p v-if="auth.userData.emailVerified">
                Hooray! Your email is already verified! Please note that
                changing it will reset your verification.
            </p>

            <button
                @click="sendVerification"
                :disabled="auth.userData.emailVerified"
            >
                Send Verification Email
            </button>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref, onMounted } from "vue";

import { useValidEmail } from "@/use/useValidate";
import { useAuthStore } from "@/stores/auth";
import { useRequest } from "@/use/useAPI";
import type { EmailUpdate, EmailSendVerify } from "@/types/api/@me/profile";

import EmailStatusItem from "@/components/profile/EmailStatusItem.vue";
import Loading from "@/components/LoadingItem.vue";
import Swal from "sweetalert2";
import ValidEmailItem from "@/components/valid/ValidEmailItem.vue";

const auth = useAuthStore();
const loading = ref(false);
const email = ref("");

const valid = computed(
    () => useValidEmail(email) && auth.userData.emailVerified
);

const updateEmail = async () => {
    const req = await useRequest<EmailUpdate>(
        "/@me/profile/email",
        {
            method: "POST",
            headers: { "Content-Type": "application/json", Token: auth.token },
            body: JSON.stringify({ email: email.value }),
        },
        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not update email!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        email.value = "";
        return;
    }

    auth.userData.email = email.value;
    auth.userData.emailVerified = false;

    Swal.fire({
        title: "Success",
        text: req.json.message,
        icon: "success",
        confirmButtonText: "Okay",
    });
};

const sendVerification = async () => {
    const req = await useRequest<EmailSendVerify>(
        "/@me/profile/verify",
        {
            method: "POST",
            headers: { Token: auth.token },
        },
        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not send email!",
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
};

onMounted(() => {
    email.value = auth.userData.email;
});
</script>
