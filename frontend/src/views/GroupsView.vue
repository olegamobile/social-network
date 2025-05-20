<template>
    <div class="groups-page-wrapper">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <GroupsInSidebar />
            </template>

            <template #main>
                <h2 class="text-3xl font-bold text-nordic-dark mb-6">Explore Groups</h2>

                <div class="max-w-lg w-full lg:max-w-xs mb-6">
                    <label for="search-groups" class="sr-only">Search</label>
                    <div class="relative">
                        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <i class="fas fa-search text-nordic-light"></i>
                        </div>
                        <input id="search-groups" title="search"
                            class="block w-full pl-10 pr-3 py-2 border border-nordic-light rounded-md leading-5 bg-white placeholder-nordic-light focus:outline-none focus:ring-2 focus:ring-nordic-secondary-accent focus:border-nordic-secondary-accent sm:text-sm"
                            placeholder="Search for groups..." type="search" v-model="search">
                    </div>
                </div>

                <div>
                    <h3 class="text-xl font-semibold text-nordic-dark mb-3">Suggested Groups</h3>
                    <ul v-if="filteredSuggestions.length > 0" class="space-y-2">
                        <li v-for="group in filteredSuggestions" :key="group.id"
                            class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150">
                            <RouterLink :to="`/groups/${group.id}`">{{ group.title }}</RouterLink>
                            <span class="text-sm text-nordic-light block ml-1">{{ group.description }}</span>
                        </li>
                    </ul>
                    <p v-else class="text-nordic-light italic">No groups match your search.</p>
                </div>
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import TopBar from '@/components/TopBar.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'
import { useErrorStore } from '@/stores/error'
import { useAuth } from '@/composables/useAuth'
import GroupsInSidebar from '@/components/GroupsInSidebar.vue'

const apiUrl = import.meta.env.VITE_API_URL
const search = ref('')
const errorStore = useErrorStore()
const router = useRouter()
const suggestedGroups = ref([])
const { logout } = useAuth()

const userGroups = ref([
    { id: 1, title: 'Class of 24' },
    { id: 2, title: 'Football Team' },
])


async function fetchGroups() {
    try {
        const res = await fetch(`${apiUrl}/api/groups`, {
            credentials: 'include'
        })

        if (res.status === 401) {
            logout();
            router.push('/login');
            return;
        }

        if (!res.ok) throw new Error(`Failed to fetch groups: ${res.status}`)

        suggestedGroups.value = await res.json()

    } catch (err) {
        errorStore.setError('Error', 'Something went wrong while loading users data.')
        router.push('/error')
    }
}

onMounted(() => {
    fetchGroups()
})

const filteredSuggestions = computed(() =>
    suggestedGroups.value.filter(g =>
        g.title.toLowerCase().includes(search.value.toLowerCase())
    )
)
</script>

<style scoped>
.groups-page-wrapper {
    min-height: 100vh;
}
</style>
