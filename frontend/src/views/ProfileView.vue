<template>
    <div class="layout">
        <TopBar :userId="userId" :username="user?.username || 'Profile'" />

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
                <div v-for="post in posts" :key="post.id" class="post">
                    <p>{{ post.content }}</p>
                    <small>Posted on {{ post.created_at }}</small>
                </div>
            </main>
        </div>
    </div>
</template>

<script setup>
import TopBar from '../components/TopBar.vue'
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const userId = parseInt(route.params.id)
const user = ref(null)
const posts = ref([])

onMounted(async () => {
    const usersRes = await fetch('http://localhost:8080/api/users')
    const allUsers = await usersRes.json()
    user.value = allUsers.find(u => u.id === userId)

    const postsRes = await fetch('http://localhost:8080/api/posts')
    const allPosts = await postsRes.json()
    posts.value = allPosts.filter(p => p.user_id === userId)
})
</script>

<style scoped>
.layout {
    display: flex;
    flex-direction: column;
    height: 100vh;
}

.top-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #333;
    color: white;
    padding: 0.5rem 1rem;
}

.nav-icons a,
.profile-link {
    margin-left: 1rem;
    color: white;
    text-decoration: none;
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