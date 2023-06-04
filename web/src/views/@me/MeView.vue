<template>
    <div class="bg-base-100 rounded-box shadow-xl p-6">
        <div class="flex md:flex-row flex-col">
            <DelayedInputItem
                :disabled="deleteMode"
                classes="input grow input-bordered mr-2 mb-2"
                v-model="query.search"
                placehold="Search..."
            />

            <button
                class="btn btn-primary mb-2"
                :disabled="deleteMode"
                @click="modals.filter = true"
            >
                Filter
            </button>
        </div>

        <div class="flex gap-1">
            <button
                class="btn btn-sm btn-accent"
                :disabled="deleteMode"
                @click="modals.add = true"
            >
                Add
            </button>

            <button
                class="btn btn-sm btn-secondary"
                @click="deleteMode = !deleteMode"
            >
                {{ deleteMode ? "Cancel" : "Remove" }}
            </button>

            <button
                v-if="deleteMode"
                class="btn btn-sm btn-accent"
                @click="purgeFiles"
                :disabled="deleteList.length <= 0"
            >
                Delete
            </button>

            <div class="grow"></div>
            <b class="hidden lg:block">Used Storage:</b>
            <span>{{ usedStorage }}</span>
        </div>

        <ModalItem :show="modals.filter" @hide="modals.filter = false">
            <h3 class="font-semibold text-2xl mb-6 text-center">File Filter</h3>
            <div class="grid justify-center align-center w-3/4 m-auto">
                <DelayedInputItem
                    v-model="query.type"
                    classes="mb-3"
                    placehold="text/html"
                    label="File Type"
                />

                <label class="input-group mb-3">
                    <button class="btn btn-active no-animation w-1/4">
                        Listing
                    </button>
                    <select
                        class="select select-bordered w-3/4"
                        v-model="query.listing"
                    >
                        <option value="any">Any</option>
                        <option value="public">Public</option>
                        <option value="private">Private</option>
                        <option value="unlisted">Unlisted</option>
                    </select>
                </label>

                <label class="input-group mb-3">
                    <button class="btn btn-active no-animation w-1/4">
                        Order
                    </button>
                    <select
                        class="select select-bordered w-3/4"
                        v-model="query.order"
                    >
                        <option value="size">Size</option>
                        <option value="type">Type</option>
                        <option value="comment">Comment</option>
                        <option value="permission">Permission</option>
                        <option value="date">Date</option>
                    </select>
                </label>

                <button @click="resetFilters" class="btn btn-error w-4/4">
                    Reset
                </button>
            </div>
        </ModalItem>

        <ModalItem :show="modals.add" @hide="modals.add = false">
            <UploadFileItem @fileUpload="uploadFile" />
        </ModalItem>

        <div class="flex my-4">
            <DisplayOrderItem
                :total="total"
                :order="query.order"
                :listing="query.listing"
                :type="query.type"
            />

            <div class="hidden lg:block grow"></div>

            <PageScrollItem
                :page="parseInt(query.page)"
                :disabled="deleteMode"
                :mobile="false"
                @clicked="(n) => (query.page = String(n))"
            />
        </div>

        <div class="divider lg:hidden"></div>

        <Loading v-if="loading" :loading="loading" />
        <div v-else class="flex flex-wrap items-center justify-center">
            <table v-if="!nothingFound" class="table w-full">
                <thead class="hidden lg:table-header-group">
                    <tr>
                        <th v-if="deleteMode">Delete</th>
                        <th>Image</th>
                        <th>File Type</th>
                        <th>Permissions</th>
                        <th>Comment</th>
                        <th>
                            <SortButtonItem
                                :mode="query.sort"
                                :disabled="deleteMode"
                                @clicked="(n) => (query.sort = n)"
                            />
                        </th>
                    </tr>
                </thead>
                <tbody class="mobile-filelist grid lg:table-row-group">
                    <FileListItem
                        v-for="d in data"
                        :key="d.id"
                        :data="d"
                        :deleteMode="deleteMode"
                        @updatePurgeList="updateDeleteList"
                        @editFile="updateFile"
                    />
                </tbody>
            </table>
            <h1 v-else class="card-title text-5xl mb-10">Nothing Found!</h1>
        </div>

        <div class="my-4">
            <div class="divider lg:hidden"></div>
            <div class="flex">
                <PageScrollItem
                    :page="parseInt(query.page)"
                    :disabled="deleteMode"
                    :mobile="true"
                    @clicked="(n) => (query.page = String(n))"
                />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { reactive, ref, watch, onMounted, computed } from "vue";

import { useRequest } from "@/use/useAPI";
import { useStorage } from "@/use/useStorage";
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

import DelayedInputItem from "@/components/@me/util/DelayedInputItem.vue";
import ModalItem from "@/components/@me/util/ModalItem.vue";
import Loading from "@/components/LoadingItem.vue";

import Swal from "sweetalert2";

// TODO: Improve DisplayOrderItem

const auth = useAuthStore();
const data = ref(Array<FileListingData>(0));
const deleteList = ref(Array<string>(0));
const total = ref(0);

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

const updateFile = (file: FileUpdate) => {
    const idx = data.value.map((n) => n.id).indexOf(file.id);

    data.value[idx].file_comment = file.comment;
    data.value[idx].permissions = ["public", "private", "unlisted"].indexOf(
        file.perm
    );
};

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

const resetFilters = () => {
    Swal.fire({
        title: "Are you sure you want to reset these filters?",
        showCancelButton: true,
        confirmButtonText: "Yes",
    }).then((result) => {
        if (!result.isConfirmed) return;

        query.page = String(0);
        query.listing = "any";
        query.type = "";
        query.order = "date";
        query.sort = "desc";
        query.search = "";
        modals.filter = false;
    });
};

const resetPage = () => (query.page = String(0));
const usedStorage = computed(() => useStorage(auth.userData.usedStorage));

// Reset page and update data
watch(() => query.listing, resetPage);
watch(() => query.type, resetPage);
watch(() => query.order, resetPage);
watch(() => query.search, resetPage);

watch(deleteMode, () => (deleteList.value.length = 0));
watch(query, filterFiles);

onMounted(filterFiles);
</script>
