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
                <PostsList :posts="posts" />
            </template>
        </TwoColumnLayout>

    </div>
</template>

<script setup>
import TopBar from '../components/TopBar.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { ref, onMounted } from 'vue'
import PostsList from '@/components/PostsList.vue'

const posts = ref([])
const apiUrl = import.meta.env.VITE_API_URL || '/api'

onMounted(async () => {
    const res = await fetch(`${apiUrl}/api/posts`);
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