<template>
    <div class="group-view">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <div v-if="group">
                    <h2 class="text-2xl font-bold mb-4">{{ group.title }}</h2>
                    <p>{{ group.description }}</p>

                    <template v-if="isMember">
                        <h4>Events</h4>
                        <EventList :events="events" small />
                        <RouterLink :to="`/chats/${group.id}`">Go to Chat</RouterLink>

                        <h4>Members</h4>
                        <ul>
                            <li v-for="member in group.members" :key="member.id">
                                {{ member.name }}
                            </li>
                        </ul>
                    </template>
                </div>
            </template>

            <template #main>
                <div v-if="isMember">
                    <PostsList :posts="group.posts" />
                </div>
                <div v-else>
                    <button @click="requestMembership">Request Membership</button>
                </div>
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import TopBar from '@/components/TopBar.vue'
import PostsList from '@/components/PostsList.vue'
import EventList from '@/components/EventList.vue'
import { useErrorStore } from '@/stores/error'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'

const route = useRoute()
const router = useRouter()
const errorStore = useErrorStore()
const apiUrl = import.meta.env.VITE_API_URL
const groupId = route.params.id
const group = ref(null)


const isMember = ref(groupId === '1') // Mock check
const events = [{
    "id": 1,
    "title": "party",
    "time": "tomorrow"
}]

function requestMembership() {
    alert('Membership requested!')
}

async function getGroup(groupId) {
    try {
        // Fetch user info
        const groupRes = await fetch(`${apiUrl}/api/group/${groupId}`, {
            credentials: 'include'
        })


        if (groupRes.status === 401) {
            // Session is invalid — logout and redirect
            logout()
            router.push('/login')
            return
        }

        if (groupRes.status === 404) {
            errorStore.setError('Group Not Found', `Group with ID ${userId} does not exist.`)
            router.push('/error')
            return
        }

        if (groupRes.status === 400) {
            errorStore.setError('Bad request', `Failed to get group with ID ${userId}.`)
            router.push('/error')
            return
        }

        if (!groupRes.ok) {
            // Generic error
            throw new Error(`Failed to fetch group: ${groupRes.status}`)
        }

        group.value = await groupRes.json()
        console.log("group value:", group.value)

    } catch (err) {
        console.log("error fetching group:", err)
        errorStore.setError('Error', 'Something went wrong while loading group data.')
        router.push('/error')
        return
    }
}

onMounted(() => {
    getGroup(route.params.id)
})
</script>
