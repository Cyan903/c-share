import HomeView from "@/views/HomeView.vue";

export default [
    {
        path: "/",
        name: "home",
        component: HomeView,
        meta: { title: "Home" },
    },

    {
        path: "/about",
        name: "about",
        component: () => import("@/views/AboutView.vue"),
        meta: { title: "About" },
    },
];
