<template>
  <div class="events-page-wrapper">
    <TopBar />

    <TwoColumnLayout>
      <template #sidebar>
        <div class="mb-8">
          <h3 class="text-xl font-semibold text-nordic-dark mb-3">Upcoming Events</h3>
          <ul v-if="upcoming && upcoming.length > 0" class="space-y-2">
            <li v-for="event in upcoming" :key="event.id" @click="select(event)" :class="[
              'cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150',
              selectedEvent?.id === event.id ? 'text-nordic-primary-accent' : ''
            ]">
              {{ event.title }}
            </li>
          </ul>
          <p v-else class="italic text-nordic-light">No upcoming events.</p>
        </div>

        <div class="mb-8">
          <h4 class="text-xl font-semibold text-nordic-dark mb-3">Invited</h4>
          <ul v-if="invited.length > 0" class="space-y-2">
            <li v-for="event in invited" :key="event.id" @click="select(event)"
              class="cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150">
              {{ event.title }}
            </li>
          </ul>
          <p v-else class="italic text-nordic-light">You're not invited to any events.</p>
        </div>

        <div>
          <h4 class="text-xl font-semibold text-nordic-dark mb-3">Past</h4>
          <ul v-if="past.length > 0" class="space-y-2">
            <li v-for="event in past" :key="event.id" @click="select(event)"
              class="cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150">
              {{ event.title }}
            </li>
          </ul>
          <p v-else class="italic text-nordic-light">No past events.</p>
        </div>
      </template>

      <template #main>
        <h2 class="text-3xl font-bold text-nordic-dark mb-6">Event Details</h2>
        <EventCard v-if="selectedEvent" :event="selectedEvent" />
      </template>
    </TwoColumnLayout>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import TopBar from '@/components/TopBar.vue'
import EventCard from '@/components/EventCard.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { useErrorStore } from '@/stores/error'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia';

const apiUrl = import.meta.env.VITE_API_URL || '/api'
const errorStore = useErrorStore()
const router = useRouter()
const now = new Date()
const allEvents = ref([])
const userStore = useUserStore()
const { user } = storeToRefs(userStore)

const selectedEvent = ref(null)

function select(event) {
  selectedEvent.value = event
}

function printStuff() {
  console.log("ups", upcoming.value)
  //console.log("inv", invited.value)
  console.log("pss", past.value)
}

const upcoming = computed(() => allEvents.value.filter(e => new Date(e.event_datetime) > now))
const past = computed(() => allEvents.value.filter(e => new Date(e.event_datetime) < now))
const invited = computed(() => allEvents.value.filter(e => e.no_response.some(invited => invited.id === user.id)))

async function getEvents() {
  try {
    const res = await fetch(`${apiUrl}/api/events/user`, {
      credentials: 'include'
    });

    if (res.status === 401) {
      // Session is invalid — logout and redirect
      console.log("events for user returned 401 status")
      errorStore.setError('Session Expired', 'Your session has expired. Please log in again.');
      logout();
      router.push('/login')
      return;
    }

    if (!res.ok) {
      // Handle other non-successful HTTP statuses (e.g., 400, 404, 500)
      console.log("error getting user events")
      const errorData = await res.json().catch(() => ({ message: 'Failed to get events and parse error.' }));
      throw new Error(errorData.message || `HTTP error ${res.status}`)
    }

    allEvents.value = await res.json()
    console.log("user events:", allEvents.value)
  } catch (error) {
    errorStore.setError('Error Loading Events', error.message || 'An unexpected error occurred while trying to load events. Please try again later.');
    router.push('/error')
    return
  }
}


onMounted(async () => {
  await getEvents()
  printStuff()
})
</script>

<style scoped>
.events-page-wrapper {
  min-height: 100vh;
}
</style>
