<template>
    <div v-if="mobile" class="lg:hidden block m-auto">
        <button
            class="btn btn-primary"
            @click="updatePage(-1)"
            :disabled="disabled || page <= 0"
        >
            &lt;
        </button>

        <span class="mx-10 text-2xl font-semibold">{{ page + 1 }}</span>

        <button
            class="btn btn-primary"
            @click="updatePage(1)"
            :disabled="disabled"
        >
            &gt;
        </button>
    </div>
    <div v-else class="hidden lg:block">
        <button
            class="btn btn-xs btn-outline btn-primary"
            @click="updatePage(-1)"
            :disabled="disabled || page <= 0"
        >
            &lt;
        </button>

        <span class="mx-2">{{ page + 1 }}</span>

        <button
            class="btn btn-xs btn-outline btn-primary"
            @click="updatePage(1)"
            :disabled="disabled"
        >
            &gt;
        </button>
    </div>
</template>

<script lang="ts" setup>
const props = defineProps<{
    page: number;
    disabled: boolean;
    mobile: boolean;
}>();

const emits = defineEmits<{
    (e: "clicked", evt: number): void;
}>();

const updatePage = (n: number) => {
    if (props.page + n < 0) {
        return;
    }

    emits("clicked", props.page + n);
};
</script>
