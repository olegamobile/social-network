<template>
    <div class="home-page">
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

                <h3>Groups</h3>
                <ul>
                    <li>Vue Fans</li>
                    <li>Go Developers</li>
                </ul>
            </template>

            <template #main>
                <h2>Home Feed</h2>
                <NewPostForm @post-submitted="handlePostSubmitted" />
                <PostsList :posts="posts" />
            </template>
        </TwoColumnLayout>

    </div>
</template>

<script setup>
import TopBar from '../components/TopBar.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import PostsList from '@/components/PostsList.vue'
import { useAuth } from '@/composables/useAuth'
import NewPostForm from '@/components/NewPostForm.vue'

const posts = ref([])
const apiUrl = import.meta.env.VITE_API_URL || '/api'
const { logout } = useAuth()
const router = useRouter()

const handlePostSubmitted = (newPost) => {
  posts.value.unshift(newPost)
}

onMounted(async () => {
    const res = await fetch(`${apiUrl}/api/posts`, {
        credentials: 'include' // This sends the session cookie with the request
    });
    if (res.status === 401) {
        // Session is invalid — logout and redirect
        logout(); // your logout function
        router.push('/login');
        return;
    }
    posts.value = await res.json()
})
</script>

<style scoped>
.home-page {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
}
</style>