import "./assets/main.css";
import "vue-toastification/dist/index.css";

import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import Toast, { POSITION } from "vue-toastification";
import router from "./router";

const app = createApp(App);

app.use(Toast, {
    position: POSITION.BOTTOM_LEFT,
    timeout: 1500,
});

app.use(createPinia());
app.use(router);

app.mount("#app");
