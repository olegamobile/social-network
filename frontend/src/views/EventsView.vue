<template>
  <div class="events-page-wrapper">
    <TopBar />

    <TwoColumnLayout>
      <template #sidebar>
        <div class="mb-8">
          <h3 class="text-xl font-semibold text-nordic-dark mb-3">Upcoming Events</h3>
          <ul v-if="upcoming.length > 0" class="space-y-2">
            <li
              v-for="event in upcoming"
              :key="event.id"
              @click="select(event)"
              :class="[
                'cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150',
                selectedEvent.id === event.id ? 'font-semibold text-nordic-primary-accent' : ''
              ]"
            >
              {{ event.title }}
            </li>
          </ul>
          <p v-else class="italic text-nordic-light">No upcoming events.</p>
        </div>

        <div class="mb-8">
          <h4 class="text-lg font-medium text-nordic-dark mb-2">Invited</h4>
          <ul v-if="invited.length > 0" class="space-y-2">
            <li
              v-for="event in invited"
              :key="event.id"
              @click="select(event)"
              class="cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150"
            >
              {{ event.title }}
            </li>
          </ul>
          <p v-else class="italic text-nordic-light">You're not invited to any events.</p>
        </div>

        <div>
          <h4 class="text-lg font-medium text-nordic-dark mb-2">Past</h4>
          <ul v-if="past.length > 0" class="space-y-2">
            <li
              v-for="event in past"
              :key="event.id"
              @click="select(event)"
              class="cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150"
            >
              {{ event.title }}
            </li>
          </ul>
          <p v-else class="italic text-nordic-light">No past events.</p>
        </div>
      </template>

      <template #main>
        <h2 class="text-3xl font-bold text-nordic-dark mb-6">Event Details</h2>
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
  {
    id: 1,
    title: 'Graduation Party',
    group: 'Class of 24',
    time: '2025-06-01T18:00',
    description: 'Celebrate graduation!',
    going: ['Omar', 'You']
  },
  {
    id: 2,
    title: 'Football Match',
    group: 'Football Team',
    time: '2025-04-20T15:00',
    description: 'Final match of the season',
    going: ['Dolgorsürengiin']
  },
  {
    id: 3,
    title: 'Reunion',
    group: 'Old Friends',
    time: '2024-01-10T19:00',
    description: 'Meet and greet',
    going: []
  }
])

const selectedEvent = ref(allEvents.value[0])

function select(event) {
  selectedEvent.value = event
}

const upcoming = allEvents.value.filter(e => new Date(e.time) > now)
const invited = allEvents.value.filter(e => e.going.includes('You'))
const past = allEvents.value.filter(e => new Date(e.time) < now)
</script>

<style scoped>
.events-page-wrapper {
  min-height: 100vh;
}
</style>
