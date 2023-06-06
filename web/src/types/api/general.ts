import type { Error } from "@/types/index";

export interface ServerInfo extends Error {
    code: number;
    count: number;
    data: Data;
}

export interface Data {
    users: number;
    storage: string;
    total_files: number;
}
