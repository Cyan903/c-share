<template>
    <div class="upload-form">
        <h3 class="font-semibold text-2xl mb-6 text-center">File Upload</h3>
        <div class="text-center">
            <input
                type="file"
                class="file-input file-input-bordered w-full max-w-xs"
                @change="updateFile"
            />

            <label class="input-group">
                <button class="btn btn-active no-animation">Permission</button>
                <select class="select select-bordered" v-model="upload.perm">
                    <option value="public">Public</option>
                    <option value="private">Private</option>
                    <option value="unlisted">Unlisted</option>
                </select>
            </label>
        </div>

        <div class="divider"></div>
        <div class="grid content-between grid-cols-1">
            <div v-if="hasPassword">
                <ValidPasswordItem v-model="upload.password" />
                <ValidPasswordConfirmItem
                    v-model="upload.confirm"
                    :check="upload.password"
                />
            </div>

            <ValidCommentItem v-model="upload.comment" />
        </div>

        <button
            class="block my-3 m-auto btn btn-primary"
            :disabled="!valid"
            @click="uploadFile"
        >
            Upload
        </button>
    </div>
</template>

<script lang="ts" setup>
import { computed, reactive, toRef } from "vue";
import { useValidPassword, useValidComment } from "@/use/useValidate";
import type { FileUpload } from "@/types/api/@me/f";

import ValidPasswordItem from "@/components/valid/ValidPasswordItem.vue";
import ValidCommentItem from "@/components/valid/ValidCommentItem.vue";
import ValidPasswordConfirmItem from "@/components/valid/ValidPasswordConfirmItem.vue";

import Swal from "sweetalert2";

let inputFile: File;

const emits = defineEmits<{
    (e: "fileUpload", evt: FileUpload): void;
}>();

const upload = reactive({
    perm: "public",
    password: "",
    confirm: "",
    comment: "",
});

const hasPassword = computed(() => upload.perm == "unlisted");
const valid = computed(() => {
    if (upload.perm == "unlisted") {
        return (
            useValidPassword(toRef(upload.password)) &&
            useValidComment(toRef(upload.comment)) &&
            upload.password == upload.confirm
        );
    }

    return useValidComment(toRef(upload.comment));
});

const updateFile = (file: Event) => {
    const files = (file.target as HTMLInputElement).files;

    if (!files || !files[0]) {
        Swal.fire({
            title: "Upload Error!",
            text: "No files provided!",
            icon: "warning",
            confirmButtonText: "Okay",
        });

        return;
    }

    inputFile = files[0];
};

const uploadFile = () => {
    if (!inputFile) {
        Swal.fire({
            title: "Upload Error!",
            text: "No files provided!",
            icon: "warning",
            confirmButtonText: "Okay",
        });

        return;
    }

    if (upload.perm == "unlisted") {
        emits("fileUpload", {
            file: inputFile,
            perm: upload.perm,
            password: upload.password,
            comment: upload.comment,
        });

        return;
    }

    emits("fileUpload", {
        file: inputFile,
        perm: upload.perm,
        comment: upload.comment,
    });
};
</script>
