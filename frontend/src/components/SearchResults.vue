<template>
    <div v-if="results && results.length > 0" class="mb-8 w-full max-w-screen-sm">
        <h3 class="text-xl font-semibold text-[var(--nordic-text-dark)] mb-3">Search Results</h3>
        <ul class="space-y-2">
            <RouterLink :to="`/profile/${user.id}`" v-for="user in results" :key="user.id">
            <li class="post-card flex items-center gap-4 mb-4 p-4 bg-[var(--nordic-primary-bg)]
                border border-[var(--nordic-border-light)] rounded-md shadow-sm cursor-pointer">
                <!-- Avatar on the left -->
                <div v-if="user.avatar_url" class="w-16 h-16 rounded-full overflow-hidden flex-shrink-0">
                  <img :src="`${apiUrl}/${user.avatar_url}`" alt="User avatar" class="w-full h-full object-cover"/>
                </div>
                <!-- Name and nickname on the right -->
                <div class="flex flex-col">
                  <span class="font-semibold break-all">
                    {{ user.first_name }} {{ user.last_name }}
                  </span>
                  <span class="text-gray-500 break-all" v-if="user.username">
                    #{{ user.username }}
                  </span>
                </div>
              </li>
            </RouterLink>
        </ul>
    </div>
    <p v-else-if="searchInitiated" class="text-[var(--nordic-text-light)] italic mb-6">No users found.</p>
</template>

<script setup>
defineProps({
    results: Array,
    searchInitiated: Boolean,
    apiUrl: String
})
</script>
