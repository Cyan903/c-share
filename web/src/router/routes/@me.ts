import MeView from "@/views/@me/MeView.vue";
import ProfileView from "@/views/@me/ProfileView.vue";

export default [
    {
        path: "",
        name: "@me",
        component: MeView,
        meta: { requiresAuth: true },
    },

    {
        path: "profile",
        name: "profile",
        component: ProfileView,
        meta: { requiresAuth: true },
        children: [
            {
                path: "email",
                name: "prof-email",
                component: () =>
                    import("@/views/@me/profile/EmailProfileView.vue"),
            },

            {
                path: "nickname",
                name: "prof-nickname",
                component: () =>
                    import("@/views/@me/profile/NickProfileView.vue"),
            },

            {
                path: "password",
                name: "prof-password",
                component: () =>
                    import("@/views/@me/profile/PasswordProfileView.vue"),
            },

            {
                path: ":id",
                name: "prof-verify",
                component: () =>
                    import("@/views/@me/profile/VerificationProfileView.vue"),
            },
        ],
    },
];
