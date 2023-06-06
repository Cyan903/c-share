<template>
    <div>
        <LoadingAuthItem :loading="loading" />

        <h2 class="font-semibold text-3xl my-4">Email Code Verification</h2>
        <p v-if="invalid">Could not verify email address!</p>
        <p v-else>
            <span class="text-success">Success!</span> Email has been verified!
        </p>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useRequest } from "@/use/useAPI";
import { useAuthStore } from "@/stores/auth";
import type { EmailCodeVerify } from "@/types/api/@me/profile";

import LoadingAuthItem from "@/components/loading/LoadingAuthItem.vue";
import Swal from "sweetalert2";

const [route, router] = [useRoute(), useRouter()];
const auth = useAuthStore();
const loading = ref(false);
const invalid = ref(false);

onMounted(async () => {
    const storage = localStorage.getItem("token");
    const req = await useRequest<EmailCodeVerify>(
        `/@me/profile/${route.params.id}`,
        {
            method: "POST",
            headers: { Token: String(storage) },
        },
        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not verify email!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        invalid.value = true;
        return;
    }

    Swal.fire({
        title: "Success",
        text: req.json.message,
        icon: "success",
        confirmButtonText: "Okay",
    });

    auth.userData.emailVerified = true;
    router.push("/@me/profile/email");
});
</script>
