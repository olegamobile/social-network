<template>
    <div class="group-view">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>

                <!-- join/leave button -->
                <div v-if="showJoinLeaveButton" class="mb-2">
                    <button @click="joinOrLeave" class="px-4 py-2 rounded text-white" :class="followButtonClass">
                        {{
                            membershipStatus === 'pending' ? 'Request Sent' :
                                membershipStatus === 'accepted' ? 'Leave Group' :
                                    membershipStatus === 'declined' ? 'Request to Join' :
                                        membershipStatus === 'admin' ? 'Remove Group' :
                                            ''
                        }}
                    </button>
                </div>

                <div v-if="group">
                    <h2 class="text-2xl font-bold mb-4">{{ group.title }}</h2>
                    <p>{{ group.description }}</p>

                    <div v-if="isMember">

                        <EventList :events="events" small class="my-4" />
                        <RouterLink :to="`/chats/${group.id}`" class="my-4">Go to Chat</RouterLink>

                        <h4 class="mt-4">Members</h4>
                        <ul>
                            <li v-for="member in members" :key="member.id">
                                {{ member.username }}
                            </li>
                        </ul>
                    </div>
                </div>
            </template>

            <template #main>
                <div v-if="isMember && posts">
                    <PostsList :posts="posts" />
                </div>
                <div v-else>
                    <button @click="requestMembership">No membership or no posts</button>
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
const group = ref(null)
const posts = ref([])
const members = ref([])
const events = ref([])

const isMember = ref(true) // Mock check
const showJoinLeaveButton = ref[true]


function requestMembership() {
    alert('Membership requested!')
}

async function getGroup(groupId) {
    try {
        // Fetch user info
        const groupRes = await fetch(`${apiUrl}/api/group/${groupId}`, {
            credentials: 'include'
        })

        if (groupRes.status === 401) { // Session is invalid — logout and redirect
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
    } catch (err) {
        console.log("error fetching group:", err)
        errorStore.setError('Error', 'Something went wrong while loading group data.')
        router.push('/error')
        return
    }
}

async function getPosts(groupId) {        // Fetch and filter posts
    try {
        const postsRes = await fetch(`${apiUrl}/api/group/posts/${groupId}`, {     //
            credentials: 'include'
        })

        if (!postsRes.ok) {
            throw new Error(`Failed to fetch posts: ${postsRes.status}`)
        }
        const groupPosts = await postsRes.json()
        if (groupPosts) posts.value.push(...groupPosts)
    } catch (error) {
        console.log("error fetching group posts:", error)
        errorStore.setError('Error', 'Something went wrong while loading group posts data.')
        router.push('/error')
        return
    }
}

async function getMembers(groupId) {        // Fetch and filter posts
    try {
        const membsRes = await fetch(`${apiUrl}/api/group/members/${groupId}`, {     //
            credentials: 'include'
        })

        if (!membsRes.ok) {
            throw new Error(`Failed to fetch posts: ${membsRes.status}`)
        }
        const groupMembs = await membsRes.json()
        if (groupMembs) members.value = groupMembs
    } catch (error) {
        console.log("error fetching group members:", error)
        errorStore.setError('Error', 'Something went wrong while loading group members data.')
        router.push('/error')
        return
    }
}

async function getEvents(groupId) {        // Fetch and filter posts
    try {
        const eventsRes = await fetch(`${apiUrl}/api/group/events/${groupId}`, {     //
            credentials: 'include'
        })

        if (!eventsRes.ok) {
            throw new Error(`Failed to fetch posts: ${eventsRes.status}`)
        }
        const groupEvs = await eventsRes.json()
        if (groupEvs) events.value = groupEvs

        console.log("events gotten:", events.value)
    } catch (error) {
        console.log("error fetching group events:", error)
        errorStore.setError('Error', 'Something went wrong while loading group events data.')
        router.push('/error')
        return
    }
}

onMounted(() => {
    getGroup(route.params.id)
    getPosts(route.params.id)
    getMembers(route.params.id)
    getEvents(route.params.id)
})
</script>
