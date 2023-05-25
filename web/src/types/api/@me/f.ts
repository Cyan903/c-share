import type { Error } from "@/types";

export interface FileListing extends Error {
    code: number;
    count: number;
    data: FileListingData[];
}

export interface FileListingData {
    id: string;
    file_size: number;
    file_type: string;
    file_comment: string;
    permissions: number;
    created_at: string;
}
