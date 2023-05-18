import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";

import Home from "@/router/routes/home";
import Auth from "@/router/routes/auth";
import Me from "@/router/routes/@me";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        ...Home,

        {
            path: "/auth",
            children: [...Auth],
        },

        {
            path: "/@me",
            children: [...Me],
        },
    ],
});

router.beforeEach((to) => {
    const auth = useAuthStore();

    // Cannot be logged in
    if (to.meta.requiresNoAuth && auth.isLoggedIn) {
        return { path: "/@me" };
    }

    // Should be logged in
    if (to.meta.requiresAuth && !auth.isLoggedIn) {
        return { path: "/auth/login" };
    }
});

export default router;
