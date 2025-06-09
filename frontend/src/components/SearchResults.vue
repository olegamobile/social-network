<template>
    <div v-if="results && results.length > 0" class="mb-8">
        <h3 class="text-xl font-semibold text-[var(--nordic-text-dark)] mb-3">Search Results</h3>
        <ul class="space-y-2">
            <li v-for="user in results" :key="user.id"
                class="text-[var(--nordic-text-light)] hover:text-[var(--nordic-primary-accent)] transition-colors duration-150 cursor-pointer break-all">
                <RouterLink :to="`/profile/${user.id}`">
                    {{ user.first_name }} {{ user.last_name }}
                    <span class="break-all" v-if="user.username"> - {{ user.username }}</span>
                </RouterLink>
            </li>
        </ul>
    </div>
    <p v-else-if="searchInitiated" class="text-[var(--nordic-text-light)] italic mb-6">No users found.</p>
</template>

<script setup>
defineProps({
    results: Array,
    searchInitiated: Boolean
})
</script>
