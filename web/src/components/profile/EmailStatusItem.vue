<template>
    <div class="my-3">
        Email Status:
        <span :class="`text-${verified[1]}`">{{ verified[0] }}</span>
    </div>
</template>

<script lang="ts" setup>
import { computed } from "vue";
import { useAuthStore } from "@/stores/auth";

const auth = useAuthStore();
const verified = computed(() =>
    auth.userData.emailVerified
        ? ["Verified", "success"]
        : ["Unverified", "error"]
);
</script>
