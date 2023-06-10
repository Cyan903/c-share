import MeView from "@/views/@me/MeView.vue";
import ProfileView from "@/views/@me/ProfileView.vue";

export default [
    {
        path: "",
        name: "@me",
        component: MeView,
        meta: { title: "Dashboard", requiresAuth: true },
    },

    {
        path: "profile",
        name: "profile",
        component: ProfileView,
        meta: { title: "Profile", requiresAuth: true },
        children: [
            {
                path: "",
                name: "prof-none",
                component: () =>
                    import("@/views/@me/profile/NoneProfileView.vue"),
                meta: { title: "Profile" },
            },

            {
                path: "email",
                name: "prof-email",
                component: () =>
                    import("@/views/@me/profile/EmailProfileView.vue"),
                meta: { title: "Profile Email" },
            },

            {
                path: "nickname",
                name: "prof-nickname",
                component: () =>
                    import("@/views/@me/profile/NickProfileView.vue"),
                meta: { title: "Profile Nickname" },
            },

            {
                path: "password",
                name: "prof-password",
                component: () =>
                    import("@/views/@me/profile/PasswordProfileView.vue"),
                meta: { title: "Profile Password" },
            },

            {
                path: "api",
                name: "prof-api",
                component: () =>
                    import("@/views/@me/profile/APIProfileView.vue"),
                meta: { title: "Profile API" },
            },

            {
                path: ":id",
                name: "prof-verify",
                component: () =>
                    import("@/views/@me/profile/VerificationProfileView.vue"),
                meta: { title: "Profile Verify" },
            },
        ],
    },
];
