<template>
    <div>
        <div>
            <input type="file" @change="updateFile" />
        </div>

        <div>
            <span>Permission</span> |
            <select v-model="upload.perm">
                <option value="public">Public</option>
                <option value="private">Private</option>
                <option value="unlisted">Unlisted</option>
            </select>

            <ValidCommentItem v-model="upload.comment" />
        </div>

        <div v-if="hasPassword">
            <ValidPasswordItem v-model="upload.password" />
            <ValidPasswordConfirmItem
                v-model="upload.confirm"
                :check="upload.password"
            />
        </div>

        <button :disabled="!valid" @click="uploadFile">Upload</button>
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
