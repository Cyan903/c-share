import { createRouter, createWebHistory } from "vue-router";

import Home from "@/router/routes/home";
import Auth from "@/router/routes/auth";
import Me from "@/router/routes/@me";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [...Home, ...Auth, ...Me],
});

export default router;
