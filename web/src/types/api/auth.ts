import type { Error } from "@/types/index";

// TODO: Reduce duplication

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

export interface PasswordReset extends Error {
    code: number;
    count: number;
    data: string;
}

export interface ResetToken extends Error {
    code: number;
    message: string;
}
