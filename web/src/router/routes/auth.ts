import LoginView from "@/views/auth/LoginView.vue";
import RegisterView from "@/views/auth/RegisterView.vue";

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
];
