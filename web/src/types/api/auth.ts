import { type Error } from "@/types/index";

export interface Login extends Error {
    code: number;
    count: number;
    data: string;
}

export interface Register extends Error {
    code: number;
    count: number;
    data: string;
}
