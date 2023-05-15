import type { Ref } from "vue";
import Swal from "sweetalert2";

type Response = {
    code: Number;
    message?: String;
    count?: Number;
    data?: String;
};

export async function useRequest(url: String, body: Object, elm: Ref<Boolean>) {
    let error = false;
    let json: Response = {
        code: 0,
        message: "",
    };

    elm.value = true;

    const response = await fetch(import.meta.env.VITE_API + url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },

        body: JSON.stringify(body),
    }).catch((err) => (error = err));

    if (error) {
        Swal.fire({
            title: "API Error",
            text: "Could not make a request to the API!",
            icon: "error",
            confirmButtonText: "Okay",
        });
    } else {
        if (response.status != 200) {
            console.warn(`[api] Error with request ${url}!`);
            console.error(response);
        }

        json = await response.json();
    }

    elm.value = false;

    return {
        error,
        response,
        json,
    };
}
