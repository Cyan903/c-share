<template>
    <h1>Hello, {{ auth.userData.nickname }}</h1>
</template>

<script lang="ts" setup>
import { useAuthStore } from "@/stores/auth";

const auth = useAuthStore();
</script>
