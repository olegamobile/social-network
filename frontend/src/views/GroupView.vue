<template>
    <div class="group-view">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <h2>{{ group.name }}</h2>
                <p>{{ group.description }}</p>

                <template v-if="isMember">
                    <h4>Events</h4>
                    <EventList :events="group.events" small />
                    <RouterLink :to="`/chats/${group.id}`">Go to Chat</RouterLink>

                    <h4>Members</h4>
                    <ul>
                        <li v-for="member in group.members" :key="member.id">
                            {{ member.name }}
                        </li>
                    </ul>
                </template>
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
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import TopBar from '@/components/TopBar.vue'
import PostsList from '@/components/PostsList.vue'
import EventList from '@/components/EventList.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'

const route = useRoute()
const groupId = route.params.id

const isMember = ref(groupId === '1') // Mock check

const group = ref({
    name: 'Class of 24',
    description: 'The official group for the class of 2024',
    members: [{ id: 1, name: 'Omar' }, { id: 2, name: 'Dolgorsürengiin' }],
    posts: [
        { id: 1, user_id: 1, content: 'Welcome to the group!' },
        { id: 2, user_id: 2, content: 'Excited about the graduation.' }
    ],
    events: [
        { id: 1, title: 'Graduation Party', time: '2025-06-01 18:00' }
    ]
})

function requestMembership() {
    alert('Membership requested!')
}
</script>
