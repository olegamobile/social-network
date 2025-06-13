<template>
  <div class="events-page-wrapper">
    <TopBar />

    <TwoColumnLayout>

      <template #sidebar>
        <div class="mb-8">
          <h3 class="text-xl font-semibold text-nordic-dark mb-3">Upcoming Events</h3>
          <ul v-if="upcoming && upcoming.length > 0" class="space-y-2">
            <li v-for="event in upcoming" :key="event.id" @click="select(event)" :class="[
              'cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 break-all',
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
              class="cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 break-all">
              {{ event.title }}
            </li>
          </ul>
          <p v-else class="italic text-nordic-light">You're not invited to any events.</p>
        </div>

        <div>
          <h4 class="text-xl font-semibold text-nordic-dark mb-3">Past</h4>
          <ul v-if="past.length > 0" class="space-y-2">
            <li v-for="event in past" :key="event.id" @click="select(event)"
              class="cursor-pointer text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 break-all">
              {{ event.title }}
            </li>
          </ul>
          <p v-else class="italic text-nordic-light">No past events.</p>
        </div>
      </template>

      <template #main>
        <h2 v-if="selectedEvent" class="text-3xl font-bold text-nordic-dark mb-6">Event Details</h2>
        <EventCard v-if="selectedEvent" :event="selectedEvent" @going="handleGoing" @notGoing="handleNotGoing" />
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
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAuth } from '@/composables/useAuth'
import { storeToRefs } from 'pinia';

const apiUrl = import.meta.env.VITE_API_URL || '/api'
const errorStore = useErrorStore()
const route = useRoute()
const router = useRouter()
const now = new Date()
const allEvents = ref([])
const userStore = useUserStore()
const { user } = storeToRefs(userStore)
const { logout } = useAuth()
const selectedEvent = ref(null)

function select(event) {
  selectedEvent.value = event
  router.push({ name: 'events', params: { id: event.id } })
}

const upcoming = computed(() => allEvents.value.filter(e => new Date(e.event_datetime) > now))
const past = computed(() => allEvents.value.filter(e => new Date(e.event_datetime) < now))
const invited = computed(() => allEvents.value.filter(e => {
  if (e.no_response && e.no_response.length > 0) {
    return e.no_response.some(invited => invited.id === user.value.id)
  } else {
    return false
  }
}))

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
  } catch (error) {
    errorStore.setError('Error Loading Events', error.message || 'An unexpected error occurred while trying to load events. Please try again later.');
    router.push('/error')
    return
  }
}

async function respondToEvent(event_id, user_id, response) {
  try {
    const res = await fetch(`${apiUrl}/api/events/respond`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ event_id, user_id, response }),
    })

    if (res.status === 401) { // Session is invalid — logout and redirect
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

    // update selectedEvent
    await getEvents()
    const updated = allEvents.value.find(e => e.id === event_id)
    if (updated) {
      selectedEvent.value = updated
    }

  } catch (error) {
    errorStore.setError('Error Loading Events', error.message || 'An unexpected error occurred while trying to load events. Please try again later.');
    router.push('/error')
    return
  }
}


function handleGoing(event_id) {
  respondToEvent(event_id, user.id, "going")
}

function handleNotGoing(event_id) {
  respondToEvent(event_id, user.id, "not_going")
}


onMounted(async () => {
  await getEvents()

  const idFromRoute = route.params.id
  if (idFromRoute) {
    const found = allEvents.value.find(e => e.id === parseInt(idFromRoute))
    if (found) {
      selectedEvent.value = found
    } else {
      errorStore.setError('Event Not Found', 'The event you are looking for does not exist or you do not have access.')
      router.push('/error')
    }
  }
})
</script>

<style scoped>
.events-page-wrapper {
  min-height: 100vh;
}
</style>
