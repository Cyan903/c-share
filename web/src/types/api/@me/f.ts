import type { Error } from "@/types";

export interface FileEditUpdate extends Error {}
export interface FileUploadRequest extends Error {}

export interface FileUpdate {
    id: string,
    comment: string;
    perm: string;
}

export interface FileUpload {
    file: File,
    comment: string;
    perm: string;
    password?: string;
}

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
