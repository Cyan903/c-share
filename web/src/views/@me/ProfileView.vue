<template>
    <div class="drawer">
        <input id="my-drawer" type="checkbox" class="drawer-toggle" />
        <div class="drawer-content">
            <div class="auth-form flex flex-wrap items-center justify-center">
                <div class="card card-normal bg-base-300 shadow-xl">
                    <div class="card-body items-center text-center">
                        <RouterView />
                    </div>

                    <label
                        for="my-drawer"
                        class="btn btn-primary drawer-button"
                    >
                        More Settings
                    </label>
                </div>
            </div>
        </div>
        <div class="drawer-side">
            <label for="my-drawer" class="drawer-overlay"></label>
            <ul class="menu p-4 w-80 h-full bg-base-200 text-base-content">
                <li @click="$router.push('/@me/profile/email')">
                    <label
                        for="my-drawer"
                        class="drawer-button"
                        :class="{ active: active('email') }"
                    >
                        Email
                    </label>
                </li>

                <li @click="$router.push('/@me/profile/nickname')">
                    <label
                        for="my-drawer"
                        class="drawer-button"
                        :class="{ active: active('nickname') }"
                    >
                        Nickname
                    </label>
                </li>

                <li @click="$router.push('/@me/profile/password')">
                    <label
                        for="my-drawer"
                        class="drawer-button"
                        :class="{ active: active('password') }"
                    >
                        Password
                    </label>
                </li>

                <li class="lg:hidden">
                    <label for="my-drawer" class="drawer-button">Close</label>
                </li>
            </ul>
        </div>
    </div>
</template>

<script lang="ts" setup>
const active = (n: string) => {
    return location.href.split("/").includes(n);
};
</script>
