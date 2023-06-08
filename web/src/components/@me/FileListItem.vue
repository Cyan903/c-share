<template>
    <tr class="hidden lg:table-row">
        <ModalItem :show="modal" @hide="modal = false">
            <h3 class="font-semibold text-2xl mb-6 text-center">
                File Details
            </h3>

            <FileViewItem :id="data.id" :type="data.file_type" />
            <div class="lg:hidden text-center">
                <h5
                    class="font-semibold text-xl hover:opacity-70"
                    @click="copyID"
                >
                    ID: {{ data.id }}
                </h5>
                <div class="text-sm my-1 opacity-50">
                    <b>{{ data.file_type }}</b>

                    <span
                        class="text-sm opacity-50 my-1"
                        :title="data.created_at"
                    >
                        | {{ date }}
                    </span>
                </div>

                <div
                    class="badge badge-ghost badge-sm"
                    :title="`${data.file_size} bytes`"
                >
                    {{ fileSize }}
                </div>
            </div>

            <div class="divider"></div>
            <LoadingFileViewItemVue :loading="loading" />

            <div class="file-list-form text-center">
                <ValidPasswordItem
                    v-model="edit.password"
                    :disabled="edit.perm != 'unlisted'"
                    placehold="File Password"
                />

                <ValidCommentItem v-model="edit.comment" />

                <label class="input-group">
                    <button class="btn btn-active no-animation">
                        Permission
                    </button>
                    <select class="select select-bordered" v-model="edit.perm">
                        <option value="public">Public</option>
                        <option value="private">Private</option>
                        <option value="unlisted">Unlisted</option>
                    </select>
                </label>
            </div>

            <div class="flex justify-center my-4">
                <button
                    class="btn btn-success m-1"
                    :disabled="loading || !valid"
                    @click="updateFile"
                >
                    Save
                </button>

                <button
                    class="btn btn-primary m-1"
                    :disabled="loading"
                    @click="copyID"
                >
                    Copy
                </button>

                <button
                    class="btn btn-error m-1"
                    :disabled="loading || !valid"
                    @click="modal = false"
                >
                    Cancel
                </button>
            </div>
        </ModalItem>

        <td v-if="deleteMode">
            <input
                ref="deleteCheckbox"
                type="checkbox"
                class="checkbox"
                :checked="selected"
                @input="updateSelect"
            />
        </td>

        <td>
            <div class="flex items-center space-x-3">
                <div class="avatar">
                    <div class="mask mask-squircle w-12 h-12 small-image">
                        <FileViewItem
                            :id="data.id"
                            :type="data.file_type"
                            :small="true"
                        />
                    </div>
                </div>
                <div>
                    <div class="font-bold hover:opacity-70" @click="copyID">
                        {{ data.id }}
                    </div>
                    <div class="text-sm opacity-50" :title="data.created_at">
                        {{ date }}
                    </div>
                </div>
            </div>
        </td>

        <td>
            {{ data.file_type }}
            <br />
            <span
                class="badge badge-ghost badge-sm"
                :title="`${data.file_size} bytes`"
            >
                {{ fileSize }}
            </span>
        </td>

        <td>{{ permissions }}</td>

        <td :title="data.file_comment">
            <textarea
                class="textarea textarea-bordered"
                readonly
                :value="shortened"
            ></textarea>
        </td>

        <th>
            <button
                class="btn btn-ghost btn-xs"
                :disabled="deleteMode"
                @click="modal = true"
            >
                Details
            </button>
        </th>
    </tr>
    <tr class="lg:hidden block file-item-mobile">
        <div class="avatar mask mask-squircle w-20">
            <FileViewItem :id="data.id" :type="data.file_type" :small="true" />
        </div>

        <div class="avatar-overlay">
            <div v-if="deleteMode">
                <input
                    ref="deleteCheckboxMobile"
                    type="checkbox"
                    class="checkbox"
                    :checked="selected"
                    @click="updateSelectMobile"
                />
            </div>

            <button
                class="rounded-box"
                :disabled="deleteMode"
                @click="modal = true"
            ></button>
        </div>
    </tr>
</template>

<script lang="ts" setup>
import { computed, ref, reactive, watch, toRef } from "vue";
import { useToast } from "vue-toastification";
import type {
    FileEditUpdate,
    FileUpdate,
    FileListingData,
} from "@/types/api/@me/f";

import { useValidComment, useValidPassword } from "@/use/useValidate";
import { useStorage } from "@/use/useStorage";
import { useRequest } from "@/use/useAPI";

import ValidPasswordItem from "@/components/valid/ValidPasswordItem.vue";
import ValidCommentItem from "@/components/valid/ValidCommentItem.vue";
import ModalItem from "@/components/@me/util/ModalItem.vue";
import FileViewItem from "@/components/@me/util/FileViewItem.vue";
import LoadingFileViewItemVue from "@/components/loading/LoadingFileViewItem.vue";

import moment from "moment";
import Swal from "sweetalert2";

const modal = ref(false);
const loading = ref(false);
const selected = ref(false);

const deleteCheckbox = ref<HTMLInputElement>();
const deleteCheckboxMobile = ref<HTMLInputElement>();

const toast = useToast();
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
const fileSize = computed(() => useStorage(props.data.file_size));
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

const copyID = () => {
    navigator.clipboard.writeText(`${import.meta.env.VITE_API}/f/${props.data.id}`).then(
        () => {
            toast.success(`Copied ${props.data.id} to clipboard`);
        },
        (err) => {
            console.warn(`[copyID] could not copy id to clipboard, ${err}`);
            toast.error(`Could not copy ${props.data.id} to clipboard`);
        }
    );
};

const updateSelect = () => {
    if (!deleteCheckbox.value) return;

    selected.value = deleteCheckbox.value.checked;
    emit("updatePurgeList", props.data.id, !selected.value);
};

const updateSelectMobile = () => {
    if (!deleteCheckboxMobile.value) return;

    selected.value = deleteCheckboxMobile.value.checked;
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
