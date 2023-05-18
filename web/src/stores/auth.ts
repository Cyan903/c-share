import type { AtMe } from "@/types/api/@me";
import { useRequest } from "@/use/useAPI";
import { defineStore } from "pinia";
import { computed, onMounted, reactive, ref, watch } from "vue";

export const useAuthStore = defineStore("auth", () => {
    const token = ref("");
    const isLoggedIn = computed(() => token.value != "");
    const userData = reactive({
        nickname: "",
        email: "",
        emailVerified: false,
        usedStorage: 0,
        createdAt: "",
    });

    let timer = 0;

    const login = async (jwt: string) => {
        const tkexpire = JSON.parse(atob(jwt.split(".")[1]));
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

        timer = setInterval(() => {
            if (Math.floor(Date.now() / 1000) > tkexpire.exp) {
                clearInterval(timer);
                logout();

                location.href = "/auth/login";
            }
        }, 60000);

        return true;
    };

    const logout = () => {
        userData.nickname = "";
        userData.email = "";
        userData.emailVerified = false;
        userData.usedStorage = 0;
        userData.createdAt = "";

        token.value = "";
        localStorage.clear();
        clearInterval(timer);
    };

    watch(token, () => localStorage.setItem("token", token.value));
    onMounted(() => {
        const storage = localStorage.getItem("token");

        if (storage) {
            console.info("[auth] token found, loading from storage");
            login(storage);
        }
    });

    return {
        token,
        userData,
        isLoggedIn,
        login,
        logout,
    };
});
