<template>
  <div class="follows-page-wrapper">
    <TopBar />

    <TwoColumnLayout>
      <template #sidebar>
        <FollowsInSidebar :userId="user.id"/>
        <br />
        <RequestsInSidebar />
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

        <div v-if="searchResults && searchResults.length > 0" class="mb-8">
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
import { useArrayUtils } from '@/composables/useArrayUtils';
import FollowsInSidebar from '@/components/FollowsInSidebar.vue'
import RequestsInSidebar from '@/components/RequestsInSidebar.vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia';

const apiUrl = import.meta.env.VITE_API_URL
const errorStore = useErrorStore()
const router = useRouter()
const users = ref([])
const searchQuery = ref('')
const searchResults = ref([])
const searchInitiated = ref(false)
const { logout } = useAuth()
const userStore = useUserStore()
const { user } = storeToRefs(userStore)

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

// suggest non-followed non-self users who either:
// - follow or are followed by a user the active user follows or is followed by
// - are in the same group as the active user
// - follow the active user
async function fetchUserSuggestions() {
  const { shuffle } = useArrayUtils();

  try {
    const res = await fetch(`${apiUrl}/api/suggest/users`, {
      credentials: 'include'
    })

    if (res.status === 401) {
      logout();
      router.push('/login');
      return;
    }

    if (!res.ok) throw new Error(`Failed to fetch suggested users: ${res.status}`)

    users.value = await res.json()
    users.value = shuffle(users.value) // random order
  } catch (err) {
    errorStore.setError('Error', 'Something went wrong while loading suggested users data.')
    router.push('/error')
  }
}

onMounted(() => {
  fetchUserSuggestions()
})
</script>

<style scoped>
.follows-page-wrapper {
  min-height: 100vh;
}
</style>
