import type { Ref } from "vue";

export function useValidEmail(email: Ref<string>) {
    return email.value.length > 6 && email.value.length < 30;
}

export function useValidPassword(password: Ref<string>) {
    return password.value.length > 6 && password.value.length < 30;
}
