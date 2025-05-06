<template>
    <div class="layout">
        <TopBar />

        <div class="content">
            <aside class="sidebar">
                <h3>User Info</h3>
                <p><strong>Username:</strong> {{ user?.username }}</p>
                <p><strong>Email:</strong> {{ user?.email }}</p>
                <p><strong>Name:</strong> {{ user?.first_name }} {{ user?.last_name }}</p>
                <p><strong>Birthday:</strong> {{ user?.birthday }}</p>
            </aside>

            <main class="main-content">
                <h2>{{ user?.first_name }}'s Posts</h2>
                <PostsList :posts="posts"/>
            </main>
        </div>
    </div>
</template>

<script setup>
import TopBar from '../components/TopBar.vue'
import { ref, onMounted, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/stores/user'
import PostsList from '@/components/PostsList.vue'

const userStore = useUserStore()
const { user, userId } = storeToRefs(userStore)  // storeToRefs() to keep it reactive

const posts = ref([])

onMounted(async () => {
    const postsRes = await fetch('http://localhost:8080/api/posts')
    const allPosts = await postsRes.json()
    posts.value = allPosts.filter(p => p.user_id === userId.value)
})

</script>

<style scoped>
.layout {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
}

.content {
    display: flex;
    flex: 1;
}

.sidebar {
    width: 250px;
    background: #f5f5f5;
    padding: 1rem;
}

.main-content {
    flex: 1;
    padding: 1rem;
}

.post {
    border-bottom: 1px solid #ccc;
    padding: 0.5rem 0;
}
</style>