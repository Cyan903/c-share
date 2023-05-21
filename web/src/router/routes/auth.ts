import LoginView from "@/views/auth/LoginView.vue";
import RegisterView from "@/views/auth/RegisterView.vue";
import ResetView from "@/views/auth/ResetView.vue";
import VerificationView from "@/views/auth/VerificationView.vue";

export default [
    {
        path: "login",
        name: "login",
        component: LoginView,
        meta: { requiresNoAuth: true },
    },

    {
        path: "register",
        name: "register",
        component: RegisterView,
        meta: { requiresNoAuth: true },
    },

    {
        path: "pwreset",
        name: "pwreset",
        component: ResetView,
        meta: { requiresNoAuth: true },
    },

    {
        path: ":id",
        name: "verification",
        component: VerificationView,
        meta: { requiresNoAuth: true },
    },
];
