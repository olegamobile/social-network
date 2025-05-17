<template>
    <div class="profile-page">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>

                <!-- follow button -->
                <div v-if="showFollowButton" class="mb-2">
                    <button :disabled="followStatus === 'pending'" @click="handleFollowAction"
                        class="px-4 py-2 rounded text-white" :class="followButtonClass">
                        {{
                            followStatus === 'not following public' ? 'Follow' :
                                followStatus === 'not following private' ? 'Request to Follow' :
                                    followStatus === 'accepted' ? 'Stop Following' :
                                        followStatus === 'pending' ? 'Follow Requested' :
                                            ''
                        }}
                    </button>
                </div>

                <!-- avatar image -->
                <div class="flex flex-col items-center mb-4">
                    <div v-if="user?.avatar_url"
                        class="profile-avatar w-24 h-24 rounded-full overflow-hidden border border-nordic-light">
                        <img :src="`${apiUrl}/${user.avatar_url}`" alt="User Avatar"
                            class="w-full h-full object-cover" />
                    </div>
                </div>

                <!-- first and last name -->
                <h2 class="text-lg font-semibold">{{ user?.first_name }} {{ user?.last_name }}</h2>
                <br />
                <!-- info: show if following or public -->
                <p v-if="user?.email"><strong>Email:</strong> {{ user?.email }}</p>
                <p v-if="formattedBirthday"><strong>Birthday:</strong> {{ formattedBirthday }}</p>
                <p v-if="user?.username"><strong>Username:</strong> {{ user?.username }}</p>
                <p v-if="user?.about_me"><strong>About:</strong> {{ user?.about_me }}</p>
                <br />
                <!-- profile type -->
                <p v-if="user?.is_public">Public profile</p>
                <p v-if="!user?.is_public">Private profile</p>
            </template>

            <template #main>
                <!-- edit profile button and form -->
                <button v-if="userStore.user && route.params.id == userStore.user.id"
                    @click="showEditForm = !showEditForm"
                    class="mb-4 px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition">
                    {{ showEditForm ? 'Close Editor' : 'Edit Profile' }}
                </button>
                <EditProfile v-if="showEditForm" />

                <!-- user's posts -->
                <h2 class="text-2xl font-bold mb-4">{{ user?.first_name }}'s Posts</h2>
                <PostsList :posts="posts" />
            </template>

        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import TopBar from '@/components/TopBar.vue'
import PostsList from '@/components/PostsList.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { useAuth } from '@/composables/useAuth'
import { useErrorStore } from '@/stores/error'
import EditProfile from '@/components/EditProfile.vue'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const user = ref(null)
const posts = ref([])
const apiUrl = import.meta.env.VITE_API_URL
const { logout } = useAuth()
const router = useRouter()
const errorStore = useErrorStore()
const showEditForm = ref(false);
const userStore = useUserStore()
const followStatus = ref('')

// Compute formatted birthday so it doesn't affect userStore
const formattedBirthday = computed(() => {
    if (user.value && user.value.birthday) {
        return new Date(user.value.birthday).toLocaleString("fi-FI", {
            dateStyle: 'short',
        });
    }
    return '';
});

const showFollowButton = computed(() =>
    ['not following private', 'not following public', 'accepted', 'pending'].includes(followStatus.value)
)

const followButtonClass = computed(() => {
    if (followStatus.value === 'not following public' || followStatus.value === 'not following private') {
        return 'bg-nordic-primary-accent hover:bg-nordic-secondary-accent text-white';
    }
    if (followStatus.value === 'accepted') {
        //return 'bg-nordic-text-light hover:bg-nordic-primary-accent text-black';  // doesn't work for some reason
        return 'bg-nordic-primary-accent hover:bg-nordic-secondary-accent text-white';
    }
    if (followStatus.value === 'pending') {
        return 'bg-nordic-secondary-bg text-nordic-light cursor-not-allowed';
    }
    return '';
});

async function fetchUserAndPosts(userId) {
    try {
        if (userId != userStore.user.id) {

            // Fetch user info
            const userRes = await fetch(`${apiUrl}/api/users/${userId}`, {
                credentials: 'include' // Necessary to send cookie all the way to backend server
            })

            if (userRes.status === 401) {
                // Session is invalid — logout and redirect
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
        } else {
            user.value = userStore.user
        }

        // Get follow status
        const followRes = await fetch(`${apiUrl}/api/following/${userId}`, {     //
            credentials: 'include'
        })

        if (!followRes.ok) {
            throw new Error(`Failed to fetch follow info: ${followRes.status}`)
        }
        followStatus.value = await followRes.json()
        console.log("follow at profile view:", followStatus.value)

        // Fetch and filter posts
        const postsRes = await fetch(`${apiUrl}/api/posts/${userId}`, {     //
            credentials: 'include'
        })

        if (!postsRes.ok) {
            throw new Error(`Failed to fetch posts: ${postsRes.status}`)
        }
        const userPosts = await postsRes.json()
        if (userPosts) posts.value.push(...userPosts)

    } catch (err) {
        console.log("error fetching posts:", err)
        errorStore.setError('Error', 'Something went wrong while loading user data.')
        router.push('/error')
        return
    }
}

async function handleFollowAction() {
    if (followStatus.value === 'pending') return

    let action = ''
    if (followStatus.value === 'not following private') action = 'request'
    else if (followStatus.value === 'not following public') action = 'follow'
    else if (followStatus.value === 'accepted') action = 'unfollow'

    try {
        const res = await fetch(`${apiUrl}/api/follow`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({
                target_id: user.value.id,
                action: action
            })
        })

        if (!res.ok) throw new Error('Failed to update follow status')

        // Get follow status again
        const followRes = await fetch(`${apiUrl}/api/following/${user.value.id}`, {     //
            credentials: 'include'
        })
        if (!followRes.ok) {
            throw new Error(`Failed to fetch follow info: ${followRes.status}`)
        }
        followStatus.value = await followRes.json()
    } catch (err) {
        console.error(err)
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

// Update own profile when userstore.user changes
watch(
    () => userStore.user,
    (newUser) => {
        //console.log("new user in Profileview:", newUser)
        if (newUser && route.params.id == newUser.id) {
            user.value = newUser
        }
    }
)

</script>

<style scoped>
.profile-page {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
}
</style>