<template>
    <div>
        <div>
            <input type="text" v-model="query.search" placeholder="Search..." />
            <button @click="modals.filter = true">Filter</button>
        </div>


        <div>
            <button @click="modals.add = true">Add</button>
            <button>Remove</button>
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

                <BarListItem
                    v-for="d in data"
                    :key="d.id"
                    :data="d"
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
import type {
    FileUpload,
    FileUpdate,
    FileListing,
    FileListingData,
    FileUploadRequest,
} from "@/types/api/@me/f";

import UploadFileItem from "@/components/@me/UploadFileItem.vue";
import SortButtonItem from "@/components/@me/SortButtonItem.vue";
import PageScrollItem from "@/components/@me/PageScrollItem.vue";
import DisplayOrderItem from "@/components/@me/DisplayOrderItem.vue";
import BarListItem from "@/components/@me/list/BarListItem.vue";

import ModalItem from "@/components/@me/util/ModalItem.vue";
import Loading from "@/components/LoadingItem.vue";

import Swal from "sweetalert2";

// TODO: Limit updates without watch()

const auth = useAuthStore();
const data = ref(Array<FileListingData>(0));
const total = ref(0);
const query = reactive({
    page: String(0),
    listing: "any",
    type: "",
    order: "date",
    sort: "desc",
    search: "",
});

const loading = ref(false);
const nothingFound = ref(false);
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

    filterFiles();
    modals.add = false;
};

watch(query, filterFiles);
onMounted(filterFiles);
</script>
