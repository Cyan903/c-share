<template>
    <div>
        <h1>Profile</h1>
        <ul>
            <li>
                <router-link to="/@me/profile/email">Email</router-link>
            </li>
            <li>
                <router-link to="/@me/profile/nickname">Nickname</router-link>
            </li>
            <li>
                <router-link to="/@me/profile/password">Password</router-link>
            </li>
        </ul>

        <div>
            <RouterView />
        </div>
    </div>
</template>

<style scoped>
ul a.router-link-exact-active {
    color: green;
}
</style>
