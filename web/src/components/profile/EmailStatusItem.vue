<template>
    <span>Email status: {{ auth.userData.emailVerified }}</span>
</template>

<script lang="ts" setup>
import { useAuthStore } from "@/stores/auth";

const auth = useAuthStore();
</script>
