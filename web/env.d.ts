/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_APP: string;
    readonly VITE_API: string;
    readonly VITE_SOURCE: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
