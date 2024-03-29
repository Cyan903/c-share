import { createRouter, createWebHistory } from "vue-router";
import NotFoundView from "@/views/NotFoundView.vue";

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

        {
            path: "/:pathMatch(.*)*",
            name: "not-found",
            component: NotFoundView,
            meta: { title: "Not Found" },
        },
    ],
});

router.beforeEach((to) => {
    const storage = localStorage.getItem("token");

    document.title = `${to.meta.title} | ${
        import.meta.env.VITE_APP || "c-share"
    }`;

    // Cannot be logged in
    if (to.meta.requiresNoAuth && storage) {
        return { path: "/@me" };
    }

    // Should be logged in
    if (to.meta.requiresAuth && !storage) {
        return { path: "/auth/login" };
    }
});

export default router;
