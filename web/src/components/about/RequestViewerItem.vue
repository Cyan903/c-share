<template>
    <div class="flex lg:flex-nowrap flex-wrap">
        <div class="m-2 w-full">
            <div class="mockup-code lg:h-4/5">
                <pre><code class="javascript">{{ data.code }}</code></pre>
            </div>

            <div class="w-full text-center">
                <div class="tabs inline-block tabs-boxed mt-2">
                    <a
                        v-for="tab in Object.keys(code)"
                        :class="{ 'tab-active': active.code == tab }"
                        :key="tab"
                        @click="setActive('code', tab)"
                        class="tab"
                    >
                        {{ tab }}
                    </a>
                </div>
            </div>
        </div>

        <div class="m-2 w-full text-center">
            <div class="mockup-code lg:h-4/5">
                <pre class="px-5">
                    <VueJsonPretty :data="data.response" />
                </pre>
            </div>

            <div class="w-full text-center">
                <div class="tabs inline-block tabs-boxed mt-2">
                    <a
                        v-for="tab in Object.keys(response)"
                        :class="{ 'tab-active': active.response == tab }"
                        :key="tab"
                        @click="setActive('response', tab)"
                        class="tab"
                    >
                        {{ tab }}
                    </a>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, reactive } from "vue";
import VueJsonPretty from "vue-json-pretty";

const props = defineProps<{
    defaultCode: string;
    code: {
        [key: string]: string;
    };

    response: {
        [key: string]: object;
    };
}>();

const active = reactive({
    code: "",
    response: "200",
});

const data = reactive({
    code: {},
    response: {},
});

const setActive = (mode: "code" | "response", code: string) => {
    active[mode] = code;
    data[mode] = props[mode][code];
};

onMounted(() => {
    setActive("code", props.defaultCode);
    setActive("response", "200");
});
</script>

<style scoped>
.mockup-code {
    min-width: 100% !important;
    max-width: 300px;
    max-height: 300px;
}
</style>
