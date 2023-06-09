import "vue-toastification/dist/index.css";
import "vue-json-pretty/lib/styles.css";
import "./assets/main.css";

import { createApp } from "vue";
import { createPinia } from "pinia";
import Toast, { POSITION } from "vue-toastification";

import App from "./App.vue";
import router from "./router";

const app = createApp(App);

app.use(Toast, {
    position: POSITION.BOTTOM_LEFT,
    timeout: 1500,
});

app.use(createPinia());
app.use(router);

app.mount("#app");
