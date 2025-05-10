<template>
    <div class="profile-page">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <h3>User Info</h3>
                <p><strong>Username:</strong> {{ user?.username }}</p>
                <p><strong>Email:</strong> {{ user?.email }}</p>
                <p><strong>Name:</strong> {{ user?.first_name }} {{ user?.last_name }}</p>
                <p><strong>Birthday:</strong> {{ user?.birthday }}</p>
            </template>

            <template #main>
                <h2>{{ user?.first_name }}'s Posts</h2>
                <PostsList :posts="posts" />
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import TopBar from '@/components/TopBar.vue'
import PostsList from '@/components/PostsList.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { useAuth } from '@/composables/useAuth'

import { useErrorStore } from '@/stores/error'

const route = useRoute()
const user = ref(null)
const posts = ref([])
const apiUrl = import.meta.env.VITE_API_URL
const { logout } = useAuth()
const router = useRouter()
const errorStore = useErrorStore()

async function fetchUserAndPosts(userId) {
    try {
        // Fetch user info
        const userRes = await fetch(`${apiUrl}/api/users/${userId}`, {
            credentials: 'include' // Necessary to send cookie all the way to backend server
        })

        if (userRes.status === 401) {
            // Session is invalid — logout and redirect
            console.log("Invalid session")
            logout()
            router.push('/login')
            return
        }

        if (userRes.status === 404) {
            errorStore.setError('User Not Found', `User with ID ${userId} does not exist.`)
            router.push('/error')
            return
        }

        if (userRes.status === 400) {
            errorStore.setError('Bad request', `Failed to get user with ID ${userId}.`)
            router.push('/error')
            return
        }

        if (!userRes.ok) {
            // Generic error
            throw new Error(`Failed to fetch user: ${userRes.status}`)
        }

        user.value = await userRes.json()

        // Fetch and filter posts
        const postsRes = await fetch(`${apiUrl}/api/posts`, {
            credentials: 'include'
        })

        if (!postsRes.ok) {
            throw new Error(`Failed to fetch posts: ${postsRes.status}`)
        }

        const allPosts = await postsRes.json()
        posts.value = allPosts.filter(p => p.user_id === Number(userId))
    } catch (err) {
        errorStore.setError('Error', 'Something went wrong while loading user data.')
        router.push('/error')
        return
    }
}


// Initial fetch
onMounted(() => {
    fetchUserAndPosts(route.params.id)
})

// React to route param changes (reload when going from one profile to another)
watch(() => route.params.id, (newId) => {
    fetchUserAndPosts(newId)
})
</script>

<style scoped>
.profile-page {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
}
</style>