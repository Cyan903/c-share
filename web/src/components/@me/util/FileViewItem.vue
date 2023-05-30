<template>
    <div v-if="small">
        <img v-if="type == 'image'" :src="source" :alt="props.id" width="35" />
    </div>
    <div class="preview" v-else>
        <img v-if="type == 'image'" :src="source" :alt="props.id" />
        <video v-else-if="type == 'video'" :src="source" controls></video>
        <audio v-else-if="type == 'audio'" :src="source" controls></audio>
        <div v-else>Unknown type!</div>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import Swal from "sweetalert2";

const source = ref("");
const type = ref("");

const props = defineProps<{
    id: string;
    type: string;
    small?: boolean;
}>();

onMounted(() => {
    const images = [
        "image/bmp",
        "image/jpeg",
        "image/x-png",
        "image/png",
        "image/gif",
    ];

    const videos = [
        "video/webm",
        "video/m4v",
        "video/avi",
        "video/mpg",
        "video/mp4",
    ];

    const audios = ["audio/mp3", "audio/x-wav", "audio/wav", "audio/ogg"];

    if (images.includes(props.type)) {
        type.value = "image";
        return display();
    }

    if (!props.small && videos.includes(props.type)) {
        type.value = "video";
        return display();
    }

    if (!props.small && audios.includes(props.type)) {
        type.value = "audio";
        return display();
    }

    if (props.small) {
        type.value = "image";
        source.value = "favicon.ico";
    }
});

const display = async () => {
    const token = localStorage.getItem("token");
    if (!token) return;

    const req = await fetch(`${import.meta.env.VITE_API}/@me/f/${props.id}`, {
        headers: {
            Token: token,
        },
    })
        .then((r) => r.blob())
        .catch(() => {
            Swal.fire({
                title: `Could not load image ${props.id}!`,
                icon: "error",
                confirmButtonText: "Okay",
            });
        });

    if (!req) return;
    source.value = URL.createObjectURL(req);
};
</script>

<style scoped>
.preview {
    text-align: center;
}

img {
    max-width: 500px;
}
</style>
