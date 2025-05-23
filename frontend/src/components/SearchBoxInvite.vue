<template>
    <div class="max-w-lg w-full lg:max-w-xs mb-6">
        <label for="search-users" class="sr-only">Search</label>
        <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <i class="fas fa-search text-[var(--nordic-text-light)]"></i>
            </div>
            <input id="search-users" name="search"
                class="block w-full pl-10 pr-3 py-2 border border-[var(--nordic-border-light)] rounded-md leading-5 bg-white placeholder-[var(--nordic-text-light)] focus:outline-none focus:ring-2 focus:ring-[var(--nordic-secondary-accent)] focus:border-[var(--nordic-secondary-accent)] sm:text-sm"
                placeholder="Search users..." type="search" v-model="searchQuery" @input="searchUsers">
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import throttle from 'lodash.throttle'

const emit = defineEmits(['results'])

const apiUrl = import.meta.env.VITE_API_URL
const searchQuery = ref('')
const searchUsers = throttle(_searchUsers, 1000)
const route = useRoute()

async function _searchUsers() {
    const query = searchQuery.value.trim()
    if (query === '') {
        emit('results', []) // Emit empty results
        return
    }
    const groupId = route.params.id

    try {
        const response = await fetch(`${apiUrl}/api/group/invite/search?query=${query}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({ group_id: groupId })
        });

        if (response.status === 401) {
            // Optional: emit a special event to trigger logout in parent if needed
            return;
        }

        if (!response.ok) {
            console.error("Search error:", response.status)
            emit('results', []) // Or emit an error object
            return
        }

        const data = await response.json()
        emit('results', data)
    } catch (err) {
        console.error("Error in search:", err)
        emit('results', [])
    }
}
</script>
