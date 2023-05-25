<template>
    <tr>
        <ModalItem :show="modal" @hide="modal = false">
            <ValidPasswordItem v-model="edit.password" placehold="File Password" />
            <input type="text" v-model="edit.comment" placeholder="File Comment" />

            <span>Listing</span> |
            <select v-model="edit.perm">
                <option value="public">Public</option>
                <option value="private">Private</option>
                <option value="unlisted">Unlisted</option>
            </select>
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
import { computed, ref, onMounted, reactive } from "vue";
import type { FileListingData } from "@/types/api/@me/f";

import ValidPasswordItem from "@/components/valid/ValidPasswordItem.vue";
import ModalItem from "@/components/@me/ModalItem.vue";

import moment from "moment";
import Swal from "sweetalert2";

// TODO: Expand image on hover
// TODO: Edit file data

const coverImage = ref("");
const modal = ref(false);
const edit = reactive({
    password: "",
    comment: "",
    perm: "",
});

const props = defineProps<{
    data: FileListingData;
}>();

const shortened = computed(() =>
    props.data.file_comment.length > 15
        ? props.data.file_comment.slice(0, 15) + "..."
        : props.data.file_comment
);

const date = computed(() => moment(props.data.created_at).fromNow());
const permissions = computed(
    () => ["Public", "Private", "Unlisted"][props.data.permissions]
);

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
</script>

<style scoped>
img {
    width: 50px;
    height: 50px;
}
</style>
