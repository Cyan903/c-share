<template>
    <tr>
        <ModalItem :show="modal" @hide="modal = false">
            <FileViewItem :id="data.id" :type="data.file_type" />

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

        <td class="small-image">
            <FileViewItem :id="data.id" :type="data.file_type" :small="true" />
        </td>

        <td>
            <a>{{ data.id }}</a>
        </td>

        <td>{{ data.file_type }}</td>
        <td :title="data.file_comment">
            <textarea readonly :value="shortened"></textarea>
        </td>
        <td>{{ permissions }}</td>
        <td :title="data.created_at">{{ date }}</td>

        <td>
            <input
                v-if="deleteMode"
                ref="deleteCheckbox"
                type="checkbox"
                :checked="selected"
                @click="updateSelect"
            />
        </td>

        <td>
            <button :disabled="deleteMode" @click="modal = true">Edit</button>
        </td>
    </tr>
</template>

<script lang="ts" setup>
import { computed, ref, reactive, watch, toRef } from "vue";
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

const deleteCheckbox = ref<HTMLInputElement>();
const modal = ref(false);
const loading = ref(false);
const selected = ref(false);

const edit = reactive({
    password: "",
    comment: "",
    perm: "",
});

const props = defineProps<{
    data: FileListingData;
    deleteMode: boolean;
}>();

const emit = defineEmits<{
    (e: "editFile", evt: FileUpdate): void;
    (e: "updatePurgeList", id: string, del: boolean): void;
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

const updateSelect = () => {
    if (!deleteCheckbox.value) return;

    selected.value = deleteCheckbox.value.checked;
    emit("updatePurgeList", props.data.id, !selected.value);
};

watch(modal, () => {
    edit.comment = props.data.file_comment;
    edit.perm = ["public", "private", "unlisted"][props.data.permissions];
});

watch(
    () => props.deleteMode,
    () => {
        selected.value = false;
    }
);
</script>

<style scoped>
img {
    width: 50px;
    height: 50px;
}

textarea {
    resize: none;
}
</style>
