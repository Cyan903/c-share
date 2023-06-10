<template>
    <div class="bg-base-300 rounded-box shadow-xl p-6">
        <h1 class="sm:text-4xl text-3xl font-bold text-center">
            ShareX Custom Uploader
        </h1>

        <h4 class="text-sm text-center font-bold opacity-60 my-5">
            Automate uploading with
            <a href="https://getsharex.com/" class="link hover:opacity-60">
                ShareX.
            </a>
        </h4>

        <div class="flex lg:flex-nowrap flex-wrap">
            <div class="w-full text-center ml-auto my-auto input-groups">
                <input
                    :class="{ 'input-error': !query.token }"
                    v-model="query.token"
                    type="text"
                    placeholder="API Token"
                    class="input input-bordered api-token"
                />

                <div v-if="query.perm == 'unlisted'" class="center-group">
                    <ValidPasswordItem
                        v-model="query.password"
                        placehold="Default Password"
                    />

                    <ValidPasswordConfirmItem
                        v-model="query.confirm"
                        :check="query.password"
                    />
                </div>

                <div class="center-group">
                    <ValidCommentItem v-model="query.comment" />
                </div>

                <label class="input-group block mb-5">
                    <button class="btn btn-active no-animation">
                        Permission
                    </button>
                    <select class="select select-bordered" v-model="query.perm">
                        <option value="public">Public</option>
                        <option value="private">Private</option>
                        <option value="unlisted">Unlisted</option>
                    </select>
                </label>

                <button class="btn btn-primary" :disabled="!valid">
                    <a :href="download" download="c-share.sxcu"
                        >Download Uploader</a
                    >
                </button>
            </div>

            <div class="w-full lg:block hidden mr-auto input-groups">
                <div class="mockup-code">
                    <pre class="px-5">
                        <VueJsonPretty :data="sxcu" />
                    </pre>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed, reactive, toRef } from "vue";
import { useValidPassword, useValidComment } from "@/use/useValidate";
import type { SXCU } from "@/types/share/config";

import ValidPasswordItem from "@/components/valid/ValidPasswordItem.vue";
import ValidPasswordConfirmItem from "@/components/valid/ValidPasswordConfirmItem.vue";
import ValidCommentItem from "@/components/valid/ValidCommentItem.vue";

import VueJsonPretty from "vue-json-pretty";

const query = reactive({
    token: "",
    perm: "public",
    comment: "",
    password: "",
    confirm: "",
});

const valid = computed(() => {
    if (query.perm == "unlisted") {
        return (
            useValidPassword(toRef(query.password)) &&
            useValidComment(toRef(query.comment)) &&
            query.password == query.confirm &&
            query.token != ""
        );
    }

    return useValidComment(toRef(query.comment)) && query.token != "";
});

const sxcu = computed(() => {
    const data: SXCU = {
        Version: "14.0.0",
        Name: import.meta.env.VITE_APP,
        DestinationType: "ImageUploader, TextUploader, FileUploader",
        RequestMethod: "POST",
        RequestURL: `${import.meta.env.VITE_API}/api/upload`,
        Parameters: {
            token: query.token,
            perm: query.perm,
            comment: query.comment,
        },

        Body: "MultipartFormData",
        FileFormName: "upload",
        URL: `${import.meta.env.VITE_API}/f/{json:data.id}`,
        ErrorMessage: "{json:code} - {json:message}",
    };

    if (query.perm == "unlisted") {
        data.Parameters.password = query.password;
    }

    return { ...data };
});

const download = computed(
    () => "data:text/plain," + encodeURIComponent(JSON.stringify(sxcu.value))
);
</script>
