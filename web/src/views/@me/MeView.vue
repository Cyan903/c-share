<template>
    <div>
        <div>
            <input
                type="text"
                :disabled="deleteMode"
                v-model="query.search"
                placeholder="Search..."
            />

            <button :disabled="deleteMode" @click="modals.filter = true">
                Filter
            </button>
        </div>

        <div>
            <button :disabled="deleteMode" @click="modals.add = true">
                Add
            </button>

            <button @click="deleteMode = !deleteMode">
                {{ deleteMode ? "Cancel" : "Remove" }}
            </button>

            <button
                v-if="deleteMode"
                @click="purgeFiles"
                :disabled="deleteList.length <= 0"
            >
                Delete
            </button>
        </div>

        <ModalItem :show="modals.filter" @hide="modals.filter = false">
            <div>
                <span>File Type</span> |
                <input
                    type="text"
                    v-model="query.type"
                    placeholder="text/html"
                />
            </div>

            <div>
                <span>Listing</span> |
                <select v-model="query.listing">
                    <option value="any">Any</option>
                    <option value="public">Public</option>
                    <option value="private">Private</option>
                    <option value="unlisted">Unlisted</option>
                </select>
            </div>

            <div>
                <span>Order</span> |
                <select v-model="query.order">
                    <option value="size">Size</option>
                    <option value="type">Type</option>
                    <option value="comment">Comment</option>
                    <option value="permission">Permission</option>
                    <option value="date">Date</option>
                </select>
            </div>
        </ModalItem>

        <ModalItem :show="modals.add" @hide="modals.add = false">
            <UploadFileItem @fileUpload="uploadFile" />
        </ModalItem>

        <div>
            <PageScrollItem
                :page="parseInt(query.page)"
                :disabled="deleteMode"
                @clicked="(n) => (query.page = String(n))"
            />

            <DisplayOrderItem
                :total="total"
                :order="query.order"
                :listing="query.listing"
                :type="query.type"
            />

            <SortButtonItem
                :mode="query.sort"
                :disabled="deleteMode"
                @clicked="(n) => (query.sort = n)"
            />
        </div>

        <div>
            <b>Used Storage: </b>
            <span>{{ usedStorage }}</span>
        </div>

        <Loading v-if="loading" :loading="loading" />
        <div v-else>
            <table v-if="!nothingFound">
                <tr>
                    <th>Image</th>
                    <th>File ID</th>
                    <th>File Type</th>
                    <th>File Comment</th>
                    <th>Permissions</th>
                    <th>Upload Date</th>
                    <th></th>
                    <th></th>
                </tr>

                <FileListItem
                    v-for="d in data"
                    :key="d.id"
                    :data="d"
                    :deleteMode="deleteMode"
                    @updatePurgeList="updateDeleteList"
                    @editFile="updateFile"
                />
            </table>
            <h1 v-else>Nothing found...</h1>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { reactive, ref, watch, onMounted, computed } from "vue";

import { useRequest } from "@/use/useAPI";
import { useAuthStore } from "@/stores/auth";
import type { FileUploadRequest, FileDeleteRequest } from "@/types/api/@me";
import type {
    FileUpload,
    FileUpdate,
    FileListing,
    FileListingData,
} from "@/types/api/@me/f";

import UploadFileItem from "@/components/@me/UploadFileItem.vue";
import SortButtonItem from "@/components/@me/SortButtonItem.vue";
import PageScrollItem from "@/components/@me/PageScrollItem.vue";
import DisplayOrderItem from "@/components/@me/DisplayOrderItem.vue";
import FileListItem from "@/components/@me/FileListItem.vue";

import ModalItem from "@/components/@me/util/ModalItem.vue";
import Loading from "@/components/LoadingItem.vue";

import Swal from "sweetalert2";

// TODO: Limit updates without watch()

const auth = useAuthStore();
const data = ref(Array<FileListingData>(0));
const total = ref(0);
const deleteList = ref(Array<string>(0));

const loading = ref(false);
const nothingFound = ref(false);
const deleteMode = ref(false);

const query = reactive({
    page: String(0),
    listing: "any",
    type: "",
    order: "date",
    sort: "desc",
    search: "",
});

const modals = reactive({
    filter: false,
    add: false,
});

const filterFiles = async () => {
    const token = localStorage.getItem("token");
    const type = query.type || "any";
    const params = new URLSearchParams({
        page: query.page,
        listing: query.listing,
        type,
        order: query.order,
        sort: query.sort,
        search: query.search,
    }).toString();

    if (!token) {
        console.warn("[@me] No token found in storage!");
        return;
    }

    const req = await useRequest<FileListing>(
        `/@me/f?${params}`,
        {
            method: "GET",
            headers: { "Content-Type": "application/json", Token: token },
        },

        loading
    );
    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not filter files!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        return;
    }

    if (req.json.data) {
        data.value = req.json.data.slice(0);
        total.value = req.json.count;
        nothingFound.value = false;

        return;
    }

    data.value = [];
    total.value = 0;
    nothingFound.value = true;
};

const usedStorage = computed(() => {
    const units = [
        "bytes",
        "KiB",
        "MiB",
        "GiB",
        "TiB",
        "PiB",
        "EiB",
        "ZiB",
        "YiB",
    ];

    let n = parseInt(String(auth.userData.usedStorage), 10) || 0;
    let l = 0;

    while (n >= 1024 && ++l) {
        n = n / 1024;
    }

    return n.toFixed(n < 10 && l > 0 ? 1 : 0) + " " + units[l];
});

const updateFile = (file: FileUpdate) => {
    const idx = data.value.map((n) => n.id).indexOf(file.id);

    data.value[idx].file_comment = file.comment;
    data.value[idx].permissions = ["public", "private", "unlisted"].indexOf(
        file.perm
    );
};

// here
const uploadFile = async (file: FileUpload) => {
    const uploadData = new FormData();
    const token = localStorage.getItem("token");

    let params = new URLSearchParams({
        perm: file.perm,
        comment: file.comment,
    }).toString();

    if (file.password) {
        params = new URLSearchParams({
            perm: file.perm,
            password: file.password,
            comment: file.comment,
        }).toString();
    }

    if (!token) {
        console.warn("[@me/upload] No token found in storage!");
        return;
    }

    uploadData.append("upload", file.file);

    const req = await useRequest<FileUploadRequest>(
        `/@me/upload?${params}`,
        {
            method: "POST",
            headers: { Token: token },
            body: uploadData,
        },

        loading
    );

    if (req.error) return;
    if (req.response.status != 200) {
        Swal.fire({
            title: "Could not upload file!",
            text: req.json.message,
            icon: "warning",
            confirmButtonText: "Okay",
        });

        return;
    }

    modals.add = false;
    auth.updateStorage(parseInt(req.json.data.storage));

    filterFiles();
};

const updateDeleteList = (id: string, del: boolean) => {
    if (del) {
        deleteList.value = deleteList.value.filter((n) => n != id);
        return;
    }

    deleteList.value.push(id);
};

// here
const purgeFiles = () => {
    Swal.fire({
        title: "Are you sure you want to delete these files?",
        text: `${deleteList.value.length} file(s) selected`,
        showCancelButton: true,
        confirmButtonText: "Delete",
    }).then(async (result) => {
        const token = localStorage.getItem("token");

        if (!token) return;
        if (!result.isConfirmed) return;

        const req = await useRequest<FileDeleteRequest>(
            `/@me/upload`,
            {
                method: "DELETE",
                headers: { Token: token },
                body: JSON.stringify(deleteList.value),
            },

            loading
        );

        deleteMode.value = false;
        deleteList.value.length = 0;

        if (req.error) return;
        if (req.response.status != 200) {
            Swal.fire({
                title: "Could not delete files!",
                text: req.json.message,
                icon: "warning",
                confirmButtonText: "Okay",
            });

            return;
        }

        auth.updateStorage(parseInt(req.json.message));
        filterFiles();
    });
};

// Reset page and update data
const resetPage = () => (query.page = String(0));

watch(() => query.listing, resetPage);
watch(() => query.type, resetPage);
watch(() => query.order, resetPage);
watch(() => query.search, resetPage);

watch(deleteMode, () => (deleteList.value.length = 0));
watch(query, filterFiles);

onMounted(filterFiles);
</script>
