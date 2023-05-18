export interface Error {
    code: number;
    message: string;
}

export interface SimpleResponse extends Error {
    code: number;
    count: number;
    data: string;
}
