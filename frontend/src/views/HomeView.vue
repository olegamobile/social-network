<template>
    <div class="layout">
        <TopBar :userId="userId" :username="username" />

        <div class="content">
            <aside class="sidebar">
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
            </aside>

            <main class="main-content">
                <h2>Home Feed</h2>
                <div v-for="post in posts" :key="post.id" class="post">
                    <p>{{ post.content }}</p>
                    <small>Posted by user ID {{ post.user_id }} on {{ post.created_at }}</small>
                </div>
            </main>
        </div>
    </div>
</template>

<script setup>
import TopBar from '../components/TopBar.vue'
import { ref, onMounted } from 'vue'

const userId = 1
const username = 'Alice'

const posts = ref([])

onMounted(async () => {
    const res = await fetch('http://localhost:8080/api/posts')
    posts.value = await res.json()
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