<template>
    <div class="groups-page">
        <TopBar />

        <TwoColumnLayout>
            <template #sidebar>
                <h3 class="text-lg font-semibold">Your Groups</h3>
                <ul>
                    <li v-for="group in userGroups" :key="group.id">
                        <RouterLink :to="`/groups/${group.id}`">{{ group.name }}</RouterLink>
                    </li>
                </ul>
            </template>

            <template #main>
                <h2 class="text-2xl font-bold mb-4">Explore Groups</h2>

                <div class="max-w-lg w-full lg:max-w-xs mb-4">
                    <label for="search" class="sr-only">Search</label>
                    <div class="relative">
                        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                            <i class="fas fa-search text-nordic-light"></i>
                        </div>
                        <input id="search" name="search"
                            class="block w-full pl-10 pr-3 py-2 border border-nordic-light rounded-md leading-5 bg-white placeholder-nordic-light focus:outline-none focus:ring-2 focus:ring-nordic-secondary-accent focus:border-nordic-secondary-accent sm:text-sm"
                            placeholder="Search for groups..." type="search" v-model="search">
                    </div>
                </div>

                <h3 class="text-lg font-semibold">Suggested Groups</h3>
                <ul>
                    <li v-for="group in filteredSuggestions" :key="group.id">
                        <RouterLink :to="`/groups/${group.id}`">{{ group.name }}</RouterLink> - {{ group.description }}
                    </li>
                </ul>
            </template>
        </TwoColumnLayout>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import TopBar from '@/components/TopBar.vue'
import TwoColumnLayout from '@/layouts/TwoColumnLayout.vue'

const search = ref('')

const userGroups = ref([
    { id: 1, name: 'Class of 24' },
    { id: 2, name: 'Football Team' },
])

const suggestedGroups = ref([
    { id: 3, name: 'Photography Club', description: 'Share and learn photography' },
    { id: 4, name: 'Study Group', description: 'Daily study sessions' },
])

const filteredSuggestions = computed(() =>
    suggestedGroups.value.filter(g =>
        g.name.toLowerCase().includes(search.value.toLowerCase())
    )
)
</script>
