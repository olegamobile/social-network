<template>
    <div class="follows-page">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <h3 class="text-lg font-semibold">Following</h3>
                <ul class="list-disc pl-6">
                    <li>@bob</li>
                    <li>@john</li>
                </ul>

                <h3 class="text-lg font-semibold">Followers</h3>
                <ul class="list-disc pl-6">
                    <li>@bob</li>
                </ul>
            </template>

            <template #main>
                <h2 class="text-2xl font-bold mb-4">Explore Users</h2>
                <div>
                    <div class="max-w-lg w-full lg:max-w-xs mb-4">
                        <label for="search" class="sr-only">Search</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <i class="fas fa-search text-nordic-light"></i>
                            </div>
                            <input id="search" name="search"
                                   class="block w-full pl-10 pr-3 py-2 border border-nordic-light rounded-md leading-5 bg-white placeholder-nordic-light focus:outline-none focus:ring-2 focus:ring-nordic-secondary-accent focus:border-nordic-secondary-accent sm:text-sm"
                                   placeholder="Search users..." type="search" v-model="searchQuery" @input="searchUsers" >
                        </div>
                    </div>

                    <div v-if="searchResults && searchResults.length > 0">
                        <h3 class="text-lg font-semibold mb-1">Search Results:</h3>
                        <ul class="list-disc pl-6">
                            <li v-for="user in searchResults" :key="user.id">
                                <RouterLink :to="`/profile/${user.id}`">{{ user.first_name }} {{ user.last_name }}</RouterLink>
                                <span v-if="user.username"> - {{ user.username }}</span>
                            </li>
                        </ul>
                    </div>
                    <p v-else-if="searchInitiated" class="text-base text-gray-700 leading-relaxed mb-4">No users found.</p>
                </div>

                <h3 class="text-lg font-semibold mv-1">Suggested Users</h3>
                <ul class="list-disc pl-6">
                    <li v-for="user in users" :key="user.id">
                        <RouterLink :to="`/profile/${user.id}`">{{ user.first_name }} {{ user.last_name }}</RouterLink>
                        <span v-if="user.username"> - {{ user.username }}</span>
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
        if (response.status === 401) {
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
