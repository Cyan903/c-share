<template>
    <div>
        <Loading :loading="loading" />

        <h1>Welcome {{ auth.userData.nickname }}</h1>
        <p>Code: {{ route.params.id }}</p>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useRequest } from "@/use/useAPI";
import { useAuthStore } from "@/stores/auth";
import type { EmailCodeVerify } from "@/types/api/@me/profile";

import Loading from "@/components/LoadingItem.vue";
import Swal from "sweetalert2";

const [route, router] = [useRoute(), useRouter()];
const auth = useAuthStore();
const loading = ref(false);

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

        return;
    }

    Swal.fire({
        title: "Success",
        text: req.json.message,
        icon: "success",
        confirmButtonText: "Okay",
    });

    auth.userData.emailVerified = true;
    router.push("/@me");
});
</script>
