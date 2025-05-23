<template>
    <div v-if="results && results.length > 0" class="mb-8">
        <h3 class="text-xl font-semibold text-nordic-dark mb-3">Search Results</h3>
        <ul class="space-y-2">
            <li v-for="user in results" :key="user.id"
                class="text-nordic-light hover:text-nordic-primary-accent transition-colors duration-150 cursor-pointer">
                <RouterLink :to="`/profile/${user.id}`">
                    {{ user.first_name }} {{ user.last_name }}
                    <span v-if="user.username"> - {{ user.username }}</span>
                </RouterLink>
            </li>
        </ul>
    </div>
    <p v-else-if="searchInitiated" class="text-nordic-light italic mb-6">No users found.</p>
</template>

<script setup>
defineProps({
    results: Array,
    searchInitiated: Boolean
})
</script>
