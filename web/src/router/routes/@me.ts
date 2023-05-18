import MeView from "@/views/@me/MeView.vue";

export default [
    {
        path: "",
        name: "@me",
        component: MeView,
        meta: { requiresAuth: true },
    },
];
