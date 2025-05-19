<template>
    <div class="home-page-wrapper">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <FollowsInSidebar :userId="user.id" v-if="user"/> <!-- not user.value.id ! -->
                <br/>
                <GroupsInSidebar/>
            </template>

            <template #main>
                <h2 class="text-3xl font-bold text-nordic-dark mb-6">Home Feed</h2>

                <NewPostForm @post-submitted="handlePostSubmitted" class="mb-8" />

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
import { useErrorStore } from '@/stores/error'
import FollowsInSidebar from '@/components/FollowsInSidebar.vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia';
import GroupsInSidebar from '@/components/GroupsInSidebar.vue'

const posts = ref([])
const apiUrl = import.meta.env.VITE_API_URL || '/api'
const { logout } = useAuth()
const router = useRouter()
const errorStore = useErrorStore()
const userStore = useUserStore()
const { user } = storeToRefs(userStore)

const handlePostSubmitted = (newPost) => {
    posts.value.unshift(newPost)
}

async function getHomeFeed() {
    try {
        //const res = await fetch(`${apiUrl}/api/posts`, {
        const res = await fetch(`${apiUrl}/api/homefeed`, {
            credentials: 'include' // This sends the session cookie with the request
        });

        if (res.status === 401) {
            // Session is invalid — logout and redirect
            console.log("homefeed returned 401 status")
            errorStore.setError('Session Expired', 'Your session has expired. Please log in again.');
            logout(); // your logout function
            router.push('/login');
            return;
        }

        if (!res.ok) {
            // Handle other non-successful HTTP statuses (e.g., 400, 404, 500)
            //const errorData = await res.json().catch(() => ({ message: 'Failed to fetch posts and parse error.' }));

            const errorData = await res.json()
            console.log("homefeed returned some error status:", res.status)
            console.log("err data:", errorData)
            throw new Error(errorData.message || `HTTP error ${res.status}`);
        }

        posts.value = await res.json()
    } catch (error) {
        errorStore.setError('Error Loading Posts', error.message || 'An unexpected error occurred while trying to load posts. Please try again later.');
        router.push('/error')
        return
    }
}

onMounted(()=>{
    getHomeFeed()
})
</script>

<style scoped>
.home-page-wrapper {
    min-height: 100vh;
}
</style>