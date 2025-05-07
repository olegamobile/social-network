<template>
    <nav class="top-bar">
        <h1>MySocial</h1>
        <div class="nav-icons">
            <router-link to="/">🏠</router-link>
            <router-link to="#">👥</router-link>
            <router-link to="#">💬</router-link>
            <router-link to="#">📅</router-link>
            <router-link to="#">🔔</router-link>
            <!-- v-if="user" condition avoids errors if the user is null before login. -->
            <router-link v-if="user" :to="`/profile/${user.id}`" class="profile-link">
                {{ user.username }}
            </router-link>
            <button class="logout-button" @click="logout">🚪</button>
        </div>
    </nav>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const { user } = storeToRefs(userStore)  // storeToRefs() ensures user is reactive when destructured
const apiUrl = import.meta.env.VITE_API_URL || '/api'

function logout() {
    //fetch('http://localhost:8080/api/logout', {
    fetch(`${apiUrl}/api/logout`, {
        method: 'POST',
        credentials: 'include',
    }).then(() => {
        userStore.clearUser()
        router.push('/login')
    })
}
</script>

<style scoped>
.top-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #333;
    color: white;
    padding: 0.5rem 1rem;
}

.nav-icons a,
.profile-link {
    margin-left: 1rem;
    color: white;
    text-decoration: none;
}

.logout-button {
    margin-left: 1rem;
    background: none;
    border: none;
    color: white;
    cursor: pointer;
    font-size: 1rem;
}
</style>