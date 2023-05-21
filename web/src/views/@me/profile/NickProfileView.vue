<template>
    <div>
        <Loading :loading="loading" />

        <h1>Nickname</h1>
        <h4>Hello {{ auth.userData.nickname }}</h4>

        <form>
            <input type="text" v-model="nick" placeholder="Nickname" v-focus />
            <input
                type="submit"
                value="Update Nickname"
                :disabled="!valid"
                @click.prevent="updateNickname"
            />
        </form>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref } from "vue";

import { vFocus } from "@/directives/vFocus";
import { useAuthStore } from "@/stores/auth";
import { useRequest } from "@/use/useAPI";
import { useValidNickname } from "@/use/useValidate";
import type { NicknameUpdate } from "@/types/api/@me/profile";

import Swal from "sweetalert2";

const auth = useAuthStore();
const loading = ref(false);
const nick = ref("");

const valid = computed(() => useValidNickname(nick));

const updateNickname = async () => {
    const req = await useRequest<NicknameUpdate>(
        "/@me/profile/nickname",
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Token: auth.token,
            },
            body: JSON.stringify({ nickname: nick.value }),
        },
        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not update nickname!",
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

    auth.userData.nickname = nick.value;
    nick.value = "";
};
</script>
