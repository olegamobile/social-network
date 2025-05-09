<template>
    <div class="events-page">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <h3>Upcoming Events</h3>
                <ul>
                    <li v-for="event in upcoming" :key="event.id" :class="{ active: selectedEvent.id === event.id }" @click="select(event)">
                        {{ event.title }}
                    </li>
                </ul>

                <h4>Invited</h4>
                <ul>
                    <li v-for="event in invited" :key="event.id" @click="select(event)">
                        {{ event.title }}
                    </li>
                </ul>

                <h4>Past</h4>
                <ul>
                    <li v-for="event in past" :key="event.id" @click="select(event)">
                        {{ event.title }}
                    </li>
                </ul>
            </template>

            <template #main>
                <EventCard :event="selectedEvent" />
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import TopBar from '@/components/TopBar.vue'
import EventCard from '@/components/EventCard.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'

const now = new Date()

const allEvents = ref([
    { id: 1, title: 'Graduation Party', group: 'Class of 24', time: '2025-06-01T18:00', description: 'Celebrate graduation!', going: ['Omar', 'You'] },
    { id: 2, title: 'Football Match', group: 'Football Team', time: '2025-04-20T15:00', description: 'Final match of the season', going: ['Dagvadorj'] },
    { id: 3, title: 'Reunion', group: 'Old Friends', time: '2024-01-10T19:00', description: 'Meet and greet', going: [] }
])

const selectedEvent = ref(allEvents.value[0])

function select(event) {
    selectedEvent.value = event
}

const upcoming = allEvents.value.filter(e => new Date(e.time) > now)
const invited = allEvents.value.filter(e => e.going.includes('You'))
const past = allEvents.value.filter(e => new Date(e.time) < now)
</script>
