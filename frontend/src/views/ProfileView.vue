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
import { useRoute } from 'vue-router'
import TopBar from '@/components/TopBar.vue'
import PostsList from '@/components/PostsList.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'

const route = useRoute()
const user = ref(null)
const posts = ref([])
const apiUrl = import.meta.env.VITE_API_URL

async function fetchUserAndPosts(userId) {
    // Fetch user info
    const userRes = await fetch(`${apiUrl}/api/users/${userId}`)
    user.value = userRes.ok ? await userRes.json() : null

    // Fetch and filter posts
    const postsRes = await fetch(`${apiUrl}/api/posts`)
    if (postsRes.ok) {
        const allPosts = await postsRes.json()
        posts.value = allPosts.filter(p => p.user_id === Number(userId))
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