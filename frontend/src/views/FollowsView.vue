<template>
    <div class="follows-page">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <h3>Following</h3>
                <ul>
                    <li>@bob</li>
                    <li>@john</li>
                </ul>

                <h3>Followers</h3>
                <ul>
                    <li>@bob</li>
                </ul>
            </template>

            <template #main>
                <h2>Explore Users</h2>
                <div>
                    <input type="text" v-model="searchQuery" placeholder="Search users..." />
                    <button @click="searchUsers">Search</button>

                    <div v-if="searchResults && searchResults.length > 0">
                        <h3>Search Results:</h3>
                        <ul>
                            <li v-for="user in searchResults" :key="user.id">
                                <RouterLink :to="`/profile/${user.id}`">{{ user.username }}</RouterLink>
                                - {{ user.first_name }} {{ user.last_name }}
                            </li>
                        </ul>
                    </div>
                    <p v-else-if="searchInitiated">No users found.</p>
                </div>

                <h3>Suggested Users</h3>
                <ul>
                    <li v-for="user in users" :key="user.id">
                        <RouterLink :to="`/profile/${user.id}`">{{ user.username }}</RouterLink>
                        - {{ user.first_name }} {{ user.last_name }}
                    </li>
                </ul>
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import TopBar from '@/components/TopBar.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { useErrorStore } from '@/stores/error'
import { useAuth } from '@/composables/useAuth'

const apiUrl = import.meta.env.VITE_API_URL
const errorStore = useErrorStore()
const router = useRouter()
const users = ref([])
const searchQuery = ref('');
const searchResults = ref([]);
const searchInitiated = ref(false);
const { logout } = useAuth()

async function searchUsers() {
    searchInitiated.value = true;
    searchResults.value = []; // Clear previous results

    if (searchQuery.value.trim() === '') {
        return; // Don't make a request for an empty query
    }

    try {
        const response = await fetch(`${apiUrl}/api/users/search?query=${searchQuery.value}`, {
            credentials: 'include'
        });
        if (res.status === 401) {
            // Session is invalid — logout and redirect
            logout();
            router.push('/login');
            return;
        }
        if (!response.ok) {
            errorStore.setError(response.status, 'Error fetching search results')
            router.push('/error')
            return
        }
        const data = await response.json();
        searchResults.value = data;
    } catch (error) {
        console.log("Error in searching users:", String(error))
        errorStore.setError('Error', 'Something went wrong while searching for users')
        router.push('/error')
        return
    }
};


async function fetchUsers() {
    try {
        // Fetch all users
        const usersRes = await fetch(`${apiUrl}/api/users`, {
            credentials: 'include' // Necessary to send cookie all the way to backend server
        })

        if (usersRes.status === 401) { // Session is invalid — logout and redirect            
            console.log("Invalid session")
            logout()
            router.push('/login')
            return
        }

        if (usersRes.status === 404) {
            errorStore.setError('User Not Found', `User with ID ${userId} does not exist.`)
            router.push('/error')
            return
        }

        if (usersRes.status === 400) {
            errorStore.setError('Bad request', `Failed to get user with ID ${userId}.`)
            router.push('/error')
            return
        }
        if (!usersRes.ok) {
            // Generic error
            throw new Error(`Failed to fetch user: ${usersRes.status}`)
        }

        users.value = await usersRes.json()
    } catch (err) {
        errorStore.setError('Error', 'Something went wrong while loading users data.')
        router.push('/error')
        return
    }
}

onMounted(() => {
    fetchUsers()
})

</script>
