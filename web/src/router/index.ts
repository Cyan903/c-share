import { createRouter, createWebHistory } from "vue-router";

import home from "@/router/routes/home";
import auth from "@/router/routes/auth";

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [...home, ...auth],
});

export default router;
