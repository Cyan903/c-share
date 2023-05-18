import LoginView from "@/views/auth/LoginView.vue";
import RegisterView from "@/views/auth/RegisterView.vue";
import ResetView from "@/views/auth/ResetView.vue";
import VerificationView from "@/views/auth/VerificationView.vue";

export default [
    {
        path: "/auth/login",
        name: "login",
        component: LoginView,
    },
    {
        path: "/auth/register",
        name: "register",
        component: RegisterView,
    },
    {
        path: "/auth/pwreset",
        name: "pwreset",
        component: ResetView,
    },
    {
        path: "/auth/:id",
        name: "verification",
        component: VerificationView,
    },
];
