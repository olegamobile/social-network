<template>
    <!-- Users I am following -->
    <div>
        <h3 class="text-xl font-semibold text-nordic-dark mb-3">Following</h3>
        <ul v-if="followedUsers && followedUsers.length > 0" class="space-y-2 mb-5">
            <li v-for="user in followedUsers" :key="user.id"
                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                <RouterLink :to="`/profile/${user.id}`">
                    {{ user.first_name }} {{ user.last_name }}<span v-if="user.username"> - {{ user.username
                    }}</span>
                </RouterLink>
            </li>
        </ul>
        <p v-else class="text-nordic-light italic mb-5">Not following anyone</p>
    </div>

    <!-- Users who follow me -->
    <div>
        <h3 class="text-xl font-semibold text-nordic-dark mb-3">Followed by</h3>
        <ul v-if="followers && followers.length > 0" class="space-y-2 mb-5">
            <li v-for="user in followers" :key="user.id"
                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                <RouterLink :to="`/profile/${user.id}`">
                    {{ user.first_name }} {{ user.last_name }}<span v-if="user.username"> - {{ user.username
                    }}</span>
                </RouterLink>
            </li>
        </ul>
        <p v-else class="text-nordic-light italic mb-5">No one following</p>
    </div>
</template>

<script setup>
import { ref, onMounted, toRefs } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useErrorStore } from '@/stores/error'
import { useAuth } from '@/composables/useAuth'

const props = defineProps({
    userId: {
        type: [String, Number], // we're using Number for now at least
        required: true
    }
})

const { userId} = toRefs(props)
const apiUrl = import.meta.env.VITE_API_URL
const errorStore = useErrorStore()
const router = useRouter()
const { logout } = useAuth()
const followedUsers = ref([])
const followers = ref([])

async function fetchFollowData() {
    try {
        const [res1, res2] = await Promise.all([
            fetch(`${apiUrl}/api/followed/${userId.value}`, { credentials: 'include' }),
            fetch(`${apiUrl}/api/followers/${userId.value}`, { credentials: 'include' })
        ])

        if (res1.status === 401 || res2.status === 401) {
            logout()
            router.push('/login')
            return
        }

        if (!res1.ok || !res2.ok) throw new Error('Failed to fetch follow data')

        const [followedJson, followersJson] = await Promise.all([
            await res1.json(),
            await res2.json()
        ])

        if (followedJson) followedUsers.value = followedJson
        if (followersJson) followers.value = followersJson

    } catch (err) {
        errorStore.setError('Error', 'Failed to load follow data.')
        router.push('/error')
    }
}

onMounted(() => {
    fetchFollowData()
})
</script>