<template>
    <!-- Sent Requests -->
    <div>
        <h3 class="text-xl font-semibold text-nordic-dark mb-3">Sent Requests</h3>
        <ul v-if="pendingSent.length > 0" class="space-y-2 mb-5">
            <li v-for="user in pendingSent" :key="user.id"
                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer break-all">
                <RouterLink :to="`/profile/${user.id}`">
                    {{ user.first_name }} {{ user.last_name }}<span v-if="user.username"> - {{ user.username }}</span>
                </RouterLink>
            </li>
        </ul>
        <p v-else class="text-nordic-light italic mb-5">You haven't sent any requests yet.</p>
    </div>

    <!-- Received Requests -->
    <div>
        <h3 class="text-xl font-semibold text-nordic-dark mb-3">Received Requests</h3>
        <ul v-if="pendingReceived.length > 0" class="space-y-2 mb-5">
            <li v-for="user in pendingReceived" :key="user.id"
                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer break-all">
                <RouterLink :to="`/profile/${user.id}`">
                    {{ user.first_name }} {{ user.last_name }}<span v-if="user.username"> - {{ user.username }}</span>
                </RouterLink>
            </li>
        </ul>
        <p v-else class="text-nordic-light italic mb-5">You have no pending requests.</p>
    </div>

</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useErrorStore } from '@/stores/error'
import { useAuth } from '@/composables/useAuth'

const apiUrl = import.meta.env.VITE_API_URL
const errorStore = useErrorStore()
const router = useRouter()
const { logout } = useAuth()
const pendingSent = ref([])     // users I've sent requests to
const pendingReceived = ref([]) // users who sent me requests

async function fetchPendingFollows() {
  try {
    const [sentRes, receivedRes] = await Promise.all([
      fetch(`${apiUrl}/api/follow/requests/sent`, { credentials: 'include' }),
      fetch(`${apiUrl}/api/follow/requests/received`, { credentials: 'include' })
    ])

    if (sentRes.status === 401 || receivedRes.status === 401) {
      logout()
      router.push('/login')
      return
    }

    if (!sentRes.ok || !receivedRes.ok) throw new Error('Failed to fetch pending requests')

    const [sentReqs, receivedReqs] = await Promise.all([
      sentRes.json(),
      receivedRes.json()
    ])

    if (sentReqs) pendingSent.value = sentReqs
    if (receivedReqs) pendingReceived.value = receivedReqs

  } catch (err) {
    console.log(err)
    errorStore.setError('Error', 'Failed to load pending follow requests.')
    router.push('/error')
  }
}

onMounted(() => {
  fetchPendingFollows()
})
</script>