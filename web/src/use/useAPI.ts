import type { Ref } from "vue";
import Swal from "sweetalert2";


// prettier-ignore
export async function useRequest<Type>(url: string, options: RequestInit, elm: Ref<Boolean>, noAlert: boolean = false) {
    let error = false;
    let json!: Type;

    elm.value = true;

    const response = await fetch(import.meta.env.VITE_API + url, options).catch(
        (err) => (error = err)
    );

    if (error) {
        if (!noAlert) {
            Swal.fire({
                title: "API Error",
                text: "Could not make a request to the API!",
                icon: "error",
                confirmButtonText: "Okay",
            });
        }
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
