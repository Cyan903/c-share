<template>
    <div>
        <h1>Login</h1>
        <div v-show="loading" class="loading">Loading...</div>
        <form>
            <input type="text" v-model="email" placeholder="Email" />
            <input type="password" v-model="password" placeholder="Password" />
            <input
                :disabled="!invalid"
                type="submit"
                value="Submit"
                @click.prevent="login"
            />
        </form>
    </div>
</template>

<script lang="ts" setup>
import { useRequest } from "@/use/useAPI";
import { ref, computed } from "vue";
import Swal from "sweetalert2";

const email = ref("");
const password = ref("");
const loading = ref(false);

const invalid = computed(
    () =>
        email.value.length > 6 &&
        email.value.length < 30 &&
        password.value.length > 6 &&
        password.value.length < 30
);

const login = async () => {
    const req = await useRequest(
        "/auth/login",
        {
            email: email.value,
            password: password.value,
        },
        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not login!",
            text: String(req.json?.message),
            icon: "warning",
            confirmButtonText: "Okay",
        });

        email.value = "";
        password.value = "";
        return;
    }

    console.log("Succcess!", req.json);
};
</script>

<style scoped>
form input {
    display: block;
}
</style>
