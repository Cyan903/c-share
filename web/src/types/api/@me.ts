import type { Error } from "@/types/index";

export interface FileUploadRequest extends Error {}
export interface FileDeleteRequest extends Error {}

export interface AtMe extends Error {
    code: number;
    count: number;
    data: AtMeData;
}

export interface AtMeData extends Error {
    nickname: string;
    email: string;
    email_verified: number;
    used_storage: number;
    created_at: string;
}
