import type { Error, SimpleResponse } from "@/types/index";

export interface Login extends SimpleResponse {}
export interface Register extends SimpleResponse {}
export interface PasswordReset extends SimpleResponse {}

export interface ResetToken extends Error {
    code: number;
    message: string;
}
