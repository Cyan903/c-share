<template>
    <tr>
        <ModalItem :show="modal" @hide="modal = false">
            <FileViewItem />

            <Loading :loading="loading" />
            <ValidPasswordItem
                v-model="edit.password"
                :disabled="edit.perm != 'unlisted'"
                placehold="File Password"
            />

            <ValidCommentItem v-model="edit.comment" />

            <span>Listing</span> |
            <select v-model="edit.perm">
                <option value="public">Public</option>
                <option value="private">Private</option>
                <option value="unlisted">Unlisted</option>
            </select>

            <div>
                <button :disabled="loading || !valid" @click="updateFile">
                    Save Changes
                </button>
                <button :disabled="loading || !valid" @click="modal = false">
                    Cancel
                </button>
            </div>
        </ModalItem>

        <td>
            <img :src="coverImage" />
        </td>

        <td>
            <a>{{ data.id }}</a>
        </td>

        <td>{{ data.file_type }}</td>
        <td>{{ shortened }}</td>
        <td>{{ permissions }}</td>
        <td :title="data.created_at">{{ date }}</td>

        <td>
            <button>Remove</button>
        </td>

        <td>
            <button @click="modal = true">Edit</button>
        </td>
    </tr>
</template>

<script lang="ts" setup>
import { computed, ref, onMounted, reactive, watch, toRef } from "vue";
import type {
    FileEditUpdate,
    FileUpdate,
    FileListingData,
} from "@/types/api/@me/f";

import { useValidComment, useValidPassword } from "@/use/useValidate";
import { useRequest } from "@/use/useAPI";

import ValidPasswordItem from "@/components/valid/ValidPasswordItem.vue";
import ValidCommentItem from "@/components/valid/ValidCommentItem.vue";
import ModalItem from "@/components/@me/util/ModalItem.vue";
import FileViewItem from "@/components/@me/util/FileViewItem.vue";
import Loading from "@/components/LoadingItem.vue";

import moment from "moment";
import Swal from "sweetalert2";

// TODO: Expand image on hover

const coverImage = ref("");
const modal = ref(false);
const loading = ref(false);
const edit = reactive({
    password: "",
    comment: "",
    perm: "",
});

const props = defineProps<{
    data: FileListingData;
}>();

const emit = defineEmits<{
    (e: "editFile", evt: FileUpdate): void;
}>();

const valid = computed(() => {
    if (edit.perm == "unlisted") {
        return (
            useValidPassword(toRef(edit.password)) &&
            useValidComment(toRef(edit.comment))
        );
    }

    return useValidComment(toRef(edit.comment));
});

const shortened = computed(() =>
    props.data.file_comment.length > 15
        ? props.data.file_comment.slice(0, 15) + "..."
        : props.data.file_comment
);

const date = computed(() => moment(props.data.created_at).fromNow());
const permissions = computed(
    () => ["Public", "Private", "Unlisted"][props.data.permissions]
);

const updateFile = async () => {
    const token = localStorage.getItem("token");
    const params = new URLSearchParams({
        password: edit.password,
        comment: edit.comment,
        perm: edit.perm,
    }).toString();

    if (!token) return;

    const req = await useRequest<FileEditUpdate>(
        `/@me/f/${props.data.id}/edit?${params}`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Token: token,
            },
        },
        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not update file!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        return;
    }

    Swal.fire({
        title: "Success",
        text: req.json.message,
        icon: "success",
        confirmButtonText: "Okay",
    });

    emit("editFile", {
        id: props.data.id,
        comment: edit.comment,
        perm: edit.perm,
    });
};

onMounted(async () => {
    const token = localStorage.getItem("token");
    const files = [
        "image/bmp",
        "image/jpeg",
        "image/x-png",
        "image/png",
        "image/gif",
    ];

    if (!files.includes(props.data.file_type)) {
        // TODO: Default icon
        coverImage.value = "favicon.ico";
        return;
    }

    if (!token) return;

    const req = await fetch(
        `${import.meta.env.VITE_API}/@me/f/${props.data.id}`,
        {
            headers: {
                Token: token,
            },
        }
    )
        .then((r) => r.blob())
        .catch(() => {
            Swal.fire({
                title: `Could not load image ${props.data.id}!`,
                icon: "error",
                confirmButtonText: "Okay",
            });
        });

    if (!req) return;
    coverImage.value = URL.createObjectURL(req);
});

watch(modal, () => {
    edit.comment = props.data.file_comment;
    edit.perm = ["public", "private", "unlisted"][props.data.permissions];
});
</script>

<style scoped>
img {
    width: 50px;
    height: 50px;
}
</style>
