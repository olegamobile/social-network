<template>
    <nav class="top-bar">
        <router-link to="/" id="home-title" aria-label="Home">
            <h1>MySocial</h1>
        </router-link>
        <div class="nav-icons" v-if="!isLoginPage">
            <router-link to="/" class="navbar-link material-icons" data-title="Home"
                aria-label="Home">home</router-link>
            <router-link to="/follows" class="navbar-link material-icons" data-title="Follows"
                aria-label="Follows">person</router-link>
            <router-link to="/groups" class="navbar-link material-icons" data-title="Groups"
                aria-label="Groups">groups</router-link>
            <router-link to="/chats" class="navbar-link material-icons" data-title="Chats"
                aria-label="Chats">chat</router-link>
            <router-link to="/events" class="navbar-link material-icons" data-title="Events"
                aria-label="Events">event</router-link>
            <router-link to="/notifications" class="navbar-link material-icons" data-title="Notifications"
                aria-label="Notifications">notifications</router-link>
            <router-link v-if="user" :to="`/profile/${user.id}`" data-title="Your Profile" aria-label="Your Profile"
                class="profile-link">
                {{ user.username }}
            </router-link>
            <button class="logout-button material-icons" @click="logout" data-title="Logout"
                aria-label="Logout">logout</button>
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
@import url('https://fonts.googleapis.com/css2?family=Material+Icons');

.top-bar {
    display: flex;
    flex-wrap: wrap;
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
    flex-wrap: wrap;
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

.material-icons,
a.material-icons {
  color: #b7d9ec; /* blue, or whatever you prefer */
}


/* Tooltip styles */
.nav-icons a:hover::after,
.profile-link:hover::after,
.logout-button:hover::after {
    content: attr(data-title);
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

    /* Override the icon font */
    font-family: sans-serif;
}

.nav-icons a::after,
.profile-link::after,
.logout-button::after {
    content: '';
    opacity: 0;
    transition: opacity 0.4s ease;
}

</style>