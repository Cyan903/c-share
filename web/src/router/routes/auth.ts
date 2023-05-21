import LoginView from "@/views/auth/LoginView.vue";
import RegisterView from "@/views/auth/RegisterView.vue";
import ResetView from "@/views/auth/ResetView.vue";
import VerificationView from "@/views/auth/VerificationView.vue";

export default [
    {
        path: "login",
        name: "login",
        component: LoginView,
        meta: { title: "Login", requiresNoAuth: true },
    },

    {
        path: "register",
        name: "register",
        component: RegisterView,
        meta: { title: "Register", requiresNoAuth: true },
    },

    {
        path: "pwreset",
        name: "pwreset",
        component: ResetView,
        meta: { title: "Password Reset", requiresNoAuth: true },
    },

    {
        path: ":id",
        name: "verification",
        component: VerificationView,
        meta: { title: "Password Verification", requiresNoAuth: true },
    },
];
