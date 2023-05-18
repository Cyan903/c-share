// /pkg/api/validate.go

import type { Ref } from "vue";

export function useValidNickname(nick: Ref<string>) {
    return (
        nick.value.search(/[^a-zA-Z]+/) === -1 &&
        nick.value.length < 10 &&
        nick.value.length > 3
    );
}

export function useValidEmail(email: Ref<string>) {
    return email.value.length > 6 && email.value.length < 30;
}

export function useValidPassword(password: Ref<string>) {
    return password.value.length > 6 && password.value.length < 30;
}
