<template>
    <div class="groups-page-wrapper">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <GroupsInSidebar />
                <AdminGroupsInSidebar />
                <RequestedGroupsInSidebar />
                <InvitedGroupsInSidebar />
            </template>

            <template #main>
                <h2 class="text-3xl font-bold text-nordic-dark mb-6">Explore Groups</h2>

                <!-- search box -->
                <div class="max-w-lg w-full lg:max-w-xs mb-6">
                    <label for="search-groups" class="sr-only">Search</label>
                    <div class="relative">
                        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <i class="fas fa-search text-nordic-light"></i>
                        </div>
                        <input id="search-groups" title="search"
                            class="block w-full pl-10 pr-3 py-2 border border-nordic-light rounded-md leading-5 bg-white placeholder-nordic-light focus:outline-none focus:ring-2 focus:ring-nordic-secondary-accent focus:border-nordic-secondary-accent sm:text-sm"
                            placeholder="Search groups..." type="search" v-model="searchQuery" @input="searchGroups">
                    </div>
                </div>

                <!-- search results -->
                <div v-if="searchResults && searchResults.length > 0" class="mb-8">
                    <h3 class="text-xl font-semibold text-nordic-dark mb-3">Search Results</h3>
                    <ul class="space-y-2">
                        <li v-for="group in searchResults" :key="group.id"
                            class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                            <RouterLink :to="`/group/${group.id}`">{{ group.title }}</RouterLink>
                            <span class="text-sm text-nordic-light block ml-1">{{ group.description }}</span>
                        </li>
                    </ul>
                </div>
                <p v-else-if="searchInitiated" class="text-nordic-light italic mb-6">No groups found</p>

                <!-- suggested groups -->
                <div>
                    <h3 class="text-xl font-semibold text-nordic-dark mb-3">Suggested Groups</h3>
                    <ul v-if="suggestedGroups && suggestedGroups.length > 0" class="space-y-2">
                        <li v-for="group in suggestedGroups" :key="group.id"
                            class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150">
                            <RouterLink :to="`/group/${group.id}`">{{ group.title }}</RouterLink>
                            <span class="text-sm text-nordic-light block ml-1">{{ group.description }}</span>
                        </li>
                    </ul>
                    <p v-else class="text-nordic-light italic">Not enough data for recommendations</p>
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
import GroupsInSidebar from '@/components/GroupsInSidebar.vue'
import RequestedGroupsInSidebar from '@/components/RequestedGroupsInSidebar.vue'
import InvitedGroupsInSidebar from '@/components/InvitedGroupsInSidebar.vue'
import AdminGroupsInSidebar from '@/components/AdminGroupsInSidebar.vue'
import { useArrayUtils } from '@/composables/useArrayUtils';
import throttle from 'lodash.throttle';

const apiUrl = import.meta.env.VITE_API_URL
const errorStore = useErrorStore()
const router = useRouter()
const suggestedGroups = ref([])
const { logout } = useAuth()
const searchResults = ref([])
const searchQuery = ref('')
const searchInitiated = ref(false)

async function _searchGroups() {
    searchInitiated.value = true;
    searchResults.value = [];

    if (searchQuery.value.trim() === '') return;

    try {
        const response = await fetch(`${apiUrl}/api/groups/search?query=${searchQuery.value}`, {
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
        console.log("Error in searching groups:", String(error))
        errorStore.setError('Error', 'Something went wrong while searching for groups')
        router.push('/error')
    }
}

async function fetchGroups() {
    const { shuffle } = useArrayUtils();
    try {
        const res = await fetch(`${apiUrl}/api/suggestgroups`, {
            credentials: 'include'
        })

        if (res.status === 401) {
            logout();
            router.push('/login');
            return;
        }

        if (!res.ok) throw new Error(`Failed to fetch groups: ${res.status}`)

        suggestedGroups.value = await res.json()
        suggestedGroups.value = shuffle(suggestedGroups.value) // random order

    } catch (err) {
        //console.log("error at suggest groups", err)
        errorStore.setError('Error', 'Something went wrong while loading groups data.')
        router.push('/error')
    }
}


const searchGroups = throttle(_searchGroups, 1000);

onMounted(() => {
    fetchGroups()
})
</script>

<style scoped>
.groups-page-wrapper {
    min-height: 100vh;
}
</style>
