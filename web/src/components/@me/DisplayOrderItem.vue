<template>
    <div class="lg:block hidden">
        <span class="font-bold">{{ total }}</span>

        <span class="font-semibold">
            Order:

            <button
                class="btn btn-xs hover:btn-error"
                :disabled="order == defaults.order"
                @click="$emit('resetFilter', 'order', defaults.order)"
            >
                {{ order }}
            </button>
        </span>

        <span class="font-semibold">
            Listing:

            <button
                class="btn btn-xs hover:btn-error"
                :disabled="listing == defaults.listing"
                @click="$emit('resetFilter', 'listing', defaults.listing)"
            >
                {{ listing }}
            </button>
        </span>

        <span class="font-semibold">
            Type:

            <button
                class="btn btn-xs hover:btn-error"
                :disabled="type == defaults.type"
                @click="$emit('resetFilter', 'type', defaults.type)"
            >
                {{ type || "Any" }}
            </button>
        </span>
    </div>
</template>

<script lang="ts" setup>
type OrdersIndex = keyof {
    page: string;
    listing: string;
    type: string;
    order: string;
    sort: string;
    search: string;
};

const defaults = {
    order: "date",
    listing: "any",
    type: "",
};

defineProps<{
    total: number;
    order: string;
    listing: string;
    type: string;
}>();

defineEmits<{
    (e: "resetFilter", filter: OrdersIndex, value: string): void;
}>();
</script>

<style scoped>
.btn {
    text-transform: capitalize;
}

span:not(span:nth-child(4)):after {
    content: " | ";
}
</style>
