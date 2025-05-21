<template>
    <div class="home-page-wrapper">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <FollowsInSidebar :userId="user.id" v-if="user" /> <!-- not user.value.id ! -->
                <GroupsInSidebar />
            </template>

            <template #main>
                <!-- title -->
                <h2 class="text-3xl font-bold text-nordic-dark mb-6">Home Feed</h2>

                <!-- new post button and form -->
                <button @click="showPostForm = !showPostForm" class="mb-4 px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition">
                    {{ showPostForm ? 'Cancel' : 'Create New Post' }}
                </button>
                <NewPostForm v-if="showPostForm" @post-submitted="handlePostSubmitted" class="mb-8" />

                <!-- feed -->
                <PostsList ref="postsListRef" :posts="posts" />
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import TopBar from '../components/TopBar.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import PostsList from '@/components/PostsList.vue'
import { useAuth } from '@/composables/useAuth'
import NewPostForm from '@/components/NewPostForm.vue'
import { useErrorStore } from '@/stores/error'
import FollowsInSidebar from '@/components/FollowsInSidebar.vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia';
import GroupsInSidebar from '@/components/GroupsInSidebar.vue'
import throttle from 'lodash.throttle';

const posts = ref([])
const apiUrl = import.meta.env.VITE_API_URL || '/api'
const { logout } = useAuth()
const router = useRouter()
const errorStore = useErrorStore()
const userStore = useUserStore()
const { user } = storeToRefs(userStore)
const showPostForm = ref(false)

let cursor = ref(null); // last post’s created_at
const limit = 10;
const lastPostId = ref(0);
const isLoading = ref(false);
const hasMore = ref(true); // To avoid loading forever
const postsListRef = ref(null); // To access the sentinel in PostsList


const handlePostSubmitted = (newPost) => {
    posts.value.unshift(newPost)
}

async function _getHomeFeed() {
    if (isLoading.value || !hasMore.value) return;
    isLoading.value = true;

    try {
        const params = new URLSearchParams();
        if (cursor.value) params.append('cursor', cursor.value);
        params.append('limit', limit);
        params.append('last_post_id', lastPostId.value);

        const res = await fetch(`${apiUrl}/api/homefeed?${params.toString()}`, {
            credentials: 'include' // This sends the session cookie with the request
        });

        if (res.status === 401) {
            // Session is invalid — logout and redirect
            console.log("homefeed returned 401 status")
            errorStore.setError('Session Expired', 'Your session has expired. Please log in again.');
            logout();
            router.push('/login')
            return;
        }

        if (!res.ok) {
            // Handle other non-successful HTTP statuses (e.g., 400, 404, 500)
            const errorData = await res.json().catch(() => ({ message: 'Failed to fetch posts and parse error.' }));
            isLoading.value = false
            throw new Error(errorData.message || `HTTP error ${res.status}`)
        }

        const newPosts = await res.json()

        // Update cursor to the created_at of the last post received
        if (newPosts && newPosts.length > 0) {
            cursor.value = newPosts[newPosts.length - 1].created_at
            lastPostId.value = newPosts[newPosts.length - 1].id
            posts.value.push(...newPosts); // append to existing posts
        } else {
            hasMore.value = false
        }

        isLoading.value = false
    } catch (error) {
        isLoading.value = false
        errorStore.setError('Error Loading Posts', error.message || 'An unexpected error occurred while trying to load posts. Please try again later.');
        router.push('/error')
        return
    }
}


// 🔁 Create a throttled version
const getHomeFeed = throttle(_getHomeFeed, 1000);

onMounted(async () => {
    await getHomeFeed()

    // Wait for posts to render so sentinel is in DOM
    await nextTick();

    const observer = new IntersectionObserver(
        (entries) => {
            const entry = entries[0];
            if (entry.isIntersecting) {
                getHomeFeed();
            }
        },
        {
            root: null, // viewport
            threshold: 0.5
        }
    );

    if (postsListRef.value?.scrollObserver) {
        observer.observe(postsListRef.value.scrollObserver);
    }
})
</script>

<style scoped>
.home-page-wrapper {
    min-height: 100vh;
}
</style>