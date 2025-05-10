<template>
    <nav class="top-bar">
        <router-link to="/" id="home-title" aria-label="Home">
            <h1>MySocial</h1>
        </router-link>
        <div class="nav-icons" v-if="!isLoginPage">
            <router-link to="/" class="navbar-link" title="Home" aria-label="Home">🏠</router-link>
            <router-link to="/groups" class="navbar-link" title="Groups" aria-label="Groups">👥</router-link>
            <router-link to="/chats" class="navbar-link" title="Chats" aria-label="Chats">💬</router-link>
            <router-link to="/events" class="navbar-link" title="Events" aria-label="Events">📅</router-link>
            <router-link to="/notifications" title="Notifications" aria-label="Notifications">🔔</router-link>
            <router-link v-if="user" :to="`/profile/${user.id}`" title="Your Profile" aria-label="Your Profile"
                class="profile-link">
                {{ user.username }}
            </router-link>
            <button class="logout-button" @click="logout" title="Logout" aria-label="Logout">🚪</button>
        </div>
    </nav>
</template>


<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/stores/user'
import { useAuth } from '@/composables/useAuth'

const route = useRoute();
const router = useRouter()
const userStore = useUserStore()

const { user } = storeToRefs(userStore)  // storeToRefs() ensures user is reactive when destructured
const apiUrl = import.meta.env.VITE_API_URL || '/api'
const isLoginPage = computed(() => route.path === '/login');
const { logout } = useAuth()

</script>

<style scoped>
.top-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #333;
    color: white;
    padding: 0.5rem 1rem;
    width: 100vw;
    box-sizing: border-box;
}

#home-title {
    color: white;
    text-decoration: none;
}

#home-title:hover {
    color: rgb(193, 193, 193);
}

.nav-icons {
    display: flex;
    gap: 1rem;
    align-items: center;
}

.nav-icons a,
.profile-link,
.logout-button {
    position: relative;
    display: inline-block;
    padding: 0.25rem 0.5rem;
    border-radius: 5px;
    color: white;
    text-decoration: none;
    /* background-color: transparent; */
    transition: background-color 0.2s ease;
}

.router-link-active.profile-link,
.router-link-active.navbar-link {
    background-color: #555;
}

/* .router-link-exact-active {
    font-weight: bold;
} */

.logout-button {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1rem;
}

/* Tooltip styles */
.nav-icons a:hover::after,
.profile-link:hover::after,
.logout-button:hover::after {
    content: attr(title);
    position: absolute;
    bottom: -2rem;
    left: 50%;
    transform: translateX(-50%);
    background: #444;
    color: #fff;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.75rem;
    white-space: nowrap;
    opacity: 1;
    pointer-events: none;
    z-index: 100;
}

.nav-icons a::after,
.profile-link::after,
.logout-button::after {
    content: '';
    opacity: 0;
    transition: opacity 0.2s ease;
}
</style>