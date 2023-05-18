import type { AtMe } from "@/types/api/@me";
import { useRequest } from "@/use/useAPI";
import { defineStore } from "pinia";
import { computed, reactive, ref } from "vue";

// TODO: Expire token
// TODO: localstorage
// TODO: Guards

export const useAuthStore = defineStore("auth", () => {
    const token = ref("");
    const userData = reactive({
        nickname: "",
        email: "",
        emailVerified: false,
        usedStorage: 0,
        createdAt: "",
    });

    const isLoggedIn = computed(() => token.value != "");

    const login = async (jwt: string) => {
        const user = await useRequest<AtMe>(
            "/@me",
            {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    Token: jwt,
                },
            },
            ref(false)
        );

        if (user.error) return false;
        if (user.response.status != 200) {
            console.warn("[auth] Could not login!");
            console.error(user);

            return false;
        }

        userData.nickname = user.json.data.nickname;
        userData.email = user.json.data.email;
        userData.emailVerified = !!user.json.data.email_verified;
        userData.usedStorage = user.json.data.used_storage;
        userData.createdAt = user.json.data.created_at;
        token.value = jwt;

        return true;
    };

    const logout = () => {
        userData.nickname = "";
        userData.email = "";
        userData.emailVerified = false;
        userData.usedStorage = 0;
        userData.createdAt = "";
        token.value = "";
    };

    return {
        token,
        userData,
        isLoggedIn,
        login,
        logout,
    };
});
