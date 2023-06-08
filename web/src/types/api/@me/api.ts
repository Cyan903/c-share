import type { Error } from "@/types";

export interface UserTokenData extends Error {
    code: number;
    count: number;
    data: TokenData;
}

export interface TokenData {
    token: string;
    user: string;
    created_at: string;
}
