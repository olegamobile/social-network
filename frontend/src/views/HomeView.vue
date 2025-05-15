<template>
    <div class="home-page-wrapper">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <div class="mb-8">
                    <h3 class="text-xl font-semibold text-nordic-dark mb-3">Following</h3>
                    <ul v-if="mockFollowing.length > 0" class="space-y-2">
                        <li v-for="user in mockFollowing" :key="user.id"
                            class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                            {{ user.name }}
                        </li>
                    </ul>
                    <p v-else class="text-nordic-light italic">Not following anyone yet.</p>
                </div>

                <div class="mb-8">
                    <h3 class="text-xl font-semibold text-nordic-dark mb-3">Followers</h3>
                    <ul v-if="mockFollowers.length > 0" class="space-y-2">
                        <li v-for="user in mockFollowers" :key="user.id"
                            class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                            {{ user.name }}
                        </li>
                    </ul>
                    <p v-else class="text-nordic-light italic">No followers yet.</p>
                </div>

                <div>
                    <h3 class="text-xl font-semibold text-nordic-dark mb-3">Groups</h3>
                    <ul v-if="mockGroups.length > 0" class="space-y-2">
                        <li v-for="group in mockGroups" :key="group.id"
                            class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                            {{ group.name }}
                        </li>
                    </ul>
                    <p v-else class="text-nordic-light italic">Not a member of any groups yet.</p>
                </div>
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

const posts = ref([])
const apiUrl = import.meta.env.VITE_API_URL || '/api'
const { logout } = useAuth()
const router = useRouter()
const errorStore = useErrorStore()

// Mock data for sidebar (replace with actual data fetching if needed)
const mockFollowing = ref([
    { id: 'f1', name: '@bob_s CoolUser' },
    { id: 'f2', name: '@john_doe_artist' },
    { id: 'f3', name: '@alice_in_wonderdev' },
]);
const mockFollowers = ref([
    { id: 'fl1', name: '@charlie_codes' },
]);
const mockGroups = ref([
    { id: 'g1', name: 'Vue Virtuosos' },
    { id: 'g2', name: 'Nordic Design Fans' },
    { id: 'g3', name: 'Tailwind CSS Masters' },
]);

const handlePostSubmitted = (newPost) => {
    posts.value.unshift(newPost)
}

onMounted(async () => {
    try {
        const res = await fetch(`${apiUrl}/api/posts`, {
            credentials: 'include' // This sends the session cookie with the request
        });

        if (res.status === 401) {
            // Session is invalid — logout and redirect
            errorStore.setError('Session Expired', 'Your session has expired. Please log in again.');
            logout(); // your logout function
            router.push('/login');
            return;
        }

        if (!res.ok) {
            // Handle other non-successful HTTP statuses (e.g., 400, 404, 500)
            const errorData = await res.json().catch(() => ({ message: 'Failed to fetch posts and parse error.' }));
            throw new Error(errorData.message || `HTTP error ${res.status}`);
        }

        posts.value = await res.json()
    } catch (error) {
        errorStore.setError('Error Loading Posts', error.message || 'An unexpected error occurred while trying to load posts. Please try again later.');
        router.push('/error')
        return
    }

})
</script>

<style scoped>


.home-page-wrapper {
  min-height: 100vh;
}

</style>