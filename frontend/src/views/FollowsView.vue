<template>
  <div class="follows-page-wrapper">
    <TopBar />

    <TwoColumnLayout>
      <template #sidebar>
        <FollowsInSidebar :userId="user.id" />
        <RequestsInSidebar />
      </template>

      <template #main>
        <h2 class="text-3xl font-bold text-nordic-dark mb-6">Explore Users</h2>

        <!-- search users -->
        <SearchBox @results="handleResults" />
        <SearchResults :results="searchResults" :searchInitiated="searchInitiated" />

        <!-- user suggestions -->
        <div>
          <h3 class="text-xl font-semibold text-nordic-dark mb-3">Suggested Users</h3>
          <ul v-if="users && users.length > 0" class="space-y-2">
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
import SearchBox from '@/components/SearchBox.vue'
import SearchResults from '@/components/SearchResults.vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia';

const apiUrl = import.meta.env.VITE_API_URL
const errorStore = useErrorStore()
const router = useRouter()
const users = ref([])
const searchResults = ref([])
const searchInitiated = ref(false)
const { logout } = useAuth()
const userStore = useUserStore()
const { user } = storeToRefs(userStore)


function handleResults(results) {
  searchInitiated.value = true
  searchResults.value = results
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
