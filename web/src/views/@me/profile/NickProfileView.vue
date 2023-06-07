<template>
    <div>
        <LoadingAuthItem :loading="loading" />

        <h2 class="font-semibold text-3xl my-4">Nickname Settings</h2>
        <p class="mx-4 my-3">
            Must be between 4 and 10 characters. Cannot use numbers or special
            characters.
        </p>

        <form>
            <ValidNicknameItem v-model="nick" />

            <input
                :disabled="!valid"
                type="submit"
                value="Update Nickname"
                class="btn btn-primary btn-outline mt-4"
                @click.prevent="updateNickname"
            />
        </form>
    </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from "vue";

import { useAuthStore } from "@/stores/auth";
import { useRequest } from "@/use/useAPI";
import { useValidNickname } from "@/use/useValidate";
import type { NicknameUpdate } from "@/types/api/@me/profile";

import ValidNicknameItem from "@/components/valid/ValidNicknameItem.vue";
import LoadingAuthItem from "@/components/loading/LoadingAuthItem.vue";

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

onMounted(() => {
    nick.value = auth.userData.nickname;
});
</script>
