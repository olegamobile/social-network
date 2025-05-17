<template>
  <div class="follows-page-wrapper">
    <TopBar />

    <TwoColumnLayout>
      <template #sidebar>

        <!-- Users I am following -->
        <div>
          <h3 class="text-xl font-semibold text-nordic-dark mb-3">Users You Follow</h3>
          <ul v-if="followedUsers.length > 0" class="space-y-2">
            <li v-for="user in followedUsers" :key="user.id"
              class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
              <RouterLink :to="`/profile/${user.id}`">
                {{ user.first_name }} {{ user.last_name }}<span v-if="user.username"> - {{ user.username }}</span>
              </RouterLink>
            </li>
          </ul>
          <p v-else class="text-nordic-light italic">You're not following anyone yet.</p>
        </div>

        <br />

        <!-- Users who follow me -->
        <div>
          <h3 class="text-xl font-semibold text-nordic-dark mb-3">Your Followers</h3>
          <ul v-if="followers.length > 0" class="space-y-2">
            <li v-for="user in followers" :key="user.id"
              class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
              <RouterLink :to="`/profile/${user.id}`">
                {{ user.first_name }} {{ user.last_name }}<span v-if="user.username"> - {{ user.username }}</span>
              </RouterLink>
            </li>
          </ul>
          <p v-else class="text-nordic-light italic">No one is following you yet.</p>
        </div>

        <!-- Sent Requests -->
        <div>
          <h3 class="text-xl font-semibold text-nordic-dark mb-3">Pending Requests You've Sent</h3>
          <ul v-if="pendingSent.length > 0" class="space-y-2">
            <li v-for="user in pendingSent" :key="user.id"
              class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
              <RouterLink :to="`/profile/${user.id}`">
                {{ user.first_name }} {{ user.last_name }}<span v-if="user.username"> - {{ user.username }}</span>
              </RouterLink>
            </li>
          </ul>
          <p v-else class="text-nordic-light italic">You haven't sent any requests yet.</p>
        </div>

        <br />

        <!-- Received Requests -->
        <div>
          <h3 class="text-xl font-semibold text-nordic-dark mb-3">Pending Requests You've Received</h3>
          <ul v-if="pendingReceived.length > 0" class="space-y-2">
            <li v-for="user in pendingReceived" :key="user.id"
              class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
              <RouterLink :to="`/profile/${user.id}`">
                {{ user.first_name }} {{ user.last_name }}<span v-if="user.username"> - {{ user.username }}</span>
              </RouterLink>
            </li>
          </ul>
          <p v-else class="text-nordic-light italic">You have no pending requests.</p>
        </div>


      </template>

      <template #main>
        <h2 class="text-3xl font-bold text-nordic-dark mb-6">Explore Users</h2>

        <div class="max-w-lg w-full lg:max-w-xs mb-6">
          <label for="search" class="sr-only">Search</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <i class="fas fa-search text-nordic-light"></i>
            </div>
            <input id="search" name="search"
              class="block w-full pl-10 pr-3 py-2 border border-nordic-light rounded-md leading-5 bg-white placeholder-nordic-light focus:outline-none focus:ring-2 focus:ring-nordic-secondary-accent focus:border-nordic-secondary-accent sm:text-sm"
              placeholder="Search users..." type="search" v-model="searchQuery" @input="searchUsers">
          </div>
        </div>

        <div v-if="searchResults.length > 0" class="mb-8">
          <h3 class="text-xl font-semibold text-nordic-dark mb-3">Search Results</h3>
          <ul class="space-y-2">
            <li v-for="user in searchResults" :key="user.id"
              class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
              <RouterLink :to="`/profile/${user.id}`">
                {{ user.first_name }} {{ user.last_name }}<span v-if="user.username"> - {{ user.username }}</span>
              </RouterLink>
            </li>
          </ul>
        </div>
        <p v-else-if="searchInitiated" class="text-nordic-light italic mb-6">No users found.</p>

        <div>
          <h3 class="text-xl font-semibold text-nordic-dark mb-3">Suggested Users</h3>
          <ul v-if="users.length > 0" class="space-y-2">
            <li v-for="user in users" :key="user.id"
              class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
              <RouterLink :to="`/profile/${user.id}`">
                {{ user.first_name }} {{ user.last_name }}<span v-if="user.username"> - {{ user.username }}</span>
              </RouterLink>
            </li>
          </ul>
          <p v-else class="text-nordic-light italic">No suggested users available.</p>
        </div>
      </template>
    </TwoColumnLayout>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import TopBar from '@/components/TopBar.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { useErrorStore } from '@/stores/error'
import { useAuth } from '@/composables/useAuth'

const apiUrl = import.meta.env.VITE_API_URL
const errorStore = useErrorStore()
const router = useRouter()
const users = ref([])
const searchQuery = ref('');
const searchResults = ref([]);
const searchInitiated = ref(false);
const { logout } = useAuth()
const followedUsers = ref([])     // users I'm following
const followers = ref([])          // users following me
const pendingSent = ref([])     // users I've sent requests to
const pendingReceived = ref([]) // users who sent me requests

async function searchUsers() {
  searchInitiated.value = true;
  searchResults.value = [];

  if (searchQuery.value.trim() === '') return;

  try {
    const response = await fetch(`${apiUrl}/api/users/search?query=${searchQuery.value}`, {
      credentials: 'include'
    });

    if (response.status === 401) {
      logout();
      router.push('/login');
      return;
    }

    if (!response.ok) {
      errorStore.setError(response.status, 'Error fetching search results')
      router.push('/error')
      return
    }

    searchResults.value = await response.json();
  } catch (error) {
    console.log("Error in searching users:", String(error))
    errorStore.setError('Error', 'Something went wrong while searching for users')
    router.push('/error')
  }
}

async function fetchUsers() {
  try {
    const res = await fetch(`${apiUrl}/api/users`, {
      credentials: 'include'
    })

    if (res.status === 401) {
      logout();
      router.push('/login');
      return;
    }

    if (!res.ok) throw new Error(`Failed to fetch users: ${res.status}`)

    users.value = await res.json()
  } catch (err) {
    errorStore.setError('Error', 'Something went wrong while loading users data.')
    router.push('/error')
  }
}

async function fetchFollowData() {
  try {
    const res1 = await fetch(`${apiUrl}/api/followed`, { credentials: 'include' })
    const res2 = await fetch(`${apiUrl}/api/followers`, { credentials: 'include' })

    if (res1.status === 401 || res2.status === 401) {
      logout()
      router.push('/login')
      return
    }

    if (!res1.ok || !res2.ok) throw new Error('Failed to fetch follow data')

    followedUsers.value = await res1.json()
    followers.value = await res2.json()
  } catch (err) {
    errorStore.setError('Error', 'Failed to load follow data.')
    router.push('/error')
  }
}

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

    /*    let sentReady
       try {
         sentReady = await sentRes.json()
         console.log(sentReady)
       } catch (error) {
         console.log("sents err:", error)
       }
   
       let receivedReady
       try {
         receivedReady = await receivedRes.json()
         console.log(receivedReady)
       } catch (error) {
         console.log("receiveds err:", error)
       } */
    let sentReady
    let receivedReady
    sentReady = await sentRes.json()
    receivedReady = await receivedRes.json()
    if (sentReady) pendingSent.value.push(...sentReady)
    if (receivedReady) pendingReceived.value.push(...receivedReady)

  } catch (err) {
    console.log(err)
    errorStore.setError('Error', 'Failed to load pending follow requests.')
    router.push('/error')
  }
}


onMounted(() => {
  fetchUsers()
  fetchFollowData()
  fetchPendingFollows()
})
</script>

<style scoped>
.follows-page-wrapper {
  min-height: 100vh;
}
</style>
